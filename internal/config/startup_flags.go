package config

import "flag"

type StartupFlags struct {
	BaseURL            string
	ResultBaseURL      string
	JSONFilePath       string
	DatabaseConnString string
}

func ParseStartupFlags() StartupFlags {
	baseURLFlag := flag.String("a", "", "Base URL for the server")
	resultBaseURL := flag.String("b", "", "Base URL for the shortened URLs")
	jsonFilePath := flag.String("f", "", "Base URL for the shortened URLs")
	databaseConnString := flag.String("d", "", "Database connection string")

	flag.Parse()

	return StartupFlags{
		BaseURL:            *baseURLFlag,
		ResultBaseURL:      *resultBaseURL,
		JSONFilePath:       *jsonFilePath,
		DatabaseConnString: *databaseConnString,
	}
}
