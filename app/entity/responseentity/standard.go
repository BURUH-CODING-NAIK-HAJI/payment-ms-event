package responseentity

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Page    interface{} `json:"page,omitempty"`
}
