package booking

import (
	"github.com/cinema-booker/internal/user"
	"time"

	"github.com/cinema-booker/internal/cinema"
	"github.com/cinema-booker/internal/event"
	"github.com/cinema-booker/internal/room"
)

type User struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type EventBasic struct {
	Id     int           `json:"id" db:"id"`
	Cinema cinema.Cinema `json:"cinema"`
	Movie  event.Movie   `json:"movie"`
}

type SessionWithEvent struct {
	Id       int        `json:"id" db:"id"`
	Price    int        `json:"price" db:"price"`
	StartsAt time.Time  `json:"starts_at" db:"starts_at"`
	Room     room.Room  `json:"room"`
	Event    EventBasic `json:"event"`
}

type Booking struct {
	Id      int              `json:"id" db:"id"`
	Place   string           `json:"place" db:"place"`
	Status  string           `json:"status" db:"status"`
	User    User             `json:"user"`
	Session SessionWithEvent `json:"session"`
}

type BookingWithUsers struct {
	Booking     Booking        `json:"booking"`
	BookingUser User           `json:"booking_user"`
	CinemaUser  user.UserBasic `json:"cinema_user"`
}
