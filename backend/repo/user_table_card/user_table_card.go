package user_table_card

import "backend/db"

func AddNewCard(tableID, userID, cardID uint) {
	userCard := db.UsersTablesCard{
		TableID: tableID,
		UserID:  userID,
		CardID:  cardID,
	}
	db.DB.Create(&userCard)
}

func FindCards(tableID, userID int) []db.UsersTablesCard {
	var cards []db.UsersTablesCard
	db.DB.Preload("Card").Where("table_id = ? AND user_id = ?", tableID, userID).Find(&cards)
	return cards
}
