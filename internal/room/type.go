package room

const (
	RoomTypeSmall  = "SMALL"
	RoomTypeMedium = "MEDIUM"
	RoomTypeLarge  = "LARGE"
)

type Room struct {
	Id     int    `json:"id" db:"id"`
	Number string `json:"number" db:"number"`
	Type   string `json:"type" db:"type"`
}
