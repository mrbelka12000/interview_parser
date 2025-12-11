package ws

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/mrbelka12000/interview_parser/internal/config"
)

// Upgrader is used to upgrade HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}

	go handleConnection(conn)
}

func handleConnection(conn *websocket.Conn) {
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message from WS connection:", err)
			break
		}
		fmt.Printf("Received: %s\n", message)

		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			fmt.Println("Error writing message to WS connection:", err)
			break
		}
	}
}

func RunServer(cfg *config.Config) error {
	http.HandleFunc("/ws", wsHandler)

	return http.ListenAndServe(fmt.Sprintf(":%v", cfg.WSServerPort), nil)
}
