package services

import (
	"api/log"
	"api/models"
	"fmt"
	"strconv"
	"time"
)

func RunLic(runs []models.LicRunSetting) {
	for _, run := range runs {
		var err error
		var licVarSet models.LicVariableSet

		DB.Preload("LicVariables").Where("id = ?", run.LicConfigurationId).Find(&licVarSet)

		//check if any of the licvariables has missing run data
		for _, licVar := range licVarSet.LicVariables {
			if licVar.RunId > 0 {
				var runType = licVarSet.RunType

				if runType == "ibnr" {
					var ibnrs models.IBNRRunSetting
					DB.Where("id = ?", licVar.RunId).Find(&ibnrs)
					if ibnrs.RunDate == "" || ibnrs.PortfolioName == "" {
						log.Info("Missing run data for ", licVar.Name)
						run.ProcessingStatus = "failed"
						run.RunFailureReason = "Missing data for the run assigned to step " + licVar.Name
						run.RunFailed = true
						run.CreationDate = time.Now()
						DB.Save(&run)
						break
					}
				}
				if runType == "cash_flows" {
					var cf models.ProjectionJob
					DB.Where("id = ?", licVar.RunId).Find(&cf)
					if cf.ID == 0 || cf.RunDate == "" {
						log.Info("Missing run data for ", licVar.Name)
						run.ProcessingStatus = "failed"
						run.RunFailureReason = "Missing data for the run assigned to step " + licVar.Name
						run.RunFailed = true
						run.CreationDate = time.Now()
						err := DB.Save(&run).Error
						if err != nil {
							log.Error(err)
						}
						break
					}
				}
			}
		}

		if run.RunFailed {
			continue
		}

		//Save the run setting
		run.CreationDate = time.Now()
		run.ProcessingStatus = "running"

		//Delete previous runs with this confguration
		DB.Where("run_date = ? and lic_configuration_name = ?",
			run.RunDate,
			run.LicConfigurationName).Delete(&models.LicRunSetting{})
		err = DB.Create(&run).Error
		if err != nil {
			log.Error(err)
		}
		start := time.Now()

		var products []string
		var ibnrs models.IBNRRunSetting
		DB.Where("id = ?", licVarSet.LicVariables[8].RunId).Find(&ibnrs)
		DB.Table("ibnr_reserve_reports").Select("distinct product_code").Where("run_date = ? and portfolio_name = ?", ibnrs.RunDate, ibnrs.PortfolioName).Pluck("product_code", &products)
		var licbuildups []models.LicBuildupResult

		oldClaimsYear, _ := strconv.Atoi(run.OpeningBalanceDate[:4])
		oldClaimsMonth, _ := strconv.Atoi(run.OpeningBalanceDate[5:])
		oldClaimsCutOff := float64(oldClaimsYear) + float64(oldClaimsMonth)/1000.0
		//prevYear, _ := strconv.Atoi(prevYearStr)
		//prevYear = prevYear - 1
		yearSearchArg := run.OpeningBalanceDate //strconv.Itoa(prevYear) + "-" + prevYearMth
		fmt.Println(yearSearchArg)

		for _, product := range products {

			var prevBsr models.BalanceSheetRecord
			err = DB.Where("date = ? and product_code = ? and measurement_type=? and ifrs17_group=?", yearSearchArg, product, "LIC", product).Find(&prevBsr).Error
			if err != nil {
				fmt.Println("query bsr error:", err)
			}

			var licParams models.Lic2Parameter
			err = DB.Where("product_code = ? and year=? and version=?", product, run.LicParameterYear, run.LicParameterVersion).Find(&licParams).Error
			if err != nil {
				fmt.Println("query lic parameter error:", err)
			}

			var ibnrBS float64
			var ibnrRiskAdjustmentBS float64

			var previbnr, previbnrAt12 float64
			var previbnrRiskAdjustment, previbnrRiskAdjustmentAt12 float64

			for i, licVar := range licVarSet.LicVariables {
				var licbuildup models.LicBuildupResult
				fmt.Println("product: ", product)
				var ibnrRunSetting models.IBNRRunSetting
				DB.Where("id = ?", licVar.RunId).Find(&ibnrRunSetting)
				fmt.Println("ibnrRunSetting", ibnrRunSetting)

				var ibnrResult models.IbnrReserveReport
				err = DB.Where("run_id=? and product_code=? ", licVar.RunId, product).Find(&ibnrResult).Error
				if err != nil {
					log.Error(err)
				}

				licbuildup.ConfigurationName = run.LicConfigurationName
				licbuildup.RunDate = run.RunDate
				licbuildup.RunId = run.ID
				licbuildup.Name = licVar.Name
				licbuildup.PortfolioId = ibnrRunSetting.PortfolioId
				licbuildup.PortfolioName = ibnrRunSetting.PortfolioName
				licbuildup.ProductCode = product
				licbuildup.IFRS17Group = product

				if licVar.RunId > 0 {

					licbuildup.PortfolioId = ibnrResult.LicPortfolioId
					licbuildup.IBNR = ibnrResult.IbnrBel
					licbuildup.RiskAdjustment = ibnrResult.IBNRRiskAdjustment
					licbuildup.LicBuildup = licbuildup.IBNR + licbuildup.RiskAdjustment
					licbuildup.IBNRAt12 = ibnrResult.IbnrBelAt12
					licbuildup.IBNRRiskAdjustmentAt12 = ibnrResult.IbnrRiskAdjustmentAt12
					ibnrBS = ibnrResult.IbnrBel
					ibnrRiskAdjustmentBS = ibnrResult.IBNRRiskAdjustment

					if i == 0 {
						previbnrAt12 = 0
						previbnrRiskAdjustmentAt12 = 0
						previbnr = 0
						previbnrRiskAdjustment = 0
						licbuildup.VariableChange = 0
						//licbuildup.IBNR = prevBsr.IBNR
						licbuildup.ReportedClaims = prevBsr.OutstandingClaimsReserve
						//licbuildup.IBNRRiskAdjustment = prevBsr.IBNRRiskAdjustment
						licbuildup.CashBack = prevBsr.CashBack
						licbuildup.UnadjustedLossAdjustmentExpenses = 0 //(licbuildup.IBNR + licbuildup.ReportedClaims*0.5) * 0.
						licbuildup.Treaty1Ibnr = prevBsr.Treaty1IBNR
						licbuildup.Treaty1IbnrRiskAdjustment = prevBsr.Treaty1IBNRRiskAdjustment
						licbuildup.Treaty2Ibnr = prevBsr.Treaty2IBNR
						licbuildup.Treaty2IbnrRiskAdjustment = prevBsr.Treaty2IBNRRiskAdjustment
						licbuildup.Treaty3Ibnr = prevBsr.Treaty3IBNR
						licbuildup.Treaty3IbnrRiskAdjustment = prevBsr.Treaty3IBNRRiskAdjustment
					} else {
						if licbuildup.Name == "Interest Accretion" {
							licbuildup.IBNR = 0
							licbuildup.ReportedClaims = 0
							licbuildup.RiskAdjustment = 0
							licbuildup.CashBack = 0
							licbuildup.UnadjustedLossAdjustmentExpenses = 0 //(licbuildup.IBNR + licbuildup.ReportedClaims*0.5) * 0.05
							licbuildup.Treaty1Ibnr = 0
							licbuildup.Treaty1IbnrRiskAdjustment = 0
							licbuildup.Treaty2Ibnr = 0
							licbuildup.Treaty2IbnrRiskAdjustment = 0
							licbuildup.Treaty3Ibnr = 0
							licbuildup.Treaty3IbnrRiskAdjustment = 0
						}

						if licbuildup.Name == "Expected IBNR Release" {
							licbuildup.IBNR = licbuildup.IBNRAt12 - licbuildup.IBNR
							licbuildup.ReportedClaims = 0 //-licbuildup.IBNR
							licbuildup.RiskAdjustment = licbuildup.IBNRRiskAdjustmentAt12 - licbuildup.RiskAdjustment
							licbuildup.CashBack = -licParams.EstimatedCashbackReserveRelease
							licbuildup.UnadjustedLossAdjustmentExpenses = 0 //(licbuildup.IBNR + licbuildup.ReportedClaims*0.5) * 0.05

							if prevBsr.IBNR > 0 {
								licbuildup.Treaty1Ibnr = licbuildup.IBNR * (prevBsr.Treaty1IBNR / prevBsr.IBNR)
								licbuildup.Treaty2Ibnr = licbuildup.IBNR * (prevBsr.Treaty2IBNR / prevBsr.IBNR)
								licbuildup.Treaty3Ibnr = licbuildup.IBNR * (prevBsr.Treaty3IBNR / prevBsr.IBNR)
							}

							if prevBsr.IBNRRiskAdjustment > 0 {
								licbuildup.Treaty1IbnrRiskAdjustment = licbuildup.RiskAdjustment * (prevBsr.Treaty1IBNRRiskAdjustment / prevBsr.IBNRRiskAdjustment)
								licbuildup.Treaty2IbnrRiskAdjustment = licbuildup.RiskAdjustment * (prevBsr.Treaty2IBNRRiskAdjustment / prevBsr.IBNRRiskAdjustment)
								licbuildup.Treaty3IbnrRiskAdjustment = licbuildup.RiskAdjustment * (prevBsr.Treaty3IBNRRiskAdjustment / prevBsr.IBNRRiskAdjustment)
							}

						}
						if licbuildup.Name == "Non Financial Assumption Change" {
							licbuildup.IBNR = licbuildup.IBNRAt12 - previbnrAt12
							licbuildup.ReportedClaims = 0
							licbuildup.RiskAdjustment = licbuildup.IBNRRiskAdjustmentAt12 - previbnrRiskAdjustmentAt12
							licbuildup.CashBack = 0
							licbuildup.UnadjustedLossAdjustmentExpenses = 0 //(licbuildup.IBNR + licbuildup.ReportedClaims*0.5) * 0.05

							if prevBsr.IBNR > 0 {
								licbuildup.Treaty1Ibnr = licbuildup.IBNR * (prevBsr.Treaty1IBNR / prevBsr.IBNR)
								licbuildup.Treaty2Ibnr = licbuildup.IBNR * (prevBsr.Treaty2IBNR / prevBsr.IBNR)
								licbuildup.Treaty3Ibnr = licbuildup.IBNR * (prevBsr.Treaty3IBNR / prevBsr.IBNR)
							}

							if prevBsr.IBNRRiskAdjustment > 0 {
								licbuildup.Treaty1IbnrRiskAdjustment = licbuildup.RiskAdjustment * (prevBsr.Treaty1IBNRRiskAdjustment / prevBsr.IBNRRiskAdjustment)
								licbuildup.Treaty2IbnrRiskAdjustment = licbuildup.RiskAdjustment * (prevBsr.Treaty2IBNRRiskAdjustment / prevBsr.IBNRRiskAdjustment)
								licbuildup.Treaty3IbnrRiskAdjustment = licbuildup.RiskAdjustment * (prevBsr.Treaty3IBNRRiskAdjustment / prevBsr.IBNRRiskAdjustment)
							}

						}

						if licbuildup.Name == "Financial Assumption Change" {
							licbuildup.IBNR = licbuildup.IBNRAt12 - previbnrAt12
							licbuildup.ReportedClaims = 0
							licbuildup.RiskAdjustment = licbuildup.IBNRRiskAdjustmentAt12 - previbnrRiskAdjustmentAt12
							licbuildup.CashBack = 0
							licbuildup.UnadjustedLossAdjustmentExpenses = 0 //(licbuildup.IBNR + licbuildup.ReportedClaims*0.5) * 0.05

							if prevBsr.IBNR > 0 {
								licbuildup.Treaty1Ibnr = licbuildup.IBNR * (prevBsr.Treaty1IBNR / prevBsr.IBNR)
								licbuildup.Treaty2Ibnr = licbuildup.IBNR * (prevBsr.Treaty2IBNR / prevBsr.IBNR)
								licbuildup.Treaty3Ibnr = licbuildup.IBNR * (prevBsr.Treaty3IBNR / prevBsr.IBNR)
							}

							if prevBsr.IBNRRiskAdjustment > 0 {
								licbuildup.Treaty1IbnrRiskAdjustment = licbuildup.RiskAdjustment * (prevBsr.Treaty1IBNRRiskAdjustment / prevBsr.IBNRRiskAdjustment)
								licbuildup.Treaty2IbnrRiskAdjustment = licbuildup.RiskAdjustment * (prevBsr.Treaty2IBNRRiskAdjustment / prevBsr.IBNRRiskAdjustment)
								licbuildup.Treaty3IbnrRiskAdjustment = licbuildup.RiskAdjustment * (prevBsr.Treaty3IBNRRiskAdjustment / prevBsr.IBNRRiskAdjustment)
							}

						}
						if licbuildup.Name == "Exchange Rate Effect" {
							licbuildup.IBNR = ibnrResult.IbnrBelAt12 - previbnrAt12
							licbuildup.ReportedClaims = 0
							licbuildup.RiskAdjustment = ibnrResult.IbnrRiskAdjustmentAt12 - previbnrRiskAdjustmentAt12
							licbuildup.CashBack = 0
							licbuildup.UnadjustedLossAdjustmentExpenses = 0 //(licbuildup.IBNR + licbuildup.ReportedClaims*0.5) * 0.05

							if prevBsr.IBNR > 0 {
								licbuildup.Treaty1Ibnr = licbuildup.IBNR * (prevBsr.Treaty1IBNR / prevBsr.IBNR)
								licbuildup.Treaty2Ibnr = licbuildup.IBNR * (prevBsr.Treaty2IBNR / prevBsr.IBNR)
								licbuildup.Treaty3Ibnr = licbuildup.IBNR * (prevBsr.Treaty3IBNR / prevBsr.IBNR)
							}

							if prevBsr.IBNRRiskAdjustment > 0 {
								licbuildup.Treaty1IbnrRiskAdjustment = licbuildup.RiskAdjustment * (prevBsr.Treaty1IBNRRiskAdjustment / prevBsr.IBNRRiskAdjustment)
								licbuildup.Treaty2IbnrRiskAdjustment = licbuildup.RiskAdjustment * (prevBsr.Treaty2IBNRRiskAdjustment / prevBsr.IBNRRiskAdjustment)
								licbuildup.Treaty3IbnrRiskAdjustment = licbuildup.RiskAdjustment * (prevBsr.Treaty3IBNRRiskAdjustment / prevBsr.IBNRRiskAdjustment)
							}

						}

						if licbuildup.Name == "Change in Non-Performance Rating" {
							licbuildup.IBNR = ibnrResult.IbnrBelAt12 - previbnrAt12
							licbuildup.ReportedClaims = 0
							licbuildup.RiskAdjustment = ibnrResult.IbnrRiskAdjustmentAt12 - previbnrRiskAdjustmentAt12
							licbuildup.CashBack = 0
							licbuildup.UnadjustedLossAdjustmentExpenses = 0

							if licParams.PrevTreaty1NonPerformanceRate != 0 {
								licbuildup.Treaty1Ibnr = ((1-licParams.CurrTreaty1NonPerformanceRate)/(1-licParams.PrevTreaty1NonPerformanceRate) - 1) * (licbuildups[0].Treaty1Ibnr + licbuildups[1].Treaty1Ibnr + licbuildups[2].Treaty1Ibnr + licbuildups[3].Treaty1Ibnr + licbuildups[4].Treaty1Ibnr + licbuildups[5].Treaty1Ibnr)
								licbuildup.Treaty1IbnrRiskAdjustment = ((1-licParams.CurrTreaty1NonPerformanceRate)/(1-licParams.PrevTreaty1NonPerformanceRate) - 1) * (licbuildups[0].Treaty1IbnrRiskAdjustment + licbuildups[1].Treaty1IbnrRiskAdjustment + licbuildups[2].Treaty1IbnrRiskAdjustment + licbuildups[3].Treaty1IbnrRiskAdjustment + licbuildups[4].Treaty1IbnrRiskAdjustment + licbuildups[5].Treaty1IbnrRiskAdjustment)

							}
							if licParams.PrevTreaty2NonPerformanceRate != 0 {
								licbuildup.Treaty2Ibnr = ((1-licParams.CurrTreaty2NonPerformanceRate)/(1-licParams.PrevTreaty2NonPerformanceRate) - 1) * (licbuildups[0].Treaty2Ibnr + licbuildups[1].Treaty2Ibnr + licbuildups[2].Treaty2Ibnr + licbuildups[3].Treaty2Ibnr + licbuildups[4].Treaty2Ibnr + licbuildups[5].Treaty2Ibnr)
								licbuildup.Treaty2IbnrRiskAdjustment = ((1-licParams.CurrTreaty2NonPerformanceRate)/(1-licParams.PrevTreaty2NonPerformanceRate) - 1) * (licbuildups[0].Treaty2IbnrRiskAdjustment + licbuildups[1].Treaty2IbnrRiskAdjustment + licbuildups[2].Treaty2IbnrRiskAdjustment + licbuildups[3].Treaty2IbnrRiskAdjustment + licbuildups[4].Treaty2IbnrRiskAdjustment + licbuildups[5].Treaty2IbnrRiskAdjustment)
							}
							if licParams.PrevTreaty3NonPerformanceRate != 0 {
								licbuildup.Treaty3Ibnr = ((1-licParams.CurrTreaty3NonPerformanceRate)/(1-licParams.PrevTreaty3NonPerformanceRate) - 1) * (licbuildups[0].Treaty3Ibnr + licbuildups[1].Treaty3Ibnr + licbuildups[2].Treaty3Ibnr + licbuildups[3].Treaty3Ibnr + licbuildups[4].Treaty3Ibnr + licbuildups[5].Treaty3Ibnr)
								licbuildup.Treaty3IbnrRiskAdjustment = ((1-licParams.CurrTreaty3NonPerformanceRate)/(1-licParams.PrevTreaty3NonPerformanceRate) - 1) * (licbuildups[0].Treaty3IbnrRiskAdjustment + licbuildups[1].Treaty3IbnrRiskAdjustment + licbuildups[2].Treaty3IbnrRiskAdjustment + licbuildups[3].Treaty3IbnrRiskAdjustment + licbuildups[4].Treaty3IbnrRiskAdjustment + licbuildups[5].Treaty3IbnrRiskAdjustment)
							}

						}
						if licbuildup.Name == "Old Claims Paid" {

							if licbuildups[0].ReportedClaims > 0 {
								var ibnribnrmp models.LicModelPoint
								err = DB.Where("run_id=? and product_code=? ", licVar.RunId, product).Find(&ibnribnrmp).Error
								if err != nil {
									log.Error(err)
								}
								var OldClaimSum float64
								query := fmt.Sprintf("SELECT sum(paid_claims) as paid_claims FROM lic_model_points where run_id = %d and reporting_year + reporting_month/1000.0 <= %.3f and settlement_year + settlement_month/1000 > %.3f and product_code = '%s'  group by product_code", licVar.RunId, oldClaimsCutOff, oldClaimsCutOff, product)
								err = DB.Raw(query).Scan(&OldClaimSum).Error
								licbuildup.ReportedClaims = -OldClaimSum
							} else {
								licbuildup.ReportedClaims = 0
							}

							licbuildup.IBNR = 0
							licbuildup.RiskAdjustment = 0
							licbuildup.CashBack = 0
							licbuildup.UnadjustedLossAdjustmentExpenses = 0 //(licbuildup.IBNR + licbuildup.ReportedClaims*0.5) * 0.05
							licbuildup.Treaty1Ibnr = 0
							licbuildup.Treaty1IbnrRiskAdjustment = 0
							licbuildup.Treaty2Ibnr = 0
							licbuildup.Treaty2IbnrRiskAdjustment = 0
							licbuildup.Treaty3Ibnr = 0
							licbuildup.Treaty3IbnrRiskAdjustment = 0
						}

						if licbuildup.Name == "New Data" {
							licbuildup.IBNR = ibnrResult.IbnrBel - previbnrAt12
							licbuildup.ReportedClaims = licParams.ReportedClaimsEstimateAdjustment
							licbuildup.RiskAdjustment = ibnrResult.IBNRRiskAdjustment - previbnrRiskAdjustmentAt12
							licbuildup.CashBack = licParams.CashbackEstimateAdjustment
							licbuildup.UnadjustedLossAdjustmentExpenses = 0 //(licbuildup.IBNR + licbuildup.ReportedClaims*0.5) * 0.05

							licbuildup.Treaty1Ibnr = licbuildup.IBNR * licParams.Treaty1ClaimsProportion * (1 - licParams.CurrTreaty1NonPerformanceRate)
							licbuildup.Treaty2Ibnr = licbuildup.IBNR * licParams.Treaty2ClaimsProportion * (1 - licParams.CurrTreaty2NonPerformanceRate)
							licbuildup.Treaty3Ibnr = licbuildup.IBNR * licParams.Treaty3ClaimsProportion * (1 - licParams.CurrTreaty3NonPerformanceRate)

							licbuildup.Treaty1IbnrRiskAdjustment = licbuildup.RiskAdjustment * licParams.Treaty1ClaimsProportion * (1 - licParams.CurrTreaty1NonPerformanceRate)
							licbuildup.Treaty2IbnrRiskAdjustment = licbuildup.RiskAdjustment * licParams.Treaty2ClaimsProportion * (1 - licParams.CurrTreaty2NonPerformanceRate)
							licbuildup.Treaty3IbnrRiskAdjustment = licbuildup.RiskAdjustment * licParams.Treaty3ClaimsProportion * (1 - licParams.CurrTreaty3NonPerformanceRate)

						}

						if licbuildup.Name == "New Claims Reported" {

							var NewClaimReportedSum float64
							query3 := fmt.Sprintf("SELECT sum(total_inflated_claim) as total_inflated_claim FROM lic_model_points where run_id = %d and reporting_year + reporting_month/1000.0 > %.3f and product_code = '%s'  group by product_code", licVar.RunId, oldClaimsCutOff, product)
							err = DB.Raw(query3).Scan(&NewClaimReportedSum).Error
							licbuildup.IBNR = 0
							licbuildup.ReportedClaims = NewClaimReportedSum
							licbuildup.RiskAdjustment = 0
							licbuildup.CashBack = 0
							licbuildup.UnadjustedLossAdjustmentExpenses = 0 //(licbuildup.IBNR + licbuildup.ReportedClaims*0.5) * 0.05
							licbuildup.Treaty1Ibnr = 0
							licbuildup.Treaty1IbnrRiskAdjustment = 0
							licbuildup.Treaty2Ibnr = 0
							licbuildup.Treaty2IbnrRiskAdjustment = 0
							licbuildup.Treaty3Ibnr = 0
							licbuildup.Treaty3IbnrRiskAdjustment = 0

						}

						if licbuildup.Name == "New Claims Paid" {
							var NewClaimSum float64
							query2 := fmt.Sprintf("SELECT sum(paid_claims) as paid_claims FROM lic_model_points where run_id = %d and reporting_year + reporting_month/1000.0 > %.3f and settlement_year + settlement_month/1000 > %.3f and product_code = '%s'  group by product_code", licVar.RunId, oldClaimsCutOff, oldClaimsCutOff, product)
							err = DB.Raw(query2).Scan(&NewClaimSum).Error
							licbuildup.IBNR = 0
							licbuildup.ReportedClaims = -NewClaimSum
							licbuildup.RiskAdjustment = 0
							licbuildup.CashBack = 0
							licbuildup.UnadjustedLossAdjustmentExpenses = 0 //(licbuildup.IBNR + licbuildup.ReportedClaims*0.5) * 0.05
							licbuildup.Treaty1Ibnr = 0
							licbuildup.Treaty1IbnrRiskAdjustment = 0
							licbuildup.Treaty2Ibnr = 0
							licbuildup.Treaty2IbnrRiskAdjustment = 0
							licbuildup.Treaty3Ibnr = 0
							licbuildup.Treaty3IbnrRiskAdjustment = 0

						}

						if licbuildup.Name == "C/F" {
							licbuildup.IBNR = ibnrResult.IbnrBel - previbnr
							licbuildup.ReportedClaims = 0
							licbuildup.RiskAdjustment = ibnrResult.IBNRRiskAdjustment - previbnrRiskAdjustment
							licbuildup.CashBack = 0
							licbuildup.UnadjustedLossAdjustmentExpenses = 0 //(licbuildup.IBNR + licbuildup.ReportedClaims*0.5) * 0.05

							licbuildup.Treaty1Ibnr = licbuildup.IBNR * licParams.Treaty1ClaimsProportion * (1 - licParams.CurrTreaty1NonPerformanceRate)
							licbuildup.Treaty2Ibnr = licbuildup.IBNR * licParams.Treaty2ClaimsProportion * (1 - licParams.CurrTreaty2NonPerformanceRate)
							licbuildup.Treaty3Ibnr = licbuildup.IBNR * licParams.Treaty3ClaimsProportion * (1 - licParams.CurrTreaty3NonPerformanceRate)

							licbuildup.Treaty1IbnrRiskAdjustment = licbuildup.RiskAdjustment * licParams.Treaty1ClaimsProportion * (1 - licParams.CurrTreaty1NonPerformanceRate)
							licbuildup.Treaty2IbnrRiskAdjustment = licbuildup.RiskAdjustment * licParams.Treaty2ClaimsProportion * (1 - licParams.CurrTreaty2NonPerformanceRate)
							licbuildup.Treaty3IbnrRiskAdjustment = licbuildup.RiskAdjustment * licParams.Treaty3ClaimsProportion * (1 - licParams.CurrTreaty3NonPerformanceRate)

						}

						licbuildup.VariableChange = licbuildup.IBNR + licbuildup.RiskAdjustment + licbuildup.ReportedClaims + licbuildup.CashBack + licbuildup.CashBack + licbuildup.UnadjustedLossAdjustmentExpenses
						if i > 1 && licbuildup.Name != "New Data" {
							licbuildup.LicBuildup = licbuildups[i-1].LicBuildup + licbuildup.VariableChange
						}

						licbuildup.Pnl = -licbuildup.VariableChange

						if licbuildup.Name == "Old Claims Paid" {
							licbuildup.Pnl = 0
						}

						if licbuildup.Name == "New Claims Paid" {
							licbuildup.Pnl = 0
						}
					}

					previbnrAt12 = ibnrResult.IbnrBelAt12
					previbnrRiskAdjustmentAt12 = ibnrResult.IbnrRiskAdjustmentAt12
					previbnr = ibnrResult.IbnrBel
					previbnrRiskAdjustment = ibnrResult.IBNRRiskAdjustment
				}
				licbuildups = append(licbuildups, licbuildup)
			}
			var ljt models.LicJournalTransactions
			ljt.IFRS17Group = product
			ljt.ProductCode = product //licbuildups[0].ProductCode
			ljt.LicRunID = licbuildups[11].RunId
			ljt.RunDate = licbuildups[11].RunDate

			ljt.LicFutureCashFlowsChange = -(licbuildups[3].Pnl + licbuildups[4].Pnl + licbuildups[5].Pnl + licbuildups[6].Pnl)
			ljt.LicExperienceVariance = -(licbuildups[2].Pnl + licbuildups[7].Pnl)
			ljt.IBNRIncurredClaims = -licbuildups[9].Pnl

			err = DB.Where("product_code = ? and ifrs17_group = ?  and run_date=?", ljt.ProductCode, ljt.IFRS17Group, ljt.RunDate).Delete(&ljt).Error
			if err != nil {
				log.Info(err.Error())
			}

			err = DB.Where("product_code = ? and ifrs17_group = ? and run_date=?", ljt.ProductCode, ljt.IFRS17Group, ljt.RunDate).Save(&ljt).Error
			if err != nil {
				log.Info(err.Error())
			}

			//Savings Balance Sheet Records
			var BalanceSheetRecord models.BalanceSheetRecord
			BalanceSheetRecord.CsmRunID = licbuildups[11].RunId
			BalanceSheetRecord.ProductCode = product
			BalanceSheetRecord.MeasurementType = "LIC"
			BalanceSheetRecord.IFRS17Group = product
			BalanceSheetRecord.Date = licbuildups[11].RunDate
			BalanceSheetRecord.OutstandingClaimsReserve = licbuildups[0].ReportedClaims + licbuildups[2].ReportedClaims + licbuildups[7].ReportedClaims + licbuildups[8].ReportedClaims + licbuildups[9].ReportedClaims + licbuildups[10].ReportedClaims // B/F- OldClaimsPaid + NewClaims Reported - NewClaimsPaid
			BalanceSheetRecord.IBNR = ibnrBS
			BalanceSheetRecord.IBNRRiskAdjustment = ibnrRiskAdjustmentBS
			BalanceSheetRecord.Treaty1IBNR = ibnrBS * licParams.Treaty1ClaimsProportion * (1 - licParams.CurrTreaty1NonPerformanceRate)
			BalanceSheetRecord.Treaty1IBNRRiskAdjustment = ibnrRiskAdjustmentBS * licParams.Treaty1ClaimsProportion * (1 - licParams.CurrTreaty1NonPerformanceRate)
			BalanceSheetRecord.Treaty2IBNR = ibnrBS * licParams.Treaty2ClaimsProportion * (1 - licParams.CurrTreaty2NonPerformanceRate)
			BalanceSheetRecord.Treaty2IBNRRiskAdjustment = ibnrRiskAdjustmentBS * licParams.Treaty2ClaimsProportion * (1 - licParams.CurrTreaty2NonPerformanceRate)
			BalanceSheetRecord.Treaty3IBNR = ibnrBS * licParams.Treaty3ClaimsProportion * (1 - licParams.CurrTreaty3NonPerformanceRate)
			BalanceSheetRecord.Treaty3IBNRRiskAdjustment = ibnrRiskAdjustmentBS * licParams.Treaty3ClaimsProportion * (1 - licParams.CurrTreaty3NonPerformanceRate)

			err = DB.Where("product_code = ? and ifrs17_group = ?  and date=? and measurement_type=?", BalanceSheetRecord.ProductCode, BalanceSheetRecord.IFRS17Group, BalanceSheetRecord.Date, "LIC").Delete(&BalanceSheetRecord).Error
			if err != nil {
				log.Info(err.Error())
			}

			err = DB.Where("product_code = ? and ifrs17_group = ? and date=? and measurement_type=?", BalanceSheetRecord.ProductCode, BalanceSheetRecord.IFRS17Group, BalanceSheetRecord.Date, "LIC").Save(&BalanceSheetRecord).Error
			if err != nil {
				log.Info(err.Error())
			}

		}

		// Disabled for now. We need to clarify this
		DB.Where("run_date=? and configuration_name=?", run.RunDate, run.LicConfigurationName).Delete(&models.LicBuildupResult{})
		DB.Save(&licbuildups)

		end := time.Since(start)
		log.Info("LicEngine took: ", end)
		run.RunTime = end.Minutes()
		run.ProcessingStatus = "completed"
		DB.Save(&run)
	}
}
