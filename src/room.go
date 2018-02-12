package main

var GameRooms map[string]*GameRoom

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
	//启动房间注册函数
	go gameRoom.run()
	//初始化运气池
	gameRoom.InitLuckCardMap()
	//初始化新闻池
	gameRoom.InitNewsCardMap()
	//初始化地图
	gameRoom.InitGameMap()
	return
}

func (r *GameRoom) InitLuckCardMap() {
	for i := int(LUCK_CARD_TYPE__MIN) + 1; i < int(LUCK_CARD_TYPE__MAX); i++ {
		r.LuckCards[LUCK_CARD_TYPE_ENUM(i)] = true
	}
}

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

func (r *GameRoom)InitGameMap(){

}

func GetGameRoomById(id string) (gameRoom *GameRoom) {
	return GameRooms[id]
}

func (r *GameRoom) run() {
	for {
		select {
		case c := <-r.Register:
			r.Connections[c] = true
			r.Money[c] = INITIAL_MONEY
		case c := <-r.Unregister:
			if _, ok := r.Connections[c]; ok {
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
