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
package examples

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	sqlcrud "github.com/TylerLafayette/go-sql-crud"
)

func main() {
	db, err := sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/test")
	if err != nil {
		log.Fatal(err)
	}

	o := sqlcrud.Init{
		Database: db,
	}

	http.HandleFunc("/getUser", o.GetRow(sqlcrud.Options{
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
				Name: "userId",
			},
			sqlcrud.Field{
				Name: "username",
			},
			sqlcrud.Field{
				Name: "fullName",
			},
		},
	}))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
```