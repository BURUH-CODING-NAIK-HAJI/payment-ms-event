package eventservice

import (
	"github.com/rizface/golang-api-template/app/entity/requestentity"
	"github.com/rizface/golang-api-template/app/entity/responseentity"
	"github.com/rizface/golang-api-template/app/errorgroup"
	"github.com/rizface/golang-api-template/app/repository/eventrepository"
	"github.com/rizface/golang-api-template/system/security"
	"gorm.io/gorm"
)

type EventServiceInterface interface {
	Create(payload map[string]interface{}) *responseentity.Event
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
