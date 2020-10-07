package mysqldb

import "time"

type (
	// Song entity
	Song struct {
		ID        int64     `json:"id"`
		Title     string    `json:"title" validate:"required"`
		Artist    string    `json:"artist" validate:"required"`
		UpdatedAt time.Time `json:"update_at"`
		CreatedAt time.Time `json:"created_at"`
	}
)
