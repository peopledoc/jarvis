package environment

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func ParseRawEnvironmentPredicate(environmentHelperPath, rawPredicate string) (*Environment, error) {
	envs, err := ParseEnvironmentFile(environmentHelperPath)
	if err != nil {
		return nil, fmt.Errorf("fatal error while listing inventories: %s", err)
	}

	predicate, err := ParsePredicate(rawPredicate)
	if err != nil {
		return nil, err
	}
	choosenEnvTree, err := FindEnvironmentTreeFromPredicate(predicate, envs)
	if err != nil {
		return nil, err
	}

	return choosenEnvTree, nil
}

func ParseEnvironmentFile(envFilePath string) ([]*Environment, error) {
	content, err := ioutil.ReadFile(envFilePath)
	if err != nil {
		return nil, err
	}
	var model []*Environment
	err = yaml.Unmarshal(content, &model)
	if err != nil {
		return nil, err
	}
	return model, nil
}
