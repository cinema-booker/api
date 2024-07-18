package session

import (
	"time"

	"github.com/cinema-booker/internal/room"
)

type Session struct {
	Id        int        `json:"id" db:"id"`
	Price     int        `json:"price" db:"price"`
	StartsAt  time.Time  `json:"starts_at" db:"starts_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	Room      room.Room  `json:"room"`
}

type FlatDashboardResponse struct {
	TotalBookings          int `json:"total_bookings"`
	TotalCinemas           int `json:"total_cinemas"`
	TotalRevenue           int `json:"total_revenue"`
	TotalEvents            int `json:"total_events"`
	TotalConfirmedBookings int `json:"total_confirmed_bookings"`
	TotalPendingBookings   int `json:"total_pending_bookings"`
	TotalManagers          int `json:"total_managers"`
	TotalViewers           int `json:"total_viewers"`
}
