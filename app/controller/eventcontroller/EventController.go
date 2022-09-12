package eventcontroller

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rizface/golang-api-template/app/middleware"
	"github.com/rizface/golang-api-template/app/service/eventservice"
)

type EventControllerInterface interface {
	Welcome(w http.ResponseWriter, r *http.Request)
}

type EventController struct {
	eventservice eventservice.EventServiceInterface
}

func New(eventservice eventservice.EventServiceInterface) EventControllerInterface {
	return &EventController{
		eventservice: eventservice,
	}
}

func Setup(router *chi.Mux, controller EventControllerInterface) {
	router.Route("/", func(r chi.Router) {
		r.Use(middleware.AuthHandler)
		r.Get("/", controller.Welcome)
	})
	router.Get("/test", controller.Welcome)
}

func (welcome *EventController) Welcome(w http.ResponseWriter, r *http.Request) {
	response := welcome.eventservice.Welcome()
	json.NewEncoder(w).Encode(response)
}
