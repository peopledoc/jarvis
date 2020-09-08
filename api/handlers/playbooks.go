package handlers

import (
	"encoding/json"
	"jarvis/internal/pkg/ansible"
	"net/http"

	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

type playbookHandler struct {
	*log.Logger
}

func (pH playbookHandler) list(w http.ResponseWriter, r *http.Request) {
	pH.Trace("playbooks path called")
	pBookPath := viper.GetString("ansible.playbook.playbooks_path")
	pbooks, err := ansible.ListPlaybooks(pBookPath)
	if err != nil {
		logErrorToResponse(err, "playbook", w)
		return
	}
	js, err := json.Marshal(pbooks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
