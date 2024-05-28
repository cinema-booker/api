package movie

type Movie struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	Language    string `json:"language" db:"language"`
	Poster      string `json:"poster" db:"poster"`
	Backdrop    string `json:"backdrop" db:"backdrop"`
}
