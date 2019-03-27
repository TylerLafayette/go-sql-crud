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
	// Set Content-Type header to application/json.
	w.Header().Set("Content-Type", "application/json")
	// Write the Response HTTP code into the header.
	w.WriteHeader(r.Code)
	// Create a new JSON encoder based on the ResponseWriter.
	j := json.NewEncoder(w)
	// Encode/send the Response data as JSON.
	j.Encode(r)
}

// getParameters returns the parameters from the GET request.
func getParameters(r *http.Request) map[string]interface{} {
	m := map[string]interface{}{}
	for k, v := range r.URL.Query() {
		m[k] = v[0]
	}

	return m
}

// postParameters returns the parameters from the POST request.
func postParameters(r *http.Request) map[string]interface{} {
	var m map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		// Return an empty map if there is an error processing the JSON.
		return map[string]interface{}{}
	}

	return m
}

// GetRow returns a handler to get one row from the database and return it to the user.
func (i *Init) GetRow(o Options) func(w http.ResponseWriter, r *http.Request) {
	// Return an HTTP hqndler.
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the parameters.
		q := map[string]interface{}{}
		if o.Mode == "GET" {
			q = getParameters(r)
		} else if o.Mode == "POST" {
			q = postParameters(r)
		}

		// Iterate through the QueryFields.
		for _, f := range o.QueryFields {
			// See if the user sent the correct values.
			if value, ok := q[f.Name]; !ok || value == nil || len(value.(string)) == 0 {
				// Send a response.
				httpWrite(w, Response{
					Code:    422,
					Message: f.Name + " is missing",
				})
				return
			}
		}
	}
}
