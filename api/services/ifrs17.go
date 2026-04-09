package services

import (
	"api/installer"
	"api/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
)

func GetAOSVariables(ctx context.Context) ([]models.BaseAosVariable, error) {
	var variables []models.BaseAosVariable
	err := DB.Find(&variables).Error
	if err != nil {
		return variables, err
	}
	return variables, nil
}

func GetAosVariableSets() ([]models.AosVariableSet, error) {
	var results []models.AosVariableSet
	err := DB.Preload("AosVariables").Find(&results).Error
	return results, err
}

func CreateFinancialReport(productCode string, runDate string) []models.IncomeStatementEntry {
	file, err := installer.Files.Open("finance_reports.json")
	//file, err := pkger.Open("/installer/finance_reports.json")
	if err != nil {
		fmt.Println(err)
	}

	body, err := io.ReadAll(file)

	var reportItems []models.IncomeStatementEntry

	err = json.Unmarshal(body, &reportItems)
	fmt.Println(err)

	//Get Journal for Product (FUN_01 in this case. Generalize later.)
	var journal models.JournalTransactions
	var licjournal models.LicJournalTransactions

	if productCode == "ALL" {
		var journals []models.JournalTransactions
		err := DB.Where("run_date=?", runDate).Find(&journals).Error
		if err != nil {
			fmt.Println(err)
		}

		var licjournals []models.LicJournalTransactions
		err = DB.Where("run_date=?", runDate).Find(&licjournals).Error
		if err != nil {
			fmt.Println(err)
		}

		for _, j := range journals {
			journal.CsmRelease += j.CsmRelease
			journal.DacRelease += j.DacRelease + j.AmortizationAcquisitionCF
			journal.Expenses += j.Expenses
			journal.AmortizationAcquisitionCF += j.AmortizationAcquisitionCF + j.DacRelease
			journal.ClaimsIncurred += j.ClaimsIncurred
			journal.ExpectedBenefits += j.ExpectedBenefits
			journal.PremiumVariance += j.PremiumVariance
			journal.ExpensesIncurred += j.ExpensesIncurred
			journal.NonAttributableExpensesIncurred += j.NonAttributableExpensesIncurred
			journal.InsuranceFinanceExpense += j.InsuranceFinanceExpense
			journal.LossComponentFutureServiceChange += j.LossComponentFutureServiceChange
			journal.LossComponentOnInitialRecog += j.LossComponentOnInitialRecog
			journal.LossComponentUnwind += j.LossComponentUnwind
			journal.RAChange += j.RAChange
			journal.PaaEarnedPremium += j.PaaEarnedPremium
			journal.PaaLossComponent += j.PaaLossComponent
			journal.PaaLossComponentAdjustment += j.PaaLossComponentAdjustment
			journal.ReinsuranceFlatCommission += j.ReinsuranceFlatCommission
			journal.ReinsuranceProvisionalCommission += j.ReinsuranceProvisionalCommission
			journal.ReinsuranceUltimateCommission += j.ReinsuranceUltimateCommission
			journal.ReinsuranceProfitCommission += j.ReinsuranceProfitCommission
			journal.ReinsuranceRecovery += j.ReinsuranceRecovery
			journal.ReinsuranceInvestmentComponent += j.ReinsuranceInvestmentComponent
			journal.PaaReinsurancePremium += j.PaaReinsurancePremium
			journal.PaaLossRecoveryComponent += j.PaaLossRecoveryComponent
			journal.PaaLossRecoveryUnwind += j.PaaLossRecoveryUnwind
			journal.PaaLossRecoveryAdjustment += j.PaaLossRecoveryAdjustment
		}

		// Lic Cashflows
		for _, licj := range licjournals {
			licjournal.LicFutureCashFlowsChange += licj.LicFutureCashFlowsChange
			licjournal.LicExperienceVariance += licj.LicExperienceVariance
			licjournal.IBNRIncurredClaims += licj.IBNRIncurredClaims
		}

	} else {
		var journals []models.JournalTransactions
		DB.Where("product_code=? and run_date=?", productCode, runDate).Find(&journals)
		journal.PaaEarnedPremium = 0
		journal.PaaLossComponent = 0

		var licjournals []models.LicJournalTransactions
		err = DB.Where("product_code=? and run_date=?", productCode, runDate).Find(&licjournals).Error
		if err != nil {
			fmt.Println(err)
		}

		for _, j := range journals {
			journal.CsmRelease += j.CsmRelease
			journal.DacRelease += j.DacRelease + j.AmortizationAcquisitionCF
			journal.Expenses += j.Expenses
			journal.AmortizationAcquisitionCF += j.AmortizationAcquisitionCF + j.DacRelease
			journal.ClaimsIncurred += j.ClaimsIncurred
			journal.ExpectedBenefits += j.ExpectedBenefits
			journal.PremiumVariance += j.PremiumVariance
			journal.ExpensesIncurred += j.ExpensesIncurred
			journal.NonAttributableExpensesIncurred += j.NonAttributableExpensesIncurred
			journal.InsuranceFinanceExpense += j.InsuranceFinanceExpense
			journal.LossComponentFutureServiceChange += j.LossComponentFutureServiceChange
			journal.LossComponentOnInitialRecog += j.LossComponentOnInitialRecog
			journal.LossComponentUnwind += j.LossComponentUnwind
			journal.RAChange += j.RAChange
			journal.PaaEarnedPremium += j.PaaEarnedPremium
			journal.PaaLossComponent += j.PaaLossComponent
			journal.PaaLossComponentAdjustment += j.PaaLossComponentAdjustment
			journal.ReinsuranceFlatCommission += j.ReinsuranceFlatCommission
			journal.ReinsuranceProvisionalCommission += j.ReinsuranceProvisionalCommission
			journal.ReinsuranceUltimateCommission += j.ReinsuranceUltimateCommission
			journal.ReinsuranceProfitCommission += j.ReinsuranceProfitCommission
			journal.ReinsuranceRecovery += j.ReinsuranceRecovery
			journal.ReinsuranceInvestmentComponent += j.ReinsuranceInvestmentComponent
			journal.PaaReinsurancePremium += j.PaaReinsurancePremium
			journal.PaaLossRecoveryComponent += j.PaaLossRecoveryComponent
			journal.PaaLossRecoveryUnwind += j.PaaLossRecoveryUnwind
			journal.PaaLossRecoveryAdjustment += j.PaaLossRecoveryAdjustment
		}

		// Lic Cashflows
		for _, licj := range licjournals {
			licjournal.LicFutureCashFlowsChange += licj.LicFutureCashFlowsChange
			licjournal.LicExperienceVariance += licj.LicExperienceVariance
			licjournal.IBNRIncurredClaims += licj.IBNRIncurredClaims
		}
	}

	//Insert calculated values in this block
	reportItems[0].CurrentYear = math.Round(journal.CsmRelease)
	reportItems[1].CurrentYear = math.Round(journal.DacRelease)
	reportItems[2].CurrentYear = math.Round(journal.RAChange)
	reportItems[3].CurrentYear = math.Round(journal.ExpectedBenefits)
	reportItems[4].CurrentYear = math.Round(journal.Expenses)
	reportItems[6].CurrentYear = math.Round(journal.PremiumVariance)
	//if journal.PaaLossComponentAdjustment < 0 {
	//	reportItems[7].CurrentYear = -math.Round(journal.PaaLossComponentAdjustment)
	//} else {
	reportItems[14].CurrentYear = math.Round(journal.PaaLossComponentAdjustment)
	//}

	reportItems[7].CurrentYear = math.Round(journal.LossComponentUnwind)
	reportItems[8].CurrentYear = -math.Round(journal.ClaimsIncurred) //+ math.Round(licjournal.IBNRIncurredClaims))
	reportItems[9].CurrentYear = -math.Round(journal.ExpensesIncurred)
	reportItems[10].CurrentYear = -math.Round(journal.AmortizationAcquisitionCF)
	reportItems[11].CurrentYear = -math.Round(licjournal.LicFutureCashFlowsChange)
	reportItems[12].CurrentYear = -math.Round(licjournal.LicExperienceVariance)
	reportItems[13].CurrentYear = -math.Round(journal.LossComponentFutureServiceChange)
	reportItems[15].CurrentYear = -(math.Round(journal.LossComponentOnInitialRecog) + math.Round(journal.PaaLossComponent)) //reportItems[30].CurrentYear = math.Round(journal.InsuranceFinanceExpense)
	reportItems[5].CurrentYear = math.Round(journal.PaaEarnedPremium)
	reportItems[16].CurrentYear = -math.Round(journal.LossComponentUnwind)
	//reportItems[16].CurrentYear = math.Round(journal.PaaLossComponent)

	reportItems[20].CurrentYear = -math.Round(journal.PaaReinsurancePremium - journal.ReinsuranceInvestmentComponent)
	reportItems[21].CurrentYear = math.Round(-journal.ReinsuranceFlatCommission - journal.ReinsuranceProvisionalCommission)
	reportItems[22].CurrentYear = 0
	reportItems[23].CurrentYear = 0
	reportItems[24].CurrentYear = math.Round(journal.ReinsuranceRecovery)

	reportItems[25].CurrentYear = math.Round(journal.ReinsuranceProfitCommission + journal.ReinsuranceUltimateCommission - journal.ReinsuranceProvisionalCommission - journal.ReinsuranceReinstatementPremium - journal.ReinsuranceInvestmentComponent)

	reportItems[26].CurrentYear = 0                                             // Changes in FCF that do not adjust underlying contracts CSM)
	reportItems[27].CurrentYear = 0                                             // Incurred Claims Adjustment Changes that Relate to Past Service
	reportItems[28].CurrentYear = math.Round(journal.PaaLossRecoveryComponent)  // Loss Recovery at initial recognition
	reportItems[29].CurrentYear = math.Round(journal.PaaLossRecoveryAdjustment) //Loss Recovery at subsequent measurements
	reportItems[30].CurrentYear = -math.Round(journal.PaaLossRecoveryUnwind)    // Reversal of Loss Recovery Component

	reportItems[31].CurrentYear = 0 //math.Round(journal.ReinsuranceInvestmentComponent)
	reportItems[32].CurrentYear = math.Round(journal.InsuranceFinanceExpense)
	reportItems[33].CurrentYear = 0
	reportItems[34].CurrentYear = -math.Round(journal.NonAttributableExpensesIncurred)

	// End of inserts
	var subTotalPastYear float64
	var subTotalCurrentYear float64

	var icrPastYear float64
	var icrCurrentYear float64

	var isePastYear float64
	var iseCurrentYear float64

	var isrPastYear float64
	var isrCurrentYear float64

	//var ifPastYear float64
	//var ifCurrentYear float64

	var tcePastYear float64
	var tceCurrentYear float64

	var reinsurancePremiumPAAPastYear float64
	var reinsurancePremiumPAACurrentYear float64

	var reinsurancePremiumGMMPastYear float64
	var reinsurancePremiumGMMCurrentYear float64

	var reinsuranceRecoveriesPastYear float64
	var reinsuranceRecoveriesCurrentYear float64

	var netReinsuranceLinePastYear float64
	var netReinsuranceLineCurrentYear float64

	var investmentIncomeFinanceExpenseLinePastYear float64
	var investmentIncomeFinanceExpenseLineCurrentYear float64

	var fcfadjustmentsandlossrecoveryComponentPastYear float64
	var fcfadjustmentsandlossrecoveryComponentCurrentYear float64

	var profitLossPastYear float64
	var profitLossCurrentYear float64

	var otherPastYear float64
	var otherCurrentYear float64

	for _, val := range reportItems {
		if val.Index < 3 {
			subTotalPastYear += val.PastYear
			subTotalCurrentYear += val.CurrentYear
		}
		if val.Index < 9 {
			icrPastYear += val.PastYear
			icrCurrentYear += val.CurrentYear
		}

		if val.Index > 8 && val.Index < 18 {
			isePastYear += val.PastYear
			iseCurrentYear += val.CurrentYear
		}

		//if val.Index > 16 && val.Index < 21 {
		//	ifPastYear += val.PastYear
		//	ifCurrentYear += val.CurrentYear
		//}

		if val.Index > 17 && val.Index < 21 {
			reinsurancePremiumGMMPastYear += val.PastYear
			reinsurancePremiumGMMCurrentYear += val.CurrentYear
		}

		if val.Index > 20 && val.Index < 23 {
			reinsurancePremiumPAAPastYear += val.PastYear
			reinsurancePremiumPAACurrentYear += val.CurrentYear
		}

		if val.Index > 22 && val.Index < 27 {
			reinsuranceRecoveriesPastYear += val.PastYear
			reinsuranceRecoveriesCurrentYear += val.CurrentYear
		}

		if val.Index > 26 && val.Index < 32 {
			fcfadjustmentsandlossrecoveryComponentPastYear += val.PastYear
			fcfadjustmentsandlossrecoveryComponentCurrentYear += val.CurrentYear
		}

		if val.Index > 31 && val.Index < 34 {
			investmentIncomeFinanceExpenseLinePastYear += val.PastYear
			investmentIncomeFinanceExpenseLineCurrentYear += val.CurrentYear
		}

		if val.Index > 33 && val.Index < 36 {
			otherPastYear += val.PastYear
			otherCurrentYear += val.CurrentYear
		}

	}

	isrCurrentYear = icrCurrentYear + iseCurrentYear
	isrPastYear = icrPastYear + isePastYear

	netReinsuranceLinePastYear = reinsuranceRecoveriesPastYear + reinsurancePremiumGMMPastYear + reinsurancePremiumPAAPastYear + fcfadjustmentsandlossrecoveryComponentPastYear
	netReinsuranceLineCurrentYear = reinsuranceRecoveriesCurrentYear + reinsurancePremiumGMMCurrentYear + reinsurancePremiumPAACurrentYear + fcfadjustmentsandlossrecoveryComponentCurrentYear

	profitLossPastYear = isrPastYear + netReinsuranceLinePastYear + investmentIncomeFinanceExpenseLinePastYear
	profitLossCurrentYear = isrCurrentYear + netReinsuranceLineCurrentYear + investmentIncomeFinanceExpenseLineCurrentYear

	tcePastYear = profitLossPastYear + otherPastYear
	tceCurrentYear = profitLossCurrentYear + otherCurrentYear

	//fmt.Println(subTotalCurrentYear, subTotalPastYear)
	//fmt.Println(icrCurrentYear, icrPastYear)
	//fmt.Println(iseCurrentYear, isePastYear)
	//fmt.Println(isrCurrentYear, isrPastYear)
	//fmt.Println(ifCurrentYear, ifPastYear)
	//fmt.Println(tceCurrentYear, tcePastYear)

	//Build Reports

	var allItems []models.IncomeStatementEntry

	for _, v := range reportItems {
		if v.Index < 3 {
			allItems = append(allItems, v)
		}

		if v.Index == 3 {
			// create a sub total entry and inject it into the array
			var item models.IncomeStatementEntry
			item.PastYear = subTotalPastYear
			item.CurrentYear = subTotalCurrentYear
			item.LineItem = "Sub Total"
			item.Type = "sub-total"
			item.Style = "bold"
			allItems = append(allItems, item, v)
		}

		if v.Index > 3 && v.Index < 8 {
			allItems = append(allItems, v)
		}

		if v.Index == 8 {
			var item models.IncomeStatementEntry
			item.LineItem = "Revenue prior to adjustment for Loss Component "
			var icrLine models.IncomeStatementEntry
			icrLine.LineItem = "Insurance Contract Revenue (A)"
			icrLine.Notes = "" //"Expected premium vs actual premium future service===CSM"
			icrLine.PastYear = icrPastYear
			icrLine.CurrentYear = icrCurrentYear
			icrLine.Type = "sub-total"
			icrLine.Style = "bold"

			allItems = append(allItems, v, icrLine, models.IncomeStatementEntry{Type: "empty"})
		}

		if v.Index > 8 && v.Index < 18 { //17 changed to 31 still need to add subtotals
			allItems = append(allItems, v)
		}

		if v.Index == 18 {
			var iseLine models.IncomeStatementEntry
			var isrLine models.IncomeStatementEntry
			var reinsuranceLine models.IncomeStatementEntry
			iseLine.LineItem = "Insurance Service Expense (B)"
			iseLine.PastYear = isePastYear
			iseLine.CurrentYear = iseCurrentYear
			iseLine.Type = "sub-total"
			iseLine.Style = "bold"

			isrLine.LineItem = "Insurance Service Results (A+B)"
			isrLine.PastYear = isrPastYear
			isrLine.CurrentYear = isrCurrentYear
			isrLine.Type = "sub-total"
			isrLine.Style = "bold"

			reinsuranceLine.LineItem = "Reinsurance contracts held"
			reinsuranceLine.Type = "sub-total-reinsurance"
			reinsuranceLine.Style = "bold"

			allItems = append(allItems, iseLine, models.IncomeStatementEntry{Type: "empty"}, isrLine, models.IncomeStatementEntry{Type: "empty"}, reinsuranceLine, v)
		}
		if v.Index > 18 && v.Index < 20 {
			allItems = append(allItems, v)
		}

		if v.Index == 20 {
			var reinsurancePremiumGMM models.IncomeStatementEntry
			reinsurancePremiumGMM.LineItem = "Allocation of Reinsurance Premium - GMM (C)"
			reinsurancePremiumGMM.PastYear = reinsurancePremiumGMMPastYear
			reinsurancePremiumGMM.CurrentYear = reinsurancePremiumGMMCurrentYear
			reinsurancePremiumGMM.Type = "sub-total"
			reinsurancePremiumGMM.Style = "bold"
			allItems = append(allItems, v, reinsurancePremiumGMM)
		}

		if v.Index > 20 && v.Index < 22 {
			allItems = append(allItems, v)
		}

		if v.Index == 22 {
			var reinsurancePremiumPAA models.IncomeStatementEntry
			reinsurancePremiumPAA.LineItem = "Allocation of Reinsurance Premium - PAA (D)"
			reinsurancePremiumPAA.PastYear = reinsurancePremiumPAAPastYear
			reinsurancePremiumPAA.CurrentYear = reinsurancePremiumPAACurrentYear
			reinsurancePremiumPAA.Type = "sub-total"
			reinsurancePremiumPAA.Style = "bold"
			allItems = append(allItems, v, reinsurancePremiumPAA)
		}

		if v.Index > 22 && v.Index < 26 {
			allItems = append(allItems, v)
		}

		if v.Index == 26 {
			var reinsuranceRecoveries models.IncomeStatementEntry
			reinsuranceRecoveries.LineItem = "Amount Recovered from Reinsurer (E)"
			reinsuranceRecoveries.PastYear = reinsuranceRecoveriesPastYear
			reinsuranceRecoveries.CurrentYear = reinsuranceRecoveriesCurrentYear
			reinsuranceRecoveries.Type = "sub-total"
			reinsuranceRecoveries.Style = "bold"
			allItems = append(allItems, v, reinsuranceRecoveries)
		}

		if v.Index > 26 && v.Index < 31 {
			allItems = append(allItems, v)
		}

		if v.Index == 31 {
			var netReinsuranceLine models.IncomeStatementEntry
			netReinsuranceLine.LineItem = "Total net income (expenses) from reinsurance contracts held"
			netReinsuranceLine.PastYear = netReinsuranceLinePastYear
			netReinsuranceLine.CurrentYear = netReinsuranceLineCurrentYear
			netReinsuranceLine.Type = "sub-total"
			netReinsuranceLine.Style = "bold"
			allItems = append(allItems, v, netReinsuranceLine)
		}

		if v.Index > 31 && v.Index < 33 {
			allItems = append(allItems, v)
		}

		if v.Index == 33 {
			var plLine models.IncomeStatementEntry
			plLine.LineItem = "Net Insurance and Investment Result" //"Profit & Loss"
			plLine.PastYear = profitLossPastYear
			plLine.CurrentYear = profitLossCurrentYear
			plLine.Type = "sub-total"
			plLine.Style = "bold"
			allItems = append(allItems, v, plLine)
		}

		if v.Index > 33 && v.Index < 36 {
			allItems = append(allItems, v)
		}
	}

	var tceLine models.IncomeStatementEntry
	tceLine.LineItem = "Total Comprehensive Income"
	tceLine.PastYear = tcePastYear
	tceLine.CurrentYear = tceCurrentYear
	tceLine.Style = "bold"
	tceLine.Type = "sub-total"

	allItems = append(allItems, tceLine)

	return allItems
}
