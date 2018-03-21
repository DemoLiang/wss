package main

import (
	"errors"
	"github.com/DemoLiang/wss/golib"
	"time"
)

//获取位置信息
func (c *Connection) GetMapLocation(p Pos, room *GameRoom) (idx int, pos Pos, err error) {
	pos = Pos{-1, -1}
	for idx, data := range room.Map.Map {
		if data.LocationX == p.LocationX && data.LocationY == p.LocationY {
			pos.LocationX = data.LocationX
			pos.LocationY = data.LocationY
			return idx, pos, nil
		}
	}
	return 0, pos, errors.New("get pos error")
}

//读取队列的内容，如果30S还不确认，则认为未收到确认消息
func (c *Connection) GetHandlerConfirmData() bool {
	timer := time.NewTicker(time.Second * 30)
	defer timer.Stop()
	for {
		select {
		case d := <-c.ConfirDataChan:
			golib.Log("%v\n", d)
			return true
		case <-timer.C:
			return false
		}
	}
	return false
}
