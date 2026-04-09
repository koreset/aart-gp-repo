package services

import (
	"api/models"
	"api/utils"
	"fmt"
	"github.com/iancoleman/strcase"
	"math"
	"strconv"
	"strings"
)

func GetAllPricingRuns() []models.PricingRun {
	var runs []models.PricingRun
	DB.Preload("PricingConfig").Order("id desc").Find(&runs)
	return runs
}

func GetPricingRun(runId int) models.PricingRun {
	var run models.PricingRun
	DB.Preload("PricingConfig").Where("id=?", runId).Find(&run)
	return run
}

func DeletePricingRun(runId int) error {
	err := DB.Where("pricing_run_id=?", runId).Delete(&models.PricingConfig{}).Error
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = DB.Where("id = ?", runId).Delete(&models.PricingRun{}).Error
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = DB.Where("run_id = ?", runId).Delete(&models.PricingPoint{}).Error
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = DB.Where("run_id = ?", runId).Delete(&models.Profitability{}).Error
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = DB.Where("pricing_run_id = ?", runId).Delete(&models.ModelPointPricing{}).Error
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func PricingRunNameDoesNotExist(runName string) bool {
	var job models.PricingRun
	DB.Where("name = ?", runName).First(&job)
	if job.Name == runName {
		return false
	} else {
		return true
	}
}

func GetProductPricingParameters(prodCode string) (models.ProductPricingParameters, []string, []string, error) {
	var pricingParam models.ProductPricingParameters
	var bases []string
	var shockBases []string
	err := DB.Where("product_code = ?", prodCode).First(&pricingParam).Error
	if err != nil {
		return pricingParam, bases, shockBases, err
	}
	err = DB.Table("product_pricing_parameters").Where("product_code = ?", prodCode).Distinct("basis").Pluck("basis", &bases).Error

	if err != nil {
		return pricingParam, bases, shockBases, err
	}

	err = DB.Table("product_pricing_shocks").Where("product_code = ?", prodCode).Distinct("shock_basis").Pluck("shock_basis", &shockBases).Error

	return pricingParam, bases, shockBases, nil
}

func GetPricingParameters(prodCode string) ([]models.PricingParameter, error) {
	var pricingParam []models.PricingParameter
	if prodCode == "" {
		err := DB.Find(&pricingParam).Error
		if err != nil {
			return pricingParam, err
		}
	} else {
		err := DB.Where("product_code = ?", prodCode).Find(&pricingParam).Error
		if err != nil {
			return pricingParam, err
		}
	}
	return pricingParam, nil
}

func GetPricingDemographicsData(prodCode string) ([]models.PricingPolicyDemographic, error) {
	var pricingDemographics []models.PricingPolicyDemographic
	err := DB.Where("product_code = ?", prodCode).Find(&pricingDemographics).Error
	if err != nil {
		return pricingDemographics, err
	}
	return pricingDemographics, nil
}

func DeletePricingParameters(productCode string) error {
	err := DB.Where("product_code = ?", productCode).Delete(&models.PricingParameter{}).Error
	if err != nil {
		return err
	}
	return nil
}

func DeletePricingDemographics(productCode string) error {
	err := DB.Where("product_code = ?", productCode).Delete(&models.PricingPolicyDemographic{}).Error
	if err != nil {
		return err
	}
	return nil
}

type SumAssuredData struct {
	SumAssured float64
}

func DeletePricingScenario(scenarioId int) error {
	//delete the scenario from pricing configs
	err := DB.Where("id=?", scenarioId).Delete(&models.PricingConfig{}).Error
	if err != nil {
		fmt.Println(err)
		return err
	}
	//delete the scenario from profitability
	err = DB.Where("scenario_id=?", scenarioId).Delete(&models.Profitability{}).Error
	if err != nil {
		fmt.Println(err)
		return err
	}

	//delete the scenario from model point pricing
	err = DB.Where("scenario_id=?", scenarioId).Delete(&models.ModelPointPricing{}).Error
	if err != nil {
		fmt.Println(err)
		return err
	}

	//delete the scenario from pricing points
	err = DB.Where("scenario_id=?", scenarioId).Delete(&models.PricingPoint{}).Error
	if err != nil {
		fmt.Println(err)
		return err
	}

	//delete the scenario from aggregated pricing points
	err = DB.Where("scenario_id=?", scenarioId).Delete(&models.AggregatedPricingPoint{}).Error
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func GetPricingData(scenarioId int) map[string]interface{} {
	var data = make(map[string]interface{})
	var profitability models.Profitability
	var pds []models.PricingDistribution
	DB.Where("scenario_id=?", scenarioId).Find(&profitability)
	data["profitability"] = profitability
	var sumAssureds []SumAssuredData

	////get the pricing run from the profitability run_id
	//var pricingRun models.PricingRun
	//DB.Where("id=?", profitability.RunID).First(&pricingRun)
	//data["pricing_run"] = pricingRun

	err := DB.Table("model_point_pricings").Distinct("sum_assured").Where("scenario_id=? ", scenarioId).Order("sum_assured asc").Find(&sumAssureds).Error
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(sumAssureds)
	var results []map[string]interface{}
	var prodFeatures models.ProductFeatures
	var mpOne models.ModelPointPricing
	DB.Where("scenario_id=?", scenarioId).First(&mpOne)
	DB.Where("product_code=?", mpOne.ProductCode).First(&prodFeatures)

	var minAge int
	var maxAge int

	//get the min and max age for a given scenario id
	err = DB.Table("model_point_pricings").Select("min(age) as min_age, max(age) as max_age").Where("scenario_id=?", scenarioId).Scan(&data).Error
	if err != nil {
		fmt.Println(err)
	}

	minAge = int(data["min_age"].(int64))
	maxAge = int(data["max_age"].(int64))

	for i := minAge; i <= maxAge; i++ {
		var mpps []models.ModelPointPricing

		var pd models.PricingDistribution

		DB.Where("scenario_id=? and age=?", scenarioId, i).Find(&mpps)

		var pdResult = make(map[string]interface{})
		pd.Age = i
		for _, mpp := range mpps {

			for _, sa := range sumAssureds {
				if mpp.SumAssured == sa.SumAssured {
					// create a key for the sum assured and gender
					sumAssuredKey := strconv.FormatFloat(sa.SumAssured, 'f', -1, 64) + " (" + mpp.Gender + ")"
					// check if the key exists in the map
					if _, ok := pdResult[sumAssuredKey]; ok {
						if prodFeatures.CreditLife {
							pdResult[sumAssuredKey] = pdResult[sumAssuredKey].(float64) + utils.FloatPrecision(mpp.CalculatedPremiumRate, AccountingPrecision)*mpp.Weighting
						}
						if !prodFeatures.CreditLife {
							pdResult[sumAssuredKey] = pdResult[sumAssuredKey].(float64) + math.Ceil(mpp.CalculatedAnnualPremium)*mpp.Weighting
						}

					} else {
						// if it doesn't exist, create the key and add the value
						if prodFeatures.CreditLife {
							pdResult[sumAssuredKey] = utils.FloatPrecision(mpp.CalculatedPremiumRate, AccountingPrecision) * mpp.Weighting
						}
						if !prodFeatures.CreditLife {
							pdResult[sumAssuredKey] = math.Ceil(mpp.CalculatedAnnualPremium) * mpp.Weighting
						}
					}
				}
			}

			if len(sumAssureds) > 0 {
				if mpp.Gender == "M" && mpp.SumAssured == sumAssureds[0].SumAssured {
					if prodFeatures.CreditLife {
						pd.Range1Male += utils.FloatPrecision(mpp.CalculatedPremiumRate, AccountingPrecision) * mpp.Weighting
					}
					if !prodFeatures.CreditLife {
						pd.Range1Male += math.Ceil(mpp.CalculatedAnnualPremium) * mpp.Weighting
					}
				}
				if mpp.Gender == "F" && mpp.SumAssured == sumAssureds[0].SumAssured {
					if prodFeatures.CreditLife {
						pd.Range1Female += utils.FloatPrecision(mpp.CalculatedPremiumRate, AccountingPrecision) * mpp.Weighting
					}
					if !prodFeatures.CreditLife {
						pd.Range1Female += math.Ceil(mpp.CalculatedAnnualPremium) * mpp.Weighting
					}

				}
			}

			if len(sumAssureds) > 1 {
				if mpp.Gender == "M" && mpp.SumAssured == sumAssureds[1].SumAssured {
					if prodFeatures.CreditLife {
						pd.Range2Male += utils.FloatPrecision(mpp.CalculatedPremiumRate, AccountingPrecision) * mpp.Weighting
					}
					if !prodFeatures.CreditLife {
						pd.Range2Male += math.Ceil(mpp.CalculatedAnnualPremium) * mpp.Weighting
					}

				}
				if mpp.Gender == "F" && mpp.SumAssured == sumAssureds[1].SumAssured {
					if prodFeatures.CreditLife {
						pd.Range2Female += utils.FloatPrecision(mpp.CalculatedPremiumRate, AccountingPrecision) * mpp.Weighting
					}
					if !prodFeatures.CreditLife {
						pd.Range2Female += math.Ceil(mpp.CalculatedAnnualPremium) * mpp.Weighting
					}

				}
			}

			if len(sumAssureds) > 2 {
				if mpp.Gender == "M" && mpp.SumAssured == sumAssureds[2].SumAssured {
					if prodFeatures.CreditLife {
						pd.Range3Male += utils.FloatPrecision(mpp.CalculatedPremiumRate, AccountingPrecision) * mpp.Weighting
					}
					if !prodFeatures.CreditLife {
						pd.Range3Male += math.Ceil(mpp.CalculatedAnnualPremium) * mpp.Weighting
					}

				}
				if mpp.Gender == "F" && mpp.SumAssured == sumAssureds[2].SumAssured {
					if prodFeatures.CreditLife {
						pd.Range3Female += utils.FloatPrecision(mpp.CalculatedPremiumRate, AccountingPrecision) * mpp.Weighting
					}
					if !prodFeatures.CreditLife {
						pd.Range3Female += math.Ceil(mpp.CalculatedAnnualPremium) * mpp.Weighting
					}

				}
			}

			if len(sumAssureds) > 3 {
				if mpp.Gender == "M" && mpp.SumAssured == sumAssureds[3].SumAssured {
					if prodFeatures.CreditLife {
						pd.Range4Male += utils.FloatPrecision(mpp.CalculatedPremiumRate, AccountingPrecision) * mpp.Weighting
					}
					if !prodFeatures.CreditLife {
						pd.Range4Male += math.Ceil(mpp.CalculatedAnnualPremium) * mpp.Weighting
					}

				}
				if mpp.Gender == "F" && mpp.SumAssured == sumAssureds[3].SumAssured {
					if prodFeatures.CreditLife {
						pd.Range4Female += utils.FloatPrecision(mpp.CalculatedPremiumRate, AccountingPrecision) * mpp.Weighting
					}
					if !prodFeatures.CreditLife {
						pd.Range4Female += math.Ceil(mpp.CalculatedAnnualPremium) * mpp.Weighting
					}

				}
			}
		}
		//results = append(results, utils.StructToMapWithNonZeroFields(pd))
		pdResult["Age/Sum Assured"] = i
		results = append(results, pdResult)
		pds = append(pds, pd)
	}
	data["pricing_distribution"] = pds
	data["filtered"] = results

	//Get aggregated data if it exists
	var aggResults []models.AggregatedPricingPoint
	err = DB.Where("scenario_id=?", scenarioId).Order("projection_month asc").Find(&aggResults).Error
	if err != nil {
		fmt.Println(err)
	}
	data["aggregated_points"] = aggResults

	var headers []string

	for i, sa := range sumAssureds {
		if i < 4 {
			headers = append(headers, fmt.Sprintf("%d (M)", int(sa.SumAssured)))
			headers = append(headers, fmt.Sprintf("%d (F)", int(sa.SumAssured)))
		}
	}

	if len(headers) < 8 {
		for i := len(headers); i < 8; i++ {
			headers = append(headers, "N/A")
		}
	}
	data["headers"] = headers
	return data
}

func GetPricingTables(tables []models.ProductPricingTable, productCode string) ([]models.ProductPricingTable, error) {
	var tableName string
	var count int64
	for i, _ := range tables {
		switch tables[i].Class {
		case "TransitionRates":
			tableName = strings.ToLower(productCode) + "_pricing_" + strings.ToLower(tables[i].Name)
			count = 0
			DB.Table(tableName).Count(&count)
			if count > 0 {
				tables[i].Populated = true
			}
		case "Margins":
			tableName = "product_pricing_" + strings.ToLower(tables[i].Name)
			count = 0
			DB.Table(tableName).Where("product_code =?", productCode).Count(&count)
			if count > 0 {
				tables[i].Populated = true
			}
		case "Global":
			tableName = "pricing_" + strings.ToLower(tables[i].Name)
			if DB.Migrator().HasTable(tableName) {
				count = 0
				DB.Table(tableName).Count(&count)
				if count > 0 {
					tables[i].Populated = true
				}
			}
		case "BenefitStructure":
			tableName = "product_pricing_" + strings.ToLower(tables[i].Name)
			if DB.Migrator().HasTable(tableName) {
				count = 0
				DB.Table(tableName).Where("product_code =?", productCode).Count(&count)
				if count > 0 {
					tables[i].Populated = true
				}
			}
		case "LapseMargins":
			tableName = "product_pricing_" + strings.ToLower(strcase.ToSnake(tables[i].Name))
			count = 0
			DB.Table(tableName).Where("product_code =?", productCode).Count(&count)
			if count > 0 {
				tables[i].Populated = true
			}

		case "Parameters":
			tableName = "product_pricing_" + strings.ToLower(tables[i].Name)
			count = 0
			DB.Table(tableName).Where("product_code =?", productCode).Count(&count)
			if count > 0 {
				tables[i].Populated = true
			}

		case "Valuations":
			tableName = "product_pricing_" + strings.ToLower(tables[i].Name)
			count = 0
			DB.Table(tableName).Where("product_code =?", productCode).Count(&count)
			if count > 0 {
				tables[i].Populated = true
			}

		case "Profitabilities":
			tableName = "product_pricing_" + strings.ToLower(tables[i].Name)
			count = 0
			DB.Table(tableName).Where("product_code =?", productCode).Count(&count)
			if count > 0 {
				tables[i].Populated = true
			}

		case "Distribution":
			tableName = "product_pricing_" + strings.ToLower(tables[i].Name)
			count = 0
			DB.Table(tableName).Where("product_code =?", productCode).Count(&count)
			if count > 0 {
				tables[i].Populated = true
			}
		}

	}

	return tables, nil
}

func GetPricingModelPoints(prodCode string) ([]models.PricingModelPointVersion, []string) {
	var mps []models.PricingModelPointVersion
	var versions []string

	tableName := strings.ToLower(prodCode) + "_pricing_modelpoints"
	DB.Table(tableName).Select("distinct pricing_mp_version").Find(&versions)
	//Use the versions to aggregate the model points
	for _, version := range versions {
		var mpSet models.PricingModelPointVersion
		var mpCount int64
		DB.Table(tableName).Where("pricing_mp_version=?", version).Count(&mpCount)
		mpSet.Version = version
		mpSet.Count = int(mpCount)
		mps = append(mps, mpSet)
	}

	return mps, versions
}

func GetPricingModelPointsForVersion(prodCode, version string) ([]models.ProductPricingModelPoint, error) {
	var mps []models.ProductPricingModelPoint
	tableName := strings.ToLower(prodCode) + "_pricing_modelpoints"
	err := DB.Table(tableName).Where("pricing_mp_version=?", version).Find(&mps).Error
	if err != nil {
		return mps, err
	}

	return mps, nil
}

func GetPricingModelPointVersions(prodCode string) []string {
	var versions []string

	tableName := strings.ToLower(prodCode) + "_pricing_modelpoints"
	DB.Table(tableName).Select("distinct pricing_mp_version").Find(&versions)
	//Use the versions to aggregate the model points

	return versions
}

func GetModelPointsForPricing(pricingRunId int) ([]models.ModelPointPricing, error) {
	var mpps []models.ModelPointPricing
	err := DB.Where("pricing_run_id = ?", pricingRunId).Find(&mpps).Error
	if err != nil {
		return mpps, err
	}

	return mpps, nil
}

func GetPricingPoints(pricingRunId, scenarioId int, productCode string) ([]models.PricingPoint, error) {
	var points []models.PricingPoint

	DB.Where("scenario_id=? and run_id=?", scenarioId, pricingRunId).Find(&points)

	return points, nil
}

func GetProductPricingTable(tableId int) models.ProductPricingTable {
	var prodTable models.ProductPricingTable
	err := DB.Where("id = ?", tableId).First(&prodTable).Error
	if err != nil {
		fmt.Println(err)
	}
	return prodTable
}

func GetPricingShock(month int, shockBasis string, productPricingShockCount int) (models.ProductPricingShock, error) {
	var shock models.ProductPricingShock
	if month > productPricingShockCount {
		month = productPricingShockCount
	}
	if month > 0 {
		key := "shock_" + "_" + shockBasis + "_" + strconv.Itoa(month)

		cacheKey := key
		cached, found := PricingCache.Get(cacheKey)

		if found {
			shock = cached.(models.ProductPricingShock)
			//return shock, nil
		} else {
			err := DB.Where("shock_basis = ? and projection_month = ?", shockBasis, month).First(&shock).Error
			if err != nil {
				return shock, err
			}
			PricingCache.Set(key, shock, 1)
			//return shock, nil
		}
	}
	return shock, nil
}
