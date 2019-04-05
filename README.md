# go-sql-crud
[![GoDoc](https://godoc.org/github.com/TylerLafayette/go-sql-crud?status.svg)](https://godoc.org/github.com/TylerLafayette/go-sql-crud)  
🗄 quickly create simple crud style apis in go

```go
import "github.com/TylerLafayette/go-sql-crud"
```

This package allows you to easily create super simple CRUD (create, read, update, and delete) APIs based around an SQL-supported database (using the built-in sql package).

## Usage
Simply connect to your database, create a configuration, and create an http handler. The example below will create a route to search for a specific user in the dataabse.
```go
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
		// The Placeholder can be any squirrel (sql query builder) compatible placeholder format
		Placeholder: sqlcrud.Dollar,
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
```