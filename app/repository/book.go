package repository

import (
	"database/sql"
	"time"
)

var (
	bookTable   = "books"
	bookColumns = []string{"id", "title", "author", "created_at"}
)

// Book represented database model
type Book struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}

func scanBook(rows *sql.Rows) (book Book, err error) {
	err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.CreatedAt)
	return
}
