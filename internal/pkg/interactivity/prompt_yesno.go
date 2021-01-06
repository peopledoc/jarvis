package interactivity

import (
	"io"

	"github.com/manifoldco/promptui"
)

type PromptYesNo struct {
	key    string
	prompt promptui.Prompt
	next   Scenario
}

func InitPromptYesNo(key, label string, stdin io.ReadCloser, stdout io.WriteCloser) *PromptString {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
		Stdin:     stdin,
		Stdout:    stdout,
	}
	return &PromptString{key, prompt, nil}
}

func (p *PromptYesNo) Execute(m map[string]string) error {
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

func (p *PromptYesNo) SetNext(s Scenario) {
	p.next = s
}
