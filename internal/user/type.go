package user

type User struct {
	Id          int    `json:"id" db:"id"`
	FirstName   string `json:"first_name" db:"first_name"`
	LastName    string `json:"last_name" db:"last_name"`
	Email       string `json:"email" db:"email"`
	Password    string `json:"password" db:"password"`
	Role        string `json:"role" db:"role"`
	Status      string `json:"status" db:"status"`
	AccountType string `json:"account_type" db:"account_type"`
}
