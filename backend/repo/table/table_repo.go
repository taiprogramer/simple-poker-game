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

func UpdateTable(table *db.Table) {
	db.DB.Save(table)
}

func FindTableByRoomID(roomID uint) db.Table {
	var table db.Table
	db.DB.Preload("Cards").Where("room_id = ?", roomID).First(&table)
	return table
}

func GetTableByID(tableID int) db.Table {
	var table db.Table
	db.DB.Preload("Cards").Where("id = ?", tableID).First(&table)
	return table
}
