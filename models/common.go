package models

type WelcomeResponse struct {
	Message string `json:"message"`
	Version string `json:"version"`
}

type GetGamesResponse struct {
	Games []Game `json:"games"`
}
type GetRoomByIDResponse struct {
	Message string `json:"message"`
	Room    *Room  `json:"room"`
}

// CreateRoomRequest CreateRoomReq 创建房间请求
type CreateRoomRequest struct {
	RoomID  string `json:"room_id"`
	GameID  int    `json:"game_id"`
	MaxSize int    `json:"max_size"`
	Creator string `json:"creator"`
}
type CreateRoomResponse struct {
	Message string `json:"message"`
	Room    *Room  `json:"room"`
}

// JoinRoomReq 加入房间请求
type JoinRoomReq struct {
	PlayerName string `json:"player_name" binding:"required,min=1,max=20" label:"玩家昵称"`
	Password   string `json:"password,omitempty" label:"房间密码"`
	Position   int    `json:"position,omitempty" binding:"omitempty,min=1" label:"指定位置"`
}

// StartGameReq 开始游戏请求
type StartGameReq struct {
	GameConfig map[string]interface{} `json:"game_config,omitempty" label:"游戏配置"`
}

// LeaveRoomReq 离开房间请求
type LeaveRoomReq struct {
	PlayerID string `json:"player_id" binding:"required" label:"玩家ID"`
}

// UpdateRoomReq 更新房间请求
type UpdateRoomReq struct {
	MaxSize     *int    `json:"max_size,omitempty" binding:"omitempty,min=2,max=20" label:"最大人数"`
	Password    *string `json:"password,omitempty" binding:"omitempty,max=20" label:"房间密码"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=100" label:"房间描述"`
}

// GetRoomsReq 获取房间列表请求
type GetRoomsReq struct {
	Status   string `form:"status" binding:"omitempty,oneof=waiting playing finished" label:"房间状态"`
	GameID   int    `form:"game_id" binding:"omitempty,min=1" label:"游戏ID"`
	Page     int    `form:"page" binding:"omitempty,min=1" label:"页码"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=50" label:"每页数量"`
	Keyword  string `form:"keyword" binding:"omitempty,max=20" label:"搜索关键词"`
}

// PlayerActionReq 玩家游戏内操作请求
type PlayerActionReq struct {
	PlayerID   string                 `json:"player_id" binding:"required" label:"玩家ID"`
	ActionType string                 `json:"action_type" binding:"required" label:"操作类型"`
	Target     string                 `json:"target,omitempty" label:"目标"`
	Data       map[string]interface{} `json:"data,omitempty" label:"操作数据"`
}

// KickPlayerReq 踢出玩家请求
type KickPlayerReq struct {
	PlayerID string `json:"player_id" binding:"required" label:"玩家ID"`
	Reason   string `json:"reason,omitempty" binding:"omitempty,max=100" label:"踢出原因"`
}

// SetPlayerReadyReq 设置玩家准备状态请求
type SetPlayerReadyReq struct {
	PlayerID string `json:"player_id" binding:"required" label:"玩家ID"`
	IsReady  bool   `json:"is_ready" label:"是否准备"`
}
