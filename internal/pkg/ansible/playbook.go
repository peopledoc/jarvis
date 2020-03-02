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
		stdout                 io.Writer
		debug                  bool
	}

	Playbook struct {
		PlaybookExecutor
		CommonArgs
	}
)

func InitPlaybookExecutor(cmdExecutor CommandExecutor, playbookBinPath, playbookPath string,
	args CommonArgs, stdout io.Writer, debug bool) *Playbook {
	return &Playbook{
		PlaybookExecutor: PlaybookExecutor{cmdExecutor, playbookBinPath, playbookPath, stdout, debug},
		CommonArgs:       args,
	}
}

func (play Playbook) Play(playbookName string) error {
	if len(play.Inventories) == 0 {
		return fmt.Errorf("playbook: no inventory to work on")
	}

	if play.debug {
		play.printInventories(play.stdout)
	}

	playbookPath, err := play.computePlaybookPath(playbookName)
	if err != nil {
		return err
	}

	for _, inventory := range play.Inventories {
		fmt.Fprintf(play.stdout, "Start running %s playbook on %s inventory...\n", playbookName, inventory)
		err = play.play(playbookPath, inventory)
		if err != nil {
			return err
		}
	}

	return nil
}

func (play Playbook) play(playbookPath string, inventory string) error {
	if len(inventory) == 0 {
		return fmt.Errorf("playbook: inventory is empty")
	}

	params := append(play.computeAnsibleOptions(inventory), playbookPath)
	err := play.commandExecutor.Run(play.playbookBinPath, params...)
	if err != nil {
		return err
	}

	return nil
}

func (play Playbook) computePlaybookPath(name string) (string, error) {
	playbookPath := filepath.Join(play.playbookRepositoryPath, name)
	info, err := os.Stat(playbookPath)

	if os.IsNotExist(err) {
		return "", fmt.Errorf("playbook: playbook does not exists:%v", playbookPath)
	}

	if info.IsDir() {
		return "", fmt.Errorf("playbook: are you kidding? it's a directory:%v", playbookPath)
	}

	return playbookPath, nil
}

func (play Playbook) computeAnsibleOptions(inventory string) []string {
	var result = play.computeCommonArgsWithInventory(inventory)
	result = append(result, play.OtherArgs...)

	return result
}
