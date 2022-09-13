package responseentity

import "time"

type Event struct {
	Id          string    `json:"id"`
	UserId      string    `json:"userId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Deadline    time.Time `json:"deadline"`
}
