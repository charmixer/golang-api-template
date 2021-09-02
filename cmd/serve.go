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

	"github.com/rs/zerolog/log"
)

type ServeCmd struct {
	Public struct {
		Port int `short:"p" long:"port" description:"Port to serve app on" default:"8080"`
		Ip string `short:"i" long:"ip" description:"IP to serve app on" default:"0.0.0.0"`
		Domain string `short:"d" long:"domain" description:"Domain to access app through" default:"127.0.0.1"`
	}
	Timeout struct {
		Write int `long:"write-timeout" description:"Timeout in seconds for write" default:"10"`
		Read int `long:"read-timeout" description:"Timeout in seconds for read" default:"5"`
		ReadHeader int `long:"read-header-timeout" description:"Timeout in seconds for read-header" default:"5"`
		Idle int `long:"idle-timeout" description:"Timeout in seconds for idle" default:"10"`
		Grace int `long:"grace-timeout" description:"Timeout in seconds before shutting down" default:"15"`
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

func (cmd *ServeCmd) Execute(args []string) error {
	app.Env.Ip     = cmd.Public.Ip
	app.Env.Port   = cmd.Public.Port
	app.Env.Domain = cmd.Public.Domain
	app.Env.Addr   = fmt.Sprintf("%s:%d", app.Env.Ip, app.Env.Port)

	oas := router.NewOas()
	oasModel := exporter.ToOasModel(oas)
	spec, err := exporter.ToYaml(oasModel)
	if err != nil {
		log.Error().Err(err)
	}
	app.Env.OpenAPI = spec

	chain := middleware.GetChain()
	route := router.NewRouter(oas)

	srv := &http.Server{
		Addr: app.Env.Addr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout:      time.Second * time.Duration(cmd.Timeout.Write),
		ReadTimeout:       time.Second * time.Duration(cmd.Timeout.Read),
		ReadHeaderTimeout: time.Second * time.Duration(cmd.Timeout.ReadHeader),
		IdleTimeout:       time.Second * time.Duration(cmd.Timeout.Idle),
		Handler:           chain.Then(route), // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Info().Msg("Listening on " + app.Env.Addr)
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(cmd.Timeout.Idle))
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
