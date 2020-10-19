package mysqldb

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
)
