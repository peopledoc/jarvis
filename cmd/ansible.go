package cmd

import (
	"jarvis/internal/pkg/ansible"
	"jarvis/internal/pkg/command"
	"jarvis/internal/pkg/environment"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	HideDiff             bool
	CheckModeDeactivated bool
	BecomeSudo           bool
)

func init() {
	rootCmd.AddCommand(ansibleCmd)
	ansibleCmd.AddCommand(playCmd)
	playCmd.Flags().BoolVar(&HideDiff, "nodiff", false,
		"Hide diff")
	playCmd.Flags().BoolVar(&CheckModeDeactivated, "nocheck", false,
		"Deactivate check mode")
	playCmd.Flags().BoolVarP(&BecomeSudo, "become", "b", false,
		"Become sudo")
}

//usage: jarvis ansible
//returns nothing, just a placeholder for sub ansible commands
var ansibleCmd = &cobra.Command{
	Use:     "ansible",
	Aliases: []string{"an"},
	Short:   "Ansible command",
	Args:    cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		//Placeholder... calling help to show usage...
		return cmd.Help()
	},
}

//usage: jarvis ansible play
//returns: output of ansible-playbook command
var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Play playbook",
	//playbook name is mandatory
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmdWorkingDir := viper.GetString("ansible_playbook.working_directory")
		envsPath := viper.GetString("environments.path")

		inventories, err := environment.GetFullPathInventoriesFromEnvironments(envsPath, *environments)
		if err != nil {
			return err
		}

		playbookArgs := ansible.PlaybookArgs{
			Inventories:          inventories,
			HideDiff:             HideDiff,
			CheckModeDeactivated: CheckModeDeactivated,
			BecomeSudo:           BecomeSudo,
			OtherArgs:            args[1:],
		}
		playbookBinPath := viper.GetString("ansible_playbook.bin_path")
		playbooksPath := viper.GetString("ansible_playbook.playbooks_path")

		runner := command.Init(os.Stdout, os.Stderr, cmdWorkingDir, isDebug)
		playbookExecutor := ansible.InitPlaybookExecutor(
			runner, playbookBinPath, playbooksPath, playbookArgs, os.Stdout, isDebug)

		err = playbookExecutor.Play(args[0])
		if err != nil {
			return err
		}

		return nil
	},
}
