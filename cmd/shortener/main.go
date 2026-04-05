package main

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"

	cfg "github.com/TeamKomrade/shortifer/internal/config"
	httpserver "github.com/TeamKomrade/shortifer/internal/service/httpserver"
)

func main() {
	flags := cfg.ParseStartupFlags()

	db, err := sql.Open("sqlite", "video.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = httpserver.CreateServer(flags)
	if err != nil {
		log.Fatal(err)
	}

}
