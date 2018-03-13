package main

import (
	"errors"
)

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
