package app

type Environment struct {
	Ip string
	Port int
	Addr string
	Domain string

	OpenAPI string
}

var Env Environment
