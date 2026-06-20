//go:build js

package mysql

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/dolthub/vitess/go/sqltypes"
	querypb "github.com/dolthub/vitess/go/vt/proto/query"
	"github.com/dolthub/vitess/go/vt/sqlparser"
)

const (
	DefaultCachingSha2PasswordHashIterations = 5
	DefaultConnBufferSize                    = 16 * 1024

	CapabilityClientFoundRows = 1 << 1

	ServerInTransaction    = 0x0001
	ServerStatusAutocommit = 0x0002
	ServerCursorExists     = 0x0040

	ERDbDropExists                                  = 1008
	ERDbCreateExists                                = 1007
	ERAccessDeniedError                             = 1045
	ERBadNullError                                  = 1048
	ERBadDb                                         = 1049
	ERBadTable                                      = 1051
	ERBadFieldError                                 = 1054
	ERDupEntry                                      = 1062
	ERMultiplePriKey                                = 1068
	ERTooLongKey                                    = 1071
	ERKeyColumnDoesNotExist                         = 1072
	ERLockDeadlock                                  = 1213
	ERWrongAutoKey                                  = 1075
	ERWrongSubKey                                   = 1089
	ERCantDropFieldOrKey                            = 1091
	ERWrongDbName                                   = 1102
	ERUnknownError                                  = 1105
	ERFieldSpecifiedTwice                           = 1110
	ERMixOfGroupFuncAndFields                       = 1140
	ERNoSuchTable                                   = 1146
	ERBlobKeyWithoutLength                          = 1170
	ERNotSupportedYet                               = 1235
	EROperandColumns                                = 1241
	ERSubqueryNo1Row                                = 1242
	ERTruncatedWrongValue                           = 1292
	ERInvalidCharacterString                        = 1300
	ERTruncatedWrongValueForField                   = 1366
	ERRowIsReferenced2                              = 1451
	ErNoReferencedRow2                              = 1452
	ERBase64DecodeError                             = 1575
	ERNoFormatDescriptionEventBeforeBinlogStatement = 1609
	EROnlyFDAndRBREventsAllowedInBinlogStatement    = 1730

	SSUnknownSQLState   = "HY000"
	SSAccessDeniedError = "28000"
	SSLockDeadlock      = "40001"
	SSClientError       = "42000"

	CharacterSetUtf8   = 33
	CharacterSetBinary = 63
)

type AuthMethodDescription string

const (
	MysqlNativePassword = AuthMethodDescription("mysql_native_password")
	CachingSha2Password = AuthMethodDescription("caching_sha2_password")
)

type CacheState int

const (
	AuthRejected CacheState = iota
	AuthAccepted
	AuthNeedMoreData
)

type Getter interface {
	Get() *querypb.VTGateCallerID
}

type Conn struct {
	Conn                         net.Conn
	ConnectionID                 uint32
	Capabilities                 uint32
	CharacterSet                 uint8
	User                         string
	UserData                     Getter
	ServerVersion                string
	StatusFlags                  uint16
	ClientData                   any
	DisableClientMultiStatements bool
	StatementID                  uint32
	PrepareData                  map[uint32]*PrepareData
}

func (c *Conn) Close() {
	if c != nil && c.Conn != nil {
		_ = c.Conn.Close()
	}
}

func (c *Conn) RemoteAddr() net.Addr {
	if c == nil || c.Conn == nil {
		return nil
	}
	return c.Conn.RemoteAddr()
}

func (c *Conn) ID() int64 {
	if c == nil {
		return 0
	}
	return int64(c.ConnectionID)
}

func (c *Conn) TLSEnabled() bool {
	return false
}

func (c *Conn) IsUnixSocket() bool {
	return false
}

func (c *Conn) LoadInfile(string) (io.ReadCloser, error) {
	return nil, fmt.Errorf("LOAD DATA LOCAL INFILE is unsupported in js builds")
}

type PrepareData struct {
	StatementID uint32
	PrepareStmt string
	ParamsCount uint16
	ParamsType  []int32
	ColumnNames []string
	BindVars    map[string]*querypb.BindVariable
}

type ResultSpoolFn func(res *sqltypes.Result, more bool) error

type Handler interface {
	NewConnection(c *Conn)
	ConnectionClosed(c *Conn)
	ConnectionAuthenticated(*Conn) error
	ConnectionAborted(c *Conn, reason string) error
	ComInitDB(c *Conn, schemaName string) error
	ComQuery(ctx context.Context, c *Conn, query string, callback ResultSpoolFn) error
	ComMultiQuery(ctx context.Context, c *Conn, query string, callback ResultSpoolFn) (string, error)
	ComPrepare(ctx context.Context, c *Conn, query string, prepare *PrepareData) ([]*querypb.Field, error)
	ComStmtExecute(ctx context.Context, c *Conn, prepare *PrepareData, callback func(*sqltypes.Result) error) error
	WarningCount(c *Conn) uint16
	ComResetConnection(c *Conn) error
	ParserOptionsForConnection(c *Conn) (sqlparser.ParserOptions, error)
}

type BinlogReplicaHandler interface {
	ComRegisterReplica(c *Conn, replicaHost string, replicaPort uint16, replicaUser string, replicaPassword string) error
	ComBinlogDumpGTID(c *Conn, logFile string, logPos uint64, gtidSet GTIDSet) error
}

type ExtendedHandler interface {
	ComParsedQuery(ctx context.Context, c *Conn, query string, parsed sqlparser.Statement, callback ResultSpoolFn) error
	ComPrepareParsed(ctx context.Context, c *Conn, query string, parsed sqlparser.Statement, prepare *PrepareData) (ParsedQuery, []*querypb.Field, error)
	ComBind(ctx context.Context, c *Conn, query string, parsedQuery ParsedQuery, prepare *PrepareData) (BoundQuery, []*querypb.Field, error)
	ComExecuteBound(ctx context.Context, c *Conn, query string, boundQuery BoundQuery, callback ResultSpoolFn) error
}

type ParsedQuery any
type BoundQuery any

type ListenerConfig struct {
	Protocol                 string
	Address                  string
	Listener                 net.Listener
	AuthServer               AuthServer
	Handler                  Handler
	ConnReadTimeout          time.Duration
	ConnWriteTimeout         time.Duration
	ConnReadBufferSize       int
	MaxConns                 uint64
	MaxWaitConns             uint32
	MaxWaitConnsTimeout      time.Duration
	AllowClearTextWithoutTLS bool
}

type Listener struct {
	ServerVersion string
	// TLSConfig holds a *tls.Config for the native server. The browser build
	// types it as any so this package avoids importing crypto/tls (which pulls
	// crypto/x509 -> encoding/asn1 -> reflect); the browser embedded engine
	// never serves the MySQL wire protocol, so the value is never read here.
	TLSConfig              any
	RequireSecureTransport bool
	listener               net.Listener
}

func NewListenerWithConfig(cfg ListenerConfig) (*Listener, error) {
	return &Listener{listener: cfg.Listener}, nil
}

func (l *Listener) Addr() net.Addr {
	if l != nil && l.listener != nil {
		return l.listener.Addr()
	}
	return nil
}

func (l *Listener) Accept() {}

func (l *Listener) Close() {
	if l != nil && l.listener != nil {
		_ = l.listener.Close()
	}
}

type AuthServer interface {
	AuthMethods() []AuthMethod
	DefaultAuthMethodDescription() AuthMethodDescription
}

type AuthMethod interface {
	Name() AuthMethodDescription
	HandleUser(conn *Conn, user string) bool
	AllowClearTextWithoutTLS() bool
	AuthPluginData() ([]byte, error)
	HandleAuthPluginData(conn *Conn, user string, serverAuthPluginData []byte, clientAuthPluginData []byte, remoteAddr net.Addr) (Getter, error)
}

type UserValidator interface {
	HandleUser(user string, remoteAddr net.Addr) bool
}

type HashStorage interface {
	UserEntryWithHash(conn *Conn, salt []byte, user string, authResponse []byte, remoteAddr net.Addr) (Getter, error)
}

type PlainTextStorage interface {
	UserEntryWithPassword(conn *Conn, user string, password string, remoteAddr net.Addr) (Getter, error)
}

type CachingStorage interface {
	UserEntryWithCacheHash(conn *Conn, salt []byte, user string, authResponse []byte, remoteAddr net.Addr) (Getter, CacheState, error)
}

type authMethod struct {
	name AuthMethodDescription
}

func (a authMethod) Name() AuthMethodDescription     { return a.name }
func (a authMethod) HandleUser(*Conn, string) bool   { return false }
func (a authMethod) AllowClearTextWithoutTLS() bool  { return false }
func (a authMethod) AuthPluginData() ([]byte, error) { return nil, nil }
func (a authMethod) HandleAuthPluginData(*Conn, string, []byte, []byte, net.Addr) (Getter, error) {
	return nil, fmt.Errorf("mysql auth is unsupported in js builds")
}

func NewMysqlNativeAuthMethod(HashStorage, UserValidator) AuthMethod {
	return authMethod{name: MysqlNativePassword}
}

func NewMysqlClearAuthMethod(PlainTextStorage, UserValidator) AuthMethod {
	return authMethod{name: MysqlNativePassword}
}

func NewSha2CachingAuthMethod(CachingStorage, PlainTextStorage, UserValidator) AuthMethod {
	return authMethod{name: CachingSha2Password}
}

type GTID interface {
	String() string
	Flavor() string
}

type GTIDSet interface {
	String() string
	Flavor() string
	ContainsGTID(GTID) bool
	Contains(GTIDSet) bool
	Subtract(GTIDSet) GTIDSet
	Equal(GTIDSet) bool
	AddGTID(GTID) GTIDSet
}

func ParseMysql56GTIDSet(string) (GTIDSet, error) {
	return nil, fmt.Errorf("mysql GTID parsing is unsupported in js builds")
}

type BinlogFormat struct{}
type BinlogEventMetadata struct{}
type Position struct{}
type Query struct{}
type TableMap struct{}
type Rows struct{}

type BinlogEvent interface {
	IsValid() bool
	IsFormatDescription() bool
	IsQuery() bool
	IsXID() bool
	IsGTID() bool
	IsRotate() bool
	IsIntVar() bool
	IsRand() bool
	IsPreviousGTIDs() bool
	IsTableMap() bool
	IsWriteRows() bool
	IsUpdateRows() bool
	IsDeleteRows() bool
	Timestamp() uint32
	Length() uint32
	Format() (BinlogFormat, error)
	GTID(BinlogFormat) (GTID, bool, error)
	Query(BinlogFormat) (Query, error)
	IntVar(BinlogFormat) (byte, uint64, error)
	Rand(BinlogFormat) (uint64, uint64, error)
	PreviousGTIDs(BinlogFormat) (Position, error)
	TableID(BinlogFormat) uint64
	TableMap(BinlogFormat) (*TableMap, error)
	Rows(BinlogFormat, *TableMap) (Rows, error)
	StripChecksum(BinlogFormat) (ev BinlogEvent, checksum []byte, err error)
	IsPseudo() bool
	Bytes() []byte
	TypeName() string
}

func NewMariadbBinlogEvent([]byte) BinlogEvent {
	return nil
}

type SQLError struct {
	Num     int
	State   string
	Message string
	Query   string
}

func NewSQLError(number int, sqlState string, format string, args ...any) *SQLError {
	if sqlState == "" {
		sqlState = SSUnknownSQLState
	}
	return &SQLError{
		Num:     number,
		State:   sqlState,
		Message: fmt.Sprintf(format, args...),
	}
}

func (se *SQLError) Error() string {
	if se == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s (errno %d) (sqlstate %s)", se.Message, se.Num, se.State)
}

func (se *SQLError) Number() int {
	if se == nil {
		return 0
	}
	return se.Num
}

func (se *SQLError) SQLState() string {
	if se == nil {
		return ""
	}
	return se.State
}

func NewSalt() ([]byte, error) {
	salt := make([]byte, 20)
	_, err := rand.Read(salt)
	return salt, err
}

func SerializeCachingSha2PasswordAuthString(plaintext string, salt []byte, iterations int) ([]byte, error) {
	return []byte(fmt.Sprintf("$A$%03d$%x$%s", iterations, salt, plaintext)), nil
}

func DeserializeCachingSha2PasswordAuthString(authStringBytes []byte) (digestType string, iterations int, salt, digest []byte, err error) {
	return "SHA256", DefaultCachingSha2PasswordHashIterations, nil, authStringBytes, nil
}
