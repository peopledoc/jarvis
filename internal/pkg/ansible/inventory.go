package ansible

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/relex/aini"
)

type (
	InventoryManipulator struct {
		inventory *aini.InventoryData
	}
)

func InitInventoryManipulator(reader io.Reader) (*InventoryManipulator, error) {
	inventoryD, err := aini.Parse(reader)
	if err != nil {
		return nil, err
	}
	return &InventoryManipulator{inventoryD}, nil
}

//GetGroupsName returns groups name sorted
func (iM InventoryManipulator) GetGroupsName(withParents bool) ([]string, error) {
	var groups []string
	for _, g := range iM.inventory.Groups {
		if !withParents && len(g.Children) > 0 {
			continue
		}
		groups = append(groups, g.Name)
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i] < groups[j]
	})
	return groups, nil
}

//GetHostsByGroupName returns hosts owned by a group, sorted
func (iM InventoryManipulator) GetHostsByGroupName(groupName string) ([]string, error) {
	group, ok := iM.inventory.Groups[groupName]
	if !ok {
		return nil, fmt.Errorf("Can't find %v group", groupName)
	}
	rawHosts := aini.HostMapListValues(group.Hosts)
	var hosts []string
	for _, h := range rawHosts {
		hosts = append(hosts, h.Name)
	}

	sort.Slice(hosts, func(i, j int) bool {
		return hosts[i] < hosts[j]
	})
	return hosts, nil
}

func BuildReadersFromInventoriesPath(invPaths [][]string) ([]io.Reader, error) {
	var invReaders []io.Reader
	for _, invs := range invPaths {
		for _, inv := range invs {
			if !fileExists(inv) {
				return nil, fmt.Errorf("the %v file does not exist", inv)
			}
			f, err := os.Open(inv)
			if err != nil {
				return nil, err
			}
			invReaders = append(invReaders, bufio.NewReader(f))
			//TODO: use logger for debug purpose
		}
	}
	return invReaders, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
