package eventrepository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rizface/golang-api-template/app/entity/requestentity"
	"github.com/rizface/golang-api-template/app/entity/responseentity"
	"github.com/rizface/golang-api-template/database/postgresql"
	"gorm.io/gorm"
)

type EventRepositoryInterface interface {
	Create(db *gorm.DB, userId string, payload *requestentity.Event) (*responseentity.Event, error)
	Get(db *gorm.DB, ctx context.Context) (*[]responseentity.Event, error)
	GetOne(db *gorm.DB, id string) (*responseentity.Event, error)
	Delete(db *gorm.DB, id string) error
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
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Deadline:    deadline,
	}
	err = db.Create(&event).Error
	if err != nil {
		return nil, err
	}
	return event.ToDomain(), nil
}

func (eventrepository *EventRepository) Get(db *gorm.DB, ctx context.Context) (*[]responseentity.Event, error) {
	var events []postgresql.Event
	var result []responseentity.Event
	query := db.Order("created_at desc").Find(&events).WithContext(ctx)
	if query.Error != nil {
		return nil, query.Error
	}

	for _, event := range events {
		result = append(result, *event.ToDomain())
	}

	return &result, nil
}

func (eventrepository *EventRepository) GetOne(db *gorm.DB, id string) (*responseentity.Event, error) {
	event := new(postgresql.Event)
	query := db.First(event, &postgresql.Event{ID: id})
	if query.Error != nil {
		return nil, query.Error
	}

	return event.ToDomain(), nil
}

func (eventrepository *EventRepository) Delete(db *gorm.DB, id string) error {
	query := db.Delete(&postgresql.Event{ID: id})
	return query.Error
}
