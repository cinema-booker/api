package api

import (
	"log"
	"net/http"
	"os"

	"github.com/cinema-booker/api/handler"
	"github.com/cinema-booker/internal/booking"
	"github.com/cinema-booker/internal/cinema"
	"github.com/cinema-booker/internal/event"
	"github.com/cinema-booker/internal/movie"
	"github.com/cinema-booker/internal/room"
	"github.com/cinema-booker/internal/user"
	"github.com/cinema-booker/third_party/tmdb"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type APIServer struct {
	address string
	db      *sqlx.DB
}

func NewAPIServer(address string, db *sqlx.DB) *APIServer {
	return &APIServer{
		address: address,
		db:      db,
	}
}

func (s *APIServer) Start() error {
	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	tmdbService := tmdb.NewTMDBService(os.Getenv("TMDB_API_KEY"))

	userStore := user.NewStore(s.db)
	userService := user.NewService(userStore)
	userHandler := handler.NewUserHandler(userService, userStore)
	userHandler.RegisterRoutes(router)

	roomStore := room.NewStore(s.db)
	roomService := room.NewService(roomStore)
	cinemaStore := cinema.NewStore(s.db)
	cinemaService := cinema.NewService(cinemaStore)
	cinemaHandler := handler.NewCinemaHandler(cinemaService, roomService, userStore)
	cinemaHandler.RegisterRoutes(router)

	eventStore := event.NewStore(s.db)
	eventService := event.NewService(eventStore, tmdbService)
	eventHandler := handler.NewEventHandler(eventService, userStore)
	eventHandler.RegisterRoutes(router)

	bookingStore := booking.NewStore(s.db)
	bookingService := booking.NewService(bookingStore)
	bookingHandler := handler.NewBookingHandler(bookingService, userStore)
	bookingHandler.RegisterRoutes(router)

	movieService := movie.NewService(tmdbService)
	movieHandler := handler.NewMovieHandler(movieService, userStore)
	movieHandler.RegisterRoutes(router)

	log.Printf("🚀 Starting server on %s", s.address)
	return http.ListenAndServe(s.address, router)
}
