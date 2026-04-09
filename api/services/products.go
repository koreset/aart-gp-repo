package services

import (
	"api/models"
	"api/utils"
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/dimchansky/utfbom"
	"github.com/iancoleman/strcase"
	"github.com/jszwec/csvutil"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

const (
	APPROVED = "approved"
	PENDING  = "pending"
)

type ProductService struct{}

func GetSelectedFeatures(features models.ProductFeatures) []models.BaseFeature {
	var selectFeatures []models.BaseFeature
	//var chosenFeatures map[string]bool
	members := reflect.ValueOf(&features).Elem()
	typeOfT := members.Type()

	for i := 0; i < members.NumField(); i++ {
		f := members.Field(i)
		if f.Type().Kind() == reflect.Bool && f.Interface() == true {
			featureName := strcase.ToScreamingSnake(typeOfT.Field(i).Name)
			fmt.Printf("%d: %s %s = %v\n", i, featureName, f.Type(), f.Interface())
			var feature models.BaseFeature
			DB.Where("name = ?", featureName).First(&feature)
			selectFeatures = append(selectFeatures, feature)
		}
	}

	return selectFeatures
}

func GetProductFamilies() ([]models.ProductFamily, error) {
	var productFamilies []models.ProductFamily
	err := DB.Find(&productFamilies).Error
	if err != nil {
		return nil, err
	}
	return productFamilies, nil
}

func GetProductsAndFamilies() ([]models.ProductFamily, error) {
	var productFamilies []models.ProductFamily
	err := DB.Preload("Products").
		Preload("Products.ProductFeatures").
		Preload("Products.ProductTransitionStates").
		Preload("Products.ProductTransitions").
		Preload("Products.ProductModelpointVariables").
		Preload("Products.GlobalTables").
		Preload("Products.ProductTables").
		Preload("Products.ProductPricingTables").
		Preload("Products.ProductParameters").
		Find(&productFamilies).Error
	if err != nil {
		fmt.Println(err)
	}
	//err := DB.Set("gorm:auto_preload", true).Find(&productFamilies).Error
	for _, family := range productFamilies {
		for i := range family.Products {
			pricingTables, err := GetPricingTables(family.Products[i].ProductPricingTables, family.Products[i].ProductCode)
			if err == nil {
				family.Products[i].ProductPricingTables = pricingTables
			}
		}
	}
	if err != nil {
		return productFamilies, err
	}
	return productFamilies, nil
}

func GetAllAvailableProducts() ([]models.Product, error) {
	var products []models.Product

	err := DB.Find(&products).Error

	if err != nil {
		return products, err
	}

	return products, nil
}

func GetMarkovStates() ([]models.MarkovState, error) {
	var states []models.MarkovState
	err := DB.Find(&states).Error
	if err != nil {
		return nil, err
	}
	return states, nil
}

func GetAssumptionVariables() []models.BaseAssumptionVariable {
	var aV []models.BaseAssumptionVariable
	err := DB.Find(&aV).Error
	if err != nil {
		return nil
	}
	return aV
}

func GetBaseFeatures() ([]models.BaseFeature, error) {
	var bf []models.BaseFeature
	err := DB.Find(&bf).Error
	if err != nil {
		return bf, err
	}
	return bf, nil
}

func GetProductById(id int) (prod models.Product, err error) {
	// Redis cache check
	cacheKey := fmt.Sprintf("ads:v2:product:id:%d", id)
	if RedisAvailable() {
		if bs, e := redisClient.Get(redisCtx, cacheKey).Bytes(); e == nil {
			if json.Unmarshal(bs, &prod) == nil {
				return prod, nil
			}
		}
	}

	if err = DB.Preload("ProductFeatures").
		Preload("ProductTransitionStates").
		Preload("ProductTransitions").
		Preload("ProductRatingFactors").
		Preload("ProductRatingFactors.Fds").
		Preload("ProductModelpointVariables").
		Preload("ProductTables").
		Preload("ProductPricingTables").
		Preload("GlobalTables").
		Preload("ProductParameters").Where("id = ?", id).Find(&prod).Error; err != nil {
		return
	}

	features := reflect.ValueOf(&prod.ProductFeatures).Elem()
	var featureList []string
	for i := 0; i < features.NumField(); i++ {
		if features.Field(i).Kind() == reflect.Bool && features.Field(i).Bool() {
			name := features.Type().Field(i).Name
			name = utils.Split(name)
			featureList = append(featureList, name)
		}
	}

	prodTables, err := CheckTablesForProduct(prod.ProductTables, prod.ProductCode)
	if err != nil {
		fmt.Println(err)
	} else {
		prod.ProductTables = prodTables
	}

	pricingTables, err := GetPricingTables(prod.ProductPricingTables, prod.ProductCode)
	if err != nil {
		fmt.Println(err)
	} else {
		prod.ProductPricingTables = pricingTables
	}

	prod.FeatureList = featureList

	// Set Redis cache
	if RedisAvailable() {
		if b, e := json.Marshal(&prod); e == nil {
			_ = redisClient.Set(redisCtx, cacheKey, b, 90*time.Minute).Err()
		}
	}

	return
}

//func GetProductByCode(code string) (prod models.Product, err error) {
//	err = DBReadWithResilience(context.Background(), func(db *gorm.DB) error {
//		return db.Preload("ProductFeatures").
//			Preload("ProductTransitionStates").
//			Preload("ProductTransitions").
//			Preload("ProductModelpointVariables").
//			Preload("ProductTables").
//			Preload("GlobalTables").
//			Preload("ProductParameters").
//			Where("product_code = ?", code).
//			Find(&prod).Error
//	})
//	if err != nil {
//		return prod, err
//	}
//
//	var pparams models.ProductParameters
//	err = DBReadWithResilience(context.Background(), func(db *gorm.DB) error {
//		return db.Where("product_code = ?", prod.ProductCode).First(&pparams).Error
//	})
//	if err != nil {
//		return prod, err
//	}
//	prod.ProductParameters = pparams
//
//	return prod, nil
//}

func GetProductByCode(code string) (prod models.Product, err error) {
	// Redis cache check
	cacheKey := fmt.Sprintf("ads:v1:product:code:%s", code)
	if RedisAvailable() {
		if bs, e := redisClient.Get(redisCtx, cacheKey).Bytes(); e == nil {
			if json.Unmarshal(bs, &prod) == nil {
				return prod, nil
			}
		}
	}

	err = DB.Preload("ProductFeatures").
		Preload("ProductTransitionStates").
		Preload("ProductTransitions").
		Preload("ProductModelpointVariables").
		Preload("ProductTables").
		Preload("GlobalTables").
		Preload("ProductParameters").
		Where("product_code = ?", code).
		Find(&prod).Error

	if err != nil {
		return prod, err
	}

	var pparams models.ProductParameters
	err = DB.Where("product_code = ?", prod.ProductCode).First(&pparams).Error

	if err != nil {
		return prod, err
	}
	prod.ProductParameters = pparams

	// Also need array of parameter fields for display - This should be conditional at some point.
	elems := reflect.ValueOf(&pparams).Elem()
	var paramArray []string
	for i := 0; i < elems.NumField(); i++ {
		name := elems.Type().Field(i).Name
		name = utils.Split(name)
		paramArray = append(paramArray, name)
	}
	if len(paramArray) > 2 && paramArray != nil {
		prod.ParametersArray = paramArray[2:]
	}

	// Set Redis cache
	if RedisAvailable() {
		if b, e := json.Marshal(&prod); e == nil {
			_ = redisClient.Set(redisCtx, cacheKey, b, 90*time.Minute).Err()
		}
	}

	return
}

func CheckTablesForProduct(tables []models.ProductTable, productCode string) ([]models.ProductTable, error) {

	var tableName string
	var count int64
	for i := range tables {
		switch tables[i].Class {
		case "TransitionRates":
			tableName = strings.ToLower(productCode) + "_" + strings.ToLower(tables[i].Name)
			count = 0
			DB.Table(tableName).Count(&count)
			if count > 0 {
				tables[i].Populated = true
			}
		case "Margins":
			tableName = "product_" + strings.ToLower(tables[i].Name)
			count = 0
			DB.Table(tableName).Where("product_code =?", productCode).Count(&count)
			if count > 0 {
				tables[i].Populated = true
			}
		case "Global":
			tableName = strings.ToLower(tables[i].Name)
			count = 0
			DB.Table(tableName).Count(&count)
			if count > 0 {
				tables[i].Populated = true
			}
		case "BenefitStructure":
			tableName = "product_" + strings.ToLower(tables[i].Name)
			count = 0
			DB.Table(tableName).Where("product_code =?", productCode).Count(&count)
			if count > 0 {
				tables[i].Populated = true
			}
		case "LapseMargins":
			tableName = "product_" + strings.ToLower(strcase.ToSnake(tables[i].Name))
			count = 0
			DB.Table(tableName).Where("product_code =?", productCode).Count(&count)
			if count > 0 {
				tables[i].Populated = true
			}

		case "Parameters":
			tableName = "product_" + strings.ToLower(tables[i].Name)
			count = 0
			DB.Table(tableName).Where("product_code =?", productCode).Count(&count)
			if count > 0 {
				tables[i].Populated = true
			}

		case "Distribution":
			tableName = "product_" + strings.ToLower(tables[i].Name)
			count = 0
			DB.Table(tableName).Where("product_code =?", productCode).Count(&count)
			if count > 0 {
				tables[i].Populated = true
			}
		case "Valuations":
			tableName = "product_" + strings.ToLower(tables[i].Name)
			count = 0
			//DB.Table(tableName).Where("product_code =?", productCode).Count(&count)
			DB.Table(tableName).Count(&count)

			if count > 0 {
				tables[i].Populated = true
			}

		}

	}

	return tables, nil
}

func GetProductTable(tableId int) models.ProductTable {
	var prodTable models.ProductTable
	err := DB.Where("id = ?", tableId).First(&prodTable).Error
	if err != nil {
		fmt.Println(err)
	}
	return prodTable
}

// validateCSVFile checks that the uploaded file is a non-empty CSV with at least one
// data row and that all requiredColumns are present in the header (case-insensitive).
func validateCSVFile(v *multipart.FileHeader, requiredColumns []string) error {
	if v.Size == 0 {
		return fmt.Errorf("uploaded file is empty")
	}

	f, err := v.Open()
	if err != nil {
		return fmt.Errorf("could not open uploaded file: %w", err)
	}
	defer f.Close()

	var delimiter rune = ','
	// Detect delimiter via a fresh read
	delimF, err2 := v.Open()
	if err2 == nil {
		delimiter, _ = utils.GetDelimiter(delimF)
		delimF.Close()
	}

	reader := csv.NewReader(f)
	reader.TrimLeadingSpace = true
	reader.Comma = delimiter

	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("could not read CSV header — file may be empty or malformed")
	}
	if len(header) == 0 {
		return fmt.Errorf("CSV file has no columns in the header row")
	}

	// Normalise header names for comparison
	headerSet := make(map[string]bool, len(header))
	for _, h := range header {
		headerSet[strings.ToLower(strings.TrimSpace(h))] = true
	}

	var missing []string
	for _, col := range requiredColumns {
		if !headerSet[strings.ToLower(col)] {
			missing = append(missing, col)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("CSV is missing required columns: %s", strings.Join(missing, ", "))
	}

	// Confirm at least one data row exists
	_, err = reader.Read()
	if err == io.EOF {
		return fmt.Errorf("CSV file has a header row but no data rows")
	}
	if err != nil {
		return fmt.Errorf("error reading CSV data: %w", err)
	}

	return nil
}

func ProcessTable(ctx context.Context,
	pt models.ProductTable,
	prodCode string,
	year int,
	file *multipart.FileHeader,
	month int, yieldCurveCode string) error {

	// Validate CSV structure before processing
	var requiredCols []string
	switch pt.Class {
	case "Margins":
		requiredCols = []string{"product_code", "basis"}
	case "Parameters":
		requiredCols = []string{"product_code", "basis"}
	case "TransitionRates":
		requiredCols = nil // header column count is checked implicitly
	default:
		requiredCols = nil
	}
	if err := validateCSVFile(file, requiredCols); err != nil {
		return err
	}

	var delimiter rune
	delimiterFile, err := file.Open()
	if err != nil {
		return err
	}
	defer delimiterFile.Close()
	delimiter, err = utils.GetDelimiter(delimiterFile)

	switch pt.Class {
	case "Margins":
		return processMarginTables(ctx, file, prodCode, year)
	case "TransitionRates":
		return processRatingTables(DB, file, delimiter, prodCode, year)
	case "Parameters":
		return processParameterTables(ctx, file, delimiter, prodCode, year)
	case "Distribution":
		return processStaticTables(pt.Name, delimiter, file, prodCode, year)
	case "BenefitStructure":
		return processStaticTables(pt.Name, delimiter, file, prodCode, year)
	case "LapseMargins":
		return processStaticTables(pt.Name, delimiter, file, prodCode, year)
	case "Valuations":
		return processStaticTables(pt.Name, delimiter, file, prodCode, year)
	case "Profitabilities":
		return processStaticTables(pt.Name, delimiter, file, prodCode, year)
	case "Global":
		return processYieldCurve(ctx, file, prodCode, year, month, yieldCurveCode)
	}
	return nil
}

func ProcessPricingTable(ctx context.Context,
	tableId int,
	prodCode string,
	file *multipart.FileHeader,
	month int, yieldCurveCode string) error {
	var delimiter rune
	delimiterFile, err := file.Open()
	if err != nil {
		return err
	}
	defer delimiterFile.Close()
	delimiter, err = utils.GetDelimiter(delimiterFile)
	pt := GetProductPricingTable(tableId)

	switch pt.Class {
	case "Margins":
		return processPricingMarginTables(ctx, file, delimiter, prodCode)
	case "TransitionRates":
		return processPricingRatingTables(DB, file, delimiter, prodCode)
	case "Parameters":
		return processPricingParameterTables(ctx, file, delimiter, prodCode)
	case "Distribution":
		return processPricingStaticTables(pt.Name, delimiter, file, prodCode)
	case "Profitabilities":
		return processPricingStaticTables(pt.Name, delimiter, file, prodCode)
	case "BenefitStructure":
		return processPricingStaticTables(pt.Name, delimiter, file, prodCode)
	case "LapseMargins":
		return processPricingStaticTables(pt.Name, delimiter, file, prodCode)
	case "Valuations":
		return processPricingStaticTables(pt.Name, delimiter, file, prodCode)
	case "Global":
		return processPricingYieldCurve(ctx, file, prodCode, month, yieldCurveCode)
	}
	return nil
}

func ProcessPricingParameters(productCode string, v *multipart.FileHeader) error {
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
	} else {
		reader := csv.NewReader(file)
		reader.TrimLeadingSpace = true
		reader.Comma = delimiter
		dec, _ := csvutil.NewDecoder(reader)

		dec.Header()

		for {
			var ra models.PricingParameter
			if err = dec.Decode(&ra); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			//ra.Year = year
			err = DB.Where("product_code = ? and basis = ?", ra.ProductCode, ra.Basis).Delete(&models.PricingParameter{}).Error
			err = DB.Create(&ra).Error
			if err != nil {
				log.Error().Err(err)
			}
		}
	}

	return nil
}

func ProcessPricingDemographics(productCode string, v *multipart.FileHeader) error {
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
	} else {
		reader := csv.NewReader(file)
		reader.TrimLeadingSpace = true
		reader.Comma = delimiter
		dec, _ := csvutil.NewDecoder(reader)

		dec.Header()

		err = DB.Where("product_code = ?", productCode).Delete(&models.PricingPolicyDemographic{}).Error
		for {
			var ra models.PricingPolicyDemographic
			if err = dec.Decode(&ra); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			ra.ProductCode = productCode

			err = DB.Create(&ra).Error
			if err != nil {
				log.Error().Err(err)
			}
		}
	}

	return nil
}

func ProcessBulkPricingParameters(v *multipart.FileHeader) error {
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
	} else {
		reader := csv.NewReader(file)
		reader.TrimLeadingSpace = true
		reader.Comma = delimiter
		dec, _ := csvutil.NewDecoder(reader)

		dec.Header()

		for {
			var ra models.PricingParameter
			if err := dec.Decode(&ra); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			err := DB.Where("product_code = ?", ra.ProductCode).Save(&ra).Error
			if err != nil {
				log.Error().Err(err)
			}
		}
	}

	return nil
}

func ProcessTables(db *gorm.DB, pc *models.Product, files []*multipart.FileHeader) error {
	//Build a table map
	var err error
	tableMap := map[string]string{}
	for _, v := range pc.ProductTables {
		tableMap[v.Name] = v.Class
	}
	for _, file := range files {
		var delimiter rune
		if tableMap[file.Filename[:len(file.Filename)-4]] == "TransitionRates" {

			delimiterFile, err := file.Open()
			if err != nil {
				return err
			}
			defer delimiterFile.Close()
			delimiter, err = utils.GetDelimiter(delimiterFile)

			date := time.Now()
			err = processRatingTables(db, file, delimiter, pc.ProductCode, date.Year())
			if err != nil {
				return err
			}
		}
	}

	//We need to create a special Yield Curve table per product that pricing would use.
	createPricingYieldCurve(db, pc.ProductCode)
	return err
}

func ProcessFeatures(db *gorm.DB, prodCode string, featureSet []models.BaseFeature) error {
	prodFeatures := models.ProductFeatures{}
	prodFeatures.ProductCode = prodCode
	p := reflect.ValueOf(&prodFeatures).Elem()
	for _, feature := range featureSet {
		newName := structify(feature.Name)
		if len(newName) > 0 {
			if p.FieldByName(newName).IsValid() {
				p.FieldByName(newName).SetBool(true)
			} else {
				return errors.New("an invalid feature type has been provided: " + newName)
			}
		}
	}

	// If a row already exists for this product code (edit mode), carry its ID
	// so GORM issues an UPDATE rather than an INSERT, avoiding the unique
	// constraint violation on product_code.
	var existing models.ProductFeatures
	if err := DB.Where("product_code = ?", prodCode).First(&existing).Error; err == nil {
		prodFeatures.ID = existing.ID
	}

	if err := DB.Save(&prodFeatures).Error; err != nil {
		return err
	}
	return nil
}

func structify(name string) string {
	caser := cases.Title(language.English)
	name = strings.Replace(name, "_", " ", -1)
	name = caser.String(strings.ToLower(name))
	name = strings.Replace(name, " ", "", -1)
	return name
}

func processSpecialTables(table string, delimiter rune, v *multipart.FileHeader, prodCode string, year int) error {
	fmt.Println("Processing Special Table", table)

	file, err := v.Open()
	defer file.Close()
	if err != nil {
		return err
	} else {
		reader := csv.NewReader(file)
		reader.TrimLeadingSpace = true
		reader.Comma = delimiter
		dec, _ := csvutil.NewDecoder(reader)
		dec.Header()
		specialCount := 0
		renewableCount := 0

		for {
			if table == "product_special_decrement_margins" {
				var pp models.ProductSpecialDecrementMargin
				if fileErr := dec.Decode(&pp); fileErr == io.EOF {
					break
				} else if fileErr != nil {
					fmt.Println(fileErr)
				}

				pp.ProductCode = prodCode
				if specialCount == 0 {
					DB.Where("special_margin_code = ? and product_code = ?", pp.SpecialMarginCode, pp.ProductCode).Delete(&pp)
				}
				err = DB.Create(&pp).Error
				if err != nil {
					log.Error().Err(err)
				}
				specialCount++
			}

			if table == "product_renewable_profit_adjustments" {
				var pp models.ProductRenewableProfitAdjustment

				if fileErr := dec.Decode(&pp); fileErr == io.EOF {
					break
				} else if fileErr != nil {
					fmt.Println(fileErr)
				}

				pp.ProductCode = prodCode
				if renewableCount == 0 {
					DB.Where("profit_adjustment_code = ? and product_code = ?", pp.ProfitAdjustmentCode, pp.ProductCode).Delete(&pp)
				}
				err = DB.Create(&pp).Error
				if err != nil {
					log.Error().Err(err)
				}
				renewableCount++
			}
		}
		//Register this activity...
		//description := "Uploaded Parameter table for " + prodCode
		//CreateActivity(ctx, "table_upload", description, 0, "product_parameters")

		return err
	}
}

func processStaticTables(table string, delimiter rune, v *multipart.FileHeader, prodCode string, year int) error {

	tableName := "product_" + strings.ToLower(strcase.ToSnake(table))

	if strings.Contains(tableName, "special_decrement_margins") || strings.Contains(tableName, "renewable_profit_adjustments") {
		return processSpecialTables(tableName, delimiter, v, prodCode, year)
	}

	file, err := v.Open()
	if err != nil {
		return fmt.Errorf("could not open uploaded file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = delimiter
	reader.TrimLeadingSpace = true
	data, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV data: %w", err)
	}

	var headerRow string
	var missingProdCode bool
	var missingYear bool

	if !strings.Contains(tableName, "special_decrement_margins") || !strings.Contains(tableName, "renewable_profit_adjustments") {
		DB.Exec(fmt.Sprintf("delete from %s where product_code='%s' and year = %d", tableName, prodCode, year))
	}

	for i, row := range data {
		if i == 0 {
			for i, item := range row {
				item = strings.Replace(strings.TrimSpace(strings.ToLower(item)), " ", "_", -1)
				headerRow += item
				if i < len(row)-1 {
					headerRow += ","
				}
			}

			if !strings.Contains(headerRow, "product_code") {
				headerRow += ",product_code"
				missingProdCode = true
			}

			if !strings.Contains(headerRow, "year") && year > 0 {
				headerRow += ",year"
				missingYear = true
			}

			fmt.Println("Header:", headerRow)
		} else {
			err = AddToStaticTable(tableName, row, headerRow, prodCode, year, missingProdCode, missingYear)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func processPricingStaticTables(table string, delimiter rune, v *multipart.FileHeader, prodCode string) error {
	tableName := "product_pricing_" + strings.ToLower(strcase.ToSnake(table))

	file, err := v.Open()
	if err != nil {
		fmt.Println(err)
	} else {
		reader := csv.NewReader(file)
		reader.Comma = delimiter
		reader.TrimLeadingSpace = true
		//dec, _ := csvutil.NewDecoder(reader)
		//dec.Header()
		data, _ := reader.ReadAll()
		var headerRow string
		var missingProdCode bool

		DB.Exec(fmt.Sprintf("delete from %s where product_code='%s'", tableName, prodCode))

		for i, row := range data {
			if i == 0 {
				for i, item := range row {
					item = strings.Replace(strings.TrimSpace(strings.ToLower(item)), " ", "_", -1)
					headerRow += item
					if i < len(row)-1 {
						headerRow += ","
					}
				}

				if !strings.Contains(headerRow, "product_code") {
					headerRow += ",product_code"
					missingProdCode = true
				}

				fmt.Println("Header:", headerRow)
			} else {
				err = AddToPricingStaticTable(tableName, row, headerRow, prodCode, missingProdCode)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func processParameterTables(ctx context.Context, v *multipart.FileHeader, delimiter rune, prodCode string, year int) error {
	file, err := v.Open()
	if err != nil {
		return err
	} else {
		reader := csv.NewReader(file)
		reader.TrimLeadingSpace = true
		reader.Comma = delimiter
		dec, _ := csvutil.NewDecoder(reader)
		dec.Header()

		for {
			if year != 0 {
				var pp models.ProductParameters
				if err := dec.Decode(&pp); err == io.EOF {
					break
				} else if err != nil {
					fmt.Println(err)
					return err
				}
				pp.Year = year
				pp.ProductCode = prodCode
				err := DB.Table("product_parameters").Create(&pp).Error
				if err != nil {
					log.Error().Err(err)
					return err
				}
			}

			if year == 0 {
				var pp models.ProductPricingParameters
				if err := dec.Decode(&pp); err == io.EOF {
					break
				} else if err != nil {
					fmt.Println(err)
					return err
				}
				if err := DB.Table("product_pricing_parameters").Save(&pp).Error; err != nil {
					return err
				}
			}
		}
	}

	//Register this activity...
	description := "Uploaded Parameter table for " + prodCode
	CreateActivity(ctx, "table_upload", description, 0, "product_parameters")
	return nil
}

func processPricingParameterTables(ctx context.Context, v *multipart.FileHeader, delimiter rune, prodCode string) error {
	file, err := v.Open()
	if err != nil {
		return err
	} else {
		reader := csv.NewReader(file)
		reader.TrimLeadingSpace = true
		reader.Comma = delimiter
		dec, _ := csvutil.NewDecoder(reader)
		dec.Header()

		for {
			var pp models.ProductPricingParameters
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
				return err
			}
			if err := DB.Table("product_pricing_parameters").Save(&pp).Error; err != nil {
				return err
			}
		}
	}

	//Register this activity...
	description := "Uploaded Pricing Parameter table for " + prodCode
	CreateActivity(ctx, "table_upload", description, 0, "product_parameters")
	return nil
}

func ProcessRiskAdjustmentFactorTables(ctx context.Context, v *multipart.FileHeader, year int, version string) error {
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
	} else {
		reader := csv.NewReader(file)
		reader.TrimLeadingSpace = true
		reader.Comma = delimiter
		dec, _ := csvutil.NewDecoder(reader)

		dec.Header()

		for {
			var ra models.RiskAdjustmentFactor
			if err := dec.Decode(&ra); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
				return err
			}
			ra.Year = year
			ra.Version = version
			ra.CreatedBy = ctx.Value("userName").(string)
			err := DB.Create(&ra).Error
			if err != nil {
				log.Error().Err(err)
				return err
			}
		}
	}
	return nil
}

func ProcessSAPFile(ctx context.Context, v *multipart.FileHeader, user models.AppUser) error {
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
	} else {
		reader := csv.NewReader(file)
		reader.TrimLeadingSpace = true
		reader.Comma = delimiter
		dec, _ := csvutil.NewDecoder(reader)

		dec.Header()

		var saps []models.ManualScopedAggregatedProjection

		for {
			var sap models.ManualScopedAggregatedProjection
			if err := dec.Decode(&sap); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			sap.User = user.UserName
			sap.UserEmail = user.UserEmail
			originalTime := time.Now()
			year, month, day := originalTime.Date()
			hour, minute, _ := originalTime.Clock()
			sap.CreatedAt = time.Date(year, month, day, hour, minute, 0, 0, originalTime.Location())
			DB.Where("run_name = ?", sap.RunName).Delete(&sap)
			saps = append(saps, sap)
		}

		err := DB.CreateInBatches(&saps, 200).Error
		if err != nil {
			log.Error().Err(err)
			return err
		}
	}
	return nil
}

//func ProcessPAAFinance(ctx context.Context, v *multipart.FileHeader, year int, version string) error {
//	var delimiter rune
//	delimiterFile, err := v.Open()
//	if err != nil {
//		return fmt.Errorf("failed to open file for delimiter detection: %w", err)
//	}
//	defer delimiterFile.Close()
//
//	delimiter, err = utils.GetDelimiter(delimiterFile)
//	if err != nil {
//		return fmt.Errorf("failed to get delimiter: %w", err)
//	}
//
//	file, err := v.Open()
//	if err != nil {
//		return fmt.Errorf("failed to open file for processing: %w", err)
//	}
//	defer file.Close()
//
//	reader := csv.NewReader(file)
//	reader.TrimLeadingSpace = true
//	reader.Comma = delimiter
//
//	// It's good practice to read the header separately to have access to it
//	// in case of errors before the decoder is fully initialized or if dec.Header() fails.
//	header, err := reader.Read()
//	if err != nil {
//		if err == io.EOF {
//			return fmt.Errorf("CSV file is empty or only contains a header")
//		}
//		return fmt.Errorf("failed to read CSV header: %w", err)
//	}
//
//	// Re-open the file and create a new reader to reset its state for the decoder
//	// Alternatively, you could seek the file back to the beginning if it supports it.
//	file.Seek(0, io.SeekStart) // Seek back to the beginning of the file
//	reader = csv.NewReader(file)
//	reader.TrimLeadingSpace = true
//	reader.Comma = delimiter
//
//	dec, err := csvutil.NewDecoder(reader, header...) // Pass the read header to the decoder
//	if err != nil {
//		return fmt.Errorf("failed to create CSV decoder: %w", err)
//	}
//
//	// dec.Header() is implicitly called by NewDecoder if header names are passed.
//	// If you don't pass header names to NewDecoder, you would call dec.Header() here.
//
//	var rowNum = 1 // Start with row 1 (after header)
//
//	for {
//		rowNum++ // Increment row number for each data row
//		var fv models.PAAFinance
//		if err := dec.Decode(&fv); err != nil {
//			if err == io.EOF {
//				break
//			}
//
//			// Check if the error is a csvutil.DecodeError
//			if decodeErr, ok := err.(*csvutil.DecodeError); ok {
//				//fieldName := decodeErr.Field // Name of the struct field
//				//columnIndex := decodeErr.Column
//				//columnName := fieldName // By default, use struct field name
//
//				// If you have the header slice, you can get the actual CSV header name
//				//if columnIndex >= 0 && columnIndex < len(header) {
//				//	columnName = header[columnIndex]
//				//}
//
//				// Get the raw field that caused the error
//				// The decoder internally reads a record (slice of strings) before decoding.
//				// We need to get the specific field value.
//				// This part can be tricky as csvutil doesn't directly expose the raw field value
//				// that caused the specific DecodeError easily.
//				// However, we can provide context with the column name and error type.
//
//				underLyingErr := decodeErr.Err
//				fmt.Println(underLyingErr.Error())
//				//return fmt.Errorf("error decoding %s on row %d, column '%s'. Please check the data in this column",
//				//	v.Filename, decodeErr.Line, columnName)
//			}
//			// Handle other types of errors
//			return fmt.Errorf("error decoding CSV on row %d: %w", rowNum, err)
//		}
//
//		fv.Year = year
//		fv.RiskRateCode = version
//
//		// Consider using a transaction for Delete and Save operations
//		// tx := DB.Begin()
//		// if err := tx.Where("year = ? and portfolio_name = ? and product_code = ? and ifrs17_group=? and version = ?", fv.Year, fv.PortfolioName, fv.ProductCode, fv.IFRS17Group, fv.RiskRateCode).Delete(&fv).Error; err != nil {
//		// 	tx.Rollback()
//		// 	return fmt.Errorf("failed to delete existing record for row %d: %w", rowNum, err)
//		// }
//		// if err := tx.Save(&fv).Error; err != nil {
//		// 	tx.Rollback()
//		// 	return fmt.Errorf("failed to save record for row %d: %w", rowNum, err)
//		// }
//		// if err := tx.Commit().Error; err != nil {
//		//      return fmt.Errorf("failed to commit transaction for row %d: %w", rowNum, err)
//		// }
//
//		// Original DB operations (ensure error handling for these as well)
//		if err := DB.Where("year = ? and portfolio_name = ? and product_code = ? and ifrs17_group=? and version = ?", fv.Year, fv.PortfolioName, fv.ProductCode, fv.IFRS17Group, fv.RiskRateCode).Delete(&fv).Error; err != nil {
//			return fmt.Errorf("failed to delete existing record for row %d (Portfolio: %s, Product: %s): %w", rowNum, fv.PortfolioName, fv.ProductCode, err)
//		}
//		if err := DB.Where("year = ? and portfolio_name = ? and product_code = ? and ifrs17_group=? and version = ?", fv.Year, fv.PortfolioName, fv.ProductCode, fv.IFRS17Group, fv.RiskRateCode).Save(&fv).Error; err != nil { // Note: Using .Save on a .Where might not be what you intend if you want to ensure it's an update or create.
//			return fmt.Errorf("failed to save record for row %d (Portfolio: %s, Product: %s): %w", rowNum, fv.PortfolioName, fv.ProductCode, err)
//		}
//		// A more common pattern for upsert might be:
//		// DB.Clauses(clause.OnConflict{
//		// Columns:   []clause.Column{{Name: "year"}, {Name: "portfolio_name"}, {Name: "product_code"}, {Name: "ifrs17_group"}, {Name: "version"}},
//		// DoUpdates: clause.AssignmentColumns([]string{"column_to_update1", "column_to_update2"}), // list columns you want to update on conflict
//		// }).Create(&fv)
//
//	}
//	return nil
//}

func ProcessPAAFinance(ctx context.Context, v *multipart.FileHeader, year int, version string) error {
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
	} else {
		reader := csv.NewReader(file)
		reader.TrimLeadingSpace = true
		reader.Comma = delimiter
		dec, _ := csvutil.NewDecoder(reader)

		dec.Header()

		for {
			var fv models.PAAFinance
			if err := dec.Decode(&fv); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
				return err
			}
			fv.Year = year
			fv.Version = version
			fv.CreatedBy = ctx.Value("userName").(string)
			fv.Created = time.Now().Unix()

			DB.Where("year = ? and portfolio_name = ? and product_code = ? and ifrs17_group=? and version = ?", fv.Year, fv.PortfolioName, fv.ProductCode, fv.IFRS17Group, fv.Version).Delete(&fv)
			err = DB.Where("year = ? and portfolio_name = ? and product_code = ? and ifrs17_group=? and version = ?", fv.Year, fv.PortfolioName, fv.ProductCode, fv.IFRS17Group, fv.Version).Save(&fv).Error
			if err != nil {
				log.Error().Err(err)
			}

		}
	}
	return nil
}

func ProcessTransitionAdjustments(ctx context.Context, v *multipart.FileHeader, year int, version string) error {
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
	} else {
		reader := csv.NewReader(utfbom.SkipOnly(file))
		reader.TrimLeadingSpace = true
		reader.Comma = delimiter
		dec, _ := csvutil.NewDecoder(reader)

		dec.Header()
		var tas []models.BalanceSheetRecord

		for {
			var ta models.BalanceSheetRecord
			if err := dec.Decode(&ta); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			ta.MeasurementType = "TransitionAdjustment"
			ta.Year = year
			ta.Version = version
			tas = append(tas, ta)
		}
		DB.Where("measurement_type = ? and date = ? and ifrs17_group = ?", tas[0].MeasurementType, tas[0].Date, tas[0].IFRS17Group).Delete(&models.BalanceSheetRecord{})
		err := DB.CreateInBatches(&tas, 200).Error
		if err != nil {
			log.Error().Err(err)
			return err
		}
	}
	return nil
}

func ProcessFinanceVariables(ctx context.Context, v *multipart.FileHeader, year int, version string) error {
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
	} else {
		reader := csv.NewReader(file)
		reader.TrimLeadingSpace = true
		reader.Comma = delimiter
		dec, _ := csvutil.NewDecoder(reader)

		dec.Header()

		for {
			var fv models.FinanceVariables
			if err := dec.Decode(&fv); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			fv.Year = year
			fv.Version = version
			fv.CreatedBy = ctx.Value("userName").(string)
			DB.Where("product_code = ? and ifrs17_group=? and year = ? and version = ?", fv.ProductCode, fv.IFRS17Group, fv.Year, fv.Version).Delete(&fv)
			//DB.Where("product_code = ? and ifrs17_group=? and year = ?", fv.ProductCode, fv.IFRS17Group, fv.Year).Save(&fv)
			err = DB.Create(&fv).Error
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func GetFinanceFile() ([]models.FinanceVariables, error) {
	var finVars []models.FinanceVariables
	err := DB.Find(&finVars).Error
	return finVars, err
}

func GetSapFileList() []models.ManualSapList {
	var saps []models.ManualSapList
	err := DB.Table("manual_scoped_aggregated_projections").Distinct("run_name").Select("run_name,run_date,created_at, user, user_email").Find(&saps).Error

	if err != nil {
		log.Error().Err(err)
	}
	return saps
}

func GetSapResultsForRunName(runName string) ([]models.ManualScopedAggregatedProjection, error) {
	var saps []models.ManualScopedAggregatedProjection
	err := DB.Where("run_name = ?", runName).Find(&saps).Error
	if err != nil {
		log.Error().Err(err)
	}
	return saps, err
}

func DeleteSapData(runName string) error {
	err := DB.Where("run_name = ?", runName).Delete(models.ManualScopedAggregatedProjection{}).Error
	if err != nil {
		log.Error().Err(err)
	}
	err = DB.Where("run_name =''").Delete(models.ManualScopedAggregatedProjection{}).Error
	return err
}

func GetAvailablePaaFinanceVersions(year string) ([]string, error) {
	var versions []string
	err := DB.Table("paa_finances").Where("year = ?", year).Distinct("version").Pluck("version", &versions).Error
	if err != nil {
		log.Error().Err(err)
	}
	return versions, err
}

func GetAvailableGmmFinanceVersions(year string) ([]string, error) {
	var versions []string
	err := DB.Table("finance_variables").Where("year = ?", year).Distinct("version").Pluck("version", &versions).Error
	if err != nil {
		log.Error().Err(err)
	}
	return versions, err

}

func GetFinanceAndRaYears(measure, paaRunId string) (map[string]interface{}, error) {
	years := make(map[string]interface{})
	paaId, err := strconv.Atoi(paaRunId)

	var financeYears []int
	if measure == "GMM" {
		err := DB.Table("finance_variables").Select("year").Distinct().Scan(&financeYears).Error
		if err != nil {
			return nil, err
		}
	}

	if measure == "VFA" {
		err := DB.Table("finance_variables").Select("year").Distinct().Scan(&financeYears).Error
		if err != nil {
			return nil, err
		}
	}

	if measure == "PAA" {
		var mgmm models.ModifiedGMMScopedAggregation
		DB.Where("run_id = ?", paaId).First(&mgmm)

		//err = DB.Table("paa_finances").Where("portfolio_name = ?", mgmm.ExpConfigurationName).Select("year").Distinct().Scan(&financeYears).Error
		err = DB.Table("paa_finances").Select("year").Distinct().Scan(&financeYears).Error
		if err != nil {
			return nil, err
		}
	}

	var raYears []int
	err = DB.Table("risk_adjustment_factors").Select("year").Distinct().Scan(&raYears).Error
	if err != nil {
		return nil, err
	}

	years["finance"] = financeYears
	years["ra"] = raYears

	return years, nil
}

func ProcessBulkParameterTables(ctx context.Context, v *multipart.FileHeader, tableName, year, yieldCurveCode, month string) error {
	var delimiter rune
	yearInt, _ := strconv.Atoi(year)
	monthInt, _ := strconv.Atoi(month)
	delimiterFile, err := v.Open()
	if err != nil {
		return err
	}
	defer delimiterFile.Close()
	delimiter, err = utils.GetDelimiter(delimiterFile)

	//var tableName = "product_" + strings.ToLower(v.Filename[:len(v.Filename)-4])

	file, err := v.Open()
	if err != nil {
		return err
	} else {
		reader := csv.NewReader(file)
		reader.TrimLeadingSpace = true
		reader.Comma = delimiter
		dec, _ := csvutil.NewDecoder(reader)

		dec.Header()

		if tableName == "product_parameters" {
			for {
				var pp models.ProductParameters
				if err := dec.Decode(&pp); err == io.EOF {
					break
				} else if err != nil {
					fmt.Println(err)
					return err
				}
				//pp.Year = yearInt
				pp.CreatedBy = ctx.Value("userName").(string)
				if err := DB.Where("year = ? and product_code = ? and yield_curve_code = ? and basis = ?", pp.Year, pp.ProductCode, pp.YieldCurveCode, pp.Basis).Delete(&pp).Error; err != nil {
					return err
				}
				err = DB.Create(&pp).Error
				if err != nil {
					fmt.Println(err)
					return err
				}
			}
		}

		if tableName == "product_margins" {
			for {
				var pm models.ProductMargins
				if err := dec.Decode(&pm); err == io.EOF {
					break
				} else if err != nil {
					fmt.Println(err)
					return err
				}
				//pm.Year = yearInt
				if err := DB.Where("year = ? and product_code = ? and basis = ?", pm.Year, pm.ProductCode, pm.Basis).Delete(&pm).Error; err != nil {
					return err
				}
				if err := DB.Create(&pm).Error; err != nil {
					return err
				}
			}
		}

		if tableName == "yield_curve" {
			for {
				var pm models.YieldCurve
				if err := dec.Decode(&pm); err == io.EOF {
					break
				} else if err != nil {
					fmt.Println(err)
					return err
				}
				pm.Year = yearInt
				pm.Month = monthInt
				pm.YieldCurveCode = yieldCurveCode
				err = DB.Where("year =? and yield_curve_code = ? and proj_time = ? and month = ?", pm.Year, pm.YieldCurveCode, pm.ProjectionTime, pm.Month).Delete(&pm).Error
				if err != nil {
					fmt.Println(err)
					return err
				}
				err = DB.Create(&pm).Error
				if err != nil {
					fmt.Println(err)
					return err
				}
			}
		}

		if tableName == "product_commission_structures" {
			for {
				var pm models.ProductCommissionStructure
				if err := dec.Decode(&pm); err == io.EOF {
					break
				} else if err != nil {
					fmt.Println(err)
					return err
				}
				//pm.Year = yearInt
				//err := DB.Where("year = ? and projection_month = ? and shock_basis = ?", pm.Year, pm.ProjectionMonth, pm.ShockBasis).Save(&pm).Error
				//err := DB.Where("projection_month = ? and shock_basis = ?", pm.ProjectionMonth, pm.ShockBasis).Save(&pm).Error
				pm.CreatedBy = ctx.Value("userName").(string)
				err := DB.Where("commission_type = ? and product_code = ?", pm.CommissionType, pm.ProductCode).Delete(&pm).Error
				if err != nil {
					return err
				}
				err = DB.Create(&pm).Error
				if err != nil {
					return err
				}
			}
		}

		if tableName == "product_shocks" {
			for {
				var pm models.ProductShock
				if err := dec.Decode(&pm); err == io.EOF {
					break
				} else if err != nil {
					fmt.Println(err)
					return err
				}
				//pm.Year = yearInt
				//err := DB.Where("year = ? and projection_month = ? and shock_basis = ?", pm.Year, pm.ProjectionMonth, pm.ShockBasis).Save(&pm).Error
				//err := DB.Where("projection_month = ? and shock_basis = ?", pm.ProjectionMonth, pm.ShockBasis).Save(&pm).Error
				pm.CreatedBy = ctx.Value("userName").(string)
				err := DB.Where("shock_basis = ?", pm.ShockBasis).Delete(&pm).Error
				if err != nil {
					return err
				}
				err = DB.Create(&pm).Error
				if err != nil {
					fmt.Println(err)
					return err
				}
			}
		}
	}

	//Register this activity...
	return nil
}

func processMarginTables(ctx context.Context, v *multipart.FileHeader, prodCode string, year int) error {
	var tableName string
	var err error

	if year == 0 {
		tableName = "product_pricing_" + strings.ToLower(v.Filename[:len(v.Filename)-4])
	} else {
		tableName = "product_" + strings.ToLower(v.Filename[:len(v.Filename)-4])
	}

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

	if tableName == "product_margins" {
		for {
			var pm models.ProductMargins
			if err := dec.Decode(&pm); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
				return err
			}
			pm.Year = year
			if err := DB.Where("year = ? and product_code = ? and basis = ?", year, pm.ProductCode, pm.Basis).Delete(&models.ProductMargins{}).Error; err != nil {
				return err
			}
			err = DB.Create(&pm).Error
			if err != nil {
				log.Error().Msg(err.Error())
				return err
			}
		}
	}

	if tableName == "product_pricing_margins" {
		for {
			var pm models.ProductPricingMargins
			if err := dec.Decode(&pm); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
				return err
			}

			if err := DB.Where("product_code = ?", pm.ProductCode).Delete(&models.ProductPricingMargins{}).Error; err != nil {
				return err
			}
			err = DB.Create(&pm).Error
			if err != nil {
				log.Error().Msg(err.Error())
				return err
			}
		}
	}

	//Register this activity...
	description := "Uploaded Product Margins for " + prodCode
	CreateActivity(ctx, "table_upload", description, 0, "product_margins")

	return nil
}

func processPricingMarginTables(ctx context.Context, v *multipart.FileHeader, delimiter rune, prodCode string) error {
	//var delimiter rune
	//delimiterFile, err := v.Open()
	//if err != nil {
	//	return err
	//}
	//defer delimiterFile.Close()
	//delimiter, err = utils.GetDelimiter(delimiterFile)

	file, err := v.Open()

	if err != nil {
		return err
	} else {
		reader := csv.NewReader(file)
		reader.TrimLeadingSpace = true
		reader.Comma = delimiter
		dec, _ := csvutil.NewDecoder(reader)

		dec.Header()

		for {
			var ppm models.ProductPricingMargins
			if err = dec.Decode(&ppm); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
				return err
			}
			//ppm.Year = year
			ppm.ProductCode = prodCode
			err = DB.Where("product_code = ?", ppm.ProductCode).Delete(&models.ProductPricingMargins{}).Error
			if err != nil {
				return err
			}
			err = DB.Create(&ppm).Error
			if err != nil {
				log.Error().Err(err)
				return err
			}
		}
	}

	return nil

	//var tableName string
	//var err error
	//
	//tableName = "product_pricing_margins"
	//
	//if err != nil {
	//	return err
	//} else {
	//	file, err := v.Open()
	//	if err != nil {
	//		return err
	//	}
	//	reader := csv.NewReader(file)
	//	reader.TrimLeadingSpace = true
	//	reader.Comma = delimiter
	//	data, _ := reader.ReadAll()
	//	fields := ""
	//	for i, row := range data {
	//		if i == 0 {
	//			fmt.Println("Header:", row)
	//			for _, field := range row {
	//				field = strings.TrimSpace(field)
	//				field = strings.ReplaceAll(field, "+", "_plus")
	//				fields += strings.ToLower(strings.Replace(field, " ", "_", -1)) + ","
	//			}
	//
	//		} else {
	//			err = AddToPricingMarginTable(tableName, row, fields, prodCode)
	//			if err != nil {
	//				return err
	//			}
	//		}
	//	}
	//}
	//
	//description := "Uploaded Product Margins for " + prodCode
	//CreateActivity(ctx, "table_upload", description, 0, "product_margins")
	//
	//return nil
}

func createPricingYieldCurve(db *gorm.DB, prodCode string) error {
	var tableName = strings.ToLower(prodCode) + "_pricing_yield_curve"

	tableQuery := fmt.Sprintf("CREATE TABLE `%s` (`proj_time` int NOT NULL, `nominal_rate` double DEFAULT NULL,  `real_rate` double DEFAULT NULL, `inflation` double DEFAULT NULL, `dividend_yield` double DEFAULT NULL, `rental_yield` double DEFAULT NULL, `realised_capital_gain` double DEFAULT NULL,`unrealised_capital_gain` double DEFAULT NULL, PRIMARY KEY (`proj_time`))", tableName)
	err := db.Exec(tableQuery).Error
	if err != nil {
		return err
	}
	return nil
}

func processYieldCurve(ctx context.Context,
	v *multipart.FileHeader,
	prodCode string,
	year int,
	month int,
	yieldCurveCode string) error {
	//var tableName = strings.ToLower(v.Filename[:len(v.Filename)-4])
	var delimiter rune
	delimiterFile, err := v.Open()
	if err != nil {
		return err
	}
	defer delimiterFile.Close()
	delimiter, err = utils.GetDelimiter(delimiterFile)

	var tableName string
	if year != 0 {
		tableName = "yield_curve"
	} else {
		tableName = "pricing_yield_curve"
	}

	file, err := v.Open()
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		reader := csv.NewReader(file)
		reader.TrimLeadingSpace = true
		reader.Comma = delimiter
		dec, _ := csvutil.NewDecoder(reader)
		dec.Header()

		//var headerRow string
		if tableName == "yield_curve" {
			if err := DB.Exec(fmt.Sprintf("delete from %s where year=%d and month=%d and yield_curve_code = %s ", tableName, year, month, yieldCurveCode)).Error; err != nil {
				return err
			}
			for {
				var pp models.YieldCurve
				if err := dec.Decode(&pp); err == io.EOF {
					break
				} else if err != nil {
					fmt.Println(err)
					return err
				}
				pp.Year = year
				pp.Month = month
				pp.YieldCurveCode = yieldCurveCode
				err = DB.Create(&pp).Error
				if err != nil {
					fmt.Println(err)
					return err
				}
			}

		}

		if tableName == "pricing_yield_curve" {
			if err := DB.Exec(fmt.Sprintf("delete from %s ", tableName)).Error; err != nil {
				return err
			}
			for {
				var pp models.PricingYieldCurve
				if err := dec.Decode(&pp); err == io.EOF {
					break
				} else if err != nil {
					fmt.Println(err)
					return err
				}
				//pp.Month = month
				//pp.YieldCurveCode = yieldCurveCode
				err = DB.Create(&pp).Error
				if err != nil {
					fmt.Println(err)
					return err
				}
			}

		}

	}
	//Register this activity...
	description := "Uploaded Yield Curve data for " + prodCode
	CreateActivity(ctx, "table_upload", description, 0, "yield_curve")
	return nil
}

func processPricingYieldCurve(ctx context.Context,
	v *multipart.FileHeader,
	prodCode string,
	month int,
	yieldCurveCode string) error {
	//var tableName = strings.ToLower(v.Filename[:len(v.Filename)-4])
	var delimiter rune
	delimiterFile, err := v.Open()
	if err != nil {
		return err
	}
	defer delimiterFile.Close()
	delimiter, err = utils.GetDelimiter(delimiterFile)

	//tableName := "pricing_yield_curve"

	file, err := v.Open()
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		reader := csv.NewReader(file)
		reader.TrimLeadingSpace = true
		reader.Comma = delimiter
		dec, _ := csvutil.NewDecoder(reader)
		dec.Header()

		//DB.Exec(fmt.Sprintf("delete from %s ", tableName))
		for {
			var pp models.PricingYieldCurve
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
				return err
			}
			//pp.Month = month
			//pp.YieldCurveCode = yieldCurveCode
			err = DB.Create(&pp).Error
			if err != nil {
				fmt.Println(err)
				return err
			}
		}

	}
	//Register this activity...
	description := "Uploaded Yield Curve data for " + prodCode
	CreateActivity(ctx, "table_upload", description, 0, "pricing_yield_curve")
	return nil
}

func AddToParameterTable(tableName string, values []string, fields string, prodCode string, year int) error {
	//tableName = "product_parameters"
	var err error
	vals := ""
	for i, v := range values {
		if i == 0 && strings.Contains(tableName, "parameters") {
			vals += "'" + prodCode + "'" + ","
		} else {
			vals += v + ","
		}
	}
	if year != 0 {
		vals += strconv.Itoa(year)
	}

	vals = strings.Trim(vals, ",")
	fields = strings.Trim(fields, ",")
	fmt.Println(vals)
	if year != 0 {
		err := DB.Exec(fmt.Sprintf("delete from %s where product_code='%s' and year=%d", tableName, prodCode, year)).Error
		if err != nil {
			return err
		}

	} else {
		err := DB.Exec(fmt.Sprintf("delete from %s where product_code='%s'", tableName, prodCode)).Error
		if err != nil {
			return err
		}
	}

	err = DB.Exec(fmt.Sprintf("insert into %s(%s) values(%s)", tableName, fields, vals)).Error
	if err != nil {
		return err
	}
	return nil
}

func BulkAddToParameterTable(tableName string, values []string, fields string) error {
	vals := ""
	for i, v := range values {
		if i == 0 && strings.Contains(tableName, "parameters") {
			vals += "'" + v + "'" + ","
		} else {
			vals += v + ","
		}
	}

	vals = strings.Trim(vals, ",")
	fields = strings.Trim(fields, ",")
	fmt.Println(vals)

	err := DB.Exec(fmt.Sprintf("insert into %s(%s) values(%s)", tableName, fields, vals)).Error
	if err != nil {
		return err
	}
	return nil
}

func AddToMarginTable(tableName string, values []string, fields string, prodCode string, year int) error {
	//tableName = "product_margins"
	vals := ""
	for i, v := range values {
		if i == 0 && strings.Contains(tableName, "margins") {
			vals += "'" + prodCode + "'" + ","
		} else {
			vals += v + ","
		}
	}
	if year != 0 {
		vals += strconv.Itoa(year)
	}

	vals = strings.Trim(vals, ",")
	fields = strings.Trim(fields, ",")
	fmt.Println("Values: ", vals)
	fmt.Println("Fields: ", fields)

	//Delete margins for the specified product code if any.
	if year != 0 {
		DB.Exec(fmt.Sprintf("delete from %s where product_code='%s' and year=%d", tableName, prodCode, year))

	} else {
		DB.Exec(fmt.Sprintf("delete from %s where product_code='%s'", tableName, prodCode))
	}
	err := DB.Exec(fmt.Sprintf("insert into %s(%s) values(%s)", tableName, fields, vals)).Error
	if err != nil {
		return err
	}

	return nil
}

func AddToPricingMarginTable(tableName string, values []string, fields string, prodCode string) error {
	//tableName = "product_margins"
	vals := ""
	for i, v := range values {
		if i == 0 && strings.Contains(tableName, "margins") {
			vals += "'" + prodCode + "'" + ","
		} else {
			vals += v + ","
		}
	}

	vals = strings.Trim(vals, ",")
	fields = strings.Trim(fields, ",")
	fmt.Println("Values: ", vals)
	fmt.Println("Fields: ", fields)

	//Delete margins for the specified product code if any.
	DB.Exec(fmt.Sprintf("delete from %s where product_code='%s'", tableName, prodCode))

	err := DB.Exec(fmt.Sprintf("insert into %s(%s) values(%s)", tableName, fields, vals)).Error
	if err != nil {
		return err
	}

	return nil
}

func AddToStaticTable(tableName string, values []string, fields string, prodCode string, year int, missingProdCode bool, missingYear bool) error {

	vals := ""
	for _, v := range values {
		_, err := strconv.ParseFloat(v, 64)
		if err != nil {
			// The value must be a string
			// This does not cover booleans though...
			vals += "'" + v + "'" + ","
		} else {
			vals += v + ","
		}
	}
	if missingProdCode {
		vals += "'" + prodCode + "'" + ","
	}

	if missingYear {
		vals += strconv.Itoa(year)
	}

	vals = strings.Trim(vals, ",")
	strings.TrimSpace(fields)
	fmt.Println(vals)
	fmt.Println(fields)

	err := DB.Exec(fmt.Sprintf("insert into %s (%s) values(%s)", tableName, fields, vals)).Error
	if err != nil {
		return err
	}

	return nil
}

func AddToPricingStaticTable(tableName string, values []string, fields string, prodCode string, missingProdCode bool) error {

	vals := ""
	for _, v := range values {
		_, err := strconv.ParseFloat(v, 64)
		if err != nil {
			// The value must be a string
			// This does not cover booleans though...
			vals += "'" + v + "'" + ","
		} else {
			vals += v + ","
		}
		//if i == 0 && strings.AggregatedProjectionsContain(tableName, "margin") {
		//	vals += "'" + v + "'" + ","
		//} else {
		//	vals += v + ","
		//}
	}
	if missingProdCode {
		vals += "'" + prodCode + "'" + ","
	}

	vals = strings.Trim(vals, ",")
	strings.TrimSpace(fields)
	fmt.Println(vals)
	fmt.Println(fields)

	err := DB.Exec(fmt.Sprintf("insert into %s (%s) values(%s)", tableName, fields, vals)).Error
	if err != nil {
		return err
	}

	return nil
}

func AddToYieldCurve(tableName string, year int, values []string, fields string) error {

	//if year == 0 {
	//	//purge pricing yield curve
	//	DB.Exec(fmt.Sprintf("delete from %s", tableName))
	//}
	vals := ""
	for i, v := range values {
		if i == 0 && strings.Contains(tableName, "margin") {
			vals += "'" + v + "'" + ","
		} else {
			vals += v + ","
		}
	}

	if !strings.Contains(fields, "year") {
		if year != 0 {
			vals += strconv.Itoa(year)
		}
	}
	vals = strings.Trim(vals, ",")
	fmt.Println(vals)
	//if year == 0 {
	//	err := DB.Exec(fmt.Sprintf("delete from %s where %s = %d", tableName, fields[0], vals[0])).Error
	//	if err != nil {
	//		fmt.Println(err)
	//		//Should return err here later....
	//	}
	//}
	err := DB.Exec(fmt.Sprintf("insert into %s (%s) values(%s)", tableName, fields, vals)).Error
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func processRatingTables(db *gorm.DB, v *multipart.FileHeader, delimiter rune, prodCode string, year int) error {
	var columnKey = ""
	var columnValue = ""
	var pricingColumnKey = ""
	var pricingColumnValue = ""

	var tableName = strings.ToLower(prodCode + "_" + v.Filename[:len(v.Filename)-4])
	var pricingTableName = strings.ToLower(prodCode + "_pricing_" + v.Filename[:len(v.Filename)-4])

	// Use the global DB for existence checks and data operations — the passed db
	// may be a transaction invalidated by prior DDL commits (edit mode).
	tableExists := DB.Migrator().HasTable(tableName)

	file, err := v.Open()
	if err != nil {
		return err
	} else {
		reader := csv.NewReader(file)
		reader.TrimLeadingSpace = true
		reader.Comma = delimiter

		data, _ := reader.ReadAll()
		if data == nil {
			return errors.New("an error was encountered while processing " + v.Filename)
		}

		startTime := time.Now()

		for i, row := range data {
			if i == 0 {
				for j := 0; j < len(row); j++ {
					if j < len(row)-1 {
						columnKey += row[j] + "_"
						pricingColumnKey += row[j] + "_"
					} else {
						columnValue = row[j]
						pricingColumnValue = row[j]
					}
				}
				columnKey = strings.ToLower(strings.Trim(columnKey, "_"))
				columnValue = strings.ToLower(strings.Replace(columnValue, " ", "_", -1))
				pricingColumnKey = strings.ToLower(strings.Trim(pricingColumnKey, "_"))
				pricingColumnValue = strings.ToLower(strings.Replace(pricingColumnValue, " ", "_", -1))

				columnKey = "year_" + columnKey

				if !tableExists {
					if err := CreateRatingTable(columnKey, columnValue, tableName); err != nil {
						return err
					}
					if err := CreateRatingTable(pricingColumnKey, pricingColumnValue, pricingTableName); err != nil {
						return err
					}
				}
			} else {
				var key = ""
				if year != 0 {
					key = strconv.Itoa(year) + "_"
				}

				var value float64
				for j := 0; j < len(row); j++ {
					if j < len(row)-1 {
						key += row[j] + "_"
					} else {
						value, _ = strconv.ParseFloat(row[j], 64)
					}
				}
				key = strings.Trim(key, "_")
				if year == 0 {
					tableName = pricingTableName
					columnKey = pricingColumnKey
				}
				err = AddToRatingTable(DB, tableName, key, value, columnKey)
				if err != nil {
					return err
				}
			}
		}

		endTime := time.Since(startTime)
		fmt.Println("Time taken to process " + v.Filename + " is " + endTime.String())
	}
	return nil
}

func processPricingRatingTables(db *gorm.DB, v *multipart.FileHeader, delimiter rune, prodCode string) error {
	var columnKey = ""
	var columnValue = ""
	var pricingColumnKey = ""
	var pricingColumnValue = ""

	var tableName = strings.ToLower(prodCode + "_" + v.Filename[:len(v.Filename)-4])
	var pricingTableName = strings.ToLower(prodCode + "_pricing_" + v.Filename[:len(v.Filename)-4])

	//check if the table exists, don't drop it. The call is an update
	tableExists := db.Migrator().HasTable(tableName)

	file, err := v.Open()
	if err != nil {
		fmt.Println(err.Error())
		return err
	} else {
		reader := csv.NewReader(file)
		reader.TrimLeadingSpace = true
		reader.Comma = delimiter

		data, _ := reader.ReadAll()
		if data == nil {
			return errors.New("an error was encountered while processing " + v.Filename)
		}

		startTime := time.Now()

		for i, row := range data {
			//var mp models.ModelPoint
			if i == 0 {
				fmt.Println("Header:", row)
				//Get the all but the last columns and concatenate them.
				for j := 0; j < len(row); j++ {
					if j < len(row)-1 {
						columnKey += row[j] + "_"
						pricingColumnKey += row[j] + "_"

					} else {
						columnValue = row[j]
						pricingColumnValue = row[j]
					}
				}
				columnKey = strings.ToLower(strings.Trim(columnKey, "_"))
				columnValue = strings.ToLower(strings.Replace(columnValue, " ", "_", -1))
				pricingColumnKey = strings.ToLower(strings.Trim(pricingColumnKey, "_"))
				pricingColumnValue = strings.ToLower(strings.Replace(pricingColumnValue, " ", "_", -1))

				columnKey = "year_" + columnKey

				if !tableExists {
					err := CreateRatingTable(columnKey, columnValue, tableName)
					if err != nil {
						fmt.Println(err)
					}

					err = CreateRatingTable(pricingColumnKey, pricingColumnValue, pricingTableName)
					if err != nil {
						fmt.Println(err)
					}

				}
			} else {
				var key = ""

				var value float64
				for j := 0; j < len(row); j++ {
					if j < len(row)-1 {
						key += row[j] + "_"
					} else {
						value, _ = strconv.ParseFloat(row[j], 64)
					}
				}
				key = strings.Trim(key, "_")
				err = AddToPricingRatingTable(db, pricingTableName, key, value, pricingColumnKey)
				if err != nil {
					return err
				}
			}
		}

		endTime := time.Since(startTime)

		fmt.Println("Time taken to process " + v.Filename + " is " + endTime.String())
	}
	return nil
}

func AddToRatingTable(db *gorm.DB, tableName string, key string, value float64, columnKey string) error {
	//We would assume that we want to delete the existing value...
	deleteQuery := fmt.Sprintf("delete from %s where %s = '%s'", tableName, columnKey, key)
	fmt.Println(deleteQuery)

	err := db.Exec(deleteQuery).Error
	if err != nil {
		fmt.Println(err)
	}

	query := fmt.Sprintf("insert into %s values('%s',%f)", tableName, key, value)
	fmt.Println(query)
	err = db.Exec(query).Error
	if err != nil {
		return err
	}
	return nil
}

func AddToPricingRatingTable(db *gorm.DB, tableName string, key string, value float64, columnKey string) error {
	//We would assume that we want to delete the existing value...
	deleteQuery := fmt.Sprintf("delete from %s where %s = '%s'", tableName, columnKey, key)
	fmt.Println(deleteQuery)

	err := db.Exec(deleteQuery).Error
	if err != nil {
		fmt.Println(err)
	}

	query := fmt.Sprintf("insert into %s values('%s',%f)", tableName, key, value)
	fmt.Println(query)
	err = db.Exec(query).Error
	if err != nil {
		return err
	}
	return nil
}

func CreateProduct(form *multipart.Form) (models.Product, error) {
	var pc models.Product
	var selectedFeatures []models.BaseFeature

	pcString := form.Value["product"][0]
	editMode, err := strconv.ParseBool(form.Value["editMode"][0])
	if err != nil {
		return pc, errors.New("no valid edit mode was provided")
	}
	pcString = strings.Replace(pcString, "\"id\":\"\"", "\"id\":0", -1)
	err = json.Unmarshal([]byte(pcString), &pc)
	if err != nil {
		return pc, errors.New("unable to parse the product configuration object: " + err.Error())
	}

	err = DB.Transaction(func(tx *gorm.DB) error {
		//Before we save, zero out all ids and ensure we add the global tables.
		if editMode {
			DeleteProduct(tx, pc.ID)
		}
		for i := range pc.ProductModelpointVariables {
			pc.ProductModelpointVariables[i].ID = 0
		}
		for i := range pc.ProductTransitionStates {
			pc.ProductTransitionStates[i].ID = 0
		}

		pc.Status = PENDING
		pc.CreatedAt = time.Now()
		err = tx.Create(&pc).Error
		if err != nil {
			return err
		}
		err = tx.Save(&pc).Error
		if err != nil {
			return err
		}
		//Generate Model Point variables....
		// DDL operations (AutoMigrate) must run on the global DB connection, not on
		// the transaction — MySQL DDL in DeleteProduct causes an implicit commit that
		// invalidates tx, so passing tx here would fail in edit mode.
		err = GenerateModelPointTable(DB, &pc)
		if err != nil {
			return errors.New("there was a problem generating the model point data:  " + err.Error())
		}

		err = json.Unmarshal([]byte(form.Value["product_features"][0]), &selectedFeatures)

		if err != nil {
			return errors.New("parsing selected features threw an error: " + err.Error())
		}
		//End of hack
		multipartFiles := form.File["files[]"]

		err = ProcessTables(tx, &pc, multipartFiles)
		if err != nil {
			return errors.New("one or more associated assumption tables returned an error: " + err.Error())
		}

		// Special cases of profitabilities
		tx.Where("product_id = ? and class = ?", pc.ID, "Profitabilities").Delete(&models.ProductTable{})

		err = ProcessFeatures(tx, pc.ProductCode, selectedFeatures)
		if err != nil {
			return errors.New("processing selected features threw an error: " + err.Error())
		}
		return nil
	})

	return pc, err
}

func GenerateModelPointTable(db *gorm.DB, pc *models.Product) error {
	tableName := strings.ToLower(pc.ProductCode) + "_modelpoints"
	pricingTableName := strings.ToLower(pc.ProductCode) + "_pricing_modelpoints"

	if err := db.Table(tableName).AutoMigrate(&models.ProductModelPoint{}); err != nil {
		return err
	}
	if err := db.Table(pricingTableName).AutoMigrate(&models.ProductPricingModelPoint{}); err != nil {
		return err
	}
	return nil
}

func DeleteProduct(db *gorm.DB, id int) error {
	var err error
	var ratingFactors []models.ProductRatingFactor
	prod, _ := GetProductById(id)

	db.Where("product_id=?", id).Find(&ratingFactors)

	for _, factor := range ratingFactors {
		tableName := strings.ToLower(prod.ProductCode) + "_" + strings.ToLower(factor.TransitionTable)
		pricingTableName := strings.ToLower(prod.ProductCode) + "_pricing_" + strings.ToLower(factor.TransitionTable)
		if db.Migrator().HasTable(tableName) {
			err := db.Migrator().DropTable(tableName)
			if err != nil {
				log.Error().Err(err)
			}
		}

		if db.Migrator().HasTable(pricingTableName) {
			err := db.Migrator().DropTable(pricingTableName)
			if err != nil {
				log.Error().Err(err)
			}
		}

		fmt.Println(err)
	}
	db.Where("product_id=?", id).Delete(&ratingFactors)

	if db.Migrator().HasTable(strings.ToLower(prod.ProductCode) + "_modelpoints") {
		err := db.Migrator().DropTable(strings.ToLower(prod.ProductCode) + "_modelpoints")
		if err != nil {
			log.Error().Err(err)
		}
	}

	if db.Migrator().HasTable(strings.ToLower(prod.ProductCode) + "_pricing_modelpoints") {
		err := db.Migrator().DropTable(strings.ToLower(prod.ProductCode) + "_pricing_modelpoints")
		if err != nil {
			log.Error().Err(err)
		}
	}

	if db.Migrator().HasTable(strings.ToLower(prod.ProductCode) + "_pricing_yield_curve") {
		err := db.Migrator().DropTable(strings.ToLower(prod.ProductCode) + "_pricing_yield_curve")
		if err != nil {
			log.Error().Err(err)
		}
	}

	err = db.Where("product_code=?", prod.ProductCode).Delete(&models.ProductFeatures{}).Error
	if err != nil {
		log.Error().Err(err)
	}
	db.Where("product_code=?", prod.ProductCode).Delete(&models.ProductMargins{})
	db.Where("product_code=?", prod.ProductCode).Delete(&models.ProductPricingMargins{})
	db.Where("product_code=?", prod.ProductCode).Delete(&models.ProductParameters{})
	db.Where("product_code=?", prod.ProductCode).Delete(&models.ProductPricingParameters{})
	db.Where("product_code=?", prod.ProductCode).Delete(&models.PricingParameter{})
	db.Where("product_id=?", prod.ID).Delete(&models.ProductModelpointVariable{})
	db.Where("product_id=?", prod.ID).Delete(&models.ProductTable{})
	db.Where("product_id=?", prod.ID).Delete(&models.ProductTransition{})
	db.Where("product_id=?", prod.ID).Delete(&models.ProductReinsurance{})
	db.Where("product_id=?", prod.ID).Delete(&models.ProductRatingFactor{})
	db.Where("product_id=?", prod.ID).Delete(&models.ProductTransitionState{})
	db.Where("product_id=?", prod.ID).Delete(&models.GlobalTable{})
	db.Where("product_code=?", prod.ProductCode).Delete(&models.ProductRider{})
	db.Where("product_code=?", prod.ProductCode).Delete(&models.ProductPricingRider{})
	db.Where("product_code=?", prod.ProductCode).Delete(&models.ProductChildAdditionalSumAssured{})
	db.Where("product_code=?", prod.ProductCode).Delete(&models.ProductPricingChildAdditionalSumAssured{})
	db.Where("product_code=?", prod.ProductCode).Delete(&models.ProductChildSumAssured{})
	db.Where("product_code=?", prod.ProductCode).Delete(&models.ProductAdditionalSumAssured{})
	db.Where("product_code=?", prod.ProductCode).Delete(&models.ProductPricingAdditionalSumAssured{})
	db.Where("product_code=?", prod.ProductCode).Delete(&models.Projection{})
	db.Where("product_code=?", prod.ProductCode).Delete(&models.ProductClawback{})
	db.Where("product_code=?", prod.ProductCode).Delete(&models.ProductPricingClawback{})
	db.Where("product_id=?", prod.ID).Delete(&models.ProductPricingTable{})
	var rfs []models.ProductRatingFactor
	db.Where("product_id=?", prod.ID).Find(&rfs)
	for _, rf := range rfs {
		db.Where("product_rating_factor_id = ?", rf.ID).Delete(&models.Fds{})
	}

	db.Where("product_id=?", prod.ID).Delete(&models.ProductRatingFactor{})
	if err != nil {
		log.Error().Err(err).Send()
	}

	err = db.Where("id = ?", id).Delete(&models.Product{}).Error
	if err != nil {
		log.Error().Err(err).Send()
	}

	//Delete Pricing tables
	if db.Migrator().HasTable(strings.ToLower(prod.ProductCode) + "pricing_modelpoints") {
		err := db.Migrator().DropTable(strings.ToLower(prod.ProductCode) + "pricing_modelpoints")
		if err != nil {
			log.Error().Err(err)
		}
	}

	//Delete All other associated data with the product ID
	var jobProducts []models.JobProduct
	db.Where("product_id=?", prod.ID).Find(&jobProducts)

	for _, jobProduct := range jobProducts {
		db.Where("job_product_id=?", jobProduct.ID).Delete(&models.AggregatedProjection{})
		db.Where("job_product_id=?", jobProduct.ID).Delete(&models.ScopedAggregatedProjection{})
	}
	db.Where("product_id=?", prod.ID).Delete(&models.JobProduct{})
	return nil
}

func DeleteProductTableData(productId, tableId int) error {
	//Get the product
	product, err := GetProductById(productId)
	if err != nil {
		return err
	}
	productTable, err := GetProductTableById(tableId)
	if err != nil {
		return err
	}

	//Check the table class type
	switch productTable.Class {
	case "TransitionRates":
		//Append product code to the table name
		tableName := strings.ToLower(product.ProductCode) + "_" + strings.ToLower(productTable.Name)
		err := DB.Exec("delete from " + tableName).Error
		if err != nil {
			log.Error().Err(err).Send()
		}
	case "Parameters":
		err := DB.Delete(&models.ProductParameters{}, "product_code = ?", product.ProductCode).Error
		if err != nil {
			log.Error().Err(err).Send()
		}
	case "Margins":
		err := DB.Delete(&models.ProductMargins{}, "product_code = ?", product.ProductCode).Error
		if err != nil {
			log.Error().Err(err).Send()
		}
	case "LapseMargins":
		err := DB.Delete(&models.ProductLapseMargin{}, "product_code = ?", product.ProductCode).Error
		if err != nil {
			log.Error().Err(err).Send()
		}
	case "Distribution":
		tableName := "product_" + strings.ToLower(productTable.Name)
		err := DB.Exec("delete from " + tableName + " where product_code = '" + product.ProductCode + "'").Error
		if err != nil {
			log.Error().Err(err).Send()
		}
	case "BenefitStructure":
		tableName := "product_" + strings.ToLower(productTable.Name)
		err := DB.Exec("delete from " + tableName + " where product_code = '" + product.ProductCode + "'").Error
		if err != nil {
			log.Error().Err(err).Send()
		}
	case "Valuations":
		tableName := "product_" + strings.ToLower(productTable.Name)
		err := DB.Exec("delete from " + tableName + " where product_code = '" + product.ProductCode + "'").Error
		if err != nil {
			log.Error().Err(err).Send()
		}
	}
	return err
}

func DeleteProductTableDatav2(productId, tableId, year int) error {
	//Get the product
	product, err := GetProductById(productId)
	if err != nil {
		return err
	}
	productTable, err := GetProductTableById(tableId)
	if err != nil {
		return err
	}

	//Check the table class type
	switch productTable.Class {
	case "TransitionRates":
		//Append product code to the table name
		tableName := strings.ToLower(product.ProductCode) + "_" + strings.ToLower(productTable.Name)

		var columns []string
		DB.Table(tableName).Raw("SELECT column_name FROM information_schema.columns WHERE table_name = ? and column_name like 'year%'", tableName).Pluck("column_name", &columns)

		// get the colums for the table
		//var columns []string
		//DB.Raw("SHOW COLUMNS FROM " + tableName).Scan(&columns)
		deleteQuery := fmt.Sprintf("delete from %s where substring(%s, 1, 4) = %d", tableName, columns[0], year)
		err = DB.Exec(deleteQuery).Error
		if err != nil {
			log.Error().Err(err).Send()
		}
	case "Parameters":
		err := DB.Delete(&models.ProductParameters{}, "product_code = ? and year = ?", product.ProductCode, year).Error
		if err != nil {
			log.Error().Err(err).Send()
		}
	case "Margins":
		err := DB.Delete(&models.ProductMargins{}, "product_code = ? and year = ?", product.ProductCode, year).Error
		if err != nil {
			log.Error().Err(err).Send()
		}
	case "LapseMargins":
		err := DB.Delete(&models.ProductLapseMargin{}, "product_code = ? and year = ?", product.ProductCode, year).Error
		if err != nil {
			log.Error().Err(err).Send()
		}
	case "Distribution":
		tableName := "product_" + strings.ToLower(productTable.Name)
		var deleteQuery string
		if year != 0 {
			deleteQuery = fmt.Sprintf("delete from %s where product_code = '%s' and year = %d", tableName, product.ProductCode, year)
		} else {
			deleteQuery = fmt.Sprintf("delete from %s where product_code = '%s'", tableName, product.ProductCode)
		}
		err := DB.Exec(deleteQuery).Error
		if err != nil {
			log.Error().Err(err).Send()
		}
	case "BenefitStructure":
		tableName := "product_" + strings.ToLower(productTable.Name)
		if year != 0 {
			deleteQuery := fmt.Sprintf("delete from %s where product_code = '%s' and year = %d", tableName, product.ProductCode, year)
			err := DB.Exec(deleteQuery).Error
			if err != nil {
				log.Error().Err(err).Send()
			}
		} else {
			deleteQuery := fmt.Sprintf("delete from %s where product_code = '%s'", tableName, product.ProductCode)
			err := DB.Exec(deleteQuery).Error
			if err != nil {
				log.Error().Err(err).Send()
			}
		}
		deleteQuery := fmt.Sprintf("delete from %s where product_code = '%s' and year = %d", tableName, product.ProductCode, year)
		err := DB.Exec(deleteQuery).Error
		if err != nil {
			log.Error().Err(err).Send()
		}
	case "Valuations":
		tableName := "product_" + strings.ToLower(productTable.Name)
		if year != 0 {
			deleteQuery := fmt.Sprintf("delete from %s where product_code = '%s' and year = %d", tableName, product.ProductCode, year)
			err := DB.Exec(deleteQuery).Error
			if err != nil {
				log.Error().Err(err).Send()
			}
		} else {
			deleteQuery := fmt.Sprintf("delete from %s where product_code = '%s'", tableName, product.ProductCode)
			err := DB.Exec(deleteQuery).Error
			if err != nil {
				log.Error().Err(err).Send()
			}
		}
	}
	return err
}

func DeleteProductPricingTableData(productId, tableId int) error {
	//Get the product
	product, err := GetProductById(productId)
	if err != nil {
		return err
	}
	productTable, err := GetProductPricingTableById(tableId)
	if err != nil {
		return err
	}

	//Check the table class type
	switch productTable.Class {
	case "TransitionRates":
		//Append product code to the table name
		tableName := strings.ToLower(product.ProductCode) + "_pricing_" + strings.ToLower(productTable.Name)
		err := DB.Exec("delete from " + tableName).Error
		if err != nil {
			log.Error().Err(err).Send()
		}
	case "Parameters":
		err := DB.Where("product_code = ?", product.ProductCode).Delete(&models.ProductPricingParameters{}).Error
		if err != nil {
			log.Error().Err(err).Send()
		}
	case "Margins":
		err := DB.Delete(&models.ProductPricingMargins{}, "product_code = ?", product.ProductCode).Error
		if err != nil {
			log.Error().Err(err).Send()
		}
	case "LapseMargins":
		err := DB.Delete(&models.ProductPricingLapseMargin{}, "product_code = ?", product.ProductCode).Error
		if err != nil {
			log.Error().Err(err).Send()
		}
	case "Distribution":
		tableName := "product_pricing_" + strings.ToLower(productTable.Name)
		err := DB.Exec("delete from " + tableName + " where product_code = '" + product.ProductCode + "'").Error
		if err != nil {
			log.Error().Err(err).Send()
		}
	case "BenefitStructure":
		tableName := "product_pricing_" + strings.ToLower(productTable.Name)
		err := DB.Exec("delete from " + tableName + " where product_code = '" + product.ProductCode + "'").Error
		if err != nil {
			log.Error().Err(err).Send()
		}
	case "Valuations":
		tableName := "product_pricing_" + strings.ToLower(productTable.Name)
		err := DB.Exec("delete from " + tableName + " where product_code = '" + product.ProductCode + "'").Error
		if err != nil {
			log.Error().Err(err).Send()
		}
	case "Global":
		tableName := "pricing_" + strings.ToLower(productTable.Name)
		fmt.Println(tableName)
		err := DB.Exec("delete from " + tableName).Error
		if err != nil {
			log.Error().Err(err).Send()
		}
	}
	return err
}

func GetProductTableById(id int) (models.ProductTable, error) {
	var productTable models.ProductTable
	err := DB.Where("id = ?", id).First(&productTable).Error
	return productTable, err
}

func GetProductPricingTableById(id int) (models.ProductPricingTable, error) {
	var productTable models.ProductPricingTable
	err := DB.Where("id = ?", id).First(&productTable).Error
	return productTable, err

}

func ActivateProduct(payload map[string]interface{}) error {
	prodId := int(payload["productId"].(float64))
	remarks := payload["description"].(string)
	prod, _ := GetProductById(prodId)
	prod.Status = APPROVED
	err := DB.Save(&prod).Error
	var task models.Task
	DB.Where("product_id = ?", prodId).First(&task)
	task.Status = "completed"
	task.Description = remarks
	DB.Save(&task)
	return err
}

func GetProductModelPointVariables(prodId int) []models.ProductModelpointVariable {
	var pMps []models.ProductModelpointVariable
	if err := DB.Where("product_id=?", prodId).Find(&pMps).Error; err != nil {
		fmt.Println(err)
	}

	return pMps
}

func GetModelPointCountForProduct(prodId int) ([]models.ModelPointCount, error) {
	// First, check if we already have cached counts in ProductModelPointCount
	var cached []models.ProductModelPointCount
	if err := DB.Where("product_id = ?", prodId).Find(&cached).Error; err == nil && len(cached) > 0 {
		// Convert to []models.ModelPointCount for existing response shape
		resp := make([]models.ModelPointCount, 0, len(cached))
		for _, c := range cached {
			resp = append(resp, models.ModelPointCount{Year: c.Year, Version: c.Version, Count: c.Count})
		}
		return resp, nil
	}

	// Otherwise compute from the modelpoints table
	prod, err := GetProductById(prodId)
	if err != nil {
		return nil, err
	}
	prod_code := prod.ProductCode
	table := strings.ToLower(prod_code) + "_modelpoints"

	counts := make([]models.ModelPointCount, 0)
	err = DB.Raw("select year, mp_version as version, count(*) as count from " + table + " group by year, mp_version").Scan(&counts).Error
	if err != nil {
		return counts, err
	}

	// Persist into ProductModelPointCount for future quick retrieval
	// Remove existing cache for this product first to avoid duplicates
	_ = DB.Where("product_id = ?", prodId).Delete(&models.ProductModelPointCount{}).Error
	for _, c := range counts {
		pmc := models.ProductModelPointCount{ProductId: prodId, ProductCode: prod.ProductCode, Year: c.Year, Version: c.Version, Count: c.Count}
		// Best-effort insert; ignore errors to not block the response
		_ = DB.Create(&pmc).Error
	}

	return counts, nil
}

func GetModelPointsForProductAndYear(prodId int, year int, version string, limit int) ([]models.ProductModelPoint, error) {
	product, err := GetProductById(prodId)
	tableName := strings.ToLower(product.ProductCode) + "_modelpoints"
	var mps []models.ProductModelPoint
	if limit > 0 {
		if version != "" {
			err = DB.Table(tableName).Where("year=? and mp_version=?", year, version).Limit(limit).Find(&mps).Error
		} else {
			err = DB.Table(tableName).Where("year=? ", year).Limit(limit).Find(&mps).Error
		}

	} else {
		if version != "" {
			err = DB.Table(tableName).Where("year=? and mp_version=?", year, version).Find(&mps).Error
		} else {
			err = DB.Table(tableName).Where("year=? ", year).Find(&mps).Error
		}
	}

	if err != nil {
		return nil, err
	}
	return mps, nil
}

func GetModelPointsForProductAndYearExcel(prodId int, year int, version string) ([]byte, error) {
	product, err := GetProductById(prodId)
	if err == nil {
		tableName := strings.ToLower(product.ProductCode) + "_modelpoints"
		var dQuery string
		if version != "" {
			dQuery = fmt.Sprintf("select * from %s where year=%d and mp_version='%s'", tableName, year, version)
		} else {
			dQuery = fmt.Sprintf("select * from %s where year=%d", tableName, year)
		}

		excelData, err := exportTableToExcel(dQuery)
		if err != nil {
			return nil, nil
		}

		return excelData, nil
	} else {
		return nil, err
	}
}

func DeleteModelPointsForProductAndYear(prodId int, year int, version string) error {
	product, err := GetProductById(prodId)
	tableName := strings.ToLower(product.ProductCode) + "_modelpoints"
	if version == "" {
		err = DB.Table(tableName).Where("year=?", year).Delete(&models.ProductModelPoint{}).Error
	} else {
		err = DB.Table(tableName).Where("year=? and mp_version = ?", year, version).Delete(&models.ProductModelPoint{}).Error
	}

	if err != nil {
		return err
	}

	// Update ProductModelPointCount cache to reflect deletions
	if version == "" {
		// Deleting all versions for the year: remove all cache rows for this product/year
		_ = DB.Where("product_id = ? AND year = ?", product.ID, year).Delete(&models.ProductModelPointCount{}).Error
	} else {
		// For a specific version, recompute the remaining count after delete
		var remaining int64
		if cerr := DB.Table(tableName).Where("year = ? AND mp_version = ?", year, version).Count(&remaining).Error; cerr == nil {
			if remaining == 0 {
				// No rows remain; remove cache entry to avoid stale non-zero counts
				_ = DB.Where("product_id = ? AND year = ? AND version = ?", product.ID, year, version).Delete(&models.ProductModelPointCount{}).Error
			} else {
				// Upsert-like behavior: replace with new count
				_ = DB.Where("product_id = ? AND year = ? AND version = ?", product.ID, year, version).Delete(&models.ProductModelPointCount{}).Error
				_ = DB.Create(&models.ProductModelPointCount{ProductId: product.ID, Year: year, Version: version, Count: int(remaining)}).Error
			}
		}
	}

	return nil
}

func GetProjectionJobRun(runId int) (models.RunParameters, error) {
	var run models.RunParameters
	err := DB.Where("projection_job_id = ?", runId).First(&run).Error
	if err != nil {
		return run, err
	}
	return run, nil
}

func GetJobProduct(jobProductId int) models.JobProduct {
	var jobProduct models.JobProduct
	DB.Where("id=?", jobProductId).Find(&jobProduct)
	return jobProduct
}

func GetBasisForValuations(productCode, year string) ([]string, error) {
	var basis []string
	basisYear, _ := strconv.Atoi(year)
	err := DB.Table("product_parameters").Distinct("basis").Where("product_code=? and year = ?", productCode, basisYear).Pluck("basis", &basis).Error
	return basis, err
}

func GetModelPointVersions(productId int, year int) ([]string, error) {
	var versions []string
	prod, _ := GetProductById(productId)
	tableName := strings.ToLower(prod.ProductCode) + "_modelpoints"
	err := DB.Table(tableName).Distinct("mp_version").Where("year=?", year).Pluck("mp_version", &versions).Error
	if err != nil {
		return nil, err
	}
	return versions, nil
}

func GetProductTableYears(productId int, tableName, tableType string, isPricing bool) []int {
	var years = make([]int, 0)
	var fullTableName string

	prod, _ := GetProductById(productId)
	if tableType == "TransitionRates" {
		if isPricing {
			fullTableName = prod.ProductCode + "_pricing_" + tableName
		} else {
			fullTableName = prod.ProductCode + "_" + tableName
		}
	} else {
		fullTableName = "product_" + tableName
	}

	fullTableName = strings.ToLower(fullTableName)

	if tableType == "TransitionRates" {
		// we need to get the columns from the table
		var columns []string
		DB.Table(fullTableName).Raw("SELECT column_name FROM information_schema.columns WHERE table_name = ? and column_name like 'year%'", fullTableName).Pluck("column_name", &columns)
		DB.Table(fullTableName).Distinct(fmt.Sprintf("substring(%s, 1, 4) as year", columns[0])).Pluck("year", &years)
	} else {
		DB.Table(fullTableName).Where("product_code = ?", prod.ProductCode).Distinct("year").Pluck("year", &years)
	}

	return years
}

func CreateModelPointStatsForProductAndYear(productCode string, year int, version string) {
	table := strings.ToLower(productCode) + "_modelpoints"
	var err error

	variables := []string{"age_at_entry", "main_member_age_at_entry", "gender", "main_member_gender", "sum_assured", "outstanding_loan", "instalment", "annual_premium", "premium_rate", "interest", "term", "outstanding_term_months"}
	for _, v := range variables {
		var stats models.ProductModelPointVariableStats
		if v != "gender" && v != "main_member_gender" {
			err = DB.Raw(fmt.Sprintf("select MIN(%s) as min, MAX(%s) as max, SUM(%s) as sum, AVG(%s) as average from %s where year = %d", v, v, v, v, table, year)).Scan(&stats).Error
			if err != nil {
				fmt.Println(err)
			}
			err = DB.Raw(fmt.Sprintf("select count(%s) as number_of_zeroes from %s where %s = 0 and year = %d", v, table, v, year)).Scan(&stats).Error
			err = DB.Raw(fmt.Sprintf("select count(%s) as number_of_lives from  %s where year = %d", v, table, year)).Scan(&stats).Error
			stats.ProductCode = productCode
			stats.Variable = v
			stats.Year = year
			stats.Version = version
		} else {
			err = DB.Raw(fmt.Sprintf("SELECT count(%s) as male FROM %s where year = %d and gender='M'", v, table, year)).Scan(&stats).Error
			err = DB.Raw(fmt.Sprintf("SELECT count(%s) as female FROM %s where year = %d and gender='F'", v, table, year)).Scan(&stats).Error
			err = DB.Raw(fmt.Sprintf("SELECT count(%s) as number_of_lives FROM %s where year = %d", v, table, year)).Scan(&stats).Error
			stats.ProductCode = productCode
			stats.Variable = v
			stats.Year = year
			stats.Version = version
		}
		DB.Where("product_code = ? and variable = ? and year = ?", productCode, v, year).FirstOrCreate(&stats)
	}

}

func GetModelPointVariableStatsForProductAndYear(prodId int, year int, version string) []models.ProductModelPointVariableStats {
	prod, err := GetProductById(prodId)
	if err != nil {
		fmt.Println(fmt.Errorf("error getting product by ID %d: %w", prodId, err))
		return nil // Or handle error as appropriate
	}

	prodCode := prod.ProductCode
	table := strings.ToLower(prodCode) + "_modelpoints" // Dynamic table name

	var allStats []models.ProductModelPointVariableStats
	// Try to fetch all stats if they already exist for this product/year/version combination
	// This initial Find might be for a scenario where if any stats exist, you return them all.
	// However, the FirstOrCreate later suggests an upsert-per-variable logic.
	// Clarify if the intent is to return immediately if *any* stats exist, or to always calculate/update.
	// For now, respecting the original structure:
	err = DB.Where("product_code = ? AND year = ? AND version = ?", prodCode, year, version).Find(&allStats).Error
	if err != nil {
		fmt.Println(fmt.Errorf("error fetching existing stats for product %s, year %d, version %s: %w", prodCode, year, version, err))
		// Decide if you should return or proceed to calculate
	}

	if len(allStats) > 0 {
		return allStats // Return if stats were found by the initial Find
	}

	// If no stats were found by the initial Find, proceed to calculate them.
	// `allStats` will be repopulated with newly calculated/FirstOrCreate'd stats.
	allStats = []models.ProductModelPointVariableStats{} // Reset allStats if we are calculating

	variables := []string{
		"age_at_entry", "main_member_age_at_entry", "gender", "main_member_gender",
		"sum_assured", "outstanding_loan", "instalment", "annual_premium",
		"premium_rate", "interest", "term", "outstanding_term_months",
	}

	baseWhereClause := fmt.Sprintf("year = %d AND mp_version='%s'", year, version)

	for _, v := range variables {
		// Initialize a new stats object for each variable
		currentStats := models.ProductModelPointVariableStats{
			ProductCode: prodCode,
			Variable:    v,
			Year:        year,
			Version:     version,
		}
		var queryErr error

		if v != "gender" && v != "main_member_gender" {
			// For numeric variables
			var mmasResult struct { // Min, Max, Sum, Average
				Min     *float64 `gorm:"column:min"` // Use pointers for safety with NULLs
				Max     *float64 `gorm:"column:max"`
				Sum     *float64 `gorm:"column:sum"`
				Average *float64 `gorm:"column:average"`
			}
			query := fmt.Sprintf("SELECT MIN(%s) AS min, MAX(%s) AS max, SUM(%s) AS sum, AVG(%s) AS average FROM %s WHERE %s", v, v, v, v, table, baseWhereClause)
			queryErr = DB.Raw(query).Scan(&mmasResult).Error
			if queryErr != nil {
				fmt.Println(fmt.Errorf("error for min/max/sum/avg on %s for product %s: %w", v, prodCode, queryErr))
			}
			if mmasResult.Min != nil {
				currentStats.Min = *mmasResult.Min
			}
			if mmasResult.Max != nil {
				currentStats.Max = *mmasResult.Max
			}
			if mmasResult.Sum != nil {
				currentStats.Sum = *mmasResult.Sum
			}
			if mmasResult.Average != nil {
				currentStats.Average = *mmasResult.Average
			}

			var zeroesResult struct {
				NumberOfZeroes int64 `gorm:"column:number_of_zeroes"`
			}
			query = fmt.Sprintf("SELECT COUNT(%s) AS number_of_zeroes FROM %s WHERE %s = 0 AND %s", v, table, v, baseWhereClause)
			queryErr = DB.Raw(query).Scan(&zeroesResult).Error
			if queryErr != nil {
				fmt.Println(fmt.Errorf("error for number_of_zeroes on %s for product %s: %w", v, prodCode, queryErr))
			}
			currentStats.NumberOfZeroes = int(zeroesResult.NumberOfZeroes)

			var livesResult struct {
				NumberOfLives int64 `gorm:"column:number_of_lives"`
			}
			// Assuming number_of_lives is a count of non-null entries for the variable
			query = fmt.Sprintf("SELECT COUNT(%s) AS number_of_lives FROM %s WHERE %s IS NOT NULL AND %s", v, table, v, baseWhereClause)
			queryErr = DB.Raw(query).Scan(&livesResult).Error
			if queryErr != nil {
				fmt.Println(fmt.Errorf("error for number_of_lives on %s for product %s: %w", v, prodCode, queryErr))
			}
			currentStats.NumberOfLives = int(livesResult.NumberOfLives)

		} else {
			// For gender variables (v == "gender" or v == "main_member_gender")
			// Assuming the column name for gender is indeed the value of 'v'
			genderColumnName := v // e.g., "gender" or "main_member_gender"

			var maleCountResult struct {
				Male int64 `gorm:"column:male"`
			}
			// Note: The original query for male/female used `count(v)`. If `v` is "gender", `count(gender)` is okay.
			// Ensure the column `v` (e.g. "gender") actually exists for this condition to be meaningful.
			// The condition `gender='M'` should ideally use the genderColumnName.
			query := fmt.Sprintf("SELECT COUNT(*) AS male FROM %s WHERE %s = 'M' AND %s", table, genderColumnName, baseWhereClause)
			queryErr = DB.Raw(query).Scan(&maleCountResult).Error
			if queryErr != nil {
				fmt.Println(fmt.Errorf("error for male count on %s for product %s: %w", v, prodCode, queryErr))
			}
			currentStats.Male = float64(maleCountResult.Male)

			var femaleCountResult struct {
				Female int64 `gorm:"column:female"`
			}
			query = fmt.Sprintf("SELECT COUNT(*) AS female FROM %s WHERE %s = 'F' AND %s", table, genderColumnName, baseWhereClause)
			queryErr = DB.Raw(query).Scan(&femaleCountResult).Error
			if queryErr != nil {
				fmt.Println(fmt.Errorf("error for female count on %s for product %s: %w", v, prodCode, queryErr))
			}
			currentStats.Female = float64(femaleCountResult.Female)

			var genderLivesResult struct {
				NumberOfLives int64 `gorm:"column:number_of_lives"`
			}
			// Number of lives for gender could be total rows where gender is M or F, or simply total relevant rows.
			// The original query `SELECT count(%s) ...` for lives might be okay if `v` is a valid column here.
			// A safer bet is COUNT(*) if we're counting rows based on gender.
			// Or, if it's counting non-null gender entries:
			query = fmt.Sprintf("SELECT COUNT(%s) AS number_of_lives FROM %s WHERE %s IS NOT NULL AND %s", genderColumnName, table, genderColumnName, baseWhereClause)
			// If it means total M + F: currentStats.NumberOfLives = currentStats.Male + currentStats.Female
			// Or if it's truly COUNT(v) from original:
			// query = fmt.Sprintf("SELECT COUNT(%s) AS number_of_lives FROM %s WHERE %s", v, table, baseWhereClause)
			queryErr = DB.Raw(query).Scan(&genderLivesResult).Error
			if queryErr != nil {
				fmt.Println(fmt.Errorf("error for number_of_lives (gender) on %s for product %s: %w", v, prodCode, queryErr))
			}
			currentStats.NumberOfLives = int(genderLivesResult.NumberOfLives) // Or sum Male + Female if that's the definition
		}

		// Persist the fully assembled currentStats object for this variable
		// FirstOrCreate will find a record matching the primary key fields (or fields in Where)
		// or create a new one if not found, using the `currentStats` object.
		// Ensure ProductCode, Variable, Year, RiskRateCode are sufficient unique keys or primary keys for ProductModelPointVariableStats.
		err = DB.Where(models.ProductModelPointVariableStats{
			ProductCode: currentStats.ProductCode,
			Variable:    currentStats.Variable,
			Year:        currentStats.Year,
			Version:     currentStats.Version,
		}).Assign(currentStats).FirstOrCreate(&currentStats).Error // Use Assign to update if found, then create if not

		if err != nil {
			fmt.Println(fmt.Errorf("error in FirstOrCreate for variable %s, product %s: %w", v, prodCode, err))
		}

		allStats = append(allStats, currentStats) // Append the (created or found & assigned) stat
	}

	return allStats
}

//func GetModelPointVariableStatsForProductAndYear(prodId int, year int, version string) []models.ProductModelPointVariableStats {
//	prod, err := GetProductById(prodId)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	prodCode := prod.ProductCode
//	table := strings.ToLower(prodCode) + "_modelpoints"
//
//	variables := []string{"age_at_entry", "main_member_age_at_entry", "gender", "main_member_gender", "sum_assured", "outstanding_loan", "instalment", "annual_premium", "premium_rate", "interest", "term", "outstanding_term_months"}
//	var allStats []models.ProductModelPointVariableStats
//	DB.Where("product_code = ? and year = ? and version = ?", prodCode, year, version).Find(&allStats)
//	if len(allStats) > 0 {
//		return allStats
//	}
//	for _, v := range variables {
//		var stats models.ProductModelPointVariableStats
//		if v != "gender" && v != "main_member_gender" {
//			err := DB.Raw(fmt.Sprintf("select MIN(%s) as min, MAX(%s) as max, SUM(%s) as sum, AVG(%s) as average from %s where year = %d and mp_version='%s'", v, v, v, v, table, year, version)).Scan(&stats).Error
//			if err != nil {
//				fmt.Println(err)
//			}
//			err = DB.Raw(fmt.Sprintf("select count(%s) as number_of_zeroes from %s where %s = 0 and year = %d and mp_version='%s'", v, table, v, year, version)).Scan(&stats).Error
//			err = DB.Raw(fmt.Sprintf("select count(%s) as number_of_lives from  %s where year = %d and mp_version='%s'", v, table, year, version)).Scan(&stats).Error
//			stats.ProductCode = prodCode
//			stats.Variable = v
//			stats.Year = year
//			stats.RiskRateCode = version
//		} else {
//			//SELECT count(IF(gender = 'M',1,NULL)) as male, count(IF(gender = 'F',1,NULL)) as female FROM
//			//SELECT count(*) filter(where gender = 'M') as male, count(*) filter (where gender = 'F') as female FROM fe_001_modelpoints
//			err = DB.Raw(fmt.Sprintf("SELECT count(%s) as male FROM %s where year = %d and gender='M' and mp_version = '%s'", v, table, year, version)).Scan(&stats).Error
//			err = DB.Raw(fmt.Sprintf("SELECT count(%s) as female FROM %s where year = %d and gender='F' and mp_version = '%s'", v, table, year, version)).Scan(&stats).Error
//			err = DB.Raw(fmt.Sprintf("SELECT count(%s) as number_of_lives FROM %s where year = %d and mp_version = '%s'", v, table, year, version)).Scan(&stats).Error
//			stats.ProductCode = prodCode
//			stats.Variable = v
//			stats.Year = year
//			stats.RiskRateCode = version
//		}
//		//DB.Where("product_code = ? and variable = ? and year = ?", prodCode, v, year).Delete(&stats)
//		DB.Where("product_code = ? and variable = ? and year = ? and version = ?", prodCode, v, year, version).FirstOrCreate(&stats)
//
//		allStats = append(allStats, stats)
//	}
//
//	return allStats
//}

func GetModelPointVersionsForProductAndYear(prodId int, year int) ([]string, error) {
	prod, err := GetProductById(prodId)
	if err != nil {
		return nil, err
	}
	tableName := strings.ToLower(prod.ProductCode) + "_modelpoints"
	var versions []string
	err = DB.Table(tableName).Distinct("mp_version").Where("year=?", year).Pluck("mp_version", &versions).Error
	if err != nil {
		return nil, err
	}
	return versions, nil
}

func GetModelPointVariableStatsForProduct(prodId int, runId int) []models.ProductModelPointVariableStats {
	prod, err := GetProductById(prodId)
	if err != nil {
		fmt.Println(err)
	}
	jobProduct := GetJobProduct(runId)
	run, err := GetProjectionJobRun(jobProduct.ProjectionJobID)
	if err != nil {
		fmt.Println(err)
	}

	prod_code := prod.ProductCode
	table := strings.ToLower(prod_code) + "_modelpoints"

	variables := []string{"age_at_entry", "main_member_age_at_entry", "gender", "main_member_gender", "sum_assured", "outstanding_loan", "instalment", "annual_premium", "premium_rate", "interest", "term", "outstanding_term_months"}
	var allStats []models.ProductModelPointVariableStats

	DB.Where("product_code = ? and year = ? and version = ?", prod_code, run.ModelpointYear, run.ModelPointVersion).Find(&allStats)
	if len(allStats) > 0 {
		return allStats
	}

	baseWhere := fmt.Sprintf("year = %d AND mp_version = '%s'", run.ModelpointYear, run.ModelPointVersion)

	for _, v := range variables {
		currentStats := models.ProductModelPointVariableStats{
			ProductCode: prod_code,
			Variable:    v,
			Year:        run.ModelpointYear,
			Version:     run.ModelPointVersion,
		}

		if v != "gender" && v != "main_member_gender" {
			var mmasResult struct {
				Min     *float64 `gorm:"column:min"`
				Max     *float64 `gorm:"column:max"`
				Sum     *float64 `gorm:"column:sum"`
				Average *float64 `gorm:"column:average"`
			}
			if err := DB.Raw(fmt.Sprintf("SELECT MIN(%s) AS min, MAX(%s) AS max, SUM(%s) AS sum, AVG(%s) AS average FROM %s WHERE %s", v, v, v, v, table, baseWhere)).Scan(&mmasResult).Error; err != nil {
				fmt.Println(err)
			}
			if mmasResult.Min != nil {
				currentStats.Min = *mmasResult.Min
			}
			if mmasResult.Max != nil {
				currentStats.Max = *mmasResult.Max
			}
			if mmasResult.Sum != nil {
				currentStats.Sum = *mmasResult.Sum
			}
			if mmasResult.Average != nil {
				currentStats.Average = *mmasResult.Average
			}

			var zeroesResult struct {
				NumberOfZeroes int64 `gorm:"column:number_of_zeroes"`
			}
			if err := DB.Raw(fmt.Sprintf("SELECT COUNT(%s) AS number_of_zeroes FROM %s WHERE %s = 0 AND %s", v, table, v, baseWhere)).Scan(&zeroesResult).Error; err != nil {
				fmt.Println(err)
			}
			currentStats.NumberOfZeroes = int(zeroesResult.NumberOfZeroes)

			var livesResult struct {
				NumberOfLives int64 `gorm:"column:number_of_lives"`
			}
			if err := DB.Raw(fmt.Sprintf("SELECT COUNT(%s) AS number_of_lives FROM %s WHERE %s IS NOT NULL AND %s", v, table, v, baseWhere)).Scan(&livesResult).Error; err != nil {
				fmt.Println(err)
			}
			currentStats.NumberOfLives = int(livesResult.NumberOfLives)
		} else {
			genderCol := v
			var maleResult struct {
				Male int64 `gorm:"column:male"`
			}
			if err := DB.Raw(fmt.Sprintf("SELECT COUNT(*) AS male FROM %s WHERE %s = 'M' AND %s", table, genderCol, baseWhere)).Scan(&maleResult).Error; err != nil {
				fmt.Println(err)
			}
			currentStats.Male = float64(maleResult.Male)

			var femaleResult struct {
				Female int64 `gorm:"column:female"`
			}
			if err := DB.Raw(fmt.Sprintf("SELECT COUNT(*) AS female FROM %s WHERE %s = 'F' AND %s", table, genderCol, baseWhere)).Scan(&femaleResult).Error; err != nil {
				fmt.Println(err)
			}
			currentStats.Female = float64(femaleResult.Female)

			var livesResult struct {
				NumberOfLives int64 `gorm:"column:number_of_lives"`
			}
			if err := DB.Raw(fmt.Sprintf("SELECT COUNT(*) AS number_of_lives FROM %s WHERE %s", table, baseWhere)).Scan(&livesResult).Error; err != nil {
				fmt.Println(err)
			}
			currentStats.NumberOfLives = int(livesResult.NumberOfLives)
		}
		allStats = append(allStats, currentStats)
	}

	return allStats

}

func GetModelPointVariableStatsForPortfolio(portfolioName string, year int, version string) []models.PaaModelPointVariableStats {
	allStats := GeneratePAAModelPointVariableStats(portfolioName, year, version)
	return allStats
}

func GeneratePAAModelPointVariableStats(portfolioName string, year int, version string) []models.PaaModelPointVariableStats {
	var allStats []models.PaaModelPointVariableStats
	// Attempt to retrieve existing stats first
	err := DB.Where("portfolio_name = ? and year = ? and mp_version = ?", portfolioName, year, version).Find(&allStats).Error
	if err != nil {
		fmt.Println(fmt.Errorf("error fetching existing stats: %w", err))
		// Decide if you want to proceed or return an error
	}

	if len(allStats) > 0 {
		return allStats
	}

	table := "modified_gmm_model_points"
	variables := []string{"sub_product_code", "term_months", "locked_in_year", "locked_in_month", "ifrs17_group", "ifrs17_group_treaty1", "ifrs17_group_treaty2", "ifrs17_group_treaty3", "annual_premium", "frequency", "duration_in_force_months", "original_loan", "annual_interest_rate", "monthly_instalment", "outstanding_loan_term_months"}

	// Pre-compile base query parts for WHERE clause for clarity and minor efficiency
	whereClauseBase := fmt.Sprintf("year = %d AND paa_portfolio_name = '%s' AND mp_version = '%s'", year, portfolioName, version)

	for _, v := range variables {
		// Initialize a new stats object for each variable
		currentStats := models.PaaModelPointVariableStats{
			PortfolioName: portfolioName,
			Year:          year,
			MpVersion:     version,
			Variable:      v,
		}
		var queryErr error

		// Number of zeroes
		if v != "ifrs17_group" && v != "sub_product_code" {
			var countResult struct {
				NumberOfZeroes int64 `gorm:"column:number_of_zeroes"`
			} // Use specific type
			query := fmt.Sprintf("SELECT COUNT(%s) AS number_of_zeroes FROM %s WHERE %s = 0 AND %s", v, table, v, whereClauseBase)
			queryErr = DB.Raw(query).Scan(&countResult).Error
			if queryErr != nil {
				fmt.Println(fmt.Errorf("error counting zeroes for %s: %w", v, queryErr))
			}
			currentStats.NumberOfZeroes = int(countResult.NumberOfZeroes) // Assign from temp struct
		}

		// Number of empty values
		if v == "ifrs17_group" {
			var emptyResult struct {
				EmptyValues int64 `gorm:"column:empty_values"`
			}
			query := fmt.Sprintf("SELECT COUNT(%s) AS empty_values FROM %s WHERE %s = '' AND %s", v, table, v, whereClauseBase)
			queryErr = DB.Raw(query).Scan(&emptyResult).Error
			if queryErr != nil {
				fmt.Println(fmt.Errorf("error counting empty values for %s: %w", v, queryErr))
			}
			currentStats.EmptyValues = int(emptyResult.EmptyValues) // Assign
		}

		// Min, Max, Average values
		if v != "sub_product_code" && v != "ifrs17_group" {
			// Use pointers or sql.NullFloat64 if MIN/MAX/AVG can be NULL to distinguish from actual zero.
			// Assuming PaaModelPointVariableStats fields Min, Max, Average are float64.
			var mmaResult struct {
				Min     *float64 `gorm:"column:min"`     // Pointer to handle potential NULL
				Max     *float64 `gorm:"column:max"`     // Pointer to handle potential NULL
				Average *float64 `gorm:"column:average"` // Pointer to handle potential NULL
			}
			query := fmt.Sprintf("SELECT MIN(%s) AS min, MAX(%s) AS max, AVG(%s) AS average FROM %s WHERE %s", v, v, v, table, whereClauseBase)
			queryErr = DB.Raw(query).Scan(&mmaResult).Error
			if queryErr != nil {
				fmt.Println(fmt.Errorf("error calculating min/max/avg for %s: %w", v, queryErr))
			}
			if mmaResult.Min != nil {
				currentStats.Min = *mmaResult.Min
			} // else currentStats.Min remains its zero value (0.0), or handle as needed
			if mmaResult.Max != nil {
				currentStats.Max = *mmaResult.Max
			}
			if mmaResult.Average != nil {
				currentStats.Average = *mmaResult.Average
			}
		}

		// Total count
		var totalCountResult struct {
			TotalCount int64 `gorm:"column:total_count"`
		}
		query := fmt.Sprintf("SELECT COUNT(%s) AS total_count FROM %s WHERE %s", v, table, whereClauseBase)
		queryErr = DB.Raw(query).Scan(&totalCountResult).Error
		if queryErr != nil {
			fmt.Println(fmt.Errorf("error getting total count for %s: %w", v, queryErr))
		}
		currentStats.TotalCount = int(totalCountResult.TotalCount)

		// Distinct values count
		var distinctCountResult struct {
			DistinctValues int64 `gorm:"column:distinct_values"`
		}
		query = fmt.Sprintf("SELECT COUNT(DISTINCT(%s)) AS distinct_values FROM %s WHERE %s", v, table, whereClauseBase)
		queryErr = DB.Raw(query).Scan(&distinctCountResult).Error
		if queryErr != nil {
			fmt.Println(fmt.Errorf("error getting distinct count for %s: %w", v, queryErr))
		}
		currentStats.DistinctValues = int(distinctCountResult.DistinctValues)

		allStats = append(allStats, currentStats)
	}

	// Delete existing stats for this portfolio, year, version before saving new ones
	// Ensure this operation is safe and intended for your workflow.
	// Consider if this delete should be part of a transaction with the Create.
	deleteErr := DB.Where("portfolio_name = ? AND year = ? AND mp_version = ?", portfolioName, year, version).Delete(&models.PaaModelPointVariableStats{}).Error
	if deleteErr != nil {
		fmt.Println(fmt.Errorf("error deleting old stats: %w", deleteErr))
		// Decide if you should proceed with Create if Delete fails
	}

	// Save all newly computed stats
	if len(allStats) > 0 { // Only create if there are stats to save
		createErr := DB.Create(&allStats).Error
		if createErr != nil {
			fmt.Println(fmt.Errorf("error saving new stats: %w", createErr))
		}
	}
	return allStats
}

func ProductCodeExists(prodCode string) bool {
	var prod models.Product
	DB.Where("product_code = ?", prodCode).First(&prod)
	if prod.ProductCode == prodCode {
		return true
	} else {
		return false
	}
}

func ProductNameExists(prodName string) bool {
	var prod models.Product
	err := DB.Where("product_name = ?", prodName).Find(&prod).Error
	if err != nil {
		fmt.Println(err)
	}
	if prod.ProductName == prodName {
		return true
	} else {
		return false
	}
}

func SavePricingModelPoints(v *multipart.FileHeader, product models.Product) error {
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

	if (len(dec.Header()) != reflect.TypeOf(models.ProductPricingModelPoint{}).NumField()-1) { //account for id field
		fmt.Println(len(dec.Header()), reflect.TypeOf(models.ProductPricingModelPoint{}).NumField())
		return errors.New("the number of columns in the file does not match the number of columns in the database")
	}

	if !utils.MatchCSVTags(dec.Header(), models.ProductPricingModelPoint{}) {
		return errors.New("the column names in the file does not match the column names in the database")
	}

	tableName := strings.ToLower(product.ProductCode) + "_pricing_modelpoints"
	//Delete existing data
	//err = DB.Table(tableName).Where("product_code = ?", product.ProductCode).Delete(models.ProductPricingModelPoint{}).Error
	if err != nil {
		fmt.Println(err)
	}

	var pps []models.ProductPricingModelPoint
	for {
		var pp models.ProductPricingModelPoint
		if err := dec.Decode(&pp); err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}
		pp.ProductCode = product.ProductCode
		DB.Table(tableName).Where("pricing_mp_version = ?", pp.PricingMpVersion).Delete(models.ProductPricingModelPoint{})
		pps = append(pps, pp)

		//err = DB.Table(tableName).Create(&pp).Error
		//if err != nil {
		//	log.Error().Msg(err.Error())
		//}
	}
	err = DB.Table(tableName).CreateInBatches(&pps, 200).Error

	return nil
}

func DeletePricingParameter(productId int) error {
	var product models.Product
	err := DB.Where("id = ?", productId).First(&product).Error
	if err != nil {
		return err
	}

	tableName := "pricing_parameters"
	err = DB.Table(tableName).Where("product_code = ?", product.ProductCode).Delete(models.PricingParameter{}).Error
	if err != nil {
		return err
	}

	return nil
}

func DeletePricingModelPoints(productCode, mpVersion string) error {

	tableName := strings.ToLower(productCode) + "_pricing_modelpoints"
	err := DB.Table(tableName).Where("pricing_mp_version = ?", mpVersion).Delete(models.ProductPricingModelPoint{}).Error
	if err != nil {
		return err
	}

	return nil
}

func SaveModelPoints(v *multipart.FileHeader, product models.Product, year int, mpVersion string) error {
	var delimiter rune
	delimiterFile, err := v.Open()
	if err != nil {
		return err
	}
	defer delimiterFile.Close()
	delimiter, err = utils.GetDelimiter(delimiterFile)

	startTime := time.Now()

	file, err := v.Open()
	if err != nil {
		return err
	}

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.Comma = delimiter
	dec, _ := csvutil.NewDecoder(reader)
	dec.Header()

	tableName := strings.ToLower(product.ProductCode) + "_modelpoints"
	err = DB.Table(tableName).Where("year = ? and mp_version = ?", year, mpVersion).Delete(&models.ProductModelPoint{}).Error
	if err != nil {
		fmt.Println(err)
	}

	var pps []models.ProductModelPoint
	for {
		var pp models.ProductModelPoint
		if err := dec.Decode(&pp); err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}
		if pp.AgeAtEntry != 0 || pp.Gender != "" || pp.MainMemberGender != "" || pp.MainMemberAgeAtEntry != 0 {
			pp.Year = year
			pp.ProductCode = product.ProductCode
			pp.MpVersion = mpVersion
			pp.CreationDate = time.Now()
			pps = append(pps, pp)
		}
	}
	err = DB.Table(tableName).CreateInBatches(&pps, 200).Error
	if err != nil {
		log.Error().Msg(err.Error())
		if strings.Contains(err.Error(), "Error 1054: Unknown column") {
			segments := strings.SplitN(err.Error(), "Unknown column '", 2)
			fmt.Println(segments[1])
			columnName := strings.SplitN(segments[1], "'", 2)[0]
			return errors.New(fmt.Sprintf("there is an unknown column '%s' in the field list. Please check the file and try again", columnName))
		}
	}

	// Update ProductModelPointCount cache for this product/year/version
	// Compute the count we just inserted to keep cache in sync
	var insertedCount int64
	if cerr := DB.Table(tableName).Where("year = ? AND mp_version = ?", year, mpVersion).Count(&insertedCount).Error; cerr == nil {
		// Upsert-like behavior: remove any existing and insert fresh
		_ = DB.Where("product_id = ? AND year = ? AND version = ?", product.ID, year, mpVersion).Delete(&models.ProductModelPointCount{}).Error
		_ = DB.Create(&models.ProductModelPointCount{ProductId: product.ID, ProductCode: product.ProductCode, Year: year, Version: mpVersion, Count: int(insertedCount)}).Error
	}

	endTime := time.Since(startTime)

	fmt.Println("Time to process: ", endTime.Seconds())

	return err
}

// SaveModelPointsFromPath is a variant that reads CSV data from a local file path.
// It mirrors SaveModelPoints logic but allows background processing after the HTTP request returns.
func SaveModelPointsFromPath(path string, product models.Product, year int, mpVersion string) error {
	startTime := time.Now()

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// We need a separate reader for delimiter detection; open again
	f2, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f2.Close()

	delimiter, err := utils.GetDelimiter(f2)
	if err != nil {
		// proceed with default if GetDelimiter returns error but provides no delimiter
	}

	reader := csv.NewReader(utfbom.SkipOnly(f))
	reader.TrimLeadingSpace = true
	reader.Comma = delimiter
	dec, _ := csvutil.NewDecoder(reader)
	dec.Header()

	tableName := strings.ToLower(product.ProductCode) + "_modelpoints"
	if derr := DB.Table(tableName).Where("year = ? and mp_version = ?", year, mpVersion).Delete(&models.ProductModelPoint{}).Error; derr != nil {
		fmt.Println(derr)
	}

	var pps []models.ProductModelPoint
	for {
		var pp models.ProductModelPoint
		if err := dec.Decode(&pp); err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}
		if pp.AgeAtEntry != 0 || pp.Gender != "" || pp.MainMemberGender != "" || pp.MainMemberAgeAtEntry != 0 {
			pp.Year = year
			pp.ProductCode = product.ProductCode
			pp.MpVersion = mpVersion
			pp.CreationDate = time.Now()
			pps = append(pps, pp)
		}
	}
	if ierr := DB.Table(tableName).CreateInBatches(&pps, 200).Error; ierr != nil {
		log.Error().Msg(ierr.Error())
		if strings.Contains(ierr.Error(), "Error 1054: Unknown column") {
			segments := strings.SplitN(ierr.Error(), "Unknown column '", 2)
			fmt.Println(segments[1])
			columnName := strings.SplitN(segments[1], "'", 2)[0]
			return errors.New(fmt.Sprintf("there is an unknown column '%s' in the field list. Please check the file and try again", columnName))
		}
	}

	// Update count cache
	var insertedCount int64
	if cerr := DB.Table(tableName).Where("year = ? AND mp_version = ?", year, mpVersion).Count(&insertedCount).Error; cerr == nil {
		_ = DB.Where("product_id = ? AND year = ? AND version = ?", product.ID, year, mpVersion).Delete(&models.ProductModelPointCount{}).Error
		_ = DB.Create(&models.ProductModelPointCount{ProductId: product.ID, ProductCode: product.ProductCode, Year: year, Version: mpVersion, Count: int(insertedCount)}).Error
	}

	fmt.Println("Time to process: ", time.Since(startTime).Seconds())
	return nil
}

func CreateReportingDisclosures(productCode string, year int) {

}

func SavePricingModelPointsBak(bytes []byte, product models.Product, prodMps []models.ProductModelpointVariable) error {
	var dat []map[string]interface{}
	if err := json.Unmarshal(bytes, &dat); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(len(dat[0]))
	tableName := strings.ToLower(product.ProductCode) + "_pricing_modelpoints"
	DB.Exec(fmt.Sprintf("delete from %s", tableName))
	for _, mp := range dat {
		insertQuery := "INSERT INTO %s(%s) values(%s)"
		queryString := ""
		fields := ""
		values := ""
		for k, val := range mp {
			//There's a key and a value and the purpose is to determine the data type of the value.
			kInflated := inflate(k)
			var productMp models.ProductModelpointVariable
			for _, mp := range prodMps {
				if mp.Code == kInflated {
					productMp = mp
					break
				}
			}

			switch productMp.Type {
			case "Integer":
				//val = int(val.(float64))
			case "String":
			case "Boolean":
				//val = val.(float64)
			case "Float":
				//val = val.(float64)

			}
			if k == "policy_number" {
				fmt.Println("Its a policy!!", k)
			}
			fields += strings.TrimSpace(k) + ","

			if reflect.TypeOf(val).Kind() == reflect.Float64 {
				if k == "weighting" {
					values += fmt.Sprintf("%.10f,", val)
				} else {
					values += fmt.Sprintf("%.0f,", val)
				}

			} else {
				if k == "product_code" {
					val = product.ProductCode //Special case for file uploads
				}
				values += "'" + val.(string) + "',"
			}

			fmt.Println(productMp.ID)

			//fmt.Printf("%s -> %s\n", k, v)
		}

		fields = strings.Trim(fields, ",")
		fields = strings.Replace(fields, " ", "", -1)
		removeBOM(fields)
		values = strings.Trim(values, ",")

		queryString = fmt.Sprintf(insertQuery, tableName, fields, values)

		fmt.Println(queryString)
		dbErr := DB.Exec(queryString).Error

		if dbErr != nil {
			log.Error().Err(dbErr).Send()
			fmt.Println(dbErr)
			if strings.Contains(dbErr.Error(), "1064") {
				return errors.New("the uploaded pricing modelpoint file does not conform to the configured structure")
			} else {
				return dbErr
			}
		}
	}
	return nil
}

func GetTableStructure(tableName string) []string {
	rows, err := DB.Table(tableName).Limit(1).Rows()
	if err != nil || rows == nil {
		return []string{}
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return []string{}
	}
	columns = remove(columns, "year")
	return columns
}

func GetAssociatedTableStructure(tableId int) []string {
	var prodTable models.ProductTable
	var tableName string
	//We need the product code...

	DB.Where("id = ?", tableId).Find(&prodTable)

	product, err := GetProductById(prodTable.ProductID)
	if err != nil {
		fmt.Println(err)
	}
	var rows *sql.Rows
	var e error
	if prodTable.Class == "TransitionRates" {
		tableName = strings.ToLower(product.ProductCode) + "_" + strings.ToLower(prodTable.Name)
		rows, e = DB.Table(tableName).Limit(1).Rows()
	} else {
		tableName = "product_" + strings.ToLower(prodTable.Name)
		rows, e = DB.Table(tableName).Where("product_code = ?", product.ProductCode).Limit(1).Rows()
	}
	if e != nil || rows == nil {
		return []string{}
	}
	defer rows.Close()
	columns, err2 := rows.Columns()
	if err2 != nil {
		return []string{}
	}
	fmt.Println(columns)
	columns = remove(columns, "year")
	return columns
}

func GetYieldCurveYears() []int {
	var years []int
	DB.Raw("select distinct year from yield_curve").Scan(&years)
	return years
}

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func removeBOM(fields string) {
	output, err := ioutil.ReadAll(utfbom.SkipOnly(bytes.NewReader([]byte(fields))))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(output))
}
func inflate(s string) string {
	s = strings.Replace(s, "_at_", "@", -1)
	s = strings.ToUpper(s)
	s = strings.TrimSpace(s)
	return s
}
