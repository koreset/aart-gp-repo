package services

import (
	"api/models"
	"api/utils"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/jszwec/csvutil"
	"github.com/rs/zerolog/log"
	"io"
	"mime/multipart"
	"sort"
	"strings"
)

func CreateLicPortfolio(licPortfolio models.LicPortfolio) (models.LicPortfolio, error) {
	err := DB.Create(&licPortfolio).Error
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062: Duplicate entry") {
			return licPortfolio, fmt.Errorf("lic portfolio already exists")
		}
	}
	return licPortfolio, err
}

func DeleteIbnrPortfolio(portfolioId int) error {
	err := DB.Where("lic_portfolio_id = ?", portfolioId).Delete(&models.LICClaimsInput{}).Error
	if err != nil {
		return err
	}
	err = DB.Where("ibnr_portfolio_id = ?", portfolioId).Delete(&models.LICEarnedPremium{}).Error
	if err != nil {
		return err
	}
	err = DB.Where("id = ?", portfolioId).Delete(&models.LicPortfolio{}).Error
	return err
}

func GetAvailableInputVersions(portfolioName, year string) ([]models.LicClaimsYearVersion, error) {
	var versions []models.LicClaimsYearVersion
	err := DB.Where("portfolio_name = ? and year = ?", portfolioName, year).Find(&versions).Error
	return versions, err
}

func GetAvailableLicParameterYears() ([]int, error) {
	var years []int
	err := DB.Table("lic2_parameters").Distinct("year").Pluck("year", &years).Error
	return years, err
}

func GetAvailableLicParameterVersions(year string) ([]string, error) {
	var versions []string
	err := DB.Table("lic2_parameters").Where("year = ?", year).Pluck("version", &versions).Error
	return versions, err
}

func GetIbnrPortfolios() ([]models.LicPortfolio, error) {
	var licPortfolios []models.LicPortfolio
	err := DB.Find(&licPortfolios).Error

	for i := range licPortfolios {
		var years []int
		err = DB.Table("lic_claims_inputs").Distinct("year").Where("portfolio_name = ?", licPortfolios[i].Name).Pluck("year", &years).Error
		fmt.Println(years)
		var claimsYearVersions []models.LicClaimsYearVersion
		DB.Where("portfolio_name = ?", licPortfolios[i].Name).Order("year desc").Find(&claimsYearVersions)
		licPortfolios[i].ClaimsYearVersion = claimsYearVersions
		//var policyCounts []models.LicPolicyCount
		//for j := range years {
		//	var count int64
		//	var policyCount models.LicPolicyCount
		//	DB.Table("lic_claims_inputs").Where("portfolio_name = ? and year = ?", licPortfolios[i].Name, years[j]).Count(&count)
		//	policyCount.Year = years[j]
		//	policyCount.Count = int(count)
		//	policyCounts = append(policyCounts, policyCount)
		//}
		//
		//licPortfolios[i].PolicyCount = policyCounts

		var earnedYearVersions []models.LicEarnedPremiumYearVersion
		DB.Where("portfolio_name = ?", licPortfolios[i].Name).Find(&earnedYearVersions)
		licPortfolios[i].EarnedPremiumYearVersion = earnedYearVersions

		//var earnedPremiumCounts []models.LicPolicyCount
		//for j := range years {
		//	var count int64
		//	var policyCount models.LicPolicyCount
		//	DB.Table("lic_earned_premiums").Where("portfolio_name = ? and year = ?", licPortfolios[i].Name, years[j]).Count(&count)
		//	policyCount.Year = years[j]
		//	policyCount.Count = int(count)
		//	earnedPremiumCounts = append(earnedPremiumCounts, policyCount)
		//}
		//
		//licPortfolios[i].EarnedPremiumCount = earnedPremiumCounts

	}
	return licPortfolios, err
}

func GetValidLicPortfolios() ([]models.LicPortfolio, error) {
	var licPortfolios []models.LicPortfolio
	err := DB.Table("lic_claims_inputs").Distinct("portfolio_name as name, lic_portfolio_id as id").Find(&licPortfolios).Error
	return licPortfolios, err
}

func GetLicPortfolioClaimYears(portfolioName string) ([]int, error) {
	var years []int
	err := DB.Table("lic_claims_inputs").Distinct("year").Where("portfolio_name = ?", portfolioName).Pluck("year", &years).Error
	return years, err
}

func GetLicPortfolioParameterYears(portfolioName string) ([]int, error) {
	var years []int
	err := DB.Table("lic_parameters").Distinct("year").Where("portfolio_name = ?", portfolioName).Pluck("year", &years).Error
	return years, err
}

func GetLicBaseVariables() ([]models.LicBaseVariable, error) {
	var licVariables []models.LicBaseVariable
	err := DB.Find(&licVariables).Error
	return licVariables, err
}

func GetLicTableMetaData() map[string]interface{} {
	// Retrieve model point counts by year, parameters and shocks data?

	var metadata []models.TableMetaData
	var results = make(map[string]interface{})

	// Test additions
	var param = models.TableMetaData{TableType: "Parameters", Data: nil, Populated: true}
	metadata = append(metadata, param)
	results["lic_tables"] = metadata

	return results
}

func GetIbnrTableMetaData() map[string]interface{} {
	// Retrieve model point counts by year, parameters and shocks data?

	var metadata []models.TableMetaData
	var results = make(map[string]interface{})

	// Test additions
	var param = models.TableMetaData{TableType: "Parameters", Data: nil, Populated: true}
	var cpi = models.TableMetaData{TableType: "CPI", Data: nil, Populated: true}
	var yieldCurve = models.TableMetaData{TableType: "Yield Curve", Data: nil, Populated: true}
	var shocks = models.TableMetaData{TableType: "Shocks", Data: nil, Populated: true}
	metadata = append(metadata, param, cpi, yieldCurve, shocks)
	results["lic_tables"] = metadata

	return results
}

func DeleteLicTable(tableType string) error {
	var err error
	switch tableType {
	case "Parameters":
		err = DB.Where("id > 0").Delete(&models.Lic2Parameter{}).Error
	}
	return err
}

func DeleteIbnrTable(tableType string) error {
	var err error
	switch tableType {
	case "Parameters":
		err = DB.Where("id > 0").Delete(&models.LICParameter{}).Error
	case "Yield Curve":
		err = DB.Where("id > 0").Delete(&models.YieldCurve{}).Error
	case "CPI":
		err = DB.Where("id > 0").Delete(&models.LicCPI{}).Error
	case "Shocks":
		err = DB.Where("id > 0").Delete(&models.IBNRShock{}).Error
	}
	return err
}

func DeleteIbnrTableData(tableType, portfolio, year, version string) error {
	var err error
	switch tableType {
	case "claims_data":
		err = DB.Where("portfolio_name = ? and year = ? and version_name = ?", portfolio, year, version).Delete(&models.LICClaimsInput{}).Error
		err = DB.Where("portfolio_name = ? and year = ? and version_name = ?", portfolio, year, version).Delete(&models.LicClaimsYearVersion{}).Error
	case "earned_premium":
		err = DB.Where("portfolio_name = ? and year_index = ? and version_name = ?", portfolio, year, version).Delete(&models.LICEarnedPremium{}).Error
		err = DB.Where("portfolio_name = ? and year = ? and version_name = ?", portfolio, year, version).Delete(&models.LicEarnedPremiumYearVersion{}).Error
	}
	return err
}

func SaveLicTables(v *multipart.FileHeader, tableType string, year int, version string) error {
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

	switch tableType {
	case "Parameters":
		for {
			var pp models.Lic2Parameter
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			pp.Year = year
			pp.Version = version
			DB.Where("year = ? and portfolio_name = ? and product_code = ? and version = ?", year, pp.PortfolioName, pp.ProductCode, pp.Version).Delete(&models.Lic2Parameter{})
			err = DB.Create(&pp).Error
			if err != nil {
				log.Error().Msg(err.Error())
			}
		}

	}

	return nil
}

func SaveIbnrTables(v *multipart.FileHeader, tableType string, year int, month int, yieldCurveCode string) error {
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
	case "Parameters":
		for {
			var pp models.LICParameter
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			pp.Year = year
			DB.Where("year = ? and portfolio_name = ? and product_code = ? and underwriting_year=? and underwriting_month=?", year, pp.PortfolioName, pp.ProductCode, pp.UnderwritingYear, pp.UnderwritingMonth).Delete(&models.LICParameter{})
			err = DB.Create(&pp).Error
			if err != nil {
				log.Error().Msg(err.Error())
			}
		}
	case "CPI":
		for {
			var pp models.LicCPI
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			DB.Where("year_index = ? and month_index = ?", pp.YearIndex, pp.MonthIndex).Delete(&models.LicCPI{})
			err = DB.Create(&pp).Error
			if err != nil {
				log.Error().Msg(err.Error())
			}
		}

	case "Shocks":
		for {
			var pp models.IBNRShock
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			DB.Where("shock_basis = ?  and year = ?", pp.ShockBasis, year).Delete(&models.IBNRShock{})
			pp.Year = year
			err = DB.Create(&pp).Error
			if err != nil {
				log.Error().Msg(err.Error())
			}
		}

	case "Yield Curve":
		for {
			var pp models.IbnrYieldCurve
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			pp.Year = year
			pp.Month = month
			pp.YieldCurveCode = yieldCurveCode
			if err != nil {
				log.Error().Msg(err.Error())
			}

			DB.Where("year=? and yield_curve_code = ? and proj_time = ? and month = ?", pp.Year, pp.YieldCurveCode, pp.ProjectionTime, pp.Month).Delete(&models.YieldCurve{})
			err = DB.Create(&pp).Error
			if err != nil {
				log.Error().Msg(err.Error())
			}
		}
	}

	return nil
}

func SaveLicClaimsData(v *multipart.FileHeader, year int, portfolioName string, portfolioId int, tableType string, versionName string) error {
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

	if tableType == "claims_data" {
		DB.Where("year = ? and portfolio_id = ? and version_name = ?", year, portfolioId, versionName).Delete(&models.LicClaimsYearVersion{})
		err = DB.Where("year = ? and lic_portfolio_id = ? and version_name = ?", year, portfolioId, versionName).Delete(&models.LICClaimsInput{}).Error
		if err != nil {
			fmt.Println(err)
		}
		var claims []models.LICClaimsInput
		for {
			var pp models.LICClaimsInput
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			pp.Year = year
			pp.LicPortfolioID = portfolioId
			pp.PortfolioName = portfolioName
			pp.VersionName = versionName
			fmt.Println(pp)
			claims = append(claims, pp)
		}
		var claimsVersion models.LicClaimsYearVersion
		claimsVersion.PortfolioName = portfolioName
		claimsVersion.Year = year
		claimsVersion.VersionName = versionName
		claimsVersion.PortfolioId = portfolioId
		claimsVersion.Count = len(claims)
		DB.Create(&claimsVersion)

		DB.CreateInBatches(&claims, 100)
	}

	if tableType == "earned_premium" {
		DB.Where("year = ? and portfolio_id = ? and version_name = ?", year, portfolioId, versionName).Delete(&models.LicEarnedPremiumYearVersion{})
		err = DB.Where("ibnr_portfolio_id = ? and year = ? and version_name = ?", portfolioId, year, versionName).Delete(&models.LICEarnedPremium{}).Error
		if err != nil {
			fmt.Println(err)
		}
		var claims []models.LICEarnedPremium
		for {
			var pp models.LICEarnedPremium
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			pp.IBNRPortfolioID = portfolioId
			pp.PortfolioName = portfolioName
			pp.VersionName = versionName
			pp.Year = year
			fmt.Println(pp)
			claims = append(claims, pp)
		}

		var epVersion models.LicEarnedPremiumYearVersion
		epVersion.PortfolioName = portfolioName
		epVersion.Year = year
		epVersion.VersionName = versionName
		epVersion.PortfolioId = portfolioId
		epVersion.Count = len(claims)
		DB.Create(&epVersion)

		DB.CreateInBatches(&claims, 100)
	}

	return nil
}

func GetLicClaimsData(portfolioName, tableType, year, version string) (map[string]interface{}, error) {
	var results = map[string]interface{}{}
	var err error
	switch tableType {
	case "claims_data":
		var licClaims []models.LICClaimsInput
		err = DB.Where("portfolio_name = ? and year = ? and version_name = ?", portfolioName, year, version).Find(&licClaims).Error
		results["data"] = licClaims
	case "earned_premium":
		var licClaims []models.LICEarnedPremium
		err = DB.Where("portfolio_name = ? and year = ? and version_name = ?", portfolioName, year, version).Find(&licClaims).Error
		results["data"] = licClaims
	}
	return results, err
}

func GetLicTableData(tableType string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	switch tableType {
	case "parameters":
		var params []models.Lic2Parameter
		DB.Find(&params)
		b, _ := json.Marshal(&params)
		err := json.Unmarshal(b, &results)
		if err != nil {
			return results, err
		}
	}
	return results, nil
}

func GetIbnrTableData(tableType string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	switch tableType {
	case "mortalities":
		var mortalities []models.ExpCurrentMortality
		DB.Find(&mortalities)
		b, _ := json.Marshal(&mortalities)
		err := json.Unmarshal(b, &results)
		if err != nil {
			return results, err
		}
	case "lapses":
		var lapses []models.ExpCurrentLapse
		DB.Find(&lapses)
		b, _ := json.Marshal(&lapses)
		err := json.Unmarshal(b, &results)
		if err != nil {
			return results, err
		}
	case "agebands":
		var ageBands []models.ExpAgeBand
		DB.Find(&ageBands)
		b, _ := json.Marshal(&ageBands)
		err := json.Unmarshal(b, &results)
		if err != nil {
			return results, err
		}
	case "parameters":
		var params []models.LICParameter
		DB.Find(&params)
		b, _ := json.Marshal(&params)
		err := json.Unmarshal(b, &results)
		if err != nil {
			return results, err
		}

	case "yieldcurve":
		var yieldCurve []models.IbnrYieldCurve
		DB.Find(&yieldCurve)
		b, _ := json.Marshal(&yieldCurve)
		err := json.Unmarshal(b, &results)
		if err != nil {
			return results, err
		}
	case "shocks":
		var shocks []models.IBNRShock
		DB.Find(&shocks)
		b, _ := json.Marshal(&shocks)
		err := json.Unmarshal(b, &results)
		if err != nil {
			return results, err
		}
	case "cpi":
		var cpi []models.LicCPI
		DB.Find(&cpi)
		b, _ := json.Marshal(&cpi)
		err := json.Unmarshal(b, &results)
		if err != nil {
			return results, err
		}

	}
	return results, nil
}

func GetAllLICRunJobs() ([]models.IBNRRunSetting, error) {
	var licRunJobs []models.IBNRRunSetting
	err := DB.Order("id desc").Find(&licRunJobs).Error
	return licRunJobs, err
}
func DeleteIBNRRunJob(id int) error {

	//TODO: Delete related tables

	err := DB.Where("run_id = ?", id).Delete(&models.LicModelPoint{}).Error
	err = DB.Where("run_id = ?", id).Delete(&models.LicTriangulation{}).Error
	err = DB.Where("run_id = ?", id).Delete(&models.LicCumulativeTriangulation{}).Error
	err = DB.Where("run_id = ?", id).Delete(&models.LicDevelopmentFactor{}).Error
	err = DB.Where("run_id = ?", id).Delete(&models.LicCumulativeProjection{}).Error
	err = DB.Where("run_id = ?", id).Delete(&models.LicIbnrReserve{}).Error
	err = DB.Where("run_id = ?", id).Delete(&models.IbnrReserveReport{}).Error
	err = DB.Where("run_id = ?", id).Delete(&models.LicTriangulationClaimCount{}).Error
	err = DB.Where("run_id = ?", id).Delete(&models.LicCumulativeTriangulationClaimCount{}).Error
	err = DB.Where("run_id = ?", id).Delete(&models.LicCumulativeTriangulationAverageClaim{}).Error
	err = DB.Where("run_id = ?", id).Delete(&models.LicDevelopmentFactorClaimCount{}).Error
	err = DB.Where("run_id = ?", id).Delete(&models.LicCumulativeProjectionClaimCount{}).Error
	err = DB.Where("run_id = ?", id).Delete(&models.LicCumulativeProjectionAverageClaim{}).Error
	err = DB.Where("run_id = ?", id).Delete(&models.LicCumulativeProjectionAveragetoTotalClaim{}).Error
	err = DB.Where("run_id = ?", id).Delete(&models.LicIncrementalProjectionAveragetoTotalClaim{}).Error
	err = DB.Where("run_id = ?", id).Delete(&models.LicIncrementalInflatedAveragetoTotalClaim{}).Error
	err = DB.Where("run_id = ?", id).Delete(&models.LicDiscountedIncrementalInflatedAveragetoTotalClaim{}).Error
	err = DB.Where("run_id = ?", id).Delete(&models.LicDevelopmentFactorAverageClaim{}).Error
	err = DB.Where("run_id = ?", id).Delete(&models.IbnrFrequency{}).Error
	err = DB.Where("run_id = ?", id).Delete(&models.LicIndividualDevelopmentFactors{}).Error
	err = DB.Where("run_id = ?", id).Delete(&models.LicBiasAdjustmentFactor{}).Error
	err = DB.Where("run_id = ?", id).Delete(&models.LicMackModelCalculatedParameters{}).Error
	err = DB.Where("run_id = ?", id).Delete(&models.LicBiasAdjustedResiduals{}).Error
	err = DB.Where("run_id = ?", id).Delete(&models.LicMeanBiasAdjustedResiduals{}).Error
	err = DB.Where("run_id = ?", id).Delete(&models.LicPseudoRatios{}).Error
	err = DB.Where("run_id = ?", id).Delete(&models.LicMackModelSimulatedDevelopmentFactor{}).Error
	err = DB.Where("run_id = ?", id).Delete(&models.LicMackCumulativeProjection{}).Error
	err = DB.Where("run_id = ?", id).Delete(&models.LicMackSimulationSummaryStats{}).Error
	err = DB.Where("run_id = ?", id).Delete(&models.MackIbnrFrequency{}).Error
	err = DB.Where("id = ?", id).Delete(&models.IBNRRunSetting{}).Error

	return err
}

func GetIBNRRunResults(id int, resultType, product string) (map[string]interface{}, error) {
	var runSettings models.IBNRRunSetting
	DB.Where("id = ?", id).Find(&runSettings)

	var results = make(map[string]interface{})
	var tableData = make([]map[string]interface{}, 0)
	switch resultType {
	case "lic_modelpoints":
		var modelPoints []models.LicModelPoint
		DB.Where("lic_portfolio_id = ? and run_date = ? and run_id = ? and version_name=?", runSettings.PortfolioId, runSettings.RunDate, runSettings.ID, runSettings.ClaimsInputVersion).Find(&modelPoints)

		tableData = append(tableData, map[string]interface{}{
			"data":       modelPoints,
			"table_name": "lic_modelpoints",
		})
	case "data_summaries":
		var claimshistory []models.IBNRClaimHistorySummary
		var licClaimsInput []models.LICClaimsInput
		var query string
		DB.Where("year=? and product_code=? and version_name=?", runSettings.ClaimsDataYear-1, product, runSettings.ClaimsInputVersion).Find(&licClaimsInput)
		if len(licClaimsInput) > 0 {
			query = fmt.Sprintf("SELECT table1.product_code,table1.damage_year as accident_year,sum(table1.claim_amount+table1.assessment_cost) as current_total_incurred_claims,sum(table2.claim_amount+table2.assessment_cost) as previous_total_incurred_claims, sum(table1.paid_claims) as current_total_paid_claims,sum(table2.paid_claims) as previous_total_paid_claims, sum(table1.claim_amount + table1.assessment_cost- table1.paid_claims) as current_total_outstanding_claims,sum(table2.claim_amount + table2.assessment_cost- table2.paid_claims) as previous_total_outstanding_claims FROM lic_claims_inputs as table1 LEFT JOIN lic_claims_inputs as table2 ON table1.damage_year = table2.damage_year where table1.year = %d and table2.year= %d and table1.product_code= '%s' and table2.product_code= '%s' and table1.version_name ='%s' and table2.version_name='%s' group by table1.damage_year, table1.product_code order by table1.damage_year asc", runSettings.ClaimsDataYear, runSettings.ClaimsDataYear-1, product, product, runSettings.ClaimsInputVersion, runSettings.ClaimsInputVersion)
		}
		if len(licClaimsInput) == 0 {
			query = fmt.Sprintf("SELECT product_code,damage_year as accident_year,sum(claim_amount+assessment_cost) as current_total_incurred_claims, sum(paid_claims) as current_total_paid_claims, sum(claim_amount + assessment_cost- paid_claims) as current_total_outstanding_claims FROM lic_claims_inputs where year = %d and product_code= '%s' and version_name ='%s' group by damage_year, product_code order by damage_year asc", runSettings.ClaimsDataYear, product, runSettings.ClaimsInputVersion)
		}
		err := DB.Raw(query).Scan(&claimshistory).Error
		if err != nil {
			fmt.Println(err)
		}
		tableData = append(tableData, map[string]interface{}{
			"data":       claimshistory,
			"table_name": "Claim History Summary",
			"chart_config": map[string]interface{}{
				"chart_type": "bar",
				"x_key":      "accident_year",
				"series": []map[string]interface{}{
					{"key": "current_total_incurred_claims", "label": "Current Incurred", "type": "bar"},
					{"key": "previous_total_incurred_claims", "label": "Prior Year Incurred", "type": "bar"},
					{"key": "current_total_paid_claims", "label": "Current Paid", "type": "bar"},
					{"key": "current_total_outstanding_claims", "label": "Current Outstanding", "type": "bar"},
				},
			},
		})

		var incurredClaims []models.IBNRIncurredClaims
		query = fmt.Sprintf("SELECT product_code,damage_year as accident_year, sum(claim_amount+assessment_cost) as total_claims_incurred, count(id) as claim_count FROM lic_claims_inputs where year = %d and product_code= '%s' and version_name ='%s' group by damage_year, product_code order by damage_year asc", runSettings.ClaimsDataYear, product, runSettings.ClaimsInputVersion)
		err = DB.Raw(query).Scan(&incurredClaims).Error
		if err != nil {
			fmt.Println(err)
		}
		tableData = append(tableData, map[string]interface{}{
			"data":       incurredClaims,
			"table_name": "Incurred Claims",
			"chart_config": map[string]interface{}{
				"chart_type": "dual_axis",
				"x_key":      "accident_year",
				"series": []map[string]interface{}{
					{"key": "total_claims_incurred", "label": "Total Incurred", "type": "bar"},
					{"key": "claim_count", "label": "Claim Count", "type": "line", "axis": "secondary"},
				},
			},
		})

		var paidOutstandingClaims []models.IBNRPaidVsOutstandingClaims
		query = fmt.Sprintf("SELECT product_code,damage_year as accident_year,sum(paid_claims) as total_claims_paid, sum(claim_amount+assessment_cost-paid_claims) as total_claims_outstanding FROM lic_claims_inputs where year = %d and product_code= '%s' and version_name ='%s' group by damage_year, product_code order by damage_year asc", runSettings.ClaimsDataYear, product, runSettings.ClaimsInputVersion)
		err = DB.Raw(query).Scan(&paidOutstandingClaims).Error
		if err != nil {
			fmt.Println(err)
		}
		tableData = append(tableData, map[string]interface{}{
			"data":       paidOutstandingClaims,
			"table_name": "Paid vs Outstanding Claims",
			"chart_config": map[string]interface{}{
				"chart_type": "stacked_bar",
				"x_key":      "accident_year",
				"series": []map[string]interface{}{
					{"key": "total_claims_paid", "label": "Paid Claims", "type": "bar"},
					{"key": "total_claims_outstanding", "label": "Outstanding Claims", "type": "bar"},
				},
			},
		})

		var percentageOutstandingClaims []models.IBNRProportionOutstandingClaims
		query = fmt.Sprintf("SELECT product_code,damage_year as accident_year,case when sum(claim_amount+assessment_cost) = 0 then 0 else sum(claim_amount+assessment_cost-paid_claims)*100/sum(claim_amount+assessment_cost) end as percentage_outstanding_claims FROM lic_claims_inputs where year = %d and product_code= '%s' and version_name ='%s' group by damage_year, product_code order by damage_year asc", runSettings.ClaimsDataYear, product, runSettings.ClaimsInputVersion)
		err = DB.Raw(query).Scan(&percentageOutstandingClaims).Error
		if err != nil {
			fmt.Println(err)
		}
		tableData = append(tableData, map[string]interface{}{
			"data":       percentageOutstandingClaims,
			"table_name": "Percentage Outstanding Claims",
			"chart_config": map[string]interface{}{
				"chart_type": "line",
				"x_key":      "accident_year",
				"series": []map[string]interface{}{
					{"key": "percentage_outstanding_claims", "label": "% Outstanding", "type": "line"},
				},
			},
		})

		var averageClaimAmount []models.IBNRAverageClaimAmount
		query = fmt.Sprintf("SELECT product_code,damage_year as accident_year,case when count(id) =0 then 0 else sum(claim_amount+assessment_cost)/count(id) end as average_claim_amount FROM lic_claims_inputs where year = %d and product_code= '%s' and version_name ='%s' group by damage_year, product_code order by damage_year asc", runSettings.ClaimsDataYear, product, runSettings.ClaimsInputVersion)
		err = DB.Raw(query).Scan(&averageClaimAmount).Error
		if err != nil {
			fmt.Println(err)
		}
		tableData = append(tableData, map[string]interface{}{
			"data":       averageClaimAmount,
			"table_name": "Average Claim Amount",
			"chart_config": map[string]interface{}{
				"chart_type": "line",
				"x_key":      "accident_year",
				"series": []map[string]interface{}{
					{"key": "average_claim_amount", "label": "Average Claim", "type": "line"},
				},
			},
		})

		// ── NEW SUMMARY 1: Loss Ratio by Accident Year (CAS / BF prerequisite) ──────────
		// Earned premium per accident year (year_index = accident year in lic_earned_premiums)
		var earnedPremiums []models.LICEarnedPremium
		DB.Where("portfolio_name = ? AND product_code = ? AND year = ? AND version_name = ?",
			runSettings.PortfolioName, product, runSettings.ClaimsDataYear, runSettings.ClaimsInputVersion).
			Find(&earnedPremiums)

		// Build accident-year → earned-premium map (sum monthly premiums)
		epByAY := map[int]float64{}
		for _, ep := range earnedPremiums {
			epByAY[ep.YearIndex] += ep.EarnedPremium
		}

		// Expected loss ratio from parameters
		var licParam models.LICParameter
		DB.Where("portfolio_name = ? AND product_code = ? AND year = ?",
			runSettings.PortfolioName, product, runSettings.ParameterYear).
			First(&licParam)

		// Per-accident-year IBNR reserves (aggregate by AY over months)
		var ibnrReserveRows []models.LicIbnrReserve
		DB.Where("lic_portfolio_id = ? AND run_id = ? AND product_code = ?",
			runSettings.PortfolioId, id, product).
			Find(&ibnrReserveRows)

		// Build AY maps from reserves
		type ayAccum struct {
			ibnrBel           float64
			paidToDate        float64
			outstanding       float64
			actualClaims      float64
			predictedTotLoss  float64
			proportionRunOff  float64
			effectiveMethod   string
			monthCount        int
		}
		ayMap := map[int]*ayAccum{}
		for _, r := range ibnrReserveRows {
			if _, ok := ayMap[r.AccidentYear]; !ok {
				ayMap[r.AccidentYear] = &ayAccum{}
			}
			a := ayMap[r.AccidentYear]
			a.ibnrBel += r.CombinedClBfIbnr
			a.actualClaims += r.ActualClaims
			a.predictedTotLoss += r.PredictedTotalLoss
			a.proportionRunOff += r.ProportionRunoff
			a.monthCount++
			if r.EffectiveMethod != "" {
				a.effectiveMethod = r.EffectiveMethod
			}
		}

		// Build paid/outstanding per AY from paidOutstandingClaims (already queried above)
		paidByAY := map[int]float64{}
		outsByAY := map[int]float64{}
		for _, po := range paidOutstandingClaims {
			paidByAY[po.AccidentYear] = po.TotalClaimsPaid
			outsByAY[po.AccidentYear] = po.TotalClaimsOutstanding
		}

		// Build incurred per AY and claim count maps
		incurredByAY := map[int]float64{}
		countByAY := map[int]int{}
		for _, ic := range incurredClaims {
			incurredByAY[ic.AccidentYear] = ic.TotalClaimsIncurred
			countByAY[ic.AccidentYear] = ic.ClaimCount
		}

		// Collect sorted AYs
		aySet := map[int]bool{}
		for ay := range epByAY { aySet[ay] = true }
		for ay := range ayMap { aySet[ay] = true }
		for ay := range incurredByAY { aySet[ay] = true }
		sortedAYs := make([]int, 0, len(aySet))
		for ay := range aySet { sortedAYs = append(sortedAYs, ay) }
		sort.Ints(sortedAYs)

		var lossRatioSummary []models.IBNRLossRatioSummary
		for _, ay := range sortedAYs {
			ep := epByAY[ay]
			inc := incurredByAY[ay]
			ibnr := 0.0
			if a, ok := ayMap[ay]; ok {
				ibnr = a.ibnrBel
			}
			actualLR, ultimateLR := 0.0, 0.0
			if ep != 0 {
				actualLR = inc / ep
				ultimateLR = (inc + ibnr) / ep
			}
			lossRatioSummary = append(lossRatioSummary, models.IBNRLossRatioSummary{
				AccidentYear:      ay,
				EarnedPremium:     ep,
				TotalIncurred:     inc,
				ActualLossRatio:   actualLR,
				ExpectedLossRatio: licParam.ExpectedLossRatio,
				UltimateLossRatio: ultimateLR,
			})
		}
		tableData = append(tableData, map[string]interface{}{
			"data":       lossRatioSummary,
			"table_name": "Loss Ratio by Accident Year",
			"chart_config": map[string]interface{}{
				"chart_type": "dual_axis",
				"x_key":      "accident_year",
				"series": []map[string]interface{}{
					{"key": "total_incurred", "label": "Total Incurred", "type": "bar"},
					{"key": "earned_premium", "label": "Earned Premium", "type": "bar"},
					{"key": "actual_loss_ratio", "label": "Actual Loss Ratio", "type": "line", "axis": "secondary"},
					{"key": "expected_loss_ratio", "label": "Expected Loss Ratio", "type": "line", "axis": "secondary"},
					{"key": "ultimate_loss_ratio", "label": "Ultimate Loss Ratio", "type": "line", "axis": "secondary"},
				},
				"reference_lines": []map[string]interface{}{
					{"value": licParam.ExpectedLossRatio, "label": "Expected LR"},
				},
			},
		})

		// ── NEW SUMMARY 2: Ultimate Loss Reconciliation (IFoA GIRO) ─────────────────────
		var ultimateLossSummary []models.IBNRUltimateLossSummary
		for _, ay := range sortedAYs {
			paid := paidByAY[ay]
			outs := outsByAY[ay]
			ibnr := 0.0
			propRunOff := 0.0
			effMethod := ""
			if a, ok := ayMap[ay]; ok {
				ibnr = a.ibnrBel
				if a.monthCount > 0 {
					propRunOff = a.proportionRunOff / float64(a.monthCount)
				}
				effMethod = a.effectiveMethod
			}
			ultimate := paid + outs + ibnr
			pctPaid := 0.0
			if ultimate != 0 {
				pctPaid = paid / ultimate * 100
			}
			ultimateLossSummary = append(ultimateLossSummary, models.IBNRUltimateLossSummary{
				AccidentYear:     ay,
				PaidToDate:       paid,
				OutstandingAtVal: outs,
				IbnrBestEstimate: ibnr,
				UltimateLoss:     ultimate,
				PctPaidToUltimate: pctPaid,
				ProportionRunOff: propRunOff,
				EffectiveMethod:  effMethod,
			})
		}
		tableData = append(tableData, map[string]interface{}{
			"data":       ultimateLossSummary,
			"table_name": "Ultimate Loss Reconciliation",
			"chart_config": map[string]interface{}{
				"chart_type": "stacked_bar",
				"x_key":      "accident_year",
				"series": []map[string]interface{}{
					{"key": "paid_to_date", "label": "Paid to Date", "type": "bar"},
					{"key": "outstanding_at_valuation", "label": "Outstanding", "type": "bar"},
					{"key": "ibnr_best_estimate", "label": "IBNR Best Estimate", "type": "bar"},
					{"key": "pct_paid_to_ultimate", "label": "% Paid to Ultimate", "type": "line", "axis": "secondary"},
				},
			},
		})

		// ── NEW SUMMARY 3: Claim Frequency & Severity Trend (Munich / Clark LDF) ────────
		var freqSevSummary []models.IBNRFrequencySeveritySummary
		for _, ay := range sortedAYs {
			ep := epByAY[ay]
			cnt := countByAY[ay]
			inc := incurredByAY[ay]
			freq, sev := 0.0, 0.0
			if ep != 0 {
				freq = float64(cnt) / ep * 1000
			}
			if cnt != 0 {
				sev = inc / float64(cnt)
			}
			freqSevSummary = append(freqSevSummary, models.IBNRFrequencySeveritySummary{
				AccidentYear:   ay,
				EarnedPremium:  ep,
				ClaimCount:     cnt,
				TotalIncurred:  inc,
				ClaimFrequency: freq,
				ClaimSeverity:  sev,
			})
		}
		tableData = append(tableData, map[string]interface{}{
			"data":       freqSevSummary,
			"table_name": "Claim Frequency & Severity Trend",
			"chart_config": map[string]interface{}{
				"chart_type": "dual_axis",
				"x_key":      "accident_year",
				"series": []map[string]interface{}{
					{"key": "claim_frequency", "label": "Frequency (per 1000 premium)", "type": "line"},
					{"key": "claim_severity", "label": "Severity (avg cost)", "type": "line", "axis": "secondary"},
				},
			},
		})

		// ── NEW SUMMARY 4: Development Maturity (ASISA SAP 5) ───────────────────────────
		var devMaturitySummary []models.IBNRDevelopmentMaturitySummary
		for _, ay := range sortedAYs {
			propRunOff := 0.0
			effMethod := ""
			if a, ok := ayMap[ay]; ok {
				if a.monthCount > 0 {
					propRunOff = a.proportionRunOff / float64(a.monthCount)
				}
				effMethod = a.effectiveMethod
			}
			devMaturitySummary = append(devMaturitySummary, models.IBNRDevelopmentMaturitySummary{
				AccidentYear:        ay,
				ProportionRunOff:    propRunOff,
				ProportionNotRunOff: 1 - propRunOff,
				EffectiveMethod:     effMethod,
				ChainLadderEligible: propRunOff >= 0.8,
			})
		}
		tableData = append(tableData, map[string]interface{}{
			"data":       devMaturitySummary,
			"table_name": "Development Maturity",
			"chart_config": map[string]interface{}{
				"chart_type": "stacked_bar",
				"x_key":      "accident_year",
				"series": []map[string]interface{}{
					{"key": "proportion_run_off", "label": "% Developed (Run-off)", "type": "bar"},
					{"key": "proportion_not_run_off", "label": "% Undeveloped", "type": "bar"},
				},
			},
		})

		// ── NEW SUMMARY 5: Actual vs Expected Loss Diagnostics (ASISA SAP 5) ────────────
		var reserveAdequacySummary []models.IBNRReserveAdequacySummary
		for _, ay := range sortedAYs {
			if a, ok := ayMap[ay]; ok {
				aeRatio := 0.0
				if a.predictedTotLoss != 0 {
					aeRatio = a.actualClaims / a.predictedTotLoss
				}
				propRunOff := 0.0
				if a.monthCount > 0 {
					propRunOff = a.proportionRunOff / float64(a.monthCount)
				}
				reserveAdequacySummary = append(reserveAdequacySummary, models.IBNRReserveAdequacySummary{
					AccidentYear:          ay,
					ActualClaims:          a.actualClaims,
					PredictedTotalLoss:    a.predictedTotLoss,
					ActualVsExpectedRatio: aeRatio,
					PredictedLossRatio:    licParam.ExpectedLossRatio,
					CredibilityWeightCL:   propRunOff,
					CredibilityWeightBF:   1 - propRunOff,
				})
			}
		}
		tableData = append(tableData, map[string]interface{}{
			"data":       reserveAdequacySummary,
			"table_name": "Actual vs Expected Loss Diagnostics",
			"chart_config": map[string]interface{}{
				"chart_type": "line",
				"x_key":      "accident_year",
				"series": []map[string]interface{}{
					{"key": "actual_vs_expected_ratio", "label": "Actual / Expected Ratio", "type": "line"},
				},
				"reference_lines": []map[string]interface{}{
					{"value": 1.0, "label": "Break-even (A/E = 1.0)"},
				},
			},
		})

		// ── NEW SUMMARY 6: IBNR Method Sensitivity (actuarial report) ───────────────────
		var reserveReport models.IbnrReserveReport
		DB.Where("lic_portfolio_id = ? AND run_id = ? AND product_code = ?",
			runSettings.PortfolioId, id, product).
			First(&reserveReport)

		pctDev := func(method, bel float64) float64 {
			if bel == 0 {
				return 0
			}
			return (method - bel) / bel * 100
		}
		methodSensitivity := models.IBNRMethodSensitivitySummary{
			ProductCode:       product,
			IbnrBel:           reserveReport.IbnrBel,
			ChainLadderIbnr:   reserveReport.ChainLadderIbnr,
			ClVsBelPct:        pctDev(reserveReport.ChainLadderIbnr, reserveReport.IbnrBel),
			BfIbnr:            reserveReport.BfIbnr,
			BfVsBelPct:        pctDev(reserveReport.BfIbnr, reserveReport.IbnrBel),
			CapeCodIbnr:       reserveReport.CapeCodIbnr,
			CcVsBelPct:        pctDev(reserveReport.CapeCodIbnr, reserveReport.IbnrBel),
			CombinedClBfIbnr:  reserveReport.CombinedClBfIbnr,
			CombinedClCcIbnr:  reserveReport.CombinedClCapeCodIbnr,
			BootstrapIbnr:     reserveReport.BootstrapIbnr,
			BootstrapVsBelPct: pctDev(reserveReport.BootstrapIbnr, reserveReport.IbnrBel),
			MackModelIbnr:     reserveReport.MackModelIbnr,
			MackVsBelPct:      pctDev(reserveReport.MackModelIbnr, reserveReport.IbnrBel),
		}
		tableData = append(tableData, map[string]interface{}{
			"data":       []models.IBNRMethodSensitivitySummary{methodSensitivity},
			"table_name": "IBNR Method Sensitivity",
			"chart_config": map[string]interface{}{
				"chart_type": "tornado",
				"x_key":      "method",
				"series": []map[string]interface{}{
					{"key": "cl_vs_bel_pct", "label": "Chain Ladder"},
					{"key": "bf_vs_bel_pct", "label": "Bornhuetter-Ferguson"},
					{"key": "cc_vs_bel_pct", "label": "Cape Cod"},
					{"key": "bootstrap_vs_bel_pct", "label": "Bootstrap"},
					{"key": "mack_vs_bel_pct", "label": "Mack Model"},
				},
				"reference_lines": []map[string]interface{}{
					{"value": 0, "label": "BEL (0% deviation)"},
				},
			},
		})

		// ── NEW SUMMARY 7: Stochastic vs Deterministic (IFRS 17 para 37) ────────────────
		var bootstrapStats models.LicBootstrappedResultSummary
		DB.Where("lic_portfolio_id = ? AND run_id = ? AND product_code = ?",
			runSettings.PortfolioId, id, product).
			First(&bootstrapStats)

		var mackStats models.LicMackSimulationSummaryStats
		DB.Where("lic_portfolio_id = ? AND run_id = ? AND product_code = ?",
			runSettings.PortfolioId, id, product).
			First(&mackStats)

		bootstrapCV := 0.0
		if bootstrapStats.Mean != 0 {
			bootstrapCV = bootstrapStats.StandardDeviation / bootstrapStats.Mean
		}
		mackCV := 0.0
		if mackStats.Mean != 0 {
			mackCV = mackStats.StandardDeviation / mackStats.Mean
		}
		stochasticSummary := models.IBNRStochasticSummary{
			ProductCode:              product,
			IbnrBel:                  reserveReport.IbnrBel,
			BootstrapMean:            bootstrapStats.Mean,
			BootstrapStdDev:          bootstrapStats.StandardDeviation,
			BootstrapCV:              bootstrapCV,
			BootstrapNthPercentile:   reserveReport.BootstrapNthPercentileIbnr,
			MackMean:                 mackStats.Mean,
			MackStdDev:               mackStats.StandardDeviation,
			MackCV:                   mackCV,
			MackNthPercentile:        reserveReport.MackModelNthPercentileIbnr,
			StochasticAdequacyMargin: reserveReport.BootstrapNthPercentileIbnr - reserveReport.IbnrBel,
		}
		tableData = append(tableData, map[string]interface{}{
			"data":       []models.IBNRStochasticSummary{stochasticSummary},
			"table_name": "Stochastic vs Deterministic Reserve Summary",
			"chart_config": map[string]interface{}{
				"chart_type": "box_plot",
				"series": []map[string]interface{}{
					{
						"label":  "Bootstrap",
						"min":    bootstrapStats.Minimum,
						"mean":   bootstrapStats.Mean,
						"median": bootstrapStats.Median,
						"max":    bootstrapStats.Maximum,
						"pct":    reserveReport.BootstrapNthPercentileIbnr,
					},
					{
						"label":  "Mack Model",
						"min":    mackStats.Minimum,
						"mean":   mackStats.Mean,
						"median": mackStats.Median,
						"max":    mackStats.Maximum,
						"pct":    reserveReport.MackModelNthPercentileIbnr,
					},
				},
				"reference_lines": []map[string]interface{}{
					{"value": reserveReport.IbnrBel, "label": "BEL"},
				},
			},
		})

	case "intermediary_base_results":
		//triangulations
		var licTriangulations []models.LicTriangulation
		DB.Where("lic_portfolio_id = ? and run_date = ? and product_code = ? and run_id = ?", runSettings.PortfolioId, runSettings.RunDate, product, id).Find(&licTriangulations)
		tableData = append(tableData, map[string]interface{}{
			"data":       licTriangulations,
			"table_name": "Triangulations",
		})

		// cumulative triangulations
		var cts []models.LicCumulativeTriangulation
		DB.Where("lic_portfolio_id = ? and run_date = ? and product_code = ? and run_id = ?", runSettings.PortfolioId, runSettings.RunDate, product, id).Find(&cts)
		tableData = append(tableData, map[string]interface{}{
			"data":       cts,
			"table_name": "Cumulative Triangulations",
		})

		// development factors
		var dfs []models.LicDevelopmentFactor
		var dfsGraphData models.LicDevelopmentFactorGraphData
		DB.Where("portfolio_name = ? and run_date = ? and product_code = ? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&dfs)

		for _, df := range dfs {
			if df.DevelopmentVariable == "Weighted Ave. Succession Ratio(rd to rd+1)" {
				err := copier.Copy(&dfsGraphData, &df)
				if err != nil {
					fmt.Println(err)
				}
			}
		}

		tableData = append(tableData, map[string]interface{}{
			"data":       dfs,
			"table_name": "Development Factors",
			"graphData":  dfsGraphData,
		})

		//cumulative projections
		var fcp []models.LicCumulativeProjection
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&fcp)

		tableData = append(tableData, map[string]interface{}{
			"data":       fcp,
			"table_name": "Cumulative Projections",
		})

		//incremental projections
		var licIncrementalProjection []models.LicIncrementalProjection
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&licIncrementalProjection)

		tableData = append(tableData, map[string]interface{}{
			"data":       licIncrementalProjection,
			"table_name": "Incremental Projections",
		})

		//incremental inflated projections
		var licIncrementalInflated []models.LicIncrementalInflated
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&licIncrementalInflated)

		tableData = append(tableData, map[string]interface{}{
			"data":       licIncrementalInflated,
			"table_name": "Incremental Inflated",
		})

		//discounted incremental inflated projections
		var licDiscountedIncrementalInflated []models.LicDiscountedIncrementalInflated
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&licDiscountedIncrementalInflated)

		tableData = append(tableData, map[string]interface{}{
			"data":       licDiscountedIncrementalInflated,
			"table_name": "Discounted Incremental Inflated",
		})
	case "intermediary_average_cost_per_claim_results":
		// triangulation claims count
		var licTriangulationClaimsCount []models.LicTriangulationClaimCount
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&licTriangulationClaimsCount)

		tableData = append(tableData, map[string]interface{}{
			"data":       licTriangulationClaimsCount,
			"table_name": "Triangulation Claims Count",
		})

		// cumulative triangulation claims count
		var lctcc []models.LicCumulativeTriangulationClaimCount
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&lctcc)

		tableData = append(tableData, map[string]interface{}{
			"data":       lctcc,
			"table_name": "Cumulative Triangulation Claim Count",
		})

		// cumulative triangulation average claim
		var licCumulativeTriangulationAverageClaim []models.LicCumulativeTriangulationAverageClaim
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&licCumulativeTriangulationAverageClaim)

		tableData = append(tableData, map[string]interface{}{
			"data":       licCumulativeTriangulationAverageClaim,
			"table_name": "Cumulative Triangulation Average Claim",
		})

		//development factor claim
		var licDevelopmentFactorClaimCount []models.LicDevelopmentFactorClaimCount
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&licDevelopmentFactorClaimCount)

		tableData = append(tableData, map[string]interface{}{
			"data":       licDevelopmentFactorClaimCount,
			"table_name": "Development Factor Claim Count",
		})

		//development factor average claim
		var licDevelopmentFactorAverageClaim []models.LicDevelopmentFactorAverageClaim
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&licDevelopmentFactorAverageClaim)

		tableData = append(tableData, map[string]interface{}{
			"data":       licDevelopmentFactorAverageClaim,
			"table_name": "Development Factor Average Claim",
		})

		// cumulative projecttion claim count
		var licCumulativeProjectionClaimCount []models.LicCumulativeProjectionClaimCount
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&licCumulativeProjectionClaimCount)

		tableData = append(tableData, map[string]interface{}{
			"data":       licCumulativeProjectionClaimCount,
			"table_name": "Cumulative Projection Claim Count",
		})

		var licCumulativeProjectionAverageClaim []models.LicCumulativeProjectionAverageClaim
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&licCumulativeProjectionAverageClaim)

		tableData = append(tableData, map[string]interface{}{
			"data":       licCumulativeProjectionAverageClaim,
			"table_name": "Cumulative Projection Average Claim",
		})

		var licCumulativeProjectionAverageToTotalClaim []models.LicCumulativeProjectionAveragetoTotalClaim
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&licCumulativeProjectionAverageToTotalClaim)

		tableData = append(tableData, map[string]interface{}{
			"data":       licCumulativeProjectionAverageToTotalClaim,
			"table_name": "cumulative_projection_average_to_total_claim",
		})

		var liiat []models.LicIncrementalInflatedAveragetoTotalClaim
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&liiat)

		tableData = append(tableData, map[string]interface{}{
			"data":       liiat,
			"table_name": "incremental_projected_inflated_average_to_total_claim",
		})

		var licIncrementalProjectionAverageToTotalClaim []models.LicIncrementalProjectionAveragetoTotalClaim
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&licIncrementalProjectionAverageToTotalClaim)

		tableData = append(tableData, map[string]interface{}{
			"data":       licIncrementalProjectionAverageToTotalClaim,
			"table_name": "incremental_projection_average_to_total_claim",
		})

		var ldiat []models.LicDiscountedIncrementalInflatedAveragetoTotalClaim
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&ldiat)

		tableData = append(tableData, map[string]interface{}{
			"data":       ldiat,
			"table_name": "discounted_incremental_inflated_average_to_total_claim",
		})
	case "ibnr_results":
		var ibnrReserves []models.LicIbnrReserve
		if product == "All Products" {
			DB.Where("portfolio_name = ? and run_date = ? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, id).Find(&ibnrReserves)
		} else {
			DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&ibnrReserves)
		}
		//DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.ExpConfigurationName, runSettings.RunDate, product, id).Find(&ibnrReserves)

		tableData = append(tableData, map[string]interface{}{
			"data":       ibnrReserves,
			"table_name": "ibnr_reserves(by accident year and accident month)",
		})
		var licIbnrReserveSummary []models.IbnrReserveReport
		//All products will only apply to ibnr reserves summary for now
		if product == "All Products" {
			DB.Where("portfolio_name = ? and run_date = ? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, id).Find(&licIbnrReserveSummary)
		} else {
			DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&licIbnrReserveSummary)
		}

		tableData = append(tableData, map[string]interface{}{
			"data":       licIbnrReserveSummary,
			"table_name": "ibnr_reserve_summary",
		})
	case "bootstrap_results":
		var licIbnrFrequency []models.IbnrFrequency
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&licIbnrFrequency)

		if len(licIbnrFrequency) > 0 {
			tableData = append(tableData, map[string]interface{}{
				"data":       licIbnrFrequency,
				"table_name": "IBNR Frequency",
			})
		}

		var licBootstrapSummaryStats []models.LicBootstrappedResultSummary
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&licBootstrapSummaryStats)

		if len(licBootstrapSummaryStats) > 0 {
			tableData = append(tableData, map[string]interface{}{
				"data":       licBootstrapSummaryStats,
				"table_name": "Bootstrap Summary Stats",
			})
		}

	case "mack_models":
		var lidf []models.LicIndividualDevelopmentFactors
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&lidf)

		if len(lidf) > 0 {
			tableData = append(tableData, map[string]interface{}{
				"data":       lidf,
				"table_name": "Individual Development Factors",
			})
		}

		var lbaf []models.LicBiasAdjustmentFactor
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&lbaf)

		if len(lbaf) > 0 {
			tableData = append(tableData, map[string]interface{}{
				"data":       lbaf,
				"table_name": "Bias Adjustment Factor",
			})
		}

		var mmcp []models.LicMackModelCalculatedParameters
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&mmcp)

		if len(mmcp) > 0 {
			tableData = append(tableData, map[string]interface{}{
				"data":       mmcp,
				"table_name": "Mack Model Calculated Parameters",
			})
		}

		var lbar []models.LicBiasAdjustedResiduals
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&lbar)

		if len(lbar) > 0 {
			tableData = append(tableData, map[string]interface{}{
				"data":       lbar,
				"table_name": "Bias Adjusted Residuals",
			})
		}

		var lmbar []models.LicMeanBiasAdjustedResiduals
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&lmbar)

		if len(lmbar) > 0 {
			tableData = append(tableData, map[string]interface{}{
				"data":       lmbar,
				"table_name": "Mean Bias Adjusted Residuals",
			})
		}

		var logNormalSigmas []models.LicLogNormalSigmas
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&logNormalSigmas)

		if len(logNormalSigmas) > 0 {
			tableData = append(tableData, map[string]interface{}{
				"data":       logNormalSigmas,
				"table_name": "LogNormal Sigma",
			})
		}

		var logNormalMeans []models.LicLogNormalMeans
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&logNormalMeans)

		if len(logNormalMeans) > 0 {
			tableData = append(tableData, map[string]interface{}{
				"data":       logNormalMeans,
				"table_name": "LogNormal Mean",
			})
		}

		var logNormalStandardDeviations []models.LicLogNormalStandardDeviations
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&logNormalStandardDeviations)

		if len(logNormalStandardDeviations) > 0 {
			tableData = append(tableData, map[string]interface{}{
				"data":       logNormalStandardDeviations,
				"table_name": "LogNormal Standard Deviation",
			})
		}

		var pRatio []models.LicPseudoRatios
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&pRatio)

		if len(pRatio) > 0 {
			tableData = append(tableData, map[string]interface{}{
				"data":       pRatio,
				"table_name": "Pseudo Ratios",
			})
		}

		var lmmsdf []models.LicMackModelSimulatedDevelopmentFactor
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&lmmsdf)

		if len(lmmsdf) > 0 {
			tableData = append(tableData, map[string]interface{}{
				"data":       lmmsdf,
				"table_name": "Mack Model Simulated Development Factor",
			})
		}

		var mackcumproj []models.LicMackCumulativeProjection
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&mackcumproj)

		if len(mackcumproj) > 0 {
			tableData = append(tableData, map[string]interface{}{
				"data":       mackcumproj,
				"table_name": "Mack Model Cumulative Projections",
			})
		}

		var mackstats []models.LicMackSimulationSummaryStats
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&mackstats)

		if len(mackstats) > 0 {
			tableData = append(tableData, map[string]interface{}{
				"data":       mackstats,
				"table_name": "Mack Model Simulation Summary Stats",
			})
		}

		var mackfreq []models.MackIbnrFrequency
		DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&mackfreq)

		if len(mackfreq) > 0 {
			tableData = append(tableData, map[string]interface{}{
				"data":       mackfreq,
				"table_name": "Mack Model Frequency Results",
			})
		}

	case "reserve_report":
		// Method comparison: returns IbnrReserveReport (one row per product) with all
		// deterministic method totals side-by-side, including the two combined methods.
		var licIbnrReserveSummary []models.IbnrReserveReport
		DB.Where("portfolio_name = ? and run_date = ? and run_id = ?", runSettings.PortfolioName, runSettings.RunDate, id).Find(&licIbnrReserveSummary)

		tableData = append(tableData, map[string]interface{}{
			"data":       licIbnrReserveSummary,
			"table_name": "Method Comparison",
		})

	case "credibility_breakdown":
		// Per-accident-year credibility weights for combined methods (CL+BF, CL+Cape Cod).
		// Returns LicIbnrReserve which contains proportion_runoff (CL weight),
		// proportion_not_runoff (BF/Cape Cod weight), and all individual + combined IBNR values
		// so the user can see exactly how much weight each method received per accident year.
		var ibnrReserves []models.LicIbnrReserve
		if product == "All Products" || product == "" {
			DB.Where("portfolio_name = ? and run_date = ? and run_id = ?",
				runSettings.PortfolioName, runSettings.RunDate, id).Find(&ibnrReserves)
		} else {
			DB.Where("portfolio_name = ? and run_date = ? and product_code=? and run_id = ?",
				runSettings.PortfolioName, runSettings.RunDate, product, id).Find(&ibnrReserves)
		}
		tableData = append(tableData, map[string]interface{}{
			"data":       ibnrReserves,
			"table_name": "Credibility Breakdown",
		})
	}
	results["run_settings"] = runSettings
	results["table_data"] = tableData

	return results, nil
}

func GetLicPortfolioProducts(portfolioID int) ([]string, error) {
	var productCodes []string
	err := DB.Model(&models.LICClaimsInput{}).
		Where("lic_portfolio_id = ?", portfolioID).
		Distinct().
		Pluck("product_code", &productCodes).Error
	return productCodes, err
}

// GetLicProductsForIbnrRun resolves the portfolio from an IBNR run and returns
// the distinct product codes present in that portfolio's claims data.
func GetLicProductsForIbnrRun(runID int) ([]string, error) {
	var run models.IBNRRunSetting
	if err := DB.First(&run, runID).Error; err != nil {
		return nil, err
	}
	return GetLicPortfolioProducts(run.PortfolioId)
}

func CheckIbnrRunName(name string) bool {
	var runSettings models.IBNRRunSetting
	DB.Where("run_name = ?", name).Find(&runSettings)
	if runSettings.ID > 0 {
		return true
	}
	return false
}

func CheckLicRunName(name string) bool {
	var runSettings models.LicRunSetting
	DB.Where("run_name = ?", name).Find(&runSettings)
	if runSettings.ID > 0 {
		return true
	}
	return false
}

func CheckLicConfigName(name string) bool {
	var config models.LicVariableSet
	DB.Where("configuration_name = ?", name).Find(&config)
	if config.ID > 0 {
		return true
	}
	return false
}

func SaveLicConfiguration(configuration models.LicVariableSet) error {
	switch configuration.RunType {
	case "ibnr":
		fmt.Println("Saving cash flow configuration")
		// We get the relevant run names from valuation runs using id
		for i, _ := range configuration.LicVariables {
			var runSettings models.IBNRRunSetting
			DB.Where("id = ?", configuration.LicVariables[i].RunId).Find(&runSettings)
			configuration.LicVariables[i].RunName = runSettings.RunName
			configuration.LicVariables[i].ID = 0

		}
		fmt.Println(configuration.LicVariables)

	case "cash_flows":
		fmt.Println("Saving ibnr configuration")
		// We get the relevant run names from valuation runs using id
		for i, _ := range configuration.LicVariables {
			var runSettings models.ProjectionJob
			DB.Where("id = ?", configuration.LicVariables[i].RunId).Find(&runSettings)
			configuration.LicVariables[i].RunName = runSettings.RunName
			configuration.LicVariables[i].ID = 0
		}
		fmt.Println(configuration.LicVariables)
	}

	err := DB.Create(&configuration).Error
	return err
}

func GetLicVariableSets() ([]models.LicVariableSet, error) {
	var variableSets []models.LicVariableSet
	err := DB.Preload("LicVariables").Order("id desc").Find(&variableSets).Error
	return variableSets, err
}

func GetModelPointRows() []models.LicModelPoint {
	var modelPoints []models.LicModelPoint
	DB.First(&modelPoints)
	return modelPoints
}

//func GetLicBuildUp(runDate, portfolioName, product string) (map[string]interface{}, error) {
//	var result map[string]interface{}
//	var list []string
//	var buildUp []models.LicBuildupResult
//	var err error
//	if runDate != "" && portfolioName != "" && product != "" {
//		err = DB.Where("portfolio_name = ? and run_date = ? and product_code = ?", portfolioName, runDate, product).Find(&buildUp).Error
//	}
//	if runDate != "" && portfolioName != "" && product == "" {
//		err = DB.Where("portfolio_name = ? and run_date = ?", portfolioName, runDate).Find(&buildUp).Error
//		DB.Table("lic_buildup_results").Distinct("product_code").Where("portfolio_name = ? and run_date = ?", portfolioName, runDate).Pluck("product_code", &list)
//	}
//	if runDate != "" && portfolioName == "" && product == "" {
//		err = DB.Where("run_date = ?", runDate).Find(&buildUp).Error
//		DB.Table("lic_buildup_results").Distinct("portfolio_name").Where("run_date = ?", runDate).Pluck("portfolio_name", &list)
//	}
//	if runDate == "" && portfolioName == "" && product == "" {
//		DB.Table("lic_buildup_results").Distinct("run_date").Pluck("run_date", &list)
//	}
//	result = map[string]interface{}{
//		"data": buildUp,
//		"list": list,
//	}
//
//	return result, err
//}

func GetLicBuildUp(runId int, portfolioName, product string) (map[string]interface{}, error) {
	var result map[string]interface{}
	var list []string
	var buildUp []models.LicBuildupResult
	var err error
	if runId > 0 && portfolioName != "" && product != "" {
		err = DB.Where("run_id = ? and portfolio_name = ? and  product_code = ?", runId, portfolioName, product).Find(&buildUp).Error
	}
	if runId > 0 && portfolioName != "" && product == "" {
		err = DB.Where("portfolio_name = ? and run_id = ?", portfolioName, runId).Find(&buildUp).Error
		DB.Table("lic_buildup_results").Distinct("product_code").Where("portfolio_name = ? and run_id = ?", portfolioName, runId).Pluck("product_code", &list)
	}
	if runId > 0 && portfolioName == "" && product == "" {
		err = DB.Where("run_id = ?", runId).Find(&buildUp).Error
		DB.Table("lic_buildup_results").Distinct("portfolio_name").Where("run_id = ?", runId).Pluck("portfolio_name", &list)
	}

	//if runDate != "" && portfolioName == "" && product == "" {
	//	err = DB.Where("run_date = ?", runDate).Find(&buildUp).Error
	//	DB.Table("lic_buildup_results").Distinct("portfolio_name").Where("run_date = ?", runDate).Pluck("portfolio_name", &list)
	//}
	//if runDate == "" && portfolioName == "" && product == "" {
	//	DB.Table("lic_buildup_results").Distinct("run_date").Pluck("run_date", &list)
	//}
	var auditTrail models.LicRunSetting

	DB.Where("id = ?", runId).Find(&auditTrail)

	result = map[string]interface{}{
		"data":        buildUp,
		"list":        list,
		"audit_trail": auditTrail,
	}

	return result, err
}

func GetLicRuns() ([]models.LicRunSetting, error) {
	var runs []models.LicRunSetting
	err := DB.Order("id desc").Find(&runs).Error
	return runs, err
}

func DeleteLicRun(id int) error {
	err := DB.Where("id = ?", id).Delete(&models.LicRunSetting{}).Error
	err = DB.Where("run_id = ?", id).Delete(&models.LicBuildupResult{}).Error
	err = DB.Where("lic_run_id = ?", id).Delete(&models.LicJournalTransactions{}).Error
	err = DB.Where("csm_run_id = ? and measurement_type=?", id, "LIC").Delete(&models.BalanceSheetRecord{}).Error
	return err
}

func DeleteLicConfig(id int) error {
	err := DB.Where("lic_variable_set_id = ?", id).Delete(&models.LicVariable{}).Error
	err = DB.Where("id = ?", id).Delete(&models.LicVariableSet{}).Error
	return err
}

func CheckExistingLicRun(run models.LicRun) bool {
	var csmRun models.LicRun
	fmt.Println(run.RunDate[:7])
	DB.Table("lic_run_settings").Where("run_date = ? and lic_configuration_name= ? and opening_balance_date = ? ", run.RunDate, run.LicConfigurationName, run.OpeningBalanceDate).Find(&csmRun)
	if csmRun.ID > 0 {
		return true
	} else {
		return false
	}
}

func GetAvailableYieldCurveYears() ([]int, error) {
	var list []int
	err := DB.Table("yield_curve").Distinct("year").Pluck("year", &list).Error
	return list, err
}

func GetAvailableIbnrYieldCurveYears() ([]int, error) {
	var list []int
	err := DB.Table("ibnr_yield_curves").Distinct("year").Pluck("year", &list).Error
	return list, err
}

func GetAvailablePAAYieldCurveYears() ([]int, error) {
	var list []int
	err := DB.Table("paa_yield_curve").Distinct("year").Pluck("year", &list).Error
	return list, err
}

func GetAvailableYieldCurveMonths(year int) ([]int, error) {
	var list []int
	err := DB.Table("yield_curve").Where("year = ?", year).Distinct("month").Pluck("month", &list).Error
	return list, err
}

func GetAvailablePAAYieldCurveMonths(year int) ([]int, error) {
	var list []int
	err := DB.Table("paa_yield_curve").Where("year = ?", year).Distinct("month").Pluck("month", &list).Error
	return list, err
}
