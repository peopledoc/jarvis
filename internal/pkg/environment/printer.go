package environment

import (
	"io"
)

type EnvironmentsPrinter interface {
	PrintEnvironments(output io.Writer, envs []*Environment)
}
