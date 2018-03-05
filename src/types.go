package main

import "github.com/gorilla/websocket"

// connection is an middleman between the websocket connection and the hub.
type Connection struct {
	// The web socket connection
	Ws *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte
}

// hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients.
	Connections map[*Connection]bool

	// Inbound messages from the clients.
	Broadcast chan []byte

	// Register requests from the clients.
	Register chan *Connection

	// Unregister requests from clients.
	Unregister chan *Connection
	//clients join to the game room
	//JoinGameRoom chan *Connection

	//GameRooms map[*GameRoom]bool
}

type GameRoom struct {
	//game room id for other client join or create
	Id string
	// Registered game user clients.
	Connections map[*Connection]bool

	// Inbound messages from the other clients. and broadcast to other's clients
	Broadcast chan []byte

	// Register clients from the pool
	Register chan *Connection

	// Unregister clinets from game room
	Unregister chan *Connection

	//运气牌池
	LuckCards map[LUCK_CARD_TYPE_ENUM]bool

	//新闻卡池
	NewsCards map[NEWS_CARD_TYPE_ENUM]bool

	//房间人数
	MaxClientNumber int

	//地图位置
	Map GameMap

	//钱
	Money map[*Connection]int64

	//银行
	Bank int64
}

type MapElement struct {
	Descript  string `json:"descript"`
	LocationX int    `json:"location_x"` //土地X坐标
	LocationY int    `json:"location_y"` //土地位置Y坐标
	Level     int    `json:"level"`      //土地星级
	Fee       int64  `json:"fee"`        //购买基础费用
	RentFee   int64  `json:"rent_fee"`   //收费租金
	Status    int    `json:"enable"`     //标记是否可用，已购买，空地，已被抵押
	Role      int    `json:"role"`       //标记是client地图元素，还是地图模块，还是运气牌模块,起点
}

type Pos struct {
	LocationX int `json:"location_x"` //土地X坐标
	LocationY int `json:"location_y"` //土地位置Y坐标
}

type GameMap struct {
	//用户当前所在的位置
	CurrentUserLocation map[*Connection]Pos

	//每个client用户拥有的地产
	ClientMap map[*Connection][]MapElement `json:"client_map"`

	//所有的地图元素
	Map []MapElement `json:"map"`
}

//基础信息
type MessageBasicInfo struct {
	MessageType MESSAGE_TYPE_ENUM `json:"message_type"` //消息类型
	GameRoomId string `json:"game_room_id"` //房间ID号
}

//创建房间
type MessageCreateRoom struct {
	MessageBasicInfo
	Number int `json:"number"` //房间人数
}

//加入房间
type MessageJoinRoom struct {
	MessageBasicInfo
	GameRoomId string `json:"game_room_id"` //房间ID号
}

//请求摇骰子
type MessageGameShakeDice struct {
	MessageBasicInfo
	GameRoomId string `json:"game_room_id"` //房间ID号
	DiceNumber int    `json:"dice_number"`  //骰子点数
}

//运气牌
type MessageGameLuckCard struct {
	MessageBasicInfo
	GameRoomId string `json:"game_room_id"` //房间ID号
	LuckCardNo int    `json:"luck_card_no"` //运气牌号
}

//新闻消息
type MessageGameNewsCard struct {
	MessageBasicInfo
	GameRoomId string `json:"game_room_id"` //房间ID号
	NewsCardNo int    `json:"luck_card_no"` //新闻卡号
}

//游戏用户移动
type MessageGameUserMove struct {
	MessageBasicInfo
	GameUserId string `json:"game_user_id"`
	GameRoomId string `json:"game_room_id"`
	MoveStep   int    `json:"move_step"`
	MovePos    Pos    `json:"move_pos"`
}

//游戏用户买地
type MessageUserBuyLand struct {
	MessageBasicInfo
	GameUserid string     `json:"game_userid"`
	GameRoomId string     `json:"game_room_id"`
	Land       MapElement `json:"land"`
}

//用户支付租金消息
type MessageUserPayRenFee struct {
	MessageBasicInfo
	GameUserid            string     `json:"game_userid"`
	GameRecvRentfeeUserid string     `json:"game_recv_renfee_userid"`
	RentFee               int64      `json:"rent_fee"`
	GameRoomId            string     `json:"game_room_id"`
	Land                  MapElement `json:"land"`
}

//用户升级地产消息确认
type MessageUserLandUpdate struct {
	MessageBasicInfo
	GameUserid string     `json:"game_userid"`
	UpdateFee  int64      `json:"update_fee"`
	GameRoomId string     `json:"game_room_id"`
	Land       MapElement `json:"land"`
}
