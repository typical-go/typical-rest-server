package repository

// Table Name
const (
	bookTable = "books"
)

// Table Column Names
const (
	idColumn        = "id"
	updatedAtColumn = "created_at"
	createdAtColumn = "created_at"

	// Book Table Column Names
	bookTitleColumn  = "title"
	bookAuthorColumn = "author"
)

// Table Columns
var (
	bookColumns = []string{idColumn, bookTitleColumn, bookAuthorColumn, updatedAtColumn, createdAtColumn}
)
