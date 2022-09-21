package eventservice

import (
	"context"
	"time"

	"github.com/rizface/golang-api-template/app/entity/requestentity"
	"github.com/rizface/golang-api-template/app/entity/responseentity"
	"github.com/rizface/golang-api-template/app/errorgroup"
	"github.com/rizface/golang-api-template/app/repository/eventrepository"
	"github.com/rizface/golang-api-template/system/security"
	"gorm.io/gorm"
)

type EventServiceInterface interface {
	Create(payload map[string]interface{}) *responseentity.Event
	Get() *[]responseentity.Event
	Delete(id string, user interface{}) *responseentity.Event
	GetOne(id string) *responseentity.Event
	UpdateOneById(map[string]interface{}) *responseentity.Event
}

type EventService struct {
	eventrepository eventrepository.EventRepositoryInterface
	db              *gorm.DB
}

func New(
	eventrepository eventrepository.EventRepositoryInterface,
	db *gorm.DB,
) EventServiceInterface {
	return &EventService{
		eventrepository: eventrepository,
		db:              db,
	}
}

func (eventservice *EventService) Create(props map[string]interface{}) *responseentity.Event {
	payload := props["payload"].(*requestentity.Event)
	user := props["user"].(security.JwtClaim).UserData
	result, err := eventservice.eventrepository.Create(eventservice.db, user.Id, payload)

	if err != nil {
		panic(errorgroup.FAILED_CREATE_EVENT)
	}

	return result
}

func (eventservice *EventService) Get() *[]responseentity.Event {
	background := context.Background()
	ctx, cancel := context.WithTimeout(background, 5*time.Second)
	defer cancel()
	result, err := eventservice.eventrepository.Get(eventservice.db, ctx)
	if err != nil {
		panic(err)
	}
	return result
}

func (eventservice *EventService) Delete(id string, user interface{}) *responseentity.Event {
	event, err := eventservice.eventrepository.GetOne(eventservice.db, id)
	if err != nil {
		panic(err)
	}

	userData := user.(security.JwtClaim).UserData
	if event.UserId != userData.Id {
		panic(errorgroup.UNAUTHORIZED)
	}

	err = eventservice.eventrepository.Delete(eventservice.db, event.Id)
	if err != nil {
		panic(err)
	}
	return event
}

func (eventservice *EventService) GetOne(id string) *responseentity.Event {
	event, err := eventservice.eventrepository.GetOne(eventservice.db, id)
	if err != nil {
		panic(err)
	}
	return event
}

func (eventservice *EventService) UpdateOneById(data map[string]interface{}) *responseentity.Event {
	id := data["id"].(string)
	payload := data["payload"].(*requestentity.Event)
	user := data["user"].(security.JwtClaim)

	event, err := eventservice.eventrepository.GetOne(eventservice.db, id)
	if err != nil {
		panic(err)
	}

	if event.UserId != user.UserData.Id {
		panic(errorgroup.UNAUTHORIZED)
	}

	result, err := eventservice.eventrepository.UpdateOneById(eventservice.db, id, payload)
	if err != nil {
		panic(err)
	}

	return result
}
