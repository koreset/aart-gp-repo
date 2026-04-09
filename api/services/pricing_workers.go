package services

import (
	"api/models"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

func GetPricingLapseMarginCount(code string) int {
	var count int64
	err := DB.Table("product_pricing_lapse_margins").Where("product_code = ?", code).Count(&count).Error
	if err != nil {
		return 0
	}
	return int(count)
}

func LoadPricingMortalityRates(productCode string) {
	gender := []string{"M", "F"}
	tablename := strings.ToLower(productCode) + "_pricing_" + strings.ToLower(pricingMortalityTableName)

	for age := 1; age <= 121; age++ {
		for _, gen := range gender {
			key := strconv.Itoa(age) + "_" + gen
			cacheKey := tablename + "_" + key
			row := DB.Table(tablename).Where("anb_gender = ? ", key).Select("qx").Row()
			var qx float64
			err := row.Scan(&qx)
			if err != nil {
				fmt.Println(err)
			}
			PricingCache.Set(cacheKey, qx, 1)
		}
	}

	//sanity check
	cached, found := PricingCache.Get(tablename + "_3_F")
	if found {
		fmt.Println("load key: ", cached)
	} else {
		fmt.Println("NOT FOUND")
	}
}

func LoadPricingAccidentalMortalityRates(productCode string) {
	gender := []string{"M", "F"}
	tablename := strings.ToLower(productCode) + "_pricing_" + strings.ToLower(pricingMortalityAccidentalTableName)

	for age := 1; age <= 121; age++ {
		for _, gen := range gender {
			key := strconv.Itoa(age) + "_" + gen
			cacheKey := tablename + "_" + key
			row := DB.Table(tablename).Where("anb_gender = ? ", key).Select("acc_qx_prop").Row()
			var qx float64
			err := row.Scan(&qx)
			if err != nil {
				fmt.Println(err)
			}
			Cache.Set(cacheKey, qx, 1)
		}
	}

	//sanity check
	cached, found := PricingCache.Get(tablename + "_3_F")
	if found {
		fmt.Println("load key: ", cached)
	} else {
		fmt.Println("NOT FOUND")
	}
}

func LoadPricingInflationFactor() {
	var count int64
	tableName := "pricing_yield_curve"
	DB.Table(tableName).Count(&count)
	fmt.Println(count)
	for i := 1; i <= int(count); i++ {
		var yieldCurve models.PricingYieldCurve
		key := "y-curve-" + strconv.Itoa(i)
		err := DB.Table(tableName).Where("proj_time = ?", i).First(&yieldCurve).Error
		if err != nil {
			fmt.Println("Yield Curve: ", errors.WithStack(err))
		} else {
			success := PricingCache.Set(key, yieldCurve.Inflation, 1)
			if !success {
				fmt.Println("Cache: key not stored")
			}
		}
	}
}

func LoadPricingForwardRate() {
	var count int64
	tableName := "pricing_yield_curve"

	DB.Table(tableName).Count(&count)
	fmt.Println(count)
	for i := 1; i <= int(count); i++ {
		var yieldCurve models.PricingYieldCurve
		key := "y-curve-fr-" + strconv.Itoa(i)
		err := DB.Table(tableName).Where("proj_time = ?", i).First(&yieldCurve).Error
		if err != nil {
			fmt.Println("Yield Curve: ", errors.WithStack(err))
		} else {
			success := PricingCache.Set(key, yieldCurve.NominalRate, 1)
			if !success {
				fmt.Println("Cache: key not stored")
			}
		}
	}
}

func GetPricingMortalityRateAccidentProportion(age int, gender string, productCode string) float64 {
	key := strconv.Itoa(age) + "_" + gender[:1]
	tableName := strings.ToLower(productCode) + "_pricing_" + strings.ToLower(pricingMortalityAccidentalTableName)
	cacheKey := strings.ToLower(tableName) + "_" + key
	cached, found := PricingCache.Get(cacheKey)

	if found {
		result := cached.(float64)
		if result > 0 {
			return result
		}
	}

	var accQxProp float64

	row := DB.Table(tableName).Where("anb_gender = ?", key).Select("acc_qx_prop").Row()
	row.Scan(&accQxProp)
	PricingCache.Set(cacheKey, accQxProp, 1)
	return accQxProp
}

func GetPricingInflationFactor(month int) float64 {
	var yieldCurve models.PricingYieldCurve
	key := "y-curve-" + strconv.Itoa(month)
	result, found := Cache.Get(key)

	if found {
		return result.(float64)
	} else {
		//TODO:Use pricing yield curve here instead...
		err := DB.Table("yield_curve").Where("proj_time = ?", month).First(&yieldCurve).Error
		if err != nil {
			fmt.Println("Yield Curve: ", errors.WithStack(err))
		}
		success := Cache.Set(key, yieldCurve.Inflation, 1)
		if !success {
			fmt.Println("Cache: key not stored")
		}
	}
	return yieldCurve.Inflation
}

func GetPricingForwardRate(month int, db *OwnDb) float64 {
	var yieldCurve models.PricingYieldCurve
	key := "y-curve-fr-" + strconv.Itoa(month)
	cached, found := Cache.Get(key)

	if found {
		result := cached.(float64)
		if result > 0 {
			return result
		}
	}
	if month == 0 {
		return 0
	}
	err := db.Table("yield_curve").Where("proj_time = ?", month).First(&yieldCurve).Error
	if err != nil {
		fmt.Println("db error: ", errors.WithStack(err))
	}
	success := Cache.Set(key, yieldCurve.NominalRate, 1)
	if !success {
		fmt.Println("Cache error: Not saved")
	}

	return yieldCurve.NominalRate
}

func GetPricingMortalityRate(age int, gender string, productCode string) float64 {
	tablename := strings.ToLower(productCode) + "_pricing_" + strings.ToLower(pricingMortalityTableName)
	key := strconv.Itoa(age) + "_" + gender[:1]
	cacheKey := tablename + "_" + key
	cached, found := Cache.Get(cacheKey)

	if found {
		result := cached.(float64)
		if result > 0 {
			return result
		}
	} else {
		fmt.Println("Cache missed: ", key)
	}

	row := DB.Table(tablename).Where("anb_gender = ?", key).Select("qx").Row()

	var qx float64
	row.Scan(&qx)
	Cache.Set(cacheKey, qx, 1)
	return qx

}

func GetPricingLapseRate(month int, prodCode string) float64 {
	tableName := strings.ToLower(prodCode) + "_pricing_" + strings.ToLower(pricingLapseTableName)

	if month == 0 {
		month = 1
	}

	if month > pricingLapseMonthCount {
		month = pricingLapseMonthCount
	}
	key := strconv.Itoa(int(month))
	cacheKey := tableName + "_" + key

	cached, found := Cache.Get(cacheKey)

	if found {
		result := cached.(float64)
		if result > 0 {
			return result
		}
	}

	var lapseRate float64

	row := DB.Table(tableName).Where("duration_if_m=?", key).Select("lapse_rate").Row()
	row.Scan(&lapseRate)
	Cache.Set(cacheKey, lapseRate, 1)

	return lapseRate
}

func GetPricingLapseMargin(month int, prodCode string) float64 {
	var margin models.ProductPricingLapseMargin
	if month == 0 { //A hack to fix the month zero issue with lapses. Discuss with M later.
		month = 1
	}

	if month > pricingLapseMonthCount {
		month = pricingLapseMonthCount
	}

	err := DB.Where("product_code = ?  and month = ?", prodCode, month).First(&margin).Error
	if err != nil {
		return 0
	}
	return margin.Margin
}

func GetPricingRetrenchmentRate(prodCode string, yearInView int) models.PricingRetrenchmentRate {
	var retrenchmentRate models.PricingRetrenchmentRate
	if yearInView == 0 {
		yearInView = 1
	}

	if yearInView > 5 {
		yearInView = 5
	}

	key := strings.ToLower(prodCode) + "-retr-" + strconv.Itoa(yearInView)
	cached, found := Cache.Get(key)
	if found {
		return cached.(models.PricingRetrenchmentRate)
	}

	tableName := strings.ToLower(prodCode) + "_retrenchment"
	query := GetColumnName(tableName) + "=?"
	yearValue := strconv.Itoa(yearInView)
	err := DB.Table(tableName).Where(query, yearValue).Find(&retrenchmentRate).Error
	if err != nil {
		fmt.Println(err)
	}
	Cache.Set(key, retrenchmentRate, 1)
	return retrenchmentRate
}

func GetPricingDisabilityIncidenceRate(prodCode string, dfs models.DisabilityIncidenceFactors) float64 {
	//Build the key
	var keyString strings.Builder
	keyString.WriteString(strconv.Itoa(dfs.Age))
	keyString.WriteString("_")
	if dfs.SocialEconomicClass > 0 {
		keyString.WriteString(strconv.Itoa(dfs.SocialEconomicClass))
		keyString.WriteString("_")
	}

	if dfs.OccupationalClass > 0 {
		keyString.WriteString(strconv.Itoa(dfs.OccupationalClass))
		keyString.WriteString("_")
	}
	if dfs.Gender != "" {
		keyString.WriteString(dfs.Gender[:1])
		keyString.WriteString("_")
	}

	key := keyString.String()
	key = strings.Trim(key, "_")
	//DisabilityTableName = "disability"
	tableName := strings.ToLower(prodCode) + "_" + strings.ToLower(ptnames[prodCode].DisabilityTableName)
	cacheKey := tableName + "_" + key
	cached, found := Cache.Get(cacheKey)

	if found {
		result := cached.(float64)
		if result > 0 {
			return result
		}
	}

	var incidenceRate float64
	if DisabilityColumnName == "" {
		DisabilityColumnName = GetColumnName(tableName)
	}

	query := DisabilityColumnName + " = ?"
	row := DB.Table(tableName).Where(query, key).Select("incidence_rate").Row()
	row.Scan(&incidenceRate)
	Cache.Set(cacheKey, incidenceRate, 1)
	return incidenceRate
}

func GetPricingChildFuneralService(productCode string, plan string, age int) float64 {
	key := fmt.Sprintf("%s-%s-%d", strings.ToLower(productCode), plan, age)
	cached, found := Cache.Get(key)
	if found {
		return cached.(float64)
	}

	var childFuneralService models.ProductPricingChildAdditionalSumAssured
	err := DB.Where("product_code =? and age=?", productCode, age).Find(&childFuneralService).Error
	if err != nil {
		fmt.Println("child funeral: ", err)
		return 0
	}
	var result float64
	switch plan {
	case "A":
		result = childFuneralService.A
	case "B":
		result = childFuneralService.B
	case "C":
		result = childFuneralService.C
	case "D":
		result = childFuneralService.D
	case "E":
		result = childFuneralService.E
	case "F":
		result = childFuneralService.F
	default:
		result = 0
	}
	Cache.Set(key, result, 1)

	return result
}

func GetPricingFuneralService(productCode string, plan string) float64 {
	key := fmt.Sprintf("%s-%s", strings.ToLower(productCode), plan)
	cached, found := Cache.Get(key)
	if found {
		return cached.(float64)
	}

	var fs models.ProductPricingAdditionalSumAssured
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
	default:
		result = 0
	}
	PricingCache.Set(key, result, 1)

	return result
}

func GetPricingChildSumAssured(productCode string, plan string, age int) float64 {
	key := fmt.Sprintf("%s-%s-%d", strings.ToLower(productCode), plan, age)
	cached, found := Cache.Get(key)
	if found {
		return cached.(float64)
	}

	var childSumAssured models.ProductPricingChildSumAssured
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
	default:
		result = 0
	}
	PricingCache.Set(key, result, 1)

	return result
}

func GetPricingClawback(month int) models.ProductPricingClawback {
	var clawback models.ProductPricingClawback
	DB.Where("duration_in_force_month = ?", month).First(&clawback)
	return clawback
}

func GetPricingRatingTable(code string, class string) string {
	prod, err := GetProductByCode(code)
	if err != nil {
		fmt.Println("product not found")
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
