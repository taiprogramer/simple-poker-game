package bet_histories

import "github.com/taiprogramer/simple-poker-game/backend/db"

func WriteBetHistory(tableID, userID, actionID, amount int, round int) {
	betHistory := db.BetHistory{
		TableID:  uint(tableID),
		UserID:   uint(userID),
		ActionID: uint(actionID),
		Amount:   amount,
		Round:    round,
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

func GetLatestByTableID(tableID int) db.BetHistory {
	var betHistory db.BetHistory
	db.DB.Where("table_id = ?", tableID).Last(&betHistory)
	return betHistory
}

func GetTotalAmountByRoundAndUserID(tableID int, round int, userID int) int {
	type Result struct {
		Total int
	}
	var result Result
	db.DB.Model(&db.BetHistory{}).
		Select("sum(amount) as total").
		Where("user_id = ? AND round = ? AND table_id = ?", userID, round, tableID).Find(&result)
	return result.Total
}
