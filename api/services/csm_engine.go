package services

import (
	"api/log"
	"api/models"
	"api/utils"
	"context"
	"errors"
	"fmt"
	"math"
	"mime/multipart"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type Result struct {
	Spcode      int    `json:"spcode"`
	Ifrs17Group string `json:"ifrs_17_group"`
}

func GetSpCodesForProduct(productCode string) []Result {
	var result []Result
	//var spcode []int
	tableName := strings.ToLower(productCode) + "_modelpoints"
	ctx := context.Background()
	err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
		return d.Table(tableName).Select("spcode as spcode").Group("spcode").Find(&result).Error
	})
	if err != nil {
		fmt.Println(err)
	}

	return result
}

func GetIfrs17GroupsForProduct(productCode string) []Result {
	var result []Result
	tableName := strings.ToLower(productCode) + "_modelpoints"
	ctx := context.Background()
	err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
		return d.Table(tableName).Where("ifrs17_group is not null").Select("ifrs17_group as ifrs17_group").Group("ifrs17_group").Find(&result).Error
	})
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func SaveRiskAdjustmentConfigs(configs []models.RAConfiguration) error {
	for _, config := range configs {
		err := DB.Create(&config).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func GetRiskAdjustmentConfigs() ([]models.RAConfiguration, error) {
	var configs []models.RAConfiguration
	ctx := context.Background()
	err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
		return d.Find(&configs).Error
	})
	if err != nil {
		return nil, err
	}

	return configs, nil
}

func SaveRiskDrivers(drivers []models.RiskDriver) error {
	for _, driver := range drivers {
		err := DB.Save(&driver).Error
		if err != nil {
			return err
		}
		LogIFRS17Event("risk_driver_saved", "risk_driver", driver.ProductCode, 0, models.AppUser{}, "")
	}
	return nil
}

func DeleteRiskDriver(productCode string) error {
	LogIFRS17Event("risk_driver_deleted", "risk_driver", productCode, 0, models.AppUser{}, "")
	err := DB.Where("product_code = ?", productCode).Delete(&models.RiskDriver{}).Error
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func GetRiskDrivers() ([]models.RiskDriver, error) {
	var drivers []models.RiskDriver
	ctx := context.Background()
	err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
		return d.Find(&drivers).Error
	})
	if err != nil {
		return nil, err
	}
	return drivers, nil
}

// GetRiskAdjustmentDriverTypes returns distinct risk_type values from the risk_adjustment_drivers table
func GetRiskAdjustmentDriverTypes() ([]string, error) {
	var types []string
	ctx := context.Background()
	err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
		return d.Table("risk_adjustment_drivers").Distinct("risk_type").Pluck("risk_type", &types).Error
	})
	if err != nil {
		return nil, err
	}
	return types, nil
}

func DeleteAosConfiguration(id int) error {
	LogIFRS17Event("config_deleted", "aos_config", "", id, models.AppUser{}, "")
	err := DB.Where("aos_variable_set_id = ?", id).Delete(&models.AosVariable{}).Error
	if err != nil {
		fmt.Println(err)
	}
	err = DB.Where("id = ?", id).Delete(&models.AosVariableSet{}).Error
	if err != nil {
		fmt.Println(err)
	}

	return err
}

func SaveAosConfigurations(configs models.AosVariableSet) (models.AosVariableSet, error) {
	//Strip out the ids from the AosVariableSet
	if !configs.ExternalSap {
		for i := range configs.AosVariables {
			configs.AosVariables[i].ID = 0
			configs.AosVariables[i].AosVariableSetName = configs.ConfigurationName
			if configs.AosVariables[i].RunName != "" {
				var runName models.ProjectionJob
				ctx := context.Background()
				err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
					return d.Where("run_name = ?", configs.AosVariables[i].RunName).Find(&runName).Error
				})
				if err != nil {
					fmt.Println(err)
				}
				configs.AosVariables[i].RunId = runName.ID
			} else {
				configs.AosVariables[i].RunId = 0
			}
		}
	} else {
		for i := range configs.AosVariables {
			configs.AosVariables[i].ID = 0
			configs.AosVariables[i].AosVariableSetName = configs.ConfigurationName
			configs.AosVariables[i].ExternalSapSource = configs.ExternalSap
		}
	}

	err := DB.Create(&configs).Error
	if err != nil {
		return configs, err
	}

	LogIFRS17Event("config_saved", "aos_config", configs.ConfigurationName, configs.ID, models.AppUser{}, "")
	return configs, nil
}

func GetAosConfigurations() ([]models.AOSConfiguration, error) {
	var configs []models.AOSConfiguration
	ctx := context.Background()
	err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
		return d.Find(&configs).Error
	})
	if err != nil {
		return nil, err
	}
	return configs, nil
}

func CreateCsmRun(csmRun *models.CsmRun) error {
	err := DB.Create(csmRun).Error
	if err == nil {
		LogIFRS17Event("run_created", "csm_run", csmRun.Name, csmRun.ID, models.AppUser{UserName: csmRun.UserName}, "")
	}
	return err
}

func CheckExistingRun(run models.CsmRun) bool {
	var csmRun models.CsmRun
	ctx := context.Background()
	if run.MeasurementType == "PAA" {
		err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
			return d.Where("run_date = ? and measurement_type= ? and paa_run_id = ?", run.RunDate, run.MeasurementType, run.PaaRunId).Find(&csmRun).Error
		})
		if err != nil {
			fmt.Println(err)
		}
		if csmRun.ID > 0 {
			return true
		} else {
			return false
		}
	}

	if run.MeasurementType == "GMM" || run.MeasurementType == "VFA" {
		err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
			return d.Where("run_date = ? and measurement_type= ? and configuration_name = ?", run.RunDate, run.MeasurementType, run.ConfigurationName).Find(&csmRun).Error
		})
		if err != nil {
			fmt.Println(err)
		}
		if csmRun.ID > 0 {
			return true
		} else {
			return false
		}
	}

	if run.MeasurementType == "VFA" {
		err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
			return d.Where("run_date = ? and measurement_type= ? and configuration_name = ?", run.RunDate, run.MeasurementType, run.ConfigurationName).Find(&csmRun).Error
		})
		if err != nil {
			fmt.Println(err)
		}
		if csmRun.ID > 0 {
			return true
		} else {
			return false
		}
	}

	return false
}

func CheckExistingCsmRunName(name string) bool {
	var csmRunName models.CsmRun
	ctx := context.Background()
	err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
		return d.Where("name = ?", name).Find(&csmRunName).Error
	})
	if err != nil {
		fmt.Println(err)
	}
	if csmRunName.ID > 0 {
		return true
	} else {
		return false
	}
}

func RunCsmEngineForPAA(csmRun models.CsmRun) map[string]interface{} {
	start := time.Now()
	csmRun.ProcessingStatus = "running"
	DB.Save(&csmRun)
	q := strings.Builder{}
	var err error

	q.WriteString("SELECT distinct ifrs17_group, product_code FROM modified_gmm_scoped_aggregations where run_id=" + strconv.Itoa(csmRun.PaaRunId))

	//var productCode = "FUN_01"
	var groups []models.GroupResults

	err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw(q.String()).Scan(&groups).Error
	})
	if err != nil {
		fmt.Println(err)
	}

	csmRun.TotalGroups = len(groups)
	DB.Save(&csmRun)

	for _, group := range groups {
		CalculatePAACsm(&csmRun, group)
		if csmRun.ProcessingStatus == "failed" {
			//We delete all calcs for this run
			DB.Where("csm_run_id = ?", csmRun.ID).Delete(&models.PAAResult{})
			DB.Where("csm_run_id = ? and measurement_type = ?", csmRun.ID, "PAA").Delete(&models.JournalTransactions{})
			DB.Where("csm_run_id = ?", csmRun.ID).Delete(&models.BalanceSheetRecord{})
			DB.Where("csm_run_id = ?", csmRun.ID).Delete(&models.PAAEligibilityTestResult{})
			break
		}
		csmRun.ProcessedGroups = csmRun.ProcessedGroups + 1
		DB.Save(&csmRun)
	}

	var res = make(map[string]interface{})
	end := time.Since(start)
	csmRun.RunTime = end.Seconds()
	if csmRun.ProcessingStatus != "failed" {
		csmRun.ProcessingStatus = "completed"
	}
	DB.Save(&csmRun)
	res["run"] = csmRun
	return res
}

func GetCsmRunById(runId int) (models.CsmRun, error) {
	var csmRun models.CsmRun

	ctx := context.Background()
	err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
		return d.Where("id = ?", runId).Find(&csmRun).Error
	})

	if err != nil {
		fmt.Println(err)
		return csmRun, err
	}
	return csmRun, nil
}

func DeleteExistingCsmRun(csmRun models.CsmRun) {
	//We will always override so delete runs that have the same parameters
	if csmRun.MeasurementType == "GMM" || csmRun.MeasurementType == "VFA" {
		err := DB.Where("run_date = ? and (measurement_type = ? or measurement_type = ?)  and configuration_name = ?", csmRun.RunDate, "GMM", "VFA", csmRun.ConfigurationName).Delete(&models.CsmRun{}).Error
		if err != nil {
			fmt.Println(err)
		}
	}

	if csmRun.MeasurementType == "PAA" {
		err := DB.Where("run_date = ? and paa_run_id = ?", csmRun.RunDate, csmRun.PaaRunId).Delete(&models.CsmRun{}).Error
		if err != nil {
			fmt.Println(err)
		}
	}

	LogIFRS17Event("run_deleted", "csm_run", "", 0, models.AppUser{}, fmt.Sprintf("date=%s type=%s", csmRun.RunDate, csmRun.MeasurementType))
}

func RerunCsmEngine(csmRun models.CsmRun) map[string]interface{} {
	// delete all calculations for the existing run. We do not want to delete the existing run as we want to keep the audit trail
	DB.Where("csm_run_id = ?", csmRun.ID).Delete(&models.CsmAosVariable{})
	DB.Where("csm_run_id = ? and run_date = ?", csmRun.ID, csmRun.RunDate).Delete(&models.JournalTransactions{})
	DB.Where("csm_run_id = ? and run_date = ?", csmRun.ID, csmRun.RunDate).Delete(&models.AOSStepResult{})
	DB.Where("csm_run_id = ? and run_date = ?", csmRun.ID, csmRun.RunDate).Delete(&models.LiabilityMovement{})
	DB.Where("csm_run_id = ? and run_date = ?", csmRun.ID, csmRun.RunDate).Delete(&models.InitialRecognition{})
	DB.Where("csm_run_id = ? and run_date = ?", csmRun.ID, csmRun.RunDate).Delete(&models.InsuranceRevenue{})
	DB.Where("csm_run_id = ? and run_date = ?", csmRun.ID, csmRun.RunDate).Delete(&models.CsmProjection{})
	DB.Where("csm_run_id = ? and date= ? and measurement_type=?", csmRun.ID, csmRun.RunDate, "GMM").Delete(&models.BalanceSheetRecord{})
	DB.Where("csm_run_id = ? and date= ? and measurement_type=?", csmRun.ID, csmRun.RunDate, "VFA").Delete(&models.BalanceSheetRecord{})
	DB.Where("id = ? and run_date =?", csmRun.ID, csmRun.RunDate).Delete(&models.CsmRun{})
	results := RunCsmEngine(csmRun, true)
	return results
}

func RunCsmEngine(csmRun models.CsmRun, rerun bool) map[string]interface{} {
	//delete existing run with same parameters
	if !rerun {
		DeleteExistingCsmRun(csmRun)
	}

	start := time.Now()
	csmRun.ProcessingStatus = "running"
	csmRun.CreationDate = time.Now()
	csmRun.FailureReason = ""
	DB.Save(&csmRun)
	var configRunIds []models.ConfigRunId
	var aosVars []models.AosVariable
	var err error
	var aosVarSet models.AosVariableSet
	err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("configuration_name = ?", csmRun.ConfigurationName).Find(&aosVarSet).Error
	})
	if err != nil {
		fmt.Println(err)
	}
	err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("aos_variable_set_id=?", aosVarSet.ID).Find(&aosVars).Error
	})
	if err != nil {
		fmt.Println(err)
	}
	if aosVarSet.ExternalSap {
		csmRun.ManualSap = true
	}

	err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw(fmt.Sprintf("SELECT distinct(run_name), run_id FROM aos_variables where aos_variable_set_id = %d;", aosVarSet.ID)).Scan(&configRunIds).Error
	})
	if err != nil {
		fmt.Println(err)
	}

	q := strings.Builder{}

	if !aosVarSet.ExternalSap {
		q.WriteString("SELECT distinct ifrs17_group, product_code FROM scoped_aggregated_projections where job_product_id in (")
		for i, item := range configRunIds {
			var jp []models.JobProduct

			//Note to self: This will throw a problem when running several products
			// in one job. Will need to refactor at later stage.
			err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
				return d.Where("projection_job_id = ?", item.RunId).Find(&jp).Error
			})
			if err != nil {
				fmt.Println(err)
			}

			if i < len(configRunIds)-1 {
				for _, j := range jp {
					q.WriteString(strconv.Itoa(j.ID) + ",")
				}
			} else {
				if len(jp) > 0 {
					for k, j := range jp {
						if k < len(jp)-1 {
							q.WriteString(strconv.Itoa(j.ID) + ",")
						} else {
							q.WriteString(strconv.Itoa(j.ID))
						}
					}
				}
				q.WriteString(")")
			}
		}
	} else {
		q.WriteString("SELECT distinct ifrs17_group, product_code FROM manual_scoped_aggregated_projections where run_name in (")
		for i, item := range configRunIds {
			if i < len(configRunIds)-1 {
				q.WriteString("'" + item.RunName + "',")
			} else {
				q.WriteString("'" + item.RunName + "'")
			}
		}
		q.WriteString(")")
	}

	fmt.Println(q.String())
	queryString := q.String()
	//var productCode = "FUN_01"
	var groups []models.GroupResults
	queryString = strings.Replace(queryString, ",)", ")", -1)
	err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw(queryString).Scan(&groups).Error
	})
	if err != nil {
		fmt.Println(err)
	}

	csmRun.TotalGroups = len(groups)
	csmRun.ProcessedGroups = 0
	DB.Save(&csmRun)

	for _, group := range groups {
		stepResults, err := CalculateCsm(csmRun, group)
		if err != nil {
			fmt.Println(err)
			csmRun.ProcessingStatus = "failed"
			csmRun.FailureReason = err.Error()
			csmRun.RunTime = time.Since(start).Seconds()
			DB.Save(&csmRun)
			DB.Where("csm_run_id = ?", csmRun.ID).Delete(&models.JournalTransactions{})
			DB.Where("csm_run_id = ?", csmRun.ID).Delete(&models.AOSStepResult{})
			DB.Where("csm_run_id = ?", csmRun.ID).Delete(&models.CsmProjection{})
			return nil
			//break
		}
		ProcessLiabilityMovements(stepResults, csmRun)
		ProcessInitialRecognitionAnalysis(stepResults)
		csmRun.ProcessedGroups++
		DB.Save(&csmRun)
	}

	var res = make(map[string]interface{})
	var prodList []models.ProductList
	err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw("select distinct product_code from aos_step_results").Scan(&prodList).Error
	})
	if err != nil {
		fmt.Println(err)
	}
	res["products"] = prodList
	res["run"] = csmRun
	end := time.Since(start)
	csmRun.RunTime = end.Seconds()
	csmRun.ProcessingStatus = "completed"
	DB.Save(&csmRun)
	return res
}

func ProcessInitialRecognitionAnalysis(stepResults []models.AOSStepResult) {
	//TODO:  build processing function
}

func GetLiabilityMovementRunList() []models.CsmRun {
	var results []models.CsmRun
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Table("liability_movements").Distinct("run_date").Find(&results).Error }); err != nil {
		fmt.Println(err)
	}
	return results
}

func ProcessLiabilityMovements(results []models.AOSStepResult, csmRun models.CsmRun) {
	// pull in all the liability_movement_lines
	var liabilityMovementLines []models.LiabilityMovementLine
	var liabilityMovements []models.LiabilityMovement
	var insuranceRevenues []models.InsuranceRevenue
	var initialRecognition []models.InitialRecognition
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Find(&liabilityMovementLines).Error
	})
	if err != nil {
		fmt.Println(err)
	}

	for _, lml := range liabilityMovementLines {
		var lm models.LiabilityMovement
		lm.Code = lml.Code
		lm.Variable = lml.Variable
		lm.CsmRunID = results[0].CsmRunID
		lm.RunDate = results[0].RunDate
		if lml.Code == 100 { //Insurance Contract Assets B/F
			lm.ProductCode = results[0].ProductCode
			lm.Ifrs17Group = results[0].IFRS17Group
			lm.Bel = results[0].BelInflow
			lm.RiskAdjustment = 0
			lm.Csm = 0

			lm.TotalLRC = lm.Bel + lm.RiskAdjustment + lm.Csm
			//if results[0].BEL < 0 {
			//	lm.Bel = results[0].BEL
			//} else {
			//	lm.Bel = 0
			//}

		}

		if lml.Code == 101 { //Insurance Contract Liabilities B/F
			lm.ProductCode = results[0].ProductCode
			lm.Ifrs17Group = results[0].IFRS17Group
			//if results[0].LossComponentChange > 0 {
			//	if (results[0].BelOutflow + results[0].RiskAdjustment) > 0 {
			//		lm.Bel = results[0].BelOutflow - results[0].LossComponentChange*results[0].BelOutflow/(results[0].BelOutflow+results[0].RiskAdjustment)
			//		lm.RiskAdjustment = results[0].RiskAdjustment - results[0].LossComponentChange*results[0].RiskAdjustment/(results[0].BelOutflow+results[0].RiskAdjustment)
			//		lm.LossComponent = results[0].LossComponentChange * results[0].BelOutflow / (results[0].BelOutflow + results[0].RiskAdjustment)
			//		lm.LossComponentRiskAdj = results[0].LossComponentChange * results[0].RiskAdjustment / (results[0].BelOutflow + results[0].RiskAdjustment)
			//	} else {
			//		lm.Bel = 0
			//		lm.RiskAdjustment = 0
			//		lm.LossComponent = results[0].LossComponentBuildup
			//	}
			//} else {
			lm.Bel = -results[0].BelOutflow
			lm.RiskAdjustment = -results[0].RiskAdjustment
			lm.LossComponent = -results[0].LossComponentBuildup
			//}

			lm.Csm = -results[0].CSMBuildup
			lm.TotalLRC = lm.Bel + lm.RiskAdjustment + lm.Csm
			//if results[0].BEL > 0 {

			//if results[0].BEL > 0 {
			//	lm.Bel = results[0].BEL
			//} else {
			//	lm.Bel = 0
			//}
			lm.IncurredClaimsBel = 0
			lm.IncurredClaimsRiskAdj = 0

			lm.FRACsm = 0
			lm.MRACsm = 0
			lm.FVACsm = 0
		}

		//CSM Release
		if lml.Code == 201 {
			lm.ProductCode = results[0].ProductCode
			lm.Ifrs17Group = results[0].IFRS17Group
			lm.Bel = 0
			lm.RiskAdjustment = 0
			lm.IncurredClaimsBel = 0
			lm.IncurredClaimsRiskAdj = 0
			lm.Csm = -results[16].CSMRelease
			//if results[15].BelOutflow+results[15].RiskAdjustment > 0 {
			//	lm.LossComponent = 0        //-results[16].LossComponentChange * results[15].BelOutflow / (results[15].BelOutflow + results[15].RiskAdjustment)
			//	lm.LossComponentRiskAdj = 0 //-results[16].LossComponentChange * (1 - results[15].BelOutflow/(results[15].BelOutflow+results[15].RiskAdjustment))
			//} else {
			//	lm.LossComponent = 0
			//	lm.LossComponentRiskAdj = 0
			//}
			lm.FRACsm = 0
			lm.MRACsm = 0
			lm.FVACsm = 0
			lm.TotalLRC = lm.Bel + lm.RiskAdjustment + lm.Csm

			var ir models.InsuranceRevenue
			ir.ProductCode = results[0].ProductCode
			ir.IFRS17Group = results[0].IFRS17Group
			ir.Code = lm.Code
			ir.RunDate = results[0].RunDate
			ir.CsmRunID = csmRun.ID
			ir.Variable = "CSM Release"
			ir.PostTransition = -lm.Csm
			ir.FullRetrospectiveApproach = -lm.FRACsm
			ir.ModifiedRetrospectiveApproach = -lm.MRACsm
			ir.FairValueApproach = -lm.FVACsm
			insuranceRevenues = append(insuranceRevenues, ir)

		}

		//Risk Adjustment Release
		if lml.Code == 202 {
			lm.ProductCode = results[0].ProductCode
			lm.Ifrs17Group = results[0].IFRS17Group
			lm.Bel = 0
			//if results[3].LossComponentBuildup > 0 {
			//	if (results[3].BEL + results[3].RiskAdjustment) > 0 {
			//		lm.RiskAdjustment = 0 //results[3].RiskAdjustmentChange * (1.0 - results[3].LossComponentBuildup/(results[3].BEL+results[3].RiskAdjustment))
			//		lm.LossComponentRiskAdj = results[3].RiskAdjustmentChange * results[3].LossComponentBuildup / (results[3].BEL + results[3].RiskAdjustment)
			//	} else {
			//		lm.RiskAdjustment = 0 //results[3].RiskAdjustmentChange
			//		lm.LossComponentRiskAdj = 0
			//	}
			//} else {
			//	lm.RiskAdjustment = results[3].RiskAdjustmentChange
			//	lm.LossComponentRiskAdj = 0
			//}

			lm.RiskAdjustment = results[3].RiskAdjustmentChange

			lm.IncurredClaimsBel = 0
			lm.IncurredClaimsRiskAdj = 0
			lm.Csm = 0
			lm.LossComponent = 0
			lm.FRACsm = 0
			lm.MRACsm = 0
			lm.FVACsm = 0
			lm.TotalLRC = lm.Bel + lm.RiskAdjustment + lm.Csm

			var ir models.InsuranceRevenue
			ir.ProductCode = results[0].ProductCode
			ir.IFRS17Group = results[0].IFRS17Group
			ir.Code = lm.Code
			ir.RunDate = results[0].RunDate
			ir.CsmRunID = csmRun.ID
			ir.Variable = "Risk Adjustment Release"
			ir.PostTransition = -lm.RiskAdjustment
			ir.FullRetrospectiveApproach = -lm.FRACsm
			ir.ModifiedRetrospectiveApproach = -lm.MRACsm
			ir.FairValueApproach = -lm.FVACsm
			insuranceRevenues = append(insuranceRevenues, ir)
		}

		//Release of FCF and Experience Adjustment
		//Release of Expected Cash flows
		//+ (BEL1atT1 - BEL0atT1) data change
		// Expected - Actual
		if lml.Code == 203 {
			lm.ProductCode = results[0].ProductCode
			lm.Ifrs17Group = results[0].IFRS17Group
			var financeVariables = getFinancialVariables(lm.ProductCode, lm.Ifrs17Group, csmRun.FinanceYear, csmRun.FinanceVersion)
			//
			//if results[8].LossComponentBuildup > 0 {
			//	if results[8].BelOutflow+results[8].RiskAdjustment > 0 {
			//		lm.Bel = results[9].BestEstimateLiabilityChange - results[9].LossComponentChange
			//		lm.RiskAdjustment = results[9].RiskAdjustmentChange * (1.0 - results[8].LossComponentBuildup/(results[8].BEL+results[8].RiskAdjustment))
			//		lm.LossComponentRiskAdj = results[9].RiskAdjustmentChange * results[8].LossComponentBuildup / (results[8].BEL + results[8].RiskAdjustment)
			//		lm.LossComponent = -(results[4].PNLChange+results[5].PNLChange+results[6].PNLChange+results[7].PNLChange)*results[4].LossComponentBuildup/(results[4].BEL+results[4].RiskAdjustment) + results[8].LossComponentChange + results[9].LossComponentChange // loss component unwind plus interest plus data change
			//	} else {
			//		lm.Bel = 0
			//		lm.RiskAdjustment = 0
			//		lm.LossComponentRiskAdj = 0
			//		lm.LossComponent = 0
			//	}
			//} else {
			//	lm.Bel = -(results[4].PNLChange + results[5].PNLChange + results[6].PNLChange + results[7].PNLChange) - results[9].BestEstimateLiabilityChange
			//	lm.RiskAdjustment = results[9].RiskAdjustmentChange
			//	lm.LossComponentRiskAdj = 0
			//	lm.LossComponent = 0
			//
			//}

			lm.Bel = -(results[4].PNLChange + results[5].PNLChange + results[6].PNLChange + results[7].PNLChange) + results[8].ExpectedCashInflow +
				financeVariables.ActualMortalityClaimsIncurred + financeVariables.ActualMorbidityClaimsIncurred + financeVariables.ActualRetrenchmentClaimsIncurred +
				financeVariables.ActualAttributableExpenses - financeVariables.ActualPremiumIncome
			lm.RiskAdjustment = 0 //results[9].RiskAdjustmentChange
			lm.Csm = 0            //results[8].CSMChange + results[9].CSMChange //interest plus data change
			lm.LossComponent = results[8].LossComponentUnwind
			lm.IncurredClaimsBel = 0
			lm.IncurredClaimsRiskAdj = 0
			lm.FRACsm = 0
			lm.MRACsm = 0
			lm.FVACsm = 0
			lm.TotalLRC = lm.Bel + lm.RiskAdjustment + lm.Csm

			//
			var ir models.InsuranceRevenue
			ir.ProductCode = results[0].ProductCode
			ir.IFRS17Group = results[0].IFRS17Group
			ir.Code = lm.Code
			ir.RunDate = results[0].RunDate
			ir.CsmRunID = csmRun.ID
			ir.Variable = "Release of FCF and Experience Adjustment"
			ir.PostTransition = results[4].PNLChange + results[5].PNLChange + results[6].PNLChange + results[7].PNLChange +
				financeVariables.ActualPremiumIncome - results[8].ExpectedCashInflow //-(lm.Bel + lm.RiskAdjustment + lm.IncurredClaimsBel + lm.IncurredClaimsRiskAdjustment)
			ir.FullRetrospectiveApproach = -lm.FRACsm
			ir.ModifiedRetrospectiveApproach = -lm.MRACsm
			ir.FairValueApproach = -lm.FVACsm
			insuranceRevenues = append(insuranceRevenues, ir)
		}

		//Other
		if lml.Code == 204 {
			lm.ProductCode = results[0].ProductCode
			lm.Ifrs17Group = results[0].IFRS17Group
			//if results[2].LossComponentChange > 0 {
			//	lm.LossComponent = results[8].ExpectedCashInflow
			//} else {
			//	lm.Bel = results[8].ExpectedCashInflow
			//}

			lm.Bel = 0
			lm.RiskAdjustment = 0
			lm.IncurredClaimsBel = 0
			lm.IncurredClaimsRiskAdj = 0
			lm.Csm = 0
			lm.FRACsm = 0
			lm.MRACsm = 0
			lm.FVACsm = 0
			lm.TotalLRC = lm.Bel + lm.RiskAdjustment + lm.Csm

			var ir models.InsuranceRevenue
			ir.ProductCode = results[0].ProductCode
			ir.IFRS17Group = results[0].IFRS17Group
			ir.Code = lm.Code
			ir.RunDate = results[0].RunDate
			ir.CsmRunID = csmRun.ID
			ir.Variable = "Other"
			ir.PostTransition = 0
			ir.FullRetrospectiveApproach = 0
			ir.ModifiedRetrospectiveApproach = 0
			ir.FairValueApproach = 0
			insuranceRevenues = append(insuranceRevenues, ir)

		}

		//Changes in estimates that affect CSM
		if lml.Code == 301 {
			var total_LC_Change float64
			var total_CSM_Change float64
			lm.ProductCode = results[0].ProductCode
			lm.Ifrs17Group = results[0].IFRS17Group
			total_CSM_Change = results[9].CSMChange + results[10].CSMChange + results[11].CSMChange + results[12].CSMChange + results[13].CSMChange
			total_LC_Change = results[9].LossComponentChange + results[10].LossComponentChange + results[11].LossComponentChange + results[12].LossComponentChange + results[13].LossComponentChange
			//if total_CSM_Change+total_LC_Change != 0 {
			//	lm.LossComponent = 0
			//	lm.Bel = (results[13].BEL - results[9].BEL) * (total_CSM_Change / (total_LC_Change + total_CSM_Change))
			//	lm.LossComponentRiskAdj = 0
			//	lm.RiskAdjustment = (results[10].RiskAdjustmentChange + results[11].RiskAdjustmentChange + results[12].RiskAdjustmentChange + results[13].RiskAdjustmentChange) * (total_CSM_Change / (total_LC_Change + total_CSM_Change))
			//	lm.Csm = total_CSM_Change
			//} else {
			//	lm.LossComponent = 0
			//	lm.Bel = 0
			//	lm.LossComponentRiskAdj = 0
			//	lm.RiskAdjustment = 0
			//	lm.Csm = 0
			//}

			if total_CSM_Change+total_LC_Change != 0 {
				lm.Bel = (results[9].BestEstimateLiabilityChange + results[10].BestEstimateLiabilityChange + results[11].BestEstimateLiabilityChange + results[12].BestEstimateLiabilityChange + results[13].BestEstimateLiabilityChange) * (total_CSM_Change / (total_LC_Change + total_CSM_Change))
				lm.RiskAdjustment = (results[9].RiskAdjustmentChange + results[10].RiskAdjustmentChange + results[11].RiskAdjustmentChange + results[12].RiskAdjustmentChange + results[13].RiskAdjustmentChange) * (total_CSM_Change / (total_LC_Change + total_CSM_Change))

			} else {
				lm.Bel = 0
				lm.RiskAdjustment = 0
			}

			lm.Csm = total_CSM_Change
			lm.TotalLRC = lm.Bel + lm.RiskAdjustment + lm.Csm

			lm.IncurredClaimsBel = 0
			lm.IncurredClaimsRiskAdj = 0
			lm.LossComponent = 0
			lm.FRACsm = 0
			lm.MRACsm = 0
			lm.FVACsm = 0
		}

		//Changes in estimates that do not affect CSM; Loss Component Reversal
		if lml.Code == 302 {
			var total_LC_Change float64
			var total_CSM_Change float64
			lm.ProductCode = results[0].ProductCode
			lm.Ifrs17Group = results[0].IFRS17Group
			total_CSM_Change = results[9].CSMChange + results[10].CSMChange + results[11].CSMChange + results[12].CSMChange + results[13].CSMChange
			total_LC_Change = results[9].LossComponentChange + results[10].LossComponentChange + results[11].LossComponentChange + results[12].LossComponentChange + results[13].LossComponentChange
			//if total_CSM_Change+total_LC_Change != 0 {
			//	if (results[13].LossComponentBuildup - results[9].LossComponentBuildup) != 0 {
			//		if (results[9].BelOutflow + results[9].RiskAdjustment) != 0 {
			//			lm.LossComponent = (results[13].LossComponentBuildup - results[9].LossComponentBuildup) * (1 - results[9].RiskAdjustment/(results[9].BelOutflow+results[9].RiskAdjustment))
			//			lm.LossComponentRiskAdj = (results[13].LossComponentBuildup - results[9].LossComponentBuildup) * results[9].RiskAdjustment / (results[9].BelOutflow + results[9].RiskAdjustment)
			//			lm.Bel = (results[13].BEL-results[9].BEL)*(1-total_CSM_Change/(total_LC_Change+total_CSM_Change)) - lm.LossComponent
			//			lm.RiskAdjustment = -(results[10].RiskAdjustmentChange+results[11].RiskAdjustmentChange+results[12].RiskAdjustmentChange+results[13].RiskAdjustmentChange)*(1.0-total_CSM_Change/(total_LC_Change+total_CSM_Change)) - lm.LossComponentRiskAdj
			//		}
			//	} else {
			//		lm.LossComponent = 0
			//		lm.Bel = (results[13].BEL - results[9].BEL) * (1.0 - total_CSM_Change/(total_LC_Change+total_CSM_Change))
			//		lm.LossComponentRiskAdj = 0
			//		lm.RiskAdjustment = (results[10].RiskAdjustmentChange + results[11].RiskAdjustmentChange + results[12].RiskAdjustmentChange + results[13].RiskAdjustmentChange) * (1.0 - total_CSM_Change/(total_LC_Change+total_CSM_Change))
			//	}
			//} else {
			//	lm.LossComponent = 0
			//	lm.Bel = 0
			//	lm.LossComponentRiskAdj = 0
			//	lm.RiskAdjustment = 0
			//}

			if total_CSM_Change+total_LC_Change != 0 {
				lm.Bel = (results[9].BestEstimateLiabilityChange + results[10].BestEstimateLiabilityChange + results[11].BestEstimateLiabilityChange + results[12].BestEstimateLiabilityChange + results[13].BestEstimateLiabilityChange) * (1 - total_CSM_Change/(total_LC_Change+total_CSM_Change))
				lm.RiskAdjustment = (results[9].RiskAdjustmentChange + results[10].RiskAdjustmentChange + results[11].RiskAdjustmentChange + results[12].RiskAdjustmentChange + results[13].RiskAdjustmentChange) * (1.0 - total_CSM_Change/(total_LC_Change+total_CSM_Change))
			} else {
				lm.Bel = 0
				lm.RiskAdjustment = 0

			}
			lm.Csm = 0
			lm.TotalLRC = lm.Bel + lm.RiskAdjustment + lm.Csm
			lm.IncurredClaimsBel = 0
			lm.IncurredClaimsRiskAdj = 0
			lm.LossComponent = results[9].LossComponentChange + results[10].LossComponentChange + results[11].LossComponentChange + results[12].LossComponentChange + results[13].LossComponentChange //lm.Bel + lm.RiskAdjustment
			lm.FRACsm = 0
			lm.MRACsm = 0
			lm.FVACsm = 0
		}

		//New Business Initial Recognition Effects
		if lml.Code == 303 {
			lm.ProductCode = results[0].ProductCode
			lm.Ifrs17Group = results[0].IFRS17Group
			lm.IncurredClaimsBel = 0
			lm.IncurredClaimsRiskAdj = 0
			//if results[2].LossComponentChange > 0 {
			//	if results[2].BelOutflow+results[2].RiskAdjustment != 0 {
			//		lm.LossComponent = results[2].LossComponentChange * results[2].BelOutflow / (results[2].BelOutflow + results[2].RiskAdjustment)
			//		lm.LossComponentRiskAdj = results[2].LossComponentChange * (1.0 - results[2].BelOutflow/(results[2].BelOutflow+results[2].RiskAdjustment))
			//	}
			//} else {
			//	lm.Bel = results[2].BEL
			//	lm.RiskAdjustment = results[2].RiskAdjustment
			//	lm.Csm = results[2].CSMChange

			//}

			lm.Bel = results[2].BEL
			lm.RiskAdjustment = results[2].RiskAdjustment
			lm.Csm = results[2].CSMChange
			lm.LossComponent = results[2].LossComponentChange
			lm.FRACsm = 0
			lm.MRACsm = 0
			lm.FVACsm = 0
			lm.TotalLRC = lm.Bel + lm.RiskAdjustment + lm.Csm

			//Initial Recognition Analysis Report
			var nb models.InitialRecognition
			nb.RunDate = results[0].RunDate
			nb.CsmRunID = csmRun.ID
			nb.ProductCode = results[0].ProductCode
			nb.IFRS17Group = results[0].IFRS17Group
			nb.Code = 1001
			nb.Variable = "BEL Outflows"

			Profitability := nb.IFRS17Group[len(nb.IFRS17Group)-1:]

			switch Profitability {
			case "N":
				nb.NoSignificantProbabilityOnerous = results[2].BelOutflow
				nb.Onerous = 0
				nb.Remaining = 0
			case "O":
				nb.NoSignificantProbabilityOnerous = 0
				nb.Onerous = results[2].BelOutflow
				nb.Remaining = 0
			case "R":
				nb.NoSignificantProbabilityOnerous = 0
				nb.Onerous = 0
				nb.Remaining = results[2].BelOutflow
			}
			nb.BusinessAcquisitionTransfer = 0
			initialRecognition = append(initialRecognition, nb)

			nb.ProductCode = results[0].ProductCode
			nb.RunDate = results[0].RunDate
			nb.CsmRunID = csmRun.ID
			nb.IFRS17Group = results[0].IFRS17Group
			nb.Code = 1002
			nb.Variable = "BEL Inflows"
			switch Profitability {
			case "N":
				nb.NoSignificantProbabilityOnerous = -results[2].BelInflow
				nb.Onerous = 0
				nb.Remaining = 0
			case "O":
				nb.NoSignificantProbabilityOnerous = 0
				nb.Onerous = -results[2].BelInflow
				nb.Remaining = 0
			case "R":
				nb.NoSignificantProbabilityOnerous = 0
				nb.Onerous = 0
				nb.Remaining = -results[2].BelInflow
			}
			nb.BusinessAcquisitionTransfer = 0
			initialRecognition = append(initialRecognition, nb)

			nb.RunDate = results[0].RunDate
			nb.CsmRunID = csmRun.ID
			nb.ProductCode = results[0].ProductCode
			nb.IFRS17Group = results[0].IFRS17Group
			nb.Code = 1003
			nb.Variable = "Risk Adjustment"
			switch Profitability {
			case "N":
				nb.NoSignificantProbabilityOnerous = results[2].RiskAdjustment
				nb.Onerous = 0
				nb.Remaining = 0
			case "O":
				nb.NoSignificantProbabilityOnerous = 0
				nb.Onerous = results[2].RiskAdjustment
				nb.Remaining = 0
			case "R":
				nb.NoSignificantProbabilityOnerous = 0
				nb.Onerous = 0
				nb.Remaining = results[2].RiskAdjustment
			}
			nb.BusinessAcquisitionTransfer = 0
			initialRecognition = append(initialRecognition, nb)

			nb.RunDate = results[0].RunDate
			nb.CsmRunID = csmRun.ID
			nb.ProductCode = results[0].ProductCode
			nb.IFRS17Group = results[0].IFRS17Group
			nb.Code = 1004
			nb.Variable = "CSM"
			switch Profitability {
			case "N":
				nb.NoSignificantProbabilityOnerous = results[2].CSMChange
				nb.Onerous = 0
				nb.Remaining = 0
			case "O":
				nb.NoSignificantProbabilityOnerous = 0
				nb.Onerous = results[2].CSMChange
				nb.Remaining = 0
			case "R":
				nb.NoSignificantProbabilityOnerous = 0
				nb.Onerous = 0
				nb.Remaining = results[2].CSMChange
			}
			nb.BusinessAcquisitionTransfer = 0
			initialRecognition = append(initialRecognition, nb)

			//1001+1002+1003+1004
			var totalnb models.InitialRecognition
			totalnb.RunDate = results[0].RunDate
			totalnb.CsmRunID = results[0].CsmRunID
			totalnb.ProductCode = results[0].ProductCode
			totalnb.IFRS17Group = results[0].IFRS17Group
			totalnb.Code = 1050
			totalnb.Variable = "Initial Recognition"
			for _, nb := range initialRecognition {
				totalnb.NoSignificantProbabilityOnerous += nb.NoSignificantProbabilityOnerous
				totalnb.Onerous += nb.Onerous
				totalnb.Remaining += nb.Remaining
				totalnb.BusinessAcquisitionTransfer += nb.BusinessAcquisitionTransfer
			}
			initialRecognition = append(initialRecognition, totalnb)
		}

		//LIC Fulfilment Cash flow Changes
		if lml.Code == 401 {
			lm.ProductCode = results[0].ProductCode
			lm.Ifrs17Group = results[0].IFRS17Group
			lm.Bel = 0
			lm.RiskAdjustment = 0
			lm.IncurredClaimsBel = 0     // to add
			lm.IncurredClaimsRiskAdj = 0 // to add
			lm.Csm = 0
			lm.LossComponent = 0
			lm.FRACsm = 0
			lm.MRACsm = 0
			lm.FVACsm = 0
			lm.TotalLRC = lm.Bel + lm.RiskAdjustment + lm.Csm
		}

		//Experience Adjustment/ New Claims
		if lml.Code == 402 {
			lm.ProductCode = results[0].ProductCode
			lm.Ifrs17Group = results[0].IFRS17Group
			lm.Bel = 0
			lm.RiskAdjustment = 0
			lm.IncurredClaimsBel = 0     // to add
			lm.IncurredClaimsRiskAdj = 0 // to add
			lm.Csm = 0
			lm.LossComponent = 0
			lm.FRACsm = 0
			lm.MRACsm = 0
			lm.FVACsm = 0
			lm.TotalLRC = lm.Bel + lm.RiskAdjustment + lm.Csm
		}

		//Interest Accretion
		if lml.Code == 601 {
			lm.ProductCode = results[0].ProductCode
			lm.Ifrs17Group = results[0].IFRS17Group
			lm.Bel = results[8].BestEstimateLiabilityChange
			lm.RiskAdjustment = results[8].RiskAdjustmentChange
			lm.IncurredClaimsBel = 0     // to add
			lm.IncurredClaimsRiskAdj = 0 // to add
			lm.Csm = 0                   //results[8].CSMChange moved to cash flows as non-SCI interest
			lm.LossComponent = results[8].LossComponentChange
			lm.FRACsm = 0
			lm.MRACsm = 0
			lm.FVACsm = 0
			lm.TotalLRC = lm.Bel + lm.RiskAdjustment + lm.Csm
		}

		//Net Finance Expense
		if lml.Code == 602 {
			lm.ProductCode = results[0].ProductCode
			lm.Ifrs17Group = results[0].IFRS17Group
			lm.Bel = results[1].BestEstimateLiabilityChange + results[16].BestEstimateLiabilityChange
			lm.RiskAdjustment = results[1].RiskAdjustmentChange + results[16].RiskAdjustmentChange
			lm.IncurredClaimsBel = 0     // to add
			lm.IncurredClaimsRiskAdj = 0 // to add
			lm.Csm = 0
			lm.LossComponent = 0
			lm.FRACsm = 0
			lm.MRACsm = 0
			lm.FVACsm = 0
			lm.TotalLRC = lm.Bel + lm.RiskAdjustment + lm.Csm
		}

		//Exchange Rate Changes
		if lml.Code == 603 {
			lm.ProductCode = results[0].ProductCode
			lm.Ifrs17Group = results[0].IFRS17Group
			lm.Bel = 0
			lm.RiskAdjustment = 0
			lm.IncurredClaimsBel = 0     // to add
			lm.IncurredClaimsRiskAdj = 0 // to add
			lm.Csm = 0
			lm.LossComponent = 0
			lm.FRACsm = 0
			lm.MRACsm = 0
			lm.FVACsm = 0
			lm.TotalLRC = lm.Bel + lm.RiskAdjustment + lm.Csm
		}

		//Premiums Received
		if lml.Code == 801 {
			lm.ProductCode = results[0].ProductCode
			lm.Ifrs17Group = results[0].IFRS17Group
			var financeVariables = getFinancialVariables(lm.ProductCode, lm.Ifrs17Group, csmRun.FinanceYear, csmRun.FinanceVersion)
			lm.Bel = financeVariables.ActualPremiumIncome
			lm.RiskAdjustment = 0
			lm.IncurredClaimsBel = 0
			lm.IncurredClaimsRiskAdj = 0
			lm.Csm = 0
			lm.LossComponent = 0
			lm.FRACsm = 0
			lm.MRACsm = 0
			lm.FVACsm = 0
			lm.TotalLRC = lm.Bel + lm.RiskAdjustment + lm.Csm

		}

		//Insurance Acquisition Cashflows
		if lml.Code == 802 {
			lm.ProductCode = results[0].ProductCode
			lm.Ifrs17Group = results[0].IFRS17Group
			//var financeVariables = getFinancialVariables(lm.ProductCode, lm.Ifrs17Group, csmRun.FinanceYear, csmRun.FinanceVersion)
			lm.Bel = -(results[8].ExpectedCashOutflow - results[4].PNLChange - results[5].PNLChange - results[6].PNLChange - results[7].PNLChange) //-financeVariables.ActualAcquisitionExpenses
			lm.RiskAdjustment = 0
			lm.IncurredClaimsBel = 0
			lm.IncurredClaimsRiskAdj = 0
			lm.Csm = 0
			lm.LossComponent = 0
			lm.FRACsm = 0
			lm.MRACsm = 0
			lm.FVACsm = 0
			lm.TotalLRC = lm.Bel + lm.RiskAdjustment + lm.Csm

		}

		//Claims and Other Insurance Service Expenses Paid
		if lml.Code == 803 {
			lm.ProductCode = results[0].ProductCode
			lm.Ifrs17Group = results[0].IFRS17Group
			var financeVariables = getFinancialVariables(lm.ProductCode, lm.Ifrs17Group, csmRun.FinanceYear, csmRun.FinanceVersion)
			lm.Bel = -(financeVariables.ActualMortalityClaimsIncurred + financeVariables.ActualMorbidityClaimsIncurred + financeVariables.ActualRetrenchmentClaimsIncurred + financeVariables.ActualAttributableExpenses)
			lm.RiskAdjustment = 0
			lm.IncurredClaimsBel = 0 //financeVariables.ActualMortalityClaims + financeVariables.ActualMorbidityClaims + financeVariables.ActualRetrenchmentClaims
			lm.IncurredClaimsRiskAdj = 0
			lm.Csm = 0
			lm.LossComponent = 0
			lm.FRACsm = 0
			lm.MRACsm = 0
			lm.FVACsm = 0
			lm.TotalLRC = lm.Bel + lm.RiskAdjustment + lm.Csm

		}

		//Claims and Other Insurance Service Expenses Paid
		if lml.Code == 804 {
			lm.ProductCode = results[0].ProductCode
			lm.Ifrs17Group = results[0].IFRS17Group
			lm.Bel = 0
			lm.RiskAdjustment = 0
			lm.IncurredClaimsBel = 0 //financeVariables.ActualMortalityClaims + financeVariables.ActualMorbidityClaims + financeVariables.ActualRetrenchmentClaims
			lm.IncurredClaimsRiskAdj = 0
			lm.Csm = results[8].CSMChange
			lm.LossComponent = 0
			lm.FRACsm = 0
			lm.MRACsm = 0
			lm.FVACsm = 0
			lm.TotalLRC = lm.Bel + lm.RiskAdjustment + lm.Csm

		}

		if lml.Code == 150 {
			lm.ProductCode = results[0].ProductCode
			lm.Ifrs17Group = results[0].IFRS17Group
			lm.Bel = -(liabilityMovements[0].Bel + liabilityMovements[1].Bel)
			lm.RiskAdjustment = -(liabilityMovements[0].RiskAdjustment + liabilityMovements[1].RiskAdjustment)
			lm.IncurredClaimsBel = 0
			lm.IncurredClaimsRiskAdj = 0
			lm.Csm = -(liabilityMovements[0].Csm + liabilityMovements[1].Csm)
			lm.LossComponent = -(liabilityMovements[0].LossComponent + liabilityMovements[1].LossComponent)
			lm.FRACsm = 0
			lm.MRACsm = 0
			lm.FVACsm = 0
			lm.TotalLRC = lm.Bel + lm.RiskAdjustment + lm.Csm

		}

		if lml.Code == 250 {
			lm.ProductCode = results[0].ProductCode
			lm.Ifrs17Group = results[0].IFRS17Group
			lm.Bel = liabilityMovements[3].Bel + liabilityMovements[4].Bel + liabilityMovements[5].Bel + liabilityMovements[6].Bel
			lm.RiskAdjustment = liabilityMovements[3].RiskAdjustment + liabilityMovements[4].RiskAdjustment + liabilityMovements[5].RiskAdjustment + liabilityMovements[6].RiskAdjustment
			lm.IncurredClaimsBel = 0
			lm.IncurredClaimsRiskAdj = 0
			lm.Csm = liabilityMovements[3].Csm + liabilityMovements[4].Csm + liabilityMovements[5].Csm + liabilityMovements[6].Csm
			lm.LossComponent = liabilityMovements[3].LossComponent + liabilityMovements[4].LossComponent + liabilityMovements[5].LossComponent + liabilityMovements[6].LossComponent
			lm.FRACsm = 0
			lm.MRACsm = 0
			lm.FVACsm = 0
			lm.TotalLRC = lm.Bel + lm.RiskAdjustment + lm.Csm

		}
		if lml.Code == 350 {
			lm.ProductCode = results[0].ProductCode
			lm.Ifrs17Group = results[0].IFRS17Group
			lm.Bel = liabilityMovements[8].Bel + liabilityMovements[9].Bel + liabilityMovements[10].Bel
			lm.RiskAdjustment = liabilityMovements[8].RiskAdjustment + liabilityMovements[9].RiskAdjustment + liabilityMovements[10].RiskAdjustment
			lm.IncurredClaimsBel = 0
			lm.IncurredClaimsRiskAdj = 0
			lm.Csm = liabilityMovements[8].Csm + liabilityMovements[9].Csm + liabilityMovements[10].Csm
			lm.LossComponent = liabilityMovements[8].LossComponent + liabilityMovements[9].LossComponent + liabilityMovements[10].LossComponent
			lm.FRACsm = 0
			lm.MRACsm = 0
			lm.FVACsm = 0
			lm.TotalLRC = lm.Bel + lm.RiskAdjustment + lm.Csm

		}

		if lml.Code == 450 {
			lm.ProductCode = results[0].ProductCode
			lm.Ifrs17Group = results[0].IFRS17Group
			lm.Bel = liabilityMovements[12].Bel + liabilityMovements[13].Bel
			lm.RiskAdjustment = liabilityMovements[12].RiskAdjustment + liabilityMovements[13].RiskAdjustment
			lm.IncurredClaimsBel = 0
			lm.IncurredClaimsRiskAdj = 0
			lm.Csm = liabilityMovements[12].Csm + liabilityMovements[13].Csm
			lm.LossComponent = liabilityMovements[12].LossComponent + liabilityMovements[13].LossComponent
			lm.FRACsm = 0
			lm.MRACsm = 0
			lm.FVACsm = 0
			lm.TotalLRC = lm.Bel + lm.RiskAdjustment + lm.Csm

		}

		if lml.Code == 550 {
			lm.ProductCode = results[0].ProductCode
			lm.Ifrs17Group = results[0].IFRS17Group
			lm.Bel = liabilityMovements[7].Bel + liabilityMovements[11].Bel + liabilityMovements[14].Bel
			lm.RiskAdjustment = liabilityMovements[7].RiskAdjustment + liabilityMovements[11].RiskAdjustment + liabilityMovements[14].RiskAdjustment
			lm.IncurredClaimsBel = 0
			lm.IncurredClaimsRiskAdj = 0
			lm.Csm = liabilityMovements[7].Csm + liabilityMovements[11].Csm + liabilityMovements[14].Csm
			lm.LossComponent = liabilityMovements[7].LossComponent + liabilityMovements[11].LossComponent + liabilityMovements[14].LossComponent
			lm.FRACsm = 0
			lm.MRACsm = 0
			lm.FVACsm = 0
			lm.TotalLRC = lm.Bel + lm.RiskAdjustment + lm.Csm

		}
		if lml.Code == 650 {
			lm.ProductCode = results[0].ProductCode
			lm.Ifrs17Group = results[0].IFRS17Group
			lm.Bel = liabilityMovements[16].Bel + liabilityMovements[17].Bel + liabilityMovements[18].Bel
			lm.RiskAdjustment = liabilityMovements[16].RiskAdjustment + liabilityMovements[17].RiskAdjustment + liabilityMovements[18].RiskAdjustment
			lm.IncurredClaimsBel = 0
			lm.IncurredClaimsRiskAdj = 0
			lm.Csm = liabilityMovements[16].Csm + liabilityMovements[17].Csm + liabilityMovements[18].Csm
			lm.LossComponent = liabilityMovements[16].LossComponent + liabilityMovements[17].LossComponent + liabilityMovements[18].LossComponent
			lm.FRACsm = 0
			lm.MRACsm = 0
			lm.FVACsm = 0
			lm.TotalLRC = lm.Bel + lm.RiskAdjustment + lm.Csm

		}
		if lml.Code == 750 {
			lm.ProductCode = results[0].ProductCode
			lm.Ifrs17Group = results[0].IFRS17Group
			lm.Bel = liabilityMovements[15].Bel + liabilityMovements[19].Bel
			lm.RiskAdjustment = liabilityMovements[15].RiskAdjustment + liabilityMovements[19].RiskAdjustment
			lm.IncurredClaimsBel = 0
			lm.IncurredClaimsRiskAdj = 0
			lm.Csm = liabilityMovements[15].Csm + liabilityMovements[19].Csm
			lm.LossComponent = liabilityMovements[15].LossComponent + liabilityMovements[19].LossComponent
			lm.FRACsm = 0
			lm.MRACsm = 0
			lm.FVACsm = 0
			lm.TotalLRC = lm.Bel + lm.RiskAdjustment + lm.Csm

		}

		if lml.Code == 850 {
			lm.ProductCode = results[0].ProductCode
			lm.Ifrs17Group = results[0].IFRS17Group
			lm.Bel = liabilityMovements[21].Bel + liabilityMovements[22].Bel + liabilityMovements[23].Bel + liabilityMovements[24].Bel
			lm.RiskAdjustment = liabilityMovements[21].RiskAdjustment + liabilityMovements[22].RiskAdjustment + liabilityMovements[23].RiskAdjustment + liabilityMovements[24].RiskAdjustment
			lm.IncurredClaimsBel = 0
			lm.IncurredClaimsRiskAdj = 0
			lm.Csm = liabilityMovements[21].Csm + liabilityMovements[22].Csm + liabilityMovements[23].Csm + liabilityMovements[24].Csm
			lm.LossComponent = liabilityMovements[21].LossComponent + liabilityMovements[22].LossComponent + liabilityMovements[23].LossComponent + liabilityMovements[24].LossComponent
			lm.FRACsm = 0
			lm.MRACsm = 0
			lm.FVACsm = 0
			lm.TotalLRC = lm.Bel + lm.RiskAdjustment + lm.Csm

		}

		if lml.Code == 900 {
			lm.ProductCode = results[0].ProductCode
			lm.Ifrs17Group = results[0].IFRS17Group
			lm.Bel = liabilityMovements[20].Bel + liabilityMovements[25].Bel + liabilityMovements[2].Bel
			lm.RiskAdjustment = liabilityMovements[20].RiskAdjustment + liabilityMovements[25].RiskAdjustment + liabilityMovements[2].RiskAdjustment
			lm.IncurredClaimsBel = 0
			lm.IncurredClaimsRiskAdj = 0
			lm.Csm = liabilityMovements[20].Csm + liabilityMovements[25].Csm + liabilityMovements[2].Csm
			lm.LossComponent = liabilityMovements[20].LossComponent + liabilityMovements[25].LossComponent + liabilityMovements[2].LossComponent
			lm.FRACsm = 0
			lm.MRACsm = 0
			lm.FVACsm = 0
			lm.TotalLRC = lm.Bel + lm.RiskAdjustment + lm.Csm
		}

		if lml.Code == 901 {
			lm.ProductCode = results[0].ProductCode
			lm.Ifrs17Group = results[0].IFRS17Group
			lm.Bel = results[16].BelInflow
			lm.RiskAdjustment = 0
			lm.IncurredClaimsBel = 0
			lm.IncurredClaimsRiskAdj = 0
			lm.Csm = 0
			lm.LossComponent = 0
			lm.FRACsm = 0
			lm.MRACsm = 0
			lm.FVACsm = 0
			lm.TotalLRC = lm.Bel + lm.RiskAdjustment + lm.Csm

		}

		if lml.Code == 902 {
			lm.ProductCode = results[0].ProductCode
			lm.Ifrs17Group = results[0].IFRS17Group
			lm.Bel = -results[16].BelOutflow
			lm.RiskAdjustment = -results[16].RiskAdjustment
			lm.IncurredClaimsBel = 0
			lm.IncurredClaimsRiskAdj = 0
			lm.Csm = -results[16].CSMBuildup
			lm.LossComponent = -results[16].LossComponentBuildup
			lm.FRACsm = 0
			lm.MRACsm = 0
			lm.FVACsm = 0
			lm.TotalLRC = lm.Bel + lm.RiskAdjustment + lm.Csm

		}

		liabilityMovements = append(liabilityMovements, lm)
	}

	//for j, _ := range *liabilityMovements {
	//	//Aggregation Lines
	//	totals := []int{150, 250, 350, 450, 550, 650, 750, 850, 950}
	//	for i := 0; i < len(totals); i++ {
	//		if (*liabilityMovements)[j].Code == totals[i] {
	//			(*liabilityMovements)[j].ProductCode = results[0].ProductCode
	//			(*liabilityMovements)[j].Ifrs17Group = results[0].IFRS17Group
	//			for r, _ := range *liabilityMovements {
	//				if (*liabilityMovements)[r].Code < totals[i] && int((*liabilityMovements)[r].Code/100) == int(totals[i]/100) {
	//					(*liabilityMovements)[j].Bel = (*liabilityMovements)[j].Bel + (*liabilityMovements)[r].Bel
	//					(*liabilityMovements)[j].RiskAdjustment = (*liabilityMovements)[j].RiskAdjustment + (*liabilityMovements)[r].RiskAdjustment
	//					(*liabilityMovements)[j].IncurredClaimsBel = (*liabilityMovements)[j].IncurredClaimsBel + (*liabilityMovements)[r].IncurredClaimsBel
	//					(*liabilityMovements)[j].IncurredClaimsRiskAdj = (*liabilityMovements)[j].IncurredClaimsRiskAdj + (*liabilityMovements)[r].IncurredClaimsRiskAdj
	//					(*liabilityMovements)[j].Csm = (*liabilityMovements)[j].Csm + (*liabilityMovements)[r].Csm
	//					(*liabilityMovements)[j].LossComponent = (*liabilityMovements)[j].LossComponent + (*liabilityMovements)[r].LossComponent
	//					(*liabilityMovements)[j].FRACsm = (*liabilityMovements)[j].FRACsm + (*liabilityMovements)[r].FRACsm
	//					(*liabilityMovements)[j].MRACsm = (*liabilityMovements)[j].MRACsm + (*liabilityMovements)[r].MRACsm
	//					(*liabilityMovements)[j].FVACsm = (*liabilityMovements)[j].FVACsm + (*liabilityMovements)[r].FVACsm
	//				}
	//			}
	//		}
	//	}
	//}

	//Add 250 line for insuranceRevenues
	var totalIR models.InsuranceRevenue
	totalIR.ProductCode = results[0].ProductCode
	totalIR.IFRS17Group = results[0].IFRS17Group
	totalIR.RunDate = results[0].RunDate
	totalIR.CsmRunID = results[0].CsmRunID
	totalIR.Code = 250
	totalIR.Variable = "Current Service Changes"

	for _, ir := range insuranceRevenues {
		totalIR.PostTransition += ir.PostTransition
		totalIR.FullRetrospectiveApproach += ir.FullRetrospectiveApproach
		totalIR.ModifiedRetrospectiveApproach += ir.ModifiedRetrospectiveApproach
		totalIR.FairValueApproach += ir.FairValueApproach
	}

	insuranceRevenues = append(insuranceRevenues, totalIR)

	err = DB.Where("run_date = ? and ifrs17_group = ? and product_code = ?", results[0].RunDate, results[0].IFRS17Group, results[0].ProductCode).Delete(&models.InsuranceRevenue{}).Error
	if err != nil {
		log.Error(err)
	}

	err = DB.Create(&insuranceRevenues).Error
	if err != nil {
		log.Error(err)
	}

	err = DB.Where("run_date = ? and ifrs17_group = ? and product_code = ?", results[0].RunDate, results[0].IFRS17Group, results[0].ProductCode).Delete(&models.LiabilityMovement{}).Error
	if err != nil {
		log.Error(err)
	}

	err = DB.Create(&liabilityMovements).Error
	if err != nil {
		log.Error(err)
	}

	err = DB.Where("run_date = ? and ifrs17_group = ? and product_code = ?", results[0].RunDate, results[0].IFRS17Group, results[0].ProductCode).Delete(&models.InitialRecognition{}).Error
	err = DB.Create(&initialRecognition).Error
	if err != nil {
		log.Error(err)
	}
}

func GetLiabilityMovements(runDate, productCode, ifrs17Group string) []models.LiabilityMovement {
	var movements []models.LiabilityMovement
	if ifrs17Group != "" {
		//DB.Where("run_date = ? and product_code = ? and ifrs17_group = ?", runDate,productCode, ifrs17Group).Find(&movements)
		query := fmt.Sprintf("SELECT product_code,ifrs17_group, run_date, code, variable, sum(bel) as bel, sum(risk_adjustment) as risk_adjustment, sum(csm) as csm, sum(total_lrc) as total_lrc, sum(incurred_claims_bel) as incurred_claims_bel, sum(incurred_claims_risk_adjustment) as incurred_claims_risk_adjustment, sum(loss_component) as loss_component, sum(fra_csm) as fra_csm, sum(mra_csm) as mra_csm, sum(fva_csm) as fva_csm from liability_movements where run_date = '%s' and product_code = '%s' and ifrs17_group = '%s' group by product_code, run_date, code, variable", runDate, productCode, ifrs17Group)
		err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Raw(query).Scan(&movements).Error
		})
		if err != nil {
			fmt.Println(err)
		}
	}
	if productCode != "" && ifrs17Group == "" {
		//remember to add aggregation of the values here like AOS Step results
		query := fmt.Sprintf("SELECT product_code, run_date, code, variable, sum(bel) as bel, sum(risk_adjustment) as risk_adjustment,sum(csm) as csm, sum(total_lrc) as total_lrc, sum(incurred_claims_bel) as incurred_claims_bel, sum(incurred_claims_risk_adjustment) as incurred_claims_risk_adjustment, sum(loss_component) as loss_component, sum(fra_csm) as fra_csm, sum(mra_csm) as mra_csm, sum(fva_csm) as fva_csm from liability_movements where run_date = '%s' and product_code = '%s' group by product_code, code, variable", runDate, productCode)
		err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Raw(query).Scan(&movements).Error
		})
		if err != nil {
			fmt.Println(err)
		}
	}
	if ifrs17Group == "" && productCode == "" {
		query := fmt.Sprintf("SELECT run_date, code, variable, sum(bel) as bel, sum(risk_adjustment) as risk_adjustment,sum(csm) as csm, sum(total_lrc) as total_lrc, sum(incurred_claims_bel) as incurred_claims_bel, sum(incurred_claims_risk_adjustment) as incurred_claims_risk_adjustment, sum(loss_component) as loss_component, sum(fra_csm) as fra_csm, sum(mra_csm) as mra_csm, sum(fva_csm) as fva_csm from liability_movements where run_date = '%s' group by code, variable", runDate)
		err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Raw(query).Scan(&movements).Error
		})
		if err != nil {
			fmt.Println(err)
		}
	}
	return reduceDecimals(movements)
}

func GetInsuranceRevenueAnalysisRuns() []models.CsmRun {
	var runs []models.CsmRun
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Table("insurance_revenues").Distinct("run_date").Where("run_date <> ''").Find(&runs).Error }); err != nil {
		fmt.Println(err)
	}
	return runs
}

func GetInitialRecognitionAnalysisRuns() []models.CsmRun {
	var runs []models.CsmRun
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Table("initial_recognitions").Distinct("run_date").Where("run_date <> ''").Find(&runs).Error }); err != nil {
		fmt.Println(err)
	}
	return runs

}

func GetInsuranceRevenueAnalysis(runDate, productCode, ifrs17Group string) map[string]interface{} {
	var results = map[string]interface{}{}
	var irs []models.InsuranceRevenue
	var list []models.InsuranceRevenue
	if ifrs17Group != "" && productCode != "" {
		if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Where("run_date = ? product_code = ? and ifrs17_group = ?", runDate, productCode, ifrs17Group).Find(&irs).Error }); err != nil {
			fmt.Println(err)
		}
		query := fmt.Sprintf("SELECT product_code,ifrs17_group, code, variable, sum(post_transition) as post_transition, sum(full_retrospective_approach) as full_retrospective_approach, sum(modified_retrospective_approach) as modified_retrospective_approach, sum(fair_value_approach) as fair_value_approach from insurance_revenues where run_date = '%s' and  product_code = '%s' and ifrs17_group = '%s' group by product_code, ifrs17_group, code, variable", runDate, productCode, ifrs17Group)
		err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Raw(query).Scan(&irs).Error
		})
		if err != nil {
			fmt.Println(err)
		}
	} else if ifrs17Group == "" && productCode != "" {
		//remember to add aggregation of the values here like AOS Step results
		query := fmt.Sprintf("SELECT run_date, product_code, code, variable, sum(post_transition) as post_transition, sum(full_retrospective_approach) as full_retrospective_approach, sum(modified_retrospective_approach) as modified_retrospective_approach, sum(fair_value_approach) as fair_value_approach from insurance_revenues where run_date = '%s' and product_code = '%s'   group by  product_code, code, variable", runDate, productCode)
		err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Raw(query).Scan(&irs).Error
		})
		if err != nil {
			fmt.Println(err)
		}
		if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Distinct("ifrs17_group").Where("run_date = ? and product_code = ?", runDate, productCode).Find(&list).Error }); err != nil {
			fmt.Println(err)
		}
	} else {
		query := fmt.Sprintf("SELECT run_date, code, variable, sum(post_transition) as post_transition, sum(full_retrospective_approach) as full_retrospective_approach, sum(modified_retrospective_approach) as modified_retrospective_approach, sum(fair_value_approach) as fair_value_approach from insurance_revenues where run_date = '%s'  group by  code, variable", runDate)
		err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Raw(query).Scan(&irs).Error
		})
		if err != nil {
			fmt.Println(err)
		}
		if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Distinct("product_code").Where("run_date = ?", runDate).Find(&list).Error }); err != nil {
			fmt.Println(err)
		}
	}
	reducedIrs := reduceInsuranceDecimals(irs)
	results["report"] = reducedIrs
	results["list"] = list
	return results
}

func GetInitialRecognitionAnalysis(runDate, productCode, ifrs17Group string) map[string]interface{} {
	var nbs []models.InitialRecognition
	var list []models.InitialRecognition
	var results = map[string]interface{}{}
	if ifrs17Group != "" && productCode != "" {
		if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Where("product_code = ? and ifrs17_group = ?", productCode, ifrs17Group).Find(&nbs).Error }); err != nil {
			fmt.Println(err)
		}
		query := fmt.Sprintf("SELECT run_date, product_code,ifrs17_group, code, variable, sum(no_significant_probability_onerous) as no_significant_probability_onerous, sum(onerous) as onerous, sum(remaining) as remaining, sum(business_acquisition_transfer) as business_acquisition_transfer from initial_recognitions where run_date = '%s' and product_code = '%s' and ifrs17_group = '%s' group by product_code, ifrs17_group, code, variable", runDate, productCode, ifrs17Group)
		err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Raw(query).Scan(&nbs).Error
		})
		if err != nil {
			fmt.Println(err)
		}
	} else if productCode != "" && ifrs17Group == "" {
		//remember to add aggregation of the values here like AOS Step results
		query := fmt.Sprintf("SELECT run_date, product_code, code, variable, sum(no_significant_probability_onerous) as no_significant_probability_onerous, sum(onerous) as onerous, sum(remaining) as remaining, sum(business_acquisition_transfer) as business_acquisition_transfer from initial_recognitions where run_date = '%s' and product_code = '%s' group by product_code, code, variable", runDate, productCode)
		err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Raw(query).Scan(&nbs).Error
		})
		if err != nil {
			fmt.Println(err)
		}
		if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Distinct("ifrs17_group").Where("run_date = ? and product_code = ?", runDate, productCode).Find(&list).Error }); err != nil {
			fmt.Println(err)
		}
	} else {
		query := fmt.Sprintf("SELECT run_date, code, variable, sum(no_significant_probability_onerous) as no_significant_probability_onerous, sum(onerous) as onerous, sum(remaining) as remaining, sum(business_acquisition_transfer) as business_acquisition_transfer from initial_recognitions where run_date = '%s' group by code, variable", runDate)
		err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Raw(query).Scan(&nbs).Error
		})
		if err != nil {
			fmt.Println(err)
		}
		if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Distinct("product_code").Where("run_date = ?", runDate).Find(&list).Error }); err != nil {
			fmt.Println(err)
		}
	}
	results["report"] = reduceInitialRecognitionDecimals(nbs)
	results["list"] = list
	return results
}

func GetCSMProjectionsRunList() []models.CsmRun {
	var runs []models.CsmRun
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Table("csm_projections").Distinct("run_date").Find(&runs).Error }); err != nil {
		fmt.Println(err)
	}
	return runs
}

func GetCsmProjections(runDate, productCode, ifrs17Group string) map[string]interface{} {
	var results = make(map[string]interface{})
	var projections []models.CsmProjection
	var list []models.CsmProjection
	if ifrs17Group != "" && productCode != "" {
		if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Where("run_date = ? and product_code = ? and ifrs17_group = ?", runDate, productCode, ifrs17Group).Find(&projections).Error }); err != nil {
			fmt.Println(err)
		}
		query := fmt.Sprintf("SELECT run_date, product_code,ifrs17_group, projection_month, sum(csm_total) as csm_total, sum(csm_release_total) as csm_release_total, sum(coverage_units) as coverage_units, sum(csm) as csm,sum(csm_release) as csm_release,sum(coverage_units_full_retrospective_approach) as coverage_units_full_retrospective_approach,sum(csm_full_retrospective_approach) as csm_full_retrospective_approach,sum(csm_release_full_retrospective_approach) as csm_release_full_retrospective_approach,sum(coverage_units_modified_retrospective_approach) as coverage_units_modified_retrospective_approach,sum(csm_modified_retrospective_approach) as csm_modified_retrospective_approach,sum(csm_release_modified_retrospective_approach) as csm_release_modified_retrospective_approach,sum(coverage_units_fair_value_approach) as coverage_units_fair_value_approach,sum(csm_fair_value_approach) as csm_fair_value_approach,sum(csm_release_fair_value_approach) as csm_release_fair_value_approach from csm_projections where run_date = '%s' and  product_code = '%s' and ifrs17_group = '%s' group by product_code, ifrs17_group, projection_month", runDate, productCode, ifrs17Group)
		err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Raw(query).Scan(&projections).Error
		})
		if err != nil {
			fmt.Println(err)
		}
	} else if ifrs17Group == "" && productCode != "" {
		//remember to add aggregation of the values here like AOS Step results
		query := fmt.Sprintf("SELECT run_date, product_code, projection_month, sum(csm_total) as csm_total, sum(csm_release_total) as csm_release_total, sum(coverage_units) as coverage_units, sum(csm) as csm,sum(csm_release) as csm_release,sum(coverage_units_full_retrospective_approach) as coverage_units_full_retrospective_approach,sum(csm_full_retrospective_approach) as csm_full_retrospective_approach,sum(csm_release_full_retrospective_approach) as csm_release_full_retrospective_approach,sum(coverage_units_modified_retrospective_approach) as coverage_units_modified_retrospective_approach,sum(csm_modified_retrospective_approach) as csm_modified_retrospective_approach,sum(csm_release_modified_retrospective_approach) as csm_release_modified_retrospective_approach,sum(coverage_units_fair_value_approach) as coverage_units_fair_value_approach,sum(csm_fair_value_approach) as csm_fair_value_approach,sum(csm_release_fair_value_approach) as csm_release_fair_value_approach from csm_projections where run_date = '%s' and product_code = '%s' group by product_code, projection_month", runDate, productCode)
		err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Raw(query).Scan(&projections).Error
		})
		if err != nil {
			fmt.Println(err)
		}
		if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Distinct("ifrs17_group").Where("run_date = ? and product_code = ?", runDate, productCode).Find(&list).Error }); err != nil {
			fmt.Println(err)
		}
	} else if ifrs17Group == "" && productCode == "" {
		query := fmt.Sprintf("SELECT run_date, projection_month, sum(csm_total) as csm_total, sum(csm_release_total) as csm_release_total, sum(coverage_units) as coverage_units, sum(csm) as csm,sum(csm_release) as csm_release,sum(coverage_units_full_retrospective_approach) as coverage_units_full_retrospective_approach,sum(csm_full_retrospective_approach) as csm_full_retrospective_approach,sum(csm_release_full_retrospective_approach) as csm_release_full_retrospective_approach,sum(coverage_units_modified_retrospective_approach) as coverage_units_modified_retrospective_approach,sum(csm_modified_retrospective_approach) as csm_modified_retrospective_approach,sum(csm_release_modified_retrospective_approach) as csm_release_modified_retrospective_approach,sum(coverage_units_fair_value_approach) as coverage_units_fair_value_approach,sum(csm_fair_value_approach) as csm_fair_value_approach,sum(csm_release_fair_value_approach) as csm_release_fair_value_approach from csm_projections where run_date = '%s'  group by  projection_month", runDate)
		err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Raw(query).Scan(&projections).Error
		})
		if err != nil {
			fmt.Println(err)
		}
		if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Distinct("product_code").Where("run_date = ?", runDate).Find(&list).Error }); err != nil {
			fmt.Println(err)
		}
	}
	results["report"] = reduceCsmProjectionDecimals(projections)
	results["list"] = list
	return results
}

func ProcessIFRS17Tables(ctx context.Context, file *multipart.FileHeader, tableType string, year string, version string) error {
	fmt.Println("Processing IFRS17 tables for projections")
	fmt.Println(file)
	fmt.Println(tableType)
	switch tableType {
	case "PAA Finance":
		financeYear, _ := strconv.Atoi(year)
		err := ProcessPAAFinance(ctx, file, financeYear, version)
		if err != nil {
			fmt.Println(err)
			return err
		}
	case "Finance Variables":
		financeYear, _ := strconv.Atoi(year)
		err := ProcessFinanceVariables(ctx, file, financeYear, version)
		if err != nil {
			fmt.Println(err)
			return err
		}
	case "RA Factors":
		financeYear, _ := strconv.Atoi(year)
		err := ProcessRiskAdjustmentFactorTables(ctx, file, financeYear, version)
		if err != nil {
			fmt.Println(err)
			return err
		}
	case "Transition Adjustments":
		financeYear, _ := strconv.Atoi(year)
		err := ProcessTransitionAdjustments(ctx, file, financeYear, version)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func reduceInsuranceDecimals(irs []models.InsuranceRevenue) []models.InsuranceRevenue {
	for i, _ := range irs {
		irs[i].PostTransition = math.Round(irs[i].PostTransition)
		irs[i].FullRetrospectiveApproach = math.Round(irs[i].FullRetrospectiveApproach)
		irs[i].ModifiedRetrospectiveApproach = math.Round(irs[i].ModifiedRetrospectiveApproach)
		irs[i].FairValueApproach = math.Round(irs[i].FairValueApproach)
	}
	return irs
}

func reduceDecimals(movements []models.LiabilityMovement) []models.LiabilityMovement {
	for i, _ := range movements {
		movements[i].Bel = math.Round(movements[i].Bel)
		movements[i].RiskAdjustment = math.Round(movements[i].RiskAdjustment)
		movements[i].IncurredClaimsBel = math.Round(movements[i].IncurredClaimsBel)
		movements[i].IncurredClaimsRiskAdj = math.Round(movements[i].IncurredClaimsRiskAdj)
		movements[i].Csm = math.Round(movements[i].Csm)
		movements[i].LossComponent = math.Round(movements[i].LossComponent)
		movements[i].FRACsm = math.Round(movements[i].FRACsm)
		movements[i].MRACsm = math.Round(movements[i].MRACsm)
		movements[i].FVACsm = math.Round(movements[i].FVACsm)
	}
	return movements
}

func reduceInitialRecognitionDecimals(nbs []models.InitialRecognition) []models.InitialRecognition {
	for i, _ := range nbs {
		nbs[i].NoSignificantProbabilityOnerous = math.Round(nbs[i].NoSignificantProbabilityOnerous)
		nbs[i].Onerous = math.Round(nbs[i].Onerous)
		nbs[i].Remaining = math.Round(nbs[i].Remaining)
		nbs[i].BusinessAcquisitionTransfer = math.Round(nbs[i].BusinessAcquisitionTransfer)
	}
	return nbs
}

func reduceCsmProjectionDecimals(csmproj []models.CsmProjection) []models.CsmProjection {
	for i, _ := range csmproj {
		csmproj[i].CsmTotal = math.Round(csmproj[i].CsmTotal)
		csmproj[i].CsmReleaseTotal = math.Round(csmproj[i].CsmReleaseTotal)
		csmproj[i].CoverageUnits = math.Round(csmproj[i].CoverageUnits)
		csmproj[i].Csm = math.Round(csmproj[i].Csm)
		csmproj[i].CsmRelease = math.Round(csmproj[i].CsmRelease)
		csmproj[i].CoverageUnitsFullRetrospectiveApproach = math.Round(csmproj[i].CoverageUnitsFullRetrospectiveApproach)
		csmproj[i].CsmFullRetrospectiveApproach = math.Round(csmproj[i].CsmFullRetrospectiveApproach)
		csmproj[i].CsmReleaseFullRetrospectiveApproach = math.Round(csmproj[i].CsmReleaseFullRetrospectiveApproach)
		csmproj[i].CoverageUnitsModifiedRetrospectiveApproach = math.Round(csmproj[i].CoverageUnitsModifiedRetrospectiveApproach)
		csmproj[i].CsmModifiedRetrospectiveApproach = math.Round(csmproj[i].CsmModifiedRetrospectiveApproach)
		csmproj[i].CsmReleaseModifiedRetrospectiveApproach = math.Round(csmproj[i].CsmReleaseModifiedRetrospectiveApproach)
		csmproj[i].CoverageUnitsFairValueApproach = math.Round(csmproj[i].CoverageUnitsFairValueApproach)
		csmproj[i].CsmFairValueApproach = math.Round(csmproj[i].CsmFairValueApproach)
		csmproj[i].CsmReleaseFairValueApproach = math.Round(csmproj[i].CsmReleaseFairValueApproach)
	}
	return csmproj
}

func GetLiabilityMovementsProducts(runDate string) map[string]interface{} {
	var result = make(map[string]interface{})
	var products []models.LiabilityMovement
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Distinct("product_code").Where("run_date = ?", runDate).Find(&products).Error }); err != nil {
		fmt.Println(err)
	}
	result["list"] = products
	result["movements"] = GetLiabilityMovements(runDate, "", "")

	return result
}

func GetLiabilityMovementsProductGroups(runDate, productCode string) map[string]interface{} {
	var result = make(map[string]interface{})
	var products []models.LiabilityMovement
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Distinct("ifrs17_group").Where("product_code = ? and run_date=?", productCode, runDate).Find(&products).Error }); err != nil {
		fmt.Println(err)
	}
	result["list"] = products
	result["movements"] = GetLiabilityMovements(runDate, productCode, "")
	return result
}

func GetCSMRunList() []models.CsmRun {
	var runList []models.CsmRun
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Distinct("run_date").Find(&runList).Error }); err != nil {
		fmt.Println(err)
	}
	return runList
}

func GetCsmTableData(tableName string) map[string]interface{} {
	var results = make(map[string]interface{})
	fmt.Println("table name: ", tableName)
	switch tableName {
	case "financevariables":
		var financeVars []models.FinanceVariables
		err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Find(&financeVars).Error
		})
		if err != nil {
			fmt.Println(err)
		}
		results["table_data"] = financeVars
	case "rafactors":
		var raFactors []models.RiskAdjustmentFactor
		err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Find(&raFactors).Error
		})
		if err != nil {
			fmt.Println(err)
		}
		results["table_data"] = raFactors
	case "paafinance":
		var paaFinance []models.PAAFinance
		err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Find(&paaFinance).Error
		})
		if err != nil {
			fmt.Println(err)
		}
		results["table_data"] = paaFinance
	case "transitionadjustments":
		var tas []models.BalanceSheetRecord
		err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Where("measurement_type = ?", "TransitionAdjustment").Find(&tas).Error
		})
		if err != nil {
			fmt.Println(err)
		}
		results["table_data"] = tas
	}

	return results
}

func DeleteCsmTableData(tableName, year, version string) error {
	var results error
	fmt.Println("table name: ", tableName)
	switch tableName {
	case "finance_variables":
		results = DB.Where("year = ? and version = ?", year, version).Delete(&models.FinanceVariables{}).Error
	case "ra_factors":
		results = DB.Where("year = ? and version = ?", year, version).Delete(&models.RiskAdjustmentFactor{}).Error
	case "paa_finance":
		results = DB.Where("year = ? and version = ?", year, version).Delete(&models.PAAFinance{}).Error
	case "transition_adjustments":
		results = DB.Where("measurement_type = ? and year = ? and version = ?", "TransitionAdjustment", year, version).Delete(&models.BalanceSheetRecord{}).Error
	}
	return results
}

func GetAosRunList() []models.CsmRun {
	var runList []models.CsmRun
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Where("measurement_type = ? or measurement_type=?", "GMM", "VFA").Distinct("run_date", "measurement_type").Find(&runList).Error }); err != nil {
		fmt.Println(err)
	}
	return runList
}

func GetPAARunList(eligibilityMode string) []models.CsmRun {
	var runList []models.CsmRun
	if eligibilityMode == "inactive" {
		if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Where("measurement_type = ? and paa_eligibility_test = ?", "PAA", 0).Distinct("run_date", "measurement_type").Find(&runList).Error }); err != nil {
			fmt.Println(err)
		}
		return runList

	}

	if eligibilityMode == "active" {
		if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Where("measurement_type = ?", "PAA").Distinct("run_date", "measurement_type").Find(&runList).Error }); err != nil {
			fmt.Println(err)
		}
		return runList
	}
	return nil
}

type DisclosuresRunList struct {
	RunDate string `json:"run_date"`
	RunId   int    `json:"run_id"`
}

func GetCombinedRunListForDisclosures() []DisclosuresRunList {
	var runList []DisclosuresRunList

	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Model(&models.CsmRun{}).Where("measurement_type = ? and paa_eligibility_test = ?", "PAA", 0).Distinct("run_date", "measurement_type").Find(&runList).Error }); err != nil {
		fmt.Println(err)
	}

	var licRunList []models.LicRunSetting
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Model(&models.LicRunSetting{}).Distinct("run_date").Find(&licRunList).Error }); err != nil {
		fmt.Println(err)
	}
	for _, licRun := range licRunList {
		runList = append(runList, DisclosuresRunList{
			RunDate: licRun.RunDate,
			RunId:   licRun.ID,
		})
	}

	distinctRunList := GetDistinctValues(runList)
	return distinctRunList //runList
}

type DisclosuresPortfolioList struct {
	RunDate       string `json:"run_date"`
	PortfolioName string `json:"portfolio_name"`
	PortfolioType string `json:"portfolio_type"`
}

type DisclosuresPortfolioProductsList struct {
	RunDate       string `json:"run_date"`
	PortfolioName string `json:"portfolio_name"`
	ProductCode   string `json:"product_code"`
}

func GetPortfoliosForDisclosures(runDate string) map[string]interface{} {
	var portfolioList []DisclosuresPortfolioList

	var results = make(map[string]interface{})

	var list []models.PAABuildUp
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Distinct("portfolio_name").Where("run_date = ?", runDate).Find(&list).Error }); err != nil {
		fmt.Println(err)
	}
	// get portfolioList from paa and lic runs
	for _, item := range list {
		portfolioList = append(portfolioList, DisclosuresPortfolioList{
			RunDate:       runDate,
			PortfolioName: item.PortfolioName,
			PortfolioType: "PAA",
		})
	}

	// get unique portfolio name from lic build up resuts
	var licList []models.LicBuildupResult
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Distinct("portfolio_name").Where("run_date = ?", runDate).Find(&licList).Error }); err != nil {
		fmt.Println(err)
	}

	for _, item := range licList {
		portfolioList = append(portfolioList, DisclosuresPortfolioList{
			RunDate:       runDate,
			PortfolioName: item.PortfolioName,
			PortfolioType: "LIC",
		})
	}

	//DB.Where("run_date = ?", runDate).Find(&portfolioList)
	distinctPortfolioList := GetDistinctValues(portfolioList)

	results["list"] = distinctPortfolioList //portfolioList

	var paaBuildUps []models.PAABuildUp
	selectBuildupQuery := "name, sum(variable_change) as variable_change,sum(paa_lrc_buildup) as paa_lrc_buildup,sum(initial_recognition_loss_component) as initial_recognition_loss_component,sum(loss_component_unwind) as loss_component_unwind,sum(loss_component_buildup) as loss_component_buildup,sum(loss_component_adjustment) as loss_component_adjustment,sum(initial_recognition_loss_recovery) as initial_recognition_loss_recovery,sum(loss_recovery_unwind) as loss_recovery_unwind, sum(loss_recovery_buildup) as loss_recovery_buildup,sum(loss_recovery_adjustment) as loss_recovery_adjustment,sum(dac_buildup) as dac_buildup,sum(modified_gmm_bel) as modified_gmm_bel,sum(modified_gmm_risk_adjustment) as modified_gmm_risk_adjustment,sum(modified_gmm_reserve) as modified_gmm_reserve, sum(insurance_revenue) as insurance_revenue,sum(insurance_service_expense) as insurance_service_expense,sum(paa_reinsurance_lrc_buildup) as paa_reinsurance_lrc_buildup, sum(paa_reinsurance_dac_buildup) as paa_reinsurance_dac_buildup, sum(loss_recovery_buildup) as loss_recovery_buildup FROM paa_build_ups"
	query := fmt.Sprintf("SELECT %s where run_date = '%s'  group by name", selectBuildupQuery, runDate)
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw(query).Scan(&paaBuildUps).Error
	})
	if err != nil {
		fmt.Println(err)
	}

	if paaBuildUps == nil || len(paaBuildUps) == 0 {
		paaBuildUps = []models.PAABuildUp{}
	}

	results["paa_build_ups"] = paaBuildUps

	var paaResults []models.PAAResult
	selectQuery := "sum(insurance_revenue) as insurance_revenue,sum(incurred_expenses) as incurred_expenses, sum(incurred_claims) as incurred_claims,sum(claims_paid) as claims_paid,sum(paa_loss_component) as paa_loss_component,sum(reinsurance_premium) as reinsurance_premium,sum(reinsurance_flat_commission) as reinsurance_flat_commission,sum(reinsurance_claims) as reinsurance_claims,sum(reinsurance_recovery) as reinsurance_recovery,sum(reinsurance_reinstatement_premium) as reinsurance_reinstatement_premium,sum(reinsurance_provisional_commission) as reinsurance_provisional_commission, sum(reinsurance_ultimate_commission) as reinsurance_ultimate_commission,sum(reinsurance_profit_commission) as reinsurance_profit_commission,sum(reinsurance_total_paid_to_cedant) as reinsurance_total_paid_to_cedant,sum(reinsurance_investment_component) as reinsurance_investment_component,sum(reinsurance_revenue) as reinsurance_revenue,sum(reinsurance_service_expense) as reinsurance_service_expense, sum(reinsurance_service_result) as reinsurance_service_result,sum(paa_liability_remaining_coverage) as paa_liability_remaining_coverage,sum(allocated_reinsurance_premium) as allocated_reinsurance_premium, sum(allocated_reinsurance_flat_commission) as allocated_reinsurance_flat_commission, sum(paa_reinsurance_lrc) as paa_reinsurance_lrc, sum(paa_reinsurance_dac) as paa_reinsurance_dac FROM paa_results"
	query = fmt.Sprintf("SELECT %s where run_date = '%s'", selectQuery, runDate)
	err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw(query).Scan(&paaResults).Error
	})
	if err != nil {
		fmt.Println(err)
	}

	// if paaResults is nil or empty, initialize it to an empty slice
	if paaResults == nil || len(paaResults) == 0 {
		paaResults = []models.PAAResult{}
	}

	results["paa_results"] = paaResults

	// LIC Buildups

	var buildUp []models.LicBuildupResult

	err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("run_date = ?", runDate).Find(&buildUp).Error
	})

	if err != nil {
		log.Error("Error fetching LIC buildups: ", err)
	}

	if len(buildUp) == 0 || buildUp == nil {
		buildUp = []models.LicBuildupResult{}
	}

	results["lic_build_ups"] = buildUp

	return results
}

func GetPortfolioProductsForDisclosures(runDate, portfolioName string) map[string]interface{} {
	var productList []DisclosuresPortfolioProductsList
	var results = make(map[string]interface{})

	var list []models.PAABuildUp
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Distinct("product_code").Where("run_date = ? and portfolio_name = ? ", runDate, portfolioName).Find(&list).Error }); err != nil {
		fmt.Println(err)
	}
	// get portfolioList from paa and lic runs
	for _, item := range list {
		productList = append(productList, DisclosuresPortfolioProductsList{
			RunDate:       runDate,
			PortfolioName: item.PortfolioName,
			ProductCode:   item.ProductCode,
		})
	}

	// get unique portfolio name from lic build up resuts
	var licList []models.LicBuildupResult
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Distinct("product_code").Where("run_date = ? and portfolio_name = ?", runDate, portfolioName).Find(&licList).Error }); err != nil {
		fmt.Println(err)
	}

	for _, item := range licList {
		productList = append(productList, DisclosuresPortfolioProductsList{
			RunDate:       runDate,
			PortfolioName: item.PortfolioName,
			ProductCode:   item.ProductCode,
		})
	}

	// distinct values
	distinctProductList := GetDistinctValues(productList)
	//DB.Where("run_date = ?", runDate).Find(&portfolioList)
	results["list"] = distinctProductList //productList

	var paaBuildUps []models.PAABuildUp
	selectBuildupQuery := "name, sum(variable_change) as variable_change,sum(paa_lrc_buildup) as paa_lrc_buildup,sum(initial_recognition_loss_component) as initial_recognition_loss_component,sum(loss_component_unwind) as loss_component_unwind,sum(loss_component_buildup) as loss_component_buildup,sum(loss_component_adjustment) as loss_component_adjustment,sum(initial_recognition_loss_recovery) as initial_recognition_loss_recovery,sum(loss_recovery_unwind) as loss_recovery_unwind, sum(loss_recovery_buildup) as loss_recovery_buildup,sum(loss_recovery_adjustment) as loss_recovery_adjustment,sum(dac_buildup) as dac_buildup,sum(modified_gmm_bel) as modified_gmm_bel,sum(modified_gmm_risk_adjustment) as modified_gmm_risk_adjustment,sum(modified_gmm_reserve) as modified_gmm_reserve, sum(insurance_revenue) as insurance_revenue,sum(insurance_service_expense) as insurance_service_expense,sum(paa_reinsurance_lrc_buildup) as paa_reinsurance_lrc_buildup, sum(paa_reinsurance_dac_buildup) as paa_reinsurance_dac_buildup, sum(loss_recovery_buildup) as loss_recovery_buildup FROM paa_build_ups"
	query := fmt.Sprintf("SELECT %s where run_date = '%s' and portfolio_name = '%s'  group by name", selectBuildupQuery, runDate, portfolioName)
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw(query).Scan(&paaBuildUps).Error
	})
	if err != nil {
		fmt.Println(err)
	}

	if paaBuildUps == nil || len(paaBuildUps) == 0 {
		paaBuildUps = []models.PAABuildUp{}
	}

	results["paa_build_ups"] = paaBuildUps

	var paaResults []models.PAAResult
	selectQuery := "sum(insurance_revenue) as insurance_revenue,sum(incurred_expenses) as incurred_expenses, sum(incurred_claims) as incurred_claims,sum(claims_paid) as claims_paid,sum(paa_loss_component) as paa_loss_component,sum(reinsurance_premium) as reinsurance_premium,sum(reinsurance_flat_commission) as reinsurance_flat_commission,sum(reinsurance_claims) as reinsurance_claims,sum(reinsurance_recovery) as reinsurance_recovery,sum(reinsurance_reinstatement_premium) as reinsurance_reinstatement_premium,sum(reinsurance_provisional_commission) as reinsurance_provisional_commission, sum(reinsurance_ultimate_commission) as reinsurance_ultimate_commission,sum(reinsurance_profit_commission) as reinsurance_profit_commission,sum(reinsurance_total_paid_to_cedant) as reinsurance_total_paid_to_cedant,sum(reinsurance_investment_component) as reinsurance_investment_component,sum(reinsurance_revenue) as reinsurance_revenue,sum(reinsurance_service_expense) as reinsurance_service_expense, sum(reinsurance_service_result) as reinsurance_service_result,sum(paa_liability_remaining_coverage) as paa_liability_remaining_coverage,sum(allocated_reinsurance_premium) as allocated_reinsurance_premium, sum(allocated_reinsurance_flat_commission) as allocated_reinsurance_flat_commission, sum(paa_reinsurance_lrc) as paa_reinsurance_lrc, sum(paa_reinsurance_dac) as paa_reinsurance_dac FROM paa_results"
	query = fmt.Sprintf("SELECT %s where run_date = '%s' and portfolio_name = '%s'", selectQuery, runDate, portfolioName)
	err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw(query).Scan(&paaResults).Error
	})
	if err != nil {
		fmt.Println(err)
	}

	// if paaResults is nil or empty, initialize it to an empty slice
	if paaResults == nil || len(paaResults) == 0 {
		paaResults = []models.PAAResult{}
	}

	results["paa_results"] = paaResults

	// LIC Buildups

	var buildUp []models.LicBuildupResult

	err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("run_date = ? and portfolio_name = ?", runDate, portfolioName).Find(&buildUp).Error
	})

	if err != nil {
		log.Error("Error fetching LIC buildups: ", err)
	}

	if len(buildUp) == 0 || buildUp == nil {
		buildUp = []models.LicBuildupResult{}
	}

	results["lic_build_ups"] = buildUp

	return results
}

func GetPortfolioProductGroupsForDisclosures(runDate, portfolioName, productCode string) map[string]interface{} {
	var groupList []models.GroupResults
	var results = make(map[string]interface{})

	var list []models.PAABuildUp
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Distinct("ifrs17_group").Where("run_date = ? and portfolio_name = ? and product_code = ?", runDate, portfolioName, productCode).Find(&list).Error }); err != nil {
		fmt.Println(err)
	}
	// get portfolioList from paa and lic runs
	for _, item := range list {
		groupList = append(groupList, models.GroupResults{
			RunDate:       runDate,
			PortfolioName: portfolioName,
			ProductCode:   productCode,
			IFRS17Group:   item.IFRS17Group,
		})
	}

	// get unique portfolio name from lic build up resuts
	var licList []models.LicBuildupResult
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Distinct("ifrs17_group").Where("run_date = ? and portfolio_name = ? and product_code = ?", runDate, portfolioName, productCode).Find(&licList).Error }); err != nil {
		fmt.Println(err)
	}
	for _, item := range licList {
		groupList = append(groupList, models.GroupResults{
			RunDate:       runDate,
			PortfolioName: portfolioName,
			ProductCode:   productCode,
			IFRS17Group:   item.IFRS17Group,
		})
	}

	distinctGroupList := GetDistinctValues(groupList)

	results["list"] = distinctGroupList // groupList

	return results
}

func GetReportGroups(runDate, productCode string) []models.GroupResults {
	//So the report groups will be pulled from step results...
	var groups []models.GroupResults
	if productCode != "All Products" {
		query := fmt.Sprintf("select distinct ifrs17_group from journal_transactions where run_date = '%s' and product_code = '%s'", runDate, productCode)
		err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Raw(query).Scan(&groups).Error
		})
		if err != nil {
			fmt.Println(err)
		}
	}
	return groups
}

func GetCSMProductList(runDate string) []models.ProductList {
	var prodList []models.ProductList
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw("select distinct product_code from journal_transactions where run_date = '" + runDate + "'").Scan(&prodList).Error
	})
	if err != nil {
		fmt.Println(err)
	}

	//if measurementType == "PAA" {
	//	DB.Raw("select distinct product_code from paa_results where csm_run_id = " + runId).Scan(&prodList)
	//}

	return prodList
}

type ResultId struct {
	ID int
}

func GetAosProductList(runId, measurementType string) []models.ProductList {
	var prodList []models.ProductList
	var csmRun models.CsmRun

	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Where("run_date = ? and measurement_type = ?", runId, measurementType).Find(&csmRun).Error }); err != nil {
		fmt.Println(err)
	}

	if measurementType == "PAA" {
		err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Raw(fmt.Sprintf("select distinct product_code from paa_results where csm_run_id = %d", csmRun.ID)).Scan(&prodList).Error
		})
		if err != nil {
			fmt.Println(err)
		}
	}

	if measurementType == "GMM" {
		err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Raw(fmt.Sprintf("select distinct product_code from aos_step_results where run_id = %d", csmRun.ID)).Scan(&prodList).Error
		})
		if err != nil {
			fmt.Println(err)
		}
	}

	return prodList
}

func InterestAccretionFactor(yieldcurvestartmonth int, lockedinyear int, lockedinmonth int, IFStatus string, prodCode string, yieldCurveCode, yieldCurveBasis string) float64 {
	var resp float64 = 0
	var factor float64 = 1
	var nbmonthcount int = 1
	nbmonthcount = int(math.Max(float64(12-lockedinmonth+1), 1.0)) // assumes a policy is written at the beginning of a month
	if IFStatus == "IF" {
		for i := 1; i <= 12; i++ {
			resp, _ = GetForwardRateWithError(yieldcurvestartmonth+i, lockedinyear, lockedinmonth, yieldCurveCode)
			factor = factor * math.Pow(1+resp, 1.0/12.0)
		}
	} else {
		for i := 1; i <= nbmonthcount; i++ {
			resp, _ = GetForwardRateWithError(yieldcurvestartmonth+i, lockedinyear, lockedinmonth, yieldCurveCode)
			factor = factor * math.Pow(1+resp, 1.0/12.0)
		}
	}
	factor = factor - 1
	return factor
}

// CalculateCsm will house the workflow for IFRS17 calculations
func CalculateCsm(csmRun models.CsmRun, group models.GroupResults) ([]models.AOSStepResult, error) {
	var configs []models.AosVariable
	var jt models.JournalTransactions
	var stepResults []models.AOSStepResult
	var csmProjections []models.CsmProjection
	var aosconfigset models.AosVariableSet
	//var liailityMovements []models.LiabilityMovement

	var financeVariables = getFinancialVariables(group.ProductCode, group.IFRS17Group, csmRun.FinanceYear, csmRun.FinanceVersion)
	if financeVariables.ProductCode == "" {
		return nil, errors.New("no finance variables found for the ifrs17 group: " + group.IFRS17Group + " and product code: " + group.ProductCode)
	}

	//checking yield curve
	//var fwRate float64
	var err error
	_, err = GetForwardRateWithError(1, financeVariables.LockedInYear, financeVariables.LockedInMonth, financeVariables.YieldCurveCode)
	if err != nil {
		return nil, errors.New("no yield curve found for the ifrs17 group: " + group.IFRS17Group + "; lockedin year: " + strconv.Itoa(financeVariables.LockedInYear) + " and yield curve month: " + strconv.Itoa(financeVariables.LockedInMonth) + "  and yield curve code: " + financeVariables.YieldCurveCode)
	}

	jt.IFRS17Group = group.IFRS17Group
	jt.ProductCode = group.ProductCode
	jt.CsmRunID = csmRun.ID
	jt.MeasurementType = csmRun.MeasurementType
	jt.RunDate = csmRun.RunDate
	err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Find(&configs).Error
	})
	if err != nil {
		fmt.Println(err)
	}
	err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Preload("AosVariables").Where("configuration_name = ?", csmRun.ConfigurationName).Find(&aosconfigset).Error
	})
	if err != nil {
		fmt.Println(err)
	}
	var expectedCashInflow float64
	var expectedCashOutflow float64
	var ReinsExpectedCashInflow float64
	var ReinsExpectedCashOutflow float64
	var riskAdjYear int
	var riskAdjVersion string
	//var riskAdjustmentAccretion float64
	var currYear, currMonth, prevYear, prevMonth int
	currYear, err = strconv.Atoi(csmRun.RunDate[:4])
	if err != nil {
		currYear = 0
	}
	prevYear, err = strconv.Atoi(csmRun.OpeningBalanceDate[:4])
	if err != nil {
		prevYear = 0
	}
	currMonth, err = strconv.Atoi(csmRun.RunDate[5:])
	if err != nil {
		currMonth = 0
	}
	prevMonth, err = strconv.Atoi(csmRun.OpeningBalanceDate[5:])
	if err != nil {
		prevMonth = 0
	}

	if financeVariables.IFStatus == IF {
		if aosconfigset.AosVariables[0].RunId > 0 {
			var riskAdjustment = getRiskAdjustmentFactors(group.ProductCode, prevYear, csmRun.OpeningRiskAdjustmentVersion)
			if riskAdjustment.ProductCode == "" {
				return nil, errors.New("no risk adjustment variables found for the ifrs17 group: " + group.IFRS17Group + " ,product code: " + group.ProductCode + ",year: " + strconv.Itoa(csmRun.RiskAdjustmentYear-1) + " and version: " + csmRun.OpeningRiskAdjustmentVersion + " - last year's ra-factors are required")
			}
		}
	}

	if group.IFRS17Group == "" {
		return nil, errors.New("missing ifrs17group for the product code: " + group.ProductCode)
	}

	for i, item := range aosconfigset.AosVariables {
		var result models.AOSStepResult
		if !csmRun.ManualSap {
			if item.RunId == 0 {
				result.ID = 0
				result.RunDate = csmRun.RunDate[:7]
				result.CsmRunID = csmRun.ID
				result.ProductCode = group.ProductCode
				result.IFRS17Group = group.IFRS17Group
				result.Name = item.Name
				result.Time = item.Time
				result.Description = item.Description
				if i > 0 {
					result.CSMBuildup = stepResults[i-1].CSMBuildup
					result.LossComponentBuildup = stepResults[i-1].LossComponentBuildup
					result.DACBuildup = 0 //stepResults[i-1].DACBuildup
					result.RiskAdjustmentBuildup = stepResults[i-1].RiskAdjustmentBuildup
				}
				stepResults = append(stepResults, result)
				continue
			}
		} else {
			if item.RunName == "" {
				result.ID = 0
				result.RunDate = csmRun.RunDate[:7]
				result.CsmRunID = csmRun.ID
				result.ProductCode = group.ProductCode
				result.IFRS17Group = group.IFRS17Group
				result.Name = item.Name
				result.Time = item.Time
				result.Description = item.Description
				if i > 0 {
					result.CSMBuildup = stepResults[i-1].CSMBuildup
					result.LossComponentBuildup = stepResults[i-1].LossComponentBuildup
					result.DACBuildup = 0 //stepResults[i-1].DACBuildup
					result.RiskAdjustmentBuildup = stepResults[i-1].RiskAdjustmentBuildup
				}
				stepResults = append(stepResults, result)
				continue
			}
		}
		var job models.JobProduct
		var sap []models.ScopedAggregatedProjection

		if !csmRun.ManualSap {
			err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
				return d.Where("projection_job_id=? and product_code=?", item.RunId, group.ProductCode).Find(&job).Error
			})
			err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
				return d.Where("job_product_id=? and ifrs17_group=? and projection_month < ?", job.ID, group.IFRS17Group, MaxProjectionMonthSap).Order("projection_month asc").Find(&sap).Error
			})
			if err != nil {
				log.Error(err)
			}
		} else {
			//Manual SAP
			err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
				return d.Table("manual_scoped_aggregated_projections").Where("ifrs17_group = ? and product_code = ? and run_name=?", group.IFRS17Group, group.ProductCode, item.RunName).Find(&sap).Error
			})
			if err != nil {
				fmt.Println(err)
			}
		}

		jt.ProductCode = group.ProductCode
		var riskAdjustmentFactors models.RiskAdjustmentFactor
		var riskDrivers models.RiskDriver
		if !csmRun.ManualSap {
			riskDrivers = getRiskDriversFor(group.ProductCode)
			if i <= 8 {
				if financeVariables.IFStatus == NB {
					riskAdjYear = csmRun.RiskAdjustmentYear
					riskAdjVersion = csmRun.RiskAdjustmentVersion
				}
				if financeVariables.IFStatus == IF {
					riskAdjYear, err = strconv.Atoi(csmRun.OpeningBalanceDate[:4])
					//riskAdjYear = csmRun.RiskAdjustmentYear - 1
					riskAdjVersion = csmRun.OpeningRiskAdjustmentVersion
				}
			}
			if i > 8 {
				riskAdjYear = csmRun.RiskAdjustmentYear
				riskAdjVersion = csmRun.RiskAdjustmentVersion
			}
			//riskAdjustmentFactors = getRiskAdjustmentFactors(group.ProductCode, csmRun.RiskAdjustmentYear)
			riskAdjustmentFactors = getRiskAdjustmentFactors(group.ProductCode, riskAdjYear, riskAdjVersion)
			if riskAdjustmentFactors.ProductCode == "" {
				return nil, errors.New("no risk adjustment variables found for the product code: " + group.ProductCode + "and year: " + strconv.Itoa(riskAdjYear) + " and version: " + riskAdjVersion)
			}
		}

		err = copier.Copy(&result, item)
		if err != nil {
			fmt.Println(err)
		}
		result.ID = 0
		result.RunDate = csmRun.RunDate[:7]
		result.ProductCode = group.ProductCode
		result.IFRS17Group = group.IFRS17Group

		var projectionTimeat12 int

		if financeVariables.IFStatus == IF { // currently coded for annual calculations.
			//projectionTimeat12 = 12
			var temp float64
			temp = utils.FloatPrecision((float64(currYear)+float64(currMonth)/1000.0)-(float64(prevYear)+float64(prevMonth)/1000.0), 3) //(float64(currYear) + float64(currMonth/1000.0)) - (float64(prevYear) + float64(prevMonth/1000.0))
			if temp == 1 {
				projectionTimeat12 = 12
			} else {
				if temp > 12 {
					projectionTimeat12 = int(math.Max(math.Min(12-(1000.0-math.Mod(temp, 1.0)*1000.0), 12), 0))
				}
				if temp <= 12 {
					projectionTimeat12 = int(math.Max(math.Min(temp*1000, 12), 0))
				}
			}
		} else {
			projectionTimeat12 = int(math.Min(float64(financeVariables.DurationInForceMonths), 12.0))
		}
		//if len(sap) > 0 {
		if len(sap) > 0 {
			if !csmRun.ManualSap {
				result.BelInflow = sap[0].DiscountedCashInflow
				result.BelOutflow = sap[0].DiscountedCashOutflow
				result.BelInflowAt12 = sap[projectionTimeat12].DiscountedCashInflow
				result.BelOutflowAt12 = sap[projectionTimeat12].DiscountedCashOutflow
				result.BelOutflowExclAcquisition = sap[0].DiscountedCashOutflowExclAcquisition
				result.BelOutflowExclAcquisitionAt12 = sap[projectionTimeat12].DiscountedCashOutflowExclAcquisition
				result.BEL = sap[0].Reserves
				result.BELAt12 = sap[projectionTimeat12].Reserves
				result.SumCoverageUnits = sap[0].SumCoverageUnits
				result.DiscountedCoverageUnits = sap[0].DiscountedCoverageUnits
				result.SumCoverageUnitsAt12 = sap[projectionTimeat12].SumCoverageUnits
				result.DiscountedCoverageUnitsAt12 = sap[projectionTimeat12].DiscountedCoverageUnits
				result.BelAcquisitionCost = sap[0].DiscountedAcquisitionCost
				result.BelAcquisitionCostAt12 = sap[projectionTimeat12].DiscountedAcquisitionCost

				result.RiskAdjustment = math.Abs(getFieldValue(&sap[0], riskDrivers.MortalityRisk)*riskAdjustmentFactors.MortalityRisk) +
					math.Abs(getFieldValue(&sap[0], riskDrivers.ExpenseRisk)*riskAdjustmentFactors.ExpenseRisk) +
					math.Abs(getFieldValue(&sap[0], riskDrivers.MorbidityRisk)*riskAdjustmentFactors.MorbidityRisk) +
					math.Abs(getFieldValue(&sap[0], riskDrivers.LapseRisk)*riskAdjustmentFactors.LapseRisk) +
					math.Abs(getFieldValue(&sap[0], riskDrivers.LongevityRisk)*riskAdjustmentFactors.LongevityRisk) +
					math.Abs(getFieldValue(&sap[0], riskDrivers.Catastrophe)*riskAdjustmentFactors.Catastrophe) +
					math.Abs(getFieldValue(&sap[0], riskDrivers.Operational)*riskAdjustmentFactors.Operational)

				result.RiskAdjustmentAt12 = math.Abs(getFieldValue(&sap[projectionTimeat12], riskDrivers.MortalityRisk)*riskAdjustmentFactors.MortalityRisk) +
					math.Abs(getFieldValue(&sap[projectionTimeat12], riskDrivers.ExpenseRisk)*riskAdjustmentFactors.ExpenseRisk) +
					math.Abs(getFieldValue(&sap[projectionTimeat12], riskDrivers.MorbidityRisk)*riskAdjustmentFactors.MorbidityRisk) +
					math.Abs(getFieldValue(&sap[projectionTimeat12], riskDrivers.LapseRisk)*riskAdjustmentFactors.LapseRisk) +
					math.Abs(getFieldValue(&sap[projectionTimeat12], riskDrivers.LongevityRisk)*riskAdjustmentFactors.LongevityRisk) +
					math.Abs(getFieldValue(&sap[projectionTimeat12], riskDrivers.Catastrophe)*riskAdjustmentFactors.Catastrophe) +
					math.Abs(getFieldValue(&sap[projectionTimeat12], riskDrivers.Operational)*riskAdjustmentFactors.Operational)
			} else {
				result.BelInflow = sap[0].DiscountedCashInflowM0 + sap[0].DiscountedCashInflowsNotVaryM0 + sap[0].DiscountedEntityShareM0
				result.BelOutflow = sap[0].DiscountedCashOutflowM0 + sap[0].TvogM0
				result.BelInflowAt12 = sap[0].DiscountedCashInflowM12 + sap[0].DiscountedCashInflowsNotVaryM12 + sap[0].DiscountedEntityShareM12
				result.BelOutflowAt12 = sap[0].DiscountedCashOutflowM12 + sap[0].TvogM12
				result.BelOutflowExclAcquisition = sap[0].DiscountedCashOutflowExclAcquisitionM0 + sap[0].TvogM0
				result.BelOutflowExclAcquisitionAt12 = sap[0].DiscountedCashOutflowExclAcquisitionM12 + sap[0].TvogM12
				result.BEL = sap[0].BelM0 + sap[0].TvogM0
				result.BELAt12 = sap[0].BelM12 + sap[0].TvogM12
				result.SumCoverageUnits = sap[0].SumCoverageUnitsM0
				result.DiscountedCoverageUnits = sap[0].DiscountedCoverageUnitsM0
				result.SumCoverageUnitsAt12 = sap[0].SumCoverageUnitsM12
				result.DiscountedCoverageUnitsAt12 = sap[0].DiscountedCoverageUnitsM12
				result.RiskAdjustment = sap[0].RiskAdjustmentM0
				result.RiskAdjustmentAt12 = sap[0].RiskAdjustmentM12
				result.BelAcquisitionCost = sap[0].DiscountedAcquisitionCostM0
				result.BelAcquisitionCostAt12 = sap[0].DiscountedAcquisitionCostM12
				result.ReinsuranceBelInflow = sap[0].Treaty1DiscountedCashInflowM0 + sap[0].Treaty2DiscountedCashInflowM0 + sap[0].Treaty3DiscountedCashInflowM0
				result.ReinsuranceBelOutflow = sap[0].Treaty1DiscountedCashOutflowM0 + sap[0].Treaty2DiscountedCashOutflowM0 + sap[0].Treaty3DiscountedCashOutflowM0
				result.ReinsuranceBel = sap[0].Treaty1BelM0 + sap[0].Treaty2BelM0 + sap[0].Treaty3BelM0
				result.ReinsuranceRiskAdjustment = sap[0].Treaty1RiskAdjustmentM0 + sap[0].Treaty2RiskAdjustmentM0 + sap[0].Treaty3RiskAdjustmentM0
				result.ReinsuranceBelInflowAt12 = sap[0].Treaty1DiscountedCashInflowM12 + sap[0].Treaty2DiscountedCashInflowM12 + sap[0].Treaty3DiscountedCashInflowM12
				result.ReinsuranceBelOutflowAt12 = sap[0].Treaty1DiscountedCashOutflowM12 + sap[0].Treaty2DiscountedCashOutflowM12 + sap[0].Treaty3DiscountedCashOutflowM12
				result.ReinsuranceBelAt12 = sap[0].Treaty1BelM12 + sap[0].Treaty2BelM12 + sap[0].Treaty3BelM12
				result.ReinsuranceRiskAdjustmentAt12 = sap[0].Treaty1RiskAdjustmentM12 + sap[0].Treaty2RiskAdjustmentM12 + sap[0].Treaty3RiskAdjustmentM12
			}
		}

		var currentPeriodExpectedCashOutflows float64
		var pvCashOutflows float64
		//var reinsuranceCurrentPeriodCashOutflows float64
		//var reinsurancePVCashOutflows float64
		if i == 3 {
			currentPeriodExpectedCashOutflows = result.BelOutflow + result.RiskAdjustment - result.BelOutflowAt12 - result.RiskAdjustmentAt12
			pvCashOutflows = result.BelOutflow + result.RiskAdjustment
			//reinsuranceCurrentPeriodCashOutflows = result.ReinsuranceBelOutflow - result.ReinsuranceBelOutflowAt12
			//reinsurancePVCashOutflows = result.ReinsuranceBelOutflow + result.ReinsuranceRiskAdjustment
		}
		if i > 3 {
			currentPeriodExpectedCashOutflows = stepResults[3].BelOutflow + stepResults[3].RiskAdjustment - stepResults[3].BelOutflowAt12 - stepResults[3].RiskAdjustmentAt12
			pvCashOutflows = stepResults[3].BelOutflow + stepResults[3].RiskAdjustment
			//reinsuranceCurrentPeriodCashOutflows = stepResults[3].ReinsuranceBelOutflow - stepResults[3].ReinsuranceBelOutflowAt12
			//reinsurancePVCashOutflows = stepResults[3].ReinsuranceBelOutflow + stepResults[3].ReinsuranceRiskAdjustment

		}

		if pvCashOutflows > 0 {
			result.LcSar = currentPeriodExpectedCashOutflows / pvCashOutflows //stepResults[2].LossComponentBuildup / (stepResults[0].BEL + stepResults[2].BEL + stepResults[0].RiskAdjustment + stepResults[2].RiskAdjustment
		} else {
			result.LcSar = 0
		}

		if item.Name == "B/F_Current" || item.Name == "B/F_Lockedin" {
			result.CSMChange = 0
			result.LossComponentChange = 0

			if i == 0 {
				result.PNLChange = 0
				result.BestEstimateLiabilityChange = 0
				result.RiskAdjustmentChange = 0

				//prevYearStr := csmRun.RunDate[:4]
				//prevYearMth := csmRun.RunDate[5:]
				//prevYear, _ := strconv.Atoi(prevYearStr)
				//prevYear = prevYear - 1
				//yearSearchArg := strconv.Itoa(prevYear) + "-" + prevYearMth
				yearSearchArg := csmRun.OpeningBalanceDate
				//fmt.Println(yearSearchArg)
				var prevBsr models.BalanceSheetRecord
				err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
					return d.Where("product_code = ? and ifrs17_group = ? and date = ? and measurement_type = ?", result.ProductCode, result.IFRS17Group, yearSearchArg, csmRun.MeasurementType).Find(&prevBsr).Error
				})
				if err != nil {
					fmt.Println("query bsr error:", err)
				}

				var transitionAdjustment models.BalanceSheetRecord
				err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
					return d.Where("measurement_type = ? and ifrs17_group = ? and date = ?", "TransitionAdjustment", result.IFRS17Group, yearSearchArg).Find(&transitionAdjustment).Error
				})
				if err != nil {
					fmt.Println("query bsr error:", err)
				}

				switch csmRun.TransitionType {
				case PostTransition:
					result.CSMBuildup = prevBsr.PostTransitionCsm + transitionAdjustment.PostTransitionCsm         //CSM_CarriedFoward
					result.LossComponentBuildup = prevBsr.PostTransitionLc + transitionAdjustment.PostTransitionLc //LossComponent_Carriedforward'
					result.DACBuildup = 0                                                                          //prevBsr.PostTransitionDAC + transitionAdjustment.PostTransitionDAC
					result.ReinsuranceCSM = prevBsr.PostTransitionTreaty1Csm + prevBsr.PostTransitionTreaty2Csm + prevBsr.PostTransitionTreaty3Csm + transitionAdjustment.PostTransitionTreaty1Csm + transitionAdjustment.PostTransitionTreaty2Csm + transitionAdjustment.PostTransitionTreaty3Csm
					result.ReinsCSMBuildup = result.ReinsuranceCSM
					result.LossRecoveryComponentBuildup = prevBsr.PostTransitionLossRecoveryComponent + transitionAdjustment.PostTransitionLossRecoveryComponent
				case FullyRetrospective:
					result.CSMBuildup = prevBsr.FullyRetrospectiveCsm + transitionAdjustment.PostTransitionCsm         //CSM_CarriedFoward
					result.LossComponentBuildup = prevBsr.FullyRetrospectiveLc + transitionAdjustment.PostTransitionLc //LossComponent_Carriedforward'
					result.DACBuildup = 0                                                                              //prevBsr.FullyRetrospectiveDAC + transitionAdjustment.PostTransitionDAC
					result.ReinsuranceCSM = prevBsr.FullyRetrospectiveTreaty1Csm + prevBsr.FullyRetrospectiveTreaty2Csm + prevBsr.FullyRetrospectiveTreaty3Csm + transitionAdjustment.PostTransitionTreaty1Csm + transitionAdjustment.PostTransitionTreaty2Csm + transitionAdjustment.PostTransitionTreaty3Csm
					result.ReinsCSMBuildup = result.ReinsuranceCSM
					result.LossRecoveryComponentBuildup = prevBsr.FRALossRecoveryComponent + transitionAdjustment.PostTransitionLossRecoveryComponent
				case ModifiedRetrospective:
					result.CSMBuildup = prevBsr.ModifiedRetrospectiveCsm + transitionAdjustment.PostTransitionCsm         //CSM_CarriedFoward
					result.LossComponentBuildup = prevBsr.ModifiedRetrospectiveLc + transitionAdjustment.PostTransitionLc //LossComponent_Carriedforward'
					result.DACBuildup = 0                                                                                 //prevBsr.ModifiedRetrospectiveDAC + transitionAdjustment.PostTransitionDAC
					result.ReinsuranceCSM = prevBsr.ModifiedRetrospectiveTreaty1Csm + prevBsr.ModifiedRetrospectiveTreaty2Csm + prevBsr.ModifiedRetrospectiveTreaty3Csm + transitionAdjustment.PostTransitionTreaty1Csm + transitionAdjustment.PostTransitionTreaty2Csm + transitionAdjustment.PostTransitionTreaty3Csm
					result.ReinsCSMBuildup = result.ReinsuranceCSM
					result.LossRecoveryComponentBuildup = prevBsr.MRALossRecoveryComponent + transitionAdjustment.PostTransitionLossRecoveryComponent
				case FairValue:
					result.CSMBuildup = prevBsr.FairValueCsm + transitionAdjustment.PostTransitionCsm         //CSM_CarriedFoward
					result.LossComponentBuildup = prevBsr.FairValueLc + transitionAdjustment.PostTransitionLc //LossComponent_Carriedforward'
					result.DACBuildup = 0                                                                     //prevBsr.FairValueDAC + transitionAdjustment.PostTransitionDAC
					result.ReinsuranceCSM = prevBsr.FairValueTreaty1Csm + prevBsr.FairValueTreaty2Csm + prevBsr.FairValueTreaty3Csm + transitionAdjustment.PostTransitionTreaty1Csm + transitionAdjustment.PostTransitionTreaty2Csm + transitionAdjustment.PostTransitionTreaty3Csm
					result.ReinsCSMBuildup = result.ReinsuranceCSM
					result.LossRecoveryComponentBuildup = prevBsr.FVALossRecoveryComponent + transitionAdjustment.PostTransitionLossRecoveryComponent
				}

			} else { //financeincome or finance expense
				result.BestEstimateLiabilityChange = result.BEL - stepResults[i-1].BEL
				result.BelInflowChange = result.BelInflow - stepResults[i-1].BelInflow    //stepResults[i-1].BelOutflowChange
				result.DACChange = 0                                                      //result.BelAcquisitionCost - stepResults[i-1].BelAcquisitionCost
				result.BelOutflowChange = result.BelOutflow - stepResults[i-1].BelOutflow //stepResults[i-1].BelOutflowChange
				result.RiskAdjustmentChange = result.RiskAdjustment - stepResults[i-1].RiskAdjustment
				result.LiabilityChange = result.BestEstimateLiabilityChange + result.RiskAdjustmentChange
				result.PNLChange = -result.BestEstimateLiabilityChange - result.RiskAdjustmentChange
				result.CSMBuildup = stepResults[i-1].CSMBuildup
				result.DACBuildup = 0 //stepResults[i-1].DACBuildup //+ result.DACChange
				result.LossComponentUnwind = 0
				result.LossComponentBuildup = stepResults[i-1].LossComponentBuildup
				result.RiskAdjustmentBuildup = stepResults[i-1].RiskAdjustmentBuildup + result.RiskAdjustmentChange
				result.ReinsBelOutflowChange = result.ReinsuranceBelOutflow - stepResults[i-1].ReinsuranceBelOutflow
				result.ReinsBelInflowChange = result.ReinsuranceBelInflow - stepResults[i-1].ReinsuranceBelInflow
				result.ReinsRAChange = result.ReinsuranceRiskAdjustment - stepResults[i-1].ReinsuranceRiskAdjustment
				result.ReinsBelChange = result.ReinsuranceBel - stepResults[i-1].ReinsuranceBel
				result.ReinsPNLChange = result.ReinsBelChange + result.ReinsRAChange
				result.ReinsCSMBuildup = stepResults[i-1].ReinsCSMChange
				result.LossRecoveryComponentBuildup = stepResults[i-1].LossComponentBuildup
			}
		}

		if item.Name == "Interest_Accretion" {
			var interestaccretionfac float64
			var durationstart int
			//var valuationyear int
			//valuationyear = getCurrentYear()
			var coverage int

			if financeVariables.IFStatus == IF { //valuationyear > financeVariables.LockedInYear {
				var temp float64
				temp = utils.FloatPrecision((float64(currYear)+float64(currMonth)/1000.0)-(float64(prevYear)+float64(prevMonth)/1000.0), 3) //(float64(currYear) + float64(currMonth/1000.0)) - (float64(prevYear) + float64(prevMonth/1000.0))
				if temp == 1 {
					coverage = 12
				} else {
					if temp > 12 {
						coverage = int(math.Max(math.Min(12-(1000.0-math.Mod(temp, 1.0)*1000.0), 12), 0))
					}
					if temp <= 12 {
						coverage = int(math.Max(math.Min(temp*1000, 12), 0))
					}
				}
				durationstart = int(math.Max(float64(financeVariables.DurationInForceMonths-coverage), 0)) //(valuationyear-financeVariables.LockedInYear)*12 - financeVariables.LockedInMonth + 1
			} else {
				durationstart = 0
			}

			if csmRun.MeasurementType == "GMM" {
				interestaccretionfac = InterestAccretionFactor(durationstart, financeVariables.LockedInYear, financeVariables.LockedInMonth, financeVariables.IFStatus, group.ProductCode, financeVariables.YieldCurveCode, financeVariables.YieldCurveBasis)
				result.ExpectedCashOutflow = expectedCashOutflow
				result.ExpectedCashInflow = expectedCashInflow
				result.PremiumDebtor = financeVariables.PremiumDebtors
				result.ReinsExpectedCashInflow = ReinsExpectedCashInflow
				result.ReinsExpectedCashOutflow = ReinsExpectedCashOutflow

				result.ActualPremium = financeVariables.ActualPremiumIncome
				result.ExperiencePremiumVariance = financeVariables.ActualPremiumIncome - (expectedCashInflow - financeVariables.PremiumDebtors)
				result.BelInflowChange = expectedCashInflow - (result.BelInflow - result.BelInflowAt12)
				result.BelOutflowChange = expectedCashOutflow - (result.BelOutflow - result.BelOutflowAt12)
				result.BestEstimateLiabilityChange = result.BelOutflowChange - result.BelInflowChange
				result.RiskAdjustmentChange = stepResults[1].RiskAdjustment * interestaccretionfac //will need to revist-- the run must include both b/f and new business
				result.LiabilityChange = result.BestEstimateLiabilityChange + result.RiskAdjustmentChange
				if stepResults[i-1].LossComponentBuildup > 0 {
					result.LossComponentUnwind = -math.Max(expectedCashOutflow+-stepResults[3].RiskAdjustmentChange-expectedCashInflow-result.LiabilityChange, 0) // riskadjustment change is negative hence and liabilityChange is interest rate on

				}
				result.CSMChange = stepResults[2].CSMBuildup * interestaccretionfac
				result.DACChange = 0 //stepResults[2].DACBuildup * interestaccretionfac

				//Reinsurance
				result.ReinsCSMChange = stepResults[2].ReinsCSMBuildup * interestaccretionfac
				result.ReinsBelInflowChange = ReinsExpectedCashInflow - (stepResults[4].ReinsuranceBelInflow - stepResults[4].ReinsuranceBelInflowAt12)
				result.ReinsBelOutflowChange = ReinsExpectedCashOutflow - (stepResults[4].ReinsuranceBelOutflow - stepResults[4].ReinsuranceBelOutflowAt12)
				result.ReinsBelChange = result.ReinsBelInflowChange - result.ReinsBelOutflowChange
				result.ReinsRAChange = 0
				result.ReinsPNLChange = result.ReinsBelChange + result.ReinsRAChange

				if stepResults[i-1].LossComponentBuildup > 0 {
					result.LossComponentChange = stepResults[i-1].LossComponentBuildup * interestaccretionfac
					result.LossRecoveryComponentChange = result.ReinsBelChange
				} else {
					result.LossComponentChange = 0
					result.LossRecoveryComponentChange = 0
				}
				result.PNLChange = -(result.BestEstimateLiabilityChange + result.RiskAdjustmentChange + result.DACChange)
				result.CSMBuildup = stepResults[i-1].CSMBuildup + result.CSMChange
				result.DACBuildup = stepResults[i-1].DACBuildup
				result.LossComponentBuildup = stepResults[i-1].LossComponentBuildup + result.LossComponentChange + result.LossComponentUnwind
				result.InterestAccretionFac = interestaccretionfac
				result.RiskAdjustmentBuildup = stepResults[i-1].RiskAdjustmentBuildup + result.RiskAdjustmentChange

				//reinsurance
				result.LossRecoveryComponentBuildup = stepResults[i-1].LossRecoveryComponentBuildup + result.LossRecoveryComponentChange
				result.ReinsCSMBuildup = stepResults[i-1].ReinsCSMBuildup + result.ReinsCSMChange

			}
			if csmRun.MeasurementType == "VFA" {
				interestaccretionfac = InterestAccretionFactor(durationstart, financeVariables.LockedInYear, financeVariables.LockedInMonth, financeVariables.IFStatus, group.ProductCode, financeVariables.YieldCurveCode, financeVariables.YieldCurveBasis)
				result.ExpectedCashOutflow = expectedCashOutflow
				result.ExpectedCashInflow = expectedCashInflow
				result.PremiumDebtor = financeVariables.PremiumDebtors
				result.ReinsExpectedCashInflow = ReinsExpectedCashInflow
				result.ReinsExpectedCashOutflow = ReinsExpectedCashOutflow

				result.ActualPremium = financeVariables.ActualPremiumIncome
				result.ExperiencePremiumVariance = financeVariables.ActualPremiumIncome - (expectedCashInflow - financeVariables.PremiumDebtors)
				result.BelInflowChange = expectedCashInflow - (stepResults[4].BelInflow - stepResults[4].BelInflowAt12)
				result.BelOutflowChange = expectedCashOutflow - (stepResults[4].BelOutflow - stepResults[4].BelOutflowAt12)
				result.BestEstimateLiabilityChange = result.BelOutflowChange - result.BelInflowChange
				result.RiskAdjustmentChange = stepResults[1].RiskAdjustment * interestaccretionfac //will need to revist-- the run must include both b/f and new business
				result.LiabilityChange = result.BestEstimateLiabilityChange + result.RiskAdjustmentChange
				if stepResults[i-1].LossComponentBuildup > 0 {
					result.LossComponentUnwind = -math.Max(expectedCashOutflow+-stepResults[3].RiskAdjustmentChange-expectedCashInflow-result.LiabilityChange, 0) // riskadjustment change is negative hence

				}
				result.CSMChange = result.BestEstimateLiabilityChange*(1-financeVariables.RiskMitigationProportion) + stepResults[2].CSMBuildup*interestaccretionfac
				result.DACChange = 0 //stepResults[2].DACBuildup * interestaccretionfac

				//Reinsurance
				result.ReinsCSMChange = -(result.ReinsBelChange+result.ReinsRAChange)*(1-financeVariables.RiskMitigationProportion) + stepResults[2].ReinsCSMBuildup*interestaccretionfac
				result.ReinsBelInflowChange = ReinsExpectedCashInflow - (stepResults[4].ReinsuranceBelInflow - stepResults[4].ReinsuranceBelInflowAt12)
				result.ReinsBelOutflowChange = ReinsExpectedCashOutflow - (stepResults[4].ReinsuranceBelOutflow - stepResults[4].ReinsuranceBelOutflowAt12)
				result.ReinsBelChange = result.ReinsBelInflowChange - result.ReinsBelOutflowChange
				result.ReinsRAChange = 0
				result.ReinsPNLChange = -(result.ReinsBelChange + result.ReinsRAChange) * financeVariables.RiskMitigationProportion

				if stepResults[i-1].LossComponentBuildup > 0 {
					result.LossComponentChange = stepResults[i-1].LossComponentBuildup * interestaccretionfac
					result.LossRecoveryComponentChange = result.ReinsBelChange
				} else {
					result.LossComponentChange = 0
					result.LossRecoveryComponentChange = 0
				}
				result.PNLChange = (result.BestEstimateLiabilityChange + result.RiskAdjustmentChange + result.DACChange) * financeVariables.RiskMitigationProportion
				result.CSMBuildup = stepResults[i-1].CSMBuildup + result.CSMChange
				result.LossComponentBuildup = stepResults[i-1].LossComponentBuildup + result.LossComponentChange + result.LossComponentUnwind
				result.InterestAccretionFac = interestaccretionfac
				result.RiskAdjustmentBuildup = stepResults[i-1].RiskAdjustmentBuildup + result.RiskAdjustmentChange

				//reinsurance
				result.LossRecoveryComponentBuildup = stepResults[i-1].LossRecoveryComponentBuildup + result.LossRecoveryComponentChange
				result.ReinsCSMBuildup = stepResults[i-1].ReinsCSMBuildup + result.ReinsCSMChange

			}
		}

		if item.Name == "Initial_Recog" {
			result.BestEstimateLiabilityChange = result.BEL
			result.BelInflowChange = result.BelInflow
			result.BelOutflowChange = result.BelOutflow
			result.RiskAdjustmentChange = result.RiskAdjustment
			result.LiabilityChange = result.BestEstimateLiabilityChange + result.RiskAdjustmentChange
			result.CSMChange = math.Max(-(result.BEL + result.RiskAdjustment), 0)
			result.LossComponentChange = -math.Min(-(result.BEL + result.RiskAdjustment), 0)
			result.PNLChange = result.LossComponentChange
			result.CSMBuildup = stepResults[i-1].CSMBuildup + result.CSMChange
			result.LossComponentBuildup = stepResults[i-1].LossComponentBuildup + result.LossComponentChange
			result.DACChange = 0 //result.BelAcquisitionCost
			result.DACBuildup = stepResults[i-1].DACBuildup + result.DACChange
			result.RiskAdjustmentBuildup = stepResults[i-1].RiskAdjustmentBuildup + result.RiskAdjustmentChange

			result.ReinsBelInflowChange = result.ReinsuranceBelInflow
			result.ReinsBelOutflowChange = result.ReinsuranceBelOutflow
			result.ReinsRAChange = result.ReinsuranceRiskAdjustment
			result.ReinsBelChange = result.ReinsuranceBel
			result.ReinsCSMChange = -(result.ReinsuranceBel + result.ReinsuranceRiskAdjustment)
			result.ReinsCSMBuildup = stepResults[i-1].ReinsCSMBuildup + result.ReinsCSMChange
			if (result.BelOutflow + result.RiskAdjustment) > 0 {
				result.LossRecoveryComponentChange = result.LossComponentChange * (result.ReinsuranceBelOutflow + result.ReinsuranceRiskAdjustment) / (result.BelOutflow + result.RiskAdjustment)

			}
			result.LossRecoveryComponentBuildup = stepResults[i-1].LossRecoveryComponentBuildup + result.LossRecoveryComponentChange

			if result.LossComponentChange != 0 {
				if result.BelOutflow+result.RiskAdjustment > 0 {
					result.BelOutflowExclLC = result.BelInflowChange * (result.BelOutflow / (result.BelOutflow + result.RiskAdjustment))
					result.RiskAdjExclLC = result.BelInflowChange * (1.0 - result.BelOutflow/(result.BelOutflow+result.RiskAdjustment))
					result.BelOutflowLC = result.BestEstimateLiabilityChange
					result.RiskAdjLC = result.RiskAdjustmentChange
				}
			}
			if result.LossComponentChange == 0 {
				result.BelOutflowExclLC = result.BelOutflowChange
				result.RiskAdjExclLC = result.RiskAdjustmentChange
				result.BelOutflowLC = 0
				result.RiskAdjLC = 0
			}

		}

		if item.Name == "Exp_Mort" {
			var expectedClaims float64
			var actualClaim = financeVariables.ActualMortalityClaimsIncurred
			var coverage int

			if financeVariables.IFStatus == IF {
				//coverage = 12
				var temp float64
				temp = utils.FloatPrecision((float64(currYear)+float64(currMonth)/1000.0)-(float64(prevYear)+float64(prevMonth)/1000.0), 3) //(float64(currYear) + float64(currMonth/1000.0)) - (float64(prevYear) + float64(prevMonth/1000.0))
				if temp == 1 {
					coverage = 12
				} else {
					if temp > 12 {
						coverage = int(math.Max(math.Min(12-(1000.0-math.Mod(temp, 1.0)*1000.0), 12), 0))
					}
					if temp <= 12 {
						coverage = int(math.Max(math.Min(temp*1000, 12), 0))
					}
				}
			} else {
				coverage = int(math.Min(float64(financeVariables.DurationInForceMonths), 12.0))
			}

			if len(sap) > 0 {
				if !csmRun.ManualSap {
					for i := 1; i <= coverage; i++ { //year end month
						//if len(sap) > 0 {
						expectedClaims += sap[i].DeathOutgo + sap[i].AccidentalDeathOutgo
						expectedCashInflow += sap[i].PremiumIncome
						expectedCashOutflow += sap[i].NetCashFlow + sap[i].PremiumIncome // computing sum of cash outflows
						//}
					}
				} else {
					switch coverage {
					case 1:
						expectedClaims = sap[0].MortalityOutgoM1
						expectedCashInflow = sap[0].CashInflowM1 + sap[0].EntityShareM1 + sap[0].CashInflowsNotVaryM1

					case 2:
						expectedClaims = sap[0].MortalityOutgoM1 + sap[0].MortalityOutgoM2
						expectedCashInflow = sap[0].CashInflowM1 + sap[0].EntityShareM1 + sap[0].CashInflowsNotVaryM1 + sap[0].CashInflowM2 + sap[0].EntityShareM2 + sap[0].CashInflowsNotVaryM2
					case 3:
						expectedClaims = sap[0].MortalityOutgoM1 + sap[0].MortalityOutgoM2 + sap[0].MortalityOutgoM3
						expectedCashInflow = sap[0].CashInflowM1 + sap[0].EntityShareM1 + sap[0].CashInflowsNotVaryM1 + sap[0].CashInflowM2 + sap[0].EntityShareM2 +
							sap[0].CashInflowM3 + sap[0].EntityShareM3 + sap[0].CashInflowsNotVaryM3
					case 4:
						expectedClaims = sap[0].MortalityOutgoM1 + sap[0].MortalityOutgoM2 + sap[0].MortalityOutgoM3 + sap[0].MortalityOutgoM4
						expectedCashInflow = sap[0].CashInflowM1 + sap[0].EntityShareM1 + sap[0].CashInflowsNotVaryM1 + sap[0].CashInflowM2 + sap[0].EntityShareM2 +
							sap[0].CashInflowM3 + sap[0].EntityShareM3 + sap[0].CashInflowsNotVaryM3 + sap[0].CashInflowM4 + sap[0].EntityShareM4 + sap[0].CashInflowsNotVaryM4
					case 5:
						expectedClaims = sap[0].MortalityOutgoM1 + sap[0].MortalityOutgoM2 + sap[0].MortalityOutgoM3 + sap[0].MortalityOutgoM4 + sap[0].MortalityOutgoM5
						expectedCashInflow = sap[0].CashInflowM1 + sap[0].EntityShareM1 + sap[0].CashInflowsNotVaryM1 + sap[0].CashInflowM2 + sap[0].EntityShareM2 +
							sap[0].CashInflowM3 + sap[0].EntityShareM3 + sap[0].CashInflowsNotVaryM3 + sap[0].CashInflowM4 + sap[0].EntityShareM4 + sap[0].CashInflowsNotVaryM4 +
							sap[0].CashInflowM5 + sap[0].EntityShareM5 + sap[0].CashInflowsNotVaryM5
					case 6:
						expectedClaims = sap[0].MortalityOutgoM1 + sap[0].MortalityOutgoM2 + sap[0].MortalityOutgoM3 + sap[0].MortalityOutgoM4 + sap[0].MortalityOutgoM5 + sap[0].MortalityOutgoM6
						expectedCashInflow = sap[0].CashInflowM1 + sap[0].EntityShareM1 + sap[0].CashInflowsNotVaryM1 + sap[0].CashInflowM2 + sap[0].EntityShareM2 +
							sap[0].CashInflowM3 + sap[0].EntityShareM3 + sap[0].CashInflowsNotVaryM3 + sap[0].CashInflowM4 + sap[0].EntityShareM4 + sap[0].CashInflowsNotVaryM4 +
							sap[0].CashInflowM5 + sap[0].EntityShareM5 + sap[0].CashInflowsNotVaryM5 + sap[0].CashInflowM6 + sap[0].EntityShareM6 + sap[0].CashInflowsNotVaryM6
					case 7:
						expectedClaims = sap[0].MortalityOutgoM1 + sap[0].MortalityOutgoM2 + sap[0].MortalityOutgoM3 + sap[0].MortalityOutgoM4 + sap[0].MortalityOutgoM5 + sap[0].MortalityOutgoM6 + sap[0].MortalityOutgoM7
						expectedCashInflow = sap[0].CashInflowM1 + sap[0].EntityShareM1 + sap[0].CashInflowsNotVaryM1 + sap[0].CashInflowM2 + sap[0].EntityShareM2 +
							sap[0].CashInflowM3 + sap[0].EntityShareM3 + sap[0].CashInflowsNotVaryM3 + sap[0].CashInflowM4 + sap[0].EntityShareM4 + sap[0].CashInflowsNotVaryM4 +
							sap[0].CashInflowM5 + sap[0].EntityShareM5 + sap[0].CashInflowsNotVaryM5 + sap[0].CashInflowM6 + sap[0].EntityShareM6 + sap[0].CashInflowsNotVaryM6 +
							sap[0].CashInflowM7 + sap[0].EntityShareM7 + sap[0].CashInflowsNotVaryM7
					case 8:
						expectedClaims = sap[0].MortalityOutgoM1 + sap[0].MortalityOutgoM2 + sap[0].MortalityOutgoM3 + sap[0].MortalityOutgoM4 + sap[0].MortalityOutgoM5 + sap[0].MortalityOutgoM6 + sap[0].MortalityOutgoM7 + sap[0].MortalityOutgoM8
						expectedCashInflow = sap[0].CashInflowM1 + sap[0].EntityShareM1 + sap[0].CashInflowsNotVaryM1 + sap[0].CashInflowM2 + sap[0].EntityShareM2 +
							sap[0].CashInflowM3 + sap[0].EntityShareM3 + sap[0].CashInflowsNotVaryM3 + sap[0].CashInflowM4 + sap[0].EntityShareM4 + sap[0].CashInflowsNotVaryM4 +
							sap[0].CashInflowM5 + sap[0].EntityShareM5 + sap[0].CashInflowsNotVaryM5 + sap[0].CashInflowM6 + sap[0].EntityShareM6 + sap[0].CashInflowsNotVaryM6 +
							sap[0].CashInflowM7 + sap[0].EntityShareM7 + sap[0].CashInflowsNotVaryM7 + sap[0].CashInflowM8 + sap[0].EntityShareM8 + sap[0].CashInflowsNotVaryM8
					case 9:
						expectedClaims = sap[0].MortalityOutgoM1 + sap[0].MortalityOutgoM2 + sap[0].MortalityOutgoM3 + sap[0].MortalityOutgoM4 + sap[0].MortalityOutgoM5 + sap[0].MortalityOutgoM6 + sap[0].MortalityOutgoM7 + sap[0].MortalityOutgoM8 + sap[0].MortalityOutgoM9
						expectedCashInflow = sap[0].CashInflowM1 + sap[0].EntityShareM1 + sap[0].CashInflowsNotVaryM1 + sap[0].CashInflowM2 + sap[0].EntityShareM2 +
							sap[0].CashInflowM3 + sap[0].EntityShareM3 + sap[0].CashInflowsNotVaryM3 + sap[0].CashInflowM4 + sap[0].EntityShareM4 + sap[0].CashInflowsNotVaryM4 +
							sap[0].CashInflowM5 + sap[0].EntityShareM5 + sap[0].CashInflowsNotVaryM5 + sap[0].CashInflowM6 + sap[0].EntityShareM6 + sap[0].CashInflowsNotVaryM6 +
							sap[0].CashInflowM7 + sap[0].EntityShareM7 + sap[0].CashInflowsNotVaryM7 + sap[0].CashInflowM8 + sap[0].EntityShareM8 + sap[0].CashInflowsNotVaryM8 +
							sap[0].CashInflowM9 + sap[0].EntityShareM9 + sap[0].CashInflowsNotVaryM9
					case 10:
						expectedClaims = sap[0].MortalityOutgoM1 + sap[0].MortalityOutgoM2 + sap[0].MortalityOutgoM3 + sap[0].MortalityOutgoM4 + sap[0].MortalityOutgoM5 + sap[0].MortalityOutgoM6 + sap[0].MortalityOutgoM7 + sap[0].MortalityOutgoM8 + sap[0].MortalityOutgoM9 + sap[0].MortalityOutgoM10
						expectedCashInflow = sap[0].CashInflowM1 + sap[0].EntityShareM1 + sap[0].CashInflowsNotVaryM1 + sap[0].CashInflowM2 + sap[0].EntityShareM2 +
							sap[0].CashInflowM3 + sap[0].EntityShareM3 + sap[0].CashInflowsNotVaryM3 + sap[0].CashInflowM4 + sap[0].EntityShareM4 + sap[0].CashInflowsNotVaryM4 +
							sap[0].CashInflowM5 + sap[0].EntityShareM5 + sap[0].CashInflowsNotVaryM5 + sap[0].CashInflowM6 + sap[0].EntityShareM6 + sap[0].CashInflowsNotVaryM6 +
							sap[0].CashInflowM7 + sap[0].EntityShareM7 + sap[0].CashInflowsNotVaryM7 + sap[0].CashInflowM8 + sap[0].EntityShareM8 + sap[0].CashInflowsNotVaryM8 +
							sap[0].CashInflowM9 + sap[0].EntityShareM9 + sap[0].CashInflowsNotVaryM9 + sap[0].CashInflowM10 + sap[0].EntityShareM10 + sap[0].CashInflowsNotVaryM10
					case 11:
						expectedClaims = sap[0].MortalityOutgoM1 + sap[0].MortalityOutgoM2 + sap[0].MortalityOutgoM3 + sap[0].MortalityOutgoM4 + sap[0].MortalityOutgoM5 + sap[0].MortalityOutgoM6 + sap[0].MortalityOutgoM7 + sap[0].MortalityOutgoM8 + sap[0].MortalityOutgoM9 + sap[0].MortalityOutgoM10 + sap[0].MortalityOutgoM11
						expectedCashInflow = sap[0].CashInflowM1 + sap[0].EntityShareM1 + sap[0].CashInflowsNotVaryM1 + sap[0].CashInflowM2 + sap[0].EntityShareM2 +
							sap[0].CashInflowM3 + sap[0].EntityShareM3 + sap[0].CashInflowsNotVaryM3 + sap[0].CashInflowM4 + sap[0].EntityShareM4 + sap[0].CashInflowsNotVaryM4 +
							sap[0].CashInflowM5 + sap[0].EntityShareM5 + sap[0].CashInflowsNotVaryM5 + sap[0].CashInflowM6 + sap[0].EntityShareM6 + sap[0].CashInflowsNotVaryM6 +
							sap[0].CashInflowM7 + sap[0].EntityShareM7 + sap[0].CashInflowsNotVaryM7 + sap[0].CashInflowM8 + sap[0].EntityShareM8 + sap[0].CashInflowsNotVaryM8 +
							sap[0].CashInflowM9 + sap[0].EntityShareM9 + sap[0].CashInflowsNotVaryM9 + sap[0].CashInflowM10 + sap[0].EntityShareM10 + sap[0].CashInflowsNotVaryM10 +
							sap[0].CashInflowM11 + sap[0].EntityShareM11 + sap[0].CashInflowsNotVaryM11
					case 12:
						expectedClaims = sap[0].MortalityOutgoM1 + sap[0].MortalityOutgoM2 + sap[0].MortalityOutgoM3 + sap[0].MortalityOutgoM4 + sap[0].MortalityOutgoM5 + sap[0].MortalityOutgoM6 + sap[0].MortalityOutgoM7 + sap[0].MortalityOutgoM8 + sap[0].MortalityOutgoM9 + sap[0].MortalityOutgoM10 + sap[0].MortalityOutgoM11 + sap[0].MortalityOutgoM12
						expectedCashInflow = sap[0].CashInflowM1 + sap[0].EntityShareM1 + sap[0].CashInflowsNotVaryM1 + sap[0].CashInflowM2 + sap[0].EntityShareM2 +
							sap[0].CashInflowM3 + sap[0].EntityShareM3 + sap[0].CashInflowsNotVaryM3 + sap[0].CashInflowM4 + sap[0].EntityShareM4 + sap[0].CashInflowsNotVaryM4 +
							sap[0].CashInflowM5 + sap[0].EntityShareM5 + sap[0].CashInflowsNotVaryM5 + sap[0].CashInflowM6 + sap[0].EntityShareM6 + sap[0].CashInflowsNotVaryM6 +
							sap[0].CashInflowM7 + sap[0].EntityShareM7 + sap[0].CashInflowsNotVaryM7 + sap[0].CashInflowM8 + sap[0].EntityShareM8 + sap[0].CashInflowsNotVaryM8 +
							sap[0].CashInflowM9 + sap[0].EntityShareM9 + sap[0].CashInflowsNotVaryM9 + sap[0].CashInflowM10 + sap[0].EntityShareM10 + sap[0].CashInflowsNotVaryM10 +
							sap[0].CashInflowM11 + sap[0].EntityShareM11 + sap[0].CashInflowsNotVaryM11 + sap[0].CashInflowM12 + sap[0].EntityShareM12 + sap[0].CashInflowsNotVaryM12
					default:
						expectedClaims = 0
						expectedCashInflow = 0
					}
					expectedCashOutflow += expectedClaims
				}
			}

			result.BestEstimateLiabilityChange = -expectedClaims
			result.BelInflowChange = 0
			result.BelOutflowChange = -expectedClaims
			result.RiskAdjustmentChange = 0
			result.CSMChange = 0
			result.LossComponentChange = 0 //-(expectedClaims * result.LcSar)
			result.LiabilityChange = -expectedClaims
			result.PNLChange = expectedClaims        //* (1 - result.LcSar)
			result.SarActualClaimNetLc = actualClaim //- expectedClaims*result.LcSar
			result.CSMBuildup = stepResults[i-1].CSMBuildup + result.CSMChange
			result.LossComponentBuildup = stepResults[i-1].LossComponentBuildup + result.LossComponentChange
			result.DACBuildup = stepResults[i-1].DACBuildup
			result.RiskAdjustmentBuildup = stepResults[i-1].RiskAdjustmentBuildup + result.RiskAdjustmentChange

			//reinsurance
			result.LossRecoveryComponentBuildup = stepResults[i-1].LossRecoveryComponentBuildup + result.LossRecoveryComponentChange
			result.ReinsCSMBuildup = stepResults[i-1].ReinsCSMBuildup + result.ReinsCSMChange

		}

		if item.Name == "Exp_Retrenchment" {
			var expectedClaims float64
			var actualRetrenchmentClaim = financeVariables.ActualRetrenchmentClaimsIncurred
			var coverage int

			if financeVariables.IFStatus == IF {
				//coverage = 12
				var temp float64
				temp = utils.FloatPrecision((float64(currYear)+float64(currMonth)/1000.0)-(float64(prevYear)+float64(prevMonth)/1000.0), 3) //(float64(currYear) + float64(currMonth/1000.0)) - (float64(prevYear) + float64(prevMonth/1000.0))
				if temp == 1 {
					coverage = 12
				} else {
					if temp > 12 {
						coverage = int(math.Max(math.Min(12-(1000.0-math.Mod(temp, 1.0)*1000.0), 12), 0))
					}
					if temp <= 12 {
						coverage = int(math.Max(math.Min(temp*1000, 12), 0))
					}
				}
			} else {
				coverage = int(math.Min(float64(financeVariables.DurationInForceMonths), 12.0))
			}

			if len(sap) > 0 {
				if !csmRun.ManualSap {
					for i := 1; i <= coverage; i++ {
						expectedClaims += sap[i].RetrenchmentOutgo
					}
				} else {
					switch coverage {
					case 1:
						expectedClaims = sap[0].RetrenchmentOutgoM1
					case 2:
						expectedClaims = sap[0].RetrenchmentOutgoM1 + sap[0].RetrenchmentOutgoM2
					case 3:
						expectedClaims = sap[0].RetrenchmentOutgoM1 + sap[0].RetrenchmentOutgoM2 + sap[0].RetrenchmentOutgoM3
					case 4:
						expectedClaims = sap[0].RetrenchmentOutgoM1 + sap[0].RetrenchmentOutgoM2 + sap[0].RetrenchmentOutgoM3 + sap[0].RetrenchmentOutgoM4
					case 5:
						expectedClaims = sap[0].RetrenchmentOutgoM1 + sap[0].RetrenchmentOutgoM2 + sap[0].RetrenchmentOutgoM3 + sap[0].RetrenchmentOutgoM4 + sap[0].RetrenchmentOutgoM5
					case 6:
						expectedClaims = sap[0].RetrenchmentOutgoM1 + sap[0].RetrenchmentOutgoM2 + sap[0].RetrenchmentOutgoM3 + sap[0].RetrenchmentOutgoM4 + sap[0].RetrenchmentOutgoM5 + sap[0].RetrenchmentOutgoM6
					case 7:
						expectedClaims = sap[0].RetrenchmentOutgoM1 + sap[0].RetrenchmentOutgoM2 + sap[0].RetrenchmentOutgoM3 + sap[0].RetrenchmentOutgoM4 + sap[0].RetrenchmentOutgoM5 + sap[0].RetrenchmentOutgoM6 + sap[0].RetrenchmentOutgoM7
					case 8:
						expectedClaims = sap[0].RetrenchmentOutgoM1 + sap[0].RetrenchmentOutgoM2 + sap[0].RetrenchmentOutgoM3 + sap[0].RetrenchmentOutgoM4 + sap[0].RetrenchmentOutgoM5 + sap[0].RetrenchmentOutgoM6 + sap[0].RetrenchmentOutgoM7 + sap[0].RetrenchmentOutgoM8
					case 9:
						expectedClaims = sap[0].RetrenchmentOutgoM1 + sap[0].RetrenchmentOutgoM2 + sap[0].RetrenchmentOutgoM3 + sap[0].RetrenchmentOutgoM4 + sap[0].RetrenchmentOutgoM5 + sap[0].RetrenchmentOutgoM6 + sap[0].RetrenchmentOutgoM7 + sap[0].RetrenchmentOutgoM8 + sap[0].RetrenchmentOutgoM9
					case 10:
						expectedClaims = sap[0].RetrenchmentOutgoM1 + sap[0].RetrenchmentOutgoM2 + sap[0].RetrenchmentOutgoM3 + sap[0].RetrenchmentOutgoM4 + sap[0].RetrenchmentOutgoM5 + sap[0].RetrenchmentOutgoM6 + sap[0].RetrenchmentOutgoM7 + sap[0].RetrenchmentOutgoM8 + sap[0].RetrenchmentOutgoM9 + sap[0].RetrenchmentOutgoM10
					case 11:
						expectedClaims = sap[0].RetrenchmentOutgoM1 + sap[0].RetrenchmentOutgoM2 + sap[0].RetrenchmentOutgoM3 + sap[0].RetrenchmentOutgoM4 + sap[0].RetrenchmentOutgoM5 + sap[0].RetrenchmentOutgoM6 + sap[0].RetrenchmentOutgoM7 + sap[0].RetrenchmentOutgoM8 + sap[0].RetrenchmentOutgoM9 + sap[0].RetrenchmentOutgoM10 + sap[0].RetrenchmentOutgoM11
					case 12:
						expectedClaims = sap[0].RetrenchmentOutgoM1 + sap[0].RetrenchmentOutgoM2 + sap[0].RetrenchmentOutgoM3 + sap[0].RetrenchmentOutgoM4 + sap[0].RetrenchmentOutgoM5 + sap[0].RetrenchmentOutgoM6 + sap[0].RetrenchmentOutgoM7 + sap[0].RetrenchmentOutgoM8 + sap[0].RetrenchmentOutgoM9 + sap[0].RetrenchmentOutgoM10 + sap[0].RetrenchmentOutgoM11 + sap[0].RetrenchmentOutgoM12
					default:
						expectedClaims = 0
					}
					expectedCashOutflow += expectedClaims
				}
			}

			result.BestEstimateLiabilityChange = -expectedClaims
			result.BelInflowChange = 0
			result.BelOutflowChange = -expectedClaims
			result.RiskAdjustmentChange = 0
			result.CSMChange = 0
			result.LossComponentChange = 0 //-(expectedClaims * result.LcSar)
			result.LiabilityChange = -expectedClaims
			result.PNLChange = expectedClaims                    //* (1 - result.LcSar)
			result.SarActualClaimNetLc = actualRetrenchmentClaim //- expectedClaims*result.LcSar
			result.CSMBuildup = stepResults[i-1].CSMBuildup + result.CSMChange
			result.LossComponentBuildup = stepResults[i-1].LossComponentBuildup + result.LossComponentChange
			result.DACBuildup = stepResults[i-1].DACBuildup
			result.RiskAdjustmentBuildup = stepResults[i-1].RiskAdjustmentBuildup + result.RiskAdjustmentChange

			//reinsurance
			result.LossRecoveryComponentBuildup = stepResults[i-1].LossRecoveryComponentBuildup + result.LossRecoveryComponentChange
			result.ReinsCSMBuildup = stepResults[i-1].ReinsCSMBuildup + result.ReinsCSMChange

		}

		if item.Name == "Exp_Morbidity" {
			var expectedClaims float64
			var actualMorbidityClaim = financeVariables.ActualMorbidityClaimsIncurred
			var actualNonLifeClaims = financeVariables.ActualNonLifeClaimsIncurred
			var coverage int

			if financeVariables.IFStatus == IF {
				//coverage = 12
				var temp float64
				temp = utils.FloatPrecision((float64(currYear)+float64(currMonth)/1000.0)-(float64(prevYear)+float64(prevMonth)/1000.0), 3) //(float64(currYear) + float64(currMonth/1000.0)) - (float64(prevYear) + float64(prevMonth/1000.0))
				if temp == 1 {
					coverage = 12
				} else {
					if temp > 12 {
						coverage = int(math.Max(math.Min(12-(1000.0-math.Mod(temp, 1.0)*1000.0), 12), 0))
					}
					if temp <= 12 {
						coverage = int(math.Max(math.Min(temp*1000, 12), 0))
					}
				}
			} else {
				coverage = int(math.Min(float64(financeVariables.DurationInForceMonths), 12.0))
			}

			if len(sap) > 0 {
				if !csmRun.ManualSap {
					for i := 1; i <= coverage; i++ {
						expectedClaims += sap[i].DisabilityOutgo + sap[i].NonLifeClaimsOutgo
					}
				} else {
					switch coverage {
					case 1:
						expectedClaims = sap[0].MorbidityOutgoM1 + sap[0].NonLifeClaimsOutgoM1

						ReinsExpectedCashOutflow = sap[0].Treaty1CashOutflowM1 +
							sap[0].Treaty2CashOutflowM1 +
							sap[0].Treaty3CashOutflowM1

						ReinsExpectedCashInflow = sap[0].Treaty1CashInflowM1 +
							sap[0].Treaty2CashInflowM1 +
							sap[0].Treaty3CashInflowM1

					case 2:
						expectedClaims = sap[0].MorbidityOutgoM1 + sap[0].MorbidityOutgoM2 + sap[0].NonLifeClaimsOutgoM1 + sap[0].NonLifeClaimsOutgoM2

						ReinsExpectedCashOutflow = sap[0].Treaty1CashOutflowM1 + sap[0].Treaty1CashOutflowM2 +
							sap[0].Treaty2CashOutflowM1 + sap[0].Treaty2CashOutflowM2 +
							sap[0].Treaty3CashOutflowM1 + sap[0].Treaty3CashOutflowM2

						ReinsExpectedCashInflow = sap[0].Treaty1CashInflowM1 + sap[0].Treaty1CashInflowM2 +
							sap[0].Treaty2CashInflowM1 + sap[0].Treaty2CashInflowM2 +
							sap[0].Treaty3CashInflowM1 + sap[0].Treaty3CashInflowM2

					case 3:
						expectedClaims = sap[0].MorbidityOutgoM1 + sap[0].MorbidityOutgoM2 + sap[0].MorbidityOutgoM3 + sap[0].NonLifeClaimsOutgoM1 + sap[0].NonLifeClaimsOutgoM2 + sap[0].NonLifeClaimsOutgoM3

						ReinsExpectedCashOutflow = sap[0].Treaty1CashOutflowM1 + sap[0].Treaty1CashOutflowM2 + sap[0].Treaty1CashOutflowM3 +
							sap[0].Treaty2CashOutflowM1 + sap[0].Treaty2CashOutflowM2 + sap[0].Treaty2CashOutflowM3 +
							sap[0].Treaty3CashOutflowM1 + sap[0].Treaty3CashOutflowM2 + sap[0].Treaty3CashOutflowM3

						ReinsExpectedCashInflow = sap[0].Treaty1CashInflowM1 + sap[0].Treaty1CashInflowM2 + sap[0].Treaty1CashInflowM3 +
							sap[0].Treaty2CashInflowM1 + sap[0].Treaty2CashInflowM2 + sap[0].Treaty2CashInflowM3 +
							sap[0].Treaty3CashInflowM1 + sap[0].Treaty3CashInflowM2 + sap[0].Treaty3CashInflowM3

					case 4:
						expectedClaims = sap[0].MorbidityOutgoM1 + sap[0].MorbidityOutgoM2 + sap[0].MorbidityOutgoM3 + sap[0].MorbidityOutgoM4 + sap[0].NonLifeClaimsOutgoM1 + sap[0].NonLifeClaimsOutgoM2 + sap[0].NonLifeClaimsOutgoM3 + sap[0].NonLifeClaimsOutgoM4

						ReinsExpectedCashOutflow = sap[0].Treaty1CashOutflowM1 + sap[0].Treaty1CashOutflowM2 + sap[0].Treaty1CashOutflowM3 + sap[0].Treaty1CashOutflowM4 +
							sap[0].Treaty2CashOutflowM1 + sap[0].Treaty2CashOutflowM2 + sap[0].Treaty2CashOutflowM3 + sap[0].Treaty2CashOutflowM4 +
							sap[0].Treaty3CashOutflowM1 + sap[0].Treaty3CashOutflowM2 + sap[0].Treaty3CashOutflowM3 + sap[0].Treaty3CashOutflowM4

						ReinsExpectedCashInflow = sap[0].Treaty1CashInflowM1 + sap[0].Treaty1CashInflowM2 + sap[0].Treaty1CashInflowM3 + sap[0].Treaty1CashInflowM4 +
							sap[0].Treaty2CashInflowM1 + sap[0].Treaty2CashInflowM2 + sap[0].Treaty2CashInflowM3 + sap[0].Treaty2CashInflowM4 +
							sap[0].Treaty3CashInflowM1 + sap[0].Treaty3CashInflowM2 + sap[0].Treaty3CashInflowM3 + sap[0].Treaty3CashInflowM4

					case 5:
						expectedClaims = sap[0].MorbidityOutgoM1 + sap[0].MorbidityOutgoM2 + sap[0].MorbidityOutgoM3 + sap[0].MorbidityOutgoM4 + sap[0].MorbidityOutgoM5 + sap[0].NonLifeClaimsOutgoM1 + sap[0].NonLifeClaimsOutgoM2 + sap[0].NonLifeClaimsOutgoM3 + sap[0].NonLifeClaimsOutgoM4 + sap[0].NonLifeClaimsOutgoM5

						ReinsExpectedCashOutflow = sap[0].Treaty1CashOutflowM1 + sap[0].Treaty1CashOutflowM2 + sap[0].Treaty1CashOutflowM3 + sap[0].Treaty1CashOutflowM4 + sap[0].Treaty1CashOutflowM5 +
							sap[0].Treaty2CashOutflowM1 + sap[0].Treaty2CashOutflowM2 + sap[0].Treaty2CashOutflowM3 + sap[0].Treaty2CashOutflowM4 + sap[0].Treaty2CashOutflowM4 +
							sap[0].Treaty3CashOutflowM1 + sap[0].Treaty3CashOutflowM2 + sap[0].Treaty3CashOutflowM3 + sap[0].Treaty3CashOutflowM4 + sap[0].Treaty3CashOutflowM4

						ReinsExpectedCashInflow = sap[0].Treaty1CashInflowM1 + sap[0].Treaty1CashInflowM2 + sap[0].Treaty1CashInflowM3 + sap[0].Treaty1CashInflowM4 + sap[0].Treaty1CashInflowM5 +
							sap[0].Treaty2CashInflowM1 + sap[0].Treaty2CashInflowM2 + sap[0].Treaty2CashInflowM3 + sap[0].Treaty2CashInflowM4 + sap[0].Treaty2CashInflowM4 +
							sap[0].Treaty3CashInflowM1 + sap[0].Treaty3CashInflowM2 + sap[0].Treaty3CashInflowM3 + sap[0].Treaty3CashInflowM4 + sap[0].Treaty3CashInflowM4

					case 6:
						expectedClaims = sap[0].MorbidityOutgoM1 + sap[0].MorbidityOutgoM2 + sap[0].MorbidityOutgoM3 + sap[0].MorbidityOutgoM4 + sap[0].MorbidityOutgoM5 + sap[0].MorbidityOutgoM6 + sap[0].NonLifeClaimsOutgoM1 + sap[0].NonLifeClaimsOutgoM2 + sap[0].NonLifeClaimsOutgoM3 + sap[0].NonLifeClaimsOutgoM4 + sap[0].NonLifeClaimsOutgoM5 + sap[0].NonLifeClaimsOutgoM6

						ReinsExpectedCashOutflow = sap[0].Treaty1CashOutflowM1 + sap[0].Treaty1CashOutflowM2 + sap[0].Treaty1CashOutflowM3 + sap[0].Treaty1CashOutflowM4 + sap[0].Treaty1CashOutflowM5 + sap[0].Treaty1CashOutflowM6 +
							sap[0].Treaty2CashOutflowM1 + sap[0].Treaty2CashOutflowM2 + sap[0].Treaty2CashOutflowM3 + sap[0].Treaty2CashOutflowM4 + sap[0].Treaty2CashOutflowM5 + sap[0].Treaty2CashOutflowM6 +
							sap[0].Treaty3CashOutflowM1 + sap[0].Treaty3CashOutflowM2 + sap[0].Treaty3CashOutflowM3 + sap[0].Treaty3CashOutflowM4 + sap[0].Treaty3CashOutflowM5 + sap[0].Treaty3CashOutflowM6

						ReinsExpectedCashInflow = sap[0].Treaty1CashInflowM1 + sap[0].Treaty1CashInflowM2 + sap[0].Treaty1CashInflowM3 + sap[0].Treaty1CashInflowM4 + sap[0].Treaty1CashInflowM5 + sap[0].Treaty1CashInflowM6 +
							sap[0].Treaty2CashInflowM1 + sap[0].Treaty2CashInflowM2 + sap[0].Treaty2CashInflowM3 + sap[0].Treaty2CashInflowM4 + sap[0].Treaty2CashInflowM5 + sap[0].Treaty2CashInflowM6 +
							sap[0].Treaty3CashInflowM1 + sap[0].Treaty3CashInflowM2 + sap[0].Treaty3CashInflowM3 + sap[0].Treaty3CashInflowM4 + sap[0].Treaty3CashInflowM5 + sap[0].Treaty3CashInflowM6

					case 7:
						expectedClaims = sap[0].MorbidityOutgoM1 + sap[0].MorbidityOutgoM2 + sap[0].MorbidityOutgoM3 + sap[0].MorbidityOutgoM4 + sap[0].MorbidityOutgoM5 + sap[0].MorbidityOutgoM6 + sap[0].MorbidityOutgoM7 + sap[0].NonLifeClaimsOutgoM1 + sap[0].NonLifeClaimsOutgoM2 + sap[0].NonLifeClaimsOutgoM3 + sap[0].NonLifeClaimsOutgoM4 + sap[0].NonLifeClaimsOutgoM5 + sap[0].NonLifeClaimsOutgoM6 + sap[0].NonLifeClaimsOutgoM7

						ReinsExpectedCashOutflow = sap[0].Treaty1CashOutflowM1 + sap[0].Treaty1CashOutflowM2 + sap[0].Treaty1CashOutflowM3 + sap[0].Treaty1CashOutflowM4 + sap[0].Treaty1CashOutflowM5 + sap[0].Treaty1CashOutflowM6 + sap[0].Treaty1CashOutflowM7 +
							sap[0].Treaty2CashOutflowM1 + sap[0].Treaty2CashOutflowM2 + sap[0].Treaty2CashOutflowM3 + sap[0].Treaty2CashOutflowM4 + sap[0].Treaty2CashOutflowM5 + sap[0].Treaty2CashOutflowM6 + sap[0].Treaty2CashOutflowM7 +
							sap[0].Treaty3CashOutflowM1 + sap[0].Treaty3CashOutflowM2 + sap[0].Treaty3CashOutflowM3 + sap[0].Treaty3CashOutflowM4 + sap[0].Treaty3CashOutflowM5 + sap[0].Treaty3CashOutflowM6 + sap[0].Treaty3CashOutflowM7

						ReinsExpectedCashInflow = sap[0].Treaty1CashInflowM1 + sap[0].Treaty1CashInflowM2 + sap[0].Treaty1CashInflowM3 + sap[0].Treaty1CashInflowM4 + sap[0].Treaty1CashInflowM5 + sap[0].Treaty1CashInflowM6 + sap[0].Treaty1CashInflowM7 +
							sap[0].Treaty2CashInflowM1 + sap[0].Treaty2CashInflowM2 + sap[0].Treaty2CashInflowM3 + sap[0].Treaty2CashInflowM4 + sap[0].Treaty2CashInflowM5 + sap[0].Treaty2CashInflowM6 + sap[0].Treaty2CashInflowM7 +
							sap[0].Treaty3CashInflowM1 + sap[0].Treaty3CashInflowM2 + sap[0].Treaty3CashInflowM3 + sap[0].Treaty3CashInflowM4 + sap[0].Treaty3CashInflowM5 + sap[0].Treaty3CashInflowM6 + sap[0].Treaty3CashInflowM7

					case 8:
						expectedClaims = sap[0].MorbidityOutgoM1 + sap[0].MorbidityOutgoM2 + sap[0].MorbidityOutgoM3 + sap[0].MorbidityOutgoM4 + sap[0].MorbidityOutgoM5 + sap[0].MorbidityOutgoM6 + sap[0].MorbidityOutgoM7 + sap[0].MorbidityOutgoM8 + sap[0].NonLifeClaimsOutgoM1 + sap[0].NonLifeClaimsOutgoM2 + sap[0].NonLifeClaimsOutgoM3 + sap[0].NonLifeClaimsOutgoM4 + sap[0].NonLifeClaimsOutgoM5 + sap[0].NonLifeClaimsOutgoM6 + sap[0].NonLifeClaimsOutgoM7 + sap[0].NonLifeClaimsOutgoM8

						ReinsExpectedCashOutflow = sap[0].Treaty1CashOutflowM1 + sap[0].Treaty1CashOutflowM2 + sap[0].Treaty1CashOutflowM3 + sap[0].Treaty1CashOutflowM4 + sap[0].Treaty1CashOutflowM5 + sap[0].Treaty1CashOutflowM6 + sap[0].Treaty1CashOutflowM7 + sap[0].Treaty1CashOutflowM8 +
							sap[0].Treaty2CashOutflowM1 + sap[0].Treaty2CashOutflowM2 + sap[0].Treaty2CashOutflowM3 + sap[0].Treaty2CashOutflowM4 + sap[0].Treaty2CashOutflowM5 + sap[0].Treaty2CashOutflowM6 + sap[0].Treaty2CashOutflowM7 + sap[0].Treaty2CashOutflowM8 +
							sap[0].Treaty3CashOutflowM1 + sap[0].Treaty3CashOutflowM2 + sap[0].Treaty3CashOutflowM3 + sap[0].Treaty3CashOutflowM4 + sap[0].Treaty3CashOutflowM5 + sap[0].Treaty3CashOutflowM6 + sap[0].Treaty3CashOutflowM7 + sap[0].Treaty3CashOutflowM8

						ReinsExpectedCashInflow = sap[0].Treaty1CashInflowM1 + sap[0].Treaty1CashInflowM2 + sap[0].Treaty1CashInflowM3 + sap[0].Treaty1CashInflowM4 + sap[0].Treaty1CashInflowM5 + sap[0].Treaty1CashInflowM6 + sap[0].Treaty1CashInflowM7 + sap[0].Treaty1CashInflowM8 +
							sap[0].Treaty2CashInflowM1 + sap[0].Treaty2CashInflowM2 + sap[0].Treaty2CashInflowM3 + sap[0].Treaty2CashInflowM4 + sap[0].Treaty2CashInflowM5 + sap[0].Treaty2CashInflowM6 + sap[0].Treaty2CashInflowM7 + sap[0].Treaty2CashInflowM8 +
							sap[0].Treaty3CashInflowM1 + sap[0].Treaty3CashInflowM2 + sap[0].Treaty3CashInflowM3 + sap[0].Treaty3CashInflowM4 + sap[0].Treaty3CashInflowM5 + sap[0].Treaty3CashInflowM6 + sap[0].Treaty3CashInflowM7 + sap[0].Treaty3CashInflowM8

					case 9:
						expectedClaims = sap[0].MorbidityOutgoM1 + sap[0].MorbidityOutgoM2 + sap[0].MorbidityOutgoM3 + sap[0].MorbidityOutgoM4 + sap[0].MorbidityOutgoM5 + sap[0].MorbidityOutgoM6 + sap[0].MorbidityOutgoM7 + sap[0].MorbidityOutgoM8 + sap[0].MorbidityOutgoM9 + sap[0].NonLifeClaimsOutgoM1 + sap[0].NonLifeClaimsOutgoM2 + sap[0].NonLifeClaimsOutgoM3 + sap[0].NonLifeClaimsOutgoM4 + sap[0].NonLifeClaimsOutgoM5 + sap[0].NonLifeClaimsOutgoM6 + sap[0].NonLifeClaimsOutgoM7 + sap[0].NonLifeClaimsOutgoM8 + sap[0].NonLifeClaimsOutgoM9

						ReinsExpectedCashOutflow = sap[0].Treaty1CashOutflowM1 + sap[0].Treaty1CashOutflowM2 + sap[0].Treaty1CashOutflowM3 + sap[0].Treaty1CashOutflowM4 + sap[0].Treaty1CashOutflowM5 + sap[0].Treaty1CashOutflowM6 + sap[0].Treaty1CashOutflowM7 + sap[0].Treaty1CashOutflowM8 + sap[0].Treaty1CashOutflowM9 +
							sap[0].Treaty2CashOutflowM1 + sap[0].Treaty2CashOutflowM2 + sap[0].Treaty2CashOutflowM3 + sap[0].Treaty2CashOutflowM4 + sap[0].Treaty2CashOutflowM5 + sap[0].Treaty2CashOutflowM6 + sap[0].Treaty2CashOutflowM7 + sap[0].Treaty2CashOutflowM8 + sap[0].Treaty2CashOutflowM9 +
							sap[0].Treaty3CashOutflowM1 + sap[0].Treaty3CashOutflowM2 + sap[0].Treaty3CashOutflowM3 + sap[0].Treaty3CashOutflowM4 + sap[0].Treaty3CashOutflowM5 + sap[0].Treaty3CashOutflowM6 + sap[0].Treaty3CashOutflowM7 + sap[0].Treaty3CashOutflowM8 + sap[0].Treaty3CashOutflowM9

						ReinsExpectedCashInflow = sap[0].Treaty1CashInflowM1 + sap[0].Treaty1CashInflowM2 + sap[0].Treaty1CashInflowM3 + sap[0].Treaty1CashInflowM4 + sap[0].Treaty1CashInflowM5 + sap[0].Treaty1CashInflowM6 + sap[0].Treaty1CashInflowM7 + sap[0].Treaty1CashInflowM8 + sap[0].Treaty1CashInflowM9 +
							sap[0].Treaty2CashInflowM1 + sap[0].Treaty2CashInflowM2 + sap[0].Treaty2CashInflowM3 + sap[0].Treaty2CashInflowM4 + sap[0].Treaty2CashInflowM5 + sap[0].Treaty2CashInflowM6 + sap[0].Treaty2CashInflowM7 + sap[0].Treaty2CashInflowM8 + sap[0].Treaty2CashInflowM9 +
							sap[0].Treaty3CashInflowM1 + sap[0].Treaty3CashInflowM2 + sap[0].Treaty3CashInflowM3 + sap[0].Treaty3CashInflowM4 + sap[0].Treaty3CashInflowM5 + sap[0].Treaty3CashInflowM6 + sap[0].Treaty3CashInflowM7 + sap[0].Treaty3CashInflowM8 + sap[0].Treaty3CashInflowM9

					case 10:
						expectedClaims = sap[0].MorbidityOutgoM1 + sap[0].MorbidityOutgoM2 + sap[0].MorbidityOutgoM3 + sap[0].MorbidityOutgoM4 + sap[0].MorbidityOutgoM5 + sap[0].MorbidityOutgoM6 + sap[0].MorbidityOutgoM7 + sap[0].MorbidityOutgoM8 + sap[0].MorbidityOutgoM9 + sap[0].MorbidityOutgoM10 + sap[0].NonLifeClaimsOutgoM1 + sap[0].NonLifeClaimsOutgoM2 + sap[0].NonLifeClaimsOutgoM3 + sap[0].NonLifeClaimsOutgoM4 + sap[0].NonLifeClaimsOutgoM5 + sap[0].NonLifeClaimsOutgoM6 + sap[0].NonLifeClaimsOutgoM7 + sap[0].NonLifeClaimsOutgoM8 + sap[0].NonLifeClaimsOutgoM9 + sap[0].NonLifeClaimsOutgoM10

						ReinsExpectedCashOutflow = sap[0].Treaty1CashOutflowM1 + sap[0].Treaty1CashOutflowM2 + sap[0].Treaty1CashOutflowM3 + sap[0].Treaty1CashOutflowM4 + sap[0].Treaty1CashOutflowM5 + sap[0].Treaty1CashOutflowM6 + sap[0].Treaty1CashOutflowM7 + sap[0].Treaty1CashOutflowM8 + sap[0].Treaty1CashOutflowM9 + sap[0].Treaty1CashOutflowM10 +
							sap[0].Treaty2CashOutflowM1 + sap[0].Treaty2CashOutflowM2 + sap[0].Treaty2CashOutflowM3 + sap[0].Treaty2CashOutflowM4 + sap[0].Treaty2CashOutflowM5 + sap[0].Treaty2CashOutflowM6 + sap[0].Treaty2CashOutflowM7 + sap[0].Treaty2CashOutflowM8 + sap[0].Treaty2CashOutflowM9 + sap[0].Treaty2CashOutflowM10 +
							sap[0].Treaty3CashOutflowM1 + sap[0].Treaty3CashOutflowM2 + sap[0].Treaty3CashOutflowM3 + sap[0].Treaty3CashOutflowM4 + sap[0].Treaty3CashOutflowM5 + sap[0].Treaty3CashOutflowM6 + sap[0].Treaty3CashOutflowM7 + sap[0].Treaty3CashOutflowM8 + sap[0].Treaty3CashOutflowM9 + sap[0].Treaty3CashOutflowM10

						ReinsExpectedCashInflow = sap[0].Treaty1CashInflowM1 + sap[0].Treaty1CashInflowM2 + sap[0].Treaty1CashInflowM3 + sap[0].Treaty1CashInflowM4 + sap[0].Treaty1CashInflowM5 + sap[0].Treaty1CashInflowM6 + sap[0].Treaty1CashInflowM7 + sap[0].Treaty1CashInflowM8 + sap[0].Treaty1CashInflowM9 + sap[0].Treaty1CashInflowM10 +
							sap[0].Treaty2CashInflowM1 + sap[0].Treaty2CashInflowM2 + sap[0].Treaty2CashInflowM3 + sap[0].Treaty2CashInflowM4 + sap[0].Treaty2CashInflowM5 + sap[0].Treaty2CashInflowM6 + sap[0].Treaty2CashInflowM7 + sap[0].Treaty2CashInflowM8 + sap[0].Treaty2CashInflowM9 + sap[0].Treaty2CashInflowM10 +
							sap[0].Treaty3CashInflowM1 + sap[0].Treaty3CashInflowM2 + sap[0].Treaty3CashInflowM3 + sap[0].Treaty3CashInflowM4 + sap[0].Treaty3CashInflowM5 + sap[0].Treaty3CashInflowM6 + sap[0].Treaty3CashInflowM7 + sap[0].Treaty3CashInflowM8 + sap[0].Treaty3CashInflowM9 + sap[0].Treaty3CashInflowM10

					case 11:
						expectedClaims = sap[0].MorbidityOutgoM1 + sap[0].MorbidityOutgoM2 + sap[0].MorbidityOutgoM3 + sap[0].MorbidityOutgoM4 + sap[0].MorbidityOutgoM5 + sap[0].MorbidityOutgoM6 + sap[0].MorbidityOutgoM7 + sap[0].MorbidityOutgoM8 + sap[0].MorbidityOutgoM9 + sap[0].MorbidityOutgoM10 + sap[0].MorbidityOutgoM11 + sap[0].NonLifeClaimsOutgoM1 + sap[0].NonLifeClaimsOutgoM2 + sap[0].NonLifeClaimsOutgoM3 + sap[0].NonLifeClaimsOutgoM4 + sap[0].NonLifeClaimsOutgoM5 + sap[0].NonLifeClaimsOutgoM6 + sap[0].NonLifeClaimsOutgoM7 + sap[0].NonLifeClaimsOutgoM8 + sap[0].NonLifeClaimsOutgoM9 + sap[0].NonLifeClaimsOutgoM10 + sap[0].NonLifeClaimsOutgoM11

						ReinsExpectedCashOutflow = sap[0].Treaty1CashOutflowM1 + sap[0].Treaty1CashOutflowM2 + sap[0].Treaty1CashOutflowM3 + sap[0].Treaty1CashOutflowM4 + sap[0].Treaty1CashOutflowM5 + sap[0].Treaty1CashOutflowM6 + sap[0].Treaty1CashOutflowM7 + sap[0].Treaty1CashOutflowM8 + sap[0].Treaty1CashOutflowM9 + sap[0].Treaty1CashOutflowM10 + sap[0].Treaty1CashOutflowM11 +
							sap[0].Treaty2CashOutflowM1 + sap[0].Treaty2CashOutflowM2 + sap[0].Treaty2CashOutflowM3 + sap[0].Treaty2CashOutflowM4 + sap[0].Treaty2CashOutflowM5 + sap[0].Treaty2CashOutflowM6 + sap[0].Treaty2CashOutflowM7 + sap[0].Treaty2CashOutflowM8 + sap[0].Treaty2CashOutflowM9 + sap[0].Treaty2CashOutflowM10 + sap[0].Treaty2CashOutflowM11 +
							sap[0].Treaty3CashOutflowM1 + sap[0].Treaty3CashOutflowM2 + sap[0].Treaty3CashOutflowM3 + sap[0].Treaty3CashOutflowM4 + sap[0].Treaty3CashOutflowM5 + sap[0].Treaty3CashOutflowM6 + sap[0].Treaty3CashOutflowM7 + sap[0].Treaty3CashOutflowM8 + sap[0].Treaty3CashOutflowM9 + sap[0].Treaty3CashOutflowM10 + sap[0].Treaty3CashOutflowM11

						ReinsExpectedCashInflow = sap[0].Treaty1CashInflowM1 + sap[0].Treaty1CashInflowM2 + sap[0].Treaty1CashInflowM3 + sap[0].Treaty1CashInflowM4 + sap[0].Treaty1CashInflowM5 + sap[0].Treaty1CashInflowM6 + sap[0].Treaty1CashInflowM7 + sap[0].Treaty1CashInflowM8 + sap[0].Treaty1CashInflowM9 + sap[0].Treaty1CashInflowM10 + sap[0].Treaty1CashInflowM11 +
							sap[0].Treaty2CashInflowM1 + sap[0].Treaty2CashInflowM2 + sap[0].Treaty2CashInflowM3 + sap[0].Treaty2CashInflowM4 + sap[0].Treaty2CashInflowM5 + sap[0].Treaty2CashInflowM6 + sap[0].Treaty2CashInflowM7 + sap[0].Treaty2CashInflowM8 + sap[0].Treaty2CashInflowM9 + sap[0].Treaty2CashInflowM10 + sap[0].Treaty2CashInflowM11 +
							sap[0].Treaty3CashInflowM1 + sap[0].Treaty3CashInflowM2 + sap[0].Treaty3CashInflowM3 + sap[0].Treaty3CashInflowM4 + sap[0].Treaty3CashInflowM5 + sap[0].Treaty3CashInflowM6 + sap[0].Treaty3CashInflowM7 + sap[0].Treaty3CashInflowM8 + sap[0].Treaty3CashInflowM9 + sap[0].Treaty3CashInflowM10 + sap[0].Treaty3CashInflowM11

					case 12:
						expectedClaims = sap[0].MorbidityOutgoM1 + sap[0].MorbidityOutgoM2 + sap[0].MorbidityOutgoM3 + sap[0].MorbidityOutgoM4 + sap[0].MorbidityOutgoM5 + sap[0].MorbidityOutgoM6 + sap[0].MorbidityOutgoM7 + sap[0].MorbidityOutgoM8 + sap[0].MorbidityOutgoM9 + sap[0].MorbidityOutgoM10 + sap[0].MorbidityOutgoM11 + sap[0].MorbidityOutgoM12 + sap[0].NonLifeClaimsOutgoM1 + sap[0].NonLifeClaimsOutgoM2 + sap[0].NonLifeClaimsOutgoM3 + sap[0].NonLifeClaimsOutgoM4 + sap[0].NonLifeClaimsOutgoM5 + sap[0].NonLifeClaimsOutgoM6 + sap[0].NonLifeClaimsOutgoM7 + sap[0].NonLifeClaimsOutgoM8 + sap[0].NonLifeClaimsOutgoM9 + sap[0].NonLifeClaimsOutgoM10 + sap[0].NonLifeClaimsOutgoM11 + sap[0].NonLifeClaimsOutgoM12

						ReinsExpectedCashOutflow = sap[0].Treaty1CashOutflowM1 + sap[0].Treaty1CashOutflowM2 + sap[0].Treaty1CashOutflowM3 + sap[0].Treaty1CashOutflowM4 + sap[0].Treaty1CashOutflowM5 + sap[0].Treaty1CashOutflowM6 + sap[0].Treaty1CashOutflowM7 + sap[0].Treaty1CashOutflowM8 + sap[0].Treaty1CashOutflowM9 + sap[0].Treaty1CashOutflowM10 + sap[0].Treaty1CashOutflowM11 + sap[0].Treaty1CashOutflowM12 +
							sap[0].Treaty2CashOutflowM1 + sap[0].Treaty2CashOutflowM2 + sap[0].Treaty2CashOutflowM3 + sap[0].Treaty2CashOutflowM4 + sap[0].Treaty2CashOutflowM5 + sap[0].Treaty2CashOutflowM6 + sap[0].Treaty2CashOutflowM7 + sap[0].Treaty2CashOutflowM8 + sap[0].Treaty2CashOutflowM9 + sap[0].Treaty2CashOutflowM10 + sap[0].Treaty2CashOutflowM11 + sap[0].Treaty2CashOutflowM12 +
							sap[0].Treaty3CashOutflowM1 + sap[0].Treaty3CashOutflowM2 + sap[0].Treaty3CashOutflowM3 + sap[0].Treaty3CashOutflowM4 + sap[0].Treaty3CashOutflowM5 + sap[0].Treaty3CashOutflowM6 + sap[0].Treaty3CashOutflowM7 + sap[0].Treaty3CashOutflowM8 + sap[0].Treaty3CashOutflowM9 + sap[0].Treaty3CashOutflowM10 + sap[0].Treaty3CashOutflowM11 + sap[0].Treaty3CashOutflowM12

						ReinsExpectedCashInflow = sap[0].Treaty1CashInflowM1 + sap[0].Treaty1CashInflowM2 + sap[0].Treaty1CashInflowM3 + sap[0].Treaty1CashInflowM4 + sap[0].Treaty1CashInflowM5 + sap[0].Treaty1CashInflowM6 + sap[0].Treaty1CashInflowM7 + sap[0].Treaty1CashInflowM8 + sap[0].Treaty1CashInflowM9 + sap[0].Treaty1CashInflowM10 + sap[0].Treaty1CashInflowM11 + sap[0].Treaty1CashInflowM12 +
							sap[0].Treaty2CashInflowM1 + sap[0].Treaty2CashInflowM2 + sap[0].Treaty2CashInflowM3 + sap[0].Treaty2CashInflowM4 + sap[0].Treaty2CashInflowM5 + sap[0].Treaty2CashInflowM6 + sap[0].Treaty2CashInflowM7 + sap[0].Treaty2CashInflowM8 + sap[0].Treaty2CashInflowM9 + sap[0].Treaty2CashInflowM10 + sap[0].Treaty2CashInflowM11 + sap[0].Treaty2CashInflowM12 +
							sap[0].Treaty3CashInflowM1 + sap[0].Treaty3CashInflowM2 + sap[0].Treaty3CashInflowM3 + sap[0].Treaty3CashInflowM4 + sap[0].Treaty3CashInflowM5 + sap[0].Treaty3CashInflowM6 + sap[0].Treaty3CashInflowM7 + sap[0].Treaty3CashInflowM8 + sap[0].Treaty3CashInflowM9 + sap[0].Treaty3CashInflowM10 + sap[0].Treaty3CashInflowM11 + sap[0].Treaty3CashInflowM12

					default:
						expectedClaims = 0
					}
					expectedCashOutflow += expectedClaims
				}

			}
			result.BestEstimateLiabilityChange = -expectedClaims
			result.BelInflowChange = 0
			result.BelOutflowChange = -expectedClaims
			result.RiskAdjustmentChange = 0
			result.CSMChange = 0
			result.LossComponentChange = 0 //-(expectedClaims * result.LcSar)
			result.LiabilityChange = -expectedClaims
			result.PNLChange = expectedClaims                                       //* (1 - result.LcSar)
			result.SarActualClaimNetLc = actualMorbidityClaim + actualNonLifeClaims //- expectedClaims*result.LcSar
			result.CSMBuildup = stepResults[i-1].CSMBuildup + result.CSMChange
			result.LossComponentBuildup = stepResults[i-1].LossComponentBuildup + result.LossComponentChange
			result.RiskAdjustmentBuildup = stepResults[i-1].RiskAdjustmentBuildup + result.RiskAdjustmentChange
			result.DACBuildup = stepResults[i-1].DACBuildup

			//reinsurance
			result.LossRecoveryComponentBuildup = stepResults[i-1].LossRecoveryComponentBuildup + result.LossRecoveryComponentChange
			result.ReinsCSMBuildup = stepResults[i-1].ReinsCSMBuildup + result.ReinsCSMChange

		}

		if item.Name == "Exp_Exp" {
			var expectedExpenses float64
			var actualExpenses = financeVariables.ActualAttributableExpenses
			var coverage int

			if financeVariables.IFStatus == IF {
				//coverage = 12
				var temp float64
				temp = utils.FloatPrecision((float64(currYear)+float64(currMonth)/1000.0)-(float64(prevYear)+float64(prevMonth)/1000.0), 3) //(float64(currYear) + float64(currMonth/1000.0)) - (float64(prevYear) + float64(prevMonth/1000.0))
				if temp == 1 {
					coverage = 12
				} else {
					if temp > 12 {
						coverage = int(math.Max(math.Min(12-(1000.0-math.Mod(temp, 1.0)*1000.0), 12), 0))
					}
					if temp <= 12 {
						coverage = int(math.Max(math.Min(temp*1000, 12), 0))
					}
				}
			} else {
				coverage = int(math.Min(float64(financeVariables.DurationInForceMonths), 12.0))
			}

			if len(sap) > 0 {
				if !csmRun.ManualSap {
					for i := 1; i <= coverage; i++ {
						//if len(sap) > 0 {
						expectedExpenses += sap[i].RenewalExpenses
						//}
					}
				} else {
					switch coverage {
					case 1:
						expectedExpenses = sap[0].RenewalExpenseOutgoM1
					case 2:
						expectedExpenses = sap[0].RenewalExpenseOutgoM1 + sap[0].RenewalExpenseOutgoM2
					case 3:
						expectedExpenses = sap[0].RenewalExpenseOutgoM1 + sap[0].RenewalExpenseOutgoM2 + sap[0].RenewalExpenseOutgoM3
					case 4:
						expectedExpenses = sap[0].RenewalExpenseOutgoM1 + sap[0].RenewalExpenseOutgoM2 + sap[0].RenewalExpenseOutgoM3 + sap[0].RenewalExpenseOutgoM4
					case 5:
						expectedExpenses = sap[0].RenewalExpenseOutgoM1 + sap[0].RenewalExpenseOutgoM2 + sap[0].RenewalExpenseOutgoM3 + sap[0].RenewalExpenseOutgoM4 + sap[0].RenewalExpenseOutgoM5
					case 6:
						expectedExpenses = sap[0].RenewalExpenseOutgoM1 + sap[0].RenewalExpenseOutgoM2 + sap[0].RenewalExpenseOutgoM3 + sap[0].RenewalExpenseOutgoM4 + sap[0].RenewalExpenseOutgoM5 + sap[0].RenewalExpenseOutgoM6
					case 7:
						expectedExpenses = sap[0].RenewalExpenseOutgoM1 + sap[0].RenewalExpenseOutgoM2 + sap[0].RenewalExpenseOutgoM3 + sap[0].RenewalExpenseOutgoM4 + sap[0].RenewalExpenseOutgoM5 + sap[0].RenewalExpenseOutgoM6 + sap[0].RenewalExpenseOutgoM7
					case 8:
						expectedExpenses = sap[0].RenewalExpenseOutgoM1 + sap[0].RenewalExpenseOutgoM2 + sap[0].RenewalExpenseOutgoM3 + sap[0].RenewalExpenseOutgoM4 + sap[0].RenewalExpenseOutgoM5 + sap[0].RenewalExpenseOutgoM6 + sap[0].RenewalExpenseOutgoM7 + sap[0].RenewalExpenseOutgoM8
					case 9:
						expectedExpenses = sap[0].RenewalExpenseOutgoM1 + sap[0].RenewalExpenseOutgoM2 + sap[0].RenewalExpenseOutgoM3 + sap[0].RenewalExpenseOutgoM4 + sap[0].RenewalExpenseOutgoM5 + sap[0].RenewalExpenseOutgoM6 + sap[0].RenewalExpenseOutgoM7 + sap[0].RenewalExpenseOutgoM8 + sap[0].RenewalExpenseOutgoM9
					case 10:
						expectedExpenses = sap[0].RenewalExpenseOutgoM1 + sap[0].RenewalExpenseOutgoM2 + sap[0].RenewalExpenseOutgoM3 + sap[0].RenewalExpenseOutgoM4 + sap[0].RenewalExpenseOutgoM5 + sap[0].RenewalExpenseOutgoM6 + sap[0].RenewalExpenseOutgoM7 + sap[0].RenewalExpenseOutgoM8 + sap[0].RenewalExpenseOutgoM9 + sap[0].RenewalExpenseOutgoM10
					case 11:
						expectedExpenses = sap[0].RenewalExpenseOutgoM1 + sap[0].RenewalExpenseOutgoM2 + sap[0].RenewalExpenseOutgoM3 + sap[0].RenewalExpenseOutgoM4 + sap[0].RenewalExpenseOutgoM5 + sap[0].RenewalExpenseOutgoM6 + sap[0].RenewalExpenseOutgoM7 + sap[0].RenewalExpenseOutgoM8 + sap[0].RenewalExpenseOutgoM9 + sap[0].RenewalExpenseOutgoM10 + sap[0].RenewalExpenseOutgoM11
					case 12:
						expectedExpenses = sap[0].RenewalExpenseOutgoM1 + sap[0].RenewalExpenseOutgoM2 + sap[0].RenewalExpenseOutgoM3 + sap[0].RenewalExpenseOutgoM4 + sap[0].RenewalExpenseOutgoM5 + sap[0].RenewalExpenseOutgoM6 + sap[0].RenewalExpenseOutgoM7 + sap[0].RenewalExpenseOutgoM8 + sap[0].RenewalExpenseOutgoM9 + sap[0].RenewalExpenseOutgoM10 + sap[0].RenewalExpenseOutgoM11 + sap[0].RenewalExpenseOutgoM12
					default:
						expectedExpenses = 0
					}
					expectedCashOutflow += expectedExpenses
				}

			}

			result.BestEstimateLiabilityChange = -expectedExpenses
			result.BelInflowChange = 0
			result.BelOutflowChange = -expectedExpenses
			result.RiskAdjustmentChange = 0
			result.CSMChange = 0
			result.LossComponentChange = 0 //-(expectedExpenses * result.LcSar)
			result.LiabilityChange = -expectedExpenses
			result.PNLChange = expectedExpenses //* (1 - result.LcSar)
			result.SarActualExpenseNetLc = actualExpenses
			result.ActualNonAttributableExpenses = financeVariables.ActualNonAttributableExpenses //- expectedExpenses*result.LcSar
			result.CSMBuildup = stepResults[i-1].CSMBuildup + result.CSMChange
			result.LossComponentBuildup = stepResults[i-1].LossComponentBuildup + result.LossComponentChange
			result.DACBuildup = stepResults[i-1].DACBuildup
			result.RiskAdjustmentBuildup = stepResults[i-1].RiskAdjustmentBuildup + result.RiskAdjustmentChange

			//reinsurance
			result.LossRecoveryComponentBuildup = stepResults[i-1].LossRecoveryComponentBuildup + result.LossRecoveryComponentChange
			result.ReinsCSMBuildup = stepResults[i-1].ReinsCSMBuildup + result.ReinsCSMChange

		}

		if item.Name == "Exp_RA_Release" {

			var interestaccretionfac float64
			var durationstart int
			//var valuationyear int
			//valuationyear = getCurrentYear()
			var coverage int

			if financeVariables.IFStatus == IF { //valuationyear > financeVariables.LockedInYear {
				var temp float64
				temp = utils.FloatPrecision((float64(currYear)+float64(currMonth)/1000.0)-(float64(prevYear)+float64(prevMonth)/1000.0), 3) //(float64(currYear) + float64(currMonth/1000.0)) - (float64(prevYear) + float64(prevMonth/1000.0))
				if temp == 1 {
					coverage = 12
				} else {
					if temp > 12 {
						coverage = int(math.Max(math.Min(12-(1000.0-math.Mod(temp, 1.0)*1000.0), 12), 0))
					}
					if temp <= 12 {
						coverage = int(math.Max(math.Min(temp*1000, 12), 0))
					}
				}
				durationstart = int(math.Max(float64(financeVariables.DurationInForceMonths-coverage), 0)) //(valuationyear-financeVariables.LockedInYear)*12 - financeVariables.LockedInMonth + 1
			} else {
				durationstart = 0
			}

			interestaccretionfac = InterestAccretionFactor(durationstart, financeVariables.LockedInYear, financeVariables.LockedInMonth, financeVariables.IFStatus, group.ProductCode, financeVariables.YieldCurveCode, financeVariables.YieldCurveBasis)

			result.BestEstimateLiabilityChange = 0
			result.BelInflowChange = 0
			result.BelOutflowChange = 0
			//result.RiskAdjustment * interestaccretionfac
			result.RiskAdjustmentChange = result.RiskAdjustmentAt12 - result.RiskAdjustment - stepResults[1].RiskAdjustment*interestaccretionfac //-math.Max(result.RiskAdjustment-result.RiskAdjustmentAt12-riskAdjustmentAccretion, 0)
			result.CSMChange = 0
			result.LossComponentChange = -result.LossComponentBuildup * result.LcSar //result.RiskAdjustmentChange * result.LcSar
			result.LiabilityChange = result.RiskAdjustmentChange
			result.PNLChange = -result.RiskAdjustmentChange
			result.ExpectedRaNetOfLc = -result.RiskAdjustmentChange * (1 - result.LcSar)
			result.CSMBuildup = stepResults[i-1].CSMBuildup + result.CSMChange - result.CSMRelease
			result.LossComponentBuildup = stepResults[i-1].LossComponentBuildup + result.LossComponentChange
			result.DACBuildup = stepResults[i-1].DACBuildup
			result.RiskAdjustmentBuildup = stepResults[i-1].RiskAdjustmentBuildup + result.RiskAdjustmentChange

			//reinsurance
			result.ReinsRAChange = result.ReinsuranceRiskAdjustmentAt12 - result.ReinsuranceRiskAdjustment - stepResults[1].ReinsuranceRiskAdjustment*interestaccretionfac
			result.ReinsPNLChange = result.ReinsRAChange
			result.LossRecoveryComponentBuildup = stepResults[i-1].LossRecoveryComponentBuildup + result.LossRecoveryComponentChange
			result.ReinsCSMBuildup = stepResults[i-1].ReinsCSMBuildup + result.ReinsCSMChange

		}

		if i > 8 {
			if item.Name == "Data_Change" {

				result.BestEstimateLiabilityChange = result.BEL - stepResults[i-1].BELAt12
				result.BelInflowChange = result.BelInflow - stepResults[i-1].BelInflowAt12
				result.BelOutflowChange = result.BelOutflow - stepResults[i-1].BelOutflowAt12
				result.RiskAdjustmentChange = result.RiskAdjustment - stepResults[i-1].RiskAdjustmentAt12
				result.LiabilityChange = result.BestEstimateLiabilityChange + result.RiskAdjustmentChange
				result.DACBuildup = stepResults[i-1].DACBuildup

				result.ReinsBelInflowChange = result.ReinsuranceBelInflow - stepResults[i-1].ReinsuranceBelInflowAt12
				result.ReinsBelOutflowChange = result.ReinsuranceBelOutflow - stepResults[i-1].ReinsuranceBelOutflowAt12
				result.ReinsBelChange = result.ReinsuranceBel - stepResults[i-1].ReinsuranceBelAt12
				result.ReinsRAChange = result.ReinsuranceRiskAdjustment - stepResults[i-1].ReinsuranceRiskAdjustmentAt12

			} else {
				if item.Name == "FA" {
					if csmRun.MeasurementType == "VFA" {
						result.BestEstimateLiabilityChange = result.BEL - stepResults[i-1].BEL
						result.BelInflowChange = result.BelInflow - stepResults[i-1].BelInflow
						result.BelOutflowChange = result.BelOutflow - stepResults[i-1].BelOutflow
						result.RiskAdjustmentChange = result.RiskAdjustment - stepResults[i-1].RiskAdjustment
						result.LiabilityChange = result.BestEstimateLiabilityChange + result.RiskAdjustmentChange
						result.DACBuildup = stepResults[i-1].DACBuildup
						result.RiskAdjustmentBuildup = stepResults[i-1].RiskAdjustmentBuildup + result.RiskAdjustmentChange

					} else {
						result.BestEstimateLiabilityChange = result.BEL - stepResults[i-1].BEL
						result.BelInflowChange = result.BelInflow - stepResults[i-1].BelInflow
						result.BelOutflowChange = result.BelOutflow - stepResults[i-1].BelOutflow
						result.RiskAdjustmentChange = result.RiskAdjustment - stepResults[i-1].RiskAdjustment
						result.LiabilityChange = result.BestEstimateLiabilityChange + result.RiskAdjustmentChange
						result.DACBuildup = stepResults[i-1].DACBuildup
						result.RiskAdjustmentBuildup = stepResults[i-1].RiskAdjustmentBuildup + result.RiskAdjustmentChange
					}
				} else {
					result.BestEstimateLiabilityChange = result.BEL - stepResults[i-1].BEL
					result.BelInflowChange = result.BelInflow - stepResults[i-1].BelInflow
					result.BelOutflowChange = result.BelOutflow - stepResults[i-1].BelOutflow
					result.RiskAdjustmentChange = result.RiskAdjustment - stepResults[i-1].RiskAdjustment
					result.LiabilityChange = result.BestEstimateLiabilityChange + result.RiskAdjustmentChange
					result.DACBuildup = stepResults[i-1].DACBuildup
					result.RiskAdjustmentBuildup = stepResults[i-1].RiskAdjustmentBuildup + result.RiskAdjustmentChange
				}
			}

			if csmRun.MeasurementType == "GMM" {
				if item.InsuranceService == "Current" { // changes affecting PNL
					result.CSMChange = 0
					result.LossComponentChange = 0
					result.PNLChange = -result.LiabilityChange
				} else if item.InsuranceService == "Future" { // changes affecting CSM
					var futureServiceAccumulatedChanges float64
					futureServiceAccumulatedChanges += -result.LiabilityChange
					if stepResults[i-1].CSMBuildup > 0 { // CSM carrying amount is positive
						if result.LiabilityChange > 0 { // reduces the CSM amount value
							result.CSMChange = -math.Min(result.LiabilityChange, stepResults[i-1].CSMBuildup)            //CSM change capped at the carrying value of the CSM
							result.LossComponentChange = math.Max(0, result.LiabilityChange-stepResults[i-1].CSMBuildup) // LossCompoenent is buildup on the proviso that the carrying CSM value is depleted
							result.PNLChange = -result.LossComponentChange
						} else {
							result.CSMChange = -result.LiabilityChange
							result.LossComponentChange = 0
							result.PNLChange = 0
						}

						//reinsurance

						result.ReinsCSMChange = -(result.ReinsBelChange + result.ReinsRAChange)
						if result.BelOutflow+result.RiskAdjustment > 0 {
							result.LossRecoveryComponentChange = result.LossComponentChange * (result.ReinsuranceBelOutflow + result.ReinsuranceRiskAdjustment) / (result.BelOutflow + result.RiskAdjustment)
						}

					} else { // CSM carrying amount is zero
						if result.LiabilityChange > 0 {
							result.LossComponentChange = result.LiabilityChange
							result.CSMChange = 0
							result.PNLChange = -result.LossComponentChange
						} else {
							result.LossComponentChange = -math.Min(-result.LiabilityChange, stepResults[i-1].LossComponentBuildup)
							result.CSMChange = math.Max(-result.LiabilityChange-stepResults[i-1].LossComponentBuildup, 0)
							result.PNLChange = -result.LossComponentChange
						}

						result.ReinsCSMChange = -(result.ReinsBelChange + result.ReinsRAChange)
						if result.BelOutflow+result.RiskAdjustment > 0 {
							result.LossRecoveryComponentChange = result.LossComponentChange * (result.ReinsuranceBelOutflow + result.ReinsuranceRiskAdjustment) / (result.BelOutflow + result.RiskAdjustment)
						}

					}
				} else {
					result.LossComponentChange = 0
					result.CSMChange = 0
				}
			} else { // if Measurement Model == VFA

				if item.Name == "FA" {
					result.PNLChange = result.LiabilityChange * financeVariables.RiskMitigationProportion
					if stepResults[i-1].CSMBuildup > 0 { // CSM carrying amount is positive
						if result.LiabilityChange > 0 { // reduces the CSM amount value
							result.CSMChange = -math.Min(result.LiabilityChange*(1-financeVariables.RiskMitigationProportion), stepResults[i-1].CSMBuildup)            //CSM change capped at the carrying value of the CSM
							result.LossComponentChange = math.Max(0, result.LiabilityChange*(1-financeVariables.RiskMitigationProportion)-stepResults[i-1].CSMBuildup) // LossCompoenent is buildup on the proviso that the carrying CSM value is depleted
							result.PNLChange = -result.LossComponentChange
						} else {
							result.CSMChange = -(result.LiabilityChange * (1 - financeVariables.RiskMitigationProportion))
							result.LossComponentChange = 0
							result.PNLChange = 0
						}
					} else { // CSM carrying amount is zero
						if result.LiabilityChange > 0 {
							result.LossComponentChange = result.LiabilityChange * (1 - financeVariables.RiskMitigationProportion)
							result.CSMChange = 0
							result.PNLChange = -result.LossComponentChange
						} else {
							result.LossComponentChange = -math.Min(-(result.LiabilityChange * (1 - financeVariables.RiskMitigationProportion)), stepResults[i-1].LossComponentBuildup)
							result.CSMChange = math.Max(-(result.LiabilityChange*(1-financeVariables.RiskMitigationProportion))-stepResults[i-1].LossComponentBuildup, 0)
							result.PNLChange = -result.LossComponentChange
						}
					}
				} else {
					if item.InsuranceService == "Current" { // changes affecting PNL
						result.CSMChange = 0
						result.LossComponentChange = 0
						result.PNLChange = -result.LiabilityChange
					} else if item.InsuranceService == "Future" { // changes affecting CSM
						var futureServiceAccumulatedChanges float64
						futureServiceAccumulatedChanges += -result.LiabilityChange
						if stepResults[i-1].CSMBuildup > 0 { // CSM carrying amount is positive
							if result.LiabilityChange > 0 { // reduces the CSM amount value
								result.CSMChange = -math.Min(result.LiabilityChange, stepResults[i-1].CSMBuildup)            //CSM change capped at the carrying value of the CSM
								result.LossComponentChange = math.Max(0, result.LiabilityChange-stepResults[i-1].CSMBuildup) // LossCompoenent is buildup on the proviso that the carrying CSM value is depleted
								result.PNLChange = -result.LossComponentChange
							} else {
								result.CSMChange = -result.LiabilityChange
								result.LossComponentChange = 0
								result.PNLChange = 0
							}
						} else { // CSM carrying amount is zero
							if result.LiabilityChange > 0 {
								result.LossComponentChange = result.LiabilityChange
								result.CSMChange = 0
								result.PNLChange = -result.LossComponentChange
							} else {
								result.LossComponentChange = -math.Min(-result.LiabilityChange, stepResults[i-1].LossComponentBuildup)
								result.CSMChange = math.Max(-result.LiabilityChange-stepResults[i-1].LossComponentBuildup, 0)
								result.PNLChange = -result.LossComponentChange
							}
						}
					} else {
						result.LossComponentChange = 0
						result.CSMChange = 0
					}

				}
			}

			//(coverageunits[start]-coverageunits[end])/coverageunits[start]
			//before any basis change
			if i >= 15 && stepResults[3].SumCoverageUnits > 0 {
				var currentService, futureService float64
				if aosconfigset.CoverageUnitOption == UndiscountedCoverageUnits {
					currentService = math.Max(stepResults[3].SumCoverageUnits-stepResults[3].SumCoverageUnitsAt12, 0)
					futureService = result.SumCoverageUnits
					result.CSMReleaseRatio = utils.FloatPrecision(currentService/(currentService+futureService), CsmDefaultPrecision) //interest accretion step
				}
				if aosconfigset.CoverageUnitOption == DiscountedCoverageUnits {

					currentService = math.Max(stepResults[3].DiscountedCoverageUnits-stepResults[3].DiscountedCoverageUnitsAt12, 0) // interestaccretion step
					futureService = result.DiscountedCoverageUnits
					result.CSMReleaseRatio = utils.FloatPrecision(currentService/(currentService+futureService), CsmDefaultPrecision)
				}

				if i == 16 {
					result.CSMRelease = utils.FloatPrecision(stepResults[len(stepResults)-1].CSMBuildup*stepResults[15].CSMReleaseRatio, 2)
					result.CSMChange = -result.CSMRelease
					result.DACChange = 0 //-utils.FloatPrecision(stepResults[len(stepResults)-1].DACBuildup*stepResults[15].CSMReleaseRatio, 2)

				} else {
					result.CSMRelease = 0
				}
			} else {
				result.CSMRelease = 0
				result.CSMReleaseRatio = 0
			}

			if i == 16 {

				result.BestEstimateLiabilityChange = result.BEL - stepResults[i-1].BEL
				result.BelInflowChange = result.BelInflow - stepResults[i-1].BelInflow
				result.BelOutflowChange = result.BelOutflow - stepResults[i-1].BelOutflow
				result.RiskAdjustmentChange = result.RiskAdjustment - stepResults[i-1].RiskAdjustment
				result.LiabilityChange = result.BestEstimateLiabilityChange + result.RiskAdjustmentChange
				result.PNLChange = -result.BestEstimateLiabilityChange - result.RiskAdjustmentChange + result.CSMRelease
				result.ReinsBelOutflowChange = result.ReinsuranceBelOutflow - stepResults[i-1].ReinsuranceBelOutflow
				result.ReinsBelInflowChange = result.ReinsuranceBelInflow - stepResults[i-1].ReinsuranceBelInflow
				result.ReinsRAChange = result.ReinsuranceRiskAdjustment - stepResults[i-1].ReinsuranceRiskAdjustment
				result.ReinsBelChange = result.ReinsuranceBel - stepResults[i-1].ReinsuranceBel
				result.ReinsPNLChange = result.ReinsBelChange + result.ReinsRAChange
				result.ReinsCSMBuildup = stepResults[i-1].ReinsCSMChange
				result.LossRecoveryComponentBuildup = stepResults[i-1].LossComponentBuildup
				result.CSMBuildup = utils.FloatPrecision(stepResults[i-1].CSMBuildup-result.CSMRelease, AccountingPrecision)
			} else {
				result.CSMBuildup = utils.FloatPrecision(stepResults[i-1].CSMBuildup+result.CSMChange-result.CSMRelease, AccountingPrecision)
			}

			result.LossComponentBuildup = utils.FloatPrecision(stepResults[i-1].LossComponentBuildup+result.LossComponentChange, AccountingPrecision)
			result.DACBuildup = utils.FloatPrecision(stepResults[i-1].DACBuildup+result.DACChange, AccountingPrecision)
			result.RiskAdjustmentBuildup = utils.FloatPrecision(stepResults[i-1].RiskAdjustmentBuildup+result.RiskAdjustmentChange, AccountingPrecision)

			result.ReinsCSMBuildup = stepResults[i-1].ReinsCSMBuildup + result.ReinsCSMChange

			if i == 16 {
				var bsr models.BalanceSheetRecord
				bsr.CsmRunID = csmRun.ID
				bsr.ProductCode = result.ProductCode
				bsr.IFRS17Group = result.IFRS17Group
				bsr.MeasurementType = csmRun.MeasurementType
				bsr.Date = csmRun.RunDate
				err = DB.Where("product_code = ? and measurement_type = ? and ifrs17_group = ? and date = ?", result.ProductCode, csmRun.MeasurementType, result.IFRS17Group, csmRun.RunDate).Delete(&models.BalanceSheetRecord{}).Error
				if err != nil {
					fmt.Println("balance sheet record delete error:", err)
				}
				if csmRun.TransitionType == PostTransition {
					bsr.PostTransitionCsm = result.CSMBuildup
					bsr.PostTransitionLc = result.LossComponentBuildup
					bsr.PostTransitionDAC = result.DACBuildup
				}

				if csmRun.TransitionType == FullyRetrospective {
					bsr.FullyRetrospectiveCsm = result.CSMBuildup
					bsr.FullyRetrospectiveLc = result.LossComponentBuildup
					bsr.FullyRetrospectiveDAC = result.DACBuildup
				}

				if csmRun.TransitionType == ModifiedRetrospective {
					bsr.ModifiedRetrospectiveCsm = result.CSMBuildup
					bsr.ModifiedRetrospectiveLc = result.LossComponentBuildup
					bsr.ModifiedRetrospectiveDAC = result.DACBuildup
				}

				if csmRun.TransitionType == FairValue {
					bsr.FairValueCsm = result.CSMBuildup
					bsr.FairValueLc = result.LossComponentBuildup
					bsr.FairValueDAC = result.DACBuildup
				}
				bsr.BELOutflow = result.BelOutflow
				bsr.BELInflow = result.BelInflow
				bsr.RiskAdjustment = result.RiskAdjustment
				bsr.Treaty1BELOutflow = 0
				bsr.Treaty1BELInflow = 0
				bsr.Treaty1RiskAdjustment = 0
				bsr.Treaty2BELOutflow = 0
				bsr.Treaty2BELInflow = 0
				bsr.Treaty2RiskAdjustment = 0
				bsr.Treaty3BELOutflow = 0
				bsr.Treaty3BELInflow = 0
				bsr.Treaty3RiskAdjustment = 0

				bsr.IBNR = 0 //MOTLATSI WILL REMEMBER TO DO THIS!!!
				err := DB.Create(&bsr).Error
				if err != nil {
					fmt.Println("balance sheet record error: ", err)
				}
				// CSM projections
				var csmProjection = make([]models.CsmProjection, 6)
				for cu := 0; cu <= MaxCsmProjectionPeriodYears; cu++ {
					cuI := cu * 12
					csmProjection[cu].RunDate = result.RunDate
					csmProjection[cu].ProjectionMonth = cuI
					if cu == 0 {
						csmProjection[cu].ProductCode = result.ProductCode
						csmProjection[cu].IFRS17Group = result.IFRS17Group
						if len(sap) > 0 {
							if !csmRun.ManualSap {
								for i, _ := range sap {
									if cuI == sap[i].ProjectionMonth {
										if aosconfigset.CoverageUnitOption == UndiscountedCoverageUnits {
											csmProjection[cu].CoverageUnits = sap[i].SumCoverageUnits
										}
										if aosconfigset.CoverageUnitOption == DiscountedCoverageUnits {
											csmProjection[cu].CoverageUnits = sap[i].DiscountedCoverageUnits
										}
									}
								}
							} else {
								if aosconfigset.CoverageUnitOption == UndiscountedCoverageUnits {
									csmProjection[cu].CoverageUnits = sap[0].SumCoverageUnitsM0
								}
								if aosconfigset.CoverageUnitOption == DiscountedCoverageUnits {
									csmProjection[cu].CoverageUnits = sap[0].DiscountedCoverageUnitsM0
								}
							}
						}
						csmProjection[cu].Csm = result.CSMBuildup
						csmProjection[cu].CsmRelease = 0
						csmProjection[cu].CoverageUnitsFullRetrospectiveApproach = 0
						csmProjection[cu].CsmFullRetrospectiveApproach = 0
						csmProjection[cu].CsmReleaseFullRetrospectiveApproach = 0
						csmProjection[cu].CoverageUnitsModifiedRetrospectiveApproach = 0
						csmProjection[cu].CsmReleaseModifiedRetrospectiveApproach = 0
						csmProjection[cu].CsmReleaseModifiedRetrospectiveApproach = 0
						csmProjection[cu].CoverageUnitsFairValueApproach = 0
						csmProjection[cu].CsmFairValueApproach = 0
						csmProjection[cu].CsmReleaseFairValueApproach = 0

					} else {
						csmProjection[cu].ProductCode = result.ProductCode
						csmProjection[cu].IFRS17Group = result.IFRS17Group
						if len(sap) > 0 {
							if !csmRun.ManualSap {
								for i, _ := range sap {
									if cuI == sap[i].ProjectionMonth {
										if aosconfigset.CoverageUnitOption == UndiscountedCoverageUnits {
											csmProjection[cu].CoverageUnits = sap[i].SumCoverageUnits
										}
										if aosconfigset.CoverageUnitOption == DiscountedCoverageUnits {
											csmProjection[cu].CoverageUnits = sap[i].DiscountedCoverageUnits
										}
									}
								}
							} else {
								if aosconfigset.CoverageUnitOption == UndiscountedCoverageUnits {
									switch cuI {
									case 0:
										csmProjection[cu].CoverageUnits = sap[0].SumCoverageUnitsM0
									case 12:
										csmProjection[cu].CoverageUnits = sap[0].SumCoverageUnitsM12
									case 24:
										csmProjection[cu].CoverageUnits = sap[0].SumCoverageUnitsM24
									case 36:
										csmProjection[cu].CoverageUnits = sap[0].SumCoverageUnitsM36
									case 48:
										csmProjection[cu].CoverageUnits = sap[0].SumCoverageUnitsM48
									case 60:
										csmProjection[cu].CoverageUnits = sap[0].SumCoverageUnitsM60
									case 72:
										csmProjection[cu].CoverageUnits = sap[0].SumCoverageUnitsM72
									case 84:
										csmProjection[cu].CoverageUnits = sap[0].SumCoverageUnitsM84
									case 96:
										csmProjection[cu].CoverageUnits = sap[0].SumCoverageUnitsM96
									case 108:
										csmProjection[cu].CoverageUnits = sap[0].SumCoverageUnitsM108
									case 120:
										csmProjection[cu].CoverageUnits = sap[0].SumCoverageUnitsM120
									default:
										csmProjection[cu].CoverageUnits = 0
									}
								}
								if aosconfigset.CoverageUnitOption == DiscountedCoverageUnits {
									switch cuI {
									case 0:
										csmProjection[cu].CoverageUnits = sap[0].DiscountedCoverageUnitsM0
									case 12:
										csmProjection[cu].CoverageUnits = sap[0].DiscountedCoverageUnitsM12
									case 24:
										csmProjection[cu].CoverageUnits = sap[0].DiscountedCoverageUnitsM24
									case 36:
										csmProjection[cu].CoverageUnits = sap[0].DiscountedCoverageUnitsM36
									case 48:
										csmProjection[cu].CoverageUnits = sap[0].DiscountedCoverageUnitsM48
									case 60:
										csmProjection[cu].CoverageUnits = sap[0].DiscountedCoverageUnitsM60
									case 72:
										csmProjection[cu].CoverageUnits = sap[0].DiscountedCoverageUnitsM72
									case 84:
										csmProjection[cu].CoverageUnits = sap[0].DiscountedCoverageUnitsM84
									case 96:
										csmProjection[cu].CoverageUnits = sap[0].DiscountedCoverageUnitsM96
									case 108:
										csmProjection[cu].CoverageUnits = sap[0].DiscountedCoverageUnitsM108
									case 120:
										csmProjection[cu].CoverageUnits = sap[0].DiscountedCoverageUnitsM120
									default:
										csmProjection[cu].CoverageUnits = 0
									}
								}
							}
						}
						if csmProjection[cu-1].CoverageUnits > 0 {
							csmProjection[cu].CsmRelease = csmProjection[cu-1].Csm * (math.Max(csmProjection[cu-1].CoverageUnits-csmProjection[cu].CoverageUnits, 0) / csmProjection[cu-1].CoverageUnits)
						} else {
							csmProjection[cu].CsmRelease = 0
						}
						csmProjection[cu].Csm = math.Max(csmProjection[cu-1].Csm-csmProjection[cu].CsmRelease, 0)
						csmProjection[cu].CoverageUnitsFullRetrospectiveApproach = 0
						csmProjection[cu].CsmFullRetrospectiveApproach = 0
						csmProjection[cu].CsmReleaseFullRetrospectiveApproach = 0
						csmProjection[cu].CoverageUnitsModifiedRetrospectiveApproach = 0
						csmProjection[cu].CsmReleaseModifiedRetrospectiveApproach = 0
						csmProjection[cu].CsmReleaseModifiedRetrospectiveApproach = 0
						csmProjection[cu].CoverageUnitsFairValueApproach = 0
						csmProjection[cu].CsmFairValueApproach = 0
						csmProjection[cu].CsmReleaseFairValueApproach = 0
					}
					csmProjection[cu].CsmTotal = csmProjection[cu].Csm +
						csmProjection[cu].CsmFullRetrospectiveApproach +
						csmProjection[cu].CsmModifiedRetrospectiveApproach +
						csmProjection[cu].CsmFairValueApproach

					csmProjection[cu].CsmReleaseTotal = csmProjection[cu].CsmRelease +
						csmProjection[cu].CsmReleaseFullRetrospectiveApproach +
						csmProjection[cu].CsmReleaseModifiedRetrospectiveApproach +
						csmProjection[cu].CsmFairValueApproach
					csmProjection[cu].CsmRunID = csmRun.ID
					csmProjection[cu].RunDate = csmRun.RunDate
					csmProjections = append(csmProjections, csmProjection[cu])
				}

				for _, csmProjection := range csmProjections {
					DB.Where("run_date = ? and product_code = ? and ifrs17_group = ? and projection_month = ?", csmProjection.RunDate, csmProjection.ProductCode, csmProjection.IFRS17Group, csmProjection.ProjectionMonth).Delete(&models.CsmProjection{})
					err = DB.Save(&csmProjection).Error
					if err != nil {
						fmt.Println(err)
					}
				}

				//err = DB.Create(&csmProjections).Error
				//if err != nil {
				//	log.Error(err)
				//}
			}

		}
		//}
		result.CsmRunID = csmRun.ID

		// Onerous contract detection
		// A group is onerous if opening CSM < 0 (loss component > 0 at initial recognition)
		// or if the CSM balance would go negative after release
		if result.CSMBuildup < 0 {
			result.IsOnerous = true
			result.OnerousReason = "opening CSM negative (loss-making group)"
		} else if result.CSMBuildup > 0 && result.CSMRelease > result.CSMBuildup {
			result.IsOnerous = true
			result.OnerousReason = "CSM release exceeds opening CSM balance"
		} else if result.LossComponentBuildup > 0 {
			result.IsOnerous = true
			result.OnerousReason = "loss component present"
		}

		stepResults = append(stepResults, result)
	}

	for _, step := range stepResults {
		DB.Where("name = ? and run_date = ? and product_code = ? and ifrs17_group = ?", step.Name, step.RunDate, step.ProductCode, step.IFRS17Group).Delete(&models.AOSStepResult{})
		err = DB.Create(&step).Error
		if err != nil {
			fmt.Println(err)
		}
	}

	jt.CsmRelease = stepResults[16].CSMRelease
	jt.PremiumVariance = stepResults[8].ExperiencePremiumVariance
	jt.DacRelease = -stepResults[16].DACChange
	jt.RAChange = -stepResults[3].RiskAdjustmentChange
	jt.ExpectedBenefits = stepResults[4].PNLChange + stepResults[6].PNLChange + stepResults[7].PNLChange
	jt.Expenses = stepResults[5].PNLChange
	jt.ClaimsIncurred = stepResults[4].SarActualClaimNetLc + stepResults[6].SarActualClaimNetLc + stepResults[7].SarActualClaimNetLc
	jt.ExpensesIncurred = stepResults[5].SarActualExpenseNetLc
	jt.AmortizationAcquisitionCF = jt.DacRelease
	jt.LossComponentFutureServiceChange = stepResults[9].LossComponentChange +
		stepResults[10].LossComponentChange +
		stepResults[11].LossComponentChange +
		stepResults[12].LossComponentChange +
		stepResults[13].LossComponentChange

	jt.LossComponentOnInitialRecog = stepResults[2].LossComponentChange
	jt.InsuranceFinanceExpense = stepResults[1].PNLChange +
		stepResults[8].PNLChange +
		stepResults[14].PNLChange +
		stepResults[16].PNLChange - stepResults[16].CSMRelease
	jt.LossComponentUnwind = stepResults[8].LossComponentUnwind
	jt.NonAttributableExpensesIncurred = stepResults[5].ActualNonAttributableExpenses

	err = DB.Where("product_code = ? and ifrs17_group = ? and run_date=? and measurement_type = ?", jt.ProductCode, jt.IFRS17Group, jt.RunDate, jt.MeasurementType).Delete(&jt).Error
	if err != nil {
		log.Info(err.Error())
	}

	err = DB.Where("product_code = ? and ifrs17_group = ? and run_date=? and measurement_type = ?", jt.ProductCode, jt.IFRS17Group, jt.RunDate, jt.MeasurementType).Save(&jt).Error
	if err != nil {
		log.Info(err.Error())
	}

	//Build data for liability_movement_lines

	return stepResults, nil
}

func getFinancialVariables(code string, group string, financeYear int, financeVersion string) models.FinanceVariables {
	var res models.FinanceVariables

	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Where("product_code = ? and ifrs17_group=? and year=? and version=?", code, group, financeYear, financeVersion).First(&res).Error }); err != nil {
		fmt.Println(err)
	}
	return res
}

func GetStepResultsForPAAProduct(prodCode string, runId string) map[string]interface{} {
	var csmRun models.CsmRun
	var paaResults []models.PAAResult
	var groups []models.GroupResults
	var res = make(map[string]interface{})
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Where("run_date = ? and measurement_type = ?", runId, "PAA").Find(&csmRun).Error }); err != nil {
		fmt.Println(err)
	}
	query := fmt.Sprintf("SELECT product_code, run_date, sum(premium_receipt) as premium_receipt,sum(nb_premium_receipt) as nb_premium_receipt,sum(total_premium_receipt) as total_premium_receipt, sum(acquisition_expenses) as acquisition_expenses , sum(earned_premium) as earned_premium,sum(insurance_revenue) as insurance_revenue,sum(incurred_claims) as incurred_claims,sum(claims_paid) as claims_paid, sum(incurred_expenses) as incurred_expenses, sum(paa_liability_remaining_coverage) as paa_liability_remaining_coverage,sum(deferred_acquisition_expenses) as deferred_acquisition_expenses,sum(amortised_acquisition_cost) as amortised_acquisition_cost ,sum(gmm_bel) as gmm_bel, sum(gmm_risk_adjustment) as gmm_risk_adjustment, sum(gmm_reserve) as gmm_reserve, sum(paa_loss_component) as paa_loss_component,sum(reinsurance_premium_paid) as reinsurance_premium_paid,sum(paa_reinsurance_lrc) as paa_reinsurance_lrc, sum(paa_reinsurance_dac) as paa_reinsurance_dac,sum(paa_loss_recovery_component) as paa_loss_recovery_component, sum(allocated_reinsurance_premium) as allocated_reinsurance_premium,sum(allocated_reinsurance_flat_commission) as allocated_reinsurance_flat_commission,sum(reinsurance_recovery) as reinsurance_recovery,sum(reinsurance_reinstatement_premium) as reinsurance_reinstatement_premium, sum(reinsurance_ultimate_loss_ratio) as reinsurance_ultimate_loss_ratio, sum(reinsurance_provisional_commission) as reinsurance_provisional_commission,sum(reinsurance_ultimate_commission) as reinsurance_ultimate_commission,sum(reinsurance_profit_commission) as reinsurance_profit_commission, sum(reinsurance_total_paid_to_cedant) as reinsurance_total_paid_to_cedant, sum(reinsurance_investment_component) as reinsurance_investment_component, sum(reinsurance_revenue) as reinsurance_revenue,sum(reinsurance_service_expense) as reinsurance_service_expense,sum(reinsurance_service_result) as reinsurance_service_result,sum(non_attributable_expenses) as non_attributable_expenses FROM paa_results where product_code = '%s' and run_date = '%s'  group by product_code", prodCode, runId)

	//query := fmt.Sprintf("SELECT product_code, run_date, sum(premium_receipt) as premium_receipt, sum(earned_premium) as earned_premium, sum(unearned_premium_reserve) as unearned_premium_reserve, sum(gmm_bel) as gmm_bel, sum(gmm_risk_adjustment) as gmm_risk_adjustment, sum(gmm_reserve) as gmm_reserve, sum(paa_loss_component) as paa_loss_component FROM paa_results where product_code = '%s' and run_date = '%s'  group by product_code", prodCode, runId)
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw(query).Scan(&paaResults).Error
	})
	if err != nil {
		fmt.Println(err)
	}
	res["steps"] = paaResults
	err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw("SELECT distinct ifrs17_group FROM paa_results where product_code = ?", prodCode).Scan(&groups).Error
	})
	if err != nil {
		fmt.Println(err)
	}
	res["groups"] = groups

	return res
}

func GetStepResultsForPAAProductGroup(prodCode string, runId string, ifrs17Group string) map[string]interface{} {
	var csmRun models.CsmRun
	var paaResults []models.PAAResult
	var res = make(map[string]interface{})
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Where("run_date = ? and measurement_type = ?", runId, "PAA").Find(&csmRun).Error }); err != nil {
		fmt.Println(err)
	}
	query := fmt.Sprintf("SELECT product_code, run_date,ifrs17_group, sum(premium_receipt) as premium_receipt,sum(nb_premium_receipt) as nb_premium_receipt,sum(total_premium_receipt) as total_premium_receipt, sum(acquisition_expenses) as acquisition_expenses , sum(earned_premium) as earned_premium,sum(insurance_revenue) as insurance_revenue, sum(incurred_claims) as incurred_claims,sum(claims_paid) as claims_paid, sum(incurred_expenses) as incurred_expenses,sum(paa_liability_remaining_coverage) as paa_liability_remaining_coverage,sum(deferred_acquisition_expenses) as deferred_acquisition_expenses,sum(amortised_acquisition_cost) as amortised_acquisition_cost, sum(gmm_bel) as gmm_bel, sum(gmm_risk_adjustment) as gmm_risk_adjustment, sum(gmm_reserve) as gmm_reserve, sum(paa_loss_component) as paa_loss_component,sum(reinsurance_premium_paid) as reinsurance_premium_paid,sum(paa_reinsurance_lrc) as paa_reinsurance_lrc, sum(paa_reinsurance_dac) as paa_reinsurance_dac,sum(paa_loss_recovery_component) as paa_loss_recovery_component, sum(allocated_reinsurance_premium) as allocated_reinsurance_premium,sum(allocated_reinsurance_flat_commission) as allocated_reinsurance_flat_commission,sum(reinsurance_recovery) as reinsurance_recovery,sum(reinsurance_reinstatement_premium) as reinsurance_reinstatement_premium, sum(reinsurance_ultimate_loss_ratio) as reinsurance_ultimate_loss_ratio, sum(reinsurance_provisional_commission) as reinsurance_provisional_commission,sum(reinsurance_ultimate_commission) as reinsurance_ultimate_commission,sum(reinsurance_profit_commission) as reinsurance_profit_commission, sum(reinsurance_total_paid_to_cedant) as reinsurance_total_paid_to_cedant, sum(reinsurance_investment_component) as reinsurance_investment_component, sum(reinsurance_revenue) as reinsurance_revenue,sum(reinsurance_service_expense) as reinsurance_service_expense,sum(reinsurance_service_result) as reinsurance_service_result,sum(non_attributable_expenses) as non_attributable_expenses FROM paa_results where product_code = '%s' and run_date = '%s' and ifrs17_group = '%s'  group by product_code", prodCode, runId, ifrs17Group)
	//query := fmt.Sprintf("SELECT product_code,run_date, ifrs17_group, sum(premium_receipt) as premium_receipt, sum(earned_premium) as earned_premium, sum(unearned_premium_reserve) as unearned_premium_reserve, sum(gmm_bel) as gmm_bel, sum(gmm_risk_adjustment) as gmm_risk_adjustment, sum(gmm_reserve) as gmm_reserve, sum(paa_loss_component) as paa_loss_component FROM paa_results where product_code = '%s' and run_date = '%s' and ifrs17_group = '%s'  group by product_code", prodCode, runId, ifrs17Group)
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw(query).Scan(&paaResults).Error
	})
	if err != nil {
		fmt.Println(err)
	}
	res["steps"] = paaResults

	return res
}

func GetStepResultsForPAASingleProductGroup(runId int, prodCode, ifrs17Group, runDate string) map[string]interface{} {
	var csmRun models.CsmRun
	var paaResults []models.PAAResult
	var res = make(map[string]interface{})
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Where("paa_run_id = ? and measurement_type = ? and run_date = ?", runId, "PAA", runDate).Find(&csmRun).Error }); err != nil {
		fmt.Println(err)
	}
	query := fmt.Sprintf("SELECT product_code, ifrs17_group, sum(premium_receipt) as premium_receipt,sum(nb_premium_receipt) as nb_premium_receipt,sum(total_premium_receipt) as total_premium_receipt ,sum(acquisition_expenses) as acquisition_expenses ,sum(insurance_revenue) as insurance_revenue,sum(incurred_claims) as incurred_claims,sum(claims_paid) as claims_paid, sum(incurred_expenses) as incurred_expenses ,sum(earned_premium) as earned_premium, sum(paa_liability_remaining_coverage) as paa_liability_remaining_coverage,sum(deferred_acquisition_expenses) as deferred_acquisition_expenses,sum(amortised_acquisition_cost) as amortised_acquisition_cost, sum(gmm_bel) as gmm_bel, sum(gmm_risk_adjustment) as gmm_risk_adjustment, sum(gmm_reserve) as gmm_reserve, sum(paa_loss_component) as paa_loss_component,sum(reinsurance_premium) as reinsurance_premium, sum(reinsurance_recovery) as reinsurance_recovery,sum(paa_reinsurance_lrc) as paa_reinsurance_lrc, sum(paa_reinsurance_dac) as paa_reinsurance_dac,sum(paa_initial_loss_recovery_component) as paa_initial_loss_recovery_component ,sum(paa_treaty1_initial_loss_recovery_component) as paa_treaty1_initial_loss_recovery_component,sum(paa_treaty2_initial_loss_recovery_component) as paa_treaty2_initial_loss_recovery_component,sum(paa_treaty3_initial_loss_recovery_component) as paa_treaty3_initial_loss_recovery_component,sum(paa_loss_recovery_component) as paa_loss_recovery_component,sum(paa_loss_recovery_unwind) as paa_loss_recovery_unwind,sum(paa_treaty1_loss_recovery_unwind) as paa_treaty1_loss_recovery_unwind,sum(paa_treaty2_loss_recovery_unwind) as paa_treaty2_loss_recovery_unwind,sum(paa_treaty3_loss_recovery_unwind) as paa_treaty3_loss_recovery_unwind ,sum(allocated_reinsurance_premium) as allocated_reinsurance_premium,sum(allocated_treaty1_premium) as allocated_treaty1_premium,sum(allocated_treaty2_premium) as allocated_treaty2_premium,sum(allocated_treaty3_premium) as allocated_treaty3_premium,sum(allocated_reinsurance_flat_commission) as allocated_reinsurance_flat_commission,sum(reinsurance_recovery) as reinsurance_recovery,sum(reinsurance_reinstatement_premium) as reinsurance_reinstatement_premium, sum(reinsurance_ultimate_loss_ratio) as reinsurance_ultimate_loss_ratio, sum(reinsurance_provisional_commission) as reinsurance_provisional_commission,sum(reinsurance_ultimate_commission) as reinsurance_ultimate_commission,sum(reinsurance_profit_commission) as reinsurance_profit_commission, sum(reinsurance_total_paid_to_cedant) as reinsurance_total_paid_to_cedant, sum(reinsurance_investment_component) as reinsurance_investment_component, sum(reinsurance_revenue) as reinsurance_revenue,sum(reinsurance_service_expense) as reinsurance_service_expense,sum(reinsurance_service_result) as reinsurance_service_result,sum(incurred_expenses) as incurred_expenses,sum(non_attributable_expenses) as non_attributable_expenses FROM paa_results where product_code = '%s' and csm_run_id = %d and ifrs17_group = '%s' and run_date = '%s'  group by product_code", prodCode, runId, ifrs17Group, csmRun.RunDate)
	//query := fmt.Sprintf("SELECT product_code, ifrs17_group, sum(premium_receipt) as premium_receipt, sum(earned_premium) as earned_premium, sum(unearned_premium_reserve) as unearned_premium_reserve, sum(gmm_bel) as gmm_bel, sum(gmm_risk_adjustment) as gmm_risk_adjustment, sum(gmm_reserve) as gmm_reserve, sum(paa_loss_component) as paa_loss_component FROM paa_results where product_code = '%s' and csm_run_id = %d and ifrs17_group = '%s'  group by product_code", prodCode, runId, ifrs17Group)
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw(query).Scan(&paaResults).Error
	})
	if err != nil {
		fmt.Println(err)
	}
	res["steps"] = paaResults

	return res
}

func GetAosStepResultsForProductByRunId(prodCode string, runId int) map[string]interface{} {
	var res = make(map[string]interface{})
	var agg []models.ProductStepResult
	var groups []models.GroupResults
	var csmRun models.CsmRun
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Where("run_id = ? and measurement_type = ?", runId, "GMM").Find(&csmRun).Error }); err != nil {
		fmt.Println(err)
	}
	query := fmt.Sprintf("SELECT product_code,  name, sum(bel) as bel, sum(bel_at12) as bel_at12, sum(risk_adjustment) as risk_adjustment, sum(risk_adjustment_at12) as risk_adjustment_at12, sum(best_estimate_liability_change) as best_estimate_liability_change, sum(risk_adjustment_change) as risk_adjustment_change, SUM(csm_change) as csm_change, SUM(liability_change) as liability_change, Sum(loss_component_change) as loss_component_change,sum(loss_component_unwind) as loss_component_unwind, SUM(pnl_change) as pnl_change, sum(csm_buildup) as csm_buildup,sum(risk_adjustment_buildup) as risk_adjustment_buildup, sum(loss_component_buildup) as loss_component_buildup, sum(sar_actual_claim_net_lc) as sar_actual_claim_net_lc,sum(sar_actual_expense_net_lc) as sar_actual_expense_net_lc, sum(expected_ra_net_of_lc) as expected_ra_net_of_lc, sum(lc_sar) as lc_sar,sum(csm_release_ratio) as csm_release_ratio, sum(csm_release) as csm_release, sum(discounted_coverage_units) as discounted_coverage_units, sum(sum_coverage_units) as sum_coverage_units, sum(interest_accretion_fac) as interest_accretion_fac,sum(bel_inflow) as bel_inflow,sum(bel_outflow) as bel_outflow, sum(bel_inflow_at12) as  bel_inflow_at12, sum(bel_outflow_at12) as bel_outflow_at12, sum(bel_inflow_change) as bel_inflow_change,sum(bel_outflow_change) as bel_outflow_change,sum(reinsurance_bel_outflow) as reinsurance_bel_outflow,sum(reinsurance_bel_inflow) as reinsurance_bel_inflow,sum(reinsurance_bel) as reinsurance_bel, sum(reinsurance_risk_adjustment) as reinsurance_risk_adjustment,sum(reinsurance_bel_outflow_at12) as reinsurance_bel_outflow_at12,sum(reinsurance_bel_inflow_at12) as reinsurance_bel_inflow_at12, sum(reinsurance_bel_at12) as reinsurance_bel_at12,sum(reinsurance_risk_adjustment_at12) as reinsurance_risk_adjustment_at12, sum(reins_bel_outflow_change) as reins_bel_outflow_change,sum(reins_bel_inflow_change) as reins_bel_inflow_change, sum(reins_bel_change) as reins_bel_change, sum(reins_ra_change) as reins_ra_change, sum(reins_csm_change) as reins_csm_change, sum(reins_csm_buildup) as reins_csm_buildup, sum(reins_pnl_change) as reins_pnl_change, sum(reinsurance_csm) as reinsurance_csm, sum(loss_recovery_component_change) as loss_recovery_component_change, sum(loss_recovery_component_buildup) as loss_recovery_component_buildup, sum(reinsurance_revenue) as reinsurance_revenue, sum(reinsurance_service_expense) as reinsurance_service_expense,sum(reinsurance_service_result) as reinsurance_service_result,sum(reins_expected_cash_outflow) as reins_expected_cash_outflow, sum(reins_expected_cash_inflow) as reins_expected_cash_inflow  FROM aos_step_results where product_code = '%s' and csm_run_id = %d  group by product_code, name", prodCode, runId)
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw(query).Scan(&agg).Error
	})
	if err != nil {
		fmt.Println(err)
	}
	res["steps"] = agg

	//Get IFRS17Groups....
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Raw("SELECT distinct ifrs17_group FROM scoped_aggregated_projections where product_code = ? and csm_run_id=?", prodCode, runId).Scan(&groups).Error }); err != nil {
		fmt.Println(err)
	}
	res["groups"] = groups
	res["run_settings"] = csmRun

	return res
}

func GetStepResultsForProduct(prodCode string, runDate string) map[string]interface{} {
	var res = make(map[string]interface{})
	var agg []models.ProductStepResult
	var groups []models.GroupResults
	var csmRun models.CsmRun
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Where("run_date = ? and measurement_type = ?", runDate, "GMM").Find(&csmRun).Error }); err != nil {
		fmt.Println(err)
	}
	query := fmt.Sprintf("SELECT product_code,  name,time, sum(bel) as bel, sum(bel_at12) as bel_at12, sum(risk_adjustment) as risk_adjustment, sum(risk_adjustment_at12) as risk_adjustment_at12, sum(best_estimate_liability_change) as best_estimate_liability_change, sum(risk_adjustment_change) as risk_adjustment_change, SUM(csm_change) as csm_change, SUM(liability_change) as liability_change, Sum(loss_component_change) as loss_component_change,sum(loss_component_unwind) as loss_component_unwind, SUM(pnl_change) as pnl_change, sum(csm_buildup) as csm_buildup,sum(risk_adjustment_buildup) as risk_adjustment_buildup, sum(loss_component_buildup) as loss_component_buildup, sum(sar_actual_claim_net_lc) as sar_actual_claim_net_lc,sum(sar_actual_expense_net_lc) as sar_actual_expense_net_lc, sum(expected_ra_net_of_lc) as expected_ra_net_of_lc, sum(lc_sar) as lc_sar,sum(actual_non_attributable_expenses) as actual_non_attributable_expenses, sum(csm_release_ratio) as csm_release_ratio,sum(csm_release) as csm_release, sum(discounted_coverage_units) as discounted_coverage_units, sum(sum_coverage_units) as sum_coverage_units, sum(interest_accretion_fac) as interest_accretion_fac, sum(bel_inflow) as bel_inflow, sum(bel_outflow) as bel_outflow, sum(bel_inflow_at12) as  bel_inflow_at12, sum(bel_outflow_at12) as bel_outflow_at12, sum(bel_inflow_change) as bel_inflow_change,sum(bel_outflow_change) as bel_outflow_change,sum(expected_cash_inflow) as expected_cash_inflow, sum(expected_cash_outflow) as expected_cash_outflow, sum(actual_premium) as actual_premium, sum(premium_debtor) as premium_debtor,sum(experience_premium_variance) as experience_premium_variance,sum(bel_outflow_excl_acquisition) as bel_outflow_excl_acquisition, sum(bel_acquisition_cost) as bel_acquisition_cost, sum(bel_outflow_excl_acquisition_at12) as bel_outflow_excl_acquisition_at12, sum(bel_acquisition_cost_at12) as bel_acquisition_cost_at12,sum(dac_change) as dac_change,sum(dac_buildup) as dac_buildup,sum(reinsurance_bel_outflow) as reinsurance_bel_outflow,sum(reinsurance_bel_inflow) as reinsurance_bel_inflow,sum(reinsurance_bel) as reinsurance_bel, sum(reinsurance_risk_adjustment) as reinsurance_risk_adjustment,sum(reinsurance_bel_outflow_at12) as reinsurance_bel_outflow_at12,sum(reinsurance_bel_inflow_at12) as reinsurance_bel_inflow_at12, sum(reinsurance_bel_at12) as reinsurance_bel_at12,sum(reinsurance_risk_adjustment_at12) as reinsurance_risk_adjustment_at12, sum(reins_bel_outflow_change) as reins_bel_outflow_change,sum(reins_bel_inflow_change) as reins_bel_inflow_change, sum(reins_bel_change) as reins_bel_change, sum(reins_ra_change) as reins_ra_change, sum(reins_csm_change) as reins_csm_change, sum(reins_csm_buildup) as reins_csm_buildup, sum(reins_pnl_change) as reins_pnl_change, sum(reinsurance_csm) as reinsurance_csm, sum(loss_recovery_component_change) as loss_recovery_component_change, sum(loss_recovery_component_buildup) as loss_recovery_component_buildup, sum(reinsurance_revenue) as reinsurance_revenue, sum(reinsurance_service_expense) as reinsurance_service_expense,sum(reinsurance_service_result) as reinsurance_service_result,sum(reins_expected_cash_outflow) as reins_expected_cash_outflow, sum(reins_expected_cash_inflow) as reins_expected_cash_inflow  FROM aos_step_results where product_code = '%s' and run_date = '%s'  group by product_code, name,time", prodCode, runDate)
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw(query).Scan(&agg).Error
	})
	if err != nil {
		fmt.Println(err)
	}
	res["steps"] = agg

	//Get IFRS17Groups....
	//DB.Raw("SELECT distinct ifrs17_group FROM scoped_aggregated_projections where product_code = ?", prodCode).Scan(&groups)
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Raw("SELECT distinct ifrs17_group FROM aos_step_results where product_code = ? and run_date = ? ", prodCode, runDate).Scan(&groups).Error }); err != nil {
		fmt.Println(err)
	}
	res["groups"] = groups

	return res
}

func GetStepResultsForAosProductByRunId(prodCode string, runId int, runDate string) map[string]interface{} {
	var res = make(map[string]interface{})
	var agg []models.ProductStepResult
	var groups []models.GroupResults
	var csmRun models.CsmRun
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Where("run_date = ? and measurement_type = ? and run_date = ?", runId, "GMM", runDate).Find(&csmRun).Error }); err != nil {
		fmt.Println(err)
	}
	query := fmt.Sprintf("SELECT product_code,  name,time, sum(bel) as bel, sum(bel_at12) as bel_at12, sum(risk_adjustment) as risk_adjustment, sum(risk_adjustment_at12) as risk_adjustment_at12, sum(best_estimate_liability_change) as best_estimate_liability_change, sum(risk_adjustment_change) as risk_adjustment_change, SUM(csm_change) as csm_change, SUM(liability_change) as liability_change, Sum(loss_component_change) as loss_component_change,sum(loss_component_unwind) as loss_component_unwind , SUM(pnl_change) as pnl_change, sum(csm_buildup) as csm_buildup,sum(risk_adjustment_buildup) as risk_adjustment_buildup, sum(loss_component_buildup) as loss_component_buildup, sum(sar_actual_claim_net_lc) as sar_actual_claim_net_lc,sum(sar_actual_expense_net_lc) as sar_actual_expense_net_lc, sum(expected_ra_net_of_lc) as expected_ra_net_of_lc, sum(lc_sar) as lc_sar,sum(actual_non_attributable_expenses) as actual_non_attributable_expenses,sum(csm_release_ratio) as csm_release_ratio, sum(csm_release) as csm_release, sum(discounted_coverage_units) as discounted_coverage_units, sum(sum_coverage_units) as sum_coverage_units, sum(interest_accretion_fac) as interest_accretion_fac,sum(bel_inflow) as bel_inflow,sum(bel_outflow) as bel_outflow, sum(bel_inflow_at12) as  bel_inflow_at12, sum(bel_outflow_at12) as bel_outflow_at12, sum(bel_inflow_change) as bel_inflow_change,sum(bel_outflow_change) as bel_outflow_change,sum(expected_cash_inflow) as expected_cash_inflow, sum(expected_cash_outflow) as expected_cash_outflow, sum(actual_premium) as actual_premium, sum(premium_debtor) as premium_debtor,sum(experience_premium_variance) as experience_premium_variance,sum(reinsurance_bel_outflow) as reinsurance_bel_outflow,sum(reinsurance_bel_inflow) as reinsurance_bel_inflow,sum(reinsurance_bel) as reinsurance_bel, sum(reinsurance_risk_adjustment) as reinsurance_risk_adjustment,sum(reinsurance_bel_outflow_at12) as reinsurance_bel_outflow_at12,sum(reinsurance_bel_inflow_at12) as reinsurance_bel_inflow_at12, sum(reinsurance_bel_at12) as reinsurance_bel_at12,sum(reinsurance_risk_adjustment_at12) as reinsurance_risk_adjustment_at12, sum(reins_bel_outflow_change) as reins_bel_outflow_change,sum(reins_bel_inflow_change) as reins_bel_inflow_change, sum(reins_bel_change) as reins_bel_change, sum(reins_ra_change) as reins_ra_change, sum(reins_csm_change) as reins_csm_change, sum(reins_csm_buildup) as reins_csm_buildup, sum(reins_pnl_change) as reins_pnl_change, sum(reinsurance_csm) as reinsurance_csm, sum(loss_recovery_component_change) as loss_recovery_component_change, sum(loss_recovery_component_buildup) as loss_recovery_component_buildup, sum(reinsurance_revenue) as reinsurance_revenue, sum(reinsurance_service_expense) as reinsurance_service_expense,sum(reinsurance_service_result) as reinsurance_service_result,sum(reins_expected_cash_outflow) as reins_expected_cash_outflow, sum(reins_expected_cash_inflow) as reins_expected_cash_inflow  FROM aos_step_results where product_code = '%s' and csm_run_id = %d  group by product_code, name,time", prodCode, runId)
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw(query).Scan(&agg).Error
	})
	if err != nil {
		fmt.Println(err)
	}
	res["steps"] = agg

	//Get IFRS17Groups....
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Raw("SELECT distinct ifrs17_group FROM aos_step_results where product_code = ? and csm_run_id = ?", prodCode, runId).Scan(&groups).Error }); err != nil {
		fmt.Println(err)
	}
	res["groups"] = groups

	return res
}

func GetRAFactors() ([]models.RiskAdjustmentFactor, error) {
	var factors []models.RiskAdjustmentFactor

	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Find(&factors).Error
	})
	if err != nil {
		return factors, err
	}
	return factors, nil
}

// GetRiskAdjustmentFactorVersions returns distinct versions available for a given year
func GetRiskAdjustmentFactorVersions(year int) ([]string, error) {
	var versions []string
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Model(&models.RiskAdjustmentFactor{}).
			Select("distinct version").
			Where("year = ?", year).
			Order("version").
			Pluck("version", &versions).Error
	})
	if err != nil {
		return nil, err
	}
	return versions, nil
}

func GetStepResultsForGroup(runDate, ifrs17Group string) map[string]interface{} {
	var res = make(map[string]interface{})
	var agg []models.ProductStepResult
	//var groups []models.GroupResults
	query := fmt.Sprintf("SELECT product_code, ifrs17_group,  name, time, sum(bel) as bel, sum(bel_at12) as bel_at12, sum(risk_adjustment) as risk_adjustment, sum(risk_adjustment_at12) as risk_adjustment_at12, sum(best_estimate_liability_change) as best_estimate_liability_change, sum(risk_adjustment_change) as risk_adjustment_change, SUM(csm_change) as csm_change, SUM(liability_change) as liability_change, Sum(loss_component_change) as loss_component_change,sum(loss_component_unwind) as loss_component_unwind, SUM(pnl_change) as pnl_change, sum(csm_buildup) as csm_buildup,sum(risk_adjustment_buildup) as risk_adjustment_buildup,sum(loss_component_unwind) as loss_component_unwind ,sum(loss_component_buildup) as loss_component_buildup, sum(sar_actual_claim_net_lc) as sar_actual_claim_net_lc,sum(sar_actual_expense_net_lc) as sar_actual_expense_net_lc, sum(expected_ra_net_of_lc) as expected_ra_net_of_lc, sum(lc_sar) as lc_sar,sum(actual_non_attributable_expenses) as actual_non_attributable_expenses, sum(csm_release_ratio) as csm_release_ratio,sum(csm_release) as csm_release, sum(discounted_coverage_units) as discounted_coverage_units, sum(sum_coverage_units) as sum_coverage_units, sum(interest_accretion_fac) as interest_accretion_fac, sum(bel_inflow) as bel_inflow, sum(bel_outflow) as bel_outflow, sum(bel_inflow_at12) as  bel_inflow_at12, sum(bel_outflow_at12) as bel_outflow_at12, sum(bel_inflow_change) as bel_inflow_change,sum(bel_outflow_change) as bel_outflow_change,sum(expected_cash_inflow) as expected_cash_inflow, sum(expected_cash_outflow) as expected_cash_outflow, sum(actual_premium) as actual_premium, sum(premium_debtor) as premium_debtor,sum(experience_premium_variance) as experience_premium_variance, sum(bel_outflow_excl_acquisition) as bel_outflow_excl_acquisition, sum(bel_acquisition_cost) as bel_acquisition_cost, sum(bel_outflow_excl_acquisition_at12) as bel_outflow_excl_acquisition_at12, sum(bel_acquisition_cost_at12) as bel_acquisition_cost_at12,sum(dac_change) as dac_change,sum(dac_buildup) as dac_buildup,sum(reinsurance_bel_outflow) as reinsurance_bel_outflow,sum(reinsurance_bel_inflow) as reinsurance_bel_inflow,sum(reinsurance_bel) as reinsurance_bel, sum(reinsurance_risk_adjustment) as reinsurance_risk_adjustment,sum(reinsurance_bel_outflow_at12) as reinsurance_bel_outflow_at12,sum(reinsurance_bel_inflow_at12) as reinsurance_bel_inflow_at12, sum(reinsurance_bel_at12) as reinsurance_bel_at12,sum(reinsurance_risk_adjustment_at12) as reinsurance_risk_adjustment_at12, sum(reins_bel_outflow_change) as reins_bel_outflow_change,sum(reins_bel_inflow_change) as reins_bel_inflow_change, sum(reins_bel_change) as reins_bel_change, sum(reins_ra_change) as reins_ra_change, sum(reins_csm_change) as reins_csm_change, sum(reins_csm_buildup) as reins_csm_buildup, sum(reins_pnl_change) as reins_pnl_change, sum(reinsurance_csm) as reinsurance_csm, sum(loss_recovery_component_change) as loss_recovery_component_change, sum(loss_recovery_component_buildup) as loss_recovery_component_buildup, sum(reinsurance_revenue) as reinsurance_revenue, sum(reinsurance_service_expense) as reinsurance_service_expense,sum(reinsurance_service_result) as reinsurance_service_result,sum(reins_expected_cash_outflow) as reins_expected_cash_outflow, sum(reins_expected_cash_inflow) as reins_expected_cash_inflow  FROM aos_step_results where ifrs17_group = '%s' and run_date = '%s'  group by product_code, ifrs17_group, name,time", ifrs17Group, runDate)
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw(query).Scan(&agg).Error
	})
	if err != nil {
		fmt.Println(err)
	}
	res["steps"] = agg

	//Get IFRS17Groups....
	//DB.Raw("SELECT distinct ifrs17_group FROM scoped_aggregated_projections where product_code = ?", ifrs17Group).Scan(&groups)
	//res["groups"] = groups

	return res
}

func GetStepResultsForAosProductGroupByRunId(runId int, runDate, prodCode string, ifrs17Group string) map[string]interface{} {
	var res = make(map[string]interface{})
	var agg []models.ProductStepResult
	var groups []models.GroupResults
	query := fmt.Sprintf("SELECT product_code, ifrs17_group,  name,time, sum(bel) as bel, sum(bel_at12) as bel_at12, sum(risk_adjustment) as risk_adjustment, sum(risk_adjustment_at12) as risk_adjustment_at12, sum(best_estimate_liability_change) as best_estimate_liability_change, sum(risk_adjustment_change) as risk_adjustment_change, SUM(csm_change) as csm_change, SUM(liability_change) as liability_change, Sum(loss_component_change) as loss_component_change,sum(loss_component_unwind) as loss_component_unwind, SUM(pnl_change) as pnl_change, sum(csm_buildup) as csm_buildup,sum(risk_adjustment_buildup) as risk_adjustment_buildup,sum(loss_component_unwind) as loss_component_unwind, sum(loss_component_buildup) as loss_component_buildup, sum(sar_actual_claim_net_lc) as sar_actual_claim_net_lc,sum(sar_actual_expense_net_lc) as sar_actual_expense_net_lc, sum(expected_ra_net_of_lc) as expected_ra_net_of_lc, sum(lc_sar) as lc_sar,sum(actual_non_attributable_expenses) as actual_non_attributable_expenses, sum(csm_release_ratio) as csm_release_ratio,sum(csm_release) as csm_release, sum(discounted_coverage_units) as discounted_coverage_units, sum(sum_coverage_units) as sum_coverage_units, sum(interest_accretion_fac) as interest_accretion_fac, sum(bel_inflow) as bel_inflow, sum(bel_outflow) as bel_outflow, sum(bel_inflow_at12) as  bel_inflow_at12, sum(bel_outflow_at12) as bel_outflow_at12, sum(bel_inflow_change) as bel_inflow_change,sum(bel_outflow_change) as bel_outflow_change,sum(expected_cash_inflow) as expected_cash_inflow, sum(expected_cash_outflow) as expected_cash_outflow, sum(actual_premium) as actual_premium, sum(premium_debtor) as premium_debtor,sum(experience_premium_variance) as experience_premium_variance,sum(bel_outflow_excl_acquisition) as bel_outflow_excl_acquisition, sum(bel_acquisition_cost) as bel_acquisition_cost, sum(bel_outflow_excl_acquisition_at12) as bel_outflow_excl_acquisition_at12, sum(bel_acquisition_cost_at12) as bel_acquisition_cost_at12,sum(dac_change) as dac_change,sum(dac_buildup) as dac_buildup, sum(bel_outflow_excl_lc) as bel_outflow_excl_lc,sum(risk_adj_excl_lc) as risk_adj_excl_lc, sum(bel_outflow_lc) as bel_outflow_lc, sum(risk_adj_lc) as risk_adj_lc,sum(reinsurance_bel_outflow) as reinsurance_bel_outflow,sum(reinsurance_bel_inflow) as reinsurance_bel_inflow,sum(reinsurance_bel) as reinsurance_bel, sum(reinsurance_risk_adjustment) as reinsurance_risk_adjustment,sum(reinsurance_bel_outflow_at12) as reinsurance_bel_outflow_at12,sum(reinsurance_bel_inflow_at12) as reinsurance_bel_inflow_at12, sum(reinsurance_bel_at12) as reinsurance_bel_at12,sum(reinsurance_risk_adjustment_at12) as reinsurance_risk_adjustment_at12, sum(reins_bel_outflow_change) as reins_bel_outflow_change,sum(reins_bel_inflow_change) as reins_bel_inflow_change, sum(reins_bel_change) as reins_bel_change, sum(reins_ra_change) as reins_ra_change, sum(reins_csm_change) as reins_csm_change, sum(reins_csm_buildup) as reins_csm_buildup, sum(reins_pnl_change) as reins_pnl_change, sum(reinsurance_csm) as reinsurance_csm, sum(loss_recovery_component_change) as loss_recovery_component_change, sum(loss_recovery_component_buildup) as loss_recovery_component_buildup, sum(reinsurance_revenue) as reinsurance_revenue, sum(reinsurance_service_expense) as reinsurance_service_expense,sum(reinsurance_service_result) as reinsurance_service_result,sum(reins_expected_cash_outflow) as reins_expected_cash_outflow, sum(reins_expected_cash_inflow) as reins_expected_cash_inflow  FROM aos_step_results where ifrs17_group = '%s' and product_code = '%s' and csm_run_id = %d and run_date = '%s'  group by product_code, ifrs17_group, name,time", ifrs17Group, prodCode, runId, runDate)
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw(query).Scan(&agg).Error
	})
	if err != nil {
		fmt.Println(err)
	}
	res["steps"] = agg

	//Get IFRS17Groups....
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Raw("SELECT distinct ifrs17_group FROM scoped_aggregated_projections where product_code = ?", ifrs17Group).Scan(&groups).Error }); err != nil {
		fmt.Println(err)
	}
	res["groups"] = groups

	return res
}

func GetAosStepResultsForAllProducts(runDate string) map[string]interface{} {
	var res = make(map[string]interface{})
	var agg []models.AOSStepResult
	var prodList []models.ProductList
	query := fmt.Sprintf("SELECT name,time, sum(bel) as bel, sum(bel_at12) as bel_at12, sum(risk_adjustment) as risk_adjustment, sum(risk_adjustment_at12) as risk_adjustment_at12, sum(best_estimate_liability_change) as best_estimate_liability_change, sum(risk_adjustment_change) as risk_adjustment_change, SUM(csm_change) as csm_change, SUM(liability_change) as liability_change, Sum(loss_component_change) as loss_component_change,sum(loss_component_unwind) as loss_component_unwind, SUM(pnl_change) as pnl_change, sum(csm_buildup) as csm_buildup,sum(risk_adjustment_buildup) as risk_adjustment_buildup,sum(loss_component_unwind) as loss_component_unwind ,sum(loss_component_buildup) as loss_component_buildup, sum(sar_actual_claim_net_lc) as sar_actual_claim_net_lc,sum(sar_actual_expense_net_lc) as sar_actual_expense_net_lc, sum(expected_ra_net_of_lc) as expected_ra_net_of_lc, sum(lc_sar) as lc_sar,sum(actual_non_attributable_expenses) as actual_non_attributable_expenses, sum(csm_release_ratio) as csm_release_ratio,sum(csm_release) as csm_release, sum(discounted_coverage_units) as discounted_coverage_units, sum(sum_coverage_units) as sum_coverage_units, sum(interest_accretion_fac) as interest_accretion_fac, sum(bel_inflow) as bel_inflow, sum(bel_outflow) as bel_outflow, sum(bel_inflow_at12) as  bel_inflow_at12, sum(bel_outflow_at12) as bel_outflow_at12, sum(bel_inflow_change) as bel_inflow_change,sum(bel_outflow_change) as bel_outflow_change,sum(expected_cash_inflow) as expected_cash_inflow, sum(expected_cash_outflow) as expected_cash_outflow,sum(actual_premium) as actual_premium,sum(premium_debtor) as premium_debtor,sum(experience_premium_variance) as experience_premium_variance,sum(bel_outflow_excl_acquisition) as bel_outflow_excl_acquisition, sum(bel_acquisition_cost) as bel_acquisition_cost, sum(bel_outflow_excl_acquisition_at12) as bel_outflow_excl_acquisition_at12, sum(bel_acquisition_cost_at12) as bel_acquisition_cost_at12,sum(dac_change) as dac_change,sum(dac_buildup) as dac_buildup,sum(bel_outflow_excl_lc) as bel_outflow_excl_lc,sum(risk_adj_excl_lc) as risk_adj_excl_lc, sum(bel_outflow_lc) as bel_outflow_lc, sum(risk_adj_lc) as risk_adj_lc, sum(reinsurance_bel_outflow) as reinsurance_bel_outflow,sum(reinsurance_bel_inflow) as reinsurance_bel_inflow,sum(reinsurance_bel) as reinsurance_bel, sum(reinsurance_risk_adjustment) as reinsurance_risk_adjustment,sum(reinsurance_bel_outflow_at12) as reinsurance_bel_outflow_at12,sum(reinsurance_bel_inflow_at12) as reinsurance_bel_inflow_at12, sum(reinsurance_bel_at12) as reinsurance_bel_at12,sum(reinsurance_risk_adjustment_at12) as reinsurance_risk_adjustment_at12, sum(reins_bel_outflow_change) as reins_bel_outflow_change,sum(reins_bel_inflow_change) as reins_bel_inflow_change, sum(reins_bel_change) as reins_bel_change, sum(reins_ra_change) as reins_ra_change, sum(reins_csm_change) as reins_csm_change, sum(reins_csm_buildup) as reins_csm_buildup, sum(reins_pnl_change) as reins_pnl_change, sum(reinsurance_csm) as reinsurance_csm, sum(loss_recovery_component_change) as loss_recovery_component_change, sum(loss_recovery_component_buildup) as loss_recovery_component_buildup, sum(reinsurance_revenue) as reinsurance_revenue, sum(reinsurance_service_expense) as reinsurance_service_expense,sum(reinsurance_service_result) as reinsurance_service_result,sum(reins_expected_cash_outflow) as reins_expected_cash_outflow, sum(reins_expected_cash_inflow) as reins_expected_cash_inflow FROM aos_step_results where run_date = '%s'  group by name,time", runDate)
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw(query).Scan(&agg).Error
	})
	if err != nil {
		fmt.Println(err)
	}

	res["steps"] = agg
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Raw("SELECT distinct product_code FROM aos_step_results where run_date = ?", runDate).Scan(&prodList).Error }); err != nil {
		fmt.Println(err)
	}
	res["products"] = prodList

	return res
}

func GetAosStepResultsForAllProductsByRunId(runId int, runDate string) map[string]interface{} {
	var res = make(map[string]interface{})
	var agg []models.AOSStepResult
	var prodList []models.ProductList
	var csmRun models.CsmRun
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Where("id = ? and run_date = ?", runId, runDate).First(&csmRun).Error }); err != nil {
		fmt.Println(err)
	}
	query := fmt.Sprintf("SELECT name,time, sum(bel) as bel, sum(bel_at12) as bel_at12, sum(risk_adjustment) as risk_adjustment, sum(risk_adjustment_at12) as risk_adjustment_at12, sum(best_estimate_liability_change) as best_estimate_liability_change, sum(risk_adjustment_change) as risk_adjustment_change, SUM(csm_change) as csm_change, SUM(liability_change) as liability_change, Sum(loss_component_change) as loss_component_change,sum(loss_component_unwind) as loss_component_unwind, SUM(pnl_change) as pnl_change, sum(csm_buildup) as csm_buildup,sum(risk_adjustment_buildup) as risk_adjustment_buildup,sum(loss_component_buildup) as loss_component_buildup, sum(sar_actual_claim_net_lc) as sar_actual_claim_net_lc,sum(sar_actual_expense_net_lc) as sar_actual_expense_net_lc, sum(expected_ra_net_of_lc) as expected_ra_net_of_lc, sum(lc_sar) as lc_sar,sum(actual_non_attributable_expenses) as actual_non_attributable_expenses, sum(csm_release_ratio) as csm_release_ratio,sum(csm_release) as csm_release, sum(discounted_coverage_units) as discounted_coverage_units, sum(sum_coverage_units) as sum_coverage_units, sum(interest_accretion_fac) as interest_accretion_fac, sum(bel_inflow) as bel_inflow, sum(bel_outflow) as bel_outflow, sum(bel_inflow_at12) as  bel_inflow_at12, sum(bel_outflow_at12) as bel_outflow_at12, sum(bel_inflow_change) as bel_inflow_change,sum(bel_outflow_change) as bel_outflow_change,sum(expected_cash_inflow) as expected_cash_inflow, sum(expected_cash_outflow) as expected_cash_outflow,sum(actual_premium) as actual_premium,sum(premium_debtor) as premium_debtor,sum(experience_premium_variance) as experience_premium_variance, sum(bel_outflow_excl_acquisition) as bel_outflow_excl_acquisition, sum(bel_acquisition_cost) as bel_acquisition_cost, sum(bel_outflow_excl_acquisition_at12) as bel_outflow_excl_acquisition_at12, sum(bel_acquisition_cost_at12) as bel_acquisition_cost_at12,sum(dac_change) as dac_change,sum(dac_buildup) as dac_buildup, sum(bel_outflow_excl_lc) as bel_outflow_excl_lc,sum(risk_adj_excl_lc) as risk_adj_excl_lc, sum(bel_outflow_lc) as bel_outflow_lc, sum(risk_adj_lc) as risk_adj_lc,sum(reinsurance_bel_outflow) as reinsurance_bel_outflow,sum(reinsurance_bel_inflow) as reinsurance_bel_inflow,sum(reinsurance_bel) as reinsurance_bel, sum(reinsurance_risk_adjustment) as reinsurance_risk_adjustment,sum(reinsurance_bel_outflow_at12) as reinsurance_bel_outflow_at12,sum(reinsurance_bel_inflow_at12) as reinsurance_bel_inflow_at12, sum(reinsurance_bel_at12) as reinsurance_bel_at12,sum(reinsurance_risk_adjustment_at12) as reinsurance_risk_adjustment_at12, sum(reins_bel_outflow_change) as reins_bel_outflow_change,sum(reins_bel_inflow_change) as reins_bel_inflow_change, sum(reins_bel_change) as reins_bel_change, sum(reins_ra_change) as reins_ra_change, sum(reins_csm_change) as reins_csm_change, sum(reins_csm_buildup) as reins_csm_buildup, sum(reins_pnl_change) as reins_pnl_change, sum(reinsurance_csm) as reinsurance_csm, sum(loss_recovery_component_change) as loss_recovery_component_change, sum(loss_recovery_component_buildup) as loss_recovery_component_buildup, sum(reinsurance_revenue) as reinsurance_revenue, sum(reinsurance_service_expense) as reinsurance_service_expense,sum(reinsurance_service_result) as reinsurance_service_result,sum(reins_expected_cash_outflow) as reins_expected_cash_outflow, sum(reins_expected_cash_inflow) as reins_expected_cash_inflow  FROM aos_step_results where csm_run_id = %d  group by name,time", runId)
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw(query).Scan(&agg).Error
	})
	if err != nil {
		fmt.Println(err)
	}

	res["run"] = csmRun
	res["steps"] = agg
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Raw("SELECT distinct product_code FROM aos_step_results where csm_run_id = ?", runId).Scan(&prodList).Error }); err != nil {
		fmt.Println(err)
	}
	res["products"] = prodList
	res["run_settings"] = csmRun
	return res
}

func GetAosStepResultsForAllProductsForDownstreamCalcs(runDate string) []models.AOSStepResult {
	var agg []models.AOSStepResult
	query := fmt.Sprintf("SELECT name, sum(bel) as bel, sum(bel_at12) as bel_at12, sum(risk_adjustment) as risk_adjustment, sum(risk_adjustment_at12) as risk_adjustment_at12, sum(best_estimate_liability_change) as best_estimate_liability_change, sum(risk_adjustment_change) as risk_adjustment_change, SUM(csm_change) as csm_change, SUM(liability_change) as liability_change, Sum(loss_component_change) as loss_component_change, SUM(pnl_change) as pnl_change, sum(csm_buildup) as csm_buildup,sum(risk_adjustment_buildup) as risk_adjustment_buildup,sum(loss_component_unwind) as loss_component_unwind, sum(loss_component_buildup) as loss_component_buildup, sum(sar_actual_claim_net_lc) as sar_actual_claim_net_lc, sum(expected_ra_net_of_lc) as expected_ra_net_of_lc, sum(lc_sar) as lc_sar, sum(csm_release_ratio) as csm_release_ratio,sum(csm_release) as csm_release, sum(discounted_coverage_units) as discounted_coverage_units, sum(sum_coverage_units) as sum_coverage_units, sum(interest_accretion_fac) as interest_accretion_fac, sum(bel_inflow) as bel_inflow, sum(bel_outflow) as bel_outflow, sum(bel_inflow_at12) as  bel_inflow_at12, sum(bel_outflow_at12) as bel_outflow_at12, sum(bel_inflow_change) as bel_inflow_change,sum(bel_outflow_change) as bel_outflow_change,sum(expected_cash_inflow) as expected_cash_inflow, sum(expected_cash_outflow) as expected_cash_outflow,sum(actual_premium) as actual_premium,sum(experience_premium_variance) as experience_premium_variance, sum(bel_outflow_excl_acquisition) as bel_outflow_excl_acquisition, sum(bel_acquisition_cost) as bel_acquisition_cost, sum(bel_outflow_excl_acquisition_at12) as bel_outflow_excl_acquisition_at12, sum(bel_acquisition_cost_at12) as bel_acquisition_cost_at12,sum(dac_change) as dac_change,sum(dac_buildup) as dac_buildup, sum(bel_outflow_excl_lc) as bel_outflow_excl_lc,sum(risk_adj_excl_lc) as risk_adj_excl_lc, sum(bel_outflow_lc) as bel_outflow_lc, sum(risk_adj_lc) as risk_adj_lc,sum(reinsurance_bel_outflow) as reinsurance_bel_outflow,sum(reinsurance_bel_inflow) as reinsurance_bel_inflow,sum(reinsurance_bel) as reinsurance_bel, sum(reinsurance_risk_adjustment) as reinsurance_risk_adjustment,sum(reinsurance_bel_outflow_at12) as reinsurance_bel_outflow_at12,sum(reinsurance_bel_inflow_at12) as reinsurance_bel_inflow_at12, sum(reinsurance_bel_at12) as reinsurance_bel_at12,sum(reinsurance_risk_adjustment_at12) as reinsurance_risk_adjustment_at12, sum(reins_bel_outflow_change) as reins_bel_outflow_change,sum(reins_bel_inflow_change) as reins_bel_inflow_change, sum(reins_bel_change) as reins_bel_change, sum(reins_ra_change) as reins_ra_change, sum(reins_csm_change) as reins_csm_change, sum(reins_csm_buildup) as reins_csm_buildup, sum(reins_pnl_change) as reins_pnl_change, sum(reinsurance_csm) as reinsurance_csm, sum(loss_recovery_component_change) as loss_recovery_component_change, sum(loss_recovery_component_buildup) as loss_recovery_component_buildup, sum(reinsurance_revenue) as reinsurance_revenue, sum(reinsurance_service_expense) as reinsurance_service_expense,sum(reinsurance_service_result) as reinsurance_service_result,sum(reins_expected_cash_outflow) as reins_expected_cash_outflow, sum(reins_expected_cash_inflow) as reins_expected_cash_inflow  FROM aos_step_results where run_date = '%s'  group by name", runDate)
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw(query).Scan(&agg).Error
	})
	if err != nil {
		fmt.Println(err)
	}

	return agg
}

func GetAosStepResultsForOneProductForDownstreamCalcs(runDate, prodCode string) []models.AOSStepResult {
	var agg []models.AOSStepResult
	query := fmt.Sprintf("SELECT product_code,name, sum(bel) as bel, sum(bel_at12) as bel_at12, sum(risk_adjustment) as risk_adjustment, sum(risk_adjustment_at12) as risk_adjustment_at12, sum(best_estimate_liability_change) as best_estimate_liability_change, sum(risk_adjustment_change) as risk_adjustment_change, SUM(csm_change) as csm_change, SUM(liability_change) as liability_change, Sum(loss_component_change) as loss_component_change, SUM(pnl_change) as pnl_change, sum(csm_buildup) as csm_buildup,sum(risk_adjustment_buildup) as risk_adjustment_buildup,sum(loss_component_unwind) as loss_component_unwind, sum(loss_component_buildup) as loss_component_buildup, sum(sar_actual_claim_net_lc) as sar_actual_claim_net_lc, sum(expected_ra_net_of_lc) as expected_ra_net_of_lc, sum(lc_sar) as lc_sar, sum(csm_release_ratio) as csm_release_ratio,sum(csm_release) as csm_release, sum(discounted_coverage_units) as discounted_coverage_units, sum(sum_coverage_units) as sum_coverage_units, sum(interest_accretion_fac) as interest_accretion_fac, sum(bel_inflow) as bel_inflow, sum(bel_outflow) as bel_outflow, sum(bel_inflow_at12) as  bel_inflow_at12, sum(bel_outflow_at12) as bel_outflow_at12, sum(bel_inflow_change) as bel_inflow_change,sum(bel_outflow_change) as bel_outflow_change,sum(expected_cash_inflow) as expected_cash_inflow, sum(expected_cash_outflow) as expected_cash_outflow, sum(actual_premium) as actual_premium,sum(experience_premium_variance) as experience_premium_variance,sum(bel_outflow_excl_acquisition) as bel_outflow_excl_acquisition, sum(bel_acquisition_cost) as bel_acquisition_cost, sum(bel_outflow_excl_acquisition_at12) as bel_outflow_excl_acquisition_at12, sum(bel_acquisition_cost_at12) as bel_acquisition_cost_at12,sum(dac_change) as dac_change,sum(dac_buildup) as dac_buildup, sum(bel_outflow_excl_lc) as bel_outflow_excl_lc,sum(risk_adj_excl_lc) as risk_adj_excl_lc, sum(bel_outflow_lc) as bel_outflow_lc, sum(risk_adj_lc) as risk_adj_lc,sum(reinsurance_bel_outflow) as reinsurance_bel_outflow,sum(reinsurance_bel_inflow) as reinsurance_bel_inflow,sum(reinsurance_bel) as reinsurance_bel, sum(reinsurance_risk_adjustment) as reinsurance_risk_adjustment,sum(reinsurance_bel_outflow_at12) as reinsurance_bel_outflow_at12,sum(reinsurance_bel_inflow_at12) as reinsurance_bel_inflow_at12, sum(reinsurance_bel_at12) as reinsurance_bel_at12,sum(reinsurance_risk_adjustment_at12) as reinsurance_risk_adjustment_at12, sum(reins_bel_outflow_change) as reins_bel_outflow_change,sum(reins_bel_inflow_change) as reins_bel_inflow_change, sum(reins_bel_change) as reins_bel_change, sum(reins_ra_change) as reins_ra_change, sum(reins_csm_change) as reins_csm_change, sum(reins_csm_buildup) as reins_csm_buildup, sum(reins_pnl_change) as reins_pnl_change, sum(reinsurance_csm) as reinsurance_csm, sum(loss_recovery_component_change) as loss_recovery_component_change, sum(loss_recovery_component_buildup) as loss_recovery_component_buildup, sum(reinsurance_revenue) as reinsurance_revenue, sum(reinsurance_service_expense) as reinsurance_service_expense,sum(reinsurance_service_result) as reinsurance_service_result,sum(reins_expected_cash_outflow) as reins_expected_cash_outflow, sum(reins_expected_cash_inflow) as reins_expected_cash_inflow  FROM aos_step_results where run_date = '%s' and product_code = '%s'  group by name", runDate, prodCode)
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw(query).Scan(&agg).Error
	})
	if err != nil {
		fmt.Println(err)
	}

	return agg
}

func GetPaaStepResultsForProductByRunId(runId int, prodCode, runDate string) map[string]interface{} {
	var csmRun models.CsmRun
	var paaResults []models.PAAResult
	var groups []models.GroupResults
	var res = make(map[string]interface{})
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Where("paa_run_id = ? and measurement_type = ? and run_date = ?", runId, "PAA", runDate).Find(&csmRun).Error }); err != nil {
		fmt.Println(err)
	}
	//query := fmt.Sprintf("SELECT product_code, sum(premium_receipt) as premium_receipt, sum(nb_premium_receipt) as nb_premium_receipt, sum(earned_premium) as earned_premium, sum(unearned_premium_reserve) as unearned_premium_reserve, sum(gmm_bel) as gmm_bel, sum(gmm_risk_adjustment) as gmm_risk_adjustment, sum(gmm_reserve) as gmm_reserve, sum(paa_loss_component) as paa_loss_component FROM paa_results where product_code = '%s' and csm_run_id = %d  group by product_code", prodCode, runId)
	query := fmt.Sprintf("SELECT product_code, sum(premium_receipt) as premium_receipt,sum(nb_premium_receipt) as nb_premium_receipt,sum(total_premium_receipt) as total_premium_receipt ,sum(acquisition_expenses) as acquisition_expenses ,sum(insurance_revenue) as insurance_revenue,sum(incurred_claims) as incurred_claims,sum(claims_paid) as claims_paid, sum(incurred_expenses) as incurred_expenses ,sum(earned_premium) as earned_premium, sum(paa_liability_remaining_coverage) as paa_liability_remaining_coverage,sum(deferred_acquisition_expenses) as deferred_acquisition_expenses,sum(amortised_acquisition_cost) as amortised_acquisition_cost, sum(gmm_bel) as gmm_bel, sum(gmm_risk_adjustment) as gmm_risk_adjustment, sum(gmm_reserve) as gmm_reserve, sum(paa_loss_component) as paa_loss_component,sum(reinsurance_premium) as reinsurance_premium, sum(reinsurance_recovery) as reinsurance_recovery,sum(paa_reinsurance_lrc) as paa_reinsurance_lrc, sum(paa_reinsurance_dac) as paa_reinsurance_dac,sum(paa_initial_loss_recovery_component) as paa_initial_loss_recovery_component ,sum(paa_treaty1_initial_loss_recovery_component) as paa_treaty1_initial_loss_recovery_component,sum(paa_treaty2_initial_loss_recovery_component) as paa_treaty2_initial_loss_recovery_component,sum(paa_treaty3_initial_loss_recovery_component) as paa_treaty3_initial_loss_recovery_component,sum(paa_loss_recovery_component) as paa_loss_recovery_component,sum(paa_loss_recovery_unwind) as paa_loss_recovery_unwind,sum(paa_treaty1_loss_recovery_unwind) as paa_treaty1_loss_recovery_unwind,sum(paa_treaty2_loss_recovery_unwind) as paa_treaty2_loss_recovery_unwind,sum(paa_treaty3_loss_recovery_unwind) as paa_treaty3_loss_recovery_unwind ,sum(allocated_reinsurance_premium) as allocated_reinsurance_premium,sum(allocated_treaty1_premium) as allocated_treaty1_premium,sum(allocated_treaty2_premium) as allocated_treaty2_premium,sum(allocated_treaty3_premium) as allocated_treaty3_premium,sum(allocated_reinsurance_flat_commission) as allocated_reinsurance_flat_commission,sum(reinsurance_recovery) as reinsurance_recovery,sum(reinsurance_reinstatement_premium) as reinsurance_reinstatement_premium, sum(reinsurance_ultimate_loss_ratio) as reinsurance_ultimate_loss_ratio, sum(reinsurance_provisional_commission) as reinsurance_provisional_commission,sum(reinsurance_ultimate_commission) as reinsurance_ultimate_commission,sum(reinsurance_profit_commission) as reinsurance_profit_commission, sum(reinsurance_total_paid_to_cedant) as reinsurance_total_paid_to_cedant, sum(reinsurance_investment_component) as reinsurance_investment_component, sum(reinsurance_revenue) as reinsurance_revenue,sum(reinsurance_service_expense) as reinsurance_service_expense,sum(reinsurance_service_result) as reinsurance_service_result,sum(incurred_expenses) as incurred_expenses,sum(non_attributable_expenses) as non_attributable_expenses FROM paa_results where product_code = '%s' and csm_run_id = %d and run_date = '%s' group by product_code", prodCode, runId, csmRun.RunDate)

	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw(query).Scan(&paaResults).Error
	})
	if err != nil {
		fmt.Println(err)
	}
	res["steps"] = paaResults
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Raw("SELECT distinct ifrs17_group FROM paa_results where product_code = ?  and csm_run_id = ?", prodCode, runId).Scan(&groups).Error }); err != nil {
		fmt.Println(err)
	}
	res["groups"] = groups
	res["run_settings"] = csmRun

	return res
}

func GetPaaStepResultsForAllProductsByRunId(runId int, runDate string) map[string]interface{} {
	var csmRun models.CsmRun
	var paaResults []models.PAAResult
	var prodList []models.ProductList
	var res = make(map[string]interface{})
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Where("paa_run_id = ? and measurement_type = ? and run_date = ?", runId, "PAA", runDate).Find(&csmRun).Error }); err != nil {
		fmt.Println(err)
	}
	var paaRun models.GMMRunSetting
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Where("mgmm_run_id = ?", runId).Find(&paaRun).Error }); err != nil {
		fmt.Println(err)
	}
	csmRun.PaaRunName = paaRun.Name
	query := fmt.Sprintf("SELECT product_code, sum(premium_receipt) as premium_receipt,sum(nb_premium_receipt) as nb_premium_receipt,sum(total_premium_receipt) as total_premium_receipt ,sum(acquisition_expenses) as acquisition_expenses ,sum(insurance_revenue) as insurance_revenue,sum(incurred_claims) as incurred_claims,sum(claims_paid) as claims_paid, sum(incurred_expenses) as incurred_expenses ,sum(earned_premium) as earned_premium, sum(paa_liability_remaining_coverage) as paa_liability_remaining_coverage,sum(deferred_acquisition_expenses) as deferred_acquisition_expenses,sum(amortised_acquisition_cost) as amortised_acquisition_cost, sum(gmm_bel) as gmm_bel, sum(gmm_risk_adjustment) as gmm_risk_adjustment, sum(gmm_reserve) as gmm_reserve, sum(paa_loss_component) as paa_loss_component,sum(reinsurance_premium) as reinsurance_premium, sum(reinsurance_recovery) as reinsurance_recovery,sum(paa_reinsurance_lrc) as paa_reinsurance_lrc, sum(paa_reinsurance_dac) as paa_reinsurance_dac,sum(paa_initial_loss_recovery_component) as paa_initial_loss_recovery_component ,sum(paa_treaty1_initial_loss_recovery_component) as paa_treaty1_initial_loss_recovery_component,sum(paa_treaty2_initial_loss_recovery_component) as paa_treaty2_initial_loss_recovery_component,sum(paa_treaty3_initial_loss_recovery_component) as paa_treaty3_initial_loss_recovery_component,sum(paa_loss_recovery_component) as paa_loss_recovery_component,sum(paa_loss_recovery_unwind) as paa_loss_recovery_unwind,sum(paa_treaty1_loss_recovery_unwind) as paa_treaty1_loss_recovery_unwind,sum(paa_treaty2_loss_recovery_unwind) as paa_treaty2_loss_recovery_unwind,sum(paa_treaty3_loss_recovery_unwind) as paa_treaty3_loss_recovery_unwind ,sum(allocated_reinsurance_premium) as allocated_reinsurance_premium,sum(allocated_treaty1_premium) as allocated_treaty1_premium,sum(allocated_treaty2_premium) as allocated_treaty2_premium,sum(allocated_treaty3_premium) as allocated_treaty3_premium,sum(allocated_reinsurance_flat_commission) as allocated_reinsurance_flat_commission,sum(reinsurance_recovery) as reinsurance_recovery,sum(reinsurance_reinstatement_premium) as reinsurance_reinstatement_premium, sum(reinsurance_ultimate_loss_ratio) as reinsurance_ultimate_loss_ratio, sum(reinsurance_provisional_commission) as reinsurance_provisional_commission,sum(reinsurance_ultimate_commission) as reinsurance_ultimate_commission,sum(reinsurance_profit_commission) as reinsurance_profit_commission, sum(reinsurance_total_paid_to_cedant) as reinsurance_total_paid_to_cedant, sum(reinsurance_investment_component) as reinsurance_investment_component, sum(reinsurance_revenue) as reinsurance_revenue,sum(reinsurance_service_expense) as reinsurance_service_expense,sum(reinsurance_service_result) as reinsurance_service_result,sum(incurred_expenses) as incurred_expenses,sum(non_attributable_expenses) as non_attributable_expenses FROM paa_results where csm_run_id = %d and run_date = '%s' group by product_code", runId, csmRun.RunDate)
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw(query).Scan(&paaResults).Error
	})
	if err != nil {
		fmt.Println(err)
	}
	res["steps"] = paaResults
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Raw("SELECT distinct product_code FROM paa_results where csm_run_id = ? and run_date =?", runId, csmRun.RunDate).Scan(&prodList).Error }); err != nil {
		fmt.Println(err)
	}
	res["products"] = prodList
	res["run_settings"] = csmRun

	//Get PAA Eligibility test results
	var paaEligibilityResults []models.PAAEligibilityTestResult
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Where("run_date = ? and csm_run_id = ?", csmRun.RunDate, runId).Find(&paaEligibilityResults).Error }); err != nil {
		fmt.Println(err)
	}
	//DB.Where("csm_run_id = ?", runId).Find(&csmRun)
	//query = fmt.Sprintf("SELECT projection_month,product_code, ifrs17_group,sum(unearned_premium_reserve) as unearned_premium_reserve,sum(earned_premium) as earned_premium,sum(deferred_acquisition_cost) as deferred_acquisition_cost ,sum(amortised_acquisition_cost) as amortised_acquisition_cost, sum(paa_loss_component) as paa_loss_component,sum(gmm_bel) as gmm_bel,sum(gmm_risk_adjustment) as gmm_risk_adjustment,sum(gmm_reserve) as gmm_reserve, sum(gmm_csm) as gmm_csm,sum(gmm_csm_release) as gmm_csm_release,sum(gmm_dac_release) as gmm_dac_release, sum(gmm_expected_outflows) as gmm_expected_outflows, sum(gmm_risk_adjustment_release) as gmm_risk_adjustment_release, sum(gmm_revenue) as gmm_revenue, sum(revenue_variance) as revenue_variance, sum(revenue_variance_proportion) as revenue_variance_proportion,sum(lrc_variance) as lrc_variance,sum(lrc_variance_proportion) as lrc_variance_proportion FROM paa_eligibility_test_results where csm_run_id = %d and run_date = '%s' group by product_code", runId, csmRun.RunDate)
	//err = DB.Raw(query).Scan(&paaEligibilityResults).Error
	//if err != nil {
	//	fmt.Println(err)
	//}
	res["eligibility_results"] = paaEligibilityResults
	return res
}

func GetPAAResultsForAllProducts(runDate string) map[string]interface{} {
	var csmRun models.CsmRun
	var paaResults []models.PAAResult
	var prodList []models.ProductList
	var res = make(map[string]interface{})
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Where("run_date = ? and measurement_type = ?", runDate, "PAA").Find(&csmRun).Error }); err != nil {
		fmt.Println(err)
	}
	query := fmt.Sprintf("SELECT product_code,run_date, sum(premium_receipt) as premium_receipt,sum(nb_premium_receipt) as nb_premium_receipt,sum(total_premium_receipt) as total_premium_receipt ,sum(acquisition_expenses) as acquisition_expenses ,sum(insurance_revenue) as insurance_revenue,sum(incurred_claims) as incurred_claims,sum(claims_paid) as claims_paid, sum(incurred_expenses) as incurred_expenses ,sum(earned_premium) as earned_premium, sum(paa_liability_remaining_coverage) as paa_liability_remaining_coverage,sum(deferred_acquisition_expenses) as deferred_acquisition_expenses,sum(amortised_acquisition_cost) as amortised_acquisition_cost, sum(gmm_bel) as gmm_bel, sum(gmm_risk_adjustment) as gmm_risk_adjustment, sum(gmm_reserve) as gmm_reserve, sum(paa_loss_component) as paa_loss_component,sum(reinsurance_premium) as reinsurance_premium, sum(reinsurance_recovery) as reinsurance_recovery,sum(paa_reinsurance_lrc) as paa_reinsurance_lrc, sum(paa_reinsurance_dac) as paa_reinsurance_dac,sum(paa_initial_loss_recovery_component) as paa_initial_loss_recovery_component ,sum(paa_treaty1_initial_loss_recovery_component) as paa_treaty1_initial_loss_recovery_component,sum(paa_treaty2_initial_loss_recovery_component) as paa_treaty2_initial_loss_recovery_component,sum(paa_treaty3_initial_loss_recovery_component) as paa_treaty3_initial_loss_recovery_component,sum(paa_loss_recovery_component) as paa_loss_recovery_component,sum(paa_loss_recovery_unwind) as paa_loss_recovery_unwind,sum(paa_treaty1_loss_recovery_unwind) as paa_treaty1_loss_recovery_unwind,sum(paa_treaty2_loss_recovery_unwind) as paa_treaty2_loss_recovery_unwind,sum(paa_treaty3_loss_recovery_unwind) as paa_treaty3_loss_recovery_unwind ,sum(allocated_reinsurance_premium) as allocated_reinsurance_premium,sum(allocated_treaty1_premium) as allocated_treaty1_premium,sum(allocated_treaty2_premium) as allocated_treaty2_premium,sum(allocated_treaty3_premium) as allocated_treaty3_premium,sum(allocated_reinsurance_flat_commission) as allocated_reinsurance_flat_commission,sum(reinsurance_recovery) as reinsurance_recovery,sum(reinsurance_reinstatement_premium) as reinsurance_reinstatement_premium, sum(reinsurance_ultimate_loss_ratio) as reinsurance_ultimate_loss_ratio, sum(reinsurance_provisional_commission) as reinsurance_provisional_commission,sum(reinsurance_ultimate_commission) as reinsurance_ultimate_commission,sum(reinsurance_profit_commission) as reinsurance_profit_commission, sum(reinsurance_total_paid_to_cedant) as reinsurance_total_paid_to_cedant, sum(reinsurance_investment_component) as reinsurance_investment_component, sum(reinsurance_revenue) as reinsurance_revenue,sum(reinsurance_service_expense) as reinsurance_service_expense,sum(reinsurance_service_result) as reinsurance_service_result,sum(incurred_expenses) as incurred_expenses,sum(non_attributable_expenses) as non_attributable_expenses FROM paa_results where run_date = '%s'  group by product_code", runDate)
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw(query).Scan(&paaResults).Error
	})
	if err != nil {
		fmt.Println(err)
	}
	res["steps"] = paaResults
	res["run"] = csmRun
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Raw("SELECT distinct product_code FROM paa_results where run_date = ?", runDate).Scan(&prodList).Error }); err != nil {
		fmt.Println(err)
	}
	res["products"] = prodList

	return res
}

func GetPAAResultsForAllProductsForDownstreamCalcs(runDate string) []models.PAAResult {
	var csmRun models.CsmRun
	var paaResults []models.PAAResult
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Where("run_date = ? and measurement_type = ?", runDate, "PAA").Find(&csmRun).Error }); err != nil {
		fmt.Println(err)
	}
	query := fmt.Sprintf("SELECT sum(premium_receipt) as premium_receipt,sum(nb_premium_receipt) as nb_premium_receipt,sum(total_premium_receipt) as total_premium_receipt, sum(earned_premium) as earned_premium,sum(insurance_revenue) as insurance_revenue,sum(incurred_claims) as incurred_claims,sum(claims_paid) as claims_paid, sum(paa_liability_remaining_coverage) as paa_liability_remaining_coverage,sum(acquisition_expenses) as acquisition_expenses,sum(amortised_acquisition_cost) as amortised_acquisition_cost, sum(gmm_bel) as gmm_bel, sum(gmm_risk_adjustment) as gmm_risk_adjustment, sum(gmm_reserve) as gmm_reserve, sum(paa_loss_component) as paa_loss_component,sum(reinsurance_premium_paid) as reinsurance_premium_paid,sum(allocated_reinsurance_premium) as allocated_reinsurance_premium,sum(allocated_reinsurance_flat_commission) as allocated_reinsurance_flat_commission,sum(reinsurance_recovery) as reinsurance_recovery,sum(reinsurance_reinstatement_premium) as reinsurance_reinstatement_premium, sum(reinsurance_ultimate_loss_ratio) as reinsurance_ultimate_loss_ratio, sum(reinsurance_provisional_commission) as reinsurance_provisional_commission,sum(reinsurance_ultimate_commission) as reinsurance_ultimate_commission,sum(reinsurance_profit_commission) as reinsurance_profit_commission, sum(reinsurance_total_paid_to_cedant) as reinsurance_total_paid_to_cedant, sum(reinsurance_investment_component) as reinsurance_investment_component, sum(reinsurance_revenue) as reinsurance_revenue,sum(reinsurance_service_expense) as reinsurance_service_expense,sum(reinsurance_service_result) as reinsurance_service_result,sum(incurred_expenses) as incurred_expenses,sum(non_attributable_expenses) as non_attributable_expenses FROM paa_results where run_date = '%s'", runDate) //group by product_code
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw(query).Scan(&paaResults).Error
	})
	if err != nil {
		fmt.Println(err)
	}
	return paaResults
}

func GetPAAResultsForOneProductForDownstreamCalcs(runDate, prodCode string) []models.PAAResult {
	var csmRun models.CsmRun
	var paaResults []models.PAAResult
	if err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error { return d.Where("run_date = ? and measurement_type = ?", runDate, "PAA").Find(&csmRun).Error }); err != nil {
		fmt.Println(err)
	}
	query := fmt.Sprintf("SELECT product_code, sum(premium_receipt) as premium_receipt,sum(nb_premium_receipt) as nb_premium_receipt,sum(total_premium_receipt) as total_premium_receipt, sum(earned_premium) as earned_premium,sum(insurance_revenue) as insurance_revenue,sum(incurred_claims) as incurred_claims,sum(claims_paid) as claims_paid ,sum(paa_liability_remaining_coverage) as paa_liability_remaining_coverage,sum(acquisition_expenses) as acquisition_expenses,sum(amortised_acquisition_cost) as amortised_acquisition_cost, sum(gmm_bel) as gmm_bel, sum(gmm_risk_adjustment) as gmm_risk_adjustment, sum(gmm_reserve) as gmm_reserve, sum(paa_loss_component) as paa_loss_component,sum(reinsurance_premium_paid) as reinsurance_premium_paid, sum(allocated_reinsurance_premium) as allocated_reinsurance_premium,sum(allocated_reinsurance_flat_commission) as allocated_reinsurance_flat_commission,sum(reinsurance_recovery) as reinsurance_recovery,sum(reinsurance_reinstatement_premium) as reinsurance_reinstatement_premium, sum(reinsurance_ultimate_loss_ratio) as reinsurance_ultimate_loss_ratio, sum(reinsurance_provisional_commission) as reinsurance_provisional_commission,sum(reinsurance_ultimate_commission) as reinsurance_ultimate_commission,sum(reinsurance_profit_commission) as reinsurance_profit_commission, sum(reinsurance_total_paid_to_cedant) as reinsurance_total_paid_to_cedant, sum(reinsurance_investment_component) as reinsurance_investment_component, sum(reinsurance_revenue) as reinsurance_revenue,sum(reinsurance_service_expense) as reinsurance_service_expense,sum(reinsurance_service_result) as reinsurance_service_result,sum(incurred_expenses) as incurred_expenses,sum(non_attributable_expenses) as non_attributable_expenses FROM paa_results where run_date = '%s' and product_code = '%s'  group by product_code", runDate, prodCode)
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw(query).Scan(&paaResults).Error
	})
	if err != nil {
		fmt.Println(err)
	}
	return paaResults
}

func GetAosStepResultsForGroup(ifrs17Group, runDate string) []models.AOSStepResult {
	var agg []models.AOSStepResult
	query := fmt.Sprintf("SELECT product_code, ifrs17_group,  name, sum(bel) as bel, sum(bel_at12) as bel_at12, sum(risk_adjustment) as risk_adjustment, sum(risk_adjustment_at12) as risk_adjustment_at12, sum(best_estimate_liability_change) as best_estimate_liability_change, sum(risk_adjustment_change) as risk_adjustment_change, SUM(csm_change) as csm_change, SUM(liability_change) as liability_change, Sum(loss_component_change) as loss_component_change, SUM(pnl_change) as pnl_change, sum(csm_buildup) as csm_buildup,sum(risk_adjustment_buildup) as risk_adjustment_buildup, sum(loss_component_buildup) as loss_component_buildup, sum(sar_actual_claim_net_lc) as sar_actual_claim_net_lc, sum(expected_ra_net_of_lc) as expected_ra_net_of_lc, sum(lc_sar) as lc_sar, sum(csm_release_ratio) as csm_release_ratio,sum(csm_release) as csm_release, sum(discounted_coverage_units) as discounted_coverage_units, sum(sum_coverage_units) as sum_coverage_units, sum(interest_accretion_fac) as interest_accretion_fac, sum(bel_inflow) as bel_inflow, sum(bel_outflow) as bel_outflow, sum(bel_inflow_at12) as  bel_inflow_at12, sum(bel_outflow_at12) as bel_outflow_at12, sum(bel_inflow_change) as bel_inflow_change,sum(bel_outflow_change) as bel_outflow_change, sum(expected_cash_inflow) as expected_cash_inflow, sum(expected_cash_outflow) as expected_cash_outflow  FROM aos_step_results where run_date = '%s' and ifrs17_group = '%s'  group by product_code, ifrs17_group, name", runDate, ifrs17Group)
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw(query).Scan(&agg).Error
	})
	if err != nil {
		fmt.Println(err)
	}

	return agg
}

func GetPAAResultsForGroup(ifrs17Group, runDate string) []models.PAAResult {
	var agg []models.PAAResult
	query := fmt.Sprintf("SELECT product_code,run_date, ifrs17_group, sum(premium_receipt) as premium_receipt,sum(nb_premium_receipt) as nb_premium_receipt,sum(total_premium_receipt) as total_premium_receipt,sum(earned_premium) as earned_premium,sum(insurance_revenue) as insurance_revenue,sum(incurred_claims) as incurred_claims,sum(claims_paid) as claims_paid, sum(paa_liability_remaining_coverage) as paa_liability_remaining_coverage,sum(acquisition_expenses) as acquisition_expenses,sum(amortised_acquisition_cost) as amortised_acquisition_cost, sum(gmm_bel) as gmm_bel, sum(gmm_risk_adjustment) as gmm_risk_adjustment, sum(gmm_reserve) as gmm_reserve, sum(paa_loss_component) as paa_loss_component, sum(allocated_reinsurance_premium) as allocated_reinsurance_premium,sum(allocated_reinsurance_flat_commission) as allocated_reinsurance_flat_commission,sum(reinsurance_recovery) as reinsurance_recovery,sum(reinsurance_reinstatement_premium) as reinsurance_reinstatement_premium, sum(reinsurance_ultimate_loss_ratio) as reinsurance_ultimate_loss_ratio, sum(reinsurance_provisional_commission) as reinsurance_provisional_commission,sum(reinsurance_ultimate_commission) as reinsurance_ultimate_commission,sum(reinsurance_profit_commission) as reinsurance_profit_commission, sum(reinsurance_total_paid_to_cedant) as reinsurance_total_paid_to_cedant, sum(reinsurance_investment_component) as reinsurance_investment_component, sum(reinsurance_revenue) as reinsurance_revenue,sum(reinsurance_service_expense) as reinsurance_service_expense,sum(reinsurance_service_result) as reinsurance_service_result FROM paa_results where ifrs17_group = '%s' and run_date = '%s'  group by product_code", ifrs17Group, runDate)
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw(query).Scan(&agg).Error
	})
	if err != nil {
		fmt.Println(err)
	}

	return agg
}

func GetLicResultsForAllProductsForDownstreamCalcs(runDate string) []models.LicBuildupResult {

	var licBuildupResults []models.LicBuildupResult
	query := fmt.Sprintf("SELECT name, sum(ibnr) as ibnr, sum(ibnr_risk_adjustment) as ibnr_risk_adjustment, sum(reported_claims) as reported_claims, sum(cash_back) as cash_back , sum(unadjusted_loss_adjustment_expenses) as unadjusted_loss_adjustment_expenses FROM lic_buildup_results where run_date = '%s' group by name", runDate) //group by product_code
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw(query).Scan(&licBuildupResults).Error
	})
	if err != nil {
		fmt.Println(err)
	}
	return licBuildupResults
}
func GetLicResultsForOneProductForDownstreamCalcs(runDate, prodCode string) []models.LicBuildupResult {

	var licBuildupResults []models.LicBuildupResult
	query := fmt.Sprintf("SELECT sum(ibnr) as ibnr, sum(ibnr_risk_adjustment) as ibnr_risk_adjustment, sum(reported_claims) as reported_claims, sum(cash_back) as cash_back , sum(unadjusted_loss_adjustment_expenses) as unadjusted_loss_adjustment_expenses FROM lic_buildup_results where run_date = '%s' and product_code = '%s' group by name", runDate, prodCode) //group by product_code
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw(query).Scan(&licBuildupResults).Error
	})
	if err != nil {
		fmt.Println(err)
	}
	return licBuildupResults
}

func GetLicResultsForOneGroupForDownstreamCalcs(runDate, ifrs17Group string) []models.LicBuildupResult {

	var licBuildupResults []models.LicBuildupResult
	query := fmt.Sprintf("SELECT sum(ibnr) as ibnr, sum(ibnr_risk_adjustment) as ibnr_risk_adjustment, sum(reported_claims) as reported_claims, sum(cash_back) as cash_back , sum(unadjusted_loss_adjustment_expenses) as unadjusted_loss_adjustment_expenses FROM lic_buildup_results where run_date = '%s' and ifrs_group = '%s' group by name", runDate, ifrs17Group) //group by product_code
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw(query).Scan(&licBuildupResults).Error
	})
	if err != nil {
		fmt.Println(err)
	}
	return licBuildupResults
}

func getRiskDriversFor(productCode string) models.RiskDriver {
	var rd models.RiskDriver

	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("product_code=?", productCode).First(&rd).Error
	})
	fmt.Println(err)
	return rd
}

func getRiskAdjustmentFactors(productCode string, riskAdjustmentYear int, version string) models.RiskAdjustmentFactor {
	var rf models.RiskAdjustmentFactor

	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("product_code = ? and year = ? and version = ?", productCode, riskAdjustmentYear, version).First(&rf).Error
	})
	fmt.Println(err)
	return rf
}

func getFieldValue(e *models.ScopedAggregatedProjection, field string) float64 {
	var field2, newField2 string
	if field != "" {
		if field == "NOT APPLICABLE" {
			return 0
		}
		if field == "DISCOUNTED_OUTGO" {
			field = "DISCOUNTED_CASH_OUTFLOW_EXCL_ACQUISITION"
		}

		newField := convertToCamel(field)

		r := reflect.ValueOf(e)
		f := reflect.Indirect(r).FieldByName(newField)
		fmt.Println(reflect.ValueOf(f))
		fmt.Println(f.Interface())
		fmt.Println(f.Float())
		if field == DISCOUNTED_DEATH_OUTGO {
			field2 = "DISCOUNTED_ACCIDENTAL_DEATH_OUTGO"
			newField2 = convertToCamel(field2)
			f2 := reflect.Indirect(r).FieldByName(newField2)
			fmt.Println(reflect.ValueOf(f2))
			fmt.Println(f2.Interface())
			fmt.Println(f2.Float())
			return f.Float() + f2.Float()
		}

		return f.Float()
	}
	return 0
}

func convertToCamel(s string) string {
	parts := strings.Split(s, "_")
	var builder strings.Builder
	for i := range parts {
		parts[i] = strings.ToLower(parts[i])
	}

	for _, v := range parts {
		builder.WriteString(v + " ")
	}

	return strings.Replace(strings.Title(builder.String()), " ", "", -1)
}

func reducestepresultsDecimals(aos []models.AOSStepResult) []models.AOSStepResult {
	for i, _ := range aos {
		aos[i].BEL = math.Round(aos[i].BEL)
		aos[i].BELAt12 = math.Round(aos[i].BELAt12)
		aos[i].RiskAdjustment = math.Round(aos[i].RiskAdjustment)
		aos[i].RiskAdjustmentAt12 = math.Round(aos[i].RiskAdjustmentAt12)
		aos[i].BestEstimateLiabilityChange = math.Round(aos[i].BestEstimateLiabilityChange)
		aos[i].RiskAdjustmentChange = math.Round(aos[i].RiskAdjustmentChange)
		aos[i].CSMChange = math.Round(aos[i].CSMChange)
		aos[i].LiabilityChange = math.Round(aos[i].LiabilityChange)
		aos[i].LossComponentChange = math.Round(aos[i].LossComponentChange)
		aos[i].PNLChange = math.Round(aos[i].PNLChange)
		aos[i].CSMBuildup = math.Round(aos[i].CSMBuildup)
		aos[i].RiskAdjustmentBuildup = math.Round(aos[i].RiskAdjustmentBuildup)
		aos[i].LossComponentBuildup = math.Round(aos[i].LossComponentBuildup)
		aos[i].SarActualClaimNetLc = math.Round(aos[i].SarActualClaimNetLc)
		aos[i].ExpectedRaNetOfLc = math.Round(aos[i].ExpectedRaNetOfLc)
		aos[i].SumCoverageUnits = math.Round(aos[i].SumCoverageUnits)
		aos[i].DiscountedCoverageUnits = math.Round(aos[i].DiscountedCoverageUnits)

	}
	return aos
}

// CalculatePAACsm will house the workflow for IFRS17 PAA calculations
func CalculatePAACsm(csmRun *models.CsmRun, group models.GroupResults) []models.PAAResult {
	var jt models.JournalTransactions
	var results []models.PAAResult
	var paabuildups []models.PAABuildUp
	var paabuildup models.PAABuildUp
	var paaEligibilityTestResults []models.PAAEligibilityTestResult
	var err error
	var result models.PAAResult

	var treaty1paaPropReinsuranceResult struct {
		ReinsurancePremiumPaid           float64
		AllocatedReinsurancePremium      float64
		InitialLossRecoveryComponent     float64
		SubsequentLossRecoveryAdjustment float64
		LossRecoveryUnwind               float64

		PaaRecovery float64
	}
	var treaty2paaPropReinsuranceResult struct {
		ReinsurancePremiumPaid           float64
		AllocatedReinsurancePremium      float64
		InitialLossRecoveryComponent     float64
		SubsequentLossRecoveryAdjustment float64
		LossRecoveryUnwind               float64

		PaaRecovery float64
	}
	var treaty3paaPropReinsuranceResult struct {
		ReinsurancePremiumPaid           float64
		AllocatedReinsurancePremium      float64
		InitialLossRecoveryComponent     float64
		SubsequentLossRecoveryAdjustment float64
		LossRecoveryUnwind               float64

		PaaRecovery float64
	}
	//var sumFutureRiskUnits float64
	//var buildups []models.PAABuildUp
	result.CsmRunID = csmRun.PaaRunId
	result.RunDate = csmRun.RunDate
	result.RunId = csmRun.RunID

	//runMonth, _ := strconv.Atoi(csmRun.RunDate[5:])

	result.ProductCode = group.ProductCode

	result.IFRS17Group = group.IFRS17Group
	//result.ID = 0
	var sap []models.ModifiedGMMScopedAggregation
	err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("ifrs17_group=? and run_id=?", group.IFRS17Group, csmRun.PaaRunId).Order("projection_month asc").Find(&sap).Error
	})
	if err != nil {
		log.Error(err)
	}

	var modifiedRunSetting models.GMMRunSetting
	err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("mgmm_run_id=?", csmRun.PaaRunId).Find(&modifiedRunSetting).Error
	})
	if err != nil {
		log.Error(err)
	}
	result.PortfolioName = modifiedRunSetting.PortfolioName

	var finance []models.PAAFinance //check if we need to add a condition
	err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("ifrs17_group=? and year=? and portfolio_name = ? and version=?", group.IFRS17Group, csmRun.FinanceYear, modifiedRunSetting.PortfolioName, csmRun.FinanceVersion).Find(&finance).Error
	})
	if err != nil {
		log.Error(err)
	}

	if len(finance) == 0 {
		csmRun.ProcessingStatus = "failed"
		csmRun.FailureReason = "Missing finance data file for " + group.IFRS17Group
		DB.Save(csmRun)
		return nil
	}

	var riparam []models.ReinsuranceParameter
	err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("product_code=? and year=?", group.ProductCode, modifiedRunSetting.ParameterYear).Find(&riparam).Error
	})
	if err != nil {
		log.Error(err)
	}

	//var mp []models.ModifiedGMMModelPoint
	//err = DB.Where("ifrs17_group=? and year=? and paa_portfolio_name = ?", group.IFRS17Group, modifiedRunSetting.ModelPoint, modifiedRunSetting.ExpConfigurationName).Find(&mp).Error
	//if err != nil {
	//	log.Error(err)
	//}

	//var paramAcquisition []models.ModifiedGMMParameter
	//err = DB.Where("product_code=? and year=? and sub_product_code=?", group.ProductCode, modifiedRunSetting.AcquisitionExpenses, mp[0].SubProductCode).Find(&paramAcquisition).Error
	//if err != nil {
	//	log.Error(err)
	//}

	var portfolio models.PaaPortfolio
	err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("name=?", modifiedRunSetting.PortfolioName).Find(&portfolio).Error
	})
	if err != nil {
		log.Error(err)
	}

	if portfolio.InsuranceType == ProportionalReinsurance {
		err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Model(&models.PAAResult{}).
				Select("SUM(allocated_treaty1_premium) AS allocated_reinsurance_premium, SUM(paa_treaty1_initial_loss_recovery_component) AS initial_loss_recovery_component, SUM(paa_treaty1_loss_recovery_unwind) AS loss_recovery_unwind, SUM(treaty1_recovery) AS paa_recovery, SUM(treaty1_premium_paid) as reinsurance_premium_paid").
				Where("treaty1_ifrs17_group = ? AND run_date = ?", group.IFRS17Group, csmRun.RunDate).
				Scan(&treaty1paaPropReinsuranceResult).Error
		})

		err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Model(&models.PAAResult{}).
				Select("SUM(allocated_treaty2_premium) AS allocated_reinsurance_premium, SUM(paa_treaty2_initial_loss_recovery_component) AS initial_loss_recovery_component, SUM(paa_treaty2_loss_recovery_unwind) AS loss_recovery_unwind, SUM(treaty2_recovery) AS paa_recovery, SUM(treaty2_premium_paid) as reinsurance_premium_paid").
				Where("treaty2_ifrs17_group = ? AND run_date = ?", group.IFRS17Group, csmRun.RunDate).
				Scan(&treaty2paaPropReinsuranceResult).Error
		})

		err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Model(&models.PAAResult{}).
				Select("SUM(allocated_treaty3_premium) AS allocated_reinsurance_premium, SUM(paa_treaty3_initial_loss_recovery_component) AS initial_loss_recovery_component, SUM(paa_treaty3_loss_recovery_unwind) AS loss_recovery_unwind, SUM(treaty3_recovery) AS paa_recovery, SUM(treaty3_premium_paid) as reinsurance_premium_paid").
				Where("treaty3_ifrs17_group = ? AND run_date = ?", group.IFRS17Group, csmRun.RunDate).
				Scan(&treaty3paaPropReinsuranceResult).Error
		})
	}
	//var userPattern []models.PremiumEarningPattern
	//err = DB.Where("product_code=? and year=?", group.ProductCode, csmRun.RunDate[:4]).Find(&userPattern).Error
	//if err != nil {
	//	log.Error(err)
	//}

	//prevYearStr := csmRun.RunDate[:4]
	//prevYearMth := csmRun.RunDate[5:]
	//prevYear, _ := strconv.Atoi(prevYearStr)
	//prevYear = prevYear - 1
	//yearSearchArg := strconv.Itoa(prevYear) + "-" + prevYearMth
	//fmt.Println(yearSearchArg)
	yearSearchArg := csmRun.OpeningBalanceDate
	var prevBsr models.BalanceSheetRecord
	err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("ifrs17_group = ? and date = ? and measurement_type = ? and product_code = ?", result.IFRS17Group, yearSearchArg, "PAA", group.ProductCode).Find(&prevBsr).Error
	})
	if err != nil {
		fmt.Println("query bsr error:", err)
	}

	//if len(mp) > 0 && portfolio.PremiumEarningPattern != PassageofTime {
	//	sumFutureRiskUnits = getsumFutureRiskUnit(modifiedRunSetting.PremiumEarning, group.ProductCode, mp[0].DurationInForceMonths)
	//}

	var tempDacRelease float64
	var paaLossComponentInitialRecognition float64
	//var paaLossRecoveryInitialRecognition float64
	var paaTreaty1LossRecoveryInitialRecognition float64
	var paaTreaty2LossRecoveryInitialRecognition float64
	var paaTreaty3LossRecoveryInitialRecognition float64
	//var paaLossComponentAdjustmentOut float64
	var paaLossRecoveryAdjustmentOut float64
	var paaTreaty1LossRecoveryAdjustmentOut float64
	var paaTreaty2LossRecoveryAdjustmentOut float64
	var paaTreaty3LossRecoveryAdjustmentOut float64
	var step3LossComponentUnwind float64
	var step3LossRecoveryUnwind float64
	var step4LossComponentUnwind float64
	//var step4LossRecoveryUnwind float64
	var step8LossComponentUnwind float64
	var step8LossRecoveryUnwind float64

	if len(finance) > 0 {
		result.PremiumReceipt = 0
		result.EarnedPremium = 0
		result.NBPremiumReceipt = 0
		result.AcquisitionExpenses = 0
		paabuildup.PaaLrcBuildup = 0
		paabuildup.DacBuildup = 0

		if finance[0].IFStatus == NB {
			result.NBPremiumReceipt = finance[0].ActualPremiumReceipt + finance[0].PremiumDebtors - finance[0].PremiumRefund
			result.PremiumReceipt = 0
			//if len(paramAcquisition) > 0 {
			result.AcquisitionExpenses = finance[0].AcquisitionCostPaid //+ paramAcquisition[0].InitialExpenseAmount + paramAcquisition[0].InitialExpenseProportion*result.NBPremiumReceipt
			//} else {
			//	result.AcquisitionExpenses = finance[0].CommissionPaid
			//}
		}
		if finance[0].IFStatus == IF {
			result.NBPremiumReceipt = 0
			result.PremiumReceipt = finance[0].ActualPremiumReceipt + finance[0].PremiumDebtors - finance[0].PremiumRefund
			result.AcquisitionExpenses = finance[0].AcquisitionCostPaid
		}
		result.TotalPremiumReceipt = result.NBPremiumReceipt + result.PremiumReceipt

		if portfolio.InsuranceType == Direct {
			//result.ReinsurancePremiumPaid = finance[0].ReinsurancePremium
			//result.AllocatedReinsurancePremium = sap[0].Treaty1CurrentPeriodEarnedPremium + sap[0].Treaty2CurrentPeriodEarnedPremium + sap[0].Treaty3CurrentPeriodEarnedPremium
			result.AllocatedTreaty1Premium = sap[0].Treaty1CurrentPeriodEarnedPremium
			result.AllocatedTreaty2Premium = sap[0].Treaty2CurrentPeriodEarnedPremium
			result.AllocatedTreaty3Premium = sap[0].Treaty3CurrentPeriodEarnedPremium
			//result.ReinsuranceRecovery = finance[0].Treaty1Recovery + finance[0].Treaty2Recovery + finance[0].Treaty3Recovery
			//result.Treaty1Recovery = finance[0].Treaty1Recovery
			//result.Treaty2Recovery = finance[0].Treaty2Recovery
			//result.Treaty3Recovery = finance[0].Treaty3Recovery
			//result.ReinsuranceRevenue = result.AllocatedReinsurancePremium
			//result.ReinsuranceServiceExpense = result.ReinsuranceRecovery
			result.Treaty1IFRS17Group = sap[0].Treaty1IFRS17Group
			result.Treaty2IFRS17Group = sap[0].Treaty2IFRS17Group
			result.Treaty3IFRS17Group = sap[0].Treaty3IFRS17Group
		}

		if portfolio.InsuranceType == ProportionalReinsurance {
			result.ReinsurancePremium = finance[0].ReinsurancePremium
			result.ReinsuranceRecovery = finance[0].ReinsuranceRecovery //treaty1paaPropReinsuranceResult.PaaRecovery + treaty2paaPropReinsuranceResult.PaaRecovery + treaty3paaPropReinsuranceResult.PaaRecovery
			result.AllocatedReinsurancePremium = treaty1paaPropReinsuranceResult.AllocatedReinsurancePremium + treaty2paaPropReinsuranceResult.AllocatedReinsurancePremium + treaty3paaPropReinsuranceResult.AllocatedReinsurancePremium
			result.PaaInitialLossRecoveryComponent = treaty1paaPropReinsuranceResult.InitialLossRecoveryComponent + treaty2paaPropReinsuranceResult.InitialLossRecoveryComponent + treaty3paaPropReinsuranceResult.InitialLossRecoveryComponent
			result.PaaLossRecoveryUnwind = treaty1paaPropReinsuranceResult.LossRecoveryUnwind + treaty2paaPropReinsuranceResult.LossRecoveryUnwind + treaty3paaPropReinsuranceResult.LossRecoveryUnwind
			result.PaaReinsuranceLrc = result.ReinsurancePremium + prevBsr.PaaReinsuranceLRC - (treaty1paaPropReinsuranceResult.AllocatedReinsurancePremium + treaty2paaPropReinsuranceResult.AllocatedReinsurancePremium + treaty3paaPropReinsuranceResult.AllocatedReinsurancePremium)

		}

		if portfolio.InsuranceType != Direct && portfolio.InsuranceType != ProportionalReinsurance {
			ulrlowerboundrate, _ := getPAAReinsuranceParameter(group.ProductCode, UlrLowerboundRate)
			ulrupperboundrate, _ := getPAAReinsuranceParameter(group.ProductCode, UlrUpperboundRate)
			slidingscaleminrate, _ := getPAAReinsuranceParameter(group.ProductCode, SlidingScaleMinRate)
			slidingscalemaxrate, _ := getPAAReinsuranceParameter(group.ProductCode, SlidingScaleMaxRate)
			profitcommissionrate, _ := getPAAReinsuranceParameter(group.ProductCode, ProfitCommissionRate)
			_, reinsuranceinwardoutward := getPAAReinsuranceParameter(group.ProductCode, ReinsuranceInwardOutward)

			result.ReinsurancePremium = result.TotalPremiumReceipt
			result.AllocatedReinsurancePremium = sap[0].CurrentPeriodEarnedPremium               //result.NBPremiumReceipt + result.PremiumReceipt
			result.AllocatedReinsuranceFlatCommission = sap[0].CurrentPeriodAmortisedAcquisition //finance[0].CommissionPaid
			result.ReinsuranceRecovery = finance[0].ReinsuranceRecovery
			result.ReinsuranceReinstatementPremium = finance[0].ReinstatementPremium
			TotalPrem := result.AllocatedReinsurancePremium + result.AllocatedReinsuranceFlatCommission //+ result.ReinsuranceReinstatementPremium
			if TotalPrem > 0 {
				result.ReinsuranceUltimateLossRatio = finance[0].ReinsuranceRecovery / TotalPrem
			}
			result.ReinsuranceProvisionalCommission = finance[0].ProvisionalCommission

			if result.ReinsuranceUltimateLossRatio > ulrupperboundrate {
				if result.ReinsuranceUltimateLossRatio < 1 {
					result.ReinsuranceUltimateCommission = TotalPrem * slidingscaleminrate
				}

			}
			if result.ReinsuranceUltimateLossRatio < ulrlowerboundrate {
				result.ReinsuranceUltimateCommission = TotalPrem * slidingscalemaxrate
			}
			if result.ReinsuranceUltimateLossRatio >= ulrlowerboundrate && result.ReinsuranceUltimateLossRatio <= ulrupperboundrate {
				result.ReinsuranceUltimateCommission = result.ReinsuranceProvisionalCommission
			}

			investmentCompultimateCommission := (result.AllocatedReinsurancePremium + result.AllocatedReinsuranceFlatCommission) * slidingscalemaxrate // calculating investment component with zero claims and zero ultimate loss ratio; reinsurance premium is net of flat commission

			reinsuranceProfit := result.AllocatedReinsurancePremium - result.ReinsuranceRecovery

			result.ReinsuranceProfitCommission = math.Max(reinsuranceProfit*profitcommissionrate, 0)

			result.ReinsuranceTotalPaidToCedant = result.ReinsuranceUltimateCommission + result.ReinsuranceRecovery + result.ReinsuranceProfitCommission - result.ReinsuranceReinstatementPremium
			result.ReinsuranceInvestmentComponent = investmentCompultimateCommission + result.AllocatedReinsurancePremium*profitcommissionrate
			result.ReinsuranceRevenue = result.AllocatedReinsurancePremium - result.ReinsuranceInvestmentComponent
			result.ReinsuranceServiceExpense = result.ReinsuranceTotalPaidToCedant - result.ReinsuranceInvestmentComponent
			var signconstant = 1
			if reinsuranceinwardoutward == Outward {
				signconstant = -1
			}
			result.ReinsuranceServiceResult = float64(signconstant) * (result.ReinsuranceRevenue - result.ReinsuranceServiceExpense)
		}

		if portfolio.InsuranceType == ProportionalReinsurance {
			result.ReinsuranceServiceResult = result.ReinsuranceRecovery - result.AllocatedReinsurancePremium
		}

		// PAA Buildup Block//
		if len(sap) > 0 && len(finance) > 0 {
			if finance[0].ActualPremiumReceipt > 0 {
				result.GmmBel = utils.FloatPrecision(sap[0].DiscountedNetCashFlowsLockedin, AccountingPrecision)
				result.GmmRiskAdjustment = utils.FloatPrecision(sap[0].RiskAdjustment, AccountingPrecision)
				result.GmmReserve = utils.FloatPrecision(result.GmmBel+result.GmmRiskAdjustment, AccountingPrecision)

				if sap[0].CurrentPeriodEarnedPremium == 0 {
					result.PaaLossComponent = utils.FloatPrecision(math.Max(result.GmmReserve, 0), AccountingPrecision)
				} else {
					if result.GmmReserve < 0 {
						result.PaaLossComponent = 0
					} else {
						result.PaaLossComponent = utils.FloatPrecision(math.Max(result.GmmReserve-math.Max(paabuildup.PaaLrcBuildup, 0), 0), AccountingPrecision)
					}
				}
			}
		} else {
			result.GmmBel = 0
			result.GmmRiskAdjustment = 0
			result.GmmReserve = 0
			result.PaaLossComponent = 0
		}

		//Line 1//
		paabuildup.Name = "B/F(Current)"
		paabuildup.RunDate = result.RunDate
		paabuildup.RunId = csmRun.PaaRunId
		paabuildup.PortfolioName = finance[0].PortfolioName
		paabuildup.ProductCode = result.ProductCode
		paabuildup.IFRS17Group = result.IFRS17Group
		paabuildup.VariableChange = 0
		if portfolio.InsuranceType == Direct {
			paabuildup.PaaLrcBuildup = prevBsr.PaaLiabilityRemainingCoverage
			paabuildup.LossComponentBuildup = prevBsr.PaaLossComponent
			paabuildup.DacBuildup = prevBsr.PaaDAC
			paabuildup.PaaReinsuranceLrcBuildup = prevBsr.PaaReinsuranceLRC
			paabuildup.PaaReinsuranceDacBuildup = prevBsr.PaaReinsuranceDAC
			paabuildup.LossRecoveryBuildup = prevBsr.PaaLossRecoveryComponent
		}
		if portfolio.InsuranceType != Direct {
			paabuildup.PaaReinsuranceLrcBuildup = prevBsr.PaaReinsuranceLRC
			paabuildup.PaaReinsuranceDacBuildup = prevBsr.PaaReinsuranceDAC
			//paabuildup.PaaLrcBuildup = prevBsr.PaaReinsuranceLRC
			//paabuildup.DacBuildup = prevBsr.PaaReinsuranceDAC
		}
		paabuildup.InsuranceRevenue = 0
		paabuildup.CoverageUnits = 0
		paabuildup.AcquisitionCostAmortizationProportion = 0
		paabuildup.EarnedPremiumProportion = 0

		paabuildups = append(paabuildups, paabuildup)

		//Line 2//
		paabuildup.Name = "B/F(Lockedin)"
		paabuildup.RunDate = result.RunDate
		paabuildup.RunId = csmRun.PaaRunId
		paabuildup.PortfolioName = finance[0].PortfolioName
		paabuildup.ProductCode = result.ProductCode
		paabuildup.IFRS17Group = result.IFRS17Group
		paabuildup.VariableChange = 0
		if portfolio.InsuranceType == Direct {
			paabuildup.PaaLrcBuildup = prevBsr.PaaLiabilityRemainingCoverage
			paabuildup.LossComponentBuildup = prevBsr.PaaLossComponent
			paabuildup.DacBuildup = prevBsr.PaaDAC
			paabuildup.PaaReinsuranceLrcBuildup = prevBsr.PaaReinsuranceLRC
			paabuildup.PaaReinsuranceDacBuildup = prevBsr.PaaReinsuranceDAC
			paabuildup.LossRecoveryBuildup = prevBsr.PaaLossRecoveryComponent
		}
		if portfolio.InsuranceType != Direct {
			paabuildup.PaaReinsuranceLrcBuildup = prevBsr.PaaReinsuranceLRC
			paabuildup.PaaReinsuranceDacBuildup = prevBsr.PaaReinsuranceDAC
			//paabuildup.PaaLrcBuildup = prevBsr.PaaReinsuranceLRC
			//paabuildup.DacBuildup = prevBsr.PaaReinsuranceDAC
		}

		paabuildup.InsuranceRevenue = 0
		paabuildup.CoverageUnits = 0
		paabuildup.AcquisitionCostAmortizationProportion = 0
		paabuildup.EarnedPremiumProportion = 0

		paabuildups = append(paabuildups, paabuildup)

		//Line 3//
		paabuildup.Name = "Premium_Receipt"
		paabuildup.RunDate = result.RunDate
		paabuildup.RunId = csmRun.PaaRunId
		paabuildup.PortfolioName = finance[0].PortfolioName
		paabuildup.ProductCode = result.ProductCode
		paabuildup.IFRS17Group = result.IFRS17Group
		paabuildup.VariableChange = result.PremiumReceipt
		paabuildup.PaaLrcBuildup += paabuildup.VariableChange

		if result.PremiumReceipt != 0 {
			paabuildup.ModifiedGMMBel = result.GmmBel
			paabuildup.ModifiedGMMRiskAdjustment = result.GmmRiskAdjustment
			paabuildup.ModifiedGMMReserve = result.GmmReserve
		}

		//inforcelossComponentUnwindDenominator := sap[0].SumFutureEarnedPremium + sap[0].CurrentPeriodEarnedPremium
		//if inforcelossComponentUnwindDenominator > 0 {
		//	paabuildup.LossComponentUnwind = utils.FloatPrecision(paabuildup.LossComponentBuildup*sap[0].CurrentPeriodEarnedPremium/inforcelossComponentUnwindDenominator, AccountingPrecision)
		//	paabuildup.LossRecoveryUnwind = utils.FloatPrecision(paabuildup.LossRecoveryBuildup*sap[0].CurrentPeriodEarnedPremium/inforcelossComponentUnwindDenominator, AccountingPrecision)
		//}

		//if result.GmmReserve < 0 {
		//	result.PaaLossComponent = 0
		//} else {
		//	result.PaaLossComponent = utils.FloatPrecision(math.Max(result.GmmReserve-math.Max(paabuildup.PaaLrcBuildup, 0), 0), AccountingPrecision)
		//}

		if result.PremiumReceipt > 0 {
			paabuildup.LossComponentBuildup += -paabuildup.LossComponentUnwind
			paabuildup.LossRecoveryBuildup += -paabuildup.LossRecoveryUnwind
		} else {
			paabuildup.LossComponentBuildup += 0
			paabuildup.LossRecoveryBuildup += 0
		}

		paabuildup.DacBuildup += 0
		paabuildup.InsuranceRevenue = 0
		paabuildup.CoverageUnits = 0
		paabuildup.AcquisitionCostAmortizationProportion = 0
		paabuildup.EarnedPremiumProportion = 0

		Step3LossComponentBuildup := paabuildup.LossComponentBuildup
		step3LossComponentUnwind = paabuildup.LossComponentUnwind
		Step3LossRecoveryBuildup := paabuildup.LossRecoveryBuildup
		step3LossRecoveryUnwind = paabuildup.LossRecoveryUnwind
		paabuildup.InsuranceServiceExpense = -step3LossComponentUnwind
		paabuildups = append(paabuildups, paabuildup)

		//Line 4//
		paabuildup.Name = "NB_Premium_Receipt"
		paabuildup.RunDate = result.RunDate
		paabuildup.RunId = csmRun.PaaRunId
		paabuildup.PortfolioName = finance[0].PortfolioName
		paabuildup.ProductCode = result.ProductCode
		paabuildup.IFRS17Group = result.IFRS17Group

		if result.NBPremiumReceipt != 0 {
			paabuildup.ModifiedGMMBel = result.GmmBel
			paabuildup.ModifiedGMMRiskAdjustment = result.GmmRiskAdjustment
			paabuildup.ModifiedGMMReserve = result.GmmReserve
		}
		paabuildup.VariableChange = result.NBPremiumReceipt

		//priorUprBuildup := paabuildup.PaaLrcBuildup
		if portfolio.InsuranceType == Direct {
			paabuildup.PaaLrcBuildup += result.NBPremiumReceipt

			paabuildup.InitialRecognitionLossComponent = utils.FloatPrecision(math.Max(sap[0].IrDiscountedNetCashFlowsLockedin+sap[0].IrRiskAdjustment, 0), AccountingPrecision)
			initialRecognitionReinsuranceClaims := sap[0].IrTreaty1DiscountedClaimsOutgo + sap[0].IrTreaty2DiscountedClaimsOutgo + sap[0].IrTreaty3DiscountedClaimsOutgo
			if sap[0].IrDiscountedCashOutflowsLockedin > 0 {
				paabuildup.InitialRecognitionLossRecovery = paabuildup.InitialRecognitionLossComponent * initialRecognitionReinsuranceClaims / sap[0].IrDiscountedCashOutflowsLockedin //sap[0].IrDiscountedClaimsOutgoLockedin
				//result.PaaLossRecoveryComponent = paabuildup.InitialRecognitionLossRecovery
				sarReinsuranceClaims := sap[0].Treaty1DiscountedCashOutflowsLockedin + sap[0].Treaty2DiscountedCashOutflowsLockedin + sap[0].Treaty3DiscountedCashOutflowsLockedin
				if sarReinsuranceClaims > 0 {
					result.PaaTreaty1InitialLossRecoveryComponent = paabuildup.InitialRecognitionLossRecovery * (sap[0].Treaty1DiscountedCashOutflowsLockedin / sarReinsuranceClaims)
					result.PaaTreaty2InitialLossRecoveryComponent = paabuildup.InitialRecognitionLossRecovery * (sap[0].Treaty2DiscountedCashOutflowsLockedin / sarReinsuranceClaims)
					result.PaaTreaty3InitialLossRecoveryComponent = paabuildup.InitialRecognitionLossRecovery * (sap[0].Treaty3DiscountedCashOutflowsLockedin / sarReinsuranceClaims)
				}

			}

			//if result.NBPremiumReceipt > 0 {
			//	NBlossComponentUnwindDenominator := sap[0].SumFutureEarnedPremium + sap[0].CurrentPeriodEarnedPremium + sap[0].CurrentPeriodAmortisedAcquisition
			//	if NBlossComponentUnwindDenominator > 0 {
			//		paabuildup.LossComponentUnwind = utils.FloatPrecision(paabuildup.InitialRecognitionLossComponent*(sap[0].CurrentPeriodEarnedPremium+sap[0].CurrentPeriodAmortisedAcquisition)/NBlossComponentUnwindDenominator, AccountingPrecision)
			//		paabuildup.LossRecoveryUnwind = utils.FloatPrecision(paabuildup.InitialRecognitionLossRecovery*(sap[0].CurrentPeriodEarnedPremium+sap[0].CurrentPeriodAmortisedAcquisition)/NBlossComponentUnwindDenominator, AccountingPrecision)
			//	}
			//}

			//Loss Component and Loss Recovery Component at Initial Recognition
			if result.NBPremiumReceipt > 0 {
				paabuildup.LossComponentBuildup += math.Max(paabuildup.InitialRecognitionLossComponent-paabuildup.LossComponentUnwind, 0) // initial recognition for new business  //utils.FloatPrecision(math.Max(result.GmmReserve, 0), AccountingPrecision)
				paabuildup.LossRecoveryBuildup += math.Max(paabuildup.InitialRecognitionLossRecovery-paabuildup.LossRecoveryUnwind, 0)
			} else {
				paabuildup.LossComponentBuildup += 0
				paabuildup.LossRecoveryBuildup += 0
			}

			paabuildup.LossComponentAdjustment = paabuildup.LossComponentBuildup + paabuildup.LossComponentUnwind - Step3LossComponentBuildup - step3LossComponentUnwind //result.PaaLossComponent
			paabuildup.LossRecoveryAdjustment = paabuildup.LossRecoveryBuildup + paabuildup.LossRecoveryUnwind - Step3LossRecoveryBuildup - step3LossRecoveryUnwind
			//if sap[0].DiscountedClaimsOutgoLockedin > 0 {
			//	result.PaaLossRecoveryComponent += paabuildup.LossComponentAdjustment * (initialRecognitionReinsuranceClaims / sap[0].DiscountedClaimsOutgoLockedin)
			//}
			paaLossComponentInitialRecognition = paabuildup.LossComponentAdjustment
			//paaLossRecoveryInitialRecognition = paabuildup.LossRecoveryAdjustment

			if initialRecognitionReinsuranceClaims > 0 {
				paaTreaty1LossRecoveryInitialRecognition = paabuildup.LossRecoveryAdjustment * sap[0].IrTreaty1DiscountedClaimsOutgo / initialRecognitionReinsuranceClaims
				paaTreaty2LossRecoveryInitialRecognition = paabuildup.LossRecoveryAdjustment * sap[0].IrTreaty2DiscountedClaimsOutgo / initialRecognitionReinsuranceClaims
				paaTreaty3LossRecoveryInitialRecognition = paabuildup.LossRecoveryAdjustment * sap[0].IrTreaty3DiscountedClaimsOutgo / initialRecognitionReinsuranceClaims
			}

			//if result.NBPremiumReceipt > 0 {
			//	paabuildup.LossRecoveryBuildup += result.PaaLossRecoveryComponent
			//}
			paabuildup.DacBuildup += 0
			paabuildup.InsuranceRevenue = 0
			paabuildup.CoverageUnits = 0
			paabuildup.AcquisitionCostAmortizationProportion = 0
			paabuildup.EarnedPremiumProportion = 0
			step4LossComponentUnwind = paabuildup.LossComponentUnwind - step3LossComponentUnwind
			//step4LossRecoveryUnwind = paabuildup.LossRecoveryUnwind - step3LossRecoveryUnwind
			paabuildup.InsuranceServiceExpense = paaLossComponentInitialRecognition - step4LossComponentUnwind //paabuildup.LossComponentUnwind

		}

		if portfolio.InsuranceType == ProportionalReinsurance {
			paabuildup.VariableChange = finance[0].ReinsurancePremium
			paabuildup.PaaReinsuranceLrcBuildup += finance[0].ReinsurancePremium
			paabuildup.InitialRecognitionLossRecovery = result.PaaInitialLossRecoveryComponent
			paabuildup.PaaReinsuranceDacBuildup += 0
			paabuildup.LossRecoveryBuildup += result.PaaInitialLossRecoveryComponent
			paabuildup.LossRecoveryUnwind = 0
		}

		if portfolio.InsuranceType != Direct && portfolio.InsuranceType != ProportionalReinsurance {
			paabuildup.PaaReinsuranceLrcBuildup += result.NBPremiumReceipt
		}

		paabuildups = append(paabuildups, paabuildup)

		//Line 5//
		paabuildup.Name = "Acquisition_Expense_Paid"
		paabuildup.RunDate = result.RunDate
		paabuildup.RunId = csmRun.PaaRunId
		paabuildup.PortfolioName = finance[0].PortfolioName
		paabuildup.ProductCode = result.ProductCode
		paabuildup.IFRS17Group = result.IFRS17Group
		if finance[0].DacBuildupIndicator {
			paabuildup.VariableChange = finance[0].AcquisitionCostPaid //result.AcquisitionExpenses
		} else {
			paabuildup.VariableChange = 0
		}

		if portfolio.InsuranceType == Direct {
			paabuildup.PaaLrcBuildup += -paabuildup.VariableChange
		}
		if portfolio.InsuranceType != Direct {
			paabuildup.PaaReinsuranceLrcBuildup += -paabuildup.VariableChange
		}
		paabuildup.LossComponentUnwind = 0
		paabuildup.InitialRecognitionLossComponent = 0
		paabuildup.LossRecoveryUnwind = 0
		paabuildup.InitialRecognitionLossRecovery = 0
		paabuildup.InsuranceServiceExpense = 0

		//if result.GmmReserve < 0 {
		//	result.PaaLossComponent = 0
		//} else {
		//	result.PaaLossComponent = utils.FloatPrecision(math.Max(result.GmmReserve-math.Max(paabuildup.PaaLrcBuildup, 0), 0), AccountingPrecision)
		//	sarReinsuranceClaims := sap[0].Treaty1DiscountedCashOutflowsLockedin + sap[0].Treaty2DiscountedCashOutflowsLockedin + sap[0].Treaty3DiscountedCashOutflowsLockedin
		//	if sap[0].DiscountedCashOutflows > 0 {
		//		result.PaaLossRecoveryComponent = result.PaaLossComponent * (sarReinsuranceClaims / sap[0].DiscountedCashOutflows)
		//	}
		//}

		if sap[0].CurrentPeriodEarnedPremium == 0 {
			paabuildup.LossComponentBuildup += 0
			paabuildup.LossRecoveryBuildup += 0
		} else {
			paabuildup.LossComponentBuildup += 0 //result.PaaLossComponent
			paabuildup.LossRecoveryBuildup += 0
		}

		if portfolio.InsuranceType == Direct {
			paabuildup.DacBuildup += paabuildup.VariableChange
		}

		if portfolio.InsuranceType != Direct {
			paabuildup.PaaReinsuranceDacBuildup += paabuildup.VariableChange
		}
		paabuildup.InsuranceRevenue = 0
		paabuildup.CoverageUnits = 0
		paabuildup.AcquisitionCostAmortizationProportion = 0
		paabuildup.EarnedPremiumProportion = 0
		paabuildup.LossComponentAdjustment = 0
		paabuildup.LossRecoveryAdjustment = 0
		paabuildups = append(paabuildups, paabuildup)

		//Line 6//
		paabuildup.Name = "Amortised_Acquisition_Expense"
		paabuildup.RunDate = result.RunDate
		paabuildup.RunId = csmRun.PaaRunId
		paabuildup.PortfolioName = finance[0].PortfolioName
		paabuildup.ProductCode = result.ProductCode
		paabuildup.IFRS17Group = result.IFRS17Group

		if finance[0].ExitStatus == Exit {
			if portfolio.InsuranceType == Direct {
				paabuildup.VariableChange = paabuildup.DacBuildup
			}
			if portfolio.InsuranceType != Direct {
				paabuildup.VariableChange = paabuildup.PaaReinsuranceDacBuildup
			}
		} else {
			if paabuildup.DacBuildup > 0 {
				paabuildup.VariableChange = sap[0].CurrentPeriodAmortisedAcquisition
			} else {
				paabuildup.VariableChange = 0
			}
		}
		result.AmortisedAcquisitionCost = paabuildup.VariableChange
		if portfolio.InsuranceType == Direct {
			paabuildup.DacBuildup = paabuildup.DacBuildup - paabuildup.VariableChange // math.Max(paabuildup.DacBuildup-paabuildup.VariableChange, 0)
		}
		if portfolio.InsuranceType != Direct {
			paabuildup.PaaReinsuranceDacBuildup = paabuildup.PaaReinsuranceDacBuildup - paabuildup.VariableChange
		}
		paabuildup.InsuranceRevenue = paabuildup.VariableChange
		paabuildup.InsuranceServiceExpense = paabuildup.VariableChange
		tempDacRelease = paabuildup.InsuranceRevenue
		paabuildup.CoverageUnits = 0
		paabuildup.AcquisitionCostAmortizationProportion = 0
		paabuildup.EarnedPremiumProportion = 0

		paabuildups = append(paabuildups, paabuildup)

		//Line 7//
		paabuildup.Name = "Interest_Accretion"
		paabuildup.RunDate = result.RunDate
		paabuildup.RunId = csmRun.PaaRunId
		paabuildup.PortfolioName = finance[0].PortfolioName
		paabuildup.ProductCode = result.ProductCode
		paabuildup.IFRS17Group = result.IFRS17Group
		paabuildup.VariableChange = 0
		paabuildup.InsuranceRevenue = paabuildup.VariableChange
		paabuildup.CoverageUnits = 0
		paabuildup.AcquisitionCostAmortizationProportion = 0
		paabuildup.EarnedPremiumProportion = 0
		paabuildup.InsuranceServiceExpense = 0

		paabuildups = append(paabuildups, paabuildup)

		//Line 8//
		paabuildup.EarnedPremiumProportion = 0.0
		paabuildup.Name = "Earned_Premium"
		paabuildup.RunDate = result.RunDate
		paabuildup.RunId = csmRun.PaaRunId
		paabuildup.PortfolioName = finance[0].PortfolioName
		paabuildup.ProductCode = result.ProductCode
		paabuildup.IFRS17Group = result.IFRS17Group

		//if paabuildup.UprBuildup > 0 && len(portfolio) > 0 {
		//	if portfolio[0].PremiumEarningPattern == PassageofTime && len(mp) > 0 {
		//		if mp[0].TermMonths != 0 {
		//			paabuildup.EarnedPremiumProportion = float64(current_period_duration) / float64(mp[0].TermMonths)
		//		} else {
		//			paabuildup.EarnedPremiumProportion = 0
		//		}
		//	}
		//	if len(userPattern) > 0 {
		//		if portfolio[0].PremiumEarningPattern == SpecifiedbyUser {
		//			var sumCoveredUnits float64
		//			if sumFutureRiskUnits > 0 {
		//				for i := duration_start_period; i <= duration_start_period+current_period_duration; i++ {
		//					sumCoveredUnits += userPattern[i].RiskUnit
		//				}
		//				paabuildup.EarnedPremiumProportion = sumCoveredUnits / sumFutureRiskUnits
		//			}
		//		}
		//	}
		//
		//}

		//paabuildup.EarnedPremiumProportion = tempEarnedPremiumProportion
		var uprvarChange float64
		//var LcvarChange float64
		if finance[0].ExitStatus == Exit {
			if portfolio.InsuranceType == Direct {
				uprvarChange = utils.FloatPrecision(paabuildup.PaaLrcBuildup, AccountingPrecision)
				//LcvarChange = utils.FloatPrecision(paabuildup.LossComponentBuildup,AccountingPrecision)
			}
			if portfolio.InsuranceType != Direct {
				uprvarChange = utils.FloatPrecision(paabuildup.PaaReinsuranceLrcBuildup, AccountingPrecision)
			}
		} else {
			if portfolio.InsuranceType == Direct {
				if paabuildup.PaaLrcBuildup > 0 {
					uprvarChange = sap[0].CurrentPeriodEarnedPremium
				}
			}
			if portfolio.InsuranceType != Direct && portfolio.InsuranceType != ProportionalReinsurance {
				if paabuildup.PaaReinsuranceLrcBuildup > 0 {
					uprvarChange = sap[0].CurrentPeriodEarnedPremium
				}
			}

			if portfolio.InsuranceType == ProportionalReinsurance {
				if paabuildup.PaaReinsuranceLrcBuildup > 0 {
					uprvarChange = result.AllocatedReinsurancePremium
				}
			}

			//utils.FloatPrecision(paabuildup.UprBuildup*paabuildup.EarnedPremiumProportion, AccountingPrecision)
			//LcvarChange = utils.FloatPrecision(paabuildup.LossComponentBuildup*paabuildup.EarnedPremiumProportion,AccountingPrecision)
		}
		paabuildup.VariableChange = uprvarChange //+ LcvarChange
		if portfolio.InsuranceType == Direct {
			paabuildup.PaaLrcBuildup = paabuildup.PaaLrcBuildup - uprvarChange
		}
		if portfolio.InsuranceType != Direct {
			paabuildup.PaaReinsuranceLrcBuildup = paabuildup.PaaReinsuranceLrcBuildup - uprvarChange
		}

		if portfolio.InsuranceType == Direct {
			//math.Max(paabuildup.UprBuildup-uprvarChange, 0)

			NBlossComponentUnwindDenominator := sap[0].SumFutureEarnedPremium + sap[0].CurrentPeriodEarnedPremium + sap[0].CurrentPeriodAmortisedAcquisition
			if NBlossComponentUnwindDenominator > 0 {
				paabuildup.LossComponentUnwind = utils.FloatPrecision(paabuildup.LossComponentBuildup*(sap[0].CurrentPeriodEarnedPremium+sap[0].CurrentPeriodAmortisedAcquisition)/NBlossComponentUnwindDenominator, AccountingPrecision)
				paabuildup.LossRecoveryUnwind = utils.FloatPrecision(paabuildup.LossRecoveryBuildup*(sap[0].CurrentPeriodEarnedPremium+sap[0].CurrentPeriodAmortisedAcquisition)/NBlossComponentUnwindDenominator, AccountingPrecision)

			}

			step8LossComponentUnwind = paabuildup.LossComponentUnwind
			step8LossRecoveryUnwind = paabuildup.LossRecoveryUnwind
			prevStepLossComponentBuildup := paabuildup.LossComponentBuildup
			prevStepLossRecoveryBuildup := paabuildup.LossRecoveryBuildup
			paabuildup.LossComponentAdjustment = -paabuildup.LossComponentUnwind //paabuildup.LossComponentAdjustment = paabuildup.LossComponentBuildup - paabuildup.LossComponentUnwind
			paabuildup.LossRecoveryAdjustment = -paabuildup.LossRecoveryUnwind
			//result.PaaLossRecoveryUnwind = step8LossRecoveryUnwind
			if result.GmmReserve < 0 {
				paabuildup.LossComponentBuildup += 0
				paabuildup.LossRecoveryBuildup += 0
			} else {
				paabuildup.LossComponentBuildup = utils.FloatPrecision(prevStepLossComponentBuildup+paabuildup.LossComponentAdjustment, AccountingPrecision)
				//paabuildup.LossComponentBuildup += utils.FloatPrecision(math.Max(result.GmmReserve-math.Max(paabuildup.PaaLrcBuildup+prevStepLossComponentBuildup, 0), 0), AccountingPrecision)
			}

			//paabuildup.LossComponentAdjustment = paabuildup.LossComponentBuildup - prevStepLossComponentBuildup

			sarReinsuranceClaims := sap[0].Treaty1DiscountedCashOutflowsLockedin + sap[0].Treaty2DiscountedCashOutflowsLockedin + sap[0].Treaty3DiscountedCashOutflowsLockedin
			if sap[0].DiscountedClaimsOutgoLockedin > 0 {
				//paabuildup.LossRecoveryAdjustment = paabuildup.LossComponentAdjustment * (sarReinsuranceClaims / sap[0].DiscountedClaimsOutgoLockedin)
				if sarReinsuranceClaims > 0 {
					//paaTreaty1LossRecoveryAdjustmentOut = paabuildup.LossRecoveryAdjustment * sap[0].Treaty1DiscountedCashOutflowsLockedin / sarReinsuranceClaims
					//paaTreaty2LossRecoveryAdjustmentOut = paabuildup.LossRecoveryAdjustment * sap[0].Treaty2DiscountedCashOutflowsLockedin / sarReinsuranceClaims
					//paaTreaty3LossRecoveryAdjustmentOut = paabuildup.LossRecoveryAdjustment * sap[0].Treaty3DiscountedCashOutflowsLockedin / sarReinsuranceClaims
					result.PaaTreaty1LossRecoveryUnwind = step8LossRecoveryUnwind * sap[0].Treaty1DiscountedCashOutflowsLockedin / sarReinsuranceClaims
					result.PaaTreaty2LossRecoveryUnwind = step8LossRecoveryUnwind * sap[0].Treaty2DiscountedCashOutflowsLockedin / sarReinsuranceClaims
					result.PaaTreaty3LossRecoveryUnwind = step8LossRecoveryUnwind * sap[0].Treaty3DiscountedCashOutflowsLockedin / sarReinsuranceClaims

				}
				paabuildup.LossRecoveryBuildup = prevStepLossRecoveryBuildup + paabuildup.LossRecoveryAdjustment
			}

			//if sap[0].CurrentPeriodEarnedPremium == 0 { //mp[0].DurationInForceMonths == 0 {
			//	paabuildup.LossComponentBuildup += 0
			//	paabuildup.LossRecoveryBuildup += 0
			//} else {
			//	paabuildup.LossComponentBuildup = result.PaaLossComponent
			//	paabuildup.LossRecoveryBuildup = result.PaaLossRecoveryComponent
			//}

			result.PaaLossComponent = paabuildup.LossComponentBuildup

			//paaLossComponentAdjustmentOut = paabuildup.LossComponentAdjustment
			//paaLossRecoveryAdjustmentOut = paabuildup.LossRecoveryAdjustment
			//if sarReinsuranceClaims > 0 {
			//	paaTreaty1LossRecoveryAdjustmentOut = paaLossRecoveryAdjustmentOut * (sap[0].Treaty1DiscountedCashOutflowsLockedin / sarReinsuranceClaims)
			//	paaTreaty2LossRecoveryAdjustmentOut = paaLossRecoveryAdjustmentOut * (sap[0].Treaty2DiscountedCashOutflowsLockedin / sarReinsuranceClaims)
			//	paaTreaty3LossRecoveryAdjustmentOut = paaLossRecoveryAdjustmentOut * (sap[0].Treaty3DiscountedCashOutflowsLockedin / sarReinsuranceClaims)
			//}
			paabuildup.CoverageUnits = 0
			paabuildup.AcquisitionCostAmortizationProportion = 0

		}

		if portfolio.InsuranceType == ProportionalReinsurance {
			paabuildup.LossRecoveryUnwind = result.PaaLossRecoveryUnwind
			paabuildup.LossRecoveryBuildup += -result.PaaLossRecoveryUnwind
		}
		result.EarnedPremium = utils.FloatPrecision(uprvarChange, AccountingPrecision)
		result.InsuranceRevenue = result.AmortisedAcquisitionCost + result.EarnedPremium
		paabuildup.InsuranceRevenue = paabuildup.VariableChange
		paabuildup.InsuranceServiceExpense = paabuildup.LossComponentAdjustment

		paabuildups = append(paabuildups, paabuildup)

		//Line 8//
		paabuildup.Name = "C/F(Lockedin)"
		paabuildup.RunDate = result.RunDate
		paabuildup.RunId = csmRun.PaaRunId
		paabuildup.PortfolioName = finance[0].PortfolioName
		paabuildup.ProductCode = result.ProductCode
		paabuildup.IFRS17Group = result.IFRS17Group
		paabuildup.VariableChange = 0
		paabuildup.CoverageUnits = 0
		paabuildup.AcquisitionCostAmortizationProportion = 0
		paabuildup.InsuranceRevenue = 0
		paabuildup.LossComponentAdjustment = 0
		paabuildup.LossComponentUnwind = 0
		paabuildup.InsuranceServiceExpense = 0
		paabuildup.LossRecoveryUnwind = 0
		paabuildup.LossRecoveryAdjustment = 0

		paabuildups = append(paabuildups, paabuildup)

		//Line 9//
		paabuildup.Name = "C/F(Current)"
		paabuildup.RunDate = result.RunDate
		paabuildup.RunId = csmRun.PaaRunId
		paabuildup.PortfolioName = finance[0].PortfolioName
		paabuildup.ProductCode = result.ProductCode
		paabuildup.IFRS17Group = result.IFRS17Group
		paabuildup.VariableChange = 0
		paabuildup.CoverageUnits = 0
		paabuildup.AcquisitionCostAmortizationProportion = 0
		paabuildup.InsuranceRevenue = 0

		paabuildups = append(paabuildups, paabuildup)

		if portfolio.InsuranceType == Direct {
			result.DeferredAcquisitionExpenses = utils.FloatPrecision(paabuildup.DacBuildup, AccountingPrecision)
			result.PaaLiabilityRemainingCoverage = utils.FloatPrecision(paabuildup.PaaLrcBuildup, AccountingPrecision)
			//reinsuranceWrittenPrem := sap[0].Treaty1WrittenPremium + sap[0].Treaty2WrittenPremium + sap[0].Treaty3WrittenPremium
			//if reinsuranceWrittenPrem > 0 {
			//	reinsuranceTreaty1Lrc = utils.FloatPrecision(result.ReinsurancePremiumPaid*sap[0].Treaty1WrittenPremium/reinsuranceWrittenPrem-sap[0].Treaty1CurrentPeriodEarnedPremium, AccountingPrecision)
			//	reinsuranceTreaty2Lrc = utils.FloatPrecision(result.ReinsurancePremiumPaid*sap[0].Treaty2WrittenPremium/reinsuranceWrittenPrem-sap[0].Treaty2CurrentPeriodEarnedPremium, AccountingPrecision)
			//	reinsuranceTreaty3Lrc = utils.FloatPrecision(result.ReinsurancePremiumPaid*sap[0].Treaty3WrittenPremium/reinsuranceWrittenPrem-sap[0].Treaty3CurrentPeriodEarnedPremium, AccountingPrecision)
			//}
			//result.PaaReinsuranceLrc = reinsuranceTreaty1Lrc + reinsuranceTreaty2Lrc + reinsuranceTreaty3Lrc
			//result.PaaTreaty1Lrc = reinsuranceTreaty1Lrc
			//result.PaaTreaty2Lrc = reinsuranceTreaty2Lrc
			//result.PaaTreaty3Lrc = reinsuranceTreaty3Lrc
			//if len(sap) > 0 && len(finance) > 0 {
			//	//if sap[0].CurrentPeriodEarnedPremium == 0 {
			//	//	result.PaaLossComponent = utils.FloatPrecision(math.Max(result.GmmReserve, 0), AccountingPrecision)
			//	//} else {
			//	result.PaaLossComponent = paabuildup.LossComponentBuildup //utils.FloatPrecision(math.Max(result.GmmReserve-math.Max(result.PaaLiabilityRemainingCoverage, 0), 0), AccountingPrecision)
			//	sarReinsuranceClaims := sap[0].Treaty1DiscountedCashOutflowsLockedin + sap[0].Treaty2DiscountedCashOutflowsLockedin + sap[0].Treaty3DiscountedCashOutflowsLockedin
			//	if sap[0].DiscountedCashOutflows > 0 {
			//		result.PaaLossRecoveryComponent = result.PaaLossComponent * (sarReinsuranceClaims / sap[0].DiscountedCashOutflows)
			//	}
			//	//	}
			//}
		}

		if portfolio.InsuranceType != Direct {
			result.PaaReinsuranceLrc = paabuildup.PaaReinsuranceLrcBuildup
			result.PaaReinsuranceDac = paabuildup.PaaReinsuranceDacBuildup
		}

		if portfolio.InsuranceType == Direct {
			result.IncurredClaims = finance[0].IncurredClaims
			result.ClaimsPaid = finance[0].ClaimsPaid
			//result.Treaty1Recovery = finance[0].Treaty1Recovery
			//result.Treaty2Recovery = finance[0].Treaty2Recovery
			//result.Treaty3Recovery = finance[0].Treaty3Recovery
			result.IncurredExpenses = finance[0].AttributableExpenses + finance[0].ClaimsExpenses
			result.NonAttributableExpenses = finance[0].NonAttributableExpenses

		}

		// PAA onerous contract detection
		if result.PaaLossComponent > 0 {
			result.IsOnerous = true
			result.OnerousReason = "PAA loss component present"
		}

		results = append(results, result)

		if !csmRun.PaaEligibilityTest {
			var bsr models.BalanceSheetRecord

			bsr.ProductCode = result.ProductCode
			bsr.IFRS17Group = result.IFRS17Group
			bsr.CsmRunID = result.CsmRunID
			bsr.MeasurementType = csmRun.MeasurementType
			bsr.Date = csmRun.RunDate
			if portfolio.InsuranceType == Direct {
				bsr.PaaLiabilityRemainingCoverage = result.PaaLiabilityRemainingCoverage
				bsr.PaaDAC = result.DeferredAcquisitionExpenses
				bsr.PaaLossComponent = result.PaaLossComponent
				//bsr.PaaReinsuranceLRC = result.PaaReinsuranceLrc
				//bsr.PaaReinsuranceDAC = result.PaaReinsuranceDac

			}
			if portfolio.InsuranceType != Direct {
				bsr.PaaReinsuranceLRC = result.PaaReinsuranceLrc
				bsr.PaaReinsuranceDAC = result.PaaReinsuranceDac
				bsr.PaaLossRecoveryComponent = result.PaaLossRecoveryComponent
				//bsr.PaaLossComponent = result.PaaLossComponent
			}

			err = DB.Where("product_code = ? and measurement_type = ? and ifrs17_group = ? and date = ?", result.ProductCode, csmRun.MeasurementType, result.IFRS17Group, csmRun.RunDate).Delete(&models.BalanceSheetRecord{}).Error
			if err != nil {
				fmt.Println("balance sheet record delete error:", err)
			}
			err = DB.Create(&bsr).Error
			if err != nil {
				fmt.Println("balance sheet record error: ", err)
			}
		} //PAA Eligibility Test Closing
	}

	// Check for PAA results run on the same run date and delete.
	DB.Where("run_date = ? and ifrs17_group = ? and csm_run_id=? ", csmRun.RunDate, group.IFRS17Group, csmRun.PaaRunId).Delete(&models.PAAResult{})
	err = DB.CreateInBatches(&results, 100).Error
	if err != nil {
		fmt.Println(err)
	}

	if !csmRun.PaaEligibilityTest {
		// Check for PAABuildup results run on the same run date and delete.
		DB.Where("run_date = ? and ifrs17_group = ? and run_id=?", csmRun.RunDate, group.IFRS17Group, csmRun.PaaRunId).Delete(&models.PAABuildUp{})
		err = DB.CreateInBatches(&paabuildups, 100).Error
		if err != nil {
			fmt.Println(err)
		}

		jt.IFRS17Group = result.IFRS17Group
		jt.ProductCode = result.ProductCode
		jt.RunDate = csmRun.RunDate
		jt.CsmRunID = csmRun.ID
		jt.MeasurementType = csmRun.MeasurementType
		jt.PaaEarnedPremium = result.EarnedPremium               //+ tempDacRelease
		jt.PaaLossComponent = paaLossComponentInitialRecognition //result.PaaLossComponent
		jt.PaaLossRecoveryComponent = result.PaaInitialLossRecoveryComponent + paaTreaty1LossRecoveryInitialRecognition*0 + paaTreaty2LossRecoveryInitialRecognition*0 + paaTreaty3LossRecoveryInitialRecognition*0
		jt.PaaLossRecoveryAdjustment = paaLossRecoveryAdjustmentOut + paaTreaty1LossRecoveryAdjustmentOut*0 + paaTreaty2LossRecoveryAdjustmentOut*0 + paaTreaty3LossRecoveryAdjustmentOut*0
		jt.PaaLossRecoveryUnwind = result.PaaLossRecoveryUnwind //step8LossRecoveryUnwind //-step3LossRecoveryUnwind - step4LossRecoveryUnwind
		//jt.LossComponentUnwind = paabuildup.LossComponentUnwind
		jt.PaaLossComponentAdjustment = step8LossComponentUnwind //-step3LossComponentUnwind - step4LossComponentUnwind + paaLossComponentAdjustmentOut
		jt.AmortizationAcquisitionCF = tempDacRelease
		jt.ClaimsIncurred = result.IncurredClaims
		jt.ExpensesIncurred = result.IncurredExpenses
		jt.NonAttributableExpensesIncurred = result.NonAttributableExpenses
		if portfolio.InsuranceType != Direct {
			jt.PaaReinsurancePremium = result.AllocatedReinsurancePremium
			jt.ReinsuranceFlatCommission = result.AllocatedReinsuranceFlatCommission
			jt.ReinsuranceRecovery = result.ReinsuranceRecovery
			jt.ReinsuranceReinstatementPremium = result.ReinsuranceReinstatementPremium
			jt.ReinsuranceProvisionalCommission = result.ReinsuranceProvisionalCommission
			jt.ReinsuranceUltimateCommission = result.ReinsuranceUltimateCommission
			jt.ReinsuranceProfitCommission = result.ReinsuranceProfitCommission
			jt.ReinsuranceInvestmentComponent = result.ReinsuranceInvestmentComponent
		}

		err = DB.Where("run_date = ? and measurement_type = ? and product_code=? and ifrs17_group = ?", jt.RunDate, jt.MeasurementType, jt.ProductCode, jt.IFRS17Group).Delete(&jt).Error
		if err != nil {
			log.Info(err.Error())
		}

		err = DB.Where("run_date = ? and measurement_type = ? and product_code=? and ifrs17_group = ?", jt.RunDate, jt.MeasurementType, jt.ProductCode, jt.IFRS17Group).Save(&jt).Error
		if err != nil {
			log.Info(err.Error())
		}
	} //PAA Eligibility test Closing Block

	// Project PAA Eligibility Cash Flows

	if csmRun.PaaEligibilityTest {
		var paaEligibilityTestResult = make([]models.PAAEligibilityTestResult, len(sap))
		var paacsmrelease float64
		var prevRiskAdjustment float64
		var respLockedin float64
		var lockedinYear int
		var lockedinMonth int
		var prevUPR float64
		var prevDAC float64
		var prevLC float64
		lockedinYear, _ = strconv.Atoi(sap[0].LockedinYear)
		lockedinMonth, _ = strconv.Atoi(sap[0].LockedinYear)

		for projMonth := 0; projMonth < len(sap); projMonth++ {
			paaEligibilityTestResult[projMonth].RunDate = result.RunDate
			paaEligibilityTestResult[projMonth].CsmRunID = csmRun.PaaRunId
			paaEligibilityTestResult[projMonth].RunId = csmRun.RunID
			paaEligibilityTestResult[projMonth].ProductCode = result.ProductCode
			paaEligibilityTestResult[projMonth].IFRS17Group = result.IFRS17Group
			paaEligibilityTestResult[projMonth].ProjectionMonth = sap[projMonth].ProjectionMonth
			if portfolio.DiscountOption == Discounted {
				respLockedin = GetPaaForwardRate(projMonth+1, lockedinYear, lockedinMonth, sap[0].YieldCurveCode)
			} else {
				respLockedin = 0
			}

			if projMonth == 0 {
				paaEligibilityTestResult[projMonth].EarnedPremium = 0
				paaEligibilityTestResult[projMonth].UnearnedPremiumReserve = 0  // result.UnearnedPremiumReserve
				paaEligibilityTestResult[projMonth].DeferredAcquisitionCost = 0 //result.DeferredAcquisitionExpenses
				paaEligibilityTestResult[projMonth].AmortisedAcquisitionCost = 0
				paaEligibilityTestResult[projMonth].PaaLossComponent = result.PaaLossComponent
				paaEligibilityTestResult[projMonth].GmmBel = result.GmmBel
				paaEligibilityTestResult[projMonth].GmmRiskAdjustment = sap[0].RiskAdjustment //result.GmmRiskAdjustment
				paaEligibilityTestResult[projMonth].GmmReserve = result.GmmReserve
				paaEligibilityTestResult[projMonth].GMMCsm = -math.Min(result.GmmReserve, 0)
				paaEligibilityTestResult[projMonth].GMMCsmRelease = 0
				paaEligibilityTestResult[projMonth].RevenueVariance = 0
				paaEligibilityTestResult[projMonth].RevenueVarianceProportion = 0
				prevUPR = paaEligibilityTestResult[projMonth].UnearnedPremiumReserve
				prevDAC = paaEligibilityTestResult[projMonth].DeferredAcquisitionCost
				prevLC = result.PaaLossComponent

			} else {

				//sap[projMonth].EarnedPremium

				//paaEligibilityTestResult[projMonth].DeferredAcquisitionCost = result.DeferredAcquisitionExpenses
				coverageunitsDenominator := sap[projMonth].SumFutureEarnedPremium + sap[projMonth].CoverageUnits
				if coverageunitsDenominator > 0 {
					paaEligibilityTestResult[projMonth].AmortisedAcquisitionCost = sap[0].SumFutureAcquisitionCost * (sap[projMonth].CoverageUnits / sap[0].SumFutureEarnedPremium)                                                      //prevDAC * (sap[projMonth].CoverageUnits / (sap[projMonth].SumFutureEarnedPremium + sap[projMonth].CoverageUnits))                           //(sap[projMonth].CoverageUnits / sap[0].SumFutureEarnedPremium)
					paaEligibilityTestResult[projMonth].EarnedPremium = (sap[0].SumFutureEarnedPremium - sap[0].SumFutureAcquisitionCost) * (sap[projMonth].CoverageUnits / sap[0].SumFutureEarnedPremium) * math.Pow(1+respLockedin, 0) //prevUPR * math.Pow(1+respLockedin, 1.0/12.0) * (sap[projMonth].CoverageUnits / (sap[projMonth].SumFutureEarnedPremium + sap[projMonth].CoverageUnits)) //(sap[projMonth].CoverageUnits / sap[0].SumFutureEarnedPremium)
				}

				if paaEligibilityTestResult[projMonth].EarnedPremium != 0 {
					paaEligibilityTestResult[projMonth].UnearnedPremiumReserve = prevUPR + sap[projMonth].NbPremiumReceipt + sap[projMonth].PremiumReceipt - sap[projMonth].InitialCommissionOutgo - sap[projMonth].InitialExpenseOutgo - paaEligibilityTestResult[projMonth].EarnedPremium //result.UnearnedPremiumReserve
					paaEligibilityTestResult[projMonth].DeferredAcquisitionCost = prevDAC + sap[projMonth].InitialExpenseOutgo + sap[projMonth].InitialCommissionOutgo - paaEligibilityTestResult[projMonth].AmortisedAcquisitionCost

				} else {
					paaEligibilityTestResult[projMonth].UnearnedPremiumReserve = 0
					paaEligibilityTestResult[projMonth].DeferredAcquisitionCost = 0
				}
				prevUPR = paaEligibilityTestResult[projMonth].UnearnedPremiumReserve
				prevDAC = paaEligibilityTestResult[projMonth].DeferredAcquisitionCost
				paaEligibilityTestResult[projMonth].GmmBel = sap[projMonth].DiscountedNetCashFlowsLockedin
				paaEligibilityTestResult[projMonth].GmmRiskAdjustment = sap[projMonth].RiskAdjustment
				paaEligibilityTestResult[projMonth].GmmReserve = paaEligibilityTestResult[projMonth].GmmBel + paaEligibilityTestResult[projMonth].GmmRiskAdjustment

				if sap[0].SumFutureEarnedPremium != 0 {
					paaEligibilityTestResult[projMonth].GMMCsmRelease = paaEligibilityTestResult[0].GMMCsm * (sap[projMonth].CoverageUnits / sap[0].SumFutureEarnedPremium)
					paaEligibilityTestResult[projMonth].GMMDacRelease = sap[0].SumFutureAcquisitionCost * (sap[projMonth].CoverageUnits / sap[0].SumFutureEarnedPremium)
					paaEligibilityTestResult[projMonth].PAALossComponentRelease = result.PaaLossComponent * (sap[projMonth].CoverageUnits / sap[0].SumFutureEarnedPremium)

				}

				paaEligibilityTestResult[projMonth].PaaLossComponent = prevLC - paaEligibilityTestResult[projMonth].PAALossComponentRelease
				prevLC = paaEligibilityTestResult[projMonth].PaaLossComponent
				paacsmrelease += paaEligibilityTestResult[projMonth].GMMCsmRelease
				paaEligibilityTestResult[projMonth].GMMCsm = math.Max(paaEligibilityTestResult[0].GMMCsm-paacsmrelease, 0)
				paaEligibilityTestResult[projMonth].GMMExpectedOutflows = sap[projMonth].ClaimsOutgo + sap[projMonth].ClaimsExpenseOutgo + sap[projMonth].RenewalCommissionOutgo + sap[projMonth].MaintenanceExpenseOutgo
				paaEligibilityTestResult[projMonth].GMMRiskAdjustmentRelease = math.Max(prevRiskAdjustment-paaEligibilityTestResult[projMonth].GmmRiskAdjustment, 0)

				paaEligibilityTestResult[projMonth].GMMRevenue = paaEligibilityTestResult[projMonth].GMMCsmRelease + paaEligibilityTestResult[projMonth].GMMExpectedOutflows + paaEligibilityTestResult[projMonth].GMMRiskAdjustmentRelease + paaEligibilityTestResult[projMonth].GMMDacRelease
				paaEligibilityTestResult[projMonth].PAARevenue = paaEligibilityTestResult[projMonth].EarnedPremium + paaEligibilityTestResult[projMonth].AmortisedAcquisitionCost + paaEligibilityTestResult[projMonth].PAALossComponentRelease

				if paaEligibilityTestResult[projMonth].GMMRevenue != 0 {
					paaEligibilityTestResult[projMonth].RevenueVariance = math.Abs(paaEligibilityTestResult[projMonth].PAARevenue - paaEligibilityTestResult[projMonth].GMMRevenue)
					paaEligibilityTestResult[projMonth].RevenueVarianceProportion = paaEligibilityTestResult[projMonth].RevenueVariance / paaEligibilityTestResult[projMonth].GMMRevenue
				}

				paaEligibilityTestResult[projMonth].GMMLrc = utils.FloatPrecision(paaEligibilityTestResult[projMonth].GmmReserve+paaEligibilityTestResult[projMonth].GMMCsm, 2)
				paaEligibilityTestResult[projMonth].PAALrc = utils.FloatPrecision(paaEligibilityTestResult[projMonth].UnearnedPremiumReserve+paaEligibilityTestResult[projMonth].PaaLossComponent, 2)

				if utils.FloatPrecision(paaEligibilityTestResult[projMonth].GmmReserve+paaEligibilityTestResult[projMonth].GMMCsm, 1) != 0 {
					paaEligibilityTestResult[projMonth].LRCVariance = math.Abs(paaEligibilityTestResult[projMonth].PAALrc - paaEligibilityTestResult[projMonth].GMMLrc)

					paaEligibilityTestResult[projMonth].LRCVarianceProportion = paaEligibilityTestResult[projMonth].LRCVariance / (paaEligibilityTestResult[projMonth].GmmReserve + paaEligibilityTestResult[projMonth].GMMCsm)
				}
			}
			prevRiskAdjustment = paaEligibilityTestResult[projMonth].GmmRiskAdjustment
			paaEligibilityTestResults = append(paaEligibilityTestResults, paaEligibilityTestResult[projMonth])
		}

		err = DB.Where("run_date = ? and run_id = ? and product_code=? and ifrs17_group = ?", paaEligibilityTestResult[0].RunDate, csmRun.PaaRunId, paaEligibilityTestResult[0].ProductCode, paaEligibilityTestResult[0].IFRS17Group).Delete(&paaEligibilityTestResults).Error
		if err != nil {
			log.Info(err.Error())
		}

		err = DB.Where("run_date = ? and run_id = ? and product_code=? and ifrs17_group = ?", csmRun.RunDate, csmRun.PaaRunId, paaEligibilityTestResult[0].ProductCode, paaEligibilityTestResult[0].IFRS17Group).Save(&paaEligibilityTestResults).Error
		if err != nil {
			log.Info(err.Error())
		}
		//err = DB.Create(&paaEligibilityTestResults).Error
		//if err != nil {
		//	log.Error(err)
		//}
	}
	return results
}

func getsumFutureRiskUnit(year int, prodCode string, durationM int) float64 {
	var sumRiskUnit float64
	key := strconv.Itoa(year) + "_" + prodCode
	cacheKey := key
	cached, found := PaaCache.Get(cacheKey)
	if found {
		return cached.(float64)
	}
	if found {
		result := cached.(float64)
		if result > 0 {
			return result
		}
	}
	if !found {

		query := fmt.Sprintf("SELECT sum(risk_unit) as sumRiskUnit from premium_earning_patterns where year='%d' and product_code= '%s' and duration_in_force> '%d'", year, prodCode, durationM)
		err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Raw(query).Find(&sumRiskUnit).Error
		})
		if err != nil {
			log.Info(err.Error())
		}

		success := PaaCache.Set(cacheKey, sumRiskUnit, 1)
		if !success {
			fmt.Println("cache: key not stored")
		}
	}
	return sumRiskUnit
}

func getYieldCurveCode(year int, prodCode, yieldCurveBasis string) string {
	key := strconv.Itoa(year) + "_" + prodCode + "_" + yieldCurveBasis
	cacheKey := key
	cached, found := Cache.Get(cacheKey)
	if found {
		return cached.(string)
	}
	if found {
		result := cached.(string)
		if result != "" {
			return result
		}
	}

	var param models.ProductParameters
	if !found {
		//Remember to update to 1FRS17 from IFRS
		err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Where("year = ? and product_code = ? and basis = ?", year, prodCode, yieldCurveBasis).First(&param).Error
		})
		if err != nil {
			log.Info(err.Error())
		}

		success := Cache.Set(cacheKey, param.YieldCurveCode, 1)
		if !success {
			fmt.Println("cache: key not stored")
		}
	}
	return param.YieldCurveCode

}
func GetDistinctValues[T comparable](input []T) []T {
	// Create a map to track seen values
	seen := make(map[T]bool)
	// Create a slice to store unique values
	result := []T{}

	// Iterate through input slice
	for _, value := range input {
		// If value hasn't been seen, add it to result
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}

	return result
}

// ─── Run Approval Workflow ───────────────────────────────────────────────────

// ReviewCsmRun transitions a completed run from draft → reviewed.
func ReviewCsmRun(runID int, notes string, user models.AppUser) error {
	var run models.CsmRun
	ctx := context.Background()
	if err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
		return d.First(&run, runID).Error
	}); err != nil {
		return err
	}
	if run.ProcessingStatus != "completed" {
		return errors.New("only completed runs can be reviewed")
	}
	if run.ApprovalStatus != "draft" && run.ApprovalStatus != "" {
		return errors.New("run must be in draft status to review")
	}
	now := time.Now()
	run.ApprovalStatus = "reviewed"
	run.ReviewedBy = user.UserName
	run.ReviewedAt = &now
	run.ReviewNotes = notes
	if err := DB.Save(&run).Error; err != nil {
		return err
	}
	LogIFRS17Event("run_reviewed", "csm_run", run.Name, run.ID, user, notes)
	go NotifyCsmRunReviewed(run, user)
	return nil
}

// ApproveCsmRun transitions a reviewed run to approved.
func ApproveCsmRun(runID int, notes string, user models.AppUser) error {
	var run models.CsmRun
	ctx := context.Background()
	if err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
		return d.First(&run, runID).Error
	}); err != nil {
		return err
	}
	if run.ApprovalStatus != "reviewed" {
		return errors.New("run must be in reviewed status to approve")
	}
	now := time.Now()
	run.ApprovalStatus = "approved"
	run.ApprovedBy = user.UserName
	run.ApprovedAt = &now
	run.ApproveNotes = notes
	if err := DB.Save(&run).Error; err != nil {
		return err
	}
	LogIFRS17Event("run_approved", "csm_run", run.Name, run.ID, user, notes)
	go NotifyCsmRunApproved(run, user)
	return nil
}

// LockCsmRun transitions an approved run to locked (immutable).
func LockCsmRun(runID int, user models.AppUser) error {
	var run models.CsmRun
	ctx := context.Background()
	if err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
		return d.First(&run, runID).Error
	}); err != nil {
		return err
	}
	if run.ApprovalStatus != "approved" {
		return errors.New("run must be approved before it can be locked")
	}
	now := time.Now()
	run.ApprovalStatus = "locked"
	run.LockedBy = user.UserName
	run.LockedAt = &now
	if err := DB.Save(&run).Error; err != nil {
		return err
	}
	LogIFRS17Event("run_locked", "csm_run", run.Name, run.ID, user, "")
	return nil
}

// ReturnCsmRunToDraft returns a reviewed or approved run back to draft.
func ReturnCsmRunToDraft(runID int, reason string, user models.AppUser) error {
	var run models.CsmRun
	ctx := context.Background()
	if err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
		return d.First(&run, runID).Error
	}); err != nil {
		return err
	}
	if run.ApprovalStatus == "locked" {
		return errors.New("locked runs cannot be returned to draft")
	}
	run.ApprovalStatus = "draft"
	run.ReturnReason = reason
	run.ReviewedBy = ""
	run.ReviewedAt = nil
	run.ApprovedBy = ""
	run.ApprovedAt = nil
	if err := DB.Save(&run).Error; err != nil {
		return err
	}
	LogIFRS17Event("run_returned_to_draft", "csm_run", run.Name, run.ID, user, reason)
	return nil
}

// CompareRuns produces a side-by-side metric comparison between two CsmRun records.
func CompareRuns(runAId, runBId int) (map[string]interface{}, error) {
	ctx := context.Background()

	var runA, runB models.CsmRun
	if err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
		return d.First(&runA, runAId).Error
	}); err != nil {
		return nil, fmt.Errorf("run A not found: %w", err)
	}
	if err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
		return d.First(&runB, runBId).Error
	}); err != nil {
		return nil, fmt.Errorf("run B not found: %w", err)
	}

	buildMetric := func(label string, aVal, bVal float64) models.RunComparisonMetric {
		variance := bVal - aVal
		var variancePct float64
		if aVal != 0 {
			variancePct = (variance / aVal) * 100
		}
		return models.RunComparisonMetric{
			Label:       label,
			RunAValue:   aVal,
			RunBValue:   bVal,
			Variance:    variance,
			VariancePct: variancePct,
		}
	}

	var metrics []models.RunComparisonMetric

	isPAA := func(r models.CsmRun) bool {
		return strings.EqualFold(r.MeasurementType, "PAA")
	}

	if !isPAA(runA) && !isPAA(runB) {
		// GMM / VFA aggregate from aos_step_results
		type aosAgg struct {
			BEL                  float64
			RiskAdjustment       float64
			CSMBuildup           float64
			LossComponentBuildup float64
			ReinsuranceBel       float64
			ReinsuranceCSM       float64
		}
		var aggA, aggB aosAgg

		if err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
			return d.Table("aos_step_results").
				Select("SUM(bel) as bel, SUM(risk_adjustment) as risk_adjustment, SUM(csm_buildup) as csm_buildup, SUM(loss_component_buildup) as loss_component_buildup, SUM(reinsurance_bel) as reinsurance_bel, SUM(reinsurance_csm) as reinsurance_csm").
				Where("csm_run_id = ?", runAId).
				Scan(&aggA).Error
		}); err != nil {
			return nil, fmt.Errorf("aggregating AOS results for run A: %w", err)
		}
		if err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
			return d.Table("aos_step_results").
				Select("SUM(bel) as bel, SUM(risk_adjustment) as risk_adjustment, SUM(csm_buildup) as csm_buildup, SUM(loss_component_buildup) as loss_component_buildup, SUM(reinsurance_bel) as reinsurance_bel, SUM(reinsurance_csm) as reinsurance_csm").
				Where("csm_run_id = ?", runBId).
				Scan(&aggB).Error
		}); err != nil {
			return nil, fmt.Errorf("aggregating AOS results for run B: %w", err)
		}

		// Insurance revenue from insurance_revenues table
		var irA, irB struct{ Total float64 }
		if err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
			return d.Table("insurance_revenues").
				Select("SUM(post_transition) as total").
				Where("csm_run_id = ?", runAId).
				Scan(&irA).Error
		}); err != nil {
			return nil, fmt.Errorf("aggregating insurance revenue for run A: %w", err)
		}
		if err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
			return d.Table("insurance_revenues").
				Select("SUM(post_transition) as total").
				Where("csm_run_id = ?", runBId).
				Scan(&irB).Error
		}); err != nil {
			return nil, fmt.Errorf("aggregating insurance revenue for run B: %w", err)
		}

		metrics = []models.RunComparisonMetric{
			buildMetric("BEL", aggA.BEL, aggB.BEL),
			buildMetric("Risk Adjustment", aggA.RiskAdjustment, aggB.RiskAdjustment),
			buildMetric("CSM Buildup", aggA.CSMBuildup, aggB.CSMBuildup),
			buildMetric("Loss Component Buildup", aggA.LossComponentBuildup, aggB.LossComponentBuildup),
			buildMetric("Reinsurance BEL", aggA.ReinsuranceBel, aggB.ReinsuranceBel),
			buildMetric("Reinsurance CSM", aggA.ReinsuranceCSM, aggB.ReinsuranceCSM),
			buildMetric("Insurance Revenue", irA.Total, irB.Total),
		}
	} else {
		// PAA aggregate from paa_results
		type paaAgg struct {
			EarnedPremium                 float64
			IncurredClaims                float64
			IncurredExpenses              float64
			PaaLiabilityRemainingCoverage float64
			InsuranceRevenue              float64
			PaaLossComponent              float64
			PaaReinsuranceLrc             float64
		}
		var aggA, aggB paaAgg

		if err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
			return d.Table("paa_results").
				Select("SUM(earned_premium) as earned_premium, SUM(incurred_claims) as incurred_claims, SUM(incurred_expenses) as incurred_expenses, SUM(paa_liability_remaining_coverage) as paa_liability_remaining_coverage, SUM(insurance_revenue) as insurance_revenue, SUM(paa_loss_component) as paa_loss_component, SUM(paa_reinsurance_lrc) as paa_reinsurance_lrc").
				Where("csm_run_id = ?", runAId).
				Scan(&aggA).Error
		}); err != nil {
			return nil, fmt.Errorf("aggregating PAA results for run A: %w", err)
		}
		if err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
			return d.Table("paa_results").
				Select("SUM(earned_premium) as earned_premium, SUM(incurred_claims) as incurred_claims, SUM(incurred_expenses) as incurred_expenses, SUM(paa_liability_remaining_coverage) as paa_liability_remaining_coverage, SUM(insurance_revenue) as insurance_revenue, SUM(paa_loss_component) as paa_loss_component, SUM(paa_reinsurance_lrc) as paa_reinsurance_lrc").
				Where("csm_run_id = ?", runBId).
				Scan(&aggB).Error
		}); err != nil {
			return nil, fmt.Errorf("aggregating PAA results for run B: %w", err)
		}

		metrics = []models.RunComparisonMetric{
			buildMetric("Earned Premium", aggA.EarnedPremium, aggB.EarnedPremium),
			buildMetric("Incurred Claims", aggA.IncurredClaims, aggB.IncurredClaims),
			buildMetric("Incurred Expenses", aggA.IncurredExpenses, aggB.IncurredExpenses),
			buildMetric("PAA Liability Remaining Coverage", aggA.PaaLiabilityRemainingCoverage, aggB.PaaLiabilityRemainingCoverage),
			buildMetric("Insurance Revenue", aggA.InsuranceRevenue, aggB.InsuranceRevenue),
			buildMetric("PAA Loss Component", aggA.PaaLossComponent, aggB.PaaLossComponent),
			buildMetric("PAA Reinsurance LRC", aggA.PaaReinsuranceLrc, aggB.PaaReinsuranceLrc),
		}
	}

	return map[string]interface{}{
		"run_a":   runA,
		"run_b":   runB,
		"metrics": metrics,
	}, nil
}

// GetConsolidatedDisclosure aggregates key IFRS 17 metrics across all approved/locked runs.
func GetConsolidatedDisclosure() (map[string]interface{}, error) {
	ctx := context.Background()

	type PortfolioRow struct {
		ProductCode           string  `json:"product_code"`
		MeasurementType       string  `json:"measurement_type"`
		TotalCSM              float64 `json:"total_csm"`
		TotalBEL              float64 `json:"total_bel"`
		TotalRA               float64 `json:"total_ra"`
		TotalLRC              float64 `json:"total_lrc"`
		TotalInsuranceRevenue float64 `json:"total_insurance_revenue"`
		OnerousContracts      int     `json:"onerous_contracts"`
		RunDate               string  `json:"run_date"`
		RunName               string  `json:"run_name"`
	}

	// Fetch all approved/locked runs
	var runs []models.CsmRun
	err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
		return d.Where("approval_status IN ?", []string{"approved", "locked"}).
			Order("run_date DESC").Find(&runs).Error
	})
	if err != nil {
		return nil, err
	}

	// For each run, fetch step results and aggregate by product+run_date
	portfolioMap := map[string]*PortfolioRow{}
	for _, run := range runs {
		var steps []models.AOSStepResult
		err2 := DBReadWithResilience(ctx, func(d *gorm.DB) error {
			return d.Where("csm_run_id = ?", run.ID).Find(&steps).Error
		})
		if err2 != nil {
			continue
		}
		for _, s := range steps {
			key := s.ProductCode + "|" + run.RunDate
			if _, ok := portfolioMap[key]; !ok {
				portfolioMap[key] = &PortfolioRow{
					ProductCode:     s.ProductCode,
					MeasurementType: run.MeasurementType,
					RunDate:         run.RunDate,
					RunName:         run.Name,
				}
			}
			row := portfolioMap[key]
			row.TotalCSM += s.CSMBuildup
			row.TotalBEL += s.BEL
			row.TotalRA += s.RiskAdjustment
			if s.IsOnerous {
				row.OnerousContracts++
			}
		}
	}

	portfolios := make([]PortfolioRow, 0, len(portfolioMap))
	var totalCSM, totalBEL, totalRA float64
	var totalOnerous int
	seenProducts := map[string]bool{}
	for _, row := range portfolioMap {
		row.TotalLRC = row.TotalCSM + row.TotalBEL + row.TotalRA
		portfolios = append(portfolios, *row)
		if !seenProducts[row.ProductCode] {
			seenProducts[row.ProductCode] = true
			totalCSM += row.TotalCSM
			totalBEL += row.TotalBEL
			totalRA += row.TotalRA
			totalOnerous += row.OnerousContracts
		}
	}

	return map[string]interface{}{
		"portfolios":      portfolios,
		"total_csm":       totalCSM,
		"total_bel":       totalBEL,
		"total_ra":        totalRA,
		"total_lrc":       totalCSM + totalBEL + totalRA,
		"total_onerous":   totalOnerous,
		"portfolio_count": len(seenProducts),
	}, nil
}
