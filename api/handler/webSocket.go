package handler

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/cinema-booker/internal/user"
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

func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	managerID := r.URL.Query().Get("managerID")
	if managerID == "" {
		http.Error(w, "Manager ID is required", http.StatusBadRequest)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	connMutex.Lock()
	managerConnections[managerID] = conn
	connMutex.Unlock()

	fmt.Printf("Manager %s connected\n", managerID, managerConnections)

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
