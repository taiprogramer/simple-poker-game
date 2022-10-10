package table

import (
	"github.com/taiprogramer/simple-poker-game/backend/db"
)

func CreateNewTable(userID uint, roomID uint, round int, pot int, done bool) uint {
	table := db.Table{
		Round:  round,
		Pot:    pot,
		Done:   done,
		UserID: userID,
		RoomID: roomID,
	}

	_ = db.DB.Create(&table).RowsAffected == 1
	return table.ID
}

func FindTableByUserIDAndRoomID(userID uint, roomID uint) db.Table {
	var table db.Table
	db.DB.Where("user_id = ? AND room_id = ?", userID, roomID).First(&table)
	return table
}
