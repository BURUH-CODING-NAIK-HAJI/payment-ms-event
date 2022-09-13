package postgresql

import (
	"time"

	"github.com/rizface/golang-api-template/app/entity/responseentity"
)

type Event struct {
	ID          string `gorm:"primaryKey"`
	UserID      string `gorm:"not null;"`
	Name        string `gorm:"not null"`
	Description string
	CreateAt    time.Time
	UpdatedAt   time.Time
	Deadline    time.Time `gorm:"not null"`
}

func (e *Event) ToDomain() *responseentity.Event {
	return &responseentity.Event{
		Id:          e.ID,
		UserId:      e.UserID,
		Name:        e.Name,
		Description: e.Description,
		CreatedAt:   e.CreateAt,
		UpdatedAt:   e.UpdatedAt,
		Deadline:    e.Deadline,
	}
}
