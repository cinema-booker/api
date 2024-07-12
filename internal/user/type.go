package user

import "time"

const (
	UserRoleAdmin   string = "ADMIN"
	UserRoleManager string = "MANAGER"
	UserRoleViewer  string = "VIEWER"
)

type User struct {
	Id        int        `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Email     string     `json:"email" db:"email"`
	Password  string     `json:"password" db:"password"`
	Role      string     `json:"role" db:"role"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}
