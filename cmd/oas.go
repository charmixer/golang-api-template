package cmd

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/charmixer/golang-api-template/app"
	"github.com/charmixer/golang-api-template/router"
	"github.com/charmixer/oas/exporter"
)

type oasCmd struct {
	// version   bool `short:"v" long:"version" description:"display version"`
}

func (v *oasCmd) Execute(args []string) error {
	router := router.NewRouter(app.Env.Build.Name, Application.Description, app.Env.Build.Version)

	oasModel := exporter.ToOasModel(router.OpenAPI)
	oasYaml, err := yaml.Marshal(&oasModel)
	if err != nil {
		return err
	}

	fmt.Println(string(oasYaml))

	return nil
}
