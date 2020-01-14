package environment

import (
	"fmt"
	"io"
)

type TreePrinter struct{}

func (p TreePrinter) PrintEnvironments(output io.Writer, envs Environments) {
	platformPrefix := "│   "
	envPrefix := "├──"
	for i, value := range envs {
		if i == len(envs)-1 {
			platformPrefix = "    "
			envPrefix = "└──"
		}
		fmt.Fprintln(output, envPrefix, value.Type)
		if value != nil {
			if len(value.Descriptions) > 0 {
				p.printEnvDescriptions(output, value.Descriptions, platformPrefix)
			}
		}
	}
}

func (p TreePrinter) printEnvDescriptions(output io.Writer, descriptions []EnvironmentDescription, prefix string) {
	for i := 0; i < len(descriptions)-1; i++ {
		fmt.Fprintln(output, prefix, "├──", descriptions[i].Name)
		p.printPlatforms(output, descriptions[i].Platforms, prefix+" │   ")
	}
	fmt.Fprintln(output, prefix, "└──", descriptions[len(descriptions)-1].Name)
	p.printPlatforms(output, descriptions[len(descriptions)-1].Platforms, prefix+"     ")
}

func (p TreePrinter) printPlatforms(output io.Writer, platform []Platform, prefix string) {
	for i := 0; i < len(platform)-1; i++ {
		fmt.Fprintln(output, prefix, "├──", platform[i].Name)
		printInventories(output, platform[i].Inventories, prefix+" │   ")
	}
	fmt.Fprintln(output, prefix, "└──", platform[len(platform)-1].Name)
	printInventories(output, platform[len(platform)-1].Inventories, prefix+"     ")
}

func printInventories(output io.Writer, inventories []string, prefix string) {
	for i := 0; i < len(inventories)-1; i++ {
		fmt.Fprintln(output, prefix, "├──", inventories[i])
	}
	fmt.Fprintln(output, prefix, "└──", inventories[len(inventories)-1])
}
