package user

import "time"

type User struct {
	Id            int        `json:"id" db:"id"`
	Name          string     `json:"name" db:"name"`
	Email         string     `json:"email" db:"email"`
	Password      string     `json:"password" db:"password"`
	Role          string     `json:"role" db:"role"`
	Code          string     `json:"code" db:"code"`
	CodeExpiresAt *time.Time `json:"code_expires_at" db:"code_expires_at"`
	DeletedAt     *time.Time `json:"deleted_at" db:"deleted_at"`
}

type UserBasic struct {
	Id       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Role     string `json:"role" db:"role"`
	CinemaId *int   `json:"cinema_id" db:"cinema_id"`
}
