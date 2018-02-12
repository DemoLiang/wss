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
	Map map[*Connection]GameMap

	//钱
	Money map[*Connection]int

	//银行
	Bank int64
}

/*
0，0，0，0，0，0，0，0，0，0，0，0，0，0，0
0，1，1，1，1，1，1，1，1，1，1，1，1，1，0，
0，1，0，0，0，0，0，0，0，0，0，0，0，1，0，
0，1，0，0，0，0，0，0，0，0，0，0，0，1，0，
0，1，0，0，0，0，0，0，0，0，0，0，0，1，0，
0，1，0，0，0，0，0，0，0，0，0，0，0，1，0，
0，1，0，0，0，0，0，0，0，0，0，0，0，1，0，
0，1，0，0，0，0，0，0，0，0，0，0，0，1，0，
0，1，0，0，0，0，0，0，0，0，0，0，0，1，0，
0，1，0，0，0，0，0，0，0，0，0，0，0，1，0，
0，1，0，0，0，0，0，0，0，0，0，0，0，1，0，
0，1，0，0，0，0，0，0，0，0，0，0，1，1，0，
0，1，0，0，0，0，0，0，0，0，0，0，1，1，0，
0，1，1，1，1，1，1，1，1，1，1，1，1，0，0
0，0，0，0，0，0，0，0，0，0，0，0，0，0，0
*/

/*




运气，青岛，南京，监狱，
武汉				宁波，
成都				重庆，
				投资，
长沙												上海
投资												税务
拉沙			投资，昆明，杭州，苏州，新闻			澳门
厦门			开封，				 苏州			台湾
长春			大连，				 天津			钓鱼岛
杭州，运气，辽宁，					苏州，西安，沈阳，桂林
													起点
*/
type MapElement struct {
	Descript  string `json:"descript"`
	LocationX int    `json:"location_x"` //土地X坐标
	LocationY int    `json:"location_y"` //土地位置Y坐标
	Level     int    `json:"level"`      //土地星级
	Fee       int    `json:"fee"`        //购买基础费用
	RentFee   int    `json:"rent_fee"`   //收费租金
	Status    int    `json:"enable"`     //标记是否可用，已购买，空地，已被抵押
	Role      int    `json:"role"`       //标记是client地图元素，还是地图模块，还是运气牌模块
}
type GameMap struct {
	Client map[*Connection][]MapElement
	Map    []MapElement
}

//基础信息
type MessageBasicInfo struct {
	MessageType MESSAGE_TYPE_ENUM `json:"message_type"` //消息类型
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

type MessageGameLuckCard struct {
	MessageBasicInfo
	GameRoomId string `json:"game_room_id"` //房间ID号
	LuckCardNo int    `json:"luck_card_no"` //运气牌号
}

type MessageGameNewsCard struct {
	MessageBasicInfo
	GameRoomId string `json:"game_room_id"` //房间ID号
	NewsCardNo int    `json:"luck_card_no"` //新闻卡号
}
