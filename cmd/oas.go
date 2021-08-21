package cmd

import (
	"fmt"
)

type OasCmd struct {
	// version   bool `short:"v" long:"version" description:"display version"`
}

func (v *OasCmd) Execute(args []string) error {
	fmt.Println("servecmd")
	fmt.Printf("%#v\n", v)
	fmt.Printf("%#v\n", Application.Config)

	return nil
}


/*
package cmd

import (
	"github.com/charmixer/golang-api-template/pkg/oas"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var oasCmd = &cobra.Command{
	Use:   "oas",
	Short: "Export OpenAPI spec",
	Long:  `Export OpenAPI spec`,
	Run:   oas.RunOas(),
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// oasCmd.Flags().StringP("ip", "i", "0.0.0.0", "The ip used for serving the api.")

	rootCmd.AddCommand(oasCmd)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}*/
