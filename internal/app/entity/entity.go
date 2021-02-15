package entity

import "time"

type (
	// Song entity
	// @entity (table:"songs" dialect:"mysql" ctor_db:"mysql")
	Song struct {
		ID        int64     `column:"id" option:"pk" json:"id"`
		Title     string    `column:"title" json:"title" validate:"required"`
		Artist    string    `column:"artist" json:"artist" validate:"required"`
		UpdatedAt time.Time `column:"updated_at" option:"now" json:"update_at"`
		CreatedAt time.Time `column:"created_at" option:"now,no_update" json:"created_at"`
	}
	// Book represented book model
	// @entity (table:"books" dialect:"postgres" ctor_db:"pg")
	Book struct {
		ID        int64     `column:"id" option:"pk" json:"id"`
		Title     string    `column:"title" json:"title" validate:"required"`
		Author    string    `column:"author" json:"author" validate:"required"`
		UpdatedAt time.Time `column:"updated_at" option:"now" json:"update_at"`
		CreatedAt time.Time `column:"created_at" option:"now,no_update" json:"created_at"`
	}
)
