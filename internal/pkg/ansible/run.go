package ansible

import (
	"fmt"
	"io"
	"strings"
)

type (
	RunModuleExecutor struct {
		commandExecutor CommandExecutor
		binPath         string
		stdout          io.Writer
		debug           bool
	}

	RunModuleArgs struct {
		CommonArgs
		ModuleName  string
		ModuleArg   string
		HostPattern string
	}

	RunModule struct {
		RunModuleExecutor
		RunModuleArgs
	}
)

func InitRunModuleExecutor(cmdExecutor CommandExecutor, binPath string, runModuleArgs RunModuleArgs,
	stdout io.Writer, debug bool) *RunModule {

	return &RunModule{
		RunModuleExecutor: RunModuleExecutor{cmdExecutor, binPath, stdout, debug},
		RunModuleArgs:     runModuleArgs,
	}
}

func (runMod RunModule) Run() error {

	if len(runMod.Inventories) == 0 {
		return fmt.Errorf("run: no inventory to work on")
	}

	if runMod.debug {
		runMod.printInventories(runMod.stdout)
	}

	if isInventoriesEmpty(runMod.Inventories) {
		return fmt.Errorf("run: inventory is empty")
	}

	for _, inventory := range runMod.Inventories {
		if runMod.debug {
			fmt.Fprintf(runMod.stdout, "Start running %s module on %s inventory...\n",
				runMod.ModuleName, inventory)
			fmt.Fprintf(runMod.stdout, strings.Join(runMod.computeAnsibleOptions(inventory), " "))
		}
		err := runMod.commandExecutor.Run(runMod.binPath, runMod.computeAnsibleOptions(inventory)...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (runMod RunModule) computeAnsibleOptions(inventory []string) []string {
	var result = runMod.buildCommonArgs()
	result = append(result, buildInventoryArg(inventory)...)
	result = append(result, "-m", runMod.ModuleName)
	if len(runMod.ModuleArg) > 0 {
		result = append(result, "--args", runMod.ModuleArg)
	}
	result = append(result, runMod.HostPattern)

	return result
}
