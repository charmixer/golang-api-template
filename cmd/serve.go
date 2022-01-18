package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/charmixer/golang-api-template/env"
	"github.com/charmixer/golang-api-template/router"

	"github.com/charmixer/oas/exporter"

	"github.com/charmixer/golang-api-template/tracing"

	"github.com/rs/zerolog/log"
)

type serveCmd struct {
	Tracing struct {
		Disabled bool   `long:"trace-disable" description:"Disable tracing"`
		Url      string `long:"trace-provider-url" description:"Trace provider endpoint to use instead of default"`
		Provider string `long:"trace-provider" description:"Provider to use for tracing" choice:"jaeger" default:"jaeger"`
	}
	Public struct {
		Port   int    `short:"p" long:"port" description:"Port to serve app on" default:"8080"`
		Ip     string `short:"i" long:"ip" description:"IP to serve app on" default:"0.0.0.0"`
		Domain string `short:"d" long:"domain" description:"Domain to access app through" default:"127.0.0.1"`
	}
	Timeout struct {
		Write      int `long:"write-timeout" description:"Timeout in seconds for write" default:"10"`
		Read       int `long:"read-timeout" description:"Timeout in seconds for read" default:"5"`
		ReadHeader int `long:"read-header-timeout" description:"Timeout in seconds for read-header" default:"5"`
		Idle       int `long:"idle-timeout" description:"Timeout in seconds for idle" default:"10"`
		Grace      int `long:"grace-timeout" description:"Timeout in seconds before shutting down" default:"15"`
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

func (cmd *serveCmd) initTracing() func() {
	if cmd.Tracing.Disabled || cmd.Tracing.Provider == "" {
		log.Debug().Msgf("Tracing is disabled")
		return nil
	}

	var err error

	exporter := tracing.SetupNilExporter()
	if cmd.Tracing.Provider == "jaeger" {
		exporter, err = tracing.SetupJaegerExporter(cmd.Tracing.Url)
	}

	if err != nil {
		log.Error().Err(err).Msg("Unable to setup trace exporter")
		return nil
	}

	if exporter == nil {
		log.Debug().Msg("No exporter was setup for tracing")
		return nil
	}

	shutdownTracing, err := tracing.SetupTracing(exporter, env.Env.Build.Name, env.Env.Build.Environment, env.Env.Build.Version)
	if err == nil {
		return shutdownTracing
	}

	// Deny by default
	log.Error().Err(err).Msg("Failed to setup tracer")
	return nil
}

func (cmd *serveCmd) Execute(args []string) error {
	env.Env.Ip = cmd.Public.Ip
	env.Env.Port = cmd.Public.Port
	env.Env.Domain = cmd.Public.Domain
	env.Env.Addr = fmt.Sprintf("%s:%d", env.Env.Ip, env.Env.Port)

	shutdown := cmd.initTracing()
	defer shutdown()

	router := router.NewRouter(env.Env.Build.Name, Application.Description, env.Env.Build.Version)

	oasModel := exporter.ToOasModel(
		router.OpenAPI,
		exporter.WithQueryTag("query"),
		exporter.WithHeaderTag("header"),
		exporter.WithCookieTag("cookie"),
		exporter.WithDescriptionTag("description"),
		exporter.WithValidationTag("validation"),
	)
	env.Env.OpenAPI = oasModel

	srv := &http.Server{
		Addr: env.Env.Addr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout:      time.Second * time.Duration(cmd.Timeout.Write),
		ReadTimeout:       time.Second * time.Duration(cmd.Timeout.Read),
		ReadHeaderTimeout: time.Second * time.Duration(cmd.Timeout.ReadHeader),
		IdleTimeout:       time.Second * time.Duration(cmd.Timeout.Idle),
		Handler:           router.Handle(),
	}

	exitCode := 0
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Info().Msg("Listening on " + env.Env.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Error().Err(err).Msg("Failed to open connection")
			exitCode = 1
			c <- os.Interrupt
		}
	}()

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
	os.Exit(exitCode)

	return nil
}
