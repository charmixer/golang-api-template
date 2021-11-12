package main

import (
	"github.com/charmixer/golang-api-template/cmd"
)

var (
	name        string = "golang-api-template"
	version     string = "0.0.0"
	environment string = "development"
	commit      string
	date        string
	tag         string
)

func main() {
	cmd.Execute(name, version, commit, date, tag, environment)
}
