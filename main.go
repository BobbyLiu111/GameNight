package main

import (
	"game-night/models"
	"game-night/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.InitWebSocket()
	r := gin.Default()
	models.InitData()
	// 基础路由
	r.GET("/", Welcome)
	r.POST("/ws", SendWebsocketMessage)
	// 游戏相关路由
	gameGroup := r.Group("/games")
	{
		gameGroup.GET("/", GetGames)
	}

	// 房间相关路由
	roomGroup := r.Group("/rooms")
	{
		roomGroup.GET("/:id", GetRoomByID)
		roomGroup.POST("/", CreateRoom)
		roomGroup.POST("/:id/join", JoinRoom)
	}
	// 玩家操作相关路由
	playerGroup := r.Group("/players")
	{
		playerGroup.POST("/play-card", PlayCard)
	}
	/*	testGroup := r.Group("/test")
		{
			testGroup.GET("/ws/:roomId", HandleWebSocket)
		}*/
	r.Run(":8080")
}
