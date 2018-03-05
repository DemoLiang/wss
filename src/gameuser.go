package main

import (
	"errors"
	"encoding/json"
)

func (c *Connection)GetMapLocation(p Pos,room *GameRoom) (idx int,pos Pos, err error) {
	pos=Pos{-1,-1}
	for idx,data := range room.Map.Map{
		if data.LocationX == p.LocationX && data.LocationY == p.LocationY{
			pos.LocationX = data.LocationX
			pos.LocationY = data.LocationY
			return idx,pos,nil
		}
	}
	return 0,pos,errors.New("get pos error")
}



//房间内的游戏用户移动其摇到的骰子的距离
func (c *Connection)GameUserMove(dice int,room *GameRoom) (err error) {
	var dstDice int = dice
	var ok bool
	var pos Pos
	if pos,ok = room.Map.CurrentUserLocation[c];ok{
		idx,_ ,err := c.GetMapLocation(pos,room)
		if err != nil{
			return errors.New("move error")
		}
		mapLen := len(room.Map.CurrentUserLocation)
		//如果已经是再次经过，则跳过起点
		if idx + dice >= mapLen+1{
			dstDice = dice +1

			//如果已经再过起点，则再给用户一部分钱
			room.BankSendMony(c,BANK_SEND_MONY)
		}
		//获取移动距离后的点的坐标
		mapPos := room.Map.Map[idx+dstDice]

		//更新坐标位置
		pos.LocationX = mapPos.LocationX
		pos.LocationY = mapPos.LocationY
	}
	//广播移动位置
	var userMove MessageGameUserMove
	userMove.MessageType = MESSAGE_TYPE__GAME_USER_MOVE
	userMove.MoveStep = dice
	userMove.MovePos = pos
	userMove.GameRoomId = room.Id
	data ,_:=json.Marshal(userMove)
	room.Broadcast <- data

	//TODO 移动到的位置，判断是否收租金，是否买地，是否不够钱需要抵押房产
	room.GameDoing(c)

	//判断输赢
	room.CheckGameDone()
	return nil
}


