package interactivity

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestExecutePromptString(t *testing.T) {
	var tests = []struct {
		name      string
		key       string
		in        string
		with_next bool
	}{
		{"with key", "key_test", "in value", true},
		{"without key", "", "in value", true},
		{"without next", "key_test", "in value", false},
	}
	label := "label test"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapObject := make(map[string]string)
			r := ioutil.NopCloser(strings.NewReader(tt.in + "\n"))
			m := NewMockScenario(ctrl)

			pS := InitPromptString(tt.key, label, r, os.Stdout)

			if tt.with_next {
				m.EXPECT().
					Execute(mapObject).
					Return(nil).
					Times(1)
				pS.SetNext(m)
			}

			err := pS.Execute(mapObject)

			if tt.key == "" {
				assert.Equal(t, 0, len(mapObject))
			} else {
				assert.Equal(t, tt.in, mapObject[tt.key])
			}
			assert.NoError(t, err)
		})
	}
}
