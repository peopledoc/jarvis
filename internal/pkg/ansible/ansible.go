package ansible

import (
	"fmt"
	"io"
)

type (
	CommandExecutor interface {
		Run(command string, params ...string) error
	}

	CommonArgs struct {
		Inventories          []string
		HideDiff, BecomeSudo bool
		CheckModeEnabled     bool
		OtherArgs            []string
	}
)

func (cargs CommonArgs) computeCommonArgsWithInventory(inventory string) []string {
	var result []string
	if cargs.BecomeSudo {
		result = append(result, "-b")
	}
	if !cargs.HideDiff {
		result = append(result, "--diff")
	}
	if cargs.CheckModeEnabled {
		result = append(result, "--check")
	}
	result = append(result, "--inventory", inventory)

	return result
}

func (cargs CommonArgs) printInventories(writer io.Writer) {
	for _, inv := range cargs.Inventories {
		fmt.Fprintln(writer, "inventory:", inv)
	}
}
