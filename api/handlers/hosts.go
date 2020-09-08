package handlers

import (
	"encoding/json"
	"io"
	"jarvis/internal/pkg/ansible"
	"jarvis/internal/pkg/environment"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type hostsHandler struct {
	*log.Logger
}

func (hH hostsHandler) list(w http.ResponseWriter, r *http.Request) {
	log.Trace("hosts path called")
	vars := mux.Vars(r)
	helperPath := path.Join(
		viper.GetString("environments.path"), viper.GetString("environments.helper"))

	//lets return selected at true by default
	selectedRaw, selectedInQuery := r.URL.Query()["selected"]
	selected := true
	if selectedInQuery && len(selectedRaw) > 0 {
		//don't care of multiple `selected` query parameter
		if s, err := strconv.ParseBool(selectedRaw[0]); err == nil {
			selected = s
		}
	}

	//env, group
	if rawEnv, ok := vars["env"]; ok {
		env, err := environment.ParseRawEnvironmentPredicate(helperPath, rawEnv)
		if err != nil {
			logErrorToResponse(err, "hosts", w)
			return
		}
		envs := []*environment.Environment{env}
		invReaders, err := inventoriesReaders(envs)
		if err != nil {
			logErrorToResponse(err, "hosts", w)
			return
		}
		r := io.MultiReader(invReaders...)
		manipulator, err := ansible.InitInventoryManipulator(r)
		if err != nil {
			logErrorToResponse(err, "hosts", w)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		if rawGroups, ok := vars["groups"]; ok {
			groups := strings.Split(rawGroups, ",")
			var hosts []string
			for _, group := range groups {
				h, err := manipulator.GetHostsByGroupName(group)
				hosts = append(hosts, h...)
				if err != nil {
					logErrorToResponse(err, "hosts", w)
					return
				}
			}
			js, err := json.Marshal(fillRundeckResult(hosts, selected))
			if err != nil {
				logErrorToResponse(err, "hosts", w)
				return
			}

			w.Write(js)
			return
		}
	}
}
