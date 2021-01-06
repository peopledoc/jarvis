package interactivity

import (
	"io"

	"github.com/manifoldco/promptui"
)

type Select struct {
	key    string
	prompt promptui.Select
	next   Scenario
}

func InitPromptSelect(key, label string, items []string, stdin io.ReadCloser, stdout io.WriteCloser) *Select {
	s := promptui.Select{
		Label:  label,
		Items:  items,
		Stdin:  stdin,
		Stdout: stdout,
	}
	return &Select{key, s, nil}
}

func (p *Select) Execute(m map[string]string) error {
	_, result, err := p.prompt.Run()

	if err != nil {
		return err
	}

	if p.key != "" {
		m[p.key] = result
	}

	if p.next != nil {
		err = p.next.Execute(m)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Select) SetNext(s Scenario) {
	p.next = s
}
