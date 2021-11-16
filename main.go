package main

import (
	"github.com/charmixer/golang-api-template/app"
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

func init() {
	app.Env.Build.Name = name
	app.Env.Build.Version = version
	app.Env.Build.Tag = tag
	app.Env.Build.Commit = commit
	app.Env.Build.Date = date
	app.Env.Build.Environment = environment
}

func main() {
	cmd.Execute()
}
