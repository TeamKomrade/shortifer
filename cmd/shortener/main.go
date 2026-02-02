package main

import (
	cfg "github.com/TeamKomrade/shortifer/internal/config"
	httpserver "github.com/TeamKomrade/shortifer/internal/service/httpserver"
)

func main() {
	flags := cfg.ParseStartupFlags()
	httpserver.CreateServer(flags)
}
