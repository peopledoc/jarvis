//go:generate mockgen -package interactivity -self_package jarvis/internal/pkg/interactivity -destination=scenario_mock.go jarvis/internal/pkg/interactivity Scenario
package interactivity

type Scenario interface {
	//use map as a generic json placeholder
	Execute(map[string]string) error
	SetNext(Scenario)
}
