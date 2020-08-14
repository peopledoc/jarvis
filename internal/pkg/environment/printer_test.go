package environment

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestTreePrinterOutput(t *testing.T) {
	normalOutput, err := ioutil.ReadFile("test-data/printer_output.txt")
	if err != nil {
		t.Errorf("Can't read sample file")
	}

	var buf bytes.Buffer
	envs := []*Environment{
		&Environment{
			Type: "toto",
			Descriptions: []EnvironmentDescription{
				EnvironmentDescription{
					Name: "chazam",
					Platforms: []Platform{
						Platform{"titi", "prov-1", []string{"titi/", "toto/"}},
					},
				},
				EnvironmentDescription{
					Name: "alan",
					Platforms: []Platform{
						Platform{"marvel", "prov-2", []string{"titi/"}},
					},
				},
			},
		},
	}

	treePrinter := TreePrinter{}
	treePrinter.PrintEnvironments(&buf, envs)
	output := buf.String()

	if output != string(normalOutput) {
		t.Errorf("Output incorrect, got: %s, want: %s", output, normalOutput)
	}
}
