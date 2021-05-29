package entity

import "time"

// TODO: generate repo code from table instead of `@dbrepo` annotation

type (
	// Book represented book model
	// @dbrepo (table:"books" dialect:"postgres" ctor_db:"pg")
	Book struct {
		ID        int64     `column:"id" option:"pk" json:"id"`
		Title     string    `column:"title" json:"title" validate:"required"`
		Author    string    `column:"author" json:"author" validate:"required"`
		UpdatedAt time.Time `column:"updated_at" option:"now" json:"update_at"`
		CreatedAt time.Time `column:"created_at" option:"now,no_update" json:"created_at"`
	}
)
