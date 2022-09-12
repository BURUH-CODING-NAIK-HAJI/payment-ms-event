package eventservice

import (
	"github.com/rizface/golang-api-template/app/repository/eventrepository"
)

type EventServiceInterface interface {
	Welcome() string
}

type EventService struct {
	eventrepository eventrepository.EventRepositoryInterface
}

func New(eventrepository eventrepository.EventRepositoryInterface) EventServiceInterface {
	return &EventService{
		eventrepository: eventrepository,
	}
}

func (eventservice *EventService) Welcome() string {
	response := eventservice.eventrepository.Welcome()
	return response
}
