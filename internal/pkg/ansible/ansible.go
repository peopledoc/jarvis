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
		Inventories          [][]string
		HideDiff, BecomeSudo bool
		CheckModeEnabled     bool
		JoinInventories      bool
		OtherArgs            []string
	}
)

func (cargs CommonArgs) buildCommonArgs() []string {
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

	result = append(result, cargs.OtherArgs...)

	return result
}

func buildInventoryArg(inventories []string) []string {
	var a []string
	for _, inv := range inventories {
		a = append(a, "--inventory", inv)
	}
	return a
}

func (cargs CommonArgs) joinInventories() []string {
	var result []string
	for _, inv := range cargs.Inventories {
		result = append(result, buildInventoryArg(inv)...)
	}
	return result
}

func (cargs CommonArgs) printInventories(writer io.Writer) {
	for _, inv := range cargs.Inventories {
		fmt.Fprintln(writer, "inventory:", inv)
	}
}

func isInventoriesEmpty(inventories [][]string) bool {
	if len(inventories) == 0 {
		return true
	}
	for _, invs := range inventories {
		if len(invs) == 0 {
			return true
		}
		for _, inv := range invs {
			if len(inv) == 0 {
				return true
			}

		}
	}

	return false
}
