package ansible

import (
	"fmt"
	"io"
	"io/ioutil"
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
	if isInventoriesEmpty(play.Inventories) {
		return fmt.Errorf("playbook: no inventory to work on")
	}

	if play.debug {
		play.printInventories(play.stdout)
	}

	playbookPath, err := play.computePlaybookPath(playbookName)
	if err != nil {
		return err
	}

	params := play.buildCommonArgs()

	if play.JoinInventories {
		invs := play.joinInventories()
		params = append(params, invs...)
		err = play.commandExecutor.Run(play.playbookBinPath,
			append(params, playbookPath)...)
	} else {
		for _, inventory := range play.Inventories {
			fmt.Fprintf(play.stdout,
				"Start running %s playbook on %s inventory...\n", playbookName, inventory)
			params := append(params, buildInventoryArg(inventory)...)
			err = play.commandExecutor.Run(play.playbookBinPath,
				append(params, playbookPath)...)
			if err != nil {
				break
			}
		}
	}

	return err
}

func ListPlaybooks(path string) ([]string, error) {
	var books []string
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if !f.IsDir() {
			books = append(books, f.Name())
		}
	}

	return books, nil
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
