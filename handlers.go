package sqlcrud

import (
	"database/sql"
	"net/http"
)

// Init contains a configuration for the database connection.
type Init struct {
	Database *sql.DB
}

// Options contains options for a single handler.
type Options struct {
	QueryFields []Field
	Fields      []Field
}

// Field represents one field to query.
type Field struct {
	InputName string
	Name      string
	Validator func(string) (string, error)
	Formatter func(string) (string, error)
}

// GetRow returns a handler to get one row from the database and return it to the user.
func (i *Init) GetRow(o Options) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
