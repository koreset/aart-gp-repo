package services

import (
	appLog "api/log"
	"api/models"
	"api/utils"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"sort"

	"github.com/jszwec/csvutil"
	"github.com/rs/zerolog/log"
)

func CreateConfiguration(configuration models.ExpConfiguration) error {
	err := DB.Create(&configuration).Error
	if err != nil {
		log.Error().Msg(err.Error())
		return errors.New("error creating configuration. The name already exists")
	}
	return nil
}

func GetConfigurations() ([]models.ExpConfiguration, error) {
	var configurations []models.ExpConfiguration
	err := DB.Find(&configurations).Error
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, errors.New("error fetching configurations")
	}

	//get the ActualDataYearVersions and ExpDataYearVersions for each configuration
	for i, configuration := range configurations {
		var actualDataYearVersions []models.ExpActualDataYearVersion
		var actualDataYears []int
		var expDataYearVersions []models.ExpExpDataYearVersion
		var expDataYears []int
		var expCurrentLapseYearVersions []models.ExpCurrentLapseYearVersion
		var expCurrentMortalityYearVersions []models.ExpCurrentMortalityYearVersion

		DB.Where("configuration_id = ?", configuration.ID).Find(&actualDataYearVersions)
		DB.Table("exp_actual_data_year_versions").Where("configuration_id = ?", configuration.ID).Distinct("year").Pluck("year", &actualDataYears)
		DB.Where("configuration_id = ?", configuration.ID).Find(&expDataYearVersions)
		DB.Table("exp_exp_data_year_versions").Where("configuration_id = ?", configuration.ID).Distinct("year").Pluck("year", &expDataYears)
		DB.Where("configuration_id = ?", configuration.ID).Find(&expCurrentLapseYearVersions)
		DB.Where("configuration_id = ?", configuration.ID).Find(&expCurrentMortalityYearVersions)

		configurations[i].ActualDataYearVersions = actualDataYearVersions
		configurations[i].ExpDataYearVersions = expDataYearVersions
		configurations[i].CurrentLapseYearVersions = expCurrentLapseYearVersions
		configurations[i].CurrentMortalityYearVersions = expCurrentMortalityYearVersions
		configurations[i].ActualDataYears = actualDataYears
		configurations[i].ExpDataYears = expDataYears
	}

	return configurations, nil
}

func GetExposureActualYears(expId string) (models.ExpExposureActualYears, error) {
	var exposureActualYears models.ExpExposureActualYears
	var expDataYears []int
	var actualDataYears []int
	var lapseYears []int
	var mortalityYears []int
	// get distinct years for exposure and actual data

	DB.Model(&models.ExpExposureData{}).Where("exp_configuration_id = ?", expId).Distinct("year").Pluck("year", &expDataYears)
	DB.Model(&models.ExpActualData{}).Where("exp_configuration_id = ?", expId).Distinct("year").Pluck("year", &actualDataYears)
	DB.Model(&models.ExpCurrentLapse{}).Distinct("year").Pluck("year", &lapseYears)
	DB.Model(&models.ExpCurrentMortality{}).Distinct("year").Pluck("year", &mortalityYears)

	exposureActualYears = models.ExpExposureActualYears{ExposureYears: expDataYears, ActualYears: actualDataYears, LapseYears: lapseYears, MortalityYears: mortalityYears}
	return exposureActualYears, nil

}

func GetExposureDataVersions(expId string, year string) ([]string, error) {
	var versions []string
	DB.Model(&models.ExpExposureData{}).Where("exp_configuration_id = ? and year = ?", expId, year).Distinct("version").Pluck("version", &versions)
	return versions, nil
}

func GetActualDataVersions(expId string, year string) ([]string, error) {
	var versions []string
	DB.Model(&models.ExpActualData{}).Where("exp_configuration_id = ? and year = ?", expId, year).Pluck("version", &versions)
	return versions, nil
}

func GetLapseDataVersions(year string) ([]string, error) {
	var versions []string
	DB.Model(&models.ExpCurrentLapse{}).Where("year = ?", year).Distinct("version").Pluck("version", &versions)
	return versions, nil
}

func GetMortalityDataVersions(year string) ([]string, error) {
	var versions []string
	DB.Model(&models.ExpCurrentMortality{}).Where("year = ?", year).Distinct("version").Pluck("version", &versions)
	return versions, nil
}

func UploadExpExposureData(v *multipart.FileHeader, year, version string, expId string, user models.AppUser) error {
	var delimiter rune
	delimiterFile, err := v.Open()
	if err != nil {
		return err
	}
	defer delimiterFile.Close()
	delimiter, err = utils.GetDelimiter(delimiterFile)

	file, err := v.Open()
	if err != nil {
		return err
	}

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.Comma = delimiter
	dec, _ := csvutil.NewDecoder(reader)
	dec.Header()

	DB.Where("year = ? and version = ? and exp_configuration_id = ? ", year, version, expId).Delete(&models.ExpExposureData{})
	DB.Where("year = ? and version = ? and configuration_id = ? ", year, version, expId).Delete(&models.ExpExpDataYearVersion{})
	var expmps []models.ExpExposureData
	for {
		var expData models.ExpExposureData
		if err := dec.Decode(&expData); err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}
		intYear := utils.StringToInt(year)
		expData.Year = intYear
		expData.Version = version
		expData.CreatedBy = user.UserName
		expData.ExpConfigurationID = utils.StringToInt(expId)
		if err != nil {
			log.Error().Msg(err.Error())
		}

		expmps = append(expmps, expData)
	}

	// save the data in batches
	err = DB.CreateInBatches(&expmps, 100).Error
	if err != nil {
		log.Error().Msg(err.Error())
	}

	//get the count
	var count int64
	DB.Model(&models.ExpExposureData{}).Where("year = ? and version = ? and exp_configuration_id = ?", year, version, expId).Count(&count)

	// get the configuration name
	var configurationName string
	DB.Model(&models.ExpConfiguration{}).Where("id = ?", expId).Pluck("name", &configurationName)

	// create a record in the year version table
	expExpDataYearVersion := models.ExpExpDataYearVersion{
		ConfigurationName: configurationName,
		ConfigurationId:   utils.StringToInt(expId),
		Year:              utils.StringToInt(year),
		Version:           version,
		Count:             int(count),
	}

	// create a record in the year version table
	err = DB.Create(&expExpDataYearVersion).Error

	return nil
}

func UploadExpActualData(v *multipart.FileHeader, year, version string, expId string, user models.AppUser) error {
	var delimiter rune
	delimiterFile, err := v.Open()
	if err != nil {
		return err
	}
	defer delimiterFile.Close()
	delimiter, err = utils.GetDelimiter(delimiterFile)

	file, err := v.Open()
	if err != nil {
		return err
	}

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.Comma = delimiter
	dec, _ := csvutil.NewDecoder(reader)
	dec.Header()

	DB.Where("year = ? and version = ? and exp_configuration_id = ? ", year, version, expId).Delete(&models.ExpActualData{})
	DB.Where("year = ? and version = ? and configuration_id = ? ", year, version, expId).Delete(&models.ExpActualDataYearVersion{})
	for {
		var expData models.ExpActualData
		if err := dec.Decode(&expData); err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}
		intYear := utils.StringToInt(year)
		expData.Year = intYear
		expData.Version = version
		expData.ExpConfigurationID = utils.StringToInt(expId)
		expData.CreatedBy = user.UserName
		err = DB.Create(&expData).Error
		if err != nil {
			log.Error().Msg(err.Error())
		}
	}

	// get the count
	var count int64
	DB.Model(&models.ExpActualData{}).Where("year = ? and version = ? and exp_configuration_id = ?", year, version, expId).Count(&count)

	// get the configuration name
	var configurationName string
	DB.Model(&models.ExpConfiguration{}).Where("id = ?", expId).Pluck("name", &configurationName)

	// create a record in the year version table
	expActualDataYearVersion := models.ExpActualDataYearVersion{
		ConfigurationName: configurationName,
		ConfigurationId:   utils.StringToInt(expId),
		Year:              utils.StringToInt(year),
		Version:           version,
		Count:             int(count),
	}

	// create a record in the year version table
	err = DB.Create(&expActualDataYearVersion).Error
	if err != nil {
		log.Error().Msg(err.Error())
	}

	return nil
}

func UploadExpCurrentLapse(v *multipart.FileHeader, year, version string, expId string, user models.AppUser) error {
	var delimiter rune
	delimiterFile, err := v.Open()
	if err != nil {
		return err
	}
	defer delimiterFile.Close()
	delimiter, err = utils.GetDelimiter(delimiterFile)

	file, err := v.Open()
	if err != nil {
		return err
	}

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.Comma = delimiter
	dec, _ := csvutil.NewDecoder(reader)
	dec.Header()

	DB.Where("year = ? and version = ? and exp_configuration_id = ? ", year, version, expId).Delete(&models.ExpCurrentLapse{})
	DB.Where("year = ? and version = ? and configuration_id = ? ", year, version, expId).Delete(&models.ExpCurrentLapseYearVersion{})
	for {
		var expCurrentLapse models.ExpCurrentLapse
		if err := dec.Decode(&expCurrentLapse); err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}
		intYear := utils.StringToInt(year)
		expCurrentLapse.Year = intYear
		expCurrentLapse.Version = version
		//expCurrentLapse.ExpConfigurationId = utils.StringToInt(expId)
		expCurrentLapse.CreatedBy = user.UserName
		err = DB.Create(&expCurrentLapse).Error
		if err != nil {
			log.Error().Msg(err.Error())
		}
	}

	// get the count
	var count int64
	DB.Model(&models.ExpCurrentLapse{}).Where("year = ? and version = ? and exp_configuration_id = ?", year, version, expId).Count(&count)

	// get the configuration name
	var configurationName string
	DB.Model(&models.ExpConfiguration{}).Where("id = ?", expId).Pluck("name", &configurationName)

	// create a record in the year version table
	expCurrentLapseYearVersion := models.ExpCurrentLapseYearVersion{
		ConfigurationName: configurationName,
		ConfigurationId:   utils.StringToInt(expId),
		Year:              utils.StringToInt(year),
		Version:           version,
		Count:             int(count),
	}

	// create a record in the year version table
	err = DB.Create(&expCurrentLapseYearVersion).Error
	if err != nil {
		log.Error().Msg(err.Error())
	}

	return nil
}

func UploadExpCurrentMortality(v *multipart.FileHeader, year, version string, expId string, user models.AppUser) error {
	var delimiter rune
	delimiterFile, err := v.Open()
	if err != nil {
		return err
	}
	defer delimiterFile.Close()
	delimiter, err = utils.GetDelimiter(delimiterFile)

	file, err := v.Open()
	if err != nil {
		return err
	}

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.Comma = delimiter
	dec, _ := csvutil.NewDecoder(reader)
	dec.Header()

	DB.Where("year = ? and version = ?", year, version).Delete(&models.ExpCurrentMortality{})
	DB.Where("year = ? and version = ?", year, version).Delete(&models.ExpCurrentMortalityYearVersion{})
	for {
		var expCurrentLapse models.ExpCurrentMortality
		if err := dec.Decode(&expCurrentLapse); err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}
		intYear := utils.StringToInt(year)
		expCurrentLapse.Year = intYear
		expCurrentLapse.Version = version
		//expCurrentLapse.ExpConfigurationId = utils.StringToInt(expId)
		expCurrentLapse.CreatedBy = user.UserName
		err = DB.Create(&expCurrentLapse).Error
		if err != nil {
			log.Error().Msg(err.Error())
		}
	}

	// get the count
	var count int64
	DB.Model(&models.ExpCurrentMortality{}).Where("year = ? and version = ?", year, version).Count(&count)

	// get the configuration name
	var configurationName string
	DB.Model(&models.ExpConfiguration{}).Where("id = ?", expId).Pluck("name", &configurationName)

	// create a record in the year version table
	expCurrentMortalityYearVersion := models.ExpCurrentMortalityYearVersion{
		ConfigurationName: configurationName,
		ConfigurationId:   utils.StringToInt(expId),
		Year:              utils.StringToInt(year),
		Version:           version,
		Count:             int(count),
	}

	// create a record in the year version table
	err = DB.Create(&expCurrentMortalityYearVersion).Error
	if err != nil {
		log.Error().Msg(err.Error())
	}

	return nil
}

func UploadAgeBands(v *multipart.FileHeader, user models.AppUser) error {
	var delimiter rune
	delimiterFile, err := v.Open()
	if err != nil {
		return err
	}
	defer delimiterFile.Close()
	delimiter, err = utils.GetDelimiter(delimiterFile)

	file, err := v.Open()
	if err != nil {
		return err
	}

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.Comma = delimiter
	dec, _ := csvutil.NewDecoder(reader)
	dec.Header()

	DB.Where("id > 0").Delete(&models.ExpAgeBand{})
	for {
		var ageBand models.ExpAgeBand
		if err := dec.Decode(&ageBand); err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}
		err = DB.Create(&ageBand).Error
		if err != nil {
			log.Error().Msg(err.Error())
		}
	}

	return nil

}

func GetAgeBandVersions() ([]string, error) {
	var versions []string
	err := DB.Model(&models.ExpAgeBand{}).Distinct("version").Pluck("version", &versions)
	if err != nil {
		appLog.Error(err.Error)
		//return nil, errors.New("error fetching age band versions")
	}
	return versions, nil
}

func GetTableData(payload models.TableDataPayload) (interface{}, error) {
	var data interface{}
	switch payload.TableName {
	case "exposure_data":
		var expData []models.ExpExposureData
		DB.Where("year = ? and version = ? and exp_configuration_id = ?", payload.Year, payload.Version, payload.ConfigurationId).Limit(2500).Find(&expData)
		data = expData
	case "actual_data":
		var actualData []models.ExpActualData
		DB.Where("year = ? and version = ? and exp_configuration_id = ?", payload.Year, payload.Version, payload.ConfigurationId).Limit(2500).Find(&actualData)
		data = actualData
	case "current_lapse":
		var lapseData []models.ExpCurrentLapse
		DB.Where("year = ? and version = ? and exp_configuration_id = ?", payload.Year, payload.Version, payload.ConfigurationId).Limit(2500).Find(&lapseData)
		data = lapseData
	case "current_mortality":
		var mortalityData []models.ExpCurrentMortality
		DB.Where("year = ? and version = ? and exp_configuration_id = ?", payload.Year, payload.Version, payload.ConfigurationId).Limit(2500).Find(&mortalityData)
		data = mortalityData
	}

	return data, nil
}

func CheckRunName(name string) (bool, error) {
	var count int64
	DB.Model(&models.ExpAnalysisRunSetting{}).Where("run_name = ?", name).Count(&count)
	if count > 0 {
		return true, errors.New("error checking run name")
	}
	return false, nil
}

func GetAnalysisRunSettings(id int) ([]models.ExpAnalysisRunSetting, error) {
	var settings []models.ExpAnalysisRunSetting

	err := DB.Order("id desc").Find(&settings).Error
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, errors.New("error fetching analysis run settings")
	}
	return settings, nil
}

func GetAnalysisRunSettingsById(id int) ([]models.ExpAnalysisRunSetting, error) {
	var settings []models.ExpAnalysisRunSetting

	err := DB.Where("exp_run_group_id = ?", id).Find(&settings).Error
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, errors.New("error fetching analysis run settings")
	}
	return settings, nil
}

func GetAnalysisRunGroups() ([]models.ExpRunGroup, error) {
	var settings []models.ExpRunGroup

	err := DB.Order("id desc").Find(&settings).Error
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, errors.New("error fetching analysis run settings")
	}
	return settings, nil
}

func GetExposureModelPointsByRunGroupId(expRunGroupId int) (map[string]interface{}, error) {
	var results = make(map[string]interface{})
	var tableData = make([]map[string]interface{}, 0)

	var runIds []int

	var expRunGroup models.ExpRunGroup
	err := DB.Where("id = ?", expRunGroupId).First(&expRunGroup).Error

	// get the distinct run ids
	var expRuns []models.ExpAnalysisRunSetting
	err = DB.Where("exp_run_group_id = ?", expRunGroupId).First(&expRuns).Error

	results["settings"] = expRuns

	err = DB.Model(&models.ExposureModelPoint{}).Where("exp_run_group_id = ?", expRunGroupId).Distinct("exp_run_id").Pluck("exp_run_id", &runIds).Error

	//err := DB.Model(&models.ExposureModelPoint{}).Where("exp_run_group_id = ?", expRunGroupId).Pluck("exp_run_id", &runIds).Error
	var modelPoints []models.ExposureModelPoint
	for _, runId := range runIds {
		var expRun models.ExpAnalysisRunSetting
		err = DB.Where("id = ?", runId).First(&expRun).Error
		err := DB.Where("exp_run_id = ?", runId).Limit(8000).Find(&modelPoints).Error

		tableData = append(tableData, map[string]interface{}{
			"data": modelPoints,
			//"table_name": fmt.Sprintf("Model Points - Run ID: %d", runId),
			"table_name": fmt.Sprintf("[ %s to %s ]", expRun.PeriodStartDate, expRun.PeriodEndDate),
		})

		if err != nil {
			log.Error().Msg(err.Error())
			return nil, errors.New("error fetching exposure model points")
		}
	}

	if err != nil {
		log.Error().Msg(err.Error())
		return nil, errors.New("error fetching exposure model points")
	}

	results["table_data"] = tableData
	//return modelPoints, nil
	return results, nil
}

func GetExposureModelPointsByRunIdAndProduct(expRunId, productName string) ([]models.ExposureModelPoint, error) {
	var modelPoints []models.ExposureModelPoint
	err := DB.Where("exp_run_id = ? and product_name = ?", expRunId, productName).Limit(8000).Find(&modelPoints).Error
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, errors.New("error fetching exposure model points")
	}
	return modelPoints, nil
}

func GetExpCrudeResultsByRunGroupId(expRunId int) (map[string]interface{}, error) {
	var results = make(map[string]interface{})
	var tableData = make([]map[string]interface{}, 0)

	var runIds []int

	// get the distinct run ids

	err := DB.Model(&models.ExpCrudeResult{}).Where("exp_run_group_id = ?", expRunId).Distinct("exp_run_id").Pluck("exp_run_id", &runIds).Error
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, errors.New("error fetching crude results")
	}

	for _, runId := range runIds {
		var crudeResults []models.ExpCrudeResult
		err := DB.Where("exp_run_id = ?", runId).Find(&crudeResults).Error
		if err != nil {
			log.Error().Msg(err.Error())
			return nil, errors.New("error fetching crude results")
		}

		var expRun models.ExpAnalysisRunSetting
		err = DB.Where("id = ?", runId).First(&expRun).Error

		tableData = append(tableData, map[string]interface{}{
			"data": crudeResults,
			//"table_name": fmt.Sprintf("Crude Results - Run ID: %d", runId),
			"table_name": fmt.Sprintf("[ %s to %s ]", expRun.PeriodStartDate, expRun.PeriodEndDate),
		})

	}

	results["table_data"] = tableData
	return results, nil
}

func GetExpLapseCrudeResultsByRunGroupId(expRunId int) (map[string]interface{}, error) {
	var results = make(map[string]interface{})
	var tableData = make([]map[string]interface{}, 0)

	var runIds []int

	// get the distinct run ids

	err := DB.Model(&models.ExpLapseCrudeResult{}).Where("exp_run_group_id = ?", expRunId).Distinct("exp_run_id").Pluck("exp_run_id", &runIds).Error
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, errors.New("error fetching lapse crude results")
	}

	for _, runId := range runIds {
		var crudeLapseResults []models.ExpLapseCrudeResult
		err := DB.Where("exp_run_id = ?", runId).Find(&crudeLapseResults).Error
		if err != nil {
			log.Error().Msg(err.Error())
			return nil, errors.New("error fetching lapse crude results")
		}

		var expRun models.ExpAnalysisRunSetting
		err = DB.Where("id = ?", runId).First(&expRun).Error

		tableData = append(tableData, map[string]interface{}{
			"data":       crudeLapseResults,
			"table_name": fmt.Sprintf("[ %s to %s ]", expRun.PeriodStartDate, expRun.PeriodEndDate),
		})

	}
	results["table_data"] = tableData

	summaryData, err := transformDataToAgGrid(expRunId)

	results["summary_data"] = summaryData

	return results, nil
}

func transformDataToAgGrid(expRunGroupId int) ([]OutputRow, error) {
	var totalSummaryResults []models.TotalLapseExpAnalysisResult
	DB.Where("exp_run_group_id = ?", expRunGroupId).Find(&totalSummaryResults)

	aggregates := make(map[string]*AggregationData)
	durations := []string{"yr1", "yr2", "yr3", "yr4", "yr5+"}

	for _, dur := range durations {
		aggregates[dur] = &AggregationData{}
	}

	for _, p := range totalSummaryResults {
		// Year 1
		aggregates["yr1"].Exposure += p.TotalYear1Exposure
		aggregates["yr1"].ActualLapses += p.TotalActualYear1Lapses
		aggregates["yr1"].ExpectedLapses += p.TotalExpectedYear1Lapses

		// Year 2
		aggregates["yr2"].Exposure += p.TotalYear2Exposure
		aggregates["yr2"].ActualLapses += p.TotalActualYear2Lapses
		aggregates["yr2"].ExpectedLapses += p.TotalExpectedYear2Lapses

		// Year 3
		aggregates["yr3"].Exposure += p.TotalYear3Exposure
		aggregates["yr3"].ActualLapses += p.TotalActualYear3Lapses
		aggregates["yr3"].ExpectedLapses += p.TotalExpectedYear3Lapses

		// Year 4
		aggregates["yr4"].Exposure += p.TotalYear4Exposure
		aggregates["yr4"].ActualLapses += p.TotalActualYear4Lapses
		aggregates["yr4"].ExpectedLapses += p.TotalExpectedYear4Lapses

		// Year 5+
		aggregates["yr5+"].Exposure += p.TotalYear5PlusExposure
		aggregates["yr5+"].ActualLapses += p.TotalActualYear5PlusLapses
		aggregates["yr5+"].ExpectedLapses += p.TotalExpectedYear5PlusLapses

	}

	var outputRows []OutputRow

	// Overall totals accumulation
	overall := AggregationData{}

	for _, dur := range durations {
		data := aggregates[dur]
		row := OutputRow{Duration: dur}
		row.Exposure = data.Exposure
		row.ActualLapses = data.ActualLapses
		row.ExpectedLapses = data.ExpectedLapses

		if data.Exposure > 0 {
			row.ActualUX = data.ActualLapses / data.Exposure
			row.ActualMonthlyRate = 1.0 - math.Exp(-row.ActualUX)
			row.ActualAnnualRate = 1.0 - math.Pow(1.0-row.ActualMonthlyRate, 12.0)

			row.ExpectedUX = data.ExpectedLapses / data.Exposure
			row.ExpectedMonthlyRate = 1.0 - math.Exp(-row.ExpectedUX)
			row.ExpectedAnnualRate = 1.0 - math.Pow(1.0-row.ExpectedMonthlyRate, 12.0)

		} else {
			// All calculated fields default to 0 if exposure is 0
		}
		outputRows = append(outputRows, row)

		overall.Exposure += data.Exposure
		overall.ActualLapses += data.ActualLapses
		overall.ExpectedLapses += data.ExpectedLapses
	}

	// Total Row
	totalRow := OutputRow{Duration: "Total"}
	totalRow.Exposure = overall.Exposure
	totalRow.ActualLapses = overall.ActualLapses
	totalRow.ExpectedLapses = overall.ExpectedLapses

	if overall.Exposure > 0 {
		totalRow.ActualUX = overall.ActualLapses / overall.Exposure
		totalRow.ActualMonthlyRate = 1.0 - math.Exp(-totalRow.ActualUX)
		totalRow.ActualAnnualRate = 1.0 - math.Pow(1.0-totalRow.ActualMonthlyRate, 12.0)

		totalRow.ExpectedUX = overall.ExpectedLapses / overall.Exposure
		totalRow.ExpectedMonthlyRate = 1.0 - math.Exp(-totalRow.ExpectedUX)
		totalRow.ExpectedAnnualRate = 1.0 - math.Pow(1.0-totalRow.ExpectedMonthlyRate, 12.0)

	}
	outputRows = append(outputRows, totalRow)

	return outputRows, nil
}

//func calculateLapseCrudeSummaries(expRunGroupId int) []models.TableRowOutput {
//	var totalSummaryResults []models.TotalLapseExpAnalysisResult
//	DB.Where("exp_run_group_id = ?", expRunGroupId).Find(&totalSummaryResults)
//
//	fmt.Println(totalSummaryResults)
//	var results []models.TableRowOutput
//	//
//	for _, record := range totalSummaryResults {
//		results = append(results, transformRecord(record))
//	}
//
//	fmt.Println(results)
//	//
//	return results
//}

// transformRecord processes a single InputRecord and returns its TransformedTableSet
func transformRecord(record models.TotalLapseExpAnalysisResult) []models.TableRowOutput {
	var actualRows []models.TableRowOutput
	//var expectedRows []models.TableRowOutput

	durations := []string{"yr1", "yr2", "yr3", "yr4", "yr5+"}

	// Data from InputRecord, organized for easier iteration
	actualLapsesArr := []float64{record.TotalActualYear1Lapses, record.TotalActualYear2Lapses, record.TotalActualYear3Lapses, record.TotalActualYear4Lapses, record.TotalActualYear5PlusLapses}
	actualUXArr := []float64{record.ActualYear1Ux, record.ActualYear2Ux, record.ActualYear3Ux, record.ActualYear4Ux, record.ActualYear5PlusUx}
	actualAnnualRateArr := []float64{record.ActualYear1AnnualRate, record.ActualYear2AnnualRate, record.ActualYear3AnnualRate, record.ActualYear4AnnualRate, record.ActualYear5PlusAnnualRate}
	actualMonthlyRateArr := []float64{record.ActualYear1MonthlyRate, record.ActualYear2MonthlyRate, record.ActualYear3MonthlyRate, record.ActualYear4MonthlyRate, record.ActualYear5PlusMonthlyRate}

	expectedLapsesArr := []float64{record.TotalExpectedYear1Lapses, record.TotalExpectedYear2Lapses, record.TotalExpectedYear3Lapses, record.TotalExpectedYear4Lapses, record.TotalExpectedYear5PlusLapses}
	expectedUXArr := []float64{record.ExpectedYear1Ux, record.ExpectedYear2Ux, record.ExpectedYear3Ux, record.ExpectedYear4Ux, record.ExpectedYear5PlusUx}
	expectedAnnualRateArr := []float64{record.ExpectedYear1AnnualRate, record.ExpectedYear2AnnualRate, record.ExpectedYear3AnnualRate, record.ExpectedYear4AnnualRate, record.ExpectedYear5PlusAnnualRate}
	expectedMonthlyRateArr := []float64{record.ExpectedYear1MonthlyRate, record.ExpectedYear2MonthlyRate, record.ExpectedYear3MonthlyRate, record.ExpectedYear4MonthlyRate, record.ExpectedYear5PlusMonthlyRate}

	exposureArr := []float64{record.TotalYear1Exposure, record.TotalYear2Exposure, record.TotalYear3Exposure, record.TotalYear4Exposure, record.TotalYear5PlusExposure}

	totalActualLapses := 0.0
	//totalExpectedLapses := 0.0
	totalExposure := 0.0

	for i, duration := range durations {
		currentExposure := exposureArr[i]

		// Actual Data Row
		actualRows = append(actualRows, models.TableRowOutput{
			Duration:                  duration,
			ActualLapses:              actualLapsesArr[i],
			ActualExposure:            currentExposure,
			ActualUX:                  actualUXArr[i],
			ActualAnnualRate:          formatProvidedRate(actualAnnualRateArr[i]),  // Use formatted provided rate
			ActualActualMonthlyRate:   formatProvidedRate(actualMonthlyRateArr[i]), // Use formatted provided rate
			ExpectedLapses:            expectedLapsesArr[i],
			ExpectedExposure:          currentExposure,
			ExpectedUX:                expectedUXArr[i],
			ExpectedAnnualRate:        formatProvidedRate(expectedAnnualRateArr[i]),  // Use formatted provided rate
			ExpectedActualMonthlyRate: formatProvidedRate(expectedMonthlyRateArr[i]), // Use formatted provided rate

		})
		totalActualLapses += actualLapsesArr[i]
		totalExposure += currentExposure

		// Expected Data Row
		//expectedRows = append(expectedRows, models.TableRowOutput{
		//	Duration:          duration,
		//	Lapses:            expectedLapsesArr[i],
		//	Exposure:          currentExposure,
		//	UX:                expectedUXArr[i],
		//	AnnualRate:        formatProvidedRate(expectedAnnualRateArr[i]),  // Use formatted provided rate
		//	ActualMonthlyRate: formatProvidedRate(expectedMonthlyRateArr[i]), // Use formatted provided rate
		//})
		//totalExpectedLapses += expectedLapsesArr[i]
	}

	// Calculate and Add Total Row for Actual Data
	//totalActualUX := 0.0
	//if totalExposure > 0 {
	//	totalActualUX = totalActualLapses / totalExposure
	//} else if totalActualLapses > 0 { // Lapses but no exposure
	//	totalActualUX = math.Inf(1)
	//} else { // 0 lapses and 0 exposure, or just 0 exposure
	//	totalActualUX = math.NaN()
	//}
	//totalActualAnnualRateStr, totalActualMonthlyRateStr := calculateRatesForTotalUX(totalActualUX) // Calculate for total
	//actualRows = append(actualRows, models.TableRowOutput{
	//	Duration:          "Total",
	//	Lapses:            totalActualLapses,
	//	Exposure:          totalExposure,
	//	UX:                totalActualUX,
	//	AnnualRate:        totalActualAnnualRateStr,
	//	ActualMonthlyRate: totalActualMonthlyRateStr,
	//})
	//
	//// Calculate and Add Total Row for Expected Data
	//totalExpectedUX := 0.0
	//if totalExposure > 0 {
	//	totalExpectedUX = totalExpectedLapses / totalExposure
	//} else if totalExpectedLapses > 0 { // Lapses but no exposure
	//	totalExpectedUX = math.Inf(1)
	//} else { // 0 lapses and 0 exposure, or just 0 exposure
	//	totalExpectedUX = math.NaN()
	//}
	//totalExpectedAnnualRateStr, totalExpectedMonthlyRateStr := calculateRatesForTotalUX(totalExpectedUX) // Calculate for total
	//expectedRows = append(expectedRows, models.TableRowOutput{
	//	Duration:          "Total",
	//	Lapses:            totalExpectedLapses,
	//	Exposure:          totalExposure, // Total exposure is the same
	//	UX:                totalExpectedUX,
	//	AnnualRate:        totalExpectedAnnualRateStr,
	//	ActualMonthlyRate: totalExpectedMonthlyRateStr,
	//})

	return actualRows
}

// formatProvidedRate formats a given rate (assumed to be a float like 0.074) into a percentage string.
func formatProvidedRate(rate float64) string {
	if math.IsNaN(rate) || math.IsInf(rate, 0) {
		return "N/A"
	}
	return fmt.Sprintf("%.2f%%", rate*100.0)
}

// calculateRatesForTotalUX calculates and formats Annual and Monthly rates based on a given UX value.
// Used for the "Total" row.
// IMPORTANT: The formula for monthly rate (ux/12) is a simplification.
// Verify if this is the correct calculation for your needs.
func calculateRatesForTotalUX(ux float64) (string, string) {
	if math.IsNaN(ux) || math.IsInf(ux, 0) {
		return "N/A", "N/A"
	}
	annualRateStr := fmt.Sprintf("%.2f%%", ux*100.0)
	// Using simple division for monthly rate. Adjust if a compound formula is needed.
	monthlyRateStr := fmt.Sprintf("%.2f%%", (ux/12.0)*100.0)
	return annualRateStr, monthlyRateStr
}

func CalculateActualsVsExpected(expRunGroupId int) (map[string]interface{}, error) {
	var results = make(map[string]interface{})
	var tableData = make([]map[string]interface{}, 0)

	// get unique run ids associated with the run group
	var runIds []int
	err := DB.Model(&models.ExpCrudeResult{}).Where("exp_run_group_id = ?", expRunGroupId).Distinct("exp_run_id").Pluck("exp_run_id", &runIds).Error

	if err != nil {
		log.Error().Msg(err.Error())
		return nil, errors.New("error fetching run ids")
	}

	var expRuns []models.ExpAnalysisRunSetting
	err = DB.Where("exp_run_group_id = ?", expRunGroupId).First(&expRuns).Error

	results["settings"] = expRuns

	for _, runId := range runIds {
		var crudeResults []models.ExpCrudeResult
		err := DB.Where("exp_run_id = ? and exp_run_group_id = ?", runId, expRunGroupId).Find(&crudeResults).Error
		//crudeResults, err := GetExpCrudeResultsByRunGroupId(expRunGroupId)
		if err != nil {
			return nil, err
		}

		var expRun models.ExpAnalysisRunSetting
		err = DB.Where("id = ?", runId).First(&expRun).Error

		//get all available age bands
		var ageBands []models.ExpAgeBand
		err = DB.Where("version =?", expRun.AgeBandVersion).Find(&ageBands).Error
		if err != nil {
			log.Error().Msg(err.Error())
			return nil, errors.New("error fetching age bands")
		}

		// for each available band, create an actual vs expected record
		// and add to a slice
		var avers []models.ExpActualsVsExpected
		for _, ageBand := range ageBands {
			var aver models.ExpActualsVsExpected
			aver.AgeNext = ageBand.Name
			aver.ExpRunId = expRunGroupId

			// get the crude result for the age band
			for _, crudeResult := range crudeResults {
				if crudeResult.Age >= ageBand.MinAge && crudeResult.Age <= ageBand.MaxAge {
					aver.NumDeathsMale += crudeResult.ActualClaimCountMale
					aver.NumDeathsFemale += crudeResult.ActualClaimCountFemale
					aver.NumDeathsTotal += crudeResult.ActualClaimCountMale + crudeResult.ActualClaimCountFemale
					aver.ExpectedNumDeathsMale += crudeResult.ExpectedClaimCountMale
					aver.ExpectedNumDeathsFemale += crudeResult.ExpectedClaimCountFemale
					aver.ExpectedNumDeathsTotal += crudeResult.ExpectedClaimCountMale + crudeResult.ExpectedClaimCountFemale
					aver.LivesExposureMale += crudeResult.ExposureCountMale
					aver.LivesExposureFemale += crudeResult.ExposureCountFemale
					aver.LivesExposureTotal += crudeResult.ExposureCountMale + crudeResult.ExposureCountFemale

				}
			}

			if aver.LivesExposureMale > 0 {
				aver.CrudeRateMale = 1 - math.Exp(-aver.NumDeathsMale/aver.LivesExposureMale)
				aver.ExpectedRateMale = 1 - math.Exp(-aver.ExpectedNumDeathsMale/aver.LivesExposureMale)
			}
			if aver.LivesExposureFemale > 0 {
				aver.CrudeRateFemale = 1 - math.Exp(-aver.NumDeathsFemale/aver.LivesExposureFemale)
				aver.ExpectedRateFemale = 1 - math.Exp(-aver.ExpectedNumDeathsFemale/aver.LivesExposureFemale)
			}

			if aver.LivesExposureTotal > 0 {
				aver.CrudeRateTotal = 1 - math.Exp(-aver.NumDeathsTotal/aver.LivesExposureTotal)
				aver.ExpectedRateTotal = 1 - math.Exp(-aver.ExpectedNumDeathsTotal/aver.LivesExposureTotal)
			}

			if aver.CrudeRateMale > 0 {
				aver.AvEMale = aver.ExpectedRateMale / aver.CrudeRateMale
			}

			if aver.CrudeRateFemale > 0 {
				aver.AvEFemale = aver.ExpectedRateFemale / aver.CrudeRateFemale
			}

			if aver.CrudeRateTotal > 0 {
				aver.AvETotal = aver.ExpectedRateTotal / aver.CrudeRateTotal
			}

			// add to slice
			avers = append(avers, aver)
		}

		tableData = append(tableData, map[string]interface{}{
			"data":       avers,
			"table_name": fmt.Sprintf("[ %s to %s ]", expRun.PeriodStartDate, expRun.PeriodEndDate),
		})

	}

	results["table_data"] = tableData

	summaryData, err := calculateSummaries(expRunGroupId)
	if err == nil {
		results["summary_data"] = summaryData
		fmt.Println(summaryData)
	} else {
		log.Error().Msg(err.Error())
		return nil, errors.New("error fetching summary data")
	}

	return results, nil

	// get crude results for id

}

func calculateSummaries(expRunGroupId int) (models.APIResponse, error) {
	var totalSummary []models.TotalMortalityExpAnalysisResult
	DB.Where("exp_run_group_id = ?", expRunGroupId).Find(&totalSummary)

	emptyResponse := models.APIResponse{
		ColumnDefs: []models.AGGridColumnDef{
			{HeaderName: "Metric", Field: "metric", Pinned: "left", MinWidth: 250},
			{HeaderName: "Total", Field: "Total", MinWidth: 120, CellStyle: map[string]string{"fontWeight": "bold"}},
		},
		MaleRowData:     []map[string]interface{}{},
		FemaleRowData:   []map[string]interface{}{},
		CombinedRowData: []map[string]interface{}{},
	}

	periodOutputDataList, err := TransformAnalysisResultToTableData(totalSummary)
	if err != nil {
		appLog.Error(err.Error())
		return emptyResponse, errors.New("error transforming analysis result to table data")
	}

	var grandTotalMale, grandTotalFemale, grandTotalCombined models.TableRowData
	for _, periodData := range periodOutputDataList {
		grandTotalMale.LivesExposure += periodData.Male.LivesExposure
		grandTotalMale.ExpectedClaims += periodData.Male.ExpectedClaims
		grandTotalMale.ActualClaims += periodData.Male.ActualClaims
		grandTotalFemale.LivesExposure += periodData.Female.LivesExposure
		grandTotalFemale.ExpectedClaims += periodData.Female.ExpectedClaims
		grandTotalFemale.ActualClaims += periodData.Female.ActualClaims
	}
	recalculateTableRowRates(&grandTotalMale)
	recalculateTableRowRates(&grandTotalFemale)
	grandTotalCombined.LivesExposure = grandTotalMale.LivesExposure + grandTotalFemale.LivesExposure
	grandTotalCombined.ExpectedClaims = grandTotalMale.ExpectedClaims + grandTotalFemale.ExpectedClaims
	grandTotalCombined.ActualClaims = grandTotalMale.ActualClaims + grandTotalFemale.ActualClaims
	recalculateTableRowRates(&grandTotalCombined)

	columnDefs := []models.AGGridColumnDef{}
	metricHeaders := []string{
		"Lives Exposure", "A/E - Lives", "Expected # Claims", "Actual # Claims",
		"Crude Rate - Lives", "Expected Crude mortality rate", "Actual Crude mortality rate", "A/E - Lives Basis",
	}

	columnDefs = append(columnDefs, models.AGGridColumnDef{HeaderName: "Metric", Field: "metric", Pinned: "left", MinWidth: 250})
	periodLabels := []string{}
	// Ensure periodLabels are unique and in order of first appearance from periodOutputDataList
	uniquePeriodLabels := make(map[string]bool)
	for _, pData := range periodOutputDataList {
		if !uniquePeriodLabels[pData.PeriodLabel] {
			periodLabels = append(periodLabels, pData.PeriodLabel)
			columnDefs = append(columnDefs, models.AGGridColumnDef{HeaderName: pData.PeriodLabel, Field: pData.PeriodLabel, MinWidth: 180})
			uniquePeriodLabels[pData.PeriodLabel] = true
		}
	}
	columnDefs = append(columnDefs, models.AGGridColumnDef{
		HeaderName: "Total", Field: "Total", MinWidth: 120, CellStyle: map[string]string{"fontWeight": "bold"},
	})

	for i := range columnDefs {
		isAELives := columnDefs[i].Field == "A/E - Lives"
		isAELivesBasis := columnDefs[i].Field == "A/E - Lives Basis"
		isCrudeRate := columnDefs[i].Field == "Crude Rate - Lives"
		isExpMortRate := columnDefs[i].Field == "Expected Crude mortality rate"
		isActMortRate := columnDefs[i].Field == "Actual Crude mortality rate"
		if isAELives || isAELivesBasis {
			columnDefs[i].ValueFormatter = "params.value == null ? '' : (params.value * 100).toFixed(2) + '%'"
		} else if isCrudeRate {
			columnDefs[i].ValueFormatter = "params.value == null ? '' : params.value.toFixed(5)"
		} else if isExpMortRate || isActMortRate {
			columnDefs[i].ValueFormatter = "params.value == null ? '' : params.value.toFixed(2)"
		} else if columnDefs[i].Field == "Lives Exposure" || columnDefs[i].Field == "Expected # Claims" || columnDefs[i].Field == "Actual # Claims" {
			columnDefs[i].ValueFormatter = "params.value == null ? '' : Math.round(params.value).toLocaleString()"
		}
	}

	maleRows := []map[string]interface{}{}
	femaleRows := []map[string]interface{}{}
	combinedRows := []map[string]interface{}{}

	populateSectionRows := func(
		targetRows *[]map[string]interface{},
		sectionPeriodDataList []models.PeriodOutputData, // Full list to iterate for periods
		sectionGrandTotal models.TableRowData,
		getPeriodMetricData func(p models.PeriodOutputData) models.TableRowData,
	) {
		for _, metricName := range metricHeaders {
			row := map[string]interface{}{"metric": metricName}
			for _, pData := range sectionPeriodDataList { // Iterate through all periods
				pLabel := pData.PeriodLabel // Use the label from the current period iteration
				periodMetricValues := getPeriodMetricData(pData)
				var val interface{}
				switch metricName {
				case "Lives Exposure":
					val = periodMetricValues.LivesExposure
				case "A/E - Lives":
					val = periodMetricValues.AELives
				case "Expected # Claims":
					val = periodMetricValues.ExpectedClaims
				case "Actual # Claims":
					val = periodMetricValues.ActualClaims
				case "Crude Rate - Lives":
					val = periodMetricValues.CrudeRateLives
				case "Expected Crude mortality rate":
					val = periodMetricValues.ExpectedCrudeMortalityRate
				case "Actual Crude mortality rate":
					val = periodMetricValues.ActualCrudeMortalityRate
				case "A/E - Lives Basis":
					val = periodMetricValues.AELivesBasis
				}
				row[pLabel] = val
			}
			var totalVal interface{}
			switch metricName {
			case "Lives Exposure":
				totalVal = sectionGrandTotal.LivesExposure
			case "A/E - Lives":
				totalVal = sectionGrandTotal.AELives
			case "Expected # Claims":
				totalVal = sectionGrandTotal.ExpectedClaims
			case "Actual # Claims":
				totalVal = sectionGrandTotal.ActualClaims
			case "Crude Rate - Lives":
				totalVal = sectionGrandTotal.CrudeRateLives
			case "Expected Crude mortality rate":
				totalVal = sectionGrandTotal.ExpectedCrudeMortalityRate
			case "Actual Crude mortality rate":
				totalVal = sectionGrandTotal.ActualCrudeMortalityRate
			case "A/E - Lives Basis":
				totalVal = sectionGrandTotal.AELivesBasis
			}
			row["Total"] = totalVal
			*targetRows = append(*targetRows, row)
		}
	}

	populateSectionRows(&maleRows, periodOutputDataList, grandTotalMale, func(p models.PeriodOutputData) models.TableRowData { return p.Male })
	populateSectionRows(&femaleRows, periodOutputDataList, grandTotalFemale, func(p models.PeriodOutputData) models.TableRowData { return p.Female })
	populateSectionRows(&combinedRows, periodOutputDataList, grandTotalCombined, func(p models.PeriodOutputData) models.TableRowData { return p.Combined })

	// Build Live Exposure by Gender table (Gender rows x Period columns)
	liveExposureColumnDefs := []models.AGGridColumnDef{{HeaderName: "Gender", Field: "gender", Pinned: "left", MinWidth: 120}}
	liveExposureRowData := []map[string]interface{}{}

	// Prepare Gender Share (percentage distribution) table structures
	genderShareColumnDefs := []models.AGGridColumnDef{{HeaderName: "Gender", Field: "gender", Pinned: "left", MinWidth: 120}}
	genderShareRowData := []map[string]interface{}{}

	// Prepare Age Band Exposure (Age rows x Period columns)
	ageBandColumnDefs := []models.AGGridColumnDef{{HeaderName: "Age", Field: "age", Pinned: "left", MinWidth: 140}}
	ageBandRowData := []map[string]interface{}{}
	// Prepare Age Band Share (percent distribution) table
	ageBandShareColumnDefs := []models.AGGridColumnDef{{HeaderName: "Age", Field: "age", Pinned: "left", MinWidth: 140}}
	ageBandShareRowData := []map[string]interface{}{}

	if len(periodOutputDataList) > 0 {
		// Determine ordered unique period labels (already built above)
		startLabel := periodLabels[0]
		endLabel := periodLabels[len(periodLabels)-1]
		combinedLabel := fmt.Sprintf("%s - %s", toStartYMD(startLabel), toEndYMD(endLabel))

		// Add combined column and each period as columns
		liveExposureColumnDefs = append(liveExposureColumnDefs, models.AGGridColumnDef{HeaderName: combinedLabel, Field: combinedLabel, MinWidth: 120, CellStyle: map[string]string{"fontWeight": "bold"}})
		genderShareColumnDefs = append(genderShareColumnDefs, models.AGGridColumnDef{HeaderName: combinedLabel, Field: combinedLabel, MinWidth: 120, CellStyle: map[string]string{"fontWeight": "bold"}})
		for _, y := range periodLabels {
			liveExposureColumnDefs = append(liveExposureColumnDefs, models.AGGridColumnDef{HeaderName: y, Field: y, MinWidth: 110})
			genderShareColumnDefs = append(genderShareColumnDefs, models.AGGridColumnDef{HeaderName: y, Field: y, MinWidth: 110})
		}
		// Format numbers with thousands separator for live exposure
		for i := range liveExposureColumnDefs {
			if liveExposureColumnDefs[i].Field != "gender" {
				liveExposureColumnDefs[i].ValueFormatter = "params.value == null ? '' : Math.round(params.value).toLocaleString()"
			}
		}
		// Format percentages for gender share
		for i := range genderShareColumnDefs {
			if genderShareColumnDefs[i].Field != "gender" {
				genderShareColumnDefs[i].ValueFormatter = "params.value == null ? '' : (params.value * 100).toFixed(0) + '%'"
			}
		}

		// Prepare Age Band columns (combined + per period)
		ageBandColumnDefs = append(ageBandColumnDefs, models.AGGridColumnDef{HeaderName: combinedLabel, Field: combinedLabel, MinWidth: 140, CellStyle: map[string]string{"fontWeight": "bold"}})
		for _, y := range periodLabels {
			ageBandColumnDefs = append(ageBandColumnDefs, models.AGGridColumnDef{HeaderName: y, Field: y, MinWidth: 120})
		}
		for i := range ageBandColumnDefs {
			if ageBandColumnDefs[i].Field != "age" {
				ageBandColumnDefs[i].ValueFormatter = "params.value == null ? '' : Math.round(params.value).toLocaleString()"
			}
		}

		// Helper to map period label to exposures
		maleRow := map[string]interface{}{"gender": "Male", combinedLabel: grandTotalMale.LivesExposure}
		femaleRow := map[string]interface{}{"gender": "Female", combinedLabel: grandTotalFemale.LivesExposure}
		totalRow := map[string]interface{}{"gender": "Total", combinedLabel: grandTotalCombined.LivesExposure}
		for _, p := range periodOutputDataList {
			maleRow[p.PeriodLabel] = p.Male.LivesExposure
			femaleRow[p.PeriodLabel] = p.Female.LivesExposure
			totalRow[p.PeriodLabel] = p.Combined.LivesExposure
		}
		liveExposureRowData = append(liveExposureRowData, maleRow, femaleRow, totalRow)

		// Build gender share rows as fractions (0..1)
		maleShareRow := map[string]interface{}{"gender": "Male", combinedLabel: toTwoDecimals(safeDivide(grandTotalMale.LivesExposure, grandTotalCombined.LivesExposure) * 100)}
		femaleShareRow := map[string]interface{}{"gender": "Female", combinedLabel: toTwoDecimals(safeDivide(grandTotalFemale.LivesExposure, grandTotalCombined.LivesExposure) * 100)}
		totalShareRow := map[string]interface{}{"gender": "Total", combinedLabel: 100.0}
		for _, p := range periodOutputDataList {
			total := p.Combined.LivesExposure
			maleShareRow[p.PeriodLabel] = toTwoDecimals(safeDivide(p.Male.LivesExposure, total) * 100)
			femaleShareRow[p.PeriodLabel] = toTwoDecimals(safeDivide(p.Female.LivesExposure, total) * 100)
			totalShareRow[p.PeriodLabel] = func() float64 {
				if total == 0 {
					return 0
				} else {
					return 100
				}
			}()
		}
		genderShareRowData = append(genderShareRowData, maleShareRow, femaleShareRow, totalShareRow)

		// Build Age Band exposure table using crude results per run and age band definitions
		// Map period label -> run id from totalSummary
		periodToRun := map[string]int{}
		for _, t := range totalSummary {
			if _, ok := periodToRun[t.Period]; !ok {
				periodToRun[t.Period] = t.ExpRunID
			}
		}
		// Determine base age bands order from the first period's run settings
		var baseBands []models.ExpAgeBand
		if firstRun, ok := periodToRun[periodLabels[0]]; ok {
			var run models.ExpAnalysisRunSetting
			_ = DB.Where("id = ?", firstRun).First(&run).Error
			_ = DB.Where("version =?", run.AgeBandVersion).Find(&baseBands).Error
			// sort by min age for consistent order
			sort.Slice(baseBands, func(i, j int) bool { return baseBands[i].MinAge < baseBands[j].MinAge })
		}
		bandNames := make([]string, 0, len(baseBands))
		for _, b := range baseBands {
			bandNames = append(bandNames, b.Name)
		}
		// Initialize rows map for each band
		rowsByBand := map[string]map[string]interface{}{}
		for _, name := range bandNames {
			row := map[string]interface{}{"age": name, combinedLabel: float64(0)}
			for _, pl := range periodLabels {
				row[pl] = float64(0)
			}
			rowsByBand[name] = row
		}
		// Total row init
		totalAgeRow := map[string]interface{}{"age": "Total", combinedLabel: float64(0)}
		for _, pl := range periodLabels {
			totalAgeRow[pl] = float64(0)
		}
		// Accumulate per period exposures per band
		for _, pl := range periodLabels {
			runID, ok := periodToRun[pl]
			if !ok {
				continue
			}
			var run models.ExpAnalysisRunSetting
			_ = DB.Where("id = ?", runID).First(&run).Error
			var bands []models.ExpAgeBand
			_ = DB.Where("version =?", run.AgeBandVersion).Find(&bands).Error
			sort.Slice(bands, func(i, j int) bool { return bands[i].MinAge < bands[j].MinAge })
			var crude []models.ExpCrudeResult
			_ = DB.Where("exp_run_id = ? and exp_run_group_id = ?", runID, expRunGroupId).Find(&crude).Error
			// For each band, sum exposure male+female
			for _, band := range bands {
				var sum float64
				for _, cr := range crude {
					if cr.Age >= band.MinAge && cr.Age <= band.MaxAge {
						sum += cr.ExposureCountMale + cr.ExposureCountFemale
					}
				}
				// ensure row exists (in case bands differ between periods)
				row, ok := rowsByBand[band.Name]
				if !ok {
					row = map[string]interface{}{"age": band.Name, combinedLabel: float64(0)}
					for _, pl2 := range periodLabels {
						row[pl2] = float64(0)
					}
					rowsByBand[band.Name] = row
					bandNames = append(bandNames, band.Name)
				}
				row[pl] = sum
				row[combinedLabel] = row[combinedLabel].(float64) + sum
				totalAgeRow[pl] = totalAgeRow[pl].(float64) + sum
				totalAgeRow[combinedLabel] = totalAgeRow[combinedLabel].(float64) + sum
			}
		}
		// Assemble ageBandRowData in the established order
		for _, name := range bandNames {
			ageBandRowData = append(ageBandRowData, rowsByBand[name])
		}
		ageBandRowData = append(ageBandRowData, totalAgeRow)

		// Build Age Band Share columns (combined + per-period) and formatters
		ageBandShareColumnDefs = append(ageBandShareColumnDefs, models.AGGridColumnDef{HeaderName: combinedLabel, Field: combinedLabel, MinWidth: 140, CellStyle: map[string]string{"fontWeight": "bold"}})
		for _, y := range periodLabels {
			ageBandShareColumnDefs = append(ageBandShareColumnDefs, models.AGGridColumnDef{HeaderName: y, Field: y, MinWidth: 120})
		}
		for i := range ageBandShareColumnDefs {
			if ageBandShareColumnDefs[i].Field != "age" {
				ageBandShareColumnDefs[i].ValueFormatter = "params.value == null ? '' : (params.value * 100).toFixed(0) + '%'"
			}
		}

		// Compute Age Band shares per period and combined
		for _, name := range bandNames {
			rowExp := rowsByBand[name]
			shareRow := map[string]interface{}{"age": name}
			// combined share
			shareRow[combinedLabel] = toTwoDecimals(safeDivide(rowExp[combinedLabel].(float64), totalAgeRow[combinedLabel].(float64)) * 100)
			// per-period shares
			for _, pl := range periodLabels {
				shareRow[pl] = toTwoDecimals(safeDivide(rowExp[pl].(float64), totalAgeRow[pl].(float64)) * 100)
			}
			ageBandShareRowData = append(ageBandShareRowData, shareRow)
		}
		// Total row for age-band shares
		ageTotalShareRow := map[string]interface{}{"age": "Total"}
		ageTotalShareRow[combinedLabel] = func() float64 {
			if totalAgeRow[combinedLabel].(float64) == 0 {
				return 0
			} else {
				return 100
			}
		}()
		for _, pl := range periodLabels {
			val := totalAgeRow[pl].(float64)
			if val == 0 {
				ageTotalShareRow[pl] = float64(0)
			} else {
				ageTotalShareRow[pl] = float64(100)
			}
		}
		ageBandShareRowData = append(ageBandShareRowData, ageTotalShareRow)
	}

	response := models.APIResponse{
		ColumnDefs:                columnDefs,
		MaleRowData:               maleRows,
		FemaleRowData:             femaleRows,
		CombinedRowData:           combinedRows,
		LiveExposureColumnDefs:    liveExposureColumnDefs,
		LiveExposureRowData:       liveExposureRowData,
		GenderShareColumnDefs:     genderShareColumnDefs,
		GenderShareRowData:        genderShareRowData,
		AgeBandExposureColumnDefs: ageBandColumnDefs,
		AgeBandExposureRowData:    ageBandRowData,
		AgeBandShareColumnDefs:    ageBandShareColumnDefs,
		AgeBandShareRowData:       ageBandShareRowData,
	}

	return response, nil
}

func TransformAnalysisResultToTableData(inputs []models.TotalMortalityExpAnalysisResult) ([]models.PeriodOutputData, error) {
	if len(inputs) == 0 {
		return []models.PeriodOutputData{}, nil
	}
	var outputs []models.PeriodOutputData
	for _, item := range inputs {
		maleData := models.TableRowData{
			LivesExposure:  item.TotalExposureCountMale,
			ExpectedClaims: item.TotalExpectedClaimCountMale,
			ActualClaims:   item.TotalActualClaimCountMale,
		}
		recalculateTableRowRates(&maleData)
		femaleData := models.TableRowData{
			LivesExposure:  item.TotalExposureCountFemale,
			ExpectedClaims: item.TotalExpectedClaimCountFemale,
			ActualClaims:   item.TotalActualClaimCountFemale,
		}
		recalculateTableRowRates(&femaleData)
		combinedData := models.TableRowData{
			LivesExposure:  maleData.LivesExposure + femaleData.LivesExposure,
			ExpectedClaims: maleData.ExpectedClaims + femaleData.ExpectedClaims,
			ActualClaims:   maleData.ActualClaims + femaleData.ActualClaims,
		}
		recalculateTableRowRates(&combinedData)
		outputs = append(outputs, models.PeriodOutputData{
			PeriodLabel: item.Period, Male: maleData, Female: femaleData, Combined: combinedData,
		})
	}
	return outputs, nil
}

func recalculateTableRowRates(data *models.TableRowData) {
	data.AELives = safeDivide(data.ActualClaims, data.ExpectedClaims)
	data.CrudeRateLives = safeDivide(data.ActualClaims, data.LivesExposure)
	data.ExpectedCrudeMortalityRate = safeDivide(data.ExpectedClaims, data.LivesExposure) * 1000
	data.ActualCrudeMortalityRate = data.CrudeRateLives * 1000
	data.AELivesBasis = data.AELives
}

func safeDivide(numerator, denominator float64) float64 {
	if denominator == 0 {
		if numerator == 0 {
			return 0
		}
		if numerator > 0 {
			return math.Inf(1)
		}
		if numerator < 0 {
			return math.Inf(-1)
		}
		return 0
	}
	return numerator / denominator
}

func GetExpTableMetaData() map[string]interface{} {
	// Retrieve model point counts by year, parameters and shocks data?

	var metadata []models.TableMetaData
	var results = make(map[string]interface{})

	// Test additions
	var ageBands = models.TableMetaData{TableType: "Age Bands", Data: nil, Populated: true, TableName: "exp_age_bands"}
	var mortalities = models.TableMetaData{TableType: "Mortalities", Data: nil, Populated: true, TableName: "exp_current_mortalities"}
	var lapses = models.TableMetaData{TableType: "Lapses", Data: nil, Populated: true, TableName: "exp_current_lapses"}
	//var mortalityTable = models.TableMetaData{TableType: "Mortality Table", Data: nil, Populated: true}
	metadata = append(metadata, ageBands, mortalities, lapses)
	results["exp_tables"] = metadata

	return results
}

func SaveExpTables(v *multipart.FileHeader, tableType string, year int, month int, version string) error {
	var delimiter rune
	delimiterFile, err := v.Open()
	if err != nil {
		return err
	}
	defer func(delimiterFile multipart.File) {
		err := delimiterFile.Close()
		if err != nil {
			log.Error().Err(err).Msg("Error closing file")
		}
	}(delimiterFile)
	delimiter, err = utils.GetDelimiter(delimiterFile)

	file, err := v.Open()
	if err != nil {
		return err
	}

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.Comma = delimiter
	dec, _ := csvutil.NewDecoder(reader)
	dec.Header()

	switch tableType {
	case "Lapses":
		DB.Where("year = ? and version = ?", year, version).Delete(&models.ExpCurrentLapse{})
		DB.Where("year = ? and version = ?", year, version).Delete(&models.ExpCurrentLapseYearVersion{})
		for {
			var expCurrentLapse models.ExpCurrentLapse
			if err := dec.Decode(&expCurrentLapse); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			expCurrentLapse.Year = year
			expCurrentLapse.Version = version
			err = DB.Create(&expCurrentLapse).Error
			if err != nil {
				log.Error().Msg(err.Error())
			}
		}
	case "Mortalities":
		DB.Where("year = ? and version = ?", year, version).Delete(&models.ExpCurrentMortality{})
		DB.Where("year = ? and version = ?", year, version).Delete(&models.ExpCurrentMortalityYearVersion{})
		for {
			var expCurrentMortality models.ExpCurrentMortality
			if err := dec.Decode(&expCurrentMortality); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			expCurrentMortality.Year = year
			expCurrentMortality.Version = version
			err = DB.Create(&expCurrentMortality).Error
			if err != nil {
				log.Error().Msg(err.Error())
			}
		}
	case "Age Bands":
		DB.Where("year = ? and version = ?", year, version).Delete(&models.ExpAgeBand{})

		for {
			var ab models.ExpAgeBand

			if err := dec.Decode(&ab); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			ab.Year = year
			ab.Version = version
			err = DB.Create(&ab).Error
			if err != nil {
				log.Error().Msg(err.Error())
			}
		}
	}

	return nil
}

func DeleteExpTableData(tableType, year, version string) error {
	switch tableType {
	case "exp_current_lapses":
		DB.Where("year = ? and version = ?", year, version).Delete(&models.ExpCurrentLapse{})
	case "exp_current_mortalities":
		DB.Where("year = ? and version = ?", year, version).Delete(&models.ExpCurrentMortality{})
	case "exp_age_bands":
		DB.Where("year = ? and version = ?", year, version).Delete(&models.ExpAgeBand{})
	}

	return nil
}

func DeleteExpConfiguration(expId int) error {
	err := DB.Where("id = ?", expId).Delete(&models.ExpConfiguration{}).Error
	if err != nil {
		log.Error().Msg(err.Error())
		return errors.New("error deleting configuration")
	}
	return nil
}

func DeleteExpConfigurationData(tableType, portfolioId, year, version string) error {
	var err error
	switch tableType {
	case "exposure_data":
		err = DB.Where("exp_configuration_id = ? and year = ? and version = ?", portfolioId, year, version).Delete(&models.ExpExposureData{}).Error
		err = DB.Where("configuration_id = ? and year = ? and version = ?", portfolioId, year, version).Delete(&models.ExpExpDataYearVersion{}).Error
	case "actual_data":
		err = DB.Where("exp_configuration_id = ? and year = ? and version = ?", portfolioId, year, version).Delete(&models.ExpActualData{}).Error
		err = DB.Where("configuration_id = ? and year = ? and version = ?", portfolioId, year, version).Delete(&models.ExpActualDataYearVersion{}).Error
	}
	return err
}

func DeleteExpRunGroup(expRunGroupId int) error {
	err := DB.Where("exp_run_group_id = ?", expRunGroupId).Delete(&models.ExpCrudeResult{}).Error
	if err != nil {
		log.Error().Msg(err.Error())
		return errors.New("error deleting crude results")
	}

	err = DB.Where("exp_run_group_id = ?", expRunGroupId).Delete(&models.ExposureModelPoint{}).Error
	if err != nil {
		log.Error().Msg(err.Error())
		return errors.New("error deleting model points")
	}

	err = DB.Where("id = ?", expRunGroupId).Delete(&models.ExpRunGroup{}).Error
	if err != nil {
		log.Error().Msg(err.Error())
		return errors.New("error deleting run group")
	}

	err = DB.Where("exp_run_group_id = ?", expRunGroupId).Delete(&models.ExpAnalysisRunSetting{}).Error
	if err != nil {
		log.Error().Msg(err.Error())
		return errors.New("error deleting ExpAnalysis run setting")
	}

	err = DB.Where("exp_run_group_id = ?", expRunGroupId).Delete(&models.TotalMortalityExpAnalysisResult{}).Error
	if err != nil {
		log.Error().Msg(err.Error())
		return errors.New("error deleting Total Mortality Exp Analysis Results")
	}

	err = DB.Where("exp_run_group_id = ?", expRunGroupId).Delete(&models.TotalLapseExpAnalysisResult{}).Error
	if err != nil {
		log.Error().Msg(err.Error())
		return errors.New("error deleting Total Lapse Exp Analysis Results")
	}

	return nil
}

func toTwoDecimals(val float64) float64 {
	return math.Round(val*100) / 100
}

func toYearMonth(label string) string {
	if len(label) >= 7 {
		return label[:7]
	}
	return label
}

// toStartYMD returns the first date (YYYY-MM-DD) from a range-like label
// e.g., '2019-05-01-2022-05-01' -> '2019-05-01'. If the label is shorter than 10, returns the label as-is.
func toStartYMD(label string) string {
	if len(label) >= 10 {
		return label[:10]
	}
	return label
}

// toEndYMD returns the last date (YYYY-MM-DD) from a range-like label
// e.g., '2019-05-01-2022-05-01' -> '2022-05-01'. If the label is shorter than 10, returns the label as-is.
func toEndYMD(label string) string {
	if len(label) >= 10 {
		return label[len(label)-10:]
	}
	return label
}
