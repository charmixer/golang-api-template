package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/creasty/defaults"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/charmixer/envconfig"
	"github.com/charmixer/go-flags"
)

type App struct {
	Name        string `long:"app-name" description:"Name of application" default:"golang-api-template"`
	Description string `long:"app-description" description:"Description of application" default:"Gives a simple blueprint for creating new api's"`
	Environment string `long:"environment" description:"Environment the application is running in, eg. prod or dev" default:"dev"`
	Version     string `long:"version" description:"Version of the application" default:"0.0.0"`

	Log struct {
		Verbose bool   `long:"verbose" short:"v" description:"Verbose logging"`
		Format  string `long:"log-format" description:"Logging format" choice:"json" choice:"plain"`
	}

	Serve serveCmd `command:"serve" description:"serves endpoints"`
	Oas   oasCmd   `command:"oas" description:"Retrieve oas document"`
}

var Application App
var parser = flags.NewParser(&Application, flags.HelpFlag|flags.PassDoubleDash)

func Execute() {
	_, err := parser.Execute()

	if err != nil {
		e := err.(*flags.Error)
		if e.Type != flags.ErrCommandRequired && e.Type != flags.ErrHelp {
			fmt.Printf("%s\n", e.Message)
		}
		parser.WriteHelp(os.Stdout)
	}

	os.Exit(0)
}

func init() {
	// 3. Priority: Config file
	parseYamlFile(os.Getenv("CFG_PATH"), &Application)

	// 2. Priority: Environment
	parseEnv("CFG", &Application)

	// 1. Priority: Flags
	parseFlags(&Application)

	// 0. Priority: Defaults, if none of above is found
	parseDefaults(&Application)

	initLogging()
}

func parseYamlFile(file string, config interface{}) {
	if file == "" {
		return
	}

	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}
}

func parseEnv(prefix string, config interface{}) {
	if err := envconfig.Process(prefix, config); err != nil {
		panic(err)
	}
}

func parseFlags(config interface{}) {
	if err := parser.ParseFlags(); err != nil {
		e := err.(*flags.Error)
		if e.Type != flags.ErrCommandRequired && e.Type != flags.ErrHelp {
			fmt.Printf("%s\n", e.Message)
		}
		parser.WriteHelp(os.Stdout)
	}
}

func parseDefaults(config interface{}) {
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
