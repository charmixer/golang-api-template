package oas

import (
	"os"
	"fmt"

	"github.com/charmixer/oas/exporter"

	"github.com/charmixer/golang-api-template/router"

	// "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// RunServe the main event loop for the service
func RunOas() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {

		oas := router.NewOas()

		oasModel := exporter.ToOasModel(oas)
		oasYaml, err := exporter.ToYaml(oasModel)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(oasYaml)

		os.Exit(0)
	}
}
