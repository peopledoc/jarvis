package cmd

import (
	"bufio"
	"jarvis/internal/pkg/environment"
	"os"

	"github.com/spf13/cobra"
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
		//For the moment we have only one available printer... tree
		w := bufio.NewWriter(os.Stdout)
		printer := environment.TreePrinter{}
		printer.PrintEnvironments(w, environments)
		w.Flush()
		return nil
	},
}
