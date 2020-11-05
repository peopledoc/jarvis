package cmd

import (
	"jarvis/internal/pkg/interactivity"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(eventCmd)

	eventCmd.AddCommand(sendCmd)
}

//usage: jarvis event
var eventCmd = &cobra.Command{
	Use:     "event",
	Aliases: []string{"ev"},
	Short:   "Event command",
	Args:    cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		//Placeholder... calling help to show usage...
		return cmd.Help()
	},
}

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send Event module",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		event := make(map[string]string)
		pYN := interactivity.InitPromptYesNo("", "Do you want to send an event?", os.Stdin, os.Stdout)

		pTitle := interactivity.InitPromptString("title", "Event's title", os.Stdin, os.Stdout)
		pYN.SetNext(pTitle)

		pText := interactivity.InitPromptString("text", "Event's text", os.Stdin, os.Stdout)
		pTitle.SetNext(pText)

		pPriority := interactivity.InitPromptSelect("priority", "Event's priority", []string{"low", "normal"}, os.Stdin, os.Stdout)
		pText.SetNext(pPriority)

		err := pYN.Execute(event)
		if err != nil {
			return err
		}
		return nil
	},
}
