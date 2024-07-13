package constants

type contextKey string

const (
	UserIDKey   contextKey = "userId"
	UserRoleKey contextKey = "userRole"
)

const (
	UserRoleAdmin   string = "ADMIN"
	UserRoleManager string = "MANAGER"
	UserRoleViewer  string = "VIEWER"
)
