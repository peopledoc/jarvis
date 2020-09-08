package handlers

import (
	"jarvis/internal/pkg/environment"
	"net/http"
	"path"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type platformsHandler struct {
	*log.Logger
}

func (pH platformsHandler) list(w http.ResponseWriter, r *http.Request) {
	log.Trace("platforms path called")
	helperPath := path.Join(
		viper.GetString("environments.path"), viper.GetString("environments.helper"))
	vars := mux.Vars(r)

	env, err := environment.ParseRawEnvironmentPredicate(helperPath, vars["env"])
	if err != nil {
		logErrorToResponse(err, "platforms", w)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	printer := environment.JsonPrinter{}
	printer.PrintEnvironments(w, []*environment.Environment{env})
}
