package serve
/*
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
	"github.com/spf13/cobra"
)

var OpenAPISpec string

// RunServe the main event loop for the service
func RunServe() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {

		ip, err := cmd.Flags().GetString("ip")
		if err != nil {
			log.Warn().Msg("Missing ip")
			os.Exit(1)
		}
		app.Env.Ip = ip

		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			log.Warn().Msg("Missing port")
			os.Exit(1)
		}
		app.Env.Port = port

		domain, err := cmd.Flags().GetString("domain")
		if err != nil {
			log.Warn().Msg("Missing domain")
			os.Exit(1)
		}
		app.Env.Domain = domain

		readTimeout, err := cmd.Flags().GetInt("read-timeout")
		if err != nil {
			log.Warn().Msg("Missing read-timeout")
			os.Exit(1)
		}

		readHeaderTimeout, err := cmd.Flags().GetInt("read-header-timeout")
		if err != nil {
			log.Warn().Msg("Missing read-header-timeout")
			os.Exit(1)
		}

		writeTimeout, err := cmd.Flags().GetInt("write-timeout")
		if err != nil {
			log.Warn().Msg("Missing write-timeout")
			os.Exit(1)
		}

		idleTimeout, err := cmd.Flags().GetInt("idle-timeout")
		if err != nil {
			log.Warn().Msg("Missing idle-timeout")
			os.Exit(1)
		}

		gracefulTimeout, err := cmd.Flags().GetInt("graceful-timeout")
		if err != nil {
			log.Warn().Msg("Missing graceful-timeout")
			os.Exit(1)
		}

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
	}
}
*/
