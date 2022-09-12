package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/rizface/golang-api-template/app/controller/eventcontroller"
	"github.com/rizface/golang-api-template/app/repository/eventrepository"
	"github.com/rizface/golang-api-template/app/service/eventservice"
)

func SetupController(router *chi.Mux) {
	eventrepository := eventrepository.New()
	eventservice := eventservice.New(eventrepository)
	controller := eventcontroller.New(eventservice)
	eventcontroller.Setup(router, controller)
}

func CreateHttpServer(router http.Handler) *http.Server {
	var appPort string

	if len(appPort) == 0 {
		appPort = "9000"
	} else {
		appPort = os.Getenv("APP_PORT")
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%v", appPort),
		Handler: router,
	}

	return httpServer
}
