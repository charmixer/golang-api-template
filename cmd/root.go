package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/charmixer/envconfig"
	"github.com/charmixer/go-flags"

	"github.com/charmixer/golang-api-template/validation"
)

type App struct {
	Description string `long:"app-description" description:"Description of application" default:"Gives a simple blueprint for creating new api's"`

	Log struct {
		Verbose bool   `long:"verbose" short:"v" description:"Verbose logging"`
		Format  string `long:"log-format" description:"Logging format" choice:"json" choice:"plain"`
	}

	Serve   serveCmd   `command:"serve" description:"serves endpoints"`
	Oas     oasCmd     `command:"oas" description:"Retrieve oas document"`
	Version versionCmd `command:"version" description:"Prints the build information from the binary"`
}

var Application App
var parser = flags.NewParser(&Application, flags.HelpFlag|flags.PassDoubleDash)

func Execute() {

	cmd := parser.GetCommand()

	if cmd != nil {
		err := validation.Validate.Struct(cmd)
		if err != nil {
			switch e := err.(type) {
			case validator.ValidationErrors:
				for _, verr := range e {
					fmt.Printf("Validation failed for %s: %s\n", verr.StructNamespace(), verr.Translate(validation.Translation))
				}
			default:
				fmt.Printf("%s\n", e.Error())
			}

			os.Exit(1)
			return
		}
	}

	_, err := parser.Execute()

	retcode := 0
	if err != nil {
		switch e := err.(type) {
		case *flags.Error:
			if e.Type != flags.ErrCommandRequired && e.Type != flags.ErrHelp {
				fmt.Printf("%s\n", e.Message)
			}
		case error:
			fmt.Printf("%s\n", err.Error())
			retcode = 1
		}

		parser.WriteHelp(os.Stdout)
	}

	os.Exit(retcode)
}

func init() {
	// 4. Priority: Defaults, used if nothing in the chain overwrites
	parseDefaults(&Application)

	// 3. Priority: Config file
	parseYamlFile(os.Getenv("CFG_PATH"), &Application)

	// 2. Priority: Environment
	parseEnv("CFG", &Application)

	// 1. Priority: Flags
	parseFlags(&Application)

	initLogging()
}

func parseYamlFile(file string, config *App) {
	if file == "" {
		return
	}

	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		panic(err)
	}
}

func parseEnv(prefix string, config *App) {
	envconfig.MustProcess(prefix, config)
}

func parseFlags(config *App) {
	if err := parser.ParseFlags(); err != nil {
		e := err.(*flags.Error)
		if e.Type != flags.ErrCommandRequired && e.Type != flags.ErrHelp {
			fmt.Printf("%s\n", e.Message)
		}
		parser.WriteHelp(os.Stdout)
	}
}

func parseDefaults(config *App) {
	if err := defaults.Set(config); err != nil {
		panic(err)
	}
}

func initLogging() {
	if Application.Log.Verbose {
		log.Logger = log.With().Caller().Logger()
	}

	if Application.Log.Format == "plain" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.TimestampFieldName = "time"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "msg"

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if Application.Log.Verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
