package librarydb

import "time"

type (
	// Book represented database model
	// @entity (table:"books" dialect:"postgres")
	Book struct {
		ID        int64     `json:"id"`
		Title     string    `json:"title" validate:"required"`
		Author    string    `json:"author" validate:"required"`
		UpdatedAt time.Time `json:"update_at"`
		CreatedAt time.Time `json:"created_at"`
	}
)
