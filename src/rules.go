package main

import (
	"encoding/json"
	"github.com/DemoLiang/wss/golib"
	"math"
)

//初始化规则处理函数
func InitRules() {
	LuckRulesFilter = map[LUCK_CARD_TYPE_ENUM]func(room *GameRoom, c *Connection) (err error){
		LUCK_CARD_TYPE__NO1:  LuckCardsFilterNO1,
		LUCK_CARD_TYPE__NO2:  LuckCardsFilterNO2,
		LUCK_CARD_TYPE__NO3:  LuckCardsFilterNO3,
		LUCK_CARD_TYPE__NO4:  LuckCardsFilterNO4,
		LUCK_CARD_TYPE__NO5:  LuckCardsFilterNO5,
		LUCK_CARD_TYPE__NO6:  LuckCardsFilterNO6,
		LUCK_CARD_TYPE__NO7:  LuckCardsFilterNO7,
		LUCK_CARD_TYPE__NO8:  LuckCardsFilterNO8,
		LUCK_CARD_TYPE__NO9:  LuckCardsFilterNO9,
		LUCK_CARD_TYPE__NO10: LuckCardsFilterNO10,
		LUCK_CARD_TYPE__NO11: LuckCardsFilterNO11,
		LUCK_CARD_TYPE__NO12: LuckCardsFilterNO12,
	}
	NewsRulesFilter = map[NEWS_CARD_TYPE_ENUM]func(room *GameRoom,c *Connection)(err error){
		NEWS_CARD_TYPE__NO1:  NewsCardsFilterNO1,
		NEWS_CARD_TYPE__NO2:  NewsCardsFilterNO2,
		NEWS_CARD_TYPE__NO3:  NewsCardsFilterNO3,
		NEWS_CARD_TYPE__NO4:  NewsCardsFilterNO4,
		NEWS_CARD_TYPE__NO5:  NewsCardsFilterNO5,
		NEWS_CARD_TYPE__NO6:  NewsCardsFilterNO6,
		NEWS_CARD_TYPE__NO7:  NewsCardsFilterNO7,
		NEWS_CARD_TYPE__NO8:  NewsCardsFilterNO8,
		NEWS_CARD_TYPE__NO9:  NewsCardsFilterNO9,
		NEWS_CARD_TYPE__NO10: NewsCardsFilterNO10,
		NEWS_CARD_TYPE__NO11: NewsCardsFilterNO11,
		NEWS_CARD_TYPE__NO12: NewsCardsFilterNO12,
	}
}



//运气卡处理函数
//遗失钱包，你失去300元，位于你后方的第一位玩家获得300元
func LuckCardsFilterNO1(room *GameRoom, c *Connection) (err error) {
	var flag bool
	for con, _ := range room.Money {
		if con == c {
			room.Money[c] -= 300
			flag = true
		}
		if con != c && flag {
			room.Money[c] += 300
		}
	}
	var ruleNO1 MessageLuckRuleFilterNO1
	ruleNO1.MessageType = MESSAGE_TYPE__LUCK_CARD__NO1
	ruleNO1.Code = c.Code
	ruleNO1.GameRoomId = room.Id
	ruleNO1.Money = -300

	//广播消息到游戏房间
	data, _ := json.Marshal(&ruleNO1)
	room.Broadcast <- data

	return nil
}

//黑历史被查，立即移动到监狱，并停留一回合
func LuckCardsFilterNO2(room *GameRoom, c *Connection) (err error) {
	var ruleNO2 MessageLuckRuleFilterNO2
	ruleNO2.MessageType = MESSAGE_TYPE__LUCK_CARD__NO2
	ruleNO2.Code = c.Code
	ruleNO2.GameRoomId = room.Id

	//广播消息到游戏房间
	data, _ := json.Marshal(&ruleNO2)
	room.Broadcast <- data
	room.StopStep[c] += len(room.Connections)
	return nil
}

//社会主义春风吹过，你可以立即免费升级一块抵偿
func LuckCardsFilterNO3(room *GameRoom, c *Connection) (err error) {
	//用户升级地产消息确认
	var updateLand MessageUserLandUpdate
	updateLand.MessageType = MESSAGE_TYPE__LAND_UPDATE
	updateLand.GameRoomId = room.Id
	updateLand.GameUserid = c.Code
	updateLand.Land = room.Map.ClientMap[c]
	for _, land := range updateLand.Land {
		updateLand.UpdateFee[&land] = GetLandUpdateFee(land)
	}
	updateLand.Number = 1
	data, _ := json.Marshal(&updateLand)
	c.Send <- data

	//收到确认消息后，确定是否升级土地
	comfirmFlag := c.GetHandlerConfirmData()
	if !comfirmFlag {
		golib.Log("收到确认信息，不升级地产\n")
	}

	return nil
}

//双十一期间，疯狂消费，支付300元
func LuckCardsFilterNO4(room *GameRoom, c *Connection) (err error) {

	var ruleNO4 MessageLuckRuleFilterNO4
	room.Money[c] -= 300

	ruleNO4.MessageType = MESSAGE_TYPE__LUCK_CARD__NO4
	ruleNO4.Code = c.Code
	ruleNO4.GameRoomId = room.Id
	ruleNO4.Money = -300

	//广播消息到游戏房间
	data, _ := json.Marshal(&ruleNO4)
	room.Broadcast <- data

	return nil
}

//前往九寨沟旅游，支付500元
func LuckCardsFilterNO5(room *GameRoom, c *Connection) (err error) {
	var ruleNO5 MessageLuckRuleFilterNO5
	room.Money[c] -= 500

	ruleNO5.MessageType = MESSAGE_TYPE__LUCK_CARD__NO5
	ruleNO5.Code = c.Code
	ruleNO5.GameRoomId = room.Id
	ruleNO5.Money = -500

	//广播消息到游戏房间
	data, _ := json.Marshal(&ruleNO5)
	room.Broadcast <- data
	return nil
}

//潜入银行 系统内部，从每位玩家手中收取300元
func LuckCardsFilterNO6(room *GameRoom, c *Connection) (err error) {
	var ruleNO6 MessageLuckRuleFilterNO6
	for con, _ := range room.Money {
		if con != c {
			room.Money[con] -= 300
			room.Money[c] += 300
			ruleNO6.Money += 500
		}
	}

	ruleNO6.MessageType = MESSAGE_TYPE__LUCK_CARD__NO5
	ruleNO6.Code = c.Code
	ruleNO6.GameRoomId = room.Id

	//广播消息到游戏房间
	data, _ := json.Marshal(&ruleNO6)
	room.Broadcast <- data

	return nil
}

//在香港乘坐豪华游轮，你可以选择支付1000元，并立即移动到起点领取奖励
func LuckCardsFilterNO7(room *GameRoom, c *Connection) (err error) {
	var moveStartPoint MessageMoveStartPoint
	moveStartPoint.MessageType = MESSAGE_TYPE__MOVE_2_START
	moveStartPoint.Money = 1000
	moveStartPoint.Code = c.Code
	moveStartPoint.GameRoomId = room.Id
	moveStartPoint.SEQ = NewID()

	data, _ := json.Marshal(&moveStartPoint)
	c.Send <- data
	//广播通知
	room.Broadcast <- data

	//收到确认消息后，确定是否花费1000元移动到起点
	comfirmFlag := c.GetHandlerConfirmData()
	if !comfirmFlag {
		golib.Log("收到确认信息，不移动到起点\n")
	}

	return nil
}

//FIXME 立即移动到你的左边手玩家的位置，并按该结果结算 需要修改，map不保证顺序性，可以改为数组形式
func LuckCardsFilterNO8(room *GameRoom, c *Connection) (err error) {
	var dice int
	var cPos,conPos Pos
	var cIdx int
	var cStepFlag bool
	//找到自己的位置，和右边下一个人的位置
	for idx,location := range room.Map.CurrentUserLocation{
		//先找到 c 的位置
		if location.C == c{
			cPos = location.Pos
			cIdx = idx
			break
		}
	}
	if cIdx < 1 {
		conPos = room.Map.CurrentUserLocation[cIdx-1].Pos
	}else {
		conPos = room.Map.CurrentUserLocation[room.MaxClientNumber-1].Pos

	}
	for _,data:=range room.Map.Map{
		if cPos.IsEqual(Pos{LocationX:data.LocationX,LocationY:data.LocationY}){
			cStepFlag = true
		}else {
			if cStepFlag{
				dice++
			}
		}
		if conPos.IsEqual(Pos{LocationX:data.LocationX,LocationY:data.LocationY}){
			break
		}
	}

	//根据计算的step移动
	room.GameUserMove(dice,c)

	return nil
}

//FIXME 立即移动到你右手边玩家的位置，并按该结果结算 需要修改，map不保证顺序性，可以改为数组形式
func LuckCardsFilterNO9(room *GameRoom, c *Connection) (err error) {
	var dice int
	var cPos,conPos Pos
	var cIdx int
	var cStepFlag bool
	//找到自己的位置，和右边下一个人的位置
	for idx,location := range room.Map.CurrentUserLocation{
		//先找到 c 的位置
		if location.C == c{
			cPos = location.Pos
			cIdx = idx
			break
		}
	}
	if cIdx < room.MaxClientNumber-1 {
		conPos = room.Map.CurrentUserLocation[cIdx+1].Pos
	}else {
		conPos = room.Map.CurrentUserLocation[0].Pos

	}
	for _,data:=range room.Map.Map{
		if cPos.IsEqual(Pos{LocationX:data.LocationX,LocationY:data.LocationY}){
			cStepFlag = true
		}else {
			if cStepFlag{
				dice++
			}
		}
		if conPos.IsEqual(Pos{LocationX:data.LocationX,LocationY:data.LocationY}){
			break
		}
	}

	//根据计算的step移动
	room.GameUserMove(dice,c)

	return nil
}

//发票刮中奖，获得400元
func LuckCardsFilterNO10(room *GameRoom, c *Connection) (err error) {
	var ruleNO10 MessageLuckRuleFilterNO10
	room.Money[c] += 400

	ruleNO10.MessageType = MESSAGE_TYPE__LUCK_CARD__NO10
	ruleNO10.Code = c.Code
	ruleNO10.GameRoomId = room.Id
	ruleNO10.Money += 400

	//广播消息到游戏房间
	data, _ := json.Marshal(&ruleNO10)
	room.Broadcast <- data

	return nil
}

//额外获得遗产，获得600元
func LuckCardsFilterNO11(room *GameRoom, c *Connection) (err error) {
	var ruleNO11 MessageLuckRuleFilterNO11
	room.Money[c] += 600

	ruleNO11.MessageType = MESSAGE_TYPE__LUCK_CARD__NO11
	ruleNO11.Code = c.Code
	ruleNO11.GameRoomId = room.Id
	ruleNO11.Money += 600

	//广播消息到游戏房间
	data, _ := json.Marshal(&ruleNO11)
	room.Broadcast <- data

	return nil
}

//TODO 购买最新款私人坐骑，支付400元，并立即额外进行一回合的行动
func LuckCardsFilterNO12(room *GameRoom, c *Connection) (err error) {

	return nil
}

//新闻卡处理函数
//TODO 投资项目分红，距离证券中心最近的玩家获得500元
func NewsCardsFilterNO1(room *GameRoom, c *Connection) (err error) {
	var cIndex ,conIndex int
	var stepList map[*Connection]int
	var minStep int
	//计算距离证券中心的位置
	land := room.GetLandByRole(GAME_ROLE__SECURITIES_CENTER)
	for _,userLocation := range room.Map.CurrentUserLocation {
		for idx, data := range room.Map.Map {
			if userLocation.Pos.IsEqual(Pos{LocationX: data.LocationX, LocationY: data.LocationY}) {
				cIndex = idx
			}
			if land.LocationX == data.LocationX && land.LocationY == data.LocationY && data.Role == GAME_ROLE__SECURITIES_CENTER {
				conIndex = idx
			}
			if minStep >= int(math.Abs(float64(cIndex-conIndex))){
				minStep = int(math.Abs(float64(cIndex-conIndex)))
			}
			stepList[userLocation.C] = int(math.Abs(float64(cIndex-conIndex)))

		}
	}
	//从银行派钱
	for c,step := range stepList{
		if step == minStep{
			room.Money[c] += 500
			room.Bank -= 500
		}
	}

	return nil
}

//社会发放福利，每位玩家获得1000元
func NewsCardsFilterNO2(room *GameRoom, c *Connection) (err error) {
	var ruleNO2List []MessageNewsRuleFilterNO2
	var index int = 0
	for con, _ := range room.Money {
		room.Money[c] += 1000
		ruleNO2List[index].MessageType = MESSAGE_TYPE__LUCK_CARD__NO11
		ruleNO2List[index].Code = con.Code
		ruleNO2List[index].GameRoomId = room.Id
		ruleNO2List[index].Money += 1000
		index++
	}

	//广播消息到游戏房间
	data, _ := json.Marshal(&ruleNO2List)
	room.Broadcast <- data
	return nil
}

//经营不善，拥有核能发电站的玩家失去300元
func NewsCardsFilterNO3(room *GameRoom, c *Connection) (err error) {
	var ruleNO3 MessageNewsRuleFilterNO3
	var flag bool = false
	for con, data := range room.Map.ClientMap {
		for _, land := range data {
			if land.Role == GAME_ROLE__NUCLEAR_POWER {
				flag = true
				room.Money[con] -= 300
				ruleNO3.Money = -300
				ruleNO3.GameRoomId = room.Id
				ruleNO3.Code = con.Code
				ruleNO3.MessageType = MESSAGE_TYPE__NEWS_CARD__NO3
			}
		}
	}

	if flag {
		golib.Log("没有找到拥有核能发电站的用户")
		return nil
	}
	//广播消息到游戏房间
	data, _ := json.Marshal(&ruleNO3)
	room.Broadcast <- data
	return nil
}

//经营不善，拥有污水处理厂的玩家失去300元
func NewsCardsFilterNO4(room *GameRoom, c *Connection) (err error) {
	var ruleNO4 MessageNewsRuleFilterNO4
	var flag bool = false
	for con, data := range room.Map.ClientMap {
		for _, land := range data {
			if land.Role == GAME_ROLE__SEWAGE_TREATMENT {
				flag = true
				room.Money[con] -= 300
				ruleNO4.Money = -300
				ruleNO4.GameRoomId = room.Id
				ruleNO4.Code = con.Code
				ruleNO4.MessageType = MESSAGE_TYPE__NEWS_CARD__NO4
			}
		}
	}

	if flag {
		golib.Log("没有找到拥有核能发电站的用户")
		return nil
	}
	//广播消息到游戏房间
	data, _ := json.Marshal(&ruleNO4)
	room.Broadcast <- data
	return nil
}

//经营不善，每位拥有运输业的玩家失去300元（大陆运输，大洋运输，空中运输）
func NewsCardsFilterNO5(room *GameRoom, c *Connection) (err error) {
	var ruleNO5 MessageNewsRuleFilterNO5
	var flag bool = false
	for con, data := range room.Map.ClientMap {
		for _, land := range data {
			if land.Role == GAME_ROLE__SEWAGE_TREATMENT {
				flag = true
				room.Money[con] -= 300
				ruleNO5.Money = -300
				ruleNO5.GameRoomId = room.Id
				ruleNO5.Code = con.Code
				ruleNO5.MessageType = MESSAGE_TYPE__NEWS_CARD__NO5
			}
		}
	}

	if flag {
		golib.Log("没有找到拥有核能发电站的用户")
		return nil
	}
	//广播消息到游戏房间
	data, _ := json.Marshal(&ruleNO5)
	room.Broadcast <- data
	return nil
	return nil
}

//政府公开补助土地少者500元
func NewsCardsFilterNO6(room *GameRoom, c *Connection) (err error) {
	var min int
	for _, data := range room.Map.ClientMap {
		tempMin := len(data)
		if tempMin < min {
			min = tempMin
		}
	}
	for con, data := range room.Map.ClientMap {
		if len(data) == min {
			room.Money[con] += 500
		}
	}

	return nil
}

//TODO 无名慈善家资助，每位玩家可以立即免费赎回一块抵偿
func NewsCardsFilterNO7(room *GameRoom, c *Connection) (err error) {
	return nil
}

//TODO 全体玩家参加狂欢节，在你下次行动结束前，所有玩家移动移动时都变为后退
func NewsCardsFilterNO8(room *GameRoom, c *Connection) (err error) {
	return nil
}

//经营不善，拥有建筑公司将自己的一个地产下降一级
func NewsCardsFilterNO9(room *GameRoom, c *Connection) (err error) {
	var ruleNO9 MessageNewRuleFilterNO9

	for con, landList := range room.Map.ClientMap {
		for _, land := range landList {
			if land.Role == GAME_ROLE__CONSTRUCTION_COMPANY {
				ruleNO9.MessageType = MESSAGE_TYPE__NEWS_CARD__NO9
				ruleNO9.Code = con.Code
				ruleNO9.GameRoomId = room.Id
				ruleNO9.LandList = landList
			}
		}
	}
	data, _ := json.Marshal(&ruleNO9)
	room.Broadcast <- data

	//收到确认消息后，确定是否升级土地
	comfirmFlag := c.GetHandlerConfirmData()
	if !comfirmFlag {
		golib.Log("收到确认信息，降级房产\n")
	}
	return nil
}

//TODO 百年一遇特大暴雨，所有玩家原地停留一回合
func NewsCardsFilterNO10(room *GameRoom, c *Connection) (err error) {
	return nil
}

//TODO 发生灵异事件，在你下次行动结束前，所有玩家都无须支付任何费用
func NewsCardsFilterNO11(room *GameRoom, c *Connection) (err error) {
	return nil
}

//所有玩家缴纳个人所得税，每块地产300元
func NewsCardsFilterNO12(room *GameRoom, c *Connection) (err error) {
	var ruleNO12List []MessageNewRuleFilterNO12
	var index int = 0
	for con, data := range room.Map.ClientMap {
		ruleNO12List[index].LandList = data
		ruleNO12List[index].GameRoomId = room.Id
		ruleNO12List[index].Code = con.Code
		ruleNO12List[index].MessageType = MESSAGE_TYPE__NEWS_CARD__NO12
		ruleNO12List[index].Money = int64(len(data) * 300)
		ruleNO12List[index].Number = int64(len(data))
		room.Money[con] -= ruleNO12List[index].Money
	}

	data, _ := json.Marshal(&ruleNO12List)
	room.Broadcast <- data
	return nil
}
