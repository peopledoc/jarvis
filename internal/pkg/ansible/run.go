package ansible

import (
	"fmt"
	"io"
	"strings"
	//	"os"
	//	"path/filepath"
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

	if runMod.debug {
		runMod.printInventories(runMod.stdout)
	}
	for _, inventory := range runMod.Inventories {
		if runMod.debug {
			fmt.Fprintf(runMod.stdout, "Start running %s module on %s inventory...\n", runMod.ModuleName, inventory)
			fmt.Fprintf(runMod.stdout, strings.Join(runMod.ansibleOptions(inventory), " "))
		}
		err := runMod.commandExecutor.Run(runMod.binPath, runMod.ansibleOptions(inventory)...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (runMod RunModule) ansibleOptions(inventory string) []string {
	var result = runMod.computeCommonArgs(inventory)
	result = append(result, "-m", runMod.ModuleName)
	result = append(result, "--args", runMod.ModuleArg)
	result = append(result, runMod.HostPattern)
	result = append(result, runMod.OtherArgs...)

	return result
}
