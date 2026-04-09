package services

import (
	appLog "api/log"
	"api/models"
	"api/utils"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/jszwec/csvutil"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func GetPhiValuationTableMetaData() map[string]interface{} {
	// Retrieve model point counts by year, parameters and shocks data?
	var metadata []models.TableMetaData
	var results = make(map[string]interface{})

	// Test additions
	var phiYieldCurve = models.TableMetaData{TableType: "PHI Yield Curve", TableName: "phi_yield_curves", Data: nil, Populated: true}
	var phiParameter = models.TableMetaData{TableType: "PHI Parameters", TableName: "phi_parameters", Data: nil, Populated: true}
	var phiShock = models.TableMetaData{TableType: "PHI Shocks", TableName: "phi_shocks", Data: nil, Populated: true}
	var phiRecoveryRate = models.TableMetaData{TableType: "PHI Recovery Rates", TableName: "phi_recovery_rates", Data: nil, Populated: true}
	var phiModelPoint = models.TableMetaData{TableType: "PHI Model Points", TableName: "phi_model_points", Data: nil, Populated: true}
	var phiMortality = models.TableMetaData{TableType: "PHI Mortality", TableName: "phi_mortalities", Data: nil, Populated: true}
	var phiReinsurance = models.TableMetaData{TableType: "PHI Reinsurance", TableName: "phi_reinsurances", Data: nil, Populated: true}

	metadata = append(metadata, phiYieldCurve, phiParameter, phiShock, phiRecoveryRate, phiMortality, phiReinsurance, phiModelPoint)
	results["table_meta_data"] = metadata
	return results
}

func SavePhiTable(v *multipart.FileHeader, tableType string, year int, version string, user models.AppUser) (error, int) {
	var count int64
	var delimiter rune
	delimiterFile, err := v.Open()
	if err != nil {
		return err, int(count)
	}
	defer delimiterFile.Close()
	delimiter, err = utils.GetDelimiter(delimiterFile)

	file, err := v.Open()
	if err != nil {
		return err, int(count)
	}

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.Comma = delimiter
	dec, _ := csvutil.NewDecoder(reader)
	dec.Header()

	switch tableType {
	case "PHI Yield Curve":
		DB.Where("year = ? AND version = ?", year, version).Delete(&models.PhiYieldCurve{})
		for {
			var pyc models.PhiYieldCurve
			if err := dec.Decode(&pyc); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			pyc.Year = year
			pyc.Version = version
			pyc.CreationDate = time.Now()
			err = DB.CreateInBatches(&pyc, 100).Error
			if err != nil {
				appLog.Error("Save Quote Tables error: ", err.Error())
			}
		}
		DB.Model(&models.PhiYieldCurve{}).Where("year = ? and version", year, version).Count(&count)

	case "PHI Parameters":
		DB.Where("year = ? AND version = ?", year, version).Delete(&models.PhiParameter{})
		var pps []models.PhiParameter
		for {
			var pp models.PhiParameter
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			pp.Year = year
			pp.Version = version
			pp.CreationDate = time.Now()
			pps = append(pps, pp)
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save PHI Tables error: ", err.Error())
		}

		DB.Model(&models.PhiParameter{}).Where("year = ? and version", year, version).Count(&count)
	case "PHI Shocks":
		DB.Where("year = ? AND version = ?", year, version).Delete(&models.PhiShock{})
		for {
			var ps models.PhiShock
			if err := dec.Decode(&ps); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			ps.CreationDate = time.Now()
			err = DB.CreateInBatches(&ps, 100).Error
			if err != nil {
				appLog.Error("Save PHI Tables error: ", err.Error())
			}
		}
		DB.Model(&models.PhiShock{}).Where("year = ? and version", year, version).Count(&count)
	case "PHI Recovery Rates":
		DB.Where("year = ? AND version = ?", year, version).Delete(&models.PhiRecoveryRate{})
		for {
			var prr models.PhiRecoveryRate
			if err := dec.Decode(&prr); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			prr.Year = year
			prr.Version = version
			prr.CreationDate = time.Now()
			err = DB.CreateInBatches(&prr, 100).Error
			if err != nil {
				appLog.Error("Save PHI Tables error: ", err.Error())
			}
		}
		DB.Model(&models.PhiRecoveryRate{}).Where("year = ? and version", year, version).Count(&count)
	case "PHI Model Points":
		DB.Where("year = ? AND version = ?", year, version).Delete(&models.PhiModelPoint{})
		for {
			var pmp models.PhiModelPoint
			if err := dec.Decode(&pmp); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			pmp.Year = year
			pmp.Version = version
			pmp.CreationDate = time.Now()
			err = DB.CreateInBatches(&pmp, 100).Error
			if err != nil {
				appLog.Error("Save PHI Tables error: ", err.Error())
			}
		}
		DB.Model(&models.PhiModelPoint{}).Where("year = ? and version", year, version).Count(&count)
	case "PHI Mortality":
		DB.Where("year = ? AND version = ?", year, version).Delete(&models.PhiMortality{})
		for {
			var pm models.PhiMortality
			if err := dec.Decode(&pm); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			pm.Year = year
			pm.Version = version
			pm.CreationDate = time.Now()
			err = DB.CreateInBatches(&pm, 100).Error
			if err != nil {
				appLog.Error("Save PHI Tables error: ", err.Error())
			}
		}
		DB.Model(&models.PhiMortality{}).Where("year = ? and version", year, version).Count(&count)
	case "PHI Reinsurance":
		DB.Where("year = ? AND version = ?", year, version).Delete(&models.PhiReinsurance{})
		for {
			var prs models.PhiReinsurance
			if err := dec.Decode(&prs); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			prs.Year = year
			prs.Version = version
			prs.CreationDate = time.Now()
			err = DB.CreateInBatches(&prs, 100).Error
			if err != nil {
				appLog.Error("Save PHI Tables error: ", err.Error())
			}
		}
		DB.Model(&models.PhiReinsurance{}).Where("year = ? and version", year, version).Count(&count)
	default:
		return fmt.Errorf("unknown table type: %s", tableType), int(count)
	}
	return nil, int(count)
}

func GetPhiValuationTableData(tableType string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	switch tableType {
	case "phiyieldcurve":
		var params []models.PhiYieldCurve
		DB.Find(&params)
		b, _ := json.Marshal(&params)
		err := json.Unmarshal(b, &results)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("error unmarshalling PHI Yield Curve data: %v", err)
		}
	case "phiparameters":
		var phiRates []models.PhiParameter
		err := DB.Find(&phiRates).Error
		if err != nil {
			fmt.Println(err)
		}
		b, _ := json.Marshal(&phiRates)
		err = json.Unmarshal(b, &results)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("error unmarshalling PHI Parameters data: %v", err)
		}

	case "phishocks":
		var phiShocks []models.PhiShock
		DB.Find(&phiShocks)
		b, _ := json.Marshal(&phiShocks)
		err := json.Unmarshal(b, &results)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("error unmarshalling PHI Shocks data: %v", err)
		}

	case "phirecoveryrates":
		var phiRecoveryRates []models.PhiRecoveryRate
		DB.Find(&phiRecoveryRates)
		b, _ := json.Marshal(&phiRecoveryRates)
		err := json.Unmarshal(b, &results)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("error unmarshalling PHI Recovery Rates data: %v", err)
		}
	case "phimodelpoints":
		var phiModelPoints []models.PhiModelPoint
		DB.Find(&phiModelPoints)
		b, _ := json.Marshal(&phiModelPoints)
		err := json.Unmarshal(b, &results)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("error unmarshalling PHI Model Points data: %v", err)
		}
	case "phimortality":
		var phiMortality []models.PhiMortality
		DB.Find(&phiMortality)
		b, _ := json.Marshal(&phiMortality)
		err := json.Unmarshal(b, &results)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("error unmarshalling PHI Mortality data: %v", err)
		}
	case "phireinsurance":
		var phiReinsurance []models.PhiReinsurance
		DB.Find(&phiReinsurance)
		b, _ := json.Marshal(&phiReinsurance)
		err := json.Unmarshal(b, &results)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("error unmarshalling PHI Reinsurance data: %v", err)
		}
	}

	return results, nil
}

func GetPhiValuationTableYears(tableType string) ([]int, error) {
	var years []int
	var err error

	err = DB.Table(tableType).Distinct("year").Pluck("year", &years).Error
	if err != nil {
		return nil, err
	}
	return years, nil
}

func GetPhiValuationTableVersions(tableType string, year int) ([]string, error) {
	var versions []string
	var err error

	err = DB.Table(tableType).Where("year = ?", year).Distinct("version").Pluck("version", &versions).Error
	if err != nil {
		return nil, err
	}
	return versions, nil
}

func SavePhiShockSetting(shockSetting models.PhiShockSetting) (models.PhiShockSetting, error) {
	DB.Save(&shockSetting)
	return shockSetting, nil
}

func GetPhiShockSettings() ([]models.PhiShockSetting, error) {
	var shockSettings []models.PhiShockSetting
	err := DB.Find(&shockSettings).Error
	return shockSettings, err
}

func GetPhiShockBases() []models.AvailableBasis {
	var bases []models.AvailableBasis
	DB.Raw("select distinct shock_basis as basis from phi_shocks").Scan(&bases)
	for i, _ := range bases {
		bases[i].Name = bases[i].Basis
	}
	return bases
}

func GetPhiModelPointYears() ([]int, error) {
	var years []int
	err := DB.Model(&models.PhiModelPoint{}).Distinct("year").Pluck("year", &years).Error
	if err != nil {
		return nil, err
	}
	return years, nil
}

func GetPhiModelPointYearVersions(year int) ([]string, error) {
	var versions []string
	err := DB.Model(&models.PhiModelPoint{}).Where("year = ?", year).Distinct("version").Pluck("version", &versions).Error
	if err != nil {
		return nil, err
	}
	return versions, nil
}

func GetPhiModelPointCount() ([]models.ModelPointCount, error) {
	var counts []models.ModelPointCount
	err := DB.Model(&models.PhiModelPoint{}).
		Select("year, version, count(*) as count").
		Group("year, version").
		Scan(&counts).Error
	if err != nil {
		return nil, err
	}
	return counts, nil
}

func GetPhiModelPointsForYear(year int, version string) ([]models.PhiModelPoint, error) {
	var mps []models.PhiModelPoint
	err := DB.Where("year = ? AND version = ?", year, version).Limit(5000).Find(&mps).Error
	if err != nil {
		return nil, err
	}
	return mps, nil
}

func GetPhiModelPointsExcel(year int, version string) ([]byte, error) {
	query := fmt.Sprintf("SELECT * FROM phi_model_points WHERE year = %d AND version = '%s'", year, version)
	return exportTableToExcel(query)
}

func DeletePhiModelPoints(year int, version string) error {
	return DB.Where("year = ? AND version = ?", year, version).Delete(&models.PhiModelPoint{}).Error
}

func GetPhiParameterYears() ([]int, error) {
	var years []int
	err := DB.Model(&models.PhiParameter{}).Distinct("year").Pluck("year", &years).Error
	if err != nil {
		return nil, err
	}
	return years, nil
}

func GetPhiParameterYearVersions(year int) ([]string, error) {
	var versions []string
	err := DB.Model(&models.PhiParameter{}).Where("year = ?", year).Distinct("version").Pluck("version", &versions).Error
	if err != nil {
		return nil, err
	}
	return versions, nil
}

func GetPhiMortalityYears() ([]int, error) {
	var years []int
	err := DB.Model(&models.PhiMortality{}).Distinct("year").Pluck("year", &years).Error
	if err != nil {
		return nil, err
	}
	return years, nil
}

func GetPhiMortalityYearVersions(year int) ([]string, error) {
	var versions []string
	err := DB.Model(&models.PhiMortality{}).Where("year = ?", year).Distinct("version").Pluck("version", &versions).Error
	if err != nil {
		return nil, err
	}
	return versions, nil
}

func GetPhiRecoveryRateYears() ([]int, error) {
	var years []int
	err := DB.Model(&models.PhiRecoveryRate{}).Distinct("year").Pluck("year", &years).Error
	if err != nil {
		return nil, err
	}
	return years, nil
}

func GetPhiRecoveryRateYearVersions(year int) ([]string, error) {
	var versions []string
	err := DB.Model(&models.PhiRecoveryRate{}).Where("year = ?", year).Distinct("version").Pluck("version", &versions).Error
	if err != nil {
		return nil, err
	}
	return versions, nil
}

func GetPhiYieldCurveYears() ([]int, error) {
	var years []int
	err := DB.Model(&models.PhiYieldCurve{}).Distinct("year").Pluck("year", &years).Error
	if err != nil {
		return nil, err
	}
	return years, nil
}

func GetPhiYieldCurveYearVersions(year int) ([]string, error) {
	var versions []string
	err := DB.Model(&models.PhiYieldCurve{}).Where("year = ?", year).Distinct("version").Pluck("version", &versions).Error
	if err != nil {
		return nil, err
	}
	return versions, nil
}

func RunPhiProjection(runJobs models.RunPhiJob) error {
	// This function is a placeholder for running the actual PHI projections.
	// The implementation would depend on the specific requirements and logic of the PHI projection process.
	// It could involve fetching necessary data, performing calculations, and storing results.

	appLog.Info("Running PHI projections with parameters: ", runJobs)
	// Add your projection logic here
	for _, job := range runJobs.Jobs {
		job.UserName = runJobs.UserName

		err := PopulatePhiProjectionsForJob(&job)
		if err != nil {
			return err
		}
	}
	return nil
}

func PopulatePhiProjectionsForJob(job *models.PhiRunParameters) error {
	var phiAggregatedProjections []models.PhiProjections
	//var phiScopedAggregatedProjections []models.PhiProjections
	aggPhiProjs := make(map[string]models.PhiProjections)
	var phiSap []models.PhiScopedAggregatedProjections
	phiSapMap := make(map[string]models.PhiProjections)

	var shockSettings models.PhiShockSetting
	DB.Where("id = ?", job.ShockSettingsID).Find(&shockSettings)
	var mps []models.PhiModelPoint
	var err error
	var jobShocks []models.PhiShock

	job.CreationDate = time.Now()

	// log the start of the run
	log.Info().Msgf("Starting projections for the Job ID: %d", job.ID)

	fmt.Println("Job: ", job)
	//Get the product
	startTime := time.Now()
	testStart := time.Now()
	projectionRange := 0

	// check if there is an existing projection job id
	//err = DB.Where("projection_job_id = ?", job.ID).First(&jobProduct).Error
	//if err != nil && jobProduct.ID > 0 {
	//	// there is an existing job product, so we will not proceed with the run
	//	log.Info().Msgf("Job already exists for  projection job id: %d", job.ID)
	//	return errors.New(fmt.Sprintf("Job already exists for projection job id: %d", job.ID))
	//}

	//Get all Modelpoints for the product

	tableName := "phi_model_points"

	//var productMargins models.ProductMargins
	var parameters models.PhiParameter

	DB.Where("year = ? and version=?", job.ParameterYear, job.ParameterVersion).Find(&parameters)

	// Load modelpoints
	if job.RunSingle {
		// TODO: revert back to 1 when done testing
		err = DB.Table(tableName).Where("year =? and version=?", job.ModelPointYear, job.ModelPointVersion).Limit(MpLimitValue).Find(&mps).Error
	} else {
		err = DB.Table(tableName).Where("year =? and version=?", job.ModelPointYear, job.ModelPointVersion).Find(&mps).Error
	}

	if len(mps) == 0 {
		return errors.New("no model points found")
	}

	//Load shocks
	DB.Where("shock_basis = ?", job.ShockSettings.ShockBasis).Find(&jobShocks)

	job.ModelPointCount = len(mps)
	job.Status = "In Progress"

	DB.Save(&job)

	if err != nil {
		log.Error().Err(err).Send()
	}

	testEnd := time.Since(testStart)
	fmt.Println("Time taken to initialize Run: ", testEnd)

	wp := workerpool.New(runtime.NumCPU())

	// Create a slice to collect errors from model points
	var modelPointErrors []*PhiModelPointError
	var errorMutex sync.Mutex

	for mpIndex, mp := range mps { //Remember to adjust here later.
		mp := mp
		mpIndex := mpIndex

		if err != nil {
			log.Error().Err(err).Send()
			break
		}

		wp.Submit(func() {
			mpErr := PopulatePhiProjectionPerModelPoint(mpIndex, mp, projectionRange,
				job, &phiSap,
				parameters, &aggPhiProjs, &phiSapMap,
				jobShocks)

			// If there was an error, add it to the collection
			if mpErr != nil {
				errorMutex.Lock()
				modelPointErrors = append(modelPointErrors, mpErr)
				errorMutex.Unlock()
			}
		})
	}
	wp.StopWait()

	// Process any errors that occurred during model point processing
	var modelPointErrorsFound error

	if len(modelPointErrors) > 0 {
		log.Info().Msgf("Encountered %d errors during model point processing", len(modelPointErrors))

		// Log each error with its associated model point information
		for _, mpErr := range modelPointErrors {
			log.Error().
				Int("SpCode", mpErr.ModelPoint.Spcode).
				Str("ProductCode", mpErr.ModelPoint.ProductCode).
				Err(mpErr.Error).
				Msg("Error processing model point")
		}

		// Store the errors in the job product for later reference
		errorMessages := make([]string, 0, 5)
		for i, mpErr := range modelPointErrors {
			if i < 5 {
				errorMessages = append(errorMessages, fmt.Sprintf("Model Point (SpCode: %d, ProductCode: %s): %s",
					mpErr.ModelPoint.Spcode, mpErr.ModelPoint.ProductCode, mpErr.Error.Error()))
			}
		}

		// Join all error messages with a separator
		errorSummary := strings.Join(errorMessages, " | ")

		// Create an error to return to the caller
		modelPointErrorsFound = fmt.Errorf("encountered %d model point errors: %s", len(modelPointErrors), errorSummary)

		// Append to JobStatusError if it already has content, otherwise set it
		if job.JobStatusError != "" {
			job.JobStatusError += " | " + errorSummary
		} else {
			job.JobStatusError = "Model point errors: " + errorSummary
		}

		err = DB.Save(&job).Error
		if err != nil {
			log.Error().Err(err).Msg("Failed to save job product with model point errors")
		}

		// We don't fail the entire job because of individual model point errors
		// The errors are collected and can be reviewed later

	}

	if job.JobStatus != "cancelled" {
		fmt.Println("Total: ", len(phiAggregatedProjections))

		for _, appProj := range aggPhiProjs {
			appProj.ID = 0
			phiAggregatedProjections = append(phiAggregatedProjections, appProj)
		}

		// if the aggregated projections is empty, return an error
		if len(phiAggregatedProjections) == 0 {
			job.JobStatus = "failed"
			job.Status = "Failed"
			job.JobStatusError = "no projections were generated"
			DB.Save(&job)
			return errors.New("no projections were generated")
		}

		err = DB.Table("phi_aggregated_projections").CreateInBatches(&phiAggregatedProjections, 100).Error
		if err != nil {
			fmt.Println(err)
		}

		duration := time.Since(startTime)

		fmt.Println("Number of records processed: ", len(mps))
		fmt.Println("Number of records skipped: ", RecordsSkipped)
		fmt.Println("Duration of run: ", duration.Seconds())
		DB.Save(&job)
		//err = DB.Preload("Products").Where("id = ?", job.ID).Find(job).Error
		//if err != nil {
		//	fmt.Println(err)
		//}

		job.RunTime = job.RunTime + duration.Minutes()

		job.JobStatus = "Complete"
		job.Status = "Completed"

		RecordsSkipped = 0
		DB.Save(&job)
		phiAggregatedProjections = nil
		phiSap = nil
		phiSapMap = nil
		aggPhiProjs = nil
	} else {
		fmt.Println("Job failed or was cancelled")
		job.JobStatus = "Cancelled"
		job.Status = "Failed"
		DB.Save(&job)
		phiAggregatedProjections = nil
		phiSap = nil
		phiSapMap = nil
		aggPhiProjs = nil
	}

	// Return any model point errors that were found
	return modelPointErrorsFound

}

func PopulatePhiProjectionPerModelPoint(mpIndex int, modelPoint models.PhiModelPoint,
	projectionRange int,
	job *models.PhiRunParameters,
	phiSap *[]models.PhiScopedAggregatedProjections,
	parameters models.PhiParameter,
	aggPhiProjs *map[string]models.PhiProjections,
	phiSapMap *map[string]models.PhiProjections, phiShocks []models.PhiShock) *PhiModelPointError {
	startLoopTime := time.Now()
	var err error
	//Don't bother running if there is no existing job product
	//var existingJob models.RunPhiJob
	//joberr := DB.Where("id = ?", job.ProjectionJobID).First(&existingJob).Error
	//if existingJob.ID == 0 {
	//	job.JobStatus = "cancelled"
	//	job.JobStatusError = "the job has been cancelled"
	//	log.Info().Msgf("Job error: %s", joberr.Error())
	//	DB.Save(&job)
	//	return &PhiModelPointError{
	//		ModelPoint: modelPoint,
	//		Error:      errors.New("will not proceeed. No existing job found"),
	//	}
	//}

	var projections []models.PhiProjections
	var reinsurance models.PhiReinsurance

	reinsurance, err = GetPhiReinsurance(modelPoint)

	var index = 0
	var p models.PhiProjections
	if projectionRange == 0 {
		projectionRange = 624
	}
	valMonth, err := strconv.Atoi(job.RunDate[5:])
	if err != nil {
		valMonth = 0
	}

	valYear, err := strconv.Atoi(job.RunDate[:4])
	if err != nil {
		valYear = 0
	}

	for ; index <= projectionRange; index++ {
		var projection models.PhiProjections
		projection.RunDate = job.RunDate
		projection.RunName = job.RunName
		projection.RunId = job.ID
		projection.IFRS17Group = modelPoint.IFRS17Group
		projection.SpCode = modelPoint.Spcode
		projection.RunBasis = job.ParameterVersion

		CalculatePhiTimeAndAge(&projection, modelPoint, index, job.ProjectionJobID)
		var shock models.PhiShock
		if projection.ProjectionMonth > 0 && job.ShockSettings.ShockBasis != "" {
			shock, err = GetPhiShock(projection.ProjectionMonth, job.ShockSettings.ShockBasis, len(phiShocks))
			if err != nil {
				log.Error().Msgf("GetShock error: %s", err.Error())
			}
		}

		if projection.AgeNextBirthday > 70 {
			break
		}

		if projection.ProjectionMonth == 1 {
			projection.CalendarMonth = valMonth
		} else {
			if p.CalendarMonth+1 == 13 {
				projection.CalendarMonth = 1
			} else {
				projection.CalendarMonth = p.CalendarMonth + 1
			}
		}
		if projection.ProjectionMonth == 0 {
			projection.CalendarMonth = 0
		}

		PhiInflationFactor(index, valYear, &projection, &p, job, modelPoint, parameters, shock)
		PhiAnnuityEscalation(index, &projection, &p, modelPoint, parameters)

		PhiBaseMortalityRate(&projection, modelPoint, job, shock, parameters)

		PhiBaseRecoveryRate(&projection, modelPoint, job, shock, parameters)

		IndependentPhiMortalityRateMonthly(&projection, modelPoint)
		IndependentPhiRecoveryRateMonthly(&projection, modelPoint)

		MonthlyDependentPhiMortality(&projection, modelPoint)
		MonthlyDependentPhiRecovery(&projection, modelPoint)

		NumberOfPhiMaturities(index, &projection, &p, modelPoint)
		NumberPhiDeathsInForce(index, &projection, &p, modelPoint)
		NumberPhiRecoveriesInForce(index, &projection, &p, modelPoint)

		InitialPhiPolicy(index, &projection, modelPoint)

		IncrementalNumberPhiDeaths(&projection, &p, index, modelPoint)
		IncrementalNumberPhiRecoveries(&projection, &p, index, modelPoint)

		PhiAnnuityIncome(&projection, modelPoint)

		PhiAnnuityOutgo(&projection, modelPoint, parameters)
		PhiExpenses(&projection, &p, modelPoint, parameters, job, shock)

		PhiNetCashFlow(&projection, modelPoint, parameters)
		PhiReinsurance(&projection, &p, modelPoint, reinsurance)

		//Set values to be used on the next loop
		projections = append(projections, projection)
		p = projection
	}

	for i := len(projections) - 1; i >= 0; i-- {
		discountedErr := CalculatePhiReservesAndDiscountedValues(i, valYear, &projections[i], job, modelPoint, parameters, &projections, phiShocks)
		if discountedErr != nil {
			return &PhiModelPointError{
				ModelPoint: modelPoint,
				Error:      errors.Wrap(discountedErr, "failed to calculate reserves and discounted values"),
			}
		}
	}

	netCashFlow := 0.0
	reserveChangeA := 0.0
	investIncomeA := 0.0
	reserveChange := 0.0
	investIncome := 0.0

	prevMonth := 0

	for i, _ := range projections {

		if i == 1 && projections[i].ValuationTimeMonth == 1 {
			//Seeding the variables
			netCashFlow = projections[i-1].NetCashFlow + projections[i].NetCashFlow // point of sale and at time 1
			projections[i].ChangeInReserves = projections[i].Reserves - projections[i-1].Reserves
			projections[i].ChangeInReservesAdjusted = projections[i].ReservesAdjusted - projections[i-1].ReservesAdjusted

			interestAccrualErr := PhiInterestAccrual(i, &projections[i], projections[i-1].Reserves, projections[i-1].ReservesAdjusted, job, modelPoint, parameters, phiShocks)
			if err != nil {
				return &PhiModelPointError{
					ModelPoint: modelPoint,
					Error:      errors.Wrap(interestAccrualErr, "failed to calculate interest accrual"),
				}
			}
			var forwardRate float64
			var forwardRateErr error

			forwardRate = GetPhiForwardRate(i, job.YieldCurveYear, job.YieldCurveVersion)

			if forwardRateErr != nil {
				return &PhiModelPointError{
					ModelPoint: modelPoint,
					Error:      errors.Wrap(forwardRateErr, "failed to get forward rate"),
				}
			}
			accrualFactor := math.Pow(1.0+forwardRate, 1.0/12.0)
			investIncomeA = projections[i].InvestmentIncomeAdjusted
			investIncome = projections[i].InvestmentIncome
			reserveChangeA = projections[i-1].ReservesAdjusted*accrualFactor + projections[i].ChangeInReservesAdjusted //point of sale and time 1
			//reserveChange = projections[i].ChangeInReserves
			reserveChange = projections[i-1].Reserves*accrualFactor + projections[i].ChangeInReserves
			prevMonth = projections[i].ProjectionMonth
			projections[i].Profit = utils.FloatPrecision(-netCashFlow-reserveChange+investIncome, 6)
			projections[i].ProfitAdjusted = utils.FloatPrecision(-netCashFlow-reserveChangeA+investIncomeA, 6)
			projections[i].CorporateTax = utils.FloatPrecision(math.Max((-netCashFlow-reserveChange+investIncome), 0), 2)
			projections[i].CorporateTaxAdjusted = utils.FloatPrecision(math.Max((-netCashFlow-reserveChangeA+investIncomeA), 0), 2)

			continue
		}
		if i > 0 {

			//DB.Save(&projections[i])
			netCashFlow = projections[i].NetCashFlow

			projections[i].ChangeInReserves = projections[i].Reserves - projections[i-1].Reserves
			projections[i].ChangeInReservesAdjusted = projections[i].ReservesAdjusted - projections[i-1].ReservesAdjusted

			reserveChange = projections[i].Reserves - projections[i-1].Reserves
			reserveChangeA = projections[i].ReservesAdjusted - projections[i-1].ReservesAdjusted

			interestAccrualErr := PhiInterestAccrual(i, &projections[i], projections[i-1].Reserves, projections[i-1].ReservesAdjusted, job, modelPoint, parameters, phiShocks)
			if err != nil {
				return &PhiModelPointError{
					ModelPoint: modelPoint,
					Error:      errors.Wrap(interestAccrualErr, "failed to calculate interest accrual"),
				}
			}
			investIncome = projections[i].InvestmentIncome
			investIncomeA = projections[i].InvestmentIncomeAdjusted

			prevMonth = projections[i].ProjectionMonth

			if prevMonth <= modelPoint.TermInMonths {
				projections[i].Profit = utils.FloatPrecision(-netCashFlow-reserveChange+investIncome, 2)
				projections[i].ProfitAdjusted = utils.FloatPrecision(-netCashFlow-reserveChangeA+investIncomeA, 2)
				projections[i].CorporateTax = utils.FloatPrecision(math.Max((-netCashFlow-reserveChange+investIncome), 0), 2)
				projections[i].CorporateTaxAdjusted = utils.FloatPrecision(math.Max((-netCashFlow-reserveChangeA+investIncomeA), 0), 2)
			}
		}

	}

	for i := len(projections) - 1; i >= 0; i-- {
		var forwardRate float64
		var rdr float64
		if i == len(projections)-1 {
			projections[i].DiscountedProfit = 0
			projections[i].DiscountedProfitAdjusted = 0
		} else {
			var fwRate float64
			var forwardRateErr error
			fwRate = GetPhiForwardRate(i+1, job.YieldCurveYear, job.YieldCurveVersion)
			if forwardRateErr != nil {
				return &PhiModelPointError{
					ModelPoint: modelPoint,
					Error:      errors.Wrap(forwardRateErr, "failed to get forward rate"),
				}
			}
			forwardRate = fwRate

			//fwRate, err = GetForwardRateWithError(i+1, run.YieldcurveYear, run.YieldcurveMonth, parameters.YieldCurveCode, run.RunBasis)
			//if err != nil {
			//	return &ModelPointError{
			//		ModelPoint: modelPoint,
			//		Error:      errors.Wrap(err, "failed to get forward rate for risk discount rate"),
			//	}
			//}
			rdr = fwRate //parameters.ShareholdersRequiredMargin + fwRate

			discountFactor := math.Pow(1.0+forwardRate, -1.0/12.0)
			//discountFactorAdjusted := math.Pow(1.0+forwardRate+productMargins.InvestmentMargin, -1.0/12.0)
			riskDiscountFactor := math.Pow(1.0+rdr, -1.0/12.0)
			//riskDiscountFactorAdjusted := math.Pow(1.0+rdr+productMargins.InvestmentMargin, -1.0/12.0)

			projections[i].RiskDiscountRate = rdr
			projections[i].RiskDiscountRateAdjusted = rdr //+ productMargins.InvestmentMargin

			projections[i].DiscountedInvestmentIncome = (projections[i+1].InvestmentIncome + projections[i+1].DiscountedInvestmentIncome) * discountFactor
			projections[i].DiscountedInvestmentIncomeAdjusted = (projections[i+1].InvestmentIncomeAdjusted + projections[i+1].DiscountedInvestmentIncomeAdjusted) * discountFactor
			projections[i].DiscountedProfit = (projections[i+1].Profit + projections[i+1].DiscountedProfit) * discountFactor
			projections[i].DiscountedProfitAdjusted = (projections[i+1].ProfitAdjusted + projections[i+1].DiscountedProfitAdjusted) * discountFactor

			// discounted using Shareholders Required Return
			projections[i].VIF = (projections[i+1].Profit + projections[i+1].VIF) * riskDiscountFactor
			projections[i].VIFAdjusted = (projections[i+1].ProfitAdjusted + projections[i+1].VIFAdjusted) * riskDiscountFactor
			projections[i].DiscountedCorporateTax = (projections[i+1].CorporateTax + projections[i+1].DiscountedCorporateTax) * riskDiscountFactor
			projections[i].DiscountedCorporateTaxAdjusted = (projections[i+1].CorporateTaxAdjusted + projections[i+1].DiscountedCorporateTaxAdjusted) * riskDiscountFactor
		}
	}

	//Only adding projections for the first modelpoint to act as a control.
	if mpIndex == 0 {
		err = DB.CreateInBatches(&projections, 100).Error
		if err != nil {
			fmt.Println(err)
		}
	}

	// define list of fields to exclude from aggregation
	//excludeFields := map[string]bool{
	//	"RunId":                                   true,
	//	"JobProductID":                            true,
	//	"SpCode":                                  true,
	//	"AccidentProportion":                      true,
	//	"InflationFactor":                         true,
	//	"InflationFactorAdjusted":                 true,
	//	"LapseMargin":                             true,
	//	"PremiumWaiverOnFactor":                   true,
	//	"PaidUpOnFactor":                          true,
	//	"MainMemberMortalityRate":                 true,
	//	"MainMemberMortalityRateAdjusted":         true,
	//	"BaseLapse":                               true,
	//	"BaseLapseAdjusted":                       true,
	//	"ContractingPartyAlivePortion":            true,
	//	"ContractingPartyAlivePortionAdjusted":    true,
	//	"ContractingPartyPolicyLapse":             true,
	//	"ContractingPartyPolicyLapseAdjusted":     true,
	//	"NonLifeMonthlyRiskRate":                  true,
	//	"NonLifeMonthlyRiskRateAdjusted":          true,
	//	"BaseMortalityRate":                       true,
	//	"BaseMortalityRateAdjusted":               true,
	//	"BaseIndependentLapse":                    true,
	//	"BaseIndependentLapseAdjusted":            true,
	//	"BaseRetrenchmentRate":                    true,
	//	"BaseRetrenchmentRateAdjusted":            true,
	//	"BaseDisabilityIncidenceRate":             true,
	//	"BaseDisabilityIncidenceRateAdjusted":     true,
	//	"MainMemberMortalityRateByMonth":          true,
	//	"MainMemberMortalityRateAdjustedByMonth":  true,
	//	"IndependentMortalityRateMonthly":         true,
	//	"IndependentMortalityRateAdjustedByMonth": true,
	//	"IndependentLapseMonthly":                 true,
	//	"IndependentLapseMonthlyAdjusted":         true,
	//	"IndependentRetrenchmentMonthly":          true,
	//	"IndependentRetrenchmentMonthlyAdjusted":  true,
	//	"IndependentDisabilityMonthly":            true,
	//	"IndependentDisabilityMonthlyAdjusted":    true,
	//	"MonthlyDependentMortality":               true,
	//	"MonthlyDependentMortalityAdjusted":       true,
	//	"MonthlyDependentLapse":                   true,
	//	"MonthlyDependentLapseAdjusted":           true,
	//	"MonthlyDependentRetrenchment":            true,
	//	"MonthlyDependentRetrenchmentAdjusted":    true,
	//	"MonthlyDependentDisability":              true,
	//	"MonthlyDependentDisabilityAdjusted":      true,
	//	"OutstandingSumAssured":                   true,
	//}

	//out <- projections[:run.AggregationPeriod]
	if len(projections) >= job.AggregationPeriod {
		UpdatePhiAggregatedProjections(projections[:job.AggregationPeriod], aggPhiProjs)
	} else {
		UpdatePhiAggregatedProjections(projections, aggPhiProjs)
	}

	// Update points_done safely without overwriting the full struct (avoids race conditions in concurrent workerpool)
	newPointsDone := mpIndex + 1
	DB.Model(&models.PhiRunParameters{}).
		Where("id = ? AND points_done < ?", job.ID, newPointsDone).
		UpdateColumn("points_done", newPointsDone)

	//job.PointsDone = job.PointsDone + 1
	//DB.Where("id = ?", job.ID).Save(&job)
	projections = nil
	endLoopTime := time.Since(startLoopTime)
	fmt.Println("Time taken for modelpoint: ", endLoopTime.Seconds())
	return nil
}

func CalculatePhiTimeAndAge(projection *models.PhiProjections, modelPoint models.PhiModelPoint, i int, jobProductId int) {
	projection.JobProductID = jobProductId
	projection.ProjectionMonth = i
	projection.ProjectionYear = int(math.Ceil(float64(i) / 12.0))
	projection.PolicyNumber = modelPoint.PolicyNumber
	projection.ValuationTimeMonth = modelPoint.DurationInForceMonths + i
	projection.ValuationTimeYear = utils.FloatPrecision(float64(projection.ValuationTimeMonth)/12, 2)
	projection.AgeNextBirthday = int(math.Ceil(float64(modelPoint.AgeAtEntry) + projection.ValuationTimeYear))
}

func GetPhiShock(month int, shockBasis string, PhiShockCount int) (models.PhiShock, error) {
	var shock models.PhiShock
	if month > PhiShockCount {
		month = PhiShockCount
	}

	key := "shock_" + "_" + shockBasis + "_" + strconv.Itoa(month)

	cacheKey := key
	cached, found := Cache.Get(cacheKey)

	if found {
		shock = cached.(models.PhiShock)
		return shock, nil
	} else {
		err := DB.Where("shock_basis = ? and projection_month = ?", shockBasis, month).First(&shock).Error
		if err != nil {
			return shock, err
		}
		Cache.Set(key, shock, 1)
		return shock, nil
	}
}

func GetPhiReinsurance(mp models.PhiModelPoint) (models.PhiReinsurance, error) {
	var reinsurance models.PhiReinsurance

	key := "phi_reinsurance" + "_" + mp.TreatyCode

	cacheKey := key
	cached, found := Cache.Get(cacheKey)

	if found {
		reinsurance = cached.(models.PhiReinsurance)
		return reinsurance, nil
	} else {
		err := DB.Where("version = ? ", mp.TreatyCode).First(&reinsurance).Error
		if err != nil {
			return reinsurance, err
		}
		Cache.Set(key, reinsurance, 1)
		return reinsurance, nil
	}
}

// InflationFactor projects inflation factor at each projection period
// Is used to inflate expenses
func PhiInflationFactor(i, valYear int, projection *models.PhiProjections, p *models.PhiProjections, job *models.PhiRunParameters, mp models.PhiModelPoint, parameters models.PhiParameter, shock models.PhiShock) {
	if i == 0 {
		projection.InflationFactor = 1
		projection.InflationFactorAdjusted = 1
	} else if projection.ValuationTimeMonth <= mp.TermInMonths {
		var inflationFactor float64 = 0
		inflationFactor = GetPhiInflationFactor(projection.ProjectionMonth, job.YieldCurveYear, job.YieldCurveVersion)
		tempInflationFactor := inflationFactor
		if job.ShockSettings.Inflation {
			inflationFactor = tempInflationFactor + math.Max(tempInflationFactor*shock.MultiplicativeInflation, shock.AdditiveInflation)
		}
		projection.InflationFactor = utils.FloatPrecision(p.InflationFactor*math.Pow(1+inflationFactor, 1/12.0), defaultPrecision)
		projection.InflationFactorAdjusted = utils.FloatPrecision(p.InflationFactorAdjusted*math.Pow(1+inflationFactor, 1/12.0), defaultPrecision)
		projection.AnnuityEscalationRate = math.Min(inflationFactor, mp.AnnuityEscalation)
	} else {
		projection.InflationFactor = 0
		projection.InflationFactorAdjusted = 0
		projection.AnnuityEscalationRate = 0
	}
}

func GetPhiInflationFactor(projectionMonth int, year int, version string) float64 {
	var yieldCurve models.PhiYieldCurve
	key := strconv.Itoa(year) + "-y-curve-" + strconv.Itoa(projectionMonth) + "-" + version
	result, found := Cache.Get(key)

	if found {
		return result.(float64)
	} else {
		err := DB.Table("phi_yield_curves").Where("year = ? and proj_time = ? and version=?", year, projectionMonth, version).First(&yieldCurve).Error
		if err != nil {
			fmt.Println("Yield Curve: ", errors.WithStack(err))
		}
		success := Cache.Set(key, yieldCurve.Inflation, 1)
		//err = redisClient.Set(key, yieldCurve.Inflation, 0).Err()
		if !success {
			fmt.Println("cache: key not stored")
		}
	}
	return yieldCurve.Inflation
}

func PhiAnnuityEscalation(i int, projection *models.PhiProjections, p *models.PhiProjections, modelPoint models.PhiModelPoint, parameters models.PhiParameter) {
	if i == 0 {
		projection.AnnuityEscalation = 1
	} else if projection.ValuationTimeMonth <= modelPoint.TermInMonths {
		if parameters.AnnuityCalendarMonthEscalationIndicator {
			if projection.CalendarMonth == modelPoint.AnnuityEscalationMonth {
				projection.AnnuityEscalation = p.AnnuityEscalation * (1 + projection.AnnuityEscalationRate)
			} else {
				projection.AnnuityEscalation = p.AnnuityEscalation
			}

		} else {
			if (projection.ValuationTimeMonth-1)%12 == 0 && projection.ValuationTimeMonth > 1 {
				projection.AnnuityEscalation = p.AnnuityEscalation * (1 + modelPoint.AnnuityEscalation)
			} else {
				projection.AnnuityEscalation = p.AnnuityEscalation
			}
		}
	} else {
		projection.AnnuityEscalation = 0
	}
}

// BaseMortalityRate reads base mortality rate by rating factors(age,gender)
func PhiBaseMortalityRate(projection *models.PhiProjections, modelPoint models.PhiModelPoint, run *models.PhiRunParameters, shock models.PhiShock, parameters models.PhiParameter) {
	if projection.ValuationTimeMonth <= modelPoint.TermInMonths {
		//if modelPoint.MemberType == "MM" {
		//	projection.BaseMortalityRate = projection.MainMemberMortalityRate
		//	projection.BaseMortalityRateAdjusted = projection.MainMemberMortalityRateAdjusted
		//} else {
		if projection.AgeNextBirthday > 70 {
			projection.BaseMortalityRate = 0
			projection.BaseMortalityRateAdjusted = 0
		} else {
			var resp float64 = 0
			resp = math.Min(GetPhiMortalityRate(projection, modelPoint, run), 1)
			if run.ShockSettings.Mortality {
				projection.BaseMortalityRate = utils.FloatPrecision(math.Max(0, math.Min(1, resp*(1+shock.MultiplicativeMortality)+shock.AdditiveMortality)), defaultPrecision)
			} else if run.ShockSettings.MortalityCatastrophe {
				monthlyResp := 1.0 - math.Pow(1.0-resp, 1/12.0)
				shockedMonthly := monthlyResp + shock.MortalityCatastropheMultiplier*math.Min(math.Max(shock.MultiplicativeMortality*monthlyResp*1000.0+shock.AdditiveMortality, shock.MortalityCatastropheFloor), shock.MortalityCatastropheCeiling)/1000.0
				shockedAnnual := 1.0 - math.Pow(1.0-shockedMonthly, 12)
				projection.BaseMortalityRate = utils.FloatPrecision(shockedAnnual, defaultPrecision)
			} else {
				projection.BaseMortalityRate = utils.FloatPrecision(resp, defaultPrecision)
			}
			projection.BaseMortalityRateAdjusted = utils.FloatPrecision(math.Min(projection.BaseMortalityRate, 1), defaultPrecision)
		}
		//	}
	} else {
		projection.BaseMortalityRate = 0
		projection.BaseMortalityRateAdjusted = 0
	}
}

func GetPhiMortalityRate(projection *models.PhiProjections, modelPoint models.PhiModelPoint, job *models.PhiRunParameters) float64 {
	tablename := "phi_mortalities"
	var mortalityRow models.PhiMortality

	key := strconv.Itoa(projection.AgeNextBirthday) + modelPoint.Gender + strconv.Itoa(job.MortalityYear) + job.MortalityVersion
	cacheKey := tablename + "_" + key
	cached, found := Cache.Get(cacheKey)

	if found {
		result := cached.(float64)
		//if result > 0 {
		return result
		//}
	} else {
		//fmt.Println("cache missed: ", key)
	}

	err := DB.Table(tablename).Where("year = ? and version = ? and anb=? and  gender=?", job.MortalityYear, job.MortalityVersion, projection.AgeNextBirthday, modelPoint.Gender).First(&mortalityRow).Error

	if err != nil {
		fmt.Println("Phi Mortality: ", errors.WithStack(err))
	}

	Cache.Set(cacheKey, mortalityRow.MortalityRate, 1)
	//time.Sleep(5 * time.Millisecond)

	return mortalityRow.MortalityRate

}

func PhiBaseRecoveryRate(projection *models.PhiProjections, modelPoint models.PhiModelPoint, run *models.PhiRunParameters, shock models.PhiShock, parameters models.PhiParameter) {
	if projection.ValuationTimeMonth <= modelPoint.TermInMonths {
		//if modelPoint.MemberType == "MM" {
		//	projection.BaseMortalityRate = projection.MainMemberMortalityRate
		//	projection.BaseMortalityRateAdjusted = projection.MainMemberMortalityRateAdjusted
		//} else {
		if projection.AgeNextBirthday > 70 {
			projection.BaseRecoveryRate = 0
			projection.BaseRecoveryRateAdjusted = 0
		} else {
			var resp float64 = 0
			resp = math.Min(GetPhiRecoveryRate(projection, modelPoint, run), 1)
			if run.ShockSettings.Recovery {
				projection.BaseRecoveryRate = utils.FloatPrecision(math.Max(0, math.Min(1, resp*(1+shock.MultiplicativeRecovery)+shock.AdditiveRecovery)), defaultPrecision)
			} else if run.ShockSettings.MorbidityCatastrophe {
				monthlyResp := 1.0 - math.Pow(1.0-resp, 1/12.0)
				shockedMonthly := monthlyResp + shock.CATScalar*monthlyResp*1000.0*shock.MorbidityCatastropheMultiplier/1000.0
				shockedAnnual := 1.0 - math.Pow(1.0-shockedMonthly, 12)
				projection.BaseRecoveryRate = utils.FloatPrecision(shockedAnnual, defaultPrecision)
			} else {
				projection.BaseRecoveryRate = utils.FloatPrecision(resp, defaultPrecision)
			}
			projection.BaseRecoveryRateAdjusted = utils.FloatPrecision(math.Min(projection.BaseRecoveryRate, 1), defaultPrecision)
		}
		//	}
	} else {
		projection.BaseMortalityRate = 0
		projection.BaseMortalityRateAdjusted = 0
	}
}

func GetPhiRecoveryRate(projection *models.PhiProjections, modelPoint models.PhiModelPoint, job *models.PhiRunParameters) float64 {
	tablename := "phi_recovery_rates"
	var recoveryRateRow models.PhiRecoveryRate

	key := strconv.Itoa(job.RecoveryYear) + "_" + job.RecoveryVersion + "_" + strconv.Itoa(projection.AgeNextBirthday) + modelPoint.Gender + "_" + strconv.Itoa(projection.ValuationTimeMonth) + "_" + strconv.Itoa(modelPoint.WaitingPeriod)
	cacheKey := tablename + "_" + key
	cached, found := Cache.Get(cacheKey)

	if found {
		result := cached.(float64)
		//if result > 0 {
		return result
		//}
	} else {
		//fmt.Println("cache missed: ", key)
	}

	err := DB.Table(tablename).Where("year = ? and version = ? and anb=? and  gender=? and duration_if_m=? and waiting_period=?", job.MortalityYear, job.MortalityVersion, projection.AgeNextBirthday, modelPoint.Gender, projection.ValuationTimeMonth, modelPoint.WaitingPeriod).First(&recoveryRateRow).Error

	if err != nil {
		fmt.Println("Phi Mortality: ", errors.WithStack(err))
	}

	Cache.Set(cacheKey, recoveryRateRow.RecoveryRate, 1)
	//time.Sleep(5 * time.Millisecond)

	return recoveryRateRow.RecoveryRate

}

func MonthlyDependentPhiMortality(projection *models.PhiProjections, mp models.PhiModelPoint) {
	if projection.ValuationTimeMonth <= mp.TermInMonths {
		projection.MonthlyDependentMortality = utils.FloatPrecision(projection.IndependentMortalityRateMonthly*(1-TimingRecoveryZero*projection.IndependentRecoveryRateMonthly), defaultPrecision)
		projection.MonthlyDependentMortalityAdjusted = utils.FloatPrecision(projection.IndependentMortalityRateAdjustedByMonth*(1-TimingRecoveryZero*projection.IndependentRecoveryRateMonthlyAdjusted), defaultPrecision)
	} else {
		projection.MonthlyDependentMortality = 0
		projection.MonthlyDependentMortalityAdjusted = 0
	}

}

func MonthlyDependentPhiRecovery(projection *models.PhiProjections, mp models.PhiModelPoint) {
	if projection.ValuationTimeMonth <= mp.TermInMonths {
		projection.MonthlyDependentRecovery = utils.FloatPrecision(projection.IndependentRecoveryRateMonthly*(1-TimingMortalityOne*projection.IndependentMortalityRateMonthly), defaultPrecision)
		projection.MonthlyDependentRecoveryAdjusted = utils.FloatPrecision(projection.IndependentRecoveryRateMonthlyAdjusted*(1-TimingMortalityOne*projection.IndependentMortalityRateAdjustedByMonth), defaultPrecision)
	} else {
		projection.MonthlyDependentRecovery = 0
		projection.MonthlyDependentRecoveryAdjusted = 0
	}

}

func IndependentPhiMortalityRateMonthly(projection *models.PhiProjections, mp models.PhiModelPoint) {
	if projection.ValuationTimeMonth <= mp.TermInMonths {
		projection.IndependentMortalityRateMonthly = utils.FloatPrecision(1-math.Pow(1-projection.BaseMortalityRate, 1/12.0), defaultPrecision)
		projection.IndependentMortalityRateAdjustedByMonth = utils.FloatPrecision(1-math.Pow(1-projection.BaseMortalityRateAdjusted, 1/12.0), defaultPrecision)
	} else {
		projection.IndependentMortalityRateMonthly = 0
		projection.IndependentMortalityRateAdjustedByMonth = 0
	}
}

func IndependentPhiRecoveryRateMonthly(projection *models.PhiProjections, mp models.PhiModelPoint) {
	if projection.ValuationTimeMonth <= mp.TermInMonths {
		projection.IndependentRecoveryRateMonthly = utils.FloatPrecision(1-math.Pow(1-projection.BaseRecoveryRate, 1/12.0), defaultPrecision)
		projection.IndependentRecoveryRateMonthlyAdjusted = utils.FloatPrecision(1-math.Pow(1-projection.BaseRecoveryRateAdjusted, 1/12.0), defaultPrecision)
	} else {
		projection.IndependentRecoveryRateMonthly = 0
		projection.IndependentRecoveryRateMonthlyAdjusted = 0
	}
}

// NumberOfMaturities computes number of cumulative maturities over the projection period
func NumberOfPhiMaturities(index int, projection *models.PhiProjections, p *models.PhiProjections, mp models.PhiModelPoint) {
	//if utils.StatesContains(&states, Maturity) {
	if index == 0 {
		projection.NumberOfMaturities = 0
		projection.NumberOfMaturitiesAdjusted = 0
		if mp.TermInMonths == 0 {
			projection.NumberOfMaturities = 1
			projection.NumberOfMaturitiesAdjusted = 1
		}
	} else {
		if projection.ValuationTimeMonth <= mp.TermInMonths+1 {
			if projection.ValuationTimeMonth == mp.TermInMonths+1 && mp.TermInMonths != 0 {
				projection.NumberOfMaturities = utils.FloatPrecision(math.Max(startingInitialPolicy-p.NumberOfMaturities-p.NumberOfDeathsInForce-p.NumberOfRecoveries, 0)+p.NumberOfMaturities, defaultdecrementPrecision)
				projection.NumberOfMaturitiesAdjusted = utils.FloatPrecision(math.Max(startingInitialPolicy-p.NumberOfMaturitiesAdjusted-p.NumberOfDeathsInForceAdjusted-p.NumberOfRecoveriesAdjusted, 0)+p.NumberOfMaturitiesAdjusted, defaultdecrementPrecision)
			} else {
				projection.NumberOfMaturities = p.NumberOfMaturities
				projection.NumberOfMaturitiesAdjusted = p.NumberOfMaturitiesAdjusted
			}

		} else {
			projection.NumberOfMaturities = p.NumberOfMaturities
			projection.NumberOfMaturitiesAdjusted = p.NumberOfMaturitiesAdjusted
		}
	}
}

func NumberPhiDeathsInForce(index int, projection *models.PhiProjections, p *models.PhiProjections, mp models.PhiModelPoint) {
	if index == 0 {
		projection.NumberOfDeathsInForce = 0
		projection.NumberOfDeathsInForceAdjusted = 0
	} else {
		if projection.ValuationTimeMonth <= mp.TermInMonths {
			//projection.NaturalDeathsInForce
			projection.NumberOfDeathsInForce = utils.FloatPrecision(math.Max(p.InitialPolicy-projection.NumberOfMaturities, 0)*projection.MonthlyDependentMortality+p.NumberOfDeathsInForce, defaultdecrementPrecision)
			projection.NumberOfDeathsInForceAdjusted = utils.FloatPrecision(math.Max(p.InitialPolicyAdjusted-projection.NumberOfMaturitiesAdjusted, 0)*projection.MonthlyDependentMortalityAdjusted+p.NumberOfDeathsInForceAdjusted, defaultdecrementPrecision)
		} else {
			projection.NumberOfDeathsInForce = p.NumberOfDeathsInForce
			projection.NumberOfDeathsInForceAdjusted = p.NumberOfDeathsInForceAdjusted
		}
	}
}

func NumberPhiRecoveriesInForce(index int, projection *models.PhiProjections, p *models.PhiProjections, mp models.PhiModelPoint) {
	if index == 0 {
		projection.NumberOfRecoveries = 0
		projection.NumberOfRecoveriesAdjusted = 0
	} else {
		if projection.ValuationTimeMonth <= mp.TermInMonths {
			//projection.NaturalDeathsInForce
			projection.NumberOfRecoveries = utils.FloatPrecision(math.Max(p.InitialPolicy-projection.NumberOfMaturities, 0)*projection.MonthlyDependentRecovery+p.NumberOfRecoveries, defaultdecrementPrecision)
			projection.NumberOfRecoveriesAdjusted = utils.FloatPrecision(math.Max(p.InitialPolicyAdjusted-projection.NumberOfMaturitiesAdjusted, 0)*projection.MonthlyDependentRecoveryAdjusted+p.NumberOfRecoveriesAdjusted, defaultdecrementPrecision)
		} else {
			projection.NumberOfRecoveries = p.NumberOfRecoveries
			projection.NumberOfRecoveriesAdjusted = p.NumberOfRecoveriesAdjusted
		}
	}
}

func InitialPhiPolicy(index int, projection *models.PhiProjections, modelPoint models.PhiModelPoint) {
	if index == 0 {
		projection.InitialPolicy = 1
		projection.InitialPolicyAdjusted = 1

		startingInitialPolicy = 1         // projection.InitialPolicy
		startingInitialPolicyAdjusted = 1 // projection.InitialPolicyAdjusted
		if modelPoint.TermInMonths == 0 || modelPoint.DurationInForceMonths > modelPoint.TermInMonths {
			projection.NumberOfMaturities = 1
			projection.NumberOfMaturitiesAdjusted = 1
			projection.InitialPolicy = 0
			projection.InitialPolicyAdjusted = 0
		}
	} else {
		if projection.ValuationTimeMonth <= modelPoint.TermInMonths {
			projection.InitialPolicy = utils.FloatPrecision(math.Max(
				startingInitialPolicy-
					projection.NumberOfMaturities-
					projection.NumberOfDeathsInForce-
					projection.NumberOfRecoveries,
				0), defaultdecrementPrecision)
			projection.InitialPolicyAdjusted = utils.FloatPrecision(math.Max(
				startingInitialPolicyAdjusted-
					projection.NumberOfMaturitiesAdjusted-
					projection.NumberOfDeathsInForceAdjusted-
					projection.NumberOfRecoveriesAdjusted,
				0), defaultdecrementPrecision)
		} else {
			projection.InitialPolicy = 0
			projection.InitialPolicyAdjusted = 0
		}
	}
}

func IncrementalNumberPhiDeaths(projection *models.PhiProjections, p *models.PhiProjections, index int, mp models.PhiModelPoint) {
	if index == 0 || projection.ValuationTimeMonth > mp.TermInMonths {
		projection.IncrementalDeaths = 0
		projection.IncrementalDeathsAdjusted = 0
		return
	}
	projection.IncrementalDeaths = utils.FloatPrecision((projection.NumberOfDeathsInForce)-
		(p.NumberOfDeathsInForce), defaultdecrementPrecision)
	projection.IncrementalDeathsAdjusted = utils.FloatPrecision((projection.NumberOfDeathsInForceAdjusted)-
		(p.NumberOfDeathsInForceAdjusted), defaultdecrementPrecision)
}

func IncrementalNumberPhiRecoveries(projection *models.PhiProjections, p *models.PhiProjections, index int, mp models.PhiModelPoint) {
	if index == 0 || projection.ValuationTimeMonth > mp.TermInMonths {
		projection.IncrementalRecoveries = 0
		projection.IncrementalRecoveriesAdjusted = 0
		return
	}
	projection.IncrementalRecoveries = utils.FloatPrecision((projection.IncrementalRecoveries)-
		(p.IncrementalRecoveries), defaultdecrementPrecision)
	projection.IncrementalRecoveriesAdjusted = utils.FloatPrecision((projection.IncrementalRecoveriesAdjusted)-
		(p.IncrementalRecoveriesAdjusted), defaultdecrementPrecision)
}

func PhiAnnuityIncome(projection *models.PhiProjections, mp models.PhiModelPoint) {
	if projection.ValuationTimeMonth <= mp.TermInMonths {
		projection.AnnuityIncome = mp.AnnuityIncome * projection.AnnuityEscalation
	} else {
		projection.AnnuityIncome = 0
	}
}

func PhiAnnuityOutgo(projection *models.PhiProjections, mp models.PhiModelPoint, params models.PhiParameter) {
	if projection.ValuationTimeMonth > mp.TermInMonths || projection.ValuationTimeMonth <= mp.WaitingPeriod {
		projection.AnnuityOutgo = 0
		projection.AnnuityOutgoAdjusted = 0
	} else {

		if mp.ClaimStatus == Notified {
			projection.AnnuityOutgo = utils.FloatPrecision(projection.InitialPolicy*projection.AnnuityIncome*params.PendingAcceptanceRatio, defaultPrecision)
			projection.AnnuityOutgoAdjusted = utils.FloatPrecision(projection.InitialPolicyAdjusted*projection.AnnuityIncome*params.PendingAcceptanceRatio, defaultPrecision)
		}

		if mp.ClaimStatus == Pending {
			projection.AnnuityOutgo = utils.FloatPrecision(projection.InitialPolicy*projection.AnnuityIncome*params.PendingAcceptanceRatio, defaultPrecision)
			projection.AnnuityOutgoAdjusted = utils.FloatPrecision(projection.InitialPolicyAdjusted*projection.AnnuityIncome*params.PendingAcceptanceRatio, defaultPrecision)

		}
		projection.AnnuityOutgo = utils.FloatPrecision(projection.InitialPolicy*projection.AnnuityIncome, defaultPrecision)
		projection.AnnuityOutgoAdjusted = utils.FloatPrecision(projection.InitialPolicyAdjusted*projection.AnnuityIncome, defaultPrecision)
	}
}

// Reinsurance
func PhiReinsurance(projection *models.PhiProjections, p *models.PhiProjections, modelPoint models.PhiModelPoint, reinsurance models.PhiReinsurance) {
	if projection.ValuationTimeMonth <= modelPoint.TermInMonths {
		if projection.ProjectionMonth == 0 || projection.ValuationTimeMonth <= modelPoint.WaitingPeriod {
			projection.LeadReAnnuityOutgo = 0
			projection.LeadReAnnuityOutgoAdjusted = 0

		} else {
			var level1value float64
			var level2value float64
			var level3value float64
			level1value = math.Min(projection.AnnuityIncome, reinsurance.Level1Upperbound-reinsurance.Level1Lowerbound)
			level2value = math.Min(projection.AnnuityIncome-level1value, reinsurance.Level2Upperbound-reinsurance.Level2Lowerbound)
			level3value = math.Min(projection.AnnuityIncome-level1value-level2value, reinsurance.Level3Upperbound-reinsurance.Level3Lowerbound)
			projection.CededAnnuityIncome = level1value*reinsurance.Level1CededProportion + level2value*reinsurance.Level2CededProportion + level3value*reinsurance.Level3CededProportion
			projection.LeadReAnnuityIncome = projection.CededAnnuityIncome * reinsurance.LeadReProportion
			projection.Re2AnnuityIncome = projection.CededAnnuityIncome * reinsurance.Re2Proportion
			projection.Re3AnnuityIncome = projection.CededAnnuityIncome * reinsurance.Re3Proportion
			projection.LeadReAnnuityOutgo = projection.LeadReAnnuityIncome * p.InitialPolicy
			projection.LeadReAnnuityOutgoAdjusted = projection.LeadReAnnuityIncome * p.InitialPolicyAdjusted
			projection.Re2AnnuityOutgo = projection.Re2AnnuityIncome * p.InitialPolicy
			projection.Re2AnnuityOutgoAdjusted = projection.Re2AnnuityIncome * p.InitialPolicyAdjusted
			projection.Re3AnnuityOutgo = projection.Re3AnnuityIncome * p.InitialPolicy
			projection.Re3AnnuityOutgoAdjusted = projection.Re3AnnuityIncome * p.InitialPolicyAdjusted

		}
		projection.NetReinsuranceCashflow = projection.LeadReAnnuityOutgo + projection.Re2AnnuityOutgo + projection.Re3AnnuityOutgo
		projection.NetReinsuranceCashflowAdjusted = projection.LeadReAnnuityOutgoAdjusted + projection.Re2AnnuityOutgoAdjusted + projection.Re3AnnuityOutgoAdjusted

	} else {
		projection.LeadReAnnuityOutgo = 0
		projection.LeadReAnnuityOutgoAdjusted = 0
		projection.Re2AnnuityOutgo = 0
		projection.Re2AnnuityOutgoAdjusted = 0
		projection.Re3AnnuityOutgo = 0
		projection.Re3AnnuityOutgoAdjusted = 0
		projection.NetReinsuranceCashflow = 0
		projection.NetReinsuranceCashflowAdjusted = 0
	}
}

func PhiExpenses(projection *models.PhiProjections, p *models.PhiProjections, mp models.PhiModelPoint, parameters models.PhiParameter, run *models.PhiRunParameters, shock models.PhiShock) {
	var exprand_shocked float64
	var exprand_shockedAdjusted float64
	var exprand float64 = 0
	var exprandAdjusted float64 = 0
	var exprand_claims float64
	var exprand_claimsadjusted float64
	var exprand_shockedclaims float64
	var exprand_shockedclaimsinflationAdjusted float64
	var exprand_shockedClaimsAdjusted float64
	var initial_attr_exp_prop float64
	var renewal_attr_exp_prop float64
	var claims_attr_exp_prop float64

	initial_attr_exp_prop = parameters.AttributableExpenseProportion
	renewal_attr_exp_prop = parameters.AttributableExpenseProportion
	claims_attr_exp_prop = parameters.AttributableExpenseProportion

	if projection.ValuationTimeMonth == 0 {
		exprand = 0 * initial_attr_exp_prop
		if run.ShockSettings.Expense {
			exprand_shocked = exprand*(1+shock.MultiplicativeExpense) + shock.AdditiveExpense/12.0
		} else {
			exprand_shocked = exprand
		}
		exprand_shockedAdjusted = exprand_shocked
	} else {
		if projection.ProjectionMonth == 0 {
			exprand_shocked = 0
		} else {
			exprand = (parameters.AnnualRenewalExpenseAmount * projection.InflationFactor / 12.0) * renewal_attr_exp_prop
			exprand_claims = (parameters.ClaimsExpenseProportion*projection.AnnuityIncome + parameters.AnnualClaimsExpenseAmount*projection.InflationFactor) * claims_attr_exp_prop
			exprandAdjusted = (parameters.AnnualRenewalExpenseAmount * projection.InflationFactorAdjusted / 12.0) * renewal_attr_exp_prop * renewal_attr_exp_prop
			exprand_claimsadjusted = (parameters.ClaimsExpenseProportion*projection.AnnuityIncome + parameters.AnnualClaimsExpenseAmount*projection.InflationFactorAdjusted) * claims_attr_exp_prop
			if run.ShockSettings.Expense {
				exprand_shocked = exprand*(1+shock.MultiplicativeExpense) + shock.AdditiveExpense*projection.InflationFactor/12.0
				exprand_shockedclaims = exprand_claims * (1 + shock.MultiplicativeExpense) //+ shock.AdditiveExpense * projection.InflationFactor/12.0  *claims expense can only be shocked proportionally
				exprand_shockedAdjusted = exprandAdjusted*(1+shock.MultiplicativeExpense) + shock.AdditiveExpense*projection.InflationFactorAdjusted/12.0
				exprand_shockedclaimsinflationAdjusted = exprand_claimsadjusted * (1 + shock.MultiplicativeExpense) //+ shock.AdditiveExpense * projection.InflationFactorAdjusted/12.0 * claims expense can only be shocked proportionally
			} else {
				exprand_shocked = exprand
				exprand_shockedclaims = exprand_claims
				exprand_shockedAdjusted = exprandAdjusted
				exprand_shockedclaimsinflationAdjusted = exprand_claimsadjusted
			}
			exprand_shockedAdjusted = exprand_shockedAdjusted
			exprand_shockedClaimsAdjusted = exprand_shockedclaimsinflationAdjusted
		}
	}

	if projection.ValuationTimeMonth <= mp.TermInMonths {

		if mp.ClaimStatus == Notified {
			projection.RenewalExpenses = utils.FloatPrecision((exprand_shocked*p.InitialPolicy+exprand_shockedclaims*p.InitialPolicy)*parameters.NotifiedAcceptanceRatio, defaultPrecision)
			projection.RenewalExpensesAdjusted = utils.FloatPrecision((exprand_shockedAdjusted*p.InitialPolicyAdjusted+exprand_shockedClaimsAdjusted*p.InitialPolicyAdjusted)*parameters.NotifiedAcceptanceRatio, defaultPrecision)
		}
		if mp.ClaimStatus == Pending {
			projection.RenewalExpenses = utils.FloatPrecision((exprand_shocked*p.InitialPolicy+exprand_shockedclaims*p.InitialPolicy)*parameters.PendingAcceptanceRatio, defaultPrecision)
			projection.RenewalExpensesAdjusted = utils.FloatPrecision((exprand_shockedAdjusted*p.InitialPolicyAdjusted+exprand_shockedClaimsAdjusted*p.InitialPolicyAdjusted)*parameters.PendingAcceptanceRatio, defaultPrecision)
		}
		if mp.ClaimStatus != Pending && mp.ClaimStatus != Notified {
			projection.RenewalExpenses = utils.FloatPrecision(exprand_shocked*p.InitialPolicy+exprand_shockedclaims*p.InitialPolicy, defaultPrecision)
			projection.RenewalExpensesAdjusted = utils.FloatPrecision(exprand_shockedAdjusted*p.InitialPolicyAdjusted+exprand_shockedClaimsAdjusted*p.InitialPolicyAdjusted, defaultPrecision)
		}
	}
}

func PhiNetCashFlow(projection *models.PhiProjections, mp models.PhiModelPoint, parameters models.PhiParameter) {
	if projection.ProjectionMonth == 0 || projection.ValuationTimeMonth > mp.TermInMonths+1 {
		projection.NetCashFlow = 0
		projection.NetCashFlowAdjusted = 0
	} else {
		projection.NetCashFlow = projection.DeathOutgo + projection.AnnuityOutgo + projection.Rider + projection.RenewalExpenses - projection.PremiumIncome
		projection.NetCashFlowAdjusted = projection.DeathOutgoAdjusted + projection.AnnuityOutgoAdjusted + projection.RiderAdjusted + projection.RenewalExpensesAdjusted - projection.PremiumIncomeAdjusted
	}
}
func CalculatePhiReservesAndDiscountedValues(index, valYear int, projection *models.PhiProjections,
	run *models.PhiRunParameters, mp models.PhiModelPoint, params models.PhiParameter,
	projections *[]models.PhiProjections, productShocks []models.PhiShock) error {
	//we are working on the previous projection results, right?
	var forwardRate float64
	var err error
	var resp float64 = 0
	var realResp float64 = 0
	var shockedReal float64 = 0
	var infl float64 = 0

	resp = GetPhiForwardRate(min(index+1, 624), run.YieldCurveYear, run.YieldCurveVersion)
	if err != nil {
		return errors.Wrap(err, "failed to get forward rate")
	}

	if run.ShockSettings.RealYieldCurve {
		infl = GetPhiInflationFactor(min(index+1, 624), run.YieldCurveYear, run.YieldCurveVersion)
		realResp = (1+resp)/(1+infl) - 1
	}

	if run.ShockSettings.NominalYieldCurve {
		var shock models.ProductShock
		if projection.ProjectionMonth > 0 && run.ShockSettings.ShockBasis != "" {
			shock, err = GetShock(projection.ProjectionMonth, run.ShockSettings.ShockBasis, len(productShocks))
			if err != nil {
				log.Error().Err(err).Send()
			}
		}
		forwardRate = math.Max(0, math.Min(1, resp*(1+shock.MultiplicativeYieldCurve)+shock.AdditiveYieldCurve))
	}

	if run.ShockSettings.RealYieldCurve {
		var shock models.ProductShock
		if projection.ProjectionMonth > 0 && run.ShockSettings.ShockBasis != "" {
			shock, err = GetShock(projection.ProjectionMonth, run.ShockSettings.ShockBasis, len(productShocks))
			if err != nil {
				log.Error().Err(err).Send()
			}
		}
		shockedReal = math.Max(0, math.Min(1, realResp*(1+shock.MultiplicativeYieldCurve)+shock.AdditiveYieldCurve))
		forwardRate = (1+shockedReal)*(1+infl) - 1
	}

	if !run.ShockSettings.NominalYieldCurve && !run.ShockSettings.RealYieldCurve {
		forwardRate = resp
	}

	discountFactor := math.Pow(1.0+forwardRate, -1/12.0)
	discountFactorAdjusted := math.Pow(1.0+forwardRate, -1/12.0)

	projection.ValuationRate = forwardRate
	projection.ValuationRateAdjusted = forwardRate //+ margins.InvestmentMargin

	if index == len(*projections)-1 {
		projection.DiscountedPremiumIncome = 0
		projection.DiscountedPremiumIncomeAdjusted = 0
		projection.DiscountedDeathOutgo = 0
		projection.DiscountedDeathOutgoAdjusted = 0
		projection.DiscountedRider = 0
		projection.DiscountedRiderAdjusted = 0
		projection.DiscountedAnnuityOutgo = 0
		projection.DiscountedAnnuityOutgoAdjusted = 0
		projection.DiscountedInitialExpenses = 0
		projection.DiscountedInitialExpensesAdjusted = 0
		projection.DiscountedRenewalExpenses = 0
		projection.DiscountedRenewalExpensesAdjusted = 0
		//projection.DiscountedReinsurancePremium = 0
		//projection.DiscountedReinsurancePremiumAdjusted = 0
		//projection.DiscountedReinsuranceCedingCommission = 0
		//projection.DiscountedReinsuranceCedingCommissionAdjusted = 0
		//projection.DiscountedReinsuranceClaims = 0
		//projection.DiscountedReinsuranceClaimsAdjusted = 0
		projection.DiscountedNetReinsuranceCashflow = 0
		projection.DiscountedNetReinsuranceCashflowAdjusted = 0
		projection.SumCoverageUnits = 0
		projection.DiscountedCoverageUnits = 0
		projection.DiscountedSurrenderOutgo = 0
		projection.DiscountedSurrenderOutgoAdjusted = 0

	} else {
		projection.DiscountedPremiumIncome = utils.FloatPrecision((*projections)[index+1].PremiumIncome+(*projections)[index+1].DiscountedPremiumIncome*discountFactor, defaultPrecision)
		projection.DiscountedPremiumIncomeAdjusted = utils.FloatPrecision((*projections)[index+1].PremiumIncomeAdjusted+(*projections)[index+1].DiscountedPremiumIncomeAdjusted*discountFactorAdjusted, defaultPrecision)

		projection.DiscountedDeathOutgo = utils.FloatPrecision(((*projections)[index+1].DeathOutgo+(*projections)[index+1].DiscountedDeathOutgo)*discountFactor, defaultPrecision)
		projection.DiscountedDeathOutgoAdjusted = utils.FloatPrecision(((*projections)[index+1].DeathOutgoAdjusted+(*projections)[index+1].DiscountedDeathOutgoAdjusted)*discountFactorAdjusted, defaultPrecision)

		projection.DiscountedRider = utils.FloatPrecision(((*projections)[index+1].Rider+(*projections)[index+1].DiscountedRider)*discountFactor, defaultPrecision)
		projection.DiscountedRiderAdjusted = utils.FloatPrecision(((*projections)[index+1].RiderAdjusted+(*projections)[index+1].DiscountedRiderAdjusted)*discountFactorAdjusted, defaultPrecision)

		projection.DiscountedAnnuityOutgo = utils.FloatPrecision(((*projections)[index+1].AnnuityOutgo+(*projections)[index+1].DiscountedAnnuityOutgo)*discountFactor, defaultPrecision)
		projection.DiscountedAnnuityOutgoAdjusted = utils.FloatPrecision(((*projections)[index+1].AnnuityOutgoAdjusted+(*projections)[index+1].DiscountedAnnuityOutgoAdjusted)*discountFactorAdjusted, defaultPrecision)

		projection.DiscountedInitialExpenses = utils.FloatPrecision((*projections)[index+1].InitialExpenses+(*projections)[index+1].DiscountedInitialExpenses*discountFactor, defaultPrecision)
		projection.DiscountedInitialExpensesAdjusted = utils.FloatPrecision((*projections)[index+1].InitialExpensesAdjusted+(*projections)[index+1].DiscountedInitialExpensesAdjusted*discountFactorAdjusted, defaultPrecision)

		projection.DiscountedRenewalExpenses = utils.FloatPrecision((*projections)[index+1].RenewalExpenses+(*projections)[index+1].DiscountedRenewalExpenses*discountFactor, defaultPrecision)
		projection.DiscountedRenewalExpensesAdjusted = utils.FloatPrecision((*projections)[index+1].RenewalExpensesAdjusted+(*projections)[index+1].DiscountedRenewalExpensesAdjusted*discountFactorAdjusted, defaultPrecision)

		//projection.DiscountedReinsurancePremium = utils.FloatPrecision(((*projections)[index+1].ReinsurancePremium+(*projections)[index+1].DiscountedReinsurancePremium)*discountFactor, defaultPrecision)
		//projection.DiscountedReinsurancePremiumAdjusted = utils.FloatPrecision(((*projections)[index+1].ReinsurancePremiumAdjusted+(*projections)[index+1].DiscountedReinsurancePremiumAdjusted)*discountFactor, defaultPrecision)
		//
		//projection.DiscountedReinsuranceCedingCommission = utils.FloatPrecision(((*projections)[index+1].ReinsuranceCedingCommission+(*projections)[index+1].DiscountedReinsuranceCedingCommission)*discountFactor, defaultPrecision)
		//projection.DiscountedReinsuranceCedingCommissionAdjusted = utils.FloatPrecision(((*projections)[index+1].ReinsuranceCedingCommissionAdjusted+(*projections)[index+1].DiscountedReinsuranceCedingCommissionAdjusted)*discountFactor, defaultPrecision)
		//
		//projection.DiscountedReinsuranceClaims = utils.FloatPrecision(((*projections)[index+1].ReinsuranceClaims+(*projections)[index+1].DiscountedReinsuranceClaims)*discountFactor, defaultPrecision)
		//projection.DiscountedReinsuranceClaimsAdjusted = utils.FloatPrecision(((*projections)[index+1].ReinsuranceClaimsAdjusted+(*projections)[index+1].DiscountedReinsuranceClaimsAdjusted)*discountFactor, defaultPrecision)

		projection.DiscountedNetReinsuranceCashflow = utils.FloatPrecision(((*projections)[index+1].NetReinsuranceCashflow+(*projections)[index+1].DiscountedNetReinsuranceCashflow)*discountFactor, defaultPrecision)
		projection.DiscountedNetReinsuranceCashflowAdjusted = utils.FloatPrecision(((*projections)[index+1].NetReinsuranceCashflowAdjusted+(*projections)[index+1].DiscountedNetReinsuranceCashflowAdjusted)*discountFactor, defaultPrecision)

		projection.SumCoverageUnits = utils.FloatPrecision((*projections)[index+1].SumCoverageUnits+(*projections)[index+1].CoverageUnits, defaultPrecision)
		projection.DiscountedCoverageUnits = utils.FloatPrecision(((*projections)[index+1].CoverageUnits+(*projections)[index+1].DiscountedCoverageUnits)*discountFactor, defaultPrecision)

		projection.DiscountedSurrenderOutgo = utils.FloatPrecision(((*projections)[index+1].SurrenderOutgo+(*projections)[index+1].DiscountedSurrenderOutgo)*discountFactor, defaultPrecision)
		projection.DiscountedSurrenderOutgoAdjusted = utils.FloatPrecision(((*projections)[index+1].SurrenderOutgoAdjusted+(*projections)[index+1].DiscountedSurrenderOutgoAdjusted)*discountFactorAdjusted, defaultPrecision)

	}

	projection.Reserves = utils.FloatPrecision(
		projection.DiscountedDeathOutgo+
			projection.DiscountedRider+
			projection.DiscountedAnnuityOutgo+
			projection.DiscountedInitialExpenses+
			projection.DiscountedRenewalExpenses+
			projection.DiscountedMaturityOutgo+
			projection.DiscountedPremiumIncome, defaultPrecision)

	projection.ReservesAdjusted = utils.FloatPrecision(
		projection.DiscountedDeathOutgoAdjusted+
			projection.DiscountedRiderAdjusted+
			projection.DiscountedAnnuityOutgoAdjusted+
			projection.DiscountedInitialExpensesAdjusted+
			projection.DiscountedRenewalExpensesAdjusted+
			projection.DiscountedMaturityOutgoAdjusted+
			projection.DiscountedPremiumIncomeAdjusted, defaultPrecision)

	projection.DiscountedCashOutflow = utils.FloatPrecision(
		projection.DiscountedDeathOutgo+
			projection.DiscountedRider+
			projection.DiscountedAnnuityOutgo+
			projection.DiscountedInitialExpenses+
			projection.DiscountedRenewalExpenses, defaultPrecision)

	projection.DiscountedCashOutflowExclAcquisition = utils.FloatPrecision(
		projection.DiscountedDeathOutgo+
			projection.DiscountedRider+
			projection.DiscountedAnnuityOutgo+
			projection.DiscountedRenewalExpenses+
			projection.DiscountedMaturityOutgo, defaultPrecision)

	projection.DiscountedAcquisitionCost = utils.FloatPrecision(
		projection.DiscountedInitialExpenses, defaultPrecision)

	projection.DiscountedCashInflow = utils.FloatPrecision(
		projection.DiscountedPremiumIncome, defaultPrecision)

	if len(*projections)-1 == index {
		projection.DiscountedProfit = 0
		projection.DiscountedProfitAdjusted = 0
		projection.DiscountedAnnuityFactor = 0

	} else {
		projection.DiscountedProfit = ((*projections)[index+1].Profit + (*projections)[index+1].DiscountedProfit) * discountFactor
		projection.DiscountedProfitAdjusted = ((*projections)[index+1].ProfitAdjusted + (*projections)[index+1].DiscountedProfitAdjusted) * discountFactorAdjusted
		projection.DiscountedAnnuityFactor = projection.DiscountedAnnuityFactor + (*projections)[index+1].DiscountedAnnuityFactor*discountFactor
	}

	//if projection.ProjectionMonth == 0 {
	//	projection.ChangeInReserves = utils.FloatPrecision(projection.Reserves-projection.NetCashFlow, defaultPrecision)
	//	projection.ChangeInReservesAdjusted = utils.FloatPrecision(projection.ReservesAdjusted-projection.NetCashFlowAdjusted, defaultPrecision)
	//
	//} else {
	//	projection.ChangeInReserves = utils.FloatPrecision(futureValues.FutureReserve-projection.Reserves, defaultPrecision)
	//	projection.ChangeInReservesAdjusted = utils.FloatPrecision(futureValues.FutureReserveAdjusted-projection.ReservesAdjusted, defaultPrecision)
	////	projection.InvestmentIncome = utils.FloatPrecision(projection.Reserves*(math.Pow(1.0+forwardRate, 1/12.0)-1), defaultPrecision)
	////	projection.InvestmentIncomeAdjusted = utils.FloatPrecision(projection.ReservesAdjusted*(math.Pow(1.0+forwardRate+(margins.InvestmentMargin/100.0), 1/12.0)-1), defaultPrecision)
	//	//projection.DiscountedPremiumIncome = (projection.)
	//}

	return nil
}

func GetPhiForwardRate(projectionMonth int, year int, version string) float64 {
	var yieldCurve models.PhiYieldCurve
	key := strconv.Itoa(year) + "-y-curve-" + strconv.Itoa(projectionMonth) + "-" + version
	result, found := Cache.Get(key)

	if found {
		return result.(float64)
	} else {
		err := DB.Table("phi_yield_curves").Where("year = ? and proj_time = ? and version=?", year, projectionMonth, version).First(&yieldCurve).Error
		if err != nil {
			fmt.Println("Yield Curve: ", errors.WithStack(err))
		}
		success := Cache.Set(key, yieldCurve.NominalRate, 1)
		//err = redisClient.Set(key, yieldCurve.Inflation, 0).Err()
		if !success {
			fmt.Println("cache: key not stored")
		}
	}
	return yieldCurve.NominalRate
}

type PhiModelPointError struct {
	ModelPoint models.PhiModelPoint
	Error      error
}

func PhiInterestAccrual(index int, projection *models.PhiProjections,
	prevReserve, prevReserveAdjusted float64,
	run *models.PhiRunParameters, mp models.PhiModelPoint,
	params models.PhiParameter, productShocks []models.PhiShock) error {
	var forwardRate float64
	var err error
	var resp float64 = 0
	var infl float64 = 0
	var realResp float64 = 0
	var shockedReal float64 = 0

	if index == 0 {
		projection.InvestmentIncome = 0
		return nil
	} else {
		resp = GetPhiForwardRate(index, run.YieldCurveYear, run.YieldCurveVersion)
		//if err != nil {
		//	return errors.Wrap(err, "failed to get forward rate")
		//}

		if run.ShockSettings.RealYieldCurve {
			infl = GetPhiInflationFactor(index, run.YieldCurveYear, run.YieldCurveVersion)
			//	if err != nil {
			//		return errors.Wrap(err, "failed to get inflation factor")
			//	}
			//	realResp = (1+resp)/(1+infl) - 1
		}

		if run.ShockSettings.NominalYieldCurve {
			var shock models.ProductShock
			if projection.ProjectionMonth > 0 && run.ShockSettings.ShockBasis != "" {
				shock, err = GetShock(projection.ProjectionMonth, run.ShockSettings.ShockBasis, len(productShocks))
				if err != nil {
					log.Error().Err(err).Send()
				}
			}
			forwardRate = math.Max(0, math.Min(1, resp*(1+shock.MultiplicativeYieldCurve)+shock.AdditiveYieldCurve))
		}

		if run.ShockSettings.RealYieldCurve {
			var shock models.ProductShock
			if projection.ProjectionMonth > 0 && run.ShockSettings.ShockBasis != "" {
				shock, err = GetShock(projection.ProjectionMonth, run.ShockSettings.ShockBasis, len(productShocks))
				if err != nil {
					log.Error().Err(err).Send()
				}
			}
			shockedReal = math.Max(0, math.Min(1, realResp*(1+shock.MultiplicativeYieldCurve)+shock.AdditiveYieldCurve))
			forwardRate = (1+shockedReal)*(1+infl) - 1
		}

		if !run.ShockSettings.NominalYieldCurve && !run.ShockSettings.RealYieldCurve {
			forwardRate = resp
		}

		projection.InvestmentIncome = (prevReserve - (projection.RenewalExpenses - (projection.PremiumIncome))) * (math.Pow(1.0+forwardRate, 1.0/12.0) - 1.0)
		projection.InvestmentIncomeAdjusted = (prevReserveAdjusted - (projection.RenewalExpensesAdjusted - (projection.PremiumIncomeAdjusted))) * (math.Pow(1.0+forwardRate, 1.0/12.0) - 1.0)

	}

	return nil
}
func UpdatePhiAggregatedProjections(projections []models.PhiProjections, aggPhiProjs *map[string]models.PhiProjections) {
	mutex.Lock()
	defer mutex.Unlock()
	for i, projection := range projections {
		// TODO: This is a hack. The key should be unique.
		key := strconv.Itoa(i) + "_" + projection.ProductCode + "_" + strconv.Itoa(projection.JobProductID) + "_" + strconv.Itoa(projection.SpCode)
		agg, exists := (*aggPhiProjs)[key]

		if exists {
			mutable := reflect.ValueOf(&agg).Elem()
			mutable2 := reflect.ValueOf(&projection).Elem()

			for j := 0; j < mutable.NumField(); j++ {
				if j == 0 {
					continue
				} else {

					if mutable.Field(j).Type().Kind() == reflect.Float64 {
						mutable.Field(j).SetFloat(mutable.Field(j).Float() + mutable2.Field(j).Float())
					}
				}
			}
			(*aggPhiProjs)[key] = agg
		} else {
			(*aggPhiProjs)[key] = projection
		}

	}
}

func GetPhiRunJobs(user models.AppUser) ([]models.PhiRunParameters, error) {
	var runJobs []models.PhiRunParameters
	var err error

	err = DB.Order("id desc").Find(&runJobs).Error
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve Phi run jobs")
		return nil, errors.Wrap(err, "failed to retrieve Phi run jobs")
	}

	return runJobs, nil
}

func GetPhiRunJob(runID int, user models.AppUser) (map[string]interface{}, error) {
	var runJob models.PhiRunParameters
	var err error
	var result = make(map[string]interface{})

	err = DB.Where("id = ?", runID).First(&runJob).Error
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve Phi run job")
		return result, errors.Wrap(err, "failed to retrieve Phi run job")
	}

	result["run_job"] = runJob

	// get phi aggregated projections for this run job
	var aggProjs []models.PhiAggregatedProjections
	err = DB.Where("run_id = ?", runID).Order("sp_code asc, ifrs17_group asc, projection_month asc").Find(&aggProjs).Error
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve Phi aggregated projections")
	}
	result["aggregated_projections"] = aggProjs

	// get phi projections for this run job
	var projs []models.PhiProjections
	err = DB.Where("run_id = ?", runID).Find(&projs).Error
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve Phi projections")
	}
	result["projections"] = projs

	// get phi scoped aggregations for this run job
	var scopedAggs []models.PhiScopedAggregatedProjections
	err = DB.Where("run_id = ?", runID).Order("sp_code asc, ifrs17_group asc, projection_month asc").Find(&scopedAggs).Error
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve Phi scoped aggregated projections")
	}
	result["scoped_aggregated_projections"] = scopedAggs

	return result, nil
}

func GetPhiRunJobControl(runID int, user models.AppUser) ([]models.PhiProjections, error) {
	var projs []models.PhiProjections
	var err error

	err = DB.Where("run_id = ?", runID).Find(&projs).Error
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve Phi projections")
		return nil, errors.Wrap(err, "failed to retrieve Phi projections")
	}

	return projs, nil

}

func DeletePhiRunJob(runID int, user models.AppUser) error {
	var err error

	// delete projections
	err = DB.Where("run_id = ?", runID).Delete(&models.PhiProjections{}).Error
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete Phi projections")
		return errors.Wrap(err, "failed to delete Phi projections")
	}

	// delete aggregated projections
	err = DB.Where("run_id = ?", runID).Delete(&models.PhiAggregatedProjections{}).Error
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete Phi aggregated projections")
		return errors.Wrap(err, "failed to delete Phi aggregated projections")
	}

	// delete scoped aggregated projections
	err = DB.Where("run_id = ?", runID).Delete(&models.PhiScopedAggregatedProjections{}).Error
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete Phi scoped aggregated projections")
		return errors.Wrap(err, "failed to delete Phi scoped aggregated projections")
	}

	// delete run job
	err = DB.Where("id = ?", runID).Delete(&models.PhiRunParameters{}).Error
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete Phi run job")
		return errors.Wrap(err, "failed to delete Phi run job")
	}

	return nil
}

func DeletePhiValuationTableData(tableType string, year int, version string) error {
	var err error
	err = DB.Table(tableType).Where("year = ? and version = ?", year, version).Delete(nil).Error
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete Phi valuation table data")
	}

	return err
}

func ListPhiRunConfigs() ([]models.PhiRunConfig, error) {
	var configs []models.PhiRunConfig
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Order("created_at DESC").Find(&configs).Error
	})
	return configs, err
}

func GetPhiRunConfig(id int) (models.PhiRunConfig, error) {
	var config models.PhiRunConfig
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.First(&config, id).Error
	})
	return config, err
}

func SavePhiRunConfig(config models.PhiRunConfig) (models.PhiRunConfig, error) {
	result := DB.Create(&config)
	return config, result.Error
}

func DeletePhiRunConfig(id int) error {
	return DB.Delete(&models.PhiRunConfig{}, id).Error
}
