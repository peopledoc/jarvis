package ansible

import (
	"errors"
	"io/ioutil"
	"jarvis/internal/pkg/mocks"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestPlaybookExecutorPlay(t *testing.T) {
	const playbookBinPath = "binPath"
	const playbookRepoPath = "test-data"
	const playbookName = "playbook.yml"
	fullPlaybookPath := filepath.Join(playbookRepoPath, playbookName)
	tests := []struct {
		name        string
		inventories []string
		args        []string
		err         error
	}{
		{"empty inventories", []string{}, []string{}, errors.New("playbook: no inventory to work on")},
		{"empty inventory", []string{""}, []string{}, errors.New("playbook: inventory is empty")},
		{
			"successfull run",
			[]string{"inventory1", "inventory2", "inventory3"},
			[]string{"--diff", "--check", "--inventory"},
			nil,
		},
	}
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mocks.NewMockCommandExecutor(ctrl)

			if tt.err != nil {
				//We don't expect any call to run if we have an error
				m.EXPECT().
					Run(playbookBinPath, tt.args).
					Times(0)
			} else {
				//We iterate over inventories to construct except
				//calls
				for _, inventory := range tt.inventories {
					completeArgs := append(tt.args, inventory, fullPlaybookPath)
					m.EXPECT().
						Run(playbookBinPath, completeArgs).
						Return(nil).
						Times(1)
				}
			}
			//ioutil.Discard because we don't care of the debug output
			//it's like > /dev/null
			pE := InitPlaybookExecutor(
				m, playbookBinPath, playbookRepoPath,
				PlaybookArgs{Inventories: tt.inventories}, ioutil.Discard, false)
			err := pE.Play(playbookName)

			if tt.err != nil {
				if tt.err.Error() != err.Error() {
					t.Errorf("Must have an error here, want:%v, have:%v", tt.err.Error(), err.Error())
				}
			}
		})
	}
}

func TestComputePlaybookPath(t *testing.T) {
	tests := []struct {
		name             string
		playbookRepoPath string
		playbookName     string
		result           string
		err              error
	}{
		{"file not exists", "test-data", "donaldknuth", "",
			errors.New("playbook: playbook does not exists:test-data/donaldknuth")},
		{"it's a directory", "test-data", "", "",
			errors.New("playbook: are you kidding? it's a directory:test-data")},
		{"valid playbook", "test-data", "playbook.yml", "test-data/playbook.yml",
			nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pE := PlaybookExecutor{playbookRepositoryPath: tt.playbookRepoPath}
			path, err := pE.computePlaybookPath(tt.playbookName)

			if tt.err != nil {
				if tt.err.Error() != err.Error() {
					t.Errorf("Must have an error here, have:%v, want:%v", err.Error(), tt.err.Error())
				}
				return
			}

			if path != tt.result {
				t.Errorf("Different playbook path, have:%v, want:%v", path, tt.result)
			}
		})
	}
}

func TestCombineWithInventory(t *testing.T) {
	tests := []struct {
		name               string
		playbookArgs       PlaybookArgs
		checkerAgainstArgs func(args []string) bool
		inventory          string
	}{
		{"show diff", PlaybookArgs{HideDiff: false},
			mustInArgs("--diff"), ""},
		{"hide diff", PlaybookArgs{HideDiff: true},
			mustNotInArgs("--diff"), ""},
		{"check mode activated", PlaybookArgs{CheckModeDeactivated: false},
			mustInArgs("--check"), ""},
		{"check mode deactivated", PlaybookArgs{CheckModeDeactivated: true},
			mustNotInArgs("--check"), ""},
		{"become", PlaybookArgs{BecomeSudo: true},
			mustInArgs("-b"), ""},
		{"not become", PlaybookArgs{BecomeSudo: false},
			mustNotInArgs("-b"), ""},
		{"inventory toto", PlaybookArgs{},
			mustInArgs("--inventory", "toto"), "toto"},
		{"other args", PlaybookArgs{OtherArgs: []string{"-t", "ok", "-l", "chazam"}},
			mustInArgs("-t", "ok", "-l", "chazam"), ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.playbookArgs.combineWithInventory(tt.inventory)

			if !tt.checkerAgainstArgs(args) {
				t.Error("error while checking args")
			}
		})
	}
}

func mustInArgs(strs ...string) func(args []string) bool {
	return func(args []string) bool {
		for _, str := range strs {
			localIn := false
			for _, arg := range args {
				if str == arg {
					localIn = true
				}
			}
			if !localIn {
				return false
			}
		}
		return true
	}
}

func mustNotInArgs(strs ...string) func(args []string) bool {
	return func(args []string) bool {
		return !mustInArgs(strs...)(args)
	}
}
