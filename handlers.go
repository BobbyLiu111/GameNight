package main

import (
	"game-night/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Welcome 欢迎页面
func Welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "欢迎来到 Game Night 聚会桌游！",
		"version": "1.0.0",
	})
}

// GetGames 获取所有游戏
func GetGames(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"games": models.Games,
	})
}

// GetGame 获取特定游戏信息
func GetGame(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "欢迎来到 PartyBoard 聚会桌游！",
		"version": "1.0.0",
	})
}

// GetRooms 获取所有房间
func GetRooms(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "欢迎来到 PartyBoard 聚会桌游！",
		"version": "1.0.0",
	})
}

// GetRoom 获取特定房间信息
func GetRoom(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "欢迎来到 PartyBoard 聚会桌游！",
		"version": "1.0.0",
	})
}

// CreateRoom 创建房间
func CreateRoom(c *gin.Context) {
	var req struct {
		GameID  int    `json:"game_id" binding:"required"`
		MaxSize int    `json:"max_size"`
		Creator string `json:"creator" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找游戏
	selectedGame := models.GetGameByID(req.GameID)
	if selectedGame == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "游戏不存在"})
		return
	}

	// 创建房间
	roomID := models.GenerateRoomID()
	maxSize := req.MaxSize
	if maxSize == 0 {
		maxSize = selectedGame.MaxPlayers
	}

	room := &models.Room{
		ID:       roomID,
		GameID:   selectedGame.ID,
		GameName: selectedGame.Name,
		Players: []models.Player{
			{ID: models.GeneratePlayerID(), Nickname: req.Creator},
		},
		Status:  "waiting",
		MaxSize: maxSize,
	}

	models.Rooms[roomID] = room

	c.JSON(http.StatusCreated, room)
}

// JoinRoom 加入房间
func JoinRoom(c *gin.Context) {
	roomID := c.Param("id")

	var req struct {
		PlayerName string `json:"player_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room, exists := models.Rooms[roomID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "房间不存在"})
		return
	}

	if len(room.Players) >= room.MaxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "房间已满"})
		return
	}

	if room.Status != "waiting" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "游戏已开始"})
		return
	}

	// 检查玩家是否已在房间内
	for _, player := range room.Players {
		if player.Nickname == req.PlayerName {
			c.JSON(http.StatusBadRequest, gin.H{"error": "玩家已在房间内"})
			return
		}
	}

	// 添加玩家
	newPlayer := models.Player{
		ID:       models.GeneratePlayerID(),
		Nickname: req.PlayerName,
	}
	room.Players = append(room.Players, newPlayer)

	c.JSON(http.StatusOK, room)
}
