package main

import (
	"github.com/gorilla/websocket"
	"sync"
)

// connection is an middleman between the websocket connection and the hub.
type Connection struct {
	// The web socket connection
	Ws *websocket.Conn

	//// wechat openid
	//OpenId string
	//
	////code wechat login response code
	//Code string
	//
	////session
	//Session string

	// Buffered channel of outbound messages.
	Send           chan []byte
	ConfirDataChan chan bool

	ClientInfo
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

	//注册游戏房间到游戏大厅
	RegisterRoom chan *GameRoom

	//房间列表
	GameRooms map[string]*GameRoom
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

	//房间状态，用于判断房间是否可用，游戏中的不可用，销毁的不可用，创建的时候可用
	RoomStatus GAMEROOM_STATUS_ENUM
	//房间状态锁
	RoomStatusLock sync.Mutex

	//TODO 游戏规则注册，目前采用全局注册规则过滤游戏 FIXME 赶紧可以用interface的接口形式实现，今后改吧
	GameRules map[GAME_RULE_ENUM]interface{}

	//某些规则会出发停留，停留的需要放入停留的
	StopStep map[*Connection]int
}

type MapElement struct {
	Descript  string          `json:"descript"`
	LocationX int             `json:"location_x"` //土地X坐标
	LocationY int             `json:"location_y"` //土地位置Y坐标
	Level     int             `json:"level"`      //土地星级
	Fee       int64           `json:"fee"`        //购买基础费用
	RentFee   int64           `json:"rent_fee"`   //收费租金
	Status    int             `json:"enable"`     //标记是否可用，已购买，空地，已被抵押
	Role      GAME_ROLES_ENUM `json:"role"`       //标记是client地图元素，还是地图模块，还是运气牌模块,起点
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

type ClientInfo struct {
	// wechat openid
	OpenId string

	//code wechat login response code
	Code string

	//session
	Session string
}

//基础信息
type MessageBasicInfo struct {
	MessageType MESSAGE_TYPE_ENUM `json:"message_type"` //消息类型
	GameRoomId  string            `json:"game_room_id"` //房间ID号
	Code string `json:"code"`	//用户登录session code
}

//创建房间
type MessageCreateRoom struct {
	MessageBasicInfo
	Number int `json:"number"` //房间人数
}

//加入房间
type MessageJoinRoom struct {
	MessageBasicInfo
	GameRoomId     string       `json:"game_room_id"` //房间ID号
	ClientInfoList []ClientInfo `json:"client_info_list"`
}

type MessageGameStart struct {
	MessageBasicInfo
	GameRoomId     string       `json:"game_room_id"` //房间ID号
	ClientInfoList []ClientInfo `json:"client_info_list"`
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
	GameUserid            string       `json:"game_userid"`
	GameRecvRentfeeUserid string       `json:"game_recv_renfee_userid"`
	RentFee               int64        `json:"rent_fee"`
	GameRoomId            string       `json:"game_room_id"`
	Land                  MapElement   `json:"land"`
	LandImpawn            []MapElement `json:"land_impawn"` //支付租金，可以抵押的房产
}

//用户升级地产消息确认
type MessageUserLandUpdate struct {
	MessageBasicInfo
	GameUserid string     `json:"game_userid"`
	UpdateFee  map[*MapElement]int64      `json:"update_fee"`
	GameRoomId string     `json:"game_room_id"`
	Land       []MapElement `json:"land"`
	Number int `json:"number"`	//可以升级的地产数目
}

//用户地产赎回消息
type MessageUserLandRedeem struct {
	MessageBasicInfo
	GameUserid string       `json:"game_userid"`
	GameRoomId string       `json:"game_room_id"`
	LandList   []MapElement `json:"land"`
}

//用户抵偿抵押消息
type MessageUserLandImpawn struct {
	MessageBasicInfo
	GameUserid string       `json:"game_userid"`
	GameRoomId string       `json:"game_room_id"`
	LandList   []MapElement `json:"land"`
}

//用户信息
type UserInfo struct {
	nickName  string `json:"nick_name"`  //用户昵称
	avatarUrl string `json:"avatar_url"` //用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。
	gender    string `json:"gender"`     //用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
	city      string `json:"city"`       //用户所在城市
	province  string `json:"province"`   //用户所在省份
	country   string `json:"country"`    //用户所在国家
	language  string `json:"language"`   //用户的语言，简体中文为zh_CN
}

//用户登录消息
type MessageLogin struct {
	MessageBasicInfo
	Code     string   `json:"code"`
	UserInfo UserInfo `json:"user_info"`
}

//购买地产确认消息
type MessageBuyLandConfirm struct {
	MessageBasicInfo
	Code     string     `json:"code"`
	Confirem bool       `json:"confirem"`
	Land     MapElement `json:"land"`
}

//游戏结束消息
type MessageGameDoenMessage struct {
	MessageBasicInfo
	Code string `json:"code"`
}

//升级地产确认消息
type MessageUpdateLandConfirm struct {
	MessageBasicInfo
	Code     string     `json:"code"`
	Confirem bool       `json:"confirem"`
	Land     MapElement `json:"land"`
}

//公园消息
type MessageMove2Park struct {
	MessageBasicInfo
	Code  string `json:"code"`
	Money int    `json:"money"`
}

//推送消息是否移动到起点
type MessageMoveStartPoint struct {
	MessageBasicInfo
	Code string `json:"code"`
	Money int64 `json:"money"`
	SEQ string `json:"seq"`
	ConfirmResult bool `json:"confirm_result"`
}


type MessageRuleFilterNO1 struct {
	MessageBasicInfo
	Money int64 `json:"money"`
}