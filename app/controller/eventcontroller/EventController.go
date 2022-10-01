package eventcontroller

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rizface/golang-api-template/app/entity/requestentity"
	"github.com/rizface/golang-api-template/app/entity/responseentity"
	"github.com/rizface/golang-api-template/app/errorgroup"
	"github.com/rizface/golang-api-template/app/middleware"
	"github.com/rizface/golang-api-template/app/schema"
	"github.com/rizface/golang-api-template/app/service/eventservice"
)

type EventControllerInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	GetOne(w http.ResponseWriter, r *http.Request)
	UpdateOneById(w http.ResponseWriter, r *http.Request)
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
		r.Get("/", controller.Get)
		r.Delete("/{id}", controller.Delete)
		r.Get("/{id}", controller.GetOne)
		r.Put("/{id}", controller.UpdateOneById)
	})
}

func (event *EventController) Create(w http.ResponseWriter, r *http.Request) {
	payload := new(requestentity.Event)
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		panic(errorgroup.BAD_REQUEST)
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
	json.NewEncoder(w).Encode(responseentity.Response{
		Code:    http.StatusOK,
		Message: "SUCCESS",
		Data:    result,
	})
}

func (event *EventController) Get(w http.ResponseWriter, r *http.Request) {
	pagination := requestentity.NewPagination(r.URL.Query())
	result := event.eventservice.Get(pagination)
	json.NewEncoder(w).Encode(responseentity.Response{
		Code:    http.StatusOK,
		Message: "SUCCESS",
		Data:    result,
		Page:    pagination,
	})
}

func (event *EventController) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := schema.ValidateEventId(id)
	if err != nil {
		panic(err)
	}
	user := r.Context().Value("user")
	result := event.eventservice.Delete(id, user)
	json.NewEncoder(w).Encode(responseentity.Response{
		Code:    http.StatusOK,
		Message: "SUCCESS",
		Data:    result,
	})
}

func (event *EventController) GetOne(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := schema.ValidateEventId(id)
	if err != nil {
		panic(err)
	}

	result := event.eventservice.GetOne(id)
	json.NewEncoder(w).Encode(responseentity.Response{
		Code:    http.StatusOK,
		Message: "SUCCESS",
		Data:    result,
	})
}

func (event *EventController) UpdateOneById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := schema.ValidateEventId(id)
	if err != nil {
		panic(err)
	}

	payload := new(requestentity.Event)
	err = json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		panic(errorgroup.BAD_REQUEST)
	}
	err = payload.Validate()
	if err != nil {
		panic(err)
	}

	user := r.Context().Value("user")
	data := map[string]interface{}{
		"user":    user,
		"payload": payload,
		"id":      id,
	}
	result := event.eventservice.UpdateOneById(data)
	json.NewEncoder(w).Encode(responseentity.Response{
		Code:    http.StatusOK,
		Message: "SUCCESS",
		Data:    result,
	})
}
