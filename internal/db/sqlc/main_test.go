package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/burakkarasel/Theatre-API/internal/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

// TestMain sets up the DB connection for testing
func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../..")

	if err != nil {
		log.Fatal("cannot load env variables:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
