package cmd

import (
	"fmt"
	"io"
	"jarvis/internal/pkg/ansible"
	"jarvis/internal/pkg/command"
	"jarvis/internal/pkg/environment"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//Ansible flags
var (
	HideDiff   bool
	BecomeSudo bool
)

//Playbook flags
var (
	CheckModeDeactivated bool
	JoinInventories      bool
)

//Run flags
var (
	CheckModeEnabled bool
	ModuleName       string
	ModuleArg        string
	HostPattern      string
)

//Inventory flags
var (
	ListGroup     bool
	HostGroupName string
	WithParent    bool
)

func init() {
	rootCmd.AddCommand(ansibleCmd)
	//AnsibleCmd
	ansibleCmd.PersistentFlags().BoolVar(&HideDiff, "nodiff", false,
		"Hide diff")
	ansibleCmd.PersistentFlags().BoolVarP(&BecomeSudo, "become", "b", false,
		"Become sudo")

	//PlaybookCmd
	playCmd.Flags().BoolVar(&CheckModeDeactivated, "nocheck", false,
		"Deactivate check mode")
	playCmd.Flags().BoolVar(&JoinInventories, "join-inventories", false,
		"Join platforms inventories")
	ansibleCmd.AddCommand(playCmd)

	//RunCmd
	runCmd.Flags().StringVarP(&ModuleName, "module", "m", "shell", "Ansible module name")
	runCmd.Flags().StringVarP(&ModuleArg, "args", "a", "", "Ansible module arg")
	runCmd.Flags().StringVarP(&HostPattern, "target", "t", "", "Ansible host-pattern")
	runCmd.Flags().BoolVar(&CheckModeEnabled, "check", false,
		"Enable check mode")
	runCmd.MarkFlagRequired("target")
	ansibleCmd.AddCommand(runCmd)

	//InventoryCmd
	inventoryCmd.Flags().BoolVarP(&WithParent, "with-parent", "W", false, "Query with parent group")
	inventoryCmd.Flags().BoolVarP(&ListGroup, "group", "G", false, "List group name (mutually exclusive with --host)")
	inventoryCmd.Flags().StringVarP(&HostGroupName, "host", "H", "", "List host by group name (mutually exclusive with --group)")
	ansibleCmd.AddCommand(inventoryCmd)
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

		inventories, err := environment.GetFullPathInventoriesFromEnvironments(envsPath, environments)
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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if JoinInventories && !environment.IsPredicateJoinInventoriesSyntaxValid(envName) {
			return fmt.Errorf("With join-inventories set, environment flag must have 'type' and 'env' set. Ex:'type.env'")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmdWorkingDir := viper.GetString("ansible.working_directory")
		envsPath := viper.GetString("environments.path")

		inventories, err := environment.GetFullPathInventoriesFromEnvironments(envsPath, environments)
		if err != nil {
			return err
		}

		playbookArgs := ansible.CommonArgs{
			Inventories:      inventories,
			HideDiff:         HideDiff,
			BecomeSudo:       BecomeSudo,
			CheckModeEnabled: !CheckModeDeactivated,
			JoinInventories:  JoinInventories,
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

//usage: jarvis ansible list
//returns:
// --group(bool) >> list groups name
// --hosts(string) >> hosts by group
var inventoryCmd = &cobra.Command{
	Use:   "list",
	Short: "Query inventory",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//checking exclusivity
		if ListGroup && HostGroupName != "" {
			return fmt.Errorf("--group and --hosts are mutually exclusive")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		invReaders, err := inventoriesReaders(environments)
		if err != nil {
			return err
		}

		//to concatenate all files to one reader
		r := io.MultiReader(invReaders...)
		manipulator, err := ansible.InitInventoryManipulator(r)
		if err != nil {
			return err
		}

		if ListGroup {
			groups, err := manipulator.GetGroupsName(WithParent)
			if err != nil {
				return err
			}
			//I know it is a bit odd but Jarvis is learning
			//be nice with him :)
			for _, g := range groups {
				fmt.Println(g)
			}

			return nil
		}

		if HostGroupName != "" {
			hosts, err := manipulator.GetHostsByGroupName(HostGroupName)
			if err != nil {
				return err
			}
			for _, h := range hosts {
				fmt.Println(h)
			}

			return nil
		}

		return nil
	},
}

func inventoriesReaders(envs []*environment.Environment) ([]io.Reader, error) {
	envsPath := viper.GetString("environments.path")

	allInventories, err := environment.GetFullPathInventoriesFromEnvironments(envsPath, envs)
	if err != nil {
		return nil, err
	}

	invReaders, err := ansible.BuildReadersFromInventoriesPath(allInventories)
	if err != nil {
		return nil, err
	}

	return invReaders, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
