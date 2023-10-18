package websocket

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func Websocket(w http.ResponseWriter, r *http.Request, c *gin.Context) {
	// Upgrade the HTTP connection to a WebSocket connection
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusBadRequest)
		return
	}
	//defer conn.Close()

	path := c.Param("path")

	// Connect to the target WebSocket server
	targetConn, _, err := websocket.DefaultDialer.Dial("ws://192.168.0.106:2342"+path, nil)
	if err != nil {
		http.Error(w, "Failed to connect to target WebSocket server", http.StatusInternalServerError)
		return
	}

	// Proxy data between the client and target server
	go func() {
		defer targetConn.Close()
		defer conn.Close()
		for {
			messageType, p, err := conn.ReadMessage()
			fmt.Printf("conn.ReadMessage : %s %s ", messageType, p)
			if err != nil {
				return
			}
			if err := targetConn.WriteMessage(messageType, p); err != nil {
				return
			}
		}
	}()
	go func() {
		for {
			messageType, p, err := targetConn.ReadMessage()
			fmt.Printf("targetConn.ReadMessage : %s %s ", messageType, p)
			if err != nil {
				return
			}
			if err := conn.WriteMessage(messageType, p); err != nil {
				return
			}
		}
	}()
}

func IsWebSocketRequest(r *http.Request) bool {
	upgradeHeader := r.Header.Get("Upgrade")
	connectionHeader := r.Header.Get("Connection")

	return upgradeHeader == "websocket" && connectionHeader == "Upgrade"
}
