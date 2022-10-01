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
	Get(db *gorm.DB, pagination *requestentity.Pagination, ctx context.Context) (*[]responseentity.Event, error)
	Count(db *gorm.DB) int64
	GetOne(db *gorm.DB, id string) (*responseentity.Event, error)
	Delete(db *gorm.DB, id string) error
	UpdateOneById(db *gorm.DB, id string, payload *requestentity.Event) (*responseentity.Event, error)
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

func (eventrepository *EventRepository) Get(db *gorm.DB, pagination *requestentity.Pagination, ctx context.Context) (*[]responseentity.Event, error) {
	var events []postgresql.Event
	var result []responseentity.Event
	offset := (pagination.Current - 1) * pagination.Limit

	query := db.Order("created_at desc").Offset(offset).Limit(pagination.Limit).Find(&events).WithContext(ctx)
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

func (eventrepository *EventRepository) UpdateOneById(db *gorm.DB, id string, payload *requestentity.Event) (*responseentity.Event, error) {
	event := &postgresql.Event{ID: id}
	deadline, _ := time.Parse(time.RFC3339, payload.Deadline)
	query := db.Model(event).Updates(postgresql.Event{Name: payload.Name, Description: payload.Description, Deadline: deadline})

	if query.Error != nil {
		return nil, query.Error
	}

	query.Scan(event)
	return event.ToDomain(), nil
}

func (eventrepository *EventRepository) Count(db *gorm.DB) int64 {
	var count int64
	db.Model(&postgresql.Event{}).Count(&count)
	return count
}
