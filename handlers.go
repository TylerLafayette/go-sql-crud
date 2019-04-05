package sqlcrud

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

// Init contains a configuration for the database connection.
type Init struct {
	Database *sql.DB
}

// Options contains options for a single handler.
type Options struct {
	Mode        string
	Table       string
	Limit       uint64
	QueryFields []Field
	Fields      []Field
}

// Field represents one field to query.
type Field struct {
	InputName string
	Name      string
	Validator func(interface{}) error
	Formatter func(interface{}) (interface{}, error)
}

// Response defines a single response.
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
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

		query := map[string]interface{}{}
		// Iterate through the QueryFields.
		for _, f := range o.QueryFields {
			// See if the user sent the correct values.
			value, ok := q[f.Name]
			if !ok || value == nil || len(value.(string)) == 0 {
				// Send a response.
				httpWrite(w, Response{
					Code:    422,
					Message: f.Name + " is missing",
				})
				return
			}

			// Check the value using the Validator.
			if f.Validator != nil {
				err := f.Validator(value)
				if err != nil {
					httpWrite(w, Response{
						Code:    422,
						Message: err.Error(),
					})
					return
				}
			}

			// Format the value using the Formatter.
			if f.Formatter != nil {
				var err error
				value, err = f.Formatter(value)
				if err != nil {
					httpWrite(w, Response{
						Code:    422,
						Message: err.Error(),
					})
					return
				}
			}

			query[f.Name] = value
		}

		selection := []string{}
		for _, v := range o.Fields {
			selection = append(selection, v.Name)
		}

		psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
		queryBuilder := psql.Select(strings.Join(selection, ",")).From(o.Table).Where(query)
		if o.Limit > 0 {
			queryBuilder = queryBuilder.Limit(o.Limit)
		}

		rows, err := queryBuilder.RunWith(i.Database).Query()
		if err != nil {
			httpWrite(w, Response{
				Code:    500,
				Message: err.Error(),
			})
			return
		}
		cols, _ := rows.Columns()

		output := []map[string]interface{}{}
		for rows.Next() {
			// Thanks: https://kylewbanks.com/blog/query-result-to-map-in-golang
			columns := make([]interface{}, len(cols))
			columnPointers := make([]interface{}, len(cols))
			for i := range columns {
				columnPointers[i] = &columns[i]
			}

			// Scan the result into the column pointers...
			if err := rows.Scan(columnPointers...); err != nil {
				httpWrite(w, Response{
					Code:    500,
					Message: err.Error(),
				})
				return
			}

			// Create our map, and retrieve the value for each column from the pointers slice,
			// storing it in the map with the name of the column as the key.
			m := map[string]interface{}{}
			for i, colName := range cols {
				val := columnPointers[i].(*interface{})
				m[colName] = *val
			}

			output = append(output, m)
		}

		if len(output) < 1 {
			httpWrite(w, Response{
				Code:    422,
				Message: "no results found",
			})
			return
		}

		httpWrite(w, Response{
			Code:    200,
			Message: "success",
			Data:    output,
		})
		return
	}

}
