package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/xzl-go/easygo/core"
	"github.com/xzl-go/easygo/logger"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleWebSocket 处理WebSocket连接
func HandleWebSocket(c *core.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			logger.Error("Failed to read message: %v", err)
			break
		}

		if err := conn.WriteMessage(messageType, message); err != nil {
			logger.Error("Failed to write message: %v", err)
			break
		}
	}
}
