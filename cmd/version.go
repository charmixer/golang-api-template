package cmd

import (
	"fmt"

	"github.com/charmixer/golang-api-template/app"
)

type versionCmd struct{}

func (v *versionCmd) Execute(args []string) error {
	fmt.Printf("Name: %s\nVersion: %s\nTag: %s\nCommit: %s\nDate: %s\nEnvironment: %s\n", app.Env.Build.Name, app.Env.Build.Version, app.Env.Build.Tag, app.Env.Build.Commit, app.Env.Build.Date, app.Env.Build.Environment)
	return nil
}
