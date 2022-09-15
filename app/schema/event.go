package schema

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func ValidateEventId(id string) error {
	err := validation.Validate(id, validation.Required, is.UUID)
	if err == nil {
		return nil
	}
	return validation.NewError("400", err.Error())
}
