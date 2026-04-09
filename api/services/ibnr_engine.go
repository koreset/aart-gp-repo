package services

import (
	"api/models"
	"api/utils"
	"encoding/csv"
	"fmt"
	"github.com/dgraph-io/ristretto"
	ztable "github.com/gregscott94/z-table-golang"
	"github.com/jinzhu/copier"
	"github.com/jszwec/csvutil"
	"github.com/montanaflynn/stats"
	log "github.com/sirupsen/logrus"
	"io"
	"math"
	"math/rand"
	"mime/multipart"
	"reflect"
	"strconv"
	"strings"

	//"strconv"
	"time"
)

var IbnrCache *ristretto.Cache

func GetCPI(year int, month int) float64 {

	key := strconv.Itoa(year) + "_" + strconv.Itoa(month)
	cached, found := IbnrCache.Get(key)

	if found {
		return cached.(float64)
	}

	var ibnrcpi models.LicCPI
	err := DB.Table("lic_cpis").Where("year_index = ? and month_index = ?", year, month).First(&ibnrcpi).Error
	if err != nil {
		return 0
		//fmt.Println("forward rate db error: ", errors.WithStack(err))
	}
	success := IbnrCache.Set(key, ibnrcpi.CpiIndex, 1)
	if !success {
		fmt.Println("cache error: Not saved")
	}

	return ibnrcpi.CpiIndex
}

func ListDates(run models.IBNRRunSetting, productCode string, interval int) []models.LicTriangulation {
	startDate, _ := time.Parse("2006-01", run.DataInputStartDate)
	endDate, _ := time.Parse("2006-01", run.DataInputEndDate)
	endDate = endDate.AddDate(0, 1, 0)
	//zero out triangulations slice
	var triangulations []models.LicTriangulation
	var earnedPremiums []models.LICEarnedPremium

	err := DB.Where("portfolio_name = ? and product_code = ? and year_index >= ? and year_index <= ? and year=? and version_name = ?", run.PortfolioName, productCode, startDate.Year(), endDate.Year(), run.ClaimsDataYear, run.ClaimsInputVersion).Find(&earnedPremiums).Error
	if err != nil {
		fmt.Println(err)
		//	return err
	}

	for startDate.Before(endDate) {
		var triangulation models.LicTriangulation
		year, _ := strconv.Atoi(startDate.Format("2006"))
		triangulation.AccidentYear = year
		month, _ := strconv.Atoi(startDate.Format("01"))
		if run.CalculationInterval == Annual {
			triangulation.AccidentMonth = 0
			var premiumSum float64
			for _, earnedPremium := range earnedPremiums {
				if earnedPremium.YearIndex == triangulation.AccidentYear {
					premiumSum += earnedPremium.EarnedPremium
				}
			}
			triangulation.EarnedPremium = premiumSum
		}
		if run.CalculationInterval != Annual {
			triangulation.AccidentMonth = int(math.Ceil(float64(month) / float64(interval)))
			var premiumSum float64
			for _, earnedPremium := range earnedPremiums {
				premiumMonth := int(math.Ceil(float64(earnedPremium.Month) / float64(interval)))
				if earnedPremium.YearIndex == triangulation.AccidentYear && premiumMonth == triangulation.AccidentMonth {
					premiumSum += earnedPremium.EarnedPremium
				}
			}

			triangulation.EarnedPremium = premiumSum

		}
		triangulation.RunDate = run.RunDate
		triangulation.RunID = run.ID
		triangulation.PortfolioName = run.PortfolioName
		triangulation.ProductCode = productCode
		triangulation.LicPortfolioId = run.PortfolioId
		triangulations = append(triangulations, triangulation)

		if run.CalculationInterval == Monthly {
			startDate = startDate.AddDate(0, 1, 0)
		}
		if run.CalculationInterval == Quarterly {
			startDate = startDate.AddDate(0, 3, 0)
		}
		if run.CalculationInterval == Annual {
			startDate = startDate.AddDate(0, 12, 0)
		}

	}
	return triangulations
}

func ListDatesClaimCount(run models.IBNRRunSetting, productCode string, interval int) []models.LicTriangulationClaimCount {
	startDate, _ := time.Parse("2006-01", run.DataInputStartDate)
	endDate, _ := time.Parse("2006-01", run.DataInputEndDate)
	endDate = endDate.AddDate(0, 1, 0)
	//zero out triangulations slice
	var triangulations []models.LicTriangulationClaimCount

	for startDate.Before(endDate) {
		var triangulation models.LicTriangulationClaimCount
		year, _ := strconv.Atoi(startDate.Format("2006"))
		triangulation.AccidentYear = year
		month, _ := strconv.Atoi(startDate.Format("01"))
		if run.CalculationInterval == Annual {
			triangulation.AccidentMonth = 0
		}
		if run.CalculationInterval != Annual {
			triangulation.AccidentMonth = int(math.Ceil(float64(month) / float64(interval)))
		}
		triangulation.RunDate = run.RunDate
		triangulation.RunID = run.ID
		triangulation.PortfolioName = run.PortfolioName
		triangulation.ProductCode = productCode
		triangulation.LicPortfolioId = run.PortfolioId
		triangulations = append(triangulations, triangulation)

		if run.CalculationInterval == Monthly {
			startDate = startDate.AddDate(0, 1, 0)
		}
		if run.CalculationInterval == Quarterly {
			startDate = startDate.AddDate(0, 3, 0)
		}
		if run.CalculationInterval == Annual {
			startDate = startDate.AddDate(0, 12, 0)
		}

	}
	return triangulations
}

func ListDatesAverageClaim(run models.IBNRRunSetting, productCode string, interval int) []models.LicTriangulationAverageClaim {
	startDate, _ := time.Parse("2006-01", run.DataInputStartDate)
	endDate, _ := time.Parse("2006-01", run.DataInputEndDate)
	endDate = endDate.AddDate(0, 1, 0)
	//zero out triangulations slice
	var triangulations []models.LicTriangulationAverageClaim

	for startDate.Before(endDate) {
		var triangulation models.LicTriangulationAverageClaim
		year, _ := strconv.Atoi(startDate.Format("2006"))
		triangulation.AccidentYear = year
		month, _ := strconv.Atoi(startDate.Format("01"))
		if run.CalculationInterval == Annual {
			triangulation.AccidentMonth = 0
		}
		if run.CalculationInterval != Annual {
			triangulation.AccidentMonth = int(math.Ceil(float64(month) / float64(interval)))
		}
		triangulation.RunDate = run.RunDate
		triangulation.RunID = run.ID
		triangulation.PortfolioName = run.PortfolioName
		triangulation.ProductCode = productCode
		triangulation.LicPortfolioId = run.PortfolioId
		triangulations = append(triangulations, triangulation)

		if run.CalculationInterval == Monthly {
			startDate = startDate.AddDate(0, 1, 0)
		}
		if run.CalculationInterval == Quarterly {
			startDate = startDate.AddDate(0, 3, 0)
		}
		if run.CalculationInterval == Annual {
			startDate = startDate.AddDate(0, 12, 0)
		}

	}
	return triangulations
}

func monthsCountSince(run models.IBNRRunSetting) int {
	startDate, _ := time.Parse("2006-01", run.DataInputStartDate)
	endDate, _ := time.Parse("2006-01", run.DataInputEndDate)
	endDate = endDate.AddDate(0, 1, 0)
	var months int
	for startDate.Before(endDate) {
		months++
		if run.CalculationInterval == Monthly {
			startDate = startDate.AddDate(0, 1, 0)
		}
		if run.CalculationInterval == Quarterly {
			startDate = startDate.AddDate(0, 3, 0)
		}
		if run.CalculationInterval == Annual {
			startDate = startDate.AddDate(0, 12, 0)
		}
		//nextMonth := startDate.Month()
		//if nextMonth != month {
		//	months++
		//}
		//month = nextMonth
	}

	return months
}

//var developmentFactors []models.LicDevelopmentFactor

func RunIbnrWithManualFactors(fileHeader *multipart.FileHeader, ibnrRunId int, manualProdCode string, user models.AppUser) error {
	var delimiter rune
	delimiterFile, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer delimiterFile.Close()
	delimiter, err = utils.GetDelimiter(delimiterFile)

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.Comma = delimiter
	dec, err := csvutil.NewDecoder(reader)
	if err != nil {
		fmt.Println(err)
	}
	dec.Header()

	var dfs []models.LicDevelopmentFactor
	for {
		var df models.LicDevelopmentFactor
		if err := dec.Decode(&df); err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}
		if strings.ToLower(df.DevelopmentVariable) == "development factor" {
			dfs = append(dfs, df)
		}
	}

	var run models.IBNRRunSetting
	DB.Where("id = ?", ibnrRunId).First(&run)
	run.RerunIndicator = true
	run.ProcessedRecords = 0
	if !strings.Contains(run.RunName, "[Smoothed]") {
		run.RunName = run.RunName + " [Smoothed]"
	}
	var ibnrRuns []models.IBNRRunSetting
	ibnrRuns = append(ibnrRuns, run)

	err = RunIBNR(ibnrRuns, user, dfs, manualProdCode)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func RunIBNR(runjobs []models.IBNRRunSetting, user models.AppUser, manualDfs []models.LicDevelopmentFactor, manualProdCode string) error {
	IbnrCache.Clear()

	// Save IBNRRunSetting to DB
	for i, _ := range runjobs {
		if runjobs[i].RerunIndicator {
			if strings.Contains(runjobs[i].RunName, "[Smoothed]") {
				// Smoothed run naming: "Base [Smoothed] N"
				parts := strings.Split(runjobs[i].RunName, " [Smoothed] ")
				if len(parts) > 1 {
					base := strings.Join(parts[:len(parts)-1], " [Smoothed] ")
					number, _ := strconv.Atoi(parts[len(parts)-1])
					runjobs[i].RunName = base + " [Smoothed] " + strconv.Itoa(number+1)
				} else {
					runjobs[i].RunName = runjobs[i].RunName + " 1"
				}
			} else {
				parts := strings.Split(runjobs[i].RunName, " Rerun ")
				if len(parts) > 1 {
					number, _ := strconv.Atoi(parts[1])
					runjobs[i].RunName = parts[0] + " Rerun " + strconv.Itoa(number+1)
				} else {
					runjobs[i].RunName = runjobs[i].RunName + " Rerun 1"
				}
			}
		}
		runjobs[i].CreationDate = time.Now()
		runjobs[i].ProcessingStatus = "queued"
		runjobs[i].ProcessedRecords = 0
		var count int64
		DB.Table("lic_claims_inputs").Where("year = ? and portfolio_name = ? and version_name = ?", runjobs[i].ClaimsDataYear, runjobs[i].PortfolioName, runjobs[i].ClaimsInputVersion).Count(&count)
		runjobs[i].TotalRecords = int(count)
		runjobs[i].UserName = user.UserName
		runjobs[i].UserEmail = user.UserEmail
		err := DB.Save(&runjobs[i]).Error
		if err != nil {
			fmt.Println("Error saving IBNRRunSetting to DB: ", err)
			return err
		}
	}

	fmt.Println("Running dataprep")
	startDataPrep := time.Now()
	for _, run := range runjobs {
		var productCodes []string
		var err error
		productCodes = getProductCodes(run)

		startRunTime := time.Now()
		run.ProcessingStatus = "processing"
		//var count int64
		//DB.Table("lic_claims_inputs").Where("year = ? and portfolio_name = ?", run.ClaimsDataYear, run.ExpConfigurationName).Count(&count)
		//run.TotalRecords = int(count)
		DB.Save(&run)

		// check for yield curve data and cpi data for the run
		// if not present, fail the run
		// get all lic parameters for the portfolio
		var licParams []models.LICParameter
		err = DB.Where("portfolio_name = ?", run.PortfolioName).Find(&licParams).Error
		var runFailed bool

		for _, licParam := range licParams {
			var yieldCurveData models.IbnrYieldCurve
			err = DB.Where("yield_curve_code = ?", licParam.YieldCurveCode).First(&yieldCurveData).Error
			if err != nil {
				fmt.Println("Error getting yield curve data for ", licParam.YieldCurveCode)
				// fail the run
				run.ProcessingStatus = "failed"
				run.FailureReason = "Yield curve data not found for " + licParam.YieldCurveCode
				DB.Save(&run)
				runFailed = true
				break
			}
		}

		if runFailed {
			continue
		}

		runDate := run.RunDate
		fmt.Println("Run Date: ", runDate)
		parts := strings.Split(runDate, "-")
		valYear, _ := strconv.Atoi(parts[0])
		valMonth, _ := strconv.Atoi(parts[1])

		// get valid cpi data
		var cpiData models.LicCPI
		err = DB.Where("year_index = ? and month_index = ?", valYear, valMonth).First(&cpiData).Error
		if err != nil {
			fmt.Println("Error getting cpi data for ", valYear, valMonth)
			// fail the run
			run.ProcessingStatus = "failed"
			run.FailureReason = "CPI data not found for valuation run date, " + runDate
			DB.Save(&run)
			continue
		}

		var inflationYear, inflationMonth int

		if run.InflationIndicator {
			inflationYear, err = strconv.Atoi(run.InflationDate[:4])
			if err != nil {
				inflationYear = 0
			}
			inflationMonth, err = strconv.Atoi(run.InflationDate[6:])
			if err != nil {
				inflationMonth = 0
			}
		}

		var interval int
		if run.CalculationInterval == Monthly {
			interval = 1
		}
		if run.CalculationInterval == Quarterly {
			interval = 3
		}
		if run.CalculationInterval == Annual {
			interval = 12
		}

		var portfolio models.LicPortfolio
		err = DB.Where("id=?", run.PortfolioId).First(&portfolio).Error
		if err != nil {
			fmt.Println(err)
			//	return err
		}
		discountOption := portfolio.DiscountOption

		projMonths := monthsCountSince(run)

		// Load all method assignments for this portfolio once before the product loop.
		// ResolveAYMethod (called inside UpdateDevelopmentFactors) handles both
		// product-level and AY-range lookups, falling back to run.IBNRMethod when
		// no assignment exists — preserving the existing single-method behaviour.
		allAssignments := BuildProductAssignmentList(run.PortfolioId)

		for _, productCode := range productCodes {
			// Shallow copy of run — safe because IBNRRunSetting has no pointer fields.
			// Per-product and per-AY method resolution now happens inside UpdateDevelopmentFactors
			// via ResolveAYMethod, so no IBNRMethod override is needed here.
			productRun := run

			startTime := time.Now()
			IbnrCache.Clear()
			triangulations := ListDates(productRun, productCode, interval)
			triangulationsClaimCount := ListDatesClaimCount(run, productCode, interval)
			trangulationsAverageClaim := ListDatesAverageClaim(run, productCode, interval)
			var cumulativeTriangulations []models.LicCumulativeTriangulation
			var cumulativeTriangulationsClaimCount []models.LicCumulativeTriangulationClaimCount
			var cumulativeTriangulationsAverageClaim []models.LicCumulativeTriangulationAverageClaim
			var developmentFactor models.LicDevelopmentFactor
			var individualDevelopmentFactors []models.LicIndividualDevelopmentFactors
			var biasAdjustmentFactors []models.LicBiasAdjustmentFactor
			var cumulativeProjections []models.LicCumulativeProjection
			var ibrmps []models.LicModelPoint
			var ibnrclaims []models.LICClaimsInput
			var params models.LICParameter
			var ibnrReserves []models.LicIbnrReserve
			var IbnrReserveReport models.IbnrReserveReport

			err := DB.Where("year = ? and portfolio_name = ? and product_code = ? and version_name = ?", run.ClaimsDataYear, run.PortfolioName, productCode, run.ClaimsInputVersion).Find(&ibnrclaims).Error
			if err != nil {
				fmt.Println(err)
				//	return err
			}

			endDate, _ := time.Parse("2006-01", run.DataInputEndDate)
			endDateYear := endDate.Year()
			var endDateMonth int
			if run.CalculationInterval == Monthly {
				endDateMonth = int(endDate.Month())
			}
			if run.CalculationInterval == Quarterly {
				endDateMonth = int(math.Ceil(float64(endDate.Month()) / 3.0))
			}
			if run.CalculationInterval == Annual {
				endDateMonth = 0
			}
			var skippedClaimsData int
			for _, mp := range ibnrclaims {

				if (mp.ReportedYear <= endDateYear && mp.ReportedMonth <= int(endDate.Month())) || (mp.ReportedYear < endDateYear && mp.ReportedMonth >= int(endDate.Month())) {
					var ibrmp models.LicModelPoint
					ibrmp.ProductCode = mp.ProductCode
					ibrmp.ProductName = mp.ProductName
					ibrmp.PortfolioName = mp.PortfolioName
					ibrmp.LicPortfolioID = mp.LicPortfolioID
					ibrmp.RunDate = run.RunDate
					ibrmp.RunID = run.ID
					ibrmp.ClaimNumber = mp.ClaimNumber
					ibrmp.IFRS17Group = mp.IFRS17Group
					ibrmp.AccidentYear = mp.DamageYear
					ibrmp.AccidentMonth = mp.DamageMonth
					ibrmp.ReportingYear = mp.ReportedYear
					ibrmp.ReportingMonth = mp.ReportedMonth
					ibrmp.UWYear = mp.UnderwritingYear
					ibrmp.VersionName = mp.VersionName
					//ibrmp.EarnedPremium = mp.PremiumIncome
					ibrmp.CPIDenominator = GetCPI(ibrmp.ReportingYear, ibrmp.ReportingMonth)
					if run.InflationIndicator {
						ibrmp.CPINumerator = GetCPI(inflationYear, inflationMonth)
					}
					ibrmp.UnInflatedClaim = mp.ClaimAmount
					ibrmp.PaidClaims = mp.PaidClaims
					ibrmp.SettlementYear = mp.SettlementYear
					ibrmp.SettlementMonth = mp.SettlementMonth
					//var calculatedAccidentMonth int
					//if interval == 12 {
					//	calculatedAccidentMonth = 0
					//}
					//if interval != 12 {
					//	calculatedAccidentMonth = int(math.Ceil(float64(ibrmp.AccidentMonth) / float64(interval)))
					//}
					params = getParams(run.ParameterYear, mp.ProductCode, run.Basis)
					ibrmp.AssessmentCost = mp.ClaimAmount*params.AllocatedClaimsExpenseProportion + params.AllocatedClaimsExpenseAmount
					ibrmp.TotalUnInflatedClaim = ibrmp.UnInflatedClaim + ibrmp.AssessmentCost

					if ibrmp.CPIDenominator > 0 && run.InflationIndicator {
						ibrmp.TotalInflatedClaim = ibrmp.TotalUnInflatedClaim * (ibrmp.CPINumerator / ibrmp.CPIDenominator)
					} else {
						ibrmp.TotalInflatedClaim = ibrmp.TotalUnInflatedClaim
					}

					if run.CalculationInterval == Monthly {
						ibrmp.ReportingDelay = int(math.Max(math.Max(float64(ibrmp.ReportingYear-ibrmp.AccidentYear), 0)*12+float64(ibrmp.ReportingMonth-ibrmp.AccidentMonth), 0))
					}
					if run.CalculationInterval == Quarterly {
						temp := int(math.Max(math.Max(float64(ibrmp.ReportingYear-ibrmp.AccidentYear), 0)*12+float64(ibrmp.ReportingMonth-ibrmp.AccidentMonth), 0))
						ibrmp.ReportingDelay = int(math.Ceil(float64(temp) / 3.0))
					}
					if run.CalculationInterval == Annual {
						temp := int(math.Max(math.Max(float64(ibrmp.ReportingYear-ibrmp.AccidentYear), 0)*12+float64(ibrmp.ReportingMonth-ibrmp.AccidentMonth), 0))
						ibrmp.ReportingDelay = int(math.Floor(float64(temp) / 12.0))
					}

					ibrmps = append(ibrmps, ibrmp)
					run.ProcessedRecords += 1
					DB.Save(&run)
					mutex.Lock()
					UpdateTriangulations(&triangulations, &triangulationsClaimCount, ibrmp, interval)
					mutex.Unlock()
				} else {
					skippedClaimsData += 1
				}
			}
			//Progress points

			// TODO: decision was taken to comment this out. Will find a solution to this later
			//DB.Where("run_date = ? and product_code = ? and version_name=?", run.RunDate, productCode, run.ClaimsInputVersion).Delete(&models.LicModelPoint{})

			err = DB.CreateInBatches(&ibrmps, 100).Error
			if err != nil {
				fmt.Println("Error saving IBRMP to DB: ", err)
				return err
			}

			if !run.RerunIndicator {
				//DB.Where("run_date=? and product_code=?", run.RunDate, productCode).Delete(&models.LicTriangulation{})
				err = DB.Save(&triangulations).Error
				if err != nil {
					fmt.Println(err)
				}

				//DB.Where("run_date=? and product_code=?", run.RunDate, productCode).Delete(&models.LicTriangulationClaimCount{})
				err = DB.Save(&triangulationsClaimCount).Error
				if err != nil {
					fmt.Println(err)
				}
			}

			UpdateCumulative(&cumulativeTriangulations, triangulations, &cumulativeTriangulationsClaimCount, triangulationsClaimCount, &cumulativeTriangulationsAverageClaim, trangulationsAverageClaim, endDateYear, endDateMonth, projMonths, interval)

			if !run.RerunIndicator {
				//DB.Where("run_date=? and product_code = ?", run.RunDate, productCode).Delete(&models.LicCumulativeTriangulation{})
				DB.Save(&cumulativeTriangulations)

				//DB.Where("run_date=? and product_code = ?", run.RunDate, productCode).Delete(&models.LicCumulativeTriangulationClaimCount{})
				DB.Save(&cumulativeTriangulationsClaimCount)

				//DB.Where("run_date=? and product_code = ?", run.RunDate, productCode).Delete(&models.LicCumulativeTriangulationAverageClaim{})
				DB.Save(&cumulativeTriangulationsAverageClaim)
			}

			// Pass allAssignments so UpdateDevelopmentFactors can resolve the effective method
			// per accident year inside the reserve loop via ResolveAYMethod.
			UpdateDevelopmentFactors(triangulations, cumulativeTriangulations, &developmentFactor, cumulativeTriangulationsClaimCount, cumulativeTriangulationsAverageClaim, &ibnrReserves, &IbnrReserveReport, productRun.RunDate, endDateYear, endDateMonth, projMonths, params, productRun, productCode, interval, discountOption, manualDfs, manualProdCode, allAssignments)

			if run.MackModelIndicator {
				UpdateMackModel(&individualDevelopmentFactors, &biasAdjustmentFactors, cumulativeTriangulations, cumulativeProjections, developmentFactor, &ibnrReserves, &IbnrReserveReport, productRun.RunDate, endDateYear, endDateMonth, projMonths, params, productRun, productCode, interval, discountOption, manualDfs, manualProdCode)
			}

			//DB.Where("run_date=? and product_code=?", runDate, IbnrReserveReport.ProductCode).Delete(&models.IbnrReserveReport{})
			err = DB.Save(&IbnrReserveReport).Error
			if err != nil {
				log.Println(err)
			}

			DB.Save(&ibnrReserves)

			err = DB.Save(&run).Error
			if err != nil {
				fmt.Println(err)
			}
			endTime := time.Since(startTime)
			fmt.Println("Time taken to run", productCode, ":", endTime.Seconds(), " seconds")

		}
		endRunTime := time.Since(startRunTime)
		run.RunTime = endRunTime.Seconds()
		run.ProcessingStatus = "complete"
		err = DB.Save(&run).Error
		if err != nil {
			fmt.Println(err)
		}
	}
	endDataPrep := time.Since(startDataPrep)

	fmt.Println("Time taken to run ibnr: ", endDataPrep.Seconds(), " seconds")
	return nil
}

func getProductCodes(run models.IBNRRunSetting) []string {
	var productCodes []string
	var ibnrclaims []models.LICClaimsInput
	err := DB.Distinct("product_code").Where("year = ? and portfolio_name = ? and version_name = ?", run.ClaimsDataYear, run.PortfolioName, run.ClaimsInputVersion).Find(&ibnrclaims).Error
	if err != nil {
		fmt.Println(err)
		//	return err
	}
	for _, mp := range ibnrclaims {
		productCodes = append(productCodes, mp.ProductCode)
	}
	return productCodes
}

func Triangulation(run models.IBNRRunSetting, mp []models.LicModelPoint) {
	//var ibrdata []models.LicModelPoint
	//ibrdata = RunIBNR(run) // Prepares and saves claims data as ibrmp
	//var triang []models.LicTriangulation

	for i := 0; i < 60; i++ {

	}

}

func UpdateTriangulations(triangulations *[]models.LicTriangulation, triangulationsClaimCount *[]models.LicTriangulationClaimCount, ibrmp models.LicModelPoint, interval int) {
	//var triangulations []models.LicTriangulation

	//if i, ok := utils.TriangulationsContain(triangulations, ibrmp); ok {

	for i, _ := range *triangulations {
		var calculatedAccidentMonth int
		if interval == 12 {
			calculatedAccidentMonth = 0
		}
		if interval != 12 {
			calculatedAccidentMonth = int(math.Ceil(float64(ibrmp.AccidentMonth) / float64(interval)))
		}
		if (*triangulations)[i].AccidentYear == ibrmp.AccidentYear && (*triangulations)[i].AccidentMonth == calculatedAccidentMonth {
			mutable := reflect.ValueOf(&(*triangulations)[i]).Elem()
			mutableClaimCount := reflect.ValueOf(&(*triangulationsClaimCount)[i]).Elem()
			ird := strconv.Itoa(ibrmp.ReportingDelay)
			fieldName := "Rd" + ird
			mutable.FieldByName(fieldName).SetFloat(mutable.FieldByName(fieldName).Float() + ibrmp.TotalInflatedClaim)
			mutableClaimCount.FieldByName(fieldName).SetFloat(mutableClaimCount.FieldByName(fieldName).Float() + 1)
			//mutable.FieldByName(EarnedPremium).SetFloat(mutable.FieldByName(EarnedPremium).Float() + ibrmp.EarnedPremium)
			mutableClaimCount.FieldByName(EarnedPremium).SetFloat(mutableClaimCount.FieldByName(EarnedPremium).Float() + ibrmp.EarnedPremium)

		}
	}
}

func UpdateCumulative(cumulativeTriangulations *[]models.LicCumulativeTriangulation, triangulations []models.LicTriangulation, cumulativeTriangulationsClaimCount *[]models.LicCumulativeTriangulationClaimCount, triangulationsClaimCount []models.LicTriangulationClaimCount, cumulativeTriangulationsAverageClaim *[]models.LicCumulativeTriangulationAverageClaim, triangulationsAverageClaim []models.LicTriangulationAverageClaim, endDateYear, endDateMonth, projMonths, interval int) {
	for i, _ := range triangulations {
		mutable := reflect.ValueOf(&triangulations[i]).Elem()
		mutableClaimCount := reflect.ValueOf(&triangulationsClaimCount[i]).Elem()
		//mutable2 := reflect.ValueOf(&cumulativeTriangulations)
		var cumulativeTriangulation models.LicCumulativeTriangulation
		var cumulativeTriangulationClaimCount models.LicCumulativeTriangulationClaimCount
		var cumulativeTriangulationAverageClaim models.LicCumulativeTriangulationAverageClaim
		m2 := reflect.ValueOf(&cumulativeTriangulation).Elem()
		m2ClaimCount := reflect.ValueOf(&cumulativeTriangulationClaimCount).Elem()
		m2AverageClaim := reflect.ValueOf(&cumulativeTriangulationAverageClaim).Elem()
		var yearadd, movingyear int
		var residualmonth, movingmonth int
		for k := 0; k <= projMonths+TriangleVariableCount; k++ {
			if k <= TriangleVariableCount { //assigning variable names
				m2.Field(k).Set(mutable.Field(k))
				m2ClaimCount.Field(k).Set(mutable.Field(k))
				m2AverageClaim.Field(k).Set(mutable.Field(k))
			} else {
				//fmt.Println(mutable.Type().Field(k).Name)
				if interval == 12 {
					yearadd = int(int64((int(mutable.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1) / (12 / interval)))
					residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mutable.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				}
				if interval != 12 {
					yearadd = int(int64((int(mutable.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1 - 1) / (12 / interval)))
					residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mutable.FieldByName(AccidentMonth).Int())-1), 12/float64(interval)))
				}
				residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mutable.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				movingyear = int(mutable.FieldByName(AccidentYear).Int()) + yearadd
				if residualmonth == 0 && interval == 3 {
					residualmonth = 4
				}
				if residualmonth == 0 && interval == 1 {
					residualmonth = 12
				}

				movingmonth = residualmonth

				if k == TriangleVariableCount+1 {
					m2.Field(k).SetFloat(mutable.Field(k).Float())
					m2ClaimCount.Field(k).SetFloat(mutableClaimCount.Field(k).Float())
					if mutableClaimCount.Field(k).Float() > 0 {
						m2AverageClaim.Field(k).SetFloat(m2.Field(k).Float() / m2ClaimCount.Field(k).Float())
					}
				} else {
					if (movingyear <= endDateYear && movingmonth <= endDateMonth) || (movingyear < endDateYear && movingmonth > endDateMonth) {
						m2.Field(k).SetFloat(m2.Field(k-1).Float() + mutable.Field(k).Float())
						m2ClaimCount.Field(k).SetFloat(m2ClaimCount.Field(k-1).Float() + mutableClaimCount.Field(k).Float())
						if m2ClaimCount.Field(k).Float() > 0 {
							m2AverageClaim.Field(k).SetFloat(m2.Field(k).Float() / m2ClaimCount.Field(k).Float())
						}
					} else {
						m2.Field(k).SetFloat(0)
						m2ClaimCount.Field(k).SetFloat(0)
						m2AverageClaim.Field(k).SetFloat(0)
					}
				}
			}
		}
		//fmt.Println(cumulativeTriangulation)
		*cumulativeTriangulations = append(*cumulativeTriangulations, cumulativeTriangulation)
		*cumulativeTriangulationsClaimCount = append(*cumulativeTriangulationsClaimCount, cumulativeTriangulationClaimCount)
		*cumulativeTriangulationsAverageClaim = append(*cumulativeTriangulationsAverageClaim, cumulativeTriangulationAverageClaim)
	}
}

func UpdateDevelopmentFactors(triangulations []models.LicTriangulation,
	cumulativeTriangulations []models.LicCumulativeTriangulation,
	developmentFactor *models.LicDevelopmentFactor,
	cumulativeTriangulationsClaimCount []models.LicCumulativeTriangulationClaimCount,
	cumulativeTriangulationsAverageClaim []models.LicCumulativeTriangulationAverageClaim,
	ibnrReserves *[]models.LicIbnrReserve,
	IbnrReserveReport *models.IbnrReserveReport,
	runDate string, endDateYear, endDateMonth, projMonths int,
	params models.LICParameter, run models.IBNRRunSetting,
	productCode string, interval int, discountOption string,
	manualDfs []models.LicDevelopmentFactor, manualProdCode string,
	assignments []models.LicIbnrMethodAssignment) {
	//if we are rerunning, we must delete previous data
	if run.RerunIndicator && productCode == manualProdCode {
		//delete previous data
		//DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.ExpConfigurationName, run.ID).Delete(models.LicTriangulation{})
		//DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.ExpConfigurationName, run.ID).Delete(models.LicCumulativeTriangulation{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicDevelopmentFactor{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicCumulativeProjection{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicIncrementalProjection{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicIncrementalInflated{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicDiscountedIncrementalInflated{})

		//DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.ExpConfigurationName, run.ID).Delete(models.LicTriangulationClaimCount{})
		//DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.ExpConfigurationName, run.ID).Delete(models.LicCumulativeTriangulationClaimCount{})
		//DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.ExpConfigurationName, run.ID).Delete(models.LicCumulativeTriangulationAverageClaim{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicDevelopmentFactorClaimCount{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicCumulativeProjectionClaimCount{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicCumulativeProjectionClaimCount{})

		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicDevelopmentFactorAverageClaim{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicCumulativeProjectionAverageClaim{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicCumulativeProjectionAveragetoTotalClaim{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicIncrementalInflatedAveragetoTotalClaim{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicIncrementalProjectionAveragetoTotalClaim{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicDiscountedIncrementalInflatedAveragetoTotalClaim{})

		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicIndividualDevelopmentFactors{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicBiasAdjustmentFactor{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicMackModelCalculatedParameters{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicBiasAdjustedResiduals{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicMeanBiasAdjustedResiduals{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicPseudoRatios{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicMackCumulativeProjection{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicMackModelSimulatedDevelopmentFactor{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicMackSimulationSummaryStats{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicMackSimulationResults{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.MackIbnrFrequency{})

		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicBootStrappedCumulativeProjection{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicBootstrappedDevelopmentFactor{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicBootStrappedIncremental{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicBootStrappedIncrementalInflatedProjection{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicBootStrappedIncrementalInflatedDiscountedProjection{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicBootstrappedResults{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicBootstrappedResultSummary{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.IbnrFrequency{})

		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.LicIbnrReserve{})
		DB.Where("run_date = ? and product_code = ? and portfolio_name = ? and run_id =?", runDate, manualProdCode, run.PortfolioName, run.ID).Delete(models.IbnrReserveReport{})
	}

	var columnSum models.LicDevelopmentFactor
	var columnSumClaimCount models.LicDevelopmentFactorClaimCount
	var columnSumAverageClaim models.LicDevelopmentFactorAverageClaim
	columnSum.PortfolioName = cumulativeTriangulations[0].PortfolioName
	columnSum.ProductCode = cumulativeTriangulations[0].ProductCode
	columnSum.RunDate = cumulativeTriangulations[0].RunDate
	columnSum.RunID = cumulativeTriangulations[0].RunID
	columnSum.DevelopmentVariable = "Column Sum"

	columnSumClaimCount.PortfolioName = cumulativeTriangulations[0].PortfolioName
	columnSumClaimCount.ProductCode = cumulativeTriangulations[0].ProductCode
	columnSumClaimCount.RunDate = cumulativeTriangulations[0].RunDate
	columnSumClaimCount.RunID = cumulativeTriangulations[0].RunID
	columnSumClaimCount.DevelopmentVariable = "Column Sum"

	columnSumAverageClaim.PortfolioName = cumulativeTriangulations[0].PortfolioName
	columnSumAverageClaim.ProductCode = cumulativeTriangulations[0].ProductCode
	columnSumAverageClaim.RunDate = cumulativeTriangulations[0].RunDate
	columnSumAverageClaim.RunID = cumulativeTriangulations[0].RunID
	columnSumAverageClaim.DevelopmentVariable = "Column Sum"

	for _, ct := range cumulativeTriangulations {
		columnSum.Rd0 += ct.Rd0
		columnSum.Rd1 += ct.Rd1
		columnSum.Rd2 += ct.Rd2
		columnSum.Rd3 += ct.Rd3
		columnSum.Rd4 += ct.Rd4
		columnSum.Rd5 += ct.Rd5
		columnSum.Rd6 += ct.Rd6
		columnSum.Rd7 += ct.Rd7
		columnSum.Rd8 += ct.Rd8
		columnSum.Rd9 += ct.Rd9
		columnSum.Rd10 += ct.Rd10
		columnSum.Rd11 += ct.Rd11
		columnSum.Rd12 += ct.Rd12
		columnSum.Rd13 += ct.Rd13
		columnSum.Rd14 += ct.Rd14
		columnSum.Rd15 += ct.Rd15
		columnSum.Rd16 += ct.Rd16
		columnSum.Rd17 += ct.Rd17
		columnSum.Rd18 += ct.Rd18
		columnSum.Rd19 += ct.Rd19
		columnSum.Rd20 += ct.Rd20
		columnSum.Rd21 += ct.Rd21
		columnSum.Rd22 += ct.Rd22
		columnSum.Rd23 += ct.Rd23
		columnSum.Rd24 += ct.Rd24
		columnSum.Rd25 += ct.Rd25
		columnSum.Rd26 += ct.Rd26
		columnSum.Rd27 += ct.Rd27
		columnSum.Rd28 += ct.Rd28
		columnSum.Rd29 += ct.Rd29
		columnSum.Rd30 += ct.Rd30
		columnSum.Rd31 += ct.Rd31
		columnSum.Rd32 += ct.Rd32
		columnSum.Rd33 += ct.Rd33
		columnSum.Rd34 += ct.Rd34
		columnSum.Rd35 += ct.Rd35
		columnSum.Rd36 += ct.Rd36
		columnSum.Rd37 += ct.Rd37
		columnSum.Rd38 += ct.Rd38
		columnSum.Rd39 += ct.Rd39
		columnSum.Rd40 += ct.Rd40
		columnSum.Rd41 += ct.Rd41
		columnSum.Rd42 += ct.Rd42
		columnSum.Rd43 += ct.Rd43
		columnSum.Rd44 += ct.Rd44
		columnSum.Rd45 += ct.Rd45
		columnSum.Rd46 += ct.Rd46
		columnSum.Rd47 += ct.Rd47
		columnSum.Rd48 += ct.Rd48
		columnSum.Rd49 += ct.Rd49
		columnSum.Rd50 += ct.Rd50
		columnSum.Rd51 += ct.Rd51
		columnSum.Rd52 += ct.Rd52
		columnSum.Rd53 += ct.Rd53
		columnSum.Rd54 += ct.Rd54
		columnSum.Rd55 += ct.Rd55
		columnSum.Rd56 += ct.Rd56
		columnSum.Rd57 += ct.Rd57
		columnSum.Rd58 += ct.Rd58
		columnSum.Rd59 += ct.Rd59
		columnSum.Rd60 += ct.Rd60
		columnSum.Rd61 += ct.Rd61
		columnSum.Rd62 += ct.Rd62
		columnSum.Rd63 += ct.Rd63
		columnSum.Rd64 += ct.Rd64
		columnSum.Rd65 += ct.Rd65
		columnSum.Rd66 += ct.Rd66
		columnSum.Rd67 += ct.Rd67
		columnSum.Rd68 += ct.Rd68
		columnSum.Rd69 += ct.Rd69
		columnSum.Rd70 += ct.Rd70
		columnSum.Rd71 += ct.Rd71
		columnSum.Rd72 += ct.Rd72
		columnSum.Rd73 += ct.Rd73
		columnSum.Rd74 += ct.Rd74
		columnSum.Rd75 += ct.Rd75
		columnSum.Rd76 += ct.Rd76
		columnSum.Rd77 += ct.Rd77
		columnSum.Rd78 += ct.Rd78
		columnSum.Rd79 += ct.Rd79
		columnSum.Rd80 += ct.Rd80
		columnSum.Rd81 += ct.Rd81
		columnSum.Rd82 += ct.Rd82
		columnSum.Rd83 += ct.Rd83
		columnSum.Rd84 += ct.Rd84
		columnSum.Rd85 += ct.Rd85
		columnSum.Rd86 += ct.Rd86
		columnSum.Rd87 += ct.Rd87
		columnSum.Rd88 += ct.Rd88
		columnSum.Rd89 += ct.Rd89
		columnSum.Rd90 += ct.Rd90
		columnSum.Rd91 += ct.Rd91
		columnSum.Rd92 += ct.Rd92
		columnSum.Rd93 += ct.Rd93
		columnSum.Rd94 += ct.Rd94
		columnSum.Rd95 += ct.Rd95
		columnSum.Rd96 += ct.Rd96
		columnSum.Rd97 += ct.Rd97
		columnSum.Rd98 += ct.Rd98
		columnSum.Rd99 += ct.Rd99
		columnSum.Rd100 += ct.Rd100
		columnSum.Rd101 += ct.Rd101
		columnSum.Rd102 += ct.Rd102
		columnSum.Rd103 += ct.Rd103
		columnSum.Rd104 += ct.Rd104
		columnSum.Rd105 += ct.Rd105
		columnSum.Rd106 += ct.Rd106
		columnSum.Rd107 += ct.Rd107
		columnSum.Rd108 += ct.Rd108
		columnSum.Rd109 += ct.Rd109
		columnSum.Rd110 += ct.Rd110
		columnSum.Rd111 += ct.Rd111
		columnSum.Rd112 += ct.Rd112
		columnSum.Rd113 += ct.Rd113
		columnSum.Rd114 += ct.Rd114
		columnSum.Rd115 += ct.Rd115
		columnSum.Rd116 += ct.Rd116
		columnSum.Rd117 += ct.Rd117
		columnSum.Rd118 += ct.Rd118
		columnSum.Rd119 += ct.Rd119
		columnSum.Rd120 += ct.Rd120
	}

	for _, ct := range cumulativeTriangulationsClaimCount {
		columnSumClaimCount.Rd0 += ct.Rd0
		columnSumClaimCount.Rd1 += ct.Rd1
		columnSumClaimCount.Rd2 += ct.Rd2
		columnSumClaimCount.Rd3 += ct.Rd3
		columnSumClaimCount.Rd4 += ct.Rd4
		columnSumClaimCount.Rd5 += ct.Rd5
		columnSumClaimCount.Rd6 += ct.Rd6
		columnSumClaimCount.Rd7 += ct.Rd7
		columnSumClaimCount.Rd8 += ct.Rd8
		columnSumClaimCount.Rd9 += ct.Rd9
		columnSumClaimCount.Rd10 += ct.Rd10
		columnSumClaimCount.Rd11 += ct.Rd11
		columnSumClaimCount.Rd12 += ct.Rd12
		columnSumClaimCount.Rd13 += ct.Rd13
		columnSumClaimCount.Rd14 += ct.Rd14
		columnSumClaimCount.Rd15 += ct.Rd15
		columnSumClaimCount.Rd16 += ct.Rd16
		columnSumClaimCount.Rd17 += ct.Rd17
		columnSumClaimCount.Rd18 += ct.Rd18
		columnSumClaimCount.Rd19 += ct.Rd19
		columnSumClaimCount.Rd20 += ct.Rd20
		columnSumClaimCount.Rd21 += ct.Rd21
		columnSumClaimCount.Rd22 += ct.Rd22
		columnSumClaimCount.Rd23 += ct.Rd23
		columnSumClaimCount.Rd24 += ct.Rd24
		columnSumClaimCount.Rd25 += ct.Rd25
		columnSumClaimCount.Rd26 += ct.Rd26
		columnSumClaimCount.Rd27 += ct.Rd27
		columnSumClaimCount.Rd28 += ct.Rd28
		columnSumClaimCount.Rd29 += ct.Rd29
		columnSumClaimCount.Rd30 += ct.Rd30
		columnSumClaimCount.Rd31 += ct.Rd31
		columnSumClaimCount.Rd32 += ct.Rd32
		columnSumClaimCount.Rd33 += ct.Rd33
		columnSumClaimCount.Rd34 += ct.Rd34
		columnSumClaimCount.Rd35 += ct.Rd35
		columnSumClaimCount.Rd36 += ct.Rd36
		columnSumClaimCount.Rd37 += ct.Rd37
		columnSumClaimCount.Rd38 += ct.Rd38
		columnSumClaimCount.Rd39 += ct.Rd39
		columnSumClaimCount.Rd40 += ct.Rd40
		columnSumClaimCount.Rd41 += ct.Rd41
		columnSumClaimCount.Rd42 += ct.Rd42
		columnSumClaimCount.Rd43 += ct.Rd43
		columnSumClaimCount.Rd44 += ct.Rd44
		columnSumClaimCount.Rd45 += ct.Rd45
		columnSumClaimCount.Rd46 += ct.Rd46
		columnSumClaimCount.Rd47 += ct.Rd47
		columnSumClaimCount.Rd48 += ct.Rd48
		columnSumClaimCount.Rd49 += ct.Rd49
		columnSumClaimCount.Rd50 += ct.Rd50
		columnSumClaimCount.Rd51 += ct.Rd51
		columnSumClaimCount.Rd52 += ct.Rd52
		columnSumClaimCount.Rd53 += ct.Rd53
		columnSumClaimCount.Rd54 += ct.Rd54
		columnSumClaimCount.Rd55 += ct.Rd55
		columnSumClaimCount.Rd56 += ct.Rd56
		columnSumClaimCount.Rd57 += ct.Rd57
		columnSumClaimCount.Rd58 += ct.Rd58
		columnSumClaimCount.Rd59 += ct.Rd59
		columnSumClaimCount.Rd60 += ct.Rd60
		columnSumClaimCount.Rd61 += ct.Rd61
		columnSumClaimCount.Rd62 += ct.Rd62
		columnSumClaimCount.Rd63 += ct.Rd63
		columnSumClaimCount.Rd64 += ct.Rd64
		columnSumClaimCount.Rd65 += ct.Rd65
		columnSumClaimCount.Rd66 += ct.Rd66
		columnSumClaimCount.Rd67 += ct.Rd67
		columnSumClaimCount.Rd68 += ct.Rd68
		columnSumClaimCount.Rd69 += ct.Rd69
		columnSumClaimCount.Rd70 += ct.Rd70
		columnSumClaimCount.Rd71 += ct.Rd71
		columnSumClaimCount.Rd72 += ct.Rd72
		columnSumClaimCount.Rd73 += ct.Rd73
		columnSumClaimCount.Rd74 += ct.Rd74
		columnSumClaimCount.Rd75 += ct.Rd75
		columnSumClaimCount.Rd76 += ct.Rd76
		columnSumClaimCount.Rd77 += ct.Rd77
		columnSumClaimCount.Rd78 += ct.Rd78
		columnSumClaimCount.Rd79 += ct.Rd79
		columnSumClaimCount.Rd80 += ct.Rd80
		columnSumClaimCount.Rd81 += ct.Rd81
		columnSumClaimCount.Rd82 += ct.Rd82
		columnSumClaimCount.Rd83 += ct.Rd83
		columnSumClaimCount.Rd84 += ct.Rd84
		columnSumClaimCount.Rd85 += ct.Rd85
		columnSumClaimCount.Rd86 += ct.Rd86
		columnSumClaimCount.Rd87 += ct.Rd87
		columnSumClaimCount.Rd88 += ct.Rd88
		columnSumClaimCount.Rd89 += ct.Rd89
		columnSumClaimCount.Rd90 += ct.Rd90
		columnSumClaimCount.Rd91 += ct.Rd91
		columnSumClaimCount.Rd92 += ct.Rd92
		columnSumClaimCount.Rd93 += ct.Rd93
		columnSumClaimCount.Rd94 += ct.Rd94
		columnSumClaimCount.Rd95 += ct.Rd95
		columnSumClaimCount.Rd96 += ct.Rd96
		columnSumClaimCount.Rd97 += ct.Rd97
		columnSumClaimCount.Rd98 += ct.Rd98
		columnSumClaimCount.Rd99 += ct.Rd99
		columnSumClaimCount.Rd100 += ct.Rd100
		columnSumClaimCount.Rd101 += ct.Rd101
		columnSumClaimCount.Rd102 += ct.Rd102
		columnSumClaimCount.Rd103 += ct.Rd103
		columnSumClaimCount.Rd104 += ct.Rd104
		columnSumClaimCount.Rd105 += ct.Rd105
		columnSumClaimCount.Rd106 += ct.Rd106
		columnSumClaimCount.Rd107 += ct.Rd107
		columnSumClaimCount.Rd108 += ct.Rd108
		columnSumClaimCount.Rd109 += ct.Rd109
		columnSumClaimCount.Rd110 += ct.Rd110
		columnSumClaimCount.Rd111 += ct.Rd111
		columnSumClaimCount.Rd112 += ct.Rd112
		columnSumClaimCount.Rd113 += ct.Rd113
		columnSumClaimCount.Rd114 += ct.Rd114
		columnSumClaimCount.Rd115 += ct.Rd115
		columnSumClaimCount.Rd116 += ct.Rd116
		columnSumClaimCount.Rd117 += ct.Rd117
		columnSumClaimCount.Rd118 += ct.Rd118
		columnSumClaimCount.Rd119 += ct.Rd119
		columnSumClaimCount.Rd120 += ct.Rd120

	}

	for _, ct := range cumulativeTriangulationsAverageClaim {
		columnSumAverageClaim.Rd0 += ct.Rd0
		columnSumAverageClaim.Rd1 += ct.Rd1
		columnSumAverageClaim.Rd2 += ct.Rd2
		columnSumAverageClaim.Rd3 += ct.Rd3
		columnSumAverageClaim.Rd4 += ct.Rd4
		columnSumAverageClaim.Rd5 += ct.Rd5
		columnSumAverageClaim.Rd6 += ct.Rd6
		columnSumAverageClaim.Rd7 += ct.Rd7
		columnSumAverageClaim.Rd8 += ct.Rd8
		columnSumAverageClaim.Rd9 += ct.Rd9
		columnSumAverageClaim.Rd10 += ct.Rd10
		columnSumAverageClaim.Rd11 += ct.Rd11
		columnSumAverageClaim.Rd12 += ct.Rd12
		columnSumAverageClaim.Rd13 += ct.Rd13
		columnSumAverageClaim.Rd14 += ct.Rd14
		columnSumAverageClaim.Rd15 += ct.Rd15
		columnSumAverageClaim.Rd16 += ct.Rd16
		columnSumAverageClaim.Rd17 += ct.Rd17
		columnSumAverageClaim.Rd18 += ct.Rd18
		columnSumAverageClaim.Rd19 += ct.Rd19
		columnSumAverageClaim.Rd20 += ct.Rd20
		columnSumAverageClaim.Rd21 += ct.Rd21
		columnSumAverageClaim.Rd22 += ct.Rd22
		columnSumAverageClaim.Rd23 += ct.Rd23
		columnSumAverageClaim.Rd24 += ct.Rd24
		columnSumAverageClaim.Rd25 += ct.Rd25
		columnSumAverageClaim.Rd26 += ct.Rd26
		columnSumAverageClaim.Rd27 += ct.Rd27
		columnSumAverageClaim.Rd28 += ct.Rd28
		columnSumAverageClaim.Rd29 += ct.Rd29
		columnSumAverageClaim.Rd30 += ct.Rd30
		columnSumAverageClaim.Rd31 += ct.Rd31
		columnSumAverageClaim.Rd32 += ct.Rd32
		columnSumAverageClaim.Rd33 += ct.Rd33
		columnSumAverageClaim.Rd34 += ct.Rd34
		columnSumAverageClaim.Rd35 += ct.Rd35
		columnSumAverageClaim.Rd36 += ct.Rd36
		columnSumAverageClaim.Rd37 += ct.Rd37
		columnSumAverageClaim.Rd38 += ct.Rd38
		columnSumAverageClaim.Rd39 += ct.Rd39
		columnSumAverageClaim.Rd40 += ct.Rd40
		columnSumAverageClaim.Rd41 += ct.Rd41
		columnSumAverageClaim.Rd42 += ct.Rd42
		columnSumAverageClaim.Rd43 += ct.Rd43
		columnSumAverageClaim.Rd44 += ct.Rd44
		columnSumAverageClaim.Rd45 += ct.Rd45
		columnSumAverageClaim.Rd46 += ct.Rd46
		columnSumAverageClaim.Rd47 += ct.Rd47
		columnSumAverageClaim.Rd48 += ct.Rd48
		columnSumAverageClaim.Rd49 += ct.Rd49
		columnSumAverageClaim.Rd50 += ct.Rd50
		columnSumAverageClaim.Rd51 += ct.Rd51
		columnSumAverageClaim.Rd52 += ct.Rd52
		columnSumAverageClaim.Rd53 += ct.Rd53
		columnSumAverageClaim.Rd54 += ct.Rd54
		columnSumAverageClaim.Rd55 += ct.Rd55
		columnSumAverageClaim.Rd56 += ct.Rd56
		columnSumAverageClaim.Rd57 += ct.Rd57
		columnSumAverageClaim.Rd58 += ct.Rd58
		columnSumAverageClaim.Rd59 += ct.Rd59
		columnSumAverageClaim.Rd60 += ct.Rd60
		columnSumAverageClaim.Rd61 += ct.Rd61
		columnSumAverageClaim.Rd62 += ct.Rd62
		columnSumAverageClaim.Rd63 += ct.Rd63
		columnSumAverageClaim.Rd64 += ct.Rd64
		columnSumAverageClaim.Rd65 += ct.Rd65
		columnSumAverageClaim.Rd66 += ct.Rd66
		columnSumAverageClaim.Rd67 += ct.Rd67
		columnSumAverageClaim.Rd68 += ct.Rd68
		columnSumAverageClaim.Rd69 += ct.Rd69
		columnSumAverageClaim.Rd70 += ct.Rd70
		columnSumAverageClaim.Rd71 += ct.Rd71
		columnSumAverageClaim.Rd72 += ct.Rd72
		columnSumAverageClaim.Rd73 += ct.Rd73
		columnSumAverageClaim.Rd74 += ct.Rd74
		columnSumAverageClaim.Rd75 += ct.Rd75
		columnSumAverageClaim.Rd76 += ct.Rd76
		columnSumAverageClaim.Rd77 += ct.Rd77
		columnSumAverageClaim.Rd78 += ct.Rd78
		columnSumAverageClaim.Rd79 += ct.Rd79
		columnSumAverageClaim.Rd80 += ct.Rd80
		columnSumAverageClaim.Rd81 += ct.Rd81
		columnSumAverageClaim.Rd82 += ct.Rd82
		columnSumAverageClaim.Rd83 += ct.Rd83
		columnSumAverageClaim.Rd84 += ct.Rd84
		columnSumAverageClaim.Rd85 += ct.Rd85
		columnSumAverageClaim.Rd86 += ct.Rd86
		columnSumAverageClaim.Rd87 += ct.Rd87
		columnSumAverageClaim.Rd88 += ct.Rd88
		columnSumAverageClaim.Rd89 += ct.Rd89
		columnSumAverageClaim.Rd90 += ct.Rd90
		columnSumAverageClaim.Rd91 += ct.Rd91
		columnSumAverageClaim.Rd92 += ct.Rd92
		columnSumAverageClaim.Rd93 += ct.Rd93
		columnSumAverageClaim.Rd94 += ct.Rd94
		columnSumAverageClaim.Rd95 += ct.Rd95
		columnSumAverageClaim.Rd96 += ct.Rd96
		columnSumAverageClaim.Rd97 += ct.Rd97
		columnSumAverageClaim.Rd98 += ct.Rd98
		columnSumAverageClaim.Rd99 += ct.Rd99
		columnSumAverageClaim.Rd100 += ct.Rd100
		columnSumAverageClaim.Rd101 += ct.Rd101
		columnSumAverageClaim.Rd102 += ct.Rd102
		columnSumAverageClaim.Rd103 += ct.Rd103
		columnSumAverageClaim.Rd104 += ct.Rd104
		columnSumAverageClaim.Rd105 += ct.Rd105
		columnSumAverageClaim.Rd106 += ct.Rd106
		columnSumAverageClaim.Rd107 += ct.Rd107
		columnSumAverageClaim.Rd108 += ct.Rd108
		columnSumAverageClaim.Rd109 += ct.Rd109
		columnSumAverageClaim.Rd110 += ct.Rd110
		columnSumAverageClaim.Rd111 += ct.Rd111
		columnSumAverageClaim.Rd112 += ct.Rd112
		columnSumAverageClaim.Rd113 += ct.Rd113
		columnSumAverageClaim.Rd114 += ct.Rd114
		columnSumAverageClaim.Rd115 += ct.Rd115
		columnSumAverageClaim.Rd116 += ct.Rd116
		columnSumAverageClaim.Rd117 += ct.Rd117
		columnSumAverageClaim.Rd118 += ct.Rd118
		columnSumAverageClaim.Rd119 += ct.Rd119
		columnSumAverageClaim.Rd120 += ct.Rd120

	}

	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, columnSum.DevelopmentVariable, columnSum.ProductCode).Delete(&models.LicDevelopmentFactor{})
	DB.Save(&columnSum)

	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, columnSumClaimCount.DevelopmentVariable, columnSumClaimCount.ProductCode).Delete(&models.LicDevelopmentFactorClaimCount{})
	DB.Save(&columnSumClaimCount)

	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, columnSumAverageClaim.DevelopmentVariable, columnSumAverageClaim.ProductCode).Delete(&models.LicDevelopmentFactorAverageClaim{})
	DB.Save(&columnSumAverageClaim)

	var lastvalue models.LicDevelopmentFactor
	var lastvalueClaimCount models.LicDevelopmentFactorClaimCount
	var lastvalueAverageClaim models.LicDevelopmentFactorAverageClaim

	lastvalue.PortfolioName = cumulativeTriangulations[0].PortfolioName
	lastvalue.ProductCode = cumulativeTriangulations[0].ProductCode
	lastvalue.RunDate = cumulativeTriangulations[0].RunDate
	lastvalue.RunID = cumulativeTriangulations[0].RunID
	lastvalue.DevelopmentVariable = "Last Value"

	lastvalueClaimCount.PortfolioName = cumulativeTriangulationsClaimCount[0].PortfolioName
	lastvalueClaimCount.ProductCode = cumulativeTriangulationsClaimCount[0].ProductCode
	lastvalueClaimCount.RunDate = cumulativeTriangulationsClaimCount[0].RunDate
	lastvalueClaimCount.RunID = cumulativeTriangulationsClaimCount[0].RunID
	lastvalueClaimCount.DevelopmentVariable = "Last Value"

	lastvalueAverageClaim.PortfolioName = cumulativeTriangulationsClaimCount[0].PortfolioName
	lastvalueAverageClaim.ProductCode = cumulativeTriangulationsClaimCount[0].ProductCode
	lastvalueAverageClaim.RunDate = cumulativeTriangulationsClaimCount[0].RunDate
	lastvalueAverageClaim.RunID = cumulativeTriangulationsClaimCount[0].RunID
	lastvalueAverageClaim.DevelopmentVariable = "Last Value"

	m2 := reflect.ValueOf(&lastvalue).Elem()
	m2ClaimCount := reflect.ValueOf(&lastvalueClaimCount).Elem()
	m2AverageClaim := reflect.ValueOf(&lastvalueAverageClaim).Elem()
	for i, _ := range cumulativeTriangulations {
		mutable := reflect.ValueOf(&cumulativeTriangulations[i]).Elem()
		mutableClaimCount := reflect.ValueOf(&cumulativeTriangulationsClaimCount[i]).Elem()
		mutableAverageClaim := reflect.ValueOf(&cumulativeTriangulationsAverageClaim[i]).Elem()

		var yearadd, movingyear int
		var residualmonth, movingmonth int

		for k := TriangleVariableCount + 1; k <= projMonths+TriangleVariableCount; k++ {

			if interval == 12 {
				yearadd = int(int64((int(mutable.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1) / (12 / interval)))
				residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mutable.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
			}
			if interval != 12 {
				yearadd = int(int64((int(mutable.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1 - 1) / (12 / interval)))
				residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mutable.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
			}
			movingyear = int(mutable.FieldByName(AccidentYear).Int()) + yearadd
			if residualmonth == 0 && interval == 3 {
				residualmonth = 4
			}

			if residualmonth == 0 && interval == 1 {
				residualmonth = 12
			}

			movingmonth = residualmonth

			if movingyear == endDateYear && movingmonth == endDateMonth {
				m2.Field(k - VariableCountDiff).Set(mutable.Field(k))
				m2ClaimCount.Field(k - VariableCountDiff).Set(mutableClaimCount.Field(k))
				m2AverageClaim.Field(k - VariableCountDiff).Set(mutableAverageClaim.Field(k))
			}
		}
	}

	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, lastvalue.DevelopmentVariable, lastvalue.ProductCode).Delete(&models.LicDevelopmentFactor{})
	DB.Save(&lastvalue)

	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, lastvalueClaimCount.DevelopmentVariable, lastvalueClaimCount.ProductCode).Delete(&models.LicDevelopmentFactorClaimCount{})
	DB.Save(&lastvalueClaimCount)

	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, lastvalueAverageClaim.DevelopmentVariable, lastvalueAverageClaim.ProductCode).Delete(&models.LicDevelopmentFactorAverageClaim{})
	DB.Save(&lastvalueAverageClaim)

	// column excluding last value
	var columnsumexcllastvalue models.LicDevelopmentFactor
	var columnsumexcllastvalueClaimCount models.LicDevelopmentFactorClaimCount
	var columnsumexcllastvalueAverageClaim models.LicDevelopmentFactorAverageClaim
	columnsumexcllastvalue.PortfolioName = cumulativeTriangulations[0].PortfolioName
	columnsumexcllastvalue.ProductCode = cumulativeTriangulations[0].ProductCode
	columnsumexcllastvalue.RunDate = cumulativeTriangulations[0].RunDate
	columnsumexcllastvalue.RunID = cumulativeTriangulations[0].RunID
	columnsumexcllastvalue.DevelopmentVariable = "Column Sum(excl.last value)"

	columnsumexcllastvalueClaimCount.PortfolioName = cumulativeTriangulations[0].PortfolioName
	columnsumexcllastvalueClaimCount.ProductCode = cumulativeTriangulations[0].ProductCode
	columnsumexcllastvalueClaimCount.RunDate = cumulativeTriangulations[0].RunDate
	columnsumexcllastvalueClaimCount.RunID = cumulativeTriangulations[0].RunID
	columnsumexcllastvalueClaimCount.DevelopmentVariable = "Column Sum(excl.last value)"

	columnsumexcllastvalueAverageClaim.PortfolioName = cumulativeTriangulations[0].PortfolioName
	columnsumexcllastvalueAverageClaim.ProductCode = cumulativeTriangulations[0].ProductCode
	columnsumexcllastvalueAverageClaim.RunDate = cumulativeTriangulations[0].RunDate
	columnsumexcllastvalueAverageClaim.RunID = cumulativeTriangulations[0].RunID
	columnsumexcllastvalueAverageClaim.DevelopmentVariable = "Column Sum(excl.last value)"

	mutable := reflect.ValueOf(&columnSum).Elem()
	mutable2 := reflect.ValueOf(&lastvalue).Elem()
	mutable3 := reflect.ValueOf(&columnsumexcllastvalue).Elem()

	mutableClaimCount := reflect.ValueOf(&columnSumClaimCount).Elem()
	mutable2ClaimCount := reflect.ValueOf(&lastvalueClaimCount).Elem()
	mutable3ClaimCount := reflect.ValueOf(&columnsumexcllastvalueClaimCount).Elem()

	mutableAverageClaim := reflect.ValueOf(&columnSumAverageClaim).Elem()
	mutable2AverageClaim := reflect.ValueOf(&lastvalueAverageClaim).Elem()
	mutable3AverageClaim := reflect.ValueOf(&columnsumexcllastvalueAverageClaim).Elem()

	for k := 0; k <= projMonths+FactorVariableCount; k++ {
		if k > FactorVariableCount {
			if mutable.Field(k).Float() > 0 {
				mutable3.Field(k).SetFloat(mutable.Field(k).Float() - mutable2.Field(k).Float())
				mutable3ClaimCount.Field(k).SetFloat(mutableClaimCount.Field(k).Float() - mutable2ClaimCount.Field(k).Float())
				mutable3AverageClaim.Field(k).SetFloat(mutableAverageClaim.Field(k).Float() - mutable2AverageClaim.Field(k).Float())
				continue
			}
		}
	}
	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, columnsumexcllastvalue.DevelopmentVariable, columnsumexcllastvalue.ProductCode).Delete(&models.LicDevelopmentFactor{})
	DB.Save(&columnsumexcllastvalue)

	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, columnsumexcllastvalueClaimCount.DevelopmentVariable, columnsumexcllastvalueClaimCount.ProductCode).Delete(&models.LicDevelopmentFactorClaimCount{})
	DB.Save(&columnsumexcllastvalueClaimCount)

	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, columnsumexcllastvalueAverageClaim.DevelopmentVariable, columnsumexcllastvalueAverageClaim.ProductCode).Delete(&models.LicDevelopmentFactorAverageClaim{})
	DB.Save(&columnsumexcllastvalueAverageClaim)

	//Weighted average succession ratios
	var successionratios models.LicDevelopmentFactor
	var successionratiosClaimCount models.LicDevelopmentFactorClaimCount
	var successionratiosAverageClaim models.LicDevelopmentFactorAverageClaim
	successionratios.PortfolioName = cumulativeTriangulations[0].PortfolioName
	successionratios.ProductCode = cumulativeTriangulations[0].ProductCode
	successionratios.RunDate = cumulativeTriangulations[0].RunDate
	successionratios.RunID = cumulativeTriangulations[0].RunID
	successionratios.DevelopmentVariable = "Weighted Ave. Succession Ratio(rd to rd+1)"

	successionratiosClaimCount.PortfolioName = cumulativeTriangulations[0].PortfolioName
	successionratiosClaimCount.ProductCode = cumulativeTriangulations[0].ProductCode
	successionratiosClaimCount.RunDate = cumulativeTriangulations[0].RunDate
	successionratiosClaimCount.RunID = cumulativeTriangulations[0].RunID
	successionratiosClaimCount.DevelopmentVariable = "Weighted Ave. Succession Ratio(rd to rd+1)"

	successionratiosAverageClaim.PortfolioName = cumulativeTriangulations[0].PortfolioName
	successionratiosAverageClaim.ProductCode = cumulativeTriangulations[0].ProductCode
	successionratiosAverageClaim.RunDate = cumulativeTriangulations[0].RunDate
	successionratiosAverageClaim.RunID = cumulativeTriangulations[0].RunID
	successionratiosAverageClaim.DevelopmentVariable = "Weighted Ave. Succession Ratio(rd to rd+1)"

	mc := reflect.ValueOf(&columnSum).Elem()
	mcl := reflect.ValueOf(&columnsumexcllastvalue).Elem()
	msr := reflect.ValueOf(&successionratios).Elem()

	mcCC := reflect.ValueOf(&columnSumClaimCount).Elem()
	mclCC := reflect.ValueOf(&columnsumexcllastvalueClaimCount).Elem()
	msrCC := reflect.ValueOf(&successionratiosClaimCount).Elem()

	mcAC := reflect.ValueOf(&columnSumAverageClaim).Elem()
	mclAC := reflect.ValueOf(&columnsumexcllastvalueAverageClaim).Elem()
	msrAC := reflect.ValueOf(&successionratiosAverageClaim).Elem()

	for k := 0; k <= projMonths+FactorVariableCount; k++ {
		if k > FactorVariableCount {
			if k == projMonths+FactorVariableCount {
				msr.Field(k).SetFloat(1)
				msrCC.Field(k).SetFloat(1)
				msrAC.Field(k).SetFloat(1)
			} else {
				if mcl.Field(k).Float() > 0 && k < projMonths+FactorVariableCount {
					msr.Field(k).SetFloat(mc.Field(k+1).Float() / mcl.Field(k).Float())
					msrCC.Field(k).SetFloat(mcCC.Field(k+1).Float() / mclCC.Field(k).Float())
					msrAC.Field(k).SetFloat(mcAC.Field(k+1).Float() / mclAC.Field(k).Float())
				}
				if mcl.Field(k).Float() == 0 && k < projMonths+FactorVariableCount {
					msr.Field(k).SetFloat(1)
					msrCC.Field(k).SetFloat(1)
					msrAC.Field(k).SetFloat(1)
				}
			}
		}
	}
	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, successionratios.DevelopmentVariable, successionratios.ProductCode).Delete(&models.LicDevelopmentFactor{})
	DB.Save(&successionratios)

	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, successionratiosClaimCount.DevelopmentVariable, successionratiosClaimCount.ProductCode).Delete(&models.LicDevelopmentFactorClaimCount{})
	DB.Save(&successionratiosClaimCount)

	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, successionratiosAverageClaim.DevelopmentVariable, successionratiosAverageClaim.ProductCode).Delete(&models.LicDevelopmentFactorAverageClaim{})
	DB.Save(&successionratiosAverageClaim)

	//Development factor
	//var developmentfactor models.LicDevelopmentFactor

	var manualDevelopmentFactor models.LicDevelopmentFactor
	var developmentfactorClaimCount models.LicDevelopmentFactorClaimCount
	var developmentfactorAverageClaim models.LicDevelopmentFactorAverageClaim

	developmentFactor.PortfolioName = cumulativeTriangulations[0].PortfolioName
	developmentFactor.ProductCode = cumulativeTriangulations[0].ProductCode
	developmentFactor.RunDate = cumulativeTriangulations[0].RunDate
	developmentFactor.RunID = cumulativeTriangulations[0].RunID
	developmentFactor.DevelopmentVariable = "Development Factor"

	manualDevelopmentFactor.PortfolioName = cumulativeTriangulations[0].PortfolioName
	manualDevelopmentFactor.ProductCode = cumulativeTriangulations[0].ProductCode
	manualDevelopmentFactor.RunDate = cumulativeTriangulations[0].RunDate
	manualDevelopmentFactor.RunID = cumulativeTriangulations[0].RunID
	manualDevelopmentFactor.DevelopmentVariable = "Manual Development Factor"

	developmentfactorClaimCount.PortfolioName = cumulativeTriangulations[0].PortfolioName
	developmentfactorClaimCount.ProductCode = cumulativeTriangulations[0].ProductCode
	developmentfactorClaimCount.RunDate = cumulativeTriangulations[0].RunDate
	developmentfactorClaimCount.RunID = cumulativeTriangulations[0].RunID
	developmentfactorClaimCount.DevelopmentVariable = "Development Factor"

	developmentfactorAverageClaim.PortfolioName = cumulativeTriangulations[0].PortfolioName
	developmentfactorAverageClaim.ProductCode = cumulativeTriangulations[0].ProductCode
	developmentfactorAverageClaim.RunDate = cumulativeTriangulations[0].RunDate
	developmentfactorAverageClaim.RunID = cumulativeTriangulations[0].RunID
	developmentfactorAverageClaim.DevelopmentVariable = "Development Factor"

	msr3 := reflect.ValueOf(&successionratios).Elem()
	mdf := reflect.ValueOf(developmentFactor).Elem()

	msr3CC := reflect.ValueOf(&successionratiosClaimCount).Elem()
	mdfCC := reflect.ValueOf(&developmentfactorClaimCount).Elem()

	msr3AC := reflect.ValueOf(&successionratiosAverageClaim).Elem()
	mdfAC := reflect.ValueOf(&developmentfactorAverageClaim).Elem()

	for k := 0; k <= projMonths+FactorVariableCount; k++ {
		if k > FactorVariableCount {
			if k == FactorVariableCount+1 {
				mdf.Field(k).SetFloat(0)
				mdfCC.Field(k).SetFloat(0)
				mdfAC.Field(k).SetFloat(0)
			} else {
				mdf.Field(k).SetFloat(msr3.Field(k - 1).Float())
				mdfCC.Field(k).SetFloat(msr3CC.Field(k - 1).Float())
				mdfAC.Field(k).SetFloat(msr3AC.Field(k - 1).Float())
			}

		}
	}
	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, developmentfactor.DevelopmentVariable, developmentfactor.ProductCode).Delete(&models.LicDevelopmentFactor{})
	err := DB.Save(&developmentFactor).Error
	if err != nil {
		log.Println(err)
	}
	if run.RerunIndicator && productCode == manualProdCode {
		if len(manualDfs) > 0 {
			err := copier.Copy(&manualDevelopmentFactor, &manualDfs[0])
			if err != nil {
				fmt.Println("return errors here")
			}
			manualDevelopmentFactor.DevelopmentVariable = "Manual Development Factor"
			manualDevelopmentFactor.RunDate = run.RunDate
			manualDevelopmentFactor.RunID = run.ID
			manualDevelopmentFactor.PortfolioName = run.PortfolioName
			manualDevelopmentFactor.ProductCode = manualProdCode
			err = DB.Save(&manualDevelopmentFactor).Error
			if err != nil {
				log.Println(err)
			}
		}
	}

	fmt.Println("manualDevelopmentFactor", manualDevelopmentFactor)

	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, developmentfactorClaimCount.DevelopmentVariable, developmentfactorClaimCount.ProductCode).Delete(&models.LicDevelopmentFactorClaimCount{})
	err = DB.Save(&developmentfactorClaimCount).Error
	if err != nil {
		log.Println(err)
	}

	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, developmentfactorAverageClaim.DevelopmentVariable, developmentfactorAverageClaim.ProductCode).Delete(&models.LicDevelopmentFactorAverageClaim{})
	err = DB.Save(&developmentfactorAverageClaim).Error
	if err != nil {
		log.Println(err)
	}

	//Cumulative Development Factors
	var cumulativedevfactors models.LicDevelopmentFactor
	var cumulativedevfactorsClaimCount models.LicDevelopmentFactorClaimCount
	var cumulativedevfactorsAverageClaim models.LicDevelopmentFactorAverageClaim
	cumulativedevfactors.PortfolioName = cumulativeTriangulations[0].PortfolioName
	cumulativedevfactors.ProductCode = cumulativeTriangulations[0].ProductCode
	cumulativedevfactors.RunDate = cumulativeTriangulations[0].RunDate
	cumulativedevfactors.RunID = cumulativeTriangulations[0].RunID
	cumulativedevfactors.DevelopmentVariable = "Cumulative Development Factors"

	cumulativedevfactorsClaimCount.PortfolioName = cumulativeTriangulations[0].PortfolioName
	cumulativedevfactorsClaimCount.ProductCode = cumulativeTriangulations[0].ProductCode
	cumulativedevfactorsClaimCount.RunDate = cumulativeTriangulations[0].RunDate
	cumulativedevfactorsClaimCount.RunID = cumulativeTriangulations[0].RunID
	cumulativedevfactorsClaimCount.DevelopmentVariable = "Cumulative Development Factors"

	cumulativedevfactorsAverageClaim.PortfolioName = cumulativeTriangulations[0].PortfolioName
	cumulativedevfactorsAverageClaim.ProductCode = cumulativeTriangulations[0].ProductCode
	cumulativedevfactorsAverageClaim.RunDate = cumulativeTriangulations[0].RunDate
	cumulativedevfactorsAverageClaim.RunID = cumulativeTriangulations[0].RunID
	cumulativedevfactorsAverageClaim.DevelopmentVariable = "Cumulative Development Factors"

	mdf2 := reflect.ValueOf(developmentFactor).Elem()
	mmdf2 := reflect.ValueOf(&manualDevelopmentFactor).Elem()
	mcdf := reflect.ValueOf(&cumulativedevfactors).Elem()

	mdf2CC := reflect.ValueOf(&developmentfactorClaimCount).Elem()
	mcdfCC := reflect.ValueOf(&cumulativedevfactorsClaimCount).Elem()

	mdf2AC := reflect.ValueOf(&developmentfactorAverageClaim).Elem()
	mcdfAC := reflect.ValueOf(&cumulativedevfactorsAverageClaim).Elem()

	for k := 0; k <= projMonths+FactorVariableCount; k++ {
		if k == FactorVariableCount+1 {
			mcdf.Field(k).SetFloat(1)
			mcdfCC.Field(k).SetFloat(1)
			mcdfAC.Field(k).SetFloat(1)
		}
		if k > FactorVariableCount+1 {
			if run.RerunIndicator && productCode == manualProdCode {
				mcdf.Field(k).SetFloat(mcdf.Field(k-1).Float() * mmdf2.Field(k).Float())
			} else {
				mcdf.Field(k).SetFloat(mcdf.Field(k-1).Float() * mdf2.Field(k).Float())
			}

			mcdfCC.Field(k).SetFloat(mcdfCC.Field(k-1).Float() * mdf2CC.Field(k).Float())
			mcdfAC.Field(k).SetFloat(mcdfAC.Field(k-1).Float() * mdf2AC.Field(k).Float())
		}
	}
	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, cumulativedevfactors.DevelopmentVariable, cumulativedevfactors.ProductCode).Delete(&models.LicDevelopmentFactor{})
	DB.Save(&cumulativedevfactors)

	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, cumulativedevfactorsClaimCount.DevelopmentVariable, cumulativedevfactorsClaimCount.ProductCode).Delete(&models.LicDevelopmentFactorClaimCount{})
	DB.Save(&cumulativedevfactorsClaimCount)

	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, cumulativedevfactorsAverageClaim.DevelopmentVariable, cumulativedevfactorsAverageClaim.ProductCode).Delete(&models.LicDevelopmentFactorAverageClaim{})
	DB.Save(&cumulativedevfactorsAverageClaim)

	//Cumulative Proportion
	var cumulativeDevelopmentProportion models.LicDevelopmentFactor
	cumulativeDevelopmentProportion.PortfolioName = cumulativeTriangulations[0].PortfolioName
	cumulativeDevelopmentProportion.ProductCode = cumulativeTriangulations[0].ProductCode
	cumulativeDevelopmentProportion.RunDate = cumulativeTriangulations[0].RunDate
	cumulativeDevelopmentProportion.RunID = cumulativeTriangulations[0].RunID
	cumulativeDevelopmentProportion.DevelopmentVariable = "Cumulative Development Proportion"

	mcdf2 := reflect.ValueOf(&cumulativedevfactors).Elem()
	mcdfp := reflect.ValueOf(&cumulativeDevelopmentProportion).Elem()

	for k := projMonths + FactorVariableCount; k > FactorVariableCount; k-- {
		if k == projMonths+FactorVariableCount {
			mcdfp.Field(k).SetFloat(1)
		}
		if k < projMonths+FactorVariableCount {
			if mcdf2.Field(projMonths+FactorVariableCount).Float() > 0 {
				mcdfp.Field(k).SetFloat(mcdf2.Field(k).Float() / mcdf2.Field(projMonths+FactorVariableCount).Float())

			}
		}
	}
	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, cumulativeDevelopmentProportion.DevelopmentVariable, cumulativeDevelopmentProportion.ProductCode).Delete(&models.LicDevelopmentFactor{})
	DB.Save(&cumulativeDevelopmentProportion)

	//Incremental Proportion
	var incrementalDevelopmentProportion models.LicDevelopmentFactor
	incrementalDevelopmentProportion.PortfolioName = cumulativeTriangulations[0].PortfolioName
	incrementalDevelopmentProportion.ProductCode = cumulativeTriangulations[0].ProductCode
	incrementalDevelopmentProportion.RunDate = cumulativeTriangulations[0].RunDate
	incrementalDevelopmentProportion.RunID = cumulativeTriangulations[0].RunID
	incrementalDevelopmentProportion.DevelopmentVariable = "Incremental Development Proportion"

	midfp := reflect.ValueOf(&incrementalDevelopmentProportion).Elem()
	mcdfp2 := reflect.ValueOf(&cumulativeDevelopmentProportion).Elem()

	for k := projMonths + FactorVariableCount; k > FactorVariableCount; k-- {
		if k <= projMonths+FactorVariableCount {
			if k == FactorVariableCount+1 {
				midfp.Field(k).SetFloat(mcdfp2.Field(k).Float())
			} else {
				midfp.Field(k).SetFloat(mcdfp2.Field(k).Float() - mcdfp2.Field(k-1).Float())
			}
		}
	}
	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, incrementalDevelopmentProportion.DevelopmentVariable, incrementalDevelopmentProportion.ProductCode).Delete(&models.LicDevelopmentFactor{})
	DB.Save(&incrementalDevelopmentProportion)

	//inverse
	var inverse models.LicDevelopmentFactor
	var inverseClaimCount models.LicDevelopmentFactorClaimCount
	var inverseAverageClaim models.LicDevelopmentFactorAverageClaim

	inverse.PortfolioName = cumulativeTriangulations[0].PortfolioName
	inverse.ProductCode = cumulativeTriangulations[0].ProductCode
	inverse.RunDate = cumulativeTriangulations[0].RunDate
	inverse.RunID = cumulativeTriangulations[0].RunID
	inverse.DevelopmentVariable = "inverse"

	inverseClaimCount.PortfolioName = cumulativeTriangulations[0].PortfolioName
	inverseClaimCount.ProductCode = cumulativeTriangulations[0].ProductCode
	inverseClaimCount.RunDate = cumulativeTriangulations[0].RunDate
	inverseClaimCount.RunID = cumulativeTriangulations[0].RunID
	inverseClaimCount.DevelopmentVariable = "inverse"

	inverseAverageClaim.PortfolioName = cumulativeTriangulations[0].PortfolioName
	inverseAverageClaim.ProductCode = cumulativeTriangulations[0].ProductCode
	inverseAverageClaim.RunDate = cumulativeTriangulations[0].RunDate
	inverseAverageClaim.RunID = cumulativeTriangulations[0].RunID
	inverseAverageClaim.DevelopmentVariable = "inverse"

	mcdf3 := reflect.ValueOf(&cumulativedevfactors).Elem()
	mi := reflect.ValueOf(&inverse).Elem()

	mcdf2CC := reflect.ValueOf(&cumulativedevfactorsClaimCount).Elem()
	miCC := reflect.ValueOf(&inverseClaimCount).Elem()

	mcdf2AC := reflect.ValueOf(&cumulativedevfactorsAverageClaim).Elem()
	miAC := reflect.ValueOf(&inverseAverageClaim).Elem()

	for k := 0; k <= projMonths+FactorVariableCount; k++ {
		if k > FactorVariableCount && mcdf3.Field(k).Float() > 0 {
			mi.Field(k).SetFloat(1 / mcdf3.Field(k).Float())
			miCC.Field(k).SetFloat(1 / mcdf2CC.Field(k).Float())
			miAC.Field(k).SetFloat(1 / mcdf2AC.Field(k).Float())
		}
	}
	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, inverse.DevelopmentVariable, inverse.ProductCode).Delete(&models.LicDevelopmentFactor{})
	DB.Save(&inverse)

	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, inverseClaimCount.DevelopmentVariable, inverseClaimCount.ProductCode).Delete(&models.LicDevelopmentFactorClaimCount{})
	DB.Save(&inverseClaimCount)

	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, inverseAverageClaim.DevelopmentVariable, inverseAverageClaim.ProductCode).Delete(&models.LicDevelopmentFactorAverageClaim{})
	DB.Save(&inverseAverageClaim)

	//year-to-year proportion

	var yearToYearRunoff models.LicDevelopmentFactor
	var yearToYearRunoffClaimCount models.LicDevelopmentFactorClaimCount
	var yearToYearRunoffAverageClaim models.LicDevelopmentFactorAverageClaim

	yearToYearRunoff.PortfolioName = cumulativeTriangulations[0].PortfolioName
	yearToYearRunoff.ProductCode = cumulativeTriangulations[0].ProductCode
	yearToYearRunoff.RunDate = cumulativeTriangulations[0].RunDate
	yearToYearRunoff.RunID = cumulativeTriangulations[0].RunID
	yearToYearRunoff.DevelopmentVariable = "Year to Year Run off"

	yearToYearRunoffClaimCount.PortfolioName = cumulativeTriangulations[0].PortfolioName
	yearToYearRunoffClaimCount.ProductCode = cumulativeTriangulations[0].ProductCode
	yearToYearRunoffClaimCount.RunDate = cumulativeTriangulations[0].RunDate
	yearToYearRunoffClaimCount.RunID = cumulativeTriangulations[0].RunID
	yearToYearRunoffClaimCount.DevelopmentVariable = "Year to Year Run off"

	yearToYearRunoffAverageClaim.PortfolioName = cumulativeTriangulations[0].PortfolioName
	yearToYearRunoffAverageClaim.ProductCode = cumulativeTriangulations[0].ProductCode
	yearToYearRunoffAverageClaim.RunDate = cumulativeTriangulations[0].RunDate
	yearToYearRunoffAverageClaim.RunID = cumulativeTriangulations[0].RunID
	yearToYearRunoffAverageClaim.DevelopmentVariable = "Year to Year Run off"

	mdf3 := reflect.ValueOf(developmentFactor).Elem()
	mmdf3 := reflect.ValueOf(&manualDevelopmentFactor).Elem()
	myr := reflect.ValueOf(&yearToYearRunoff).Elem()

	mdf3CC := reflect.ValueOf(&developmentfactorClaimCount).Elem()
	myrCC := reflect.ValueOf(&yearToYearRunoffClaimCount).Elem()

	mdf3AC := reflect.ValueOf(&developmentfactorAverageClaim).Elem()
	myrAC := reflect.ValueOf(&yearToYearRunoffAverageClaim).Elem()

	for k := 0; k <= projMonths+FactorVariableCount; k++ {
		if k > FactorVariableCount {
			if k == projMonths+FactorVariableCount {
				myr.Field(k).SetFloat(1)
				myrCC.Field(k).SetFloat(1)
				myrAC.Field(k).SetFloat(1)
			}
			if k < projMonths+FactorVariableCount && mdf2.Field(k+1).Float() > 0 {
				if run.RerunIndicator && productCode == manualProdCode {
					if mmdf3.Field(k+1).Float() != 0 {
						myr.Field(k).SetFloat(1 / mmdf3.Field(k+1).Float()) //references manual development factors
					}
					if mdf3CC.Field(k+1).Float() != 0 {
						myrCC.Field(k).SetFloat(1 / mdf3CC.Field(k+1).Float())
					}
					if mdf3AC.Field(k+1).Float() != 0 {
						myrAC.Field(k).SetFloat(1 / mdf3AC.Field(k+1).Float())
					}
				} else {
					if mdf3.Field(k+1).Float() != 0 {
						myr.Field(k).SetFloat(1 / mdf3.Field(k+1).Float())
					}
					if mdf3CC.Field(k+1).Float() != 0 {
						myrCC.Field(k).SetFloat(1 / mdf3CC.Field(k+1).Float())
					}
					if mdf3AC.Field(k+1).Float() != 0 {
						myrAC.Field(k).SetFloat(1 / mdf3AC.Field(k+1).Float())
					}
				}
			}
		}
	}
	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, yearToYearRunoff.DevelopmentVariable, yearToYearRunoff.ProductCode).Delete(&models.LicDevelopmentFactor{})
	err = DB.Save(&yearToYearRunoff).Error
	if err != nil {
		log.Println(err)
	}

	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, yearToYearRunoffClaimCount.DevelopmentVariable, yearToYearRunoffClaimCount.ProductCode).Delete(&models.LicDevelopmentFactorClaimCount{})
	err = DB.Save(&yearToYearRunoffClaimCount).Error
	if err != nil {
		log.Println(err)
	}

	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, yearToYearRunoffAverageClaim.DevelopmentVariable, yearToYearRunoffAverageClaim.ProductCode).Delete(&models.LicDevelopmentFactorAverageClaim{})
	err = DB.Save(&yearToYearRunoffAverageClaim).Error
	if err != nil {
		log.Println(err)
	}

	//proportion runoff
	var proportionRunoff models.LicDevelopmentFactor
	var proportionRunoffClaimCount models.LicDevelopmentFactorClaimCount
	var proportionRunoffAverageClaim models.LicDevelopmentFactorAverageClaim

	proportionRunoff.PortfolioName = cumulativeTriangulations[0].PortfolioName
	proportionRunoff.ProductCode = cumulativeTriangulations[0].ProductCode
	proportionRunoff.RunDate = cumulativeTriangulations[0].RunDate
	proportionRunoff.RunID = cumulativeTriangulations[0].RunID
	proportionRunoff.DevelopmentVariable = "proportion runoff"

	proportionRunoffClaimCount.PortfolioName = cumulativeTriangulations[0].PortfolioName
	proportionRunoffClaimCount.ProductCode = cumulativeTriangulations[0].ProductCode
	proportionRunoffClaimCount.RunDate = cumulativeTriangulations[0].RunDate
	proportionRunoffClaimCount.RunID = cumulativeTriangulations[0].RunID
	proportionRunoffClaimCount.DevelopmentVariable = "proportion runoff"

	proportionRunoffAverageClaim.PortfolioName = cumulativeTriangulations[0].PortfolioName
	proportionRunoffAverageClaim.ProductCode = cumulativeTriangulations[0].ProductCode
	proportionRunoffAverageClaim.RunDate = cumulativeTriangulations[0].RunDate
	proportionRunoffAverageClaim.RunID = cumulativeTriangulations[0].RunID
	proportionRunoffAverageClaim.DevelopmentVariable = "proportion runoff"

	mdf4 := reflect.ValueOf(developmentFactor).Elem()
	mmdf4 := reflect.ValueOf(&manualDevelopmentFactor).Elem()

	mpr := reflect.ValueOf(&proportionRunoff).Elem()

	mdf4CC := reflect.ValueOf(&developmentfactorClaimCount).Elem()
	mprCC := reflect.ValueOf(&proportionRunoffClaimCount).Elem()

	mdf4AC := reflect.ValueOf(&developmentfactorAverageClaim).Elem()
	mprAC := reflect.ValueOf(&proportionRunoffAverageClaim).Elem()

	for k := projMonths + FactorVariableCount; k > FactorVariableCount; k-- {
		if k == projMonths+FactorVariableCount {
			mpr.Field(k).SetFloat(1)
			mprCC.Field(k).SetFloat(1)
			mprAC.Field(k).SetFloat(1)
		} else {
			// manual development factor
			if run.RerunIndicator && productCode == manualProdCode {
				if mmdf4.Field(k).Float() > 0 && mmdf4.Field(k+1).Float() == 0 {
					mpr.Field(k).SetFloat(1)
					mprCC.Field(k).SetFloat(1)
					mprAC.Field(k).SetFloat(1)
				}
				if mmdf4.Field(k+1).Float() > 0 {
					mpr.Field(k).SetFloat(utils.FloatPrecision(mpr.Field(k+1).Float()/mmdf4.Field(k+1).Float(), ibnrRatioPrecision))
					mprCC.Field(k).SetFloat(utils.FloatPrecision(mprCC.Field(k+1).Float()/mdf4CC.Field(k+1).Float(), ibnrRatioPrecision))
					mprAC.Field(k).SetFloat(utils.FloatPrecision(mprAC.Field(k+1).Float()/mdf4AC.Field(k+1).Float(), ibnrRatioPrecision))
				}
			}

			// base run, not a manual run
			if !run.RerunIndicator {
				if mdf4.Field(k).Float() > 0 && mdf4.Field(k+1).Float() == 0 {
					mpr.Field(k).SetFloat(1)
					mprCC.Field(k).SetFloat(1)
					mprAC.Field(k).SetFloat(1)
				}
				if mdf4.Field(k+1).Float() > 0 {
					mpr.Field(k).SetFloat(utils.FloatPrecision(mpr.Field(k+1).Float()/mdf4.Field(k+1).Float(), ibnrRatioPrecision))
					mprCC.Field(k).SetFloat(utils.FloatPrecision(mprCC.Field(k+1).Float()/mdf4CC.Field(k+1).Float(), ibnrRatioPrecision))
					mprAC.Field(k).SetFloat(utils.FloatPrecision(mprAC.Field(k+1).Float()/mdf4AC.Field(k+1).Float(), ibnrRatioPrecision))
				}
			}

		}
	}
	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, proportionRunoff.DevelopmentVariable, proportionRunoff.ProductCode).Delete(&models.LicDevelopmentFactor{})
	DB.Save(&proportionRunoff)

	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, proportionRunoffClaimCount.DevelopmentVariable, proportionRunoffClaimCount.ProductCode).Delete(&models.LicDevelopmentFactorClaimCount{})
	DB.Save(&proportionRunoffClaimCount)

	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, proportionRunoffAverageClaim.DevelopmentVariable, proportionRunoffAverageClaim.ProductCode).Delete(&models.LicDevelopmentFactorAverageClaim{})
	DB.Save(&proportionRunoffAverageClaim)

	//proportion not runoff
	var proportionNotRunoff models.LicDevelopmentFactor
	var proportionNotRunoffClaimCount models.LicDevelopmentFactorClaimCount
	var proportionNotRunoffAverageClaim models.LicDevelopmentFactorAverageClaim

	proportionNotRunoff.PortfolioName = cumulativeTriangulations[0].PortfolioName
	proportionNotRunoff.ProductCode = cumulativeTriangulations[0].ProductCode
	proportionNotRunoff.RunDate = cumulativeTriangulations[0].RunDate
	proportionNotRunoff.RunID = cumulativeTriangulations[0].RunID
	proportionNotRunoff.DevelopmentVariable = "proportion not runoff"

	proportionNotRunoffClaimCount.PortfolioName = cumulativeTriangulations[0].PortfolioName
	proportionNotRunoffClaimCount.ProductCode = cumulativeTriangulations[0].ProductCode
	proportionNotRunoffClaimCount.RunDate = cumulativeTriangulations[0].RunDate
	proportionNotRunoffClaimCount.RunID = cumulativeTriangulations[0].RunID
	proportionNotRunoffClaimCount.DevelopmentVariable = "proportion not runoff"

	proportionNotRunoffAverageClaim.PortfolioName = cumulativeTriangulations[0].PortfolioName
	proportionNotRunoffAverageClaim.ProductCode = cumulativeTriangulations[0].ProductCode
	proportionNotRunoffAverageClaim.RunDate = cumulativeTriangulations[0].RunDate
	proportionNotRunoffAverageClaim.RunID = cumulativeTriangulations[0].RunID
	proportionNotRunoffAverageClaim.DevelopmentVariable = "proportion not runoff"

	mpr2 := reflect.ValueOf(&proportionRunoff).Elem()
	mpnr := reflect.ValueOf(&proportionNotRunoff).Elem()

	mpr2CC := reflect.ValueOf(&proportionRunoffClaimCount).Elem()
	mpnrCC := reflect.ValueOf(&proportionNotRunoffClaimCount).Elem()

	mpr2AC := reflect.ValueOf(&proportionRunoffAverageClaim).Elem()
	mpnrAC := reflect.ValueOf(&proportionNotRunoffAverageClaim).Elem()

	for k := 0; k <= projMonths+FactorVariableCount; k++ {
		if k > FactorVariableCount {
			mpnr.Field(k).SetFloat(1 - mpr2.Field(k).Float())
			mpnrCC.Field(k).SetFloat(1 - mpr2CC.Field(k).Float())
			mpnrAC.Field(k).SetFloat(1 - mpr2AC.Field(k).Float())
		}
	}

	DB.Where("run_date=? and development_variable=? and product_code=?", runDate, proportionNotRunoff.DevelopmentVariable, proportionNotRunoff.ProductCode).Delete(&models.LicDevelopmentFactor{})
	err = DB.Save(&proportionNotRunoff).Error
	if err != nil {
		log.Println(err)
	}

	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, proportionNotRunoffClaimCount.DevelopmentVariable, proportionNotRunoffClaimCount.ProductCode).Delete(&models.LicDevelopmentFactorClaimCount{})
	err = DB.Save(&proportionNotRunoffClaimCount).Error
	if err != nil {
		log.Println(err)
	}

	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, proportionNotRunoffAverageClaim.DevelopmentVariable, proportionNotRunoffAverageClaim.ProductCode).Delete(&models.LicDevelopmentFactorAverageClaim{})
	err = DB.Save(&proportionNotRunoffAverageClaim).Error
	if err != nil {
		log.Println(err)
	}

	//projections
	var cumulativeProjections = []models.LicCumulativeProjection{}
	var cumulativeProjectionsClaimCount = []models.LicCumulativeProjectionClaimCount{}
	var cumulativeProjectionsAverageClaim = []models.LicCumulativeProjectionAverageClaim{}
	var cumulativeProjectionsAveragetoTotalClaim = []models.LicCumulativeProjectionAveragetoTotalClaim{}   //
	var incrementalProjectionsAveragetoTotalClaim = []models.LicIncrementalProjectionAveragetoTotalClaim{} //
	var inflatedincrementalProjectionsTotalClaim = []models.LicIncrementalInflatedAveragetoTotalClaim{}
	var discountedinflatedincrementalProjectionsTotalClaim = []models.LicDiscountedIncrementalInflatedAveragetoTotalClaim{}
	for i, _ := range cumulativeTriangulations {
		var cumulativeProjection models.LicCumulativeProjection
		var cumulativeProjectionClaimCount models.LicCumulativeProjectionClaimCount
		var cumulativeProjectionAverageClaim models.LicCumulativeProjectionAverageClaim
		var cumulativeProjectionAveragetoTotalClaim models.LicCumulativeProjectionAveragetoTotalClaim
		var incrementalProjectionAveragetoTotalClaim models.LicIncrementalProjectionAveragetoTotalClaim
		var inflatedincrementalProjectionTotalClaim models.LicIncrementalInflatedAveragetoTotalClaim
		var discountedinflatedincrementalProjectionTotalClaim models.LicDiscountedIncrementalInflatedAveragetoTotalClaim
		mct := reflect.ValueOf(&cumulativeTriangulations[i]).Elem()
		msr4 := reflect.ValueOf(&successionratios).Elem()
		mcp := reflect.ValueOf(&cumulativeProjection).Elem()

		mctCC := reflect.ValueOf(&cumulativeTriangulationsClaimCount[i]).Elem()
		msr4CC := reflect.ValueOf(&successionratiosClaimCount).Elem()
		mcpCC := reflect.ValueOf(&cumulativeProjectionClaimCount).Elem()

		mctAC := reflect.ValueOf(&cumulativeTriangulationsAverageClaim[i]).Elem()
		msr4AC := reflect.ValueOf(&successionratiosAverageClaim).Elem()
		mcpAC := reflect.ValueOf(&cumulativeProjectionAverageClaim).Elem() // projected average claim count

		mcpATC := reflect.ValueOf(&cumulativeProjectionAveragetoTotalClaim).Elem() // projected average claim

		mipATC := reflect.ValueOf(&incrementalProjectionAveragetoTotalClaim).Elem() // projected total claims based on projected average claim and projected claim count

		miipATC := reflect.ValueOf(&inflatedincrementalProjectionTotalClaim).Elem() //inflated projected total claims

		mdiipATC := reflect.ValueOf(&discountedinflatedincrementalProjectionTotalClaim).Elem() //discounted inflated  projected total claims

		var yearadd, movingyear int
		var residualmonth, movingmonth int
		for k := 0; k <= projMonths+TriangleVariableCount; k++ {
			if k <= TriangleVariableCount {
				mcp.Field(k).Set(mct.Field(k))
				mcpCC.Field(k).Set(mctCC.Field(k))
				mcpAC.Field(k).Set(mctAC.Field(k))
				mcpATC.Field(k).Set(mctAC.Field(k))
				mipATC.Field(k).Set(mctAC.Field(k))
				miipATC.Field(k).Set(mctAC.Field(k))
				mdiipATC.Field(k).Set(mctAC.Field(k))

			} else {
				//fmt.Println(mct.Type().Field(k).Name)
				//fmt.Println(projMonths)
				if interval == 12 {
					yearadd = int(int64((int(mct.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1) / (12 / interval)))
					residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mct.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				}
				if interval != 12 {
					yearadd = int(int64((int(mct.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1 - 1) / (12 / interval)))
					residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mct.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				}
				movingyear = int(mct.FieldByName(AccidentYear).Int()) + yearadd
				if residualmonth == 0 && interval == 3 {
					residualmonth = 4
				}
				if residualmonth == 0 && interval == 1 {
					residualmonth = 12
				}

				movingmonth = residualmonth

				if k == TriangleVariableCount+1 {
					mcp.Field(k).SetFloat(utils.FloatPrecision(mct.Field(k).Float(), AccountingPrecision))
					mcpCC.Field(k).SetFloat(utils.FloatPrecision(mctCC.Field(k).Float(), AccountingPrecision))
					mcpAC.Field(k).SetFloat(utils.FloatPrecision(mctAC.Field(k).Float(), AccountingPrecision))
					mcpATC.Field(k).SetFloat(utils.FloatPrecision(mcpCC.Field(k).Float()*mcpAC.Field(k).Float(), AccountingPrecision))
					mipATC.Field(k).SetFloat(utils.FloatPrecision(mcpATC.Field(k).Float(), AccountingPrecision))
				} else {
					if (movingyear <= endDateYear && movingmonth <= endDateMonth) || (movingyear < endDateYear && movingmonth > endDateMonth) {
						mcp.Field(k).SetFloat(utils.FloatPrecision(mct.Field(k).Float(), AccountingPrecision))
						mcpCC.Field(k).SetFloat(utils.FloatPrecision(mctCC.Field(k).Float(), AccountingPrecision))
						mcpAC.Field(k).SetFloat(utils.FloatPrecision(mctAC.Field(k).Float(), AccountingPrecision))
						mcpATC.Field(k).SetFloat(utils.FloatPrecision(mcpCC.Field(k).Float()*mcpAC.Field(k).Float(), AccountingPrecision))
						mipATC.Field(k).SetFloat(utils.FloatPrecision(mcpATC.Field(k).Float()-mcpATC.Field(k-1).Float(), AccountingPrecision))
					} else {
						if k <= projMonths+TriangleVariableCount {
							//fmt.Println(mcp.Field(k - 1).Float())
							//fmt.Println(msr4.Field(k - 3 - 1).Float())
							mcp.Field(k).SetFloat(utils.FloatPrecision(mcp.Field(k-1).Float()*msr4.Field(k-VariableCountDiff-1).Float(), AccountingPrecision))
							mcpCC.Field(k).SetFloat(utils.FloatPrecision(mcpCC.Field(k-1).Float()*msr4CC.Field(k-VariableCountDiff-1).Float(), AccountingPrecision))
							mcpAC.Field(k).SetFloat(utils.FloatPrecision(mcpAC.Field(k-1).Float()*msr4AC.Field(k-VariableCountDiff-1).Float(), AccountingPrecision))
							mcpATC.Field(k).SetFloat(utils.FloatPrecision(mcpCC.Field(k).Float()*mcpAC.Field(k).Float(), AccountingPrecision))
							mipATC.Field(k).SetFloat(utils.FloatPrecision(mcpATC.Field(k).Float()-mcpATC.Field(k-1).Float(), AccountingPrecision))
							miipATC.Field(k).SetFloat(utils.FloatPrecision(mipATC.Field(k).Float()*1, AccountingPrecision))
							mdiipATC.Field(k).SetFloat(utils.FloatPrecision(miipATC.Field(k).Float()*1, AccountingPrecision))
						}

					}
					//}else{
					////if (movingyear == endDateYear && movingmonth > endDateMonth) || (movingyear > endDateYear) {
					//	mcp.Field(k).SetFloat(mct.Field(k-1).Float() * msr4.Field(k-1).Float())
					////}
					//}
				}
			}
		}
		cumulativeProjections = append(cumulativeProjections, cumulativeProjection)
		cumulativeProjectionsClaimCount = append(cumulativeProjectionsClaimCount, cumulativeProjectionClaimCount)
		cumulativeProjectionsAverageClaim = append(cumulativeProjectionsAverageClaim, cumulativeProjectionAverageClaim)
		cumulativeProjectionsAveragetoTotalClaim = append(cumulativeProjectionsAveragetoTotalClaim, cumulativeProjectionAveragetoTotalClaim)
		incrementalProjectionsAveragetoTotalClaim = append(incrementalProjectionsAveragetoTotalClaim, incrementalProjectionAveragetoTotalClaim)
		inflatedincrementalProjectionsTotalClaim = append(inflatedincrementalProjectionsTotalClaim, inflatedincrementalProjectionTotalClaim)
		discountedinflatedincrementalProjectionsTotalClaim = append(discountedinflatedincrementalProjectionsTotalClaim, discountedinflatedincrementalProjectionTotalClaim)
	}

	// Incremental Projections, inflated and discounted accordingly where applicable
	adjustedCumulativeProjections := InflatedDiscountedProjections(&cumulativeProjections, run, runDate, projMonths, endDateYear, endDateMonth, interval, discountOption, params.YieldCurveCode)

	//DB.Where("run_date=? and product_code=?", runDate, adjustedCumulativeProjections[0].ProductCode).Delete(&models.LicCumulativeProjection{})
	err = DB.Save(&adjustedCumulativeProjections).Error
	if err != nil {
		log.Println(err)
	}

	//DB.Where("run_date=? and product_code=?", runDate, cumulativeProjectionsClaimCount[0].ProductCode).Delete(&models.LicCumulativeProjectionClaimCount{})
	err = DB.Save(&cumulativeProjectionsClaimCount).Error
	if err != nil {
		log.Println(err)
	}

	//DB.Where("run_date=? and product_code=?", runDate, cumulativeProjectionsAverageClaim[0].ProductCode).Delete(&models.LicCumulativeProjectionAverageClaim{})
	err = DB.Save(&cumulativeProjectionsAverageClaim).Error
	if err != nil {
		log.Println(err)
	}

	//DB.Where("run_date=? and product_code=?", runDate, cumulativeProjectionsAveragetoTotalClaim[0].ProductCode).Delete(&models.LicCumulativeProjectionAveragetoTotalClaim{})
	err = DB.Save(&cumulativeProjectionsAveragetoTotalClaim).Error
	if err != nil {
		log.Println(err)
	}

	//DB.Where("run_date=? and product_code=?", runDate, incrementalProjectionsAveragetoTotalClaim[0].ProductCode).Delete(&models.LicIncrementalProjectionAveragetoTotalClaim{})
	err = DB.Save(&incrementalProjectionsAveragetoTotalClaim).Error
	if err != nil {
		log.Println(err)
	}

	//DB.Where("run_date=? and product_code=?", runDate, inflatedincrementalProjectionsTotalClaim[0].ProductCode).Delete(&models.LicIncrementalInflatedAveragetoTotalClaim{})
	err = DB.Save(&inflatedincrementalProjectionsTotalClaim).Error
	if err != nil {
		log.Println(err)
	}

	//DB.Where("run_date=? and product_code=?", runDate, discountedinflatedincrementalProjectionsTotalClaim[0].ProductCode).Delete(&models.LicDiscountedIncrementalInflatedAveragetoTotalClaim{})
	err = DB.Save(&discountedinflatedincrementalProjectionsTotalClaim).Error
	if err != nil {
		log.Println(err)
	}
	var bootstrapmeanOut float64
	var NthPercentileOut float64
	if run.BootStrapIndicator {

		// block's global variables

		// LicExpectedSimulation; cumulative
		var expectedSimulations = []models.LicExpectedSimulation{}
		for i, _ := range adjustedCumulativeProjections {
			var expectedSimulation models.LicExpectedSimulation
			mcp2 := reflect.ValueOf(&adjustedCumulativeProjections[i]).Elem()
			mes := reflect.ValueOf(&expectedSimulation).Elem()
			mcdfp3 := reflect.ValueOf(&cumulativeDevelopmentProportion).Elem()

			var yearadd, movingyear int
			var residualmonth, movingmonth int

			for k := 0; k <= projMonths+TriangleVariableCount; k++ {
				if k <= TriangleVariableCount { //assigning variable names
					mes.Field(k).Set(mcp2.Field(k))
				} else {
					if interval == 12 {
						yearadd = int(int64((int(mcp2.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1) / (12 / interval)))
						residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mcp2.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
					}
					if interval != 12 {
						yearadd = int(int64((int(mcp2.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1 - 1) / (12 / interval)))
						residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mcp2.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
					}
					movingyear = int(mcp2.FieldByName(AccidentYear).Int()) + yearadd
					if residualmonth == 0 && interval == 3 {
						residualmonth = 4
					}

					if residualmonth == 0 && interval == 1 {
						residualmonth = 12
					}

					movingmonth = residualmonth

					if (movingyear <= endDateYear && movingmonth <= endDateMonth) || (movingyear < endDateYear && movingmonth > endDateMonth) {
						mes.Field(k).SetFloat(mcp2.Field(projMonths+TriangleVariableCount).Float() * mcdfp3.Field(k-VariableCountDiff).Float())
					} else {
						mes.Field(k).SetFloat(0)
					}
				}
			}

			expectedSimulations = append(expectedSimulations, expectedSimulation)
		}

		//DB.Where("run_date=? and product_code=?", runDate, expectedSimulations[0].ProductCode).Delete(&models.LicExpectedSimulation{})
		err = DB.Save(&expectedSimulations).Error
		if err != nil {
			log.Println(err)
		}

		// LicExpectedSimulation; incremental
		var incrementalexpectedSimulations = []models.LicIncrementalExpectedSimulation{}
		for i, _ := range expectedSimulations {
			var incrementalexpectedSimulation models.LicIncrementalExpectedSimulation
			mes2 := reflect.ValueOf(&expectedSimulations[i]).Elem()
			mies := reflect.ValueOf(&incrementalexpectedSimulation).Elem()

			var yearadd, movingyear int
			var residualmonth, movingmonth int

			for k := 0; k <= projMonths+TriangleVariableCount; k++ {
				if k <= TriangleVariableCount { //assigning variable names
					mies.Field(k).Set(mes2.Field(k))
				}
				if k == TriangleVariableCount+1 {
					mies.Field(k).SetFloat(mes2.Field(k).Float())
				}
				if k > TriangleVariableCount+1 {

					if interval == 12 {
						yearadd = int(int64((int(mes2.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1) / (12 / interval)))
						residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mes2.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
					}
					if interval != 12 {
						yearadd = int(int64((int(mes2.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1 - 1) / (12 / interval)))
						residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mes2.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
					}
					movingyear = int(mes2.FieldByName(AccidentYear).Int()) + yearadd
					if residualmonth == 0 && interval == 3 {
						residualmonth = 4
					}
					if residualmonth == 0 && interval == 1 {
						residualmonth = 12
					}
					movingmonth = residualmonth

					if (movingyear <= endDateYear && movingmonth <= endDateMonth) || (movingyear < endDateYear && movingmonth > endDateMonth) {
						mies.Field(k).SetFloat(mes2.Field(k).Float() - mes2.Field(k-1).Float())
					} else {
						mies.Field(k).SetFloat(0)
					}
				}
			}

			incrementalexpectedSimulations = append(incrementalexpectedSimulations, incrementalexpectedSimulation)
		}

		//DB.Where("run_date=? and product_code=?", runDate, expectedSimulations[0].ProductCode).Delete(&models.LicExpectedSimulation{})
		err = DB.Save(&incrementalexpectedSimulations).Error
		if err != nil {
			log.Println(err)
		}

		// Standardised Residuals

		var standardisedResiduals = []models.LicStandardisedResiduals{}
		for i, _ := range incrementalexpectedSimulations {
			var standardisedResidual models.LicStandardisedResiduals
			mies2 := reflect.ValueOf(&incrementalexpectedSimulations[i]).Elem()
			mt2 := reflect.ValueOf(&triangulations[i]).Elem()
			mstdr := reflect.ValueOf(&standardisedResidual).Elem()

			var yearadd, movingyear int
			var residualmonth, movingmonth int

			for k := 0; k <= projMonths+TriangleVariableCount; k++ {
				if k <= TriangleVariableCount { //assigning variable names
					mstdr.Field(k).Set(mies2.Field(k))
				} else {
					if interval == 12 {
						yearadd = int(int64((int(mies2.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1) / (12 / interval)))
						residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mies2.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
					}
					if interval != 12 {
						yearadd = int(int64((int(mies2.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1 - 1) / (12 / interval)))
						residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mies2.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
					}
					movingyear = int(mies2.FieldByName(AccidentYear).Int()) + yearadd
					if residualmonth == 0 && interval == 3 {
						residualmonth = 4
					}
					if residualmonth == 0 && interval == 1 {
						residualmonth = 12
					}
					movingmonth = residualmonth
					if (movingyear <= endDateYear && movingmonth <= endDateMonth) || (movingyear < endDateYear && movingmonth > endDateMonth) {
						if mies2.Field(k).Float() > 0 {
							mstdr.Field(k).SetFloat((mt2.Field(k).Float() - mies2.Field(k).Float()) / math.Sqrt(mies2.Field(k).Float()))
						}
					} else {
						mstdr.Field(k).SetFloat(0)
					}
				}
			}
			standardisedResiduals = append(standardisedResiduals, standardisedResidual)
		}

		//DB.Where("run_date=? and product_code=?", runDate, standardisedResiduals[0].ProductCode).Delete(&models.LicStandardisedResiduals{})
		err = DB.Save(&standardisedResiduals).Error
		if err != nil {
			log.Println(err)
		}

		// Tabulating Residual Data

		var tabulatedRandomResiduals = []models.LicGeneratedRandomResiduals{}
		counter := 1
		for i, _ := range standardisedResiduals {
			var tabulatedRandomResidual models.LicGeneratedRandomResiduals
			mstdr2 := reflect.ValueOf(&standardisedResiduals[i]).Elem()
			mtrr := reflect.ValueOf(&tabulatedRandomResidual).Elem()
			var yearadd, movingyear int
			var residualmonth, movingmonth int

			for k := TriangleVariableCount + 1; k <= projMonths+TriangleVariableCount; k++ {

				if interval == 12 {
					yearadd = int(int64((int(mstdr2.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1) / (12 / interval)))
					residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mstdr2.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				}
				if interval != 12 {
					yearadd = int(int64((int(mstdr2.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1 - 1) / (12 / interval)))
					residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mstdr2.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				}
				movingyear = int(mstdr2.FieldByName(AccidentYear).Int()) + yearadd
				if residualmonth == 0 && interval == 3 {
					residualmonth = 4
				}
				if residualmonth == 0 && interval == 1 {
					residualmonth = 12
				}
				movingmonth = residualmonth
				if (movingyear <= endDateYear && movingmonth <= endDateMonth) || (movingyear < endDateYear && movingmonth > endDateMonth) {
					tabulatedRandomResidual.RunDate = run.RunDate
					tabulatedRandomResidual.RunId = run.ID
					tabulatedRandomResidual.PortfolioName = run.PortfolioName
					tabulatedRandomResidual.LicPortfolioId = run.PortfolioId
					tabulatedRandomResidual.ProductCode = standardisedResiduals[0].ProductCode
					mtrr.FieldByName(ResidualNumber).SetInt(int64(counter))
					mtrr.FieldByName(Residual).SetFloat(mstdr2.Field(k).Float())
					counter += 1
					tabulatedRandomResiduals = append(tabulatedRandomResiduals, tabulatedRandomResidual)
				}
			}
		}
		var bootStrappedResults = []models.LicBootstrappedResults{}
		rand.Seed(time.Now().UnixNano()) // Is this global....?
		for sim := 1; sim <= run.Simulations; sim++ {
			// Randomising Residuals
			randList := rand.Perm(counter - 1)
			for x, _ := range tabulatedRandomResiduals {
				pointer := randList[x]
				temp1 := reflect.ValueOf(&tabulatedRandomResiduals[pointer]).Elem()
				tempRandomResidual := reflect.ValueOf(&tabulatedRandomResiduals[x]).Elem()
				tempRandomResidual.FieldByName(RandomResidual).SetFloat(temp1.FieldByName(Residual).Float())
			}

			// Random Residuals Triangulated

			var randomResiduals = []models.LicRandomResiduals{}
			counter2 := 0
			for i, _ := range standardisedResiduals {
				var randomResidual models.LicRandomResiduals
				mstdr3 := reflect.ValueOf(&standardisedResiduals[i]).Elem()
				mrr := reflect.ValueOf(&randomResidual).Elem()
				var yearadd, movingyear int
				var residualmonth, movingmonth int

				for k := 0; k <= projMonths+TriangleVariableCount; k++ {

					if k <= TriangleVariableCount {
						mrr.Field(k).Set(mstdr3.Field(k))
					} else {
						if interval == 12 {
							yearadd = int(int64((int(mstdr3.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1) / (12 / interval)))
							residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mstdr3.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
						}
						if interval != 12 {
							yearadd = int(int64((int(mstdr3.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1 - 1) / (12 / interval)))
							residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mstdr3.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
						}
						movingyear = int(mstdr3.FieldByName(AccidentYear).Int()) + yearadd
						if residualmonth == 0 && interval == 3 {
							residualmonth = 4
						}
						if residualmonth == 0 && interval == 1 {
							residualmonth = 12
						}
						movingmonth = residualmonth
						if (movingyear <= endDateYear && movingmonth <= endDateMonth) || (movingyear < endDateYear && movingmonth > endDateMonth) {
							mtrr3 := reflect.ValueOf(&tabulatedRandomResiduals[counter2]).Elem()
							mrr.Field(k).SetFloat(mtrr3.FieldByName(RandomResidual).Float())
							counter2 += 1
						}
					}
				}
				randomResiduals = append(randomResiduals, randomResidual)
			}

			//saving triangulated generated random residuals
			if sim == run.Simulations {
				DB.Where("run_date=? and product_code=? and run_id=?", runDate, randomResiduals[0].ProductCode, randomResiduals[0].RunID).Delete(&models.LicRandomResiduals{})
				err = DB.Save(&randomResiduals).Error
				if err != nil {
					log.Println(err)
				}
			}

			// BootStrapped Incremental

			var bootStrappedIncrementals = []models.LicBootStrappedIncremental{}
			for i, _ := range randomResiduals {
				var bootStrappedIncremental models.LicBootStrappedIncremental
				mrr2 := reflect.ValueOf(&randomResiduals[i]).Elem()
				mies3 := reflect.ValueOf(&incrementalexpectedSimulations[i]).Elem()
				mbsi := reflect.ValueOf(&bootStrappedIncremental).Elem()
				var yearadd, movingyear int
				var residualmonth, movingmonth int

				for k := 0; k <= projMonths+TriangleVariableCount; k++ {

					if k <= TriangleVariableCount {
						mbsi.Field(k).Set(mrr2.Field(k))
					} else {
						if interval == 12 {
							yearadd = int(int64((int(mrr2.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1) / (12 / interval)))
							residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mrr2.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
						}
						if interval != 12 {
							yearadd = int(int64((int(mrr2.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1 - 1) / (12 / interval)))
							residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mrr2.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
						}
						movingyear = int(mrr2.FieldByName(AccidentYear).Int()) + yearadd
						if residualmonth == 0 && interval == 3 {
							residualmonth = 4
						}
						if residualmonth == 0 && interval == 1 {
							residualmonth = 12
						}
						movingmonth = residualmonth
						if (movingyear <= endDateYear && movingmonth <= endDateMonth) || (movingyear < endDateYear && movingmonth > endDateMonth) {
							mbsi.Field(k).SetFloat(mrr2.Field(k).Float()*math.Sqrt(mies3.Field(k).Float()) + mies3.Field(k).Float())
						}
					}
				}
				bootStrappedIncrementals = append(bootStrappedIncrementals, bootStrappedIncremental)
			}

			// BootStrapped Cumulative

			var bootStrappedCumulatives = []models.LicBootStrappedCumulative{}
			for i, _ := range bootStrappedIncrementals {
				var bootStrappedCumulative models.LicBootStrappedCumulative
				mbsi2 := reflect.ValueOf(&bootStrappedIncrementals[i]).Elem()
				mbsc := reflect.ValueOf(&bootStrappedCumulative).Elem()
				var yearadd, movingyear int
				var residualmonth, movingmonth int

				for k := 0; k <= projMonths+TriangleVariableCount; k++ {

					if k <= TriangleVariableCount {
						mbsc.Field(k).Set(mbsi2.Field(k))
					}
					if k == TriangleVariableCount+1 {
						mbsc.Field(k).SetFloat(mbsi2.Field(k).Float())
					}

					if k > TriangleVariableCount+1 {
						if interval == 12 {
							yearadd = int(int64((int(mbsi2.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1) / (12 / interval)))
							residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mbsi2.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
						}
						if interval != 12 {
							yearadd = int(int64((int(mbsi2.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1 - 1) / (12 / interval)))
							residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mbsi2.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
						}
						movingyear = int(mbsi2.FieldByName(AccidentYear).Int()) + yearadd
						if residualmonth == 0 && interval == 3 {
							residualmonth = 4
						}
						if residualmonth == 0 && interval == 1 {
							residualmonth = 12
						}
						movingmonth = residualmonth
						if (movingyear <= endDateYear && movingmonth <= endDateMonth) || (movingyear < endDateYear && movingmonth > endDateMonth) {
							mbsc.Field(k).SetFloat(mbsc.Field(k-1).Float() + mbsi2.Field(k).Float())
						}
					}
				}
				bootStrappedCumulatives = append(bootStrappedCumulatives, bootStrappedCumulative)
			}

			if sim == run.Simulations {
				//DB.Where("run_date=? and product_code=? and run_id=?", runDate, bootStrappedCumulatives[0].ProductCode, bootStrappedCumulatives[0].RunID).Delete(&models.LicBootStrappedCumulative{})
				err = DB.Save(&bootStrappedCumulatives).Error
				if err != nil {
					log.Println(err)
				}
			}

			var bootStrappedResult models.LicBootstrappedResults
			bootStrappedResult = BootstrappedDevelopmentFactor(bootStrappedCumulatives, runDate, endDateYear, endDateMonth, projMonths, params, run, sim, interval)

			bootStrappedResults = append(bootStrappedResults, bootStrappedResult)

			if sim == run.Simulations {
				DB.Where("run_date=? and product_code=? and run_id=?", runDate, bootStrappedResults[0].ProductCode, bootStrappedResults[0].RunID).Delete(&models.LicBootstrappedResults{})
				err = DB.Save(&bootStrappedResults).Error
				if err != nil {
					log.Println(err)
				}
			}

		}

		// Summary Stats
		var err error

		query := fmt.Sprintf("SELECT run_date,run_id,portfolio_name,lic_portfolio_id,product_code,avg(reserve) as mean, stddev(reserve) as standard_deviation, min(reserve) as minimum, max(reserve) as maximum FROM lic_bootstrapped_results where run_date = '%s' and product_code = '%s' and portfolio_name = '%s' and lic_portfolio_id = %d and run_id = %d", run.RunDate, productCode, run.PortfolioName, run.PortfolioId, run.ID)
		var results []models.LicBootstrappedResultSummary
		err = DB.Raw(query).Scan(&results).Error
		if err != nil {
			fmt.Println(err)
		}

		Query := fmt.Sprintf("select reserve FROM lic_bootstrapped_results where run_date = '%s' and product_code = '%s' and portfolio_name = '%s' and lic_portfolio_id = %d and run_id = %d order by reserve", run.RunDate, productCode, run.PortfolioName, run.PortfolioId, run.ID)
		var medianData []float64
		var nthpercentile float64
		var median, mean float64
		err = DB.Raw(Query).Scan(&medianData).Error
		median, _ = stats.Median(medianData)
		if params.RiskAdjustmentConfidenceLevel == 0 {
			nthpercentile = 0
		} else {
			nthpercentile, _ = stats.Percentile(medianData, params.RiskAdjustmentConfidenceLevel*100.0)
		}
		mean, _ = stats.Mean(medianData)

		results[0].Percentile = nthpercentile
		results[0].Median = median
		bootstrapmeanOut = mean
		NthPercentileOut = nthpercentile

		//DB.Where("run_date=? and product_code=?", runDate, results[0].ProductCode).Delete(&models.LicBootstrappedResultSummary{})
		err = DB.Save(&results).Error
		if err != nil {
			log.Println(err)
		}

		//Frequency Table
		var frequencytable = []models.IbnrFrequency{}
		groupStep := (results[0].Maximum - results[0].Minimum) / 15
		for freqGroup := 0; freqGroup <= 15; freqGroup++ {
			var frequency models.IbnrFrequency
			frequency.RunDate = results[0].RunDate
			frequency.RunID = results[0].RunID
			frequency.PortfolioName = results[0].PortfolioName
			frequency.LicPortfolioId = results[0].LicPortfolioId
			frequency.ProductCode = results[0].ProductCode
			frequency.Group = freqGroup + 1
			if groupStep > 0 {
				if freqGroup == 0 {
					frequency.Reserve = math.Floor(results[0].Minimum)
				} else {
					if freqGroup == 15 {
						frequency.Reserve = math.Ceil(results[0].Maximum)
					} else {
						frequency.Reserve = frequencytable[freqGroup-1].Reserve + groupStep
					}

				}
				frequencytable = append(frequencytable, frequency)
			}
			if groupStep == 0 && freqGroup == 0 { //adds one line if gap is zero else no
				frequencytable = append(frequencytable, frequency)
				freqGroup = 16
			}

		}

		if groupStep > 0 {
			for freqGroup := 0; freqGroup < 15; freqGroup++ {
				lowerBound := frequencytable[freqGroup].Reserve
				upperBound := frequencytable[freqGroup+1].Reserve
				query := fmt.Sprintf("SELECT count(reserve) as count FROM lic_bootstrapped_results where reserve >= '%f' and reserve < '%f' and run_date = '%s' and product_code = '%s' and run_id = %d", lowerBound, upperBound, run.RunDate, productCode, results[0].RunID)
				var count int
				err = DB.Raw(query).Scan(&count).Error
				if err != nil {
					fmt.Println(err)
				}
				frequencytable[freqGroup].Frequency = count
			}
		}

		//DB.Where("run_date=? and product_code=?", runDate, frequencytable[0].ProductCode).Delete(&models.IbnrFrequency{})
		err = DB.Save(&frequencytable).Error
		if err != nil {
			log.Println(err)
		}

	}

	//IBNR Reserve by Accident Year and Accident Month

	// Pre-pass: resolve the effective IBNR method for every accident year row.
	// Precedence: AY-range assignment (narrowest wins) > product-level assignment > run.IBNRMethod.
	// When assignments is empty (no per-product overrides configured), every row gets run.IBNRMethod
	// and the behaviour is identical to the original single-method path.
	effectiveMethods := make([]string, len(adjustedCumulativeProjections))
	for i := range adjustedCumulativeProjections {
		effectiveMethods[i] = ResolveAYMethod(assignments, productCode, adjustedCumulativeProjections[i].AccidentYear, run.IBNRMethod)
	}

	var actualclaimstotal, premiumrunofftotal float64
	for i, _ := range adjustedCumulativeProjections {
		em := effectiveMethods[i]
		var ibnrReserve models.LicIbnrReserve
		macp := reflect.ValueOf(&adjustedCumulativeProjections[i]).Elem()
		mpnr3 := reflect.ValueOf(&proportionNotRunoff).Elem()
		mpr3 := reflect.ValueOf(&proportionRunoff).Elem()
		mdiipATC2 := reflect.ValueOf(&discountedinflatedincrementalProjectionsTotalClaim[i]).Elem()
		//mibnr := reflect.ValueOf(&ibnrReserve).Elem()
		var yearadd, movingyear int
		var residualmonth, movingmonth int
		ibnrReserve.ID = adjustedCumulativeProjections[i].ID
		ibnrReserve.RunDate = adjustedCumulativeProjections[i].RunDate
		ibnrReserve.RunID = adjustedCumulativeProjections[i].RunID
		ibnrReserve.PortfolioName = adjustedCumulativeProjections[i].PortfolioName
		ibnrReserve.LicPortfolioId = adjustedCumulativeProjections[i].LicPortfolioId
		ibnrReserve.ProductCode = adjustedCumulativeProjections[i].ProductCode
		ibnrReserve.AccidentYear = adjustedCumulativeProjections[i].AccidentYear
		ibnrReserve.AccidentMonth = adjustedCumulativeProjections[i].AccidentMonth
		ibnrReserve.EffectiveMethod = em
		ibnrReserve.EarnedPremium = adjustedCumulativeProjections[i].EarnedPremium
		params2 := getParams(run.ParameterYear, productCode, run.Basis)
		ibnrReserve.ExpectedTotalLoss = utils.FloatPrecision(ibnrReserve.EarnedPremium*params2.ExpectedLossRatio, AccountingPrecision)
		ibnrReserve.ProportionNotRunoff = mpnr3.Field(projMonths + FactorVariableCount - i).Float()
		ibnrReserve.BfIbnr = ibnrReserve.ExpectedTotalLoss * ibnrReserve.ProportionNotRunoff
		if i > 0 {
			ibnrReserve.BfIbnrAt12 = ibnrReserve.ExpectedTotalLoss * mpnr3.Field(projMonths+FactorVariableCount-i+1).Float()
		}
		ibnrReserve.ProportionRunoff = mpr3.Field(projMonths + FactorVariableCount - i).Float()
		ibnrReserve.PremiumRunoff = ibnrReserve.EarnedPremium * ibnrReserve.ProportionRunoff
		// Cape Cod pool: only accumulate premium run-off for AYs assigned to Cape Cod methods.
		// This ensures the empirical loss ratio is computed from the correct subset of data.
		if em == CapeCod || em == ChainLadderCapeCod {
			premiumrunofftotal += ibnrReserve.PremiumRunoff
		}
		for k := projMonths + TriangleVariableCount; k >= 0; k-- {
			if k > TriangleVariableCount {
				if interval == 12 {
					yearadd = int(int64((int(macp.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1) / (12 / interval)))
					residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(macp.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				}
				if interval != 12 {
					yearadd = int(int64((int(macp.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1 - 1) / (12 / interval)))
					residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(macp.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				}
				movingyear = int(macp.FieldByName(AccidentYear).Int()) + yearadd
				if residualmonth == 0 && interval == 3 {
					residualmonth = 4
				}
				if residualmonth == 0 && interval == 1 {
					residualmonth = 12
				}
				movingmonth = residualmonth
				if movingyear == endDateYear && movingmonth == endDateMonth {
					ibnrReserve.ChainLadderIbnr = macp.Field(projMonths+TriangleVariableCount).Float() - macp.Field(k).Float()
					if k < projMonths+TriangleVariableCount {
						ibnrReserve.ChainLadderIbnrAt12 = macp.Field(projMonths+TriangleVariableCount).Float() - macp.Field(k+1).Float() // to update +1 with projectedtimeAt12
					}
					ibnrReserve.ActualClaims = macp.Field(k).Float()
					// Cape Cod pool: only include actual claims for AYs using Cape Cod methods.
					if em == CapeCod || em == ChainLadderCapeCod {
						actualclaimstotal += ibnrReserve.ActualClaims
					}
				}
				if movingyear > endDateYear || (movingyear == endDateYear && movingmonth > endDateMonth) {
					ibnrReserve.ChainLadderAverageCostPerClaimIbnr += mdiipATC2.Field(k).Float()
					if k < projMonths+TriangleVariableCount {
						ibnrReserve.ChainLadderAverageCostPerClaimIbnrAt12 += mdiipATC2.Field(k + 1).Float() // to update +1 with projectedtimeAt12
					}
				}
			}

			if em == ChainLadder {
				ibnrReserve.RiskAdjustment = ibnrReserve.ChainLadderIbnr * params2.RiskAdjustmentFactor
			}
			if em == BornHuetterFerguson {
				ibnrReserve.RiskAdjustment = ibnrReserve.BfIbnr * params2.RiskAdjustmentFactor
			}

			if em == ChainLadderAverageCostperClaim {
				ibnrReserve.RiskAdjustment = ibnrReserve.ChainLadderAverageCostPerClaimIbnr * params2.RiskAdjustmentFactor
			}

		}

		// Combined CL+BF: credibility-weighted by proportion run-off.
		// Mature accident years (high run-off) rely on observed Chain Ladder data;
		// immature years (low run-off) rely on the BF a priori expectation.
		ibnrReserve.CombinedClBfIbnr = ibnrReserve.ProportionRunoff*ibnrReserve.ChainLadderIbnr +
			ibnrReserve.ProportionNotRunoff*ibnrReserve.BfIbnr
		if i > 0 {
			pnrAt12 := mpnr3.Field(projMonths + FactorVariableCount - i + 1).Float()
			ibnrReserve.CombinedClBfIbnrAt12 = (1.0-pnrAt12)*ibnrReserve.ChainLadderIbnrAt12 +
				pnrAt12*ibnrReserve.BfIbnrAt12
		}
		if em == ChainLadderBF {
			ibnrReserve.RiskAdjustment = ibnrReserve.CombinedClBfIbnr * params2.RiskAdjustmentFactor
		}

		*ibnrReserves = append(*ibnrReserves, ibnrReserve)

		// ibnr reserve report — use em (per-AY effective method) so mixed-method runs
		// accumulate the correct IBNR into the portfolio BEL total.
		if em == ChainLadder {
			(*IbnrReserveReport).IbnrBel += ibnrReserve.ChainLadderIbnr
			(*IbnrReserveReport).IbnrBelAt12 += ibnrReserve.ChainLadderIbnrAt12
			(*IbnrReserveReport).IbnrRiskAdjustmentAt12 += ibnrReserve.ChainLadderIbnrAt12 * params2.RiskAdjustmentFactor

		}
		if em == BornHuetterFerguson {
			(*IbnrReserveReport).IbnrBel += ibnrReserve.BfIbnr
			(*IbnrReserveReport).IbnrBelAt12 += ibnrReserve.BfIbnrAt12
			(*IbnrReserveReport).IbnrRiskAdjustmentAt12 += ibnrReserve.BfIbnrAt12 * params2.RiskAdjustmentFactor
		}

		if em == ChainLadderAverageCostperClaim {
			(*IbnrReserveReport).IbnrBel += ibnrReserve.ChainLadderAverageCostPerClaimIbnr
			(*IbnrReserveReport).IbnrBelAt12 += ibnrReserve.ChainLadderAverageCostPerClaimIbnrAt12
			(*IbnrReserveReport).IbnrRiskAdjustmentAt12 += ibnrReserve.ChainLadderAverageCostPerClaimIbnrAt12 * params2.RiskAdjustmentFactor
		}
		if em == ChainLadderBF {
			(*IbnrReserveReport).IbnrBel += ibnrReserve.CombinedClBfIbnr
			(*IbnrReserveReport).IbnrBelAt12 += ibnrReserve.CombinedClBfIbnrAt12
			(*IbnrReserveReport).IbnrRiskAdjustmentAt12 += ibnrReserve.CombinedClBfIbnrAt12 * params2.RiskAdjustmentFactor
		}
		(*IbnrReserveReport).ChainLadderIbnr += ibnrReserve.ChainLadderIbnr
		(*IbnrReserveReport).ChainLadderAverageCostPerClaimIbnr += ibnrReserve.ChainLadderAverageCostPerClaimIbnr
		(*IbnrReserveReport).BfIbnr += ibnrReserve.BfIbnr
		(*IbnrReserveReport).IBNRRiskAdjustment += ibnrReserve.RiskAdjustment
		// Always accumulate combined totals so the comparison view shows all methods side-by-side
		(*IbnrReserveReport).CombinedClBfIbnr += ibnrReserve.CombinedClBfIbnr

	}

	for i, _ := range *ibnrReserves {
		em2 := effectiveMethods[i]
		params3 := getParams(run.ParameterYear, productCode, run.Basis)
		if premiumrunofftotal > 0 {
			// Cape Cod IBNR fields are always computed so the comparison view can show
			// all methods side-by-side. Risk adjustment and BEL accumulation use em2.
			(*ibnrReserves)[i].PredictedLossRatio = actualclaimstotal / premiumrunofftotal
			(*ibnrReserves)[i].PredictedTotalLoss = (*ibnrReserves)[i].EarnedPremium * (*ibnrReserves)[i].PredictedLossRatio
			(*ibnrReserves)[i].CapeCodIbnr = (*ibnrReserves)[i].PredictedTotalLoss * (*ibnrReserves)[i].ProportionNotRunoff
			if i > 0 {
				(*ibnrReserves)[i].CapeCodIbnrAt12 = (*ibnrReserves)[i].PredictedTotalLoss * (*ibnrReserves)[i-1].ProportionNotRunoff
			}

			// Combined CL+Cape Cod: credibility-weighted by proportion run-off.
			// Cape Cod uses an empirical portfolio loss ratio rather than an a priori assumption.
			(*ibnrReserves)[i].CombinedClCapeCodIbnr = (*ibnrReserves)[i].ProportionRunoff*(*ibnrReserves)[i].ChainLadderIbnr +
				(*ibnrReserves)[i].ProportionNotRunoff*(*ibnrReserves)[i].CapeCodIbnr
			if i > 0 {
				pnrAt12 := (*ibnrReserves)[i-1].ProportionNotRunoff
				(*ibnrReserves)[i].CombinedClCapeCodIbnrAt12 = (1.0-pnrAt12)*(*ibnrReserves)[i].ChainLadderIbnrAt12 +
					pnrAt12*(*ibnrReserves)[i].CapeCodIbnrAt12
			}

			if em2 == CapeCod {
				(*ibnrReserves)[i].RiskAdjustment = (*ibnrReserves)[i].CapeCodIbnr * params3.RiskAdjustmentFactor
			}
			if em2 == ChainLadderCapeCod {
				(*ibnrReserves)[i].RiskAdjustment = (*ibnrReserves)[i].CombinedClCapeCodIbnr * params3.RiskAdjustmentFactor
			}
		}
		if em2 == CapeCod {
			(*IbnrReserveReport).IbnrBel += (*ibnrReserves)[i].CapeCodIbnr
			(*IbnrReserveReport).IBNRRiskAdjustment += (*ibnrReserves)[i].RiskAdjustment
			(*IbnrReserveReport).IbnrBelAt12 += (*ibnrReserves)[i].CapeCodIbnrAt12
			(*IbnrReserveReport).IbnrRiskAdjustmentAt12 = (*ibnrReserves)[i].CapeCodIbnrAt12 * params3.RiskAdjustmentFactor

		}
		if em2 == ChainLadderCapeCod {
			(*IbnrReserveReport).IbnrBel += (*ibnrReserves)[i].CombinedClCapeCodIbnr
			(*IbnrReserveReport).IBNRRiskAdjustment += (*ibnrReserves)[i].CombinedClCapeCodIbnr * params3.RiskAdjustmentFactor
			(*IbnrReserveReport).IbnrBelAt12 += (*ibnrReserves)[i].CombinedClCapeCodIbnrAt12
			(*IbnrReserveReport).IbnrRiskAdjustmentAt12 += (*ibnrReserves)[i].CombinedClCapeCodIbnrAt12 * params3.RiskAdjustmentFactor
		}
		(*IbnrReserveReport).CapeCodIbnr += (*ibnrReserves)[i].CapeCodIbnr
		// Always accumulate combined totals so the comparison view shows all methods side-by-side
		(*IbnrReserveReport).CombinedClCapeCodIbnr += (*ibnrReserves)[i].CombinedClCapeCodIbnr
	}

	//DB.Where("run_date=? and product_code=?", runDate, ibnrReserves[0].ProductCode).Delete(&models.LicIbnrReserve{})
	//err = DB.Save(&ibnrReserves).Error
	//if err != nil {
	//	log.Println(err)
	//}

	//ibnr reserve report
	(*IbnrReserveReport).RunID = run.ID
	(*IbnrReserveReport).RunDate = (*ibnrReserves)[0].RunDate
	(*IbnrReserveReport).RunID = (*ibnrReserves)[0].RunID
	(*IbnrReserveReport).LicPortfolioId = (*ibnrReserves)[0].LicPortfolioId
	(*IbnrReserveReport).PortfolioName = (*ibnrReserves)[0].PortfolioName
	(*IbnrReserveReport).ProductCode = (*ibnrReserves)[0].ProductCode
	//IbnrReserveReport.IBNRRiskAdjustment = IbnrReserveReport.IbnrBel * params.RiskAdjustmentFactor
	(*IbnrReserveReport).IbnrReserve = IbnrReserveReport.IbnrBel + IbnrReserveReport.IBNRRiskAdjustment
	(*IbnrReserveReport).BootstrapIbnr = bootstrapmeanOut
	(*IbnrReserveReport).BootstrapNthPercentileIbnr = NthPercentileOut

}

func UpdateMackModel(individualDevelopmentFactors *[]models.LicIndividualDevelopmentFactors, biasAdjustmentFactors *[]models.LicBiasAdjustmentFactor, cumulativeTriangulations []models.LicCumulativeTriangulation, cumulativeProjections []models.LicCumulativeProjection, developmentFactor models.LicDevelopmentFactor, ibnrReserves *[]models.LicIbnrReserve, IbnrReserveReport *models.IbnrReserveReport, runDate string, endDateYear, endDateMonth, projMonths int, params models.LICParameter, run models.IBNRRunSetting, productCode string, interval int, discountOption string, manualDfs []models.LicDevelopmentFactor, manualProdCode string) {

	var individualDevelopmentFactorsResults []models.LicIndividualDevelopmentFactors

	for i, _ := range cumulativeTriangulations {
		mutable := reflect.ValueOf(&cumulativeTriangulations[i]).Elem()
		var individualDevelopmentFactor models.LicIndividualDevelopmentFactors
		var biasAdjustmentFactor models.LicBiasAdjustmentFactor
		mack := reflect.ValueOf(&individualDevelopmentFactor).Elem()
		df := reflect.ValueOf(&developmentFactor).Elem()
		var manualdf reflect.Value

		if len(manualDfs) > 0 {
			manualdf = reflect.ValueOf(&manualDfs[0]).Elem()
		}

		biasAdjFac := reflect.ValueOf(&biasAdjustmentFactor).Elem()

		//Computing Individual Development Factors
		var yearadd, movingyear int
		var residualmonth, movingmonth int
		for k := 0; k <= projMonths+TriangleVariableCount; k++ {
			if k <= TriangleVariableCount { //assigning variable names
				mack.Field(k).Set(mutable.Field(k))
				//biasAdjFac.Field(k).Set(mutable.Field(k))
			} else {
				//fmt.Println(mutable.Type().Field(k).Name)
				if interval == 12 {
					yearadd = int(int64((int(mutable.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1) / (12 / interval)))
					residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mutable.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				}
				if interval != 12 {
					yearadd = int(int64((int(mutable.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1 - 1) / (12 / interval)))
					residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mutable.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				}
				residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mutable.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				movingyear = int(mutable.FieldByName(AccidentYear).Int()) + yearadd
				if residualmonth == 0 && interval == 3 {
					residualmonth = 1
				}
				if residualmonth == 0 && interval == 1 {
					residualmonth = 12
				}
				movingmonth = residualmonth
				if k == TriangleVariableCount+1 {
					mack.Field(k).SetFloat(0)
					biasAdjFac.Field(k).SetFloat(math.Pow(mutable.Field(k).Float(), 2.0)) ///????
				} else {
					if (movingyear <= endDateYear && movingmonth <= endDateMonth) || (movingyear < endDateYear && movingmonth > endDateMonth) {
						if mutable.Field(k-1).Float() > 0 {
							mack.Field(k).SetFloat(mutable.Field(k).Float() / mutable.Field(k-1).Float())
						}
					} else {
						mack.Field(k).SetFloat(0)
					}
				}
			}
		}

		*individualDevelopmentFactors = append(*individualDevelopmentFactors, individualDevelopmentFactor)
		individualDevelopmentFactorsResults = *individualDevelopmentFactors
		DB.Save(*individualDevelopmentFactors)

		//Computing Bias Adjustment Factor
		for k := 0; k <= projMonths+TriangleVariableCount; k++ {
			if k <= TriangleVariableCount { //assigning variable names
				biasAdjFac.Field(k).Set(mutable.Field(k))
			} else {
				//fmt.Println(mutable.Type().Field(k).Name)
				if interval == 12 {
					yearadd = int(int64((int(mutable.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1) / (12 / interval)))
					residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mutable.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				}
				if interval != 12 {
					yearadd = int(int64((int(mutable.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1 - 1) / (12 / interval)))
					residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mutable.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				}
				residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mutable.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				movingyear = int(mutable.FieldByName(AccidentYear).Int()) + yearadd
				if residualmonth == 0 && interval == 3 {
					residualmonth = 1
				}
				if residualmonth == 0 && interval == 1 {
					residualmonth = 12
				}
				movingmonth = residualmonth
				if k == TriangleVariableCount+1 {
					biasAdjFac.Field(k).SetFloat(0) //biasAdjFac.Field(k).SetFloat(math.Pow(mack.Field(k+1).Float()-df.Field(k+1-VariableCountDiff).Float(), 2.0) * mutable.Field(k).Float())
				} else {
					if (movingyear <= endDateYear && movingmonth <= endDateMonth) || (movingyear < endDateYear && movingmonth > endDateMonth) {
						if run.RerunIndicator && productCode == manualProdCode && manualdf.IsValid() {
							biasAdjFac.Field(k).SetFloat(math.Pow(mack.Field(k).Float()-manualdf.Field(k-VariableCountDiff).Float(), 2.0) * mutable.Field(k-1).Float())
						} else {
							biasAdjFac.Field(k).SetFloat(math.Pow(mack.Field(k).Float()-df.Field(k-VariableCountDiff).Float(), 2.0) * mutable.Field(k-1).Float())
						}
					} else {
						biasAdjFac.Field(k).SetFloat(0)
					}
				}
			}
		}

		//fmt.Println(cumulativeTriangulation)
		*biasAdjustmentFactors = append(*biasAdjustmentFactors, biasAdjustmentFactor)
	}
	DB.Save(*biasAdjustmentFactors)

	//Adding Mack Model Parameters
	biasAdjustmentFactorsResult := *biasAdjustmentFactors

	var volumeWeightedDevFac models.LicMackModelCalculatedParameters
	var biasAdjcolumnSum models.LicMackModelCalculatedParameters
	var includedFactors models.LicMackModelCalculatedParameters
	var totalNumberOfIncludedFactors models.LicMackModelCalculatedParameters
	var numberOfDevelopmentFactors models.LicMackModelCalculatedParameters
	var numberOfParametersInDFM models.LicMackModelCalculatedParameters
	var sigmaSquared models.LicMackModelCalculatedParameters
	var sigma models.LicMackModelCalculatedParameters

	volumeWeightedDevFac.RunDate = biasAdjustmentFactorsResult[0].RunDate
	volumeWeightedDevFac.RunID = biasAdjustmentFactorsResult[0].RunID
	volumeWeightedDevFac.PortfolioName = biasAdjustmentFactorsResult[0].PortfolioName
	volumeWeightedDevFac.ProductCode = biasAdjustmentFactorsResult[0].ProductCode
	volumeWeightedDevFac.MackModelVariable = "Volume Weighted Development Factor"

	biasAdjcolumnSum.RunDate = biasAdjustmentFactorsResult[0].RunDate
	biasAdjcolumnSum.RunID = biasAdjustmentFactorsResult[0].RunID
	biasAdjcolumnSum.PortfolioName = biasAdjustmentFactorsResult[0].PortfolioName
	biasAdjcolumnSum.ProductCode = biasAdjustmentFactorsResult[0].ProductCode
	biasAdjcolumnSum.MackModelVariable = "Bias Adjustment Column Sum"

	includedFactors.RunDate = biasAdjustmentFactorsResult[0].RunDate
	includedFactors.RunID = biasAdjustmentFactorsResult[0].RunID
	includedFactors.PortfolioName = biasAdjustmentFactorsResult[0].PortfolioName
	includedFactors.ProductCode = biasAdjustmentFactorsResult[0].ProductCode
	includedFactors.MackModelVariable = "Number of Included Factors"

	totalNumberOfIncludedFactors.RunDate = biasAdjustmentFactorsResult[0].RunDate
	totalNumberOfIncludedFactors.RunID = biasAdjustmentFactorsResult[0].RunID
	totalNumberOfIncludedFactors.PortfolioName = biasAdjustmentFactorsResult[0].PortfolioName
	totalNumberOfIncludedFactors.ProductCode = biasAdjustmentFactorsResult[0].ProductCode
	totalNumberOfIncludedFactors.MackModelVariable = "Total Number of Included Factors"

	numberOfDevelopmentFactors.RunDate = biasAdjustmentFactorsResult[0].RunDate
	numberOfDevelopmentFactors.RunID = biasAdjustmentFactorsResult[0].RunID
	numberOfDevelopmentFactors.PortfolioName = biasAdjustmentFactorsResult[0].PortfolioName
	numberOfDevelopmentFactors.ProductCode = biasAdjustmentFactorsResult[0].ProductCode
	numberOfDevelopmentFactors.MackModelVariable = "Included Number of Development Factors"

	numberOfParametersInDFM.RunDate = biasAdjustmentFactorsResult[0].RunDate
	numberOfParametersInDFM.RunID = biasAdjustmentFactorsResult[0].RunID
	numberOfParametersInDFM.PortfolioName = biasAdjustmentFactorsResult[0].PortfolioName
	numberOfParametersInDFM.ProductCode = biasAdjustmentFactorsResult[0].ProductCode
	numberOfParametersInDFM.MackModelVariable = "Number of DFM Parameters"

	sigmaSquared.RunDate = biasAdjustmentFactorsResult[0].RunDate
	sigmaSquared.RunID = biasAdjustmentFactorsResult[0].RunID
	sigmaSquared.PortfolioName = biasAdjustmentFactorsResult[0].PortfolioName
	sigmaSquared.ProductCode = biasAdjustmentFactorsResult[0].ProductCode
	sigmaSquared.MackModelVariable = "Sigma Squared"

	sigma.RunDate = biasAdjustmentFactorsResult[0].RunDate
	sigma.RunID = biasAdjustmentFactorsResult[0].RunID
	sigma.PortfolioName = biasAdjustmentFactorsResult[0].PortfolioName
	sigma.ProductCode = biasAdjustmentFactorsResult[0].ProductCode
	sigma.MackModelVariable = "Sigma"

	vwDevFac := reflect.ValueOf(&volumeWeightedDevFac).Elem() //volume weighted Development Factor
	biasAdjcolSum := reflect.ValueOf(&biasAdjcolumnSum).Elem()
	includedFac := reflect.ValueOf(&includedFactors).Elem()
	totalFac := reflect.ValueOf(&totalNumberOfIncludedFactors).Elem()
	devFac := reflect.ValueOf(&numberOfDevelopmentFactors).Elem()
	DFMParams := reflect.ValueOf(&numberOfParametersInDFM).Elem()
	sigSq := reflect.ValueOf(&sigmaSquared).Elem()
	sig := reflect.ValueOf(&sigma).Elem()
	var totalFacCount float64
	for i, _ := range biasAdjustmentFactorsResult {
		biasAdjustFac := reflect.ValueOf(&biasAdjustmentFactorsResult[i]).Elem()
		indivDevFac := reflect.ValueOf(&individualDevelopmentFactorsResults[i]).Elem()
		finalDevFacResult := reflect.ValueOf(&developmentFactor).Elem()
		var manualfinalDevFacResult reflect.Value

		if len(manualDfs) > 0 {
			manualfinalDevFacResult = reflect.ValueOf(&manualDfs[0]).Elem()
		}

		for k := TriangleVariableCount + 1; k <= projMonths+TriangleVariableCount; k++ {
			if k == TriangleVariableCount+1 {
				if i == 0 {
					if run.RerunIndicator && productCode == manualProdCode && manualfinalDevFacResult.IsValid() {
						vwDevFac.Field(k - VariableCountDiff).SetFloat(manualfinalDevFacResult.Field(k - VariableCountDiff).Float())
					} else {
						vwDevFac.Field(k - VariableCountDiff).SetFloat(finalDevFacResult.Field(k - VariableCountDiff).Float())
					}

				}
				biasAdjcolSum.Field(k - VariableCountDiff).SetFloat(biasAdjcolSum.Field(k-VariableCountDiff).Float() + biasAdjustFac.Field(k).Float())
				includedFac.Field(k - VariableCountDiff).SetFloat(0)
				totalFac.Field(k - VariableCountDiff).SetFloat(0)
				devFac.Field(k - VariableCountDiff).SetFloat(0)
				DFMParams.Field(k - VariableCountDiff).SetFloat(0)

			} else {
				if i == 0 { //initialise and read volume weighed development factor from upstream
					if run.RerunIndicator && productCode == manualProdCode && manualfinalDevFacResult.IsValid() {
						vwDevFac.Field(k - VariableCountDiff).SetFloat(manualfinalDevFacResult.Field(k - VariableCountDiff).Float())
					} else {
						vwDevFac.Field(k - VariableCountDiff).SetFloat(finalDevFacResult.Field(k - VariableCountDiff).Float())
					}
				}
				biasAdjcolSum.Field(k - VariableCountDiff).SetFloat(biasAdjcolSum.Field(k-VariableCountDiff).Float() + biasAdjustFac.Field(k).Float())

				if indivDevFac.Field(k).Float() != 0 {
					includedFac.Field(k - VariableCountDiff).SetFloat(includedFac.Field(k-VariableCountDiff).Float() + 1)
					totalFacCount = totalFacCount + 1
				}
				totalFac.Field(k - VariableCountDiff).SetFloat(totalFacCount)
				devFac.Field(k - VariableCountDiff).SetFloat(float64(projMonths - 1))
				DFMParams.Field(k - VariableCountDiff).SetFloat(float64(projMonths))

				//calculating sigma squared and sigma once parameters are calculated
				if i == len(biasAdjustmentFactorsResult)-1 {
					n := totalFacCount
					nd := includedFac.Field(k - VariableCountDiff).Float()
					Z := devFac.Field(FactorVariableCount + 2).Float()
					p := DFMParams.Field(FactorVariableCount + 2).Float()
					if includedFac.Field(k-VariableCountDiff).Float() == 1 || includedFac.Field(k-VariableCountDiff).Float() == 0 {
						sigSq.Field(k - VariableCountDiff).SetFloat(0)
					} else {
						sigSq.Field(k - VariableCountDiff).SetFloat(biasAdjcolSum.Field(k-VariableCountDiff).Float() * (1 / (nd - 1)) * ((n - Z) / (n - p)))
						fmt.Println(biasAdjcolSum.Field(k - VariableCountDiff).Float())
					}
					sig.Field(k - VariableCountDiff).SetFloat(math.Sqrt(sigSq.Field(k - VariableCountDiff).Float()))
				}

			}

		}
	}

	DB.Save(&vwDevFac)
	DB.Save(&includedFactors)
	DB.Save(&totalNumberOfIncludedFactors)
	DB.Save(&numberOfDevelopmentFactors)
	DB.Save(&numberOfParametersInDFM)
	DB.Save(&biasAdjcolumnSum)
	DB.Save(&sigmaSquared)
	DB.Save(&sigma)

	// End of Mack Model Block

	//Bias Adjusted Residuals
	var biasAdjustedResisuals []models.LicBiasAdjustedResiduals
	var nbrNonZeroResiduals, sumResiduals, meanBiasAdjustedResiduals float64
	for i, _ := range cumulativeTriangulations {

		var biasAdjustedResidual models.LicBiasAdjustedResiduals
		cumTria := reflect.ValueOf(&cumulativeTriangulations[i]).Elem()
		indDevFac := reflect.ValueOf(&individualDevelopmentFactorsResults[i]).Elem()
		biasAdjRes := reflect.ValueOf(&biasAdjustedResidual).Elem()
		var yearadd, movingyear int
		var residualmonth, movingmonth int
		for k := 0; k <= projMonths+TriangleVariableCount; k++ {
			var temp2 float64
			if k <= TriangleVariableCount {
				biasAdjRes.Field(k).Set(cumTria.Field(k))
			} else {

				if interval == 12 {
					yearadd = int(int64((int(cumTria.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1) / (12 / interval)))
					residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(cumTria.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				}
				if interval != 12 {
					yearadd = int(int64((int(cumTria.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1 - 1) / (12 / interval)))
					residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(cumTria.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				}
				residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(cumTria.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				movingyear = int(cumTria.FieldByName(AccidentYear).Int()) + yearadd
				if residualmonth == 0 && interval == 3 {
					residualmonth = 1
				}
				if residualmonth == 0 && interval == 1 {
					residualmonth = 12
				}
				movingmonth = residualmonth
				if (movingyear <= endDateYear && movingmonth <= endDateMonth) || (movingyear < endDateYear && movingmonth > endDateMonth) {
					n := totalFac.Field(k - VariableCountDiff).Float()
					nd := includedFac.Field(k - VariableCountDiff).Float()
					Z := devFac.Field(k - VariableCountDiff).Float()
					p := DFMParams.Field(k - VariableCountDiff).Float()
					//temp1 := indDevFac.Field(k+1).Float() - vwDevFac.Field(k-VariableCountDiff+1).Float()
					if nd > 1 && (n-p) != 0 {
						temp2 = math.Sqrt((nd / (nd - 1)) * (n - Z) / (n - p))
					}
					if sig.Field(k-VariableCountDiff).Float() != 0 {
						biasAdjRes.Field(k).SetFloat((math.Sqrt(cumTria.Field(k-1).Float()) / sig.Field(k-VariableCountDiff).Float() * (indDevFac.Field(k).Float() - vwDevFac.Field(k-VariableCountDiff).Float())) * temp2)
						if biasAdjRes.Field(k).Float() != 0 {
							nbrNonZeroResiduals += 1
							sumResiduals += biasAdjRes.Field(k).Float()
						}
					}

				} else {
					biasAdjRes.Field(k).SetFloat(0)
				}
			}
		}

		biasAdjustedResisuals = append(biasAdjustedResisuals, biasAdjustedResidual)

	} // End of Bias Adjusted Residuals
	if nbrNonZeroResiduals > 0 {
		meanBiasAdjustedResiduals = sumResiduals / nbrNonZeroResiduals
	}
	DB.Save(&biasAdjustedResisuals)

	//Mean Adjusted Residuals
	var meanAdjustedResiduals []models.LicMeanBiasAdjustedResiduals
	for i, _ := range biasAdjustedResisuals {
		var meanAdjustedResidual models.LicMeanBiasAdjustedResiduals
		biasAdjRes := reflect.ValueOf(&biasAdjustedResisuals[i]).Elem()
		meanAdjRes := reflect.ValueOf(&meanAdjustedResidual).Elem()

		var yearadd, movingyear int
		var residualmonth, movingmonth int
		for k := 0; k <= projMonths+TriangleVariableCount; k++ {
			if k <= TriangleVariableCount {
				meanAdjRes.Field(k).Set(biasAdjRes.Field(k))
			} else {
				if interval == 12 {
					yearadd = int(int64((int(biasAdjRes.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1) / (12 / interval)))
					residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(biasAdjRes.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				}
				if interval != 12 {
					yearadd = int(int64((int(biasAdjRes.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1 - 1) / (12 / interval)))
					residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(biasAdjRes.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				}
				residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(biasAdjRes.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				movingyear = int(biasAdjRes.FieldByName(AccidentYear).Int()) + yearadd
				if residualmonth == 0 && interval == 3 {
					residualmonth = 1
				}
				if residualmonth == 0 && interval == 1 {
					residualmonth = 12
				}
				movingmonth = residualmonth
				if (movingyear <= endDateYear && movingmonth <= endDateMonth) || (movingyear < endDateYear && movingmonth > endDateMonth) {
					if k == TriangleVariableCount+1 {
						meanAdjRes.Field(k).SetFloat(0)
					} else {
						meanAdjRes.Field(k).SetFloat(biasAdjRes.Field(k).Float() - meanBiasAdjustedResiduals)
					}
				} else {
					meanAdjRes.Field(k).SetFloat(0)
				}
			}
		}
		meanAdjustedResiduals = append(meanAdjustedResiduals, meanAdjustedResidual)
	}
	DB.Save(&meanAdjustedResiduals)

	//if run.DistributionModel == Lognormal {
	//sigmas
	var logNormalSigmas []models.LicLogNormalSigmas
	for i, _ := range cumulativeTriangulations {
		var logNormalSigma models.LicLogNormalSigmas
		cumTria200 := reflect.ValueOf(&cumulativeTriangulations[i]).Elem()
		logNormalSig := reflect.ValueOf(&logNormalSigma).Elem()

		//var yearadd, movingyear int
		//var residualmonth, movingmonth int
		for k := 0; k <= projMonths+TriangleVariableCount; k++ {
			if k <= TriangleVariableCount {
				logNormalSig.Field(k).Set(cumTria200.Field(k))
			} else {
				if k == TriangleVariableCount+1 {
					logNormalSig.Field(k).SetFloat(0)
				} else {
					if cumTria200.Field(k-1).Float() > 0 {
						logNormalSig.Field(k).SetFloat(math.Sqrt(sigSq.Field(k-VariableCountDiff).Float() / cumTria200.Field(k-1).Float()))
					} else {
						logNormalSig.Field(k).SetFloat(0)
					}

				}
			}
		}
		logNormalSigmas = append(logNormalSigmas, logNormalSigma)
	}

	if run.DistributionModel == Lognormal || run.DistributionModel == Gamma {
		DB.Save(&logNormalSigmas)
	}

	//mean
	var logNormalMeans []models.LicLogNormalMeans
	for i, _ := range cumulativeTriangulations {
		var logNormalMean models.LicLogNormalMeans
		cumTria300 := reflect.ValueOf(&cumulativeTriangulations[i]).Elem()
		logNormalSig2 := reflect.ValueOf(&logNormalSigmas[i]).Elem()
		logNormalMU := reflect.ValueOf(&logNormalMean).Elem()
		vwDevFac200 := reflect.ValueOf(&volumeWeightedDevFac).Elem()

		for k := 0; k <= projMonths+TriangleVariableCount; k++ {
			if k <= TriangleVariableCount {
				logNormalMU.Field(k).Set(cumTria300.Field(k))
			} else {
				if k == TriangleVariableCount+1 {
					logNormalMU.Field(k).SetFloat(0)
				} else {
					if vwDevFac200.Field(k-VariableCountDiff).Float() > 0 && logNormalSig2.Field(k).Float() > 0 {
						logNormalMU.Field(k).SetFloat(math.Log(vwDevFac200.Field(k-VariableCountDiff).Float()) - math.Pow(logNormalSig2.Field(k).Float(), 2.0)/2.0)

					} else {
						logNormalMU.Field(k).SetFloat(0)
					}
				}
			}
		}
		logNormalMeans = append(logNormalMeans, logNormalMean)
	}

	if run.DistributionModel == Lognormal || run.DistributionModel == Gamma {
		DB.Save(&logNormalMeans)
	}

	// LogNormal Standard Deviation
	var logNormalStandardDeviations []models.LicLogNormalStandardDeviations
	for i, _ := range cumulativeTriangulations {
		var logNormalStandardDeviation models.LicLogNormalStandardDeviations
		cumTria400 := reflect.ValueOf(&cumulativeTriangulations[i]).Elem()
		logNormalMU2 := reflect.ValueOf(&logNormalMeans[i]).Elem()
		logNormalStd := reflect.ValueOf(&logNormalStandardDeviation).Elem()
		vwDevFac300 := reflect.ValueOf(&volumeWeightedDevFac).Elem()

		for k := 0; k <= projMonths+TriangleVariableCount; k++ {
			if k <= TriangleVariableCount {
				logNormalStd.Field(k).Set(cumTria400.Field(k))
			} else {
				if k == TriangleVariableCount+1 {
					logNormalStd.Field(k).SetFloat(0)
				} else {
					if logNormalMU2.Field(k).Float() > 0 {
						temp1 := math.Exp(math.Pow(logNormalMU2.Field(k).Float(), 2.0)) - 1.0
						temp2 := math.Exp(2.0*vwDevFac300.Field(k-VariableCountDiff).Float() + math.Pow(logNormalMU2.Field(k).Float(), 2.0))
						logNormalStd.Field(k).SetFloat(math.Sqrt(temp1 * temp2))
					} else {
						logNormalStd.Field(k).SetFloat(0)
					}
				}
			}
		}
		logNormalStandardDeviations = append(logNormalStandardDeviations, logNormalStandardDeviation)
	}

	if run.DistributionModel == Lognormal || run.DistributionModel == Gamma {
		DB.Save(&logNormalStandardDeviations)
	}

	//Pseudo Ratios
	var MackSimulationResults []models.LicMackSimulationResults
	zScoreTable := ztable.NewZTable(nil)
	rand.Seed(time.Now().UnixNano()) // Is this global....?
	for mackSim := 1; mackSim <= run.MackModelSimulations; mackSim++ {
		var pseudoRatios []models.LicPseudoRatios
		rand.Seed(time.Now().UnixNano())
		for i, _ := range cumulativeTriangulations {
			var pseudoRatio models.LicPseudoRatios
			cumTria2 := reflect.ValueOf(&cumulativeTriangulations[i]).Elem()
			logNormalMU3 := reflect.ValueOf(&logNormalMeans[i]).Elem()
			logNormalStd2 := reflect.ValueOf(&logNormalStandardDeviations[i]).Elem()
			pseudoR := reflect.ValueOf(&pseudoRatio).Elem()
			vwDevFac2 := reflect.ValueOf(&volumeWeightedDevFac).Elem()
			var yearadd, movingyear int
			var residualmonth, movingmonth int
			for k := 0; k <= projMonths+TriangleVariableCount; k++ {
				if k <= TriangleVariableCount {
					pseudoR.Field(k).Set(cumTria2.Field(k))
				} else if k == TriangleVariableCount+1 {
					pseudoR.Field(k).SetFloat(1)
				} else {
					if interval == 12 {
						yearadd = int(int64((int(cumTria2.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1 - 1) / (12 / interval)))
						residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1-1+int(cumTria2.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
					}
					if interval != 12 {
						yearadd = int(int64((int(cumTria2.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1 - 1 - 1) / (12 / interval)))
						residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1-1+int(cumTria2.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
					}
					residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1-1+int(cumTria2.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
					movingyear = int(cumTria2.FieldByName(AccidentYear).Int()) + yearadd
					if residualmonth == 0 && interval == 3 {
						residualmonth = 1
					}

					if residualmonth == 0 && interval == 1 {
						residualmonth = 12
					}
					movingmonth = residualmonth
					if (movingyear <= endDateYear && movingmonth <= endDateMonth) || (movingyear < endDateYear && movingmonth > endDateMonth) {
						if cumTria2.Field(k-1).Float() > 0 {
							if run.DistributionModel == Normal {
								x := utils.NormalInverse(vwDevFac2.Field(k-VariableCountDiff).Float(), math.Sqrt(sigSq.Field(k-VariableCountDiff).Float()/cumTria2.Field(k-1).Float()), zScoreTable) //x = μ + z σ
								pseudoR.Field(k).SetFloat(x)
							} else if run.DistributionModel == Lognormal {
								x := utils.LogNormalInverse(logNormalMU3.Field(k).Float(), logNormalStd2.Field(k).Float(), zScoreTable) //x = μ + z σ
								pseudoR.Field(k).SetFloat(x)
							} else if run.DistributionModel == Gamma {
								x := utils.NormalInverse(vwDevFac2.Field(k-VariableCountDiff).Float(), math.Sqrt(sigSq.Field(k-VariableCountDiff).Float()/cumTria2.Field(k-1).Float()), zScoreTable)
								pseudoR.Field(k).SetFloat(x)

							} else { //Resampling
								x := utils.NormalInverse(vwDevFac2.Field(k-VariableCountDiff).Float(), math.Sqrt(sigSq.Field(k-VariableCountDiff).Float()/cumTria2.Field(k-1).Float()), zScoreTable) //x = μ + z σ
								pseudoR.Field(k).SetFloat(x)
							}

						}
					} else {
						pseudoR.Field(k).SetFloat(1)
					}
				}
			}
			pseudoRatios = append(pseudoRatios, pseudoRatio)
		}

		if mackSim == run.MackModelSimulations {
			DB.Save(&pseudoRatios)
		}

		// LicMackModelSimulatedDevelopmentFactor

		var mackModelSimulatedDevelopmentFactor models.LicMackModelSimulatedDevelopmentFactor
		var cumulativeTriagulationColumnSum models.LicMackModelSimulatedDevelopmentFactor
		mackSimDF := reflect.ValueOf(&mackModelSimulatedDevelopmentFactor).Elem()
		cumTriaColSum := reflect.ValueOf(&cumulativeTriagulationColumnSum).Elem()
		for i, _ := range cumulativeTriangulations {
			cumTria3 := reflect.ValueOf(&cumulativeTriangulations[i]).Elem()
			pseudoR2 := reflect.ValueOf(&pseudoRatios[i]).Elem()
			cumulativeTriagulationColumnSum.RunDate = biasAdjustmentFactorsResult[0].RunDate
			cumulativeTriagulationColumnSum.RunID = biasAdjustmentFactorsResult[0].RunID
			cumulativeTriagulationColumnSum.PortfolioName = biasAdjustmentFactorsResult[0].PortfolioName
			cumulativeTriagulationColumnSum.ProductCode = biasAdjustmentFactorsResult[0].ProductCode
			cumulativeTriagulationColumnSum.MackModelVariable = "Cumulative Triangulation Column Sum"

			mackModelSimulatedDevelopmentFactor.RunDate = biasAdjustmentFactorsResult[0].RunDate
			mackModelSimulatedDevelopmentFactor.RunID = biasAdjustmentFactorsResult[0].RunID
			mackModelSimulatedDevelopmentFactor.PortfolioName = biasAdjustmentFactorsResult[0].PortfolioName
			mackModelSimulatedDevelopmentFactor.ProductCode = biasAdjustmentFactorsResult[0].ProductCode
			mackModelSimulatedDevelopmentFactor.MackModelVariable = "Mack Simulated Volume Weighted Development Factor"

			for k := 0; k <= projMonths+TriangleVariableCount; k++ {
				if k > TriangleVariableCount {
					cumTriaColSum.Field(k - VariableCountDiff).SetFloat(cumTriaColSum.Field(k-VariableCountDiff).Float() + cumTria3.Field(k).Float())
					if k == TriangleVariableCount+1 {
						mackSimDF.Field(k - VariableCountDiff).SetFloat(0)
					} else {
						mackSimDF.Field(k - VariableCountDiff).SetFloat(mackSimDF.Field(k-VariableCountDiff).Float() + pseudoR2.Field(k).Float()*cumTria3.Field(k-1).Float())

					}
				}
			}

			//Mack Volume Weighted Development Factor
			if i == len(cumulativeTriangulations)-1 {
				for k := TriangleVariableCount + 1; k <= projMonths+TriangleVariableCount; k++ {
					if k == TriangleVariableCount+1 {
						mackSimDF.Field(k - VariableCountDiff).SetFloat(0)
					} else {
						mackSimDF.Field(k - VariableCountDiff).SetFloat(mackSimDF.Field(k-VariableCountDiff).Float() / cumTriaColSum.Field(k-VariableCountDiff-1).Float())
					}
				}
			}
		}

		if mackSim == run.MackModelSimulations {
			DB.Save(&cumulativeTriagulationColumnSum)
			DB.Save(&mackModelSimulatedDevelopmentFactor)
		}

		//Cumulative Projections using Mack Simulated Development Factors
		//𝑆_(𝜔,𝑑) 𝑆_(𝜔_, 𝑑+1)=𝑢_(𝜔,𝑑+1)+𝑟𝜎_(𝜔,𝑑+1),    𝑢_(𝜔,𝑑+1)=𝑆_(𝑤,𝑑) 𝑅_𝑑
		//Where S is cumulated claims
		var mackCumulativeProjections []models.LicMackCumulativeProjection
		for i, _ := range cumulativeTriangulations {
			var mackCumulativeProjection models.LicMackCumulativeProjection
			cumTria4 := reflect.ValueOf(&cumulativeTriangulations[i]).Elem()
			//pseudoR := reflect.ValueOf(&pseudoRatio).Elem()
			mackSimDF := reflect.ValueOf(&mackModelSimulatedDevelopmentFactor).Elem()
			mackCumProj := reflect.ValueOf(&mackCumulativeProjection).Elem()
			//vwDevFac2 := reflect.ValueOf(&volumeWeightedDevFac).Elem()
			var yearadd, movingyear int
			var residualmonth, movingmonth int
			for k := 0; k <= projMonths+TriangleVariableCount; k++ {
				if k <= TriangleVariableCount {
					mackCumProj.Field(k).Set(cumTria4.Field(k))
				} else {
					if interval == 12 {
						yearadd = int(int64((int(cumTria4.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1) / (12 / interval)))
						residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(cumTria4.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
					}
					if interval != 12 {
						yearadd = int(int64((int(cumTria4.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1 - 1) / (12 / interval)))
						residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(cumTria4.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
					}
					residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(cumTria4.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
					movingyear = int(cumTria4.FieldByName(AccidentYear).Int()) + yearadd
					if residualmonth == 0 && interval == 3 {
						residualmonth = 1
					}
					if residualmonth == 0 && interval == 1 {
						residualmonth = 12
					}
					movingmonth = residualmonth
					if (movingyear <= endDateYear && movingmonth <= endDateMonth) || (movingyear < endDateYear && movingmonth > endDateMonth) {
						mackCumProj.Field(k).SetFloat(cumTria4.Field(k).Float())
					} else {
						//if k <= projMonths+TriangleVariableCount {
						mackCumProj.Field(k).SetFloat(mackCumProj.Field(k-1).Float() * mackSimDF.Field(k-VariableCountDiff).Float())
						//} else {
						//	mackCumProj.Field(k).SetFloat(0)
						//}
					}
				}
			}
			mackCumulativeProjections = append(mackCumulativeProjections, mackCumulativeProjection)
		}

		if mackSim == run.MackModelSimulations {
			DB.Save(&mackCumulativeProjections)
		}

		//Mack Reserve
		var mackReserve float64
		for i, _ := range mackCumulativeProjections {
			mackCumProj2 := reflect.ValueOf(&mackCumulativeProjections[i]).Elem()
			var yearadd, movingyear int
			var residualmonth, movingmonth int
			for k := TriangleVariableCount + 1; k <= projMonths+TriangleVariableCount; k++ {
				if interval == 12 {
					yearadd = int(int64((int(mackCumProj2.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1) / (12 / interval)))
					residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mackCumProj2.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				}
				if interval != 12 {
					yearadd = int(int64((int(mackCumProj2.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1 - 1) / (12 / interval)))
					residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mackCumProj2.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				}
				residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mackCumProj2.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				movingyear = int(mackCumProj2.FieldByName(AccidentYear).Int()) + yearadd
				if residualmonth == 0 && interval == 3 {
					residualmonth = 1
				}
				if residualmonth == 0 && interval == 1 {
					residualmonth = 12
				}
				movingmonth = residualmonth
				if movingyear == endDateYear && movingmonth == endDateMonth {
					//if mackCumProj2.Field(k).Kind() == reflect.Int{
					//
					//}
					fmt.Println(mackCumProj2.Field(k).Kind())
					fmt.Println(mackCumProj2.Field(projMonths + TriangleVariableCount).Kind())
					fmt.Println(k, projMonths+TriangleVariableCount)
					(*ibnrReserves)[i].MackModelIbnr = mackCumProj2.Field(projMonths+TriangleVariableCount).Float() - mackCumProj2.Field(k).Float()
					mackReserve += (*ibnrReserves)[i].MackModelIbnr
				}
			}

		}

		//Simulations
		var MackSimulationResult models.LicMackSimulationResults

		MackSimulationResult.RunDate = mackCumulativeProjections[0].RunDate
		MackSimulationResult.RunID = mackCumulativeProjections[0].RunID
		MackSimulationResult.PortfolioName = mackCumulativeProjections[0].PortfolioName
		MackSimulationResult.LicPortfolioId = mackCumulativeProjections[0].LicPortfolioId
		MackSimulationResult.ProductCode = mackCumulativeProjections[0].ProductCode
		MackSimulationResult.SimulationNumber = mackSim
		MackSimulationResult.Reserve = mackReserve

		MackSimulationResults = append(MackSimulationResults, MackSimulationResult)

		if mackSim == run.MackModelSimulations { // saving mack model simulation results
			DB.Where("run_date=? and product_code=? and run_id=?", runDate, MackSimulationResults[0].ProductCode, MackSimulationResults[0].RunID).Delete(&models.LicMackSimulationResults{})
			err := DB.Save(&MackSimulationResults).Error
			if err != nil {
				log.Println(err)
			}
		}
	}

	// Summary Stats
	var err error

	query := fmt.Sprintf("SELECT run_date,run_id,portfolio_name,lic_portfolio_id,product_code,avg(reserve) as mean, stddev(reserve) as standard_deviation, min(reserve) as minimum, max(reserve) as maximum FROM lic_mack_simulation_results where run_date = '%s' and product_code = '%s' and portfolio_name = '%s' and lic_portfolio_id = %d and run_id = %d", run.RunDate, productCode, run.PortfolioName, run.PortfolioId, run.ID)
	var mackSimulationSummaryStats []models.LicMackSimulationSummaryStats
	err = DB.Raw(query).Scan(&mackSimulationSummaryStats).Error
	if err != nil {
		fmt.Println(err)
	}

	medianQuery2 := fmt.Sprintf("select reserve FROM lic_mack_simulation_results where run_date = '%s' and product_code = '%s' and portfolio_name = '%s' and lic_portfolio_id = %d and run_id = %d order by reserve", run.RunDate, productCode, run.PortfolioName, run.PortfolioId, run.ID)
	var medianData []float64
	var nthpercentile float64
	var median, mean float64
	err = DB.Raw(medianQuery2).Scan(&medianData).Error
	median, _ = stats.Median(medianData)
	if params.RiskAdjustmentConfidenceLevel == 0 {
		nthpercentile = 0
	} else {
		nthpercentile, _ = stats.Percentile(medianData, params.RiskAdjustmentConfidenceLevel*100.0)
	}

	mean, _ = stats.Mean(medianData)

	mackSimulationSummaryStats[0].Percentile = nthpercentile
	mackSimulationSummaryStats[0].Median = median
	mackSimulationSummaryStats[0].Mean = mean
	(*IbnrReserveReport).MackModelIbnr = mackSimulationSummaryStats[0].Mean
	(*IbnrReserveReport).MackModelNthPercentileIbnr = nthpercentile

	//DB.Where("run_date=? and product_code=?", runDate, results[0].ProductCode).Delete(&models.LicBootstrappedResultSummary{})
	err = DB.Save(&mackSimulationSummaryStats).Error
	if err != nil {
		log.Println(err)
	}

	//Frequency Table
	var mackFrequencyTable []models.MackIbnrFrequency
	groupStep := (mackSimulationSummaryStats[0].Maximum - mackSimulationSummaryStats[0].Minimum) / 15
	for freqGroup := 0; freqGroup <= 15; freqGroup++ {
		var mackFrequency models.MackIbnrFrequency
		mackFrequency.RunDate = mackSimulationSummaryStats[0].RunDate
		mackFrequency.RunID = mackSimulationSummaryStats[0].RunID
		mackFrequency.PortfolioName = mackSimulationSummaryStats[0].PortfolioName
		mackFrequency.LicPortfolioId = mackSimulationSummaryStats[0].LicPortfolioId
		mackFrequency.ProductCode = mackSimulationSummaryStats[0].ProductCode
		mackFrequency.Group = freqGroup + 1
		if groupStep > 0 {
			if freqGroup == 0 {
				mackFrequency.Reserve = math.Floor(mackSimulationSummaryStats[0].Minimum)
			} else {
				if freqGroup == 15 {
					mackFrequency.Reserve = math.Ceil(mackSimulationSummaryStats[0].Maximum)
				} else {
					mackFrequency.Reserve = mackFrequencyTable[freqGroup-1].Reserve + groupStep
				}

			}
			mackFrequencyTable = append(mackFrequencyTable, mackFrequency)
		}
		if groupStep == 0 && freqGroup == 0 { //adds one line if gap is zero else no
			mackFrequencyTable = append(mackFrequencyTable, mackFrequency)
			freqGroup = 16
		}

	}

	if groupStep > 0 {
		for freqGroup := 0; freqGroup < 15; freqGroup++ {
			lowerBound := mackFrequencyTable[freqGroup].Reserve
			upperBound := mackFrequencyTable[freqGroup+1].Reserve
			query := fmt.Sprintf("SELECT count(reserve) as count FROM lic_mack_simulation_results where reserve >= '%f' and reserve < '%f' and run_date = '%s' and product_code = '%s' and run_id = %d", lowerBound, upperBound, run.RunDate, productCode, mackSimulationSummaryStats[0].RunID)
			var count int
			err = DB.Raw(query).Scan(&count).Error
			if err != nil {
				fmt.Println(err)
			}
			mackFrequencyTable[freqGroup].Frequency = count
		}
	}

	//DB.Where("run_date=? and product_code=?", runDate, frequencytable[0].ProductCode).Delete(&models.IbnrFrequency{})
	err = DB.Save(&mackFrequencyTable).Error
	if err != nil {
		log.Println(err)
	}

}

func SaveIBNRShockSetting(shockSetting models.IBNRShockSetting) (models.IBNRShockSetting, error) {
	DB.Save(&shockSetting)
	return shockSetting, nil
}

func GetIBNRShockSettings() ([]models.IBNRShockSetting, error) {
	var shockSettings []models.IBNRShockSetting
	err := DB.Find(&shockSettings).Error
	return shockSettings, err
}

func DeleteIBNRShockSetting(id int) error {
	err := DB.Where("id = ?", id).Delete(&models.IBNRShockSetting{}).Error
	if err != nil {
		return err
	}
	return nil
}

func BootstrappedDevelopmentFactor(bootStrappedCumulatives []models.LicBootStrappedCumulative, runDate string, endDateYear, endDateMonth, projMonths int, params models.LICParameter, run models.IBNRRunSetting, sim, interval int) models.LicBootstrappedResults {

	var columnSum models.LicBootstrappedDevelopmentFactor

	columnSum.PortfolioName = bootStrappedCumulatives[0].PortfolioName
	columnSum.ProductCode = bootStrappedCumulatives[0].ProductCode
	columnSum.RunDate = bootStrappedCumulatives[0].RunDate
	columnSum.RunID = bootStrappedCumulatives[0].RunID
	columnSum.DevelopmentVariable = "Column Sum"

	for _, ct := range bootStrappedCumulatives {
		columnSum.Rd0 += ct.Rd0
		columnSum.Rd1 += ct.Rd1
		columnSum.Rd2 += ct.Rd2
		columnSum.Rd3 += ct.Rd3
		columnSum.Rd4 += ct.Rd4
		columnSum.Rd5 += ct.Rd5
		columnSum.Rd6 += ct.Rd6
		columnSum.Rd7 += ct.Rd7
		columnSum.Rd8 += ct.Rd8
		columnSum.Rd9 += ct.Rd9
		columnSum.Rd10 += ct.Rd10
		columnSum.Rd11 += ct.Rd11
		columnSum.Rd12 += ct.Rd12
		columnSum.Rd13 += ct.Rd13
		columnSum.Rd14 += ct.Rd14
		columnSum.Rd15 += ct.Rd15
		columnSum.Rd16 += ct.Rd16
		columnSum.Rd17 += ct.Rd17
		columnSum.Rd18 += ct.Rd18
		columnSum.Rd19 += ct.Rd19
		columnSum.Rd20 += ct.Rd20
		columnSum.Rd21 += ct.Rd21
		columnSum.Rd22 += ct.Rd22
		columnSum.Rd23 += ct.Rd23
		columnSum.Rd24 += ct.Rd24
		columnSum.Rd25 += ct.Rd25
		columnSum.Rd26 += ct.Rd26
		columnSum.Rd27 += ct.Rd27
		columnSum.Rd28 += ct.Rd28
		columnSum.Rd29 += ct.Rd29
		columnSum.Rd30 += ct.Rd30
		columnSum.Rd31 += ct.Rd31
		columnSum.Rd32 += ct.Rd32
		columnSum.Rd33 += ct.Rd33
		columnSum.Rd34 += ct.Rd34
		columnSum.Rd35 += ct.Rd35
		columnSum.Rd36 += ct.Rd36
		columnSum.Rd37 += ct.Rd37
		columnSum.Rd38 += ct.Rd38
		columnSum.Rd39 += ct.Rd39
		columnSum.Rd40 += ct.Rd40
		columnSum.Rd41 += ct.Rd41
		columnSum.Rd42 += ct.Rd42
		columnSum.Rd43 += ct.Rd43
		columnSum.Rd44 += ct.Rd44
		columnSum.Rd45 += ct.Rd45
		columnSum.Rd46 += ct.Rd46
		columnSum.Rd47 += ct.Rd47
		columnSum.Rd48 += ct.Rd48
		columnSum.Rd49 += ct.Rd49
		columnSum.Rd50 += ct.Rd50
		columnSum.Rd51 += ct.Rd51
		columnSum.Rd52 += ct.Rd52
		columnSum.Rd53 += ct.Rd53
		columnSum.Rd54 += ct.Rd54
		columnSum.Rd55 += ct.Rd55
		columnSum.Rd56 += ct.Rd56
		columnSum.Rd57 += ct.Rd57
		columnSum.Rd58 += ct.Rd58
		columnSum.Rd59 += ct.Rd59
		columnSum.Rd60 += ct.Rd60
		columnSum.Rd61 += ct.Rd61
		columnSum.Rd62 += ct.Rd62
		columnSum.Rd63 += ct.Rd63
		columnSum.Rd64 += ct.Rd64
		columnSum.Rd65 += ct.Rd65
		columnSum.Rd66 += ct.Rd66
		columnSum.Rd67 += ct.Rd67
		columnSum.Rd68 += ct.Rd68
		columnSum.Rd69 += ct.Rd69
		columnSum.Rd70 += ct.Rd70
		columnSum.Rd71 += ct.Rd71
		columnSum.Rd72 += ct.Rd72
		columnSum.Rd73 += ct.Rd73
		columnSum.Rd74 += ct.Rd74
		columnSum.Rd75 += ct.Rd75
		columnSum.Rd76 += ct.Rd76
		columnSum.Rd77 += ct.Rd77
		columnSum.Rd78 += ct.Rd78
		columnSum.Rd79 += ct.Rd79
		columnSum.Rd80 += ct.Rd80
		columnSum.Rd81 += ct.Rd81
		columnSum.Rd82 += ct.Rd82
		columnSum.Rd83 += ct.Rd83
		columnSum.Rd84 += ct.Rd84
		columnSum.Rd85 += ct.Rd85
		columnSum.Rd86 += ct.Rd86
		columnSum.Rd87 += ct.Rd87
		columnSum.Rd88 += ct.Rd88
		columnSum.Rd89 += ct.Rd89
		columnSum.Rd90 += ct.Rd90
		columnSum.Rd91 += ct.Rd91
		columnSum.Rd92 += ct.Rd92
		columnSum.Rd93 += ct.Rd93
		columnSum.Rd94 += ct.Rd94
		columnSum.Rd95 += ct.Rd95
		columnSum.Rd96 += ct.Rd96
		columnSum.Rd97 += ct.Rd97
		columnSum.Rd98 += ct.Rd98
		columnSum.Rd99 += ct.Rd99
		columnSum.Rd100 += ct.Rd100
		columnSum.Rd101 += ct.Rd101
		columnSum.Rd102 += ct.Rd102
		columnSum.Rd103 += ct.Rd103
		columnSum.Rd104 += ct.Rd104
		columnSum.Rd105 += ct.Rd105
		columnSum.Rd106 += ct.Rd106
		columnSum.Rd107 += ct.Rd107
		columnSum.Rd108 += ct.Rd108
		columnSum.Rd109 += ct.Rd109
		columnSum.Rd110 += ct.Rd110
		columnSum.Rd111 += ct.Rd111
		columnSum.Rd112 += ct.Rd112
		columnSum.Rd113 += ct.Rd113
		columnSum.Rd114 += ct.Rd114
		columnSum.Rd115 += ct.Rd115
		columnSum.Rd116 += ct.Rd116
		columnSum.Rd117 += ct.Rd117
		columnSum.Rd118 += ct.Rd118
		columnSum.Rd119 += ct.Rd119
		columnSum.Rd120 += ct.Rd120
	}

	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, columnSum.DevelopmentVariable, columnSum.ProductCode).Delete(&models.LicBootstrappedDevelopmentFactor{})
	DB.Save(&columnSum)

	var lastvalue models.LicBootstrappedDevelopmentFactor

	lastvalue.PortfolioName = bootStrappedCumulatives[0].PortfolioName
	lastvalue.ProductCode = bootStrappedCumulatives[0].ProductCode
	lastvalue.RunDate = bootStrappedCumulatives[0].RunDate
	lastvalue.RunID = bootStrappedCumulatives[0].RunID
	lastvalue.DevelopmentVariable = "Last Value"

	mlv := reflect.ValueOf(&lastvalue).Elem()

	for i, _ := range bootStrappedCumulatives {
		mbsc2 := reflect.ValueOf(&bootStrappedCumulatives[i]).Elem()

		var yearadd, movingyear int
		var residualmonth, movingmonth int

		for k := TriangleVariableCount + 1; k <= projMonths+TriangleVariableCount; k++ {
			if interval == 12 {
				yearadd = int(int64((int(mbsc2.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1) / (12 / interval)))
				residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mbsc2.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
			}
			if interval != 12 {
				yearadd = int(int64((int(mbsc2.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1 - 1) / (12 / interval)))
				residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mbsc2.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
			}
			movingyear = int(mbsc2.FieldByName(AccidentYear).Int()) + yearadd
			if residualmonth == 0 && interval == 3 {
				residualmonth = 1
			}
			if residualmonth == 0 && interval == 1 {
				residualmonth = 12
			}
			movingmonth = residualmonth

			if movingyear == endDateYear && movingmonth == endDateMonth {
				mlv.Field(k - VariableCountDiff).Set(mbsc2.Field(k))
			}
		}
	}

	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, lastvalue.DevelopmentVariable, lastvalue.ProductCode).Delete(&models.LicBootstrappedDevelopmentFactor{})
	DB.Save(&lastvalue)

	// column excluding last value
	var columnsumexcllastvalue models.LicBootstrappedDevelopmentFactor

	columnsumexcllastvalue.PortfolioName = bootStrappedCumulatives[0].PortfolioName
	columnsumexcllastvalue.ProductCode = bootStrappedCumulatives[0].ProductCode
	columnsumexcllastvalue.RunDate = bootStrappedCumulatives[0].RunDate
	columnsumexcllastvalue.RunID = bootStrappedCumulatives[0].RunID
	columnsumexcllastvalue.DevelopmentVariable = "Column Sum(excl.last value)"

	mcs := reflect.ValueOf(&columnSum).Elem()
	mlv2 := reflect.ValueOf(&lastvalue).Elem()
	mcselv := reflect.ValueOf(&columnsumexcllastvalue).Elem()

	for k := 0; k <= projMonths+FactorVariableCount; k++ {
		if k > FactorVariableCount {
			if mcs.Field(k).Float() > 0 {
				mcselv.Field(k).SetFloat(mcs.Field(k).Float() - mlv2.Field(k).Float())
				continue
			}
		}
	}
	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, columnsumexcllastvalue.DevelopmentVariable, columnsumexcllastvalue.ProductCode).Delete(&models.LicBootstrappedDevelopmentFactor{})
	DB.Save(&columnsumexcllastvalue)

	//Weighted average succession ratios
	var successionratios models.LicBootstrappedDevelopmentFactor
	successionratios.PortfolioName = bootStrappedCumulatives[0].PortfolioName
	successionratios.ProductCode = bootStrappedCumulatives[0].ProductCode
	successionratios.RunDate = bootStrappedCumulatives[0].RunDate
	successionratios.RunID = bootStrappedCumulatives[0].RunID
	successionratios.DevelopmentVariable = "Weighted Ave. Succession Ratio(rd to rd+1)"

	mcs2 := reflect.ValueOf(&columnSum).Elem()
	mcselv2 := reflect.ValueOf(&columnsumexcllastvalue).Elem()
	msr3 := reflect.ValueOf(&successionratios).Elem()

	for k := 0; k <= projMonths+FactorVariableCount; k++ {
		if k > FactorVariableCount {
			if k == projMonths+FactorVariableCount {
				msr3.Field(k).SetFloat(1)
			} else {
				if mcselv2.Field(k).Float() > 0 && k < projMonths+FactorVariableCount {
					msr3.Field(k).SetFloat(mcs2.Field(k+1).Float() / mcselv2.Field(k).Float())
				}
			}
		}
	}
	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, successionratios.DevelopmentVariable, successionratios.ProductCode).Delete(&models.LicBootstrappedDevelopmentFactor{})
	DB.Save(&successionratios)

	//development factor
	var developmentfactor models.LicBootstrappedDevelopmentFactor
	developmentfactor.PortfolioName = bootStrappedCumulatives[0].PortfolioName
	developmentfactor.ProductCode = bootStrappedCumulatives[0].ProductCode
	developmentfactor.RunDate = bootStrappedCumulatives[0].RunDate
	developmentfactor.RunID = bootStrappedCumulatives[0].RunID
	developmentfactor.DevelopmentVariable = "development factor"

	msr4 := reflect.ValueOf(&successionratios).Elem()
	mdf := reflect.ValueOf(&developmentfactor).Elem()

	for k := 0; k <= projMonths+FactorVariableCount; k++ {
		if k > FactorVariableCount {
			if k == FactorVariableCount+1 {
				mdf.Field(k).SetFloat(0)
			} else {
				mdf.Field(k).SetFloat(msr4.Field(k - 1).Float())
			}

		}
	}
	//DB.Where("run_date=? and development_variable=? and product_code=?", runDate, developmentfactor.DevelopmentVariable, developmentfactor.ProductCode).Delete(&models.LicBootstrappedDevelopmentFactor{})
	err := DB.Save(&developmentfactor).Error
	if err != nil {
		log.Println(err)
	}

	//Bootstrapped projections
	var BootstrappedCumulativeProjections = []models.LicBootStrappedCumulativeProjection{}
	var incrementalProjections = []models.LicBootStrappedIncrementalProjection{}
	var inflatedincrementalProjections = []models.LicBootStrappedIncrementalInflatedProjection{}
	var discountedinflatedincrementalProjections = []models.LicBootStrappedIncrementalInflatedDiscountedProjection{}
	for i, _ := range bootStrappedCumulatives {
		var BootstrappedCumulativeProjection models.LicBootStrappedCumulativeProjection
		var incrementalProjection models.LicBootStrappedIncrementalProjection
		var inflatedincrementalProjection models.LicBootStrappedIncrementalInflatedProjection
		var discountedinflatedincrementalProjection models.LicBootStrappedIncrementalInflatedDiscountedProjection
		mbsc3 := reflect.ValueOf(&bootStrappedCumulatives[i]).Elem()
		msr5 := reflect.ValueOf(&successionratios).Elem()
		mbscp := reflect.ValueOf(&BootstrappedCumulativeProjection).Elem()

		mip := reflect.ValueOf(&incrementalProjection).Elem() // projected total claims based on projected average claim and projected claim count

		miip := reflect.ValueOf(&inflatedincrementalProjection).Elem() //inflated projected total claims

		mdiip := reflect.ValueOf(&discountedinflatedincrementalProjection).Elem() //discounted inflated  projected total claims

		var yearadd, movingyear int
		var residualmonth, movingmonth int
		for k := 0; k <= projMonths+TriangleVariableCount; k++ {
			if k <= TriangleVariableCount {
				mbscp.Field(k).Set(mbsc3.Field(k))
				mip.Field(k).Set(mbsc3.Field(k))
				miip.Field(k).Set(mbsc3.Field(k))
				mdiip.Field(k).Set(mbsc3.Field(k))

			} else {
				if interval == 12 {
					yearadd = int(int64((int(mbsc3.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1) / (12 / interval)))
					residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mbsc3.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				}
				if interval != 12 {
					yearadd = int(int64((int(mbsc3.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1 - 1) / (12 / interval)))
					residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mbsc3.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				}
				movingyear = int(mbsc3.FieldByName(AccidentYear).Int()) + yearadd
				if residualmonth == 0 && interval == 3 {
					residualmonth = 1
				}
				if residualmonth == 0 && interval == 1 {
					residualmonth = 12
				}
				movingmonth = residualmonth
				if k == TriangleVariableCount+1 {
					mbscp.Field(k).SetFloat(mbsc3.Field(k).Float())
					mip.Field(k).SetFloat(mbsc3.Field(k).Float())
				} else {
					if (movingyear <= endDateYear && movingmonth <= endDateMonth) || (movingyear < endDateYear && movingmonth > endDateMonth) {
						mbscp.Field(k).SetFloat(mbsc3.Field(k).Float())
						mip.Field(k).SetFloat(mbscp.Field(k).Float() - mbscp.Field(k-1).Float())
					} else {
						if k <= projMonths+TriangleVariableCount {
							mbscp.Field(k).SetFloat(mbscp.Field(k-1).Float() * msr5.Field(k-VariableCountDiff-1).Float())
							mip.Field(k).SetFloat(mbscp.Field(k).Float() - mbscp.Field(k-1).Float())
							miip.Field(k).SetFloat(mip.Field(k).Float() * 1)   // inflated
							mdiip.Field(k).SetFloat(miip.Field(k).Float() * 1) // discounted
						}

					}
				}
			}
		}
		BootstrappedCumulativeProjections = append(BootstrappedCumulativeProjections, BootstrappedCumulativeProjection)
		incrementalProjections = append(incrementalProjections, incrementalProjection)
		inflatedincrementalProjections = append(inflatedincrementalProjections, inflatedincrementalProjection)
		discountedinflatedincrementalProjections = append(discountedinflatedincrementalProjections, discountedinflatedincrementalProjection)
	}

	//DB.Where("run_date=? and product_code=?", runDate, BootstrappedCumulativeProjections[0].ProductCode).Delete(&models.LicBootStrappedCumulativeProjection{})
	err = DB.Save(&BootstrappedCumulativeProjections).Error
	if err != nil {
		log.Println(err)
	}

	//DB.Where("run_date=? and product_code=?", runDate, incrementalProjections[0].ProductCode).Delete(&models.LicBootStrappedIncrementalProjection{})
	err = DB.Save(&incrementalProjections).Error
	if err != nil {
		log.Println(err)
	}

	//DB.Where("run_date=? and product_code=?", runDate, inflatedincrementalProjections[0].ProductCode).Delete(&models.LicBootStrappedIncrementalInflatedProjection{})
	err = DB.Save(&inflatedincrementalProjections).Error
	if err != nil {
		log.Println(err)
	}

	//DB.Where("run_date=? and product_code=?", runDate, discountedinflatedincrementalProjections[0].ProductCode).Delete(&models.LicBootStrappedIncrementalInflatedDiscountedProjection{})
	err = DB.Save(&discountedinflatedincrementalProjections).Error
	if err != nil {
		log.Println(err)
	}

	var bootStrappedResult models.LicBootstrappedResults
	var aggResults float64
	for i, _ := range discountedinflatedincrementalProjections {
		mdiipx := reflect.ValueOf(&discountedinflatedincrementalProjections[i]).Elem()
		var yearadd, movingyear int
		var residualmonth, movingmonth int

		for k := 0; k <= projMonths+TriangleVariableCount; k++ {

			if k >= TriangleVariableCount+1 {
				if interval == 12 {
					yearadd = int(int64((int(mdiipx.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1) / (12 / interval)))
					residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mdiipx.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				}
				if interval != 12 {
					yearadd = int(int64((int(mdiipx.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1 - 1) / (12 / interval)))
					residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mdiipx.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				}
				movingyear = int(mdiipx.FieldByName(AccidentYear).Int()) + yearadd
				if residualmonth == 0 && interval == 3 {
					residualmonth = 1
				}
				if residualmonth == 0 && interval == 1 {
					residualmonth = 12
				}
				movingmonth = residualmonth
				if (movingyear <= endDateYear && movingmonth <= endDateMonth) || (movingyear < endDateYear && movingmonth > endDateMonth) {
					continue
				} else {
					aggResults += mdiipx.Field(k).Float()
				}
			}
		}

	}

	bootStrappedResult.RunDate = discountedinflatedincrementalProjections[0].RunDate
	bootStrappedResult.RunID = discountedinflatedincrementalProjections[0].RunID
	bootStrappedResult.PortfolioName = discountedinflatedincrementalProjections[0].PortfolioName
	bootStrappedResult.LicPortfolioId = discountedinflatedincrementalProjections[0].LicPortfolioId
	bootStrappedResult.ProductCode = discountedinflatedincrementalProjections[0].ProductCode
	bootStrappedResult.SimulationNumber = sim
	bootStrappedResult.Reserve = aggResults

	return bootStrappedResult
}

func CumulativeAveragetoTotalClaim(cumulativeProjectionsClaimCount []models.LicCumulativeProjectionClaimCount, cumulativeProjectionsAverageClaim []models.LicCumulativeProjectionAverageClaim, runDate, projMonths, endDateYear, endDateMonth, interval int) {
	var cumulativeProjectionsAveragetoTotalClaim = []models.LicCumulativeProjectionAveragetoTotalClaim{} //

	for i, _ := range cumulativeProjectionsClaimCount {
		var cumulativeProjectionAveragetoTotalClaim models.LicCumulativeProjectionAveragetoTotalClaim
		mcpCC := reflect.ValueOf(&cumulativeProjectionsClaimCount[i]).Elem()
		mcpAC := reflect.ValueOf(&cumulativeProjectionsAverageClaim[i]).Elem() // projected average claim count

		mcpATC := reflect.ValueOf(&cumulativeProjectionAveragetoTotalClaim).Elem() // projected average claim
		var yearadd, movingyear int
		var residualmonth, movingmonth int
		for k := 0; k <= projMonths+FactorVariableCount; k++ {
			if k <= FactorVariableCount {
				mcpATC.Field(k).Set(mcpCC.Field(k))
			} else {
				if interval == 12 {
					yearadd = int(int64((int(mcpCC.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1) / (12 / interval)))
					residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mcpCC.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				}
				if interval != 12 {
					yearadd = int(int64((int(mcpCC.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1 - 1) / (12 / interval)))
					residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mcpCC.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				}
				movingyear = int(mcpCC.FieldByName(AccidentYear).Int()) + yearadd
				if residualmonth == 0 && interval == 3 {
					residualmonth = 1
				}
				if residualmonth == 0 && interval == 1 {
					residualmonth = 12
				}
				movingmonth = residualmonth
				if k == FactorVariableCount+1 {
					mcpATC.Field(k).SetFloat(utils.FloatPrecision(mcpCC.Field(k).Float()*mcpAC.Field(k).Float(), AccountingPrecision))
				} else {
					if (movingyear <= endDateYear && movingmonth <= endDateMonth) || (movingyear < endDateYear && movingmonth > endDateMonth) {
						mcpATC.Field(k).SetFloat(utils.FloatPrecision(mcpCC.Field(k).Float()*mcpAC.Field(k).Float(), AccountingPrecision))
					} else {
						if k <= projMonths+FactorVariableCount {
							mcpATC.Field(k).SetFloat(utils.FloatPrecision(mcpCC.Field(k).Float()*mcpAC.Field(k).Float(), AccountingPrecision))
						}

					}
				}
			}
		}
		cumulativeProjectionsAveragetoTotalClaim = append(cumulativeProjectionsAveragetoTotalClaim, cumulativeProjectionAveragetoTotalClaim)
	}

	DB.Where("run_date=? and product_code=?", runDate, cumulativeProjectionsAveragetoTotalClaim[0].ProductCode).Delete(&models.LicCumulativeProjectionAveragetoTotalClaim{})
	err := DB.Save(&cumulativeProjectionsAveragetoTotalClaim).Error
	if err != nil {
		log.Println(err)
	}
}

func InflatedDiscountedProjections(cumulativeProjections *[]models.LicCumulativeProjection, run models.IBNRRunSetting, runDate string, projMonths, endDateYear, endDateMonth, interval int, discountOption string, yieldCurveCode string) []models.LicCumulativeProjection {
	var incrementalProjections = []models.LicIncrementalProjection{} //
	var inflatedincrementalProjections = []models.LicIncrementalInflated{}
	var discountedinflatedincrementalProjections = []models.LicDiscountedIncrementalInflated{}
	for i, _ := range *cumulativeProjections {
		var incrementalProjection models.LicIncrementalProjection
		var inflatedincrementalProjection models.LicIncrementalInflated
		var discountedinflatedincrementalProjection models.LicDiscountedIncrementalInflated

		mcp := reflect.ValueOf(&(*cumulativeProjections)[i]).Elem()

		mip := reflect.ValueOf(&incrementalProjection).Elem() // projected total claims based on projected average claim and projected claim count

		miip := reflect.ValueOf(&inflatedincrementalProjection).Elem() //inflated projected total claims

		mdiip := reflect.ValueOf(&discountedinflatedincrementalProjection).Elem() //discounted inflated  projected total claims

		var yearadd, movingyear int
		var residualmonth, movingmonth int
		for k := 0; k <= projMonths+TriangleVariableCount; k++ {
			if k <= TriangleVariableCount {
				mip.Field(k).Set(mcp.Field(k))
				miip.Field(k).Set(mcp.Field(k))
				mdiip.Field(k).Set(mcp.Field(k))

			} else {
				//fmt.Println(mct.Type().Field(k).Name)
				//fmt.Println(projMonths)
				if interval == 12 {
					yearadd = int(int64((int(mcp.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1) / (12 / interval)))
					residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mcp.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				}
				if interval != 12 {
					yearadd = int(int64((int(mcp.FieldByName(AccidentMonth).Int()) + k - TriangleVariableCount - 1 - 1) / (12 / interval)))
					residualmonth = int(math.Mod(float64(k-TriangleVariableCount-1+int(mcp.FieldByName(AccidentMonth).Int())), 12/float64(interval)))
				}
				movingyear = int(mcp.FieldByName(AccidentYear).Int()) + yearadd
				if residualmonth == 0 && interval == 3 {
					residualmonth = 1
				}
				if residualmonth == 0 && interval == 1 {
					residualmonth = 12
				}
				movingmonth = residualmonth
				if k == TriangleVariableCount+1 {
					mip.Field(k).SetFloat(utils.FloatPrecision(mcp.Field(k).Float(), AccountingPrecision))
					miip.Field(k).SetFloat(mip.Field(k).Float())
					mdiip.Field(k).SetFloat(miip.Field(k).Float())
				} else {
					if (movingyear <= endDateYear && movingmonth <= endDateMonth) || (movingyear < endDateYear && movingmonth > endDateMonth) {
						mip.Field(k).SetFloat(utils.FloatPrecision(mcp.Field(k).Float()-mcp.Field(k-1).Float(), AccountingPrecision))
						miip.Field(k).SetFloat(mip.Field(k).Float())
						mdiip.Field(k).SetFloat(miip.Field(k).Float())
					} else {
						if k <= projMonths+TriangleVariableCount {
							var inflationFactor float64
							var discountFactor float64
							if run.InflationIndicator {
								inflationFactor = IbnrInflationFactor(i, k-TriangleVariableCount-1, projMonths, run, interval, yieldCurveCode)
							} else {
								inflationFactor = 1
							}
							if discountOption == Discounted {
								discountFactor = IbnrDiscountFactor(i, k-TriangleVariableCount-1, projMonths, run, interval, yieldCurveCode, run.Basis)
							} else {
								discountFactor = 1
							}

							mip.Field(k).SetFloat(utils.FloatPrecision(mcp.Field(k).Float()-mcp.Field(k-1).Float(), AccountingPrecision))
							miip.Field(k).SetFloat(utils.FloatPrecision(mip.Field(k).Float()*inflationFactor, AccountingPrecision))
							mdiip.Field(k).SetFloat(utils.FloatPrecision(miip.Field(k).Float()*discountFactor, AccountingPrecision))
						}

					}
				}
			}
		}
		incrementalProjections = append(incrementalProjections, incrementalProjection)
		inflatedincrementalProjections = append(inflatedincrementalProjections, inflatedincrementalProjection)
		discountedinflatedincrementalProjections = append(discountedinflatedincrementalProjections, discountedinflatedincrementalProjection)
	}

	//DB.Where("run_date=? and product_code=?", runDate, incrementalProjections[0].ProductCode).Delete(&models.LicIncrementalProjection{})
	err := DB.Save(&incrementalProjections).Error
	if err != nil {
		log.Println(err)
	}

	//DB.Where("run_date=? and product_code=?", runDate, inflatedincrementalProjections[0].ProductCode).Delete(&models.LicIncrementalInflated{})
	err = DB.Save(&inflatedincrementalProjections).Error
	if err != nil {
		log.Println(err)
	}

	//DB.Where("run_date=? and product_code=?", runDate, discountedinflatedincrementalProjections[0].ProductCode).Delete(&models.LicDiscountedIncrementalInflated{})
	err = DB.Save(&discountedinflatedincrementalProjections).Error
	if err != nil {
		log.Println(err)
	}

	var timeValueAdjustedCumulativeProjections []models.LicCumulativeProjection
	for i, _ := range discountedinflatedincrementalProjections {
		var timeValueAdjustedCumulativeProjection models.LicCumulativeProjection
		mdiip2 := reflect.ValueOf(&discountedinflatedincrementalProjections[i]).Elem()
		mtvacp := reflect.ValueOf(&timeValueAdjustedCumulativeProjection).Elem()
		for k := 0; k <= projMonths+TriangleVariableCount; k++ {
			if k <= TriangleVariableCount {
				mtvacp.Field(k).Set(mdiip2.Field(k))
			}
			if k == TriangleVariableCount+1 {
				mtvacp.Field(k).SetFloat(mdiip2.Field(k).Float())
			}

			if k > TriangleVariableCount+1 {
				mtvacp.Field(k).SetFloat(mtvacp.Field(k-1).Float() + mdiip2.Field(k).Float())
			}
		}
		timeValueAdjustedCumulativeProjections = append(timeValueAdjustedCumulativeProjections, timeValueAdjustedCumulativeProjection)
	}

	return timeValueAdjustedCumulativeProjections
}

func getParams(parameterYear int, prodCode, basis string) models.LICParameter {
	var params models.LICParameter

	key := strconv.Itoa(parameterYear) + "_" + prodCode + "_" + basis
	cacheKey := key
	cached, found := IbnrCache.Get(cacheKey)
	if found {
		return cached.(models.LICParameter)
	}
	if found {
		result := cached.(models.LICParameter)
		if result.ProductCode != "" {
			return result
		}
	}
	if !found {
		err := DB.Where("year=? and product_code=? and basis=?", parameterYear, prodCode, basis).Find(&params).Error
		if err != nil {
			fmt.Println(err)
			//	return err
		}
	}
	return params
}
