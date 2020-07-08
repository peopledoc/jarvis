package ansible

import (
	"bufio"
	"errors"
	"os"
	"testing"
)

func TestGetGroupName(t *testing.T) {
	tests := []struct {
		name        string
		withParents bool
		goldenFile  string
		result      []string
		err         error
	}{
		{"withoutParents", false, "inventory1", []string{"g1", "g2", "ungrouped"}, nil},
		{"withParents", true, "inventory1", []string{"all", "g1", "g2", "p1", "ungrouped"}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := os.Open("test-data/" + tt.goldenFile)
			if err != nil {
				t.Fatalf("Error loading golden file: %s", err)
			}
			reader := bufio.NewReader(content)
			iM := InitInventoryManipulator(reader)
			res, err := iM.GetGroupsName(tt.withParents)

			if tt.err != nil {
				if tt.err.Error() != err.Error() {
					t.Errorf("Must have an error here, have:%v, want:%v", err.Error(), tt.err.Error())
				}
				return
			}

			if !Equal(res, tt.result) {
				t.Errorf("Groups name are different, have %v; want %v", res, tt.result)
			}
		})
	}
}

func TestGetHostsByGroupName(t *testing.T) {
	tests := []struct {
		name       string
		goldenFile string
		group      string
		result     []string
		err        error
	}{
		{"children group", "inventory1", "g1", []string{"h1", "h2"}, nil},
		{"parent group", "inventory1", "p1", []string{"h1", "h2", "h3"}, nil},
		{"cant find group", "inventory1", "x1", nil, errors.New("Can't find x1 group")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := os.Open("test-data/" + tt.goldenFile)
			if err != nil {
				t.Fatalf("Error loading golden file: %s", err)
			}
			reader := bufio.NewReader(content)
			iM := InitInventoryManipulator(reader)
			res, err := iM.GetHostsByGroupName(tt.group)

			if tt.err != nil {
				if tt.err.Error() != err.Error() {
					t.Errorf("Must have an error here, have:%v, want:%v", err.Error(), tt.err.Error())
				}
				return
			}

			if !Equal(res, tt.result) {
				t.Errorf("Hosts name are different, have %v; want %v", res, tt.result)
			}
		})
	}
}

func Equal(a, b []string) bool {
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
