// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package db

import (
	"database/sql"
)

type Todo struct {
	ID          int64
	Description string
	UserID      sql.NullInt64
}

type User struct {
	ID        int64
	FirstName sql.NullString
	LastName  sql.NullString
	Password  string
	Email     string
}
