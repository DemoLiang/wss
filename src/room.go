package main

import (
	"encoding/json"
	"errors"
	"github.com/DemoLiang/wss/golib"
	"github.com/segmentio/ksuid"
	"sync"
)

//var GameRooms map[string]*GameRoom

func NewGameRoom(number int) (gameRoom *GameRoom) {
	gameRoom = &GameRoom{
		Broadcast:       make(chan []byte),
		Register:        make(chan *Connection),
		Unregister:      make(chan *Connection),
		Connections:     make(map[*Connection]bool, number),
		LuckCards:       make(map[LUCK_CARD_TYPE_ENUM]bool, int(LUCK_CARD_TYPE__MAX)),
		NewsCards:       make(map[NEWS_CARD_TYPE_ENUM]bool, int(NEWS_CARD_TYPE__MAX)),
		MaxClientNumber: number,
		Bank:            INITIAL_BANK_MONEY,
	}
	//初始化roomid
	gameRoom.Id = newID()
	//启动房间注册函数
	go gameRoom.run()
	//初始化运气池
	gameRoom.InitLuckCardMap()
	//初始化新闻池
	gameRoom.InitNewsCardMap()
	//初始化地图
	gameRoom.InitGameMap()

	gameRoom.SetRoomStatus(GAMEROOM_STATUS__ENABLE)
	return
}
//初始化运气卡
func (r *GameRoom) InitLuckCardMap() {
	for i := int(LUCK_CARD_TYPE__MIN) + 1; i < int(LUCK_CARD_TYPE__MAX); i++ {
		r.LuckCards[LUCK_CARD_TYPE_ENUM(i)] = true
	}
}
//初始化新闻卡
func (r *GameRoom) InitNewsCardMap() {
	for i := int(NEWS_CARD_TYPE__MIN) + 1; i < int(NEWS_CARD_TYPE__MAX); i++ {
		r.NewsCards[NEWS_CARD_TYPE_ENUM(i)] = true
	}
}

func (r *GameRoom) LuckCard() (cardNo int) {
	cardNo = RandNumber() % int(LUCK_CARD_TYPE__MAX-1)
	if r.LuckCards[LUCK_CARD_TYPE_ENUM(cardNo)] == true {
		r.LuckCards[LUCK_CARD_TYPE_ENUM(cardNo)] = false
		return cardNo
	} else {
		var flag bool = false
		for idx, d := range r.LuckCards {
			if d == true {
				flag = true
				return int(idx)
			}
		}
		if flag == false {
			r.InitLuckCardMap()
			return r.LuckCard()
		}
	}
	return cardNo
}

func (r *GameRoom) NewsCard() (cardNo int) {
	cardNo = RandNumber() % int(NEWS_CARD_TYPE__MAX-1)
	if r.NewsCards[NEWS_CARD_TYPE_ENUM(cardNo)] == true {
		r.NewsCards[NEWS_CARD_TYPE_ENUM(cardNo)] = false
		return cardNo
	} else {
		var flag bool = false
		for idx, d := range r.NewsCards {
			if d == true {
				flag = true
				return int(idx)
			}
		}
		if flag == false {
			r.InitNewsCardMap()
			return r.NewsCard()
		}
	}
	return cardNo
}

//初始化地图
func (r *GameRoom) InitGameMap() {
	r.Map.Map = InitialGameMap.Map
}

//根据游戏房间ID号，获取游戏房间信息
func GetGameRoomById(id string) (gameRoom *GameRoom) {
	return h.GameRooms[id]
}

//根据游戏房间ID号，获取游戏房间用户的基础信息 FIXME ： 需要删除session信息
func GetGameRoomClientInfo(id string) (ClientInfoList []ClientInfo) {
	for c, _ := range h.GameRooms[id].Connections {
		ClientInfoList = append(ClientInfoList, c.ClientInfo)
	}
	return ClientInfoList
}

//房间内的游戏用户移动其摇到的骰子的距离
func (room *GameRoom) GameUserMove(dice int, c *Connection) (err error) {
	var dstDice int = dice
	var ok bool
	var pos Pos

	//获取移动的目标位置点
	if pos, ok = room.Map.CurrentUserLocation[c]; ok {
		idx, _, err := c.GetMapLocation(pos, room)
		if err != nil {
			return errors.New("move error")
		}
		mapLen := len(room.Map.CurrentUserLocation)
		//如果已经是再次经过，则跳过起点
		if idx+dice >= mapLen+1 {
			dstDice = dice + 1

			//如果已经再过起点，则再给用户一部分钱
			room.BankSendMony(c, BANK_SEND_MONY)
		}
		//获取移动距离后的点的坐标
		mapPos := room.Map.Map[idx+dstDice]

		//更新坐标点位置
		pos.LocationX = mapPos.LocationX
		pos.LocationY = mapPos.LocationY
	}
	//广播移动的目标位置点
	var userMove MessageGameUserMove
	userMove.MessageType = MESSAGE_TYPE__GAME_USER_MOVE
	userMove.MoveStep = dice
	userMove.MovePos = pos
	userMove.GameRoomId = room.Id
	data, _ := json.Marshal(userMove)
	room.Broadcast <- data

	//TODO 移动到的位置，判断是否收租金，是否买地，是否不够钱需要抵押房产
	room.Map.CurrentUserLocation[c] = pos
	room.GameDoing(c)

	//判断输赢
	gameFlag, c, _ := room.CheckGameDone()
	if gameFlag {
		var gameDone MessageGameDoenMessage
		gameDone.MessageType = MESSAGE_TYPE__GAME_DONE
		gameDone.GameRoomId = room.Id
		gameDone.Code = c.Code
		data, _ := json.Marshal(&gameDone)
		//广播游戏结束
		room.Broadcast <- data
	}

	return nil
}

//银行给用户派送钱
func (room *GameRoom) BankSendMony(c *Connection, mony int64) (err error) {
	var lock sync.Mutex
	lock.Lock()
	defer lock.Unlock()
	//银行支出
	room.Bank = room.Bank - mony
	//用户增加钱
	room.Money[c] = room.Money[c] + mony

	return nil
}

//获取房间状态
func (room *GameRoom) GetRoomStatus() (roomStatus GAMEROOM_STATUS_ENUM) {
	room.RoomStatusLock.Lock()
	defer room.RoomStatusLock.Unlock()
	return room.RoomStatus
}

//设置房间状态
func (room *GameRoom) SetRoomStatus(roomStatus GAMEROOM_STATUS_ENUM) {
	room.RoomStatusLock.Lock()
	defer room.RoomStatusLock.Unlock()
	room.RoomStatus = roomStatus
}

//用户掷完骰子后，检查需要做的动作，比如：付租金，买地，升级地产，抵押地产来付租金，
func (room *GameRoom) GameDoing(c *Connection) (err error) {
	var confirmData []byte
	for con, data := range room.Map.ClientMap {
		for index, mapLand := range data {
			switch mapLand.Role {
			case GAME_ROLE__LAND:
				//过路/自己的地，需要支付租金/升级地产
				if room.Map.CurrentUserLocation[c].IsEqual(Pos{mapLand.LocationX, mapLand.LocationY}) {
					if con == c {
						//自己的地，确认是否升级地产
						if room.Map.ClientMap[con][index].Level == int(LAND_LEVEL__MAX) {
							golib.Log("已经是最高等级地产%v\n")
							//最高级别地产，不能继续升级
							return
						}
						var land MessageUserLandUpdate
						land.Land = room.Map.ClientMap[con][index]
						land.UpdateFee = land.Land.Fee + int64(float64(land.Land.Level)*0.2+float64(land.Land.Level))*land.Land.Fee
						land.GameRoomId = room.Id
						land.MessageType = MESSAGE_TYPE__LAND_UPDATE
						confirmData, _ = json.Marshal(land)

						//发送消息，确认是否操作
						c.Send <- confirmData

						//收到确认消息后，确定是否升级土地
						comfirmFlag := c.GetHandlerConfirmData()
						if !comfirmFlag {
							golib.Log("收到确认信息，不升级地产\n")
						}
					} else {
						//TODO 路过别人的地，需要支付租金，支付租金不需要确认，如果需要抵押，则需要确认抵押的房产
						room.Money[c] = room.Money[c] - room.Map.ClientMap[con][index].RentFee
						room.Money[con] = room.Money[con] + room.Map.ClientMap[con][index].RentFee
						var land MessageUserPayRenFee
						land.RentFee = room.Map.ClientMap[con][index].RentFee
						land.GameRoomId = room.Id
						land.Land = room.Map.ClientMap[con][index]
						land.MessageType = MESSAGE_TYPE__PAY_RENT_FEE
						land.LandImpawn = room.Map.ClientMap[c]
						confirmData, _ = json.Marshal(land)
						c.Send <- confirmData

						//如果缴纳租金后，余额小于0，则需要抵押
						if room.Money[c] < 0 {
							//收到确认消息后，确定是否抵押完成，则进行下一步
							comfirmFlag := c.GetHandlerConfirmData()
							if comfirmFlag {
								golib.Log("收到确认信息，抵押地产\n")
							}
							//TODO 如果用户点击抵押那块地，则认为它不选，则系统随机抵押
						}
					}
				} else {
					//空地，发送消息是否买地
					var land MessageUserBuyLand
					land.Land = room.Map.ClientMap[con][index]
					land.GameRoomId = room.Id
					land.MessageType = MESSAGE_TYPE__BUY_LAND

					confirmData, _ = json.Marshal(land)
					c.Send <- confirmData
					//收到确认消息后，确定购买地产
					comfirmFlag := c.GetHandlerConfirmData()
					if comfirmFlag {
						golib.Log("收到确认信息，确定购买地产\n")
					}
				}
			case GAME_ROLE__LUCK:
				//运气
				var luckCard MessageGameLuckCard
				luckCard.GameRoomId = room.Id
				luckCardNo := room.LuckCard()
				luckCard.LuckCardNo = luckCardNo
				luckCard.MessageType = MESSAGE_TYPE__LUCK_CARD
				luckData, _ := json.Marshal(&luckCard)
				//处理运气牌
				room.HandlerLuckCards(c, luckCard.LuckCardNo)
				//广播给房间其它的小伙伴
				room.Broadcast <- luckData
			case GAME_ROLE__NEWS:
				//新闻
				var newsCard MessageGameNewsCard
				newsCard.GameRoomId = room.Id
				newsCardNo := room.NewsCard()
				newsCard.NewsCardNo = newsCardNo
				newsCard.MessageType = MESSAGE_TYPE__NEWS_CARD
				//处理新闻牌
				room.HandlerNewsCards(c, newsCard.NewsCardNo)
				//广播给房间其它的小伙伴
				newsData, _ := json.Marshal(&newsCard)
				room.Broadcast <- newsData
			case GAME_ROLE__SECURITIES_CENTER:
				//证券
			case GAME_ROLE__PRISION:
				//监狱
			case GAME_ROLE__JAIL:
				//入狱
			case GAME_ROLE_PARK:
				//公园
				var park MessageMove2Park
				park.GameRoomId = room.Id
				park.MessageType = MESSAGE_TYPE__MOVE_2_PARK
				park.Code = c.Code
				park.Money = 300
				room.Money[c] += 300
				//广播给房间其它的小伙伴
				parkData, _ := json.Marshal(&park)
				room.Broadcast <- parkData
			case GAME_ROLE_TAX_CENTER:
				//税务
			case GAME_ROLE__NUCLEAR_POWER:
				//核能发电
			case GAME_ROLE__CONSTRUCTION_COMPANY:
				//建筑公司
			case GAME_ROLE__CONTINENTAL_TRANSPORTION:
				//大陆运输
			case GAME_ROLE__TV_STATION:
				//电视台
			case GAME_ROLE__AIR_TRANSPORTION:
				//航空运输
			case GAME_ROLE__SEWAGE_TREATMENT:
				//污水处理
			case GAME_ROLE__OCEAN_TRANSPORTION:
				//大洋运输
			default:
			}
		}
	}

	return nil
}

//购买地产
func (room *GameRoom) BuyLand(c *Connection, land MapElement) (err error) {
	room.Map.ClientMap[c] = append(room.Map.ClientMap[c], land)
	return nil
}

//升级地产
func (room *GameRoom) UpdateLand(c *Connection, land MapElement) (err error) {
	for _, data := range room.Map.ClientMap[c] {
		if data.IsEqual(land) {
			//扣除钱 FIXME 需要判断钱够不够
			room.Money[c] = room.Money[c] - int64(float64(data.Level)*0.2+float64(data.Level))*data.Fee
			//升级土地级别 FIXME 今后需要用枚举值，判断土地是否到最大级别
			data.Level += 1
		}
	}
	return nil
}

//用户地产抵押
func (room *GameRoom) LandImpawn(c *Connection, mapList []MapElement) (err error) {
	for idx, data := range mapList {
		if room.Map.ClientMap[c][idx].IsEqual(data) {
			//抵押清算费用
			room.Money[c] = room.Money[c] + room.Map.ClientMap[c][idx].Fee
			//把地产变为不可用
			room.Map.ClientMap[c][idx].Status = 0
		}
	}

	return nil
}

//用户地产赎回
func (room *GameRoom) LandRedeem(c *Connection, mapList []MapElement) {
	for idx, data := range mapList {
		if room.Map.ClientMap[c][idx].IsEqual(data) {
			//判断地产是否是同一个地产，如果是同一个地产，则把地产赎回，并根据地产计算费用
			//支付费用，TODO 如果钱不够，不能赎回
			room.Money[c] = room.Money[c] - room.Map.ClientMap[c][idx].Fee
			//把地产变为可用
			room.Map.ClientMap[c][idx].Status = 1
		}
	}

	return
}

//处理运气卡
func (room *GameRoom) HandlerLuckCards(c *Connection, luckCardNo int) (err error) {

	return nil
}

//处理新闻卡
func (room *GameRoom) HandlerNewsCards(c *Connection, newCardNo int) (err error) {

	return nil
}

//判断游戏是否结束
func (r *GameRoom) CheckGameDone() (done bool, con *Connection, err error) {
	var count int = 0
	for c, data := range r.Money {
		if data >= GAME_DOEN_MONY {
			return true, c, nil
		}
		if data <= 0 {
			count++
		}
		//已经是最后的一个用户
		if count >= len(r.Connections)-1 {
			return true, c, nil
		}
	}
	var clienFlag int
	for c, data := range r.Map.ClientMap {
		if len(data) <= 0 {
			clienFlag += 1
		}
		if clienFlag >= len(r.Connections)-1 {
			return true, c, nil
		}
	}

	return false, nil, nil
}

//判断一个点是否跟自己属于同一个点
func (this Pos) IsEqual(pos Pos) bool {
	if this.LocationX == pos.LocationX && this.LocationY == pos.LocationY {
		return true
	}

	return false
}

//判断一个地图元素是否跟自己是同一个地图元素
func (m MapElement) IsEqual(m1 MapElement) bool {
	if m.LocationX == m1.LocationX && m.LocationY == m1.LocationY {
		return true
	}
	return false
}

func newID() string {
	return ksuid.New().String()
}

func (r *GameRoom) run() {
	for {
		select {
		case c := <-r.Register:
			//用户连接变为可用
			r.Connections[c] = true
			//初始化钱
			r.Money[c] = INITIAL_MONEY
			//初始化位置到起点
			r.Map.CurrentUserLocation[c] = Pos{LocationX: r.Map.Map[0].LocationX, LocationY: r.Map.Map[0].LocationY}
		case c := <-r.Unregister:
			if _, ok := r.Connections[c]; ok {
				r.Connections[c] = false
				//FIXME 回收地产
				//FIXME 重置为起点
				//FIXME 回收钱
				delete(r.Connections, c)
				close(c.Send)
			}
		case m := <-r.Broadcast:
			for c := range r.Connections {
				select {
				case c.Send <- m:
				default:
					delete(r.Connections, c)
					close(c.Send)
				}
			}
		}
	}
}
