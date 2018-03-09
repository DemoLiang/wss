package main

import (
	"encoding/json"
	"github.com/DemoLiang/wss/golib"
	"github.com/DemoLiang/wss/golib/wechat"
	"math/rand"
	"time"
)

func ShakeDice() (dice int) {
	dice = RandNumber() % 6
	return dice
}

func RandNumber() (number int) {
	source := rand.NewSource(time.Now().Unix())
	r := rand.New(source)
	number = r.Intn(9999)
	return
}

func (c *Connection) HandlerMessage(data []byte) {
	var messageBasicInfo MessageBasicInfo
	json.Unmarshal(data, &messageBasicInfo)
	if c.Session == "" {
		c.Send <- []byte("请先登录")
	}

	switch messageBasicInfo.MessageType {
	case MESSAGE_TYPE__LOGIN_SERVER:
		var login MessageLogin
		var session_key string
		json.Unmarshal(data, &login)
		c.Code = login.Code
		c.OpenId, session_key = wechat.GetWeChatOpenIdByCode(c.Code)
		c.Session = golib.MD5Sum(session_key)
		golib.Log("c.Code:%v c.Openid:%v c.Session:%v\n", c.Code, c.OpenId, c.Session)
	case MESSAGE_TYPE__CREATE_ROOM:
		var createRoom MessageCreateRoom
		json.Unmarshal(data, &createRoom)
		gameRoom := NewGameRoom(createRoom.Number)
		gameRoom.Register <- c
		GameRooms[gameRoom.Id] = gameRoom
	case MESSAGE_TYPE__JOIN_ROOM:
		var joinRoom MessageJoinRoom
		json.Unmarshal(data, &joinRoom)
		gameRoom := GetGameRoomById(joinRoom.GameRoomId)
		gameRoom.Register <- c
	case MESSAGE_TYPE__SHAKE_DICE:
		var shakeDice MessageGameShakeDice
		json.Unmarshal(data, &shakeDice)
		gameRoom := GetGameRoomById(shakeDice.GameRoomId)
		dice := ShakeDice()
		shakeDice.DiceNumber = dice
		data, _ := json.Marshal(&shakeDice)
		//广播给房间其它的小伙伴
		gameRoom.Broadcast <- data

		//掷完骰子后，就自动移动
		c.GameUserMove(dice, gameRoom)
	case MESSAGE_TYPE__LUCK_CARD:
		var luckCard MessageGameLuckCard
		json.Unmarshal(data, &luckCard)
		gameRoom := GetGameRoomById(luckCard.GameRoomId)
		luckCardNo := gameRoom.LuckCard()
		luckCard.LuckCardNo = luckCardNo
		data, _ := json.Marshal(&luckCard)
		//广播给房间其它的小伙伴
		gameRoom.Broadcast <- data
	case MESSAGE_TYPE__NEWS_CARD:
		var newsCard MessageGameNewsCard
		json.Unmarshal(data, &newsCard)
		gameRoom := GetGameRoomById(newsCard.GameRoomId)
		newsCardNo := gameRoom.NewsCard()
		newsCard.NewsCardNo = newsCardNo
		data, _ := json.Marshal(&newsCard)
		//广播给房间其它的小伙伴
		gameRoom.Broadcast <- data
	case MESSAGE_TYPE__GAME_USER_MOVE:
		//
	case MESSAGE_TYPE__LAND_IMPAWN:
		var landImpawn MessageUserLandImpawn
		json.Unmarshal(data, &landImpawn)
		gameRoom := GetGameRoomById(landImpawn.GameRoomId)
		gameRoom.LandImpawn(c, landImpawn.LandList)

		//广播给房间的其它小伙伴，其进行了地产抵押
		gameRoom.Broadcast <- data
	case MESSAGE_TYPE__LAND_REDEEM:
		var landRedeem MessageUserLandRedeem
		json.Unmarshal(data, &landRedeem)
		gameRoom := GetGameRoomById(landRedeem.GameRoomId)
		gameRoom.LandRedeem(c, landRedeem.LandList)

		//广播给房间的其它小伙伴，其进行了地产赎回
		gameRoom.Broadcast <- data

	default:
		//golib.Log("default unknown message")
	}
}
