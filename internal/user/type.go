package user

import (
	"database/sql"
)

const (
	UserRoleAdmin   string = "ADMIN"
	UserRoleManager string = "MANAGER"
	UserRoleViewer  string = "VIEWER"
)

type User struct {
	Id        int          `json:"id" db:"id"`
	FirstName string       `json:"first_name" db:"first_name"`
	LastName  string       `json:"last_name" db:"last_name"`
	Email     string       `json:"email" db:"email"`
	Password  string       `json:"password" db:"password"`
	Role      string       `json:"role" db:"role"`
	DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at"`
}
