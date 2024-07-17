package handler

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/cinema-booker/api/middleware"
	"github.com/cinema-booker/internal/user"
	"github.com/cinema-booker/pkg/errors"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	upgrader  websocket.Upgrader
	userStore user.UserStore
}

var (
	managerConnections = make(map[string]*websocket.Conn)
	connMutex          sync.Mutex
)

func NewWebSocketHandler(userStore user.UserStore) *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Adjust this to your security needs
			},
		},
		userStore: userStore,
	}
}

func (h *WebSocketHandler) RegisterRoutes(mux *mux.Router) {
	mux.Handle("/ws/{managerID}", errors.ErrorHandler(middleware.IsAuth(h.HandleWebSocket, h.userStore))).Methods(http.MethodGet)
}

func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	managerID := vars["managerID"]

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return err
	}
	defer conn.Close()

	connMutex.Lock()
	managerConnections[managerID] = conn
	connMutex.Unlock()

	fmt.Printf("Manager %s connected\n", managerID)

	defer func() {
		connMutex.Lock()
		delete(managerConnections, managerID)
		connMutex.Unlock()
		fmt.Printf("Manager %s disconnected\n", managerID)
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}
		fmt.Printf("Received from %s: %s\n", managerID, message)
	}

	return nil
}

func NotifyManager(managerID, message string) {
	connMutex.Lock()
	conn, ok := managerConnections[managerID]
	connMutex.Unlock()

	if ok {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			fmt.Println("Write error:", err)
		}
	} else {
		fmt.Printf("Manager %s not connected\n", managerID)
	}
}
