package services

import (
	"api/log"
	"api/models"
	"api/utils"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	//"github.com/rs/zerolog/log"
)

func GetLapseMargin(month int, year int, prodCode, basis string) float64 {
	var margin models.ProductLapseMargin
	if month > ptnames[prodCode].LapseMarginMonthCount {
		month = ptnames[prodCode].LapseMarginMonthCount
	}

	key := "lp-" + prodCode + "_" + basis + "_" + strconv.Itoa(year) + "-" + strconv.Itoa(month)
	result, found := Cache.Get(key)
	if found {
		return result.(float64)
	} else {
		err := DB.Where("product_code = ? and year = ? and month = ? and basis = ?", prodCode, year, month, basis).Find(&margin).Error
		if err != nil {
			return 0
		}
		Cache.Set(key, margin.Margin, 1)
		return margin.Margin
	}
}

func LoadInflationFactor(year int) {
	var count int64
	DB.Table("yield_curve").Where("year = ? ", year).Count(&count)
	fmt.Println(count)
	for i := 1; i <= int(count); i++ {
		var yieldCurve models.YieldCurve
		key := strconv.Itoa(year) + "-y-curve-" + strconv.Itoa(i)
		err := DB.Table("yield_curve").Where("year = ? and proj_time = ?", year, i).First(&yieldCurve).Error
		if err != nil {
			fmt.Println("Yield Curve: ", errors.WithStack(err))
		} else {
			success := Cache.Set(key, yieldCurve.Inflation, 1)
			// Also set in Redis (TTL 24h) if enabled
			RedisSetFloat(key, yieldCurve.Inflation, 24*time.Hour)
			if !success {
				fmt.Println("cache: key not stored")
			}
		}
	}
}

func LoadForwardRate(year int) {
	var count int64
	DB.Table("yield_curve").Where("year = ? ", year).Count(&count)
	fmt.Println(count)
	for i := 1; i <= int(count); i++ {
		var yieldCurve models.YieldCurve
		key := strconv.Itoa(year) + "-y-curve-fr-" + strconv.Itoa(i)
		err := DB.Table("yield_curve").Where("year = ? and proj_time = ?", year, i).First(&yieldCurve).Error
		if err != nil {
			fmt.Println("Yield Curve: ", errors.WithStack(err))
		} else {
			success := Cache.Set(key, yieldCurve.NominalRate, 1)
			// Also set in Redis (TTL 24h) if enabled
			RedisSetFloat(key, yieldCurve.NominalRate, 24*time.Hour)
			if !success {
				fmt.Println("cache: key not stored")
			}
		}
	}
}

func GetInflationFactor(projectionMonth int, year int, month int, yieldCurveCode string) float64 {
	var yieldCurve models.YieldCurve
	key := strconv.Itoa(year) + "-" + strconv.Itoa(month) + "-y-curve-" + strconv.Itoa(projectionMonth) + "-" + yieldCurveCode

	// 1) Try in-memory cache
	if result, found := Cache.Get(key); found {
		return result.(float64)
	}
	// 2) Try Redis cache and hydrate in-memory cache on hit
	if val, ok := RedisGetFloat(key); ok {
		Cache.Set(key, val, 1)
		return val
	}

	// 3) Load from DB and populate caches
	err := DB.Table("yield_curve").Where("year = ? and proj_time = ? and month = ? and yield_curve_code = ?", year, projectionMonth, month, yieldCurveCode).First(&yieldCurve).Error
	if err != nil {
		fmt.Println("Yield Curve: ", errors.WithStack(err))
	}
	_ = Cache.Set(key, yieldCurve.Inflation, 1)
	// Store in Redis with a sensible TTL (e.g., 24h); if Redis disabled, this is a no-op
	RedisSetFloat(key, yieldCurve.Inflation, 24*time.Hour)
	return yieldCurve.Inflation
}

func GetIbnrInflationFactor(projectionMonth int, year int, month int, yieldCurveCode string) float64 {
	var yieldCurve models.IbnrYieldCurve
	key := strconv.Itoa(year) + "-" + strconv.Itoa(month) + "-y-curve-" + strconv.Itoa(projectionMonth) + "-" + yieldCurveCode

	// 1) Try in-memory cache
	if result, found := Cache.Get(key); found {
		return result.(float64)
	}
	// 2) Try Redis cache and hydrate in-memory cache on hit
	if val, ok := RedisGetFloat(key); ok {
		Cache.Set(key, val, 1)
		return val
	}

	// 3) Load from DB and populate caches
	err := DB.Table("ibnr_yield_curve").Where("year = ? and proj_time = ? and month = ? and yield_curve_code = ?", year, projectionMonth, month, yieldCurveCode).First(&yieldCurve).Error
	if err != nil {
		fmt.Println("Yield Curve: ", errors.WithStack(err))
	}
	_ = Cache.Set(key, yieldCurve.Inflation, 1)
	// Store in Redis with a sensible TTL (e.g., 24h); if Redis disabled, this is a no-op
	RedisSetFloat(key, yieldCurve.Inflation, 24*time.Hour)
	return yieldCurve.Inflation
}

func GetLockedinInflationFactor(projectionMonth int, year int, month int) float64 {
	var yieldCurve models.YieldCurve
	key := strconv.Itoa(year) + "-" + strconv.Itoa(month) + "-y-curve-" + strconv.Itoa(projectionMonth)

	// 1) Try in-memory cache
	if result, found := Cache.Get(key); found {
		return result.(float64)
	}
	// 2) Try Redis cache and hydrate in-memory cache on hit
	if val, ok := RedisGetFloat(key); ok {
		Cache.Set(key, val, 1)
		return val
	}

	// 3) Load from DB and populate caches
	err := DB.Table("yield_curve").Where("year = ? and proj_time = ? and month = ?", year, projectionMonth, month).First(&yieldCurve).Error
	if err != nil {
		fmt.Println("Yield Curve: ", errors.WithStack(err))
	}
	_ = Cache.Set(key, yieldCurve.Inflation, 1)
	// Store in Redis with a sensible TTL (e.g., 24h)
	RedisSetFloat(key, yieldCurve.Inflation, 24*time.Hour)
	return yieldCurve.Inflation
}

// GetInflationFactorWithError is a version of GetInflationFactor that returns an error
func GetInflationFactorWithError(projectionMonth int, year int, month int, yieldCurveCode string) (float64, error) {
	var yieldCurve models.YieldCurve
	key := strconv.Itoa(year) + "-" + strconv.Itoa(month) + "-y-curve-" + strconv.Itoa(projectionMonth) + "-" + yieldCurveCode
	result, found := Cache.Get(key)

	if found {
		return result.(float64), nil
	}

	err := DB.Table("yield_curve").Where("year = ? and proj_time = ? and month = ? and yield_curve_code = ?", year, projectionMonth, month, yieldCurveCode).First(&yieldCurve).Error
	if err != nil {
		return 0, errors.Wrap(err, "failed to get inflation factor")
	}
	success := Cache.Set(key, yieldCurve.Inflation, 1)
	if !success {
		log.WithField("year", year).WithField("yield_curve_code", yieldCurveCode).Error("cache: key not stored")
	} else {
		log.WithField("year", year).WithField("yield_curve_code", yieldCurveCode).Info(("cache: key stored"))
	}

	return yieldCurve.Inflation, nil
}

// GetLockedinInflationFactorWithError is a version of GetInflationFactor that returns an error
func GetLockedinInflationFactorWithError(projectionMonth int, year int, month int) (float64, error) {
	var yieldCurve models.YieldCurve
	key := strconv.Itoa(year) + "-" + strconv.Itoa(month) + "-y-curve-" + strconv.Itoa(projectionMonth)
	result, found := Cache.Get(key)

	if found {
		return result.(float64), nil
	}

	err := DB.Table("yield_curve").Where("year = ? and proj_time = ? and month = ?", year, projectionMonth, month).First(&yieldCurve).Error
	if err != nil {
		return 0, errors.Wrap(err, "failed to get inflation factor")
	}
	success := Cache.Set(key, yieldCurve.Inflation, 1)
	if !success {
		log.Error("cache: key not stored")
	}

	return yieldCurve.Inflation, nil
}

//func GetForwardRate(projectionMonth int, year int, month int, yieldCurveCode string, basis string) float64 {
//	var yieldCurve models.YieldCurve
//	key := strconv.Itoa(year) + "-" + strconv.Itoa(month) + "-y-curve-fr-" + strconv.Itoa(projectionMonth) + yieldCurveCode + basis
//	cached, found := Cache.Get(key)
//
//	if found {
//		result := cached.(float64)
//		//if result > 0 {
//		return result
//		//}
//	}
//	if projectionMonth == 0 {
//		return 0
//	}
//	err := DB.Table("yield_curve").Where("year = ? and proj_time = ? and month = ? and yield_curve_code=? and basis=?", year, projectionMonth, month, yieldCurveCode, basis).First(&yieldCurve).Error
//	if err != nil {
//		fmt.Println("forward rate db error: ", errors.WithStack(err))
//	}
//	success := Cache.Set(key, yieldCurve.NominalRate, 1)
//	if !success {
//		fmt.Println("cache error: Not saved")
//	}
//
//	return yieldCurve.NominalRate
//}

// GetForwardRateWithError is a version of GetForwardRate that returns an error
func GetForwardRateWithError(projectionMonth int, year int, month int, yieldCurveCode string) (float64, error) {
	var yieldCurve models.YieldCurve
	key := strconv.Itoa(year) + "-" + strconv.Itoa(month) + "-y-curve-fr-" + strconv.Itoa(projectionMonth) + yieldCurveCode
	cached, found := Cache.Get(key)

	if found {
		result := cached.(float64)
		return result, nil
	}
	if projectionMonth == 0 {
		return 0, nil
	}
	err := DB.Table("yield_curve").Where("year = ? and proj_time = ? and month = ? and yield_curve_code=?", year, projectionMonth, month, yieldCurveCode).First(&yieldCurve).Error
	if err != nil {
		return 0, errors.Wrap(err, "forward rate db error")
	}
	success := Cache.Set(key, yieldCurve.NominalRate, 1)
	if !success {
		log.Error("cache error: Not saved")
	}

	return yieldCurve.NominalRate, nil
}

// GetForwardRateWithError is a version of GetForwardRate that returns an error
func GetIbnrForwardRateWithError(projectionMonth int, year int, month int, yieldCurveCode string) (float64, error) {
	var yieldCurve models.IbnrYieldCurve
	key := strconv.Itoa(year) + "-" + strconv.Itoa(month) + "-y-curve-fr-" + strconv.Itoa(projectionMonth) + yieldCurveCode
	cached, found := Cache.Get(key)

	if found {
		result := cached.(float64)
		return result, nil
	}
	if projectionMonth == 0 {
		return 0, nil
	}
	err := DB.Table("ibnr_yield_curve").Where("year = ? and proj_time = ? and month = ? and yield_curve_code=?", year, projectionMonth, month, yieldCurveCode).First(&yieldCurve).Error
	if err != nil {
		return 0, errors.Wrap(err, "forward rate db error")
	}
	success := Cache.Set(key, yieldCurve.NominalRate, 1)
	if !success {
		log.Error("cache error: Not saved")
	}

	return yieldCurve.NominalRate, nil
}

// GetLockedinForwardRateWithError is a version of GetForwardRate that returns an error
func GetLockedinForwardRateWithError(projectionMonth int, year int, month int) (float64, error) {
	var yieldCurve models.YieldCurve
	key := strconv.Itoa(year) + "-" + strconv.Itoa(month) + "-y-curve-fr-" + strconv.Itoa(projectionMonth)

	cached, found := Cache.Get(key)

	if found {
		result := cached.(float64)
		return result, nil
	}
	if projectionMonth == 0 {
		return 0, nil
	}
	err := DB.Table("yield_curve").Where("year = ? and proj_time = ? and month = ?", year, projectionMonth, month).First(&yieldCurve).Error
	if err != nil {
		return 0, errors.Wrap(err, "forward rate db error")
	}
	success := Cache.Set(key, yieldCurve.NominalRate, 1)
	if !success {
		log.Error("cache error: Not saved")
	}

	return yieldCurve.NominalRate, nil
}

func GetPaaForwardRate(projectionMonth int, year int, month int, yieldCurveCode string) float64 {
	var yieldCurve models.PaaYieldCurve
	key := strconv.Itoa(year) + "-" + strconv.Itoa(month) + "-y-curve-fr-" + strconv.Itoa(projectionMonth) + "-" + yieldCurveCode
	cached, found := PaaCache.Get(key)

	if found {
		result := cached.(float64)
		//if result > 0 {
		return result
		//}
	}
	//if projectionMonth == 0 {
	//	return 0
	//}
	err := DB.Table("paa_yield_curve").Where("year = ? and proj_time = ? and month = ? and yield_curve_code = ?", year, projectionMonth, month, yieldCurveCode).First(&yieldCurve).Error
	if err != nil {
		//fmt.Println("forward rate db error: ", errors.WithStack(err))
	}
	success := PaaCache.Set(key, yieldCurve.NominalRate, 1)
	if !success {
		fmt.Println("cache error: Not saved")
	}

	return yieldCurve.NominalRate
}

// GetMortalityRate factors used to read a mortality table
// Age, Gender, Smoking_Status,Select_Period,Income,Education,SEC, OCC_Class
func GetMortalityRate(arg models.TransitionRateArguments) float64 {
	tablename := strings.ToLower(arg.ProductCode) + "_" + strings.ToLower(ptnames[arg.ProductCode].MortalityTableName)

	key := keyBuilder(arg, "Death")
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

	var qx float64
	if ptnames[arg.ProductCode].MortalityColumnName != "" {
		query := ptnames[arg.ProductCode].MortalityColumnName + " = ?"
		row := DB.Table(tablename).Where(query, key).Select("qx").Row()
		row.Scan(&qx)
		Cache.Set(cacheKey, qx, 1)
		//time.Sleep(5 * time.Millisecond)
	}
	return qx

}

func GetBenefitMultiplier(mp models.ProductModelPoint) float64 {
	tablename := "product_benefit_multiplier"

	key := "annuity_income_multiplier" + "_" + mp.ProductCode + "_" + mp.Plan
	//cacheKey := tablename + "_" + key
	cached, found := Cache.Get(key)

	if found {
		result := cached.(float64)
		//if result > 0 {
		return result
		//}
	} else {
		//fmt.Println("cache missed: ", key)
	}

	var annuityIncomeMultiplier float64
	row := DB.Table(tablename).Where("product_code=? and plan=?", mp.ProductCode, mp.Plan).Select("annuity_income_multiplier").Row()
	row.Scan(&annuityIncomeMultiplier)
	Cache.Set(key, annuityIncomeMultiplier, 1)
	return annuityIncomeMultiplier
}

// GetBenefitMultiplierWithError is a version of GetBenefitMultiplier that returns an error
func GetBenefitMultiplierWithError(mp models.ProductModelPoint) (float64, error) {
	tablename := "product_benefit_multiplier"

	key := "annuity_income_multiplier" + "_" + mp.ProductCode + "_" + mp.Plan
	//cacheKey := tablename + "_" + key
	cached, found := Cache.Get(key)

	if found {
		result := cached.(float64)
		return result, nil
	}

	var annuityIncomeMultiplier float64
	row := DB.Table(tablename).Where("product_code=? and plan=?", mp.ProductCode, mp.Plan).Select("annuity_income_multiplier").Row()
	err := row.Scan(&annuityIncomeMultiplier)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get benefit multiplier")
	}
	Cache.Set(key, annuityIncomeMultiplier, 1)
	return annuityIncomeMultiplier, nil
}

func LoadMortalityRates(productCode string, year int) {
	//gender := []string{"M", "F"}
	if ptnames[productCode].MortalityTableName != "" {
		tablename := strings.ToLower(productCode) + "_" + strings.ToLower(ptnames[productCode].MortalityTableName)

		columnName := ptnames[productCode].MortalityColumnName
		query := "left(" + columnName + ",4) = ? "
		var qx float64
		var cname string
		key := columnName
		cacheKey := tablename + "_" + key
		//Select(fmt.Sprintf("%s as cname, qx",columnName))
		rows, err := DB.Table(tablename).Where(query, strconv.Itoa(year)).Select(fmt.Sprintf("%s as cname, qx", columnName)).Rows()
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&cname, &qx)
			if err != nil {
				fmt.Println(err)
			}
			Cache.Set(cacheKey, qx, 1)
		}

		//sanity check
		cached, found := Cache.Get(tablename + "_" + strconv.Itoa(year) + "_3_F")
		if found {
			fmt.Println("load key: ", cached)
		} else {
			fmt.Println("NOT FOUND")
		}
	}
}

func LoadAccidentalMortalityRates(productCode string, year int) {
	if ptnames[productCode].MortalityAccidentalTableName != "" {
		//gender := []string{"M", "F"}
		tablename := strings.ToLower(productCode) + "_" + strings.ToLower(ptnames[productCode].MortalityAccidentalTableName)

		columnName := ptnames[productCode].MortalityAccidentalColumnName
		query := "left(" + columnName + ",4) = ? "
		var acc_qx_prop float64
		var cname string
		key := columnName
		cacheKey := tablename + "_" + key
		//Select(fmt.Sprintf("%s as cname, qx",columnName))
		rows, err := DB.Table(tablename).Where(query, strconv.Itoa(year)).Select(fmt.Sprintf("%s as cname, acc_qx_prop", columnName)).Rows()
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&cname, &acc_qx_prop)
			if err != nil {
				fmt.Println(err)
			}
			Cache.Set(cacheKey, acc_qx_prop, 1)
		}

		//sanity check
		cached, found := Cache.Get(tablename + "_" + strconv.Itoa(year) + "_3_F")
		if found {
			fmt.Println("load key: ", cached)
		} else {
			fmt.Println("NOT FOUND")
		}
	}
}

func GetMortalityRateAccidentProportion(arg models.TransitionRateArguments) float64 {
	key := keyBuilder(arg, "Accidental Death")
	//key := strconv.Itoa(year) + "_" + strconv.Itoa(age) + "_" + gender[:1]
	tableName := strings.ToLower(arg.ProductCode) + "_" + strings.ToLower(ptnames[arg.ProductCode].MortalityAccidentalTableName)
	cacheKey := strings.ToLower(tableName) + "_" + key
	cached, found := Cache.Get(cacheKey)

	if found {
		result := cached.(float64)
		//if result > 0 {
		return result
		//}
	}

	var accQxProp float64

	if ptnames[arg.ProductCode].MortalityAccidentalColumnName != "" {
		query := ptnames[arg.ProductCode].MortalityAccidentalColumnName + " = ?"
		row := DB.Table(tableName).Where(query, key).Select("acc_qx_prop").Row()
		row.Scan(&accQxProp)
		Cache.Set(cacheKey, accQxProp, 1)
	}
	return accQxProp
}

func LoadDisabilityRates(productCode string, year int) {
	if ptnames[productCode].DisabilityTableName != "" {
		//gender := []string{"M", "F"}
		tablename := strings.ToLower(productCode) + "_" + strings.ToLower(ptnames[productCode].DisabilityTableName)

		columnName := ptnames[productCode].DisabilityColumnName
		query := "left(" + columnName + ",4) = ? "
		var incidence_rate float64
		var cname string
		//key := columnName
		//cacheKey := tablename + "_" + cname
		//Select(fmt.Sprintf("%s as cname, qx",columnName))
		rows, err := DB.Table(tablename).Where(query, strconv.Itoa(year)).Select(fmt.Sprintf("%s as cname, incidence_rate", columnName)).Rows()
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&cname, &incidence_rate)
			if err != nil {
				fmt.Println(err)
			}
			Cache.Set(tablename+"_"+cname, incidence_rate, 1)
		}
	}
}

func LoadRetrenchmentRates(productCode string, year int) {
	if ptnames[productCode].RetrenchmentTableName != "" {
		tablename := strings.ToLower(productCode) + "_" + strings.ToLower(ptnames[productCode].RetrenchmentTableName)

		columnName := ptnames[productCode].RetrenchmentColumnName
		query := "left(" + columnName + ",4) = ? "
		var retr_rate float64
		var cname string
		key := columnName
		cacheKey := tablename + "_" + key
		rows, err := DB.Table(tablename).Where(query, strconv.Itoa(year)).Select(fmt.Sprintf("%s as cname, retr_rate", columnName)).Rows()
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&cname, &retr_rate)
			if err != nil {
				fmt.Println(err)
			}
			Cache.Set(cacheKey, retr_rate, 1)
		}
	}
}

func LoadLapseRates(productCode string, year int) {
	if ptnames[productCode].LapseTableName != "" {
		tablename := strings.ToLower(productCode) + "_" + strings.ToLower(ptnames[productCode].LapseTableName)

		columnName := ptnames[productCode].LapseColumnName
		query := "left(" + columnName + ",4) = ? "
		var lapse_rate float64
		var cname string
		key := columnName
		cacheKey := tablename + "_" + key
		rows, err := DB.Table(tablename).Where(query, strconv.Itoa(year)).Select(fmt.Sprintf("%s as cname, lapse_rate", columnName)).Rows()
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&cname, &lapse_rate)
			if err != nil {
				fmt.Println(err)
			}
			Cache.Set(cacheKey, lapse_rate, 1)
		}
	}
}

func GetLapseRate(dfs models.TransitionRateArguments) float64 {
	//if dfs.DurationIfM > LapseMaxMonthDimension {
	//	dfs.DurationIfM = LapseMaxMonthDimension
	//}
	if dfs.DurationIfM > ptnames[dfs.ProductCode].LapseMonthCount {
		dfs.DurationIfM = ptnames[dfs.ProductCode].LapseMonthCount
	}

	key := keyBuilder(dfs, "Lapse")
	tableName := strings.ToLower(dfs.ProductCode) + "_" + strings.ToLower(ptnames[dfs.ProductCode].LapseTableName)

	//if month == 0 {
	//	month = 1
	//}
	//
	//if month > ptnames[prodCode].LapseMonthCount {
	//	month = ptnames[prodCode].LapseMonthCount
	//}
	cacheKey := tableName + "_" + key

	cached, found := Cache.Get(cacheKey)

	if found {
		result := cached.(float64)
		//if result > 0 {
		return result
		//}
	}

	var lapseRate float64

	if ptnames[dfs.ProductCode].LapseColumnName != "" {
		query := ptnames[dfs.ProductCode].LapseColumnName + " = ?"
		row := DB.Table(tableName).Where(query, key).Select("lapse_rate").Row()
		row.Scan(&lapseRate)
		Cache.Set(cacheKey, lapseRate, 1)
	}
	return lapseRate
}

func GetSpecialMargins(anb int, projection *models.Projection, prodCode, memberType string, run models.RunParameters, parameters models.ProductParameters, decrement string) float64 {
	//anb := projection.AgeNextBirthday

	key := strings.ToLower(prodCode) + "_" + strconv.Itoa(anb) + "_" + strings.ToLower(memberType) + parameters.SpecialMarginCode + "_" + run.RunBasis + decrement
	tableName := strings.ToLower("product") + "_" + strings.ToLower("Special_Decrement_Margins")

	cacheKey := tableName + "_" + key

	cached, found := Cache.Get(cacheKey)

	if found {
		result := cached.(float64)
		//if result != 0 {
		return result
		//}
	}

	var specialMargins models.ProductSpecialDecrementMargin
	var specialMargin float64
	err := DB.Where("product_code =? and anb=? and member_type=? and basis=? and special_margin_code=?", prodCode, anb, memberType, run.RunBasis, parameters.SpecialMarginCode).Find(&specialMargins).Error
	if err != nil {
		fmt.Println("specialMargins: ", err)
		switch decrement {
		case "mortality_margin":
			Cache.Set(cacheKey, 0, 1)
			return 0
		case "morbidity_margin":
			Cache.Set(cacheKey, 0, 1)
			return 0
		case "retrenchment_margin":
			Cache.Set(cacheKey, 0, 1)
			return 0
		case "mortality_table_prop":
			Cache.Set(cacheKey, 1, 1)
			return 1
		case "disability_table_prop":
			Cache.Set(cacheKey, 1, 1)
			return 1
		case "lapse_table_prop":
			Cache.Set(cacheKey, 1, 1)
			return 1
		default:
			Cache.Set(cacheKey, 0, 1)
			return 0
		}
		//return 0
	}
	switch decrement {
	case "mortality_margin":
		specialMargin = specialMargins.MortalityMargin
	case "morbidity_margin":
		specialMargin = specialMargins.MorbidityMargin
	case "retrenchment_margin":
		specialMargin = specialMargins.RetrenchmentMargin
	case "mortality_table_prop":
		specialMargin = specialMargins.MortalityTableProp
	case "disability_table_prop":
		specialMargin = specialMargins.DisabilityTableProp
	case "lapse_table_prop":
		specialMargin = specialMargins.LapseTableProp
	default:
		specialMargin = 0
	}

	Cache.Set(cacheKey, specialMargin, 1)
	return specialMargin
}

func GetProfitRenewalAdjustment(projection *models.Projection, mp models.ProductModelPoint, parameters models.ProductParameters) float64 {
	durationYear := int(projection.ValuationTimeYear) + 1
	key := strings.ToLower(mp.ProductCode) + "_" + strconv.Itoa(durationYear) + parameters.ProfitAdjustmentCode
	tableName := strings.ToLower("product") + "_" + strings.ToLower("renewable_profit_adjustments")

	cacheKey := tableName + "_" + key

	cached, found := Cache.Get(cacheKey)

	if found {
		result := cached.(float64)
		//if result != 0 {
		return result
		//}
	}

	var profitAdjustments models.ProductRenewableProfitAdjustment
	var profitAdjustment float64
	err := DB.Where("product_code =? and duration_if_y = ? and profit_adjustment_code = ?", mp.ProductCode, durationYear, parameters.ProfitAdjustmentCode).Find(&profitAdjustments).Error
	if err != nil {
		fmt.Println("profit adjustment: ", err)
		return 0
	}

	profitAdjustment = profitAdjustments.RenewableProfitAdjustment

	Cache.Set(cacheKey, profitAdjustment, 1)
	return profitAdjustment
}

func GetLapseMarginCount(year int, code, basis string) int {
	var count int64
	var err error
	if basis == "" {
		err = DB.Table("product_lapse_margins").Where("year = ? and product_code = ?", year, code).Count(&count).Error

	}
	if basis != "" {
		err = DB.Table("product_lapse_margins").Where("year = ? and product_code = ? and basis = ?", year, code, basis).Count(&count).Error

	}
	if err != nil {
		log.Errorf("GetLapseMarginCount: %s", err.Error())
		return 0
	}
	return int(count)
}

func GetRetrenchmentCount(year int, RetrenchmentTable, RetrenchmentColumnName string) int {
	var count int64
	query := fmt.Sprintf("select count(*) from %s where %s like '%%%d%%'", RetrenchmentTable, RetrenchmentColumnName, year)
	err := DB.Raw(query).Count(&count).Error
	if err != nil {
		log.Error("GetRetrenchmentCount: ", err.Error())
		return 0
	}
	return int(count)
}

func GetLapseCount(year int, code string) int {
	var count int64
	var tableName = strings.ToLower(code) + "_" + strings.ToLower(ptnames[code].LapseTableName)

	// Convert the integer year to a string prefix (e.g., "2025_")
	yearPrefix := fmt.Sprintf("%d_%%", year)

	if ptnames[code].LapseColumnName != "" {
		query := ptnames[code].LapseColumnName + " LIKE ?"
		err := DB.Table(tableName).Where(query, yearPrefix).Count(&count).Error
		if err != nil {
			log.Error("GetLapseCount: ", err.Error())
			return 0
		}
	}
	return int(count)
}

func GetRatingTable(code, class string) string {
	prod, err := GetProductByCode(code)
	if err != nil {
		log.Error("Product not found", code)
	}
	table := ""
	for _, trans := range prod.ProductTransitions {
		if trans.EndState == class {
			table = trans.AssociatedTable
			break
		}
	}
	return table
}

func GetCommissionStructure(mp models.ProductModelPoint) models.ProductCommissionStructure {
	var commStructure models.ProductCommissionStructure
	key := fmt.Sprintf("commissiontype%s_productcode_%s", mp.CommissionType, mp.ProductCode)
	cached, found := Cache.Get(key)
	if found {
		commStructure = cached.(models.ProductCommissionStructure)
	} else {
		err := DB.Where("commission_type = ? and product_code = ?", mp.CommissionType, mp.ProductCode).First(&commStructure).Error
		if err == nil {
			Cache.Set(key, commStructure, 1)
			//time.Sleep(2 * time.Millisecond)
		}
		if err != nil {
			Cache.Set(key, commStructure, 1)
		}
	}
	return commStructure
}

func GetClawback(month int, mp models.ProductModelPoint) models.ProductClawback {
	var clawback models.ProductClawback
	key := fmt.Sprintf("clawback_%d_productcode_%s", month, mp.ProductCode)
	cached, found := Cache.Get(key)
	if found {
		clawback = cached.(models.ProductClawback)
	} else {
		err := DB.Where("duration_in_force_month = ? and product_code = ?", month, mp.ProductCode).First(&clawback).Error
		if err == nil {
			Cache.Set(key, clawback, 1)
			//time.Sleep(2 * time.Millisecond)
		}
	}
	return clawback
}

//func GetClawback(month int) models.ProductClawback {
//	var clawback models.ProductClawback
//	DB.Where("duration_in_force_month = ?", month).First(&clawback)
//	return clawback
//}

func GetChildAdditionalSumAssured(productCode string, plan string, age int) float64 {
	key := fmt.Sprintf("%s-%s-%d", strings.ToLower(productCode), plan, age)
	cached, found := Cache.Get(key)
	if found {
		return cached.(float64)
	}

	var childAdditionalSumassured models.ProductChildAdditionalSumAssured
	err := DB.Where("product_code =? and age=?", productCode, age).Find(&childAdditionalSumassured).Error
	if err != nil {
		fmt.Println("child funeral: ", err)
		return 0
	}
	var result float64
	switch plan {
	case "A":
		result = childAdditionalSumassured.A
	case "B":
		result = childAdditionalSumassured.B
	case "C":
		result = childAdditionalSumassured.C
	case "D":
		result = childAdditionalSumassured.D
	case "E":
		result = childAdditionalSumassured.E
	case "F":
		result = childAdditionalSumassured.F
	case "G":
		result = childAdditionalSumassured.G
	case "H":
		result = childAdditionalSumassured.H
	case "I":
		result = childAdditionalSumassured.I
	case "J":
		result = childAdditionalSumassured.J
	case "K":
		result = childAdditionalSumassured.K
	case "L":
		result = childAdditionalSumassured.L
	case "M":
		result = childAdditionalSumassured.M
	case "N":
		result = childAdditionalSumassured.N
	case "O":
		result = childAdditionalSumassured.O
	default:
		result = 0
	}
	Cache.Set(key, result, 1)

	return result
}

func GetChildSumAssured(productCode string, plan string, age int) float64 {
	key := fmt.Sprintf("%s-%s-%d", strings.ToLower(productCode), plan, age)
	cached, found := Cache.Get(key)
	if found {
		return cached.(float64)
	}

	var childSumAssured models.ProductChildSumAssured
	err := DB.Where("product_code =? and age=?", productCode, age).Find(&childSumAssured).Error
	if err != nil {
		fmt.Println("child sum assured: ", err)
		return 0
	}
	var result float64
	switch plan {
	case "A":
		result = childSumAssured.A
	case "B":
		result = childSumAssured.B
	case "C":
		result = childSumAssured.C
	case "D":
		result = childSumAssured.D
	case "E":
		result = childSumAssured.E
	case "F":
		result = childSumAssured.F
	case "G":
		result = childSumAssured.G
	case "H":
		result = childSumAssured.H
	case "I":
		result = childSumAssured.I
	case "J":
		result = childSumAssured.J
	case "K":
		result = childSumAssured.K
	case "L":
		result = childSumAssured.L
	case "M":
		result = childSumAssured.M
	case "N":
		result = childSumAssured.N
	case "O":
		result = childSumAssured.O
	default:
		result = 0
	}
	Cache.Set(key, result, 1)

	return result
}

func GetAdditionalSumAssured(productCode string, plan string) float64 {
	key := fmt.Sprintf("%s-%s", strings.ToLower(productCode), plan)
	cached, found := Cache.Get(key)
	if found {
		return cached.(float64)
	}

	var fs models.ProductAdditionalSumAssured
	err := DB.Where("product_code =?", productCode).Find(&fs).Error
	if err != nil {
		fmt.Println("funeral: ", err)
		return 0
	}
	var result float64
	switch plan {
	case "A":
		result = fs.A
	case "B":
		result = fs.B
	case "C":
		result = fs.C
	case "D":
		result = fs.D
	case "E":
		result = fs.E
	case "F":
		result = fs.F
	case "G":
		result = fs.G
	case "H":
		result = fs.H
	case "I":
		result = fs.I
	case "J":
		result = fs.J
	case "K":
		result = fs.K
	case "L":
		result = fs.L
	case "M":
		result = fs.M
	case "N":
		result = fs.N
	case "O":
		result = fs.O
	default:
		result = 0
	}
	Cache.Set(key, result, 1)

	return result
}

func GetColumnName(tableName string) string {
	rows, err := DB.Table(tableName).Exec(fmt.Sprintf("describe %s", tableName)).Rows()
	if err != nil {
		log.WithField("error", err).Error("GetColumnName: ", err)
		return ""
	}
	if rows == nil {
		return ""
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil || len(cols) == 0 {
		if err != nil {
			log.WithField("error", err).Error("GetColumnName Columns: ", err)
		}
		return ""
	}
	return cols[0]
}

func GetDisabilityIncidenceRate(dfs models.TransitionRateArguments) float64 {
	//Build the key
	key := keyBuilder(dfs, "Permanent Disability")

	tableName := strings.ToLower(dfs.ProductCode) + "_" + strings.ToLower(ptnames[dfs.ProductCode].DisabilityTableName)
	cacheKey := tableName + "_" + key
	cached, found := Cache.Get(cacheKey)

	if found {
		result := cached.(float64)
		//if result > 0 {
		return result
		//}
	}

	var incidenceRate float64

	if ptnames[dfs.ProductCode].DisabilityColumnName != "" {
		query := ptnames[dfs.ProductCode].DisabilityColumnName + " = ?"
		row := DB.Table(tableName).Where(query, key).Select("incidence_rate").Row()
		row.Scan(&incidenceRate)
		Cache.Set(cacheKey, incidenceRate, 1)
	}
	return incidenceRate
}

func GetRetrenchmentRate(dfs models.TransitionRateArguments) float64 {
	//Build the key
	//if dfs.ProjectionMonth > RetrMaxYearDimension {
	//	dfs.ProjectionMonth = RetrMaxYearDimension
	//}
	if dfs.ProjectionMonth > ptnames[dfs.ProductCode].RetrenchmentRowCount {
		dfs.ProjectionMonth = ptnames[dfs.ProductCode].RetrenchmentRowCount
	}

	key := keyBuilder(dfs, "Retrenchment")

	tableName := strings.ToLower(dfs.ProductCode) + "_" + strings.ToLower(ptnames[dfs.ProductCode].RetrenchmentTableName)
	cacheKey := tableName + "_" + key
	cached, found := Cache.Get(cacheKey)

	if found {
		result := cached.(float64)
		//if result > 0 {
		return result
		//}
	}
	var retrRate float64
	if ptnames[dfs.ProductCode].RetrenchmentColumnName != "" {
		query := ptnames[dfs.ProductCode].RetrenchmentColumnName + " = ?"
		row := DB.Table(tableName).Where(query, key).Select("retr_rate").Row()
		row.Scan(&retrRate)
		Cache.Set(cacheKey, retrRate, 1)
	}
	return retrRate
}

func keyBuilder(dfs models.TransitionRateArguments, state string) string {
	var trans models.ProductTransition
	var err error

	if dfs.Age < 1 || dfs.Gender == "" {
		//something is off... return
		return ""
	}

	associatedTableKey := strconv.Itoa(dfs.ProductId) + utils.Snakify(state)
	associatedTable, found := Cache.Get(associatedTableKey)
	if !found {
		err = DB.Where("product_id = ? and end_state =?", dfs.ProductId, state).First(&trans).Error
		if err != nil {
			log.Error(err.Error())
		} else {
			associatedTable = trans.AssociatedTable
			Cache.Set(associatedTableKey, associatedTable, 1)
		}
	}

	var ratingFactor models.ProductRatingFactor

	ratingFactorKey := "rf_" + strconv.Itoa(dfs.ProductId) + "_" + associatedTable.(string)

	cached, found := Cache.Get(ratingFactorKey)

	if !found {
		err = DB.Preload("Fds").Where("product_id = ? and transition_table =?", dfs.ProductId, associatedTable.(string)).First(&ratingFactor).Error
		if err != nil {
			log.Error(err.Error())
		} else {
			Cache.Set(ratingFactorKey, ratingFactor, 1)
		}
	} else {
		ratingFactor = cached.(models.ProductRatingFactor)
	}

	var keyString strings.Builder
	keyString.WriteString(strconv.Itoa(dfs.Year) + "_")

	if utils.FactorsContains(&ratingFactor.Fds, "ANB") {
		keyString.WriteString(strconv.Itoa(dfs.Age) + "_")
	}

	if utils.FactorsContains(&ratingFactor.Fds, "GENDER") {
		keyString.WriteString(dfs.Gender[:1] + "_")
	}

	if utils.FactorsContains(&ratingFactor.Fds, "SMOKER_STATUS") {
		keyString.WriteString(dfs.SmokerStatus + "_")
	}

	if utils.FactorsContains(&ratingFactor.Fds, "INCOME") {
		keyString.WriteString(strconv.Itoa(dfs.Income))
		keyString.WriteString("_")
	}

	if utils.FactorsContains(&ratingFactor.Fds, "SEC") {
		keyString.WriteString(strconv.Itoa(dfs.SocioEconomicClass))
		keyString.WriteString("_")
	}

	if utils.FactorsContains(&ratingFactor.Fds, "OCC_CLASS") {
		keyString.WriteString(dfs.OccupationalClass)
		keyString.WriteString("_")
	}

	if utils.FactorsContains(&ratingFactor.Fds, "SELECT_PERIOD") {
		keyString.WriteString(strconv.Itoa(dfs.SelectPeriod))
		keyString.WriteString("_")
	}

	if utils.FactorsContains(&ratingFactor.Fds, "EDUCATION_LEVEL") {
		keyString.WriteString(strconv.Itoa(dfs.EducationLevel))
		keyString.WriteString("_")
	}

	if utils.FactorsContains(&ratingFactor.Fds, "DURATION_IF_M") {
		keyString.WriteString(strconv.Itoa(dfs.DurationIfM))
		keyString.WriteString("_")
	}

	if utils.FactorsContains(&ratingFactor.Fds, "PROJECTION_MONTH") {
		keyString.WriteString(strconv.Itoa(int(math.Min(float64(dfs.ProjectionMonth), 60.0))))
		keyString.WriteString("_")
	}

	if utils.FactorsContains(&ratingFactor.Fds, "DISTRIBUTION_CHANNEL") {
		keyString.WriteString(dfs.DistributionChannel)
		keyString.WriteString("_")
	}

	key := keyString.String()
	key = strings.Trim(key, "_")
	//cacheKey := dfs.ProductCode + "_" + trans.AssociatedTable + "_" + key

	return key
}

func GetRiders(mp models.ProductModelPoint, features models.ProductFeatures) float64 {

	var riderBenefitSA = 0.0
	if features.RiderBenefit {
		if mp.Repatriation {
			riderBenefitSA += GetRiderValue("Repatriation", mp.ProductCode, mp.Plan)
		}
		if mp.EducatorOption {
			riderBenefitSA += GetRiderValue("Educator", mp.ProductCode, mp.Plan)
		}
		if mp.Tombstone {
			riderBenefitSA += GetRiderValue("Tombstone", mp.ProductCode, mp.Plan)
		}
		if mp.CowBenefit {
			riderBenefitSA += GetRiderValue("Cow", mp.ProductCode, mp.Plan)
		}
		if mp.AdditionalSumAssuredIndicator {
			riderBenefitSA += GetRiderValue("Additional SumAssured", mp.ProductCode, mp.Plan)
		}
		if mp.Grocery {
			riderBenefitSA += GetRiderValue("Grocery", mp.ProductCode, mp.Plan)
		}
	}
	return riderBenefitSA
}

// GetRidersWithError is a version of GetRiders that returns an error
func GetRidersWithError(mp models.ProductModelPoint, features models.ProductFeatures) (float64, error) {
	var riderBenefitSA = 0.0
	if features.RiderBenefit {
		if mp.Repatriation {
			value, err := GetRiderValueWithError("Repatriation", mp.ProductCode, mp.Plan)
			if err != nil {
				return 0, errors.Wrap(err, "failed to get Repatriation rider value")
			}
			riderBenefitSA += value
		}
		if mp.EducatorOption {
			value, err := GetRiderValueWithError("Educator", mp.ProductCode, mp.Plan)
			if err != nil {
				return 0, errors.Wrap(err, "failed to get Educator rider value")
			}
			riderBenefitSA += value
		}
		if mp.Tombstone {
			value, err := GetRiderValueWithError("Tombstone", mp.ProductCode, mp.Plan)
			if err != nil {
				return 0, errors.Wrap(err, "failed to get Tombstone rider value")
			}
			riderBenefitSA += value
		}
		if mp.CowBenefit {
			value, err := GetRiderValueWithError("Cow", mp.ProductCode, mp.Plan)
			if err != nil {
				return 0, errors.Wrap(err, "failed to get Cow rider value")
			}
			riderBenefitSA += value
		}
		if mp.AdditionalSumAssuredIndicator {
			value, err := GetRiderValueWithError("Additional SumAssured", mp.ProductCode, mp.Plan)
			if err != nil {
				return 0, errors.Wrap(err, "failed to get Additional SumAssured rider value")
			}
			riderBenefitSA += value
		}
		if mp.Grocery {
			value, err := GetRiderValueWithError("Grocery", mp.ProductCode, mp.Plan)
			if err != nil {
				return 0, errors.Wrap(err, "failed to get Grocery rider value")
			}
			riderBenefitSA += value
		}
	}
	return riderBenefitSA, nil
}

func GetRiderValue(benefit string, prodCode string, plan string) float64 {
	var prodRider models.ProductRider
	key := "rider_" + prodCode + "_" + benefit + "_" + plan

	rider, found := Cache.Get(key)
	if found {
		prodRider = rider.(models.ProductRider)
	} else {
		err := DB.Where("rider_benefit =? and product_code = ?", benefit, prodCode).First(&prodRider).Error
		if err == nil {
			Cache.Set(key, prodRider, 1)
		}
	}

	switch plan {
	case "A":
		return prodRider.A
	case "B":
		return prodRider.B
	case "C":
		return prodRider.C
	case "D":
		return prodRider.D
	case "E":
		return prodRider.E
	case "F":
		return prodRider.F
	case "G":
		return prodRider.G
	case "H":
		return prodRider.H
	case "I":
		return prodRider.I
	case "J":
		return prodRider.J
	case "K":
		return prodRider.K
	case "L":
		return prodRider.L
	case "M":
		return prodRider.M
	case "N":
		return prodRider.N
	case "O":
		return prodRider.O
	}
	return 0
}

// GetRiderValueWithError is a version of GetRiderValue that returns an error
func GetRiderValueWithError(benefit string, prodCode string, plan string) (float64, error) {
	var prodRider models.ProductRider
	key := "rider_" + prodCode + "_" + benefit + "_" + plan

	rider, found := Cache.Get(key)
	if found {
		prodRider = rider.(models.ProductRider)
	} else {
		err := DB.Where("rider_benefit =? and product_code = ?", benefit, prodCode).First(&prodRider).Error
		if err != nil {
			return 0, errors.Wrap(err, "failed to get rider value")
		}
		Cache.Set(key, prodRider, 1)
	}

	switch plan {
	case "A":
		return prodRider.A, nil
	case "B":
		return prodRider.B, nil
	case "C":
		return prodRider.C, nil
	case "D":
		return prodRider.D, nil
	case "E":
		return prodRider.E, nil
	case "F":
		return prodRider.F, nil
	case "G":
		return prodRider.G, nil
	case "H":
		return prodRider.H, nil
	case "I":
		return prodRider.I, nil
	case "J":
		return prodRider.J, nil
	case "K":
		return prodRider.K, nil
	case "L":
		return prodRider.L, nil
	case "M":
		return prodRider.M, nil
	case "N":
		return prodRider.N, nil
	case "O":
		return prodRider.O, nil
	}
	return 0, nil
}

func GetPricingRiderValue(benefit string, prodCode string, plan string) float64 {
	var prodRider models.ProductPricingRider
	key := "rider_" + prodCode + "_" + benefit + "_" + plan

	rider, found := Cache.Get(key)
	if found {
		prodRider = rider.(models.ProductPricingRider)
	} else {
		err := DB.Where("rider_benefit =? and product_code = ?", benefit, prodCode).First(&prodRider).Error
		if err == nil {
			Cache.Set(key, prodRider, 1)
		}
	}

	switch plan {
	case "A":
		return prodRider.A
	case "B":
		return prodRider.B
	case "C":
		return prodRider.C
	case "D":
		return prodRider.D
	case "E":
		return prodRider.E
	case "F":
		return prodRider.F
	case "G":
		return prodRider.G
	case "H":
		return prodRider.H
	case "I":
		return prodRider.I
	case "J":
		return prodRider.J
	case "K":
		return prodRider.K
	case "L":
		return prodRider.L
	case "M":
		return prodRider.M
	case "N":
		return prodRider.N
	case "O":
		return prodRider.O
	}
	return 0
}

func GetShock(month int, shockBasis string, productShockCount int) (models.ProductShock, error) {
	var shock models.ProductShock
	if month > productShockCount {
		month = productShockCount
	}

	key := "shock_" + "_" + shockBasis + "_" + strconv.Itoa(month)

	cacheKey := key
	cached, found := Cache.Get(cacheKey)

	if found {
		shock = cached.(models.ProductShock)
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

func GetReinsuranceStructure(prodCode string, year int, variable string) float64 {
	var prodreinsurance models.ProductReinsurance
	key := prodCode + "_" + strconv.Itoa(year)

	reInsuranceStructure, found := Cache.Get(key)
	if found {
		prodreinsurance = reInsuranceStructure.(models.ProductReinsurance)
	} else {
		err := DB.Where("product_code =? and treaty_year= ?", prodCode, year).First(&prodreinsurance).Error
		if err == nil {
			Cache.Set(key, prodreinsurance, 1)
		}
		if err != nil {
			fmt.Println(err)
		}
	}

	switch variable {
	case "flat_annual_reins_prem_rate":
		return prodreinsurance.FlatAnnualReinsPremRate
	case "level1_ceded_proportion":
		return prodreinsurance.Level1CededProportion
	case "level1_lowerbound":
		return prodreinsurance.Level1Lowerbound
	case "level1_upperbound":
		return prodreinsurance.Level1Upperbound
	case "level2_ceded_proportion":
		return prodreinsurance.Level2CededProportion
	case "level2_lowerbound":
		return prodreinsurance.Level2Lowerbound
	case "level2_upperbound":
		return prodreinsurance.Level2Upperbound
	case "level3_ceded_proportion":
		return prodreinsurance.Level3CededProportion
	case "level3_lowerbound":
		return prodreinsurance.Level3Lowerbound
	case "level3_upperbound":
		return prodreinsurance.Level3Upperbound
	case "ceding_commission":
		return prodreinsurance.CedingCommission
	}
	return 0
}

// GetReinsuranceStructureWithError is a version of GetReinsuranceStructure that returns an error
func GetReinsuranceStructureWithError(prodCode string, year int, variable string) (float64, error) {
	var prodreinsurance models.ProductReinsurance
	key := prodCode + "_" + strconv.Itoa(year)

	reInsuranceStructure, found := Cache.Get(key)
	if found {
		prodreinsurance = reInsuranceStructure.(models.ProductReinsurance)
	} else {
		err := DB.Where("product_code =? and treaty_year= ?", prodCode, year).First(&prodreinsurance).Error
		if err != nil {
			return 0, errors.Wrap(err, "failed to get reinsurance structure")
		}
		Cache.Set(key, prodreinsurance, 1)
	}

	switch variable {
	case "flat_annual_reins_prem_rate":
		return prodreinsurance.FlatAnnualReinsPremRate, nil
	case "level1_ceded_proportion":
		return prodreinsurance.Level1CededProportion, nil
	case "level1_lowerbound":
		return prodreinsurance.Level1Lowerbound, nil
	case "level1_upperbound":
		return prodreinsurance.Level1Upperbound, nil
	case "level2_ceded_proportion":
		return prodreinsurance.Level2CededProportion, nil
	case "level2_lowerbound":
		return prodreinsurance.Level2Lowerbound, nil
	case "level2_upperbound":
		return prodreinsurance.Level2Upperbound, nil
	case "level3_ceded_proportion":
		return prodreinsurance.Level3CededProportion, nil
	case "level3_lowerbound":
		return prodreinsurance.Level3Lowerbound, nil
	case "level3_upperbound":
		return prodreinsurance.Level3Upperbound, nil
	case "ceding_commission":
		return prodreinsurance.CedingCommission, nil
	}
	return 0, nil
}

// GetNonLife Risk Rates
func GetNonLifeRiskRate(arg models.TransitionRateArguments) float64 {
	tablename := "product_non_life_ratings"
	var riskRateTable models.ProductNonLifeRating

	key := "rf_" + arg.ProductCode + "_" + strconv.Itoa(arg.DurationIfM) + "_" + strconv.Itoa(arg.Year)
	cacheKey := tablename + "_" + key
	cached, found := Cache.Get(cacheKey)

	if found {
		result := cached.(float64)
		return result
	} else {
		//fmt.Println("cache missed: ", key)
	}

	err := DB.Where("year=? and product_code =? and duration_if_m=?", arg.Year, arg.ProductCode, arg.DurationIfM).First(&riskRateTable).Error
	if err == nil {
		Cache.Set(key, riskRateTable, 1)
	}
	return riskRateTable.AnnualRiskRate
}

func GetUnitFundCharge(month int, mp models.ProductModelPoint, run models.RunParameters) models.ProductUnitFundCharge {

	durationYear := int((month-1)/12) + 1
	var unitFundCharge models.ProductUnitFundCharge
	key := fmt.Sprintf("unitcharge_%d_%s", durationYear, mp.FundCode)
	cached, found := Cache.Get(key)
	if found {
		unitFundCharge = cached.(models.ProductUnitFundCharge)
	} else {
		err := DB.Where("product_code = ? and duration_if_y = ? and fund_code =?", mp.ProductCode, durationYear, mp.FundCode).First(&unitFundCharge).Error
		if err == nil && unitFundCharge.ID > 0 {
			Cache.Set(key, unitFundCharge, 1)
		} else {
			var jobProduct models.JobProduct
			DB.Where("product_code = ? and projection_job_id = ?", mp.ProductCode, run.ProjectionJobID).First(&jobProduct)
			var errString string

			if mp.FundCode == "" {
				errString = fmt.Sprintf("missing base fund code in the model point")
			} else {
				errString = fmt.Sprintf("unit fund charge data was not found for fund code %s and duration year %d", mp.FundCode, durationYear)
			}

			reportError(jobProduct, "Unit Fund Charge", errString)
			//Looks like we also need to cache this zero value
			Cache.Set(key, unitFundCharge, 1)

		}
	}
	return unitFundCharge
}

func GetFundInvestmentReturns(fundYear int, mp models.ProductModelPoint, run models.RunParameters) models.ProductInvestmentReturn {

	var fundReturns = models.ProductInvestmentReturn{}
	key := fmt.Sprintf("investmentreturns%d_%s", fundYear, mp.FundCode)
	cached, found := Cache.Get(key)
	if found {
		fundReturns = cached.(models.ProductInvestmentReturn)
	} else {
		err := DB.Where("product_code = ? and fund_year = ? and fund_code = ?", mp.ProductCode, fundYear, mp.FundCode).First(&fundReturns).Error
		if err == nil && fundReturns.ID > 0 {
			Cache.Set(key, fundReturns, 1)
		} else {
			var jobProduct models.JobProduct
			DB.Where("product_code = ? and projection_job_id = ?", mp.ProductCode, run.ProjectionJobID).First(&jobProduct)
			var errString = ""
			if mp.FundCode == "" {
				errString = fmt.Sprintf("missing base fund code in the model point")
			} else {
				errString = fmt.Sprintf("fund investment return data was not found for fund code %s and year %d", mp.FundCode, fundYear)
			}

			reportError(jobProduct, "Fund Investment Returns", errString)

			//Looks like we also need to cache this zero value
			Cache.Set(key, fundReturns, 1)

		}
	}
	return fundReturns
}

func GetBsaFundInvestmentReturns(fundYear int, mp models.ProductModelPoint, run models.RunParameters) models.ProductInvestmentReturn {

	var fundReturns = models.ProductInvestmentReturn{}
	key := fmt.Sprintf("investmentreturns%d_%s", fundYear, mp.BsaFundCode)
	cached, found := Cache.Get(key)
	if found {
		fundReturns = cached.(models.ProductInvestmentReturn)
	} else {
		err := DB.Where("product_code = ? and fund_year = ? and fund_code = ?", mp.ProductCode, fundYear, mp.BsaFundCode).First(&fundReturns).Error
		if err == nil && fundReturns.ID > 0 {
			Cache.Set(key, fundReturns, 1)
		} else {
			//get the associated jobProduct
			var jobProduct models.JobProduct
			DB.Where("product_code = ? and projection_job_id = ?", mp.ProductCode, run.ProjectionJobID).First(&jobProduct)
			var errString string
			if mp.BsaFundCode == "" {
				errString = fmt.Sprintf("missing bsa fund code in the model point")
			} else {
				errString = fmt.Sprintf("base investment return not found for base fund code %s and year %d", mp.BsaFundCode, fundYear)
			}
			reportError(jobProduct, "Base Fund Investment Returns", errString)
			//Looks like we also need to cache this zero value
			Cache.Set(key, fundReturns, 1)
		}
	}
	return fundReturns
}

func GetAssetDistribution(year int, mp models.ProductModelPoint, run models.RunParameters) models.ProductFundAssetDistribution {
	var assetDistribution = models.ProductFundAssetDistribution{}
	key := fmt.Sprintf("assetdistribution_%s_%s", mp.ProductCode, mp.FundCode)
	cached, found := Cache.Get(key)
	if found {
		assetDistribution = cached.(models.ProductFundAssetDistribution)
	} else {
		err := DB.Where("product_code = ? and fund_code = ?", mp.ProductCode, mp.FundCode).First(&assetDistribution).Error
		if err == nil && assetDistribution.ID > 0 {
			Cache.Set(key, assetDistribution, 1)
		} else {
			var jobProduct models.JobProduct
			DB.Where("product_code = ? and projection_job_id = ?", mp.ProductCode, run.ProjectionJobID).First(&jobProduct)
			var errString string

			if mp.FundCode == "" {
				errString = fmt.Sprintf("missing base fund code in the model point")
			} else {
				errString = fmt.Sprintf("base investment return not found for fund code %s and year %d", mp.FundCode, year)
			}

			reportError(jobProduct, "Asset Distribution", errString)

			//Looks like we also need to cache this zero value
			Cache.Set(key, assetDistribution, 1)
		}
	}
	return assetDistribution
}

func reportError(jobProduct models.JobProduct, failurePoint string, error string) {
	var jobProductRunError = models.JobProductRunError{}
	jobProductRunError.JobProductID = jobProduct.ID
	jobProductRunError.ProductCode = jobProduct.ProductCode
	jobProductRunError.ProjectionJobID = jobProduct.ProjectionJobID
	jobProductRunError.FailurePoint = failurePoint
	jobProductRunError.Error = error
	//Delete any existing errors for this jobProduct
	//DB.Where("job_product_id = ? and projection_job_id = ? and product_code = ? and error = ? and failure_point = ? ", jobProductRunError.JobProductID, jobProductRunError.ProjectionJobID, jobProductRunError.ProductCode, jobProductRunError.Error, jobProductRunError.FailurePoint).Delete(models.JobProductRunError{})
	DB.Save(&jobProductRunError)
}

func GetSurrenderQuadraticCoefficients(mp models.ProductModelPoint, run models.RunParameters) models.ProductSurrenderValueCoefficient {

	var quadraticCoefficients = models.ProductSurrenderValueCoefficient{}
	key := fmt.Sprintf("surrendervaluecoeffients_%s", mp.SurrenderValueCode)
	cached, found := Cache.Get(key)
	if found {
		quadraticCoefficients = cached.(models.ProductSurrenderValueCoefficient)
	} else {
		err := DB.Where("surrender_value_code = ?", mp.SurrenderValueCode).First(&quadraticCoefficients).Error
		if err == nil && quadraticCoefficients.ID > 0 {
			Cache.Set(key, quadraticCoefficients, 1)
		} else {
			//get the associated jobProduct
			var jobProduct models.JobProduct
			DB.Where("product_code = ? and projection_job_id = ?", mp.ProductCode, run.ProjectionJobID).First(&jobProduct)
			var errString string
			if mp.SurrenderValueCode == "" {
				errString = fmt.Sprintf("missing surrender value code in the model point")
			} else {
				errString = fmt.Sprintf("surrender value coefficients not found for surrender value code %s", mp.SurrenderValueCode)
			}
			reportError(jobProduct, "Surrender Value Coefficients", errString)
			//Looks like we also need to cache this zero value
			Cache.Set(key, quadraticCoefficients, 1)
		}
	}
	return quadraticCoefficients
}

func GetMaturityBenefitPattern(remainingTermYear int, mp models.ProductModelPoint, run models.RunParameters) models.ProductMaturityPattern {
	var maturityBenefitPattern = models.ProductMaturityPattern{}
	key := fmt.Sprintf("maturityPattern_%s_%d", mp.MaturityBenefitCode, remainingTermYear)
	cached, found := Cache.Get(key)
	if found {
		maturityBenefitPattern = cached.(models.ProductMaturityPattern)
	} else {
		err := DB.Where("maturity_pattern_code = ? and remaining_term_year = ?", mp.MaturityBenefitCode, remainingTermYear).First(&maturityBenefitPattern).Error
		if err == nil && maturityBenefitPattern.ID > 0 {
			Cache.Set(key, maturityBenefitPattern, 1)
		} else {
			//get the associated jobProduct
			var jobProduct models.JobProduct
			DB.Where("product_code = ? and projection_job_id = ?", mp.ProductCode, run.ProjectionJobID).First(&jobProduct)
			var errString string
			if mp.SurrenderValueCode == "" {
				errString = fmt.Sprintf("missing maturity benefit code in the model point")
			} else {
				errString = fmt.Sprintf("maturity benefit pattern not found for maturity benefit code %s and remaining term %d", mp.MaturityBenefitCode, remainingTermYear)
			}
			reportError(jobProduct, "Maturity Benefit Pattern", errString)
			//Looks like we also need to cache this zero value
			Cache.Set(key, maturityBenefitPattern, 1)
		}
	}
	return maturityBenefitPattern
}
