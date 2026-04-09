package services

import (
	"api/models"
	"api/utils"
	"fmt"
	"github.com/dgraph-io/ristretto"
	"github.com/gammazero/workerpool"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"math"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// var currentYear = time.Now().Year()
var totalRun = 0.0

var PricingCache *ristretto.Cache

var totalConfigs = 0
var configsDone = 0
var totalPoints = 0
var totalPointsDone = 0

var pricingMortalityTableName string
var pricingMortalityAccidentalTableName string
var pricingLapseTableName string
var pricingRetrenchmentTableName string
var pricingDisabilityTableName string
var pricingLapseMonthCount int

var pricingTables = make(map[string]models.PricingProductTableNames)

func RunPricing(pricingRun models.PricingRun, user models.AppUser) error {
	var states []models.ProductTransitionState
	var features models.ProductFeatures
	var margins models.ProductPricingMargins

	pricingRun.RunDate = time.Now()
	pricingRun.User = user.UserName
	pricingRun.Status = "queued"
	pricingRun.UserEmail = user.UserEmail
	DB.Create(&pricingRun)
	DB.Save(&pricingRun)

	for {
		var count int64
		DB.Model(&models.PricingRun{}).Where("status = ?", "in progress").Count(&count)
		fmt.Println("polling for jobs to finish")
		if count == 0 {
			fmt.Println("no jobs in progress")
			break
		}
		time.Sleep(10 * time.Second)
	}

	PricingCache.Clear()

	var mps []models.ProductPricingModelPoint
	var columnNames models.Columnname
	var pricingTableNames models.PricingProductTableNames

	// Saving the pricingRun first?
	totalConfigs = 0
	configsDone = 0
	totalPoints = 0
	totalPointsDone = 0

	runStart := time.Now()
	prod, err := GetProductById(pricingRun.ProductId)
	if err != nil {
		pricingRun.Status = "failed"
		pricingRun.FailureReason = "No pricing product was found for the given id"
		pricingRun.RunTime = time.Since(runStart).Seconds()
		DB.Save(&pricingRun)
		return err
	}
	pricingRun.Status = "in progress"
	pricingRun.RunDate = time.Now()
	DB.Save(&pricingRun)

	pricingModelPointTableName := strings.ToLower(prod.ProductCode) + "_pricing_modelpoints"

	pricingMortalityTableName = getPricingRatingTable(prod.ProductCode, "Death")
	pricingTableNames.MortalityTableName = strings.ToLower(prod.ProductCode) + "_" + "pricing_" + strings.ToLower(pricingMortalityTableName)
	if pricingTableNames.MortalityTableName != "" {
		pricingTableNames.MortalityColumnName = GetColumnName(strings.ToLower(pricingTableNames.MortalityTableName))
		columnNames.MortalityColumnName = GetColumnName(strings.ToLower(prod.ProductCode) + "_" + "pricing_" + strings.ToLower(pricingMortalityTableName))
	}

	pricingMortalityAccidentalTableName = getPricingRatingTable(prod.ProductCode, "Accidental Death")
	pricingTableNames.MortalityAccidentalTableName = strings.ToLower(prod.ProductCode) + "_" + "pricing_" + strings.ToLower(pricingMortalityAccidentalTableName)
	if pricingTableNames.MortalityAccidentalTableName != "" {
		pricingTableNames.MortalityAccidentalColumnName = GetColumnName(strings.ToLower(pricingTableNames.MortalityAccidentalTableName))
		columnNames.AccidentalMortalityColumnName = GetColumnName(strings.ToLower(prod.ProductCode) + "_" + "pricing_" + strings.ToLower(pricingMortalityAccidentalTableName))
	}

	pricingLapseTableName = getPricingRatingTable(prod.ProductCode, "Lapse")
	pricingTableNames.LapseTableName = strings.ToLower(prod.ProductCode) + "_" + "pricing_" + strings.ToLower(pricingLapseTableName)
	if pricingTableNames.LapseTableName != "" {
		pricingTableNames.LapseColumnName = GetColumnName(strings.ToLower(pricingTableNames.LapseTableName))
		columnNames.LapseColumnName = GetColumnName(strings.ToLower(prod.ProductCode) + "_" + "pricing_" + strings.ToLower(pricingLapseTableName))
	}

	pricingRetrenchmentTableName = getPricingRatingTable(prod.ProductCode, "Retrenchment")
	pricingTableNames.RetrenchmentTableName = strings.ToLower(prod.ProductCode) + "_" + "pricing_" + strings.ToLower(pricingRetrenchmentTableName)
	if pricingTableNames.RetrenchmentTableName != "" {
		pricingTableNames.RetrenchmentColumnName = GetColumnName(strings.ToLower(pricingTableNames.RetrenchmentTableName))
		columnNames.RetrenchmentColumnName = GetColumnName(strings.ToLower(prod.ProductCode) + "_" + "pricing_" + strings.ToLower(pricingRetrenchmentTableName))
	}

	pricingDisabilityTableName = getPricingRatingTable(prod.ProductCode, "Permanent Disability")
	pricingTableNames.DisabilityTableName = strings.ToLower(prod.ProductCode) + "_" + "pricing_" + strings.ToLower(pricingDisabilityTableName)
	if pricingTableNames.DisabilityTableName != "" {
		pricingTableNames.DisabilityColumnName = GetColumnName(strings.ToLower(pricingTableNames.RetrenchmentTableName))
		columnNames.DisabilityColumnName = GetColumnName(strings.ToLower(prod.ProductCode) + "_" + "pricing_" + strings.ToLower(pricingDisabilityTableName))
	}

	pricingTableNames.RetrenchmentTableRowCount = getPricingRetrenchmentCount(pricingTableNames.RetrenchmentTableName)

	DB.Where("product_id = ?", prod.ID).Find(&states)
	DB.Where("product_code=?", prod.ProductCode).First(&features)
	DB.Where("product_code=?", prod.ProductCode).First(&margins)

	loadPricingMortalityRates(prod.ProductCode)
	loadPricingAccidentalMortalityRates(prod.ProductCode)
	//loadPricingInflationFactor() //deliberately hard coded to avoid instantiating currentYear
	//loadPricingForwardRate()

	var pricingParams models.PricingParameter
	DB.Where("product_code=?", prod.ProductCode).First(&pricingParams)

	pricingTableNames.LapseMarginMonthCount = getPricingLapseMarginCount(prod.ProductCode, pricingParams)
	pricingTableNames.LapseMonthCount = getPricingLapseCount(pricingTableNames.LapseTableName)
	pricingTables[prod.ProductCode] = pricingTableNames

	// Check for valid yield curve data.
	var ycd models.PricingYieldCurve
	err = DB.Where("yield_curve_code=? ", pricingParams.YieldCurveCode).First(&ycd).Error
	if err != nil && ycd.YieldCurveCode == "" {
		pricingRun.Status = "failed"
		pricingRun.FailureReason = "Yield Curve Code, " + pricingParams.YieldCurveCode + " not found in the available yield curve data"
		pricingRun.RunTime = time.Since(runStart).Seconds()
		DB.Save(&pricingRun)
		return err
	}

	// check if values in pricing yield curve are equal to or greater than 1. Specifically nominal rate and inflation rate
	var ycds []models.PricingYieldCurve
	DB.Where("yield_curve_code=? ", pricingParams.YieldCurveCode).Limit(20).Find(&ycds)
	for _, ycd := range ycds {
		if ycd.NominalRate >= 1 || ycd.Inflation >= 1 {
			pricingRun.Status = "failed"
			pricingRun.FailureReason = "yield curve data must be in rates and not percentages. 0.07 as opposed to 7% for example"
			pricingRun.RunTime = time.Since(runStart).Seconds()
			DB.Save(&pricingRun)
			return errors.New("yield curve data must be in rates and not percentages. 0.07 as opposed to 7% for example")
		}
	}

	var multipliers models.ProductPricingAccidentalBenefitMultiplier
	if features.AccidentalDeathBenefit {
		DB.Where("product_code=?", prod.ProductCode).First(&multipliers)
	} else {
		multipliers.MainMember = 1
		multipliers.Spouse = 1
		multipliers.Child = 1
	}

	for _, cfg := range pricingRun.PricingConfig {
		var count int64
		DB.Table(pricingModelPointTableName).Where("pricing_mp_version =?", cfg.MpVersion).Count(&count)
		totalPoints += int(count)
	}

	fmt.Println(runtime.NumCPU())
	wp := workerpool.New(runtime.NumCPU() * 10)
	//Cycle through the pricing configs and run mps for each config.
	totalConfigs = len(pricingRun.PricingConfig)
	for _, config := range pricingRun.PricingConfig {
		config := config
		wp.Submit(func() {
			runScenario(config, mps,
				pricingParams, multipliers,
				prod, pricingRun.ID, columnNames,
				pricingTableNames, pricingRun.RunGoalSeek, pricingRun,
				features, states, margins,
			)

		})
		//if config.Description == "MM"{
		//	runScenario(config, mps, pricingParams, prod, pricingRun.ID)
		//}

	}
	wp.StopWait()

	fmt.Println("Pricing Run Duration (seconds): ", time.Since(runStart).Seconds(), " seconds")
	fmt.Println("Pricing Run Duration (minutes: ", time.Since(runStart).Minutes(), " minutes")
	return nil
}

func runScenario(config models.PricingConfig,
	mps []models.ProductPricingModelPoint,
	pricingParams models.PricingParameter, multipliers models.ProductPricingAccidentalBenefitMultiplier,
	prod models.Product, runId int, columnNames models.Columnname,
	pricingTableNames models.PricingProductTableNames, runGoalSeek bool, pricingRun models.PricingRun,
	features models.ProductFeatures,
	states []models.ProductTransitionState,
	margins models.ProductPricingMargins) {
	scenarioStart := time.Now()
	profitability := models.Profitability{}
	profitability.ScenarioDescription = config.Description
	profitability.ScenarioID = config.ID
	profitability.RunID = runId
	aggpricingProjs := make(map[string]models.AggregatedPricingPoint)
	var productPricingShocks []models.ProductPricingShock
	var productPricingParameters models.ProductPricingParameters

	pricingModelPointTableName := strings.ToLower(prod.ProductCode) + "_pricing_modelpoints"

	DB.Where("product_code=? and basis=?", prod.ProductCode, config.ParameterBasis).First(&productPricingParameters)
	DB.Where("shock_basis = ?", config.ShockBasis).Find(&productPricingShocks)
	DB.Where("product_code=? and basis=?", prod.ProductCode, config.ParameterBasis).First(&pricingParams)
	productPricingShockCount := len(productPricingShocks)

	//Missing variables/arguments
	if pricingRun.RunSingle {
		//err := DB.Table(pricingModelPointTableName).Limit(1).Find(&mps).Error
		err := DB.Table(pricingModelPointTableName).Where("pricing_mp_version =?", config.MpVersion).Limit(1).Find(&mps).Error
		if err != nil {
			fmt.Println(err)
			//return err
		}
	} else {
		//err = DB.Table(pricingModelPointTableName).Find(&mps).Error
		err := DB.Table(pricingModelPointTableName).Where("pricing_mp_version =?", config.MpVersion).Find(&mps).Error
		if err != nil {
			fmt.Println(err)
			//return err
		}
	}

	//if len(mps) == 0 {
	//	return errors.New("No records were found satisfying the given criteria")
	//}

	for k, mp := range mps {
		//store calc premium, discounted cashflows at month zero, sensitivityid
		var variance = 1000.0
		var technicalPremiumAfterProfitMargin = 0.0
		var technicalPremiumRateAfterProfitMargin = 0.0
		var calculatedMinimumPremium = 0.0
		var childMinimumCalculatedPremium = 0.0
		var childFinalCalculatedPremium = 0.0
		var previousAnnualPremium = mp.AnnualPremium
		var previousPremiumRate = mp.PremiumRate
		for i := 1; i <= pricingParams.GoalSeekMaxIterations; i++ {
			if i == 1 && runGoalSeek {
				if mp.AnnualPremium == 0 {
					mp.AnnualPremium = 1000
					previousAnnualPremium = mp.AnnualPremium
				}
				if mp.PremiumRate == 0 && features.CreditLife {
					mp.PremiumRate = 3
					previousPremiumRate = mp.PremiumRate
				}
			}
			if math.Abs(variance) > 0.5 {
				if i > 2 {
					if variance > 0 {
						mp.AnnualPremium = previousAnnualPremium * (1 + pricingParams.GoalSeekStep)
						mp.PremiumRate = previousPremiumRate * (1 + pricingParams.GoalSeekStep)
					} else {
						mp.AnnualPremium = previousAnnualPremium / (1 + pricingParams.GoalSeekStep)
						mp.PremiumRate = previousPremiumRate / (1 + pricingParams.GoalSeekStep)
					}
				} else {
					mp.AnnualPremium = previousAnnualPremium
					mp.PremiumRate = previousPremiumRate
				}

				RunPricingPerModelPoint(i, mp, prod, config,
					pricingParams, multipliers, &variance,
					&technicalPremiumAfterProfitMargin, &technicalPremiumRateAfterProfitMargin,
					&calculatedMinimumPremium, &childMinimumCalculatedPremium,
					&childFinalCalculatedPremium, &previousAnnualPremium,
					&previousPremiumRate, &profitability,
					columnNames, pricingTableNames, runGoalSeek, pricingRun.ProfitSignature,
					features, productPricingParameters, states, margins, &aggpricingProjs, productPricingShockCount, config.ShockBasis, k)
			} else {
				calculatedMinimumPremium = mp.AnnualPremium
				break
			}
		}

		totalPointsDone += 1
		percentagePointsDone := float64(totalPointsDone) / float64(totalPoints) * 100
		updateRunProgress(percentagePointsDone, runId, 0)
	}
	profitability.TotalValue = profitability.WeightedDiscountedRisk + profitability.WeightedDiscountedRider + profitability.WeightedDiscountedEducator + profitability.WeightedDiscountedProfit + profitability.WeightedDiscountedPremiumNotReceived + profitability.WeightedDiscountedExpenses + profitability.WeightedDiscountedCommission + profitability.WeightedDiscountedCashBackOnSurvival + profitability.WeightedDiscountedCashBackOnDeath //+ profitability.WeightedDiscountedChangeInReserve - profitability.WeightedDiscountedInvestmentIncome
	DB.Create(&profitability)
	var aggregatedPricingProjections []models.AggregatedPricingPoint
	for _, appProj := range aggpricingProjs {
		aggregatedPricingProjections = append(aggregatedPricingProjections, appProj)
	}
	DB.CreateInBatches(&aggregatedPricingProjections, 100)

	configsDone += 1
	fmt.Println("configs done: ", configsDone)
	fmt.Println("total configs: ", totalConfigs)
	percentageDone := float64(configsDone) / float64(totalConfigs) * 100
	fmt.Println(percentageDone, "%")
	updateRunProgress(percentageDone, runId, time.Since(scenarioStart).Seconds())
	totalRun += time.Since(scenarioStart).Minutes()
	fmt.Println("Scenario Run: ", time.Since(scenarioStart).Minutes(), " minutes")
}

func updateRunProgress(percentageDone float64, runId int, runtime float64) {
	var pricingRun models.PricingRun
	DB.Where("id = ?", runId).First(&pricingRun)
	status := "in progress"
	if percentageDone > 99 {
		status = "complete"
	}
	err := DB.Model(&pricingRun).Where("id = ?", runId).Updates(map[string]interface{}{"status": status, "progress": percentageDone, "run_time": runtime}).Error
	if err != nil {
		fmt.Println(err)
	}
}

func RunPricingPerModelPoint(
	gskIndex int,
	mp models.ProductPricingModelPoint,
	prod models.Product,
	pricingConfig models.PricingConfig,
	pricingParams models.PricingParameter,
	multipliers models.ProductPricingAccidentalBenefitMultiplier,
	variance *float64,
	technicalPremiumAfterProfitmargin *float64,
	technicalPremiumRateAfterProfitMargin *float64,
	calculatedMinimumPremium *float64,
	childMinimumCalculatedPremium *float64,
	childFinalCalculatedPremium *float64,
	previousAnnualPremium *float64,
	previousPremiumRate *float64,
	pf *models.Profitability,
	columnnames models.Columnname,
	pricingTableNames models.PricingProductTableNames,
	runGoalSeek bool, profitSignature bool,
	features models.ProductFeatures,
	parameters models.ProductPricingParameters,
	states []models.ProductTransitionState,
	margins models.ProductPricingMargins,
	aggpricingProjs *map[string]models.AggregatedPricingPoint, productPricingShockCount int, pricingShockBasis string,
	k int,
) {
	fmt.Println("Running for: ", mp.PolicyNumber)
	var err error
	var pricingPoints []models.PricingPoint
	var p models.PricingPoint
	var startingInitialPolicy float64 = 0
	var startingInitialPolicyAdjusted float64 = 0

	if features.FuneralCover {
		//Validate the age range of the life. No use running if the age falls outside the range.
		if mp.MainMemberAgeAtEntry < 18 || mp.MainMemberAgeAtEntry > 120 {
			fmt.Println("The Main member age is outside the required range. Exiting this run...", mp.MainMemberAgeAtEntry)
			RecordsSkipped += 1
			return
		}

		if mp.AgeAtEntry < 18 || mp.AgeAtEntry > 120 {
			fmt.Println("The age for this life is outside the required range. Exiting this run...")
			RecordsSkipped += 1
			return
		}
	}

	if mp.MemberType != "MM" {
		RecordsSkipped += 1
		return
	}
	calculateProductPricingConfig(&parameters, mp, features, multipliers)
	var index = 0
	projectionRange := 1440 // Should be 1440...

	for ; index <= projectionRange; index++ {
		var pp models.PricingPoint
		var shock models.ProductPricingShock
		fmt.Println("running: ", index)
		pp.ProductCode = prod.ProductCode
		pp.JobProductID = prod.ID
		pp.SensitivityID = pricingConfig.ID
		pp.RunId = pricingConfig.PricingRunID
		pp.ScenarioId = pricingConfig.ID
		pricingTimeAndAge(&pp, mp, index, features)

		//End the loop. No need to continue
		if pp.AgeNextBirthday > 121 {
			break
		}

		//PricingBenefitInForce(&pp, mp, features, parameters)
		PricingAccidentProportion(&pp, mp, states, parameters, columnnames.AccidentalMortalityColumnName, pricingTableNames.MortalityAccidentalTableName)
		PricingInflationFactor(index, &pp, p, margins, pricingParams)
		PricingPremiumEscalation(index, &pp, p, mp, features, pricingConfig)
		PricingLapseMargin(&pp, pricingParams, pricingTables[mp.ProductCode].LapseMarginMonthCount)
		if pricingShockBasis != "N/A" {
			shock, _ = GetPricingShock(index, pricingShockBasis, productPricingShockCount)
		}

		//not required for pricing
		//PricingPremiumWaiverOnFactor(&pp, mp, parameters)
		//PricingPaidUpOnFactor(&pp, mp, parameters)
		//*
		PricingMainMemberMortalityRate(&pp, mp, margins, features, states, parameters, columnnames.MortalityColumnName, pricingTableNames.MortalityTableName, shock, pricingShockBasis)
		PricingBaseLapse(&pp, mp, features, prod.ProductCode, parameters, pricingParams, columnnames.LapseColumnName, pricingTableNames.LapseTableName, shock, pricingShockBasis)
		PricingContractingPartyAlivePortion(&pp, p, mp, features, parameters) //revisit
		//PricingContractingPartyPolicyLapse(&pp, p, mp, features, parameters) // revisit
		PricingBaseMortalityRate(&pp, mp, margins, states, parameters, columnnames.MortalityColumnName, pricingTableNames.MortalityTableName, shock, pricingShockBasis)

		if pricingConfig.Retrenchment {
			//utils.StatesContains(&states, Retrenchment)
			if pp.AgeNextBirthday <= 65 {
				PricingBaseRetrenchmentRate(&pp, features, mp, margins, states, parameters, columnnames.RetrenchmentColumnName, pricingTableNames.RetrenchmentTableName, shock, pricingShockBasis)
			}
		}

		if pricingConfig.PermDisability {
			//utils.StatesContains(&states, PermanentDisability)
			if pp.AgeNextBirthday <= 65 {
				PricingBaseDisabilityIncrement(&pp, mp, features, margins, states, parameters, columnnames.DisabilityColumnName, pricingTableNames.DisabilityTableName, shock, pricingShockBasis)
			}
		}

		PricingMainMemberMortalityRateByMonth(&pp, features, parameters)
		PricingIndependentLapseMonthly(&pp, p, pricingParams, parameters)
		PricingIndependentRetrenchmentMonthly(&pp, parameters)
		PricingIndependentDisabilityMonthly(&pp, parameters)

		PricingMonthlyDependentMortality(&pp, parameters)
		PricingMonthlyDependentLapse(&pp, parameters)
		PricingMonthlyDependentRetrenchment(&pp, parameters)
		PricingMonthlyDependentDisability(&pp, parameters)

		// should only run if there is maturity benefit
		PricingNumberOfMaturities(&pp, p, mp, features, parameters)
		PricingNaturalDeathsInForce(&pp, p, parameters)
		PricingNumberOfDeathsAccident(&pp, p, parameters)

		//main member's deaths calculation is not required for premiumwaiver on death since the premiums are waived upon death of the main member
		//PricingNaturalDeathsPremiumWaiver(&pp, p, parameters)
		//PricingNaturalDeathsTemporaryWaivers(&pp, p, parameters)
		//PricingAccidentDeathsPremiumWaiver(&pp, p, parameters)
		//PricingAccidentDeathsTemporaryPremiumWaiver(&pp, p, parameters)

		PricingNumberOfLapses(&pp, p, parameters, mp, pricingConfig)
		PricingNumberOfDisabilities(&pp, p, parameters)
		PricingNumberOfRetrenchments(&pp, p, parameters)

		//revisit the block below
		//PricingIncrementalNaturalDeaths(&pp, p, features,parameters)

		PricingTotalIncrementalLapses(&pp, p, parameters)
		PricingNaturalDeathsPaidUp(&pp, p, parameters)
		PricingAccidentDeathsPaidUp(&pp, p, parameters)
		PricingTotalIncrementalNaturalDeaths(&pp, p, parameters)
		PricingTotalIncrementalAccidentalDeaths(&pp, mp, p, parameters)
		PricingTotalIncrementalDisabilities(&pp, p, parameters)
		PricingTotalIncrementalRetrenchments(&pp, p, parameters)

		PricingIncrementalPaidUp(&pp, p, pricingConfig, parameters, features)
		PricingIncrementalPremiumWaivers(&pp, pricingConfig, parameters, features)

		if pricingConfig.SpouseIndicator {
			PricingSpouseAgeNextBirthday(&pp, mp, features)
			PricingSpouseMortalityRate(&pp, mp, pricingConfig, margins, columnnames.MortalityColumnName, pricingTableNames.MortalityTableName)
			PricingBaseSpouseIndependentLapse(&pp, pricingConfig)
			PricingIndependentSpouseMortalityRateByMonth(&pp)
			PricingIndependentSpouseLapseMonthly(&pp)
			PricingMonthlySpouseDependentMortality(&pp)
			PricingMonthlySpouseDependentLapse(&pp)
			PricingSpouseNumberPolicies(&pp, p, pricingConfig)
		}

		if pricingConfig.ChildIndicator {
			PricingChildMortalityRate(&pp, mp, pricingConfig, pricingParams, margins, columnnames.MortalityColumnName, pricingTableNames.MortalityTableName)
			PricingChildAccidentalProportion(&pp, mp, pricingConfig, pricingParams, columnnames.AccidentalMortalityColumnName, pricingTableNames.MortalityAccidentalTableName)
			PricingIndependentChildMonthlyMortalityRate(&pp, pricingConfig)
			PricingMonthlyChildDependentMortality(&pp)
			PricingChildIndependentLapse(&pp, pricingConfig, parameters)
			PricingIndependentChildLapseMonthly(&pp, pricingConfig, parameters)
			//PricingChildNumberPolicies(&pp, p, pricingConfig, pricingParams)
			PricingMonthlyChildDependentLapse(&pp)
			PricingChildNumberPolicies(&pp, p, pricingConfig, pricingParams, parameters)
		}

		PricingNumberPaidup(&pp, p, parameters, pricingConfig, pricingParams)

		//since premium waiver applies on the death of the main member, calculation of mainmember deaths is not required
		//PricingAccidentDeathsPremiumWaiver(&pp, p, parameters)
		//PricingAccidentDeathsTemporaryPremiumWaiver(&pp, p, parameters)

		PricingInitialPolicy(&pp, mp, &startingInitialPolicy, &startingInitialPolicyAdjusted, pricingConfig, parameters)
		//*

		// Stuff all spouse related params here
		if pricingConfig.SpouseIndicator {
			PricingSpouseNumberOfPaidUps(&pp, p, pricingConfig)
			PricingSpouseNumberOfPremiumWaivers(&pp, p, pricingConfig)
			PricingTotalSpouseIncrementalNaturalDeaths(&pp, p, mp, columnnames.AccidentalMortalityColumnName, pricingTableNames.MortalityAccidentalTableName)
			PricingTotalSpouseIncrementalAccidentalDeaths(&pp, p, mp, columnnames.AccidentalMortalityColumnName, pricingTableNames.MortalityAccidentalTableName)
			PricingSpouseNumberOfPaidUpNaturalDeaths(&pp, p, mp, parameters, columnnames.AccidentalMortalityColumnName, pricingTableNames.MortalityAccidentalTableName)
			PricingSpouseNumberOfPaidUpAccidentalDeaths(&pp, p, mp, parameters, columnnames.AccidentalMortalityColumnName, pricingTableNames.MortalityAccidentalTableName)
			PricingSpouseNumberOfPremiumWaiverNaturalDeaths(&pp, p, mp, parameters, columnnames.AccidentalMortalityColumnName, pricingTableNames.MortalityAccidentalTableName)
			PricingSpouseNumberOfPremiumWaiverAccidentalDeaths(&pp, p, mp, parameters, columnnames.AccidentalMortalityColumnName, pricingTableNames.MortalityAccidentalTableName)
		}

		// Stuff all child related params here
		if pricingConfig.ChildIndicator {
			PricingChildNumberOfPaidUps(&pp, p, pricingConfig, pricingParams, parameters)
			PricingChildNumberOfPremiumWaivers(&pp, p, pricingConfig, pricingParams, parameters)
			PricingTotalChildIncrementalNaturalDeaths(&pp, p, pricingConfig, pricingParams, parameters)
			PricingTotalChildIncrementalAccidentalDeaths(&pp, p, pricingParams, parameters)
			PricingChildNumberOfPaidUpNaturalDeaths(&pp, p, parameters)
			PricingChildNumberOfPaidUpAccidentalDeaths(&pp, p, parameters)
			PricingChildNumberOfPremiumWaiverNaturalDeaths(&pp, p, parameters)
			PricingChildNumberOfPremiumWaiverAccidentalDeaths(&pp, p, parameters)
		}

		//***
		//Only runs if credit_life or retrenchment benefit structure is chosen
		if features.CreditLife || features.RetrenchmentBenefit {
			PricingCalculatedInstalment(&pp, &p, mp, parameters)
		}

		//***
		//Only runs if credit_life benefit structure is chosen
		if features.CreditLife {
			PricingOutstandingSumAssured(&pp, &p, mp, parameters, features)
		}

		PricingSumAssured(&pp, mp, parameters, pricingParams, features, pricingConfig)
		PricingAdditionalSumAssured(&pp, mp, parameters, pricingParams) //revisit
		PricingPremium(&pp, mp, features, parameters, pricingConfig, *childFinalCalculatedPremium)
		if pricingConfig.ChildIndicator {
			PricingChildSumAssured(&pp, mp, parameters, pricingConfig, pricingParams)
			PricingChildAdditionalSumAssured(&pp, mp, parameters, pricingConfig, pricingParams)
		}

		PricingPremiumIncome(&pp, p, mp, parameters, pricingConfig, pricingParams)
		PricingPremiumsNotReceived(&pp, mp, parameters, features)
		PricingCommissions(mp, &pp, p, parameters, pricingParams)
		PricingClawback(mp, &pp, parameters, pricingParams)
		PricingDeathOutgo(&pp, features, mp, parameters, pricingConfig)

		PricingAccidentalDeathOutgo(&pp, parameters, multipliers, features, pricingConfig)
		if pricingConfig.ChildIndicator && pricingConfig.EducatorIndicator {
			PricingEducator(&pp, p, pricingParams, parameters, mp)
		}

		PricingDisabilityOutgo(&pp, mp, parameters, features)

		if pricingConfig.Retrenchment {
			PricingRetrenchmentOutgo(&pp, mp, parameters, features)
		}

		if pricingConfig.CashBackOnSurvival || pricingConfig.CashBack {
			PricingCashBackOnSurvival(&pp, mp, parameters, features, pricingConfig)
		}

		if pricingConfig.CashBackOnDeath {
			PricingCashBackOnDeath(&pp, mp, parameters, pricingConfig, pricingParams, features)
		}

		if pricingConfig.Funeral {
			PricingFuneralRider(&pp, mp, parameters, pricingConfig, pricingParams, features)
		}

		PricingRider(&pp, features, mp, pricingConfig, parameters)
		PricingExpenses(&pp, p, mp, parameters, margins, features, shock, pricingShockBasis)
		PricingNetCashFlow(&pp, p, mp, parameters) //**
		pricingPoints = append(pricingPoints, pp)
		//fmt.Println("PricingPoints:", len(pricingPoints))
		p = pp
	}

	// Discounted Values
	//projPeriod := len(pricingPoints) - 1
	for i := len(pricingPoints) - 1; i >= 0; i-- {
		var shock models.ProductPricingShock
		if pricingShockBasis != "N/A" {
			shock, _ = GetPricingShock(i, pricingShockBasis, productPricingShockCount)
		}
		pricingDiscountedValues(i, &pricingPoints[i], mp, margins, pricingParams, &pricingPoints, features, shock, pricingShockBasis, parameters) // future values discounted hence [i+1]
	}
	netCashFlow := 0.0
	reserveChangeA := 0.0
	reserveChange := 0.0

	prevMonth := 0

	for i := range pricingPoints {
		var shock models.ProductPricingShock
		if pricingShockBasis != "N/A" {
			shock, _ = GetPricingShock(i, pricingShockBasis, productPricingShockCount)
		}
		if i == 0 {
			//Seeding the variables
			netCashFlow = pricingPoints[i].NetCashFlow
			reserveChange = pricingPoints[i].Reserves
			reserveChangeA = pricingPoints[i].ReservesAdjusted
			prevMonth = pricingPoints[i].ProjectionMonth
			pricingPoints[i].ProfRiskAdjustment = math.Abs(pricingPoints[i].Reserves * pricingParams.RiskAdjustmentProp)
			pricingPoints[i].ProfCSM = math.Max(-(pricingPoints[i].Reserves + pricingPoints[i].ProfRiskAdjustment), 0)
			pricingPoints[i].ProfLossComponent = math.Max(pricingPoints[i].Reserves+pricingPoints[i].ProfRiskAdjustment, 0)
			if pricingPoints[i].ValuationTimeMonth == 1 {
				pricingPoints[i].ChangeInReserves = pricingPoints[i].Reserves
				pricingPoints[i].ChangeInReservesAdjusted = pricingPoints[i].ReservesAdjusted
			}
			continue
		} else {

			netCashFlow = pricingPoints[i].NetCashFlow

			pricingPoints[i].ChangeInReserves = pricingPoints[i].Reserves - pricingPoints[i-1].Reserves
			pricingPoints[i].ChangeInReservesAdjusted = pricingPoints[i].ReservesAdjusted - pricingPoints[i-1].ReservesAdjusted

			reserveChange = pricingPoints[i].ChangeInReserves
			reserveChangeA = pricingPoints[i].ChangeInReservesAdjusted

			if prevMonth <= parameters.CalculatedTerm {
				pricingInterestAccrual(i, &pricingPoints[i], &pricingPoints[i-1], margins, pricingParams, shock, pricingShockBasis)
				forwardRate := getPricingForwardRate(i, pricingParams.YieldCurveCode, shock, pricingShockBasis)
				accrualFactor := math.Pow(1.0+forwardRate, 1.0/12.0)
				accrualFactorAdjusted := math.Pow(1.0+forwardRate+margins.InvestmentMargin, 1.0/12.0)
				if pricingPoints[i].ValuationTimeMonth == 1 {
					pricingPoints[i].Profit = utils.FloatPrecision((-pricingPoints[i-1].NetCashFlow-pricingPoints[i-1].Reserves)*accrualFactor+
						-netCashFlow-reserveChange+pricingPoints[i].InvestmentIncome, AccountingPrecision)

					pricingPoints[i].ProfitAdjusted = utils.FloatPrecision((-pricingPoints[i-1].NetCashFlowAdjusted-pricingPoints[i-1].ReservesAdjusted)*accrualFactorAdjusted+
						-netCashFlow-reserveChangeA+pricingPoints[i].InvestmentIncomeAdjusted, AccountingPrecision)
				} else {
					pricingPoints[i].Profit = utils.FloatPrecision(-netCashFlow-reserveChange+pricingPoints[i].InvestmentIncome, AccountingPrecision)
					pricingPoints[i].ProfitAdjusted = utils.FloatPrecision(-netCashFlow-reserveChangeA+pricingPoints[i].InvestmentIncomeAdjusted, AccountingPrecision)
				}

				//pricingPoints[i].Profit = utils.FloatPrecision(-netCashFlow-reserveChange+investIncome, 6)
				//pricingPoints[i].ProfitAdjusted = utils.FloatPrecision(-netCashFlow-reserveChangeA+investIncomeA, 6)
			}
		}

		//DB.Save(&projections[i])
		prevMonth = pricingPoints[i].ProjectionMonth
	}

	for i := len(pricingPoints) - 1; i >= 0; i-- {
		var shock models.ProductPricingShock
		if pricingShockBasis != "N/A" {
			shock, _ = GetPricingShock(i, pricingShockBasis, productPricingShockCount)
		}
		if i == len(pricingPoints)-1 {
			pricingPoints[i].DiscountedProfit = 0
			pricingPoints[i].DiscountedProfitAdjusted = 0
			pricingPoints[i].ProfDiscountedChangeInReserve = 0
			pricingPoints[i].ProfDiscountedInvestmentIncome = 0
		} else {
			forwardRate := getPricingForwardRate(i+1, pricingParams.YieldCurveCode, shock, pricingShockBasis)
			discountFactor := math.Pow(1.0+forwardRate, -1.0/12.0)
			discountFactorAdjusted := math.Pow(1.0+forwardRate+margins.InvestmentMargin, -1.0/12.0)
			pricingPoints[i].DiscountedInvestmentIncome = (pricingPoints[i+1].InvestmentIncome + pricingPoints[i+1].DiscountedInvestmentIncome) * discountFactor
			pricingPoints[i].DiscountedInvestmentIncomeAdjusted = 0 //(pricingPoints[i+1].InvestmentIncomeAdjusted + pricingPoints[i+1].DiscountedInvestmentIncomeAdjusted) * discountFactorAdjusted
			pricingPoints[i].DiscountedProfit = (pricingPoints[i+1].Profit + pricingPoints[i+1].DiscountedProfit) * discountFactor
			pricingPoints[i].DiscountedProfitAdjusted = (pricingPoints[i+1].ProfitAdjusted + pricingPoints[i+1].DiscountedProfitAdjusted) * discountFactorAdjusted
			pricingPoints[i].ProfDiscountedInvestmentIncome = (pricingPoints[i+1].InvestmentIncome + pricingPoints[i+1].ProfDiscountedInvestmentIncome) * discountFactor
			if i == 0 {
				pricingPoints[i].ProfDiscountedChangeInReserve = pricingPoints[i].Reserves + (pricingPoints[i+1].ChangeInReserves+pricingPoints[i+1].ProfDiscountedChangeInReserve)*discountFactor
			}
			if i > 0 {
				pricingPoints[i].ProfDiscountedChangeInReserve = (pricingPoints[i+1].ChangeInReserves + pricingPoints[i+1].ProfDiscountedChangeInReserve) * discountFactor
			}

			if pricingConfig.ChildIndicator {
				if pricingPoints[i+1].ChildNumberPolicies == 0 || pricingPoints[i].ChildNumberPolicies == 0 {
					pricingPoints[i].ChildAnnuityFactor = 0
				} else {
					pricingPoints[i].ChildAnnuityFactor = pricingPoints[i].ChildNumberPolicies*pricingPoints[i+1].PremiumEscalation/pricingPoints[0].ChildNumberPolicies + pricingPoints[i+1].ChildAnnuityFactor*discountFactor
				}
			}
		}
	}

	firstPP := pricingPoints[0]
	var profitMargin float64
	if firstPP.DiscountedPremiumIncome == 0 {
		profitMargin = 0
	} else {
		profitMargin = firstPP.DiscountedProfit / firstPP.DiscountedPremiumIncome
	}

	if runGoalSeek {
		*variance = (getProfitMargin(mp.AgeAtEntry, int(firstPP.SumAssured), mp.ProductCode) - profitMargin) * 100.0
	} else {
		*variance = 0
	}

	//checkValue := firstPP.DiscountedPremiumIncome / mp.AnnualPremium
	//controlValue := firstPP.AnnuityFactor / 12
	//checkAnnualPremiumValue := controlValue * mp.AnnualPremium
	//fmt.Println(checkValue, controlValue, checkAnnualPremiumValue)
	if gskIndex == 1 && runGoalSeek {
		var temp1 float64
		if pricingParams.InitialCommissionPercentage2 > 0 {
			if pricingParams.InitialCommissionPercentage1 > 0 {
				temp1 = 1 + pricingParams.InitialCommissionPercentage2/pricingParams.InitialCommissionPercentage1
			}
		} else {
			temp1 = 1
		}

		if mp.CommissionType == Initial {

			if pricingConfig.CashBackOnSurvival {
				*calculatedMinimumPremium = (firstPP.DiscountedDeathOutgo + firstPP.ChildDiscountedDeathOutgo + firstPP.SpouseDiscountedDeathOutgo + firstPP.DiscountedAccidentalDeathOutgo + firstPP.DiscountedDisabilityOutgo + firstPP.DiscountedRetrenchmentOutgo + firstPP.ChildDiscountedAccidentalDeathOutgo + firstPP.SpouseDiscountedAccidentalDeathOutgo + firstPP.DiscountedEducator + firstPP.DiscountedRider + firstPP.DiscountedExpenses + firstPP.DiscountedRiderFuneral + pricingParams.InitialCommissionRand) / ((firstPP.AnnuityFactor / 12) - (pricingParams.InitialCommissionPercentage1*temp1 + parameters.CashbackOnSurvivalRatio))
				if pricingConfig.ChildIndicator {
					*childMinimumCalculatedPremium = (firstPP.ChildDiscountedDeathOutgo + firstPP.ChildDiscountedAccidentalDeathOutgo) / ((firstPP.ChildAnnuityFactor / 12) - (pricingParams.InitialCommissionPercentage1*temp1 + parameters.CashbackOnSurvivalRatio))
				}
			} else {
				*calculatedMinimumPremium = (firstPP.DiscountedDeathOutgo + firstPP.ChildDiscountedDeathOutgo + firstPP.SpouseDiscountedDeathOutgo + firstPP.DiscountedAccidentalDeathOutgo + firstPP.ChildDiscountedAccidentalDeathOutgo + firstPP.DiscountedDisabilityOutgo + firstPP.DiscountedRetrenchmentOutgo + firstPP.SpouseDiscountedAccidentalDeathOutgo + firstPP.DiscountedEducator + firstPP.DiscountedRider + firstPP.DiscountedExpenses + firstPP.DiscountedRiderFuneral + pricingParams.InitialCommissionRand) / ((firstPP.AnnuityFactor / 12) - (pricingParams.InitialCommissionPercentage1 * temp1))
				if pricingConfig.ChildIndicator {
					*childMinimumCalculatedPremium = (firstPP.ChildDiscountedDeathOutgo + firstPP.ChildDiscountedAccidentalDeathOutgo) / ((firstPP.ChildAnnuityFactor / 12) - (pricingParams.InitialCommissionPercentage1 * temp1))
				}
			}

		}

		if mp.CommissionType == Renewal {
			if pricingConfig.CashBackOnSurvival {
				*calculatedMinimumPremium = (firstPP.DiscountedDeathOutgo + firstPP.ChildDiscountedDeathOutgo + firstPP.SpouseDiscountedDeathOutgo + firstPP.DiscountedAccidentalDeathOutgo + firstPP.DiscountedDisabilityOutgo + firstPP.DiscountedRetrenchmentOutgo + firstPP.ChildDiscountedAccidentalDeathOutgo + firstPP.SpouseDiscountedAccidentalDeathOutgo + firstPP.DiscountedEducator + firstPP.DiscountedRider + firstPP.DiscountedExpenses + firstPP.DiscountedRiderFuneral + firstPP.RenewalCommissionAnnuityFactor*pricingParams.RenewalCommissionRand) / ((firstPP.AnnuityFactor / 12.0) * (1 - pricingParams.RenewalCommissionPercentage - parameters.CashbackOnSurvivalRatio))
				if pricingConfig.ChildIndicator {
					*childMinimumCalculatedPremium = (firstPP.ChildDiscountedDeathOutgo + firstPP.ChildDiscountedAccidentalDeathOutgo) / ((firstPP.ChildAnnuityFactor / 12) * (1 - pricingParams.RenewalCommissionPercentage - parameters.CashbackOnSurvivalRatio))
				}
			} else {
				*calculatedMinimumPremium = (firstPP.DiscountedDeathOutgo + firstPP.ChildDiscountedDeathOutgo + firstPP.SpouseDiscountedDeathOutgo + firstPP.DiscountedAccidentalDeathOutgo + firstPP.DiscountedDisabilityOutgo + firstPP.DiscountedRetrenchmentOutgo + firstPP.ChildDiscountedAccidentalDeathOutgo + firstPP.SpouseDiscountedAccidentalDeathOutgo + firstPP.DiscountedEducator + firstPP.DiscountedRider + firstPP.DiscountedExpenses + firstPP.DiscountedRiderFuneral + firstPP.RenewalCommissionAnnuityFactor*pricingParams.RenewalCommissionRand) / ((firstPP.AnnuityFactor / 12.0) * (1 - pricingParams.RenewalCommissionPercentage))
				if pricingConfig.ChildIndicator {
					*childMinimumCalculatedPremium = (firstPP.ChildDiscountedDeathOutgo + firstPP.ChildDiscountedAccidentalDeathOutgo) / ((firstPP.ChildAnnuityFactor / 12) * (1 - pricingParams.RenewalCommissionPercentage))
				}
			}
		} //We need to create enums here...

		if mp.CommissionType == Hybrid {
			if pricingConfig.CashBackOnSurvival {
				*calculatedMinimumPremium = (firstPP.DiscountedDeathOutgo + firstPP.ChildDiscountedDeathOutgo + firstPP.SpouseDiscountedDeathOutgo + firstPP.DiscountedAccidentalDeathOutgo + firstPP.DiscountedDisabilityOutgo + firstPP.DiscountedRetrenchmentOutgo + firstPP.ChildDiscountedAccidentalDeathOutgo + firstPP.SpouseDiscountedAccidentalDeathOutgo + firstPP.DiscountedEducator + firstPP.DiscountedRider + firstPP.DiscountedExpenses + firstPP.DiscountedRiderFuneral + pricingParams.InitialCommissionRand + firstPP.RenewalCommissionAnnuityFactor*pricingParams.RenewalCommissionRand) / (((firstPP.AnnuityFactor - pricingParams.RenewalCommissionPercentage*firstPP.RenewalCommissionAnnuityFactor - parameters.CashbackOnSurvivalRatio) / 12.0) - pricingParams.InitialCommissionPercentage1*temp1)
				if pricingConfig.ChildIndicator {
					*childMinimumCalculatedPremium = (firstPP.ChildDiscountedDeathOutgo + firstPP.ChildDiscountedAccidentalDeathOutgo) / (((firstPP.AnnuityFactor - pricingParams.RenewalCommissionPercentage*firstPP.RenewalCommissionAnnuityFactor - parameters.CashbackOnSurvivalRatio) / 12.0) - pricingParams.InitialCommissionPercentage1*temp1)
				}
			} else {
				*calculatedMinimumPremium = (firstPP.DiscountedDeathOutgo + firstPP.ChildDiscountedDeathOutgo + firstPP.SpouseDiscountedDeathOutgo + firstPP.DiscountedAccidentalDeathOutgo + firstPP.DiscountedDisabilityOutgo + firstPP.DiscountedRetrenchmentOutgo + firstPP.ChildDiscountedAccidentalDeathOutgo + firstPP.SpouseDiscountedAccidentalDeathOutgo + firstPP.DiscountedEducator + firstPP.DiscountedRider + firstPP.DiscountedExpenses + firstPP.DiscountedRiderFuneral + pricingParams.InitialCommissionRand + firstPP.RenewalCommissionAnnuityFactor*pricingParams.RenewalCommissionRand) / (((firstPP.AnnuityFactor - pricingParams.RenewalCommissionPercentage*firstPP.RenewalCommissionAnnuityFactor) / 12.0) - pricingParams.InitialCommissionPercentage1*temp1)
				if pricingConfig.ChildIndicator {
					*childMinimumCalculatedPremium = (firstPP.ChildDiscountedDeathOutgo + firstPP.ChildDiscountedAccidentalDeathOutgo) / ((firstPP.ChildAnnuityFactor/12)*(1-pricingParams.RenewalCommissionPercentage) - pricingParams.InitialCommissionPercentage1*temp1 + pricingParams.RenewalCommissionPercentage)
				}
			}
		}

		if mp.CommissionType != Initial && mp.CommissionType != Renewal && mp.CommissionType != Hybrid {
			if pricingConfig.CashBackOnSurvival {
				*calculatedMinimumPremium = (firstPP.DiscountedDeathOutgo + firstPP.ChildDiscountedDeathOutgo + firstPP.SpouseDiscountedDeathOutgo + firstPP.DiscountedAccidentalDeathOutgo + firstPP.DiscountedDisabilityOutgo + firstPP.DiscountedRetrenchmentOutgo + firstPP.ChildDiscountedAccidentalDeathOutgo + firstPP.SpouseDiscountedAccidentalDeathOutgo + firstPP.DiscountedEducator + firstPP.DiscountedRider + firstPP.DiscountedExpenses + firstPP.DiscountedRiderFuneral) / ((firstPP.AnnuityFactor / 12.0) * (1 - parameters.CashbackOnSurvivalRatio))
				if pricingConfig.ChildIndicator {
					*childMinimumCalculatedPremium = (firstPP.ChildDiscountedDeathOutgo + firstPP.ChildDiscountedAccidentalDeathOutgo) / ((firstPP.ChildAnnuityFactor / 12) * (1 - parameters.CashbackOnSurvivalRatio))
				}
			} else {
				*calculatedMinimumPremium = (firstPP.DiscountedDeathOutgo + firstPP.ChildDiscountedDeathOutgo + firstPP.SpouseDiscountedDeathOutgo + firstPP.DiscountedAccidentalDeathOutgo + firstPP.DiscountedDisabilityOutgo + firstPP.DiscountedRetrenchmentOutgo + firstPP.ChildDiscountedAccidentalDeathOutgo + firstPP.SpouseDiscountedAccidentalDeathOutgo + firstPP.DiscountedEducator + firstPP.DiscountedRider + firstPP.DiscountedExpenses + firstPP.DiscountedRiderFuneral) / (firstPP.AnnuityFactor / 12.0)
				if pricingConfig.ChildIndicator {
					*childMinimumCalculatedPremium = (firstPP.ChildDiscountedDeathOutgo + firstPP.ChildDiscountedAccidentalDeathOutgo) / (firstPP.ChildAnnuityFactor / 12)
				}
			}
		}

		if mp.CommissionType == Renewal {
			*technicalPremiumAfterProfitmargin = *calculatedMinimumPremium / (1 - getProfitMargin(mp.AgeAtEntry, int(firstPP.SumAssured), mp.ProductCode) - pricingParams.RenewalCommissionPercentage*getProfitMargin(mp.AgeAtEntry, int(firstPP.SumAssured), mp.ProductCode)) //pricingParams.TargetProfitMargin)
		}
		if mp.CommissionType == Hybrid {
			*technicalPremiumAfterProfitmargin = *calculatedMinimumPremium / (1 - getProfitMargin(mp.AgeAtEntry, int(firstPP.SumAssured), mp.ProductCode)) //pricingParams.TargetProfitMargin)
		}
		if mp.CommissionType != Renewal && mp.CommissionType != Hybrid {
			*technicalPremiumAfterProfitmargin = *calculatedMinimumPremium / (1 - getProfitMargin(mp.AgeAtEntry, int(firstPP.SumAssured), mp.ProductCode)) //pricingParams.TargetProfitMargin)
		}

		if pricingConfig.ChildIndicator {
			*childFinalCalculatedPremium = *childMinimumCalculatedPremium / (1 - getProfitMargin(mp.AgeAtEntry, int(firstPP.SumAssured), mp.ProductCode)) //pricingParams.TargetProfitMargin)
		}
		*previousAnnualPremium = *technicalPremiumAfterProfitmargin
		if mp.OutstandingLoan > 0 {
			*previousPremiumRate = utils.FloatPrecision(*technicalPremiumAfterProfitmargin*1000.0/(mp.OutstandingLoan), AccountingPrecision) // Annual premium rate
		}
	} else {
		*previousAnnualPremium = mp.AnnualPremium
		*previousPremiumRate = mp.PremiumRate
	}

	if math.Abs(*variance) < 0.5 {
		// Save calculatedMinimumPremium and all discounted values at month zero for the modelpoint and for each scenario
		//fmt.Println("running pricing for Model Point: ", mp.PolicyNumber, " and ", pricingConfig.Description, " scenario")
		//fmt.Println("RevenueVariance: ", *variance)
		//fmt.Println("GoalSeekStep: ", gskIndex)

		var mpp models.ModelPointPricing
		mpp.RunID = pricingPoints[0].RunId
		mpp.ProductCode = mp.ProductCode
		mpp.PolicyNumber = mp.PolicyNumber
		mpp.ScenarioID = pricingConfig.ID
		mpp.PricingRunId = pricingConfig.PricingRunID
		mpp.ScenarioDescription = pricingConfig.Description
		mpp.CalculatedMinimumPremium = *calculatedMinimumPremium
		mpp.CalculatedMinimumChildPremium = *childMinimumCalculatedPremium
		mpp.CalculatedAnnualPremium = mp.AnnualPremium
		mpp.CalculatedPremiumRate = mp.PremiumRate

		//**replacing reserve variables with profitability variables
		mpp.DiscountedCommission = pricingPoints[0].ProfDiscountedCommission //pricingPoints[0].DiscountedCommission
		mpp.DiscountedProfit = pricingPoints[0].DiscountedProfit
		mpp.DiscountedExpenses = pricingPoints[0].ProfDiscountedExpenses //pricingPoints[0].DiscountedExpenses
		mpp.DiscountedRider = pricingPoints[0].ProfDiscountedRider
		mpp.DiscountedRiderFuneral = pricingPoints[0].ProfDiscountedRiderFuneral
		mpp.DiscountedCashBackOnSurvival = pricingPoints[0].ProfDiscountedCashBackOnSurvival            //pricingPoints[0].DiscountedCashBackOnSurvival
		mpp.DiscountedEducator = pricingPoints[0].ProfDiscountedEducator                                //pricingPoints[0].DiscountedEducator
		mpp.DiscountedAccidentalDeathOutgo = pricingPoints[0].ProfDiscountedAccidentalDeath             // pricingPoints[0].DiscountedAccidentalDeathOutgo
		mpp.ChildDiscountedAccidentalDeathOutgo = pricingPoints[0].ProfChildDiscountedAccidentalDeath   //pricingPoints[0].ChildDiscountedAccidentalDeathOutgo
		mpp.SpouseDiscountedAccidentalDeathOutgo = pricingPoints[0].ProfSpouseDiscountedAccidentalDeath //pricingPoints[0].SpouseDiscountedAccidentalDeathOutgo
		mpp.DiscountedDeathOutgo = pricingPoints[0].ProfDiscountedDeath                                 //pricingPoints[0].DiscountedDeathOutgo
		mpp.DiscountedDisabilityOutgo = pricingPoints[0].ProfDiscountedDisability                       //pricingPoints[0].DiscountedDisabilityOutgo
		mpp.DiscountedRetrenchmentOutgo = pricingPoints[0].ProfDiscountedRetrenchment                   //pricingPoints[0].DiscountedRetrenchmentOutgo
		mpp.ChildDiscountedDeathOutgo = pricingPoints[0].ProfChildDiscountedDeath                       //pricingPoints[0].ChildDiscountedDeathOutgo
		mpp.SpouseDiscountedDeathOutgo = pricingPoints[0].ProfSpouseDiscountedDeath                     //pricingPoints[0].SpouseDiscountedDeathOutgo
		mpp.DiscountedPremiumIncome = pricingPoints[0].ProfDiscountedPremium                            //pricingPoints[0].DiscountedPremiumIncome
		mpp.DiscountedPremiumNotReceived = pricingPoints[0].ProfDiscountedPremiumNotReceived            //pricingPoints[0].DiscountedPremiumNotReceived
		mpp.DiscountedClawBack = pricingPoints[0].ProfDiscountedClawback                                //pricingPoints[0].DiscountedClawBack
		mpp.DiscountedChangeInReserve = pricingPoints[0].ProfDiscountedChangeInReserve
		mpp.DiscountedInvestmentIncome = pricingPoints[0].ProfDiscountedInvestmentIncome
		mpp.AnnuityFactor = pricingPoints[0].AnnuityFactor
		mpp.ProfRiskAdjustment = pricingPoints[0].ProfRiskAdjustment
		mpp.ProfCSM = pricingPoints[0].ProfCSM
		mpp.ProfLossComponent = pricingPoints[0].ProfLossComponent
		mpp.Weighting = mp.Weighting
		mpp.Age = mp.AgeAtEntry
		mpp.SumAssured = mp.SumAssured
		mpp.Gender = mp.MainMemberGender

		err = DB.Create(&mpp).Error
		if err != nil {
			fmt.Println(errors.WithStack(err))
		} else {
			pf.WeightedDiscountedCashBackOnSurvival += mpp.DiscountedCashBackOnSurvival * mpp.Weighting
			pf.WeightedDiscountedCashBackOnDeath += mpp.DiscountedCashBackOnDeath * mpp.Weighting
			pf.WeightedDiscountedCommission += (mpp.DiscountedCommission - mpp.DiscountedClawBack) * mpp.Weighting
			pf.WeightedDiscountedExpenses += mpp.DiscountedExpenses * mpp.Weighting
			//pf.WeightedDiscountedChangeInReserve += (mpp.DiscountedChangeInReserve - mpp.DiscountedInvestmentIncome) * mpp.Weighting
			pf.WeightedDiscountedPremiumNotReceived += mpp.DiscountedPremiumNotReceived * mpp.Weighting
			pf.WeightedDiscountedProfit += mpp.DiscountedProfit * mpp.Weighting
			pf.WeightedDiscountedRider += (mpp.DiscountedRider + mpp.DiscountedRiderFuneral) * mpp.Weighting
			//pf.WeightedDiscountedInvestmentIncome += mpp.DiscountedInvestmentIncome * mpp.Weighting
			pf.WeightedDiscountedRisk += (mpp.DiscountedDeathOutgo + mpp.DiscountedDisabilityOutgo + mpp.DiscountedRetrenchmentOutgo + mpp.ChildDiscountedDeathOutgo + mpp.SpouseDiscountedDeathOutgo + mpp.DiscountedAccidentalDeathOutgo + mpp.ChildDiscountedAccidentalDeathOutgo + mpp.SpouseDiscountedAccidentalDeathOutgo + (mpp.DiscountedChangeInReserve - mpp.DiscountedInvestmentIncome)) * mpp.Weighting
		}

		// save only the first 120 months of the pricing points
		if k == 0 {
			err = DB.CreateInBatches(&pricingPoints, 100).Error
			if err != nil {
				fmt.Println(err)
			}
		}

		if profitSignature && !runGoalSeek {
			UpdatePricingAggregatedProjections(pricingPoints[:120], aggpricingProjs, mp.Weighting)
		}
	}
}

func calculateProductPricingConfig(parameters *models.ProductPricingParameters, mp models.ProductPricingModelPoint, features models.ProductFeatures, multipliers models.ProductPricingAccidentalBenefitMultiplier) {

	if mp.MemberType == "CH" {
		parameters.CalculatedTerm = int(math.Min(math.Max(float64(parameters.ChildExitAge-mp.AgeAtEntry), 0.0)*12.0, float64(mp.Term)))
	} else if features.CreditLife {
		parameters.CalculatedTerm = mp.OutstandingTermMonths + mp.DurationInForceMonths + 1
	} else if features.WholeOfLife {
		parameters.CalculatedTerm = (120 - mp.AgeAtEntry) * 12
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

}

func pricingTimeAndAge(pp *models.PricingPoint, mp models.ProductPricingModelPoint, i int, features models.ProductFeatures) {
	pp.ProjectionMonth = i
	pp.ProjectionYear = int(math.Ceil(float64(i) / 12.0))
	pp.PolicyNumber = mp.PolicyNumber
	pp.ValuationTimeMonth = mp.DurationInForceMonths + i
	pp.ValuationTimeYear = utils.FloatPrecision(float64(pp.ValuationTimeMonth)/12, 2)
	//if features.FuneralCover {
	pp.MainMemberAgeNextBirthday = int(math.Ceil(float64(mp.MainMemberAgeAtEntry) + pp.ValuationTimeYear))
	//}
	//if features.CreditLife {
	//	pp.MainMemberAgeNextBirthday = int(math.Ceil(float64(mp.AgeAtEntry) + pp.ValuationTimeYear))
	//}
	pp.AgeNextBirthday = int(math.Ceil(float64(mp.AgeAtEntry) + pp.ValuationTimeYear))
}

func pricingDiscountedValues(index int, pp *models.PricingPoint,
	mp models.ProductPricingModelPoint, margins models.ProductPricingMargins,
	pricingParams models.PricingParameter, pricingPoints *[]models.PricingPoint,
	features models.ProductFeatures, shock models.ProductPricingShock, pricingShockBasis string, params models.ProductPricingParameters) {
	//we are working on the previous projection results, right?
	forwardRate := getPricingForwardRate(index+1, pricingParams.YieldCurveCode, shock, pricingShockBasis)

	discountFactor := math.Pow(1.0+forwardRate, -1/12.0)
	discountFactorAdjusted := math.Pow(1.0+forwardRate+margins.InvestmentMargin, -1/12.0)

	//disabilityDiscountFactor := math.Pow(1.0+forwardRate, -(1.0+float64(mp.DeferredPeriod))/12.0)
	//disabilityDiscountFactorAdjusted := math.Pow(1.0+forwardRate+margins.InvestmentMargin, -(1.0+float64(mp.DeferredPeriod))/12.0)

	pp.DiscountFactor = discountFactor
	pp.DiscountFactorAdjusted = discountFactorAdjusted

	if index == len(*pricingPoints)-1 {
		pp.DiscountedPremiumIncome = 0
		pp.DiscountedPremiumNotReceived = 0
		pp.DiscountedCommission = 0
		pp.DiscountedClawBack = 0
		pp.DiscountedDeathOutgo = 0
		pp.ChildDiscountedDeathOutgo = 0
		pp.SpouseDiscountedDeathOutgo = 0
		pp.DiscountedAccidentalDeathOutgo = 0
		pp.ChildDiscountedAccidentalDeathOutgo = 0
		pp.SpouseDiscountedAccidentalDeathOutgo = 0
		pp.DiscountedCashBackOnSurvival = 0
		pp.DiscountedCashBackOnDeath = 0
		pp.DiscountedRider = 0
		pp.DiscountedEducator = 0
		pp.DiscountedExpenses = 0
		pp.DiscountedProfit = 0
		pp.DiscountedProfitAdjusted = 0
		pp.DiscountedDisabilityOutgo = 0
		pp.DiscountedRetrenchmentOutgo = 0

		pp.ProfDiscountedPremium = 0
		pp.ProfDiscountedPremiumNotReceived = 0
		pp.ProfDiscountedDeath = 0
		pp.ProfDiscountedDisability = 0
		pp.ProfDiscountedRetrenchment = 0
		pp.ProfDiscountedCommission = 0
		pp.ProfDiscountedClawback = 0
		pp.ProfDiscountedAccidentalDeath = 0
		pp.ProfDiscountedExpenses = 0
		pp.ProfDiscountedChangeInReserve = 0
		pp.ProfDiscountedInvestmentIncome = 0
		pp.ProfSpouseDiscountedDeath = 0
		pp.ProfSpouseDiscountedAccidentalDeath = 0
		pp.ProfChildDiscountedDeath = 0
		pp.ProfChildDiscountedAccidentalDeath = 0
		pp.ProfDiscountedEducator = 0
		pp.ProfDiscountedCashBackOnSurvival = 0
		pp.ProfDiscountedCashBackOnDeath = 0
		pp.ProfDiscountedRider = 0
		pp.ProfDiscountedRiderFuneral = 0

		pp.DiscountedRenewalCommissionAnnuityFeeder = 0

	} else {
		pp.DiscountedPremiumIncome = (*pricingPoints)[index+1].PremiumIncome + (*pricingPoints)[index+1].DiscountedPremiumIncome*discountFactor

		pp.DiscountedPremiumNotReceived = (*pricingPoints)[index+1].PremiumNotReceived + (*pricingPoints)[index+1].DiscountedPremiumNotReceived*discountFactor

		pp.DiscountedCommission = (*pricingPoints)[index+1].Commission + (*pricingPoints)[index+1].DiscountedCommission*discountFactor

		pp.DiscountedRenewalCommissionAnnuityFeeder = (*pricingPoints)[index+1].RenewalCommissionAnnuityFeeder + (*pricingPoints)[index+1].DiscountedRenewalCommissionAnnuityFeeder*discountFactor

		pp.DiscountedClawBack = ((*pricingPoints)[index+1].ClawBack + (*pricingPoints)[index+1].DiscountedClawBack) * discountFactor

		pp.DiscountedDeathOutgo = ((*pricingPoints)[index+1].DeathOutgo + (*pricingPoints)[index+1].DiscountedDeathOutgo) * discountFactor //ann_prem

		pp.ChildDiscountedDeathOutgo = ((*pricingPoints)[index+1].ChildDeathOutgo + (*pricingPoints)[index+1].ChildDiscountedDeathOutgo) * discountFactor

		pp.SpouseDiscountedDeathOutgo = ((*pricingPoints)[index+1].SpouseDeathOutgo + (*pricingPoints)[index+1].SpouseDiscountedDeathOutgo) * discountFactor

		pp.DiscountedAccidentalDeathOutgo = ((*pricingPoints)[index+1].AccidentalDeathOutgo + (*pricingPoints)[index+1].DiscountedAccidentalDeathOutgo) * discountFactor //ann_prem

		pp.ChildDiscountedAccidentalDeathOutgo = ((*pricingPoints)[index+1].ChildAccidentalDeathOutgo + (*pricingPoints)[index+1].ChildDiscountedAccidentalDeathOutgo) * discountFactor //ann_prem

		pp.SpouseDiscountedAccidentalDeathOutgo = ((*pricingPoints)[index+1].SpouseAccidentalDeathOutgo + (*pricingPoints)[index+1].SpouseDiscountedAccidentalDeathOutgo) * discountFactor //ann_prem

		pp.DiscountedCashBackOnSurvival = ((*pricingPoints)[index+1].CashBackOnSurvival + (*pricingPoints)[index+1].DiscountedCashBackOnSurvival) * discountFactor //ann_prem

		pp.DiscountedCashBackOnDeath = ((*pricingPoints)[index+1].CashBackOnDeath + (*pricingPoints)[index+1].DiscountedCashBackOnDeath) * discountFactor //ann_prem

		pp.DiscountedRiderFuneral = ((*pricingPoints)[index+1].RiderFuneral + (*pricingPoints)[index+1].DiscountedRiderFuneral) * discountFactor //ann_prem

		pp.DiscountedRider = ((*pricingPoints)[index+1].Rider + (*pricingPoints)[index+1].DiscountedRider) * discountFactor //ann_prem

		pp.DiscountedEducator = ((*pricingPoints)[index+1].Educator + (*pricingPoints)[index+1].DiscountedEducator) * discountFactor //ann_prem

		pp.DiscountedExpenses = (*pricingPoints)[index+1].Expenses + (*pricingPoints)[index+1].DiscountedExpenses*discountFactor //ann_prem

		pp.DiscountedDisabilityOutgo = ((*pricingPoints)[index+1].DisabilityOutgo + (*pricingPoints)[index+1].DiscountedDisabilityOutgo) * discountFactor // disabilityDiscountFactor

		pp.DiscountedRetrenchmentOutgo = ((*pricingPoints)[index+1].RetrenchmentOutgo + (*pricingPoints)[index+1].DiscountedRetrenchmentOutgo) * discountFactor

		pp.ProfDiscountedPremium = ((*pricingPoints)[index+1].PremiumIncome + (*pricingPoints)[index+1].ProfDiscountedPremium) * discountFactor
		pp.ProfDiscountedPremiumNotReceived = ((*pricingPoints)[index+1].PremiumNotReceived + (*pricingPoints)[index+1].ProfDiscountedPremiumNotReceived) * discountFactor
		pp.ProfDiscountedDeath = ((*pricingPoints)[index+1].DeathOutgo + (*pricingPoints)[index+1].ProfDiscountedDeath) * discountFactor
		pp.ProfDiscountedDisability = ((*pricingPoints)[index+1].DisabilityOutgo + (*pricingPoints)[index+1].ProfDiscountedDisability) * discountFactor
		pp.ProfDiscountedRetrenchment = ((*pricingPoints)[index+1].RetrenchmentOutgo + (*pricingPoints)[index+1].ProfDiscountedRetrenchment) * discountFactor
		pp.ProfDiscountedCommission = ((*pricingPoints)[index+1].Commission + (*pricingPoints)[index+1].ProfDiscountedCommission) * discountFactor
		pp.ProfDiscountedClawback = ((*pricingPoints)[index+1].ClawBack + (*pricingPoints)[index+1].ProfDiscountedClawback) * discountFactor
		pp.ProfDiscountedAccidentalDeath = ((*pricingPoints)[index+1].AccidentalDeathOutgo + (*pricingPoints)[index+1].ProfDiscountedAccidentalDeath) * discountFactor
		pp.ProfDiscountedExpenses = ((*pricingPoints)[index+1].Expenses + (*pricingPoints)[index+1].ProfDiscountedExpenses) * discountFactor
		pp.ProfDiscountedChangeInReserve = ((*pricingPoints)[index+1].ChangeInReserves + (*pricingPoints)[index+1].ProfDiscountedChangeInReserve) * discountFactor
		pp.ProfDiscountedInvestmentIncome = ((*pricingPoints)[index+1].InvestmentIncome + (*pricingPoints)[index+1].ProfDiscountedInvestmentIncome) * discountFactor
		pp.ProfSpouseDiscountedDeath = ((*pricingPoints)[index+1].SpouseDeathOutgo + (*pricingPoints)[index+1].ProfSpouseDiscountedDeath) * discountFactor
		pp.ProfSpouseDiscountedAccidentalDeath = ((*pricingPoints)[index+1].SpouseAccidentalDeathOutgo + (*pricingPoints)[index+1].ProfSpouseDiscountedAccidentalDeath) * discountFactor
		pp.ProfChildDiscountedDeath = ((*pricingPoints)[index+1].ChildDeathOutgo + (*pricingPoints)[index+1].ProfChildDiscountedDeath) * discountFactor
		pp.ProfChildDiscountedAccidentalDeath = ((*pricingPoints)[index+1].ChildAccidentalDeathOutgo + (*pricingPoints)[index+1].ProfChildDiscountedAccidentalDeath) * discountFactor
		pp.ProfDiscountedEducator = ((*pricingPoints)[index+1].Educator + (*pricingPoints)[index+1].ProfDiscountedEducator) * discountFactor
		pp.ProfDiscountedCashBackOnSurvival = ((*pricingPoints)[index+1].CashBackOnSurvival + (*pricingPoints)[index+1].ProfDiscountedCashBackOnSurvival) * discountFactor
		pp.ProfDiscountedCashBackOnDeath = ((*pricingPoints)[index+1].CashBackOnDeath + (*pricingPoints)[index+1].ProfDiscountedCashBackOnDeath) * discountFactor
		pp.ProfDiscountedRider = ((*pricingPoints)[index+1].Rider + (*pricingPoints)[index+1].ProfDiscountedRider) * discountFactor
		pp.ProfDiscountedRiderFuneral = ((*pricingPoints)[index+1].RiderFuneral + (*pricingPoints)[index+1].ProfDiscountedRiderFuneral) * discountFactor
		if features.CreditLifeDecreasingPremium {
			if mp.OutstandingLoan > 0 {
				pp.AnnuityFactor = pp.InitialPolicy*(*pricingPoints)[index+1].PremiumEscalation*pp.OutstandingSumAssured/mp.OutstandingLoan + (*pricingPoints)[index+1].AnnuityFactor*discountFactor
			}
		}

	}

	if mp.AnnualPremium > 0 {
		if !features.CreditLifeDecreasingPremium {
			pp.AnnuityFactor = 12.0 * pp.DiscountedPremiumIncome / mp.AnnualPremium
			if pricingParams.RenewalCommissionPercentage > 0 {
				pp.RenewalCommissionAnnuityFactor = 12.0 * pp.DiscountedRenewalCommissionAnnuityFeeder / (mp.AnnualPremium * pricingParams.RenewalCommissionPercentage)
			}

			if pricingParams.RenewalCommissionRand > 0 {
				pp.RenewalCommissionAnnuityFactor = pp.DiscountedRenewalCommissionAnnuityFeeder / (pricingParams.RenewalCommissionRand)
			}

			if pricingParams.RenewalCommissionPercentage > 0 && pricingParams.RenewalCommissionRand > 0 {
				pp.RenewalCommissionAnnuityFactor = pp.DiscountedRenewalCommissionAnnuityFeeder / ((mp.AnnualPremium * pricingParams.RenewalCommissionPercentage / 12.0) + pricingParams.RenewalCommissionRand)
			}
		}
	} else {
		pp.AnnuityFactor = 0
	}

	pp.Reserves = pp.DiscountedPremiumNotReceived +
		pp.DiscountedCommission +
		pp.DiscountedDeathOutgo + pp.DiscountedDisabilityOutgo + pp.DiscountedRetrenchmentOutgo +
		pp.ChildDiscountedDeathOutgo + pp.SpouseDiscountedDeathOutgo +
		pp.DiscountedAccidentalDeathOutgo + pp.ChildDiscountedAccidentalDeathOutgo + pp.SpouseDiscountedAccidentalDeathOutgo +
		pp.DiscountedEducator +
		pp.DiscountedCashBackOnDeath +
		pp.DiscountedCashBackOnSurvival +
		pp.DiscountedRiderFuneral +
		pp.DiscountedRider +
		pp.DiscountedExpenses -
		pp.DiscountedClawBack -
		pp.DiscountedPremiumIncome

	pp.ReservesAdjusted = pp.DiscountedPremiumNotReceivedAdjusted +
		pp.DiscountedCommissionAdjusted +
		pp.DiscountedDeathOutgoAdjusted + pp.DiscountedDisabilityOutgoAdjusted + pp.DiscountedRetrenchmentOutgoAdjusted +
		pp.ChildDiscountedDeathOutgoAdjusted + pp.SpouseDiscountedDeathOutgoAdjusted +
		pp.DiscountedAccidentalDeathOutgoAdjusted + pp.ChildDiscountedAccidentalDeathOutgoAdjusted + pp.SpouseDiscountedAccidentalDeathOutgoAdjusted +
		pp.DiscountedEducatorAdjusted +
		pp.DiscountedCashBackOnDeathAdjusted +
		pp.DiscountedCashBackOnSurvivalAdjusted +
		pp.DiscountedRiderFuneralAdjusted +
		pp.DiscountedRiderAdjusted +
		pp.DiscountedExpensesAdjusted -
		pp.DiscountedClawBackAdjusted -
		pp.DiscountedPremiumIncomeAdjusted

}

func pricingInterestAccrual(index int, p *models.PricingPoint, prev *models.PricingPoint, margins models.ProductPricingMargins, pricingParams models.PricingParameter, shock models.ProductPricingShock, pricingShockBasis string) {
	forwardRate := getPricingForwardRate(index, pricingParams.YieldCurveCode, shock, pricingShockBasis)
	if p.ProjectionMonth == 0 {
		p.InvestmentIncome = 0
		p.InvestmentIncomeAdjusted = 0
	} else {
		p.InvestmentIncome = (prev.Reserves - (p.PremiumNotReceived + p.Commission + p.Expenses - p.PremiumIncome)) * (math.Pow(1.0+forwardRate, 1.0/12.0) - 1.0)
		p.InvestmentIncomeAdjusted = (prev.ReservesAdjusted - (p.PremiumNotReceivedAdjusted + p.CommissionAdjusted + p.ExpensesAdjusted - p.PremiumIncomeAdjusted)) * (math.Pow(1.0+forwardRate+margins.InvestmentMargin, 1.0/12.0) - 1.0)
	}
}

func getPricingRatingTable(code string, class string) string {
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

func getPricingLapseMarginCount(code string, pricingparams models.PricingParameter) int {
	var count int64
	err := DB.Table("product_pricing_lapse_margins").Where("product_code = ? and basis=?", code, pricingparams.Basis).Count(&count).Error
	if err != nil {
		return 0
	}
	return int(count)
}

func getPricingLapseCount(lapseTableName string) int {
	var count int64
	err := DB.Table(lapseTableName).Count(&count).Error
	if err != nil {
		return 0
	}
	return int(count)
}

func getPricingRetrenchmentCount(RetrenchmentTableName string) int {
	var count int64
	err := DB.Table(RetrenchmentTableName).Count(&count).Error
	if err != nil {
		return 0
	}
	return int(count)
}

func loadPricingMortalityRates(productCode string) {
	if pricingTables[productCode].MortalityTableName != "" {
		//gender := []string{"M", "F"}
		tablename := pricingTables[productCode].MortalityTableName //strings.ToLower(productCode) + "_pricing_" + strings.ToLower(pricingMortalityTableName)
		columnName := pricingTables[productCode].MortalityColumnName

		var qx float64
		var cname string
		key := columnName
		cacheKey := tablename + "_" + key
		//Select(fmt.Sprintf("%s as cname, qx",columnName))
		rows, err := DB.Table(tablename).Select(fmt.Sprintf("%s as cname, qx", columnName)).Rows()
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
		cached, found := PricingCache.Get(tablename + "_3_F")
		if found {
			fmt.Println("load key: ", cached)
		} else {
			fmt.Println("NOT FOUND")
		}
	}
}

func loadPricingAccidentalMortalityRates(productCode string) {
	if pricingTables[productCode].MortalityAccidentalTableName != "" {
		tablename := pricingTables[productCode].MortalityAccidentalTableName //strings.ToLower(productCode) + "_pricing_" + strings.ToLower(pricingMortalityAccidentalTableName)
		columnName := pricingTables[productCode].MortalityAccidentalColumnName
		var qx float64
		var cname string
		key := columnName
		cacheKey := tablename + "_" + key
		//Select(fmt.Sprintf("%s as cname, qx",columnName))
		rows, err := DB.Table(tablename).Select(fmt.Sprintf("%s as cname, acc_qx_prop", columnName)).Rows()
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
		cached, found := PricingCache.Get(tablename + "_3_F")
		if found {
			fmt.Println("load key: ", cached)
		} else {
			fmt.Println("NOT FOUND")
		}
	}

}

func loadPricingInflationFactor() {
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

func loadPricingForwardRate() {
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

func getPricingMortalityRateAccidentProportion(arg models.TransitionRateArguments, columnname, tableName string) float64 { //age int, gender string, productCode string) float64 {
	//key := strconv.Itoa(age) + "_" + gender[:1]
	key := PricingkeyBuilder(arg, AccidentalDeath)
	//tableName := strings.ToLower(arg.ProductCode) + "_pricing_" + strings.ToLower(pricingMortalityAccidentalTableName)
	columnName := columnname //"anb_gender" //GetColumnName(tableName)
	cacheKey := strings.ToLower(tableName) + "_" + key
	cached, found := PricingCache.Get(cacheKey)

	if found {
		result := cached.(float64)
		if result > 0 {
			return result
		}
	}

	var accQxProp float64
	query := columnName + " = ?"
	row := DB.Table(tableName).Where(query, key).Select("acc_qx_prop").Row()
	err := row.Scan(&accQxProp)
	if err != nil {
		fmt.Println("db error: ", errors.WithStack(err))
	}
	PricingCache.Set(cacheKey, accQxProp, 1)
	return accQxProp
}

func getPricingInflationFactor(month int, yieldCurveCode string) float64 {
	var yieldCurve models.PricingYieldCurve
	key := "y-curve-" + strconv.Itoa(month) + yieldCurveCode
	result, found := PricingCache.Get(key)

	if found {
		return result.(float64)
	} else {
		//TODO:Use pricing yield curve here instead...
		err := DB.Table("pricing_yield_curve").Where("proj_time = ? and yield_curve_code = ?", month, yieldCurveCode).First(&yieldCurve).Error
		if err != nil {
			fmt.Println("Yield Curve: ", errors.WithStack(err))
		}
		success := PricingCache.Set(key, yieldCurve.Inflation, 1)
		if !success {
			fmt.Println("Cache: key not stored")
		}
	}
	return yieldCurve.Inflation
}

func getPricingForwardRate(month int, yieldCurveCode string, shock models.ProductPricingShock, pricingShockBasis string) float64 {
	var yieldCurve models.PricingYieldCurve
	var tempNominalRate float64
	key := "y-curve-fr-" + strconv.Itoa(month) + yieldCurveCode
	cached, found := PricingCache.Get(key)

	if found {
		result := cached.(float64)
		tempNominalRate = result
		if pricingShockBasis != "N/A" {
			tempNominalRate = tempNominalRate*(1+shock.MultiplicativeYieldCurve) + shock.AdditiveYieldCurve
		}
		return tempNominalRate
	}
	if month == 0 {
		return 0
	}
	err := DB.Table("pricing_yield_curve").Where("proj_time = ? and yield_curve_code = ?", month, yieldCurveCode).First(&yieldCurve).Error
	if err != nil {
		fmt.Println("db error: ", errors.WithStack(err))
	}
	success := PricingCache.Set(key, yieldCurve.NominalRate, 1)
	if !success {
		fmt.Println("Cache error: Not saved")
	}

	if pricingShockBasis != "N/A" {
		tempNominalRate = yieldCurve.NominalRate*(1+shock.MultiplicativeYieldCurve) + shock.AdditiveYieldCurve
	}

	if pricingShockBasis == "N/A" {
		tempNominalRate = yieldCurve.NominalRate
	}

	return tempNominalRate
}

func getPricingMortalityRate(arg models.TransitionRateArguments, columnname, tableName string) float64 {
	//tablename := strings.ToLower(arg.ProductCode) + "_pricing_" + strings.ToLower(pricingMortalityTableName)
	//key := strconv.Itoa(age) + "_" + gender[:1]
	key := PricingkeyBuilder(arg, Death)
	columnName := columnname //"anb_gender"//GetColumnName(tablename)
	cacheKey := tableName + "_" + key
	cached, found := PricingCache.Get(cacheKey)

	if found {
		result := cached.(float64)
		//if result > 0 {
		return result
		//}
	} else {
		fmt.Println("Cache missed: ", key)
	}
	query := columnName + " = ?"
	row := DB.Table(tableName).Where(query, key).Select("qx").Row()

	var qx float64
	err := row.Scan(&qx)
	if err != nil {
		fmt.Println("db error: ", errors.WithStack(err))
	}
	PricingCache.Set(cacheKey, qx, 1)
	return qx

}

func getPricingLapseRate(arg models.TransitionRateArguments, columnName, tableName string) float64 {
	//tableName := strings.ToLower(arg.ProductCode) + "_pricing_" + strings.ToLower(pricingLapseTableName)
	//if arg.DurationIfM > LapseMaxMonthDimension {
	//	arg.DurationIfM = LapseMaxMonthDimension
	//}

	if arg.DurationIfM > pricingTables[arg.ProductCode].LapseMonthCount {
		arg.DurationIfM = pricingTables[arg.ProductCode].LapseMonthCount
	}

	key := PricingkeyBuilder(arg, Lapse)
	//columnName := "duration_if_m" // GetColumnName(tableName)

	//if month == 0 {
	//	month = 1
	//}
	//
	//if month > pricingLapseMonthCount {
	//	month = pricingLapseMonthCount
	//}
	cacheKey := tableName + "_" + key

	cached, found := PricingCache.Get(cacheKey)

	if found {
		result := cached.(float64)
		//if result > 0 {
		return result
		//}
	}

	var lapseRate float64
	query := columnName + " = ?"
	row := DB.Table(tableName).Where(query, key).Select("lapse_rate").Row()
	err := row.Scan(&lapseRate)
	if err != nil {
		fmt.Println("db error: ", errors.WithStack(err))
	}
	PricingCache.Set(cacheKey, lapseRate, 1)

	return lapseRate
}

func getPricingLapseMargin(month int, prodCode string, basis string, LapseMarginMonthCount int) float64 {
	var margin models.ProductPricingLapseMargin
	if month == 0 { //A hack to fix the month zero issue with lapses. Discuss with M later.
		month = 1
	}

	if month > LapseMarginMonthCount {
		month = LapseMarginMonthCount
	}

	key := "lp-" + prodCode + "_" + basis + "_" + "-" + strconv.Itoa(month)
	result, found := PricingCache.Get(key)
	if found {
		return result.(float64)
	} else {
		err := DB.Where("product_code = ?  and month = ? and basis=?", prodCode, month, basis).First(&margin).Error
		if err != nil {
			PricingCache.Set(key, 0.0, 1)
			return 0.0

		}
		PricingCache.Set(key, margin.Margin, 1)
	}

	return margin.Margin
}

func getPricingRetrenchmentRate(arg models.TransitionRateArguments, columnname string, tableName string) models.PricingRetrenchmentRate {
	var retrenchmentRate models.PricingRetrenchmentRate
	if arg.ProjectionMonth > pricingTables[arg.ProductCode].RetrenchmentTableRowCount {
		arg.ProjectionMonth = pricingTables[arg.ProductCode].RetrenchmentTableRowCount
	}
	//if yearInView == 0 {
	//	yearInView = 1
	//}
	//
	//if yearInView > 5 {
	//	yearInView = 5
	//}

	//key := strings.ToLower(arg.ProductCode) + "-retr-" + strconv.Itoa(yearInView)
	key := PricingkeyBuilder(arg, Retrenchment)
	cacheKey := tableName + "_" + key
	cached, found := PricingCache.Get(cacheKey)
	if found {
		return cached.(models.PricingRetrenchmentRate)
	}

	query := columnname + "=?"
	err := DB.Table(tableName).Where(query, key).Find(&retrenchmentRate).Error
	if err != nil {
		fmt.Println(err)
	}
	PricingCache.Set(key, retrenchmentRate, 1)
	return retrenchmentRate
}

func getPricingDisabilityIncidenceRate(arg models.TransitionRateArguments, columnname, tableName string) float64 {
	//Build the key
	key := PricingkeyBuilder(arg, PermanentDisability)

	cacheKey := tableName + "_" + key
	cached, found := PricingCache.Get(cacheKey)

	if found {
		result := cached.(float64)
		//if result > 0 {
		return result
		//}
	}

	var incidenceRate float64
	if DisabilityColumnName == "" {
		DisabilityColumnName = columnname //GetColumnName(tableName)
	}

	query := DisabilityColumnName + " = ?"
	row := DB.Table(tableName).Where(query, key).Select("incidence_rate").Row()
	err := row.Scan(&incidenceRate)
	if err != nil {
		fmt.Println("db error: ", errors.WithStack(err))
	}
	PricingCache.Set(cacheKey, incidenceRate, 1)
	return incidenceRate
}

func getPricingChildAdditionalSumAssured(productCode string, plan string, age int) float64 {
	key := fmt.Sprintf("%s-%s-%d", strings.ToLower(productCode), plan, age)
	cached, found := PricingCache.Get(key)
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
	case "G":
		result = childFuneralService.G
	case "H":
		result = childFuneralService.H
	case "I":
		result = childFuneralService.I
	case "J":
		result = childFuneralService.J
	case "K":
		result = childFuneralService.K
	case "L":
		result = childFuneralService.L
	case "M":
		result = childFuneralService.M
	case "N":
		result = childFuneralService.N
	case "O":
		result = childFuneralService.O
	default:
		result = 0
	}
	PricingCache.Set(key, result, 1)

	return result
}

func getPricingAdditionalSumAssured(productCode string, plan string) float64 {
	key := fmt.Sprintf("%s-%s", strings.ToLower(productCode), plan)
	cached, found := PricingCache.Get(key)
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
	PricingCache.Set(key, result, 1)

	return result
}

func getPricingChildSumAssured(productCode string, plan string, age int) float64 {
	key := fmt.Sprintf("%s-%s-%d", strings.ToLower(productCode), plan, age)
	cached, found := PricingCache.Get(key)
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
	PricingCache.Set(key, result, 1)

	return result
}

func getPricingClawback(month int, mp models.ProductPricingModelPoint) models.ProductPricingClawback {
	var clawback models.ProductPricingClawback
	key := fmt.Sprintf("clawback_%d_productcode_%s", month, mp.ProductCode)
	cached, found := PricingCache.Get(key)
	if found {
		clawback = cached.(models.ProductPricingClawback)
	} else {
		err := DB.Where("duration_in_force_month = ? and product_code = ?", month, mp.ProductCode).First(&clawback).Error
		if err == nil {
			PricingCache.Set(key, clawback, 1)
		}
	}
	return clawback
	DB.Where("duration_in_force_month = ? and product_code=?", month, mp.ProductCode).First(&clawback)
	return clawback
}

func getPricingRiders(mp models.ProductPricingModelPoint, features models.ProductFeatures, pricingConfig models.PricingConfig, parameters models.ProductPricingParameters) float64 {
	var riderBenefitSA = 0.0
	//var prodPricingParam models.ProductPricingParameters

	if features.RiderBenefit {
		if pricingConfig.Repatriation {
			riderBenefitSA += GetPricingRiderValue("Repatriation", mp.ProductCode, mp.Plan) * parameters.RepatriationAssumption
		}
		if pricingConfig.Tombstone {
			riderBenefitSA += GetPricingRiderValue("Tombstone", mp.ProductCode, mp.Plan)
		}
		if pricingConfig.Cow {
			riderBenefitSA += GetPricingRiderValue("Cow", mp.ProductCode, mp.Plan)
		}
		if pricingConfig.Grocery {
			riderBenefitSA += GetPricingRiderValue("Grocery", mp.ProductCode, mp.Plan) * float64(parameters.GroceryBenefitMultiplier)
		}
	}
	return riderBenefitSA
}

func getProfitMargin(anb int, SA int, productCode string) float64 {
	var profitMargin models.ProductPricingProfitMargin
	key := "profit-margin" + strconv.Itoa(SA) + strconv.Itoa(anb) + productCode
	cached, found := PricingCache.Get(key)

	if found {
		result := cached.(float64)
		return result
	}
	err := DB.Table("product_pricing_profit_margins").Where("anb = ? and sum_assured = ? and product_code= ?", anb, SA, productCode).First(&profitMargin).Error
	if err != nil {
		fmt.Println("db error: ", errors.WithStack(err))
		profitMargin.ProfitMargin = 0
	}
	success := PricingCache.Set(key, profitMargin.ProfitMargin, 1)
	if !success {
		fmt.Println("Cache error: Not saved")
	}

	return profitMargin.ProfitMargin
}

func getEscalationRate(anb int, productCode, escalationVariable, basis string) float64 {
	var escalation models.ProductPricingProductLevelEscalation
	var rate float64
	key := "escalations" + escalationVariable + strconv.Itoa(anb) + productCode + basis
	cached, found := PricingCache.Get(key)

	if found {
		result := cached.(float64)
		return result
	}
	err := DB.Table("product_pricing_product_level_escalations").Where("anb = ? and product_code= ? and basis=?", anb, productCode, basis).First(&escalation).Error
	if err != nil {
		fmt.Println("db error: ", errors.WithStack(err))
		escalation.PremiumEscalationRate = 0
		escalation.SumAssuredEscalationRate = 0
	}

	switch escalationVariable {
	case "SumAssuredEscalationRate":
		rate = escalation.SumAssuredEscalationRate
	case "PremiumEscalationRate":
		rate = escalation.PremiumEscalationRate
	default:
		rate = 0
	}

	success := PricingCache.Set(key, rate, 1)
	if !success {
		fmt.Println("Cache error: Not saved")
	}

	return rate
}

func UpdatePricingAggregatedProjections(pricingPoints []models.PricingPoint, aggpricingProjs *map[string]models.AggregatedPricingPoint, weighting float64) {
	for i, pricingPoint := range pricingPoints {
		key := strconv.Itoa(i) + "_" + pricingPoint.ProductCode + "_" + strconv.Itoa(pricingPoint.JobProductID) + "_" + strconv.Itoa(pricingPoint.ScenarioId)
		var aggProj = models.AggregatedPricingPoint{}
		err := copier.Copy(&aggProj, &pricingPoint)
		aggProj.ID = 0

		if err != nil {
			fmt.Println(err)
		}
		mutex.Lock()
		agg, exists := (*aggpricingProjs)[key]

		if exists {
			mutable := reflect.ValueOf(&agg).Elem()
			mutable2 := reflect.ValueOf(&aggProj).Elem()

			for j := 0; j < mutable.NumField(); j++ {
				if j == 0 {
					continue
				} else {
					if mutable.Field(j).Type().Kind() == reflect.Float64 {
						mutable.Field(j).SetFloat(math.Round(mutable.Field(j).Float()*100)/100 + math.Round(mutable2.Field(j).Float()*weighting*100)/100)
					}
				}
			}
			(*aggpricingProjs)[key] = agg
		} else {
			(*aggpricingProjs)[key] = aggProj
		}

		mutex.Unlock()
	}
}

//	func
//	getPricingChildAdditionalSumAssured(code interface {}) {
//
//}
