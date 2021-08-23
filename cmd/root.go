package cmd

import (
	"os"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/charmixer/envconfig"
	"github.com/charmixer/go-flags"
)

type App struct {
	Log struct {
		Verbose  bool `long:"verbose" description:"Verbose logging"`
		Format string `long:"log-format" description:"Logging format" choice:"json" choice:"plain"`
	}

	Serve `command:"serve" description:"serves endpoints"`
	Oas   `command:"oas" description:"Retrieve oas document"`
}

var Application App

func Execute(){
	var parser = flags.NewParser(&Application, flags.HelpFlag | flags.PassDoubleDash)
	_,err := parser.Parse()

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

	// Only option which isn't available as normal - path to config file
	configFile := os.Getenv("CFG_PATH")

	// Collect config files (just preparation for multiple files)
	files := []string{}
	if configFile != "" {
		files = append(files, configFile)
	}

	// Parse all config files into struct (config files has lowest priority)
	for _, file := range files {
		yamlFile, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(yamlFile, &Application)
		if err != nil {
			panic(err)
		}
	}

	// Parse environment into struct (2. priority, flags has priority 1)
	err := envconfig.Process("CFG", &Application)
  if err != nil {
		panic(err)
  }

	initLogging() // FIXME too early, haven't parsed flags
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
