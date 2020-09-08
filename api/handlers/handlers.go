package handlers

import (
	"io"
	"jarvis/internal/pkg/ansible"
	"jarvis/internal/pkg/environment"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func API(log *log.Logger) http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	//check (health)
	cH := checkHandler{log}
	router.HandleFunc("/health", cH.check).
		Methods("GET")

	//playbooks
	pH := playbookHandler{log}
	router.HandleFunc("/playbooks", pH.list).
		Methods("GET")

	//platforms
	ptH := platformsHandler{log}
	router.HandleFunc("/platforms/{env}", ptH.list).
		Methods("GET")

	//groups
	gH := groupsHandler{log}
	router.HandleFunc("/groups/{env}", gH.list).
		Methods("GET")

	//hosts
	hH := hostsHandler{log}
	router.HandleFunc("/hosts/{env}/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}).
		Methods("GET")
	router.HandleFunc("/hosts/{env}/{groups}", hH.list).
		Methods("GET")

	return router
}

func logErrorToResponse(err error, handler string, w http.ResponseWriter) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	log.WithFields(log.Fields{
		"handler": handler,
	}).Error(err.Error())
}

func inventoriesReaders(envs []*environment.Environment) ([]io.Reader, error) {
	envsPath := viper.GetString("environments.path")

	allInventories, err := environment.GetFullPathInventoriesFromEnvironments(envsPath, envs)
	if err != nil {
		return nil, err
	}

	invReaders, err := ansible.BuildReadersFromInventoriesPath(allInventories)
	if err != nil {
		return nil, err
	}

	return invReaders, nil
}
