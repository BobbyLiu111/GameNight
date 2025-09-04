package models

import (
	"game-night/constant"
)

// WebsocketMessage WebSocket消息基础结构
type WebsocketMessage struct {
	Type      constant.MessageType `json:"type"`           // 消息类型
	From      string               `json:"from,omitempty"` // 发送者ID
	RoomID    string               `json:"room_id"`        // 房间ID
	Content   interface{}          `json:"content"`        // 消息内容
	Timestamp int64                `json:"timestamp"`      // 时间戳
}

// ChatContent 聊天消息内容
type ChatContent struct {
	Text     string `json:"text"`
	Nickname string `json:"nickname"`
}

// GameActionContent 游戏操作内容
type GameActionContent struct {
	Action string                 `json:"action"`
	Target string                 `json:"target,omitempty"`
	Data   map[string]interface{} `json:"data,omitempty"`
}

// PlayerJoinContent 玩家加入内容
type PlayerJoinContent struct {
	PlayerID string `json:"player_id"`
	Nickname string `json:"nickname"`
	Message  string `json:"message"`
}

// PlayerLeaveContent 玩家离开内容
type PlayerLeaveContent struct {
	PlayerID string `json:"player_id"`
	Nickname string `json:"nickname"`
	Message  string `json:"message"`
}

// GameStateContent 游戏状态内容
type GameStateContent struct {
	Phase    string                 `json:"phase"`
	Round    int                    `json:"round"`
	Players  []interface{}          `json:"players,omitempty"`
	GameData map[string]interface{} `json:"game_data,omitempty"`
}

// NoticeContent 通知内容
type NoticeContent struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Level   string `json:"level"` // info, warning, error
}

// ErrorContent 错误内容
type ErrorContent struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
