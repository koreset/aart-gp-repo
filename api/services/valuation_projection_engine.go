package services

import (
	"api/models"
	"api/utils"
	"context"
	"fmt"
	"math"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"api/log"

	"github.com/dgraph-io/ristretto"
	"github.com/gammazero/workerpool"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var Cache *ristretto.Cache
var GroupPricingCache *ristretto.Cache

var mutex = &sync.Mutex{}
var aggMutex = &sync.Mutex{}

type TestResult struct {
	Value int64
}

type ProductTableNames struct {
	MortalityTableName            string
	MortalityColumnName           string
	MortalityAccidentalTableName  string
	MortalityAccidentalColumnName string
	LapseTableName                string
	LapseColumnName               string
	DisabilityTableName           string
	DisabilityColumnName          string
	RetrenchmentTableName         string
	RetrenchmentColumnName        string
	LapseMonthCount               int
	LapseMarginMonthCount         int
	RetrenchmentRowCount          int
}

var ptnames = make(map[string]ProductTableNames)

var DisabilityColumnName = ""

//var RetrenchmentTableName string
//var LapseMonthCount int

var RecordsSkipped = 0

func init() {
	Cache, _ = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of Cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	PaaCache, _ = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of Cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	IbnrCache, _ = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of Cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	PricingCache, _ = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of Cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	ExpAnalysisCache, _ = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of Cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	GroupPricingCache, _ = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of Cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
}

// PopulateProjectionsForProduct kickstarts the process of calculating cash flows
func PopulateProjectionsForProduct(prodId int, job *models.ProjectionJob) error {
	// log the start of the run status

	// Get the available DB stats
	//sqlDb, _ := DB.DB()

	//stats := sqlDb.Stats()
	// 2. Get the underlying *sql.DB object
	sqlDB, _ := DB.DB()

	// 3. Get and print the connection pool stats
	stats := sqlDB.Stats()
	log.Info("MaxOpenConnections: ", stats.MaxOpenConnections)
	log.Info("OpenConnections: ", stats.OpenConnections)
	log.Info("InUse: ", stats.InUse)
	log.Info("Idle: ", stats.Idle)
	log.Info("WaitCount: ", stats.WaitCount)
	log.Info("WaitDuration: ", stats.WaitDuration)
	log.Info("MaxIdleClosed: ", stats.MaxIdleClosed)
	log.Info("MaxLifetimeClosed: ", stats.MaxLifetimeClosed)

	// keep checking for stats.OpenConnections until we get a value greater than 0

	if stats.WaitCount > 2000 {
		log.Info("Waiting for connections to be released")
		time.Sleep(time.Second * 120)
		//return errors.New("Waiting for connections to be released")
	}

	log.Info("Resuming Execution")

	Cache.Clear()
	var aggregatedProjections []models.Projection
	var scopedAggregatedProjections []models.Projection
	aggProjs := make(map[string]models.Projection)
	var sap []models.ScopedAggregatedProjection
	sapMap := make(map[string]models.Projection)
	var lic []models.LICAggregatedProjections
	var run models.RunParameters
	// Redis: Try get RunParameters by job id (includes ShockSettings if previously stored)
	if RedisAvailable() {
		if RedisGetJSON(fmt.Sprintf("ads:v1:run:%d", job.ID), &run) {
			// ensure ProjectionJobID is set
			run.ProjectionJobID = job.ID
		} else {
			// cache miss - load from DB
			DB.Where("projection_job_id = ?", job.ID).Find(&run)
			var ss models.ShockSetting
			if RedisAvailable() {
				// also try cache shock settings
				if !RedisGetJSON(fmt.Sprintf("ads:v1:shockSettings:%d", run.ShockSettingsID), &ss) {
					DB.Where("id = ?", run.ShockSettingsID).Find(&ss)
					RedisSetJSON(fmt.Sprintf("ads:v1:shockSettings:%d", run.ShockSettingsID), &ss, 120*time.Minute)
				}
			} else {
				DB.Where("id = ?", run.ShockSettingsID).Find(&ss)
			}
			run.ShockSettings = ss
			run.ProjectionJobID = job.ID
			// set run with embedded shock settings
			RedisSetJSON(fmt.Sprintf("ads:v1:run:%d", job.ID), &run, 120*time.Minute)
		}
	} else {
		// Redis unavailable
		DB.Where("projection_job_id = ?", job.ID).Find(&run)
		var shockSettings models.ShockSetting
		DB.Where("id = ?", run.ShockSettingsID).Find(&shockSettings)
		run.ShockSettings = shockSettings
		run.ProjectionJobID = job.ID
	}
	var prod models.Product
	var jobProduct = models.JobProduct{}
	var mps []models.ProductModelPoint
	var err error
	var productShocks []models.ProductShock

	var prodTableNames = ProductTableNames{}

	// log the start of the run
	log.WithField("prodId", prodId).WithField("jobID", job.ID).Info(fmt.Sprintf("Starting projections for product ID: %d, Job ID: %d", prodId, job.ID))

	//Get the product
	startTime := time.Now()
	testStart := time.Now()
	projectionRange := 0

	//Load product
	prod, err = GetProductById(prodId)
	if err != nil {
		//log.Error().Err(err)
		log.WithFields(map[string]interface{}{
			"prodId": prodId,
			"jobID":  job.ID,
		}).Error(errors.Wrap(err, "failed to get product"))
		return errors.Wrap(err, "failed to get product")
	}

	log.WithField("prodId", prodId).Info("Retrieved product from database")

	// moved the population of transition states for a products here...
	var states []models.ProductTransitionState
	// Redis first, then local ristretto cache
	redisStateKey := fmt.Sprintf("ads:v1:product:%d:states", prod.ID)
	if RedisAvailable() && RedisGetJSON(redisStateKey, &states) {
		// ok from Redis
	} else {
		key := fmt.Sprintf("states_%d", prod.ID)
		if cached, found := Cache.Get(key); found {
			states = cached.([]models.ProductTransitionState)
		} else {
			err := DB.Where("product_id = ?", prod.ID).Find(&states).Error
			if err == nil {
				Cache.Set(key, states, 1)
				if RedisAvailable() {
					RedisSetJSON(redisStateKey, &states, 120*time.Minute)
				}
			}
		}
	}

	log.WithField("prodId", prodId).Info("Retrieved markov states from database")

	// check if there is an existing job product for the product code and projection job id
	//err = DB.Where("product_code = ? and projection_job_id = ?", prod.ProductCode, job.ID).First(&jobProduct).Error
	jobErr := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("product_code = ? and projection_job_id = ?", prod.ProductCode, job.ID).First(&jobProduct).Error
	})

	if jobErr == nil && jobProduct.ID > 0 {
		// there is an existing job product, so we will not proceed with the run
		log.WithFields(map[string]interface{}{
			"prodId":          prodId,
			"jobID":           job.ID,
			"productCode":     prod.ProductCode,
			"projectionJobId": job.ID,
		}).Error(errors.New(fmt.Sprintf("Job product already exists for product code: %s and projection job id: %d", prod.ProductCode, job.ID)))
		return errors.New(fmt.Sprintf("Job product already exists for product code: %s and projection job id: %d", prod.ProductCode, job.ID))
	}
	log.WithField("prodId", prodId).Info("Successful check for non existent job product")

	//Get all Modelpoints for the product

	tableName := strings.ToLower(prod.ProductCode) + "_modelpoints"

	prodTableNames.MortalityTableName = GetRatingTable(prod.ProductCode, "Death")
	if prodTableNames.MortalityTableName != "" {
		prodTableNames.MortalityColumnName = GetColumnName(strings.ToLower(prod.ProductCode) + "_" + strings.ToLower(prodTableNames.MortalityTableName))
		log.WithField("prodId", prodId).Info("Mortality Table Name: ", prodTableNames.MortalityTableName)
	}
	prodTableNames.MortalityAccidentalTableName = GetRatingTable(prod.ProductCode, "Accidental Death")
	if prodTableNames.MortalityAccidentalTableName != "" {
		log.WithField("prodId", prodId).Info("Mortality Accidental Table Name: ", prodTableNames.MortalityAccidentalTableName)
		prodTableNames.MortalityAccidentalColumnName = GetColumnName(strings.ToLower(prod.ProductCode) + "_" + strings.ToLower(prodTableNames.MortalityAccidentalTableName))
	}

	prodTableNames.DisabilityTableName = GetRatingTable(prod.ProductCode, "Permanent Disability")
	if prodTableNames.DisabilityTableName != "" {
		log.WithField("prodId", prodId).Info("Disability Table Name: ", prodTableNames.DisabilityTableName)
		prodTableNames.DisabilityColumnName = GetColumnName(strings.ToLower(prod.ProductCode) + "_" + strings.ToLower(prodTableNames.DisabilityTableName))
	}

	prodTableNames.RetrenchmentTableName = GetRatingTable(prod.ProductCode, "Retrenchment")
	if prodTableNames.RetrenchmentTableName != "" {
		prodTableNames.RetrenchmentColumnName = GetColumnName(strings.ToLower(prod.ProductCode) + "_" + strings.ToLower(prodTableNames.RetrenchmentTableName))
		log.WithField("prodId", prodId).Info("Retrenchment Table Name: ", prodTableNames.RetrenchmentTableName)
	}

	prodTableNames.LapseTableName = GetRatingTable(prod.ProductCode, "Lapse")
	if prodTableNames.LapseTableName != "" {
		prodTableNames.LapseColumnName = GetColumnName(strings.ToLower(prod.ProductCode) + "_" + strings.ToLower(prodTableNames.LapseTableName))
		log.WithField("prodId", prodId).Info("Lapse Table Name: ", prodTableNames.LapseTableName)

	}

	ptnames[prod.ProductCode] = prodTableNames

	// if states contains retrenchment, then we need to get the retrenchment count
	if utils.StatesContains(&states, "Retrenchment") {
		log.WithField("prodId", prodId).Info("Retrenchment Active")
		prodTableNames.RetrenchmentRowCount = GetRetrenchmentCount(run.RetrenchmentYear, prodTableNames.RetrenchmentTableName, prodTableNames.RetrenchmentColumnName)
	}

	jobProduct.ProductName = prod.ProductName
	jobProduct.ProductID = prod.ID
	jobProduct.ProjectionJobID = job.ID
	jobProduct.ProductCode = prod.ProductCode

	LoadMortalityRates(prod.ProductCode, run.MortalityYear)
	log.WithField("prodId", prodId).Info("Mortality Rates Loaded")
	LoadAccidentalMortalityRates(prod.ProductCode, run.MortalityYear)
	log.WithField("prodId", prodId).Info("Accidental Mortality Rates Loaded")
	LoadDisabilityRates(prod.ProductCode, run.MorbidityYear)
	log.WithField("prodId", prodId).Info("Disability Rates Loaded")
	LoadRetrenchmentRates(prod.ProductCode, run.RetrenchmentYear)
	log.WithField("prodId", prodId).Info("Retrenchment Rates Loaded")
	LoadLapseRates(prod.ProductCode, run.LapseYear)
	log.WithField("prodId", prodId).Info("Lapse Rates Loaded")

	var productMargins models.ProductMargins
	var parameters models.ProductParameters
	var features models.ProductFeatures
	var multipliers models.ProductAccidentalBenefitMultiplier
	var benefitMultipliers models.ProductBenefitMultiplier

	// Product parameters (cache by prodId:year:basis)
	paramKey := fmt.Sprintf("ads:v1:product:%d:parameters:%d:%s", prod.ID, run.ParameterYear, run.RunBasis)
	if !(RedisAvailable() && RedisGetJSON(paramKey, &parameters)) {
		err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Where("year = ? and product_code=? and basis=?", run.ParameterYear, prod.ProductCode, run.RunBasis).Find(&parameters).Error
		})
		if err == nil && RedisAvailable() {
			RedisSetJSON(paramKey, &parameters, 120*time.Minute)
		}
	}
	if err != nil {
		log.WithField("prodId", prodId).WithField("jobID", job.ID).Error(fmt.Sprintf("Product Parameters not found for product code: %s and year: %d and basis: %s", prod.ProductCode, run.ParameterYear, run.RunBasis))
		jobProduct.JobStatus = "failed"
		jobProduct.JobStatusError = fmt.Sprintf("Product Parameters not found for product code: %s and year: %d and basis: %s", prod.ProductCode, run.ParameterYear, run.RunBasis)
		DB.Save(&jobProduct)

	} else {
		log.WithField("prodId", prodId).Info("Parameters Loaded")
	}

	prodTableNames.LapseMarginMonthCount = GetLapseMarginCount(run.LapseMarginYear, prod.ProductCode, parameters.MarginBasis)
	log.WithField("prodId", prodId).Info("Lapse margin month Count: ", prodTableNames.LapseMarginMonthCount)

	prodTableNames.LapseMonthCount = GetLapseCount(run.LapseYear, prod.ProductCode)
	log.WithField("prodId", prodId).Info("Lapse month Count: ", prodTableNames.LapseMonthCount)

	ptnames[prod.ProductCode] = prodTableNames

	// Product features (cache by prodId)
	featuresKey := fmt.Sprintf("ads:v1:product:%d:features", prod.ID)
	if !(RedisAvailable() && RedisGetJSON(featuresKey, &features)) {
		err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Where("product_code=?", prod.ProductCode).Find(&features).Error
		})
		if err == nil && RedisAvailable() {
			RedisSetJSON(featuresKey, &features, 120*time.Minute)
		}
	}
	if err != nil {
		log.WithField("prodId", prodId).WithField("jobID", job.ID).Error(fmt.Sprintf("Product Features not found for product code: %s", prod.ProductCode))
	} else {
		log.WithField("prodId", prodId).Info("Features Loaded")
	}

	// Accidental benefit multipliers (cache by prodId)
	abmKey := fmt.Sprintf("ads:v1:product:%d:multipliers:accidental", prod.ID)
	if !(RedisAvailable() && RedisGetJSON(abmKey, &multipliers)) {
		err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Where("product_code = ? ", prod.ProductCode).Find(&multipliers).Error
		})
		if err == nil && RedisAvailable() {
			RedisSetJSON(abmKey, &multipliers, 120*time.Minute)
		}
	}
	if err != nil {
		log.WithField("prodId", prodId).WithField("jobID", job.ID).Error(fmt.Sprintf("Product Accidental Benefit Multipliers not found for product code: %s", prod.ProductCode))

	} else {
		log.WithField("prodId", prodId).Info("Accidental Multipliers Loaded")
	}

	// Benefit multipliers (cache by prodId)
	bmKey := fmt.Sprintf("ads:v1:product:%d:multipliers:benefit", prod.ID)
	if !(RedisAvailable() && RedisGetJSON(bmKey, &benefitMultipliers)) {
		err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Where("product_code = ? ", prod.ProductCode).Find(&benefitMultipliers).Error
		})
		if err == nil && RedisAvailable() {
			RedisSetJSON(bmKey, &benefitMultipliers, 120*time.Minute)
		}
	}
	if err != nil {
		log.WithField("prodId", prodId).WithField("jobID", job.ID).Error(fmt.Sprintf("Product Benefit Multipliers not found for product code: %s", prod.ProductCode))
		jobProduct.JobStatus = "failed"
		jobProduct.JobStatusError = fmt.Sprintf("Product Benefit Multipliers not found for product code: %s", prod.ProductCode)
		DB.Save(&jobProduct)
	} else {
		log.WithField("prodId", prodId).Info("Benefit Multipliers Loaded")
	}

	// Product margins (cache by prodId:year:basis)
	margKey := fmt.Sprintf("ads:v1:productcode:%s:margins:%d:%s", prod.ProductCode, run.ParameterYear, parameters.MarginBasis)
	if !(RedisAvailable() && RedisGetJSON(margKey, &productMargins)) {
		err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Where("product_code=? and year=? and basis=?", prod.ProductCode, run.ParameterYear, parameters.MarginBasis).First(&productMargins).Error
		})
		if err == nil && RedisAvailable() {
			RedisSetJSON(margKey, &productMargins, 120*time.Minute)
		}
	}
	if err != nil {
		log.WithFields(map[string]interface{}{
			"prodId":   prodId,
			"jobID":    job.ID,
			"prodCode": prod.ProductCode,
			"runYear":  run.ParameterYear,
			"basis":    parameters.MarginBasis,
		}).Error("Product Margins not found for product code")
		jobProduct.JobStatus = "failed"
		jobProduct.JobStatusError = fmt.Sprintf("Product Margins not found for product code: %s and year: %d and margin basis: %s", prod.ProductCode, run.ParameterYear, parameters.MarginBasis)
		DB.Save(&jobProduct)

		return errors.New(fmt.Sprintf("Product Margins not found for product code: %s and year: %d and margin basis: %s", prod.ProductCode, run.ParameterYear, parameters.MarginBasis))
	} else {
		log.WithField("prodId", prodId).Info("Margins Loaded")
	}

	// Load modelpoints
	if run.RunSingle {
		// TODO: revert back to 1 when done testing
		err = DB.Table(tableName).Where("year =? and mp_version=?", run.ModelpointYear, run.ModelPointVersion).Limit(MpLimitValue).Find(&mps).Error
	} else {
		err = DB.Table(tableName).Where("year =? and mp_version=?", run.ModelpointYear, run.ModelPointVersion).Find(&mps).Error
	}

	if len(mps) == 0 {
		// return message to include the product code and year and version
		//log.Error().Msgf("No model points found for product code: %s, year: %d, version: %d", prod.ProductCode, run.ModelpointYear, run.ModelPointVersion)
		log.WithField("prodId", prodId).WithField("jobID", job.ID).Error(fmt.Sprintf("No model points found for product code: %s, year: %d, version: %s", prod.ProductCode, run.ModelpointYear, run.ModelPointVersion))
		jobProduct.JobStatus = "failed"
		jobProduct.JobStatusError = fmt.Sprintf("No model points found for product code: %s, year: %d, version: %s", prod.ProductCode, run.ModelpointYear, run.ModelPointVersion)
		DB.Save(&jobProduct)

		return errors.New(fmt.Sprintf("No model points found for product code: %s, year: %d, version: %s", prod.ProductCode, run.ModelpointYear, run.ModelPointVersion))
	}
	//Load product shocks (cache by prodId and shock basis)
	shocksKey := fmt.Sprintf("ads:v1:product:%d:shocks:%s", prod.ID, run.ShockSettings.ShockBasis)
	if !(RedisAvailable() && RedisGetJSON(shocksKey, &productShocks)) {
		err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Where("shock_basis = ?", run.ShockSettings.ShockBasis).Find(&productShocks).Error
		})
		if err == nil && RedisAvailable() {
			RedisSetJSON(shocksKey, &productShocks, 120*time.Minute)
		}
	}
	if err != nil {
		jobProduct.JobStatus = "failed"
		jobProduct.JobStatusError = fmt.Sprintf("missing shock basis: %s in the shocks table", run.ShockSettings.ShockBasis)
		DB.Save(&jobProduct)
		return errors.Wrap(err, fmt.Sprintf("missing shock basis: %s in the shocks table", run.ShockSettings.ShockBasis))
	} else {
		log.WithField("prodId", prodId).Info("Shocks Loaded")
	}

	//Do yield curve check here

	jobProduct.ModelPointCount = len(mps)
	jobProduct.ProductName = prod.ProductName
	jobProduct.ProductID = prod.ID
	jobProduct.ProjectionJobID = job.ID
	jobProduct.ProductCode = prod.ProductCode

	yieldCurveCode := parameters.YieldCurveCode
	basis := parameters.Basis
	var yieldCurveYear int
	var month int

	if run.YieldCurveBasis == Current {
		yieldCurveYear = run.YieldcurveYear
		month = run.YieldcurveMonth // valMonth
	} else {
		yieldCurveYear = mps[0].LockedInYear
		month = mps[0].LockedInMonth
	}

	projMonth01 := 1
	var fwRate float64
	var yieldCurve models.YieldCurve
	if run.YieldCurveBasis == LockedInRates {
		//err = DB.Table("yield_curve").Where("year = ? and proj_time = ? and month = ? and basis=?", yieldCurveYear, projMonth01, month, basis).First(&yieldCurve).Error
		err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Table("yield_curve").Where("year = ? and proj_time = ? and month = ?", yieldCurveYear, projMonth01, month).First(&yieldCurve).Error
		})
		if err != nil {
			log.WithField("prodId", prodId).WithField("jobID", job.ID).Error(fmt.Sprintf("yield curve data for yield curve code: %s, yield curve year: %d, month: %d was not found: %v", yieldCurveCode, yieldCurveYear, month, err))
		} else {
			log.WithField("prodId", prodId).Info("Yield Curve Loaded")
		}
	}

	if run.YieldCurveBasis == Current {
		//err = DB.Table("yield_curve").Where("year = ? and proj_time = ? and month = ? and yield_curve_code=? and basis=?", yieldCurveYear, projMonth01, month, yieldCurveCode, basis).First(&yieldCurve).Error
		err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Table("yield_curve").Where("year = ? and proj_time = ? and month = ? and yield_curve_code=?", yieldCurveYear, projMonth01, month, yieldCurveCode).First(&yieldCurve).Error
		})
		if err != nil {
			log.WithField("prodId", prodId).WithField("jobID", job.ID).Error(fmt.Sprintf("yield curve data for yield curve code: %s, yield curve year: %d, month: %d was not found: %v", yieldCurveCode, yieldCurveYear, month, err))
		} else {
			log.WithField("prodId", prodId).Info("Yield Curve Loaded")
		}
	}

	fwRate = yieldCurve.NominalRate

	if err != nil {
		jobProduct.JobStatus = "failed"
		jobProduct.JobStatusError = fmt.Sprintf("yield curve data for yield curve code: %s, yield curve year: %d, month: %d and basis: %s was not found: %v", yieldCurveCode, yieldCurveYear, month, basis, err)
		DB.Save(&jobProduct)
		log.WithField("prodId", prodId).WithField("jobID", job.ID).Error(
			fmt.Sprintf("yield curve data for yield curve code: %s, yield curve year: %d, month: %d was not found: %v", yieldCurveCode, yieldCurveYear, month, err),
		)
		return errors.Wrap(err, fmt.Sprintf("yield curve data for yield curve code: %s, yield curve year: %d, month: %d  was not found: %v", yieldCurveCode, yieldCurveYear, month, err))
	}
	if fwRate == 0 {
		jobProduct.JobStatus = "failed"
		jobProduct.JobStatusError = fmt.Sprintf("yield curve data for yield curve code: %s, yield curve year: %d, month: %d  returned zero rate", yieldCurveCode, yieldCurveYear, month)
		DB.Save(&jobProduct)
		log.WithField("prodId", prodId).WithField("jobID", job.ID).Error()
		return errors.New(fmt.Sprintf("yield curve data for yield curve code: %s, yield curve year: %d, month: %d returned zero rate", yieldCurveCode, yieldCurveYear, month))
	}

	//The run parameter year must be available in the SpecialDecrementMargins table for the product code if the special decrement feature is turned on
	//If not, then the run cannot proceed

	if features.SpecialDecrementMargin {
		var specialMargins models.ProductSpecialDecrementMargin

		//err = DB.Where("product_code = ? and special_margin_code = ?", prod.ProductCode, parameters.SpecialMarginCode).First(&specialMargins).Error
		err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
			return d.Where("product_code = ? and special_margin_code = ?", prod.ProductCode, parameters.SpecialMarginCode).First(&specialMargins).Error
		})
		if err != nil && specialMargins.ID == 0 {
			log.WithField("prodId", prodId).WithField("jobID", job.ID).Error(fmt.Sprintf("Special decrement margin data for product code: %s and special_margin_code: %s was not found", prod.ProductCode, parameters.SpecialMarginCode))
			jobProduct.JobStatus = "failed"
			jobProduct.JobStatusError = fmt.Sprintf("special decrement margin data for product code: %s and special_margin_code: %s was not found", prod.ProductCode, parameters.SpecialMarginCode)
			DB.Save(&jobProduct)
			return errors.New(fmt.Sprintf("special decrement margin data for product code: %s and special_margin_code: %s was not found", prod.ProductCode, parameters.SpecialMarginCode))
		}
	}

	// the run parameter year must be available in Margins table for the product code
	// if not, then the run cannot proceed
	var margin models.ProductMargins
	//err = DB.Where("product_code = ? and year = ? and basis=?", prod.ProductCode, run.ParameterYear, parameters.MarginBasis).First(&margin).Error
	err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("product_code = ? and year = ? and basis=?", prod.ProductCode, run.ParameterYear, parameters.MarginBasis).First(&margin).Error
	})

	if err != nil && margin.Year == 0 {
		//log.Error().Msgf("Product Code: %s", errors.WithStack(err))
		log.WithField("prodId", prodId).WithField("jobID", job.ID).Error(fmt.Sprintf("Margin data for product code: %s and year: %d and basis:%s was not found", prod.ProductCode, run.ParameterYear, parameters.MarginBasis))
		jobProduct.JobStatus = "failed"
		jobProduct.JobStatusError = fmt.Sprintf("margin data for product code: %s and year: %d and basis:%s was not found", prod.ProductCode, run.ParameterYear, parameters.MarginBasis)
		DB.Save(&jobProduct)
		return errors.New(fmt.Sprintf("margin data for product code: %s and year: %d and basis:%s was not found", prod.ProductCode, run.ParameterYear, parameters.MarginBasis))
	} else {
		log.WithField("prodId", prodId).Info("Margins Loaded")
	}

	//The run parameter year must be available in RenewableProfitAdjustments table for the product code if the renewable profit feature is turned on
	//If not, then the run cannot proceed
	if features.RenewableProfitAdjustment {
		var renewableProfit models.ProductRenewableProfitAdjustment
		err = DB.Where("product_code = ? and profit_adjustment_code = ?", prod.ProductCode, parameters.ProfitAdjustmentCode).First(&renewableProfit).Error
		//err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		//	return d.Where("product_code = ? and profit_adjustment_code = ?", prod.ProductCode, parameters.ProfitAdjustmentCode).First(&renewableProfit).Error
		//})

		if err != nil && renewableProfit.ID == 0 {
			//log.Error().Msgf("Product Code: %s", errors.WithStack(err))
			log.WithField("prodId", prodId).WithField("jobID", job.ID).Error(fmt.Sprintf("Renewable profit adjustment data for product code: %s and profit_adjustment_code: %s was not found", prod.ProductCode, parameters.ProfitAdjustmentCode))
			jobProduct.JobStatus = "failed"
			jobProduct.JobStatusError = fmt.Sprintf("renewable profit adjustment data for product code: %s and profit_adjustment_code: %s was not found", prod.ProductCode, parameters.ProfitAdjustmentCode)
			DB.Save(&jobProduct)
			return errors.New(fmt.Sprintf("renewable profit adjustment data for product code: %s and profit_adjustment_code: %s was not found", prod.ProductCode, parameters.ProfitAdjustmentCode))
		} else {
			log.WithField("prodId", prodId).Info("Renewable Profit Adjustments Loaded")
		}
	} else {
		log.WithField("prodId", prodId).Info("Renewable Profit Adjustments Not Selected")
	}

	if features.ProportionalReinsurance {
		//Get a distinct list of treaty year values from the model points table
		var treatyYears []int
		mpTableName := strings.ToLower(prod.ProductCode) + "_modelpoints"
		DB.Table(mpTableName).Where("year = ? and mp_version=?", run.ModelpointYear, run.ModelPointVersion).Pluck("distinct treaty_year", &treatyYears)

		for _, treatyYear := range treatyYears {
			var reinsurance models.ProductReinsurance
			err = DB.Where("product_code = ? and treaty_year = ?", prod.ProductCode, treatyYear).First(&reinsurance).Error
			if err != nil && reinsurance.ID == 0 {
				//log.Error().Msgf("Product Code: %s", errors.WithStack(err))
				log.WithField("prodId", prodId).WithField("jobID", job.ID).Error(fmt.Sprintf("Reinsurance data for product code: %s and treaty year: %d was not found", prod.ProductCode, treatyYear))
				jobProduct.JobStatus = "failed"
				jobProduct.JobStatusError = fmt.Sprintf("reinsurance data for product code: %s and year: %d was not found", prod.ProductCode, treatyYear)
				DB.Save(&jobProduct)
				return errors.New(fmt.Sprintf("reinsurance data for product code: %s and treaty year: %d was not found", prod.ProductCode, treatyYear))
			}
		}
		log.WithField("prodId", prodId).Info("Reinsurance Loaded")
	} else {
		log.WithField("prodId", prodId).Info("Reinsurance Not Selected")
	}

	// checking that the product accidental mortality year is available in the AccidentalMortalityRates table for the product code
	// if not, then the run cannot proceed.
	// get the product accidental mortality table name

	//prodAccMortTableName := GetRatingTable(prod.ProductCode, "Accidental Death")
	//log.WithField("prodId", prodId).WithField("jobID", job.ID).Info("Product Accidental Mortality Table Name: ", prodAccMortTableName)
	//if prodAccMortTableName != "" {
	//	columnName := GetColumnName(strings.ToLower(prod.ProductCode) + "_" + strings.ToLower(prodAccMortTableName))
	//	//get a count where the accidental mortality year is equal to the run accidental mortality year
	//	var count int64
	//
	//	tName := strings.ToLower(prod.ProductCode) + "_" + strings.ToLower(prodAccMortTableName)
	//	//DB.Table(tName).Where("year = ?", run.MortalityYear).Count(&count)
	//
	//	query := fmt.Sprintf("select count(*) from %s where %s like '%%%d%%'", tName, columnName, run.MortalityYear)
	//	DB.Raw(query).Row()
	//
	//	if scanErr := DB.Raw(query).Scan(&count).Error; scanErr != nil {
	//		log.WithField("prodId", prodId).WithField("jobID", job.ID).Error(fmt.Sprintf("Error: %s", scanErr))
	//	}
	//
	//	if count == 0 {
	//		jobProduct.JobStatus = "failed"
	//		jobProduct.JobStatusError = fmt.Sprintf("accidental mortality data for product code: %s and year: %d was not found", prod.ProductCode, run.MortalityYear)
	//		DB.Save(&jobProduct)
	//		return errors.New(fmt.Sprintf("accidental mortality data for product code: %s and year: %d was not found", prod.ProductCode, run.MortalityYear))
	//	} else {
	//		log.WithField("prodId", prodId).Info("Accidental Mortality Loaded")
	//	}
	//}

	//checking that the product mortality year is available in the MortalityRates table for the product code
	//if not, then the run cannot proceed

	log.WithField("prodId", prodId).WithField("jobID", job.ID).Info("Updating Job Product")
	err = DB.Save(&jobProduct).Error

	if err != nil {
		//log.Error().Err(err).Send()
		log.WithField("prodId", prodId).WithField("jobID", job.ID).Error(fmt.Sprintf("Job Product was not saved: %s ", err.Error()))
	}
	log.WithField("prodId", prodId).WithField("jobID", job.ID).Info("Updated Job Product")

	testEnd := time.Since(testStart)
	log.WithField("prodId", prodId).WithField("jobID", job.ID).Info("Time taken to initialize Run: ", testEnd)

	// Setup cancellation context and poller to detect job cancellation without per-worker DB queries
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Poll job status periodically and cancel context if job is marked as Cancelled
	go func(jobID int) {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				var cur models.ProjectionJob
				if err := DB.Select("id, status").Where("id = ?", jobID).First(&cur).Error; err == nil {
					if strings.EqualFold(cur.Status, "Cancelled") {
						cancel()
						return
					}
				}
			}
		}
	}(job.ID)

	wp := workerpool.New(min(runtime.NumCPU(), 8))

	// Create a slice to collect errors from model points
	var modelPointErrors []*ModelPointError
	var errorMutex sync.Mutex

	for mpIndex, mp := range mps { //Remember to adjust here later.
		mp := mp
		mpIndex := mpIndex

		wp.Submit(func() {
			mpErr := PopulateProjectionPerModelPoint(mpIndex, mp, projectionRange,
				&jobProduct, run,
				prod, &sap, productMargins,
				features, parameters, multipliers,
				&lic, &aggProjs, &sapMap,
				productShocks, states, job.ID, ctx)

			// If there was an error, add it to the collection
			if mpErr != nil {
				errorMutex.Lock()
				modelPointErrors = append(modelPointErrors, mpErr)
				errorMutex.Unlock()
			}
		})
	}
	wp.StopWait()

	// Process any errors that occurred during model point processing
	var modelPointErrorsFound error

	if len(modelPointErrors) > 0 {
		//log.Info().Msgf("Encountered %d errors during model point processing", len(modelPointErrors))
		log.WithField("prodId", prodId).WithField("jobID", job.ID).Error(fmt.Sprintf("Encountered %d errors during model point processing", len(modelPointErrors)))

		// Log each error with its associated model point information
		for _, mpErr := range modelPointErrors {
			log.WithField("prodId", prodId).WithField("jobID", job.ID).Error(fmt.Sprintf("Model point error: %s", mpErr.Error))
		}

		// Store the errors in the job product for later reference
		errorMessages := make([]string, 0, 5)
		for i, mpErr := range modelPointErrors {
			if i < 5 {
				errorMessages = append(errorMessages, fmt.Sprintf("Model Point (SpCode: %d, ProductCode: %s): %s",
					mpErr.ModelPoint.Spcode, mpErr.ModelPoint.ProductCode, mpErr.Error.Error()))
			}
		}

		// Join all error messages with a separator
		errorSummary := strings.Join(errorMessages, " | ")

		// Create an error to return to the caller
		modelPointErrorsFound = fmt.Errorf("encountered %d model point errors: %s", len(modelPointErrors), errorSummary)

		// Append to JobStatusError if it already has content, otherwise set it
		if jobProduct.JobStatusError != "" {
			jobProduct.JobStatusError += " | " + errorSummary
		} else {
			jobProduct.JobStatusError = "Model point errors: " + errorSummary
		}

		err = DB.Save(&jobProduct).Error
		if err != nil {
			log.WithField("prodId", prodId).WithField("jobID", job.ID).Error(fmt.Sprintf("Failed to save job product with model point errors"))
		}

		// We don't fail the entire job because of individual model point errors
		// The errors are collected and can be reviewed later

	}

	if jobProduct.JobStatus != "cancelled" {

		for _, appProj := range aggProjs {
			appProj.ID = 0
			aggregatedProjections = append(aggregatedProjections, appProj)
		}

		// if the aggregated projections is empty, return an error
		if len(aggregatedProjections) == 0 {
			jobProduct.JobStatus = "failed"
			jobProduct.JobStatusError = "no projections were generated"
			DB.Save(&jobProduct)
			return errors.New("no projections were generated")
		}

		err = DB.Table("aggregated_projections").CreateInBatches(&aggregatedProjections, 100).Error
		if err != nil {
			log.WithField("prodId", prodId).WithField("jobID", job.ID).Error(fmt.Sprintf("Failed to save aggregated projections"))
		}

		if run.IFRS17Indicator {
			for _, appProj := range sapMap {
				appProj.ID = 0
				scopedAggregatedProjections = append(scopedAggregatedProjections, appProj)
			}
			err = DB.Table("scoped_aggregated_projections").CreateInBatches(&scopedAggregatedProjections, 100).Error
			if err != nil {
				log.WithField("prodId", prodId).WithField("jobID", job.ID).Error("Failed to save scoped aggregated projections")
			}
		}

		if run.LICIndicator {
			err = DB.CreateInBatches(&lic, 100).Error
			if err != nil {
				log.WithField("prodId", prodId).WithField("jobID", job.ID).Error("Failed to save LIC projections")
			}
		}

		duration := time.Since(startTime)

		log.WithField("prodId", prodId).WithField("jobID", job.ID).Info("Number of records processed: ", len(mps))
		log.WithField("prodId", prodId).WithField("jobID", job.ID).Info("Number of records processed: ", len(mps))
		log.WithField("prodId", prodId).WithField("jobID", job.ID).Info("Number of records processed: ", len(mps))
		DB.Save(&jobProduct)
		err = DB.Preload("Products").Where("id = ?", job.ID).Find(job).Error
		if err != nil {
			log.WithField("prodId", prodId).WithField("jobID", job.ID).Error("Failed to load job")
		}

		job.RunTime = job.RunTime + duration.Minutes()
		job.PointsDone = job.PointsDone + jobProduct.PointsDone

		totalPointsCompleted := 0

		for _, jp := range job.Products {
			totalPointsCompleted += jp.PointsDone
		}

		if job.TotalPoints == totalPointsCompleted {
			job.Status = "Complete"
		}
		RecordsSkipped = 0
		DB.Save(&job)
		aggregatedProjections = nil
		sap = nil
		sapMap = nil
		aggProjs = nil
	} else {
		log.WithField("prodId", prodId).WithField("jobID", job.ID).Info("Job failed or was cancelled")
		job.Status = "Cancelled"
		DB.Save(&job)
		aggregatedProjections = nil
		sap = nil
		sapMap = nil
		aggProjs = nil
	}
	log.WithField("prodId", prodId).WithField("jobID", job.ID).Info("Job completed")
	// Release the job-scoped mutex to avoid leaking locks for completed jobs
	ReleaseJobMutex(job.ID)

	// Return any model point errors that were found
	return modelPointErrorsFound
}

// PopulateProjectionPerModelPoint will run projections for all months specified in the time range.
// Good candidate for multithreading?? No. It relies on results from previous runs.
// ModelPointError represents an error associated with a specific model point
type ModelPointError struct {
	ModelPoint models.ProductModelPoint
	Error      error
}

func PopulateProjectionPerModelPoint(mpIndex int, modelPoint models.ProductModelPoint,
	projectionRange int,
	jobProduct *models.JobProduct,
	run models.RunParameters,
	prod models.Product,
	sap *[]models.ScopedAggregatedProjection,
	productMargins models.ProductMargins,
	features models.ProductFeatures,
	parameters models.ProductParameters,
	multipliers models.ProductAccidentalBenefitMultiplier,
	lic *[]models.LICAggregatedProjections,
	aggProjs *map[string]models.Projection,
	sapMap *map[string]models.Projection, productShocks []models.ProductShock, states []models.ProductTransitionState, jobID int, ctx context.Context) *ModelPointError {
	//startLoopTime := time.Now()
	var err error
	// Check for cancellation via context (avoid per-worker DB query)
	select {
	case <-ctx.Done():
		jobProduct.JobStatus = "cancelled"
		jobProduct.JobStatusError = "the job has been cancelled"
		log.WithField("prodId", prod.ID).WithField("jobID", jobProduct.ID).Error("Job cancelled")
		DB.Save(&jobProduct)
		return &ModelPointError{
			ModelPoint: modelPoint,
			Error:      errors.New("the job has been cancelled"),
		}
	default:
	}

	var projections []models.Projection

	prodConfigErr := CalculateProductConfig(&parameters, modelPoint, features, multipliers, run)
	if prodConfigErr != nil {
		return &ModelPointError{
			ModelPoint: modelPoint,
			Error:      errors.Wrap(prodConfigErr, "failed to calculate product config"),
		}
	}

	var index = 0
	var p models.Projection
	if projectionRange == 0 {
		projectionRange = 1440
	}
	valMonth, err := strconv.Atoi(run.RunDate[5:])
	if err != nil {
		valMonth = 0
	}

	valYear, err := strconv.Atoi(run.RunDate[:4])
	if err != nil {
		valYear = 0
	}

	for ; index <= projectionRange; index++ {
		// Periodically check for cancellation to avoid heavy DB polling per worker
		if index%50 == 0 {
			select {
			case <-ctx.Done():
				jobProduct.JobStatus = "cancelled"
				jobProduct.JobStatusError = "the job has been cancelled"
				log.WithField("prodId", prod.ID).WithField("jobID", jobProduct.ID).Error("Job cancelled during projection loop")
				DB.Save(&jobProduct)
				return &ModelPointError{
					ModelPoint: modelPoint,
					Error:      errors.New("the job has been cancelled"),
				}
			default:
			}
		}
		var projection models.Projection
		projection.RunType = run.RunType
		projection.RunDate = run.RunDate
		projection.RunName = run.RunName
		projection.RunId = jobProduct.ProjectionJobID
		projection.ProductCode = prod.ProductCode
		projection.ProductID = prod.ID
		projection.IFRS17Group = modelPoint.IFRS17Group
		projection.SpCode = modelPoint.Spcode
		projection.RunBasis = run.RunBasis

		CalculateTimeAndAge(&projection, modelPoint, index, jobProduct.ID, features)
		var shock models.ProductShock
		if projection.ProjectionMonth > 0 && run.ShockSettings.ShockBasis != "" {
			shock, err = GetShock(projection.ProjectionMonth, run.ShockSettings.ShockBasis, len(productShocks))
			if err != nil {
				log.WithField("prodId", prod.ID).WithField("jobID", jobProduct.ID).Error(fmt.Sprintf("GetShock error: %s", err.Error()))
				return &ModelPointError{
					ModelPoint: modelPoint,
					Error:      errors.Wrap(err, "failed to get shock"),
				}
			}
		}

		if projection.AgeNextBirthday > 121 {
			break
		}

		if projection.ProjectionMonth == 1 {
			projection.CalendarMonth = valMonth
		} else {
			if p.CalendarMonth+1 == 13 {
				projection.CalendarMonth = 1
			} else {
				projection.CalendarMonth = p.CalendarMonth + 1
			}
		}
		if projection.ProjectionMonth == 0 {
			projection.CalendarMonth = 0
		}

		// Only runs if Accidental Death is one of the transition states selected
		// and if accidental benefit is one of the features.
		if features.AccidentalDeathBenefit {
			AccidentProportion(&projection, modelPoint, run, states, parameters)
		}

		InflationFactor(index, valYear, run.YieldcurveMonth, &projection, &p, productMargins, run, modelPoint, parameters, features, shock)
		PremiumEscalation(index, &projection, &p, modelPoint, parameters) // sum assured and premium escalation rates

		if features.AnnuityIncome {
			AnnuityEscalation(index, &projection, &p, modelPoint, parameters, features)
		}

		//Assumes the lapse state as one of the default transition states
		LapseMargin(&projection, modelPoint, run, parameters)

		PremiumWaiverOnFactor(&projection, modelPoint, parameters)
		PaidUpOnFactor(&projection, modelPoint, parameters)

		// Only run if death is one of the transition states selected
		var specialMortalityMargin float64
		var mortalityTableProp float64
		var disabilityTableProp float64
		var lapseTableProp float64
		var mainMemberspecialMortalityMargin float64
		var mainMembermortalityTableProp float64

		if features.SpecialDecrementMargin {
			specialMortalityMargin = GetSpecialMargins(projection.AgeNextBirthday, &projection, modelPoint.ProductCode, modelPoint.MemberType, run, parameters, "mortality_margin")
			mortalityTableProp = GetSpecialMargins(projection.AgeNextBirthday, &projection, modelPoint.ProductCode, modelPoint.MemberType, run, parameters, "mortality_table_prop")
			disabilityTableProp = GetSpecialMargins(projection.AgeNextBirthday, &projection, modelPoint.ProductCode, modelPoint.MemberType, run, parameters, "disability_table_prop")
			lapseTableProp = GetSpecialMargins(projection.AgeNextBirthday, &projection, modelPoint.ProductCode, modelPoint.MemberType, run, parameters, "lapse_table_prop")

			//use correct main member age next birthday if a model point being processed is not a main member. This is important if lapse is dependent on the  main member's mortality rate.
			mainMemberspecialMortalityMargin = GetSpecialMargins(projection.MainMemberAgeNextBirthday, &projection, modelPoint.ProductCode, "MM", run, parameters, "mortality_margin")
			mainMembermortalityTableProp = GetSpecialMargins(projection.MainMemberAgeNextBirthday, &projection, modelPoint.ProductCode, "MM", run, parameters, "mortality_table_prop")

		} else {
			specialMortalityMargin = 0
			mortalityTableProp = 1
			disabilityTableProp = 1
			lapseTableProp = 1
			mainMemberspecialMortalityMargin = 0
			mainMembermortalityTableProp = 1
		}
		MainMemberMortalityRate(&projection, modelPoint, productMargins, run, states, shock, parameters, features, mainMemberspecialMortalityMargin, mainMembermortalityTableProp)

		//Only runs if lapse state is one of the transition states selected
		BaseLapses(&projection, p, modelPoint, run, states, shock, parameters, lapseTableProp)

		ContractingPartyAlivePortion(index, &projection, p, modelPoint, features, parameters)

		// Only runs if death is one of the transition states selected
		BaseMortalityRate(&projection, modelPoint, run, productMargins, states, shock, parameters, features, specialMortalityMargin, mortalityTableProp)

		//Only runs if nonlife feature is chosen
		if features.NonLife {
			NonLifeRiskRate(&projection, modelPoint, run, productMargins, states, shock, parameters)
		}

		BaseIndependentLapse(&projection, parameters)

		// Only runs if retrenchment is one of the transition states selected
		if utils.StatesContains(&states, Retrenchment) {
			if projection.AgeNextBirthday <= 65 {
				BaseRetrenchmentRate(&projection, modelPoint, productMargins, run, states, shock, parameters)
			}
		}

		//Only runs if disability is one of the transition states selected
		if utils.StatesContains(&states, PermanentDisability) {
			if projection.AgeNextBirthday <= 65 {
				BaseDisabilityIncidenceRates(&projection, modelPoint, productMargins, run, states, shock, parameters, disabilityTableProp)
			}
		}

		MainMemberMortalityRateByMonth(&projection, features, parameters)
		IndependentMortalityRateMonthly(&projection, parameters)
		IndependentLapseMonthly(&projection, parameters)
		IndependentRetrenchmentMonthly(&projection, parameters)
		IndependentDisabilityMonthly(&projection, parameters)

		MonthlyDependentMortality(&projection, parameters)
		MonthlyDependentLapse(&projection, parameters)
		MonthlyDependentRetrenchment(&projection, parameters)
		MonthlyDependentDisability(&projection, parameters)

		NumberOfMaturities(index, &projection, &p, modelPoint, features, states, parameters, run)
		NaturalDeathsInForce(index, &projection, &p, parameters)
		NaturalDeathsPaidUp(index, &projection, &p, parameters)
		NaturalDeathsPremiumWaiver(index, &projection, &p, parameters)
		NaturalDeathsTemporaryWaivers(index, &projection, &p, parameters)

		NumberOfDeathsAccident(index, &projection, &p, parameters)
		AccidentDeathsPaidUp(&projection, &p, index)
		AccidentDeathsPremiumWaiver(index, &projection, &p, parameters)
		AccidentDeathsTemporaryPremiumWaiver(index, &projection, &p, parameters)

		NumberOfLapses(index, &projection, &p, modelPoint, parameters)
		NumberOfDisabilities(index, &projection, &p, modelPoint, parameters)
		NumberOfRetrenchments(index, &projection, &p, modelPoint, parameters)

		InitialPaidUp(&projection, &p, modelPoint, features, index, parameters)
		InitialPremiumWaivers(index, &projection, &p, modelPoint, features, parameters)
		InitialPolicy(index, &projection, modelPoint, parameters)

		TotalIncrementalLapses(index, &projection, &p, parameters)
		TotalIncrementalNaturalDeaths(&projection, &p, index, parameters)
		TotalIncrementalDisabilities(&projection, &p, index, parameters)
		TotalIncrementalRetrenchments(&projection, &p, index, parameters)
		TotalIncrementalAccidentalDeaths(&projection, &p, index, parameters)

		//calculates if credit life
		if features.CreditLife || features.RetrenchmentBenefit {
			CalculatedInstalment(&projection, &p, modelPoint, parameters)
		}

		//Only runs if credit_life benefit structure is chosen
		if features.CreditLife {
			OutstandingSumAssured(&projection, &p, modelPoint, parameters, features)
		}

		//selected feature directs whether its fixed base lump sum or outstanding loan.Default is 0
		SumAssured(&projection, &p, modelPoint, parameters, features)

		if features.RiderBenefit {
			RiderSumAssured(&projection, modelPoint, parameters, features)
		}
		if features.StandardAdditionalLumpSum {
			StandardAdditionalLumpSum(&projection, modelPoint, parameters)
		}
		if features.AnnuityIncome {
			AnnuityIncome(&projection, modelPoint, parameters)
		}
		Premium(&projection, modelPoint, features, parameters)
		PremiumIncome(&projection, &p, modelPoint, parameters)
		PremiumNotReceived(&projection, features, parameters, modelPoint)
		// Unit fund

		if features.UnitFund || features.WithProfit {
			UnitFund(&projection, &p, modelPoint, parameters, features, valYear, valMonth, run)
		}

		Commission(&projection, &p, modelPoint, features, parameters)
		ClawBack(modelPoint, &projection, features, parameters)
		NonLifeRiskOutgo(&projection, &p, features, parameters, modelPoint)
		DeathOutgo(&projection, &p, features, parameters, modelPoint)
		AccidentalDeathOutgo(&projection, features, modelPoint, multipliers, parameters)
		CashBackOnSurvival(&projection, modelPoint, parameters, features)
		CashBackOnDeath(&projection, modelPoint, parameters, features)
		Rider(&projection, features, modelPoint, parameters)
		DisabilityOutgo(&projection, &p, modelPoint, parameters, features)
		RetrenchmentOutgo(&projection, modelPoint, parameters, features)

		if features.AnnuityIncome {
			AnnuityOutgo(&projection, &p, modelPoint, parameters)
		}

		MaturityOutgo(&projection, &p, modelPoint, parameters, features, valYear, valMonth, run)
		Expenses(&projection, &p, modelPoint, parameters, features, run, productMargins, shock)
		CoverageUnits(&projection, modelPoint, features, parameters)
		NetCashFlow(&projection, &p, modelPoint, parameters, features)
		Reinsurance(&projection, &p, features, modelPoint, parameters)

		//Set values to be used on the next loop
		projections = append(projections, projection)
		p = projection
	}

	for i := len(projections) - 1; i >= 0; i-- {
		discountedErr := CalculateReservesAndDiscountedValues(i, valYear, run.YieldcurveMonth, &projections[i], productMargins, run, modelPoint, parameters, features, &projections, productShocks)
		if discountedErr != nil {
			return &ModelPointError{
				ModelPoint: modelPoint,
				Error:      errors.Wrap(discountedErr, "failed to calculate reserves and discounted values"),
			}
		}
	}

	netCashFlow := 0.0
	reserveChangeA := 0.0
	investIncomeA := 0.0
	reserveChange := 0.0
	investIncome := 0.0

	prevMonth := 0

	for i := range projections {

		var profitAdjustment float64
		if features.RenewableProfitAdjustment {
			profitAdjustment = GetProfitRenewalAdjustment(&projections[i], modelPoint, parameters)
		} else {
			profitAdjustment = 0
		}

		if i == 1 && projections[i].ValuationTimeMonth == 1 {
			//Seeding the variables
			netCashFlow = projections[i-1].NetCashFlow + projections[i].NetCashFlow // point of sale and at time 1
			projections[i].ChangeInReserves = projections[i].Reserves - projections[i-1].Reserves
			projections[i].ChangeInReservesAdjusted = projections[i].ReservesAdjusted - projections[i-1].ReservesAdjusted

			interestAccrualErr := InterestAccrual(i, valYear, run.YieldcurveMonth, &projections[i], projections[i-1].Reserves, projections[i-1].ReservesAdjusted, run, modelPoint, parameters, productMargins, productShocks)
			if err != nil {
				return &ModelPointError{
					ModelPoint: modelPoint,
					Error:      errors.Wrap(interestAccrualErr, "failed to calculate interest accrual"),
				}
			}
			var forwardRate float64
			var forwardRateErr error
			if run.YieldCurveBasis == LockedInRates {
				forwardRate, forwardRateErr = GetLockedinForwardRateWithError(i, modelPoint.LockedInYear, modelPoint.LockedInMonth)
			}
			if run.YieldCurveBasis == Current {
				forwardRate, forwardRateErr = GetForwardRateWithError(i, run.YieldcurveYear, run.YieldcurveMonth, parameters.YieldCurveCode)
			}

			if forwardRateErr != nil {
				return &ModelPointError{
					ModelPoint: modelPoint,
					Error:      errors.Wrap(forwardRateErr, "failed to get forward rate"),
				}
			}
			accrualFactor := math.Pow(1.0+forwardRate, 1.0/12.0)
			investIncomeA = projections[i].InvestmentIncomeAdjusted
			investIncome = projections[i].InvestmentIncome
			reserveChangeA = projections[i-1].ReservesAdjusted*accrualFactor + projections[i].ChangeInReservesAdjusted //point of sale and time 1
			//reserveChange = projections[i].ChangeInReserves
			reserveChange = projections[i-1].Reserves*accrualFactor + projections[i].ChangeInReserves
			prevMonth = projections[i].ProjectionMonth
			projections[i].ProfitAdjustment = (-netCashFlow - reserveChange + investIncome) * profitAdjustment
			projections[i].ProfitAdjustmentAdjusted = (-netCashFlow - reserveChangeA + investIncomeA) * profitAdjustment
			projections[i].Profit = utils.FloatPrecision(-netCashFlow-reserveChange+investIncome+projections[i].ProfitAdjustment, 6)
			projections[i].ProfitAdjusted = utils.FloatPrecision(-netCashFlow-reserveChangeA+investIncomeA+projections[i].ProfitAdjustmentAdjusted, 6)
			projections[i].CorporateTax = utils.FloatPrecision(math.Max((-netCashFlow-reserveChange+investIncome+projections[i].ProfitAdjustment)*parameters.CorporateTaxRate, 0), 2)
			projections[i].CorporateTaxAdjusted = utils.FloatPrecision(math.Max((-netCashFlow-reserveChangeA+investIncomeA+projections[i].ProfitAdjustmentAdjusted)*parameters.CorporateTaxRate, 0), 2)

			continue
		}
		if i > 0 {

			//DB.Save(&projections[i])
			netCashFlow = projections[i].NetCashFlow

			projections[i].ChangeInReserves = projections[i].Reserves - projections[i-1].Reserves
			projections[i].ChangeInReservesAdjusted = projections[i].ReservesAdjusted - projections[i-1].ReservesAdjusted

			reserveChange = projections[i].Reserves - projections[i-1].Reserves
			reserveChangeA = projections[i].ReservesAdjusted - projections[i-1].ReservesAdjusted

			interestAccrualErr := InterestAccrual(i, valYear, run.YieldcurveMonth, &projections[i], projections[i-1].Reserves, projections[i-1].ReservesAdjusted, run, modelPoint, parameters, productMargins, productShocks)
			if err != nil {
				return &ModelPointError{
					ModelPoint: modelPoint,
					Error:      errors.Wrap(interestAccrualErr, "failed to calculate interest accrual"),
				}
			}
			investIncome = projections[i].InvestmentIncome
			investIncomeA = projections[i].InvestmentIncomeAdjusted

			prevMonth = projections[i].ProjectionMonth

			if prevMonth <= parameters.CalculatedTerm {
				projections[i].ProfitAdjustment = (-netCashFlow - reserveChange + investIncome) * profitAdjustment
				projections[i].ProfitAdjustmentAdjusted = (-netCashFlow - reserveChangeA + investIncomeA) * profitAdjustment
				projections[i].Profit = utils.FloatPrecision(-netCashFlow-reserveChange+investIncome+projections[i].ProfitAdjustment, 2)
				projections[i].ProfitAdjusted = utils.FloatPrecision(-netCashFlow-reserveChangeA+investIncomeA+projections[i].ProfitAdjustmentAdjusted, 2)
				projections[i].CorporateTax = utils.FloatPrecision(math.Max((-netCashFlow-reserveChange+investIncome+projections[i].ProfitAdjustment)*parameters.CorporateTaxRate, 0), 2)
				projections[i].CorporateTaxAdjusted = utils.FloatPrecision(math.Max((-netCashFlow-reserveChangeA+investIncomeA+projections[i].ProfitAdjustmentAdjusted)*parameters.CorporateTaxRate, 0), 2)
			}
		}

	}

	for i := len(projections) - 1; i >= 0; i-- {
		var forwardRate float64
		var rdr float64
		if i == len(projections)-1 {
			projections[i].DiscountedProfit = 0
			projections[i].DiscountedProfitAdjusted = 0
		} else {
			if run.YieldCurveBasis == Current {
				var fwRate float64
				var forwardRateErr error
				fwRate, forwardRateErr = GetForwardRateWithError(i+1, run.YieldcurveYear, run.YieldcurveMonth, parameters.YieldCurveCode)
				if forwardRateErr != nil {
					return &ModelPointError{
						ModelPoint: modelPoint,
						Error:      errors.Wrap(forwardRateErr, "failed to get forward rate"),
					}
				}
				forwardRate = fwRate

				//fwRate, err = GetForwardRateWithError(i+1, run.YieldcurveYear, run.YieldcurveMonth, parameters.YieldCurveCode, run.RunBasis)
				//if err != nil {
				//	return &ModelPointError{
				//		ModelPoint: modelPoint,
				//		Error:      errors.Wrap(err, "failed to get forward rate for risk discount rate"),
				//	}
				//}
				rdr = parameters.ShareholdersRequiredMargin + fwRate
			} else {
				var fwRate float64
				var forwardRateErr error
				fwRate, forwardRateErr = GetLockedinForwardRateWithError(i+1, modelPoint.LockedInYear, modelPoint.LockedInMonth)
				if forwardRateErr != nil {
					return &ModelPointError{
						ModelPoint: modelPoint,
						Error:      errors.Wrap(forwardRateErr, "failed to get forward rate"),
					}
				}
				forwardRate = fwRate

				rdr = parameters.ShareholdersRequiredMargin + fwRate
			}
			discountFactor := math.Pow(1.0+forwardRate, -1.0/12.0)
			riskDiscountFactor := math.Pow(1.0+rdr, -1.0/12.0)

			projections[i].RiskDiscountRate = rdr
			projections[i].RiskDiscountRateAdjusted = rdr //+ productMargins.InvestmentMargin

			projections[i].DiscountedInvestmentIncome = (projections[i+1].InvestmentIncome + projections[i+1].DiscountedInvestmentIncome) * discountFactor
			projections[i].DiscountedInvestmentIncomeAdjusted = (projections[i+1].InvestmentIncomeAdjusted + projections[i+1].DiscountedInvestmentIncomeAdjusted) * discountFactor
			projections[i].DiscountedProfit = (projections[i+1].Profit + projections[i+1].DiscountedProfit) * discountFactor
			projections[i].DiscountedProfitAdjusted = (projections[i+1].ProfitAdjusted + projections[i+1].DiscountedProfitAdjusted) * discountFactor

			// discounted using Shareholders Required Return
			projections[i].DiscountedProfitAdjustment = (projections[i+1].ProfitAdjustment + projections[i+1].DiscountedProfitAdjustment) * riskDiscountFactor
			projections[i].DiscountedProfitAdjustmentAdjusted = (projections[i+1].ProfitAdjustmentAdjusted + projections[i+1].DiscountedProfitAdjustmentAdjusted) * riskDiscountFactor

			projections[i].VIF = (projections[i+1].Profit + projections[i+1].VIF) * riskDiscountFactor
			projections[i].VIFAdjusted = (projections[i+1].ProfitAdjusted + projections[i+1].VIFAdjusted) * riskDiscountFactor
			projections[i].DiscountedCorporateTax = (projections[i+1].CorporateTax + projections[i+1].DiscountedCorporateTax) * riskDiscountFactor
			projections[i].DiscountedCorporateTaxAdjusted = (projections[i+1].CorporateTaxAdjusted + projections[i+1].DiscountedCorporateTaxAdjusted) * riskDiscountFactor
		}
	}

	//Only adding projections for the first modelpoint to act as a control.
	if mpIndex == 0 {
		err = DB.CreateInBatches(&projections, 100).Error
		if err != nil {
			log.WithField("modelpoint", modelPoint).Error("failed to save projections")
		}
	}

	//out <- projections[:run.AggregationPeriod]
	if len(projections) >= run.AggregationPeriod {
		UpdateAggregatedProjections(projections[:run.AggregationPeriod], aggProjs, jobID)
	} else {
		UpdateAggregatedProjections(projections, aggProjs, jobID)
	}

	if run.IFRS17Indicator {
		if len(projections) >= ScopedAggregatedProjectionPeriod {
			UpdateScopedAggregatedProjections(projections[:ScopedAggregatedProjectionPeriod], sapMap, jobID)
		} else {
			UpdateScopedAggregatedProjections(projections, sapMap, jobID)
		}
	}

	if run.LICIndicator {
		jobMu := GetJobMutex(jobID)
		jobMu.Lock()
		UpdateLICProjections(lic, projections)
		jobMu.Unlock()
	}

	//DB.Where("id = ?", job.ID).Find(&job)
	//if mpIndex+1 > jobProduct.PointsDone {
	//	jobProduct.PointsDone = mpIndex + 1
	//	//fmt.Println(jobProduct.PointsDone)
	//}
	//jobProduct.PointsDone = mpIndex //+= 1

	//err = DB.Save(jobProduct).Error

	mpIndex = mpIndex + 1
	if mpIndex > jobProduct.PointsDone {
		jobProduct.PointsDone = mpIndex
	}

	err = DB.Model(&models.JobProduct{}).
		Where("id = ? AND points_done < ?", jobProduct.ID, mpIndex).
		UpdateColumn("points_done", mpIndex).Error

	if err != nil {
		return &ModelPointError{
			ModelPoint: modelPoint,
			Error:      err,
		}
	}

	//job.PointsDone = job.PointsDone + 1
	//DB.Where("id = ?", job.ID).Save(&job)
	projections = nil
	//endLoopTime := time.Since(startLoopTime)
	return nil
}

func UpdateLICProjections(lic *[]models.LICAggregatedProjections, projections []models.Projection) {
	//TODO: Update LIC projections
	for k, projection := range projections {
		if k < LICAggregatedProjectionPeriod {
			var aggProj = models.LICAggregatedProjections{}
			err := copier.Copy(&aggProj, &projection)
			if err != nil {
				log.WithField("error", err).Error("failed to copy lic projection")
			}

			if i, ok := utils.LICAggregatedProjectionsContain(*lic, aggProj); ok {

				mutable := reflect.ValueOf((*lic)[i]).Elem()
				mutable2 := reflect.ValueOf(&aggProj).Elem()
				for j := 0; j < mutable.NumField(); j++ {
					if j == 0 {
						continue
					} else {
						if mutable.Field(j).Type().Kind() == reflect.Float64 {
							mutable.Field(j).SetFloat(mutable.Field(j).Float() + mutable2.Field(j).Float())
						}
					}
				}
			} else {
				*lic = append(*lic, aggProj)
			}
		}
	}
}

// current implementation of UpdateAggregatedProjections
func UpdateAggregatedProjections(projections []models.Projection, aggProjs *map[string]models.Projection, jobID int) {
	jobMu := GetJobMutex(jobID)
	jobMu.Lock()
	defer jobMu.Unlock()
	for _, projection := range projections {
		// TODO: This is a hack. The key should be unique.
		//key := strconv.Itoa(i) + "_" + projection.ProductCode + "_" + strconv.Itoa(projection.JobProductID) + "_" + strconv.Itoa(projection.SpCode)
		key := strings.Join([]string{
			projection.ProductCode,
			strconv.Itoa(projection.JobProductID),
			strconv.Itoa(projection.SpCode),
			strconv.Itoa(projection.ProjectionMonth),
		}, "_")
		agg, exists := (*aggProjs)[key]

		if exists {
			mutable := reflect.ValueOf(&agg).Elem()
			mutable2 := reflect.ValueOf(&projection).Elem()

			for j := 0; j < mutable.NumField(); j++ {
				if j == 0 {
					continue
				} else {

					if mutable.Field(j).Type().Kind() == reflect.Float64 {
						mutable.Field(j).SetFloat(mutable.Field(j).Float() + mutable2.Field(j).Float())
					}
				}
			}
			(*aggProjs)[key] = agg
		} else {
			(*aggProjs)[key] = projection
		}

	}
}

func UpdateScopedAggregatedProjections(projections []models.Projection, sapMap *map[string]models.Projection, jobID int) {
	jobMu := GetJobMutex(jobID)
	jobMu.Lock()
	defer jobMu.Unlock()
	for k, projection := range projections {
		if k < ScopedAggregatedProjectionPeriod {
			key := projection.RunDate + "_" + strconv.Itoa(projection.ProjectionMonth) + "_" + projection.IFRS17Group
			//var aggProj = models.ScopedAggregatedProjection{}
			//err := copier.Copy(&aggProj, &projection)
			//aggProj.ID = 0
			//agg, exists := (*sapMap)[key]
			//
			//if err != nil {
			//	fmt.Println(err)
			//}

			sm, exists := (*sapMap)[key]

			if exists {
				mutable := reflect.ValueOf(&sm).Elem()
				mutable2 := reflect.ValueOf(&projection).Elem()

				for j := 0; j < mutable.NumField(); j++ {
					if j == 0 {
						continue
					} else {
						if mutable.Field(j).Type().Kind() == reflect.Float64 {
							mutable.Field(j).SetFloat(mutable.Field(j).Float() + mutable2.Field(j).Float())
						}
					}
				}
				(*sapMap)[key] = sm
			} else {
				(*sapMap)[key] = projection
			}
		}
	}
}

func CalculateProductConfig(parameters *models.ProductParameters, mp models.ProductModelPoint, features models.ProductFeatures, multipliers models.ProductAccidentalBenefitMultiplier, run models.RunParameters) error {
	//Calculated term
	var err error

	// TODO: found that the product parameters are already being passed into the function.
	//var params models.ProductParameters
	//key := fmt.Sprintf("%s-%s-%d", strings.ToLower(mp.ProductCode), strings.ToLower(run.RunBasis), run.ParameterYear)
	//cached, found := Cache.Get(key)
	//if found {
	//	params = cached.(models.ProductParameters)
	//} else {
	//	err := DB.Where("product_code=? and basis=? and year=?", mp.ProductCode, run.RunBasis, run.ParameterYear).First(&params).Error
	//	if err == nil {
	//		Cache.Set(key, params, 1)
	//	}
	//}

	if mp.MemberType == "CH" {
		parameters.CalculatedTerm = int(math.Min(math.Max(float64(parameters.ChildExitAge-mp.AgeAtEntry), 0.0)*12.0, float64(mp.Term)))
	} else if features.CreditLife {
		parameters.CalculatedTerm = mp.OutstandingTermMonths + mp.DurationInForceMonths + 1
	} else if features.WholeOfLife {
		parameters.CalculatedTerm = mp.Term //(120 - mp.AgeAtEntry) * 12
	} else {
		parameters.CalculatedTerm = mp.Term
	}

	if features.AccidentalDeathBenefit {
		switch mp.MemberType {
		case "MM":
			parameters.AccidentalBenefitMultiplier = multipliers.MainMember
		case "SP":
			parameters.AccidentalBenefitMultiplier = multipliers.Spouse
		case "CH":
			parameters.AccidentalBenefitMultiplier = multipliers.Child
		case "PAR":
			parameters.AccidentalBenefitMultiplier = multipliers.Parent
		case "EXT":
			parameters.AccidentalBenefitMultiplier = multipliers.ExtendedFamily
		}
	} else {
		parameters.AccidentalBenefitMultiplier = 0
	}

	if mp.MemberType == "MM" {
		parameters.MainMemberIndicator = true
	} else {
		parameters.MainMemberIndicator = false
	}
	if parameters.MainMemberIndicator {
		parameters.OtherLivesIndicator = false
	} else {
		parameters.OtherLivesIndicator = true
	}

	if features.ProportionalReinsurance {
		var val float64

		val, err = GetReinsuranceStructureWithError(mp.ProductCode, mp.TreatyYear, "flat_annual_reins_prem_rate")
		if err != nil {
			return errors.Wrap(err, "failed to get flat annual reins prem rate")
		}
		parameters.FlatAnnualReinsPremRate = val

		val, err = GetReinsuranceStructureWithError(mp.ProductCode, mp.TreatyYear, "level1_ceded_proportion")
		if err != nil {
			return errors.Wrap(err, "failed to get level1 ceded proportion")
		}
		parameters.Level1CededProp = val

		val, err = GetReinsuranceStructureWithError(mp.ProductCode, mp.TreatyYear, "level1_lowerbound")
		if err != nil {
			return errors.Wrap(err, "failed to get level1 lowerbound")
		}
		parameters.Level1Lowerbound = val

		val, err = GetReinsuranceStructureWithError(mp.ProductCode, mp.TreatyYear, "level1_upperbound")
		if err != nil {
			return errors.Wrap(err, "failed to get level1 upperbound")
		}
		parameters.Level1Upperbound = val

		val, err = GetReinsuranceStructureWithError(mp.ProductCode, mp.TreatyYear, "level2_ceded_proportion")
		if err != nil {
			return errors.Wrap(err, "failed to get level2 ceded proportion")
		}
		parameters.Level2CededProp = val

		val, err = GetReinsuranceStructureWithError(mp.ProductCode, mp.TreatyYear, "level2_lowerbound")
		if err != nil {
			return errors.Wrap(err, "failed to get level2 lowerbound")
		}
		parameters.Level2Lowerbound = val

		val, err = GetReinsuranceStructureWithError(mp.ProductCode, mp.TreatyYear, "level2_upperbound")
		if err != nil {
			return errors.Wrap(err, "failed to get level2 upperbound")
		}
		parameters.Level2Upperbound = val

		val, err = GetReinsuranceStructureWithError(mp.ProductCode, mp.TreatyYear, "level3_ceded_proportion")
		if err != nil {
			return errors.Wrap(err, "failed to get level3 ceded proportion")
		}
		parameters.Level3CededProp = val

		val, err = GetReinsuranceStructureWithError(mp.ProductCode, mp.TreatyYear, "level3_lowerbound")
		if err != nil {
			return errors.Wrap(err, "failed to get level3 lowerbound")
		}
		parameters.Level3Lowerbound = val

		val, err = GetReinsuranceStructureWithError(mp.ProductCode, mp.TreatyYear, "level3_upperbound")
		if err != nil {
			return errors.Wrap(err, "failed to get level3 upperbound")
		}
		parameters.Level3Upperbound = val

		val, err = GetReinsuranceStructureWithError(mp.ProductCode, mp.TreatyYear, "ceding_commission")
		if err != nil {
			return errors.Wrap(err, "failed to get ceding commission")
		}
		parameters.CedingCommission = val
	} else {
		parameters.FlatAnnualReinsPremRate = 0
		parameters.Level1CededProp = 0
		parameters.Level1Lowerbound = 0
		parameters.Level1Upperbound = 0
		parameters.Level2CededProp = 0
		parameters.Level2Lowerbound = 0
		parameters.Level2Upperbound = 0
		parameters.Level3CededProp = 0
		parameters.Level3Lowerbound = 0
		parameters.Level3Upperbound = 0
		parameters.CedingCommission = 0
	}

	if features.RiderBenefit {
		riderSumAssured, err := GetRidersWithError(mp, features)
		if err != nil {
			return errors.Wrap(err, "failed to get riders")
		}
		parameters.RiderSumAssured = riderSumAssured
	}

	if features.AnnuityIncome {
		benefitMultiplier, err := GetBenefitMultiplierWithError(mp)
		if err != nil {
			return errors.Wrap(err, "failed to get benefit multiplier")
		}
		parameters.BenefitMultiplier = benefitMultiplier
	}

	return nil
}

func CalculateTimeAndAge(projection *models.Projection, modelPoint models.ProductModelPoint, i int, jobProductId int, features models.ProductFeatures) {
	projection.JobProductID = jobProductId
	projection.ProjectionMonth = i
	projection.ProjectionYear = int(math.Ceil(float64(i) / 12.0))
	projection.PolicyNumber = modelPoint.PolicyNumber
	projection.ValuationTimeMonth = modelPoint.DurationInForceMonths + i
	projection.ValuationTimeYear = utils.FloatPrecision(float64(projection.ValuationTimeMonth)/12, 2)
	//if features.FuneralCover {
	projection.MainMemberAgeNextBirthday = int(math.Ceil(float64(modelPoint.MainMemberAgeAtEntry) + projection.ValuationTimeYear))
	//}
	//if features.CreditLife {
	//	projection.MainMemberAgeNextBirthday = int(math.Ceil(float64(modelPoint.AgeAtEntry) + projection.ValuationTimeYear))
	//}
	projection.AgeNextBirthday = int(math.Ceil(float64(modelPoint.AgeAtEntry) + projection.ValuationTimeYear))
}

func CalculateReservesAndDiscountedValues(index, valYear,
	valMonth int, projection *models.Projection, margins models.ProductMargins,
	run models.RunParameters, mp models.ProductModelPoint, params models.ProductParameters,
	features models.ProductFeatures, projections *[]models.Projection, productShocks []models.ProductShock) error {
	//we are working on the previous projection results, right?
	var forwardRate float64
	var err error
	var resp float64 = 0
	var realResp float64 = 0
	var shockedReal float64 = 0
	var infl float64 = 0
	if run.YieldCurveBasis == Current {
		resp, err = GetForwardRateWithError(min(index+1, 1440), run.YieldcurveYear, run.YieldcurveMonth, params.YieldCurveCode)
		if err != nil {
			return errors.Wrap(err, "failed to get forward rate")
		}

		if run.ShockSettings.RealYieldCurve {
			infl, err = GetInflationFactorWithError(min(index+1, 1440), run.YieldcurveYear, run.YieldcurveMonth, params.YieldCurveCode)
			if err != nil {
				return errors.Wrap(err, "failed to get inflation factor")
			}
			realResp = (1+resp)/(1+infl) - 1
		}

	} else {

		resp, err = GetLockedinForwardRateWithError(min(index+1, 1440), mp.LockedInYear, mp.LockedInMonth)
		if err != nil {
			return errors.Wrap(err, "failed to get forward rate")
		}

		if run.ShockSettings.RealYieldCurve {
			infl, err = GetLockedinInflationFactorWithError(index+1, mp.LockedInYear, mp.LockedInMonth)
			if err != nil {
				return errors.Wrap(err, "failed to get inflation factor")
			}
			realResp = (1+resp)/(1+infl) - 1
		}
	}

	if run.ShockSettings.NominalYieldCurve {
		var shock models.ProductShock
		if projection.ProjectionMonth > 0 && run.ShockSettings.ShockBasis != "" {
			shock, err = GetShock(projection.ProjectionMonth, run.ShockSettings.ShockBasis, len(productShocks))
			if err != nil {
				log.WithField("projection_month", projection.ProjectionMonth).
					WithField("shock_basis", run.ShockSettings.ShockBasis).
					WithField("product_shocks", productShocks).
					Error("failed to get shock")
			} else {
				forwardRate = math.Max(0, math.Min(1, resp*(1+shock.MultiplicativeYieldCurve)+shock.AdditiveYieldCurve))
			}
		}
	}

	if run.ShockSettings.RealYieldCurve {
		var shock models.ProductShock
		if projection.ProjectionMonth > 0 && run.ShockSettings.ShockBasis != "" {
			shock, err = GetShock(projection.ProjectionMonth, run.ShockSettings.ShockBasis, len(productShocks))
			if err != nil {
				log.WithField("projection_month", projection.ProjectionMonth).WithField("shock_basis", run.ShockSettings.ShockBasis).Error("failed to get shock")
			}

		}
		shockedReal = math.Max(0, math.Min(1, realResp*(1+shock.MultiplicativeYieldCurve)+shock.AdditiveYieldCurve))
		forwardRate = (1+shockedReal)*(1+infl) - 1
	}

	if !run.ShockSettings.NominalYieldCurve && !run.ShockSettings.RealYieldCurve {
		forwardRate = resp
	}

	discountFactor := math.Pow(1.0+forwardRate, -1/12.0)
	discountFactorAdjusted := math.Pow(1.0+forwardRate+(margins.InvestmentMargin), -1/12.0)

	projection.ValuationRate = forwardRate
	projection.ValuationRateAdjusted = forwardRate + margins.InvestmentMargin

	if index == len(*projections)-1 {
		projection.DiscountedPremiumIncome = 0
		projection.DiscountedPremiumIncomeAdjusted = 0
		projection.DiscountedPremiumNotReceived = 0
		projection.DiscountedPremiumNotReceivedAdjusted = 0
		projection.DiscountedCommission = 0
		projection.DiscountedCommissionAdjusted = 0
		projection.DiscountedRenewalCommission = 0
		projection.DiscountedRenewalCommissionAdjusted = 0
		projection.DiscountedClawBack = 0
		projection.DiscountedClawBackAdjusted = 0
		projection.DiscountedDeathOutgo = 0
		projection.DiscountedDeathOutgoAdjusted = 0
		projection.DiscountedAccidentalDeathOutgo = 0
		projection.DiscountedAccidentalDeathOutgoAdjusted = 0
		projection.DiscountedCashBackOnSurvival = 0
		projection.DiscountedCashBackOnSurvivalAdjusted = 0
		projection.DiscountedCashBackOnDeath = 0
		projection.DiscountedCashBackOnDeathAdjusted = 0
		projection.DiscountedRider = 0
		projection.DiscountedRiderAdjusted = 0
		projection.DiscountedDisabilityOutgo = 0
		projection.DiscountedDisabilityOutgoAdjusted = 0
		projection.DiscountedRetrenchmentOutgo = 0
		projection.DiscountedRetrenchmentOutgoAdjusted = 0
		projection.DiscountedAnnuityOutgo = 0
		projection.DiscountedAnnuityOutgoAdjusted = 0
		projection.DiscountedInitialExpenses = 0
		projection.DiscountedInitialExpensesAdjusted = 0
		projection.DiscountedRenewalExpenses = 0
		projection.DiscountedRenewalExpensesAdjusted = 0
		projection.DiscountedReinsurancePremium = 0
		projection.DiscountedReinsurancePremiumAdjusted = 0
		projection.DiscountedReinsuranceCedingCommission = 0
		projection.DiscountedReinsuranceCedingCommissionAdjusted = 0
		projection.DiscountedReinsuranceClaims = 0
		projection.DiscountedReinsuranceClaimsAdjusted = 0
		projection.DiscountedNetReinsuranceCashflow = 0
		projection.DiscountedNetReinsuranceCashflowAdjusted = 0
		projection.SumCoverageUnits = 0
		projection.DiscountedCoverageUnits = 0
		projection.DiscountedSurrenderOutgo = 0
		projection.DiscountedSurrenderOutgoAdjusted = 0

	} else {
		projection.DiscountedPremiumIncome = utils.FloatPrecision((*projections)[index+1].PremiumIncome+(*projections)[index+1].DiscountedPremiumIncome*discountFactor, defaultPrecision)
		projection.DiscountedPremiumIncomeAdjusted = utils.FloatPrecision((*projections)[index+1].PremiumIncomeAdjusted+(*projections)[index+1].DiscountedPremiumIncomeAdjusted*discountFactorAdjusted, defaultPrecision)

		projection.DiscountedPremiumNotReceived = utils.FloatPrecision((*projections)[index+1].PremiumNotReceived+(*projections)[index+1].DiscountedPremiumNotReceived*discountFactor, defaultPrecision)
		projection.DiscountedPremiumNotReceivedAdjusted = utils.FloatPrecision((*projections)[index+1].PremiumNotReceivedAdjusted+(*projections)[index+1].DiscountedPremiumNotReceivedAdjusted*discountFactorAdjusted, defaultPrecision)

		projection.DiscountedCommission = utils.FloatPrecision((*projections)[index+1].Commission+(*projections)[index+1].DiscountedCommission*discountFactor, defaultPrecision)
		projection.DiscountedCommissionAdjusted = utils.FloatPrecision((*projections)[index+1].CommissionAdjusted+(*projections)[index+1].DiscountedCommissionAdjusted*discountFactorAdjusted, defaultPrecision)

		projection.DiscountedRenewalCommission = utils.FloatPrecision((*projections)[index+1].RenewalCommission+(*projections)[index+1].DiscountedRenewalCommission*discountFactor, defaultPrecision)
		projection.DiscountedRenewalCommissionAdjusted = utils.FloatPrecision((*projections)[index+1].RenewalCommissionAdjusted+(*projections)[index+1].DiscountedRenewalCommissionAdjusted*discountFactorAdjusted, defaultPrecision)

		projection.DiscountedClawBack = utils.FloatPrecision(((*projections)[index+1].ClawBack+(*projections)[index+1].DiscountedClawBack)*discountFactor, defaultPrecision)
		projection.DiscountedClawBackAdjusted = utils.FloatPrecision(((*projections)[index+1].ClawBackAdjusted+(*projections)[index+1].DiscountedClawBackAdjusted)*discountFactorAdjusted, defaultPrecision)

		projection.DiscountedDeathOutgo = utils.FloatPrecision(((*projections)[index+1].DeathOutgo+(*projections)[index+1].DiscountedDeathOutgo)*discountFactor, defaultPrecision)
		projection.DiscountedDeathOutgoAdjusted = utils.FloatPrecision(((*projections)[index+1].DeathOutgoAdjusted+(*projections)[index+1].DiscountedDeathOutgoAdjusted)*discountFactorAdjusted, defaultPrecision)

		projection.DiscountedAccidentalDeathOutgo = utils.FloatPrecision(((*projections)[index+1].AccidentalDeathOutgo+(*projections)[index+1].DiscountedAccidentalDeathOutgo)*discountFactor, defaultPrecision)
		projection.DiscountedAccidentalDeathOutgoAdjusted = utils.FloatPrecision(((*projections)[index+1].AccidentalDeathOutgoAdjusted+(*projections)[index+1].DiscountedAccidentalDeathOutgoAdjusted)*discountFactorAdjusted, defaultPrecision)

		projection.DiscountedCashBackOnSurvival = utils.FloatPrecision(((*projections)[index+1].CashBackOnSurvival+(*projections)[index+1].DiscountedCashBackOnSurvival)*discountFactor, defaultPrecision)
		projection.DiscountedCashBackOnSurvivalAdjusted = utils.FloatPrecision(((*projections)[index+1].CashBackOnSurvivalAdjusted+(*projections)[index+1].DiscountedCashBackOnSurvivalAdjusted)*discountFactorAdjusted, defaultPrecision)

		projection.DiscountedCashBackOnDeath = utils.FloatPrecision(((*projections)[index+1].CashBackOnDeath+(*projections)[index+1].DiscountedCashBackOnDeath)*discountFactor, defaultPrecision)
		projection.DiscountedCashBackOnDeathAdjusted = utils.FloatPrecision(((*projections)[index+1].CashBackOnDeathAdjusted+(*projections)[index+1].DiscountedCashBackOnDeathAdjusted)*discountFactorAdjusted, defaultPrecision)

		projection.DiscountedRider = utils.FloatPrecision(((*projections)[index+1].Rider+(*projections)[index+1].DiscountedRider)*discountFactor, defaultPrecision)
		projection.DiscountedRiderAdjusted = utils.FloatPrecision(((*projections)[index+1].RiderAdjusted+(*projections)[index+1].DiscountedRiderAdjusted)*discountFactorAdjusted, defaultPrecision)

		projection.DiscountedDisabilityOutgo = utils.FloatPrecision(((*projections)[index+1].DisabilityOutgo+(*projections)[index+1].DiscountedDisabilityOutgo)*discountFactor, defaultPrecision)
		projection.DiscountedDisabilityOutgoAdjusted = utils.FloatPrecision(((*projections)[index+1].DisabilityOutgoAdjusted+(*projections)[index+1].DiscountedDisabilityOutgoAdjusted)*discountFactorAdjusted, defaultPrecision)

		projection.DiscountedRetrenchmentOutgo = utils.FloatPrecision(((*projections)[index+1].RetrenchmentOutgo+(*projections)[index+1].DiscountedRetrenchmentOutgo)*discountFactor, defaultPrecision)
		projection.DiscountedRetrenchmentOutgoAdjusted = utils.FloatPrecision(((*projections)[index+1].RetrenchmentOutgoAdjusted+(*projections)[index+1].DiscountedRetrenchmentOutgoAdjusted)*discountFactorAdjusted, defaultPrecision)

		projection.DiscountedAnnuityOutgo = utils.FloatPrecision(((*projections)[index+1].AnnuityOutgo+(*projections)[index+1].DiscountedAnnuityOutgo)*discountFactor, defaultPrecision)
		projection.DiscountedAnnuityOutgoAdjusted = utils.FloatPrecision(((*projections)[index+1].AnnuityOutgoAdjusted+(*projections)[index+1].DiscountedAnnuityOutgoAdjusted)*discountFactorAdjusted, defaultPrecision)

		projection.DiscountedInitialExpenses = utils.FloatPrecision((*projections)[index+1].InitialExpenses+(*projections)[index+1].DiscountedInitialExpenses*discountFactor, defaultPrecision)
		projection.DiscountedInitialExpensesAdjusted = utils.FloatPrecision((*projections)[index+1].InitialExpensesAdjusted+(*projections)[index+1].DiscountedInitialExpensesAdjusted*discountFactorAdjusted, defaultPrecision)

		projection.DiscountedRenewalExpenses = utils.FloatPrecision((*projections)[index+1].RenewalExpenses+(*projections)[index+1].DiscountedRenewalExpenses*discountFactor, defaultPrecision)
		projection.DiscountedRenewalExpensesAdjusted = utils.FloatPrecision((*projections)[index+1].RenewalExpensesAdjusted+(*projections)[index+1].DiscountedRenewalExpensesAdjusted*discountFactorAdjusted, defaultPrecision)

		projection.DiscountedReinsurancePremium = utils.FloatPrecision(((*projections)[index+1].ReinsurancePremium+(*projections)[index+1].DiscountedReinsurancePremium)*discountFactor, defaultPrecision)
		projection.DiscountedReinsurancePremiumAdjusted = utils.FloatPrecision(((*projections)[index+1].ReinsurancePremiumAdjusted+(*projections)[index+1].DiscountedReinsurancePremiumAdjusted)*discountFactor, defaultPrecision)

		projection.DiscountedReinsuranceCedingCommission = utils.FloatPrecision(((*projections)[index+1].ReinsuranceCedingCommission+(*projections)[index+1].DiscountedReinsuranceCedingCommission)*discountFactor, defaultPrecision)
		projection.DiscountedReinsuranceCedingCommissionAdjusted = utils.FloatPrecision(((*projections)[index+1].ReinsuranceCedingCommissionAdjusted+(*projections)[index+1].DiscountedReinsuranceCedingCommissionAdjusted)*discountFactor, defaultPrecision)

		projection.DiscountedReinsuranceClaims = utils.FloatPrecision(((*projections)[index+1].ReinsuranceClaims+(*projections)[index+1].DiscountedReinsuranceClaims)*discountFactor, defaultPrecision)
		projection.DiscountedReinsuranceClaimsAdjusted = utils.FloatPrecision(((*projections)[index+1].ReinsuranceClaimsAdjusted+(*projections)[index+1].DiscountedReinsuranceClaimsAdjusted)*discountFactor, defaultPrecision)

		projection.DiscountedNetReinsuranceCashflow = utils.FloatPrecision(((*projections)[index+1].NetReinsuranceCashflow+(*projections)[index+1].DiscountedNetReinsuranceCashflow)*discountFactor, defaultPrecision)
		projection.DiscountedNetReinsuranceCashflowAdjusted = utils.FloatPrecision(((*projections)[index+1].NetReinsuranceCashflowAdjusted+(*projections)[index+1].DiscountedNetReinsuranceCashflowAdjusted)*discountFactor, defaultPrecision)

		projection.SumCoverageUnits = utils.FloatPrecision((*projections)[index+1].SumCoverageUnits+(*projections)[index+1].CoverageUnits, defaultPrecision)
		projection.DiscountedCoverageUnits = utils.FloatPrecision(((*projections)[index+1].CoverageUnits+(*projections)[index+1].DiscountedCoverageUnits)*discountFactor, defaultPrecision)

		if features.SurrenderBenefit {
			projection.DiscountedSurrenderOutgo = utils.FloatPrecision(((*projections)[index+1].SurrenderOutgo+(*projections)[index+1].DiscountedSurrenderOutgo)*discountFactor, defaultPrecision)
			projection.DiscountedSurrenderOutgoAdjusted = utils.FloatPrecision(((*projections)[index+1].SurrenderOutgoAdjusted+(*projections)[index+1].DiscountedSurrenderOutgoAdjusted)*discountFactorAdjusted, defaultPrecision)
		}

	}
	if features.UnitFund || features.WithProfit {

		if index == len(*projections)-1 {
			projection.DiscountedUnallocatedPremiumIncome = 0
			projection.DiscountedUnallocatedPremiumIncomeAdjusted = 0
			projection.DiscountedPolicyFee = 0
			projection.DiscountedPolicyFeeAdjusted = 0
			projection.DiscountedPremiumAdvisoryFeeAdjusted = 0
			projection.DiscountedPremiumAdvisoryFee = 0
			projection.DiscountedPremiumAdvisoryFeeAdjusted = 0
			projection.DiscountedFundAdvisoryFee = 0
			projection.DiscountedFundAdvisoryFeeAdjusted = 0
			projection.DiscountedPremiumAdvisoryCost = 0
			projection.DiscountedPremiumAdvisoryCostAdjusted = 0
			projection.DiscountedFundAdvisoryCost = 0
			projection.DiscountedFundAdvisoryCost = 0
			projection.DiscountedMaturityOutgo = 0
			projection.DiscountedMaturityOutgoAdjusted = 0
			projection.DiscountedFundRiskCharge = 0
			projection.DiscountedFundRiskChargeAdjusted = 0
			projection.DiscountedDeathOutgo = 0
			projection.DiscountedDeathOutgoAdjusted = 0
			projection.DiscountedAccidentalDeathOutgo = 0
			projection.DiscountedAccidentalDeathOutgoAdjusted = 0
			projection.DiscountedSurrenderPenaltyCharge = 0
			projection.DiscountedSurrenderPenaltyChargeAdjusted = 0
			projection.DiscountedFundAssetManagementCharges = 0
			projection.DiscountedFundAssetManagementChargesAdjusted = 0
			projection.DiscountedBsaShareholderCharge = 0
			projection.DiscountedBsaShareholderChargeAdjusted = 0
			projection.DiscountedGuaranteeCost = 0
			projection.DiscountedGuaranteeCostAdjusted = 0

		} else {
			// Start of the Month
			projection.DiscountedUnallocatedPremiumIncome = utils.FloatPrecision((*projections)[index+1].PremiumIncome-(*projections)[index+1].EAllocatedPremiumIncome+(*projections)[index+1].DiscountedUnallocatedPremiumIncome*discountFactor, defaultPrecision)
			projection.DiscountedUnallocatedPremiumIncomeAdjusted = utils.FloatPrecision((*projections)[index+1].PremiumIncomeAdjusted-(*projections)[index+1].EAllocatedPremiumIncomeAdjusted+(*projections)[index+1].DiscountedUnallocatedPremiumIncomeAdjusted*discountFactorAdjusted, defaultPrecision)

			projection.DiscountedPolicyFee = utils.FloatPrecision((*projections)[index+1].EPolicyFee+(*projections)[index+1].DiscountedPolicyFee*discountFactor, defaultPrecision)
			projection.DiscountedPolicyFeeAdjusted = utils.FloatPrecision((*projections)[index+1].EPolicyFeeAdjusted+(*projections)[index+1].DiscountedPolicyFeeAdjusted*discountFactorAdjusted, defaultPrecision)

			projection.DiscountedPremiumAdvisoryFee = utils.FloatPrecision((*projections)[index+1].EPremiumAdvisoryFee+(*projections)[index+1].DiscountedPremiumAdvisoryFee*discountFactor, defaultPrecision)
			projection.DiscountedPremiumAdvisoryFeeAdjusted = utils.FloatPrecision((*projections)[index+1].EPremiumAdvisoryFeeAdjusted+(*projections)[index+1].DiscountedPremiumAdvisoryFeeAdjusted*discountFactorAdjusted, defaultPrecision)

			projection.DiscountedFundAdvisoryFee = utils.FloatPrecision((*projections)[index+1].EFundAdvisoryFee+(*projections)[index+1].DiscountedFundAdvisoryFee*discountFactor, defaultPrecision)
			projection.DiscountedFundAdvisoryFeeAdjusted = utils.FloatPrecision((*projections)[index+1].EFundAdvisoryFeeAdjusted+(*projections)[index+1].DiscountedFundAdvisoryFeeAdjusted*discountFactorAdjusted, defaultPrecision)

			projection.DiscountedPremiumAdvisoryCost = utils.FloatPrecision((*projections)[index+1].EPremiumAdvisoryCost+(*projections)[index+1].DiscountedPremiumAdvisoryCost*discountFactor, defaultPrecision)
			projection.DiscountedPremiumAdvisoryCostAdjusted = utils.FloatPrecision((*projections)[index+1].EPremiumAdvisoryCostAdjusted+(*projections)[index+1].DiscountedPremiumAdvisoryCostAdjusted*discountFactorAdjusted, defaultPrecision)

			projection.DiscountedFundAdvisoryCost = utils.FloatPrecision((*projections)[index+1].EFundAdvisoryCost+(*projections)[index+1].DiscountedFundAdvisoryCost*discountFactor, defaultPrecision)
			projection.DiscountedFundAdvisoryCostAdjusted = utils.FloatPrecision((*projections)[index+1].EFundAdvisoryCostAdjusted+(*projections)[index+1].DiscountedFundAdvisoryCostAdjusted*discountFactorAdjusted, defaultPrecision)

			//End of the Month

			projection.DiscountedFundRiskCharge = utils.FloatPrecision(((*projections)[index+1].EFundRiskCharge+(*projections)[index+1].DiscountedFundRiskCharge)*discountFactor, defaultPrecision)
			projection.DiscountedFundRiskChargeAdjusted = utils.FloatPrecision(((*projections)[index+1].EFundRiskChargeAdjusted+(*projections)[index+1].DiscountedFundRiskChargeAdjusted)*discountFactorAdjusted, defaultPrecision)

			projection.DiscountedDeathOutgo = utils.FloatPrecision(((*projections)[index+1].DeathOutgo+(*projections)[index+1].DiscountedDeathOutgo)*discountFactor, defaultPrecision)
			projection.DiscountedDeathOutgoAdjusted = utils.FloatPrecision(((*projections)[index+1].DeathOutgoAdjusted+(*projections)[index+1].DiscountedDeathOutgoAdjusted)*discountFactorAdjusted, defaultPrecision)

			projection.DiscountedAccidentalDeathOutgo = utils.FloatPrecision(((*projections)[index+1].AccidentalDeathOutgo+(*projections)[index+1].DiscountedAccidentalDeathOutgo)*discountFactor, defaultPrecision)
			projection.DiscountedAccidentalDeathOutgoAdjusted = utils.FloatPrecision(((*projections)[index+1].AccidentalDeathOutgoAdjusted+(*projections)[index+1].DiscountedAccidentalDeathOutgoAdjusted)*discountFactorAdjusted, defaultPrecision)

			projection.DiscountedSurrenderPenaltyCharge = utils.FloatPrecision(((*projections)[index+1].ESurrenderPenaltyCharge+(*projections)[index+1].DiscountedSurrenderPenaltyCharge)*discountFactor, defaultPrecision)
			projection.DiscountedSurrenderPenaltyChargeAdjusted = utils.FloatPrecision(((*projections)[index+1].ESurrenderPenaltyChargeAdjusted+(*projections)[index+1].DiscountedSurrenderPenaltyChargeAdjusted)*discountFactorAdjusted, defaultPrecision)

			projection.DiscountedFundAssetManagementCharges = utils.FloatPrecision(((*projections)[index+1].EFundAssetManagementCharge+(*projections)[index+1].DiscountedFundAssetManagementCharges)*discountFactor, defaultPrecision)
			projection.DiscountedFundAssetManagementChargesAdjusted = utils.FloatPrecision(((*projections)[index+1].EFundAssetManagementChargeAdjusted+(*projections)[index+1].DiscountedFundAssetManagementChargesAdjusted)*discountFactorAdjusted, defaultPrecision)

			projection.DiscountedBsaShareholderCharge = utils.FloatPrecision(((*projections)[index+1].EBsaShareholderCharge+(*projections)[index+1].DiscountedBsaShareholderCharge)*discountFactor, defaultPrecision)
			projection.DiscountedBsaShareholderChargeAdjusted = utils.FloatPrecision(((*projections)[index+1].EBsaShareholderChargeAdjusted+(*projections)[index+1].DiscountedBsaShareholderChargeAdjusted)*discountFactorAdjusted, defaultPrecision)

			projection.DiscountedMaturityOutgo = utils.FloatPrecision(((*projections)[index+1].MaturityOutgo+(*projections)[index+1].DiscountedMaturityOutgo)*discountFactor, defaultPrecision)
			projection.DiscountedMaturityOutgoAdjusted = utils.FloatPrecision(((*projections)[index+1].MaturityOutgo+(*projections)[index+1].DiscountedMaturityOutgo)*discountFactorAdjusted, defaultPrecision)

			projection.DiscountedGuaranteeCost = utils.FloatPrecision(((*projections)[index+1].EGuaranteeCost+(*projections)[index+1].DiscountedGuaranteeCost)*discountFactor, defaultPrecision)
			projection.DiscountedGuaranteeCostAdjusted = utils.FloatPrecision(((*projections)[index+1].EGuaranteeCostAdjusted+(*projections)[index+1].DiscountedGuaranteeCostAdjusted)*discountFactorAdjusted, defaultPrecision)
		}
	}

	var variableIndicator float64
	if features.UnitFund || features.WithProfit {
		variableIndicator = 1.0
	} else {
		variableIndicator = 0.0
	}
	projection.Reserves = utils.FloatPrecision(
		projection.DiscountedPremiumNotReceived+
			projection.DiscountedCommission+
			projection.DiscountedRenewalCommission+
			projection.DiscountedDeathOutgo+
			projection.DiscountedAccidentalDeathOutgo+
			projection.DiscountedCashBackOnSurvival+
			projection.DiscountedCashBackOnDeath+
			projection.DiscountedRider+
			projection.DiscountedDisabilityOutgo+
			projection.DiscountedRetrenchmentOutgo+
			projection.DiscountedAnnuityOutgo+
			projection.DiscountedInitialExpenses+
			projection.DiscountedRenewalExpenses+
			projection.DiscountedPremiumAdvisoryCost+
			projection.DiscountedFundAdvisoryCost+
			projection.DiscountedMaturityOutgo+
			projection.DiscountedGuaranteeCost-
			projection.DiscountedUnallocatedPremiumIncome-
			projection.DiscountedPremiumAdvisoryFee-
			projection.DiscountedFundAdvisoryFee-
			projection.DiscountedPolicyFee-
			projection.DiscountedFundRiskCharge-
			projection.DiscountedSurrenderPenaltyCharge-
			projection.DiscountedFundAssetManagementCharges-
			projection.DiscountedBsaShareholderCharge-
			projection.DiscountedPremiumIncome*(1-variableIndicator)-
			projection.DiscountedClawBack, defaultPrecision)

	projection.ReservesAdjusted = utils.FloatPrecision(
		projection.DiscountedPremiumNotReceivedAdjusted+
			projection.DiscountedCommissionAdjusted+
			projection.DiscountedRenewalCommissionAdjusted+
			projection.DiscountedDeathOutgoAdjusted+
			projection.DiscountedAccidentalDeathOutgoAdjusted+
			projection.DiscountedCashBackOnSurvivalAdjusted+
			projection.DiscountedCashBackOnDeathAdjusted+
			projection.DiscountedRiderAdjusted+
			projection.DiscountedDisabilityOutgoAdjusted+
			projection.DiscountedRetrenchmentOutgoAdjusted+
			projection.DiscountedAnnuityOutgoAdjusted+
			projection.DiscountedInitialExpensesAdjusted+
			projection.DiscountedRenewalExpensesAdjusted+
			projection.DiscountedPremiumAdvisoryCostAdjusted+
			projection.DiscountedFundAdvisoryCostAdjusted+
			projection.DiscountedMaturityOutgoAdjusted+
			projection.DiscountedGuaranteeCostAdjusted-
			projection.DiscountedUnallocatedPremiumIncomeAdjusted-
			projection.DiscountedPremiumAdvisoryFeeAdjusted-
			projection.DiscountedFundAdvisoryFeeAdjusted-
			projection.DiscountedPolicyFeeAdjusted-
			projection.DiscountedFundRiskChargeAdjusted-
			projection.DiscountedSurrenderPenaltyChargeAdjusted-
			projection.DiscountedFundAssetManagementChargesAdjusted-
			projection.DiscountedBsaShareholderChargeAdjusted-
			projection.DiscountedPremiumIncomeAdjusted*(1-variableIndicator)-
			projection.DiscountedClawBackAdjusted, defaultPrecision)

	projection.DiscountedCashOutflow = utils.FloatPrecision(
		projection.DiscountedPremiumNotReceived+
			projection.DiscountedCommission+
			projection.DiscountedRenewalCommission-
			projection.DiscountedClawBack+
			projection.DiscountedDeathOutgo+
			projection.DiscountedAccidentalDeathOutgo+
			projection.DiscountedCashBackOnSurvival+
			projection.DiscountedCashBackOnDeath+
			projection.DiscountedRider+
			projection.DiscountedDisabilityOutgo+
			projection.DiscountedRetrenchmentOutgo+
			projection.DiscountedAnnuityOutgo+
			projection.DiscountedInitialExpenses+
			projection.DiscountedRenewalExpenses+
			projection.DiscountedPremiumAdvisoryCost+
			projection.DiscountedFundAdvisoryCost, defaultPrecision)

	projection.DiscountedCashOutflowExclAcquisition = utils.FloatPrecision(
		projection.DiscountedPremiumNotReceived+
			projection.DiscountedDeathOutgo+
			projection.DiscountedRenewalCommission+
			projection.DiscountedAccidentalDeathOutgo+
			projection.DiscountedCashBackOnSurvival+
			projection.DiscountedCashBackOnDeath+
			projection.DiscountedRider+
			projection.DiscountedDisabilityOutgo+
			projection.DiscountedRetrenchmentOutgo+
			projection.DiscountedAnnuityOutgo+
			projection.DiscountedRenewalExpenses+
			projection.DiscountedMaturityOutgo+
			projection.DiscountedGuaranteeCost, defaultPrecision)

	projection.DiscountedAcquisitionCost = utils.FloatPrecision(
		projection.DiscountedCommission+
			projection.DiscountedInitialExpenses, defaultPrecision)

	projection.DiscountedCashInflow = utils.FloatPrecision(
		projection.DiscountedPremiumIncome*(1-variableIndicator)+
			projection.DiscountedUnallocatedPremiumIncome+
			projection.DiscountedPolicyFee+
			projection.DiscountedPremiumAdvisoryFee+
			projection.DiscountedFundAdvisoryFee+
			projection.DiscountedFundRiskCharge+
			projection.DiscountedFundAssetManagementCharges+
			projection.DiscountedSurrenderPenaltyCharge+
			projection.DiscountedBsaShareholderCharge, defaultPrecision)

	if len(*projections)-1 == index {
		projection.DiscountedProfit = 0
		projection.DiscountedProfitAdjusted = 0
		projection.DiscountedAnnuityFactor = 0

	} else {
		projection.DiscountedProfit = ((*projections)[index+1].Profit + (*projections)[index+1].DiscountedProfit) * discountFactor
		projection.DiscountedProfitAdjusted = ((*projections)[index+1].ProfitAdjusted + (*projections)[index+1].DiscountedProfitAdjusted) * discountFactorAdjusted
		projection.DiscountedAnnuityFactor = projection.DiscountedAnnuityFactor + (*projections)[index+1].DiscountedAnnuityFactor*discountFactor
	}

	//if projection.ProjectionMonth == 0 {
	//	projection.ChangeInReserves = utils.FloatPrecision(projection.Reserves-projection.NetCashFlow, defaultPrecision)
	//	projection.ChangeInReservesAdjusted = utils.FloatPrecision(projection.ReservesAdjusted-projection.NetCashFlowAdjusted, defaultPrecision)
	//
	//} else {
	//	projection.ChangeInReserves = utils.FloatPrecision(futureValues.FutureReserve-projection.Reserves, defaultPrecision)
	//	projection.ChangeInReservesAdjusted = utils.FloatPrecision(futureValues.FutureReserveAdjusted-projection.ReservesAdjusted, defaultPrecision)
	////	projection.InvestmentIncome = utils.FloatPrecision(projection.Reserves*(math.Pow(1.0+forwardRate, 1/12.0)-1), defaultPrecision)
	////	projection.InvestmentIncomeAdjusted = utils.FloatPrecision(projection.ReservesAdjusted*(math.Pow(1.0+forwardRate+(margins.InvestmentMargin/100.0), 1/12.0)-1), defaultPrecision)
	//	//projection.DiscountedPremiumIncome = (projection.)
	//}

	return nil
}

func InterestAccrual(index, valYear, valMonth int, projection *models.Projection,
	prevReserve, prevReserveAdjusted float64,
	run models.RunParameters, mp models.ProductModelPoint,
	params models.ProductParameters, margins models.ProductMargins, productShocks []models.ProductShock) error {
	var forwardRate float64
	var err error
	var resp float64 = 0
	var infl float64 = 0
	var realResp float64 = 0
	var shockedReal float64 = 0

	if index == 0 {
		projection.InvestmentIncome = 0
		return nil
	} else {
		if run.YieldCurveBasis == Current {
			resp, err = GetForwardRateWithError(index, run.YieldcurveYear, run.YieldcurveMonth, params.YieldCurveCode)
			if err != nil {
				return errors.Wrap(err, "failed to get forward rate")
			}

			if run.ShockSettings.RealYieldCurve {
				infl, err = GetInflationFactorWithError(index, run.YieldcurveYear, run.YieldcurveMonth, params.YieldCurveCode)
				if err != nil {
					return errors.Wrap(err, "failed to get inflation factor")
				}
				realResp = (1+resp)/(1+infl) - 1
			}
		} else {
			resp, err = GetLockedinForwardRateWithError(index, mp.LockedInYear, mp.LockedInMonth)
			if err != nil {
				return errors.Wrap(err, "failed to get forward rate")
			}

			if run.ShockSettings.RealYieldCurve {
				infl, err = GetLockedinInflationFactorWithError(index, mp.LockedInYear, mp.LockedInMonth)
				if err != nil {
					return errors.Wrap(err, "failed to get inflation factor")
				}
				realResp = (1+resp)/(1+infl) - 1
			}
		}

		if run.ShockSettings.NominalYieldCurve {
			var shock models.ProductShock
			if projection.ProjectionMonth > 0 && run.ShockSettings.ShockBasis != "" {
				shock, err = GetShock(projection.ProjectionMonth, run.ShockSettings.ShockBasis, len(productShocks))
				if err != nil {
					log.WithField("projection_month", projection.ProjectionMonth).
						WithField("shock_basis", run.ShockSettings.ShockBasis).
						WithField("product_shocks", len(productShocks)).
						Error("failed to get shock")
				}
			}
			forwardRate = math.Max(0, math.Min(1, resp*(1+shock.MultiplicativeYieldCurve)+shock.AdditiveYieldCurve))
		}

		if run.ShockSettings.RealYieldCurve {
			var shock models.ProductShock
			if projection.ProjectionMonth > 0 && run.ShockSettings.ShockBasis != "" {
				shock, err = GetShock(projection.ProjectionMonth, run.ShockSettings.ShockBasis, len(productShocks))
				if err != nil {
					log.WithField("projection_month", projection.ProjectionMonth).WithField("shock_basis", run.ShockSettings.ShockBasis).Error("failed to get shock")

				}
			}
			shockedReal = math.Max(0, math.Min(1, realResp*(1+shock.MultiplicativeYieldCurve)+shock.AdditiveYieldCurve))
			forwardRate = (1+shockedReal)*(1+infl) - 1
		}

		if !run.ShockSettings.NominalYieldCurve && !run.ShockSettings.RealYieldCurve {
			forwardRate = resp
		}

		projection.InvestmentIncome = (prevReserve - (projection.PremiumNotReceived + projection.Commission + projection.RenewalCommission + projection.RenewalExpenses + projection.EPremiumAdvisoryCost + projection.EFundAdvisoryCost - (projection.PremiumIncome - projection.EAllocatedPremiumIncome) - projection.EPolicyFee - projection.EPremiumAdvisoryFee - projection.EFundAdvisoryFee)) * (math.Pow(1.0+forwardRate, 1.0/12.0) - 1.0)
		projection.InvestmentIncomeAdjusted = (prevReserveAdjusted - (projection.PremiumNotReceived + projection.CommissionAdjusted + projection.RenewalCommissionAdjusted + projection.RenewalExpensesAdjusted + projection.EPremiumAdvisoryCostAdjusted + projection.EFundAdvisoryCostAdjusted - (projection.PremiumIncomeAdjusted - projection.EAllocatedPremiumIncomeAdjusted) - projection.EPolicyFeeAdjusted - projection.EPremiumAdvisoryFeeAdjusted - projection.EFundAdvisoryFeeAdjusted)) * (math.Pow(1.0+forwardRate+margins.InvestmentMargin, 1.0/12.0) - 1.0)

	}

	return nil
}
