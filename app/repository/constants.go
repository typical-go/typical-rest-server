package repository

// Table Name
const (
	bookTable = "books"
)

// Table Column Names
const (
	idColumn        = "id"
	updatedAtColumn = "updated_at"
	createdAtColumn = "created_at"

	// Book Table Column Names
	bookTitleColumn  = "title"
	bookAuthorColumn = "author"
)

// Table Columns
var (
	BookColumns = []string{idColumn, bookTitleColumn, bookAuthorColumn, updatedAtColumn, createdAtColumn}
)
