package card

import "github.com/taiprogramer/simple-poker-game/backend/db"

func FindCardIDByNumberAndSuit(number, suit int) uint {
	var card db.Card
	db.DB.Where("number = ? AND suit = ?", number, suit).Find(&card)
	return card.ID
}
