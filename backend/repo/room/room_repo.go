package room

import "github.com/taiprogramer/simple-poker-game/backend/db"

func FindRoomByID(id int) (*db.Room, bool) {
	var room db.Room
	result := db.DB.First(&room, id)
	if result.RowsAffected == 0 {
		return nil, false

	}
	return &room, true
}

func FindWaitingListsByRoomID(id int) (*[]db.WaitingList, bool) {
	var waitingLists []db.WaitingList
	db.DB.Where("room_id", id).Find(&waitingLists)
	return &waitingLists, true
}

func UpdateRoom(room *db.Room) {
	db.DB.Save(room)
}

func DeleteRoom(room *db.Room) {
	db.DB.Where("id = ?", room.ID).Delete(room)
}

func DeleteWaitingListByUserID(userID int, waiting *db.WaitingList) {
	db.DB.Where("user_id = ?", userID).Delete(waiting)
}

func DeleteWaitingListsByRoomID(roomID int) {
	db.DB.Where("room_id = ?", roomID).Delete(&db.WaitingList{})
}

func GetWaitingListsByRoomID(roomID int) []db.WaitingList {
	var waitingLists []db.WaitingList
	db.DB.Where("room_id = ?", roomID).Find(&waitingLists)
	return waitingLists
}

func IncreaseUserAmount(userID, roomID, amount int) {
	var waitingList db.WaitingList
	db.DB.Where("room_id = ? AND user_id = ?", roomID, userID).First(&waitingList)
	waitingList.AvailableMoney = waitingList.AvailableMoney + amount
	db.DB.Save(&waitingList)
}
