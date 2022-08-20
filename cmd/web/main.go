package main

import (
	"database/sql"
	"log"

	"github.com/burakkarasel/Theatre-API/internal/api"
	db "github.com/burakkarasel/Theatre-API/internal/db/sqlc"
	"github.com/burakkarasel/Theatre-API/internal/util"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	// first i load the env variables
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("cannot load env variables:", err)
	}

	// then i connect to DB
	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to DB:", err)
	}

	log.Println("connected to DB")

	// then i run my migrations
	runDBMigration("file://internal/db/migration", config.DBSource)

	// then i create a new store to create a new server
	store := db.NewStore(conn)

	log.Println("created a new store instance")

	// then i create a new server to run the server
	server := api.NewServer(store)

	log.Println("created a new server instance")

	err = server.Start(config.ServerAddress)

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
