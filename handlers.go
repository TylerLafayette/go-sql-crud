package sqlcrud

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// Init contains a configuration for the database connection.
type Init struct {
	Database *sql.DB
}

// Options contains options for a single handler.
type Options struct {
	Mode        string
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

// Response defines a single response.
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

// httpWrite writes a json response and sets the code in the http header.
func httpWrite(w http.ResponseWriter, r Response) {
	w.WriteHeader(r.Code)
	j := json.NewEncoder(w)
	j.Encode(r)
}

// GetRow returns a handler to get one row from the database and return it to the user.
func (i *Init) GetRow(o Options) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		for _, f := range o.QueryFields {
			if value, ok := q[f.Name]; !ok || len(value) < 1 {
				httpWrite(w, Response{
					Code:    404,
					Message: f.Name + " is missing",
				})
				return
			}
		}
	}
}
