package bet_histories

import "github.com/taiprogramer/simple-poker-game/backend/db"

func WriteBetHistory(tableID, userID, actionID, amount int) {
	betHistory := db.BetHistory{
		TableID:  uint(tableID),
		UserID:   uint(userID),
		ActionID: uint(actionID),
		Amount:   amount,
	}
	db.DB.Save(&betHistory)
}

func FindActionIDByName(name string) int {
	var action db.Action
	db.DB.Where("name = ?", name).First(&action)
	return int(action.ID)
}

func GetLatest(userID int) db.BetHistory {
	var betHistory db.BetHistory
	db.DB.Where("user_id = ?", userID).Last(&betHistory)
	return betHistory
}
