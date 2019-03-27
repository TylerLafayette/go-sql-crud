package main

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
	db, err := sql.Open("postgres", "user=tyler dbname=tyler sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	o := sqlcrud.Init{
		Database: db,
	}

	http.HandleFunc("/getUser", o.GetRow(sqlcrud.Options{
		Mode:  "GET",
		Table: "users",
		QueryFields: []sqlcrud.Field{
			sqlcrud.Field{
				Name: "username",
				Validator: func(i interface{}) error {
					if len(i.(string)) < 4 || len(i.(string)) > 56 {
						// If the username is too short or too long, return an error to stop the request.
						return errors.New("username invalid")
					}

					return nil
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
