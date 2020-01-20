package environment

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func ParseEnvironmentFile(envFilePath string) (*Environments, error) {
	content, err := ioutil.ReadFile(envFilePath)
	if err != nil {
		return nil, err
	}
	var model Environments
	err = yaml.Unmarshal(content, &model)
	if err != nil {
		return nil, err
	}
	return &model, nil
}
