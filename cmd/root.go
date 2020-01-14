package cmd

import (
	"fmt"
	"jarvis/internal/pkg/config"
	"jarvis/internal/pkg/funquotes"

	"github.com/spf13/cobra"
)

const SilenceUsageOnError = true

var (
	configFile string
	envName    string

	rootCmd = &cobra.Command{
		Use:          "jarvis",
		SilenceUsage: SilenceUsageOnError,
		Short:        "jarvis is our automation CLI",
		Long: `jarvis is the main command, used to facilitate SRE's life.
		
jarvis is smart, jarvis is beautiful,
jarvis is made with love by the SRE team.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(funquotes.GiveMeAQuote())
		},
	}
)

func init() {
	rootCmd.PersistentFlags().
		StringVar(&configFile, "config", config.DEFAULT_CONFIG_PATH, "config file path")
	rootCmd.PersistentFlags().
		StringVarP(&envName, "env", "e", "", "Environment name")

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	configReader := config.InitConfigurationReader(configFile)

	err := configReader.ParseConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error while parsing config file: %s", err))
	}
}

func Execute() error {
	return rootCmd.Execute()
}
