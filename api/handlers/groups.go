package handlers

import (
	"encoding/json"
	"io"
	"jarvis/internal/pkg/ansible"
	"jarvis/internal/pkg/environment"
	"net/http"
	"path"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type groupsHandler struct {
	*log.Logger
}

func (gH groupsHandler) list(w http.ResponseWriter, r *http.Request) {
	log.Trace("hosts path called")
	vars := mux.Vars(r)
	helperPath := path.Join(
		viper.GetString("environments.path"), viper.GetString("environments.helper"))

	//env, group
	if rawEnv, ok := vars["env"]; ok {
		env, err := environment.ParseRawEnvironmentPredicate(helperPath, rawEnv)
		if err != nil {
			logErrorToResponse(err, "groups", w)
			return
		}
		envs := []*environment.Environment{env}
		invReaders, err := inventoriesReaders(envs)
		if err != nil {
			logErrorToResponse(err, "groups", w)
			return
		}
		r := io.MultiReader(invReaders...)
		manipulator, err := ansible.InitInventoryManipulator(r)
		if err != nil {
			logErrorToResponse(err, "groups", w)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		groups, err := manipulator.GetGroupsName(false)
		if err != nil {
			logErrorToResponse(err, "groups", w)
			return
		}

		js, err := json.Marshal(groups)
		if err != nil {
			logErrorToResponse(err, "groups", w)
			return
		}

		w.Write(js)
		return
	}
}
