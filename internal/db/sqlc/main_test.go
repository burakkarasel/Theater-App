package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/burakkarasel/Theatre-API/internal/dsn"
	_ "github.com/lib/pq"
)

const dbDriver = "postgres"

var testQueries *Queries
var testDB *sql.DB

// TestMain sets up the DB connection for testing
func TestMain(m *testing.M) {
	var err error

	testDB, err = sql.Open(dbDriver, dsn.DSN)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
