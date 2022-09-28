package user_table_card

import "github.com/taiprogramer/simple-poker-game/backend/db"

func AddNewCard(tableID, userID, cardID uint) {
	userCard := db.UsersTablesCard{
		TableID: tableID,
		UserID:  userID,
		CardID:  cardID,
	}
	db.DB.Create(&userCard)
}
