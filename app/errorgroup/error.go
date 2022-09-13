package errorgroup

type Error struct {
	Code               int         `json:"code"`
	Message            string      `json:"message"`
	ActualErrorMessage interface{} `json:",omitempty"`
}

func (error *Error) Error() string {
	return error.Message
}
