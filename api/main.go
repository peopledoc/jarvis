package api

import (
	"fmt"
	"jarvis/api/handlers"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type Api struct {
	logger *log.Logger
	port   int
}

func InitApi(logger *log.Logger, port int) *Api {
	return &Api{logger, port}
}

func (api *Api) Run() error {
	var httpSrv *http.Server

	h := handlers.API(api.logger)
	httpSrv = &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      h,
		Addr:         fmt.Sprintf(":%d", api.port),
	}

	err := httpSrv.ListenAndServe()

	return err
}
