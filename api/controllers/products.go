package controllers

import (
	"api/log"
	"api/models"
	"api/services"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetProductFamilies godoc
// @Summary Get a list of available product family types
// @Description Get available product family types
// @Accept json
// @Produce json
// @Success 200 {array} models.ProductFamily
// @Router /getproducts [get]
func GetProductFamilies(c *gin.Context) {
	products, err := services.GetProductFamilies()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "The requested product family was not found"})
	}

	c.JSON(http.StatusOK, products)
}

// GetMarkovStates godoc
// @Summary Get a list of available transition states
// @Description Get available Transition (Markov) states
// @Accept json
// @Produce json
// @Success 200 {array} models.MarkovState
// @Router /markovstates [get]
func GetMarkovStates(c *gin.Context) {
	states, err := services.GetMarkovStates()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "The requested states were not found"})
	}
	c.JSON(http.StatusOK, states)
}

func PostTransitionStates(c *gin.Context) {
	prodId := c.Param("prodId")
	var states []models.TransitionState
	c.Bind(&states)

	for _, state := range states {
		state.ProductCode = strings.ToUpper(prodId)
		services.DB.Save(&state)
	}

	c.JSON(http.StatusOK, states)

}

// CreateProduct godoc
// @Summary Upload a valid product configuration
// @Description Use this endpoint to upload a completed product configuration
// @Accept mpfd
// @Produce json
// @Success 201 {object} models.Product
// @Router /upload [post]
func CreateProduct(c *gin.Context) {

	log.SetReportCaller(true)

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	pc, err := services.CreateProduct(form)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	creator := form.Value["user"][0]
	reviewer := form.Value["reviewer"][0]

	var task models.Task
	task.ProductCode = pc.ProductCode
	task.ProductID = pc.ID
	task.ProductName = pc.ProductName
	task.Creator = creator
	task.Approver = reviewer
	task.Assignee = reviewer
	task.Status = "active"
	task.Comments = "Requesting a review and approval of product configuration for " + task.ProductName
	services.GenerateTask(task)
	c.JSON(http.StatusCreated, pc)
	//c.JSON(http.StatusBadRequest, pc)
}

func UpdateRatingTable(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form: " + err.Error()})
		return
	}
	values := form.Value
	tableIdString := values["table_id"][0]
	prodCode := values["product_code"][0]
	year, err := strconv.Atoi(values["year"][0])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year: " + err.Error()})
		return
	}
	month := 0
	yieldCurveCode := ""
	if values["month"] != nil {
		month, err = strconv.Atoi(values["month"][0])
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month: " + err.Error()})
			return
		}
	}

	if values["yield_curve_code"] != nil {
		yieldCurveCode = values["yield_curve_code"][0]
	}

	file := form.File["file"]
	fmt.Println(file)

	tableId, err := strconv.Atoi(tableIdString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid table_id: " + err.Error()})
		return
	}

	productTable := services.GetProductTable(tableId)
	fmt.Println(productTable.Name)
	ctx := context.WithValue(c, "keys", c.Keys)
	err = services.ProcessTable(ctx, productTable, prodCode, year, file[0], month, yieldCurveCode)
	if err != nil {
		fmt.Println("Error processing table:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process table: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Table updated successfully"})
}

func UploadBulkAssumptionTables(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	file := form.File["file"]
	values := form.Value
	assumptionType := values["assumption_type"][0]
	year := values["year"][0]

	var yieldCurveCode string
	var month string
	if assumptionType == "Yield Curve" {
		yieldCurveCode = values["yield_curve_code"][0]
		month = values["month"][0]
	}

	ctx := context.WithValue(c, "keys", c.Keys)

	var tableName = ""

	switch assumptionType {
	case "Parameters":
		tableName = "product_parameters"
	case "Yield Curve":
		tableName = "yield_curve"
	case "Margins":
		tableName = "product_margins"
	case "Commission Structures":
		tableName = "product_commission_structures"
	case "Shocks":
		tableName = "product_shocks"
	}

	if len(tableName) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No assumption type was supplied"})
		return
	}

	err = services.ProcessBulkParameterTables(ctx, file[0], tableName, year, yieldCurveCode, month)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

func DeleteProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	services.DeleteProduct(services.DB, id)
	c.JSON(http.StatusOK, nil)
}

func DeleteProductPricingTableData(c *gin.Context) {
	productId, _ := strconv.Atoi(c.Param("id"))
	tableId, _ := strconv.Atoi(c.Param("tableId"))

	services.DeleteProductPricingTableData(productId, tableId)
	c.JSON(http.StatusOK, nil)
}

func DeleteProductTableData(c *gin.Context) {
	productId, _ := strconv.Atoi(c.Param("id"))
	tableId, _ := strconv.Atoi(c.Param("tableId"))

	services.DeleteProductTableData(productId, tableId)
	c.JSON(http.StatusOK, nil)
}

func DeleteProductTableDatav2(c *gin.Context) {
	productId, _ := strconv.Atoi(c.Param("id"))
	tableId, _ := strconv.Atoi(c.Param("tableId"))
	year, _ := strconv.Atoi(c.Param("year"))

	err := services.DeleteProductTableDatav2(productId, tableId, year)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func ActivateProduct(c *gin.Context) {
	var body map[string]interface{}
	err := c.Bind(&body)
	if err != nil {
		fmt.Println(err)
	}

	err = services.ActivateProduct(body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func GetProductsAndFamilies(c *gin.Context) {
	products, err := services.GetProductsAndFamilies()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, products)
}

func GetAvailableProducts(c *gin.Context) {
	products, err := services.GetAllAvailableProducts()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, products)
}

func GetProductById(c *gin.Context) {
	prodId, err := strconv.Atoi(c.Param("id"))
	product, err := services.GetProductById(prodId)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"product": product})
		return
	}
	c.JSON(http.StatusBadRequest, nil)
}

func GetTableData(c *gin.Context) {
	reqbody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error(err)
	}
	var body map[string]interface{}
	err = json.Unmarshal(reqbody, &body)
	if err != nil {
		fmt.Println(err)
	}
	tableId := int(body["table_id"].(float64))
	results := services.GetProductTableSample(body["product_code"].(string), tableId)
	c.JSON(http.StatusOK, results)
}

func GetPricingTableData(c *gin.Context) {
	reqbody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error(err)
	}
	var body map[string]interface{}
	err = json.Unmarshal(reqbody, &body)
	if err != nil {
		fmt.Println(err)
	}

	tableId := int(body["table_id"].(float64))
	results := services.GetProductPricingTableSample(body["product_code"].(string), tableId)
	c.JSON(http.StatusOK, results)
}

func CheckProductCode(c *gin.Context) {
	prodCode := c.Param("id")
	found := services.ProductCodeExists(prodCode)
	c.JSON(http.StatusOK, gin.H{"exists": found})
}

func CheckProductName(c *gin.Context) {
	prodCode := c.Param("id")
	found := services.ProductNameExists(prodCode)
	c.JSON(http.StatusOK, gin.H{"exists": found})
}

func GetSelectedFeatures(c *gin.Context) {
	var features models.ProductFeatures
	c.BindJSON(&features)
	results := services.GetSelectedFeatures(features)
	c.JSON(http.StatusOK, results)
}

func GetRatingFactors(c *gin.Context) {
	results := services.GetRatingFactors()
	c.JSON(http.StatusOK, results)
}

func GetDistinctCommissionProductCodes(c *gin.Context) {
	results, err := services.GetDistinctCommissionProductCodes()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

func GetGlobalTableData(c *gin.Context) {
	tableName := c.Param("table_name")
	results, err := services.GetGlobalTableData(tableName)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, results)

}

func DeleteGlobalTableData(c *gin.Context) {
	tableName := c.Param("table_name")
	key := c.Param("key")

	services.DeleteGlobalTableData(tableName, key)
	c.JSON(http.StatusOK, "table data deleted")
}

func GetYieldCurveYears(c *gin.Context) {
	results := services.GetYieldCurveYears()
	c.JSON(http.StatusOK, results)
}

func GetProductTableYears(c *gin.Context) {
	productId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tableName := c.Param("table-name")
	dataType := c.Param("data-type")
	tableType := c.Param("table-type")

	var isPricingTable bool
	if dataType == "pricing" {
		isPricingTable = true
	} else {
		isPricingTable = false
	}
	results := services.GetProductTableYears(productId, tableName, tableType, isPricingTable)
	c.JSON(http.StatusOK, results)
}
