package models

// Game 游戏基础结构
type Game struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MinPlayers  int    `json:"min_players"`
	MaxPlayers  int    `json:"max_players"`
	Category    string `json:"category"`
}

// Room 房间结构
type Room struct {
	ID       string   `json:"id"`
	GameID   int      `json:"game_id"`
	GameName string   `json:"game_name"`
	Players  []Player `json:"players"`
	Status   string   `json:"status"` // waiting, playing, finished
	MaxSize  int      `json:"max_size"`
}

// Player 玩家结构
type Player struct {
	Nickname string `json:"nickname"`
	Role     string `json:"role,omitempty"`   // 角色
	Status   string `json:"status,omitempty"` // 状态
	Score    int    `json:"score,omitempty"`  // 分数
}

// 全局数据存储
var (
	Games = []Game{}
	Rooms = make(map[string]*Room)
)

// InitData 初始化数据
func InitData() {
	Games = []Game{
		{ID: 1, Name: "狼人杀", Description: "经典推理游戏", MinPlayers: 6, MaxPlayers: 20, Category: "推理"},
		{ID: 2, Name: "谁是卧底", Description: "找出卧底", MinPlayers: 3, MaxPlayers: 12, Category: "推理"},
		{ID: 3, Name: "UNO", Description: "经典纸牌游戏", MinPlayers: 2, MaxPlayers: 10, Category: "纸牌"},
		{ID: 4, Name: "剧本杀", Description: "角色扮演推理", MinPlayers: 4, MaxPlayers: 8, Category: "推理"},
	}
}

// GetGameByID 根据ID查找游戏
func GetGameByID(id int) *Game {
	for i, game := range Games {
		if game.ID == id {
			return &Games[i]
		}
	}
	return nil
}

// GeneratePlayerID 生成玩家ID
func GeneratePlayerID() string {
	return "player_" + string(rune(len(Rooms)*10+65))
}
