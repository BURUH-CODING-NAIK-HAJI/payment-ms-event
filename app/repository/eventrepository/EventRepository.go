package eventrepository

import (
	"time"

	"github.com/google/uuid"
	"github.com/rizface/golang-api-template/app/entity/requestentity"
	"github.com/rizface/golang-api-template/app/entity/responseentity"
	"github.com/rizface/golang-api-template/database/postgresql"
	"gorm.io/gorm"
)

type EventRepositoryInterface interface {
	Create(db *gorm.DB, userId string, payload *requestentity.Event) (*responseentity.Event, error)
}

type EventRepository struct {
}

func New() EventRepositoryInterface {
	return &EventRepository{}
}

func (eventrepository *EventRepository) Create(db *gorm.DB, userId string, payload *requestentity.Event) (*responseentity.Event, error) {
	deadline, err := time.Parse(time.RFC3339, payload.Deadline)
	if err != nil {
		return nil, err
	}
	event := &postgresql.Event{
		ID:          uuid.NewString(),
		UserID:      userId,
		Name:        payload.Name,
		Description: payload.Description,
		CreateAt:    time.Now(),
		UpdatedAt:   time.Now(),
		Deadline:    deadline,
	}
	err = db.Create(&event).Error
	if err != nil {
		return nil, err
	}
	return event.ToDomain(), nil
}
