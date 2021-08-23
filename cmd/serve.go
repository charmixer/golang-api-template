package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/charmixer/golang-api-template/middleware"
	"github.com/charmixer/golang-api-template/router"
	"github.com/charmixer/golang-api-template/app"

	"github.com/charmixer/oas/exporter"

	"github.com/justinas/alice"
	"github.com/rs/zerolog/log"
)

type Serve struct {
	Public struct {
		Port int `short:"p" long:"port" description:"Port to serve app on" default:"8080"`
		Ip string `short:"i" long:"ip" description:"IP to serve app on" default:"0.0.0.0"`
		Domain string `short:"d" long:"domain" description:"Domain to access app through" default:"127.0.0.1"`
	}
	Timeout struct {
		Write int `long:"write-timeout" description:"Timeout in seconds for write"`
		Read int `long:"read-timeout" description:"Timeout in seconds for read"`
		ReadHeader int `long:"read-header-timeout" description:"Timeout in seconds for read-header"`
		Idle int `long:"idle-timeout" description:"Timeout in seconds for idle"`
		Grace int `long:"grace-timeout" description:"Timeout in seconds before shutting down"`
	}
	TLS struct {
		Cert struct {
			Path string
		}
		Key struct {
			Path string
		}
	}
}

func (v *Serve) Execute(args []string) error {
	ip := Application.Serve.Public.Ip
	port := Application.Serve.Public.Port
	domain := Application.Serve.Public.Domain

	app.Env.Ip = ip
	app.Env.Port = port
	app.Env.Domain = domain

	readTimeout := Application.Timeout.Read
	readHeaderTimeout := Application.Timeout.ReadHeader
	writeTimeout := Application.Timeout.Write
	idleTimeout := Application.Timeout.Idle
	gracefulTimeout := Application.Timeout.Grace

	addr := fmt.Sprintf("%s:%d", ip, port)
	app.Env.Addr = addr

	oas := router.NewOas()
	oasModel := exporter.ToOasModel(oas)
	spec, err := exporter.ToYaml(oasModel)
	if err != nil {
		log.Error().Err(err)
	}
	app.Env.OpenAPI = spec
	route := router.NewRouter(oas)

	chain := alice.New(middleware.Context, middleware.Logging)

	srv := &http.Server{
		Addr: addr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout:      time.Second * time.Duration(writeTimeout),
		ReadTimeout:       time.Second * time.Duration(readTimeout),
		ReadHeaderTimeout: time.Second * time.Duration(readHeaderTimeout),
		IdleTimeout:       time.Second * time.Duration(idleTimeout),
		Handler:           chain.Then(route), // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Info().Msg("Listening on " + addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Error().Err(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(gracefulTimeout))
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Info().Msg("shutting down")
	os.Exit(0)

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
