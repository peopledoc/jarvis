package cmd

import (
	"fmt"
	"jarvis/internal/pkg/config"
	"jarvis/internal/pkg/environment"
	"jarvis/internal/pkg/funquotes"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const SilenceUsageOnError = true

var (
	configFile   string
	envName      string
	environments *environment.Environments
	isDebug      bool

	rootCmd = &cobra.Command{
		Use:          "jarvis",
		SilenceUsage: SilenceUsageOnError,
		Short:        "jarvis is our automation CLI",
		Long: `jarvis is the main command, used to facilitate SRE's life.
		
jarvis is smart, jarvis is beautiful,
jarvis is made with love by the SRE team.`,

		//PersistentPreRunE is inherited for children commands :)
		PersistentPreRunE: rootPersistentPreRunE,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(funquotes.GiveMeAQuote())
			cmd.Help()
		},
	}
)

func init() {
	rootCmd.PersistentFlags().
		StringVar(&configFile, "config", config.DEFAULT_CONFIG_PATH, "config file path")
	rootCmd.PersistentFlags().
		StringVarP(&envName, "env", "e", "", "Environment name, syntax(type.env.platform)")
	rootCmd.PersistentFlags().
		BoolVar(&isDebug, "debug", false, "debug mode")

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	configReader := config.InitConfigurationReader(configFile)

	err := configReader.ParseConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error while parsing config file: %s", err))
	}
}

func rootPersistentPreRunE(cmd *cobra.Command, args []string) error {
	helperPath := path.Join(
		viper.GetString("environments.path"), viper.GetString("environments.helper"))

	envs, err := environment.ParseEnvironmentFile(helperPath)
	if err != nil {
		return fmt.Errorf("fatal error while listing inventories: %s", err)
	}

	//Do we need to handle a specific environment?
	if len(envName) > 0 {
		choosenEnv, err := handleEnvArgument(envName, envs)
		if err != nil {
			return err
		}
		envs = &environment.Environments{choosenEnv}
	}

	environments = envs
	return nil
}

func handleEnvArgument(envName string, envs *environment.Environments) (*environment.Environment, error) {
	envPredicate, err := environment.ParsePredicate(envName)
	if err != nil {
		return nil, err
	}
	choosenEnv, err := environment.FindEnvironmentTreeFromPredicate(envPredicate, envs)
	if err != nil {
		return nil, err
	}
	return choosenEnv, nil
}

func Execute() error {
	return rootCmd.Execute()
}
