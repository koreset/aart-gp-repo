package services

import "api/models"

func GetCSMRuns() ([]models.CsmRun, error) {
	var runs []models.CsmRun
	err := DB.Order("id desc").Find(&runs).Error
	return runs, err

}

func DeleteCsmRun(id int, measurementType, runDate string) error {
	var err error
	if measurementType == "GMM" || measurementType == "VFA" {
		err = DB.Where("csm_run_id = ? and run_date = ?", id, runDate).Delete(&models.CsmProjection{}).Error
		if err != nil {
			return err
		}

		DB.Where("csm_run_id = ?", id).Delete(&models.CsmAosVariable{})
		DB.Where("csm_run_id = ? and run_date = ?", id, runDate).Delete(&models.JournalTransactions{})
		DB.Where("csm_run_id = ? and run_date = ?", id, runDate).Delete(&models.AOSStepResult{})
		DB.Where("csm_run_id = ? and run_date = ?", id, runDate).Delete(&models.LiabilityMovement{})
		DB.Where("csm_run_id = ? and run_date = ?", id, runDate).Delete(&models.InitialRecognition{})
		DB.Where("csm_run_id = ? and run_date = ?", id, runDate).Delete(&models.InsuranceRevenue{})
		DB.Where("csm_run_id = ? and run_date = ?", id, runDate).Delete(&models.CsmProjection{})
		err = DB.Where("csm_run_id = ? and date= ? and measurement_type=?", id, runDate, "GMM").Delete(&models.BalanceSheetRecord{}).Error
		err = DB.Where("csm_run_id = ? and date= ? and measurement_type=?", id, runDate, "VFA").Delete(&models.BalanceSheetRecord{}).Error
		err = DB.Where("id = ? and run_date =?", id, runDate).Delete(&models.CsmRun{}).Error

	}

	if measurementType == "PAA" {
		var csmRun models.CsmRun
		err = DB.Where("id = ? and run_date = ?", id, runDate).First(&csmRun).Error
		if err != nil {
			return err
		}
		DB.Where("csm_run_id = ? and run_date = ?", csmRun.PaaRunId, runDate).Delete(&models.PAAResult{})
		DB.Where("run_id = ? and run_date = ?", csmRun.PaaRunId, runDate).Delete(&models.PAABuildUp{})
		DB.Where("csm_run_id = ? and measurement_type = ? and run_date =?", csmRun.PaaRunId, "PAA", runDate).Delete(&models.BalanceSheetRecord{})
		DB.Where("csm_run_id = ? and run_date = ?", id, runDate).Delete(&models.JournalTransactions{})
		DB.Where("csm_run_id = ? and run_date = ?", csmRun.PaaRunId, runDate).Delete(&models.PAAEligibilityTestResult{})
		err = DB.Where("id = ? and run_date =?", id, runDate).Delete(&models.CsmRun{}).Error
	}
	return err
}
