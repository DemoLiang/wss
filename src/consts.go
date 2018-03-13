package main

type MESSAGE_TYPE_ENUM int

const (
	_ MESSAGE_TYPE_ENUM = iota
	//登录服务器
	MESSAGE_TYPE__LOGIN_SERVER
	//创建房间
	MESSAGE_TYPE__CREATE_ROOM
	//加入房间
	MESSAGE_TYPE__JOIN_ROOM
	//开始游戏
	MESSAGE_TYPE__GAME_START
	//摇骰子
	MESSAGE_TYPE__SHAKE_DICE
	//移动消息
	MESSAGE_TYPE__GAME_USER_MOVE
	//运气卡消息
	MESSAGE_TYPE__LUCK_CARD
	//新闻卡消息
	MESSAGE_TYPE__NEWS_CARD
	//付租金消息
	MESSAGE_TYPE__PAY_RENT_FEE
	//升级地产消息
	MESSAGE_TYPE__LAND_UPDATE
	//购买地产消息
	MESSAGE_TYPE__BUY_LAND
	//购买地产确认消息
	MESSAGE_TYPE__BUY_LAND_CONFIRM
	//抵押地产消息
	MESSAGE_TYPE__LAND_IMPAWN
	//赎回地产消息
	MESSAGE_TYPE__LAND_REDEEM
)

type LUCK_CARD_TYPE_ENUM int

const (
	LUCK_CARD_TYPE__MIN LUCK_CARD_TYPE_ENUM = iota
	LUCK_CARD_TYPE__NO1
	LUCK_CARD_TYPE__NO2
	LUCK_CARD_TYPE__MAX
)

type NEWS_CARD_TYPE_ENUM int

const (
	NEWS_CARD_TYPE__MIN NEWS_CARD_TYPE_ENUM = iota
	NEWS_CARD_TYPE__NO1
	_
	NEWS_CARD_TYPE__NO2 = 3 + iota
	NEWS_CARD_TYPE__MAX
)

const (
	//初始化游戏，每个用户的钱
	INITIAL_MONEY = 15000
	//初始化游戏，银行的钱
	INITIAL_BANK_MONEY = 50000000
	//初始化游戏，每次过起点，送给用户的钱
	BANK_SEND_MONY = 3000
	//判断游戏结束的钱，当个人用户的钱达到50000时就算用户赢了
	GAME_DOEN_MONY = 50000
)

type GAMEROOM_STATUS_ENUM int

const (
	_ GAMEROOM_STATUS_ENUM = iota
	//房间可用
	GAMEROOM_STATUS__ENABLE
	//房间游戏开始
	GAMEROOM_STATUS__GAMESTART
	//房间不可用
	GAMEROOM_STATUS__DISABLE
)

//游戏规则枚举
type GAME_RULE_ENUM int

const (
	_ GAME_RULE_ENUM = iota
	GAME_RULE__THROUGH_START_POINT
)

//地图元素角色
type GAME_ROLES_ENUM int

const (
	_ GAME_ROLES_ENUM = iota
	//土地
	GAME_ROLE__LAND
	//起点
	GAME_ROLE__START_POINT
	//运气
	GAME_ROLE__LUCK
	//新闻
	GAME_ROLE__NEWS
	//用户
	GAME_ROLE__USER
	//证券
	GAME_ROLE__SECURITIES_CENTER
	//监狱
	GAME_ROLE__PRISION
	//入狱
	GAME_ROLE__JAIL
	//公园
	GAME_ROLE_PARK
	//税务
	GAME_ROLE_TAX_CENTER
)
const (
	//投资
	GAME_ROLE__INVESTMENT_START = 100 + iota

	//核能发电
	GAME_ROLE__NUCLEAR_POWER
	//建筑公司
	GAME_ROLE__CONSTRUCTION_COMPANY
	//大陆运输
	GAME_ROLE__CONTINENTAL_TRANSPORTION
	//电视台
	GAME_ROLE__TV_STATION
	//航空运输
	GAME_ROLE__AIR_TRANSPORTION
	//污水处理
	GAME_ROLE__SEWAGE_TREATMENT
	//大洋运输
	GAME_ROLE__OCEAN_TRANSPORTION

	GAME_ROLE__INVESTMENT_END
)

//土地星级
type LAND_LEVELS_ENUM int

const (
	//最小星级
	LAND_LEVEL__MIN LAND_LEVELS_ENUM = iota
	//一个星级
	LAND_LEVEL__START1
	//二星级
	LAND_LEVEL__START2
	//三星级
	LAND_LEVEL__START3
	//最大星级
	LAND_LEVEL__MAX = LAND_LEVEL__START3
)
