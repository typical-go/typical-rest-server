package repository

import (
	"database/sql"
	"time"

	validator "gopkg.in/go-playground/validator.v9"
)

var (
	bookTable   = "books"
	bookColumns = []string{"id", "title", "author", "created_at"}
)

// Book represented database model
type Book struct {
	ID        int       `json:"id"`
	Title     string    `json:"title" validate:"required"`
	Author    string    `json:"author" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}

func scanBook(rows *sql.Rows) (book Book, err error) {
	err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.CreatedAt)
	return
}

// Validate book
func (b *Book) Validate() error {
	validate := validator.New()
	return validate.Struct(b)
}
