package cmd

import (
	"os"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"

	"github.com/kelseyhightower/envconfig"
	"github.com/jessevdk/go-flags"
)

type App struct {
	Cli struct {
		Serve `command:"serve" description:"serves the application"`
	}

	Config struct {
		Log struct {
			Debug  bool
			Format string
		}

		Serve struct {
			Public struct {
				Port int `short:"p" long:"port" description:"Port to serve app on"`
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

	}
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
}

func init() {

	configFile := os.Getenv("CFG_PATH")

	files := []string{}
	if configFile != "" {
		files = append(files, configFile)
	}

	for _, file := range files {
		yamlFile, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(yamlFile, &Application.Config)
		if err != nil {
			panic(err)
		}
	}

	err := envconfig.Process("CFG", &Application.Config)
  if err != nil {
		panic(err)
  }
}

/*
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var cfgFile string
var verbose bool
var console bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "golang-api-template",
	Short: "Template api",
	Long:  `Template api`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initLogging)
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.golang-api-template.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
	rootCmd.PersistentFlags().BoolVarP(&console, "console", "c", false, "enable human readable console output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Error().Err(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".golang-api-template" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".golang-api-template")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info().Msgf("Using config file: %s", viper.ConfigFileUsed())
	}
}

func initLogging() {

	if verbose {
		log.Logger = log.With().Caller().Logger()
	}

	if console {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.TimestampFieldName = "time"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "msg"

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
*/
