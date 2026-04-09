package services

import (
	"api/models"
	"fmt"
)

func GetOlympicWinners(request map[string]interface{}) (map[string]interface{}, error) {
	var results map[string]interface{}
	fmt.Println("GetOlympicWinners: ", request)
	var winners []models.OlympicWinner
	DB.Table("olympic_winners").Limit(100).Find(&winners)
	results = make(map[string]interface{})
	results["rows"] = winners
	results["total"] = 8000

	return results, nil
}
