package cmd

import (
	"bufio"
	"fmt"
	"jarvis/internal/pkg/environment"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(environmentCmd)
	environmentCmd.AddCommand(environmentListCmd)
}

//usage: jarvis environment
//returns nothing, just a placeholder for sub environment commands
var environmentCmd = &cobra.Command{
	Use:     "environment",
	Aliases: []string{"env"},
	Short:   "Environment command",
	RunE: func(cmd *cobra.Command, args []string) error {
		//Placeholder... calling help to show usage...
		return cmd.Help()
	},
}

var environmentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		helperPath := path.Join(
			viper.GetString("environments_path"), viper.GetString("environments_helper"))

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
		//For the moment we have only one available printer... tree
		printer := environment.TreePrinter{}
		printEnvsToOutput(envs, printer)
		return nil
	},
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

func printEnvsToOutput(envs *environment.Environments, printer environment.EnvironmentsPrinter) {
	w := bufio.NewWriter(os.Stdout)
	printer.PrintEnvironments(w, *envs)
	w.Flush()
}
