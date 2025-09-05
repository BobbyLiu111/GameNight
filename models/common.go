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

// JoinRoomRequest 加入房间请求
type JoinRoomRequest struct {
	PlayerName string `json:"player_name"`
	RoomID     string `json:"room_id"`
}

// JoinRoomResponse 加入房间请求
type JoinRoomResponse struct {
	PlayerName string `json:"player_name"`
	RoomID     string `json:"room_id"`
}
