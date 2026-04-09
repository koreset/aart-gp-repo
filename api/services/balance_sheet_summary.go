package services

func GetBalanceSheetDates() ([]string, error) {
	var dates []string
	err := DB.Table("balance_sheet_records").Select("DISTINCT date").Pluck("date", &dates)
	if err.Error != nil {
		return nil, err.Error
	}
	return dates, nil
}
