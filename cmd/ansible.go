package cmd

import (
	"fmt"
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
	CheckModeEnabled     bool
	BecomeSudo           bool
)

var (
	ModuleName  string
	ModuleArg   string
	HostPattern string
)

func init() {
	rootCmd.AddCommand(ansibleCmd)
	ansibleCmd.AddCommand(playCmd)
	ansibleCmd.PersistentFlags().BoolVar(&HideDiff, "nodiff", false,
		"Hide diff")
	ansibleCmd.PersistentFlags().BoolVarP(&BecomeSudo, "become", "b", false,
		"Become sudo")
	playCmd.Flags().BoolVar(&CheckModeDeactivated, "nocheck", false,
		"Deactivate check mode")

	ansibleCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&ModuleName, "module", "m", "shell", "Ansible module name")
	runCmd.Flags().StringVarP(&ModuleArg, "args", "a", "", "Ansible module name")
	runCmd.Flags().StringVarP(&HostPattern, "target", "t", "", "Ansible host-pattern")
	runCmd.Flags().BoolVar(&CheckModeEnabled, "check", false,
		"Enable check mode")
	runCmd.MarkFlagRequired("target")
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

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run Ansible module",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		cmdWorkingDir := viper.GetString("ansible.working_directory")
		envsPath := viper.GetString("environments.path")

		inventories, err := environment.GetFullPathInventoriesFromEnvironments(envsPath, *environments)
		if err != nil {
			return err
		}
		ModuleArgs := ansible.RunModuleArgs{
			ModuleName:  ModuleName,
			ModuleArg:   ModuleArg,
			HostPattern: HostPattern,
			CommonArgs: ansible.CommonArgs{
				Inventories:      inventories,
				HideDiff:         HideDiff,
				BecomeSudo:       BecomeSudo,
				CheckModeEnabled: CheckModeEnabled,
				OtherArgs:        args,
			},
		}
		if isDebug {
			fmt.Printf("Ansible module args: %v\n", args)
		}
		ansibleBinPath := viper.GetString("ansible.run.bin_path")

		runner := command.Init(os.Stdout, os.Stderr, cmdWorkingDir, isDebug)
		runExecutor := ansible.InitRunModuleExecutor(
			runner, ansibleBinPath, ModuleArgs, os.Stdout, isDebug)

		err = runExecutor.Run()
		if err != nil {
			return err
		}

		return nil
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
		cmdWorkingDir := viper.GetString("ansible.working_directory")
		envsPath := viper.GetString("environments.path")

		inventories, err := environment.GetFullPathInventoriesFromEnvironments(envsPath, *environments)
		if err != nil {
			return err
		}

		if CheckModeDeactivated {
			CheckModeEnabled = false
		}
		playbookArgs := ansible.CommonArgs{
			Inventories:      inventories,
			HideDiff:         HideDiff,
			BecomeSudo:       BecomeSudo,
			CheckModeEnabled: CheckModeEnabled,
			OtherArgs:        args[1:],
		}
		playbookBinPath := viper.GetString("ansible.playbook.bin_path")
		playbooksPath := viper.GetString("ansible.playbook.playbooks_path")

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
