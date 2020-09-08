package environment

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	predicateValidator                *regexp.Regexp
	predicateValidatorJoinInventories *regexp.Regexp
)

func init() {
	//compile only once the regex
	predicateValidator = regexp.MustCompile(`(^[a-z-A-Z]+)(\.[a-z-A-Z]+)?(\.[a-z-A-Z]+)?$`)
	predicateValidatorJoinInventories = regexp.MustCompile(`(^[a-z-A-Z]+)(\.[a-z-A-Z]+)(\.[a-z-A-Z]+)?$`)
}

//FindEnvironmentTreeFromPredicate return an environment tree corresponding to
// the predicate parameter
func FindEnvironmentTreeFromPredicate(predicate *ParsedPredicate, envs []*Environment) (*Environment, error) {
	if predicate == nil || envs == nil {
		return nil, errors.New("environment: predicate or environment are null")
	}
	if len(predicate.Type) == 0 {
		return nil, errors.New("environment: environment type predicate is empty")
	}

	var env *Environment
	if env = findEnvironment(predicate.Type, envs); env == nil {
		return nil, fmt.Errorf("environment: can't find desired environment type: %s", predicate.Type)
	}

	if len(predicate.Environment) == 0 {
		return env, nil
	}

	var envDescription *EnvironmentDescription
	if envDescription = findEnvDescription(predicate.Environment, env.Descriptions); envDescription == nil {
		return nil, fmt.Errorf("environment: can't find desired environment: %s", predicate.Environment)
	}

	if len(predicate.Platform) == 0 {
		env.Descriptions = []EnvironmentDescription{
			*envDescription,
		}
		return env, nil
	}

	var platform *Platform
	if platform = findPlatform(predicate.Platform, envDescription.Platforms); platform == nil {
		return nil, fmt.Errorf("environment: can't find desired platform: %s", predicate.Platform)
	}

	envDescription.Platforms = []Platform{
		*platform,
	}
	env.Descriptions = []EnvironmentDescription{
		*envDescription,
	}

	return env, nil
}

func GetFullPathInventoriesFromEnvironments(basePath string, envs []*Environment) ([][]string, error) {
	if envs == nil {
		return nil, fmt.Errorf("environment: environment tree is nil")
	}
	var globalInventories [][]string

	for _, env := range envs {
		if env.Descriptions == nil {
			return nil, fmt.Errorf("environment: no environment descriptions")
		}
		for _, description := range env.Descriptions {
			if description.Platforms == nil {
				return nil, fmt.Errorf("environment: no platforms")
			}

			for _, platform := range description.Platforms {
				if platform.Inventories == nil && len(platform.Inventories) == 0 {
					return nil, fmt.Errorf("environment: no inventories")
				}

				invs := make([]string, len(platform.Inventories))
				for i, relativePath := range platform.Inventories {
					invs[i] = filepath.Join(basePath, relativePath)
				}

				globalInventories = append(globalInventories, invs)
			}
		}

	}
	return globalInventories, nil
}

//ParsePredicate validate that the predicate follow `type.env.platform' syntax
//then return a parsedpredicate corresponding to the predicate string parameter
func ParsePredicate(predicate string) (*ParsedPredicate, error) {
	if !isPredicateSyntaxValid(predicate) {
		return nil, fmt.Errorf("environment: predicate is not valid: %s", predicate)
	}
	matches := predicateValidator.FindStringSubmatch(predicate)
	//Index 0 holds the text of the leftmost match and following by the matches

	//Index 1 is mandatory and contains an environment type
	result := &ParsedPredicate{}
	result.Type = strings.TrimSpace(matches[1])

	//Because the last two groups capture a dot following by a word
	//we need to remove the first character

	//Index 2 are optional, it contains a environment name
	if len(matches[2]) > 0 {
		result.Environment = strings.TrimSpace(matches[2][1:])
	}

	//Index 3 are optional too, it contains a platform
	if len(matches[3]) > 0 {
		result.Platform = strings.TrimSpace(matches[3][1:])
	}

	return result, nil
}

func IsPredicateJoinInventoriesSyntaxValid(predicate string) bool {
	return predicateValidatorJoinInventories.MatchString(predicate)
}

func findPlatform(name string, platforms []Platform) *Platform {
	for _, platform := range platforms {
		if platform.Name == name {
			return &platform
		}
	}
	return nil
}

func findEnvDescription(name string, envs []EnvironmentDescription) *EnvironmentDescription {
	for _, env := range envs {
		if env.Name == name {
			return &env
		}
	}
	return nil
}

func findEnvironment(envName string, envs []*Environment) *Environment {
	for _, env := range envs {
		if env.Type == envName {
			return env
		}
	}
	return nil
}

func isPredicateSyntaxValid(predicate string) bool {
	return predicateValidator.MatchString(predicate)
}
