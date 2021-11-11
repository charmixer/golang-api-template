package app

import (
	oas "github.com/charmixer/oas/exporter"
)

type Environment struct {
	Ip     string
	Port   int
	Addr   string
	Domain string

	OpenAPI oas.Openapi
}

var Env Environment
