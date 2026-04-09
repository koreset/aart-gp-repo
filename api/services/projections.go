package services

import (
	"api/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/copier"
	"github.com/rs/zerolog/log"
)

// var aggPeriod int

//func SaveProjectionJob(job models.ProjectionJob) (models.ProjectionJob, error) {
//	err := DB.Save(&job).Error
//	return job, err
//}

func SaveProjectionJob(job models.ProjectionJob) (models.ProjectionJob, error) {
	// Serialize the ProdIds slice to JSON before saving
	if len(job.ProdIds) > 0 {
		prodIdsJSON, err := json.Marshal(job.ProdIds)
		if err != nil {
			return job, err
		}
		job.ProdIdsJson = string(prodIdsJSON)
	}

	err := DB.Save(&job).Error
	return job, err
}

func GetValidJobs() []models.ProjectionJob {
	var jobs []models.ProjectionJob
	DB.Preload("RunParameters").Where("run_type = ? and status = ?", 0, "Complete").Order("id desc").Find(&jobs)
	return jobs
}

func GetProductsForJob(jobId int) []models.JobProduct {
	var products []models.JobProduct
	DB.Where("projection_job_id = ?", jobId).Find(&products)
	return products
}

func GetSpCodesForRunProduct(jobId int, productCode string) []string {
	var spCodes []string
	DB.Table("aggregated_projections").Where("run_id = ? AND product_code = ?", jobId, productCode).Distinct("sp_code").Pluck("sp_code", &spCodes)
	return spCodes
}

func GetSpCodesForRun(runId int) []string {
	var spCodes []string
	DB.Table("aggregated_projections").Where("run_id = ?", runId).Distinct("sp_code").Pluck("sp_code", &spCodes)
	return spCodes
}

func GetJobs(runType int) []models.ProjectionJob {
	var jobs []models.ProjectionJob
	start := time.Now()
	DB.Preload("RunParameters").Preload("Products").Where("run_type = ?", runType).Order("id desc").Find(&jobs)

	for i, job := range jobs {
		if job.Status == "In Progress" {
			jobs[i].PointsDone = 0
			for _, jp := range job.Products {
				jobs[i].PointsDone += jp.PointsDone
			}
		}
	}

	end := time.Since(start)
	fmt.Println("Query time: ", end)
	return jobs
}

func DeleteProjection(projectionJobId int) {
	var jps models.JobProduct
	DB.Where("projection_job_id = ?", projectionJobId).Find(&jps)
	DB.Model(&models.BaseAosVariable{}).Where("run_id = ?", projectionJobId).Update("run_id", 0)
	DB.Where("job_product_id = ?", jps.ID).Delete(&models.AggregatedProjection{})
	DB.Where("job_product_id = ?", jps.ID).Delete(&models.ScopedAggregatedProjection{})

	err := DB.Where("projection_job_id = ?", projectionJobId).Delete(&jps).Error
	if err != nil {
		fmt.Println(err)
	}

	err = DB.Where("projection_job_id = ?", projectionJobId).Delete(&models.RunParameters{}).Error
	if err != nil {
		fmt.Println(err)
	}

	err = DB.Where("projection_job_id = ?", projectionJobId).Delete(&models.JobProduct{}).Error
	if err != nil {
		fmt.Println(err)
	}

	err = DB.Where("id = ?", projectionJobId).Delete(&models.ProjectionJob{}).Error
	if err != nil {
		fmt.Println(err)
	}

	err = DB.Where("projection_job_id = ?", projectionJobId).Delete(&models.JobProductRunError{}).Error
	if err != nil {
		fmt.Println(err)
	}
}

func GetMostRecentJob() models.ProjectionJob {
	var job models.ProjectionJob
	DB.Preload("Products").Order("id desc").First(&job)
	return job
}

type ProjectionJobList struct {
	RunJob        models.RunPayload
	RunParameters models.RunParameters
	ProjectionJob models.ProjectionJob
}

func RunValuations(runJobs models.RunJob) {
	//Looks like we first setup the runs so that they show and then execute them sequentially
	var joblist []ProjectionJobList
	for _, run := range runJobs.Jobs {
		run.UserName = runJobs.User.UserName
		run.UserEmail = runJobs.User.UserEmail

		var job = models.ProjectionJob{}
		var runParams models.RunParameters

		job.CreationDate = time.Now()
		job.RunName = run.RunName
		job.UserEmail = runJobs.User.UserEmail
		job.UserName = runJobs.User.UserName
		job.YieldCurveBasis = run.YieldCurveBasis
		job.YieldCurveMonth = run.YieldcurveMonth
		job.RunBasis = run.RunBasis
		job.ProdIds = run.ProdIds
		job.AggregationPeriod = run.AggregationPeriod

		job.RunDescription = run.RunDescription
		job.JobsTemplateID = run.JobsTemplateID
		job.Products = []models.JobProduct{}
		if run.RunType == 1 {
			job.RunType = 1
		} else {
			job.RunType = 0
		}
		valDate := run.RunDate
		if valDate != "" {
			run.RunDate = valDate
			job.RunDate = valDate
		}
		job.Status = "Queued"
		job.ShockSettingID = run.ShockSettings.ID

		for _, v := range job.ProdIds {
			prod, _ := GetProductById(v)
			tableName := strings.ToLower(prod.ProductCode) + "_modelpoints"
			var count int64
			DB.Table(tableName).Where("year =? and mp_version = ?", run.ModelpointYear, run.ModelPointVersion).Count(&count)
			if run.RunSingle {
				job.TotalPoints += MpLimitValue
			} else {
				job.TotalPoints += int(count)
			}
		}

		job, err := SaveProjectionJob(job)
		if err != nil {
			log.Error().Err(err)
		}

		copier.Copy(&runParams, run)
		runParams.ProjectionJobID = job.ID
		runParams.ShockSettingsID = run.ShockSettings.ID
		DB.Save(&runParams)
		var tmpJobList ProjectionJobList
		tmpJobList.RunJob = run
		tmpJobList.RunParameters = runParams
		tmpJobList.ProjectionJob = job
		joblist = append(joblist, tmpJobList)
	}

	// we return here after queuing the jobs
	// This is where we would normally kick off the jobs to run

	//fmt.Println(joblist)
	////Keep polling here until there are no jobs in the queue
	//for {
	//	var count int64
	//	DB.Model(&models.ProjectionJob{}).Where("status = ?", "In Progress").Count(&count)
	//	fmt.Println("polling for jobs to finish")
	//	if count == 0 {
	//		fmt.Println("no jobs in progress")
	//		break
	//	}
	//	time.Sleep(3 * time.Second)
	//}
	//
	//for _, item := range joblist {
	//	for _, id := range item.ProjectionJob.ProdIds {
	//		item.ProjectionJob.Status = "In Progress"
	//		DB.Save(&item.ProjectionJob)
	//		err := PopulateProjectionsForProduct(id, &item.ProjectionJob)
	//		if err != nil {
	//			log.Error().Err(err)
	//			item.ProjectionJob.Status = "Failed"
	//			item.ProjectionJob.StatusError = err.Error()
	//		}
	//		DB.Save(&item.ProjectionJob)
	//	}
	//	Cache.Clear()
	//	PaaCache.Clear()
	//}

	//At this point, clear the cache?

}

func GetAllAggregatedReserves(jobProductIds []interface{}, resultRange int, runVariable string) models.AllAggregatedReserves {
	var allAggregatedReserves models.AllAggregatedReserves
	var err error
	//var resultMap = make([]map[int]interface{}, 0)
	for _, jobProductId := range jobProductIds {
		var cachedReserves models.CachedReserveResults
		var reserves models.AggregatedReserves
		var jpId int
		switch reflect.TypeOf(jobProductId).Kind() {
		case reflect.Float64:
			jpId = int(jobProductId.(float64))
		case reflect.String:
			jpId, err = strconv.Atoi(jobProductId.(string))
		case reflect.Int:
			jpId = jobProductId.(int)
		}

		DB.Where("job_product_id = ? and result_range = ? and variable = ?", jpId, resultRange, runVariable).First(&cachedReserves)
		if cachedReserves.JobProductId == jpId {
			reservesJson := cachedReserves.Result
			json.Unmarshal(reservesJson, &reserves)
		} else {
			reserves = GetAggregatedReserves(jpId, resultRange, runVariable)
			aggJson, _ := json.Marshal(reserves)
			cachedReserves.JobProductId = jpId
			cachedReserves.ResultRange = resultRange
			cachedReserves.Result = aggJson
			cachedReserves.Variable = runVariable
			DB.Save(&cachedReserves)

		}
		if len(reserves.AggregatedReserves) > 0 {
			allAggregatedReserves.AllAggregatedReserves = append(allAggregatedReserves.AllAggregatedReserves, reserves)
		}

		if err != nil {
			fmt.Println(err)
		}
	}
	return allAggregatedReserves
}

func GetAggregatedReserves(jobProductId int, resultRange int, runVariable string) models.AggregatedReserves {
	var reserves = models.AggregatedReserves{}
	var jobProduct models.JobProduct
	DB.Where("id=?", jobProductId).First(&jobProduct)
	reserves.ProjectionJobProduct = jobProduct
	var rows *sql.Rows
	var err error
	selectQuery := fmt.Sprintf("projection_month,SUM(%s) as reserves", runVariable)

	spanStart := time.Now()
	//revert this later
	resultRange = 0

	if resultRange > 0 {
		rows, err = DB.Table("aggregated_projections").Where("job_product_id=? and projection_month < ?", jobProductId, resultRange).Select(selectQuery).Group("projection_month").Rows()
	} else {
		rows, err = DB.Table("aggregated_projections").Where("job_product_id=?", jobProductId).Select(selectQuery).Group("projection_month").Order("projection_month asc").Rows()
	}
	if err != nil || rows == nil {
		if err != nil {
			fmt.Println(err)
		}
		return reserves
	}
	defer rows.Close()
	duration := time.Since(spanStart)
	fmt.Println("Query time: ", duration)
	spanStart = time.Now()
	for rows.Next() {
		//var projRow models.Projection
		var agg = models.AggregatedReserve{}
		_ = DB.ScanRows(rows, &agg)
		//_ = DB.ScanRows(rows, &projRow)
		//_ = copier.Copy(&agg, &projRow)
		reserves.AggregatedReserves = append(reserves.AggregatedReserves, agg)
	}
	duration = time.Since(spanStart)
	fmt.Println("Append time : ", duration)
	return reserves
}
func GetProjections(jobProductId int) []models.Projection {
	var projections []models.Projection
	dbErr := DB.Where("job_product_id=?", jobProductId).Find(&projections).Error
	if dbErr != nil {
		fmt.Println(dbErr)
	}

	return projections

}

func GetSpCodes(jobProductId int) ([]int, error) {
	var spCodes []int

	err := DB.Table("aggregated_projections").Distinct("sp_code").Where("job_product_id = ?", jobProductId).Pluck("sp_code", &spCodes).Error

	return spCodes, err
}

func RestartStalledJob(projectionJobId int) error {
	// Load the specific job efficiently
	var job models.ProjectionJob
	if err := DB.Where("id = ?", projectionJobId).First(&job).Error; err != nil {
		log.Info().Msgf("No stalled job found with ID: %d (err=%v)", projectionJobId, err)
		return err
	}

	log.Warn().Msgf("Recovering stalled job ID: %d, Name: %s. Marking as queued and clearing partial results.", job.ID, job.RunName)

	var jobProducts []models.JobProduct
	if err := DB.Where("projection_job_id = ?", job.ID).Find(&jobProducts).Error; err != nil {
	}

	go func() {
		for _, jobProduct := range jobProducts {
			if err := DB.Where("run_id = ? and job_product_id = ?", job.ID, jobProduct.ID).Delete(&models.AggregatedProjection{}).Error; err != nil {
				log.Info().Msgf("Failed to delete stalled job results for job ID: %d, Product ID: %d", job.ID, jobProduct.ID)

			}
			if err := DB.Where("job_product_id = ? ", jobProduct.ID).Delete(&models.ScopedAggregatedProjection{}).Error; err != nil {
				log.Info().Msgf("Failed to delete stalled job results for job ID: %d, Product ID: %d", job.ID, jobProduct.ID)

			}
			if err := DB.Where("run_id = ? and job_product_id = ?", job.ID, jobProduct.ID).Delete(&models.Projection{}).Error; err != nil {
				log.Info().Msgf("Failed to delete stalled job results for job ID: %d, Product ID: %d", job.ID, jobProduct.ID)

			}
		}

	}()

	// Wrap cleanup + status reset in a single transaction for atomicity and fewer round-trips
	if err := DB.Where("projection_job_id = ?", job.ID).Delete(&models.JobProduct{}).Error; err != nil {
		return err
	}
	// If AggregatedReserves needs clearing for this job, enable the following line once the column criteria is confirmed
	// _ = tx.Where("projection_job_id = ?", job.ID).Delete(&models.AggregatedReserves{}).Error

	job.Status = "Queued"
	job.StatusError = "Reprocessing Job. Marked as Queued for reprocessing."
	job.RunTime = 0
	job.PointsDone = 0
	job.CreationDate = time.Now()
	if err := DB.Save(&job).Error; err != nil {
		log.Error().Err(err).Msgf("Failed to update stalled job status for job ID: %d", job.ID)
		return err
	}

	return nil
}

func GetAggregatedProjections(jobProductId int, spCode string) []models.AggregatedProjection {
	var projections []models.AggregatedProjection
	var err error

	dbErr := DB.Where("job_product_id=? and sp_code = ?", jobProductId, spCode).Order("sp_code asc").Order("projection_month asc").Find(&projections).Error
	if dbErr != nil {
		fmt.Println(err)
	}

	return projections
}

func GetExcelAggregatedProjections(jobId, jobProductId int, spCode string, variableGroupId int) []map[string]interface{} {
	var projections []models.AggregatedProjection
	//var err error

	// get the variable group
	var variableGroup models.AggregatedVariableGroup
	DB.Where("id = ?", variableGroupId).First(&variableGroup)

	if jobProductId == 0 && spCode == "" {
		// user has chosen all products for results
		// do results for all products
		// first we get all agg variables and the query results aggregation code
		// var aggVars []models.AggregatedVariableGroup
		var variables []string
		if variableGroup.ID == 0 {
			variables = GetAggregationVariables()
		} else {
			variables = variableGroup.Variables
		}
		aggregations, err := GetAggregations(jobId, "", spCode, variables)
		if err != nil {
			fmt.Println(err)
		}
		jsonData, err := json.Marshal(aggregations)
		err = json.Unmarshal(jsonData, &projections)
		return aggregations
	}

	if jobProductId != 0 && spCode == "" {
		// user has chosen a product for results
		// do results for a single product
		// first we get all agg variables and the query results aggregation code
		// var aggVars []models.AggregatedVariableGroup
		var jobProduct models.JobProduct
		DB.Where("id = ?", jobProductId).First(&jobProduct)
		var variables []string
		if variableGroup.ID == 0 {
			variables = GetAggregationVariables()
		} else {
			variables = variableGroup.Variables
		}
		aggregations, err := GetAggregations(jobId, jobProduct.ProductCode, spCode, variables)
		if err != nil {
			fmt.Println(err)
		}
		jsonData, err := json.Marshal(aggregations)
		err = json.Unmarshal(jsonData, &projections)
		return aggregations
	}

	if jobProductId != 0 && spCode != "" {
		// user has chosen a product and a sp code for results
		// do results for a single product
		// first we get all agg variables and the query results aggregation code
		// var aggVars []models.AggregatedVariableGroup
		var jobProduct models.JobProduct
		DB.Where("id = ?", jobProductId).First(&jobProduct)
		var variables []string
		if variableGroup.ID == 0 {
			variables = GetAggregationVariables()
		} else {
			variables = variableGroup.Variables
		}
		if spCode == "All" {
			spCode = ""
		}

		aggregations, err := GetAggregations(jobId, jobProduct.ProductCode, spCode, variables)
		if err != nil {
			fmt.Println(err)
		}
		jsonData, err := json.Marshal(aggregations)
		err = json.Unmarshal(jsonData, &projections)
		return aggregations
	}

	if jobProductId == 0 && spCode == "All" {
		// user has chosen a product and a sp code for results
		// do results for a single product
		// first we get all agg variables and the query results aggregation code
		// var aggVars []models.AggregatedVariableGroup
		var jobProduct models.JobProduct
		DB.Where("id = ?", jobProductId).First(&jobProduct)
		var variables []string
		if variableGroup.ID == 0 {
			variables = GetAggregationVariables()
		} else {
			variables = variableGroup.Variables
		}
		spCode = ""
		aggregations, err := GetAggregations(jobId, "", spCode, variables)
		if err != nil {
			fmt.Println(err)
		}
		jsonData, err := json.Marshal(aggregations)
		err = json.Unmarshal(jsonData, &projections)
		return aggregations
	}

	if jobProductId == 0 && spCode != "" && spCode != "All" {
		// user has chosen a product and a sp code for results
		// do results for a single product
		// first we get all agg variables and the query results aggregation code
		// var aggVars []models.AggregatedVariableGroup
		var jobProduct models.JobProduct
		DB.Where("id = ?", jobProductId).First(&jobProduct)
		var variables []string
		if variableGroup.ID == 0 {
			variables = GetAggregationVariables()
		} else {
			variables = variableGroup.Variables
		}

		aggregations, err := GetAggregations(jobId, "", spCode, variables)
		if err != nil {
			fmt.Println(err)
		}
		jsonData, err := json.Marshal(aggregations)
		err = json.Unmarshal(jsonData, &projections)
		return aggregations
	}
	return nil

	//dbErr := DB.Where("job_product_id=? and sp_code = ?", jobProductId, spCode).Order("sp_code asc").Order("projection_month asc").Find(&projections).Error
	//if dbErr != nil {
	//	fmt.Println(err)
	//}
	//
	//return projections
}

// getAggregatedProjectionExcludedFields returns a list of field names to exclude from SQL queries
// based on the JSON struct tags of the AggregatedProjection struct.
// Fields with JSON tag "-" or empty JSON tag will be included in the result.
func getAggregatedProjectionExcludedFields() []string {
	var excludedFields []string
	t := reflect.TypeOf(models.AggregatedProjection{})

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")

		// If the JSON tag is "-" or empty, add the field name to the excluded list
		if jsonTag == "-" || jsonTag == "" {
			// Convert field name to snake_case for SQL column name
			fieldName := field.Name
			snakeCase := ""
			for j, c := range fieldName {
				if j > 0 && c >= 'A' && c <= 'Z' {
					snakeCase += "_"
				}
				snakeCase += strings.ToLower(string(c))
			}
			excludedFields = append(excludedFields, snakeCase)
		}
	}

	return excludedFields
}

func GetProjectionJobExcel(projectionJobId int) ([]byte, error) {
	var err error
	dQuery := fmt.Sprintf("SELECT run_name,product_code,sp_code,ifrs17_group,projection_month,initial_policy,premium,sum_assured,discounted_premium_income,discounted_commission,discounted_renewal_commission,discounted_claw_back,discounted_death_outgo,discounted_accidental_death_outgo,discounted_investment_income,discounted_initial_expenses,discounted_renewal_expenses,discounted_cash_outflow,discounted_cash_outflow_excl_acquisition,discounted_profit,reserves,vif FROM aggregated_projections where run_id=%d order by product_code,sp_code, projection_month ", projectionJobId)
	excelData, err := exportTableToExcel(dQuery)
	if err != nil {
		return nil, nil
	}

	return excelData, nil
}

func GetProjectionJobExcelScoped(projectionJobId int) ([]byte, error) {
	var err error
	dQuery := fmt.Sprintf("SELECT run_name,product_code,ifrs17_group,projection_month,initial_policy,premium,sum_assured,premium_income,premium_not_received_lapse,commission,renewal_commission,claw_back,non_life_claims_outgo,death_outgo,accidental_death_outgo,cash_back_on_survival,cash_back_on_death,disability_outgo,retrenchment_outgo,rider,initial_expenses,renewal_expenses,discounted_premium_income,discounted_commission,discounted_renewal_commission,discounted_claw_back,discounted_death_outgo,discounted_accidental_death_outgo,discounted_investment_income,discounted_initial_expenses,discounted_renewal_expenses,discounted_cash_inflow,discounted_cash_outflow,discounted_cash_outflow_excl_acquisition,discounted_profit,reserves,vif,sum_coverage_units,discounted_coverage_units FROM scoped_aggregated_projections where run_id=%d order by product_code,ifrs17_group, projection_month ", projectionJobId)
	excelData, err := exportTableToExcel(dQuery)
	if err != nil {
		return nil, nil
	}

	return excelData, nil
}

func GetAggregatedProjectionsForDownload(jobProductId int) ([]byte, error) {
	//var projections []models.AggregatedProjectionData
	var err error
	var selectItems string

	DB.Exec("SET SESSION group_concat_max_len = 1000000")
	dbName, _ := GetDBName(DB)

	// Get excluded fields based on JSON struct tags
	excludedFields := getAggregatedProjectionExcludedFields()

	// Convert the excluded fields slice to a comma-separated string for the SQL query
	excludedFieldsStr := "'" + strings.Join(excludedFields, "', '") + "'"

	// add job_product_id, run_id and id to the excluded fields
	excludedFieldsStr += ", 'job_product_id', 'run_id', 'id'"

	columnQuery := fmt.Sprintf("SELECT GROUP_CONCAT(column_name SEPARATOR ', ') FROM INFORMATION_SCHEMA.COLUMNS WHERE table_schema = '%s' and table_name = 'aggregated_projections' AND column_name NOT IN (%s)", dbName, excludedFieldsStr)
	DB.Raw(columnQuery).Scan(&selectItems)

	dQuery := fmt.Sprintf("select %s from aggregated_projections where job_product_id = %d order by sp_code, projection_month", selectItems, jobProductId)
	excelData, err := exportTableToExcel(dQuery)
	if err != nil {
		return nil, nil
	}

	return excelData, nil
}

// toTitleCase converts a snake_case string to Title Case
func toTitleCase(input string) string {
	words := strings.Split(input, "_")
	for i, word := range words {
		words[i] = strings.ToTitle(word)
	}
	return strings.Join(words, " ")
}

func GetGMMScopedProjections(jobProductId int) []models.ScopedAggregatedProjection {
	var sap []models.ScopedAggregatedProjection
	var err error
	dbErr := DB.Where("job_product_id=?", jobProductId).Order("ifrs17_group asc").Order("projection_month asc").Find(&sap).Error
	if dbErr != nil {
		fmt.Println(err)
	}
	return sap
}

func GetJobRunSettings(projectionJobId int) models.ProductRunParameter {
	//Need to get the projectionJobId from the job_product table
	var jobProduct models.JobProduct
	DB.Where("id=?", projectionJobId).First(&jobProduct)

	var result models.ProductRunParameter
	DB.Model(&models.RunParameters{}).Select("run_parameters.*, job_products.product_name, job_products.product_code").Joins("join job_products on job_products.projection_job_id = run_parameters.projection_job_id").Where("job_products.projection_job_id=?", jobProduct.ProjectionJobID).First(&result)
	var shockSettings models.ShockSetting
	DB.Where("id = ?", result.ShockSettingsID).First(&shockSettings)
	result.ShockSettingName = shockSettings.Name
	return result
}

func GetJobRunErrors(jobProductId int) []models.JobProductRunError {
	var jobProduct models.JobProduct
	DB.Where("id=?", jobProductId).First(&jobProduct)
	var runErrors []models.JobProductRunError

	DB.Where("projection_job_id=? and product_code = ?", jobProduct.ProjectionJobID, jobProduct.ProductCode).Distinct("failure_point", "error", "job_product_id", "projection_job_id", "product_code").Limit(25).Find(&runErrors)
	//DB.Where("projection_job_id=? and product_code = ?", jobProduct.ProjectionJobID, jobProduct.ProductCode).Find(&runErrors)
	return runErrors
}

func GetJobsForProduct(prodCode string) []models.ProjectionJob {
	var jobProducts []models.JobProduct
	DB.Where("product_code=?", prodCode).Find(&jobProducts)
	var jobs []models.ProjectionJob
	for _, jobProduct := range jobProducts {
		var job models.ProjectionJob
		DB.Preload("Products").Where("id=?", jobProduct.ProjectionJobID).Find(&job)
		jobs = append(jobs, job)
	}
	return jobs
}

func ValuationRunNameDoesNotExist(runName string) bool {
	var job models.ProjectionJob
	DB.Where("run_name = ?", runName).First(&job)
	if job.RunName == runName {
		return false
	} else {
		return true
	}
}

func GetAvailableYieldYears() []models.AvailableYieldResult {
	var years []models.AvailableYieldResult = []models.AvailableYieldResult{}
	DB.Raw("select distinct year from yield_curve").Scan(&years)
	return years
}

func GetYieldCurveCodes(year string) []string {
	var codes []string

	DB.Raw("select distinct yield_curve_code from yield_curve where year = ?", year).Scan(&codes)
	return codes
}

func GetIbnrYieldCurveCodes(year string) []string {
	var codes []string

	DB.Raw("select distinct yield_curve_code from ibnr_yield_curves where year = ?", year).Scan(&codes)
	return codes
}

func GetYieldCurveMonths(year string, code string) []string {
	var months []string
	DB.Raw("select distinct month from yield_curve where year = ? and yield_curve_code = ?", year, code).Scan(&months)
	return months
}

func GetYieldCurveMonthsv2(productCode, yieldYear, parameterYear, basis string) []string {
	var yieldCurveCode string
	DB.Raw("select distinct yield_curve_code from product_parameters where product_code = ? and year = ? and basis= ?", productCode, parameterYear, basis).Scan(&yieldCurveCode)

	var months []string
	tableName := "yield_curve"
	DB.Raw("select distinct month from "+tableName+" where year = ? and yield_curve_code = ?", yieldYear, yieldCurveCode).Scan(&months)
	return months
}

func DeleteYieldCurveData(year string, code string, month string) error {
	err := DB.Where("year = ? and yield_curve_code = ? and month = ?", year, code, month).Delete(&models.YieldCurve{}).Error
	return err
}

func DeleteIbnrYieldCurveData(year string, code string, month string) error {
	err := DB.Where("year = ? and yield_curve_code = ? and month = ?", year, code, month).Delete(&models.IbnrYieldCurve{}).Error
	return err
}

func GetAvailableParameterYears(productCode string) []models.AvailableYieldResult {
	var years []models.AvailableYieldResult = []models.AvailableYieldResult{}
	if productCode == "" {
		DB.Raw("select distinct year from product_parameters").Scan(&years)
	} else {
		DB.Raw("select distinct year from product_parameters where product_code = ?", productCode).Scan(&years)
	}
	return years
}

func GetAvailableMarginYears() []models.AvailableYieldResult {
	var years []models.AvailableYieldResult = []models.AvailableYieldResult{}
	DB.Raw("select distinct year from product_margins").Scan(&years)
	return years
}

func GetAvailableModelPointYears(productCode string) []models.AvailableYieldResult {
	var years []models.AvailableYieldResult = []models.AvailableYieldResult{}
	//tableName := strings.ToLower(productCode) + "_modelpoints"
	//tableName := "product_model_point_counts"

	//DB.Raw("select distinct year from "+tableName+" where product_code = ?", productCode).Scan(&years)
	//DB.Raw("select distinct year from " + tableName).Scan(&years)
	DB.Model(&models.ProductModelPointCount{}).Distinct("year").Where("product_code = ?", productCode).Scan(&years)
	return years
}

func GetAvailableLapseYears(productCode string) []models.AvailableYieldResult {
	var years []models.AvailableLapseResult
	var results []models.AvailableYieldResult = []models.AvailableYieldResult{}
	tableName := strings.ToLower(productCode + "_" + GetRatingTable(productCode, "Lapse"))
	if tableName == "" {
		return results
	}

	lapseColumnName := GetColumnName(tableName)
	DB.Raw("select distinct left(" + lapseColumnName + ", 4) as year from " + tableName).Scan(&years)
	if len(years) > 0 {
		for _, v := range years {
			var result models.AvailableYieldResult
			result.Year, _ = strconv.Atoi(v.Year)
			results = append(results, result)
		}
	}
	return results
}

func GetAvailableMortalityYears(productCode string) []models.AvailableYieldResult {
	var years []models.AvailableLapseResult
	var results []models.AvailableYieldResult = []models.AvailableYieldResult{}
	tableName := strings.ToLower(productCode + "_" + GetRatingTable(productCode, "Death"))
	//DB.Raw("select distinct left(year_age_gender, 4) as year from " + tableName).Scan(&years)
	if tableName == "" {
		return results
	}
	mortalityColumnName := GetColumnName(tableName)
	err := DB.Raw("select distinct left(" + mortalityColumnName + ", 4) as year from " + tableName).Scan(&years).Error
	if err != nil {
		fmt.Println(err)
	}
	if len(years) > 0 {
		for _, v := range years {
			var result models.AvailableYieldResult
			result.Year, _ = strconv.Atoi(v.Year)
			results = append(results, result)
		}
	}
	return results
}

func GetAvailableMortalityAccidentalYears(productCode string) []models.AvailableYieldResult {
	var years []models.AvailableLapseResult
	var results []models.AvailableYieldResult = []models.AvailableYieldResult{}
	tableName := strings.ToLower(productCode + "_" + GetRatingTable(productCode, "Accidental Death"))
	if tableName == "" {
		return results
	}
	//DB.Raw("select distinct left(year_age_gender, 4) as year from " + tableName).Scan(&years)
	mortalityColumnName := GetColumnName(tableName)
	err := DB.Raw("select distinct left(" + mortalityColumnName + ", 4) as year from " + tableName).Scan(&years).Error
	if err != nil {
		fmt.Println(err)
	}
	if len(years) > 0 {
		for _, v := range years {
			var result models.AvailableYieldResult
			result.Year, _ = strconv.Atoi(v.Year)
			results = append(results, result)
		}
	}
	return results
}

func GetAvailableRetrenchmentYears(productCode string) []models.AvailableYieldResult {
	// We need to revisit this for the column name at some stage.
	var years []models.AvailableLapseResult
	var results []models.AvailableYieldResult = []models.AvailableYieldResult{}
	tableName := strings.ToLower(productCode + "_" + GetRatingTable(productCode, "Retrenchment"))
	if tableName == "" {
		return results
	}

	retrenchmentColumnName := GetColumnName(tableName)
	DB.Raw("select distinct left(" + retrenchmentColumnName + ", 4) as year from " + tableName).Scan(&years)

	if len(years) > 0 {
		for _, v := range years {
			var result models.AvailableYieldResult
			result.Year, _ = strconv.Atoi(v.Year)
			results = append(results, result)
		}
	}
	return results

}

func GetAvailableDisabilityYears(productCode string) []models.AvailableYieldResult {
	// We need to revisit this for the column name at some stage.
	var years []models.AvailableLapseResult
	var results []models.AvailableYieldResult = []models.AvailableYieldResult{}
	tableName := strings.ToLower(productCode + "_" + GetRatingTable(productCode, "Permanent Disability"))
	if tableName == "" {
		return results
	}
	disabilityColumnName := GetColumnName(tableName)
	//DB.Raw("select distinct left(year_anb_sec_occ_class_gender, 4) as year from " + tableName).Scan(&years)
	DB.Raw("select distinct left(" + disabilityColumnName + ", 4) as year from " + tableName).Scan(&years)

	if len(years) > 0 {
		for _, v := range years {
			var result models.AvailableYieldResult
			result.Year, _ = strconv.Atoi(v.Year)
			results = append(results, result)
		}
	}
	return results
}

func GetAvailableLapseMarginYears(productCode string) []models.AvailableYieldResult {
	var years []models.AvailableLapseResult
	var results []models.AvailableYieldResult = []models.AvailableYieldResult{}
	tableName := "product_lapse_margins"
	DB.Raw("select distinct year as year from " + tableName + " where product_code = '" + productCode + "'").Scan(&years)
	if len(years) > 0 {
		for _, v := range years {
			var result models.AvailableYieldResult
			result.Year, _ = strconv.Atoi(v.Year)
			results = append(results, result)
		}
	}
	return results
}

func GetAvailableBases(productCode string) []models.AvailableBasis {
	var bases []models.AvailableBasis
	DB.Raw("select distinct basis as basis from product_parameters where product_code = ?", productCode).Scan(&bases)
	return bases
}

func GetGMMShockBases() []models.AvailableBasis {
	var bases []models.AvailableBasis
	DB.Raw("select distinct shock_basis as basis from product_shocks").Scan(&bases)
	for i, _ := range bases {
		bases[i].Name = bases[i].Basis
	}
	return bases

}

//func GetShockBases() []models.AvailableBasis {
//	var bases []models.AvailableBasis
//	DB.Raw("select distinct shock_basis as basis from product_shocks").Scan(&bases)
//	for i, _ := range bases {
//		bases[i].Name = bases[i].Basis
//	}
//	return bases
//}

func GetIBNRBases(portfolio, year string) []models.AvailableBasis {
	parameterYear, err := strconv.Atoi(year)
	if err != nil {
		fmt.Println(err)
	}
	var bases []models.AvailableBasis
	DB.Raw("select distinct basis as basis from lic_parameters where portfolio_name = ? and year = ?", portfolio, parameterYear).Scan(&bases)
	return bases
}

func GetIBNRShockBases() []models.AvailableBasis {
	var bases []models.AvailableBasis
	DB.Raw("select distinct shock_basis as basis from ibnr_shocks").Scan(&bases)
	return bases
}

func SaveShockSetting(shockSetting models.ShockSetting) (models.ShockSetting, error) {
	DB.Save(&shockSetting)
	return shockSetting, nil
}

func GetShockSettings() ([]models.ShockSetting, error) {
	var shockSettings []models.ShockSetting

	// Get distinct valid shock bases from product_shocks
	var validBases []string
	if err := DB.Model(&models.ProductShock{}).Distinct("shock_basis").Pluck("shock_basis", &validBases).Error; err != nil {
		return shockSettings, err
	}

	// If there are no valid bases, return an empty slice (no settings should be returned)
	if len(validBases) == 0 {
		return []models.ShockSetting{}, nil
	}

	// Return only shock settings whose ShockBasis exists in product_shocks
	err := DB.Where("shock_basis IN (?)", validBases).Find(&shockSettings).Error
	return shockSettings, err
}

func DeleteShockSetting(id int) error {

	err := DB.Where("id = ?", id).Delete(&models.PhiShockSetting{}).Error
	if err != nil {
		return err
	}
	return nil
}

func GetJobsTemplates() []models.JobsTemplate {
	var templates []models.JobsTemplate
	DB.Find(&templates)

	return templates
}

func GetJobsTemplate(templateId int) models.JobTemplateContent {
	var templateContent models.JobTemplateContent

	DB.Where("jobs_template_id = ?", templateId).Find(&templateContent)
	return templateContent
}

func DeleteJobsTemplate(templateId int) models.JobTemplateContent {
	var templateContent models.JobTemplateContent

	err := DB.Where("id = ?", templateId).Delete(&models.JobsTemplate{}).Error
	if err != nil {
		fmt.Println(err)
	}
	err = DB.Where("jobs_template_id = ?", templateId).Delete(&templateContent).Error
	if err != nil {
		fmt.Println(err)
	}
	return templateContent
}

func SaveJobsTemplate(template models.JobsTemplate) models.JobsTemplate {
	err := DB.Save(&template).Error
	if err != nil {
		fmt.Println(err)
	}
	return template
}

func SaveJobTemplateContent(body string, id int) error {
	var content models.JobTemplateContent
	content.JobsTemplateID = id
	content.Content = body
	return DB.Where("jobs_template_id =?", id).Save(&content).Error
}

func SaveVariableGroup(variableGroup models.AggregatedVariableGroup) (models.AggregatedVariableGroup, error) {
	err := DB.Save(&variableGroup).Error
	return variableGroup, err
}

func GetVariableGroups() []models.AggregatedVariableGroup {
	var variableGroups []models.AggregatedVariableGroup
	DB.Find(&variableGroups)
	return variableGroups
}

func DeleteVariableGroup(id int) error {
	return DB.Delete(&models.AggregatedVariableGroup{}, id).Error
}

// projectionJobWorkerRunning is used to ensure only one worker is running
var projectionJobWorkerRunning = false
var projectionJobWorkerMutex sync.Mutex

// RecoverStalledJobs checks for and recovers any stalled projection jobs
// This should be called once during application startup
// It finds any jobs with "In Progress" status that may have been interrupted by a server shutdown
// and marks them as "Failed" so they can be requeued if needed
func RecoverStalledJobs() {
	var stalledJobs []models.ProjectionJob

	err := DB.Where("status = ?", "In Progress").Find(&stalledJobs).Error
	if err == nil && len(stalledJobs) > 0 {
		log.Info().Msgf("Found %d stalled jobs during server startup. Recovering...", len(stalledJobs))

		for _, stalledJob := range stalledJobs {
			log.Warn().Msgf("Recovering stalled job ID: %d, Name: %s. Marking as queued.", stalledJob.ID, stalledJob.RunName)
			// somewhere here, we need to find any results that were generated and clear them
			// This could be done by checking the aggregated_projections table for this job ID
			// and deleting any rows that match this job ID
			// Clear any results for this job
			DB.Where("projection_job_id = ?", stalledJob.ID).Delete(&models.AggregatedProjection{})
			DB.Where("projection_job_id = ?", stalledJob.ID).Delete(&models.ScopedAggregatedProjection{})
			DB.Where("projection_job_id = ?", stalledJob.ID).Delete(&models.Projection{})
			DB.Where("projection_job_id = ?", stalledJob.ID).Delete(&models.JobProduct{})
			DB.Where("projection_job_id = ?", stalledJob.ID).Delete(&models.AggregatedReserves{})

			stalledJob.Status = "Queued"
			stalledJob.StatusError = "Reprocessing interrupted Job. Marked as Queued for reprocessing."
			_, err = SaveProjectionJob(stalledJob)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to update stalled job status for job ID: %d", stalledJob.ID)
			}
		}
	}
}

// StartProjectionJobWorker starts the background worker that processes projection jobs
// It should be called once during application startup
func StartProjectionJobWorker() {
	projectionJobWorkerMutex.Lock()
	defer projectionJobWorkerMutex.Unlock()

	// Recover any stalled jobs before starting the worker
	//RecoverStalledJobs()

	if projectionJobWorkerRunning {
		log.Info().Msg("Projection job worker is already running")
		return
	}

	projectionJobWorkerRunning = true
	go func() {
		log.Info().Msg("Starting projection job worker")
		for {
			ProcessQueuedProjectionJobs()
			// Sleep for a short time before checking again and randomize the sleep duration
			// to avoid all workers waking up at the same time
			randomSleep := time.Duration(rand.Intn(5)+5) * time.Second // Sleep
			//log.Info().Msgf("Projection job worker sleeping for %s", randomSleep)
			time.Sleep(randomSleep)
			//time.Sleep(10 * time.Second)
		}
	}()
}

// ProcessQueuedProjectionJobs processes any queued projection jobs
// It checks for stalled in-progress jobs and processes new jobs concurrently
func ProcessQueuedProjectionJobs() {
	var err error

	// Try to fetch the next queued job
	var job models.ProjectionJob
	err = DB.Raw("SELECT * FROM projection_jobs WHERE status = ? ORDER BY creation_date ASC LIMIT 1", "Queued").Scan(&job).Error
	if err != nil {
		log.Error().Err(err).Msg("Failed to query for queued projection job")
		return
	}

	if job.ID == 0 {
		// No queued jobs
		return
	}

	// Attempt to atomically claim the job. If another instance already claimed it,
	// RowsAffected will be 0 and we'll safely return without processing.
	update := DB.Model(&models.ProjectionJob{}).Where("id = ? AND status = ?", job.ID, "Queued").Update("status", "In Progress")
	if update.Error != nil {
		log.Error().Err(update.Error).Msgf("Failed to claim job ID: %d", job.ID)
		return
	}
	if update.RowsAffected == 0 {
		// Lost the race to another instance
		log.Debug().Msgf("Job ID %d was claimed by another worker", job.ID)
		return
	}

	// We successfully claimed the job; now prepare its payload
	if err := json.Unmarshal([]byte(job.ProdIdsJson), &job.ProdIds); err != nil {
		log.Error().Err(err).Msgf("Failed to unmarshal ProdIdsJson for job ID: %d", job.ID)
		// mark as failed since we already claimed it
		job.Status = "Failed"
		job.StatusError = err.Error()
		if _, saveErr := SaveProjectionJob(job); saveErr != nil {
			log.Error().Err(saveErr).Msgf("Failed to update job status to Failed for job ID: %d", job.ID)
		}
		return
	}

	// Process each product in the job
	for _, id := range job.ProdIds {
		if productErr := PopulateProjectionsForProduct(id, &job); productErr != nil {
			log.Error().Err(productErr).Msgf("Failed to process product ID: %d for job ID: %d", id, job.ID)
			job.Status = "Failed"
			job.StatusError = productErr.Error()
			if _, saveErr := SaveProjectionJob(job); saveErr != nil {
				log.Error().Err(saveErr).Msgf("Failed to update job status to Failed for job ID: %d", job.ID)
			} else {
				log.Info().Msgf("Job ID: %d marked as Failed", job.ID)
			}
			return
		}
	}

	// Mark complete
	job.Status = "Complete"
	job.StatusError = ""
	if _, err := SaveProjectionJob(job); err != nil {
		log.Error().Err(err).Msgf("Failed to update job status to Complete for job ID: %d", job.ID)
		return
	}

	log.Info().Msgf("Completed processing job ID: %d", job.ID)

}
