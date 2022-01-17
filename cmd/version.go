package cmd

import (
	"fmt"

	"github.com/charmixer/golang-api-template/env"
)

type versionCmd struct{}

func (v *versionCmd) Execute(args []string) error {
	fmt.Printf("Name: %s\nVersion: %s\nTag: %s\nCommit: %s\nDate: %s\nEnvironment: %s\n", env.Env.Build.Name, env.Env.Build.Version, env.Env.Build.Tag, env.Env.Build.Commit, env.Env.Build.Date, env.Env.Build.Environment)
	return nil
}
