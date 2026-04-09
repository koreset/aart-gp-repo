package controllers

import (
	"api/models"
	"api/services"
	"api/utils"
	"compress/gzip"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func UploadModelPoints(c *gin.Context) {
	prodId, err := strconv.Atoi(c.Param("id"))
	product, err := services.GetProductById(prodId)
	if err != nil {
		log.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form"})
		return
	}

	year := strings.TrimSpace(form.Value["year"][0])
	mpVersion := strings.TrimSpace(form.Value["mp_version"][0])
	file := form.File["file"][0]
	if year == "" || mpVersion == "" || file == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing year, mp_version or file"})
		return
	}
	yr, _ := strconv.Atoi(year)

	// Prepare temp dir
	tmpDir := filepath.Join("tmp", "uploads")
	_ = os.MkdirAll(tmpDir, 0o755)

	// Generate upload id
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	uploadID := hex.EncodeToString(b)

	// Destination path
	origName := filepath.Base(file.Filename)
	destPath := filepath.Join(tmpDir, uploadID+"_"+origName)

	// Save uploaded file to disk (this ensures we can process it after response returns)
	if err := c.SaveUploadedFile(file, destPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save upload"})
		return
	}

	// Create status entry and start background processing
	total := file.Size
	// If the uploaded file is gzipped, set total to original (uncompressed) size using gzip ISIZE footer
	if strings.HasSuffix(strings.ToLower(origName), ".gz") || strings.HasSuffix(strings.ToLower(origName), ".gzip") {
		if f, err := os.Open(destPath); err == nil {
			if fi, _ := f.Stat(); fi != nil && fi.Size() >= 4 {
				if _, err := f.Seek(-4, io.SeekEnd); err == nil {
					var tail [4]byte
					if _, err := f.Read(tail[:]); err == nil {
						unzSize := binary.LittleEndian.Uint32(tail[:])
						if unzSize > 0 {
							total = int64(unzSize)
						}
					}
				}
			}
			_ = f.Close()
		}
	}
	services.CreateUploadStatus(uploadID, origName, product.ID, yr, mpVersion, total)

	go func(path string, fname string) {
		services.SetUploadProcessing(uploadID)
		start := time.Now()
		procPath := path
		// If gzip, decompress to a temp CSV file first
		if strings.HasSuffix(strings.ToLower(fname), ".gz") || strings.HasSuffix(strings.ToLower(fname), ".gzip") {
			csvPath := strings.TrimSuffix(path, filepath.Ext(path)) + ".csv"
			in, e1 := os.Open(path)
			if e1 != nil {
				services.SetUploadFailed(uploadID, e1)
				return
			}
			defer in.Close()
			gzr, e2 := gzip.NewReader(in)
			if e2 != nil {
				services.SetUploadFailed(uploadID, e2)
				return
			}
			defer gzr.Close()
			out, e3 := os.Create(csvPath)
			if e3 != nil {
				services.SetUploadFailed(uploadID, e3)
				return
			}
			defer out.Close()
			buf := make([]byte, 32*1024)
			for {
				n, er := gzr.Read(buf)
				if n > 0 {
					if _, ew := out.Write(buf[:n]); ew != nil {
						services.SetUploadFailed(uploadID, ew)
						return
					}
					services.AddUploadProgress(uploadID, int64(n))
				}
				if er == io.EOF {
					break
				}
				if er != nil {
					services.SetUploadFailed(uploadID, er)
					return
				}
			}
			procPath = csvPath
		}

		if e := services.SaveModelPointsFromPath(procPath, product, yr, mpVersion); e != nil {
			services.SetUploadFailed(uploadID, e)
			return
		}
		// generate stats asynchronously but not blocking status completion
		go services.CreateModelPointStatsForProductAndYear(product.ProductCode, yr, mpVersion)
		services.SetUploadCompleted(uploadID)
		// Cleanup temp files after a grace period
		go func(p1, p2 string) {
			time.Sleep(2 * time.Minute)
			_ = os.Remove(p1)
			if p2 != "" && p2 != p1 {
				_ = os.Remove(p2)
			}
		}(path, procPath)
		_ = start // for potential future timing metrics
	}(destPath, origName)

	c.JSON(http.StatusAccepted, gin.H{"upload_id": uploadID, "message": "File received. Import is processing in background."})
}

func UploadPricingModelPoints(c *gin.Context) {
	prodId, err := strconv.Atoi(c.Param("id"))
	product, err := services.GetProductById(prodId)
	if err != nil {
		log.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	form, err := c.MultipartForm()
	file := form.File["file"][0]
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	fmt.Println(form)
	// We must add weighting variable manually here for now.
	if err != nil {
		log.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	bytes, str := utils.ReadCSV(file)
	utils.SaveFile(bytes, str)

	err = services.SavePricingModelPoints(file, product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, nil)
}

func DeletePricingParameter(c *gin.Context) {
	prodId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error().Err(err).Send()
	}
	err = services.DeletePricingParameter(prodId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func DeletePricingModelPoints(c *gin.Context) {
	productCode := c.Param("prod_code")
	mpVersion := c.Param("mp_version")
	err := services.DeletePricingModelPoints(productCode, mpVersion)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func RunValuations(c *gin.Context) {
	//The calling user must be a valid license holder

	fmt.Println(c.Keys)

	user := c.Keys["user"].(models.AppUser)

	reqbody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Err(err).Send()
	}
	var runJobs models.RunJob

	fmt.Println(string(reqbody))
	err = json.Unmarshal(reqbody, &runJobs)
	if err != nil {
		fmt.Println(err)
	}
	runJobs.UserName = user.UserName
	runJobs.UserEmail = user.UserEmail
	runJobs.User = user

	err = services.SaveJobTemplateContent(string(reqbody), runJobs.JobsTemplateID)
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		services.RunValuations(runJobs)
	}()

	c.JSON(http.StatusOK, gin.H{"message": "The job is being processed. An email will be sent on completion"})
}

func DeleteJobTemplate(c *gin.Context) {

}

func RunPricing(c *gin.Context) {
	reqbody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Err(err).Send()
	}
	user := c.Keys["user"].(models.AppUser)

	var pricingRun models.PricingRun
	err = json.Unmarshal(reqbody, &pricingRun)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	go func(pricingRun models.PricingRun, user models.AppUser) {
		err = services.RunPricing(pricingRun, user)
		if err != nil {
			fmt.Println("Error occurred: ", err)
		}
	}(pricingRun, user)

	c.JSON(http.StatusOK, gin.H{"message": "The pricing job has been successfully submitted. A notification will follow on completion"})
}

func ReRunPricing(c *gin.Context) {
	runId, err := strconv.Atoi(c.Param("job_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something wrong with the ID"})
		return
	}
	user := c.Keys["user"].(models.AppUser)

	pricingRun := services.GetPricingRun(runId)

	// we need to delete the existing pricing run
	err = services.DeletePricingRun(runId)

	// zero the pricing run id
	pricingRun.ID = 0

	go func(pricingRun models.PricingRun, user models.AppUser) {
		err = services.RunPricing(pricingRun, user)
		if err != nil {
			fmt.Println("Error occurred: ", err)
		}
	}(pricingRun, user)

	c.JSON(http.StatusOK, gin.H{"message": "The pricing job has been successfully submitted. A notification will follow on completion"})
}

func GetAggregatedReserves(c *gin.Context) {

	reqbody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Err(err).Send()
	}
	var prods map[string]interface{}
	json.Unmarshal(reqbody, &prods)

	var jobIds []interface{}
	var monthrange int
	var runVariable string
	//var isJobIds bool
	if prods != nil {
		jobIds = prods["jobProductIds"].([]interface{})
		monthrange = int(prods["monthRange"].(float64))
		monthrange = 0
		runVariable = prods["variable"].(string)
		//isJobIds = prods["jobIds"].(bool)
	} else {
		//We get the most recent run from the db
		projectionJob := services.GetMostRecentJob()
		monthrange = 0
		runVariable = "reserves"
		//isJobIds = false
		for _, prod := range projectionJob.Products {
			jobIds = append(jobIds, prod.ID)
		}
	}

	fmt.Println(jobIds)

	//jobId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something wrong with the ID"})
	}
	reserveList := services.GetAllAggregatedReserves(jobIds, monthrange, runVariable)
	c.JSON(http.StatusOK, reserveList)
}

func GetAllValuationJobs(c *gin.Context) {
	jobs := services.GetJobs(0)
	c.JSON(http.StatusOK, jobs)
}

func GetAllValidValuationJobs(c *gin.Context) {
	jobs := services.GetValidJobs()
	c.JSON(http.StatusOK, jobs)
}

func GetAggResults(c *gin.Context) {
	var payload models.ExcelAggPayload
	err := c.BindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something wrong with the ID"})
	}
	productCode := payload.ProductCode
	fmt.Println(productCode)
	spCode := payload.SpCode
	jobProductId := payload.JobProductId
	jobId := payload.RunId
	variableGroupId := payload.VariableGroupId
	projections := services.GetExcelAggregatedProjections(jobId, jobProductId, spCode, variableGroupId)
	c.JSON(http.StatusOK, projections)
}

func GetProductsForValuationJob(c *gin.Context) {
	jobId, err := strconv.Atoi(c.Param("run-id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something wrong with the ID"})
	}
	products := services.GetProductsForJob(jobId)
	c.JSON(http.StatusOK, products)
}

func GetSpCodesForRunProduct(c *gin.Context) {
	jobId, err := strconv.Atoi(c.Param("run-id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something wrong with the ID"})
	}
	prodCode := c.Param("product-code")
	spcodes := services.GetSpCodesForRunProduct(jobId, prodCode)
	c.JSON(http.StatusOK, spcodes)
}

func GetSpCodesForRun(c *gin.Context) {
	runId, err := strconv.Atoi(c.Param("run-id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something wrong with the ID"})
	}
	spcodes := services.GetSpCodesForRun(runId)
	c.JSON(http.StatusOK, spcodes)
}

func GetAllPricingRuns(c *gin.Context) {
	jobs := services.GetAllPricingRuns()
	c.JSON(http.StatusOK, jobs)
}

func GetValuationJob(c *gin.Context) {
	spCode := c.Param("sp-code")

	jobId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something wrong with the ID"})
	}

	if spCode != "" {
		projections := services.GetAggregatedProjections(jobId, spCode)
		c.JSON(http.StatusOK, gin.H{"projections": projections})
		return
	}

	spcodes, err := services.GetSpCodes(jobId)
	//sort spcodes
	sort.Ints(spcodes)
	projections := services.GetAggregatedProjections(jobId, strconv.Itoa(spcodes[0]))
	scopedProjections := services.GetGMMScopedProjections(jobId)
	runSettings := services.GetJobRunSettings(jobId)
	runErrors := services.GetJobRunErrors(jobId)
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{"projections": projections,
		"scopedProjections": scopedProjections,
		"settings":          runSettings,
		"spcodes":           spcodes,
		"errors":            runErrors})
}

func RestartStalledJob(c *gin.Context) {
	jobId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something wrong with the ID"})
		return
	}

	err = services.RestartStalledJob(jobId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to restart job: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job restarted successfully"})
}

func GetValuationJobControl(c *gin.Context) {
	jobId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something wrong with the ID"})
	}
	var results = make(map[string]interface{})
	projections := services.GetProjections(jobId)
	results["projections"] = projections
	runErrors := services.GetJobRunErrors(jobId)
	results["errors"] = runErrors
	c.JSON(http.StatusOK, results)
}

func DeleteValuationJob(c *gin.Context) {
	jobId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something wrong with the ID"})
	}
	services.DeleteProjection(jobId)
	c.JSON(http.StatusOK, nil)
}

func DeleteValuationJobs(c *gin.Context) {
	var runIds []int

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Err(err).Send()
	}

	err = json.Unmarshal(body, &runIds)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(runIds)

	for _, runId := range runIds {
		services.DeleteProjection(runId)
	}

	c.JSON(http.StatusOK, nil)
}

func GetPricingJob(c *gin.Context) {
	jobId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something wrong with the ID"})
	}
	pricingRun := services.GetPricingRun(jobId)
	c.JSON(http.StatusOK, pricingRun)
}

func DeletePricingJob(c *gin.Context) {
	jobId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something wrong with the ID"})
	}
	err = services.DeletePricingRun(jobId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, nil)
}

func GetValuationJobExcel(c *gin.Context) {
	runId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something wrong with the ID"})
	}

	if strings.Contains(c.Request.URL.Path, "control") {
		projections := services.GetProjections(runId)
		services.MakeControlExcelFile("Job-Control.xlsx", projections)
		c.FileAttachment("Job-Control.xlsx", "job-control.xlsx")
		// delete the file after sending
		defer func() {
			err := os.Remove("Job-Control.xlsx")
			if err != nil {
				fmt.Println(err)
			}
		}()

	} else {
		excelData, err := services.GetAggregatedProjectionsForDownload(runId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to zip files: %v", err)})
			return
		}

		c.Header("Content-Disposition", `attachment; filename="export.xlsx"`)
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelData)
	}
}

func GetProjectionJobExcel(c *gin.Context) {
	jobId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something wrong with the ID"})
	}

	excelData, err := services.GetProjectionJobExcel(jobId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to zip files: %v", err)})
		return
	}

	c.Header("Content-Disposition", `attachment; filename="export.xlsx"`)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelData)
}

func GetProjectionJobExcelScoped(c *gin.Context) {
	jobId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something wrong with the ID"})
	}

	excelData, err := services.GetProjectionJobExcelScoped(jobId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to zip files: %v", err)})
		return
	}

	c.Header("Content-Disposition", `attachment; filename="export.xlsx"`)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelData)
}

func GetAOSStepResultsExcel(c *gin.Context) {
	prodCode := c.Param("prod_code")
	services.GetAOSStepResultsExcelByProductCode("product-aosresult.xlsx", prodCode)
	c.FileAttachment("product-aosresult.xlsx", "product-aosresult.xlsx")
}

func GetAOSStepResultsExcelByCsmId(c *gin.Context) {
	runId := c.Param("run_id")

	excelData, err := services.GetAOSStepResultsExcelByCsmRun(runId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to zip files: %v", err)})
		return
	}

	//services.GetAOSStepResultsExcelByCsmRun("ifrs17result.xlsx", runId)
	//c.FileAttachment("ifrs17result.xlsx", "ifrs17result.xlsx")

	c.Header("Content-Disposition", `attachment; filename="ifrs17-export.xlsx"`)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelData)

}

func GetAOSStepResultsByGroupExcel(c *gin.Context) {
	group := c.Param("group")
	services.GetAOSStepResultsExcelByGroup("product-aosresult.xlsx", group)
	c.FileAttachment("product-aosresult.xlsx", "product-aosresult.xlsx")
}

func GetJobsForProduct(c *gin.Context) {
	prodCode := c.Param("prodcode")
	jobs := services.GetJobsForProduct(prodCode)
	c.JSON(http.StatusOK, jobs)
}

func CheckValuationRunName(c *gin.Context) {
	runName := c.Param("run_name")
	found := services.ValuationRunNameDoesNotExist(runName)
	c.JSON(http.StatusOK, found)
}

func GetAvailableYieldYears(c *gin.Context) {
	results := services.GetAvailableYieldYears()
	c.JSON(http.StatusOK, results)
}
func GetYieldCurveCodes(c *gin.Context) {
	year := c.Param("year")
	results := services.GetYieldCurveCodes(year)
	c.JSON(http.StatusOK, results)
}

func GetYieldCurveMonths(c *gin.Context) {
	year := c.Param("year")
	code := c.Param("code")
	results := services.GetYieldCurveMonths(year, code)
	c.JSON(http.StatusOK, results)
}

func GetYieldCurveMonthsv2(c *gin.Context) {
	productCode := c.Param("product-code")
	yieldYear := c.Param("yield-year")
	parameterYear := c.Param("parameter-year")
	basis := c.Param("basis")

	results := services.GetYieldCurveMonthsv2(productCode, yieldYear, parameterYear, basis)
	c.JSON(http.StatusOK, results)
}

func DeleteYieldCurveData(c *gin.Context) {
	year := c.Param("year")
	code := c.Param("code")
	month := c.Param("month")
	err := services.DeleteYieldCurveData(year, code, month)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, nil)
}

func DeleteIbnrYieldCurveData(c *gin.Context) {
	year := c.Param("year")
	code := c.Param("code")
	month := c.Param("month")
	err := services.DeleteIbnrYieldCurveData(year, code, month)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, nil)
}

func GetAvailableIbnrYieldCurveYears(c *gin.Context) {
	results, err := services.GetAvailableIbnrYieldCurveYears()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, results)
}

func GetAvailableBases(c *gin.Context) {
	productCode := c.Param("prod-code")
	results := services.GetAvailableBases(productCode)
	c.JSON(http.StatusOK, results)
}

func GetAvailableYieldCurveMonths(c *gin.Context) {
	year := c.Param("year")
	code := c.Param("code")
	results := services.GetPAAYieldCurveMonths(year, code)
	c.JSON(http.StatusOK, results)
}

func GetAvailableIbnrYieldCurveMonths(c *gin.Context) {
	year := c.Param("year")
	basis := c.Param("basis")
	parameterYear := c.Param("parameter-year")
	portfolioId := c.Param("portfolio-id")
	results := services.GetIbnrYieldCurveMonths(year, basis, parameterYear, portfolioId)
	c.JSON(http.StatusOK, results)
}

func GetAvailableIbnrYieldCurveCodes(c *gin.Context) {
	year := c.Param("year")
	results := services.GetIbnrYieldCurveCodes(year)
	c.JSON(http.StatusOK, results)
}

func GetIBNRBasis(c *gin.Context) {
	portfolio := c.Param("portfolio")
	year := c.Param("year")
	results := services.GetIBNRBases(portfolio, year)
	c.JSON(http.StatusOK, results)
}

func GetIBNRShockBases(c *gin.Context) {
	results := services.GetIBNRShockBases()
	c.JSON(http.StatusOK, results)
}

func GetAvailableParameterYears(c *gin.Context) {
	prodCode := c.Param("prod-code")
	results := services.GetAvailableParameterYears(prodCode)
	c.JSON(http.StatusOK, results)
}

func GetAvailableMarginYears(c *gin.Context) {
	results := services.GetAvailableMarginYears()
	c.JSON(http.StatusOK, results)
}

func GetAvailableModelPointYears(c *gin.Context) {
	prodCode := c.Param("prod-code")
	results := services.GetAvailableModelPointYears(prodCode)
	c.JSON(http.StatusOK, results)
}

func GetAvailableLapseYears(c *gin.Context) {
	prodCode := c.Param("prod-code")
	results := services.GetAvailableLapseYears(prodCode)
	c.JSON(http.StatusOK, results)
}

func GetAvailableMortalityYears(c *gin.Context) {
	prodCode := c.Param("prod-code")
	results := services.GetAvailableMortalityYears(prodCode)
	c.JSON(http.StatusOK, results)
}

func GetAvailableRetrenchmentYears(c *gin.Context) {
	prodCode := c.Param("prod-code")
	results := services.GetAvailableRetrenchmentYears(prodCode)
	c.JSON(http.StatusOK, results)
}

func GetAvailableDisabilityYears(c *gin.Context) {
	prodCode := c.Param("prod-code")
	results := services.GetAvailableDisabilityYears(prodCode)
	c.JSON(http.StatusOK, results)
}

func GetAvailableLapseMarginYears(c *gin.Context) {
	prodCode := c.Param("prod-code")
	results := services.GetAvailableLapseMarginYears(prodCode)
	c.JSON(http.StatusOK, results)
}

func CreateShockSetting(c *gin.Context) {
	var shockSetting models.ShockSetting
	err := c.BindJSON(&shockSetting)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shockSetting, err = services.SaveShockSetting(shockSetting)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, shockSetting)
}

func GetShockSettings(c *gin.Context) {
	settings, err := services.GetShockSettings()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}

func DeleteShockSetting(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("setting_name"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	err = services.DeleteShockSetting(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func GetGMMShockBases(c *gin.Context) {
	results := services.GetGMMShockBases()
	c.JSON(http.StatusOK, results)
}

func GetJobsTemplate(c *gin.Context) {
	templateId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	templates := services.GetJobsTemplate(templateId)
	c.JSON(http.StatusOK, templates)
}

func DeleteJobsTemplate(c *gin.Context) {
	templateId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	templates := services.DeleteJobsTemplate(templateId)
	c.JSON(http.StatusOK, templates)
}

func GetAllJobsTemplates(c *gin.Context) {
	templates := services.GetJobsTemplates()
	c.JSON(http.StatusOK, templates)
}

func SaveJobsTemplate(c *gin.Context) {
	var template models.JobsTemplate
	c.BindJSON(&template)
	templates := services.SaveJobsTemplate(template)
	c.JSON(http.StatusOK, templates)
}

func SaveVariableGroups(c *gin.Context) {
	var group models.AggregatedVariableGroup
	err := c.BindJSON(&group)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group, err = services.SaveVariableGroup(group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, group)
}

func GetVariableGroups(c *gin.Context) {
	results := services.GetVariableGroups()
	c.JSON(http.StatusOK, results)
}

func UpdateVariableGroup(c *gin.Context) {
	var group models.AggregatedVariableGroup
	err := c.BindJSON(&group)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	group.ID = id

	group, err = services.SaveVariableGroup(group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, group)
}

func DeleteVariableGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = services.DeleteVariableGroup(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// GetModelPointUploadStatus returns the current status of an async modelpoint upload/import
func GetModelPointUploadStatus(c *gin.Context) {
	uploadID := c.Param("upload_id")
	if uploadID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing upload_id"})
		return
	}
	if st, ok := services.GetUploadStatus(uploadID); ok {
		c.JSON(http.StatusOK, st)
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "upload not found"})
}
