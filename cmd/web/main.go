package main

import (
	"database/sql"
	"log"

	"github.com/burakkarasel/Theatre-API/internal/api"
	db "github.com/burakkarasel/Theatre-API/internal/db/sqlc"
	"github.com/burakkarasel/Theatre-API/internal/dsn"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = dsn.DSN
	port     = "localhost:8080"
)

func main() {
	// first i connect to DB
	conn, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("cannot connect to DB:", err)
	}

	log.Println("connected to DB")

	// then i run my migrations
	runDBMigration("file://internal/db/migration", dsn.DSN)

	// then i create a new store to create a new server
	store := db.NewStore(conn)

	log.Println("created a new store instance")

	// then i create a new server to run the server
	server := api.NewServer(store)

	log.Println("created a new server instance")

	err = server.Start(port)

	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

// runDBMigration runs the migrations at the start of the program
func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create new migrate instance:", err)
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up:", err)
	}

	log.Println("db migrated succesfully")
}
