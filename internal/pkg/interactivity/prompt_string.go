package interactivity

import (
	"io"

	"github.com/manifoldco/promptui"
)

type PromptString struct {
	key    string
	prompt promptui.Prompt
	next   Scenario
}

func InitPromptString(key, label string, stdin io.ReadCloser, stdout io.WriteCloser) *PromptString {
	prompt := promptui.Prompt{
		Label:  label,
		Stdin:  stdin,
		Stdout: stdout,
	}
	return &PromptString{key, prompt, nil}
}

func (p *PromptString) Execute(m map[string]string) error {
	result, err := p.prompt.Run()

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

func (p *PromptString) SetNext(s Scenario) {
	p.next = s
}
