package main

import (
	"log"

	cfg "github.com/TeamKomrade/shortifer/internal/config"
	httpserver "github.com/TeamKomrade/shortifer/internal/service/httpserver"
)

func main() {
	flags := cfg.ParseStartupFlags()
	err := httpserver.CreateServer(flags)

	if err != nil {
		log.Fatal(err)
	}
}
