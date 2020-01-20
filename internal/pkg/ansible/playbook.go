//go:generate mockgen -source=playbook.go -destination=../mocks/mock_commandexecutor.go -package=mocks CommandExecutor
package ansible

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type (
	PlaybookExecutor struct {
		commandExecutor        CommandExecutor
		playbookBinPath        string
		playbookRepositoryPath string
		args                   PlaybookArgs
		stdout                 io.Writer
		debug                  bool
	}

	PlaybookArgs struct {
		Inventories                    []string
		HideDiff, CheckModeDeactivated bool
		BecomeSudo                     bool
		OtherArgs                      []string
	}

	CommandExecutor interface {
		Run(command string, params ...string) error
	}
)

func InitPlaybookExecutor(cmdExecutor CommandExecutor, playbookBinPath, playbookPath string,
	args PlaybookArgs, stdout io.Writer, debug bool) *PlaybookExecutor {
	return &PlaybookExecutor{cmdExecutor, playbookBinPath, playbookPath, args, stdout, debug}
}

func (pE PlaybookExecutor) Play(playbookName string) error {
	if len(pE.args.Inventories) == 0 {
		return fmt.Errorf("playbook: no inventory to work on")
	}

	if pE.debug {
		pE.printInventories()
	}

	playbookPath, err := pE.computePlaybookPath(playbookName)
	if err != nil {
		return err
	}

	for _, inventory := range pE.args.Inventories {
		fmt.Fprintf(pE.stdout, "Start running %s playbook on %s inventory...\n", playbookName, inventory)
		err = pE.play(playbookPath, inventory)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pE PlaybookExecutor) play(playbookPath string, inventory string) error {
	if len(inventory) == 0 {
		return fmt.Errorf("playbook: inventory is empty")
	}

	params := append(pE.args.combineWithInventory(inventory), playbookPath)
	err := pE.commandExecutor.Run(pE.playbookBinPath, params...)
	if err != nil {
		return err
	}

	return nil
}

func (pE PlaybookExecutor) computePlaybookPath(name string) (string, error) {
	playbookPath := filepath.Join(pE.playbookRepositoryPath, name)
	info, err := os.Stat(playbookPath)

	if os.IsNotExist(err) {
		return "", fmt.Errorf("playbook: playbook does not exists:%v", playbookPath)
	}

	if info.IsDir() {
		return "", fmt.Errorf("playbook: are you kidding? it's a directory:%v", playbookPath)
	}

	return playbookPath, nil
}

func (pA PlaybookArgs) combineWithInventory(inventory string) []string {
	var result []string
	if !pA.HideDiff {
		result = append(result, "--diff")
	}
	if !pA.CheckModeDeactivated {
		result = append(result, "--check")
	}
	if pA.BecomeSudo {
		result = append(result, "-b")
	}

	result = append(result, "--inventory", inventory)

	result = append(result, pA.OtherArgs...)

	return result
}

func (pE PlaybookExecutor) printInventories() {
	for _, inv := range pE.args.Inventories {
		fmt.Fprintln(pE.stdout, "inventory:", inv)
	}
}
