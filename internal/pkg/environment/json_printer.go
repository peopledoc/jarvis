package environment

import (
	"encoding/json"
	"fmt"
	"io"
)

type JsonPrinter struct{}

//TODO: must add an error return type
func (p JsonPrinter) PrintEnvironments(output io.Writer, envs []*Environment) {
	var rawEnvs []string
	for _, env := range envs {
		for _, desc := range env.Descriptions {
			for _, plat := range desc.Platforms {
				raw := fmt.Sprintf("%s.%s.%s", env.Type, desc.Name, plat.Name)
				rawEnvs = append(rawEnvs, raw)
			}
		}
	}

	js, _ := json.Marshal(rawEnvs)
	output.Write(js)
}
