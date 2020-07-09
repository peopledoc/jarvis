package ansible

import (
	"fmt"
	"io"
	"sort"

	"github.com/relex/aini"
)

type (
	InventoryManipulator struct {
		io.Reader
	}
)

func InitInventoryManipulator(reader io.Reader) *InventoryManipulator {
	return &InventoryManipulator{reader}
}

//GetGroupsName returns groups name sorted
func (iM InventoryManipulator) GetGroupsName(withParents bool) ([]string, error) {
	inventoryD, err := aini.Parse(iM)
	if err != nil {
		return nil, err
	}

	var groups []string
	for _, g := range inventoryD.Groups {
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
	inventoryD, err := aini.Parse(iM)
	if err != nil {
		return nil, err
	}

	group, ok := inventoryD.Groups[groupName]
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
