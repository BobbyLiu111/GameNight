package main

import (
	"game-night/constant"
	"game-night/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Welcome(c *gin.Context) {
	c.JSON(http.StatusOK, models.WelcomeResponse{
		Message: "Welcome to Game Night!",
		Version: "v1.0.0",
	})
}

// GetGames 获取所有游戏
func GetGames(c *gin.Context) {
	c.JSON(http.StatusOK, models.GetGamesResponse{
		Games: models.Games,
	})
}

// GetRoomByID 获取房间
func GetRoomByID(c *gin.Context) {
	roomID := c.Param("id")
	room := models.Rooms[roomID]
	if room == nil {
		c.JSON(http.StatusNotFound, models.GetRoomByIDResponse{Room: nil, Message: "未找到房间"})
		return
	}
	c.JSON(http.StatusOK, models.GetRoomByIDResponse{
		Room:    room,
		Message: "Let's Game Night!",
	},
	)
}

// CreateRoom 创建房间
func CreateRoom(c *gin.Context) {
	req := models.CreateRoomRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	roomID := req.RoomID
	//判断房间是否存在
	if models.Rooms[roomID] != nil {
		c.JSON(http.StatusCreated, models.CreateRoomResponse{
			Message: "房间号已存在",
		})
		return
	}
	// 查找游戏
	selectedGame := models.GetGameByID(req.GameID)
	if selectedGame == nil {
		c.JSON(http.StatusBadRequest, models.CreateRoomResponse{
			Message: "游戏不存在",
		})
		return
	}

	// 创建房间
	maxSize := req.MaxSize
	if maxSize == 0 {
		maxSize = selectedGame.MaxPlayers
	}

	room := &models.Room{
		ID:       roomID,
		GameID:   selectedGame.ID,
		GameName: selectedGame.Name,
		Players: []models.Player{
			{Nickname: req.Creator},
		},
		Status:  constant.RoomStatusWaiting,
		MaxSize: maxSize,
	}

	models.Rooms[roomID] = room

	c.JSON(http.StatusCreated, models.CreateRoomResponse{Message: "创建成功", Room: room})
}

// JoinRoom 加入房间
func JoinRoom(c *gin.Context) {
	req := models.JoinRoomRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room, exists := models.Rooms[req.RoomID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "房间不存在"})
		return
	}

	if len(room.Players) >= room.MaxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "房间已满"})
		return
	}

	if room.Status != constant.RoomStatusWaiting {
		c.JSON(http.StatusBadRequest, gin.H{"error": "游戏已开始"})
		return
	}

	// 判断昵称是否重复
	for _, player := range room.Players {
		if player.Nickname == req.PlayerName {
			c.JSON(http.StatusBadRequest, gin.H{"error": "昵称不可用"})
			return
		}
	}

	// 添加玩家
	newPlayer := models.Player{
		Nickname: req.PlayerName,
	}
	room.Players = append(room.Players, newPlayer)

	c.JSON(http.StatusOK, models.JoinRoomResponse{
		PlayerName: newPlayer.Nickname,
		RoomID:     room.ID,
	})
}

func PlayCard(c *gin.Context) {

}

func SendWebsocketMessage(c *gin.Context) {
	msg := models.WebsocketMessage{}
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}
