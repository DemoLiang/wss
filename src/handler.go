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
		//连接server,微信登录换取服务器登录
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
		//创建房间
		var createRoom MessageCreateRoom
		json.Unmarshal(data, &createRoom)
		//新建房间
		gameRoom := NewGameRoom(createRoom.Number)
		//向房间注册用户
		gameRoom.Register <- c
		//初始化房主
		gameRoom.Homeowner = c.Session
		//向游戏大厅注册房间
		h.RegisterRoom <- gameRoom
		//返回游戏房间信息给前端
		c.SendMessage(&createRoom)

	case MESSAGE_TYPE__JOIN_ROOM:
		//加入房间
		var joinRoom MessageJoinRoom
		json.Unmarshal(data, &joinRoom)
		//获取房间
		gameRoom := GetGameRoomById(joinRoom.GameRoomId)
		//向房间注册用户
		gameRoom.Register <- c
		//返回前端房间信息,客户端信息
		joinRoom.ClientInfoList = GetGameRoomClientInfo(joinRoom.GameRoomId)
		gameRoom.BroadcastMessage(&joinRoom)
	case MESSAGE_TYPE__GAME_START:
		//开始游戏
		var gameStart MessageGameStart
		json.Unmarshal(data, &gameStart)
		//获取房间
		gameRoom := GetGameRoomById(gameStart.GameRoomId)
		//判断请求开始游戏的人是谁
		if c.Session != gameRoom.Homeowner {
			golib.Log("不是房主，不能开始游戏")
			return nil
		}
		//把房间变为不可用，游戏开始
		gameRoom.SetRoomStatus(GAMEROOM_STATUS__GAMESTART)
		//返回前端房间信息,客户端信息
		gameStart.ClientInfoList = GetGameRoomClientInfo(gameStart.GameRoomId)
		gameRoom.BroadcastMessage(&gameStart)
	case MESSAGE_TYPE__SHAKE_DICE:
		//摇骰子
		var shakeDice MessageGameShakeDice
		json.Unmarshal(data, &shakeDice)
		//获取房间信息
		gameRoom := GetGameRoomById(shakeDice.GameRoomId)
		//TODO 判断用户是否在监狱，如果在监狱，则不允许摇动骰子，如果已经在监狱
		if value, ok := gameRoom.Prision[c]; ok {
			if value > 0 {
				golib.Log("用户还在监狱，不能摇骰子，不能移动：%v", value)
				var msgErr MessageError
				msgErr.GameRoomId = gameRoom.Id
				msgErr.Code = c.Code
				msgErr.MessageType = MESSAGE_TYPE__ERROR
				msgErr.ErrorDesc = "位于监狱，不能摇骰子"
				gameRoom.BroadcastMessage(&msgErr)
				return
			}
		} else {
			//减暂停回合计数
			gameRoom.DescPrision(c)
		}
		//TODO 如果是此时的规则为后退，则把骰子置为负数，则往后退
		//摇动骰子
		dice := ShakeDice()
		shakeDice.DiceNumber = dice
		data, _ := json.Marshal(&shakeDice)
		//广播给房间其它的小伙伴
		gameRoom.Broadcast <- data
		gameRoom.SetRoomStatus(GAMEROOM_STATUS__DICE_DISAVAILABLE)
		//增加判断，如果方向为反向，则将骰子职位负数
		if gameRoom.Direction == GAME_DIRETION__LEFT {
			dice = -dice
		}
		//掷完骰子后，就自动移动
		gameRoom.GameUserMove(dice, c)
		//摇动完骰子，置为可用
		gameRoom.SetRoomStatus(GAMEROOM_STATUS__DICE_AVAILABLE)
	case MESSAGE_TYPE__LAND_IMPAWN:
		//抵押地产
		var landImpawn MessageUserLandImpawn
		json.Unmarshal(data, &landImpawn)
		gameRoom := GetGameRoomById(landImpawn.GameRoomId)
		gameRoom.LandImpawn(c, landImpawn.LandList)

		//广播给房间的其它小伙伴，其进行了地产抵押
		gameRoom.Broadcast <- data
		//通知正在处理业务的端，已经确认抵押，可以进行下一步
		c.ConfirDataChan <- true
	case MESSAGE_TYPE__LAND_REDEEM:
		//赎回地产
		var landRedeem MessageUserLandRedeem
		json.Unmarshal(data, &landRedeem)
		gameRoom := GetGameRoomById(landRedeem.GameRoomId)
		gameRoom.LandRedeem(c, landRedeem.LandList)

		//广播给房间的其它小伙伴，其进行了地产赎回
		gameRoom.Broadcast <- data
	case MESSAGE_TYPE__BUY_LAND_CONFIRM:
		//确认买地产
		var buyLandConfirm MessageBuyLandConfirm
		json.Unmarshal(data, &buyLandConfirm)
		room := GetGameRoomById(buyLandConfirm.GameRoomId)
		room.BuyLand(c, buyLandConfirm.Land)
		c.ConfirDataChan <- buyLandConfirm.Confirem
	case MESSAGE_TYPE__LAND_UPDATE_CONFIRM:
		//确认升级地产
		var updateLandConfirm MessageUpdateLandConfirm
		json.Unmarshal(data, &updateLandConfirm)
		room := GetGameRoomById(updateLandConfirm.GameRoomId)
		room.UpdateLand(c, updateLandConfirm.Land)
		//广播消息
		room.BroadcastMessage(&updateLandConfirm)
	case MESSAGE_TYPE__LOGOUT:
		//退出房间
		var logout MessageLogoutRoom
		json.Unmarshal(data, &logout)
		room := GetGameRoomById(logout.GameRoomId)
		room.RoomStatusLock.Lock()
		defer room.RoomStatusLock.Unlock()
		//归还地产
		room.Map.Map = append(room.Map.Map, room.Map.ClientMap[c]...)
		//归还钱
		room.Bank += room.Money[c]
		//广播消息给其它客户端
		room.BroadcastMessage(&logout)
	default:
		golib.Log("default unknown message：%s\n", data)
	}
	return nil
}
