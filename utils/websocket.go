package utils

import (
	"encoding/json"
	"fmt"
	"game-night/models"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 全局WebSocket管理器
var wsManager *WebSocketManager

// WebSocketManager WebSocket管理器
type WebSocketManager struct {
	connections map[string]*WSConnection
	groups      map[string]map[string]bool
	upgrader    websocket.Upgrader
	mutex       sync.RWMutex
}

// WSConnection WebSocket连接
type WSConnection struct {
	ID       string
	conn     *websocket.Conn
	send     chan []byte
	closed   bool
	mutex    sync.Mutex
	metadata map[string]interface{}
}

// InitWebSocket 初始化WebSocket
func InitWebSocket() {
	wsManager = &WebSocketManager{
		connections: make(map[string]*WSConnection),
		groups:      make(map[string]map[string]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	log.Println("WebSocket管理器已初始化")
}

// CreateWebSocketConnection 创建WebSocket连接
func CreateWebSocketConnection(c *gin.Context, connectionID string, metadata map[string]interface{}) (*WSConnection, error) {
	if wsManager == nil {
		return nil, fmt.Errorf("WebSocket管理器未初始化")
	}

	conn, err := wsManager.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil, fmt.Errorf("WebSocket升级失败: %w", err)
	}

	wsConn := &WSConnection{
		ID:       connectionID,
		conn:     conn,
		send:     make(chan []byte, 256),
		metadata: metadata,
	}

	wsManager.mutex.Lock()
	wsManager.connections[connectionID] = wsConn
	wsManager.mutex.Unlock()

	go wsConn.writePump()

	log.Printf("WebSocket连接创建: %s", connectionID)
	return wsConn, nil
}

// AddToGroup 将连接加入组
func AddToGroup(groupID, connectionID string) error {
	if wsManager == nil {
		return fmt.Errorf("WebSocket管理器未初始化")
	}

	wsManager.mutex.Lock()
	defer wsManager.mutex.Unlock()

	if _, exists := wsManager.connections[connectionID]; !exists {
		return fmt.Errorf("连接不存在: %s", connectionID)
	}

	if wsManager.groups[groupID] == nil {
		wsManager.groups[groupID] = make(map[string]bool)
	}

	wsManager.groups[groupID][connectionID] = true
	log.Printf("连接 %s 加入组 %s", connectionID, groupID)
	return nil
}

// RemoveFromGroup 从组中移除连接
func RemoveFromGroup(groupID, connectionID string) {
	if wsManager == nil {
		return
	}

	wsManager.mutex.Lock()
	defer wsManager.mutex.Unlock()

	if group, exists := wsManager.groups[groupID]; exists {
		delete(group, connectionID)
		if len(group) == 0 {
			delete(wsManager.groups, groupID)
		}
		log.Printf("连接 %s 从组 %s 移除", connectionID, groupID)
	}
}

// SendToConnection 发送消息到指定连接
func SendToConnection(connectionID string, message models.WebsocketMessage) error {
	if wsManager == nil {
		return fmt.Errorf("WebSocket管理器未初始化")
	}

	wsManager.mutex.RLock()
	conn, exists := wsManager.connections[connectionID]
	wsManager.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("连接不存在: %s", connectionID)
	}

	return conn.sendMessage(message)
}

// SendToGroup 发送消息到组内所有连接
func SendToGroup(groupID string, message models.WebsocketMessage, excludeConnIDs ...string) error {
	if wsManager == nil {
		return fmt.Errorf("WebSocket管理器未初始化")
	}

	wsManager.mutex.RLock()
	group, exists := wsManager.groups[groupID]
	if !exists {
		wsManager.mutex.RUnlock()
		return fmt.Errorf("组不存在: %s", groupID)
	}

	excludeMap := make(map[string]bool)
	for _, id := range excludeConnIDs {
		excludeMap[id] = true
	}

	var targetConnections []*WSConnection
	for connID := range group {
		if !excludeMap[connID] {
			if conn, exists := wsManager.connections[connID]; exists && !conn.isClosed() {
				targetConnections = append(targetConnections, conn)
			}
		}
	}
	wsManager.mutex.RUnlock()

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("消息序列化失败: %w", err)
	}

	for _, conn := range targetConnections {
		select {
		case conn.send <- messageBytes:
		default:
			log.Printf("发送消息失败，连接 %s 通道已满", conn.ID)
		}
	}

	return nil
}

// GetGroupConnectionCount 获取组内连接数量（业务层调用）
func GetGroupConnectionCount(groupID string) int {
	if wsManager == nil {
		return 0
	}

	wsManager.mutex.RLock()
	defer wsManager.mutex.RUnlock()

	if group, exists := wsManager.groups[groupID]; exists {
		return len(group)
	}
	return 0
}

// GetGroupConnectionIDs 获取组内连接ID列表（业务层调用）
func GetGroupConnectionIDs(groupID string) []string {
	if wsManager == nil {
		return []string{}
	}

	wsManager.mutex.RLock()
	defer wsManager.mutex.RUnlock()

	if group, exists := wsManager.groups[groupID]; exists {
		ids := make([]string, 0, len(group))
		for connID := range group {
			ids = append(ids, connID)
		}
		return ids
	}
	return []string{}
}

// RemoveConnection 移除连接（内部调用）
func RemoveConnection(connectionID string) {
	if wsManager == nil {
		return
	}

	wsManager.mutex.Lock()
	defer wsManager.mutex.Unlock()

	if conn, exists := wsManager.connections[connectionID]; exists {
		conn.close()
		delete(wsManager.connections, connectionID)

		// 从所有组中移除
		for groupID, connections := range wsManager.groups {
			delete(connections, connectionID)
			if len(connections) == 0 {
				delete(wsManager.groups, groupID)
			}
		}

		log.Printf("WebSocket连接移除: %s", connectionID)
	}
}

func (c *WSConnection) sendMessage(message models.WebsocketMessage) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.closed {
		return fmt.Errorf("连接已关闭")
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("消息序列化失败: %w", err)
	}

	select {
	case c.send <- messageBytes:
		return nil
	default:
		return fmt.Errorf("发送通道已满")
	}
}

func (c *WSConnection) close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.closed {
		c.closed = true
		close(c.send)
		c.conn.Close()
	}
}

func (c *WSConnection) isClosed() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.closed
}

func (c *WSConnection) writePump() {
	defer func() {
		c.conn.Close()
		RemoveConnection(c.ID)
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("WebSocket写入错误: %v", err)
				return
			}
		}
	}
}

// GetConnectionMetadata 获取连接元数据
func GetConnectionMetadata(connectionID, key string) (interface{}, bool) {
	if wsManager == nil {
		return nil, false
	}

	wsManager.mutex.RLock()
	conn, exists := wsManager.connections[connectionID]
	wsManager.mutex.RUnlock()

	if !exists {
		return nil, false
	}

	conn.mutex.Lock()
	defer conn.mutex.Unlock()

	if conn.metadata == nil {
		return nil, false
	}

	value, exists := conn.metadata[key]
	return value, exists
}
