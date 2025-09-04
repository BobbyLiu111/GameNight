package constant

// MessageType 消息类型
type MessageType string

const (
	// 连接管理
	MsgTypeJoin      MessageType = "join"      // 加入房间
	MsgTypeLeave     MessageType = "leave"     // 离开房间
	MsgTypeHeartbeat MessageType = "heartbeat" // 心跳

	// 聊天消息
	MsgTypeChat MessageType = "chat" // 聊天

	// 游戏消息
	MsgTypeGameStart  MessageType = "game_start"  // 游戏开始
	MsgTypeGameAction MessageType = "game_action" // 游戏操作
	MsgTypeGameState  MessageType = "game_state"  // 游戏状态更新

	// 玩家状态
	MsgTypePlayerJoin  MessageType = "player_join"  // 玩家加入
	MsgTypePlayerLeave MessageType = "player_leave" // 玩家离开

	// 系统消息
	MsgTypeNotice MessageType = "notice" // 系统通知
	MsgTypeError  MessageType = "error"  // 错误消息
)
