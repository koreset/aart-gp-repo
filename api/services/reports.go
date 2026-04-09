package services

import (
	"api/models"
	"api/utils"
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func SaveConsolidatedResult(conResult models.ConsolidateResult) error {
	return nil
}

func GetMonthlyConsolidatedResults() []models.ConsolidateResult {
	var conResults []models.ConsolidateResult
	DB.Find(&conResults)
	return conResults
}

func GetAnnualConsolidatedResults() []models.AnnualConsolidatedResult {
	var conResults []models.AnnualConsolidatedResult
	DB.Find(&conResults)
	return conResults
}

type QxValue struct {
	Qx    float64
	AccQx float64
}

type LapseValue struct {
	LapseRate float64 `json:"lapse_rate"`
}

func GetLapseDisclosure(productCode string, currentYear int, pastYear int) []models.LapseDisclosureItem {
	LapseTableName := GetRatingTable(productCode, "Lapse")

	var items []models.LapseDisclosureItem
	tableName := strings.ToLower(productCode) + "_" + strings.ToLower(LapseTableName)
	currentMonthCount := GetLapseCount(currentYear, productCode)
	pastMonthCount := GetLapseCount(pastYear, productCode)
	if currentMonthCount > 0 || pastMonthCount > 0 {
		//Let's determine which value to use
		monthCount := 0
		if currentMonthCount > 0 {
			monthCount = currentMonthCount
		}
		if monthCount == 0 {
			monthCount = pastMonthCount
		}
		for i := 1; i <= monthCount; i++ {
			var currentQx LapseValue
			var pastQx LapseValue
			var item models.LapseDisclosureItem
			item.DurationInForceMonths = i
			err := DB.Table(tableName).Where("year_duration_if_m=?", strconv.Itoa(currentYear)+"_"+strconv.Itoa(i)).First(&currentQx)
			if err != nil {
				fmt.Println(err)
			}
			err = DB.Table(tableName).Where("year_duration_if_m=?", strconv.Itoa(pastYear)+"_"+strconv.Itoa(i)).First(&pastQx)
			if err != nil {
				fmt.Println(err)
			}
			item.CurrentYear = currentQx.LapseRate
			item.PastYear = pastQx.LapseRate
			item.Variance = utils.FloatPrecision(item.CurrentYear-item.PastYear, 2)
			if item.Variance != 0.0 {
				if item.PastYear > 0 {
					item.Change = utils.FloatPrecision((item.CurrentYear-item.PastYear)/item.PastYear, 2)
				} else {
					item.Change = 0
				}
			}

			items = append(items, item)
		}
	}

	return items
}

func GetMortalityDisclosure(productCode string, currentYear int, pastYear int) []models.MortalityDisclosureItem {
	tableName := strings.ToLower(productCode + "_" + GetRatingTable(productCode, "Death"))

	genderFemale := "F"
	genderMale := "M"
	var items []models.MortalityDisclosureItem

	//Lets run for males first
	for i := 1; i <= 120; i++ {
		var currentQx QxValue
		var pastQx QxValue
		var item models.MortalityDisclosureItem
		err := DB.Table(tableName).Where("year_anb_gender=?", strconv.Itoa(currentYear)+"_"+strconv.Itoa(i)+"_"+genderFemale).First(&currentQx).Error
		if err != nil {
			fmt.Println(err)
		}
		err = DB.Table(tableName).Where("year_anb_gender=?", strconv.Itoa(pastYear)+"_"+strconv.Itoa(i)+"_"+genderFemale).First(&pastQx).Error

		item.CurrentYear = currentQx.Qx
		item.PastYear = pastQx.Qx
		item.ANB = i
		item.Gender = genderFemale
		fmt.Println(item)
		items = append(items, item)
	}

	//Now for males
	for i := 1; i <= 120; i++ {
		var currentQx QxValue
		var pastQx QxValue
		var item models.MortalityDisclosureItem
		err := DB.Table(tableName).Where("year_anb_gender=?", strconv.Itoa(currentYear)+"_"+strconv.Itoa(i)+"_"+genderMale).First(&currentQx).Error
		if err != nil {
			fmt.Println(err)
		}
		err = DB.Table(tableName).Where("year_anb_gender=?", strconv.Itoa(pastYear)+"_"+strconv.Itoa(i)+"_"+genderMale).First(&pastQx).Error

		item.CurrentYear = currentQx.Qx
		item.PastYear = pastQx.Qx
		item.ANB = i
		item.Gender = genderMale
		fmt.Println(item)
		items = append(items, item)
	}
	return items
}

func GetMortalityAccidentalDisclosure(productCode string, currentYear int, pastYear int) []models.MortalityDisclosureItem {
	tableName := strings.ToLower(productCode + "_" + GetRatingTable(productCode, "Accidental Death"))
	genderFemale := "F"
	genderMale := "M"
	var items []models.MortalityDisclosureItem

	//Lets run for males first
	for i := 1; i <= 120; i++ {
		var currentQx QxValue
		var pastQx QxValue
		var item models.MortalityDisclosureItem
		err := DB.Table(tableName).Where("year_anb_gender=?", strconv.Itoa(currentYear)+"_"+strconv.Itoa(i)+"_"+genderFemale).First(&currentQx).Error
		if err != nil {
			fmt.Println(err)
		}
		err = DB.Table(tableName).Where("year_anb_gender=?", strconv.Itoa(pastYear)+"_"+strconv.Itoa(i)+"_"+genderFemale).First(&pastQx).Error

		item.CurrentYear = currentQx.AccQx
		item.PastYear = pastQx.AccQx
		item.ANB = i
		item.Gender = genderFemale
		fmt.Println(item)
		items = append(items, item)
	}

	//Now for males
	for i := 1; i <= 120; i++ {
		var currentQx QxValue
		var pastQx QxValue
		var item models.MortalityDisclosureItem
		err := DB.Table(tableName).Where("year_anb_gender=?", strconv.Itoa(currentYear)+"_"+strconv.Itoa(i)+"_"+genderMale).First(&currentQx).Error
		if err != nil {
			fmt.Println(err)
		}
		err = DB.Table(tableName).Where("year_anb_gender=?", strconv.Itoa(pastYear)+"_"+strconv.Itoa(i)+"_"+genderMale).First(&pastQx).Error

		item.CurrentYear = currentQx.AccQx
		item.PastYear = pastQx.AccQx
		item.ANB = i
		item.Gender = genderMale
		fmt.Println(item)
		items = append(items, item)
	}
	return items
}

// GetParameterDisclosure returns a comparison of ProductParameters fields between two years.
// When basis is non-empty the query is scoped to that specific basis; otherwise it picks the
// first row matching product_code + year (backward-compatible with the initial page load).
// Internal fields (json tag "-") are excluded. String/bool fields are flagged with IsNumeric=false.
func GetParameterDisclosure(productCode string, currentYear int, pastYear int, basis string) []models.ParameterDisclosureItem {
	var currentParams models.ProductParameters
	var pastParams models.ProductParameters

	currentQ := DB.Where("product_code=? and year=?", productCode, currentYear)
	pastQ := DB.Where("product_code=? and year=?", productCode, pastYear)
	if basis != "" {
		currentQ = currentQ.Where("basis=?", basis)
		pastQ = pastQ.Where("basis=?", basis)
	}
	currentQ.First(&currentParams)
	pastQ.First(&pastParams)

	t := reflect.TypeOf(currentParams)
	elems1 := reflect.ValueOf(&currentParams).Elem()
	elems2 := reflect.ValueOf(&pastParams).Elem()

	var items []models.ParameterDisclosureItem
	for i := 0; i < elems1.NumField(); i++ {
		fieldType := t.Field(i)

		// Skip internal computed fields and identifier/metadata fields that are not actuarial parameters
		jsonTag := fieldType.Tag.Get("json")
		skipTags := map[string]bool{
			"-": true, "year": true, "product_code": true,
			"basis": true, "created": true, "created_by": true,
		}
		if skipTags[jsonTag] {
			continue
		}

		var item models.ParameterDisclosureItem
		item.Variable = fieldType.Name
		item.CurrentYear = elems1.Field(i).Interface()
		item.PreviousYear = elems2.Field(i).Interface()

		kind := elems1.Field(i).Type().Kind()
		switch kind {
		case reflect.Float64:
			item.IsNumeric = true
			cur := item.CurrentYear.(float64)
			prev := item.PreviousYear.(float64)
			item.Variance = utils.FloatPrecision(cur-prev, 2)
			if item.Variance != 0 && prev != 0 {
				item.Change = utils.FloatPrecision((cur-prev)/prev, 2)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			item.IsNumeric = true
			// Cast to int64 for safe arithmetic across all int kinds
			cur := elems1.Field(i).Int()
			prev := elems2.Field(i).Int()
			item.Variance = utils.FloatPrecision(float64(cur-prev), 2)
			if item.Variance != 0 && prev != 0 {
				item.Change = utils.FloatPrecision(float64(cur-prev)/float64(prev), 2)
			}
		default:
			// String, bool, etc. — display value but skip numeric computations
			item.IsNumeric = false
		}

		items = append(items, item)
	}
	return items
}

// GetParameterMeta returns the distinct bases and years available for a product's parameters.
func GetParameterMeta(productCode string) map[string]interface{} {
	results := make(map[string]interface{})
	var bases []string
	var years []int
	_ = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Model(&models.ProductParameters{}).
			Where("product_code=?", productCode).
			Distinct("basis").
			Pluck("basis", &bases).Error
	})
	_ = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Model(&models.ProductParameters{}).
			Where("product_code=?", productCode).
			Distinct("year").
			Order("year desc").
			Pluck("year", &years).Error
	})
	results["bases"] = bases
	results["years"] = years
	return results
}

// GetMortalityForYears fetches mortality and accidental-mortality disclosure data for two arbitrary years.
func GetMortalityForYears(productCode string, year1 int, year2 int) map[string]interface{} {
	results := make(map[string]interface{})
	results["mortality"] = GetMortalityDisclosure(productCode, year1, year2)
	results["mortality_accidental"] = GetMortalityAccidentalDisclosure(productCode, year1, year2)
	return results
}

// GetAllYieldCurveCodes returns every distinct yield_curve_code in the yield_curve table,
// sorted alphabetically. Used to populate both comparison dropdowns on the Reporting screen.
func GetAllYieldCurveCodes() []string {
	var codes []string
	_ = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw("SELECT DISTINCT yield_curve_code FROM yield_curve ORDER BY yield_curve_code").
			Scan(&codes).Error
	})
	return codes
}

// LatestYieldSnapshot holds the most recent year+month for a given yield_curve_code.
type LatestYieldSnapshot struct {
	MaxYear  int `gorm:"column:max_year"`
	MaxMonth int `gorm:"column:max_month"`
}

// GetYieldCurveDisclosureByCodes compares the nominal rate series of two yield curve codes,
// each taken from their most recent (year, month) upload to avoid mixing vintages.
func GetYieldCurveDisclosureByCodes(code1, code2 string) []models.YieldDisclosureItem {
	// Fetch most-recent snapshot for code1
	var snap1 LatestYieldSnapshot
	_ = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw("SELECT MAX(year) AS max_year, MAX(month) AS max_month FROM yield_curve WHERE yield_curve_code = ?", code1).
			Scan(&snap1).Error
	})

	var series1 []models.YieldCurve
	_ = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("yield_curve_code=? AND year=? AND month=?", code1, snap1.MaxYear, snap1.MaxMonth).
			Order("proj_time asc").
			Find(&series1).Error
	})

	// Fetch most-recent snapshot for code2
	var snap2 LatestYieldSnapshot
	_ = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Raw("SELECT MAX(year) AS max_year, MAX(month) AS max_month FROM yield_curve WHERE yield_curve_code = ?", code2).
			Scan(&snap2).Error
	})

	var series2 []models.YieldCurve
	_ = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("yield_curve_code=? AND year=? AND month=?", code2, snap2.MaxYear, snap2.MaxMonth).
			Order("proj_time asc").
			Find(&series2).Error
	})

	// Build a lookup map for code2 data by projection_time
	code2Map := make(map[int]models.YieldCurve, len(series2))
	for _, row := range series2 {
		code2Map[row.ProjectionTime] = row
	}

	var items []models.YieldDisclosureItem
	for _, row1 := range series1 {
		item := models.YieldDisclosureItem{
			ProjectionTime:       row1.ProjectionTime,
			CurrentYearNominalRate: row1.NominalRate,
			CurrentYearInflation:   row1.Inflation,
		}
		if row2, ok := code2Map[row1.ProjectionTime]; ok {
			item.PastYearNominalRate = row2.NominalRate
			item.PastYearInflation = row2.Inflation
		}
		items = append(items, item)
	}
	return items
}

func GetYieldCurveDisclosure(currentYear int, pastYear int) []models.YieldDisclosureItem {
	var items []models.YieldDisclosureItem
	var currentYield []models.YieldCurve
	var pastYield []models.YieldCurve

	err := DB.Where("year=?", currentYear).Find(&currentYield).Error
	if err != nil {
		fmt.Println(err)
	}
	err = DB.Where("year=?", pastYear).Find(&pastYield).Error
	if err != nil {
		fmt.Println(err)
	}

	//Several scenarios:
	// 1. currentYear and PastYear are both available
	// 2. currentYear available, pastYear not
	// 3. pastYear available, currentYear not available

	if len(currentYield) > 0 && len(pastYield) > 0 {
		for i, _ := range currentYield {
			var item models.YieldDisclosureItem
			item.ProjectionTime = currentYield[i].ProjectionTime
			item.CurrentYearNominalRate = currentYield[i].NominalRate
			item.CurrentYearInflation = currentYield[i].Inflation
			if len(pastYield) > 0 {
				item.PastYearNominalRate = pastYield[i].NominalRate
				item.PastYearInflation = pastYield[i].Inflation
			}
			items = append(items, item)
		}
	}

	if !(len(currentYield) > 0) && len(pastYield) > 0 {
		for i, _ := range pastYield {
			var item models.YieldDisclosureItem
			item.ProjectionTime = pastYield[i].ProjectionTime
			item.PastYearNominalRate = pastYield[i].NominalRate
			item.PastYearInflation = pastYield[i].Inflation
			items = append(items, item)
		}
	}

	if (len(currentYield) > 0) && !(len(pastYield) > 0) {
		for i, _ := range currentYield {
			var item models.YieldDisclosureItem
			item.ProjectionTime = currentYield[i].ProjectionTime
			item.CurrentYearNominalRate = currentYield[i].NominalRate
			item.CurrentYearInflation = currentYield[i].Inflation
			items = append(items, item)
		}
	}

	return items
}

func GetReserveSummaryRunlist() map[string]interface{} {
	var results = make(map[string]interface{})
	var runlist []models.RunList
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Table("aggregated_projections").Distinct("run_date").Find(&runlist).Error
	})
	if err != nil {
		fmt.Println(err)
	}
	results["runlist"] = runlist
	return results
}

func GetPAAReserveSummaryRunlist() map[string]interface{} {
	var results = make(map[string]interface{})
	var runlist []models.RunList
	err := DB.Table("aggregated_modified_gmm_projections").Distinct("run_date").Find(&runlist).Error
	if err != nil {
		fmt.Println(err)
	}
	results["runlist"] = runlist

	return results
}

func GetReserveSummary(runDate string) map[string]interface{} {
	var results = make(map[string]interface{})

	var projections []models.AggregatedProjection
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("run_date=? and projection_month = ?", runDate, 0).Find(&projections).Error
	})
	if err != nil {
		fmt.Println(err)
	}

	// Preload month-1 premiums in a single query to avoid N+1
	var month1Rows []models.Projection
	_ = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("run_date=? and projection_month=?", runDate, 1).Find(&month1Rows).Error
	})
	premiumByProduct := make(map[string]float64)
	for _, r := range month1Rows {
		premiumByProduct[r.ProductCode] = r.Premium * 12.0
	}

	var items []models.ReserveSummary
	for _, projection := range projections {
		var item models.ReserveSummary
		item.BEL = projection.Reserves
		item.BELAdjusted = projection.ReservesAdjusted
		item.ReinsuranceBEL = projection.DiscountedNetReinsuranceCashflow
		item.ReinsuranceBELAdjusted = projection.DiscountedReinsuranceClaimsAdjusted
		item.ProductCode = projection.ProductCode
		item.SpCode = projection.SpCode
		item.ValuationDate = projection.RunDate
		item.SumAssured = projection.SumAssured
		item.UnfundedUnitFund = projection.UnfundedUnitFundEom
		item.BonusStabilisationAccount = projection.BonusStabilisationAccount
		item.AnnuityIncome = projection.AnnuityIncome
		item.PolicyCount = projection.InitialPolicy + projection.NumberOfMaturities
		item.InitialPolicy = projection.InitialPolicy
		item.RunName = projection.RunName
		item.AnnualPremium = premiumByProduct[projection.ProductCode]
		item.Vif = projection.VIF
		item.VifAdjusted = projection.VIFAdjusted
		item.DiscountedCorporateTax = projection.DiscountedCorporateTax
		item.DiscountedCorporateTaxAdjusted = projection.DiscountedCorporateTaxAdjusted
		item.DiscountedPremiumIncome = projection.DiscountedPremiumIncome
		items = append(items, item)
	}

	results["items"] = items

	var bases []string
	_ = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Table("aggregated_projections").Distinct("run_basis").Where("run_date=?", runDate).Find(&bases).Error
	})
	results["bases"] = bases

	return results
}

func GetPAAReserveSummary(runDate string) map[string]interface{} {
	var results = make(map[string]interface{})
	var items []models.PAAReserveSummary
	var projections []models.AggregatedModifiedGMMProjection
	err := DB.Where("run_date=? and projection_month = ? ", runDate, 0).Find(&projections).Error
	if err != nil {
		fmt.Println(err)
	}

	for _, projection := range projections {
		var item models.PAAReserveSummary
		item.ValuationDate = projection.RunDate
		item.RunName = projection.RunName
		item.ProductCode = projection.ProductCode
		item.PortfolioName = projection.PortfolioName
		item.PolicyCount = projection.InforcePolicyCountSM
		item.SumFutureEarnedPremium = projection.SumFutureEarnedPremium
		item.CurrentPeriodEarnedPremium = projection.CurrentPeriodEarnedPremium
		item.CurrentPeriodAmortisedAcquisition = projection.CurrentPeriodAmortisedAcquisition
		item.CurrentPeriodInsuranceRevenue = projection.CurrentPeriodInsuranceRevenue
		item.GMMBel = projection.DiscountedNetCashFlowsLockedin
		item.RiskAdjustment = projection.RiskAdjustment

		items = append(items, item)
	}

	results["items"] = items
	var prodList []string
	DB.Table("aggregated_modified_gmm_projections").Where("run_date = ?", runDate).Distinct("portfolio_name").Find(&prodList)
	results["portfolios"] = prodList
	return results

}

func GetReserveSummaryForPortfolio(runDate string, portfolioName string) map[string]interface{} {
	var results = make(map[string]interface{})
	var items []models.PAAReserveSummary
	var projections []models.AggregatedModifiedGMMProjection
	err := DB.Where("run_date=? and projection_month = ? and portfolio_name = ?", runDate, 0, portfolioName).Find(&projections).Error
	if err != nil {
		fmt.Println(err)
	}

	for _, projection := range projections {
		var item models.PAAReserveSummary
		item.ValuationDate = projection.RunDate
		item.RunName = projection.RunName
		item.ProductCode = projection.ProductCode
		item.PortfolioName = projection.PortfolioName
		item.PolicyCount = projection.InforcePolicyCountSM
		item.SumFutureEarnedPremium = projection.SumFutureEarnedPremium
		item.CurrentPeriodEarnedPremium = projection.CurrentPeriodEarnedPremium
		item.CurrentPeriodAmortisedAcquisition = projection.CurrentPeriodAmortisedAcquisition
		item.CurrentPeriodInsuranceRevenue = projection.CurrentPeriodInsuranceRevenue
		item.GMMBel = projection.DiscountedNetCashFlowsLockedin
		item.RiskAdjustment = projection.RiskAdjustment
		items = append(items, item)
	}

	results["items"] = items
	var prodList []string
	DB.Table("aggregated_modified_gmm_projections").Where("run_date = ? and portfolio_name = ?", runDate, portfolioName).Distinct("product_code").Find(&prodList)
	results["products"] = prodList
	return results
}

func GetReserveSummaryForPortfolioProduct(runDate string, portfolioName string, productCode string) map[string]interface{} {
	var results = make(map[string]interface{})
	var items []models.PAAReserveSummary
	var projections []models.AggregatedModifiedGMMProjection
	err := DB.Where("run_date=? and projection_month = ? and portfolio_name = ? and product_code = ?", runDate, 0, portfolioName, productCode).Find(&projections).Error
	if err != nil {
		fmt.Println(err)
	}

	for _, projection := range projections {
		var item models.PAAReserveSummary
		item.ValuationDate = projection.RunDate
		item.RunName = projection.RunName
		item.ProductCode = projection.ProductCode
		item.PortfolioName = projection.PortfolioName
		item.PolicyCount = projection.InforcePolicyCountSM
		item.SumFutureEarnedPremium = projection.SumFutureEarnedPremium
		item.CurrentPeriodEarnedPremium = projection.CurrentPeriodEarnedPremium
		item.CurrentPeriodAmortisedAcquisition = projection.CurrentPeriodAmortisedAcquisition
		item.CurrentPeriodInsuranceRevenue = projection.CurrentPeriodInsuranceRevenue
		item.GMMBel = projection.DiscountedNetCashFlowsLockedin
		item.RiskAdjustment = projection.RiskAdjustment
		items = append(items, item)
	}

	results["items"] = items

	var groupList []string
	DB.Table("aggregated_modified_gmm_projections").Where("run_date = ? and portfolio_name = ? and product_code = ?", runDate, portfolioName, productCode).Distinct("ifrs17_group").Find(&groupList)
	results["groups"] = groupList

	return results
}

func GetReserveSummaryForPortfolioProductGroup(runDate, portfolioName, productCode, group string) map[string]interface{} {
	var results = make(map[string]interface{})
	var items []models.PAAReserveSummary
	var projections []models.AggregatedModifiedGMMProjection
	err := DB.Where("run_date=? and projection_month = ? and portfolio_name = ? and product_code = ? and ifrs17_group = ?", runDate, 0, portfolioName, productCode, group).Find(&projections).Error
	if err != nil {
		fmt.Println(err)
	}

	for _, projection := range projections {
		var item models.PAAReserveSummary
		item.ValuationDate = projection.RunDate
		item.ProductCode = projection.ProductCode
		item.PortfolioName = projection.PortfolioName
		item.IFRS17Group = projection.IFRS17Group
		item.PolicyCount = projection.InforcePolicyCountSM
		item.SumFutureEarnedPremium = projection.SumFutureEarnedPremium
		item.CurrentPeriodEarnedPremium = projection.CurrentPeriodEarnedPremium
		item.CurrentPeriodAmortisedAcquisition = projection.CurrentPeriodAmortisedAcquisition
		item.CurrentPeriodInsuranceRevenue = projection.CurrentPeriodInsuranceRevenue
		item.GMMBel = projection.DiscountedNetCashFlowsLockedin
		item.RiskAdjustment = projection.RiskAdjustment
		items = append(items, item)
	}

	results["items"] = items
	return results

}

func GetReserveSummaryForProduct(runDate string, productCode string) map[string]interface{} {
	var results = make(map[string]interface{})
	var items = make([]models.ReserveSummary, 0)
	var projections []models.AggregatedProjection
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("run_date=? and product_code=? and projection_month=?", runDate, productCode, 0).Find(&projections).Error
	})
	if err != nil {
		fmt.Println(err)
	}

	// Preload month-1 premiums in a single query to avoid N+1
	var month1Rows []models.Projection
	_ = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("run_date=? and product_code=? and projection_month=?", runDate, productCode, 1).Find(&month1Rows).Error
	})
	premiumByProduct := make(map[string]float64)
	for _, r := range month1Rows {
		premiumByProduct[r.ProductCode] = r.Premium * 12.0
	}

	for _, projection := range projections {
		var item models.ReserveSummary
		item.BEL = projection.Reserves
		item.BELAdjusted = projection.ReservesAdjusted
		item.ReinsuranceBEL = projection.DiscountedNetReinsuranceCashflow
		item.ReinsuranceBELAdjusted = projection.DiscountedReinsuranceClaimsAdjusted
		item.ProductCode = projection.ProductCode
		item.SpCode = projection.SpCode
		item.ValuationDate = projection.RunDate
		item.SumAssured = projection.SumAssured
		item.UnfundedUnitFund = projection.UnfundedUnitFundEom
		item.BonusStabilisationAccount = projection.BonusStabilisationAccount
		item.AnnuityIncome = projection.AnnuityIncome
		item.PolicyCount = projection.InitialPolicy + projection.NumberOfMaturities
		item.InitialPolicy = projection.InitialPolicy
		item.RunName = projection.RunName
		item.AnnualPremium = premiumByProduct[projection.ProductCode]
		item.Vif = projection.VIF
		item.VifAdjusted = projection.VIFAdjusted
		item.DiscountedCorporateTax = projection.DiscountedCorporateTax
		item.DiscountedCorporateTaxAdjusted = projection.DiscountedCorporateTaxAdjusted
		item.DiscountedPremiumIncome = projection.DiscountedPremiumIncome
		items = append(items, item)
	}

	results["items"] = items

	var bases []string
	_ = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Table("aggregated_projections").Where("product_code=? and run_date=?", productCode, runDate).Distinct("run_basis").Find(&bases).Error
	})
	results["bases"] = bases

	return results
}

func GetReserveSummaryForProductAndBasis(runDate string, productCode string, basis string) map[string]interface{} {
	var results = make(map[string]interface{})
	var items []models.ReserveSummary
	var projections []models.AggregatedProjection
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("run_date=? and product_code=? and run_basis=? and projection_month=?", runDate, productCode, basis, 0).Find(&projections).Error
	})
	if err != nil {
		fmt.Println(err)
	}

	// Preload month-1 premiums in a single query to avoid N+1
	var month1Rows []models.Projection
	_ = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("run_date=? and product_code=? and run_basis=? and projection_month=?", runDate, productCode, basis, 1).Find(&month1Rows).Error
	})
	premiumByProduct := make(map[string]float64)
	for _, r := range month1Rows {
		premiumByProduct[r.ProductCode] = r.Premium * 12.0
	}

	for _, projection := range projections {
		var item models.ReserveSummary
		item.BEL = projection.Reserves
		item.BELAdjusted = projection.ReservesAdjusted
		item.ReinsuranceBEL = projection.DiscountedNetReinsuranceCashflow
		item.ReinsuranceBELAdjusted = projection.DiscountedReinsuranceClaimsAdjusted
		item.ProductCode = projection.ProductCode
		item.SpCode = projection.SpCode
		item.ValuationDate = projection.RunDate
		item.SumAssured = projection.SumAssured
		item.UnfundedUnitFund = projection.UnfundedUnitFundEom
		item.BonusStabilisationAccount = projection.BonusStabilisationAccount
		item.AnnuityIncome = projection.AnnuityIncome
		item.PolicyCount = projection.InitialPolicy + projection.NumberOfMaturities
		item.InitialPolicy = projection.InitialPolicy
		item.RunName = projection.RunName
		item.AnnualPremium = premiumByProduct[projection.ProductCode]
		item.Vif = projection.VIF
		item.VifAdjusted = projection.VIFAdjusted
		item.DiscountedCorporateTax = projection.DiscountedCorporateTax
		item.DiscountedCorporateTaxAdjusted = projection.DiscountedCorporateTaxAdjusted
		item.DiscountedPremiumIncome = projection.DiscountedPremiumIncome
		items = append(items, item)
	}

	results["items"] = items

	var spcodes []int
	_ = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Table("aggregated_projections").Where("run_date=? and product_code=? and run_basis=?", runDate, productCode, basis).Distinct("sp_code").Find(&spcodes).Error
	})
	results["spcodes"] = spcodes

	return results
}

func GetReserveSummaryForProductBasisAndSpcode(runDate string, productCode string, spCode int, basis string) map[string]interface{} {
	var results = make(map[string]interface{})
	var items []models.ReserveSummary
	var projections []models.AggregatedProjection
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("run_date=? and product_code=? and sp_code=? and projection_month=? and run_basis=?", runDate, productCode, spCode, 0, basis).Find(&projections).Error
	})
	if err != nil {
		fmt.Println(err)
	}

	// Preload month-1 premiums in a single query to avoid N+1
	var month1Rows []models.Projection
	_ = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("run_date=? and product_code=? and sp_code=? and projection_month=?", runDate, productCode, spCode, 1).Find(&month1Rows).Error
	})
	premiumBySP := make(map[int]float64)
	for _, r := range month1Rows {
		premiumBySP[r.SpCode] = r.Premium * 12.0
	}

	for _, projection := range projections {
		var item models.ReserveSummary
		item.BEL = projection.Reserves
		item.BELAdjusted = projection.ReservesAdjusted
		item.ReinsuranceBEL = projection.DiscountedNetReinsuranceCashflow
		item.ReinsuranceBELAdjusted = projection.DiscountedReinsuranceClaimsAdjusted
		item.ProductCode = projection.ProductCode
		item.SpCode = projection.SpCode
		item.ValuationDate = projection.RunDate
		item.SumAssured = projection.SumAssured
		item.UnfundedUnitFund = projection.UnfundedUnitFundEom
		item.BonusStabilisationAccount = projection.BonusStabilisationAccount
		item.AnnuityIncome = projection.AnnuityIncome
		item.PolicyCount = projection.InitialPolicy + projection.NumberOfMaturities
		item.InitialPolicy = projection.InitialPolicy
		item.RunName = projection.RunName
		item.AnnualPremium = premiumBySP[projection.SpCode]
		item.Vif = projection.VIF
		item.VifAdjusted = projection.VIFAdjusted
		item.DiscountedCorporateTax = projection.DiscountedCorporateTax
		item.DiscountedCorporateTaxAdjusted = projection.DiscountedCorporateTaxAdjusted
		item.DiscountedPremiumIncome = projection.DiscountedPremiumIncome
		items = append(items, item)
	}

	results["items"] = items
	return results
}

// GetReserveSummaryForBasis returns all projection_month=0 rows for a given run_date and run_basis,
// plus a list of distinct product codes available for further filtering.
func GetReserveSummaryForBasis(runDate string, basis string) map[string]interface{} {
	var results = make(map[string]interface{})
	var projections []models.AggregatedProjection
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("run_date=? and run_basis=? and projection_month=?", runDate, basis, 0).Find(&projections).Error
	})
	if err != nil {
		fmt.Println(err)
	}

	var month1Rows []models.Projection
	_ = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("run_date=? and run_basis=? and projection_month=?", runDate, basis, 1).Find(&month1Rows).Error
	})
	premiumByProduct := make(map[string]float64)
	for _, r := range month1Rows {
		premiumByProduct[r.ProductCode] = r.Premium * 12.0
	}

	var items []models.ReserveSummary
	for _, projection := range projections {
		item := mapProjectionToReserveSummary(projection, premiumByProduct[projection.ProductCode])
		items = append(items, item)
	}
	results["items"] = items

	var products []string
	_ = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Table("aggregated_projections").Where("run_date=? and run_basis=?", runDate, basis).Distinct("product_code").Find(&products).Error
	})
	products = append([]string{"All Products"}, products...)
	results["products"] = products

	return results
}

// GetReserveSummaryForBasisAndProduct returns rows filtered by run_date, run_basis, and product_code,
// plus a list of distinct sp_codes for further filtering.
func GetReserveSummaryForBasisAndProduct(runDate string, basis string, productCode string) map[string]interface{} {
	var results = make(map[string]interface{})
	var projections []models.AggregatedProjection
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("run_date=? and run_basis=? and product_code=? and projection_month=?", runDate, basis, productCode, 0).Find(&projections).Error
	})
	if err != nil {
		fmt.Println(err)
	}

	var month1Rows []models.Projection
	_ = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("run_date=? and run_basis=? and product_code=? and projection_month=?", runDate, basis, productCode, 1).Find(&month1Rows).Error
	})
	premiumByProduct := make(map[string]float64)
	for _, r := range month1Rows {
		premiumByProduct[r.ProductCode] = r.Premium * 12.0
	}

	var items []models.ReserveSummary
	for _, projection := range projections {
		item := mapProjectionToReserveSummary(projection, premiumByProduct[projection.ProductCode])
		items = append(items, item)
	}
	results["items"] = items

	var spcodes []int
	_ = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Table("aggregated_projections").Where("run_date=? and run_basis=? and product_code=?", runDate, basis, productCode).Distinct("sp_code").Find(&spcodes).Error
	})
	results["spcodes"] = spcodes

	return results
}

// GetReserveSummaryForBasisProductAndSpcode returns rows filtered by all four dimensions.
func GetReserveSummaryForBasisProductAndSpcode(runDate string, basis string, productCode string, spCode int) map[string]interface{} {
	var results = make(map[string]interface{})
	var projections []models.AggregatedProjection
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("run_date=? and run_basis=? and product_code=? and sp_code=? and projection_month=?", runDate, basis, productCode, spCode, 0).Find(&projections).Error
	})
	if err != nil {
		fmt.Println(err)
	}

	var month1Rows []models.Projection
	_ = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("run_date=? and run_basis=? and product_code=? and sp_code=? and projection_month=?", runDate, basis, productCode, spCode, 1).Find(&month1Rows).Error
	})
	premiumBySP := make(map[int]float64)
	for _, r := range month1Rows {
		premiumBySP[r.SpCode] = r.Premium * 12.0
	}

	var items []models.ReserveSummary
	for _, projection := range projections {
		item := mapProjectionToReserveSummary(projection, premiumBySP[projection.SpCode])
		items = append(items, item)
	}
	results["items"] = items
	return results
}

// mapProjectionToReserveSummary converts an AggregatedProjection row to a ReserveSummary,
// using the pre-fetched annualPremium to avoid per-row DB lookups.
func mapProjectionToReserveSummary(projection models.AggregatedProjection, annualPremium float64) models.ReserveSummary {
	return models.ReserveSummary{
		ValuationDate:                  projection.RunDate,
		ProductCode:                    projection.ProductCode,
		SpCode:                         projection.SpCode,
		RunName:                        projection.RunName,
		PolicyCount:                    projection.InitialPolicy + projection.NumberOfMaturities,
		InitialPolicy:                  projection.InitialPolicy,
		SumAssured:                     projection.SumAssured,
		UnfundedUnitFund:               projection.UnfundedUnitFundEom,
		BonusStabilisationAccount:      projection.BonusStabilisationAccount,
		AnnuityIncome:                  projection.AnnuityIncome,
		AnnualPremium:                  annualPremium,
		BEL:                            projection.Reserves,
		BELAdjusted:                    projection.ReservesAdjusted,
		ReinsuranceBEL:                 projection.DiscountedNetReinsuranceCashflow,
		ReinsuranceBELAdjusted:         projection.DiscountedReinsuranceClaimsAdjusted,
		Vif:                            projection.VIF,
		VifAdjusted:                    projection.VIFAdjusted,
		DiscountedCorporateTax:         projection.DiscountedCorporateTax,
		DiscountedCorporateTaxAdjusted: projection.DiscountedCorporateTaxAdjusted,
		DiscountedPremiumIncome:        projection.DiscountedPremiumIncome,
	}
}
