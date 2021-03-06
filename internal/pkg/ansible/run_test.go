package ansible

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/golang/mock/gomock"

	"jarvis/internal/pkg/mocks"
)

func TestRunModule(t *testing.T) {
	const runBinPath = "binPath"
	const runModuleName = "modulename"
	const runModuleArg = "someargs"

	tests := []struct {
		name        string
		module      string
		modulearg   string
		hosttarget  string
		inventories [][]string
		args        []string
		err         error
	}{
		{"no inventories", "modulename", "fooarg", "targets", [][]string{}, []string{"someargs"}, errors.New("run: no inventory to work on")},
		{"empty inventory", "modulename", "fooarg", "targets", [][]string{{""}}, []string{"someargs"}, errors.New("run: inventory is empty")},
		{"valid inventories", "modulename", "fooarg", "targets",
			[][]string{{"inv1"}, {"inv2"}},
			[]string{"--diff", "--inventory"}, nil},
	}
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mocks.NewMockCommandExecutor(ctrl)
			if tt.err != nil {
				m.EXPECT().
					Run(runBinPath, []string{}).
					Times(0)
			} else {
				for _, inventories := range tt.inventories {
					for _, inventory := range inventories {
						completeArgs := append(
							tt.args,
							inventory,
							"-m", tt.module,
							"--args", tt.modulearg,
							tt.hosttarget,
						)
						m.EXPECT().
							Run(runBinPath, completeArgs).
							Times(1)
					}
				}
			}
			rE := InitRunModuleExecutor(
				m, runBinPath,
				RunModuleArgs{
					CommonArgs:  CommonArgs{Inventories: tt.inventories},
					ModuleName:  tt.module,
					ModuleArg:   tt.modulearg,
					HostPattern: tt.hosttarget,
				},
				ioutil.Discard,
				false)
			ret := rE.Run()

			if tt.err != nil {
				if tt.err.Error() != ret.Error() {
					t.Errorf("Expected error doesn't match, want:'%v' vs have:'%v'", tt.err, ret)
				}
			}

		})
	}
}

func TestComputeRunAnsibleOptions(t *testing.T) {
	tests := []struct {
		name         string
		runArgs      RunModule
		inventory    []string
		expectedArgs []string
	}{
		{
			"default module arg",
			RunModule{
				RunModuleArgs: RunModuleArgs{
					ModuleName:  "fooMod",
					ModuleArg:   "fooArg",
					HostPattern: "fooHost",
				}},
			[]string{"inv1"},
			[]string{"--diff", "--inventory", "inv1", "-m", "fooMod", "--args", "fooArg", "fooHost"},
		},
		{
			"nodiff",
			RunModule{
				RunModuleArgs: RunModuleArgs{

					CommonArgs: CommonArgs{
						HideDiff:  true,
						OtherArgs: []string{"other1", "other2"},
					},
					ModuleName:  "barMod",
					ModuleArg:   "barArg",
					HostPattern: "barHost",
				}},
			[]string{"inv2"},
			[]string{"other1", "other2", "--inventory", "inv2", "-m", "barMod", "--args", "barArg", "barHost"},
		},
		{
			"noarg",
			RunModule{
				RunModuleArgs: RunModuleArgs{

					CommonArgs: CommonArgs{
						BecomeSudo: true,
						OtherArgs:  []string{"other1", "other2"},
					},
					ModuleName:  "ping",
					HostPattern: "barHost",
				}},
			[]string{"inv3"},
			[]string{"-b", "--diff", "other1", "other2", "--inventory", "inv3", "-m", "ping", "barHost"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.runArgs.computeAnsibleOptions(tt.inventory)
			if !equal(args, tt.expectedArgs) {
				t.Errorf("Wrong args want:%v, got:%v", tt.expectedArgs, args)
			}

		})
	}
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
