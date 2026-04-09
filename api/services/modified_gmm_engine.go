package services

import (
	"api/models"
	"api/utils"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/ristretto"
	"github.com/gammazero/workerpool"
	"github.com/jinzhu/copier"
	"github.com/jszwec/csvutil"
	"github.com/rs/zerolog/log"
	"io"
	"math"
	"mime/multipart"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var PaaCache *ristretto.Cache

func RunGMMProjectionGroups(runJobGroups []models.MgmmRunPayload, user models.AppUser) {
	for k := range runJobGroups {
		var mgmmRun = models.MgmmRun{}
		mgmmRun.CreationDate = time.Now()
		mgmmRun.UserName = user.UserName
		mgmmRun.UserEmail = user.UserEmail
		mgmmRun.Name = runJobGroups[k].Name
		mgmmRun.RunDate = runJobGroups[k].RunDate
		mgmmRun.Ifrs17Ready = true
		mgmmRun.ProcessingStatus = "queued"
		DB.Create(&mgmmRun)
		for i := range runJobGroups[k].Portfolios {
			runJobGroups[k].Portfolios[i].MgmmRunID = mgmmRun.ID
			runJobGroups[k].Portfolios[i].Name = runJobGroups[k].Portfolios[i].Name + "_" + runJobGroups[k].Portfolios[i].PortfolioName
			runJobGroups[k].Portfolios[i].CreationDate = time.Now()
			runJobGroups[k].Portfolios[i].UserName = user.UserName
			runJobGroups[k].Portfolios[i].UserEmail = user.UserEmail
			runJobGroups[k].Portfolios[i].ProcessingStatus = "queued"
			err := DB.Create(&runJobGroups[k].Portfolios[i]).Error
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	for _, runJobGroup := range runJobGroups {
		RunGMMProjections(runJobGroup.Portfolios, user)
	}
}

func RunGMMProjections(runs []models.GMMRunSetting, user models.AppUser) {

	//Keep polling here until there are no jobs in the queue
	for {
		var count int64
		DB.Model(&models.GMMRunSetting{}).Where("processing_status = ?", "processing").Count(&count)
		fmt.Println("polling for jobs to finish")
		if count == 0 {
			fmt.Println("no jobs in progress")
			break
		}
		time.Sleep(3 * time.Second)
	}

	var mgmmRun models.MgmmRun
	DB.Where("id=?", runs[0].MgmmRunID).Find(&mgmmRun)

	for _, run := range runs {
		aggModPAA := make(map[string]models.AggregatedModifiedGMMProjection)
		scopedModPAA := make(map[string]models.ModifiedGMMScopedAggregation)

		run.ProcessingStatus = "processing"
		mgmmRun.ProcessingStatus = "processing"
		err := DB.Save(&run).Error
		if err != nil {
			fmt.Println(err)
		}
		err = RunGMMProjection(run, &aggModPAA, &scopedModPAA, mgmmRun)
		if err != nil {
			run.ProcessingStatus = "failed"
			mgmmRun.ProcessingStatus = "failed"
			run.FailureReason = err.Error()
			//var mgmmRun models.MgmmRun
			//DB.Where("id=?", run.MgmmRunID).Find(&mgmmRun)
			mgmmRun.Ifrs17Ready = false
			DB.Save(&run)
			DB.Save(&mgmmRun)
			break
		} else {
			var modifiedGMMScopedAggregations []models.ModifiedGMMScopedAggregation
			var aggregatedModifiedGMMProjections []models.AggregatedModifiedGMMProjection

			for _, appProj := range aggModPAA {
				aggregatedModifiedGMMProjections = append(aggregatedModifiedGMMProjections, appProj)
			}
			errSave := DB.CreateInBatches(aggregatedModifiedGMMProjections, 100).Error
			if errSave != nil {
				fmt.Println(errSave)
			}
			//DB.Where("run_date=?", run.RunDate).Delete(&models.ModifiedGMMScopedAggregation{})
			if run.IFRS17Aggregation {
				for _, scopedAgg := range scopedModPAA {
					modifiedGMMScopedAggregations = append(modifiedGMMScopedAggregations, scopedAgg)
				}
				errSave = DB.CreateInBatches(modifiedGMMScopedAggregations, 100).Error
				if errSave != nil {
					fmt.Println(errSave)
				}
			}

			fmt.Println("agg length::::::", len(aggregatedModifiedGMMProjections))
			fmt.Println("scoped length::::::", len(modifiedGMMScopedAggregations))

		}
	}
	PaaCache.Clear()
}

func RunGMMProjection(runSettings models.GMMRunSetting, aggModPAA *map[string]models.AggregatedModifiedGMMProjection, scopedModPAA *map[string]models.ModifiedGMMScopedAggregation, mgmmRun models.MgmmRun) error {
	PaaCache.Clear()

	startTime := time.Now()
	var gmmMps []models.ModifiedGMMModelPoint
	var portfolio models.PaaPortfolio
	var errRun error
	var sumRiskUnits float64
	//var totalExpectedPremiumFromInception float64
	var respLockedin float64
	var respCurrent float64
	var valuationDate, prevValuationDate time.Time
	var prevValuationYear int

	//var initialRecogCSM float64
	if runSettings.RunSingle {
		errRun = DB.Where("year = ? and mp_version=? and paa_portfolio_id = ?", runSettings.ModelPoint, runSettings.ModelPointVersion, runSettings.PortfolioId).First(&gmmMps).Error
	} else {
		errRun = DB.Where("year = ? and mp_version=? and paa_portfolio_id = ?", runSettings.ModelPoint, runSettings.ModelPointVersion, runSettings.PortfolioId).Find(&gmmMps).Error
	}

	if errRun != nil {
		return errRun
	}

	errRun = DB.Where("name = ?", runSettings.PortfolioName).Find(&portfolio).Error

	if errRun != nil {
		return errRun
	}
	runYear, err := strconv.Atoi(runSettings.RunDate[:4])
	if err != nil {
		runYear = 0
	}

	runMonth, err := strconv.Atoi(runSettings.RunDate[5:])
	if err != nil {
		runMonth = 0
	}

	valuationDate, err = utils.ParseDateString(runSettings.RunDate + "-01")
	valuationDate = endOfMonth(valuationDate, 0)
	if int(valuationDate.Month()) <= runSettings.YearEndMonth {
		prevValuationYear = valuationDate.Year() - 1
	} else {
		prevValuationYear = valuationDate.Year()
	}
	prevValuationDate, err = utils.ParseDateString(strconv.Itoa(prevValuationYear) + "-" + strconv.Itoa(runSettings.YearEndMonth) + "-01")
	prevValuationDate = endOfMonth(prevValuationDate, 0)

	yieldCurveCode := getModifiedGMMParameterString(gmmMps[0].ProductCode, gmmMps[0].SubProductCode, runSettings.ParameterYear, "YieldCurveCode")

	if portfolio.DiscountOption == Discounted {
		respLockedin = GetPaaForwardRate(1, gmmMps[0].LockedInYear, gmmMps[0].LockedInMonth, yieldCurveCode)
		respCurrent = GetPaaForwardRate(1, runYear, runMonth, yieldCurveCode)
	}

	if portfolio.DiscountOption == Discounted && respLockedin == 0 {
		return fmt.Errorf("no valid yield curve data was found for yield curve code, %s, lockedinyear, %d, and lockedinmonth, %d,with discount option selected", yieldCurveCode, gmmMps[0].LockedInYear, gmmMps[0].LockedInMonth)
	}

	if respCurrent == 0 && portfolio.DiscountOption == Discounted {
		return fmt.Errorf("no valid yield curve data was found for yield curve code, %s, run year, %d, and month, %d", yieldCurveCode, runYear, runMonth)
	}

	runSettings.TotalRecords = len(gmmMps)
	runSettings.ProcessingStatus = "processing"
	runSettings.CreationDate = startTime
	mgmmRun.ProcessingStatus = "processing"
	mgmmRun.TotalRecords += runSettings.TotalRecords
	DB.Save(&runSettings)
	DB.Save(&mgmmRun)

	wp2 := workerpool.New(runtime.NumCPU())
	var processedCount int
	for _, mp := range gmmMps {
		mp := mp
		wp2.Submit(func() {
			startLoopTime := time.Now()
			var mgpls []models.AggregatedModifiedGMMProjection
			var prevInforcePolicyCount float64
			var IrprevInforcePolicyCount float64 //Initial Recognition variable
			var shock models.ModifiedGMMShock
			var yieldCurveCode string
			var totalExpectedPremiumFromInception float64
			var commDate, endCoverDate time.Time

			if portfolio.PremiumEarningPattern == SpecifiedbyUser {
				sumRiskUnits = getsumRiskUnit(runSettings.ParameterYear, mp.ProductCode, mp.SubProductCode, int(utils.RoundUp(mp.TermMonths)))
			} else {
				sumRiskUnits = 0
			}

			if portfolio.PremiumEarningPattern == DailyPassageofTime {
				commDate, err = utils.ParseDateString(mp.CommencementDate)
				endCoverDate, err = utils.ParseDateString(mp.CoverEndDate)
				tempterm := endCoverDate.Sub(commDate).Hours() / 24
				mp.TermMonths = float64(int(tempterm / 30))
			}

			if mp.Frequency == 0 {
				totalExpectedPremiumFromInception = mp.AnnualPremium
			}
			if mp.Frequency != 0 {
				totalExpectedPremiumFromInception = mp.AnnualPremium * float64(mp.TermMonths) / 12.0
			}

			yieldCurveCode = getModifiedGMMParameterString(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "YieldCurveCode")

			DacBuildupIndicator := getModifiedGMMParameterBoolean(gmmMps[0].ProductCode, gmmMps[0].SubProductCode, runSettings.ParameterYear, "DacBuildupIndicator")

			for i := 0; i <= int(utils.RoundUp(mp.TermMonths)); i++ { //runSettings.ProjectionPeriod; i++ {
				var projection models.AggregatedModifiedGMMProjection
				projection.ProjectionMonth = i
				projection.PortfolioName = runSettings.PortfolioName
				projection.YieldCurveCode = yieldCurveCode
				projection.PortfolioId = runSettings.PortfolioId
				projection.RunID = runSettings.MgmmRunID
				projection.JobRunId = runSettings.ID
				projection.RunName = runSettings.Name
				projection.RunDate = runSettings.RunDate[:7]

				projection.ValuationMonth = projection.ProjectionMonth + mp.DurationInForceMonths

				projection.ValuationDate = endOfMonth(valuationDate, projection.ProjectionMonth).String()

				if endOfMonth(valuationDate, projection.ProjectionMonth).After(endCoverDate) {
					projection.DaysCovered = 0
				} else {
					if projection.ProjectionMonth == 0 {
						projection.DaysCovered = int(valuationDate.Sub(maxDate(prevValuationDate, commDate)).Hours() / 24)
					}
					if projection.ProjectionMonth > 0 {
						currProjDate := endOfMonth(valuationDate, projection.ProjectionMonth)
						prevProjDate := endOfMonth(valuationDate, projection.ProjectionMonth-1)
						tempDuration := currProjDate.Sub(maxDate(commDate, prevProjDate))
						projection.DaysCovered = int(tempDuration.Hours() / 24)
					}
				}

				projection.IrValuationMonth = projection.ProjectionMonth

				projection.LapseRate = getPAALapseRate(runSettings.ParameterYear, mp.ProductCode, projection.ValuationMonth, mp.DistributionChannel)
				if mp.Status == NB {
					projection.IrLapseRate = getPAALapseRate(runSettings.ParameterYear, mp.ProductCode, projection.IrValuationMonth, mp.DistributionChannel)
				}
				if projection.ProjectionMonth == 0 {
					projection.InforcePolicyCountSM = 1
					projection.InforcePolicyCount = 1
					projection.IrInforcePolicyCountSM = 1
				}
				if projection.ProjectionMonth == 1 {
					projection.InforcePolicyCountSM = 1
					projection.InforcePolicyCount = projection.InforcePolicyCountSM * math.Pow(float64(1-projection.LapseRate), 1/12.0)
					projection.IrInforcePolicyCountSM = 1
					IrprevInforcePolicyCount = projection.InforcePolicyCountSM * math.Pow(float64(1-projection.IrLapseRate), 1/12.0)
				}
				if projection.ProjectionMonth > 1 {
					projection.InforcePolicyCountSM = prevInforcePolicyCount
					projection.InforcePolicyCount = projection.InforcePolicyCountSM * math.Pow(float64(1-projection.LapseRate), 1/12.0)
					projection.IrInforcePolicyCountSM = IrprevInforcePolicyCount
					IrprevInforcePolicyCount = projection.IrInforcePolicyCountSM * math.Pow(float64(1-projection.IrLapseRate), 1/12.0)
				}
				prevInforcePolicyCount = projection.InforcePolicyCount
				projection.Year = mp.Year
				projection.ProductCode = mp.ProductCode
				projection.IFRS17Group = mp.IFRS17Group
				projection.Treaty1IFRS17Group = mp.IFRS17GroupTreaty1
				projection.Treaty2IFRS17Group = mp.IFRS17GroupTreaty2
				projection.Treaty3IFRS17Group = mp.IFRS17GroupTreaty3
				projection.PolicyNumber = mp.PolicyNumber
				projection.LockedinYear = strconv.Itoa(mp.LockedInYear)
				projection.LockedinMonth = strconv.Itoa(mp.LockedInMonth)
				if projection.ProjectionMonth > 0 && runSettings.ShockSettings.ShockBasis != "" {
					shock, err = GetModifiedGMMShock(projection.ProjectionMonth, runSettings.ShockSettings.ShockBasis)
					if err != nil {
						log.Error().Err(err).Send()
					}
				}
				ModifiedGMMPremiumReceipt(&projection, mp, runSettings, shock)
				ModifiedGMMEarnedPremium(&projection, mp, portfolio, runSettings, sumRiskUnits, totalExpectedPremiumFromInception, shock, DacBuildupIndicator)
				if mp.AnnualInterestRate > 0 {
					projection.OutstandingLoan = CalculatePV(mp.AnnualInterestRate, math.Max(float64(mp.OutstandingLoanTermMonths-projection.ProjectionMonth+1), 0), mp.MonthlyInstalment)
				}

				ModifiedGMMCashOutflow(&projection, mp, runSettings, shock)
				ModifiedGMMNetCashflow(&projection, mp)

				// Initial Recognition cashflows for the purpose of quantifying loss component
				if mp.Status == NB {
					InitialRecognitionModifiedGMMPremiumReceipt(&projection, mp, runSettings, shock)
					InitialRecognitionModifiedGMMEarnedPremium(&projection, mp, portfolio, runSettings, sumRiskUnits, totalExpectedPremiumFromInception, shock)
					InitialRecognitionModifiedGMMCashOutflow(&projection, mp, runSettings, shock)
					InitialRecognitionModifiedGMMNetCashflow(&projection, mp)
				}

				mgpls = append(mgpls, projection)
			}

			ModifiedGMMFutureValues := models.ModifiedDiscountedValues{
				EarnedPremium:                          0,
				SumFutureEarnedPremium:                 0,
				DiscountedEarnedPremiumCurrent:         0,
				DiscountedEarnedPremiumLockedin:        0,
				CashOutflow:                            0,
				SumFutureCashOutflows:                  0,
				DiscountedCashOutflows:                 0,
				DiscountedCashOutflowsLockedin:         0,
				NetCashFlow:                            0,
				SumFutureNetCashFlows:                  0,
				SumFutureAcquisitionCost:               0,
				DiscountedNetCashFlowsCurrent:          0,
				DiscountedNetCashFlowsLockedin:         0,
				DiscountedDACLockedin:                  0,
				RiskAdjustment:                         0,
				Treaty1EarnedPremium:                   0,
				Treaty1SumFutureEarnedPremium:          0,
				Treaty1DiscountedEarnedPremiumCurrent:  0,
				Treaty1DiscountedEarnedPremiumLockedin: 0,
				Treaty1CashOutflow:                     0,
				Treaty1SumFutureCashOutflows:           0,
				Treaty1DiscountedCashOutflows:          0,
				Treaty1DiscountedCashOutflowsLockedin:  0,
				Treaty1NetCashFlow:                     0,
				Treaty1SumFutureNetCashFlows:           0,
				Treaty1DiscountedNetCashFlowsCurrent:   0,
				Treaty1DiscountedNetCashFlowsLockedin:  0,
				Treaty1RiskAdjustment:                  0,
				Treaty2EarnedPremium:                   0,
				Treaty2CashOutflow:                     0,
				Treaty2NetCashFlow:                     0,
				Treaty2SumFutureEarnedPremium:          0,
				Treaty2DiscountedEarnedPremiumCurrent:  0,
				Treaty2DiscountedEarnedPremiumLockedin: 0,
				Treaty2SumFutureCashOutflows:           0,
				Treaty2DiscountedCashOutflows:          0,
				Treaty2DiscountedCashOutflowsLockedin:  0,
				Treaty2SumFutureNetCashFlows:           0,
				Treaty2DiscountedNetCashFlowsCurrent:   0,
				Treaty2DiscountedNetCashFlowsLockedin:  0,
				Treaty2RiskAdjustment:                  0,
				Treaty3EarnedPremium:                   0,
				Treaty3CashOutflow:                     0,
				Treaty3NetCashFlow:                     0,
				Treaty3SumFutureEarnedPremium:          0,
				Treaty3DiscountedEarnedPremiumCurrent:  0,
				Treaty3DiscountedEarnedPremiumLockedin: 0,
				Treaty3SumFutureCashOutflows:           0,
				Treaty3DiscountedCashOutflows:          0,
				Treaty3DiscountedCashOutflowsLockedin:  0,
				Treaty3SumFutureNetCashFlows:           0,
				Treaty3DiscountedNetCashFlowsCurrent:   0,
				Treaty3DiscountedNetCashFlowsLockedin:  0,
				Treaty3RiskAdjustment:                  0,
				CededEarnedPremium:                     0,
				CededCashOutflow:                       0,
				CededNetCashFlow:                       0,
				CededSumFutureEarnedPremium:            0,
				CededDiscountedEarnedPremiumCurrent:    0,
				CededSumFutureCashOutflows:             0,
				CededDiscountedCashOutflows:            0,
				CededDiscountedCashOutflowsLockedin:    0,
				CededSumFutureNetCashFlows:             0,
				CededDiscountedNetCashFlowsCurrent:     0,
				CededDiscountedNetCashFlowsLockedin:    0,
				CededRiskAdjustment:                    0,
				IrNbPremiumReceipt:                     0,
				IrClaimsOutgo:                          0,
				IrClaimsExpenseOutgo:                   0,
				IrInitialCommissionOutgo:               0,
				IrRenewalCommissionOutgo:               0,
				IrInitialExpenseOutgo:                  0,
				IrMaintenanceExpenseOutgo:              0,
				IrDiscountedEarnedPremium:              0,
				IrDiscountedNetCashFlowsLockedin:       0,
				IrEarnedPremium:                        0,
				IrDiscountedClaimsOutgoLockedin:        0,
				IrTreaty1DiscountedClaimsOutgo:         0,
				IrTreaty2DiscountedClaimsOutgo:         0,
				IrTreaty3DiscountedClaimsOutgo:         0,
			}

			for i := len(mgpls) - 1; i >= 0; i-- {
				CalculateModifiedGMMDiscountedValues(i, runYear, runMonth, yieldCurveCode, runSettings, portfolio, &mgpls[i], mp, ModifiedGMMFutureValues, shock)
				ModifiedGMMFutureValues.TotalPremiumReceipt = mgpls[i].TotalPremiumReceipt
				ModifiedGMMFutureValues.DiscountedTotalPremiumReceiptCurrent = mgpls[i].DiscountedTotalPremiumReceiptCurrent
				ModifiedGMMFutureValues.DiscountedTotalPremiumReceiptLockedin = mgpls[i].DiscountedTotalPremiumReceiptLockedin
				ModifiedGMMFutureValues.EarnedPremium = mgpls[i].EarnedPremium
				ModifiedGMMFutureValues.SumFutureEarnedPremium = mgpls[i].SumFutureEarnedPremium
				ModifiedGMMFutureValues.DiscountedEarnedPremiumCurrent = mgpls[i].DiscountedEarnedPremiumCurrent
				ModifiedGMMFutureValues.DiscountedEarnedPremiumLockedin = mgpls[i].DiscountedEarnedPremiumLockedin
				ModifiedGMMFutureValues.CashOutflow = mgpls[i].CashOutflow
				ModifiedGMMFutureValues.ClaimsOutgo = mgpls[i].ClaimsOutgo
				ModifiedGMMFutureValues.ClaimsExpenseOutgo = mgpls[i].ClaimsExpenseOutgo
				ModifiedGMMFutureValues.InitialCommissionOutgo = mgpls[i].InitialCommissionOutgo
				ModifiedGMMFutureValues.InitialExpenseOutgo = mgpls[i].InitialExpenseOutgo
				ModifiedGMMFutureValues.SumFutureCashOutflows = mgpls[i].SumFutureCashOutflows
				ModifiedGMMFutureValues.SumFutureAcquisitionCost = mgpls[i].SumFutureAcquisitionCost
				ModifiedGMMFutureValues.DiscountedCashOutflows = mgpls[i].DiscountedCashOutflows
				ModifiedGMMFutureValues.DiscountedCashOutflowsLockedin = mgpls[i].DiscountedCashOutflowsLockedin
				ModifiedGMMFutureValues.NetCashFlow = mgpls[i].NetCashFlow
				ModifiedGMMFutureValues.SumFutureNetCashFlows = mgpls[i].SumFutureNetCashFlows
				ModifiedGMMFutureValues.DiscountedNetCashFlowsCurrent = mgpls[i].DiscountedNetCashFlowsCurrent
				ModifiedGMMFutureValues.DiscountedNetCashFlowsLockedin = mgpls[i].DiscountedNetCashFlowsLockedin
				ModifiedGMMFutureValues.RiskAdjustment = mgpls[i].RiskAdjustment
				ModifiedGMMFutureValues.DiscountedClaimsOutgoLockedin = mgpls[i].DiscountedClaimsOutgoLockedin
				ModifiedGMMFutureValues.DiscountedDACLockedin = mgpls[i].DiscountedDACLockedin
				ModifiedGMMFutureValues.Treaty1PremiumReceipt = mgpls[i].Treaty1PremiumReceipt
				ModifiedGMMFutureValues.Treaty1EarnedPremium = mgpls[i].Treaty1EarnedPremium
				ModifiedGMMFutureValues.Treaty1SumFutureEarnedPremium = mgpls[i].Treaty1SumFutureEarnedPremium
				ModifiedGMMFutureValues.Treaty1DiscountedEarnedPremiumCurrent = mgpls[i].Treaty1DiscountedEarnedPremiumCurrent
				ModifiedGMMFutureValues.Treaty1DiscountedEarnedPremiumLockedin = mgpls[i].Treaty1DiscountedEarnedPremiumLockedin
				ModifiedGMMFutureValues.Treaty1CashOutflow = mgpls[i].Treaty1CashOutflow
				ModifiedGMMFutureValues.Treaty1SumFutureCashOutflows = mgpls[i].Treaty1SumFutureCashOutflows
				ModifiedGMMFutureValues.Treaty1DiscountedCashOutflows = mgpls[i].Treaty1DiscountedCashOutflows
				ModifiedGMMFutureValues.Treaty1DiscountedCashOutflowsLockedin = mgpls[i].Treaty1DiscountedCashOutflowsLockedin
				ModifiedGMMFutureValues.Treaty1NetCashFlow = mgpls[i].Treaty1NetCashFlow
				ModifiedGMMFutureValues.Treaty1SumFutureNetCashFlows = mgpls[i].Treaty1SumFutureNetCashFlows
				ModifiedGMMFutureValues.Treaty1DiscountedNetCashFlowsCurrent = mgpls[i].Treaty1DiscountedNetCashFlowsCurrent
				ModifiedGMMFutureValues.Treaty1DiscountedNetCashFlowsLockedin = mgpls[i].Treaty1DiscountedNetCashFlowsLockedin
				ModifiedGMMFutureValues.Treaty1RiskAdjustment = mgpls[i].Treaty1RiskAdjustment
				ModifiedGMMFutureValues.Treaty2PremiumReceipt = mgpls[i].Treaty2NbPremiumReceipt
				ModifiedGMMFutureValues.Treaty2EarnedPremium = mgpls[i].Treaty2EarnedPremium
				ModifiedGMMFutureValues.Treaty2CashOutflow = mgpls[i].Treaty2CashOutflow
				ModifiedGMMFutureValues.Treaty2NetCashFlow = mgpls[i].Treaty2NetCashFlow
				ModifiedGMMFutureValues.Treaty2SumFutureEarnedPremium = mgpls[i].Treaty2SumFutureEarnedPremium
				ModifiedGMMFutureValues.Treaty2DiscountedEarnedPremiumCurrent = mgpls[i].Treaty2DiscountedEarnedPremiumCurrent
				ModifiedGMMFutureValues.Treaty2DiscountedEarnedPremiumLockedin = mgpls[i].Treaty2DiscountedEarnedPremiumLockedin
				ModifiedGMMFutureValues.Treaty2SumFutureCashOutflows = mgpls[i].Treaty2SumFutureCashOutflows
				ModifiedGMMFutureValues.Treaty2DiscountedCashOutflows = mgpls[i].Treaty2DiscountedCashOutflows
				ModifiedGMMFutureValues.Treaty2DiscountedCashOutflowsLockedin = mgpls[i].Treaty2DiscountedCashOutflowsLockedin
				ModifiedGMMFutureValues.Treaty2SumFutureNetCashFlows = mgpls[i].Treaty2SumFutureNetCashFlows
				ModifiedGMMFutureValues.Treaty2DiscountedNetCashFlowsCurrent = mgpls[i].Treaty2DiscountedNetCashFlowsCurrent
				ModifiedGMMFutureValues.Treaty2DiscountedNetCashFlowsLockedin = mgpls[i].Treaty2DiscountedNetCashFlowsLockedin
				ModifiedGMMFutureValues.Treaty2RiskAdjustment = mgpls[i].Treaty2RiskAdjustment
				ModifiedGMMFutureValues.Treaty3PremiumReceipt = mgpls[i].Treaty3PremiumReceipt
				ModifiedGMMFutureValues.Treaty3EarnedPremium = mgpls[i].Treaty3EarnedPremium
				ModifiedGMMFutureValues.Treaty3CashOutflow = mgpls[i].Treaty3CashOutflow
				ModifiedGMMFutureValues.Treaty3NetCashFlow = mgpls[i].Treaty3NetCashFlow
				ModifiedGMMFutureValues.Treaty3SumFutureEarnedPremium = mgpls[i].Treaty3SumFutureEarnedPremium
				ModifiedGMMFutureValues.Treaty3DiscountedEarnedPremiumCurrent = mgpls[i].Treaty3DiscountedEarnedPremiumCurrent
				ModifiedGMMFutureValues.Treaty3DiscountedEarnedPremiumLockedin = mgpls[i].Treaty3DiscountedEarnedPremiumLockedin
				ModifiedGMMFutureValues.Treaty3SumFutureCashOutflows = mgpls[i].Treaty3SumFutureCashOutflows
				ModifiedGMMFutureValues.Treaty3DiscountedCashOutflows = mgpls[i].Treaty3DiscountedCashOutflows
				ModifiedGMMFutureValues.Treaty3DiscountedCashOutflowsLockedin = mgpls[i].Treaty3DiscountedCashOutflowsLockedin
				ModifiedGMMFutureValues.Treaty3SumFutureNetCashFlows = mgpls[i].Treaty3SumFutureNetCashFlows
				ModifiedGMMFutureValues.Treaty3DiscountedNetCashFlowsCurrent = mgpls[i].Treaty3DiscountedNetCashFlowsCurrent
				ModifiedGMMFutureValues.Treaty3DiscountedNetCashFlowsLockedin = mgpls[i].Treaty3DiscountedNetCashFlowsLockedin
				ModifiedGMMFutureValues.Treaty3RiskAdjustment = mgpls[i].Treaty3RiskAdjustment
				ModifiedGMMFutureValues.TotalCededPremiumReceipt = mgpls[i].TotalCededPremiumReceipt
				ModifiedGMMFutureValues.CededEarnedPremium = mgpls[i].CededEarnedPremium
				ModifiedGMMFutureValues.CededCashOutflow = mgpls[i].CededCashOutflow
				ModifiedGMMFutureValues.CededNetCashFlow = mgpls[i].CededNetCashFlow
				ModifiedGMMFutureValues.CededSumFutureEarnedPremium = mgpls[i].CededSumFutureEarnedPremium
				ModifiedGMMFutureValues.CededDiscountedEarnedPremiumCurrent = mgpls[i].CededDiscountedEarnedPremiumCurrent
				ModifiedGMMFutureValues.CededSumFutureCashOutflows = mgpls[i].CededSumFutureCashOutflows
				ModifiedGMMFutureValues.CededDiscountedCashOutflows = mgpls[i].CededDiscountedCashOutflows
				ModifiedGMMFutureValues.CededDiscountedCashOutflowsLockedin = mgpls[i].CededDiscountedCashOutflowsLockedin
				ModifiedGMMFutureValues.CededSumFutureNetCashFlows = mgpls[i].CededSumFutureNetCashFlows
				ModifiedGMMFutureValues.CededDiscountedNetCashFlowsCurrent = mgpls[i].CededDiscountedNetCashFlowsCurrent
				ModifiedGMMFutureValues.CededDiscountedNetCashFlowsLockedin = mgpls[i].CededDiscountedNetCashFlowsLockedin
				ModifiedGMMFutureValues.CededRiskAdjustment = mgpls[i].CededRiskAdjustment
				ModifiedGMMFutureValues.IrClaimsOutgo = mgpls[i].IrClaimsOutgo
				ModifiedGMMFutureValues.IrClaimsExpenseOutgo = mgpls[i].IrClaimsExpenseOutgo
				ModifiedGMMFutureValues.IrInitialExpenseOutgo = mgpls[i].IrInitialExpenseOutgo
				ModifiedGMMFutureValues.IrMaintenanceExpenseOutgo = mgpls[i].IrMaintenanceExpenseOutgo
				ModifiedGMMFutureValues.IrInitialCommissionOutgo = mgpls[i].IrInitialCommissionOutgo
				ModifiedGMMFutureValues.IrRenewalCommissionOutgo = mgpls[i].IrRenewalCommissionOutgo
				ModifiedGMMFutureValues.IrNbPremiumReceipt = mgpls[i].IrNbPremiumReceipt
				ModifiedGMMFutureValues.IrDiscountedNetCashFlowsLockedin = mgpls[i].IrDiscountedNetCashFlowsLockedin
				ModifiedGMMFutureValues.IrDiscountedCashOutflowsLockedin = mgpls[i].IrDiscountedCashOutflowsLockedin
				ModifiedGMMFutureValues.IrDiscountedEarnedPremium = mgpls[i].IrDiscountedEarnedPremium
				ModifiedGMMFutureValues.IrEarnedPremium = mgpls[i].IrEarnedPremium
				ModifiedGMMFutureValues.IrDiscountedClaimsOutgoLockedin = mgpls[i].IrDiscountedClaimsOutgoLockedin
				ModifiedGMMFutureValues.IrTreaty1ClaimsOutgo = mgpls[i].IrTreaty1ClaimsOutgo
				ModifiedGMMFutureValues.IrTreaty2ClaimsOutgo = mgpls[i].IrTreaty2ClaimsOutgo
				ModifiedGMMFutureValues.IrTreaty3ClaimsOutgo = mgpls[i].IrTreaty3ClaimsOutgo
				ModifiedGMMFutureValues.IrTreaty1DiscountedClaimsOutgo = mgpls[i].IrTreaty1DiscountedClaimsOutgo
				ModifiedGMMFutureValues.IrTreaty2DiscountedClaimsOutgo = mgpls[i].IrTreaty2DiscountedClaimsOutgo
				ModifiedGMMFutureValues.IrTreaty3DiscountedClaimsOutgo = mgpls[i].IrTreaty3DiscountedClaimsOutgo
			}

			if runSettings.IndividualResults {
				err = DB.CreateInBatches(&mgpls, 100).Error
				if err != nil {
					fmt.Println(err)
					//return err
				}
			} else {
				UpdateAggregatedModifiedGMMProjection(mgpls[:int(utils.RoundUp(mp.TermMonths))+1], aggModPAA)

				//Add balance sheet records. Unearned Premium Reserve to the Balance Sheet Records

			}

			//Do scoped aggregations here...
			UpdateScopedModifiedGMMProjections(mgpls[:int(utils.RoundUp(mp.TermMonths))+1], scopedModPAA)

			processedCount += 1

			runSettings.ProcessedRecords = processedCount

			DB.Save(&runSettings)

			endLoopTime := time.Since(startLoopTime)
			fmt.Println("End Loop Time: ", endLoopTime.Seconds())

		})
	}

	wp2.StopWait()

	fmt.Println("End of Projection Loop", processedCount, runSettings.ProcessedRecords)
	if runSettings.TotalRecords == processedCount {
		runSettings.ProcessedRecords = processedCount
		runSettings.ProcessingStatus = "completed"
	}

	endTime := time.Since(startTime)
	runSettings.RunTime = endTime.Seconds()
	mgmmRun.ProcessingStatus = "completed"
	mgmmRun.RunTime += runSettings.RunTime
	DB.Save(runSettings)
	DB.Save(&mgmmRun)

	return nil
}

func getPAALapseRate(year int, prodCode string, durationInForce int, distributionChannel string) float64 {
	var paal models.PAALapse
	key := strconv.Itoa(year) + "_" + prodCode + "_" + strconv.Itoa(durationInForce) + "_" + distributionChannel
	cacheKey := key
	cached, found := PaaCache.Get(cacheKey)
	if found {
		return cached.(float64)
	}
	if found {
		result := cached.(float64)
		//if result > 0 {
		return result
		//}
	}
	if !found {
		DB.Where("year = ? and product_code = ? and duration_in_force_month = ? and distribution_channel = ?", year, prodCode, durationInForce, distributionChannel).Find(&paal)
		success := PaaCache.Set(cacheKey, paal.LapseRate, 1)
		if !success {
			fmt.Println("cache: key not stored")
		}
	}
	return paal.LapseRate
}

func getsumRiskUnit(year int, prodCode string, subProductCode string, termMonths int) float64 {

	keyparam := "gmm-" + strings.ToLower(prodCode) + "-" + subProductCode + "-" + strconv.Itoa(year)
	var param models.ModifiedGMMParameter
	result, found := PaaCache.Get(keyparam)

	if found {
		param = result.(models.ModifiedGMMParameter)
	} else {
		err := DB.Where("product_code = ? and sub_product_code=? and year = ?", prodCode, subProductCode, year).Find(&param).Error
		if err != nil {
			return 0
		}
		PaaCache.Set(keyparam, param, 1)
	}

	var sumRiskUnit float64
	key := strconv.Itoa(year) + "_" + param.PremiumEarningPatternCode + "_" + strconv.Itoa(termMonths) + "_sumrisk"
	cacheKey := key
	cached, found := PaaCache.Get(cacheKey)
	if found {
		return cached.(float64)
	}
	if found {
		result := cached.(float64)
		//if result > 0 {
		return result
		//}
	}
	if !found {

		query := fmt.Sprintf("SELECT sum(risk_unit) as sumRiskUnit from premium_earning_patterns where year='%d' and premium_earning_pattern_code= '%s' and duration_in_force<='%d'", year, param.PremiumEarningPatternCode, termMonths)
		DB.Raw(query).Find(&sumRiskUnit)
		success := PaaCache.Set(cacheKey, sumRiskUnit, 1)
		if !success {
			fmt.Println("cache: key not stored")
		}
	}
	return sumRiskUnit
}

func getRiskUnit(year int, prodCode string, subProductCode string, durationInForce int) float64 {

	keyparam := "gmm-" + strings.ToLower(prodCode) + "-" + subProductCode + "-" + strconv.Itoa(year)
	var param models.ModifiedGMMParameter
	result, found := PaaCache.Get(keyparam)

	if found {
		param = result.(models.ModifiedGMMParameter)
	} else {
		err := DB.Where("product_code = ? and sub_product_code=? and year = ?", prodCode, subProductCode, year).Find(&param).Error
		if err != nil {
			return 0
		}
		PaaCache.Set(keyparam, param, 1)
	}

	var premiumEarningPattern models.PremiumEarningPattern
	key := strconv.Itoa(year) + "_" + param.PremiumEarningPatternCode + "_" + strconv.Itoa(durationInForce)
	cacheKey := key
	cached, found := PaaCache.Get(cacheKey)
	if found {
		return cached.(float64)
	}
	//if found {
	//	result := cached.(float64)
	//	if result > 0 {
	//		return result
	//	}
	//}
	if !found {
		DB.Where("year = ? and premium_earning_pattern_code = ? and duration_in_force = ?", year, param.PremiumEarningPatternCode, durationInForce).Find(&premiumEarningPattern)
		success := PaaCache.Set(cacheKey, premiumEarningPattern.RiskUnit, 1)
		if !success {
			fmt.Println("cache: key not stored")
		}
	}
	return premiumEarningPattern.RiskUnit
}

func CalculateModifiedGMMDiscountedValues(index int, runYear, runMonth int, yieldCurveCode string, run models.GMMRunSetting, portfolio models.PaaPortfolio, mgp *models.AggregatedModifiedGMMProjection, mp models.ModifiedGMMModelPoint, modifiedFutureValues models.ModifiedDiscountedValues, shock models.ModifiedGMMShock) {
	var forwardRateCurrent float64
	var forwardRateLockedin float64
	var respCurrent float64 = 0
	var respLockedin float64 = 0

	//if portfolio.DiscountOption == Discounted {

	if portfolio.DiscountOption == Discounted {
		respLockedin = GetPaaForwardRate(mgp.ValuationMonth+1, mp.LockedInYear, mp.LockedInMonth, yieldCurveCode)
		respCurrent = GetPaaForwardRate(index+1, runYear, runMonth, yieldCurveCode)
	}

	if run.ShockSettings.DiscountCurve {
		forwardRateCurrent = respCurrent*(1+shock.MultiplicativeDiscountCurve) + shock.AdditiveDiscountCurve
		forwardRateLockedin = respLockedin*(1+shock.MultiplicativeDiscountCurve) + shock.AdditiveDiscountCurve
	} else {
		forwardRateCurrent = respCurrent
		forwardRateLockedin = respLockedin
	}

	mgp.CurrentRate = forwardRateCurrent
	mgp.LockedInRate = forwardRateLockedin
	mgp.Treaty1CurrentDiscount = mgp.CurrentRate
	mgp.Treaty1LockedInDiscount = mgp.LockedInRate
	mgp.Treaty2CurrentDiscount = mgp.CurrentRate
	mgp.Treaty2LockedInDiscount = mgp.LockedInRate
	mgp.Treaty3CurrentDiscount = mgp.CurrentRate
	mgp.Treaty3LockedInDiscount = mgp.LockedInRate

	discountFactorCurrent := math.Pow(1.0+forwardRateCurrent, -1/12.0)
	discountFactorLockedin := math.Pow(1.0+forwardRateLockedin, -1/12.0)

	mgp.DiscountedTotalPremiumReceiptCurrent = modifiedFutureValues.DiscountedTotalPremiumReceiptCurrent*discountFactorCurrent + modifiedFutureValues.TotalPremiumReceipt
	mgp.DiscountedTotalPremiumReceiptLockedin = modifiedFutureValues.DiscountedTotalPremiumReceiptLockedin*discountFactorLockedin + modifiedFutureValues.TotalPremiumReceipt

	mgp.SumFutureEarnedPremium = modifiedFutureValues.SumFutureEarnedPremium + modifiedFutureValues.EarnedPremium
	mgp.SumFutureAcquisitionCost = modifiedFutureValues.SumFutureAcquisitionCost + modifiedFutureValues.InitialCommissionOutgo + modifiedFutureValues.InitialExpenseOutgo
	mgp.DiscountedEarnedPremiumCurrent = (modifiedFutureValues.DiscountedEarnedPremiumCurrent + modifiedFutureValues.EarnedPremium) * discountFactorCurrent
	mgp.DiscountedEarnedPremiumLockedin = (modifiedFutureValues.DiscountedEarnedPremiumLockedin + modifiedFutureValues.EarnedPremium) * discountFactorLockedin

	mgp.SumFutureCashOutflows = modifiedFutureValues.SumFutureCashOutflows + modifiedFutureValues.CashOutflow
	mgp.DiscountedCashOutflows = (modifiedFutureValues.DiscountedCashOutflows+modifiedFutureValues.ClaimsOutgo+modifiedFutureValues.ClaimsExpenseOutgo)*discountFactorCurrent + (modifiedFutureValues.CashOutflow - modifiedFutureValues.ClaimsOutgo - modifiedFutureValues.ClaimsExpenseOutgo)
	mgp.DiscountedCashOutflowsLockedin = (modifiedFutureValues.DiscountedCashOutflowsLockedin+modifiedFutureValues.ClaimsOutgo+modifiedFutureValues.ClaimsExpenseOutgo)*discountFactorLockedin + (modifiedFutureValues.CashOutflow - modifiedFutureValues.ClaimsOutgo - modifiedFutureValues.ClaimsExpenseOutgo)
	mgp.DiscountedClaimsOutgoLockedin = (modifiedFutureValues.DiscountedClaimsOutgoLockedin + modifiedFutureValues.ClaimsOutgo) * discountFactorLockedin

	mgp.IrDiscountedClaimsOutgoLockedin = (modifiedFutureValues.IrDiscountedClaimsOutgoLockedin + modifiedFutureValues.IrClaimsOutgo) * discountFactorLockedin

	mgp.DiscountedDACLockedin = (modifiedFutureValues.DiscountedDACLockedin + modifiedFutureValues.InitialCommissionOutgo + modifiedFutureValues.InitialExpenseOutgo) * discountFactorLockedin

	mgp.SumFutureNetCashFlows = modifiedFutureValues.SumFutureNetCashFlows + modifiedFutureValues.NetCashFlow
	mgp.DiscountedNetCashFlowsCurrent = mgp.DiscountedCashOutflows - mgp.DiscountedTotalPremiumReceiptCurrent
	mgp.DiscountedNetCashFlowsLockedin = mgp.DiscountedCashOutflowsLockedin - mgp.DiscountedTotalPremiumReceiptLockedin

	mgp.IrDiscountedNetCashFlowsLockedin = (modifiedFutureValues.IrDiscountedNetCashFlowsLockedin+modifiedFutureValues.IrClaimsOutgo+modifiedFutureValues.IrClaimsExpenseOutgo)*discountFactorLockedin +
		modifiedFutureValues.IrInitialCommissionOutgo + modifiedFutureValues.IrRenewalCommissionOutgo + modifiedFutureValues.IrInitialExpenseOutgo + modifiedFutureValues.IrMaintenanceExpenseOutgo -
		modifiedFutureValues.IrNbPremiumReceipt

	mgp.IrDiscountedCashOutflowsLockedin = (modifiedFutureValues.IrDiscountedCashOutflowsLockedin+modifiedFutureValues.IrClaimsOutgo+modifiedFutureValues.IrClaimsExpenseOutgo)*discountFactorLockedin +
		modifiedFutureValues.IrInitialCommissionOutgo + modifiedFutureValues.IrRenewalCommissionOutgo + modifiedFutureValues.IrInitialExpenseOutgo + modifiedFutureValues.IrMaintenanceExpenseOutgo

	mgp.IrDiscountedEarnedPremium = modifiedFutureValues.IrDiscountedEarnedPremium*discountFactorLockedin + modifiedFutureValues.IrEarnedPremium

	mgp.RiskAdjustment = mgp.DiscountedEarnedPremiumLockedin * getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, run.ParameterYear, "RiskAdjustmentProportion", shock, run)

	mgp.IrRiskAdjustment = mgp.IrDiscountedEarnedPremium * getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, run.ParameterYear, "RiskAdjustmentProportion", shock, run)

	if portfolio.DiscountOption == Discounted {
		mgp.DiscountedCoverageUnits = mgp.DiscountedEarnedPremiumLockedin
	} else {
		mgp.DiscountedCoverageUnits = mgp.SumFutureEarnedPremium
	}

	if mgp.ProjectionMonth == 0 {
		mgp.ModifiedGMMCsm = -math.Min(mgp.DiscountedNetCashFlowsLockedin+mgp.RiskAdjustment, 0)
	}
	if mgp.ProjectionMonth > 0 && (mgp.CoverageUnits+mgp.DiscountedCoverageUnits) > 0 {
		mgp.CSMAllocationRatio = mgp.CoverageUnits / (mgp.CoverageUnits + mgp.DiscountedCoverageUnits)
	}

	mgp.Treaty1SumFutureEarnedPremium = modifiedFutureValues.Treaty1SumFutureEarnedPremium + modifiedFutureValues.Treaty1EarnedPremium
	mgp.Treaty1DiscountedEarnedPremiumCurrent = (modifiedFutureValues.Treaty1DiscountedEarnedPremiumCurrent + modifiedFutureValues.Treaty1EarnedPremium) * discountFactorCurrent
	mgp.Treaty1DiscountedEarnedPremiumLockedin = (modifiedFutureValues.Treaty1DiscountedEarnedPremiumLockedin + modifiedFutureValues.Treaty1EarnedPremium) * discountFactorLockedin

	mgp.Treaty1SumFutureCashOutflows = modifiedFutureValues.Treaty1SumFutureCashOutflows + modifiedFutureValues.Treaty1CashOutflow
	mgp.Treaty1DiscountedCashOutflows = (modifiedFutureValues.Treaty1DiscountedCashOutflows + modifiedFutureValues.Treaty1CashOutflow) * discountFactorCurrent
	mgp.Treaty1DiscountedCashOutflowsLockedin = (modifiedFutureValues.Treaty1DiscountedCashOutflowsLockedin + modifiedFutureValues.Treaty1CashOutflow) * discountFactorLockedin
	mgp.IrTreaty1DiscountedClaimsOutgo = (modifiedFutureValues.IrTreaty1DiscountedClaimsOutgo + modifiedFutureValues.IrTreaty1ClaimsOutgo) * discountFactorLockedin

	mgp.Treaty1SumFutureNetCashFlows = modifiedFutureValues.Treaty1SumFutureNetCashFlows + modifiedFutureValues.Treaty1NetCashFlow
	mgp.Treaty1DiscountedNetCashFlowsCurrent = (modifiedFutureValues.Treaty1DiscountedNetCashFlowsCurrent+modifiedFutureValues.Treaty1NetCashFlow)*discountFactorCurrent - modifiedFutureValues.Treaty1TotalPremiumReceipt
	mgp.Treaty1DiscountedNetCashFlowsLockedin = (modifiedFutureValues.Treaty1DiscountedNetCashFlowsLockedin+modifiedFutureValues.Treaty1NetCashFlow)*discountFactorLockedin - modifiedFutureValues.Treaty1TotalPremiumReceipt

	mgp.Treaty1RiskAdjustment = mgp.Treaty1DiscountedEarnedPremiumCurrent * getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, run.ParameterYear, "RiskAdjustmentProportion", shock, run)

	mgp.Treaty2SumFutureEarnedPremium = modifiedFutureValues.Treaty2SumFutureEarnedPremium + modifiedFutureValues.Treaty2EarnedPremium
	mgp.Treaty2DiscountedEarnedPremiumCurrent = (modifiedFutureValues.Treaty2DiscountedEarnedPremiumCurrent + modifiedFutureValues.Treaty2EarnedPremium) * discountFactorCurrent
	mgp.Treaty2DiscountedEarnedPremiumLockedin = (modifiedFutureValues.Treaty2DiscountedEarnedPremiumLockedin + modifiedFutureValues.Treaty2EarnedPremium) * discountFactorLockedin

	mgp.Treaty2SumFutureCashOutflows = modifiedFutureValues.Treaty2SumFutureCashOutflows + modifiedFutureValues.Treaty2CashOutflow
	mgp.Treaty2DiscountedCashOutflows = (modifiedFutureValues.Treaty2DiscountedCashOutflows + modifiedFutureValues.Treaty2CashOutflow) * discountFactorCurrent
	mgp.Treaty2DiscountedCashOutflowsLockedin = (modifiedFutureValues.Treaty2DiscountedCashOutflowsLockedin + modifiedFutureValues.Treaty2CashOutflow) * discountFactorLockedin
	mgp.IrTreaty2DiscountedClaimsOutgo = (modifiedFutureValues.IrTreaty2DiscountedClaimsOutgo + modifiedFutureValues.IrTreaty2ClaimsOutgo) * discountFactorLockedin

	mgp.Treaty2SumFutureNetCashFlows = modifiedFutureValues.Treaty2SumFutureNetCashFlows + modifiedFutureValues.Treaty2NetCashFlow
	mgp.Treaty2DiscountedNetCashFlowsCurrent = (modifiedFutureValues.Treaty2DiscountedNetCashFlowsCurrent+modifiedFutureValues.Treaty2NetCashFlow)*discountFactorCurrent - modifiedFutureValues.Treaty2TotalPremiumReceipt
	mgp.Treaty2DiscountedNetCashFlowsLockedin = (modifiedFutureValues.Treaty2DiscountedNetCashFlowsLockedin+modifiedFutureValues.Treaty2NetCashFlow)*discountFactorLockedin - modifiedFutureValues.Treaty2TotalPremiumReceipt

	mgp.Treaty2RiskAdjustment = mgp.Treaty2DiscountedEarnedPremiumCurrent * getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, run.ParameterYear, "RiskAdjustmentProportion", shock, run)

	mgp.Treaty3SumFutureEarnedPremium = modifiedFutureValues.Treaty3SumFutureEarnedPremium + modifiedFutureValues.Treaty3EarnedPremium
	mgp.Treaty3DiscountedEarnedPremiumCurrent = (modifiedFutureValues.Treaty3DiscountedEarnedPremiumCurrent + modifiedFutureValues.Treaty3EarnedPremium) * discountFactorCurrent
	mgp.Treaty3DiscountedEarnedPremiumLockedin = (modifiedFutureValues.Treaty3DiscountedEarnedPremiumLockedin + modifiedFutureValues.Treaty3EarnedPremium) * discountFactorLockedin
	mgp.IrTreaty3DiscountedClaimsOutgo = (modifiedFutureValues.IrTreaty3DiscountedClaimsOutgo + modifiedFutureValues.IrTreaty3ClaimsOutgo) * discountFactorLockedin

	mgp.Treaty3SumFutureCashOutflows = modifiedFutureValues.Treaty3SumFutureCashOutflows + modifiedFutureValues.Treaty3CashOutflow
	mgp.Treaty3DiscountedCashOutflows = (modifiedFutureValues.Treaty3DiscountedCashOutflows + modifiedFutureValues.Treaty3CashOutflow) * discountFactorCurrent
	mgp.Treaty3DiscountedCashOutflowsLockedin = (modifiedFutureValues.Treaty3DiscountedCashOutflowsLockedin + modifiedFutureValues.Treaty3CashOutflow) * discountFactorLockedin

	mgp.Treaty3SumFutureNetCashFlows = modifiedFutureValues.Treaty3SumFutureNetCashFlows + modifiedFutureValues.Treaty3NetCashFlow
	mgp.Treaty3DiscountedNetCashFlowsCurrent = (modifiedFutureValues.Treaty3DiscountedNetCashFlowsCurrent+modifiedFutureValues.Treaty3NetCashFlow)*discountFactorCurrent - modifiedFutureValues.Treaty3TotalPremiumReceipt
	mgp.Treaty3DiscountedNetCashFlowsLockedin = (modifiedFutureValues.Treaty3DiscountedNetCashFlowsLockedin+modifiedFutureValues.Treaty3NetCashFlow)*discountFactorLockedin - modifiedFutureValues.Treaty3TotalPremiumReceipt

	mgp.Treaty3RiskAdjustment = mgp.Treaty3DiscountedEarnedPremiumCurrent * getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, run.ParameterYear, "RiskAdjustmentProportion", shock, run)

	mgp.CededSumFutureEarnedPremium = modifiedFutureValues.CededSumFutureEarnedPremium + modifiedFutureValues.CededEarnedPremium
	mgp.CededDiscountedEarnedPremiumCurrent = (modifiedFutureValues.CededDiscountedEarnedPremiumCurrent + modifiedFutureValues.CededEarnedPremium) * discountFactorCurrent
	mgp.CededDiscountedEarnedPremiumLockedin = (modifiedFutureValues.CededDiscountedEarnedPremiumLockedin + modifiedFutureValues.CededEarnedPremium) * discountFactorLockedin

	mgp.CededSumFutureCashOutflows = modifiedFutureValues.CededSumFutureCashOutflows + modifiedFutureValues.CededCashOutflow
	mgp.CededDiscountedCashOutflows = (modifiedFutureValues.CededDiscountedCashOutflows + modifiedFutureValues.CededCashOutflow) * discountFactorCurrent
	mgp.CededDiscountedCashOutflowsLockedin = (modifiedFutureValues.CededDiscountedCashOutflowsLockedin + modifiedFutureValues.CededCashOutflow) * discountFactorLockedin

	mgp.CededSumFutureNetCashFlows = modifiedFutureValues.CededSumFutureNetCashFlows + modifiedFutureValues.CededNetCashFlow
	mgp.CededDiscountedNetCashFlowsCurrent = (modifiedFutureValues.CededDiscountedNetCashFlowsCurrent+modifiedFutureValues.CededCashOutflow)*discountFactorCurrent - modifiedFutureValues.TotalCededPremiumReceipt
	mgp.CededDiscountedNetCashFlowsLockedin = (modifiedFutureValues.CededDiscountedNetCashFlowsLockedin+modifiedFutureValues.CededCashOutflow)*discountFactorLockedin - modifiedFutureValues.TotalCededPremiumReceipt

	mgp.CededRiskAdjustment = mgp.Treaty1RiskAdjustment + mgp.Treaty2RiskAdjustment + mgp.Treaty3RiskAdjustment

}

func ModifiedGMMPremiumReceipt(mgp *models.AggregatedModifiedGMMProjection, mp models.ModifiedGMMModelPoint, runSettings models.GMMRunSetting, shock models.ModifiedGMMShock) {
	if mgp.ValuationMonth > int(utils.RoundUp(mp.TermMonths)) || mgp.ProjectionMonth == 0 || mgp.ValuationMonth == 0 {
		mgp.PremiumReceipt = 0
	} else {

		//var annualPremAdj float64
		//if mp.TermMonths < 12 && mp.Frequency != 0 {
		//	annualPremAdj = float64(mp.TermMonths) / 12.0
		//} else {
		//	annualPremAdj = 1.0
		//}

		if mp.Status == IF {
			if mp.Frequency == 0 { //Single Premium
				if mgp.ValuationMonth == 1 {
					mgp.PremiumReceipt = utils.FloatPrecision(mp.AnnualPremium*mgp.InforcePolicyCountSM, AccountingPrecision)
					mgp.Treaty1PremiumReceipt = utils.FloatPrecision(mgp.PremiumReceipt*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty1PremiumProportion", shock, runSettings), AccountingPrecision)
					mgp.Treaty2PremiumReceipt = utils.FloatPrecision(mgp.PremiumReceipt*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty2PremiumProportion", shock, runSettings), AccountingPrecision)
					mgp.Treaty3PremiumReceipt = utils.FloatPrecision(mgp.PremiumReceipt*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty3PremiumProportion", shock, runSettings), AccountingPrecision)
					mgp.CededPremiumReceipt = mgp.Treaty1PremiumReceipt + mgp.Treaty2PremiumReceipt + mgp.Treaty3PremiumReceipt
				} else {
					mgp.PremiumReceipt = 0
					mgp.Treaty1PremiumReceipt = 0
					mgp.Treaty2PremiumReceipt = 0
					mgp.Treaty3PremiumReceipt = 0
					mgp.CededPremiumReceipt = 0
				}
			} else if mp.Frequency != 0 {
				if math.Mod(float64(mgp.ValuationMonth-1), float64(12/mp.Frequency)) == 0 {
					mgp.PremiumReceipt = utils.FloatPrecision(mp.AnnualPremium*mgp.InforcePolicyCountSM/float64(mp.Frequency), AccountingPrecision)
					mgp.Treaty1PremiumReceipt = utils.FloatPrecision(mgp.PremiumReceipt*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty1PremiumProportion", shock, runSettings), AccountingPrecision)
					mgp.Treaty2PremiumReceipt = utils.FloatPrecision(mgp.PremiumReceipt*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty2PremiumProportion", shock, runSettings), AccountingPrecision)
					mgp.Treaty3PremiumReceipt = utils.FloatPrecision(mgp.PremiumReceipt*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty3PremiumProportion", shock, runSettings), AccountingPrecision)
					mgp.CededPremiumReceipt = mgp.Treaty1PremiumReceipt + mgp.Treaty2PremiumReceipt + mgp.Treaty3PremiumReceipt
				} else {
					mgp.PremiumReceipt = 0
					mgp.Treaty1PremiumReceipt = 0
					mgp.Treaty2PremiumReceipt = 0
					mgp.Treaty3PremiumReceipt = 0
					mgp.CededPremiumReceipt = 0
				}
			}
			mgp.NbPremiumReceipt = 0
			mgp.Treaty1NbPremiumReceipt = 0
			mgp.Treaty2NbPremiumReceipt = 0
			mgp.Treaty3NbPremiumReceipt = 0
			mgp.CededNbPremiumReceipt = 0
		}
		if mp.Status == NB {
			if mp.Frequency == 0 { // Single Premium
				if mgp.ValuationMonth == 1 {
					mgp.NbPremiumReceipt = utils.FloatPrecision(mp.AnnualPremium*mgp.InforcePolicyCountSM, AccountingPrecision)
					mgp.Treaty1NbPremiumReceipt = utils.FloatPrecision(mgp.NbPremiumReceipt*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty1PremiumProportion", shock, runSettings), AccountingPrecision)
					mgp.Treaty2NbPremiumReceipt = utils.FloatPrecision(mgp.NbPremiumReceipt*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty2PremiumProportion", shock, runSettings), AccountingPrecision)
					mgp.Treaty3NbPremiumReceipt = utils.FloatPrecision(mgp.NbPremiumReceipt*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty3PremiumProportion", shock, runSettings), AccountingPrecision)
					mgp.CededNbPremiumReceipt = mgp.Treaty1NbPremiumReceipt + mgp.Treaty2NbPremiumReceipt + mgp.Treaty3NbPremiumReceipt
				} else {
					mgp.NbPremiumReceipt = 0
					mgp.Treaty1NbPremiumReceipt = 0
					mgp.Treaty2NbPremiumReceipt = 0
					mgp.Treaty3NbPremiumReceipt = 0
					mgp.CededNbPremiumReceipt = 0
				}
			} else if mp.Frequency != 0 {
				if math.Mod(float64(mgp.ValuationMonth-1), float64(12/mp.Frequency)) == 0 {
					mgp.NbPremiumReceipt = utils.FloatPrecision(mp.AnnualPremium*mgp.InforcePolicyCountSM/float64(mp.Frequency), AccountingPrecision)
					mgp.Treaty1NbPremiumReceipt = utils.FloatPrecision(mgp.NbPremiumReceipt*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty1PremiumProportion", shock, runSettings), AccountingPrecision)
					mgp.Treaty2NbPremiumReceipt = utils.FloatPrecision(mgp.NbPremiumReceipt*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty2PremiumProportion", shock, runSettings), AccountingPrecision)
					mgp.Treaty3NbPremiumReceipt = utils.FloatPrecision(mgp.NbPremiumReceipt*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty3PremiumProportion", shock, runSettings), AccountingPrecision)
					mgp.CededNbPremiumReceipt = mgp.Treaty1NbPremiumReceipt + mgp.Treaty2NbPremiumReceipt + mgp.Treaty3NbPremiumReceipt
				} else {
					mgp.NbPremiumReceipt = 0
					mgp.Treaty1NbPremiumReceipt = 0
					mgp.Treaty2NbPremiumReceipt = 0
					mgp.Treaty3NbPremiumReceipt = 0
					mgp.CededNbPremiumReceipt = 0
				}
			}
			mgp.PremiumReceipt = 0
			mgp.Treaty1PremiumReceipt = 0
			mgp.Treaty2PremiumReceipt = 0
			mgp.Treaty3PremiumReceipt = 0
			mgp.CededPremiumReceipt = 0
		}

	}
	mgp.TotalPremiumReceipt = mgp.PremiumReceipt + mgp.NbPremiumReceipt
	mgp.Treaty1TotalPremiumReceipt = mgp.Treaty1PremiumReceipt + mgp.Treaty1NbPremiumReceipt
	mgp.Treaty2TotalPremiumReceipt = mgp.Treaty2PremiumReceipt + mgp.Treaty2NbPremiumReceipt
	mgp.Treaty3TotalPremiumReceipt = mgp.Treaty3PremiumReceipt + mgp.Treaty3NbPremiumReceipt
	mgp.TotalCededPremiumReceipt = mgp.Treaty1TotalPremiumReceipt + mgp.Treaty2TotalPremiumReceipt + mgp.Treaty3TotalPremiumReceipt
}

func ModifiedGMMEarnedPremium(mgp *models.AggregatedModifiedGMMProjection, mp models.ModifiedGMMModelPoint, portfolio models.PaaPortfolio, runSettings models.GMMRunSetting, sumRiskUnits, totalExpectedPremiumFromInception float64, shock models.ModifiedGMMShock, DacBuildupIndicator bool) {
	if mgp.ValuationMonth > int(utils.RoundUp(mp.TermMonths)) || mgp.ProjectionMonth == 0 || mgp.ValuationMonth == 0 || mp.TermMonths == 0 {
		mgp.EarnedPremium = 0
		mgp.Treaty1EarnedPremium = 0
		mgp.Treaty2EarnedPremium = 0
		mgp.Treaty3EarnedPremium = 0
		mgp.CededEarnedPremium = 0

		treaty1PremProp := getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty1PremiumProportion", shock, runSettings)
		treaty2PremProp := getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty2PremiumProportion", shock, runSettings)
		treaty3PremProp := getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty3PremiumProportion", shock, runSettings)

		if mgp.ProjectionMonth == 0 {
			var currentPeriodDuration float64
			var durationStartPeriod int
			var tempEarnedProportion float64
			var sumCoveredUnits float64
			var monthIndex int
			var commDate, endCoverDate, dailyValDate time.Time

			runMonth, _ := strconv.Atoi(runSettings.RunDate[5:])

			//if mp.Status == NB {
			//	durationStartPeriod = int(math.Max(float64(mp.DurationInForceMonths-monthIndex), 0))
			//	currentPeriodDuration = int(math.Min(float64(mp.DurationInForceMonths), float64(mp.TermMonths))) //int(math.Max(float64(runMonth-mp[0].CommencementMonth+1), 0))
			//
			//}
			//if mp.Status == IF {
			if runMonth == runSettings.YearEndMonth {
				monthIndex = 12
			} else if runMonth < runSettings.YearEndMonth {
				monthIndex = 12 - (runSettings.YearEndMonth - runMonth)
			} else {
				monthIndex = runMonth - runSettings.YearEndMonth
			}
			durationStartPeriod = int(math.Max(float64(mp.DurationInForceMonths-monthIndex), 0)) //int(math.Max(float64(mp.DurationInForceMonths-runMonth), 0))
			outstandingTermTemp := math.Max(mp.TermMonths-float64(durationStartPeriod), 0)
			//if outstandingTermTemp < 12 {
			//	currentPeriodDuration = int(math.Min(float64(runMonth), outstandingTermTemp))
			//} else {
			//	currentPeriodDuration = 12
			//}
			currentPeriodDuration = math.Min(math.Min(outstandingTermTemp, float64(monthIndex)), float64(mp.DurationInForceMonths))
			//}
			if portfolio.PremiumEarningPattern == PassageofTime {
				if mp.TermMonths > 0 {
					tempEarnedProportion = float64(currentPeriodDuration) / float64(mp.TermMonths)
				}
			}

			if portfolio.PremiumEarningPattern == DailyPassageofTime {
				var err error
				dailyValDate, err = utils.ParseDateString(runSettings.RunDate + "-01")
				dailyValDate = endOfMonth(dailyValDate, 0)
				commDate, err = utils.ParseDateString(mp.CommencementDate)
				endCoverDate, err = utils.ParseDateString(mp.CoverEndDate)
				tempterm := endCoverDate.Sub(commDate).Hours() / 24
				currentPeriodDuration = dailyValDate.Sub(commDate).Hours() / 24
				if err != nil {
					fmt.Println(err)
				}
				if tempterm > 0 {
					tempEarnedProportion = float64(currentPeriodDuration) / float64(tempterm)
				}
			}

			if portfolio.PremiumEarningPattern == SpecifiedbyUser {

				if sumRiskUnits > 0 {
					sumCoveredUnits = 0
					for i := durationStartPeriod + 1; i <= durationStartPeriod+int(utils.RoundUp(currentPeriodDuration)); i++ {
						sumCoveredUnits += getRiskUnit(runSettings.ParameterYear, mp.ProductCode, mp.SubProductCode, i)
					}

					tempEarnedProportion = sumCoveredUnits / sumRiskUnits
				}
			}

			//Acquisition Expenses
			var tempcalc float64
			var annualPremAdj float64
			if mp.TermMonths < 12 && mp.Frequency != 0 {
				annualPremAdj = float64(mp.TermMonths) / 12.0
			} else {
				annualPremAdj = 1.0
			}

			if DacBuildupIndicator {
				tempcalc = getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "InitialExpenseAmount", shock, runSettings) +
					mp.AnnualPremium*annualPremAdj*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "InitialExpenseProportion", shock, runSettings) +
					getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "InitialCommissionAmount", shock, runSettings) +
					mp.AnnualPremium*annualPremAdj*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "InitialYr1CommissionProportion", shock, runSettings)
			} else {
				tempcalc = 0
			}

			mgp.CurrentPeriodAmortisedAcquisition = tempcalc * tempEarnedProportion
			mgp.CurrentPeriodEarnedPremium = (totalExpectedPremiumFromInception - tempcalc) * tempEarnedProportion //math.Max(totalExpectedPremiumFromInception-tempcalc, 0) * tempEarnedProportion
			mgp.Treaty1CurrentPeriodEarnedPremium = totalExpectedPremiumFromInception * tempEarnedProportion * treaty1PremProp
			mgp.Treaty2CurrentPeriodEarnedPremium = totalExpectedPremiumFromInception * tempEarnedProportion * treaty2PremProp
			mgp.Treaty3CurrentPeriodEarnedPremium = totalExpectedPremiumFromInception * tempEarnedProportion * treaty3PremProp
			mgp.WrittenPremium = totalExpectedPremiumFromInception
			mgp.Treaty1WrittenPremium = totalExpectedPremiumFromInception * treaty1PremProp
			mgp.Treaty2WrittenPremium = totalExpectedPremiumFromInception * treaty2PremProp
			mgp.Treaty3WrittenPremium = totalExpectedPremiumFromInception * treaty3PremProp
		}

	} else {
		if portfolio.PremiumEarningPattern == PassageofTime {
			if mp.TermMonths > 0 {
				mgp.EarnedPremium = utils.FloatPrecision(mgp.InforcePolicyCountSM*totalExpectedPremiumFromInception/float64(mp.TermMonths), AccountingPrecision)
			}
		}

		if portfolio.PremiumEarningPattern == SpecifiedbyUser {
			if sumRiskUnits > 0 {
				mgp.EarnedPremium = utils.FloatPrecision(mgp.InforcePolicyCountSM*totalExpectedPremiumFromInception*getRiskUnit(runSettings.ParameterYear, mp.ProductCode, mp.SubProductCode, mp.DurationInForceMonths+mgp.ProjectionMonth)/sumRiskUnits, AccountingPrecision)
			}

		}

		mgp.Treaty1EarnedPremium = utils.FloatPrecision(mgp.EarnedPremium*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty1PremiumProportion", shock, runSettings), AccountingPrecision)
		mgp.Treaty2EarnedPremium = utils.FloatPrecision(mgp.EarnedPremium*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty2PremiumProportion", shock, runSettings), AccountingPrecision)
		mgp.Treaty3EarnedPremium = utils.FloatPrecision(mgp.EarnedPremium*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty3PremiumProportion", shock, runSettings), AccountingPrecision)

		mgp.CededEarnedPremium = mgp.Treaty1EarnedPremium + mgp.Treaty2EarnedPremium + mgp.Treaty3EarnedPremium
	}
	mgp.CoverageUnits = mgp.EarnedPremium
	mgp.CurrentPeriodInsuranceRevenue = mgp.CurrentPeriodEarnedPremium + mgp.CurrentPeriodAmortisedAcquisition
}

func ModifiedGMMCashOutflow(mgp *models.AggregatedModifiedGMMProjection, mp models.ModifiedGMMModelPoint, runSettings models.GMMRunSetting, shock models.ModifiedGMMShock) {
	if mgp.ValuationMonth > int(utils.RoundUp(mp.TermMonths)) || mgp.ProjectionMonth == 0 || mgp.ValuationMonth == 0 {
		mgp.ClaimsOutgo = 0
		mgp.ClaimsExpenseOutgo = 0
		mgp.MaintenanceExpenseOutgo = 0
		//mgp.ExpenseOutgo = 0
		mgp.InitialCommissionOutgo = 0
		mgp.InitialExpenseOutgo = 0
		//mgp.AcquisitionExpenseOutgo = 0
	} else {
		mgp.ClaimsOutgo = utils.FloatPrecision(mgp.EarnedPremium*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ClaimsProportion", shock, runSettings), AccountingPrecision)
		if mgp.Treaty1EarnedPremium > 0 {
			mgp.Treaty1ClaimsOutgo = utils.FloatPrecision(mgp.ClaimsOutgo*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty1ClaimsProportion", shock, runSettings), AccountingPrecision)
		}
		if mgp.Treaty2EarnedPremium > 0 {
			mgp.Treaty2ClaimsOutgo = utils.FloatPrecision(mgp.ClaimsOutgo*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty2ClaimsProportion", shock, runSettings), AccountingPrecision)
		}
		if mgp.Treaty3EarnedPremium > 0 {
			mgp.Treaty3ClaimsOutgo = utils.FloatPrecision(mgp.ClaimsOutgo*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty3ClaimsProportion", shock, runSettings), AccountingPrecision)
		}
		mgp.ClaimsExpenseOutgo = utils.FloatPrecision(mgp.EarnedPremium*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ClaimsExpenseProportion", shock, runSettings), AccountingPrecision)

		var annualPremAdj float64
		if mp.TermMonths < 12 && mp.Frequency != 0 {
			annualPremAdj = float64(mp.TermMonths) / 12.0
		} else {
			annualPremAdj = 1.0
		}

		if mgp.ValuationMonth == 1 {
			mgp.InitialCommissionOutgo = utils.FloatPrecision(mp.AnnualPremium*annualPremAdj*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "InitialYr1CommissionProportion", shock, runSettings)+getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "InitialCommissionAmount", shock, runSettings), AccountingPrecision)
			mgp.InitialExpenseOutgo = utils.FloatPrecision(mp.AnnualPremium*annualPremAdj*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "InitialExpenseProportion", shock, runSettings)+getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "InitialExpenseAmount", shock, runSettings), AccountingPrecision)
		}
		if mgp.ValuationMonth == 12 {
			mgp.InitialCommissionOutgo = utils.FloatPrecision(mp.AnnualPremium*annualPremAdj*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "InitialYr2CommissionProportion", shock, runSettings), AccountingPrecision)
		}

		if mp.Frequency != 0 && mgp.TotalPremiumReceipt > 0 {
			mgp.RenewalCommissionOutgo = utils.FloatPrecision(mgp.TotalPremiumReceipt*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "RenewalCommissionProportion", shock, runSettings)+getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "RenewalAnnualCommissionAmount", shock, runSettings)/float64(mp.Frequency), AccountingPrecision)
		}
		if mp.Frequency == 0 && mgp.TotalPremiumReceipt > 0 {
			mgp.RenewalCommissionOutgo = utils.FloatPrecision(mgp.TotalPremiumReceipt*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "RenewalCommissionProportion", shock, runSettings)+getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "RenewalAnnualCommissionAmount", shock, runSettings), AccountingPrecision)
		}
		if mgp.EarnedPremium > 0 {
			mgp.MaintenanceExpenseOutgo = utils.FloatPrecision(mgp.EarnedPremium*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "MaintenanceExpenseProportion", shock, runSettings)+getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "MaintenanceAnnualExpenseAmount", shock, runSettings)/12.0, AccountingPrecision)

		}

		//mgp.AcquisitionExpenseOutgo = mgp.InitialCommissionOutgo + mgp.InitialExpenseOutgo + mgp.RenewalCommissionOutgo
		mgp.CashOutflow = mgp.ClaimsOutgo + mgp.ClaimsExpenseOutgo + mgp.MaintenanceExpenseOutgo + mgp.InitialExpenseOutgo + mgp.RenewalCommissionOutgo + mgp.InitialCommissionOutgo
		mgp.Treaty1CashOutflow = mgp.Treaty1ClaimsOutgo
		mgp.Treaty2CashOutflow = mgp.Treaty2ClaimsOutgo
		mgp.Treaty3CashOutflow = mgp.Treaty3ClaimsOutgo
		mgp.CededCashOutflow = mgp.Treaty1CashOutflow + mgp.Treaty2CashOutflow + mgp.Treaty3CashOutflow
		mgp.CededClaimsOutgo = mgp.Treaty1ClaimsOutgo + mgp.Treaty2ClaimsOutgo + mgp.Treaty3ClaimsOutgo
	}
}

func InitialRecognitionModifiedGMMPremiumReceipt(mgp *models.AggregatedModifiedGMMProjection, mp models.ModifiedGMMModelPoint, runSettings models.GMMRunSetting, shock models.ModifiedGMMShock) {
	if mgp.IrValuationMonth > int(utils.RoundUp(mp.TermMonths)) || mgp.ProjectionMonth == 0 || mgp.IrValuationMonth == 0 {
		mgp.IrNbPremiumReceipt = 0
	} else {

		//var annualPremAdj float64
		//if mp.TermMonths < 12 && mp.Frequency != 0 {
		//	annualPremAdj = float64(mp.TermMonths) / 12.0
		//} else {
		//	annualPremAdj = 1.0
		//}
		if mp.Frequency == 0 { // Single Premium
			if mgp.IrValuationMonth == 1 {
				mgp.IrNbPremiumReceipt = utils.FloatPrecision(mp.AnnualPremium*mgp.IrInforcePolicyCountSM, AccountingPrecision)
				//mgp.Treaty1NbPremiumReceipt = utils.FloatPrecision(mgp.NbPremiumReceipt*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty1PremiumProportion", shock, runSettings), AccountingPrecision)
				//mgp.Treaty2NbPremiumReceipt = utils.FloatPrecision(mgp.NbPremiumReceipt*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty2PremiumProportion", shock, runSettings), AccountingPrecision)
				//mgp.Treaty3NbPremiumReceipt = utils.FloatPrecision(mgp.NbPremiumReceipt*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty3PremiumProportion", shock, runSettings), AccountingPrecision)
				//mgp.CededNbPremiumReceipt = mgp.Treaty1NbPremiumReceipt + mgp.Treaty2NbPremiumReceipt + mgp.Treaty3NbPremiumReceipt
			} else {
				mgp.IrNbPremiumReceipt = 0
				//mgp.Treaty1NbPremiumReceipt = 0
				//mgp.Treaty2NbPremiumReceipt = 0
				//mgp.Treaty3NbPremiumReceipt = 0
				//mgp.CededNbPremiumReceipt = 0
			}
		} else if mp.Frequency != 0 {
			if math.Mod(float64(mgp.IrValuationMonth-1), float64(12/mp.Frequency)) == 0 {
				mgp.IrNbPremiumReceipt = utils.FloatPrecision(mp.AnnualPremium*mgp.IrInforcePolicyCountSM/float64(mp.Frequency), AccountingPrecision)
				//mgp.Treaty1NbPremiumReceipt = utils.FloatPrecision(mgp.NbPremiumReceipt*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty1PremiumProportion", shock, runSettings), AccountingPrecision)
				//mgp.Treaty2NbPremiumReceipt = utils.FloatPrecision(mgp.NbPremiumReceipt*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty2PremiumProportion", shock, runSettings), AccountingPrecision)
				//mgp.Treaty3NbPremiumReceipt = utils.FloatPrecision(mgp.NbPremiumReceipt*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty3PremiumProportion", shock, runSettings), AccountingPrecision)
				//mgp.CededNbPremiumReceipt = mgp.Treaty1NbPremiumReceipt + mgp.Treaty2NbPremiumReceipt + mgp.Treaty3NbPremiumReceipt
			} else {
				mgp.IrNbPremiumReceipt = 0
				//mgp.Treaty1NbPremiumReceipt = 0
				//mgp.Treaty2NbPremiumReceipt = 0
				//mgp.Treaty3NbPremiumReceipt = 0
				//mgp.CededNbPremiumReceipt = 0
			}
		}
	}
}

func InitialRecognitionModifiedGMMEarnedPremium(mgp *models.AggregatedModifiedGMMProjection, mp models.ModifiedGMMModelPoint, portfolio models.PaaPortfolio, runSettings models.GMMRunSetting, sumRiskUnits, totalExpectedPremiumFromInception float64, shock models.ModifiedGMMShock) {
	if mgp.IrValuationMonth > int(utils.RoundUp(mp.TermMonths)) || mgp.ProjectionMonth == 0 || mgp.IrValuationMonth == 0 || mp.TermMonths == 0 {
		mgp.IrEarnedPremium = 0
		mgp.IrTreaty1EarnedPremium = 0
		mgp.IrTreaty2EarnedPremium = 0
		mgp.IrTreaty3EarnedPremium = 0
	} else {
		if portfolio.PremiumEarningPattern == PassageofTime {
			if mp.TermMonths > 0 {
				mgp.IrEarnedPremium = utils.FloatPrecision(mgp.IrInforcePolicyCountSM*totalExpectedPremiumFromInception/float64(mp.TermMonths), AccountingPrecision)
			}
		} else {
			if sumRiskUnits > 0 {
				mgp.IrEarnedPremium = utils.FloatPrecision(mgp.IrInforcePolicyCountSM*totalExpectedPremiumFromInception*getRiskUnit(runSettings.ParameterYear, mp.ProductCode, mp.SubProductCode, mgp.ProjectionMonth)/sumRiskUnits, AccountingPrecision)
			}

		}
		mgp.IrTreaty1EarnedPremium = utils.FloatPrecision(mgp.IrEarnedPremium*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty1PremiumProportion", shock, runSettings), AccountingPrecision)
		mgp.IrTreaty2EarnedPremium = utils.FloatPrecision(mgp.IrEarnedPremium*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty2PremiumProportion", shock, runSettings), AccountingPrecision)
		mgp.IrTreaty3EarnedPremium = utils.FloatPrecision(mgp.IrEarnedPremium*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty3PremiumProportion", shock, runSettings), AccountingPrecision)

	}
}

func InitialRecognitionModifiedGMMCashOutflow(mgp *models.AggregatedModifiedGMMProjection, mp models.ModifiedGMMModelPoint, runSettings models.GMMRunSetting, shock models.ModifiedGMMShock) {
	if mgp.IrValuationMonth > int(utils.RoundUp(mp.TermMonths)) || mgp.ProjectionMonth == 0 || mgp.IrValuationMonth == 0 {
		mgp.IrClaimsOutgo = 0
		mgp.IrClaimsExpenseOutgo = 0
		mgp.IrMaintenanceExpenseOutgo = 0
		//mgp.ExpenseOutgo = 0
		mgp.IrInitialCommissionOutgo = 0
		mgp.IrInitialExpenseOutgo = 0
		//mgp.AcquisitionExpenseOutgo = 0
	} else {
		mgp.IrClaimsOutgo = utils.FloatPrecision(mgp.IrEarnedPremium*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ClaimsProportion", shock, runSettings), AccountingPrecision)
		if mgp.IrTreaty1EarnedPremium > 0 {
			mgp.IrTreaty1ClaimsOutgo = utils.FloatPrecision(mgp.IrClaimsOutgo*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty1ClaimsProportion", shock, runSettings), AccountingPrecision)
		}
		if mgp.IrTreaty2EarnedPremium > 0 {
			mgp.IrTreaty2ClaimsOutgo = utils.FloatPrecision(mgp.IrClaimsOutgo*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty2ClaimsProportion", shock, runSettings), AccountingPrecision)
		}
		if mgp.IrTreaty3EarnedPremium > 0 {
			mgp.IrTreaty3ClaimsOutgo = utils.FloatPrecision(mgp.IrClaimsOutgo*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ReinsuranceTreaty3ClaimsProportion", shock, runSettings), AccountingPrecision)
		}
		mgp.IrClaimsExpenseOutgo = utils.FloatPrecision(mgp.IrEarnedPremium*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "ClaimsExpenseProportion", shock, runSettings), AccountingPrecision)

		var annualPremAdj float64
		if mp.TermMonths < 12 && mp.Frequency != 0 {
			annualPremAdj = float64(mp.TermMonths) / 12.0
		} else {
			annualPremAdj = 1.0
		}

		if mgp.IrValuationMonth == 1 {
			mgp.IrInitialCommissionOutgo = utils.FloatPrecision(mp.AnnualPremium*annualPremAdj*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "InitialYr1CommissionProportion", shock, runSettings)+getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "InitialCommissionAmount", shock, runSettings), AccountingPrecision)
			mgp.IrInitialExpenseOutgo = utils.FloatPrecision(mp.AnnualPremium*annualPremAdj*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "InitialExpenseProportion", shock, runSettings)+getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "InitialExpenseAmount", shock, runSettings), AccountingPrecision)
		}
		if mgp.IrValuationMonth == 12 {
			mgp.IrInitialCommissionOutgo = utils.FloatPrecision(mp.AnnualPremium*annualPremAdj*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "InitialYr2CommissionProportion", shock, runSettings), AccountingPrecision)
		}

		if mp.Frequency != 0 && mgp.IrNbPremiumReceipt > 0 {
			mgp.IrRenewalCommissionOutgo = utils.FloatPrecision(mgp.IrNbPremiumReceipt*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "RenewalCommissionProportion", shock, runSettings)+getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "RenewalAnnualCommissionAmount", shock, runSettings)/float64(mp.Frequency), AccountingPrecision)
		}
		if mp.Frequency == 0 && mgp.IrNbPremiumReceipt > 0 {
			mgp.IrRenewalCommissionOutgo = utils.FloatPrecision(mgp.IrNbPremiumReceipt*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "RenewalCommissionProportion", shock, runSettings)+getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "RenewalAnnualCommissionAmount", shock, runSettings), AccountingPrecision)
		}
		if mgp.IrEarnedPremium > 0 {
			mgp.IrMaintenanceExpenseOutgo = utils.FloatPrecision(mgp.IrEarnedPremium*getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "MaintenanceExpenseProportion", shock, runSettings)+getModifiedGMMParameter(mp.ProductCode, mp.SubProductCode, runSettings.ParameterYear, "MaintenanceAnnualExpenseAmount", shock, runSettings)/12.0, AccountingPrecision)

		}
	}
}

func getModifiedGMMParameter(productCode string, subProductCode string, year int, paramType string, shock models.ModifiedGMMShock, runSettings models.GMMRunSetting) float64 {
	var param models.ModifiedGMMParameter
	var claimMult, claimAdd, expMult, expAdd float64
	if runSettings.ShockSettings.Expenses {
		expMult = shock.MultiplicativeExpenses
		expAdd = shock.AdditiveExpenses
	}
	if runSettings.ShockSettings.Claims {
		claimMult = shock.MultiplicativeClaims
		claimAdd = shock.AdditiveClaims
	}

	key := "gmm-" + strings.ToLower(productCode) + "-" + subProductCode + "-" + strconv.Itoa(year)

	result, found := PaaCache.Get(key)

	if found {
		param = result.(models.ModifiedGMMParameter)
	} else {
		err := DB.Where("product_code = ? and sub_product_code=? and year = ?", productCode, subProductCode, year).Find(&param).Error
		if err != nil {
			return 0
		}
		PaaCache.Set(key, param, 1)
	}

	switch paramType {
	case "ClaimsProportion":
		return param.ClaimsProportion*(1+claimMult) + claimAdd
	case "ClaimsExpenseProportion":
		return param.ClaimsExpenseProportion*(1+expMult) + expAdd
	case "MaintenanceExpenseProportion":
		return param.MaintenanceExpenseProportion*(1+expMult) + expAdd
	case "InitialCommissionAmount":
		return param.InitialCommissionAmount
	case "InitialYr1CommissionProportion":
		return param.InitialYr1CommissionProportion
	case "InitialYr2CommissionProportion":
		return param.InitialYr2CommissionProportion
	case "RenewalAnnualCommissionAmount":
		return param.RenewalAnnualCommissionAmount
	case "RenewalCommissionProportion":
		return param.RenewalCommissionProportion
	//case "AcquisitionExpenseProportion":
	//	return math.Min(param.InitialExpenseProportion*(1+shock.MultiplicativeExpenses)+shock.AdditiveExpenses, 1)
	case "InitialExpenseProportion":
		return param.InitialExpenseProportion*(1+expMult) + expAdd //return math.Min(param.InitialExpenseProportion*(1+expMult+expAdd), 1)
	case "InitialExpenseAmount":
		return param.InitialExpenseAmount*(1+expMult) + expAdd
	case "MaintenanceAnnualExpenseAmount":
		return param.MaintenanceAnnualExpenseAmount*(1+expMult) + expAdd
	case "IBNRProportion":
		return param.IBNRProportion
	case "ReinsuranceTreaty1ClaimsProportion":
		return param.ReinsuranceTreaty1ClaimsProportion
	case "ReinsuranceTreaty1PremiumProportion":
		return param.ReinsuranceTreaty1PremiumProportion
	case "ReinsuranceTreaty2ClaimsProportion":
		return param.ReinsuranceTreaty2ClaimsProportion
	case "ReinsuranceTreaty2PremiumProportion":
		return param.ReinsuranceTreaty2PremiumProportion
	case "ReinsuranceTreaty3ClaimsProportion":
		return param.ReinsuranceTreaty3ClaimsProportion
	case "ReinsuranceTreaty3PremiumProportion":
		return param.ReinsuranceTreaty3PremiumProportion
	case "RiskAdjustmentProportion":
		return param.RiskAdjustmentProportion
	case "ReinsuranceTreaty1RiskAdjustmentProportion":
		return param.ReinsuranceTreaty1RiskAdjustmentProportion
	case "ReinsuranceTreaty2RiskAdjustmentProportion":
		return param.ReinsuranceTreaty2RiskAdjustmentProportion
	case "ReinsuranceTreaty3RiskAdjustmentProportion":
		return param.ReinsuranceTreaty3RiskAdjustmentProportion
	default:
		return 0
	}
}

func getModifiedGMMParameterString(productCode string, subProductCode string, year int, paramType string) string {
	var param models.ModifiedGMMParameter

	key := "gmm-" + strings.ToLower(productCode) + "-" + subProductCode + "-" + strconv.Itoa(year)

	result, found := PaaCache.Get(key)

	if found {
		param = result.(models.ModifiedGMMParameter)
	} else {
		err := DB.Where("product_code = ? and sub_product_code=? and year = ?", productCode, subProductCode, year).Find(&param).Error
		if err != nil {
			return "0"
		}
		PaaCache.Set(key, param, 1)
	}
	switch paramType {
	case "YieldCurveCode":
		return param.YieldCurveCode
	default:
		return "0"
	}
}

func getModifiedGMMParameterBoolean(productCode string, subProductCode string, year int, paramType string) bool {
	var param models.ModifiedGMMParameter

	key := "gmm-" + strings.ToLower(productCode) + "-" + subProductCode + "-" + strconv.Itoa(year)

	result, found := PaaCache.Get(key)

	if found {
		param = result.(models.ModifiedGMMParameter)
	} else {
		err := DB.Where("product_code = ? and sub_product_code=? and year = ?", productCode, subProductCode, year).Find(&param).Error
		if err != nil {
			return false
		}
		PaaCache.Set(key, param, 1)
	}
	switch paramType {
	case "DacBuildupIndicator":
		return param.DacBuildupIndicator
	default:
		return false
	}
}

func getPAAReinsuranceParameter(productCode string, paramType string) (float64, string) {

	var param models.ReinsuranceParameter

	key := strings.ToLower(productCode)

	result, found := PaaCache.Get(key)

	if found {
		param = result.(models.ReinsuranceParameter)
	} else {
		err := DB.Where("product_code = ?", productCode).Find(&param).Error
		if err != nil {
			return 0, ""
		}
		PaaCache.Set(key, param, 1)
	}

	switch paramType {
	case "ReinsuranceInwardOutward":
		return 0, param.ReinsuranceInwardOutward

	case "UlrLowerboundRate":
		return param.UlrLowerboundRate, ""
	case "UlrUpperboundRate":
		return param.UlrUpperboundRate, ""

	case "SlidingScaleMinRate":
		return param.SlidingScaleMinRate, ""
	case "SlidingScaleMaxRate":
		return param.SlidingScaleMaxRate, ""
	case "ProfitCommissionModel":
		return 0, param.ProfitCommissionModel
	case "ProfitCommissionRate":
		return param.ProfitCommissionRate, ""
	default:
		return 0, ""
	}
}

func ModifiedGMMNetCashflow(mgp *models.AggregatedModifiedGMMProjection, mp models.ModifiedGMMModelPoint) {
	if mgp.ValuationMonth > int(utils.RoundUp(mp.TermMonths)) || mgp.ProjectionMonth == 0 || mgp.ValuationMonth == 0 {
		mgp.NetCashFlow = 0
	} else {
		mgp.NetCashFlow = mgp.CashOutflow - mgp.TotalPremiumReceipt
		mgp.Treaty1NetCashFlow = mgp.Treaty1TotalPremiumReceipt - mgp.Treaty1ClaimsOutgo
		mgp.Treaty2NetCashFlow = mgp.Treaty2TotalPremiumReceipt + mgp.Treaty2NbPremiumReceipt - mgp.Treaty2ClaimsOutgo
		mgp.Treaty3NetCashFlow = mgp.Treaty3TotalPremiumReceipt + mgp.Treaty3NbPremiumReceipt - mgp.Treaty3ClaimsOutgo
		mgp.CededNetCashFlow = mgp.TotalCededPremiumReceipt - mgp.CededClaimsOutgo

	}
}

func InitialRecognitionModifiedGMMNetCashflow(mgp *models.AggregatedModifiedGMMProjection, mp models.ModifiedGMMModelPoint) {
	if mgp.IrValuationMonth > int(utils.RoundUp(mp.TermMonths)) || mgp.ProjectionMonth == 0 || mgp.IrValuationMonth == 0 {
		mgp.IrNetCashFlow = 0
	} else {
		mgp.IrNetCashFlow = mgp.IrClaimsOutgo + mgp.IrClaimsExpenseOutgo + mgp.IrInitialExpenseOutgo + mgp.IrMaintenanceExpenseOutgo + mgp.IrInitialCommissionOutgo +
			mgp.IrRenewalCommissionOutgo - mgp.IrNbPremiumReceipt
	}
}
func GetModifiedGMMShock(month int, shockBasis string) (models.ModifiedGMMShock, error) {
	var shock models.ModifiedGMMShock
	if month > 24 {
		month = 24
	}

	key := "shock_" + shockBasis + "_" + strconv.Itoa(month)

	cacheKey := key
	cached, found := PaaCache.Get(cacheKey)

	if found {
		shock = cached.(models.ModifiedGMMShock)
		return shock, nil
	} else {
		err := DB.Where("shock_basis = ? and projection_month = ?", shockBasis, month).First(&shock).Error
		if err != nil {
			return shock, err
		}
		PaaCache.Set(key, shock, 1)
		return shock, nil
	}
}

func SaveMGMMShockSetting(shockSetting models.ModifiedGMMShockSetting) (models.ModifiedGMMShockSetting, error) {
	DB.Save(&shockSetting)
	return shockSetting, nil
}

func GetMGMMShockSettings() ([]models.ModifiedGMMShockSetting, error) {
	var shockSettings []models.ModifiedGMMShockSetting
	err := DB.Find(&shockSettings).Error
	return shockSettings, err
}

func GetMGMMShockBases() ([]string, error) {
	var bases []string
	err := DB.Table("modified_gmm_shocks").Distinct("shock_basis").Pluck("shock_basis", &bases).Error
	return bases, err
}

func DeleteMGMMShockSetting(id int) error {
	err := DB.Where("id = ?", id).Delete(&models.ModifiedGMMShockSetting{}).Error
	if err != nil {
		return err
	}
	return nil
}

func SaveGMMModelPoints(v *multipart.FileHeader, year int, portfolioName string, portfolioId int, mpVersion string, user models.AppUser) (models.PAAYearVersion, error) {
	var paaYearVersion models.PAAYearVersion
	var delimiter rune
	delimiterFile, err := v.Open()
	if err != nil {
		return paaYearVersion, err
	}
	defer delimiterFile.Close()
	delimiter, err = utils.GetDelimiter(delimiterFile)

	file, err := v.Open()
	if err != nil {
		return paaYearVersion, err
	}

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.Comma = delimiter
	dec, _ := csvutil.NewDecoder(reader)
	dec.Header()

	err = DB.Where("year = ? and paa_portfolio_id = ? and mp_version = ?", year, portfolioId, mpVersion).Delete(&models.ModifiedGMMModelPoint{}).Error
	if err != nil {
		fmt.Println(err)
	}

	var pps []models.ModifiedGMMModelPoint
	for {
		var pp models.ModifiedGMMModelPoint
		if err := dec.Decode(&pp); err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}
		pp.Year = year
		pp.CreatedBy = user.UserName
		pp.PaaPortfolioID = portfolioId
		pp.PaaPortfolioName = portfolioName
		pp.MpVersion = mpVersion
		fmt.Println(pp)
		pps = append(pps, pp)
		//err = DB.Create(&pp).Error
		//if err != nil {
		//	log.Error().Msg(err.Error())
		//}
	}
	paaYearVersion.Year = year
	paaYearVersion.PortfolioId = portfolioId
	paaYearVersion.PortfolioName = portfolioName
	paaYearVersion.MpVersion = mpVersion
	paaYearVersion.Count = len(pps)
	DB.Save(&paaYearVersion)
	err = DB.CreateInBatches(&pps, 200).Error
	if err != nil {
		return paaYearVersion, err
	}
	go func() {
		GeneratePAAModelPointVariableStats(portfolioName, year, mpVersion)
	}()

	return paaYearVersion, nil
}

func GetAllGMMTableMetaData() map[string]interface{} {
	// Retrieve model point counts by year, parameters and shocks data?
	table := "modified_gmm_model_points"
	var metadata []models.TableMetaData
	var counts []models.ModelPointCount
	var results = make(map[string]interface{})

	query := fmt.Sprintf("select year, count(*) as count from %s  group by year", table)
	DB.Raw(query).Scan(&counts)

	var mpd models.TableMetaData

	mpd.TableType = "Model Points"
	jsonString, _ := json.Marshal(&counts)
	var countData []map[string]interface{}
	err := json.Unmarshal(jsonString, &countData)
	if err != nil {
		fmt.Println(err)
	}
	if len(counts) > 0 {
		mpd.Populated = true
	}
	mpd.Data = countData

	// Test additions
	var param = models.TableMetaData{TableType: "Parameters", Data: nil, Populated: true, TableName: "modified_gmm_parameters"}
	var shock = models.TableMetaData{TableType: "Shocks", Data: nil, Populated: true, TableName: "modified_gmm_shocks"}
	var yieldCurve = models.TableMetaData{TableType: "Yield Curve", Data: nil, Populated: true, TableName: "paa_yield_curves"}
	var premiumEarning = models.TableMetaData{TableType: "Premium Earning Pattern", Data: nil, Populated: true, TableName: "premium_earning_patterns"}
	var reinsurance = models.TableMetaData{TableType: "ReinsuranceBEL Parameters", Data: nil, Populated: true, TableName: "reinsurance_parameters"}
	var paaLapse = models.TableMetaData{TableType: "PAA Lapse", Data: nil, Populated: true, TableName: "paa_lapses"}
	metadata = append(metadata, param, shock, yieldCurve, premiumEarning, reinsurance, paaLapse)
	results["associated_tables"] = metadata
	results["model_points"] = mpd
	return results
}

func DeleteGMMTables(tableType string, yieldCurveCode string, yieldYear int, yieldMonth int, year int) error {
	var err error
	switch tableType {
	case "Model Points":
		err = DB.Delete(&models.ModifiedGMMModelPoint{}).Error
	case "Parameters":
		err = DB.Where("year = ?", year).Delete(&models.ModifiedGMMParameter{}).Error
	case "Shocks":
		err = DB.Where("year = ?").Delete(&models.ModifiedGMMShock{}).Error
	case "Yield Curve":
		err = DB.Where("yield_curve_code = ? and year = ? and  month = ?", yieldCurveCode, yieldYear, yieldMonth).Delete(&models.PaaYieldCurve{}).Error
	case "Premium Earning Pattern":
		err = DB.Where("year = ?").Delete(&models.PremiumEarningPattern{}).Error
	case "ReinsuranceBEL Parameters":
		err = DB.Where("year = ?").Delete(&models.ReinsuranceParameter{}).Error
	case "PAA Lapse":
		err = DB.Where("year = ?").Delete(&models.PAALapse{}).Error
	case "PAA Finance":
		err = DB.Where("year = ?").Delete(&models.PAAFinance{}).Error
	}
	return err
}

func SaveGMMTables(v *multipart.FileHeader, tableType string, year, month int, portfolioName, yieldCurveCode string, user models.AppUser) error {
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

	switch tableType {
	case "Parameters":
		DB.Where("year = ? ", year).Delete(&models.ModifiedGMMParameter{})
		for {
			var pp models.ModifiedGMMParameter
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			//pp.Year = year
			pp.CreatedBy = user.UserName
			err = DB.Create(&pp).Error
			if err != nil {
				log.Error().Msg(err.Error())
			}
		}
	case "ReinsuranceBEL Parameters":
		for {
			var pp models.ReinsuranceParameter
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			pp.Year = year
			pp.CreatedBy = user.UserName
			err = DB.Create(&pp).Error
			if err != nil {
				log.Error().Msg(err.Error())
			}
		}

	case "Shocks":
		DB.Where("id > 0").Delete(&models.ModifiedGMMShock{})
		for {
			var pp models.ModifiedGMMShock
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			pp.CreatedBy = user.UserName
			err = DB.Create(&pp).Error
			if err != nil {
				log.Error().Msg(err.Error())
			}
		}

	case "Yield Curve":
		for {
			var pp models.PaaYieldCurve
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			pp.Year = year
			pp.CreatedBy = user.UserName
			pp.Month = month
			//pp.ExpConfigurationName = portfolioName
			pp.YieldCurveCode = yieldCurveCode
			err = DB.Where("year = ? and proj_time = ? and month = ? and yield_curve_code = ?", pp.Year, pp.ProjectionTime, pp.Month, pp.YieldCurveCode).Delete(&models.PaaYieldCurve{}).Error
			if err != nil {
				log.Error().Msg(err.Error())
			}

			err = DB.Create(&pp).Error
			if err != nil {
				log.Error().Msg(err.Error())
			}
		}

	case "Premium Earning Pattern":
		for {
			var pp models.PremiumEarningPattern
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			pp.Year = year
			pp.CreatedBy = user.UserName
			if err != nil {
				log.Error().Msg(err.Error())
			}

			err = DB.Where("year = ? and duration_in_force = ? and premium_earning_pattern_code = ?", pp.Year, pp.DurationInForce, pp.PremiumEarningPatternCode).Delete(&models.PremiumEarningPattern{}).Error
			err = DB.Create(&pp).Error
			if err != nil {
				log.Error().Msg(err.Error())
			}
		}
	case "PAA Lapse":
		for {
			var pp models.PAALapse
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			pp.Year = year
			pp.CreatedBy = user.UserName
			if err != nil {
				log.Error().Msg(err.Error())
			}

			err = DB.Where("year = ? and duration_in_force_month = ? and product_code = ?", pp.Year, pp.DurationInForceMonth, pp.ProductCode).Delete(&models.PAALapse{}).Error
			err = DB.Create(&pp).Error
			if err != nil {
				log.Error().Msg(err.Error())
			}
		}
	case "PAA Finance":
		err = DB.Where("year = ?", year).Delete(&models.PAAFinance{}).Error
		for {
			var pp models.PAAFinance
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			pp.Year = year
			pp.CreatedBy = user.UserName
			pp.PortfolioName = portfolioName
			if err != nil {
				log.Error().Msg(err.Error())
			}

			err = DB.Create(&pp).Error
			if err != nil {
				log.Error().Msg(err.Error())
			}
		}
	}

	return nil
}

func SaveGMMParameters(v *multipart.FileHeader, year int) error {
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

	err = DB.Where("year = ?", year).Delete(&models.ModifiedGMMParameter{}).Error
	if err != nil {
		fmt.Println(err)
	}

	for {
		var pp models.ModifiedGMMParameter
		if err := dec.Decode(&pp); err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}
		pp.Year = year
		fmt.Println(pp)
		err = DB.Create(&pp).Error
		if err != nil {
			log.Error().Msg(err.Error())
		}
	}

	return nil
}

func GetGMMParameterYears() []int {
	type Result struct {
		Year int
	}
	var result []Result

	var years []int
	err := DB.Table("modified_gmm_parameters").Distinct("year").Find(&result).Error
	if err != nil {
		fmt.Println(err)
	}

	for _, r := range result {
		years = append(years, r.Year)
	}

	return years
}

func GetAllRunJobs() []models.GMMRunSetting {
	var jobs []models.GMMRunSetting

	err := DB.Order("id desc").Find(&jobs).Error
	if err != nil {
		fmt.Println(err)
	}
	return jobs
}

func GetAllRunJobsV2() []models.MgmmRun {
	var jobs []models.MgmmRun

	err := DB.Preload("GMMRunSettings").Order("id desc").Find(&jobs).Error
	if err != nil {
		fmt.Println(err)
	}
	return jobs

}

func GetMGMMProjectionsByRun(id int) []models.AggregatedModifiedGMMProjection {
	//var runSetting models.GMMRunSetting
	//DB.Where("id = ?", id).Find(&runSetting)

	var projections []models.AggregatedModifiedGMMProjection
	DB.Where("job_run_id = ?", id).Order("product_code").Order("policy_number").Order("projection_month asc").Limit(5000).Find(&projections)
	return projections
}

func GetGMMProjectionsForGroup(id int, group string) []models.AggregatedModifiedGMMProjection {
	var projections []models.AggregatedModifiedGMMProjection
	DB.Where("run_id = ? and ifrs17_group = ?", id, group).Find(&projections)
	return projections
}

func GetMGMMIFRS17GroupsForProjection(id int) []string {
	var groups []string

	DB.Table("modified_gmm_scoped_aggregations").Distinct("ifrs17_group").Where("job_run_id = ?", id).Pluck("ifrs17_group", &groups)

	return groups
}

func GetGMMRunJob(id int) models.GMMRunSetting {
	var job models.GMMRunSetting
	DB.Where("id = ?", id).Find(&job)
	return job
}

func GetMGMMScopedProjectionsByRun(id int) []models.ModifiedGMMScopedAggregation {
	var projections []models.ModifiedGMMScopedAggregation
	DB.Where("job_run_id = ?", id).Order("ifrs17_group asc").Order("projection_month asc").Limit(5000).Find(&projections)
	return projections
}

func GetMGMMScopedProjectionsByRunAndGroup(id int, group string) []models.ModifiedGMMScopedAggregation {
	var projections []models.ModifiedGMMScopedAggregation
	DB.Where("job_run_id = ? and ifrs17_group = ?", id, group).Order("projection_month asc").Limit(5000).Find(&projections)
	return projections
}

func GetGMMRunSettings(runId int) models.GMMRunSetting {
	var run models.GMMRunSetting
	DB.Where("id = ?", runId).Find(&run)
	// Get the shock setting data
	var shockSettings models.ModifiedGMMShockSetting
	DB.Where("id = ?", run.ShockSettingID).Find(&shockSettings)
	run.ShockSettings = shockSettings
	return run
}

func DeleteMGMMProjectionsRun(id int) error {
	var run models.MgmmRun
	err := DB.Where("id = ?", id).Find(&run).Error

	if err != nil {
		fmt.Println(err)
	}

	err = DB.Where("mgmm_run_id = ?", run.ID).Delete(&models.GMMRunSetting{}).Error
	//Delete related projection runs and scoped run results
	err = DB.Where("run_id=?", run.ID).Delete(&models.ModifiedGMMScopedAggregation{}).Error
	err = DB.Where("run_id=?", run.ID).Delete(&models.AggregatedModifiedGMMProjection{}).Error
	err = DB.Where("id = ?", run.ID).Delete(&models.MgmmRun{}).Error

	return err
}

func GetPortfolioModelPoints(portfolioId, year int, version string) ([]models.ModifiedGMMModelPoint, error) {
	var results []models.ModifiedGMMModelPoint
	if year == 0 {
		err := DB.Where("paa_portfolio_id = ?", portfolioId).Limit(2500).Find(&results).Error
		return results, err
	}
	err := DB.Where("paa_portfolio_id = ? and year = ? and mp_version = ? ", portfolioId, year, version).Limit(2500).Find(&results).Error
	return results, err
}

func GetAvailableModelPointVersions(portfolioName, year string) []string {
	//get distinct model point versions for a portfolio and year
	var versions []string
	err := DB.Table("modified_gmm_model_points").Distinct("mp_version").Where("paa_portfolio_name = ? and year = ?", portfolioName, year).Pluck("mp_version", &versions).Error
	if err != nil {
		fmt.Println(err)
	}
	return versions
}

func DeleteModelPointsForPortfolioAndYear(portfolioId int, year int, version string) error {
	if version == "null" {
		version = ""
	}
	err := DB.Where("paa_portfolio_id = ? and year = ? and mp_version = ?", portfolioId, year, version).Delete(&models.ModifiedGMMModelPoint{}).Error
	//delete paayearversion
	err = DB.Where("portfolio_id = ? and year = ? and mp_version = ?", portfolioId, year, version).Delete(&models.PAAYearVersion{}).Error

	err = DB.Where("portfolio_id = ? and year = ? and mp_version = ?", portfolioId, year, version).Delete(&models.PaaModelPointVariableStats{}).Error
	return err
}

func CheckFinanceYearExists(portfolioName string, year int) bool {
	var finance models.PAAFinance
	DB.Where("portfolio_name = ? and year = ?", portfolioName, year).First(&finance)

	if finance.ID == 0 {
		return false
	} else {
		return true
	}
}

func GetGMMTableData(tableType string) []map[string]interface{} {
	var results []map[string]interface{}

	switch tableType {
	case "parameters":
		var params []models.ModifiedGMMParameter
		DB.Table("modified_gmm_parameters").Find(&params)
		b, _ := json.Marshal(&params)
		err := json.Unmarshal(b, &results)
		if err != nil {
			fmt.Println(err)
		}
	case "reinsurancebelparameters":
		var params []models.ReinsuranceParameter
		DB.Table("reinsurance_parameters").Find(&params)
		b, _ := json.Marshal(&params)
		err := json.Unmarshal(b, &results)
		if err != nil {
			fmt.Println(err)
		}

	case "shocks":
		var shocks []models.ModifiedGMMShock
		DB.Table("modified_gmm_shocks").Find(&shocks)
		b, _ := json.Marshal(&shocks)
		err := json.Unmarshal(b, &results)
		if err != nil {
			fmt.Println(err)
		}

	case "yieldcurve":
		var yieldCurve []models.PaaYieldCurve
		DB.Find(&yieldCurve)
		b, _ := json.Marshal(&yieldCurve)
		err := json.Unmarshal(b, &results)
		if err != nil {
			fmt.Println(err)
		}

	case "premiumearningpattern":
		var premiumEarningPattern []models.PremiumEarningPattern
		DB.Find(&premiumEarningPattern)
		b, _ := json.Marshal(&premiumEarningPattern)
		err := json.Unmarshal(b, &results)
		if err != nil {
			fmt.Println(err)
		}
	case "paalapse":
		var paalapse []models.PAALapse
		DB.Find(&paalapse)
		b, _ := json.Marshal(&paalapse)
		err := json.Unmarshal(b, &results)
		if err != nil {
			fmt.Println(err)
		}
	case "paafinance":
		var paalapse []models.PAAFinance
		DB.Find(&paalapse)
		b, _ := json.Marshal(&paalapse)
		err := json.Unmarshal(b, &results)
		if err != nil {
			fmt.Println(err)
		}

	case "modelpoints":
		var modelpoints []models.ModifiedGMMModelPoint
		DB.Table("modified_gmm_model_points").Find(&modelpoints)
		b, _ := json.Marshal(&modelpoints)
		err := json.Unmarshal(b, &results)
		if err != nil {
			fmt.Println(err)
		}
	}

	return results
}

func UpdateScopedModifiedGMMProjections(modifiedGMMProjections []models.AggregatedModifiedGMMProjection, scopedModPAA *map[string]models.ModifiedGMMScopedAggregation) {
	for _, modifiedgmmprojection := range modifiedGMMProjections {
		key := strconv.Itoa(modifiedgmmprojection.ProjectionMonth) + "_" + modifiedgmmprojection.IFRS17Group
		var aggProj = models.ModifiedGMMScopedAggregation{}
		err := copier.Copy(&aggProj, &modifiedgmmprojection)

		aggProj.ID = 0
		if err != nil {
			fmt.Println(err)
		}

		mutex.Lock()

		agg, exists := (*scopedModPAA)[key]

		//if i, ok := utils.ModifiedGMMScopedAggregatedProjectionsContain(modifiedGMMScopedAggregations, aggProj); ok {
		if exists {
			mutable := reflect.ValueOf(&agg).Elem()
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
			(*scopedModPAA)[key] = agg

		} else {
			(*scopedModPAA)[key] = aggProj
			//modifiedGMMScopedAggregations = append(modifiedGMMScopedAggregations, aggProj)
		}

		mutex.Unlock()
	}
}

func UpdateAggregatedModifiedGMMProjection(modifiedGMMProjections []models.AggregatedModifiedGMMProjection, aggModPAA *map[string]models.AggregatedModifiedGMMProjection) {
	for _, modifiedGMMProjection := range modifiedGMMProjections {
		key := strconv.Itoa(modifiedGMMProjection.ProjectionMonth) + "_" + modifiedGMMProjection.ProductCode
		var aggProj = models.AggregatedModifiedGMMProjection{}
		err := copier.Copy(&aggProj, &modifiedGMMProjection)

		aggProj.ID = 0
		if err != nil {
			fmt.Println(err)
		}
		mutex.Lock()
		agg, exists := (*aggModPAA)[key]

		//if i, ok := utils.AggregatedModifiedGMMProjectionsContain(aggregatedModifiedGMMProjections, aggProj); ok {
		if exists {
			mutable := reflect.ValueOf(&agg).Elem()
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
			(*aggModPAA)[key] = agg
		} else {
			(*aggModPAA)[key] = aggProj
			//aggregatedModifiedGMMProjections = append(aggregatedModifiedGMMProjections, aggProj)
		}
		mutex.Unlock()
	}
}

func CreatePaaPortfolio(portfolio models.PaaPortfolio) models.PaaPortfolio {
	err := DB.Create(&portfolio).Error
	if err != nil {
		fmt.Println(err)
	}
	return portfolio
}

func DeletePortfolio(portfolioId int) error {
	err := DB.Where("paa_portfolio_id = ?", portfolioId).Delete(&models.ModifiedGMMModelPoint{}).Error
	if err != nil {
		fmt.Println(err)
	}
	err = DB.Where("id = ?", portfolioId).Delete(&models.PaaPortfolio{}).Error
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func GetPortfolios() ([]models.PaaPortfolio, error) {
	var portfolios = make([]models.PaaPortfolio, 0)
	var err error
	err = DB.Find(&portfolios).Error

	for i, _ := range portfolios {
		var yearMps []models.PAAYearVersion
		DB.Where("portfolio_id = ?", portfolios[i].ID).Find(&yearMps)
		var mpYears []int
		err := DB.Table("paa_year_versions").Distinct("year").Where("portfolio_id = ?", portfolios[i].ID).Order("year desc").Pluck("year as model_point_year", &mpYears).Error
		if err != nil {
			fmt.Println(err)
		}

		portfolios[i].YearVersions = yearMps
		portfolios[i].ModelPointYears = mpYears
	}
	return portfolios, err
}

func GetRunPortfolioProductGroup(runDate, portfolioName, productCode, group string) map[string]interface{} {
	var result = make(map[string]interface{})
	var buildUps []models.PAABuildUp
	var list []models.PAABuildUp
	DB.Distinct("ifrs17_group").Where("run_date = ? and portfolio_name = ? and product_code = ? and ifrs17_group = ?", runDate, portfolioName, productCode, group).Find(&list)
	result["list"] = list

	DB.Where("run_date = ? AND portfolio_name = ? and product_code = ? and ifrs17_group = ?", runDate, portfolioName, productCode, group).Find(&buildUps)
	result["portfolios"] = buildUps

	return result
}

func GetRunPortfolioProductGroups(runDate, portfolioName, productCode string) map[string]interface{} {
	var result = make(map[string]interface{})
	var buildUps []models.PAABuildUp
	var list []models.PAABuildUp
	DB.Distinct("ifrs17_group").Where("run_date = ? and portfolio_name = ? and product_code = ?", runDate, portfolioName, productCode).Find(&list)
	result["list"] = list
	query := fmt.Sprintf("SELECT name,portfolio_name,product_code, sum(variable_change) as variable_change,sum(paa_lrc_buildup) as paa_lrc_buildup,sum(initial_recognition_loss_component) as initial_recognition_loss_component,sum(loss_component_unwind) as loss_component_unwind,sum(loss_component_buildup) as loss_component_buildup,sum(loss_component_adjustment) as loss_component_adjustment,sum(initial_recognition_loss_recovery) as initial_recognition_loss_recovery,sum(loss_recovery_unwind) as loss_recovery_unwind, sum(loss_recovery_buildup) as loss_recovery_buildup,sum(loss_recovery_adjustment) as loss_recovery_adjustment,sum(dac_buildup) as dac_buildup,sum(modified_gmm_bel) as modified_gmm_bel,sum(modified_gmm_risk_adjustment) as modified_gmm_risk_adjustment,sum(modified_gmm_reserve) as modified_gmm_reserve, sum(insurance_revenue) as insurance_revenue,sum(insurance_service_expense) as insurance_service_expense,sum(paa_reinsurance_lrc_buildup) as paa_reinsurance_lrc_buildup, sum(paa_reinsurance_dac_buildup) as paa_reinsurance_dac_buildup, sum(loss_recovery_buildup) as loss_recovery_buildup FROM paa_build_ups where run_date = '%s' and portfolio_name =  '%s' and product_code = '%s'  group by name, portfolio_name, product_code", runDate, portfolioName, productCode)
	DB.Raw(query).Find(&buildUps)
	result["portfolios"] = buildUps

	return result
}

func GetRunPortfolioProducts(runDate string, portfolioName string) map[string]interface{} {
	var result = make(map[string]interface{})
	var buildUps []models.PAABuildUp
	var list []models.PAABuildUp
	DB.Distinct("product_code").Where("run_date = ? and portfolio_name = ?", runDate, portfolioName).Find(&list)
	result["list"] = list
	query := fmt.Sprintf("SELECT name,portfolio_name, sum(variable_change) as variable_change,sum(paa_lrc_buildup) as paa_lrc_buildup,sum(initial_recognition_loss_component) as initial_recognition_loss_component,sum(loss_component_unwind) as loss_component_unwind,sum(loss_component_buildup) as loss_component_buildup,sum(loss_component_adjustment) as loss_component_adjustment,sum(initial_recognition_loss_recovery) as initial_recognition_loss_recovery,sum(loss_recovery_unwind) as loss_recovery_unwind, sum(loss_recovery_buildup) as loss_recovery_buildup,sum(loss_recovery_adjustment) as loss_recovery_adjustment,sum(dac_buildup) as dac_buildup,sum(modified_gmm_bel) as modified_gmm_bel,sum(modified_gmm_risk_adjustment) as modified_gmm_risk_adjustment,sum(modified_gmm_reserve) as modified_gmm_reserve, sum(insurance_revenue) as insurance_revenue,sum(insurance_service_expense) as insurance_service_expense,sum(paa_reinsurance_lrc_buildup) as paa_reinsurance_lrc_buildup, sum(paa_reinsurance_dac_buildup) as paa_reinsurance_dac_buildup, sum(loss_recovery_buildup) as loss_recovery_buildup FROM paa_build_ups where run_date = '%s' and portfolio_name =  '%s'  group by name, portfolio_name", runDate, portfolioName)
	DB.Raw(query).Find(&buildUps)
	result["portfolios"] = buildUps

	return result
}

func GetRunPortfolios(runDate string) map[string]interface{} {
	var result = make(map[string]interface{})

	var list []models.PAABuildUp
	var portfolios []models.PAABuildUp
	DB.Distinct("portfolio_name").Where("run_date = ?", runDate).Find(&list)
	result["list"] = list
	query := fmt.Sprintf("SELECT name, sum(variable_change) as variable_change,sum(paa_lrc_buildup) as paa_lrc_buildup,sum(initial_recognition_loss_component) as initial_recognition_loss_component,sum(loss_component_unwind) as loss_component_unwind,sum(loss_component_buildup) as loss_component_buildup,sum(loss_component_adjustment) as loss_component_adjustment,sum(initial_recognition_loss_recovery) as initial_recognition_loss_recovery,sum(loss_recovery_unwind) as loss_recovery_unwind, sum(loss_recovery_buildup) as loss_recovery_buildup,sum(loss_recovery_adjustment) as loss_recovery_adjustment,sum(dac_buildup) as dac_buildup,sum(modified_gmm_bel) as modified_gmm_bel,sum(modified_gmm_risk_adjustment) as modified_gmm_risk_adjustment,sum(modified_gmm_reserve) as modified_gmm_reserve, sum(insurance_revenue) as insurance_revenue,sum(insurance_service_expense) as insurance_service_expense,sum(paa_reinsurance_lrc_buildup) as paa_reinsurance_lrc_buildup, sum(paa_reinsurance_dac_buildup) as paa_reinsurance_dac_buildup, sum(loss_recovery_buildup) as loss_recovery_buildup FROM paa_build_ups where run_date = '%s'  group by name", runDate)
	DB.Raw(query).Scan(&portfolios)
	result["portfolios"] = portfolios
	return result
}

func GetValidPortfolios() []models.PaaPortfolio {
	var portfolios []models.PaaPortfolio
	var validPortfolios []models.PaaPortfolio
	DB.Find(&portfolios)
	for _, portfolio := range portfolios {
		var count int64
		DB.Table("modified_gmm_model_points").Where("paa_portfolio_id=?", portfolio.ID).Count(&count)
		if count > 0 {
			validPortfolios = append(validPortfolios, portfolio)
		}
	}
	return validPortfolios
}

func GetModelpointYearsForPortfolio(portfolioId int) []int {
	type Result struct {
		Year int
	}
	var result []Result

	var years []int
	err := DB.Table("modified_gmm_model_points").Where("paa_portfolio_id=?", portfolioId).Distinct("year").Find(&result).Error
	if err != nil {
		fmt.Println(err)
	}

	for _, r := range result {
		years = append(years, r.Year)
	}

	return years
}

func GetPAAYieldCurveCodes(year string) []string {
	var codes []string

	DB.Raw("select distinct yield_curve_code from paa_yield_curve where year = ?", year).Scan(&codes)
	return codes
}

func GetPAAYieldCurveMonths(year string, code string) []string {
	var months []string
	DB.Raw("select distinct month from paa_yield_curve where year = ? and yield_curve_code = ?", year, code).Scan(&months)
	return months
}

func GetIbnrYieldCurveMonths(year string, basis string, parameterYear string, portfolioId string) []string {
	var ibnrPortfolio models.LicPortfolio
	DB.Where("id = ?", portfolioId).Find(&ibnrPortfolio)

	var licParameter models.LICParameter
	DB.Where("year = ? and portfolio_name = ? and basis = ?", parameterYear, ibnrPortfolio.Name, basis).Find(&licParameter)

	var months []string
	DB.Table("ibnr_yield_curves").Distinct("month").Where("year = ? and yield_curve_code  = ?", year, licParameter.YieldCurveCode).Pluck("month", &months)
	return months
}

func countCalendarMonths(start, end time.Time) int {
	// Ensure start is before end
	if start.After(end) {
		start, end = end, start
	}

	// If start and end are in the same month and year, count 1 month
	if start.Year() == end.Year() && start.Month() == end.Month() {
		return 1
	}

	// Otherwise, count inclusive months
	months := (end.Year()-start.Year())*12 + int(end.Month()) - int(start.Month()) + 1
	return months
}

func endOfMonth(t time.Time, months int) time.Time {
	// Move to the first of the month to avoid rollover issues
	firstOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())

	// Add months
	target := firstOfMonth.AddDate(0, months+1, 0) // move to first of next month after target
	endOfTargetMonth := target.AddDate(0, 0, -1)   // subtract one day

	return endOfTargetMonth
}

func maxDate(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}
