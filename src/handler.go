package main

import (
	"encoding/json"
	"github.com/DemoLiang/wss/golib"
	"math/rand"
	"time"
)

func ShakeDice() (dice int) {
	dice = RandNumber() % 6
	return dice
}

func RandNumber() (number int) {
	source :=rand.NewSource(time.Now().Unix())
	r:=rand.New(source)
	number = r.Intn(9999)
	return
}


func (c *Connection)HandlerMessage(data []byte)  {
	var messageBasicInfo MessageBasicInfo
	json.Unmarshal(data,&messageBasicInfo)
	switch messageBasicInfo.MessageType {
	case MESSAGE_TYPE__CREATE_ROOM:
		var createRoom MessageCreateRoom
		json.Unmarshal(data,&createRoom)
		gameRoom := NewGameRoom(createRoom.Number)
		gameRoom.Register <- c
		GameRooms[gameRoom.Id] = gameRoom
	case MESSAGE_TYPE__JOIN_ROOM:
		var joinRoom MessageJoinRoom
		json.Unmarshal(data,&joinRoom)
		gameRoom := GetGameRoomById(joinRoom.GameRoomId)
		gameRoom.Register <-c
	case MESSAGE_TYPE__SHAKE_DICE:
		var shakeDice MessageGameShakeDice
		json.Unmarshal(data,&shakeDice)
		gameRoom := GetGameRoomById(shakeDice.GameRoomId)
		dice := ShakeDice()
		shakeDice.DiceNumber = dice
		data,_:=json.Marshal(&shakeDice)
		//广播给房间其它的小伙伴
		gameRoom.Broadcast <- data
	case MESSAGE_TYPE__LUCK_CARD:
		var luckCard MessageGameLuckCard
		gameRoom := GetGameRoomById(luckCard.GameRoomId)
		luckCardNo := gameRoom.LuckCard()
		luckCard.LuckCardNo = luckCardNo
		data,_:=json.Marshal(&luckCard)
		//广播给房间其它的小伙伴
		gameRoom.Broadcast <- data
	case MESSAGE_TYPE__NEWS_CARD:
		var newsCard MessageGameNewsCard
		gameRoom := GetGameRoomById(newsCard.GameRoomId)
		newsCardNo := gameRoom.NewsCard()
		newsCard.NewsCardNo = newsCardNo
		data,_:=json.Marshal(&newsCard)
		//广播给房间其它的小伙伴
		gameRoom.Broadcast <- data
	case MESSAGE_TYPE_GAME_USER_MOVE:

	default:
		golib.Log("default unknown message")
	}
}
