package handlers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type checkHandler struct {
	*log.Logger
}

func (cH checkHandler) check(w http.ResponseWriter, r *http.Request) {
	log.Trace("health path called")
	okResp := map[string]bool{"ok": true}

	js, err := json.Marshal(okResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return
}
