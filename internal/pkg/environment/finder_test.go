package environment

import (
	"errors"
	"testing"
)

func TestFindEnvironmentTree(t *testing.T) {
	envs := Environments{
		&Environment{
			Type: "toto",
			Descriptions: []EnvironmentDescription{
				EnvironmentDescription{
					"chazam",
					[]Platform{
						Platform{"there", "pub", []string{"here/"}},
					},
				},
			},
		},
	}

	tests := []struct {
		name        string
		predicate   *ParsedPredicate
		environment *Environment
		err         error
	}{
		{"toto", &ParsedPredicate{"toto", "", ""}, envs[0], nil},
		{"toto.chazam", &ParsedPredicate{"toto", "chazam", ""}, envs[0], nil},
		{"toto.chazam.there", &ParsedPredicate{"toto", "chazam", "there"}, envs[0], nil},
		{"titi", &ParsedPredicate{"titi", "", ""}, nil,
			errors.New("environment: can't find desired environment type: titi")},
		{"toto.titi", &ParsedPredicate{"toto", "titi", ""}, nil,
			errors.New("environment: can't find desired environment: titi")},
		{"toto.chazam.not-there", &ParsedPredicate{"toto", "chazam", "not-there"}, nil,
			errors.New("environment: can't find desired platform: not-there")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FindEnvironmentTreeFromPredicate(tt.predicate, &envs)
			if tt.err != nil {
				if err == nil {
					t.Errorf("Should have returned an error, %v", tt.err)
					return
				}
				if tt.err.Error() != err.Error() {
					t.Errorf("Something is wrong while finding env tree")
					return
				}
			}
		})
	}
}

func TestParsePredicate(t *testing.T) {
	tests := []struct {
		name            string
		predicate       string
		parsedPredicate *ParsedPredicate
		err             error
	}{
		{"only env", "env", &ParsedPredicate{"env", "", ""}, nil},
		{"env and location", "env.location", &ParsedPredicate{"env", "location", ""}, nil},
		{"env,& location and sublocation", "env.location.subloc",
			&ParsedPredicate{"env", "location", "subloc"}, nil},
		{"wrong predicate", "env,21312", nil,
			errors.New("environment: predicate is not valid: env,21312")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ParsePredicate(tt.predicate)

			if tt.err != nil {
				if err.Error() != tt.err.Error() {
					t.Errorf("Something bad happened, got: %v, want: %v", err, tt.err)
				}
			} else {
				if *res != *tt.parsedPredicate {
					t.Errorf("parsed predicate not equal, 1%s 2%s 3%s", res.Environment, res.Platform,
						res.Platform)
				}
			}
		})
	}
}
