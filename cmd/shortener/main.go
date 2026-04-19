package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	cfg "github.com/TeamKomrade/shortifer/internal/config"
	httpserver "github.com/TeamKomrade/shortifer/internal/service/httpserver"
	"github.com/golang-migrate/migrate/v4"
)

func main() {
	flags := cfg.ParseStartupFlags()
	databaseConnectionString := flags.DatabaseConnString

	envDatabaseConnString := os.Getenv("DATABASE_CONN_STRING")
	if envDatabaseConnString != "" {
		databaseConnectionString = envDatabaseConnString
	}

	db, err := sql.Open("pgx", databaseConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	migrationsPath := "file://migrations"

	m, err := migrate.New(
		migrationsPath,
		databaseConnectionString,
	)
	if err != nil {
		log.Fatal("Error initializing migrate:", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Error applying migrations:", err)
	}
	log.Print("Database migrations applied successfully!")

	err = httpserver.CreateServer(flags)
	if err != nil {
		log.Fatal(err)
	}

}
