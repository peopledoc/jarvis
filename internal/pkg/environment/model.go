package environment

type (
	Environment struct {
		Type         string                   `yaml:"type"`
		Descriptions []EnvironmentDescription `yaml:"envs"`
	}

	EnvironmentDescription struct {
		Name      string     `yaml:"name"`
		Platforms []Platform `yaml:"platforms"`
	}

	Platform struct {
		Name        string   `yaml:"name"`
		Provider    string   `yaml:"provider"`
		Inventories []string `yaml:"inventories"`
	}

	ParsedPredicate struct {
		Type, Environment, Platform string
	}
)
