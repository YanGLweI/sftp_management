package models

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebSocketManager 管理WebSocket连接
type WebSocketManager struct {
	upgrader websocket.Upgrader
}

// NewWebSocketManager 创建新的WebSocket管理器
func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // 生产环境中应限制允许的域名
			},
		},
	}
}

// UpgradeConnection 升级HTTP连接到WebSocket
func (wm *WebSocketManager) UpgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := wm.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return nil, err
	}
	return conn, nil
}

// SendMessage 发送消息到WebSocket连接
func (wm *WebSocketManager) SendMessage(conn *websocket.Conn, message string) error {
	err := conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}
	return nil
}

// CloseConnection 关闭WebSocket连接
func (wm *WebSocketManager) CloseConnection(conn *websocket.Conn) {
	if conn != nil {
		conn.Close()
	}
}
