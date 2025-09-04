package main

import (
	"game-night/models"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	models.InitData()
	// 基础路由
	r.GET("/", Welcome)

	// 游戏相关路由
	gameGroup := r.Group("/games")
	{
		gameGroup.GET("/", GetGames)
		gameGroup.GET("/:id", GetGame)
	}

	// 房间相关路由
	roomGroup := r.Group("/rooms")
	{
		roomGroup.GET("/", GetRooms)
		roomGroup.GET("/:id", GetRoom)
		roomGroup.POST("/", CreateRoom)
		roomGroup.POST("/:id/join", JoinRoom)
	}

	r.Run(":8080")
}
