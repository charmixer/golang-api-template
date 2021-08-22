package cmd

import (
	"fmt"
)

type Serve struct {
	Version bool `short:"v" long:"version" description:"display version"`
	Port int
}

func (v *Serve) Execute(args []string) error {
	fmt.Println("servecmd")
	fmt.Printf("%#v\n", v)
	fmt.Printf("%#v\n", Application.Config)

	return nil
}

/*
import (
	"github.com/charmixer/golang-api-template/pkg/serve"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve endpoints",
	Long:  `Serve endpoints`,
	Run:   serve.RunServe(),
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	serveCmd.Flags().IntP("port", "p", 8080, "The port used for serving the api.")
	serveCmd.Flags().StringP("ip", "i", "0.0.0.0", "The ip used for serving the api.")
	serveCmd.Flags().StringP("domain", "d", "localhost", "The domain used to access the api.")

	serveCmd.Flags().IntP("write-timeout", "", 10, "Timeout in seconds when writing response.")
	serveCmd.Flags().IntP("read-timeout", "", 10, "Timeout in seconds when reading request headers and body.")
	serveCmd.Flags().IntP("read-header-timeout", "", 5, "Timeout in seconds when reading request headers.")
	serveCmd.Flags().IntP("idle-timeout", "", 15, "Timeout in seconds between requests when keep-alive is enabled. If 0 read-timeout is used.")
	serveCmd.Flags().IntP("graceful-timeout", "", 15, "Timeout in seconds when shutting down.")

	rootCmd.AddCommand(serveCmd)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}*/
