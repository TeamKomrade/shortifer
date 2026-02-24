package config

import "flag"

type StartupFlags struct {
	BaseURL       string
	ResultBaseURL string
}

func ParseStartupFlags() StartupFlags {
	baseURLFlag := flag.String("a", "", "Base URL for the server")
	resultBaseURL := flag.String("b", "", "Base URL for the shortened URLs")
	flag.Parse()

	return StartupFlags{
		BaseURL:       *baseURLFlag,
		ResultBaseURL: *resultBaseURL,
	}
}
