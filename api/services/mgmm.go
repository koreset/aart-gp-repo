package services

import "api/models"

func GetAvailablePAAParameterYears(portfolioName string) ([]int, error) {
	var list []int
	err := DB.Model(&models.ModifiedGMMParameter{}).Where("portfolio_name = ?", portfolioName).Distinct("year").Pluck("year", &list).Error
	return list, err
}

func GetAvailablePAAPremiumPatternYears(portfolioName string) ([]int, error) {
	//First get the unique premium earning pattern code for the given portfolio
	var param models.ModifiedGMMParameter
	var list []int

	err := DB.Model(&models.ModifiedGMMParameter{}).Where("portfolio_name = ?", portfolioName).First(&param).Error
	if err != nil {
		return list, err
	}

	err = DB.Model(&models.PremiumEarningPattern{}).Where("premium_earning_pattern_code = ?", param.PremiumEarningPatternCode).Distinct("year").Pluck("year", &list).Error
	return list, err
}

func CheckExistingGMMRunName(name string) (bool, error) {
	var count int64
	err := DB.Model(&models.MgmmRun{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

func CheckExistingPortfolioName(name string) (bool, error) {
	var count int64
	err := DB.Model(&models.PaaPortfolio{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}
