// Copyright 2020-2021 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package driver

import (
	"context"
	"database/sql/driver"
	"fmt"
	"net/url"
	"sync"
	"sync/atomic"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/analyzer"
)

// ScanKind indicates how values should be scanned.
type ScanKind int

const (
	// ScanAsString indicates values should be scanned as strings.
	//
	// Applies to JSON columns.
	ScanAsString ScanKind = iota

	// ScanAsBytes indicates values should be scanned as byte arrays.
	//
	// Applies to JSON columns.
	ScanAsBytes

	// ScanAsObject indicates values should be scanned as objects.
	//
	// Applies to JSON columns.
	ScanAsObject

	// ScanAsStored indicates values should not be modified during scanning.
	//
	// Applies to JSON columns.
	ScanAsStored
)

// Options for the driver.
type Options struct {
	// Engine is an existing engine to borrow instead of constructing one per DSN.
	// Its catalog must resolve each connection's context. The caller closes it
	// after every borrowing driver and connection has finished.
	Engine *sqle.Engine
	// EngineVariables supplies the variable bundle shared by this driver's engines.
	// A nil bundle gives each new engine independent default variable state.
	// A borrowed Engine already owns its variables and ignores this option.
	EngineVariables *sql.EngineVariables
	// JSON indicates how JSON row values should be scanned
	JSON ScanKind
}

// Provider resolves SQL catalogs from DSNs.
//
// Provider can optionally implement ProviderWithSessionBuilder and ProviderWithSessionBuilder.
type Provider interface {
	// Resolve determines the server name and DatabaseProvider for the given dsn.
	// The will create a separate sql.Engine for each "server name" value.
	Resolve(dsn string, options *Options) (string, sql.DatabaseProvider, error)
}

// ProviderWithContextBuilder is a Provider that also provides a base sql.Context.
// Provider should parse the DSN and use it to adjust the context in NewContext.
type ProviderWithContextBuilder interface {
	Provider
	ContextBuilder
}

// ProviderWithSessionBuilder is a Provider that also constructs sessions.
type ProviderWithSessionBuilder interface {
	Provider
	SessionBuilder
}

// A Driver exposes an engine as a stdlib SQL driver.
type Driver struct {
	provider Provider
	sessions SessionBuilder
	contexts ContextBuilder
	options  *Options
	dbs      map[string]*dbConn
	mu       sync.Mutex
}

// New returns a driver using the specified provider.
func New(provider Provider, options *Options) *Driver {
	var sessionBuilder SessionBuilder
	if provWithSessionBuilder, ok := provider.(ProviderWithSessionBuilder); ok {
		sessionBuilder = provWithSessionBuilder
	}

	var contextBuilder ContextBuilder
	if provWithCtxBuilder, ok := provider.(ProviderWithContextBuilder); ok {
		contextBuilder = provWithCtxBuilder
	} else {
		contextBuilder = DefaultContextBuilder{}
	}

	return &Driver{
		provider: provider,
		options:  options,
		sessions: sessionBuilder,
		contexts: contextBuilder,
		dbs:      map[string]*dbConn{},
	}
}

// Open returns a new connection to the database.
func (d *Driver) Open(name string) (driver.Conn, error) {
	conn, err := d.OpenConnector(name)
	if err != nil {
		return nil, err
	}
	return conn.Connect(context.Background())
}

// OpenConnector calls the driver factory and returns a new connector.
func (d *Driver) OpenConnector(dsn string) (driver.Connector, error) {
	// Normalize this connector's scan options and DSN.
	options := d.options // copy
	if options == nil {
		options = &Options{}
	}

	dsnURI, err := url.Parse(dsn)
	if err == nil {
		query := dsnURI.Query()
		qJSON := query.Get("jsonAs")
		switch qJSON {
		case "":
			// default
		case "object":
			options.JSON = ScanAsObject
		case "string":
			options.JSON = ScanAsString
		case "bytes":
			options.JSON = ScanAsBytes
		default:
			return nil, fmt.Errorf("%q is not a valid option for 'jsonAs'", qJSON)
		}

		query.Del("jsonAs")
		dsnURI.RawQuery = query.Encode()
		dsn = dsnURI.String()
	}

	// Resolve the catalog and connection-local session builder.
	serverName, pro, err := d.provider.Resolve(dsn, options)
	if err != nil {
		return nil, err
	}

	if d.sessions == nil {
		d.sessions = NewDefaultSessionBuilder(pro)
	}

	// Own a new engine or borrow the caller's engine for this server.
	d.mu.Lock()
	db, ok := d.dbs[serverName]
	if !ok {
		engine := options.Engine
		if engine == nil {
			anlz := analyzer.NewDefault(pro)
			engine = sqle.New(anlz, &sqle.Config{EngineVariables: options.EngineVariables})
		}
		db = &dbConn{engine: engine, borrowed: options.Engine != nil}
		d.dbs[serverName] = db
	}
	d.mu.Unlock()

	// Keep the connection bound to its selected driver and engine.
	return &Connector{
		driver:     d,
		options:    options,
		serverName: serverName,
		dsn:        dsn,
		dbConn:     db,
	}, nil
}

func (d *Driver) Close() error {
	var firstErr error
	for _, conn := range d.dbs {
		if firstErr == nil {
			firstErr = conn.close()
		} else {
			conn.close()
		}
	}
	return firstErr
}

// driverConnectionIDs and driverProcessIDs keep identities distinct when several
// drivers borrow the same engine. Allocating an identity costs one atomic add.
var (
	driverConnectionIDs atomic.Uint32
	driverProcessIDs    atomic.Uint64
)

// dbConn owns or borrows the engine used by one resolved server.
type dbConn struct {
	engine   *sqle.Engine
	borrowed bool
}

// nextConnectionID allocates a driver session identity across borrowed engines.
func (c *dbConn) nextConnectionID() uint32 { return driverConnectionIDs.Add(1) }

// nextProcessID allocates an execution identity across borrowed engines.
func (c *dbConn) nextProcessID() uint64 { return driverProcessIDs.Add(1) }

// close releases only an engine this driver constructed.
func (c *dbConn) close() error {
	if c.borrowed {
		return nil
	}
	return c.engine.Close()
}

// A Connector represents a driver in a fixed configuration
// and can create any number of equivalent Conns for use
// by multiple goroutines.
type Connector struct {
	driver     *Driver
	options    *Options
	dbConn     *dbConn
	serverName string
	dsn        string
}

// Driver returns the driver.
func (c *Connector) Driver() driver.Driver { return c.driver }

// Server returns the server name.
func (c *Connector) Server() string { return c.serverName }

// DSN returns the original DSN passed by the client.
func (c *Connector) DSN() string { return c.dsn }

// EngineVariables returns the state shared by connections to this engine.
func (c *Connector) EngineVariables() *sql.EngineVariables { return c.dbConn.engine.Variables() }

// Connect returns a connection to the database.
func (c *Connector) Connect(ctx context.Context) (driver.Conn, error) {
	id := c.dbConn.nextConnectionID()

	session, err := c.driver.sessions.NewSession(ctx, id, c)
	if err != nil {
		return nil, err
	}

	indexes := sql.NewIndexRegistry()
	views := sql.NewViewRegistry()
	return &Conn{
		dsn:      c.dsn,
		options:  c.options,
		dbConn:   c.dbConn,
		session:  session,
		contexts: c.driver.contexts,
		indexes:  indexes,
		views:    views,
	}, nil
}
