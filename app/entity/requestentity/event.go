package requestentity

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Event struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
}

func (e *Event) Validate() error {
	payload := *e
	return validation.ValidateStruct(
		&payload,
		validation.Field(&payload.Name, validation.Required),
		validation.Field(&payload.Name, validation.Required),
		validation.Field(&payload.Deadline, validation.Required, validation.Date(time.RFC3339)),
	)
}
