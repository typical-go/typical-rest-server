package repository

import (
	"database/sql"
	"time"

	validator "gopkg.in/go-playground/validator.v9"
)

// Book represented database model
type Book struct {
	ID        int       `json:"id"`
	Title     string    `json:"title" validate:"required"`
	Author    string    `json:"author" validate:"required"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func scanBook(rows *sql.Rows) (book Book, err error) {
	err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.UpdatedAt, &book.CreatedAt)
	return
}

// Validate book
func (b *Book) Validate() error {
	validate := validator.New()
	return validate.Struct(b)
}
