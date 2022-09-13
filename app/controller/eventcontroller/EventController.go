package eventcontroller

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rizface/golang-api-template/app/entity/requestentity"
	"github.com/rizface/golang-api-template/app/middleware"
	"github.com/rizface/golang-api-template/app/service/eventservice"
)

type EventControllerInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
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
		r.Post("/", controller.Create)
	})
}

func (event *EventController) Create(w http.ResponseWriter, r *http.Request) {
	payload := new(requestentity.Event)
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		panic(err)
	}

	err = payload.Validate()
	if err != nil {
		panic(err)
	}

	props := map[string]interface{}{
		"user":    r.Context().Value("user"),
		"payload": payload,
	}
	result := event.eventservice.Create(props)
	json.NewEncoder(w).Encode(result)
}
