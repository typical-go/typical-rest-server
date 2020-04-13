package repository

// BookColumns contain columns of book entity
var BookColumns = bookColumns{
	ID:        "id",
	Title:     "title",
	Author:    "author",
	UpdatedAt: "updated_at",
	CreatedAt: "created_at",
}

type bookColumns struct {
	ID        string
	Title     string
	Author    string
	UpdatedAt string
	CreatedAt string
}

func (c bookColumns) All() []string {
	return []string{c.ID, c.Title, c.Author, c.UpdatedAt, c.CreatedAt}
}
