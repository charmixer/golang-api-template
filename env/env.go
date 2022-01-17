package env

import (
	oas "github.com/charmixer/oas/exporter"
)

type Environment struct {
	Ip     string
	Port   int
	Addr   string
	Domain string

	Build struct {
		Name        string
		Version     string
		Commit      string
		Date        string
		Tag         string
		Environment string // eg. prod or dev
	}

	OpenAPI oas.Openapi
}

var Env Environment
