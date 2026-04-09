package services

import "api/models"

func GetPAAPortfolioNames() ([]string, error) {
	var portfolioNames []string
	err := DB.Table("modified_gmm_scoped_aggregations").Distinct("portfolio_name").Select("portfolio_name").Find(&portfolioNames).Error
	return portfolioNames, err
}

func GetPAARuns() ([]models.MgmmRun, error) {
	var runs []models.MgmmRun
	//err := DB.Debug().Table("modified_gmm_scoped_aggregations").Distinct("run_id").Select("run_id, run_name, portfolio_name").Find(&runs).Error
	//err := DB.Table("modified_gmm_scoped_aggregations").Distinct("run_id").Select("run_id, run_name, portfolio_name").Find(&runs).Error
	err := DB.Where("ifrs17_ready = ? ", true).Order("id desc").Find(&runs).Error
	return runs, err
}
