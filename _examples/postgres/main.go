package postgres

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	sqlcrud "github.com/TylerLafayette/go-sql-crud"
	// Postgres driver
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "user=tyler dbname=tyler ssl-mode=disable")
	if err != nil {
		log.Fatal(err)
	}

	o := sqlcrud.Init{
		Database: db,
	}

	http.HandleFunc("/getUser", o.GetRow(sqlcrud.Options{
		Mode: "GET",
		QueryFields: []sqlcrud.Field{
			sqlcrud.Field{
				Name: "username",
				Validator: func(i string) (string, error) {
					if len(i) < 4 || len(i) > 56 {
						// If the username is too short or too long, return an error to stop the request.
						return i, errors.New("username invalid")
					}

					return i, nil
				},
			},
		},
		Fields: []sqlcrud.Field{
			sqlcrud.Field{
				Name: "id",
			},
			sqlcrud.Field{
				Name: "username",
			},
			sqlcrud.Field{
				Name: "preferences",
			},
		},
	}))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
