package main

import (
	"encoding/json"
	"github.com/DemoLiang/wss/golib"
	"github.com/DemoLiang/wss/golib/wechat"
	"math/rand"
	"time"
)

//摇骰子，生成小于等于6的随机数
func ShakeDice() (dice int) {
	dice = RandNumber() % 6
	return dice
}

//生成随机数
func RandNumber() (number int) {
	source := rand.NewSource(time.Now().Unix())
	r := rand.New(source)
	number = r.Intn(9999)
	return
}

//总体处理所有消息
func (c *Connection) HandlerMessage(data []byte) (err error) {
	var messageBasicInfo MessageBasicInfo
	json.Unmarshal(data, &messageBasicInfo)
	//if c.Session == "" && messageBasicInfo.MessageType != MESSAGE_TYPE__LOGIN_SERVER {
	//	c.Send <- []byte("请先登录")
	//	return golib.EInternalError
	//}

	switch messageBasicInfo.MessageType {
	case MESSAGE_TYPE__LOGIN_SERVER:
		var login MessageLogin
		var openid, session_key string
		json.Unmarshal(data, &login)
		golib.Log("code:%v\n", login.Code)
		openid, session_key = wechat.GetWeChatOpenIdByCode(login.Code)
		if openid == "" || session_key == "" {
			return golib.EInternalError
		}
		c.ClientInfo = ClientInfo{
			Code:   login.Code,
			OpenId: openid,
		}
		c.ClientInfo.Session = golib.MD5Sum(session_key)

		h.Register <- c
		golib.Log("c.Code:%v c.Openid:%v c.Session:%v\n", c.Code, c.OpenId, c.Session)
	case MESSAGE_TYPE__CREATE_ROOM:
		var createRoom MessageCreateRoom
		json.Unmarshal(data, &createRoom)
		//新建房间
		gameRoom := NewGameRoom(createRoom.Number)
		//向房间注册用户
		gameRoom.Register <- c
		//向游戏大厅注册房间
		h.RegisterRoom <- gameRoom
		//返回游戏房间信息给前端
		data, _ := json.Marshal(createRoom)
		c.Send <- data
	case MESSAGE_TYPE__JOIN_ROOM:
		var joinRoom MessageJoinRoom
		json.Unmarshal(data, &joinRoom)
		//获取房间
		gameRoom := GetGameRoomById(joinRoom.GameRoomId)
		//向房间注册用户
		gameRoom.Register <- c
		//返回前端房间信息,客户端信息
		joinRoom.ClientInfoList = GetGameRoomClientInfo(joinRoom.GameRoomId)
		data, _ := json.Marshal(joinRoom)
		c.Send <- data
		//广播消息到房间所有的用户
		gameRoom.Broadcast <- data
	case MESSAGE_TYPE__GAME_START:
		var gameStart MessageGameStart
		json.Unmarshal(data, &gameStart)

		//TODO 需要判断是否房主发起的开始游戏

		//获取房间
		gameRoom := GetGameRoomById(gameStart.GameRoomId)

		//把房间变为不可用，游戏开始
		gameRoom.SetRoomStatus(GAMEROOM_STATUS__GAMESTART)

		//返回前端房间信息,客户端信息
		gameStart.ClientInfoList = GetGameRoomClientInfo(gameStart.GameRoomId)
		data, _ := json.Marshal(gameStart)
		c.Send <- data
		//广播消息到房间的所有用户
		gameRoom.Broadcast <- data
	case MESSAGE_TYPE__SHAKE_DICE:
		var shakeDice MessageGameShakeDice
		json.Unmarshal(data, &shakeDice)
		//获取房间信息
		gameRoom := GetGameRoomById(shakeDice.GameRoomId)
		//摇动骰子
		dice := ShakeDice()
		shakeDice.DiceNumber = dice
		data, _ := json.Marshal(&shakeDice)
		//广播给房间其它的小伙伴
		gameRoom.Broadcast <- data
		gameRoom.SetRoomStatus(GAMEROOM_STATUS__DICE_DISAVAILABLE)
		//掷完骰子后，就自动移动
		gameRoom.GameUserMove(dice, c)
		//摇动完骰子，置为可用
		gameRoom.SetRoomStatus(GAMEROOM_STATUS__DICE_AVAILABLE)
	//case MESSAGE_TYPE__LUCK_CARD:
	//	var luckCard MessageGameLuckCard
	//	json.Unmarshal(data, &luckCard)
	//	gameRoom := GetGameRoomById(luckCard.GameRoomId)
	//	luckCardNo := gameRoom.LuckCard()
	//	luckCard.LuckCardNo = luckCardNo
	//	data, _ := json.Marshal(&luckCard)
	//	//广播给房间其它的小伙伴
	//	gameRoom.Broadcast <- data
	//case MESSAGE_TYPE__NEWS_CARD:
	//	var newsCard MessageGameNewsCard
	//	json.Unmarshal(data, &newsCard)
	//	gameRoom := GetGameRoomById(newsCard.GameRoomId)
	//	newsCardNo := gameRoom.NewsCard()
	//	newsCard.NewsCardNo = newsCardNo
	//	data, _ := json.Marshal(&newsCard)
	//	//广播给房间其它的小伙伴
	//	gameRoom.Broadcast <- data
	//case MESSAGE_TYPE__GAME_USER_MOVE:
	//	//
	case MESSAGE_TYPE__LAND_IMPAWN:
		var landImpawn MessageUserLandImpawn
		json.Unmarshal(data, &landImpawn)
		gameRoom := GetGameRoomById(landImpawn.GameRoomId)
		gameRoom.LandImpawn(c, landImpawn.LandList)

		//广播给房间的其它小伙伴，其进行了地产抵押
		gameRoom.Broadcast <- data
		//通知正在处理业务的端，已经确认抵押，可以进行下一步
		c.ConfirDataChan <- true
	case MESSAGE_TYPE__LAND_REDEEM:
		var landRedeem MessageUserLandRedeem
		json.Unmarshal(data, &landRedeem)
		gameRoom := GetGameRoomById(landRedeem.GameRoomId)
		gameRoom.LandRedeem(c, landRedeem.LandList)

		//广播给房间的其它小伙伴，其进行了地产赎回
		gameRoom.Broadcast <- data
	case MESSAGE_TYPE__BUY_LAND_CONFIRM:
		var buyLandConfirm MessageBuyLandConfirm
		json.Unmarshal(data, &buyLandConfirm)
		room := GetGameRoomById(buyLandConfirm.GameRoomId)
		room.BuyLand(c, buyLandConfirm.Land)
		c.ConfirDataChan <- buyLandConfirm.Confirem
	case MESSAGE_TYPE__LAND_UPDATE_CONFIRM:
		var updateLandConfirm MessageUpdateLandConfirm
		json.Unmarshal(data, &updateLandConfirm)
		room := GetGameRoomById(updateLandConfirm.GameRoomId)
		room.UpdateLand(c, updateLandConfirm.Land)
	default:
		//golib.Log("default unknown message")
	}
	return nil
}
