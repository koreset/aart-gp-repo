package services

import (
	"api/models"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
)

func CreateRatingTable(key, value, tableName string) error {
	var err error
	fmt.Println(key, value, tableName)
	key = strings.ToLower(key)
	value = strings.ToLower(value)
	tableName = strings.ToLower(tableName)

	exists := DB.Migrator().HasTable(key)
	if !exists {
		if DbBackend == "postgresql" {
			err := DB.Exec(fmt.Sprintf("create table %s(%s varchar(100) primary key, %s numeric)", tableName, key,
				value)).Error
			if err != nil {
				return err
			}
		}
		if DbBackend == "mysql" {
			err := DB.Exec(fmt.Sprintf("create table %s(%s varchar(100) primary key, %s double)", tableName, key,
				value)).Error
			if err != nil {
				return err
			}
		}
		return nil
	} else {
		return err
	}
}

func GetProductTableSample(productCode string, tableID int) []map[string]interface{} {
	var productTable models.ProductTable
	var results []map[string]interface{}

	DB.Where("id = ? ", tableID).First(&productTable)

	// We need to get the actual table,
	// based on the  TransitionTable type. (TransitionRate, Global, Margins, Parameters)
	var tableName string

	switch productTable.Class {
	case "TransitionRates":
		tableName = strings.ToLower(productCode) + "_" + strings.ToLower(productTable.Name)
		results = getTransitionRateData(tableName, results)
	case "Global":
		tableName = strings.ToLower(productTable.Name)
		results = getYieldData(tableName, results)
	case "BenefitStructure":
		tableName = "product_" + strings.ToLower(strcase.ToSnake(productTable.Name))
		results = getStaticData(productCode, tableName, results)
	case "LapseMargins":
		tableName = "product_" + strings.ToLower(strcase.ToSnake(productTable.Name))
		results = getStaticData(productCode, tableName, results)

	case "Distribution":
		tableName = "product_" + strings.ToLower(productTable.Name)
		results = getStaticData(productCode, tableName, results)
	case "Margins":
		tableName = "product_" + strings.ToLower(productTable.Name)
		results = getData(productCode, tableName, results)
	case "Parameters":
		tableName = "product_" + strings.ToLower(productTable.Name)
		results = getData(productCode, tableName, results)
	case "Valuations":
		tableName = "product_" + strings.ToLower(productTable.Name)
		results = getData(productCode, tableName, results)

	}

	fmt.Println(tableName)

	return results

}

func GetProductPricingTableSample(productCode string, tableID int) []map[string]interface{} {
	var productTable models.ProductPricingTable
	var results []map[string]interface{}

	DB.Where("id = ? ", tableID).First(&productTable)

	// We need to get the actual table,
	// based on the  TransitionTable type. (TransitionRate, Global, Margins, Parameters)
	var tableName string

	switch productTable.Class {
	case "TransitionRates":
		tableName = strings.ToLower(productCode) + "_pricing_" + strings.ToLower(productTable.Name)
		results = getTransitionRateData(tableName, results)
	case "Global":
		tableName = "pricing_" + strings.ToLower(productTable.Name)
		results = getYieldData(tableName, results)
	case "BenefitStructure":
		tableName = "product_pricing_" + strings.ToLower(productTable.Name)
		results = getStaticData(productCode, tableName, results)
	case "LapseMargins":
		tableName = "product_pricing_" + strings.ToLower(strcase.ToSnake(productTable.Name))
		results = getStaticData(productCode, tableName, results)
	case "Profitabilities":
		tableName = "product_pricing_" + strings.ToLower(strcase.ToSnake(productTable.Name))
		results = getStaticData(productCode, tableName, results)
	case "Distribution":
		tableName = "product_pricing_" + strings.ToLower(productTable.Name)
		results = getStaticData(productCode, tableName, results)
	case "Margins":
		tableName = "product_pricing_" + strings.ToLower(productTable.Name)
		results = getData(productCode, tableName, results)
	case "Parameters":
		tableName = "product_pricing_" + strings.ToLower(productTable.Name)
		results = getData(productCode, tableName, results)
	case "Valuations":
		tableName = "product_pricing_" + strings.ToLower(productTable.Name)
		results = getData(productCode, tableName, results)

	}

	fmt.Println(tableName)

	return results

}

func getTransitionRateData(tableName string, results []map[string]interface{}) []map[string]interface{} {
	rows, err := DB.Raw("select * from " + tableName).Rows()
	if err != nil {
		fmt.Println(err)
	}
	if rows != nil {
		defer rows.Close()
		cols, _ := rows.Columns()

		fmt.Println(cols)
		for rows.Next() {
			data := make(map[string]interface{})
			columns := make([]string, len(cols))
			columnPointers := make([]interface{}, len(cols))
			for i, _ := range columns {
				columnPointers[i] = &columns[i]
			}

			rows.Scan(columnPointers...)
			for i, colName := range cols {
				if i == 0 {
					colValues := strings.Split(columns[i], "_")
					colNames := strings.SplitN(colName, "_", len(colValues))
					for j, v := range colNames {
						data[v] = colValues[j]
					}
				} else {
					data[colName] = columns[i]
				}
			}
			results = append(results, data)
		}
	}
	return results
}

func getData(productCode string, tableName string, results []map[string]interface{}) []map[string]interface{} {
	var err error

	if productCode != "" {
		err = DB.Table(tableName).Where("product_code = ?", productCode).Find(&results).Error
		if err != nil {
			fmt.Println(err)
		}

	} else {
		err = DB.Table(tableName).Find(&results).Error
		if err != nil {
			fmt.Println(err)
		}
	}
	return results
}

func getYieldData(tableName string, results []map[string]interface{}) []map[string]interface{} {
	err := DB.Table(tableName).Find(&results).Error
	if err != nil {
		fmt.Println(err)
	}
	return results
}

func getStaticData(productCode string, tableName string, results []map[string]interface{}) []map[string]interface{} {
	DB.Table(tableName).Where("product_code = ?", productCode).Find(&results)
	return results
}

func GetRatingFactors() []models.RatingFactor {
	var factors []models.RatingFactor
	DB.Find(&factors)
	return factors
}

func GetDistinctCommissionProductCodes() ([]string, error) {
	var productCodes []string
	err := DB.Model(&models.ProductCommissionStructure{}).Distinct("product_code").Pluck("product_code", &productCodes).Error
	return productCodes, err
}

func GetGlobalTableData(tableName string) (map[string]interface{}, error) {
	var results = make(map[string]interface{})
	var count int64

	switch tableName {
	case "Yield Curve":
		var yieldCurves []models.YieldCurve
		DB.Find(&yieldCurves)
		results["data"] = yieldCurves
	case "Margins":
		var margins []models.ProductMargins
		DB.Find(&margins)
		DB.Model(&models.ProductMargins{}).Count(&count)
		results["data"] = margins
		results["row_count"] = count
	case "Parameters":
		var parameters []models.ProductParameters
		DB.Find(&parameters)
		DB.Model(&models.ProductParameters{}).Count(&count)
		results["data"] = parameters
		results["row_count"] = count
	case "Commission Structures":
		var commStructure []models.ProductCommissionStructure
		DB.Find(&commStructure)
		DB.Model(&models.ProductCommissionStructure{}).Count(&count)
		results["data"] = commStructure
		results["row_count"] = count
	case "Shocks":
		var shocks []models.ProductShock
		DB.Find(&shocks)
		DB.Model(&models.ProductShock{}).Count(&count)
		results["data"] = shocks
		results["row_count"] = count
	}

	return results, nil
}

func DeleteGlobalTableData(tableName, key string) (map[string]interface{}, error) {
	var results = make(map[string]interface{})
	var err error
	switch tableName {
	case "Yield Curve":
		err = DB.Exec("delete from yield_curve").Error
	case "Margins":
		if key != "" {
			err = DB.Exec("delete from product_margins where year = ?", key).Error
		} else {
			err = DB.Exec("delete from product_margins").Error
		}
	case "Parameters":
		if key != "" {
			err = DB.Exec("delete from product_parameters where year = ?", key).Error
		} else {
			err = DB.Exec("delete from product_parameters").Error
		}
	case "Shocks":
		if key != "" {
			err = DB.Exec("delete from product_shocks where shock_basis = ?", key).Error
		} else {
			err = DB.Exec("delete from product_shocks").Error
		}
	case "Commission Structures":
		if key != "" {
			err = DB.Exec("delete from product_commission_structures where product_code = ?", key).Error
		} else {
			err = DB.Exec("delete from product_commission_structures").Error
		}
	}

	return results, err
}
