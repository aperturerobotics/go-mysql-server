package driver_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/driver"
)

// TestSharedEngineDrivers keeps session IDs and engine ownership across drivers.
func TestSharedEngineDrivers(t *testing.T) {
	// Use the real memory catalog and one engine borrowed by separate drivers.
	provider, _ := personMemTable("db", "person")
	_, catalog, err := provider.Resolve("db", nil)
	require.NoError(t, err)
	engine := sqle.NewDefault(catalog)
	defer engine.Close()
	var ids []uint64
	for range 2 {
		wrapped := driver.New(provider, &driver.Options{Engine: engine})
		connector, err := wrapped.OpenConnector("db")
		require.NoError(t, err)
		db := sql.OpenDB(connector)
		var id uint64
		require.NoError(t, db.QueryRow("SELECT CONNECTION_ID()").Scan(&id))
		ids = append(ids, id)
		require.NoError(t, db.Close())
		require.NoError(t, wrapped.Close())
	}
	require.NotEqual(t, ids[0], ids[1])
}
