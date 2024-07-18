package api

import (
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/cinema-booker/api/handler"
	"github.com/cinema-booker/internal/booking"
	"github.com/cinema-booker/internal/cinema"
	"github.com/cinema-booker/internal/event"
	"github.com/cinema-booker/internal/room"
	"github.com/cinema-booker/internal/session"
	"github.com/cinema-booker/internal/user"
	"github.com/gorilla/handlers"
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
	}).Methods(http.MethodGet)

	sessionStore := session.NewStore(s.db)
	sessionService := session.NewService(sessionStore)

	userStore := user.NewStore(s.db)
	userService := user.NewService(userStore)
	userHandler := handler.NewUserHandler(userService, userStore, sessionService)
	userHandler.RegisterRoutes(router)

	roomStore := room.NewStore(s.db)
	roomService := room.NewService(roomStore)
	cinemaStore := cinema.NewStore(s.db)
	cinemaService := cinema.NewService(cinemaStore)
	cinemaHandler := handler.NewCinemaHandler(cinemaService, roomService, userStore)
	cinemaHandler.RegisterRoutes(router)

	eventStore := event.NewStore(s.db)
	eventService := event.NewService(eventStore)
	eventHandler := handler.NewEventHandler(eventService, sessionService, userStore)
	eventHandler.RegisterRoutes(router)

	bookingStore := booking.NewStore(s.db)
	bookingService := booking.NewService(bookingStore, sessionStore)
	bookingHandler := handler.NewBookingHandler(bookingService, userStore)
	bookingHandler.RegisterRoutes(router)

	router.PathPrefix("/docs/swagger.json").Handler(http.StripPrefix("/docs", http.FileServer(http.Dir("./docs"))))

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		//TODO Change URL before put it in prod
		httpSwagger.URL("http://localhost:3000/docs/swagger.json"),
	))

	router.HandleFunc("/webhook", handler.HandleWebhook(bookingService)).Methods(http.MethodPost)

	websocketHandler := handler.NewWebSocketHandler()
	router.HandleFunc("/ws", websocketHandler.HandleWebSocket).Methods(http.MethodGet)

	// Configure CORS
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	log.Printf("🚀 Starting server on %s", s.address)
	return http.ListenAndServe(s.address, cors(router))
}
