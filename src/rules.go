package main


func InitRules(){
	RulesFilter = map[LUCK_CARD_TYPE_ENUM]interface{}{
		LUCK_CARD_TYPE__NO1:LuckCardsFilterNO1,
		LUCK_CARD_TYPE__NO2:LuckCardsFilterNO2,

		NEWS_CARD_TYPE__NO1:NewsCardsFilterNO1,
		NEWS_CARD_TYPE__NO2:NewsCardsFilterNO2,
	}
}

//运气卡处理函数
func LuckCardsFilterNO1(room *GameRoom,c *Connection)(err error) {

	return nil
}

func LuckCardsFilterNO2(room *GameRoom,c *Connection)(err error) {

	return nil
}



//新闻卡处理函数
func NewsCardsFilterNO1(room *GameRoom,c *Connection)(err error){
	return nil
}

func NewsCardsFilterNO2(room *GameRoom, c *Connection) (err error) {
	return nil
}