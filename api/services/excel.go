package services

import (
	"api/models"
	"api/utils"
	"fmt"
	"github.com/xuri/excelize/v2"
	"reflect"
	"strconv"
)

func MakeExcelFile(outputPath string, projections []models.AggregatedProjection) {
	var aProjection models.AggregatedProjection
	var fieldHeaders []string
	var err error
	var excelColumns = utils.GenerateColumns()

	s := reflect.ValueOf(&aProjection).Elem()
	for i := 0; i < s.NumField(); i++ {
		fieldHeaders = append(fieldHeaders, s.Type().Field(i).Name)
		//fmt.Print(s.Type().Field(i).Status, " ")
	}
	f := excelize.NewFile()
	for k, v := range fieldHeaders {
		f.SetCellValue("Sheet1", excelColumns[k]+"1", v)
	}

	var rowNumber int = 2

	for _, v := range projections {
		s := reflect.ValueOf(&v).Elem()
		for i := 0; i < s.NumField(); i++ {
			err = f.SetCellValue("Sheet1", excelColumns[i]+strconv.Itoa(rowNumber), s.Field(i).Interface())
			if err != nil {
				fmt.Println(err)
			}
		}
		rowNumber++
	}

	err = f.SaveAs(outputPath)
	if err != nil {
		fmt.Println(err)
	}
}

func GetAOSStepResultsExcelByProductCode(outputPath string, prodCode string) {
	var aosResults []models.AOSStepResult
	err := DB.Where("product_code = ?", prodCode).Find(&aosResults).Error
	if err != nil {
		fmt.Println(err)
	}

	var aosResult models.AOSStepResult

	var fieldHeaders []string
	var excelColumns = utils.GenerateColumns()

	s := reflect.ValueOf(&aosResult).Elem()
	for i := 0; i < s.NumField(); i++ {
		fieldHeaders = append(fieldHeaders, s.Type().Field(i).Name)
		//fmt.Print(s.Type().Field(i).Status, " ")
	}
	f := excelize.NewFile()
	for k, v := range fieldHeaders {
		f.SetCellValue("Sheet1", excelColumns[k]+"1", v)
	}

	var rowNumber int = 2

	for _, v := range aosResults {
		s := reflect.ValueOf(&v).Elem()
		for i := 0; i < s.NumField(); i++ {
			err = f.SetCellValue("Sheet1", excelColumns[i]+strconv.Itoa(rowNumber), s.Field(i).Interface())
			if err != nil {
				fmt.Println(err)
			}
		}
		rowNumber++
	}

	err = f.SaveAs(outputPath)
	if err != nil {
		fmt.Println(err)
	}
}

func GetAOSStepResultsExcelByCsmRun(runId string) ([]byte, error) {
	var err error
	id := utils.StringToInt(runId)
	dQuery := fmt.Sprintf("select * from aos_step_results where csm_run_id = %d", id)

	excelData, err := exportTableToExcel(dQuery)
	if err != nil {
		return nil, err
	}
	return excelData, nil

}

func GetAOSStepResultsExcelByGroup(outputPath string, group string) {
	var aosResults []models.AOSStepResult
	err := DB.Where("ifrs17_group = ?", group).Find(&aosResults).Error
	if err != nil {
		fmt.Println(err)
	}

	var aosResult models.AOSStepResult

	var fieldHeaders []string
	var excelColumns = utils.GenerateColumns()

	s := reflect.ValueOf(&aosResult).Elem()
	for i := 0; i < s.NumField(); i++ {
		fieldHeaders = append(fieldHeaders, s.Type().Field(i).Name)
		//fmt.Print(s.Type().Field(i).Status, " ")
	}
	f := excelize.NewFile()
	for k, v := range fieldHeaders {
		f.SetCellValue("Sheet1", excelColumns[k]+"1", v)
	}

	var rowNumber int = 2

	for _, v := range aosResults {
		s := reflect.ValueOf(&v).Elem()
		for i := 0; i < s.NumField(); i++ {
			err = f.SetCellValue("Sheet1", excelColumns[i]+strconv.Itoa(rowNumber), s.Field(i).Interface())
			if err != nil {
				fmt.Println(err)
			}
		}
		rowNumber++
	}

	err = f.SaveAs(outputPath)
	if err != nil {
		fmt.Println(err)
	}
}

func MakePricingFile(outputPath string, projections []models.ModelPointPricing) {
	var pricingPoint models.ModelPointPricing
	var fieldHeaders []string
	var err error
	var excelColumns = utils.GenerateColumns()

	s := reflect.ValueOf(&pricingPoint).Elem()
	for i := 0; i < s.NumField(); i++ {
		fieldHeaders = append(fieldHeaders, s.Type().Field(i).Name)
		//fmt.Print(s.Type().Field(i).Status, " ")
	}
	f := excelize.NewFile()
	for k, v := range fieldHeaders {
		f.SetCellValue("Sheet1", excelColumns[k]+"1", v)
	}

	var rowNumber int = 2

	for _, v := range projections {
		s := reflect.ValueOf(&v).Elem()
		for i := 0; i < s.NumField(); i++ {
			err = f.SetCellValue("Sheet1", excelColumns[i]+strconv.Itoa(rowNumber), s.Field(i).Interface())
			if err != nil {
				fmt.Println(err)
			}
		}
		rowNumber++
	}

	err = f.SaveAs(outputPath)
	if err != nil {
		fmt.Println(err)
	}
}

func MakeControlExcelFile(outputPath string, projections []models.Projection) {
	var aProjection models.Projection
	var fieldHeaders []string
	var excelColumns = utils.GenerateColumns()
	var err error

	s := reflect.ValueOf(&aProjection).Elem()
	for i := 0; i < s.NumField(); i++ {
		fieldHeaders = append(fieldHeaders, s.Type().Field(i).Name)
		//fmt.Print(s.Type().Field(i).Status, " ")
	}
	f := excelize.NewFile()
	for k, v := range fieldHeaders {
		err := f.SetCellValue("Sheet1", excelColumns[k]+"1", v)
		fmt.Println(err)
	}

	var rowNumber int = 2

	for _, v := range projections {
		s := reflect.ValueOf(&v).Elem()
		for i := 0; i < s.NumField(); i++ {
			err = f.SetCellValue("Sheet1", excelColumns[i]+strconv.Itoa(rowNumber), s.Field(i).Interface())
			if err != nil {
				fmt.Println(err)
			}
		}
		rowNumber++
	}

	err = f.SaveAs(outputPath)
	if err != nil {
		fmt.Println(err)
	}
}

func MakePricingControlExcelFile(outputPath string, projections []models.PricingPoint) {
	var points models.PricingPoint
	var fieldHeaders []string
	var err error
	var columns = utils.GenerateColumns()

	s := reflect.ValueOf(&points).Elem()
	for i := 0; i < s.NumField(); i++ {
		fieldHeaders = append(fieldHeaders, s.Type().Field(i).Name)
		//fmt.Print(s.Type().Field(i).Status, " ")
	}
	f := excelize.NewFile()
	for k, v := range fieldHeaders {
		err := f.SetCellValue("Sheet1", columns[k]+"1", v)
		fmt.Println(err)
	}

	var rowNumber int = 2

	for _, v := range projections {
		s := reflect.ValueOf(&v).Elem()
		for i := 0; i < s.NumField(); i++ {
			err = f.SetCellValue("Sheet1", columns[i]+strconv.Itoa(rowNumber), s.Field(i).Interface())
			if err != nil {
				fmt.Println(err)
			}
		}
		rowNumber++
	}

	err = f.SaveAs(outputPath)
	if err != nil {
		fmt.Println(err)
	}
}
