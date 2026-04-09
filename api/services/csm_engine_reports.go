package services

import (
	"api/models"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

func GenerateJournalsforAll(runDate string) map[string]interface{} {
	var res = make(map[string]interface{})
	var journals []models.JTReportEntry
	results := GetAosStepResultsForAllProductsForDownstreamCalcs(runDate)
	if len(results) == 0 {
		mockResults := getMockResults()
		results = append(results, mockResults...)
	}
	paaresults := GetPAAResultsForAllProductsForDownstreamCalcs(runDate)
	if len(paaresults) == 0 {
		mockResult := models.PAAResult{}
		paaresults = append(paaresults, mockResult)
	}

	licBuildupResults := GetLicResultsForAllProductsForDownstreamCalcs(runDate)
	if len(licBuildupResults) == 0 {
		mockResult := getLICMockResults()
		licBuildupResults = append(licBuildupResults, mockResult...)
	}

	journals = GenerateJournalTransactionReportByGroup(results, paaresults, licBuildupResults)
	prodList := GetJournalEntryProducts(runDate)

	res["steps"] = journals
	res["products"] = prodList

	return res
}

func GenerateJournalsforOneProduct(runDate, prodCode string) map[string]interface{} {

	var res = make(map[string]interface{})
	var journals []models.JTReportEntry
	results := GetAosStepResultsForOneProductForDownstreamCalcs(runDate, prodCode)
	if len(results) == 0 {
		mockResults := getMockResults()
		results = append(results, mockResults...)
	}
	paaresults := GetPAAResultsForOneProductForDownstreamCalcs(runDate, prodCode)
	if len(paaresults) == 0 {
		mockResult := models.PAAResult{}
		paaresults = append(paaresults, mockResult)
	}
	licBuildupResults := GetLicResultsForAllProductsForDownstreamCalcs(runDate)

	if len(licBuildupResults) == 0 {
		mockResult := getLICMockResults()
		licBuildupResults = append(licBuildupResults, mockResult...)
	}

	journals = GenerateJournalTransactionReportByGroup(results, paaresults, licBuildupResults)
	groupList := GetGroups(runDate, prodCode)

	res["steps"] = journals
	res["groups"] = groupList

	return res
}

func getMockResults() []models.AOSStepResult {
	var aosVariables []models.BaseAosVariable
	var mockResults []models.AOSStepResult

	DB.Find(&aosVariables)

	for _, variable := range aosVariables {
		var step models.AOSStepResult
		step.Name = variable.Name
		mockResults = append(mockResults, step)
	}
	return mockResults
}

func getLICMockResults() []models.LicBuildupResult {
	var licVariables []models.LicBaseVariable
	var mockResults []models.LicBuildupResult

	DB.Find(&licVariables)

	for _, variable := range licVariables {
		var step models.LicBuildupResult
		step.Name = variable.Name
		mockResults = append(mockResults, step)
	}
	return mockResults
}

func GenerateTransactionReports(ifrs17Group, runDate string) []models.JTReportEntry {
	// Get step results by product code and group
	// TODO: Check if group is available on PAA table or GMM Table

	results := GetAosStepResultsForGroup(ifrs17Group, runDate)
	if len(results) == 0 {
		mockResults := getMockResults()
		results = append(results, mockResults...)
	}
	paaresults := GetPAAResultsForGroup(ifrs17Group, runDate)
	if len(paaresults) == 0 {
		mockResult := models.PAAResult{}
		paaresults = append(paaresults, mockResult)
	}

	licBuildupResults := GetLicResultsForOneGroupForDownstreamCalcs(runDate, ifrs17Group)
	if len(licBuildupResults) == 0 {
		mockResult := getLICMockResults()
		licBuildupResults = append(licBuildupResults, mockResult...)
	}

	report := GenerateJournalTransactionReportByGroup(results, paaresults, licBuildupResults)
	return report
}

func GenerateSubLedgerReports(ifrs17Group, runDate string) []models.SubLedgerReportEntry {
	results := GetAosStepResultsForGroup(ifrs17Group, runDate)
	if len(results) == 0 {
		mockResults := getMockResults()
		results = append(results, mockResults...)
	}
	paaresults := GetPAAResultsForGroup(ifrs17Group, runDate)
	if len(paaresults) == 0 {
		mockResult := models.PAAResult{}
		paaresults = append(paaresults, mockResult)
	}
	licBuildupResults := GetLicResultsForOneGroupForDownstreamCalcs(runDate, ifrs17Group)
	if len(licBuildupResults) == 0 {
		mockResult := getLICMockResults()
		licBuildupResults = append(licBuildupResults, mockResult...)
	}

	res := GenerateJournalTransactionReportByGroup(results, paaresults, licBuildupResults)
	report := GenerateSubLedgerReportByGroup(res, results, paaresults)

	return report
}

func GenerateSubLedgerReportsforAllProducts(runDate string) map[string]interface{} {

	var res = make(map[string]interface{})
	var report []models.SubLedgerReportEntry

	results := GetAosStepResultsForAllProductsForDownstreamCalcs(runDate)
	if len(results) == 0 {
		mockResults := getMockResults()
		results = append(results, mockResults...)
	}
	paaresults := GetPAAResultsForAllProductsForDownstreamCalcs(runDate)
	if len(paaresults) == 0 {
		mockResult := models.PAAResult{}
		paaresults = append(paaresults, mockResult)
	}
	licBuildupResults := GetLicResultsForAllProductsForDownstreamCalcs(runDate)
	if len(licBuildupResults) == 0 {
		mockResult := getLICMockResults()
		licBuildupResults = append(licBuildupResults, mockResult...)
	}
	journals := GenerateJournalTransactionReportByGroup(results, paaresults, licBuildupResults)
	report = GenerateSubLedgerReportByGroup(journals, results, paaresults)

	prodList := GetJournalEntryProducts(runDate)

	res["steps"] = report
	res["products"] = prodList

	return res
}

func GenerateSubLedgerReportsforOneProduct(runDate, prodCode string) map[string]interface{} {

	var res = make(map[string]interface{})
	var ledgers []models.SubLedgerReportEntry
	results := GetAosStepResultsForOneProductForDownstreamCalcs(runDate, prodCode)
	if len(results) == 0 {
		mockResults := getMockResults()
		results = append(results, mockResults...)
	}
	paaresults := GetPAAResultsForOneProductForDownstreamCalcs(runDate, prodCode)
	if len(paaresults) == 0 {
		mockResult := models.PAAResult{}
		paaresults = append(paaresults, mockResult)
	}
	licBuildupResults := GetLicResultsForOneProductForDownstreamCalcs(runDate, prodCode)
	if len(licBuildupResults) == 0 {
		mockResult := getLICMockResults()
		licBuildupResults = append(licBuildupResults, mockResult...)
	}
	journals := GenerateJournalTransactionReportByGroup(results, paaresults, licBuildupResults)
	ledgers = GenerateSubLedgerReportByGroup(journals, results, paaresults)
	groupList := GetGroups(runDate, prodCode)

	res["steps"] = ledgers
	res["groups"] = groupList

	return res
}

func GenerateTrialBalanceReports(ifrs17Group, runDate string) []models.TrialBalanceReportEntry {
	results := GetAosStepResultsForGroup(ifrs17Group, runDate)
	if len(results) == 0 {
		mockResults := getMockResults()
		results = append(results, mockResults...)
	}
	paaresults := GetPAAResultsForGroup(ifrs17Group, runDate)
	if len(paaresults) == 0 {
		mockResult := models.PAAResult{}
		paaresults = append(paaresults, mockResult)
	}
	licBuildupResults := GetLicResultsForAllProductsForDownstreamCalcs(runDate)
	if len(licBuildupResults) == 0 {
		mockResult := getLICMockResults()
		licBuildupResults = append(licBuildupResults, mockResult...)
	}
	res := GenerateJournalTransactionReportByGroup(results, paaresults, licBuildupResults)
	res2 := GenerateSubLedgerReportByGroup(res, results, paaresults)
	report := GenerateTrialBalanceReportByGroup(res2, results, paaresults)

	return report
}

func GenerateBalanceSheetSummariesForAll(runDate, productCode, group string) map[string]interface{} {
	var results = make(map[string]interface{})
	var records []models.BalanceSheetRecord
	var prevRecords []models.BalanceSheetRecord
	var balanceSheetSummaries []models.BalanceSheetSummaryRecord

	parts := strings.Split(runDate, "-")
	valYear, _ := strconv.Atoi(parts[0])
	prevValYear := strconv.Itoa(valYear - 1)
	prevrunDate := prevValYear + "-" + parts[1]

	if runDate != "" && productCode == "" && group == "" {
		err := DB.Where("date = ?", runDate).Find(&records).Error
		if err != nil {
			fmt.Println(err)
		}
		err = DB.Where("date = ?", prevrunDate).Find(&prevRecords).Error
		if err != nil {
			fmt.Println(err)
		}
	}

	if runDate != "" && productCode != "" && group == "" {
		err := DB.Where("date = ? AND product_code = ?", runDate, productCode).Find(&records).Error
		if err != nil {
			fmt.Println(err)
		}

		err = DB.Where("date = ? AND product_code = ?", prevrunDate, productCode).Find(&prevRecords).Error
		if err != nil {
			fmt.Println(err)
		}
	}

	if runDate != "" && productCode != "" && group != "" {
		err := DB.Where("date = ? AND product_code = ? AND ifrs17_group = ?", runDate, productCode, group).Find(&records).Error
		if err != nil {
			fmt.Println(err)
		}

		err = DB.Where("date = ? AND product_code = ? AND ifrs17_group = ?", prevrunDate, productCode, group).Find(&prevRecords).Error
		if err != nil {
			fmt.Println(err)
		}
	}

	// Start calculations

	// Assets

	var assetLine models.BalanceSheetSummaryRecord
	var contractCostsLine models.BalanceSheetSummaryRecord
	var insuranceContractAssetsLine models.BalanceSheetSummaryRecord
	var reinsuranceContractAssetsLine models.BalanceSheetSummaryRecord
	var deferredTaxAssetsLine models.BalanceSheetSummaryRecord
	var otherAssetsLine models.BalanceSheetSummaryRecord
	var liabilitiesLine models.BalanceSheetSummaryRecord
	var insuranceContractLiabilitiesLine models.BalanceSheetSummaryRecord
	var reinsuranceContractLiabilitiesLine models.BalanceSheetSummaryRecord
	var investmentContractsLine models.BalanceSheetSummaryRecord
	var deferredTaxLiabilitiesLine models.BalanceSheetSummaryRecord
	var otherLiabilitiesLine models.BalanceSheetSummaryRecord

	for i, _ := range records {

		// Contract costs for investment management services
		contractCostsLine.AccountName = "Contract costs for investment management services"
		contractCostsLine.CurrentYear += 0
		contractCostsLine.Notes = "DAC"

		// Insurance contract assets
		insuranceContractAssetsLine.AccountName = "Insurance contract assets"
		insuranceContractAssetsLine.CurrentYear += records[i].BELInflow
		insuranceContractAssetsLine.Notes = "Insurance contract assets"

		// Reinsurance contract assets
		reinsuranceContractAssetsLine.AccountName = "Reinsurance contract assets"
		reinsuranceContractAssetsLine.CurrentYear += records[i].Treaty1BELOutflow + records[i].Treaty1BELOutflow + records[i].Treaty1BELOutflow +
			records[i].Treaty1RiskAdjustment + records[i].Treaty2RiskAdjustment + records[i].Treaty3RiskAdjustment
		reinsuranceContractAssetsLine.Notes = "Reinsurance contract assets"

		// Deffered tax assets
		deferredTaxAssetsLine.AccountName = "Deferred tax assets"
		deferredTaxAssetsLine.CurrentYear += 0
		deferredTaxAssetsLine.Notes = "Deferred tax assets"

		// Other assets
		otherAssetsLine.AccountName = "Other assets"
		otherAssetsLine.CurrentYear += 0
		otherAssetsLine.Notes = "Other assets"

		assetLine.AccountName = "Assets"
		assetLine.CurrentYear = contractCostsLine.CurrentYear + insuranceContractAssetsLine.CurrentYear + reinsuranceContractAssetsLine.CurrentYear +
			deferredTaxAssetsLine.CurrentYear + otherAssetsLine.CurrentYear
		assetLine.Notes = "Total Assets"

		// Insurance contract liabilities
		insuranceContractLiabilitiesLine.AccountName = "Insurance contract liabilities"
		// Do calculations for the following lines
		insuranceContractLiabilitiesLine.CurrentYear += records[i].BELOutflow + records[i].RiskAdjustment + records[i].PostTransitionCsm + records[i].FullyRetrospectiveCsm + records[i].ModifiedRetrospectiveCsm + records[i].FairValueCsm
		insuranceContractLiabilitiesLine.Notes = "Insurance contract liabilities"

		// Reinsurance contract liabilities
		reinsuranceContractLiabilitiesLine.AccountName = "Reinsurance contract liabilities"
		reinsuranceContractLiabilitiesLine.CurrentYear += records[i].Treaty1BELInflow + records[i].Treaty2BELInflow + records[i].Treaty3BELInflow
		reinsuranceContractLiabilitiesLine.Notes = "Reinsurance contract liabilities"

		// Investment contracts
		investmentContractsLine.AccountName = "Investment contracts"
		investmentContractsLine.CurrentYear += 0
		investmentContractsLine.Notes = "Investment contracts"

		// Deferred tax liabilities
		deferredTaxLiabilitiesLine.AccountName = "Deferred tax liabilities"
		deferredTaxLiabilitiesLine.CurrentYear += 0
		deferredTaxLiabilitiesLine.Notes = "Deferred tax liabilities"

		// Other liabilities
		otherLiabilitiesLine.AccountName = "Other liabilities"
		otherLiabilitiesLine.CurrentYear += 0
		otherLiabilitiesLine.Notes = "Other liabilities"

		// Liabilities
		liabilitiesLine.AccountName = "Liabilities"
		liabilitiesLine.CurrentYear = insuranceContractLiabilitiesLine.CurrentYear + reinsuranceContractLiabilitiesLine.CurrentYear +
			investmentContractsLine.CurrentYear + deferredTaxLiabilitiesLine.CurrentYear + otherLiabilitiesLine.CurrentYear
		liabilitiesLine.Notes = "Total Liabilities"

	}

	//Previous records
	for i, _ := range prevRecords {

		// Contract costs for investment management services
		contractCostsLine.PreviousYear += 0

		// Insurance contract assets
		insuranceContractAssetsLine.PreviousYear += prevRecords[i].BELInflow

		// Reinsurance contract assets
		reinsuranceContractAssetsLine.PreviousYear += prevRecords[i].Treaty1BELOutflow + prevRecords[i].Treaty1BELOutflow + prevRecords[i].Treaty1BELOutflow +
			prevRecords[i].Treaty1RiskAdjustment + prevRecords[i].Treaty2RiskAdjustment + prevRecords[i].Treaty3RiskAdjustment

		// Deffered tax assets
		deferredTaxAssetsLine.PreviousYear += 0

		// Other assets
		otherAssetsLine.PreviousYear += 0

		//Assets
		assetLine.PreviousYear += contractCostsLine.PreviousYear + insuranceContractAssetsLine.PreviousYear + reinsuranceContractAssetsLine.PreviousYear +
			deferredTaxAssetsLine.PreviousYear + otherAssetsLine.PreviousYear

		// Insurance contract liabilities
		insuranceContractLiabilitiesLine.PreviousYear += prevRecords[i].BELOutflow + prevRecords[i].RiskAdjustment + prevRecords[i].PostTransitionCsm + prevRecords[i].FullyRetrospectiveCsm + prevRecords[i].ModifiedRetrospectiveCsm + prevRecords[i].FairValueCsm

		// Reinsurance contract liabilities
		reinsuranceContractLiabilitiesLine.PreviousYear += prevRecords[i].Treaty1BELInflow + prevRecords[i].Treaty2BELInflow + prevRecords[i].Treaty3BELInflow

		// Investment contracts
		investmentContractsLine.PreviousYear += 0
		// Deferred tax liabilities
		deferredTaxLiabilitiesLine.PreviousYear += 0

		// Other liabilities
		otherLiabilitiesLine.PreviousYear += 0

		// Liabilities
		liabilitiesLine.PreviousYear = insuranceContractLiabilitiesLine.PreviousYear + reinsuranceContractLiabilitiesLine.PreviousYear +
			investmentContractsLine.PreviousYear + deferredTaxLiabilitiesLine.PreviousYear + otherLiabilitiesLine.PreviousYear

	}

	balanceSheetSummaries = append(balanceSheetSummaries, assetLine)
	balanceSheetSummaries = append(balanceSheetSummaries, contractCostsLine)
	balanceSheetSummaries = append(balanceSheetSummaries, insuranceContractAssetsLine)
	balanceSheetSummaries = append(balanceSheetSummaries, reinsuranceContractAssetsLine)
	balanceSheetSummaries = append(balanceSheetSummaries, deferredTaxAssetsLine)
	balanceSheetSummaries = append(balanceSheetSummaries, otherAssetsLine)
	balanceSheetSummaries = append(balanceSheetSummaries, liabilitiesLine)
	balanceSheetSummaries = append(balanceSheetSummaries, insuranceContractLiabilitiesLine)
	balanceSheetSummaries = append(balanceSheetSummaries, reinsuranceContractLiabilitiesLine)
	balanceSheetSummaries = append(balanceSheetSummaries, investmentContractsLine)
	balanceSheetSummaries = append(balanceSheetSummaries, deferredTaxLiabilitiesLine)
	balanceSheetSummaries = append(balanceSheetSummaries, otherLiabilitiesLine)

	results["summaries"] = balanceSheetSummaries

	var prodlist []string
	DB.Table("balance_sheet_records").Where("date = ?", runDate).Distinct("product_code").Find(&prodlist)

	results["products"] = prodlist

	if productCode != "" {
		var groupList []string
		DB.Table("balance_sheet_records").Where("date = ? AND product_code = ?", runDate, productCode).Distinct("ifrs17_group").Find(&groupList)
		results["groups"] = groupList
	}

	results["headers"] = []string{"Account Name", runDate, prevrunDate, "Notes"}
	return results
}

func GenerateBalanceSheetReportsforAllProducts(runDate string) map[string]interface{} {
	var results = make(map[string]interface{})
	var balanceSheets []models.BalanceSheetRecord

	err := DB.Where("date = ?", runDate).Find(&balanceSheets).Error
	if err != nil {
		fmt.Println(err)
	}

	results["reports"] = balanceSheets

	var prodlist []string
	DB.Table("balance_sheet_records").Where("date = ?", runDate).Distinct("product_code").Find(&prodlist)

	results["products"] = prodlist
	return results
}

func GenerateBalanceSheetReportForProduct(runDate, prodCode string) map[string]interface{} {
	var results = make(map[string]interface{})
	var balanceSheets []models.BalanceSheetRecord

	err := DB.Where("date = ? AND product_code = ?", runDate, prodCode).Find(&balanceSheets).Error
	if err != nil {
		fmt.Println(err)
	}

	results["reports"] = balanceSheets

	var prodlist []string
	DB.Table("balance_sheet_records").Where("date = ? AND product_code = ?", runDate, prodCode).Distinct("ifrs17_group").Find(&prodlist)

	results["groups"] = prodlist
	return results
}

func GenerateBalanceSheetReportForProductGroup(runDate, prodCode, ifrs17Group string) map[string]interface{} {
	var results = make(map[string]interface{})
	var balanceSheets []models.BalanceSheetRecord

	err := DB.Where("date = ? AND product_code = ? AND ifrs17_group = ?", runDate, prodCode, ifrs17Group).Find(&balanceSheets).Error
	if err != nil {
		fmt.Println(err)
	}

	results["reports"] = balanceSheets

	var prodlist []string
	DB.Table("balance_sheet_records").Where("date = ? AND product_code = ? AND ifrs17_group = ?", runDate, prodCode, ifrs17Group).Distinct("ifrs17_group").Find(&prodlist)

	results["groups"] = prodlist
	return results
}

func GenerateTrialBalanceReportsforAllProducts(runDate string) map[string]interface{} {

	var res = make(map[string]interface{})
	var report []models.TrialBalanceReportEntry

	results := GetAosStepResultsForAllProductsForDownstreamCalcs(runDate)
	if len(results) == 0 {
		mockResults := getMockResults()
		results = append(results, mockResults...)
	}
	paaresults := GetPAAResultsForAllProductsForDownstreamCalcs(runDate)
	if len(paaresults) == 0 {
		mockResult := models.PAAResult{}
		paaresults = append(paaresults, mockResult)
	}
	licBuildupResults := GetLicResultsForAllProductsForDownstreamCalcs(runDate)
	if len(licBuildupResults) == 0 {
		mockResult := getLICMockResults()
		licBuildupResults = append(licBuildupResults, mockResult...)
	}
	journals := GenerateJournalTransactionReportByGroup(results, paaresults, licBuildupResults)
	ledgers := GenerateSubLedgerReportByGroup(journals, results, paaresults)
	report = GenerateTrialBalanceReportByGroup(ledgers, results, paaresults)

	prodList := GetJournalEntryProducts(runDate)

	res["steps"] = report
	res["products"] = prodList

	return res
}

func GenerateTrialBalanceReportsforOneProduct(runDate, prodCode string) map[string]interface{} {

	var res = make(map[string]interface{})
	var trialBalance []models.TrialBalanceReportEntry
	results := GetAosStepResultsForOneProductForDownstreamCalcs(runDate, prodCode)
	if len(results) == 0 {
		mockResults := getMockResults()
		results = append(results, mockResults...)
	}

	paaresults := GetPAAResultsForOneProductForDownstreamCalcs(runDate, prodCode)
	if len(paaresults) == 0 {
		mockResult := models.PAAResult{}
		paaresults = append(paaresults, mockResult)
	}

	licBuildupResults := GetLicResultsForOneProductForDownstreamCalcs(runDate, prodCode)
	if len(licBuildupResults) == 0 {
		mockResult := getLICMockResults()
		licBuildupResults = append(licBuildupResults, mockResult...)
	}
	journals := GenerateJournalTransactionReportByGroup(results, paaresults, licBuildupResults)
	ledgers := GenerateSubLedgerReportByGroup(journals, results, paaresults)
	trialBalance = GenerateTrialBalanceReportByGroup(ledgers, results, paaresults)
	groupList := GetGroups(runDate, prodCode)

	res["steps"] = trialBalance
	res["groups"] = groupList

	return res
}

func GetJournalEntryProducts(runDate string) []string {
	var results []string
	DB.Table("journal_transactions").Where("run_date = ?", runDate).Distinct("product_code").Find(&results)
	return results
}

func GetGroups(runDate, prodCode string) []string {
	var results []string
	DB.Table("journal_transactions").Where("run_date = ? and product_code = ?", runDate, prodCode).Distinct("ifrs17_group").Find(&results)
	return results
}

func GenerateJournalTransactionReportByGroup(stepResults []models.AOSStepResult, paaresults []models.PAAResult, licBuildupResults []models.LicBuildupResult) []models.JTReportEntry {
	// The assumption we will work with is that the stepResults have been
	// filtered by product_code and ifrs17_group

	var reportEntries []models.JTReportEntry
	//For now we calculate each entry manually...

	var tempBelOutflowBuildup float64
	var tempBelInflowBuildup float64
	var tempRiskAdjustmentBuildup float64
	var tempCSMChangeBuildup float64
	var tempLossComponentChangeBuildup float64
	var tempInvestmentFinanceExpense float64

	if len(stepResults) > 0 {
		for _, step := range stepResults {
			var csmRun models.CsmRun
			err := DB.Where("id = ?", step.CsmRunID).First(&csmRun).Error
			if err != nil {
				fmt.Println(err)
			}

			if step.Name == "C/F_Lockedin" { // skip these two steps.
				continue
			}

			if step.Name == "B/F_Current" {
				var entry1, entry2, entry3, entry4, entry5, entry6, entry7 models.JTReportEntry
				entry1.ProductCode = step.ProductCode
				entry1.IFRS17Group = step.IFRS17Group
				populateEntryFields(1, "GMM/VFA", "BEL_Inflow", 1001, "B/S", "Debit", "Asset", 1, &entry1)
				entry1.Debit = 0 //step.BelInflow

				entry2.ProductCode = step.ProductCode
				entry2.IFRS17Group = step.IFRS17Group
				populateEntryFields(1, "GMM/VFA", "DAC", 1002, "B/S", "Debit", "Asset", 1, &entry2)
				entry2.Debit = 0 //step.DACBuildup

				entry3.ProductCode = step.ProductCode
				entry3.IFRS17Group = step.IFRS17Group
				populateEntryFields(1, "GMM/VFA", "BEL_Outflow", 2001, "B/S", "Credit", "Liability", 1, &entry3)
				entry3.Credit = 0 //step.BelOutflow

				entry4.ProductCode = step.ProductCode
				entry4.IFRS17Group = step.IFRS17Group
				populateEntryFields(1, "GMM/VFA", "Risk_Adjustment", 2002, "B/S", "Credit", "Liability", 1, &entry4)
				entry4.Credit = 0 //step.RiskAdjustment

				entry5.ProductCode = step.ProductCode
				entry5.IFRS17Group = step.IFRS17Group
				populateEntryFields(1, "GMM/VFA", "CSM", 2003, "B/S", "Credit", "Liability", 1, &entry5)
				entry5.Credit = 0 //step.CSMBuildup

				entry6.ProductCode = step.ProductCode
				entry6.IFRS17Group = step.IFRS17Group
				populateEntryFields(1, "GMM/VFA", "Acquisition_Cashflows", 2007, "B/S", "Credit", "Liability", 1, &entry6)
				entry6.Credit = 0

				entry7.ProductCode = step.ProductCode
				entry7.IFRS17Group = step.IFRS17Group
				populateEntryFields(1, "GMM/VFA", "Loss_Component", 5007, "B/S", "Credit", "Liability", 1, &entry7)
				entry7.Credit = 0 //step.LossComponentBuildup

				reportEntries = append(reportEntries, entry1, entry2, entry3, entry4, entry5, entry6, entry7)
			}

			if step.Name == "B/F_Lockedin" { // Step 2
				var entry8, entry9, entry10, entry11, entry12, entry13 models.JTReportEntry

				//Investment Finance or Expense
				entry8.ProductCode = step.ProductCode
				entry8.IFRS17Group = step.IFRS17Group
				populateEntryFields(2, "GMM/VFA", "BEL_Outflow", 2001, "B/S", "Credit", "Liability", 2, &entry8)
				entry8.Debit = math.Max(step.BelOutflowChange, 0)
				entry8.Credit = math.Max(-step.BelOutflowChange, 0)

				entry9.ProductCode = step.ProductCode
				entry9.IFRS17Group = step.IFRS17Group
				populateEntryFields(2, "GMM/VFA", "Insurance_Finance_Income_or_Expense", 420, "PnL", "Credit", "Insurance_Finance_Income_or_(Expense)", 2, &entry9)
				entry9.Debit = math.Max(-step.BelOutflowChange, 0)
				entry9.Credit = math.Max(step.BelOutflowChange, 0)

				entry10.ProductCode = step.ProductCode
				entry10.IFRS17Group = step.IFRS17Group
				populateEntryFields(2, "GMM/VFA", "Risk_Adjustment", 2002, "B/S", "Debit", "Liability", 3, &entry10)
				entry10.Debit = math.Max(step.RiskAdjustmentChange, 0)
				entry10.Credit = math.Max(-step.RiskAdjustmentChange, 0)

				entry11.ProductCode = step.ProductCode
				entry11.IFRS17Group = step.IFRS17Group
				populateEntryFields(2, "GMM/VFA", "Insurance_Finance_Income_or_Expense", 420, "PnL", "Credit", "Insurance_Finance_Income_or_(Expense)", 3, &entry11)
				entry11.Debit = math.Max(-step.RiskAdjustmentChange, 0)
				entry11.Credit = math.Max(step.RiskAdjustmentChange, 0)

				entry12.ProductCode = step.ProductCode
				entry12.IFRS17Group = step.IFRS17Group
				populateEntryFields(2, "GMM/VFA", "BEL_Inflow", 1001, "B/S", "Debit", "Asset", 4, &entry12)
				entry12.Credit = math.Max(step.BelInflowChange, 0)
				entry12.Debit = math.Max(-step.BelInflowChange, 0)

				entry13.ProductCode = step.ProductCode
				entry13.IFRS17Group = step.IFRS17Group
				populateEntryFields(2, "GMM/VFA", "Insurance_Finance_Income_or_Expense", 420, "PnL", "Credit", "Insurance_Finance_Income_or_(Expense)", 4, &entry13)
				entry13.Credit = math.Max(-step.BelInflowChange, 0)
				entry13.Debit = math.Max(step.BelInflowChange, 0)

				tempInvestmentFinanceExpense += step.LiabilityChange

				reportEntries = append(reportEntries, entry8, entry9, entry10, entry11, entry12, entry13)
			}

			if step.Name == "Initial_Recog" {
				var entry1, entry2, entry3, entry4, entry5, entry6, entry7, entry8, entry9 models.JTReportEntry

				entry1.ProductCode = step.ProductCode
				entry1.IFRS17Group = step.IFRS17Group
				populateEntryFields(3, "GMM/VFA", "BEL_Inflow", 1001, "B/S", "Debit", "Asset", 5, &entry1)
				entry1.Debit = step.BelInflow

				entry2.ProductCode = step.ProductCode
				entry2.IFRS17Group = step.IFRS17Group
				populateEntryFields(3, "GMM/VFA", "BEL_Outflow", 2001, "B/S", "Credit", "Liability", 5, &entry2)
				entry2.Credit = step.BelOutflowExclLC

				entry3.ProductCode = step.ProductCode
				entry3.IFRS17Group = step.IFRS17Group
				populateEntryFields(3, "GMM/VFA", "Risk_Adjustment", 2002, "B/S", "Credit", "Liability", 5, &entry3)
				entry3.Credit = step.RiskAdjustment
				entry3.Credit = step.RiskAdjExclLC

				entry4.ProductCode = step.ProductCode
				entry4.IFRS17Group = step.IFRS17Group
				populateEntryFields(3, "GMM/VFA", "CSM", 2003, "B/S", "Credit", "Liability", 5, &entry4)
				entry4.Credit = step.CSMChange

				entry5.ProductCode = step.ProductCode
				entry5.IFRS17Group = step.IFRS17Group
				populateEntryFields(3, "GMM/VFA", "Loss_Component_NB", 201, "PnL", "Debit", "Insurance_Service_Expense", 6, &entry5)
				entry5.Debit = step.LossComponentChange

				entry6.ProductCode = step.ProductCode
				entry6.IFRS17Group = step.IFRS17Group
				populateEntryFields(3, "GMM/VFA", "BEL_Outflow", 2001, "B/S", "Credit", "Liability", 6, &entry6)
				if step.LossComponentChange != 0 {
					entry6.Credit = math.Max(step.BelOutflow-step.BelOutflowExclLC, 0)
				} else {
					entry6.Credit = 0
				}

				entry7.ProductCode = step.ProductCode
				entry7.IFRS17Group = step.IFRS17Group
				populateEntryFields(3, "GMM/VFA", "Risk_Adjustment", 2002, "B/S", "Credit", "Liability", 6, &entry7)
				if step.LossComponentChange > 0 {
					entry7.Credit = math.Max(step.RiskAdjustment-step.RiskAdjExclLC, 0)
				} else {
					entry7.Credit = 0
				}

				//DAC
				entry8.ProductCode = step.ProductCode
				entry8.IFRS17Group = step.IFRS17Group
				populateEntryFields(3, "GMM/VFA", "DAC", 1002, "B/S", "Debit", "Asset", 7, &entry8)
				entry8.Debit = 0 //?

				entry9.ProductCode = step.ProductCode
				entry9.IFRS17Group = step.IFRS17Group
				populateEntryFields(3, "GMM/VFA", "Acquisition_Cashflows", 2007, "B/S", "Credit", "Liability", 7, &entry9)
				entry9.Credit = 0 //?

				reportEntries = append(reportEntries, entry1, entry2, entry3, entry4, entry5, entry6, entry7, entry8, entry9)
			}

			if step.Name == "Exp_RA_Release" {
				var entry1, entry2, entry3, entry4 models.JTReportEntry

				entry1.ProductCode = step.ProductCode
				entry1.IFRS17Group = step.IFRS17Group
				populateEntryFields(4, "GMM/VFA", "Risk_Adjustment", 2002, "B/S", "Credit", "Liability", 8, &entry1)
				entry1.Debit = math.Abs(step.RiskAdjustmentChange)

				entry2.ProductCode = step.ProductCode
				entry2.IFRS17Group = step.IFRS17Group
				populateEntryFields(4, "GMM/VFA", "Risk_Adjustment_Release", 101, "PnL", "Credit", "Insurance Revenue", 8, &entry2)
				entry2.Credit = math.Abs(step.RiskAdjustmentChange)

				entry3.ProductCode = step.ProductCode
				entry3.IFRS17Group = step.IFRS17Group
				populateEntryFields(4, "GMM/VFA", "Loss_Component", 5007, "Suspense", "Credit", "Suspense", 9, &entry3)
				entry3.Debit = math.Abs(step.RiskAdjustmentChange * step.LcSar)

				entry4.ProductCode = step.ProductCode
				entry4.IFRS17Group = step.IFRS17Group
				populateEntryFields(4, "GMM/VFA", "Risk_Adjustment_Release_SAR_LC", 5001, "Suspense", "Credit", "Suspense", 9, &entry4)
				entry4.Credit = math.Abs(step.RiskAdjustmentChange * step.LcSar)

				reportEntries = append(reportEntries, entry1, entry2, entry3, entry4)
			}

			if step.Name == "Exp_Mort" {
				var entry1, entry2, entry3, entry4 models.JTReportEntry

				entry1.ProductCode = step.ProductCode
				entry1.IFRS17Group = step.IFRS17Group
				populateEntryFields(5, "GMM/VFA", "BEL_Outflow", 2001, "B/S", "Credit", "Liability", 10, &entry1)
				entry1.Debit = step.PNLChange

				entry2.ProductCode = step.ProductCode
				entry2.IFRS17Group = step.IFRS17Group
				populateEntryFields(5, "GMM/VFA", "Expected_Mortality_Claims", 102, "PnL", "Credit", "Insurance Revenue", 10, &entry2)
				entry2.Credit = step.PNLChange

				entry3.ProductCode = step.ProductCode
				entry3.IFRS17Group = step.IFRS17Group
				populateEntryFields(5, "GMM/VFA", "Loss_Component", 5007, "Suspense", "Credit", "Suspense", 11, &entry3)
				entry3.Debit = step.PNLChange * step.LcSar

				entry4.ProductCode = step.ProductCode
				entry4.IFRS17Group = step.IFRS17Group
				populateEntryFields(5, "GMM/VFA", "Expected_Mortality_Claims_SAR_LC", 5002, "Suspense", "Credit", "Suspense", 11, &entry4)
				entry4.Credit = step.PNLChange * step.LcSar

				reportEntries = append(reportEntries, entry1, entry2, entry3, entry4)
			}

			if step.Name == "Exp_Exp" {
				var entry1, entry2, entry3, entry4 models.JTReportEntry

				entry1.ProductCode = step.ProductCode
				entry1.IFRS17Group = step.IFRS17Group
				populateEntryFields(6, "GMM/VFA", "BEL_Outflow", 2001, "B/S", "Credit", "Liability", 12, &entry1)
				entry1.Debit = step.PNLChange

				entry2.ProductCode = step.ProductCode
				entry2.IFRS17Group = step.IFRS17Group
				populateEntryFields(6, "GMM/VFA", "Expected_Expenses", 103, "PnL", "Credit", "Insurance Revenue", 12, &entry2)
				entry2.Credit = step.PNLChange

				entry3.ProductCode = step.ProductCode
				entry3.IFRS17Group = step.IFRS17Group
				populateEntryFields(6, "GMM/VFA", "Loss_Component", 5007, "LC_Suspense", "Credit", "Suspense", 13, &entry3)
				entry3.Debit = step.PNLChange * step.LcSar

				entry4.ProductCode = step.ProductCode
				entry4.IFRS17Group = step.IFRS17Group
				populateEntryFields(6, "GMM/VFA", "Expected_Expenses_SAR_LC", 5003, "Suspense", "Credit", "Suspense", 13, &entry4)
				entry4.Credit = step.PNLChange * step.LcSar

				reportEntries = append(reportEntries, entry1, entry2, entry3, entry4)
			}

			if step.Name == "Exp_Retrenchment" {
				var entry1, entry2, entry3, entry4 models.JTReportEntry

				entry1.ProductCode = step.ProductCode
				entry1.IFRS17Group = step.IFRS17Group
				populateEntryFields(7, "GMM/VFA", "BEL_Outflow", 2001, "B/S", "Credit", "Liability", 14, &entry1)
				entry1.Debit = step.PNLChange

				entry2.ProductCode = step.ProductCode
				entry2.IFRS17Group = step.IFRS17Group
				populateEntryFields(7, "GMM/VFA", "Expected_Retrenchment", 104, "PnL", "Credit", "Insurance Revenue", 14, &entry2)
				entry2.Credit = step.PNLChange

				entry3.ProductCode = step.ProductCode
				entry3.IFRS17Group = step.IFRS17Group
				populateEntryFields(7, "GMM/VFA", "Loss_Component", 5007, "Suspense", "Credit", "Suspense", 15, &entry3)
				entry3.Debit = step.PNLChange * step.LcSar

				entry4.ProductCode = step.ProductCode
				entry4.IFRS17Group = step.IFRS17Group
				populateEntryFields(7, "GMM/VFA", "Expected_Retrenchment_SAR_LC", 5004, "Suspense", "Credit", "Suspense", 15, &entry4)
				entry4.Credit = step.PNLChange * step.LcSar

				reportEntries = append(reportEntries, entry1, entry2, entry3, entry4)
			}

			if step.Name == "Exp_Morbidity" {
				var entry1, entry2, entry3, entry4 models.JTReportEntry

				entry1.ProductCode = step.ProductCode
				entry1.IFRS17Group = step.IFRS17Group
				populateEntryFields(8, "GMM/VFA", "BEL_Outflow", 2001, "B/S", "Credit", "Liability", 16, &entry1)
				entry1.Debit = step.PNLChange

				entry2.ProductCode = step.ProductCode
				entry2.IFRS17Group = step.IFRS17Group
				populateEntryFields(8, "GMM/VFA", "Expected_Morbidity", 105, "PnL", "Credit", "Insurance Revenue", 16, &entry2)
				entry2.Credit = step.PNLChange

				entry3.ProductCode = step.ProductCode
				entry3.IFRS17Group = step.IFRS17Group
				populateEntryFields(8, "GMM/VFA", "Loss_Component", 5007, "LC_Suspense", "Credit", "Suspense", 17, &entry3)
				entry3.Debit = step.PNLChange * step.LcSar

				entry4.ProductCode = step.ProductCode
				entry4.IFRS17Group = step.IFRS17Group
				populateEntryFields(8, "GMM/VFA", "Expected_Morbidity_SAR_LC", 5005, "Suspense", "Credit", "Suspense", 17, &entry4)
				entry4.Credit = step.PNLChange * step.LcSar

				reportEntries = append(reportEntries, entry1, entry2, entry3, entry4)
			}

			if step.Name == "Interest_Accretion" {
				var entry1, entry2, entry3, entry4, entry5, entry6, entry7, entry8, entry9, entry10, entry11, entry12, entry13, entry14, entry15, entry16, entry17, entry18 models.JTReportEntry

				entry1.ProductCode = step.ProductCode
				entry1.IFRS17Group = step.IFRS17Group
				populateEntryFields(9, "GMM/VFA", "Cash_and_Cash_Equivalent", 1005, "B/S", "Debit", "Asset", 18, &entry1)
				entry1.Debit = step.CSMChange

				entry2.ProductCode = step.ProductCode
				entry2.IFRS17Group = step.IFRS17Group
				populateEntryFields(9, "GMM/VFA", "CSM", 2003, "B/S", "Credit", "Liability", 18, &entry2)
				entry2.Credit = step.CSMChange

				entry3.ProductCode = step.ProductCode
				entry3.IFRS17Group = step.IFRS17Group
				populateEntryFields(9, "GMM/VFA", "BEL_Inflow", 1001, "B/S", "Debit", "Asset", 19, &entry3)
				entry3.Debit = step.BelInflowChange
				//entry3.Debit = step.BelOutflowChange

				entry4.ProductCode = step.ProductCode
				entry4.IFRS17Group = step.IFRS17Group
				populateEntryFields(9, "GMM/VFA", "Insurance_Finance_Income_or_(Expense)", 420, "PnL", "Credit", "Investment_Income", 19, &entry4)
				entry4.Credit = step.BelInflowChange
				//entry4.Credit = step.BelOutflowChange

				entry5.ProductCode = step.ProductCode
				entry5.IFRS17Group = step.IFRS17Group
				populateEntryFields(9, "GMM/VFA", "Cash_and_Cash_Equivalent", 1005, "B/S", "Debit", "Asset", 20, &entry5)
				entry5.Debit = step.RiskAdjustmentChange

				entry6.ProductCode = step.ProductCode
				entry6.IFRS17Group = step.IFRS17Group
				populateEntryFields(9, "GMM/VFA", "Investment_Income", 410, "PnL", "Credit", "Investment_Income", 20, &entry6)
				entry6.Credit = step.RiskAdjustmentChange

				entry7.ProductCode = step.ProductCode
				entry7.IFRS17Group = step.IFRS17Group
				populateEntryFields(9, "GMM/VFA", "Insurance_Finance_Income_or_(Expense)", 420, "PnL", "Credit", "Investment_Income", 21, &entry7)
				entry7.Debit = step.BelOutflowChange

				entry8.ProductCode = step.ProductCode
				entry8.IFRS17Group = step.IFRS17Group
				populateEntryFields(9, "GMM/VFA", "BEL_Outflow", 2001, "B/S", "Credit", "Liability", 21, &entry8)
				entry8.Credit = step.BelOutflowChange

				//need to be added
				//entry7.ProductCode = step.ProductCode
				//entry7.IFRS17Group = step.IFRS17Group
				//populateEntryFields(9, "GMM/VFA", "Loss_Component_Int", 5006, "Suspense", "Credit", "Suspense", 17, &entry7)
				//entry7.Debit = step.LossComponentChange
				//
				//entry8.ProductCode = step.ProductCode
				//entry8.IFRS17Group = step.IFRS17Group
				//populateEntryFields(9, "GMM/VFA", "Loss_Component", 5007, "B/S", "Credit", "Suspense", 17, &entry8)
				//entry8.Credit = step.LossComponentChange

				//
				entry9.ProductCode = step.ProductCode
				entry9.IFRS17Group = step.IFRS17Group
				populateEntryFields(9, "GMM/VFA", "Premium_Debtors", 1006, "B/S", "Debit", "Asset", 22, &entry9)
				entry9.Debit = step.ExpectedCashInflow

				entry10.ProductCode = step.ProductCode
				entry10.IFRS17Group = step.IFRS17Group
				populateEntryFields(9, "GMM/VFA", "BEL_Inflow", 1001, "B/S", "Debit", "Asset", 22, &entry10)
				entry10.Credit = step.ExpectedCashInflow

				entry11.ProductCode = step.ProductCode
				entry11.IFRS17Group = step.IFRS17Group
				populateEntryFields(9, "GMM/VFA", "Premium_Debtors", 1006, "B/S", "Debit", "Asset", 23, &entry11)
				entry11.Debit = math.Max(step.ExperiencePremiumVariance, 0)
				entry11.Credit = math.Max(-step.ExperiencePremiumVariance, 0)

				entry12.ProductCode = step.ProductCode
				entry12.IFRS17Group = step.IFRS17Group
				populateEntryFields(9, "GMM/VFA", "Experience_Premium_Variance", 113, "PnL", "Credit", "Insurance Revenue", 23, &entry12)
				entry12.Debit = math.Max(-step.ExperiencePremiumVariance, 0)
				entry12.Credit = math.Max(step.ExperiencePremiumVariance, 0)

				entry13.ProductCode = step.ProductCode
				entry13.IFRS17Group = step.IFRS17Group
				populateEntryFields(9, "GMM/VFA", "Cash_and_Cash_Equivalent", 1005, "B/S", "Debit", "Asset", 24, &entry13)
				entry13.Debit = step.ActualPremium

				entry14.ProductCode = step.ProductCode
				entry14.IFRS17Group = step.IFRS17Group
				populateEntryFields(9, "GMM/VFA", "Premium_Debtors", 1006, "B/S", "Debit", "Asset", 24, &entry14)
				entry14.Credit = step.ActualPremium

				entry15.ProductCode = step.ProductCode
				entry15.IFRS17Group = step.IFRS17Group
				populateEntryFields(9, "GMM/VFA", "BEL_Outflow", 2001, "B/S", "Credit", "Liability", 24, &entry15)
				entry15.Debit = step.ExpectedCashOutflow - stepResults[4].PNLChange - stepResults[5].PNLChange - stepResults[6].PNLChange - stepResults[7].PNLChange

				entry16.ProductCode = step.ProductCode
				entry16.IFRS17Group = step.IFRS17Group
				populateEntryFields(9, "GMM/VFA", "Other_CashOutflows", 111, "B/S", "Credit", "Insurance Revenue", 24, &entry16)
				entry16.Credit = step.ExpectedCashOutflow - stepResults[4].PNLChange - stepResults[5].PNLChange - stepResults[6].PNLChange - stepResults[7].PNLChange

				entry17.ProductCode = step.ProductCode
				entry17.IFRS17Group = step.IFRS17Group
				populateEntryFields(9, "GMM/VFA", "Insurance_Finance_Income_or_(Expense)", 420, "PnL", "Credit", "Investment_Income", 25, &entry17)
				entry17.Debit = step.RiskAdjustmentChange

				entry18.ProductCode = step.ProductCode
				entry18.IFRS17Group = step.IFRS17Group
				populateEntryFields(9, "GMM/VFA", "Risk_Adjustment", 2002, "B/S", "Credit", "Liability", 25, &entry18)
				entry18.Credit = step.RiskAdjustmentChange

				tempInvestmentFinanceExpense += step.BestEstimateLiabilityChange + step.RiskAdjustmentChange

				reportEntries = append(reportEntries, entry1, entry2, entry3, entry4, entry5, entry6, entry7, entry8, entry9, entry10, entry11, entry12, entry13, entry14, entry15, entry16, entry17, entry18)
			}

			if step.Name == "Data_Change" {
				var entry1, entry2, entry3, entry4, entry5, entry6 models.JTReportEntry

				entry1.ProductCode = step.ProductCode
				entry1.IFRS17Group = step.IFRS17Group
				populateEntryFields(10, "GMM/VFA", "BEL_Outflow_Change", 2004, "B/S", "Credit", "Liability", 25, &entry1)
				entry1.Debit = math.Max(step.BelOutflowChange, 0)
				entry1.Credit = math.Max(-step.BelOutflowChange, 0)

				entry2.ProductCode = step.ProductCode
				entry2.IFRS17Group = step.IFRS17Group
				populateEntryFields(10, "GMM/VFA", "BEL_Outflow_Buildup", 2008, "B/S", "Credit", "Liability", 25, &entry2)
				entry2.Debit = math.Max(-step.BelOutflowChange, 0)
				entry2.Credit = math.Max(step.BelOutflowChange, 0)

				entry3.ProductCode = step.ProductCode
				entry3.IFRS17Group = step.IFRS17Group
				populateEntryFields(10, "GMM/VFA", "RA_Change", 2005, "B/S", "Credit", "Liability", 26, &entry3)
				entry3.Debit = math.Max(step.RiskAdjustmentChange, 0)
				entry3.Credit = math.Max(-step.RiskAdjustmentChange, 0)

				entry4.ProductCode = step.ProductCode
				entry4.IFRS17Group = step.IFRS17Group
				populateEntryFields(10, "GMM/VFA", "RA_Buildup", 2009, "B/S", "Credit", "Liability", 26, &entry4)
				entry4.Debit = math.Max(-step.RiskAdjustmentChange, 0)
				entry4.Credit = math.Max(step.RiskAdjustmentChange, 0)

				entry5.ProductCode = step.ProductCode
				entry5.IFRS17Group = step.IFRS17Group
				populateEntryFields(10, "GMM/VFA", "BEL_Inflow_Change", 1003, "B/S", "Debit", "Asset", 27, &entry5)
				entry5.Debit = math.Max(-step.BelInflowChange, 0)
				entry5.Credit = math.Max(step.BelInflowChange, 0)

				entry6.ProductCode = step.ProductCode
				entry6.IFRS17Group = step.IFRS17Group
				populateEntryFields(10, "GMM/VFA", "BEL_Inflow_Buildup", 1004, "B/S", "Debit", "Asset", 27, &entry6)
				entry6.Debit = math.Max(step.BelInflowChange, 0)
				entry6.Credit = math.Max(-step.BelInflowChange, 0)

				tempBelOutflowBuildup += step.BelOutflowChange
				tempRiskAdjustmentBuildup += step.RiskAdjustmentChange
				tempBelInflowBuildup += step.BelInflowChange
				tempCSMChangeBuildup += step.CSMChange
				tempLossComponentChangeBuildup += step.LossComponentChange

				reportEntries = append(reportEntries, entry1, entry2, entry3, entry4, entry5, entry6)
			}

			if step.Name == "NFA_Mort" {
				var entry1, entry2, entry3, entry4, entry5, entry6 models.JTReportEntry

				entry1.ProductCode = step.ProductCode
				entry1.IFRS17Group = step.IFRS17Group
				populateEntryFields(11, "GMM/VFA", "BEL_Outflow_Change", 2004, "B/S", "Credit", "Liability", 28, &entry1)
				entry1.Debit = math.Max(step.BelOutflowChange, 0)
				entry1.Credit = math.Max(-step.BelOutflowChange, 0)

				entry2.ProductCode = step.ProductCode
				entry2.IFRS17Group = step.IFRS17Group
				populateEntryFields(11, "GMM/VFA", "BEL_Outflow_Buildup", 2008, "B/S", "Credit", "Liability", 28, &entry2)
				entry2.Debit = math.Max(-step.BelOutflowChange, 0)
				entry2.Credit = math.Max(step.BelOutflowChange, 0)

				entry3.ProductCode = step.ProductCode
				entry3.IFRS17Group = step.IFRS17Group
				populateEntryFields(11, "GMM/VFA", "RA_Change", 2005, "B/S", "Credit", "Liability", 29, &entry3)
				entry3.Debit = math.Max(step.RiskAdjustmentChange, 0)
				entry3.Credit = math.Max(-step.RiskAdjustmentChange, 0)

				entry4.ProductCode = step.ProductCode
				entry4.IFRS17Group = step.IFRS17Group
				populateEntryFields(11, "GMM/VFA", "RA_Buildup", 2009, "B/S", "Credit", "Liability", 29, &entry4)
				entry4.Debit = math.Max(-step.RiskAdjustmentChange, 0)
				entry4.Credit = math.Max(step.RiskAdjustmentChange, 0)

				entry5.ProductCode = step.ProductCode
				entry5.IFRS17Group = step.IFRS17Group
				populateEntryFields(11, "GMM/VFA", "BEL_Inflow_Change", 1003, "B/S", "Debit", "Asset", 30, &entry5)
				entry5.Debit = math.Max(-step.BelInflowChange, 0)
				entry5.Credit = math.Max(step.BelInflowChange, 0)

				entry6.ProductCode = step.ProductCode
				entry6.IFRS17Group = step.IFRS17Group
				populateEntryFields(11, "GMM/VFA", "BEL_Inflow_Buildup", 1004, "B/S", "Debit", "Asset", 30, &entry6)
				entry6.Debit = math.Max(step.BelInflowChange, 0)
				entry6.Credit = math.Max(-step.BelInflowChange, 0)

				tempBelOutflowBuildup += step.BelOutflowChange
				tempRiskAdjustmentBuildup += step.RiskAdjustmentChange
				tempBelInflowBuildup += step.BelInflowChange
				tempCSMChangeBuildup += step.CSMChange
				tempLossComponentChangeBuildup += step.LossComponentChange

				reportEntries = append(reportEntries, entry1, entry2, entry3, entry4, entry5, entry6)
			}

			if step.Name == "NFA_Exp" {
				var entry1, entry2, entry3, entry4, entry5, entry6 models.JTReportEntry

				entry1.ProductCode = step.ProductCode
				entry1.IFRS17Group = step.IFRS17Group
				populateEntryFields(12, "GMM/VFA", "BEL_Outflow_Change", 2004, "B/S", "Credit", "Liability", 31, &entry1)
				entry1.Debit = math.Max(step.BelOutflowChange, 0)
				entry1.Credit = math.Max(-step.BelOutflowChange, 0)

				entry2.ProductCode = step.ProductCode
				entry2.IFRS17Group = step.IFRS17Group
				populateEntryFields(12, "GMM/VFA", "BEL_Outflow_Buildup", 2008, "B/S", "Credit", "Liability", 31, &entry2)
				entry2.Debit = math.Max(-step.BelOutflowChange, 0)
				entry2.Credit = math.Max(step.BelOutflowChange, 0)

				entry3.ProductCode = step.ProductCode
				entry3.IFRS17Group = step.IFRS17Group
				populateEntryFields(12, "GMM/VFA", "RA_Change", 2005, "B/S", "Credit", "Liability", 32, &entry3)
				entry3.Debit = math.Max(step.RiskAdjustmentChange, 0)
				entry3.Credit = math.Max(-step.RiskAdjustmentChange, 0)

				entry4.ProductCode = step.ProductCode
				entry4.IFRS17Group = step.IFRS17Group
				populateEntryFields(12, "GMM/VFA", "RA_Buildup", 2009, "B/S", "Credit", "Liability", 32, &entry4)
				entry4.Debit = math.Max(-step.RiskAdjustmentChange, 0)
				entry4.Credit = math.Max(step.RiskAdjustmentChange, 0)

				entry5.ProductCode = step.ProductCode
				entry5.IFRS17Group = step.IFRS17Group
				populateEntryFields(12, "GMM/VFA", "BEL_Inflow_Change", 1003, "B/S", "Debit", "Asset", 33, &entry5)
				entry5.Debit = math.Max(-step.BelInflowChange, 0)
				entry5.Credit = math.Max(step.BelInflowChange, 0)

				entry6.ProductCode = step.ProductCode
				entry6.IFRS17Group = step.IFRS17Group
				populateEntryFields(12, "GMM/VFA", "BEL_Inflow_Buildup", 1004, "B/S", "Debit", "Asset", 33, &entry6)
				entry6.Debit = math.Max(step.BelInflowChange, 0)
				entry6.Credit = math.Max(-step.BelInflowChange, 0)

				tempBelOutflowBuildup += step.BelOutflowChange
				tempRiskAdjustmentBuildup += step.RiskAdjustmentChange
				tempBelInflowBuildup += step.BelInflowChange
				tempCSMChangeBuildup += step.CSMChange
				tempLossComponentChangeBuildup += step.LossComponentChange

				reportEntries = append(reportEntries, entry1, entry2, entry3, entry4, entry5, entry6)
			}

			if step.Name == "NFA_Retrenchment" {
				var entry1, entry2, entry3, entry4, entry5, entry6 models.JTReportEntry

				entry1.ProductCode = step.ProductCode
				entry1.IFRS17Group = step.IFRS17Group
				populateEntryFields(13, "GMM/VFA", "BEL_Outflow_Change", 2004, "B/S", "Credit", "Liability", 34, &entry1)
				entry1.Debit = math.Max(step.BelOutflowChange, 0)
				entry1.Credit = math.Max(-step.BelOutflowChange, 0)

				entry2.ProductCode = step.ProductCode
				entry2.IFRS17Group = step.IFRS17Group
				populateEntryFields(13, "GMM/VFA", "BEL_Outflow_Buildup", 2008, "B/S", "Credit", "Liability", 34, &entry2)
				entry2.Debit = math.Max(-step.BelOutflowChange, 0)
				entry2.Credit = math.Max(step.BelOutflowChange, 0)

				entry3.ProductCode = step.ProductCode
				entry3.IFRS17Group = step.IFRS17Group
				populateEntryFields(13, "GMM/VFA", "RA_Change", 2005, "B/S", "Credit", "Liability", 35, &entry3)
				entry3.Debit = math.Max(step.RiskAdjustmentChange, 0)
				entry3.Credit = math.Max(-step.RiskAdjustmentChange, 0)

				entry4.ProductCode = step.ProductCode
				entry4.IFRS17Group = step.IFRS17Group
				populateEntryFields(13, "GMM/VFA", "RA_Buildup", 2009, "B/S", "Credit", "Liability", 35, &entry4)
				entry4.Debit = math.Max(-step.RiskAdjustmentChange, 0)
				entry4.Credit = math.Max(step.RiskAdjustmentChange, 0)

				entry5.ProductCode = step.ProductCode
				entry5.IFRS17Group = step.IFRS17Group
				populateEntryFields(13, "GMM/VFA", "BEL_Inflow_Change", 1003, "B/S", "Debit", "Asset", 36, &entry5)
				entry5.Debit = math.Max(-step.BelInflowChange, 0)
				entry5.Credit = math.Max(step.BelInflowChange, 0)

				entry6.ProductCode = step.ProductCode
				entry6.IFRS17Group = step.IFRS17Group
				populateEntryFields(13, "GMM/VFA", "BEL_Inflow_Buildup", 1004, "B/S", "Debit", "Asset", 36, &entry6)
				entry6.Debit = math.Max(step.BelInflowChange, 0)
				entry6.Credit = math.Max(-step.BelInflowChange, 0)

				tempBelOutflowBuildup += step.BelOutflowChange
				tempRiskAdjustmentBuildup += step.RiskAdjustmentChange
				tempBelInflowBuildup += step.BelInflowChange
				tempCSMChangeBuildup += step.CSMChange
				tempLossComponentChangeBuildup += step.LossComponentChange

				reportEntries = append(reportEntries, entry1, entry2, entry3, entry4, entry5, entry6)
			}

			if step.Name == "NFA_Morbidity" {
				var entry1, entry2, entry3, entry4, entry5, entry6, entry7, entry8, entry9, entry10, entry11 models.JTReportEntry

				entry1.ProductCode = step.ProductCode
				entry1.IFRS17Group = step.IFRS17Group
				populateEntryFields(14, "GMM/VFA", "BEL_Outflow_Change", 2004, "B/S", "Credit", "Liability", 37, &entry1)
				entry1.Debit = math.Max(step.BelOutflowChange, 0)
				entry1.Credit = math.Max(-step.BelOutflowChange, 0)

				entry2.ProductCode = step.ProductCode
				entry2.IFRS17Group = step.IFRS17Group
				populateEntryFields(14, "GMM/VFA", "BEL_Outflow_Buildup", 2008, "B/S", "Credit", "Liability", 37, &entry2)
				entry2.Debit = math.Max(-step.BelOutflowChange, 0)
				entry2.Credit = math.Max(step.BelOutflowChange, 0)

				entry3.ProductCode = step.ProductCode
				entry3.IFRS17Group = step.IFRS17Group
				populateEntryFields(14, "GMM/VFA", "RA_Change", 2005, "B/S", "Credit", "Liability", 38, &entry3)
				entry3.Debit = math.Max(step.RiskAdjustmentChange, 0)
				entry3.Credit = math.Max(-step.RiskAdjustmentChange, 0)

				entry4.ProductCode = step.ProductCode
				entry4.IFRS17Group = step.IFRS17Group
				populateEntryFields(14, "GMM/VFA", "RA_Buildup", 2009, "B/S", "Credit", "Liability", 38, &entry4)
				entry4.Debit = math.Max(-step.RiskAdjustmentChange, 0)
				entry4.Credit = math.Max(step.RiskAdjustmentChange, 0)

				entry5.ProductCode = step.ProductCode
				entry5.IFRS17Group = step.IFRS17Group
				populateEntryFields(14, "GMM/VFA", "BEL_Inflow_Change", 1003, "B/S", "Debit", "Asset", 39, &entry5)
				entry5.Debit = math.Max(-step.BelInflowChange, 0)
				entry5.Credit = math.Max(step.BelInflowChange, 0)

				entry6.ProductCode = step.ProductCode
				entry6.IFRS17Group = step.IFRS17Group
				populateEntryFields(14, "GMM/VFA", "BEL_Inflow_Buildup", 1004, "B/S", "Debit", "Asset", 39, &entry6)
				entry6.Debit = math.Max(step.BelInflowChange, 0)
				entry6.Credit = math.Max(-step.BelInflowChange, 0)

				tempBelOutflowBuildup += step.BelOutflowChange
				tempRiskAdjustmentBuildup += step.RiskAdjustmentChange
				tempBelInflowBuildup += step.BelInflowChange
				tempCSMChangeBuildup += step.CSMChange
				tempLossComponentChangeBuildup += step.LossComponentChange

				entry7.ProductCode = step.ProductCode
				entry7.IFRS17Group = step.IFRS17Group
				populateEntryFields(14, "GMM/VFA", "BEL_Outflow_Buildup", 2008, "B/S", "Credit", "Liability", 40, &entry7)
				entry7.Debit = math.Max(-tempBelOutflowBuildup, 0)
				entry7.Credit = math.Max(tempBelOutflowBuildup, 0)

				entry8.ProductCode = step.ProductCode
				entry8.IFRS17Group = step.IFRS17Group
				populateEntryFields(14, "GMM/VFA", "RA_Buildup", 2009, "B/S", "Credit", "Liability", 40, &entry8)
				entry8.Debit = math.Max(-tempRiskAdjustmentBuildup, 0)
				entry8.Credit = math.Max(tempRiskAdjustmentBuildup, 0)

				entry9.ProductCode = step.ProductCode
				entry9.IFRS17Group = step.IFRS17Group
				populateEntryFields(14, "GMM/VFA", "BEL_Inflow_Buildup", 1004, "B/S", "Debit", "Asset", 40, &entry9)
				entry9.Debit = math.Max(tempBelInflowBuildup, 0)
				entry9.Credit = math.Max(-tempBelInflowBuildup, 0)

				entry10.ProductCode = step.ProductCode
				entry10.IFRS17Group = step.IFRS17Group
				populateEntryFields(14, "GMM/VFA", "CSM", 2003, "B/S", "Credit", "Liability", 40, &entry10)
				entry10.Debit = math.Max(-tempCSMChangeBuildup, 0)
				entry10.Credit = math.Max(tempCSMChangeBuildup, 0)

				entry11.ProductCode = step.ProductCode
				entry11.IFRS17Group = step.IFRS17Group
				populateEntryFields(14, "GMM/VFA", "Loss_Component_Adjustment", 202, "PnL", "Credit", "Insurance_Service_Expense", 40, &entry11)
				//entry11.Debit = math.Max(tempLossComponentChangeBuildup, 0)
				entry11.Debit = math.Max(tempLossComponentChangeBuildup, 0)
				entry11.Credit = math.Max(-tempLossComponentChangeBuildup, 0)

				//same as loss component adjustment
				//entry12.ProductCode = step.ProductCode
				//entry12.IFRS17Group = step.IFRS17Group
				//populateEntryFields(14, "GMM/VFA", "Loss_Component_Reversal", 107, "PnL", "Debit", "Insurance_Revenue", 23, &entry12)
				////entry12.Credit = math.Max(-tempLossComponentChangeBuildup, 0)
				//entry12.Debit = math.Max(tempLossComponentChangeBuildup, 0)
				//entry12.Credit = math.Max(-tempLossComponentChangeBuildup, 0)

				reportEntries = append(reportEntries, entry1, entry2, entry3, entry4, entry5, entry6, entry7, entry8, entry9, entry10, entry11)
			}

			if step.Name == "FA" {
				var entry1, entry2, entry3, entry4, entry5, entry6 models.JTReportEntry

				entry1.ProductCode = step.ProductCode
				entry1.IFRS17Group = step.IFRS17Group
				populateEntryFields(15, "GMM/VFA", "BEL_Outflow", 2001, "B/S", "Credit", "Liability", 41, &entry1)
				entry1.Debit = math.Max(-step.BelOutflowChange, 0)
				entry1.Credit = math.Max(step.BelOutflowChange, 0)

				entry2.ProductCode = step.ProductCode
				entry2.IFRS17Group = step.IFRS17Group
				populateEntryFields(15, "GMM/VFA", "Insurance_Finance_Income_or_Expense", 420, "PnL", "Credit", "Insurance_Finance_Income_or_(Expense)", 41, &entry2)
				entry2.Debit = math.Max(step.BelOutflowChange, 0)
				entry2.Credit = math.Max(-step.BelOutflowChange, 0)

				entry3.ProductCode = step.ProductCode
				entry3.IFRS17Group = step.IFRS17Group
				populateEntryFields(15, "GMM/VFA", "Risk_Adjustment", 2002, "B/S", "Credit", "Liability", 42, &entry3)
				entry3.Debit = math.Max(-step.RiskAdjustmentChange, 0)
				entry3.Credit = math.Max(step.RiskAdjustmentChange, 0)

				entry4.ProductCode = step.ProductCode
				entry4.IFRS17Group = step.IFRS17Group
				populateEntryFields(15, "GMM/VFA", "Insurance_Finance_Income_or_Expense", 420, "PnL", "Credit", "Insurance_Finance_Income_or_(Expense)", 42, &entry4)
				entry4.Debit = math.Max(step.RiskAdjustmentChange, 0)
				entry4.Credit = math.Max(-step.RiskAdjustmentChange, 0)

				entry5.ProductCode = step.ProductCode
				entry5.IFRS17Group = step.IFRS17Group
				populateEntryFields(15, "GMM/VFA", "BEL_Inflow", 2001, "B/S", "Debit", "Asset", 43, &entry5)
				entry5.Debit = math.Max(step.BelInflowChange, 0)
				entry5.Credit = math.Max(-step.BelInflowChange, 0)

				entry6.ProductCode = step.ProductCode
				entry6.IFRS17Group = step.IFRS17Group
				populateEntryFields(15, "GMM/VFA", "Insurance_Finance_Income_or_Expense", 420, "PnL", "Credit", "Insurance_Finance_Income_or_(Expense)", 43, &entry6)
				entry6.Debit = math.Max(-step.BelInflowChange, 0)
				entry6.Credit = math.Max(step.BelInflowChange, 0)

				tempInvestmentFinanceExpense += step.BelOutflowChange + step.RiskAdjustmentChange - step.BelInflowChange

				reportEntries = append(reportEntries, entry1, entry2, entry3, entry4, entry5, entry6)
			}

			if step.Name == "C/F_Current" {
				var entry1, entry2, entry3, entry4, entry5, entry6, entry7, entry8, entry9, entry10, entry11, entry12 models.JTReportEntry

				entry1.ProductCode = step.ProductCode
				entry1.IFRS17Group = step.IFRS17Group
				populateEntryFields(17, "GMM/VFA", "BEL_Outflow", 2001, "B/S", "Credit", "Liability", 44, &entry1)
				entry1.Debit = math.Max(-step.BelOutflowChange, 0)
				entry1.Credit = math.Max(step.BelOutflowChange, 0)

				entry2.ProductCode = step.ProductCode
				entry2.IFRS17Group = step.IFRS17Group
				populateEntryFields(17, "GMM/VFA", "Insurance_Finance_Income_or_Expense", 420, "PnL", "Credit", "Insurance_Finance_Income_or_(Expense)", 44, &entry2)
				entry2.Debit = math.Max(step.BelOutflowChange, 0)
				entry2.Credit = math.Max(-step.BelOutflowChange, 0)

				entry3.ProductCode = step.ProductCode
				entry3.IFRS17Group = step.IFRS17Group
				populateEntryFields(17, "GMM/VFA", "Risk_Adjustment", 2002, "B/S", "Credit", "Liability", 45, &entry3)
				entry3.Debit = math.Max(-step.RiskAdjustmentChange, 0)
				entry3.Credit = math.Max(step.RiskAdjustmentChange, 0)

				entry4.ProductCode = step.ProductCode
				entry4.IFRS17Group = step.IFRS17Group
				populateEntryFields(17, "GMM/VFA", "Insurance_Finance_Income_or_Expense", 420, "PnL", "Credit", "Insurance_Finance_Income_or_(Expense)", 45, &entry4)
				entry4.Debit = math.Max(step.RiskAdjustmentChange, 0)
				entry4.Credit = math.Max(-step.RiskAdjustmentChange, 0)

				entry5.ProductCode = step.ProductCode
				entry5.IFRS17Group = step.IFRS17Group
				populateEntryFields(17, "GMM/VFA", "BEL_Inflow", 1001, "B/S", "Debit", "Asset", 46, &entry5)
				entry5.Debit = math.Max(step.BelInflowChange, 0)
				entry5.Credit = math.Max(-step.BelInflowChange, 0)

				entry6.ProductCode = step.ProductCode
				entry6.IFRS17Group = step.IFRS17Group
				populateEntryFields(17, "GMM/VFA", "Insurance_Finance_Income_or_Expense", 420, "PnL", "Credit", "Insurance_Finance_Income_or_(Expense)", 46, &entry6)
				entry6.Debit = math.Max(-step.BelInflowChange, 0)
				entry6.Credit = math.Max(step.BelInflowChange, 0)

				//entry4.ProductCode = step.ProductCode
				//entry4.IFRS17Group = step.IFRS17Group
				//populateEntryFields(17, "GMM/VFA", "Insurance_Finance_Income_or_Expense", 420, "PnL", "Credit", "Insurance_Finance_Income_or_(Expense)", 28, &entry4)
				//entry4.Credit = math.Max(-step.LiabilityChange, 0)
				//entry4.Debit = math.Max(step.LiabilityChange, 0)

				entry7.ProductCode = step.ProductCode
				entry7.IFRS17Group = step.IFRS17Group
				populateEntryFields(17, "GMM/VFA", "CSM", 2003, "B/S", "Credit", "Liability", 47, &entry7)
				entry7.Debit = step.CSMRelease

				entry8.ProductCode = step.ProductCode
				entry8.IFRS17Group = step.IFRS17Group
				populateEntryFields(17, "GMM/VFA", "CSM_Release", 106, "PnL", "Credit", "Insurance_Revenue", 47, &entry8)
				entry8.Credit = step.CSMRelease

				entry9.ProductCode = step.ProductCode
				entry9.IFRS17Group = step.IFRS17Group
				populateEntryFields(17, "GMM/VFA", "DAC_Release", 203, "PnL", "Debit", "Insurance_Service_Expense", 48, &entry9)
				entry9.Debit = 0 //?

				entry10.ProductCode = step.ProductCode
				entry10.IFRS17Group = step.IFRS17Group
				populateEntryFields(17, "GMM/VFA", "DAC", 1002, "B/S", "Debit", "Asset", 48, &entry10)
				entry10.Credit = 0 //?

				entry11.ProductCode = step.ProductCode
				entry11.IFRS17Group = step.IFRS17Group
				populateEntryFields(17, "GMM/VFA", "Acquisition_Cashflows", 2007, "B/S", "Credit", "Liability", 49, &entry11)
				entry11.Debit = 0 //?

				entry12.ProductCode = step.ProductCode
				entry12.IFRS17Group = step.IFRS17Group
				populateEntryFields(17, "GMM/VFA", "Amortisation_of_Acquisition_Cashflows", 108, "PnL", "Credit", "Insurance_Revenue", 49, &entry12)
				entry12.Credit = 0 //?

				reportEntries = append(reportEntries, entry1, entry2, entry3, entry4, entry5, entry6, entry7, entry8, entry9, entry10, entry11, entry12)
			}

			// Actuals
			if step.Name == "C/F_Current" {
				var entry1, entry2, entry3, entry4, entry5, entry6, entry7, entry8, entry9, entry10, entry11, entry12, entry13, entry14, entry15, entry16 models.JTReportEntry
				var financeVariables = getFinancialVariables(step.ProductCode, step.IFRS17Group, csmRun.FinanceYear, csmRun.FinanceVersion)

				entry1.ProductCode = step.ProductCode
				entry1.IFRS17Group = step.IFRS17Group
				populateEntryFields(100, "GMM/VFA", "Actual_Mortality_Claims", 205, "PnL", "Debit", "Insurance_Service_Expense", 50, &entry1)
				entry1.Debit = financeVariables.ActualMortalityClaimsIncurred

				entry2.ProductCode = step.ProductCode
				entry2.IFRS17Group = step.IFRS17Group
				populateEntryFields(100, "GMM/VFA", "OCR", 2012, "PnL", "Debit", "Asset", 50, &entry2)
				entry2.Credit = financeVariables.ActualMortalityClaimsIncurred

				entry3.ProductCode = step.ProductCode
				entry3.IFRS17Group = step.IFRS17Group
				populateEntryFields(101, "GMM/VFA", "Actual_Expenses", 206, "PnL", "Debit", "Insurance_Service_Expense", 51, &entry3)
				entry3.Debit = financeVariables.ActualAttributableExpenses

				entry4.ProductCode = step.ProductCode
				entry4.IFRS17Group = step.IFRS17Group
				populateEntryFields(101, "GMM/VFA", "Cash_and_Cash_Equivalent", 1005, "PnL", "Debit", "Asset", 51, &entry4)
				entry4.Credit = financeVariables.ActualAttributableExpenses

				entry5.ProductCode = step.ProductCode
				entry5.IFRS17Group = step.IFRS17Group
				populateEntryFields(102, "GMM/VFA", "Actual_Retrenchment_Claims", 207, "B/S", "Debit", "Insurance_Service_Expense", 52, &entry5)
				entry5.Debit = financeVariables.ActualRetrenchmentClaimsIncurred

				entry6.ProductCode = step.ProductCode
				entry6.IFRS17Group = step.IFRS17Group
				populateEntryFields(102, "GMM/VFA", "Cash_and_Cash_Equivalent", 1005, "PnL", "Debit", "Asset", 52, &entry6)
				entry6.Credit = financeVariables.ActualRetrenchmentClaimsIncurred

				entry7.ProductCode = step.ProductCode
				entry7.IFRS17Group = step.IFRS17Group
				populateEntryFields(103, "GMM/VFA", "Actual_Morbidity_Claims", 208, "PnL", "Debit", "Insurance_Service_Expense", 53, &entry7)
				entry7.Debit = financeVariables.ActualMorbidityClaimsIncurred

				entry8.ProductCode = step.ProductCode
				entry8.IFRS17Group = step.IFRS17Group
				populateEntryFields(103, "GMM/VFA", "Cash_and_Cash_Equivalent", 1005, "B/S", "Debit", "Asset", 53, &entry8)
				entry8.Credit = financeVariables.ActualMorbidityClaimsIncurred

				entry9.ProductCode = step.ProductCode
				entry9.IFRS17Group = step.IFRS17Group
				populateEntryFields(104, "GMM/VFA", "BEL_Outflow", 2001, "B/S", "Credit", "Liability", 54, &entry9)
				//entry9.Debit = tempBelOutflowBuildup
				entry9.Debit = math.Max(-tempBelOutflowBuildup, 0)
				entry9.Credit = math.Max(tempBelOutflowBuildup, 0)

				entry10.ProductCode = step.ProductCode
				entry10.IFRS17Group = step.IFRS17Group
				populateEntryFields(104, "GMM/VFA", "BEL_Outflow_Buildup", 2008, "PnL", "Credit", "Liability", 54, &entry10)
				//entry10.Credit = tempBelOutflowBuildup
				entry10.Debit = math.Max(tempBelOutflowBuildup, 0)
				entry10.Credit = math.Max(-tempBelOutflowBuildup, 0)

				entry11.ProductCode = step.ProductCode
				entry11.IFRS17Group = step.IFRS17Group
				populateEntryFields(105, "GMM/VFA", "Risk_Adjustment", 2002, "B/S", "Credit", "Liability", 55, &entry11)
				//entry11.Debit = tempRiskAdjustmentBuildup
				entry11.Debit = math.Max(-tempRiskAdjustmentBuildup, 0)
				entry11.Credit = math.Max(tempRiskAdjustmentBuildup, 0)

				entry12.ProductCode = step.ProductCode
				entry12.IFRS17Group = step.IFRS17Group
				populateEntryFields(105, "GMM/VFA", "Risk_Adjustment_Buildup", 2009, "PnL", "Credit", "Liability", 55, &entry12)
				//entry12.Credit = tempRiskAdjustmentBuildup
				entry12.Debit = math.Max(tempRiskAdjustmentBuildup, 0)
				entry12.Credit = math.Max(-tempRiskAdjustmentBuildup, 0)

				entry13.ProductCode = step.ProductCode
				entry13.IFRS17Group = step.IFRS17Group
				populateEntryFields(106, "GMM/VFA", "BEL_Inflow_Buildup", 1004, "B/S", "Debit", "Asset", 56, &entry13)
				//entry13.Debit = tempBelInflowBuildup
				entry13.Debit = math.Max(-tempBelInflowBuildup, 0)
				entry13.Credit = math.Max(tempBelInflowBuildup, 0)

				entry14.ProductCode = step.ProductCode
				entry14.IFRS17Group = step.IFRS17Group
				populateEntryFields(106, "GMM/VFA", "BEL_Inflow", 1001, "PnL", "Debit", "Asset", 56, &entry14)
				//entry14.Credit = tempBelInflowBuildup
				entry14.Debit = math.Max(tempBelInflowBuildup, 0)
				entry14.Credit = math.Max(-tempBelInflowBuildup, 0)

				entry15.ProductCode = step.ProductCode
				entry15.IFRS17Group = step.IFRS17Group
				populateEntryFields(107, "GMM/VFA", "Cash_and_Cash_Equivalent", 1005, "B/S", "Debit", "Asset", 57, &entry15)
				//entry15.Debit = tempInvestmentFinanceExpense
				entry15.Debit = 0  // math.Max(tempInvestmentFinanceExpense, 0) already allowed for in individual transactions above
				entry15.Credit = 0 // math.Max(-tempInvestmentFinanceExpense, 0)

				entry16.ProductCode = step.ProductCode
				entry16.IFRS17Group = step.IFRS17Group
				populateEntryFields(107, "GMM/VFA", "Insurance_Finance_Income_or_Expense", 420, "PnL", "Credit", "Insurance_Finance_Income_or_(Expense)", 57, &entry16)
				//entry16.Credit = tempInvestmentFinanceExpense
				entry16.Debit = 0  //math.Max(-tempInvestmentFinanceExpense, 0)
				entry16.Credit = 0 // math.Max(tempInvestmentFinanceExpense, 0)

				reportEntries = append(reportEntries, entry1, entry2, entry3, entry4, entry5, entry6, entry7, entry8, entry9, entry10, entry11, entry12, entry13, entry14, entry15, entry16)
			}
		}
	} else {
		//populateEntryFields(2, "GMM", "BEL_Inflow", 1001, "B/S", "Debit", "Asset", 1, &entry1)
	}

	if len(paaresults) > 0 {
		var entry1, entry2, entry3, entry4, entry5, entry6, entry7, entry8, entry9, entry10, entry11, entry12, entry13, entry14, entry15, entry16, entry17, entry18, entry19, entry20, entry21, entry22, entry23, entry24, entry25, entry26, entry27, entry28, entry29, entry30, entry31, entry32, entry33, entry34, entry35 models.JTReportEntry
		for _, groupresult := range paaresults {
			entry1.ProductCode = groupresult.ProductCode
			entry1.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "Cash_and_Cash_Equivalent", 1005, "B/S", "Debit", "Asset", 58, &entry1)
			entry1.Debit = groupresult.TotalPremiumReceipt

			entry2.ProductCode = groupresult.ProductCode
			entry2.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "Premium_Debtors", 1006, "B/S", "Debit", "Asset", 58, &entry2)
			entry2.Debit = 0

			entry3.ProductCode = groupresult.ProductCode
			entry3.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "PAA_Unearned_Premium", 2010, "B/S", "Credit", "Liability", 58, &entry3)
			entry3.Credit = groupresult.TotalPremiumReceipt - groupresult.AcquisitionExpenses

			entry4.ProductCode = groupresult.ProductCode
			entry4.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "Acquisition_Cashflows", 2007, "B/S", "Credit", "Liability", 58, &entry4)
			entry4.Credit = groupresult.AcquisitionExpenses

			entry5.ProductCode = groupresult.ProductCode
			entry5.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "PAA_Unearned_Premium", 2010, "B/S", "Credit", "Liability", 59, &entry5)
			entry5.Debit = groupresult.EarnedPremium

			entry6.ProductCode = groupresult.ProductCode
			entry6.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "PAA_Earned_Premium", 109, "PnL", "Credit", "Insurance_Revenue", 59, &entry6)
			entry6.Credit = groupresult.EarnedPremium

			entry7.ProductCode = groupresult.ProductCode
			entry7.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "Acquisition_Cashflows", 2007, "B/S", "Credit", "Liability", 60, &entry7)
			entry7.Debit = groupresult.AmortisedAcquisitionCost

			entry8.ProductCode = groupresult.ProductCode
			entry8.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "amortised_acquisition_cost", 108, "PnL", "Credit", "Insurance_Revenue", 60, &entry8)
			entry8.Credit = groupresult.AmortisedAcquisitionCost

			entry9.ProductCode = groupresult.ProductCode
			entry9.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "PAA_Loss_Component", 209, "PnL", "Debit", "Insurance_Service_Expense", 61, &entry9)
			entry9.Debit = groupresult.PaaLossComponent

			entry10.ProductCode = groupresult.ProductCode
			entry10.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "Loss_Component_Reserve", 2011, "B/S", "Credit", "Liability", 61, &entry10)
			entry10.Credit = groupresult.PaaLossComponent

			entry11.ProductCode = groupresult.ProductCode
			entry11.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "PAA_Incurred_Claims", 210, "PnL", "Debit", "Insurance_Service_Expense", 62, &entry11)
			entry11.Debit = groupresult.IncurredClaims

			entry12.ProductCode = groupresult.ProductCode
			entry12.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "OCR", 2012, "B/S", "Credit", "Liability", 62, &entry12)
			entry12.Credit = groupresult.IncurredClaims

			entry13.ProductCode = groupresult.ProductCode
			entry13.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "OCR", 2012, "B/S", "Credit", "Liability", 63, &entry13)
			entry13.Debit = groupresult.ClaimsPaid

			entry14.ProductCode = groupresult.ProductCode
			entry14.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "Cash_and_Cash_Equivalent", 1005, "B/S", "Debit", "Asset", 63, &entry14)
			entry14.Credit = groupresult.ClaimsPaid

			entry15.ProductCode = groupresult.ProductCode
			entry15.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "PAA_Loss_Component_Adjustment", 211, "PnL", "Debit", "Insurance_Service_Expense", 64, &entry15)
			entry15.Credit = 0

			entry16.ProductCode = groupresult.ProductCode
			entry16.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "Loss_Component_Reserve", 2011, "B/S", "Credit", "Liability", 64, &entry16)
			entry16.Credit = 0

			entry17.ProductCode = groupresult.ProductCode
			entry17.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "PAA_Reinsurance_LRC", 31001, "B/S", "Debit", "Asset", 65, &entry17)
			entry17.Credit = groupresult.PaaReinsuranceDac

			entry18.ProductCode = groupresult.ProductCode
			entry18.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "PAA_Reinsurance_Deferred_Ceding_Commission", 32001, "B/S", "Credit", "Liability", 65, &entry18)
			entry18.Credit = groupresult.PaaReinsuranceDac

			entry19.ProductCode = groupresult.ProductCode
			entry19.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "Cash_and_Cash_Equivalent", 1005, "B/S", "Debit", "Asset", 65, &entry19)
			entry19.Credit = groupresult.ReinsurancePremium

			entry20.ProductCode = groupresult.ProductCode
			entry20.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "PAA_Allocated_Reinsurance_Premium", 3201, "PnL", "Debit", "Insurance_Service_Expense", 66, &entry20)
			entry20.Debit = groupresult.AllocatedReinsurancePremium

			entry21.ProductCode = groupresult.ProductCode
			entry21.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "PAA_Reinsurance_LRC", 31001, "B/S", "Debit", "Asset", 66, &entry21)
			entry21.Credit = groupresult.AllocatedReinsurancePremium

			entry22.ProductCode = groupresult.ProductCode
			entry22.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "PAA_Reinsurance_deferred_Ceding_Commission", 32001, "B/S", "Credit", "Liability", 67, &entry22)
			entry22.Debit = groupresult.AllocatedReinsuranceFlatCommission

			entry23.ProductCode = groupresult.ProductCode
			entry23.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "PAA_Allocated_Ceding_Commission", 3101, "PnL", "Credit", "Insurance_Revenue", 67, &entry23)
			entry23.Credit = groupresult.AllocatedReinsuranceFlatCommission

			entry24.ProductCode = groupresult.ProductCode
			entry24.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "PAA_Loss_Recovery_Asset", 31002, "B/S", "Debit", "Asset", 68, &entry24)
			entry24.Debit = groupresult.PaaInitialLossRecoveryComponent

			entry25.ProductCode = groupresult.ProductCode
			entry25.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "Initial_Recognition_Loss_Recovery_Income", 3102, "PnL", "Credit", "Insurance_Revenue", 68, &entry25)
			entry25.Credit = groupresult.PaaInitialLossRecoveryComponent

			entry26.ProductCode = groupresult.ProductCode
			entry26.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "PAA_Loss_Recovery_Asset", 31002, "B/S", "Debit", "Asset", 69, &entry26)
			entry26.Debit = groupresult.PaaLossRecoveryComponent

			entry27.ProductCode = groupresult.ProductCode
			entry27.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "Loss_Recovery_Component_Income", 3103, "PnL", "Credit", "Insurance_Revenue", 69, &entry27)
			entry27.Credit = groupresult.PaaLossRecoveryComponent

			entry28.ProductCode = groupresult.ProductCode
			entry28.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "Cash_and_Cash_Equivalent", 1005, "B/S", "Debit", "Asset", 70, &entry28)
			entry28.Credit = 0

			entry29.ProductCode = groupresult.ProductCode
			entry29.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "Reinsurance_Incurred_Claims_Adjustment", 3105, "PnL", "Credit", "Insurance_Revenue", 70, &entry29)
			entry29.Credit = 0

			entry30.ProductCode = groupresult.ProductCode
			entry30.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "Reversal_Loss_Recovery_Component", 3202, "PnL", "Credit", "Insurance_Service_Expense", 71, &entry30)
			entry30.Debit = 0

			entry31.ProductCode = groupresult.ProductCode
			entry31.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "PAA_Loss_Recovery_Asset", 31002, "PnL", "Debit", "Asset", 71, &entry31)
			entry31.Credit = 0

			entry32.ProductCode = groupresult.ProductCode
			entry32.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "Non_Attributable_Expenses", 212, "PnL", "Debit", "Insurance_Service_Expense", 72, &entry32)
			entry32.Debit = groupresult.NonAttributableExpenses

			entry33.ProductCode = groupresult.ProductCode
			entry33.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "Cash_and_Cash_Equivalent", 1005, "PnL", "Debit", "Asset", 72, &entry33)
			entry33.Credit = groupresult.NonAttributableExpenses

			entry34.ProductCode = groupresult.ProductCode
			entry34.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "Actual_Expenses", 206, "PnL", "Debit", "Insurance_Service_Expense", 73, &entry34)
			entry34.Debit = groupresult.IncurredExpenses

			entry35.ProductCode = groupresult.ProductCode
			entry35.IFRS17Group = groupresult.IFRS17Group
			populateEntryFields(201, "PAA", "Cash_and_Cash_Equivalent", 1005, "PnL", "Debit", "Asset", 73, &entry35)
			entry35.Credit = groupresult.IncurredExpenses

			reportEntries = append(reportEntries, entry1, entry2, entry3, entry4, entry5, entry6, entry7, entry8, entry9, entry10, entry11, entry12, entry13, entry14, entry15, entry16, entry17, entry18, entry19, entry20, entry21, entry22, entry23, entry24, entry25, entry26, entry27, entry28, entry29, entry30, entry31, entry32, entry33, entry34, entry35)
		}
	}

	if len(licBuildupResults) > 0 {
		var entry1, entry2, entry3, entry4, entry5, entry6, entry7, entry8 models.JTReportEntry

		entry1.ProductCode = licBuildupResults[8].ProductCode
		entry1.IFRS17Group = licBuildupResults[8].ProductCode
		populateEntryFields(201, "LIC", "LIC_FCF_Changes", 215, "PnL", "Debit", "Insurance_Service_Expense", 73, &entry1)
		entry1.Debit = licBuildupResults[3].IBNR + licBuildupResults[4].IBNR + licBuildupResults[5].IBNR + licBuildupResults[6].IBNR

		entry2.ProductCode = licBuildupResults[8].ProductCode
		entry2.IFRS17Group = licBuildupResults[8].ProductCode
		populateEntryFields(201, "LIC", "IBNR_Reserve", 42001, "B/S", "Credit", "Liability", 73, &entry2)
		entry2.Credit = entry1.Debit

		entry3.ProductCode = licBuildupResults[8].ProductCode
		entry3.IFRS17Group = licBuildupResults[8].ProductCode
		populateEntryFields(201, "LIC", "LIC_FCF_Changes", 215, "PnL", "Debit", "Insurance_Service_Expense", 74, &entry3)
		entry3.Debit = licBuildupResults[3].RiskAdjustment + licBuildupResults[4].RiskAdjustment + licBuildupResults[5].RiskAdjustment + licBuildupResults[6].RiskAdjustment

		entry4.ProductCode = licBuildupResults[8].ProductCode
		entry4.IFRS17Group = licBuildupResults[8].ProductCode
		populateEntryFields(201, "LIC", "LIC_Risk_Adjustment", 42002, "B/S", "Credit", "Liability", 74, &entry4)
		entry4.Credit = entry3.Debit

		entry5.ProductCode = licBuildupResults[8].ProductCode
		entry5.IFRS17Group = licBuildupResults[8].ProductCode
		populateEntryFields(201, "LIC", "LIC_Experience_Adjustments", 216, "PnL", "Debit", "Insurance_Service_Expense", 75, &entry5)
		entry5.Debit = licBuildupResults[2].IBNR + licBuildupResults[7].IBNR + licBuildupResults[8].IBNR + licBuildupResults[9].IBNR + licBuildupResults[10].IBNR

		entry6.ProductCode = licBuildupResults[8].ProductCode
		entry6.IFRS17Group = licBuildupResults[8].ProductCode
		populateEntryFields(201, "LIC", "IBNR_Reserve", 42001, "B/S", "Credit", "Liability", 75, &entry6)
		entry6.Credit = entry5.Debit

		entry7.ProductCode = licBuildupResults[8].ProductCode
		entry7.IFRS17Group = licBuildupResults[8].ProductCode
		populateEntryFields(201, "LIC", "LIC_Experience_Adjustments", 216, "PnL", "Debit", "Insurance_Service_Expense", 76, &entry7)
		entry7.Debit = licBuildupResults[2].RiskAdjustment + licBuildupResults[7].RiskAdjustment + licBuildupResults[8].RiskAdjustment + licBuildupResults[9].RiskAdjustment + licBuildupResults[10].RiskAdjustment

		entry8.ProductCode = licBuildupResults[8].ProductCode
		entry8.IFRS17Group = licBuildupResults[8].ProductCode
		populateEntryFields(201, "LIC", "LIC_Risk_Adjustment", 42002, "B/S", "Credit", "Liability", 76, &entry8)
		entry8.Credit = entry7.Debit

		reportEntries = append(reportEntries, entry1, entry2, entry3, entry4, entry5, entry6, entry7, entry8)
	}

	sort.Slice(reportEntries, func(i, j int) bool {
		var sortedByReportBundle bool
		var sortedByAccountNumber bool
		sortedByReportBundle = reportEntries[i].ReportBundle < reportEntries[j].ReportBundle

		if reportEntries[i].ReportBundle == reportEntries[j].ReportBundle {
			sortedByAccountNumber = reportEntries[i].AccountNumber < reportEntries[j].AccountNumber
			return sortedByAccountNumber
		}
		return sortedByReportBundle
	})
	return reportEntries
}

func GenerateSubLedgerReportByGroup(journalTransactions []models.JTReportEntry, results []models.AOSStepResult, paaresults []models.PAAResult) []models.SubLedgerReportEntry {
	// The assumption we will work with is that the stepResults have been
	// filtered by product_code and ifrs17_group

	var subledgerEntries []models.SubLedgerReportEntry
	//For now we calculate each entry manually...

	if len(results) > 0 {

		var entry100, entry101, entry102, entry103, entry104, entry105, entry106, entry107, entry108, entry109, entry100999 models.SubLedgerReportEntry //1

		entry100.ProductCode = journalTransactions[1].ProductCode
		entry100.IFRS17Group = journalTransactions[1].IFRS17Group
		entry100.LedgerName = "Bel_Inflow"
		entry100.MasterAccountType = journalTransactions[0].MasterAccountType
		entry100.AccountNumber = 1001
		entry100.PostingDate = ""
		entry100.ContraAccountTransactionDescription = "Bel_Inflow_B/F"
		entry100.PostReference = "J" + strconv.Itoa(journalTransactions[0].AosStep)
		entry100.Debit = results[0].BelInflow
		entry100.Credit = 0
		if journalTransactions[0].NormalAccountBalance == "Debit" {
			entry100.Balance = -(entry100.Credit - entry100.Debit)
		} else {
			entry100.Balance = entry100.Credit - entry100.Debit
		}

		entry101.ProductCode = journalTransactions[1].ProductCode
		entry101.IFRS17Group = journalTransactions[1].IFRS17Group
		entry101.LedgerName = "Bel_Inflow"
		entry101.MasterAccountType = journalTransactions[12].MasterAccountType
		entry101.AccountNumber = 1001
		entry101.PostingDate = ""
		entry101.ContraAccountTransactionDescription = "Insurance_Finance_Income_or_(Expense)"
		entry101.PostReference = "J" + strconv.Itoa(journalTransactions[12].AosStep)
		entry101.Debit = journalTransactions[12].Debit
		entry101.Credit = journalTransactions[12].Credit
		if journalTransactions[12].NormalAccountBalance == "Debit" {
			entry101.Balance = -(entry101.Credit - entry101.Debit) + entry100.Balance
		} else {
			entry101.Balance = entry101.Credit - entry101.Debit + entry100.Balance
		}

		entry102.ProductCode = journalTransactions[1].ProductCode
		entry102.IFRS17Group = journalTransactions[1].IFRS17Group
		entry102.LedgerName = "Bel_Inflow"
		entry102.MasterAccountType = journalTransactions[13].MasterAccountType
		entry102.AccountNumber = 1001
		entry102.PostingDate = ""
		entry102.ContraAccountTransactionDescription = "Bel_Outflow"
		entry102.PostReference = "J" + strconv.Itoa(journalTransactions[13].AosStep)
		entry102.Debit = journalTransactions[14].Credit //80
		entry102.Credit = journalTransactions[14].Debit
		if journalTransactions[13].NormalAccountBalance == "Debit" {
			entry102.Balance = -(entry102.Credit - entry102.Debit) + entry101.Balance
		} else {
			entry102.Balance = entry102.Credit - entry102.Debit + entry101.Balance
		}

		entry103.ProductCode = journalTransactions[1].ProductCode
		entry103.IFRS17Group = journalTransactions[1].IFRS17Group
		entry103.LedgerName = "Bel_Inflow"
		entry103.MasterAccountType = journalTransactions[13].MasterAccountType
		entry103.AccountNumber = 1001
		entry103.PostingDate = ""
		entry103.ContraAccountTransactionDescription = "CSM"
		entry103.PostReference = "J" + strconv.Itoa(journalTransactions[13].AosStep)
		entry103.Debit = journalTransactions[16].Credit
		entry103.Credit = journalTransactions[16].Debit
		if journalTransactions[13].NormalAccountBalance == "Debit" {
			entry103.Balance = -(entry103.Credit - entry103.Debit) + entry102.Balance
		} else {
			entry103.Balance = entry103.Credit - entry103.Debit + entry102.Balance
		}

		entry104.ProductCode = journalTransactions[1].ProductCode
		entry104.IFRS17Group = journalTransactions[1].IFRS17Group
		entry104.LedgerName = "Bel_Inflow"
		entry104.MasterAccountType = journalTransactions[13].MasterAccountType
		entry104.AccountNumber = 1001
		entry104.PostingDate = ""
		entry104.ContraAccountTransactionDescription = "Risk_Adjustment"
		entry104.PostReference = "J" + strconv.Itoa(journalTransactions[13].AosStep)
		entry104.Debit = journalTransactions[15].Credit //10
		entry104.Credit = journalTransactions[15].Debit
		if journalTransactions[13].NormalAccountBalance == "Debit" {
			entry104.Balance = -(entry104.Credit - entry104.Debit) + entry103.Balance
		} else {
			entry104.Balance = entry104.Credit - entry104.Debit + entry103.Balance
		}

		entry105.ProductCode = journalTransactions[1].ProductCode
		entry105.IFRS17Group = journalTransactions[1].IFRS17Group
		entry105.LedgerName = "Bel_Inflow"
		entry105.MasterAccountType = journalTransactions[50].MasterAccountType
		entry105.AccountNumber = 1001
		entry105.PostingDate = ""
		entry105.ContraAccountTransactionDescription = "Premium_Debtors"
		entry105.PostReference = "J" + strconv.Itoa(journalTransactions[50].AosStep)
		entry105.Debit = journalTransactions[50].Debit
		entry105.Credit = journalTransactions[50].Credit
		if journalTransactions[50].NormalAccountBalance == "Debit" {
			entry105.Balance = -(entry105.Credit - entry105.Debit) + entry104.Balance
		} else {
			entry105.Balance = entry105.Credit - entry105.Debit + entry104.Balance
		}

		entry106.ProductCode = journalTransactions[1].ProductCode
		entry106.IFRS17Group = journalTransactions[1].IFRS17Group
		entry106.LedgerName = "Bel_Inflow"
		entry106.MasterAccountType = journalTransactions[45].MasterAccountType
		entry106.AccountNumber = 1001
		entry106.PostingDate = ""
		entry106.ContraAccountTransactionDescription = "Insurance_Finance_Income_or_(Expense)"
		entry106.PostReference = "J" + strconv.Itoa(journalTransactions[45].AosStep)
		entry106.Debit = journalTransactions[45].Debit
		entry106.Credit = journalTransactions[45].Credit
		if journalTransactions[45].NormalAccountBalance == "Debit" {
			entry106.Balance = -(entry106.Credit - entry106.Debit) + entry105.Balance
		} else {
			entry106.Balance = entry106.Credit - entry106.Debit + entry105.Balance
		}

		entry107.ProductCode = journalTransactions[1].ProductCode
		entry107.IFRS17Group = journalTransactions[1].IFRS17Group
		entry107.LedgerName = "Bel_Inflow"
		entry107.MasterAccountType = journalTransactions[125].MasterAccountType
		entry107.AccountNumber = 1001
		entry107.PostingDate = ""
		entry107.ContraAccountTransactionDescription = "Bel_Inflow_Buildup"
		entry107.PostReference = "J" + strconv.Itoa(journalTransactions[125].AosStep)
		entry107.Debit = journalTransactions[125].Debit
		entry107.Credit = journalTransactions[125].Credit
		if journalTransactions[125].NormalAccountBalance == "Debit" {
			entry107.Balance = -(entry107.Credit - entry107.Debit) + entry106.Balance
		} else {
			entry107.Balance = entry107.Credit - entry107.Debit + entry106.Balance
		}

		entry108.ProductCode = journalTransactions[1].ProductCode
		entry108.IFRS17Group = journalTransactions[1].IFRS17Group
		entry108.LedgerName = "Bel_Inflow"
		entry108.MasterAccountType = journalTransactions[100].MasterAccountType
		entry108.AccountNumber = 1001
		entry108.PostingDate = ""
		entry108.ContraAccountTransactionDescription = "Insurance_Finance_Income_or_Expense"
		entry108.PostReference = "J" + strconv.Itoa(journalTransactions[100].AosStep)
		entry108.Debit = journalTransactions[100].Debit
		entry108.Credit = journalTransactions[100].Credit
		if journalTransactions[100].NormalAccountBalance == "Debit" {
			entry108.Balance = -(entry108.Credit - entry108.Debit) + entry107.Balance
		} else {
			entry108.Balance = entry108.Credit - entry108.Debit + entry107.Balance
		}

		entry109.ProductCode = journalTransactions[1].ProductCode
		entry109.IFRS17Group = journalTransactions[1].IFRS17Group
		entry109.LedgerName = "Bel_Inflow"
		entry109.MasterAccountType = journalTransactions[106].MasterAccountType
		entry109.AccountNumber = 1001
		entry109.PostingDate = ""
		entry109.ContraAccountTransactionDescription = "Insurance_Finance_Income_or_Expense"
		entry109.PostReference = "J" + strconv.Itoa(journalTransactions[106].AosStep)
		entry109.Debit = journalTransactions[106].Debit
		entry109.Credit = journalTransactions[106].Credit
		if journalTransactions[106].NormalAccountBalance == "Debit" {
			entry109.Balance = -(entry109.Credit - entry109.Debit) + entry108.Balance
		} else {
			entry109.Balance = entry109.Credit - entry109.Debit + entry108.Balance
		}

		entry100999.ProductCode = journalTransactions[1].ProductCode
		entry100999.IFRS17Group = journalTransactions[1].IFRS17Group
		entry100999.LedgerName = "Bel_Inflow_B/F"
		entry100999.MasterAccountType = journalTransactions[0].MasterAccountType
		entry100999.AccountNumber = 1001999
		entry100999.PostingDate = ""
		entry100999.ContraAccountTransactionDescription = "Bel_Inflow"
		entry100999.PostReference = "J" + strconv.Itoa(journalTransactions[0].AosStep)
		entry100999.Debit = 0
		entry100999.Credit = results[0].BelInflow
		if journalTransactions[0].NormalAccountBalance == "Debit" {
			entry100999.Balance = -(entry100999.Credit - entry100999.Debit)
		} else {
			entry100999.Balance = entry100999.Credit - entry100999.Debit
		}

		subledgerEntries = append(subledgerEntries, entry100, entry101, entry102, entry103, entry104, entry105, entry106, entry107, entry108, entry109, entry100999)

		// Bel_Outflow
		var entry201999, entry201, entry202, entry203, entry204, entry205, entry206, entry207, entry208, entry209, entry210, entry211, entry212, entry213 models.SubLedgerReportEntry //2

		//...999 opening balance offset account
		entry201999.ProductCode = journalTransactions[1].ProductCode
		entry201999.IFRS17Group = journalTransactions[1].IFRS17Group
		entry201999.LedgerName = "Bel_Outflow_B/F"
		entry201999.MasterAccountType = journalTransactions[2].MasterAccountType
		entry201999.AccountNumber = 2001999
		entry201999.PostingDate = ""
		entry201999.ContraAccountTransactionDescription = "Bel_Outflow"
		entry201999.PostReference = "J" + strconv.Itoa(journalTransactions[2].AosStep)
		entry201999.Debit = results[0].BelOutflow
		entry201999.Credit = 0
		if journalTransactions[2].NormalAccountBalance == "Debit" {
			entry201999.Balance = -(entry201999.Credit - entry201999.Debit)
		} else {
			entry201999.Balance = entry201999.Credit - entry201999.Debit
		}

		entry201.ProductCode = journalTransactions[1].ProductCode
		entry201.IFRS17Group = journalTransactions[1].IFRS17Group
		entry201.LedgerName = "Bel_Outflow"
		entry201.MasterAccountType = journalTransactions[2].MasterAccountType
		entry201.AccountNumber = 2001
		entry201.PostingDate = ""
		entry201.ContraAccountTransactionDescription = "Bel_Outflow_B/F"
		entry201.PostReference = "J" + strconv.Itoa(journalTransactions[2].AosStep)
		entry201.Debit = 0                      //journalTransactions[2].Debit
		entry201.Credit = results[0].BelOutflow //journalTransactions[2].Credit
		if journalTransactions[2].NormalAccountBalance == "Debit" {
			entry201.Balance = -(entry201.Credit - entry201.Debit)
		} else {
			entry201.Balance = entry201.Credit - entry201.Debit
		}

		entry202.ProductCode = journalTransactions[1].ProductCode
		entry202.IFRS17Group = journalTransactions[1].IFRS17Group
		entry202.LedgerName = "Bel_Outflow"
		entry202.MasterAccountType = journalTransactions[8].MasterAccountType
		entry202.AccountNumber = 2001
		entry202.PostingDate = ""
		entry202.ContraAccountTransactionDescription = "Insurance_Finance_Income_or_Expense"
		entry202.PostReference = "J" + strconv.Itoa(journalTransactions[8].AosStep)
		entry202.Debit = journalTransactions[8].Debit //7
		entry202.Credit = journalTransactions[8].Credit
		if journalTransactions[8].NormalAccountBalance == "Debit" {
			entry202.Balance = -(entry202.Credit - entry202.Debit) + entry201.Balance
		} else {
			entry202.Balance = entry202.Credit - entry202.Debit + entry201.Balance
		}

		entry203.ProductCode = journalTransactions[1].ProductCode
		entry203.IFRS17Group = journalTransactions[1].IFRS17Group
		entry203.LedgerName = "Bel_Outflow"
		entry203.MasterAccountType = journalTransactions[14].MasterAccountType
		entry203.AccountNumber = 2001
		entry203.PostingDate = ""
		entry203.ContraAccountTransactionDescription = "Bel_Inflow" //Loss_Component_NB
		entry203.PostReference = "J" + strconv.Itoa(journalTransactions[14].AosStep)
		entry203.Debit = journalTransactions[14].Debit //***
		entry203.Credit = journalTransactions[14].Credit
		if journalTransactions[14].NormalAccountBalance == "Debit" {
			entry203.Balance = -(entry203.Credit - entry203.Debit) + entry202.Balance
		} else {
			entry203.Balance = entry203.Credit - entry203.Debit + entry202.Balance
		}

		entry204.ProductCode = journalTransactions[1].ProductCode
		entry204.IFRS17Group = journalTransactions[1].IFRS17Group
		entry204.LedgerName = "Bel_Outflow"
		entry204.MasterAccountType = journalTransactions[14].MasterAccountType
		entry204.AccountNumber = 2001
		entry204.PostingDate = ""
		entry204.ContraAccountTransactionDescription = "Loss_Component_NB" //Loss_Component_NB
		entry204.PostReference = "J" + strconv.Itoa(journalTransactions[18].AosStep)
		entry204.Debit = journalTransactions[18].Debit //***
		entry204.Credit = journalTransactions[18].Credit
		if journalTransactions[18].NormalAccountBalance == "Debit" {
			entry204.Balance = -(entry204.Credit - entry204.Debit) + entry203.Balance
		} else {
			entry204.Balance = entry204.Credit - entry204.Debit + entry203.Balance
		}

		entry205.ProductCode = journalTransactions[1].ProductCode
		entry205.IFRS17Group = journalTransactions[1].IFRS17Group
		entry205.LedgerName = "Bel_Outflow"
		entry205.MasterAccountType = journalTransactions[27].MasterAccountType
		entry205.AccountNumber = 2001
		entry205.PostingDate = ""
		entry205.ContraAccountTransactionDescription = "Expected_Mortality_Claims"
		entry205.PostReference = "J" + strconv.Itoa(journalTransactions[27].AosStep)
		entry205.Debit = journalTransactions[27].Debit
		entry205.Credit = journalTransactions[27].Credit //***
		if journalTransactions[27].NormalAccountBalance == "Debit" {
			entry205.Balance = -(entry205.Credit - entry205.Debit) + entry204.Balance
		} else {
			entry205.Balance = entry205.Credit - entry205.Debit + entry204.Balance
		}

		entry206.ProductCode = journalTransactions[1].ProductCode
		entry206.IFRS17Group = journalTransactions[1].IFRS17Group
		entry206.LedgerName = "Bel_Outflow"
		entry206.MasterAccountType = journalTransactions[31].MasterAccountType
		entry206.AccountNumber = 2001
		entry206.PostingDate = ""
		entry206.ContraAccountTransactionDescription = "Expected_Expenses"
		entry206.PostReference = "J" + strconv.Itoa(journalTransactions[31].AosStep)
		entry206.Debit = journalTransactions[31].Debit //***
		entry206.Credit = journalTransactions[31].Credit
		if journalTransactions[31].NormalAccountBalance == "Debit" {
			entry206.Balance = -(entry206.Credit - entry206.Debit) + entry205.Balance
		} else {
			entry206.Balance = entry206.Credit - entry206.Debit + entry205.Balance
		}

		entry207.ProductCode = journalTransactions[1].ProductCode
		entry207.IFRS17Group = journalTransactions[1].IFRS17Group
		entry207.LedgerName = "Bel_Outflow"
		entry207.MasterAccountType = journalTransactions[35].MasterAccountType
		entry207.AccountNumber = 2001
		entry207.PostingDate = ""
		entry207.ContraAccountTransactionDescription = "Expected_Retrenchment"
		entry207.PostReference = "J" + strconv.Itoa(journalTransactions[35].AosStep)
		entry207.Debit = journalTransactions[35].Debit //***
		entry207.Credit = journalTransactions[35].Credit
		if journalTransactions[35].NormalAccountBalance == "Debit" {
			entry207.Balance = -(entry207.Credit - entry207.Debit) + entry206.Balance
		} else {
			entry207.Balance = entry207.Credit - entry207.Debit + entry206.Balance
		}

		entry208.ProductCode = journalTransactions[1].ProductCode
		entry208.IFRS17Group = journalTransactions[1].IFRS17Group
		entry208.LedgerName = "Bel_Outflow"
		entry208.MasterAccountType = journalTransactions[39].MasterAccountType
		entry208.AccountNumber = 2001
		entry208.PostingDate = ""
		entry208.ContraAccountTransactionDescription = "Expected_Morbidity"
		entry208.PostReference = "J" + strconv.Itoa(journalTransactions[39].AosStep)
		entry208.Debit = journalTransactions[39].Debit //***
		entry208.Credit = journalTransactions[39].Credit
		if journalTransactions[39].NormalAccountBalance == "Debit" {
			entry208.Balance = -(entry208.Credit - entry208.Debit) + entry207.Balance
		} else {
			entry208.Balance = entry208.Credit - entry208.Debit + entry207.Balance
		}

		entry209.ProductCode = journalTransactions[1].ProductCode
		entry209.IFRS17Group = journalTransactions[1].IFRS17Group
		entry209.LedgerName = "Bel_Outflow"
		entry209.MasterAccountType = journalTransactions[49].MasterAccountType
		entry209.AccountNumber = 2001
		entry209.PostingDate = ""
		entry209.ContraAccountTransactionDescription = "Insurance_Finance_Income_or_Expense"
		entry209.PostReference = "J" + strconv.Itoa(journalTransactions[49].AosStep)
		entry209.Debit = journalTransactions[49].Debit //***
		entry209.Credit = journalTransactions[49].Credit
		if journalTransactions[49].NormalAccountBalance == "Debit" {
			entry209.Balance = -(entry209.Credit - entry209.Debit) + entry208.Balance
		} else {
			entry209.Balance = entry209.Credit - entry209.Debit + entry208.Balance
		}

		entry210.ProductCode = journalTransactions[1].ProductCode
		entry210.IFRS17Group = journalTransactions[1].IFRS17Group
		entry210.LedgerName = "Bel_Outflow"
		entry210.MasterAccountType = journalTransactions[57].MasterAccountType
		entry210.AccountNumber = 2001
		entry210.PostingDate = ""
		entry210.ContraAccountTransactionDescription = "Other_CashOutflows"
		entry210.PostReference = "J" + strconv.Itoa(journalTransactions[57].AosStep)
		entry210.Debit = journalTransactions[57].Debit //***
		entry210.Credit = journalTransactions[57].Credit
		if journalTransactions[57].NormalAccountBalance == "Debit" {
			entry210.Balance = -(entry210.Credit - entry210.Debit) + entry209.Balance
		} else {
			entry210.Balance = entry210.Credit - entry210.Debit + entry209.Balance
		}

		entry211.ProductCode = journalTransactions[1].ProductCode
		entry211.IFRS17Group = journalTransactions[1].IFRS17Group
		entry211.LedgerName = "Bel_Outflow"
		entry211.MasterAccountType = journalTransactions[121].MasterAccountType
		entry211.AccountNumber = 2001
		entry211.PostingDate = ""
		entry211.ContraAccountTransactionDescription = "Bel_Outflow_Buildup"
		entry211.PostReference = "J" + strconv.Itoa(journalTransactions[121].AosStep)
		entry211.Debit = journalTransactions[121].Debit //***
		entry211.Credit = journalTransactions[121].Credit
		if journalTransactions[121].NormalAccountBalance == "Debit" {
			entry211.Balance = -(entry211.Credit - entry211.Debit) + entry210.Balance
		} else {
			entry211.Balance = entry211.Credit - entry211.Debit + entry210.Balance
		}

		entry212.ProductCode = journalTransactions[1].ProductCode
		entry212.IFRS17Group = journalTransactions[1].IFRS17Group
		entry212.LedgerName = "Bel_Outflow"
		entry212.MasterAccountType = journalTransactions[96].MasterAccountType
		entry212.AccountNumber = 2001
		entry212.PostingDate = ""
		entry212.ContraAccountTransactionDescription = "Insurance_Finance_Income_or_Expense"
		entry212.PostReference = "J" + strconv.Itoa(journalTransactions[94].AosStep)
		entry212.Debit = journalTransactions[96].Debit //***
		entry212.Credit = journalTransactions[96].Credit
		if journalTransactions[96].NormalAccountBalance == "Debit" {
			entry212.Balance = -(entry212.Credit - entry212.Debit) + entry211.Balance
		} else {
			entry212.Balance = entry212.Credit - entry212.Debit + entry211.Balance
		}

		entry213.ProductCode = journalTransactions[1].ProductCode
		entry213.IFRS17Group = journalTransactions[1].IFRS17Group
		entry213.LedgerName = "Bel_Outflow"
		entry213.MasterAccountType = journalTransactions[102].MasterAccountType
		entry213.AccountNumber = 2001
		entry213.PostingDate = ""
		entry213.ContraAccountTransactionDescription = "Insurance_Finance_Income_or_Expense"
		entry213.PostReference = "J" + strconv.Itoa(journalTransactions[102].AosStep)
		entry213.Debit = journalTransactions[102].Debit //***
		entry213.Credit = journalTransactions[102].Credit
		if journalTransactions[102].NormalAccountBalance == "Debit" {
			entry213.Balance = -(entry213.Credit - entry213.Debit) + entry212.Balance
		} else {
			entry213.Balance = entry213.Credit - entry213.Debit + entry212.Balance
		}

		subledgerEntries = append(subledgerEntries, entry201999, entry201, entry202, entry203, entry204, entry205, entry206, entry207, entry208, entry209, entry210, entry211, entry212, entry213)

		// DAC Release
		var entry301 models.SubLedgerReportEntry //3

		entry301.ProductCode = journalTransactions[1].ProductCode
		entry301.IFRS17Group = journalTransactions[1].IFRS17Group
		entry301.LedgerName = "Dac_Release"
		entry301.MasterAccountType = journalTransactions[109].MasterAccountType
		entry301.AccountNumber = 203
		entry301.PostingDate = ""
		entry301.ContraAccountTransactionDescription = "Dac"
		entry301.PostReference = "J" + strconv.Itoa(journalTransactions[109].AosStep)
		entry301.Debit = journalTransactions[109].Debit
		entry301.Credit = journalTransactions[109].Credit
		if journalTransactions[109].NormalAccountBalance == "Debit" {
			entry301.Balance = -(entry301.Credit - entry301.Debit)
		} else {
			entry301.Balance = entry301.Credit - entry301.Debit
		}
		subledgerEntries = append(subledgerEntries, entry301)

		//Risk Adjustment
		var entry401999, entry401, entry402, entry403, entry404, entry405, entry406, entry407, entry408, entry409 models.SubLedgerReportEntry //4

		entry401999.ProductCode = journalTransactions[1].ProductCode
		entry401999.IFRS17Group = journalTransactions[1].IFRS17Group
		entry401999.LedgerName = "Risk_Adjustment_B/F"
		entry401999.MasterAccountType = journalTransactions[3].MasterAccountType
		entry401999.AccountNumber = 2002999
		entry401999.PostingDate = ""
		entry401999.ContraAccountTransactionDescription = "Risk_Adjustment"
		entry401999.PostReference = "J" + strconv.Itoa(journalTransactions[3].AosStep)
		entry401999.Debit = results[0].RiskAdjustment // journalTransactions[3].Debit
		entry401999.Credit = 0                        //journalTransactions[3].Credit
		if journalTransactions[3].NormalAccountBalance == "Debit" {
			entry401999.Balance = -(entry401999.Credit - entry401999.Debit)
		} else {
			entry401999.Balance = entry401999.Credit - entry401999.Debit
		}

		entry401.ProductCode = journalTransactions[1].ProductCode
		entry401.IFRS17Group = journalTransactions[1].IFRS17Group
		entry401.LedgerName = "Risk_Adjustment"
		entry401.MasterAccountType = journalTransactions[3].MasterAccountType
		entry401.AccountNumber = 2002
		entry401.PostingDate = ""
		entry401.ContraAccountTransactionDescription = "Risk_Adjustment_B/F"
		entry401.PostReference = "J" + strconv.Itoa(journalTransactions[3].AosStep)
		entry401.Debit = 0                          // journalTransactions[3].Debit
		entry401.Credit = results[0].RiskAdjustment //journalTransactions[3].Credit
		if journalTransactions[3].NormalAccountBalance == "Debit" {
			entry401.Balance = -(entry401.Credit - entry401.Debit)
		} else {
			entry401.Balance = entry401.Credit - entry401.Debit
		}

		entry402.ProductCode = journalTransactions[1].ProductCode
		entry402.IFRS17Group = journalTransactions[1].IFRS17Group
		entry402.LedgerName = "Risk_Adjustment"
		entry402.MasterAccountType = journalTransactions[10].MasterAccountType
		entry402.AccountNumber = 2002
		entry402.PostingDate = ""
		entry402.ContraAccountTransactionDescription = "Insurance_Finance_Income_or_Expense"
		entry402.PostReference = "J" + strconv.Itoa(journalTransactions[10].AosStep)
		entry402.Debit = journalTransactions[10].Debit
		entry402.Credit = journalTransactions[10].Credit
		if journalTransactions[10].NormalAccountBalance == "Debit" {
			entry402.Balance = -(entry402.Credit - entry402.Debit) + entry401.Balance
		} else {
			entry402.Balance = entry402.Credit - entry402.Debit + entry401.Balance
		}

		entry403.ProductCode = journalTransactions[1].ProductCode
		entry403.IFRS17Group = journalTransactions[1].IFRS17Group
		entry403.LedgerName = "Risk_Adjustment"
		entry403.MasterAccountType = journalTransactions[15].MasterAccountType
		entry403.AccountNumber = 2002
		entry403.PostingDate = ""
		entry403.ContraAccountTransactionDescription = "Bel_Inflow"
		entry403.PostReference = "J" + strconv.Itoa(journalTransactions[15].AosStep)
		entry403.Debit = journalTransactions[15].Debit
		entry403.Credit = journalTransactions[15].Credit
		if journalTransactions[15].NormalAccountBalance == "Debit" {
			entry403.Balance = -(entry403.Credit - entry403.Debit) + entry402.Balance
		} else {
			entry403.Balance = entry403.Credit - entry403.Debit + entry402.Balance
		}

		entry404.ProductCode = journalTransactions[1].ProductCode
		entry404.IFRS17Group = journalTransactions[1].IFRS17Group
		entry404.LedgerName = "Risk_Adjustment"
		entry404.MasterAccountType = journalTransactions[19].MasterAccountType
		entry404.AccountNumber = 2002
		entry404.PostingDate = ""
		entry404.ContraAccountTransactionDescription = "Loss_Component_NB"
		entry404.PostReference = "J" + strconv.Itoa(journalTransactions[19].AosStep)
		entry404.Debit = journalTransactions[19].Debit
		entry404.Credit = journalTransactions[19].Credit
		if journalTransactions[19].NormalAccountBalance == "Debit" {
			entry404.Balance = -(entry404.Credit - entry404.Debit) + entry403.Balance
		} else {
			entry404.Balance = entry404.Credit - entry404.Debit + entry403.Balance
		}

		entry405.ProductCode = journalTransactions[1].ProductCode
		entry405.IFRS17Group = journalTransactions[1].IFRS17Group
		entry405.LedgerName = "Risk_Adjustment"
		entry405.MasterAccountType = journalTransactions[23].MasterAccountType
		entry405.AccountNumber = 2002
		entry405.PostingDate = ""
		entry405.ContraAccountTransactionDescription = "Risk_Adjustment_Release"
		entry405.PostReference = "J" + strconv.Itoa(journalTransactions[23].AosStep)
		entry405.Debit = journalTransactions[23].Debit
		entry405.Credit = journalTransactions[23].Credit
		if journalTransactions[23].NormalAccountBalance == "Debit" {
			entry405.Balance = -(entry405.Credit - entry405.Debit) + entry404.Balance
		} else {
			entry405.Balance = entry405.Credit - entry405.Debit + entry404.Balance
		}

		entry406.ProductCode = journalTransactions[1].ProductCode
		entry406.IFRS17Group = journalTransactions[1].IFRS17Group
		entry406.LedgerName = "Risk_Adjustment"
		entry406.MasterAccountType = journalTransactions[59].MasterAccountType
		entry406.AccountNumber = 2002
		entry406.PostingDate = ""
		entry406.ContraAccountTransactionDescription = "Insurance_Finance_Income_or_Expense"
		entry406.PostReference = "J" + strconv.Itoa(journalTransactions[59].AosStep)
		entry406.Debit = journalTransactions[59].Debit
		entry406.Credit = journalTransactions[59].Credit
		if journalTransactions[59].NormalAccountBalance == "Debit" {
			entry406.Balance = -(entry406.Credit - entry406.Debit) + entry405.Balance
		} else {
			entry406.Balance = entry406.Credit - entry406.Debit + entry405.Balance
		}

		entry407.ProductCode = journalTransactions[1].ProductCode
		entry407.IFRS17Group = journalTransactions[1].IFRS17Group
		entry407.LedgerName = "Risk_Adjustment"
		entry407.MasterAccountType = journalTransactions[98].MasterAccountType
		entry407.AccountNumber = 2002
		entry407.PostingDate = ""
		entry407.ContraAccountTransactionDescription = "Insurance_Finance_Income_or_Expense"
		entry407.PostReference = "J" + strconv.Itoa(journalTransactions[98].AosStep)
		entry407.Debit = journalTransactions[98].Debit
		entry407.Credit = journalTransactions[98].Credit
		if journalTransactions[98].NormalAccountBalance == "Debit" {
			entry407.Balance = -(entry407.Credit - entry407.Debit) + entry406.Balance
		} else {
			entry407.Balance = entry407.Credit - entry407.Debit + entry406.Balance
		}

		entry408.ProductCode = journalTransactions[1].ProductCode
		entry408.IFRS17Group = journalTransactions[1].IFRS17Group
		entry408.LedgerName = "Risk_Adjustment"
		entry408.MasterAccountType = journalTransactions[104].MasterAccountType
		entry408.AccountNumber = 2002
		entry408.PostingDate = ""
		entry408.ContraAccountTransactionDescription = "Insurance_Finance_Income_or_Expense"
		entry408.PostReference = "J" + strconv.Itoa(journalTransactions[104].AosStep)
		entry408.Debit = journalTransactions[104].Debit
		entry408.Credit = journalTransactions[104].Credit
		if journalTransactions[104].NormalAccountBalance == "Debit" {
			entry408.Balance = -(entry408.Credit - entry408.Debit) + entry407.Balance
		} else {
			entry408.Balance = entry408.Credit - entry408.Debit + entry407.Balance
		}

		entry409.ProductCode = journalTransactions[1].ProductCode
		entry409.IFRS17Group = journalTransactions[1].IFRS17Group
		entry409.LedgerName = "Risk_Adjustment"
		entry409.MasterAccountType = journalTransactions[123].MasterAccountType
		entry409.AccountNumber = 2002
		entry409.PostingDate = ""
		entry409.ContraAccountTransactionDescription = "Risk_Adjustment_Buildup"
		entry409.PostReference = "J" + strconv.Itoa(journalTransactions[123].AosStep)
		entry409.Debit = journalTransactions[123].Debit
		entry409.Credit = journalTransactions[123].Credit
		if journalTransactions[123].NormalAccountBalance == "Debit" {
			entry409.Balance = -(entry409.Credit - entry409.Debit) + entry408.Balance
		} else {
			entry409.Balance = entry409.Credit - entry409.Debit + entry408.Balance
		}

		subledgerEntries = append(subledgerEntries, entry401999, entry401, entry402, entry403, entry404, entry405, entry406, entry407, entry408, entry409)

		//CSM
		var entry500, entry501, entry502, entry503, entry504, entry505 models.SubLedgerReportEntry //5

		entry500.ProductCode = journalTransactions[1].ProductCode
		entry500.IFRS17Group = journalTransactions[1].IFRS17Group
		entry500.LedgerName = "CSM_B/F(Offset)"
		entry500.MasterAccountType = journalTransactions[4].MasterAccountType
		entry500.AccountNumber = 2003999
		entry500.PostingDate = ""
		entry500.ContraAccountTransactionDescription = "CSM"
		entry500.PostReference = "J" + strconv.Itoa(journalTransactions[4].AosStep)
		entry500.Debit = results[0].CSMBuildup //journalTransactions[4].Debit
		entry500.Credit = 0                    //journalTransactions[4].Credit
		if journalTransactions[4].NormalAccountBalance == "Debit" {
			entry500.Balance = -(entry500.Credit - entry500.Debit)
		} else {
			entry500.Balance = entry500.Credit - entry500.Debit
		}

		entry501.ProductCode = journalTransactions[1].ProductCode
		entry501.IFRS17Group = journalTransactions[1].IFRS17Group
		entry501.LedgerName = "CSM"
		entry501.MasterAccountType = journalTransactions[4].MasterAccountType
		entry501.AccountNumber = 2003
		entry501.PostingDate = ""
		entry501.ContraAccountTransactionDescription = "CSM_B/F"
		entry501.PostReference = "J" + strconv.Itoa(journalTransactions[4].AosStep)
		entry501.Debit = 0                      //journalTransactions[4].Debit
		entry501.Credit = results[0].CSMBuildup //journalTransactions[4].Credit
		if journalTransactions[4].NormalAccountBalance == "Debit" {
			entry501.Balance = -(entry501.Credit - entry501.Debit)
		} else {
			entry501.Balance = entry501.Credit - entry501.Debit
		}

		entry502.ProductCode = journalTransactions[1].ProductCode
		entry502.IFRS17Group = journalTransactions[1].IFRS17Group
		entry502.LedgerName = "CSM"
		entry502.MasterAccountType = journalTransactions[16].MasterAccountType
		entry502.AccountNumber = 2003
		entry502.PostingDate = ""
		entry502.ContraAccountTransactionDescription = "Bel_Inflow"
		entry502.PostReference = "J" + strconv.Itoa(journalTransactions[16].AosStep)
		entry502.Debit = journalTransactions[16].Debit
		entry502.Credit = journalTransactions[16].Credit
		if journalTransactions[16].NormalAccountBalance == "Debit" {
			entry502.Balance = -(entry502.Credit - entry502.Debit) + entry501.Balance
		} else {
			entry502.Balance = entry502.Credit - entry502.Debit + entry501.Balance
		}

		entry503.ProductCode = journalTransactions[1].ProductCode
		entry503.IFRS17Group = journalTransactions[1].IFRS17Group
		entry503.LedgerName = "CSM"
		entry503.MasterAccountType = journalTransactions[43].MasterAccountType
		entry503.AccountNumber = 2003
		entry503.PostingDate = ""
		entry503.ContraAccountTransactionDescription = "Cash_and_Cash_Equivalent"
		entry503.PostReference = "J" + strconv.Itoa(journalTransactions[43].AosStep)
		entry503.Debit = journalTransactions[43].Debit
		entry503.Credit = journalTransactions[43].Credit
		if journalTransactions[43].NormalAccountBalance == "Debit" {
			entry503.Balance = -(entry503.Credit - entry503.Debit) + entry502.Balance
		} else {
			entry503.Balance = entry503.Credit - entry503.Debit + entry502.Balance
		}

		entry504.ProductCode = journalTransactions[1].ProductCode
		entry504.IFRS17Group = journalTransactions[1].IFRS17Group
		entry504.LedgerName = "CSM"
		entry504.MasterAccountType = journalTransactions[92].MasterAccountType
		entry504.AccountNumber = 2003
		entry504.PostingDate = ""
		entry504.ContraAccountTransactionDescription = "Bel_Inflow_Buildup"
		entry504.PostReference = "J" + strconv.Itoa(journalTransactions[92].AosStep)
		entry504.Debit = journalTransactions[92].Debit
		entry504.Credit = journalTransactions[92].Credit
		if journalTransactions[92].NormalAccountBalance == "Debit" {
			entry504.Balance = -(entry504.Credit - entry504.Debit) + entry503.Balance
		} else {
			entry504.Balance = entry504.Credit - entry504.Debit + entry503.Balance
		}

		entry505.ProductCode = journalTransactions[1].ProductCode
		entry505.IFRS17Group = journalTransactions[1].IFRS17Group
		entry505.LedgerName = "CSM"
		entry505.MasterAccountType = journalTransactions[108].MasterAccountType
		entry505.AccountNumber = 2003
		entry505.PostingDate = ""
		entry505.ContraAccountTransactionDescription = "CSM_Release"
		entry505.PostReference = "J" + strconv.Itoa(journalTransactions[108].AosStep)
		entry505.Debit = journalTransactions[108].Debit
		entry505.Credit = journalTransactions[108].Credit
		if journalTransactions[108].NormalAccountBalance == "Debit" {
			entry505.Balance = -(entry505.Credit - entry505.Debit) + entry504.Balance
		} else {
			entry505.Balance = entry505.Credit - entry505.Debit + entry504.Balance
		}
		subledgerEntries = append(subledgerEntries, entry500, entry501, entry502, entry503, entry504, entry505)

		//Risk Adjustment Release
		var entry601 models.SubLedgerReportEntry //6

		entry601.ProductCode = journalTransactions[1].ProductCode
		entry601.IFRS17Group = journalTransactions[1].IFRS17Group
		entry601.LedgerName = "Risk_Adjustment_Release"
		entry601.MasterAccountType = journalTransactions[22].MasterAccountType
		entry601.AccountNumber = 101
		entry601.PostingDate = ""
		entry601.ContraAccountTransactionDescription = "Risk_Adjustment"
		entry601.PostReference = "J" + strconv.Itoa(journalTransactions[22].AosStep)
		entry601.Debit = journalTransactions[22].Debit
		entry601.Credit = journalTransactions[22].Credit
		if journalTransactions[22].NormalAccountBalance == "Debit" {
			entry601.Balance = -(entry601.Credit - entry601.Debit)
		} else {
			entry601.Balance = entry601.Credit - entry601.Debit
		}
		subledgerEntries = append(subledgerEntries, entry601)

		//CSM Release
		var entry701 models.SubLedgerReportEntry //7

		entry701.ProductCode = journalTransactions[1].ProductCode
		entry701.IFRS17Group = journalTransactions[1].IFRS17Group
		entry701.LedgerName = "CSM_Release"
		entry701.MasterAccountType = journalTransactions[107].MasterAccountType
		entry701.AccountNumber = 106
		entry701.PostingDate = ""
		entry701.ContraAccountTransactionDescription = "CSM"
		entry701.PostReference = "J" + strconv.Itoa(journalTransactions[107].AosStep)
		entry701.Debit = journalTransactions[107].Debit
		entry701.Credit = journalTransactions[107].Credit
		if journalTransactions[107].NormalAccountBalance == "Debit" {
			entry701.Balance = -(entry701.Credit - entry701.Debit)
		} else {
			entry701.Balance = entry701.Credit - entry701.Debit
		}

		subledgerEntries = append(subledgerEntries, entry701)

		//Amortised Acquisition Cash flows
		var entry801 models.SubLedgerReportEntry //8

		entry801.ProductCode = journalTransactions[1].ProductCode
		entry801.IFRS17Group = journalTransactions[1].IFRS17Group
		entry801.LedgerName = "Amortisation_of_Acquisition_Cashflows"
		entry801.MasterAccountType = journalTransactions[111].MasterAccountType
		entry801.AccountNumber = 108
		entry801.PostingDate = ""
		entry801.ContraAccountTransactionDescription = "Acquisition_Cashflows"
		entry801.PostReference = "J" + strconv.Itoa(journalTransactions[111].AosStep)
		entry801.Debit = journalTransactions[111].Debit
		entry801.Credit = journalTransactions[111].Credit
		if journalTransactions[111].NormalAccountBalance == "Debit" {
			entry801.Balance = -(entry801.Credit - entry801.Debit)
		} else {
			entry801.Balance = entry801.Credit - entry801.Debit
		}

		subledgerEntries = append(subledgerEntries, entry801)

		//Expected Mortality Claims
		var entry901 models.SubLedgerReportEntry //9

		entry901.ProductCode = journalTransactions[1].ProductCode
		entry901.IFRS17Group = journalTransactions[1].IFRS17Group
		entry901.LedgerName = "Expected_Mortality_Claims"
		entry901.MasterAccountType = journalTransactions[26].MasterAccountType
		entry901.AccountNumber = 102
		entry901.PostingDate = ""
		entry901.ContraAccountTransactionDescription = "Bel_Outflow"
		entry901.PostReference = "J" + strconv.Itoa(journalTransactions[26].AosStep)
		entry901.Debit = journalTransactions[26].Debit
		entry901.Credit = journalTransactions[26].Credit
		if journalTransactions[26].NormalAccountBalance == "Debit" {
			entry901.Balance = -(entry901.Credit - entry901.Debit)
		} else {
			entry901.Balance = entry901.Credit - entry901.Debit
		}

		subledgerEntries = append(subledgerEntries, entry901)

		//Expected Expenses
		var entry1001 models.SubLedgerReportEntry //10

		entry1001.ProductCode = journalTransactions[1].ProductCode
		entry1001.IFRS17Group = journalTransactions[1].IFRS17Group
		entry1001.LedgerName = "Expected_Expenses"
		entry1001.MasterAccountType = journalTransactions[30].MasterAccountType
		entry1001.AccountNumber = 103
		entry1001.PostingDate = ""
		entry1001.ContraAccountTransactionDescription = "Bel_Outflow"
		entry1001.PostReference = "J" + strconv.Itoa(journalTransactions[30].AosStep)
		entry1001.Debit = journalTransactions[30].Debit
		entry1001.Credit = journalTransactions[30].Credit
		if journalTransactions[30].NormalAccountBalance == "Debit" {
			entry1001.Balance = -(entry1001.Credit - entry1001.Debit)
		} else {
			entry1001.Balance = entry1001.Credit - entry1001.Debit
		}

		subledgerEntries = append(subledgerEntries, entry1001)

		//Expected Retrenchment
		var entry1101 models.SubLedgerReportEntry //11

		entry1101.ProductCode = journalTransactions[1].ProductCode
		entry1101.IFRS17Group = journalTransactions[1].IFRS17Group
		entry1101.LedgerName = "Expected_Retrenchment"
		entry1101.MasterAccountType = journalTransactions[34].MasterAccountType
		entry1101.AccountNumber = 104
		entry1101.PostingDate = ""
		entry1101.ContraAccountTransactionDescription = "Bel_Outflow"
		entry1101.PostReference = "J" + strconv.Itoa(journalTransactions[34].AosStep)
		entry1101.Debit = journalTransactions[34].Debit
		entry1101.Credit = journalTransactions[34].Credit
		if journalTransactions[34].NormalAccountBalance == "Debit" {
			entry1101.Balance = -(entry1101.Credit - entry1101.Debit)
		} else {
			entry1101.Balance = entry1101.Credit - entry1101.Debit
		}

		subledgerEntries = append(subledgerEntries, entry1101)

		//Expected Morbidity
		var entry1201, entryX1201, entryY1201, entryZ1201, entryZA1201 models.SubLedgerReportEntry //12

		entry1201.ProductCode = journalTransactions[1].ProductCode
		entry1201.IFRS17Group = journalTransactions[1].IFRS17Group
		entry1201.LedgerName = "Expected_Morbidity"
		entry1201.MasterAccountType = journalTransactions[38].MasterAccountType
		entry1201.AccountNumber = 105
		entry1201.PostingDate = ""
		entry1201.ContraAccountTransactionDescription = "Bel_Outflow"
		entry1201.PostReference = "J" + strconv.Itoa(journalTransactions[38].AosStep)
		entry1201.Debit = journalTransactions[38].Debit
		entry1201.Credit = journalTransactions[38].Credit
		if journalTransactions[38].NormalAccountBalance == "Debit" {
			entry1201.Balance = -(entry1201.Credit - entry1201.Debit)
		} else {
			entry1201.Balance = entry1201.Credit - entry1201.Debit
		}

		entryX1201.ProductCode = journalTransactions[1].ProductCode
		entryX1201.IFRS17Group = journalTransactions[1].IFRS17Group
		entryX1201.LedgerName = "Other_CashOutflows"
		entryX1201.MasterAccountType = journalTransactions[54].MasterAccountType
		entryX1201.AccountNumber = 105
		entryX1201.PostingDate = ""
		entryX1201.ContraAccountTransactionDescription = "Bel_Outflow"
		entryX1201.PostReference = "J" + strconv.Itoa(journalTransactions[54].AosStep)
		entryX1201.Debit = journalTransactions[54].Debit
		entryX1201.Credit = journalTransactions[54].Credit
		if journalTransactions[54].NormalAccountBalance == "Debit" {
			entryX1201.Balance = -(entryX1201.Credit - entryX1201.Debit)
		} else {
			entryX1201.Balance = entryX1201.Credit - entryX1201.Debit
		}

		entryY1201.ProductCode = journalTransactions[1].ProductCode
		entryY1201.IFRS17Group = journalTransactions[1].IFRS17Group
		entryY1201.LedgerName = "Premium_Debtors"
		entryY1201.MasterAccountType = journalTransactions[56].MasterAccountType
		entryY1201.AccountNumber = 1006
		entryY1201.PostingDate = ""
		entryY1201.ContraAccountTransactionDescription = "Cash_and_Cash_Equivalent"
		entryY1201.PostReference = "J" + strconv.Itoa(journalTransactions[56].AosStep)
		entryY1201.Debit = journalTransactions[56].Debit
		entryY1201.Credit = journalTransactions[56].Credit
		if journalTransactions[56].NormalAccountBalance == "Debit" {
			entryY1201.Balance = -(entryY1201.Credit - entryY1201.Debit)
		} else {
			entryY1201.Balance = entryY1201.Credit - entryY1201.Debit
		}

		entryZ1201.ProductCode = journalTransactions[1].ProductCode
		entryZ1201.IFRS17Group = journalTransactions[1].IFRS17Group
		entryZ1201.LedgerName = "Premium_Debtors"
		entryZ1201.MasterAccountType = journalTransactions[51].MasterAccountType
		entryZ1201.AccountNumber = 1006
		entryZ1201.PostingDate = ""
		entryZ1201.ContraAccountTransactionDescription = "BEL_Inflow"
		entryZ1201.PostReference = "J" + strconv.Itoa(journalTransactions[51].AosStep)
		entryZ1201.Debit = journalTransactions[51].Debit
		entryZ1201.Credit = journalTransactions[51].Credit
		if journalTransactions[51].NormalAccountBalance == "Debit" {
			entryZ1201.Balance = -(entryZ1201.Credit - entryZ1201.Debit) + entryY1201.Balance
		} else {
			entryZ1201.Balance = entryZ1201.Credit - entryZ1201.Debit + entryY1201.Balance
		}

		entryZA1201.ProductCode = journalTransactions[1].ProductCode
		entryZA1201.IFRS17Group = journalTransactions[1].IFRS17Group
		entryZA1201.LedgerName = "Premium_Debtors"
		entryZA1201.MasterAccountType = journalTransactions[53].MasterAccountType
		entryZA1201.AccountNumber = 1006
		entryZA1201.PostingDate = ""
		entryZA1201.ContraAccountTransactionDescription = "Experience_Premium_Variance"
		entryZA1201.PostReference = "J" + strconv.Itoa(journalTransactions[53].AosStep)
		entryZA1201.Debit = journalTransactions[53].Debit
		entryZA1201.Credit = journalTransactions[53].Credit
		if journalTransactions[53].NormalAccountBalance == "Debit" {
			entryZA1201.Balance = -(entryZA1201.Credit - entryZA1201.Debit) + entryZ1201.Balance
		} else {
			entryZA1201.Balance = entryZA1201.Credit - entryZA1201.Debit + entryZ1201.Balance
		}

		subledgerEntries = append(subledgerEntries, entry1201, entryX1201, entryY1201, entryZ1201, entryZA1201)
		//DAC
		var entry1301, entry1302, entry1303 models.SubLedgerReportEntry //13

		entry1301.ProductCode = journalTransactions[1].ProductCode
		entry1301.IFRS17Group = journalTransactions[1].IFRS17Group
		entry1301.LedgerName = "DAC"
		entry1301.MasterAccountType = journalTransactions[1].MasterAccountType
		entry1301.AccountNumber = 1002
		entry1301.PostingDate = ""
		entry1301.ContraAccountTransactionDescription = "Acquisition_Cashflows"
		entry1301.PostReference = "J" + strconv.Itoa(journalTransactions[1].AosStep)
		entry1301.Debit = journalTransactions[1].Debit
		entry1301.Credit = journalTransactions[1].Credit
		if journalTransactions[1].NormalAccountBalance == "Debit" {
			entry1301.Balance = -(entry1301.Credit - entry1301.Debit)
		} else {
			entry1301.Balance = entry1301.Credit - entry1301.Debit
		}

		entry1302.ProductCode = journalTransactions[1].ProductCode
		entry1302.IFRS17Group = journalTransactions[1].IFRS17Group
		entry1302.LedgerName = "DAC"
		entry1302.MasterAccountType = journalTransactions[20].MasterAccountType
		entry1302.AccountNumber = 1002
		entry1302.PostingDate = ""
		entry1302.ContraAccountTransactionDescription = "Acquisition_Cashflows"
		entry1302.PostReference = "J" + strconv.Itoa(journalTransactions[20].AosStep)
		entry1302.Debit = journalTransactions[20].Debit
		entry1302.Credit = journalTransactions[20].Credit
		if journalTransactions[20].NormalAccountBalance == "Debit" {
			entry1302.Balance = -(entry1302.Credit - entry1302.Debit) + entry1301.Balance
		} else {
			entry1302.Balance = entry1302.Credit - entry1302.Debit + entry1301.Balance
		}

		entry1303.ProductCode = journalTransactions[1].ProductCode
		entry1303.IFRS17Group = journalTransactions[1].IFRS17Group
		entry1303.LedgerName = "DAC"
		entry1303.MasterAccountType = journalTransactions[110].MasterAccountType
		entry1303.AccountNumber = 1002
		entry1303.PostingDate = ""
		entry1303.ContraAccountTransactionDescription = "Dac_Release"
		entry1303.PostReference = "J" + strconv.Itoa(journalTransactions[110].AosStep)
		entry1303.Debit = journalTransactions[110].Debit
		entry1303.Credit = journalTransactions[110].Credit
		if journalTransactions[110].NormalAccountBalance == "Debit" {
			entry1303.Balance = -(entry1303.Credit - entry1303.Debit) + entry1302.Balance
		} else {
			entry1303.Balance = entry1303.Credit - entry1303.Debit + entry1302.Balance
		}

		subledgerEntries = append(subledgerEntries, entry1301, entry1302, entry1303)

		//Acquisition Cashflows
		var entry1401, entry1402 models.SubLedgerReportEntry //14

		entry1401.ProductCode = journalTransactions[1].ProductCode
		entry1401.IFRS17Group = journalTransactions[1].IFRS17Group
		entry1401.LedgerName = "Acquisition_Cashflows"
		entry1401.MasterAccountType = journalTransactions[21].MasterAccountType
		entry1401.AccountNumber = 2007
		entry1401.PostingDate = ""
		entry1401.ContraAccountTransactionDescription = "DAC"
		entry1401.PostReference = "J" + strconv.Itoa(journalTransactions[21].AosStep)
		entry1401.Debit = journalTransactions[21].Debit
		entry1401.Credit = journalTransactions[21].Credit
		if journalTransactions[21].NormalAccountBalance == "Debit" {
			entry1401.Balance = -(entry1401.Credit - entry1401.Debit)
		} else {
			entry1401.Balance = entry1401.Credit - entry1401.Debit
		}

		entry1402.ProductCode = journalTransactions[1].ProductCode
		entry1402.IFRS17Group = journalTransactions[1].IFRS17Group
		entry1402.LedgerName = "Acquisition_Cashflows"
		entry1402.MasterAccountType = journalTransactions[112].MasterAccountType
		entry1402.AccountNumber = 2007
		entry1402.PostingDate = ""
		entry1402.ContraAccountTransactionDescription = "Amortisation_of_Acquisition_Cashflows"
		entry1402.PostReference = "J" + strconv.Itoa(journalTransactions[112].AosStep)
		entry1402.Debit = journalTransactions[112].Debit
		entry1402.Credit = journalTransactions[112].Credit
		if journalTransactions[112].NormalAccountBalance == "Debit" {
			entry1402.Balance = -(entry1402.Credit - entry1402.Debit) + entry1401.Balance
		} else {
			entry1402.Balance = entry1402.Credit - entry1402.Debit + entry1401.Balance
		}

		subledgerEntries = append(subledgerEntries, entry1401, entry1402)

		//DAC Release
		var entry1501 models.SubLedgerReportEntry //15

		entry1501.ProductCode = journalTransactions[1].ProductCode
		entry1501.IFRS17Group = journalTransactions[1].IFRS17Group
		entry1501.LedgerName = "DAC_Release"
		entry1501.MasterAccountType = journalTransactions[109].MasterAccountType
		entry1501.AccountNumber = 106
		entry1501.PostingDate = ""
		entry1501.ContraAccountTransactionDescription = "DAC"
		entry1501.PostReference = "J" + strconv.Itoa(journalTransactions[109].AosStep)
		entry1501.Debit = journalTransactions[109].Debit
		entry1501.Credit = journalTransactions[109].Credit
		if journalTransactions[109].NormalAccountBalance == "Debit" {
			entry1501.Balance = -(entry1501.Credit - entry1501.Debit)
		} else {
			entry1501.Balance = entry1501.Credit - entry1501.Debit
		}

		subledgerEntries = append(subledgerEntries, entry1501)

		//Actual Mortality Claims
		var entry1601 models.SubLedgerReportEntry //16

		entry1601.ProductCode = journalTransactions[1].ProductCode
		entry1601.IFRS17Group = journalTransactions[1].IFRS17Group
		entry1601.LedgerName = "Actual_Mortality_Claims"
		entry1601.MasterAccountType = journalTransactions[113].MasterAccountType
		entry1601.AccountNumber = 205
		entry1601.PostingDate = ""
		entry1601.ContraAccountTransactionDescription = "OCR"
		entry1601.PostReference = "J" + strconv.Itoa(journalTransactions[113].AosStep)
		entry1601.Debit = journalTransactions[113].Debit
		entry1601.Credit = journalTransactions[113].Credit
		if journalTransactions[113].NormalAccountBalance == "Debit" {
			entry1601.Balance = -(entry1601.Credit - entry1601.Debit)
		} else {
			entry1601.Balance = entry1601.Credit - entry1601.Debit
		}

		subledgerEntries = append(subledgerEntries, entry1601)

		//Actual Expenses
		var entry1701 models.SubLedgerReportEntry //17

		entry1701.ProductCode = journalTransactions[1].ProductCode
		entry1701.IFRS17Group = journalTransactions[1].IFRS17Group
		entry1701.LedgerName = "Actual_Expenses"
		entry1701.MasterAccountType = journalTransactions[115].MasterAccountType
		entry1701.AccountNumber = 206
		entry1701.PostingDate = ""
		entry1701.ContraAccountTransactionDescription = "Cash_and_Cash_Equivalent"
		entry1701.PostReference = "J" + strconv.Itoa(journalTransactions[115].AosStep)
		entry1701.Debit = journalTransactions[115].Debit
		entry1701.Credit = journalTransactions[115].Credit
		if journalTransactions[115].NormalAccountBalance == "Debit" {
			entry1701.Balance = -(entry1701.Credit - entry1701.Debit)
		} else {
			entry1701.Balance = entry1701.Credit - entry1701.Debit
		}

		subledgerEntries = append(subledgerEntries, entry1701)

		//Actual Retrenchment Claims
		var entry1801 models.SubLedgerReportEntry //18

		entry1801.ProductCode = journalTransactions[1].ProductCode
		entry1801.IFRS17Group = journalTransactions[1].IFRS17Group
		entry1801.LedgerName = "Actual_Retrenchment_Claims"
		entry1801.MasterAccountType = journalTransactions[117].MasterAccountType
		entry1801.AccountNumber = 207
		entry1801.PostingDate = ""
		entry1801.ContraAccountTransactionDescription = "Cash_and_Cash_Equivalent"
		entry1801.PostReference = "J" + strconv.Itoa(journalTransactions[117].AosStep)
		entry1801.Debit = journalTransactions[117].Debit
		entry1801.Credit = journalTransactions[117].Credit
		if journalTransactions[117].NormalAccountBalance == "Debit" {
			entry1801.Balance = -(entry1801.Credit - entry1801.Debit)
		} else {
			entry1801.Balance = entry1801.Credit - entry1801.Debit
		}

		subledgerEntries = append(subledgerEntries, entry1801)

		//Actual Morbidity Claims
		var entry1901 models.SubLedgerReportEntry //19

		entry1901.ProductCode = journalTransactions[1].ProductCode
		entry1901.IFRS17Group = journalTransactions[1].IFRS17Group
		entry1901.LedgerName = "Actual_Morbidity_Claims"
		entry1901.MasterAccountType = journalTransactions[119].MasterAccountType
		entry1901.AccountNumber = 208
		entry1901.PostingDate = ""
		entry1901.ContraAccountTransactionDescription = "OCR"
		entry1901.PostReference = "J" + strconv.Itoa(journalTransactions[119].AosStep)
		entry1901.Debit = journalTransactions[119].Debit
		entry1901.Credit = journalTransactions[119].Credit
		if journalTransactions[119].NormalAccountBalance == "Debit" {
			entry1901.Balance = -(entry1901.Credit - entry1901.Debit)
		} else {
			entry1901.Balance = entry1901.Credit - entry1901.Debit
		}

		subledgerEntries = append(subledgerEntries, entry1901)

		//Bel Inflow Buildup
		var entry2001, entry2002, entry2003, entry2004, entry2005, entry2006 models.SubLedgerReportEntry //20

		entry2001.ProductCode = journalTransactions[1].ProductCode
		entry2001.IFRS17Group = journalTransactions[1].IFRS17Group
		entry2001.LedgerName = "Bel_Inflow_Buildup"
		entry2001.MasterAccountType = journalTransactions[91].MasterAccountType
		entry2001.AccountNumber = 1004
		entry2001.PostingDate = ""
		entry2001.ContraAccountTransactionDescription = "RA_Buildup"
		entry2001.PostReference = "J" + strconv.Itoa(journalTransactions[91].AosStep)
		entry2001.Debit = journalTransactions[94].Credit
		entry2001.Credit = journalTransactions[94].Debit
		if journalTransactions[91].NormalAccountBalance == "Debit" {
			entry2001.Balance = -(entry2001.Credit - entry2001.Debit)
		} else {
			entry2001.Balance = entry2001.Credit - entry2001.Debit
		}

		entry2002.ProductCode = journalTransactions[1].ProductCode
		entry2002.IFRS17Group = journalTransactions[1].IFRS17Group
		entry2002.LedgerName = "Bel_Inflow_Buildup"
		entry2002.MasterAccountType = journalTransactions[91].MasterAccountType
		entry2002.AccountNumber = 1004
		entry2002.PostingDate = ""
		entry2002.ContraAccountTransactionDescription = "CSM"
		entry2002.PostReference = "J" + strconv.Itoa(journalTransactions[91].AosStep)
		entry2002.Debit = journalTransactions[92].Credit
		entry2002.Credit = journalTransactions[92].Debit
		if journalTransactions[91].NormalAccountBalance == "Debit" {
			entry2002.Balance = -(entry2002.Credit - entry2002.Debit) + entry2001.Balance
		} else {
			entry2002.Balance = entry2002.Credit - entry2002.Debit + entry2001.Balance
		}

		entry2003.ProductCode = journalTransactions[1].ProductCode
		entry2003.IFRS17Group = journalTransactions[1].IFRS17Group
		entry2003.LedgerName = "Bel_Inflow_Buildup"
		entry2003.MasterAccountType = journalTransactions[91].MasterAccountType
		entry2003.AccountNumber = 1004
		entry2003.PostingDate = ""
		entry2003.ContraAccountTransactionDescription = "Bel_Outflow_Buildup"
		entry2003.PostReference = "J" + strconv.Itoa(journalTransactions[91].AosStep)
		entry2003.Debit = journalTransactions[93].Credit
		entry2003.Credit = journalTransactions[93].Debit
		if journalTransactions[91].NormalAccountBalance == "Debit" {
			entry2003.Balance = -(entry2003.Credit - entry2003.Debit) + entry2002.Balance
		} else {
			entry2003.Balance = entry2003.Credit - entry2003.Debit + entry2002.Balance
		}

		entry2004.ProductCode = journalTransactions[1].ProductCode /// check the journal transactions
		entry2004.IFRS17Group = journalTransactions[1].IFRS17Group
		entry2004.LedgerName = "Bel_Inflow_Buildup"
		entry2004.MasterAccountType = journalTransactions[91].MasterAccountType
		entry2004.AccountNumber = 1004
		entry2004.PostingDate = ""
		entry2004.ContraAccountTransactionDescription = "Loss_Component_Adjustment"
		entry2004.PostReference = "J" + strconv.Itoa(journalTransactions[91].AosStep)
		entry2004.Debit = journalTransactions[90].Credit
		entry2004.Credit = journalTransactions[90].Debit
		if journalTransactions[91].NormalAccountBalance == "Debit" {
			entry2004.Balance = -(entry2004.Credit - entry2004.Debit) + entry2003.Balance
		} else {
			entry2004.Balance = entry2004.Credit - entry2004.Debit + entry2003.Balance
		}

		entry2005.ProductCode = journalTransactions[1].ProductCode
		entry2005.IFRS17Group = journalTransactions[1].IFRS17Group
		entry2005.LedgerName = "Bel_Inflow_Buildup"
		entry2005.MasterAccountType = journalTransactions[91].MasterAccountType
		entry2005.AccountNumber = 1004
		entry2005.PostingDate = ""
		entry2005.ContraAccountTransactionDescription = "Loss_Component_Reversal"
		entry2005.PostReference = "J" + strconv.Itoa(journalTransactions[91].AosStep)
		entry2005.Debit = 0  //journalTransactions[74].Credit
		entry2005.Credit = 0 //journalTransactions[74].Debit
		if journalTransactions[91].NormalAccountBalance == "Debit" {
			entry2005.Balance = -(entry2005.Credit - entry2005.Debit) + entry2004.Balance
		} else {
			entry2005.Balance = entry2005.Credit - entry2005.Debit + entry2004.Balance
		}

		entry2006.ProductCode = journalTransactions[1].ProductCode
		entry2006.IFRS17Group = journalTransactions[1].IFRS17Group
		entry2006.LedgerName = "Bel_Inflow_Buildup"
		entry2006.MasterAccountType = journalTransactions[126].MasterAccountType
		entry2006.AccountNumber = 1004
		entry2006.PostingDate = ""
		entry2006.ContraAccountTransactionDescription = "Bel_Inflow"
		entry2006.PostReference = "J" + strconv.Itoa(journalTransactions[126].AosStep)
		entry2006.Debit = journalTransactions[126].Debit
		entry2006.Credit = journalTransactions[126].Credit
		if journalTransactions[126].NormalAccountBalance == "Debit" {
			entry2006.Balance = -(entry2006.Credit - entry2006.Debit) + entry2005.Balance
		} else {
			entry2006.Balance = entry2006.Credit - entry2006.Debit + entry2005.Balance
		}

		subledgerEntries = append(subledgerEntries, entry2001, entry2002, entry2003, entry2004, entry2005, entry2006)

		//Bel Outflow Buildup
		var entry2101, entry2102 models.SubLedgerReportEntry //21

		entry2101.ProductCode = journalTransactions[1].ProductCode
		entry2101.IFRS17Group = journalTransactions[1].IFRS17Group
		entry2101.LedgerName = "Bel_Outlfow_Buildup"
		entry2101.MasterAccountType = journalTransactions[93].MasterAccountType
		entry2101.AccountNumber = 2008
		entry2101.PostingDate = ""
		entry2101.ContraAccountTransactionDescription = "Bel_Inflow_Buildup"
		entry2101.PostReference = "J" + strconv.Itoa(journalTransactions[93].AosStep)
		entry2101.Debit = journalTransactions[93].Debit
		entry2101.Credit = journalTransactions[93].Credit
		if journalTransactions[93].NormalAccountBalance == "Debit" {
			entry2101.Balance = -(entry2101.Credit - entry2101.Debit)
		} else {
			entry2101.Balance = entry2101.Credit - entry2101.Debit
		}

		entry2102.ProductCode = journalTransactions[1].ProductCode
		entry2102.IFRS17Group = journalTransactions[1].IFRS17Group
		entry2102.LedgerName = "Bel_Outlfow_Buildup"
		entry2102.MasterAccountType = journalTransactions[122].MasterAccountType
		entry2102.AccountNumber = 2008
		entry2102.PostingDate = ""
		entry2102.ContraAccountTransactionDescription = "Bel_Outflow"
		entry2102.PostReference = "J" + strconv.Itoa(journalTransactions[122].AosStep)
		entry2102.Debit = journalTransactions[122].Debit
		entry2102.Credit = journalTransactions[122].Credit
		if journalTransactions[122].NormalAccountBalance == "Debit" {
			entry2102.Balance = -(entry2102.Credit - entry2102.Debit) + entry2101.Balance
		} else {
			entry2102.Balance = entry2102.Credit - entry2102.Debit + entry2101.Balance
		}

		subledgerEntries = append(subledgerEntries, entry2101, entry2102)

		//Risk Adjustment Buildup
		var entry2201, entry2202 models.SubLedgerReportEntry //22
		entry2201.ProductCode = journalTransactions[1].ProductCode
		entry2201.IFRS17Group = journalTransactions[1].IFRS17Group
		entry2201.LedgerName = "Risk_Adjustment_Buildup"
		entry2201.MasterAccountType = journalTransactions[94].MasterAccountType
		entry2201.AccountNumber = 2009
		entry2201.PostingDate = ""
		entry2201.ContraAccountTransactionDescription = "Bel_Inflow_Buildup"
		entry2201.PostReference = "J" + strconv.Itoa(journalTransactions[94].AosStep)
		entry2201.Debit = journalTransactions[94].Debit
		entry2201.Credit = journalTransactions[94].Credit
		if journalTransactions[94].NormalAccountBalance == "Debit" {
			entry2201.Balance = -(entry2201.Credit - entry2201.Debit)
		} else {
			entry2201.Balance = entry2201.Credit - entry2201.Debit
		}

		entry2202.ProductCode = journalTransactions[1].ProductCode
		entry2202.IFRS17Group = journalTransactions[1].IFRS17Group
		entry2202.LedgerName = "Risk_Adjustment_Buildup"
		entry2202.MasterAccountType = journalTransactions[124].MasterAccountType
		entry2202.AccountNumber = 2009
		entry2202.PostingDate = ""
		entry2202.ContraAccountTransactionDescription = "Risk_Adjustment"
		entry2202.PostReference = "J" + strconv.Itoa(journalTransactions[124].AosStep)
		entry2202.Debit = journalTransactions[124].Debit
		entry2202.Credit = journalTransactions[124].Credit
		if journalTransactions[124].NormalAccountBalance == "Debit" {
			entry2202.Balance = -(entry2202.Credit - entry2202.Debit) + entry2201.Balance
		} else {
			entry2202.Balance = entry2202.Credit - entry2202.Debit + entry2201.Balance
		}
		subledgerEntries = append(subledgerEntries, entry2201, entry2202)

		//Insurance Finance Income or Expense
		var entry2301, entry2302, entry2303, entry2304, entry2305, entry2306, entry2307, entry2308, entry2309, entry2310, entry2311, entry2312, entry2313 models.SubLedgerReportEntry //23

		entry2301.ProductCode = journalTransactions[1].ProductCode
		entry2301.IFRS17Group = journalTransactions[1].IFRS17Group
		entry2301.LedgerName = "Insurance_Finance_Income_or_Expense"
		entry2301.MasterAccountType = journalTransactions[11].MasterAccountType
		entry2301.AccountNumber = 420
		entry2301.PostingDate = ""
		entry2301.ContraAccountTransactionDescription = "BEL_Inflow"
		entry2301.PostReference = "J" + strconv.Itoa(journalTransactions[11].AosStep)
		entry2301.Debit = journalTransactions[11].Debit
		entry2301.Credit = journalTransactions[11].Credit
		if journalTransactions[11].NormalAccountBalance == "Debit" {
			entry2301.Balance = -(entry2301.Credit - entry2301.Debit)
		} else {
			entry2301.Balance = entry2301.Credit - entry2301.Debit
		}

		entry2302.ProductCode = journalTransactions[1].ProductCode
		entry2302.IFRS17Group = journalTransactions[1].IFRS17Group
		entry2302.LedgerName = "Insurance_Finance_Income_or_Expense"
		entry2302.MasterAccountType = journalTransactions[7].MasterAccountType
		entry2302.AccountNumber = 420
		entry2302.PostingDate = ""
		entry2302.ContraAccountTransactionDescription = "BEL_Outflow"
		entry2302.PostReference = "J" + strconv.Itoa(journalTransactions[7].AosStep)
		entry2302.Debit = journalTransactions[7].Debit
		entry2302.Credit = journalTransactions[7].Credit
		if journalTransactions[7].NormalAccountBalance == "Debit" {
			entry2302.Balance = -(entry2302.Credit - entry2302.Debit) + entry2301.Balance
		} else {
			entry2302.Balance = entry2302.Credit - entry2302.Debit + entry2301.Balance
		}

		entry2303.ProductCode = journalTransactions[1].ProductCode
		entry2303.IFRS17Group = journalTransactions[1].IFRS17Group
		entry2303.LedgerName = "Insurance_Finance_Income_or_Expense"
		entry2303.MasterAccountType = journalTransactions[9].MasterAccountType
		entry2303.AccountNumber = 420
		entry2303.PostingDate = ""
		entry2303.ContraAccountTransactionDescription = "Risk_adjustment"
		entry2303.PostReference = "J" + strconv.Itoa(journalTransactions[9].AosStep)
		entry2303.Debit = journalTransactions[9].Debit
		entry2303.Credit = journalTransactions[9].Credit
		if journalTransactions[9].NormalAccountBalance == "Debit" {
			entry2303.Balance = -(entry2303.Credit - entry2303.Debit) + entry2302.Balance
		} else {
			entry2303.Balance = entry2303.Credit - entry2303.Debit + entry2302.Balance
		}

		entry2304.ProductCode = journalTransactions[1].ProductCode
		entry2304.IFRS17Group = journalTransactions[1].IFRS17Group
		entry2304.LedgerName = "Insurance_Finance_Income_or_Expense"
		entry2304.MasterAccountType = journalTransactions[44].MasterAccountType
		entry2304.AccountNumber = 420
		entry2304.PostingDate = ""
		entry2304.ContraAccountTransactionDescription = "BEL_Inflow"
		entry2304.PostReference = "J" + strconv.Itoa(journalTransactions[44].AosStep)
		entry2304.Debit = journalTransactions[44].Debit
		entry2304.Credit = journalTransactions[44].Credit
		if journalTransactions[44].NormalAccountBalance == "Debit" {
			entry2304.Balance = -(entry2304.Credit - entry2304.Debit) + entry2303.Balance
		} else {
			entry2304.Balance = entry2304.Credit - entry2304.Debit + entry2303.Balance
		}

		entry2305.ProductCode = journalTransactions[1].ProductCode
		entry2305.IFRS17Group = journalTransactions[1].IFRS17Group
		entry2305.LedgerName = "Insurance_Finance_Income_or_Expense"
		entry2305.MasterAccountType = journalTransactions[48].MasterAccountType
		entry2305.AccountNumber = 420
		entry2305.PostingDate = ""
		entry2305.ContraAccountTransactionDescription = "BEL_Outflow"
		entry2305.PostReference = "J" + strconv.Itoa(journalTransactions[48].AosStep)
		entry2305.Debit = journalTransactions[48].Debit
		entry2305.Credit = journalTransactions[48].Credit
		if journalTransactions[48].NormalAccountBalance == "Debit" {
			entry2305.Balance = -(entry2305.Credit - entry2305.Debit) + entry2304.Balance
		} else {
			entry2305.Balance = entry2305.Credit - entry2305.Debit + entry2304.Balance
		}

		entry2306.ProductCode = journalTransactions[1].ProductCode
		entry2306.IFRS17Group = journalTransactions[1].IFRS17Group
		entry2306.LedgerName = "Insurance_Finance_Income_or_Expense"
		entry2306.MasterAccountType = journalTransactions[58].MasterAccountType
		entry2306.AccountNumber = 420
		entry2306.PostingDate = ""
		entry2306.ContraAccountTransactionDescription = "Risk_Adjustment"
		entry2306.PostReference = "J" + strconv.Itoa(journalTransactions[58].AosStep)
		entry2306.Debit = journalTransactions[58].Debit
		entry2306.Credit = journalTransactions[58].Credit
		if journalTransactions[58].NormalAccountBalance == "Debit" {
			entry2306.Balance = -(entry2306.Credit - entry2306.Debit) + entry2305.Balance
		} else {
			entry2306.Balance = entry2306.Credit - entry2306.Debit + entry2305.Balance
		}

		entry2307.ProductCode = journalTransactions[1].ProductCode
		entry2307.IFRS17Group = journalTransactions[1].IFRS17Group
		entry2307.LedgerName = "Insurance_Finance_Income_or_Expense"
		entry2307.MasterAccountType = journalTransactions[95].MasterAccountType
		entry2307.AccountNumber = 420
		entry2307.PostingDate = ""
		entry2307.ContraAccountTransactionDescription = "Bel_Outflow"
		entry2307.PostReference = "J" + strconv.Itoa(journalTransactions[95].AosStep)
		entry2307.Debit = journalTransactions[95].Debit
		entry2307.Credit = journalTransactions[95].Credit
		if journalTransactions[95].NormalAccountBalance == "Debit" {
			entry2307.Balance = -(entry2307.Credit - entry2307.Debit) + entry2306.Balance
		} else {
			entry2307.Balance = entry2307.Credit - entry2307.Debit + entry2306.Balance
		}

		entry2308.ProductCode = journalTransactions[1].ProductCode
		entry2308.IFRS17Group = journalTransactions[1].IFRS17Group
		entry2308.LedgerName = "Insurance_Finance_Income_or_Expense"
		entry2308.MasterAccountType = journalTransactions[99].MasterAccountType
		entry2308.AccountNumber = 420
		entry2308.PostingDate = ""
		entry2308.ContraAccountTransactionDescription = "BEL_Inflow"
		entry2308.PostReference = "J" + strconv.Itoa(journalTransactions[99].AosStep)
		entry2308.Debit = journalTransactions[99].Debit
		entry2308.Credit = journalTransactions[99].Credit
		if journalTransactions[99].NormalAccountBalance == "Debit" {
			entry2308.Balance = -(entry2308.Credit - entry2308.Debit) + entry2307.Balance
		} else {
			entry2308.Balance = entry2308.Credit - entry2308.Debit + entry2307.Balance
		}

		entry2309.ProductCode = journalTransactions[1].ProductCode
		entry2309.IFRS17Group = journalTransactions[1].IFRS17Group
		entry2309.LedgerName = "Insurance_Finance_Income_or_Expense"
		entry2309.MasterAccountType = journalTransactions[97].MasterAccountType
		entry2309.AccountNumber = 420
		entry2309.PostingDate = ""
		entry2309.ContraAccountTransactionDescription = "Risk_Adjustment"
		entry2309.PostReference = "J" + strconv.Itoa(journalTransactions[97].AosStep)
		entry2309.Debit = journalTransactions[97].Debit
		entry2309.Credit = journalTransactions[97].Credit
		if journalTransactions[97].NormalAccountBalance == "Debit" {
			entry2309.Balance = -(entry2309.Credit - entry2309.Debit) + entry2308.Balance
		} else {
			entry2309.Balance = entry2309.Credit - entry2309.Debit + entry2308.Balance
		}

		entry2310.ProductCode = journalTransactions[1].ProductCode
		entry2310.IFRS17Group = journalTransactions[1].IFRS17Group
		entry2310.LedgerName = "Insurance_Finance_Income_or_Expense"
		entry2310.MasterAccountType = journalTransactions[101].MasterAccountType
		entry2310.AccountNumber = 420
		entry2310.PostingDate = ""
		entry2310.ContraAccountTransactionDescription = "Bel_Outflow"
		entry2310.PostReference = "J" + strconv.Itoa(journalTransactions[101].AosStep)
		entry2310.Debit = journalTransactions[101].Debit
		entry2310.Credit = journalTransactions[101].Credit
		if journalTransactions[101].NormalAccountBalance == "Debit" {
			entry2310.Balance = -(entry2310.Credit - entry2310.Debit) + entry2309.Balance
		} else {
			entry2310.Balance = entry2310.Credit - entry2310.Debit + entry2309.Balance
		}

		entry2311.ProductCode = journalTransactions[1].ProductCode
		entry2311.IFRS17Group = journalTransactions[1].IFRS17Group
		entry2311.LedgerName = "Insurance_Finance_Income_or_Expense"
		entry2311.MasterAccountType = journalTransactions[105].MasterAccountType
		entry2311.AccountNumber = 420
		entry2311.PostingDate = ""
		entry2311.ContraAccountTransactionDescription = "Bel_Inflow"
		entry2311.PostReference = "J" + strconv.Itoa(journalTransactions[105].AosStep)
		entry2311.Debit = journalTransactions[105].Debit
		entry2311.Credit = journalTransactions[105].Credit
		if journalTransactions[105].NormalAccountBalance == "Debit" {
			entry2311.Balance = -(entry2311.Credit - entry2311.Debit) + entry2310.Balance
		} else {
			entry2311.Balance = entry2311.Credit - entry2311.Debit + entry2310.Balance
		}

		entry2312.ProductCode = journalTransactions[1].ProductCode
		entry2312.IFRS17Group = journalTransactions[1].IFRS17Group
		entry2312.LedgerName = "Insurance_Finance_Income_or_Expense"
		entry2312.MasterAccountType = journalTransactions[101].MasterAccountType
		entry2312.AccountNumber = 420
		entry2312.PostingDate = ""
		entry2312.ContraAccountTransactionDescription = "Risk_Adjustment"
		entry2312.PostReference = "J" + strconv.Itoa(journalTransactions[103].AosStep)
		entry2312.Debit = journalTransactions[103].Debit
		entry2312.Credit = journalTransactions[103].Credit
		if journalTransactions[103].NormalAccountBalance == "Debit" {
			entry2312.Balance = -(entry2312.Credit - entry2312.Debit) + entry2311.Balance
		} else {
			entry2312.Balance = entry2312.Credit - entry2312.Debit + entry2311.Balance
		}

		entry2313.ProductCode = journalTransactions[1].ProductCode
		entry2313.IFRS17Group = journalTransactions[1].IFRS17Group
		entry2313.LedgerName = "Insurance_Finance_Income_or_Expense"
		entry2313.MasterAccountType = journalTransactions[127].MasterAccountType
		entry2313.AccountNumber = 420
		entry2313.PostingDate = ""
		entry2313.ContraAccountTransactionDescription = "Cash_and_Cash_Equivalent"
		entry2313.PostReference = "J" + strconv.Itoa(journalTransactions[127].AosStep)
		entry2313.Debit = journalTransactions[127].Debit
		entry2313.Credit = journalTransactions[127].Credit
		if journalTransactions[127].NormalAccountBalance == "Debit" {
			entry2313.Balance = -(entry2313.Credit - entry2313.Debit) + entry2312.Balance
		} else {
			entry2313.Balance = entry2313.Credit - entry2313.Debit + entry2312.Balance
		}

		subledgerEntries = append(subledgerEntries, entry2301, entry2302, entry2303, entry2304, entry2305, entry2306, entry2307, entry2308, entry2309, entry2310, entry2311, entry2312, entry2313)

		//loss component adjustment
		var entry2401 models.SubLedgerReportEntry //23

		entry2401.ProductCode = journalTransactions[1].ProductCode
		entry2401.IFRS17Group = journalTransactions[1].IFRS17Group
		entry2401.LedgerName = "Loss_Component_Adjustment"
		entry2401.MasterAccountType = journalTransactions[90].MasterAccountType
		entry2401.AccountNumber = 202
		entry2401.PostingDate = ""
		entry2401.ContraAccountTransactionDescription = "BEL_Inflow_Buildup"
		entry2401.PostReference = "J" + strconv.Itoa(journalTransactions[90].AosStep)
		entry2401.Debit = journalTransactions[90].Debit
		entry2401.Credit = journalTransactions[90].Credit
		if journalTransactions[90].NormalAccountBalance == "Debit" {
			entry2401.Balance = -(entry2401.Credit - entry2401.Debit)
		} else {
			entry2401.Balance = entry2401.Credit - entry2401.Debit
		}

		subledgerEntries = append(subledgerEntries, entry2401)

		var entry2601, entry2602 models.SubLedgerReportEntry //23

		entry2601.ProductCode = journalTransactions[1].ProductCode
		entry2601.IFRS17Group = journalTransactions[1].IFRS17Group
		entry2601.LedgerName = "Loss_Component_NB"
		entry2601.MasterAccountType = journalTransactions[17].MasterAccountType
		entry2601.AccountNumber = 201
		entry2601.PostingDate = ""
		entry2601.ContraAccountTransactionDescription = "Bel_Outflow"
		entry2601.PostReference = "J" + strconv.Itoa(journalTransactions[18].AosStep)
		entry2601.Debit = journalTransactions[18].Credit
		entry2601.Credit = journalTransactions[18].Debit
		if journalTransactions[18].NormalAccountBalance == "Debit" {
			entry2601.Balance = -(entry2601.Credit - entry2601.Debit)
		} else {
			entry2601.Balance = entry2601.Credit - entry2601.Debit
		}

		entry2602.ProductCode = journalTransactions[1].ProductCode
		entry2602.IFRS17Group = journalTransactions[1].IFRS17Group
		entry2602.LedgerName = "Loss_Component_NB"
		entry2602.MasterAccountType = journalTransactions[17].MasterAccountType
		entry2602.AccountNumber = 201
		entry2602.PostingDate = ""
		entry2602.ContraAccountTransactionDescription = "Risk_Adjustment"
		entry2602.PostReference = "J" + strconv.Itoa(journalTransactions[19].AosStep)
		entry2602.Debit = journalTransactions[19].Credit
		entry2602.Credit = journalTransactions[19].Debit
		if journalTransactions[19].NormalAccountBalance == "Debit" {
			entry2602.Balance = -(entry2602.Credit - entry2602.Debit) + entry2601.Balance
		} else {
			entry2602.Balance = entry2602.Credit - entry2602.Debit + entry2601.Balance
		}

		subledgerEntries = append(subledgerEntries, entry2601, entry2602)

		var entry2801 models.SubLedgerReportEntry //23

		entry2801.ProductCode = journalTransactions[1].ProductCode
		entry2801.IFRS17Group = journalTransactions[1].IFRS17Group
		entry2801.LedgerName = "Experience_Premium_Variance"
		entry2801.MasterAccountType = journalTransactions[52].MasterAccountType
		entry2801.AccountNumber = 113
		entry2801.PostingDate = ""
		entry2801.ContraAccountTransactionDescription = "Premium_Debtors"
		entry2801.PostReference = "J" + strconv.Itoa(journalTransactions[52].AosStep)
		entry2801.Debit = journalTransactions[52].Debit
		entry2801.Credit = journalTransactions[52].Credit
		if journalTransactions[52].NormalAccountBalance == "Debit" {
			entry2801.Balance = -(entry2801.Credit - entry2801.Debit)
		} else {
			entry2801.Balance = entry2801.Credit - entry2801.Debit
		}
		subledgerEntries = append(subledgerEntries, entry2801)
	}

	if len(paaresults) > 0 {
		// PAA_Unearned_Premium
		var entry10001, entry10002, entry10003, entry20001, entry20002, entry20003, entry30001, entry40001, entry50001, entry60001, entry6000X1, entry70001, entry70002, entry70003, entry80001, entry80002, entry80003, entry80004, entry80005, entry80006, entry80007, entry80008, entry80009, entry90001, entry100001 models.SubLedgerReportEntry

		//opening balance
		entry10001.ProductCode = journalTransactions[1].ProductCode
		entry10001.IFRS17Group = journalTransactions[1].IFRS17Group
		entry10001.LedgerName = "PAA_Unearned_Premium"
		entry10001.MasterAccountType = journalTransactions[128].MasterAccountType
		entry10001.AccountNumber = 2010
		entry10001.PostingDate = ""
		entry10001.ContraAccountTransactionDescription = "PAA_Unearned_Premium"
		entry10001.PostReference = "J0"
		entry10001.Debit = 0
		entry10001.Credit = 0
		entry10001.Balance = 0

		entry10002.ProductCode = journalTransactions[1].ProductCode
		entry10002.IFRS17Group = journalTransactions[1].IFRS17Group
		entry10002.LedgerName = "PAA_Unearned_Premium"
		entry10002.MasterAccountType = journalTransactions[132].MasterAccountType
		entry10002.AccountNumber = 2010
		entry10002.PostingDate = ""
		entry10002.ContraAccountTransactionDescription = "Cash_and_Cash_Equivalent"
		entry10002.PostReference = "J" + strconv.Itoa(journalTransactions[132].AosStep)
		entry10002.Debit = journalTransactions[132].Debit
		entry10002.Credit = journalTransactions[132].Credit
		if journalTransactions[132].NormalAccountBalance == "Debit" {
			entry10002.Balance = -(entry10002.Credit - entry10002.Debit) + entry10001.Balance
		} else {
			entry10002.Balance = entry10002.Credit - entry10002.Debit + entry10001.Balance
		}

		entry10003.ProductCode = journalTransactions[1].ProductCode
		entry10003.IFRS17Group = journalTransactions[1].IFRS17Group
		entry10003.LedgerName = "PAA_Unearned_Premium"
		entry10003.MasterAccountType = journalTransactions[134].MasterAccountType
		entry10003.AccountNumber = 2010
		entry10003.PostingDate = ""
		entry10003.ContraAccountTransactionDescription = "PAA_Earned_Premium"
		entry10003.PostReference = "J" + strconv.Itoa(journalTransactions[134].AosStep)
		entry10003.Debit = journalTransactions[134].Debit
		entry10003.Credit = journalTransactions[134].Credit
		if journalTransactions[134].NormalAccountBalance == "Debit" {
			entry10003.Balance = -(entry10003.Credit - entry10003.Debit) + entry10002.Balance
		} else {
			entry10003.Balance = entry10003.Credit - entry10003.Debit + entry10002.Balance
		}

		//Acquisition Cashflows

		entry20001.ProductCode = journalTransactions[1].ProductCode
		entry20001.IFRS17Group = journalTransactions[1].IFRS17Group
		entry20001.LedgerName = "Acquisition_Cashflows"
		entry20001.MasterAccountType = journalTransactions[131].MasterAccountType
		entry20001.AccountNumber = 2007
		entry20001.PostingDate = ""
		entry20001.ContraAccountTransactionDescription = "Acquisition_Cashflows"
		entry20001.PostReference = "J0"
		entry20001.Debit = 0
		entry20001.Credit = 0
		entry20001.Balance = 0 //add upr from the buildup

		entry20002.ProductCode = journalTransactions[1].ProductCode
		entry20002.IFRS17Group = journalTransactions[1].IFRS17Group
		entry20002.LedgerName = "Acquisition_Cashflows"
		entry20002.MasterAccountType = journalTransactions[131].MasterAccountType
		entry20002.AccountNumber = 2007
		entry20002.PostingDate = ""
		entry20002.ContraAccountTransactionDescription = "Cash_and_Cash_Equivalent"
		entry20002.PostReference = "J" + strconv.Itoa(journalTransactions[131].AosStep)
		entry20002.Debit = journalTransactions[131].Debit
		entry20002.Credit = journalTransactions[131].Credit
		if journalTransactions[131].NormalAccountBalance == "Debit" {
			entry20002.Balance = -(entry20002.Credit - entry20002.Debit) + entry20001.Balance
		} else {
			entry20002.Balance = entry20002.Credit - entry20002.Debit + entry20001.Balance
		}

		//
		entry20003.ProductCode = journalTransactions[1].ProductCode
		entry20003.IFRS17Group = journalTransactions[1].IFRS17Group
		entry20003.LedgerName = "Acquisition_Cashflows"
		entry20003.MasterAccountType = journalTransactions[136].MasterAccountType
		entry20003.AccountNumber = 2007
		entry20003.PostingDate = ""
		entry20003.ContraAccountTransactionDescription = "Amortised_Acquisition_Cost"
		entry20003.PostReference = "J" + strconv.Itoa(journalTransactions[136].AosStep)
		entry20003.Debit = journalTransactions[136].Debit
		entry20003.Credit = journalTransactions[136].Credit
		if journalTransactions[136].NormalAccountBalance == "Debit" {
			entry20003.Balance = -(entry20003.Credit - entry20003.Debit) + entry20002.Balance
		} else {
			entry20003.Balance = entry20003.Credit - entry20003.Debit + entry20002.Balance
		}

		// PAA_Earned_Premium
		entry30001.ProductCode = journalTransactions[1].ProductCode
		entry30001.IFRS17Group = journalTransactions[1].IFRS17Group
		entry30001.LedgerName = "PAA_Earned_Premium"
		entry30001.MasterAccountType = journalTransactions[133].MasterAccountType
		entry30001.AccountNumber = 109
		entry30001.PostingDate = ""
		entry30001.ContraAccountTransactionDescription = "PAA_Unearned_Premium"
		entry30001.PostReference = "J" + strconv.Itoa(journalTransactions[133].AosStep)
		entry30001.Debit = journalTransactions[133].Debit
		entry30001.Credit = journalTransactions[133].Credit
		if journalTransactions[133].NormalAccountBalance == "Debit" {
			entry30001.Balance = -(entry30001.Credit - entry30001.Debit)
		} else {
			entry30001.Balance = entry30001.Credit - entry30001.Debit
		}

		//Acquisition Cost
		entry40001.ProductCode = journalTransactions[1].ProductCode
		entry40001.IFRS17Group = journalTransactions[1].IFRS17Group
		entry40001.LedgerName = "Amortised_Acquisition_Cost"
		entry40001.MasterAccountType = journalTransactions[135].MasterAccountType
		entry40001.AccountNumber = 108
		entry40001.PostingDate = ""
		entry40001.ContraAccountTransactionDescription = "Acquisition_Cashflows"
		entry40001.PostReference = "J" + strconv.Itoa(journalTransactions[135].AosStep)
		entry40001.Debit = journalTransactions[135].Debit
		entry40001.Credit = journalTransactions[135].Credit
		if journalTransactions[135].NormalAccountBalance == "Debit" {
			entry40001.Balance = -(entry40001.Credit - entry40001.Debit)
		} else {
			entry40001.Balance = entry40001.Credit - entry40001.Debit
		}

		// PAA Loss Component
		entry50001.ProductCode = journalTransactions[1].ProductCode
		entry50001.IFRS17Group = journalTransactions[1].IFRS17Group
		entry50001.LedgerName = "PAA_Loss_Component"
		entry50001.MasterAccountType = journalTransactions[137].MasterAccountType
		entry50001.AccountNumber = 209
		entry50001.PostingDate = ""
		entry50001.ContraAccountTransactionDescription = "Loss_Component_Reserve"
		entry50001.PostReference = "J" + strconv.Itoa(journalTransactions[137].AosStep)
		entry50001.Debit = journalTransactions[137].Debit
		entry50001.Credit = journalTransactions[137].Credit
		if journalTransactions[137].NormalAccountBalance == "Debit" {
			entry50001.Balance = -(entry50001.Credit - entry50001.Debit)
		} else {
			entry50001.Balance = entry50001.Credit - entry50001.Debit
		}

		// Loss_Component_Reserve
		entry60001.ProductCode = journalTransactions[1].ProductCode
		entry60001.IFRS17Group = journalTransactions[1].IFRS17Group
		entry60001.LedgerName = "Loss_Component_Reserve"
		entry60001.MasterAccountType = journalTransactions[138].MasterAccountType
		entry60001.AccountNumber = 2011
		entry60001.PostingDate = ""
		entry60001.ContraAccountTransactionDescription = "PAA_Loss_Component"
		entry60001.PostReference = "J" + strconv.Itoa(journalTransactions[138].AosStep)
		entry60001.Debit = journalTransactions[138].Debit
		entry60001.Credit = journalTransactions[138].Credit
		if journalTransactions[138].NormalAccountBalance == "Debit" {
			entry60001.Balance = -(entry60001.Credit - entry60001.Debit)
		} else {
			entry60001.Balance = entry60001.Credit - entry60001.Debit
		}

		//PAA Incurred Claims
		entry6000X1.ProductCode = journalTransactions[1].ProductCode
		entry6000X1.IFRS17Group = journalTransactions[1].IFRS17Group
		entry6000X1.LedgerName = "PAA_Incurred_Claims"
		entry6000X1.MasterAccountType = journalTransactions[139].MasterAccountType
		entry6000X1.AccountNumber = 210
		entry6000X1.PostingDate = ""
		entry6000X1.ContraAccountTransactionDescription = "OCR"
		entry6000X1.PostReference = "J" + strconv.Itoa(journalTransactions[139].AosStep)
		entry6000X1.Debit = journalTransactions[139].Debit
		entry6000X1.Credit = journalTransactions[139].Credit
		if journalTransactions[139].NormalAccountBalance == "Debit" {
			entry6000X1.Balance = -(entry6000X1.Credit - entry6000X1.Debit)
		} else {
			entry6000X1.Balance = entry6000X1.Credit - entry6000X1.Debit
		}

		// OCR
		entry70001.ProductCode = journalTransactions[1].ProductCode
		entry70001.IFRS17Group = journalTransactions[1].IFRS17Group
		entry70001.LedgerName = "OCR"
		entry70001.MasterAccountType = journalTransactions[140].MasterAccountType
		entry70001.AccountNumber = 2012
		entry70001.PostingDate = ""
		entry70001.ContraAccountTransactionDescription = "Cash_and_Cash_Equivalent"
		entry70001.PostReference = "J0"
		entry70001.Debit = 0
		entry70001.Credit = 0
		entry70001.Balance = 0

		entry70002.ProductCode = journalTransactions[1].ProductCode
		entry70002.IFRS17Group = journalTransactions[1].IFRS17Group
		entry70002.LedgerName = "OCR"
		entry70002.MasterAccountType = journalTransactions[140].MasterAccountType
		entry70002.AccountNumber = 2012
		entry70002.PostingDate = ""
		entry70002.ContraAccountTransactionDescription = "PAA_Incurred_Claims"
		entry70002.PostReference = "J" + strconv.Itoa(journalTransactions[140].AosStep)
		entry70002.Debit = journalTransactions[140].Debit
		entry70002.Credit = journalTransactions[140].Credit
		if journalTransactions[140].NormalAccountBalance == "Debit" {
			entry70002.Balance = -(entry70002.Credit - entry70002.Debit) + entry70001.Balance
		} else {
			entry70002.Balance = entry70002.Credit - entry70002.Debit + entry70001.Balance
		}

		entry70003.ProductCode = journalTransactions[1].ProductCode
		entry70003.IFRS17Group = journalTransactions[1].IFRS17Group
		entry70003.LedgerName = "OCR"
		entry70003.MasterAccountType = journalTransactions[142].MasterAccountType
		entry70003.AccountNumber = 2012
		entry70003.PostingDate = ""
		entry70003.ContraAccountTransactionDescription = "Cash_and_Cash_Equivalent"
		entry70003.PostReference = "J" + strconv.Itoa(journalTransactions[142].AosStep)
		entry70003.Debit = journalTransactions[142].Debit
		entry70003.Credit = journalTransactions[142].Credit
		if journalTransactions[142].NormalAccountBalance == "Debit" {
			entry70003.Balance = -(entry70003.Credit - entry70003.Debit) + entry70002.Balance
		} else {
			entry70003.Balance = entry70003.Credit - entry70003.Debit + entry70002.Balance
		}

		//entry70004.ProductCode = journalTransactions[1].ProductCode
		//entry70004.IFRS17Group = journalTransactions[1].IFRS17Group
		//entry70004.LedgerName = "OCR"
		//entry70004.MasterAccountType = journalTransactions[157].MasterAccountType
		//entry70004.AccountNumber = 2012
		//entry70004.ReportingDate = ""
		//entry70004.TransactionDescription = "Other_Claims"
		//entry70004.PostReference = "J" + strconv.Itoa(journalTransactions[157].AosStep)
		//entry70004.Debit = journalTransactions[157].Debit
		//entry70004.Credit = journalTransactions[157].Credit
		//if journalTransactions[157].NormalAccountBalance == "Debit" {
		//	entry70004.Balance = -(entry70004.Credit - entry70004.Debit) + entry70003.Balance
		//} else {
		//	entry70004.Balance = entry70004.Credit - entry70004.Debit + entry70003.Balance
		//}

		entry80001.ProductCode = journalTransactions[1].ProductCode
		entry80001.IFRS17Group = journalTransactions[1].IFRS17Group
		entry80001.LedgerName = "Cash_and_Cash_Equivalent"
		entry80001.MasterAccountType = journalTransactions[127].MasterAccountType
		entry80001.AccountNumber = 1005
		entry80001.PostingDate = ""
		entry80001.ContraAccountTransactionDescription = " Cash_and_Cash_Equivalent"
		entry80001.PostReference = "J0"
		entry80001.Credit = 0
		entry80001.Balance = 0

		entry80002.ProductCode = journalTransactions[1].ProductCode
		entry80002.IFRS17Group = journalTransactions[1].IFRS17Group
		entry80002.LedgerName = "Cash_and_Cash_Equivalent"
		entry80002.MasterAccountType = journalTransactions[129].MasterAccountType
		entry80002.AccountNumber = 1005
		entry80002.PostingDate = ""
		entry80002.ContraAccountTransactionDescription = "PAA_Unearned_Premium"
		entry80002.PostReference = "J" + strconv.Itoa(journalTransactions[129].AosStep)
		entry80002.Debit = journalTransactions[132].Credit
		entry80002.Credit = journalTransactions[132].Debit
		if journalTransactions[129].NormalAccountBalance == "Debit" {
			entry80002.Balance = -(entry80002.Credit - entry80002.Debit) + entry80001.Balance
		} else {
			entry80002.Balance = entry80002.Credit - entry80002.Debit + entry80001.Balance
		}

		entry80003.ProductCode = journalTransactions[1].ProductCode
		entry80003.IFRS17Group = journalTransactions[1].IFRS17Group
		entry80003.LedgerName = "Cash_and_Cash_Equivalent"
		entry80003.MasterAccountType = journalTransactions[129].MasterAccountType
		entry80003.AccountNumber = 1005
		entry80003.PostingDate = ""
		entry80003.ContraAccountTransactionDescription = "Acquisition_Cashflows"
		entry80003.PostReference = "J" + strconv.Itoa(journalTransactions[129].AosStep)
		entry80003.Debit = journalTransactions[131].Credit
		entry80003.Credit = journalTransactions[131].Debit
		if journalTransactions[129].NormalAccountBalance == "Debit" {
			entry80003.Balance = -(entry80003.Credit - entry80003.Debit) + entry80002.Balance
		} else {
			entry80003.Balance = entry80003.Credit - entry80003.Debit + entry80002.Balance
		}

		entry80004.ProductCode = journalTransactions[1].ProductCode
		entry80004.IFRS17Group = journalTransactions[1].IFRS17Group
		entry80004.LedgerName = "Cash_and_Cash_Equivalent"
		entry80004.MasterAccountType = journalTransactions[141].MasterAccountType
		entry80004.AccountNumber = 1005
		entry80004.PostingDate = ""
		entry80004.ContraAccountTransactionDescription = "OCR"
		entry80004.PostReference = "J" + strconv.Itoa(journalTransactions[141].AosStep)
		entry80004.Debit = journalTransactions[141].Debit
		entry80004.Credit = journalTransactions[141].Credit
		if journalTransactions[141].NormalAccountBalance == "Debit" {
			entry80004.Balance = -(entry80004.Credit - entry80004.Debit) + entry80003.Balance
		} else {
			entry80004.Balance = entry80004.Credit - entry80004.Debit + entry80003.Balance
		}

		entry80005.ProductCode = journalTransactions[1].ProductCode
		entry80005.IFRS17Group = journalTransactions[1].IFRS17Group
		entry80005.LedgerName = "Cash_and_Cash_Equivalent"
		entry80005.MasterAccountType = journalTransactions[161].MasterAccountType
		entry80005.AccountNumber = 1005
		entry80005.PostingDate = ""
		entry80005.ContraAccountTransactionDescription = "Non_Attributable_Expenses"
		entry80005.PostReference = "J" + strconv.Itoa(journalTransactions[161].AosStep)
		entry80005.Debit = journalTransactions[161].Debit
		entry80005.Credit = journalTransactions[161].Credit
		if journalTransactions[161].NormalAccountBalance == "Debit" {
			entry80005.Balance = -(entry80005.Credit - entry80005.Debit) + entry80004.Balance
		} else {
			entry80005.Balance = entry80005.Credit - entry80005.Debit + entry80004.Balance
		}

		entry80006.ProductCode = journalTransactions[1].ProductCode
		entry80006.IFRS17Group = journalTransactions[1].IFRS17Group
		entry80006.LedgerName = "Cash_and_Cash_Equivalent"
		entry80006.MasterAccountType = journalTransactions[164].MasterAccountType
		entry80006.AccountNumber = 1005
		entry80006.PostingDate = ""
		entry80006.ContraAccountTransactionDescription = "Actual_Expenses"
		entry80006.PostReference = "J" + strconv.Itoa(journalTransactions[164].AosStep)
		entry80006.Debit = journalTransactions[162].Credit
		entry80006.Credit = journalTransactions[162].Debit
		if journalTransactions[164].NormalAccountBalance == "Debit" {
			entry80006.Balance = -(entry80006.Credit - entry80006.Debit) + entry80005.Balance
		} else {
			entry80006.Balance = entry80006.Credit - entry80006.Debit + entry80005.Balance
		}

		entry80007.ProductCode = journalTransactions[1].ProductCode
		entry80007.IFRS17Group = journalTransactions[1].IFRS17Group
		entry80007.LedgerName = "Cash_and_Cash_Equivalent"
		entry80007.MasterAccountType = journalTransactions[42].MasterAccountType
		entry80007.AccountNumber = 1005
		entry80007.PostingDate = ""
		entry80007.ContraAccountTransactionDescription = "CSM"
		entry80007.PostReference = "J" + strconv.Itoa(journalTransactions[42].AosStep)
		entry80007.Debit = journalTransactions[42].Debit
		entry80007.Credit = journalTransactions[42].Credit
		if journalTransactions[42].NormalAccountBalance == "Debit" {
			entry80007.Balance = -(entry80007.Credit - entry80007.Debit) + entry80006.Balance
		} else {
			entry80007.Balance = entry80007.Credit - entry80007.Debit + entry80006.Balance
		}

		entry80008.ProductCode = journalTransactions[1].ProductCode
		entry80008.IFRS17Group = journalTransactions[1].IFRS17Group
		entry80008.LedgerName = "Cash_and_Cash_Equivalent"
		entry80008.MasterAccountType = journalTransactions[128].MasterAccountType
		entry80008.AccountNumber = 1005
		entry80008.PostingDate = ""
		entry80008.ContraAccountTransactionDescription = "Insurance_Finance_Income_or_Expense"
		entry80008.PostReference = "J" + strconv.Itoa(journalTransactions[128].AosStep)
		entry80008.Debit = journalTransactions[128].Debit
		entry80008.Credit = journalTransactions[128].Credit
		if journalTransactions[128].NormalAccountBalance == "Debit" {
			entry80008.Balance = -(entry80008.Credit - entry80008.Debit) + entry80007.Balance
		} else {
			entry80008.Balance = entry80008.Credit - entry80008.Debit + entry80007.Balance
		}

		entry80009.ProductCode = journalTransactions[1].ProductCode
		entry80009.IFRS17Group = journalTransactions[1].IFRS17Group
		entry80009.LedgerName = "Cash_and_Cash_Equivalent"
		entry80009.MasterAccountType = journalTransactions[55].MasterAccountType
		entry80009.AccountNumber = 1005
		entry80009.PostingDate = ""
		entry80009.ContraAccountTransactionDescription = "Premium_Debtors"
		entry80009.PostReference = "J" + strconv.Itoa(journalTransactions[55].AosStep)
		entry80009.Debit = journalTransactions[55].Debit
		entry80009.Credit = journalTransactions[55].Credit
		if journalTransactions[55].NormalAccountBalance == "Debit" {
			entry80009.Balance = -(entry80009.Credit - entry80009.Debit) + entry80008.Balance
		} else {
			entry80009.Balance = entry80009.Credit - entry80009.Debit + entry80008.Balance
		}

		entry90001.ProductCode = journalTransactions[1].ProductCode
		entry90001.IFRS17Group = journalTransactions[1].IFRS17Group
		entry90001.LedgerName = "Non_Attributable_Expenses"
		entry90001.MasterAccountType = journalTransactions[160].MasterAccountType
		entry90001.AccountNumber = 212
		entry90001.PostingDate = ""
		entry90001.ContraAccountTransactionDescription = "Cash_and_Cash_Equivalent"
		entry90001.PostReference = "J" + strconv.Itoa(journalTransactions[158].AosStep)
		entry90001.Debit = journalTransactions[160].Debit
		entry90001.Credit = journalTransactions[160].Credit
		if journalTransactions[160].NormalAccountBalance == "Debit" {
			entry90001.Balance = -(entry90001.Credit - entry90001.Debit)
		} else {
			entry90001.Balance = entry90001.Credit - entry90001.Debit
		}

		entry100001.ProductCode = journalTransactions[1].ProductCode
		entry100001.IFRS17Group = journalTransactions[1].IFRS17Group
		entry100001.LedgerName = "Actual_Expenses"
		entry100001.MasterAccountType = journalTransactions[162].MasterAccountType
		entry100001.AccountNumber = 206
		entry100001.PostingDate = ""
		entry100001.ContraAccountTransactionDescription = "Cash_and_Cash_Equivalent"
		entry100001.PostReference = "J" + strconv.Itoa(journalTransactions[162].AosStep)
		entry100001.Debit = journalTransactions[162].Debit
		entry100001.Credit = journalTransactions[162].Credit
		if journalTransactions[162].NormalAccountBalance == "Debit" {
			entry100001.Balance = -(entry100001.Credit - entry100001.Debit)
		} else {
			entry100001.Balance = entry100001.Credit - entry100001.Debit
		}

		var entry110001, entry110002, entry110003, entry110004, entry110005, entry110006, entry110007, entry110008, entry110009, entry110010 models.SubLedgerReportEntry
		entry110001.ProductCode = journalTransactions[1].ProductCode
		entry110001.IFRS17Group = journalTransactions[1].IFRS17Group
		entry110001.LedgerName = "PAA_Reinsurance_LRC"
		entry110001.MasterAccountType = journalTransactions[146].MasterAccountType
		entry110001.AccountNumber = 31001
		entry110001.PostingDate = ""
		entry110001.ContraAccountTransactionDescription = "B/F"
		entry110001.PostReference = "J" + strconv.Itoa(journalTransactions[146].AosStep)
		entry110001.Debit = 0
		entry110001.Credit = 0
		if journalTransactions[146].NormalAccountBalance == "Debit" {
			entry110001.Balance = -(entry110001.Credit - entry110001.Debit)
		} else {
			entry110001.Balance = entry110001.Credit - entry110001.Debit
		}

		entry110002.ProductCode = journalTransactions[1].ProductCode
		entry110002.IFRS17Group = journalTransactions[1].IFRS17Group
		entry110002.LedgerName = "PAA_Reinsurance_LRC"
		entry110002.MasterAccountType = journalTransactions[146].MasterAccountType
		entry110002.AccountNumber = 31001
		entry110002.PostingDate = ""
		entry110002.ContraAccountTransactionDescription = "Cash_and_Cash_Equivalent"
		entry110002.PostReference = "J" + strconv.Itoa(journalTransactions[146].AosStep)
		entry110002.Debit = journalTransactions[146].Debit
		entry110002.Credit = journalTransactions[146].Credit
		if journalTransactions[146].NormalAccountBalance == "Debit" {
			entry110002.Balance = -(entry110002.Credit - entry110002.Debit) + entry110001.Balance
		} else {
			entry110002.Balance = entry110002.Credit - entry110002.Debit + entry110001.Balance
		}

		entry110003.ProductCode = journalTransactions[1].ProductCode
		entry110003.IFRS17Group = journalTransactions[1].IFRS17Group
		entry110003.LedgerName = "PAA_Reinsurance_LRC"
		entry110003.MasterAccountType = journalTransactions[149].MasterAccountType
		entry110003.AccountNumber = 31001
		entry110003.PostingDate = ""
		entry110003.ContraAccountTransactionDescription = "PAA_Allocated_Reinsurance_Premium"
		entry110003.PostReference = "J" + strconv.Itoa(journalTransactions[149].AosStep)
		entry110003.Debit = journalTransactions[149].Debit
		entry110003.Credit = journalTransactions[149].Credit
		if journalTransactions[149].NormalAccountBalance == "Debit" {
			entry110003.Balance = -(entry110003.Credit - entry110003.Debit) + entry110002.Balance
		} else {
			entry110003.Balance = entry110003.Credit - entry110003.Debit + entry110002.Balance
		}

		entry110004.ProductCode = journalTransactions[1].ProductCode
		entry110004.IFRS17Group = journalTransactions[1].IFRS17Group
		entry110004.LedgerName = "PAA_Loss_Recovery_Asset"
		entry110004.MasterAccountType = journalTransactions[153].MasterAccountType
		entry110004.AccountNumber = 31001
		entry110004.PostingDate = ""
		entry110004.ContraAccountTransactionDescription = "B/F"
		entry110004.PostReference = "J" + strconv.Itoa(journalTransactions[153].AosStep)
		entry110004.Debit = 0
		entry110004.Credit = 0
		if journalTransactions[153].NormalAccountBalance == "Debit" {
			entry110004.Balance = -(entry110004.Credit - entry110004.Debit)
		} else {
			entry110004.Balance = entry110004.Credit - entry110004.Debit
		}

		entry110005.ProductCode = journalTransactions[1].ProductCode
		entry110005.IFRS17Group = journalTransactions[1].IFRS17Group
		entry110005.LedgerName = "PAA_Loss_Recovery_Asset"
		entry110005.MasterAccountType = journalTransactions[153].MasterAccountType
		entry110005.AccountNumber = 31002
		entry110005.PostingDate = ""
		entry110005.ContraAccountTransactionDescription = "Initial_Recognition_Loss_Recovery_Income"
		entry110005.PostReference = "J" + strconv.Itoa(journalTransactions[153].AosStep)
		entry110005.Debit = journalTransactions[153].Debit
		entry110005.Credit = journalTransactions[153].Credit
		if journalTransactions[153].NormalAccountBalance == "Debit" {
			entry110005.Balance = -(entry110005.Credit - entry110005.Debit) + entry110004.Balance
		} else {
			entry110005.Balance = entry110005.Credit - entry110005.Debit + entry110004.Balance
		}

		entry110006.ProductCode = journalTransactions[1].ProductCode
		entry110006.IFRS17Group = journalTransactions[1].IFRS17Group
		entry110006.LedgerName = "PAA_Loss_Recovery_Asset"
		entry110006.MasterAccountType = journalTransactions[155].MasterAccountType
		entry110006.AccountNumber = 31002
		entry110006.PostingDate = ""
		entry110006.ContraAccountTransactionDescription = "Loss_Recovery_Component_Income"
		entry110006.PostReference = "J" + strconv.Itoa(journalTransactions[155].AosStep)
		entry110006.Debit = journalTransactions[155].Debit
		entry110006.Credit = journalTransactions[155].Credit
		if journalTransactions[155].NormalAccountBalance == "Debit" {
			entry110006.Balance = -(entry110006.Credit - entry110006.Debit) + entry110005.Balance
		} else {
			entry110006.Balance = entry110006.Credit - entry110006.Debit + entry110005.Balance
		}

		entry110007.ProductCode = journalTransactions[1].ProductCode
		entry110007.IFRS17Group = journalTransactions[1].IFRS17Group
		entry110007.LedgerName = "PAA_Loss_Recovery_Asset"
		entry110007.MasterAccountType = journalTransactions[159].MasterAccountType
		entry110007.AccountNumber = 31002
		entry110007.PostingDate = ""
		entry110007.ContraAccountTransactionDescription = "Reversal_Loss_Recovery_Component"
		entry110007.PostReference = "J" + strconv.Itoa(journalTransactions[159].AosStep)
		entry110007.Debit = journalTransactions[159].Debit
		entry110007.Credit = journalTransactions[159].Credit
		if journalTransactions[159].NormalAccountBalance == "Debit" {
			entry110007.Balance = -(entry110007.Credit - entry110007.Debit) + entry110006.Balance
		} else {
			entry110007.Balance = entry110007.Credit - entry110007.Debit + entry110006.Balance
		}

		entry110008.ProductCode = journalTransactions[1].ProductCode
		entry110008.IFRS17Group = journalTransactions[1].IFRS17Group
		entry110008.LedgerName = "Initial_Recognition_Loss_Recovery_Income"
		entry110008.MasterAccountType = journalTransactions[152].MasterAccountType
		entry110008.AccountNumber = 3102
		entry110008.PostingDate = ""
		entry110008.ContraAccountTransactionDescription = "PAA_Loss_Recovery_Asset"
		entry110008.PostReference = "J" + strconv.Itoa(journalTransactions[152].AosStep)
		entry110008.Debit = journalTransactions[152].Debit
		entry110008.Credit = journalTransactions[152].Credit
		if journalTransactions[152].NormalAccountBalance == "Debit" {
			entry110008.Balance = -(entry110008.Credit - entry110008.Debit)
		} else {
			entry110008.Balance = entry110008.Credit - entry110008.Debit
		}

		entry110009.ProductCode = journalTransactions[1].ProductCode
		entry110009.IFRS17Group = journalTransactions[1].IFRS17Group
		entry110009.LedgerName = "Loss_Recovery_Component_Income"
		entry110009.MasterAccountType = journalTransactions[154].MasterAccountType
		entry110009.AccountNumber = 3103
		entry110009.PostingDate = ""
		entry110009.ContraAccountTransactionDescription = "PAA_Loss_Recovery_Asset"
		entry110009.PostReference = "J" + strconv.Itoa(journalTransactions[154].AosStep)
		entry110009.Debit = journalTransactions[154].Debit
		entry110009.Credit = journalTransactions[154].Credit
		if journalTransactions[154].NormalAccountBalance == "Debit" {
			entry110009.Balance = -(entry110009.Credit - entry110009.Debit)
		} else {
			entry110009.Balance = entry110009.Credit - entry110009.Debit
		}

		entry110010.ProductCode = journalTransactions[1].ProductCode
		entry110010.IFRS17Group = journalTransactions[1].IFRS17Group
		entry110010.LedgerName = "Reversal_Loss_Recovery_Component"
		entry110010.MasterAccountType = journalTransactions[158].MasterAccountType
		entry110010.AccountNumber = 3202
		entry110010.PostingDate = ""
		entry110010.ContraAccountTransactionDescription = "PAA_Loss_Recovery_Asset"
		entry110010.PostReference = "J" + strconv.Itoa(journalTransactions[158].AosStep)
		entry110010.Debit = journalTransactions[158].Debit
		entry110010.Credit = journalTransactions[158].Credit
		if journalTransactions[158].NormalAccountBalance == "Debit" {
			entry110010.Balance = -(entry110010.Credit - entry110010.Debit)
		} else {
			entry110010.Balance = entry110010.Credit - entry110010.Debit
		}

		subledgerEntries = append(subledgerEntries, entry10001, entry10002, entry10003, entry20001, entry20002, entry20003, entry30001, entry40001, entry50001, entry60001, entry6000X1, entry70001, entry70002, entry70003, entry80001, entry80002, entry80003, entry80004, entry80005, entry80006, entry80007, entry80008, entry80009, entry90001, entry100001)
		subledgerEntries = append(subledgerEntries, entry110001, entry110002, entry110003, entry110004, entry110005, entry110006, entry110007, entry110008, entry110009, entry110010)
	}

	var entryX101, entryX102, entryY101, entryY102, entryX201, entryX202, entryX203, entryX204, entryX205, entryX206 models.SubLedgerReportEntry
	entryX101.ProductCode = journalTransactions[1].ProductCode
	entryX101.IFRS17Group = journalTransactions[1].IFRS17Group
	entryX101.LedgerName = "LIC_FCF_Changes"
	entryX101.MasterAccountType = journalTransactions[163].MasterAccountType
	entryX101.AccountNumber = 215
	entryX101.PostingDate = ""
	entryX101.ContraAccountTransactionDescription = "IBNR_Reserve"
	entryX101.PostReference = "J" + strconv.Itoa(journalTransactions[163].AosStep)
	entryX101.Debit = journalTransactions[165].Debit
	entryX101.Credit = journalTransactions[165].Credit
	if journalTransactions[163].NormalAccountBalance == "Debit" {
		entryX101.Balance = -(entryX101.Credit - entryX101.Debit)
	} else {
		entryX101.Balance = entryX101.Credit - entryX101.Debit
	}

	entryX102.ProductCode = journalTransactions[1].ProductCode
	entryX102.IFRS17Group = journalTransactions[1].IFRS17Group
	entryX102.LedgerName = "LIC_FCF_Changes"
	entryX102.MasterAccountType = journalTransactions[166].MasterAccountType
	entryX102.AccountNumber = 215
	entryX102.PostingDate = ""
	entryX102.ContraAccountTransactionDescription = "LIC_Risk_Adjustment"
	entryX102.PostReference = "J" + strconv.Itoa(journalTransactions[166].AosStep)
	entryX102.Debit = journalTransactions[166].Debit
	entryX102.Credit = journalTransactions[166].Credit
	if journalTransactions[166].NormalAccountBalance == "Debit" {
		entryX102.Balance = -(entryX102.Credit - entryX102.Debit) + entryX101.Balance
	} else {
		entryX102.Balance = entryX102.Credit - entryX102.Debit + entryX101.Balance
	}

	entryY101.ProductCode = journalTransactions[1].ProductCode
	entryY101.IFRS17Group = journalTransactions[1].IFRS17Group
	entryY101.LedgerName = "LIC_Experience_Adjustments"
	entryY101.MasterAccountType = journalTransactions[168].MasterAccountType
	entryY101.AccountNumber = 216
	entryY101.PostingDate = ""
	entryY101.ContraAccountTransactionDescription = "IBNR_Reserve"
	entryY101.PostReference = "J" + strconv.Itoa(journalTransactions[168].AosStep)
	entryY101.Debit = journalTransactions[168].Debit
	entryY101.Credit = journalTransactions[168].Credit
	if journalTransactions[168].NormalAccountBalance == "Debit" {
		entryY101.Balance = -(entryY101.Credit - entryY101.Debit)
	} else {
		entryY101.Balance = entryY101.Credit - entryY101.Debit
	}

	entryY102.ProductCode = journalTransactions[1].ProductCode
	entryY102.IFRS17Group = journalTransactions[1].IFRS17Group
	entryY102.LedgerName = "LIC_Experience_Adjustments"
	entryY102.MasterAccountType = journalTransactions[170].MasterAccountType
	entryY102.AccountNumber = 216
	entryY102.PostingDate = ""
	entryY102.ContraAccountTransactionDescription = "LIC_Risk_Adjustment"
	entryY102.PostReference = "J" + strconv.Itoa(journalTransactions[170].AosStep)
	entryY102.Debit = journalTransactions[170].Debit
	entryY102.Credit = journalTransactions[170].Credit
	if journalTransactions[170].NormalAccountBalance == "Debit" {
		entryY102.Balance = -(entryY102.Credit - entryY102.Debit) + entryY101.Balance
	} else {
		entryY102.Balance = entryY102.Credit - entryY102.Debit + entryY101.Balance
	}

	entryX201.ProductCode = journalTransactions[1].ProductCode
	entryX201.IFRS17Group = journalTransactions[1].IFRS17Group
	entryX201.LedgerName = "IBNR_Reserve"
	entryX201.MasterAccountType = journalTransactions[165].MasterAccountType
	entryX201.AccountNumber = 42001
	entryX201.PostingDate = ""
	entryX201.ContraAccountTransactionDescription = "B/F"
	entryX201.PostReference = "J00"
	entryX201.Debit = 0
	entryX201.Credit = 0
	if journalTransactions[165].NormalAccountBalance == "Debit" {
		entryX201.Balance = -(entryX201.Credit - entryX201.Debit)
	} else {
		entryX201.Balance = entryX201.Credit - entryX201.Debit
	}

	entryX202.ProductCode = journalTransactions[1].ProductCode
	entryX202.IFRS17Group = journalTransactions[1].IFRS17Group
	entryX202.LedgerName = "IBNR_Reserve"
	entryX202.MasterAccountType = journalTransactions[165].MasterAccountType
	entryX202.AccountNumber = 42001
	entryX202.PostingDate = ""
	entryX202.ContraAccountTransactionDescription = "LIC_FCF_Changes"
	entryX202.PostReference = "J" + strconv.Itoa(journalTransactions[165].AosStep)
	entryX202.Debit = journalTransactions[165].Debit
	entryX202.Credit = journalTransactions[165].Credit
	if journalTransactions[165].NormalAccountBalance == "Debit" {
		entryX202.Balance = -(entryX202.Credit - entryX202.Debit) + entryX201.Balance
	} else {
		entryX202.Balance = entryX202.Credit - entryX202.Debit + entryX201.Balance
	}

	entryX203.ProductCode = journalTransactions[1].ProductCode
	entryX203.IFRS17Group = journalTransactions[1].IFRS17Group
	entryX203.LedgerName = "IBNR_Reserve"
	entryX203.MasterAccountType = journalTransactions[169].MasterAccountType
	entryX203.AccountNumber = 42001
	entryX203.PostingDate = ""
	entryX203.ContraAccountTransactionDescription = "LIC_Experience_Adjustments"
	entryX203.PostReference = "J" + strconv.Itoa(journalTransactions[169].AosStep)
	entryX203.Debit = journalTransactions[169].Debit
	entryX203.Credit = journalTransactions[169].Credit
	if journalTransactions[169].NormalAccountBalance == "Debit" {
		entryX203.Balance = -(entryX203.Credit - entryX203.Debit) + entryX202.Balance
	} else {
		entryX203.Balance = entryX203.Credit - entryX203.Debit + entryX202.Balance
	}

	entryX204.ProductCode = journalTransactions[1].ProductCode
	entryX204.IFRS17Group = journalTransactions[1].IFRS17Group
	entryX204.LedgerName = "LIC_Risk_Adjustment"
	entryX204.MasterAccountType = journalTransactions[167].MasterAccountType
	entryX204.AccountNumber = 42002
	entryX204.PostingDate = ""
	entryX204.ContraAccountTransactionDescription = "B/F"
	entryX204.PostReference = "J00"
	entryX204.Debit = 0
	entryX204.Credit = 0
	if journalTransactions[167].NormalAccountBalance == "Debit" {
		entryX204.Balance = -(entryX204.Credit - entryX204.Debit)
	} else {
		entryX204.Balance = entryX204.Credit - entryX204.Debit
	}

	entryX205.ProductCode = journalTransactions[1].ProductCode
	entryX205.IFRS17Group = journalTransactions[1].IFRS17Group
	entryX205.LedgerName = "LIC_Risk_Adjustment"
	entryX205.MasterAccountType = journalTransactions[167].MasterAccountType
	entryX205.AccountNumber = 42002
	entryX205.PostingDate = ""
	entryX205.ContraAccountTransactionDescription = "LIC_FCF_Changes"
	entryX205.PostReference = "J00"
	entryX205.Debit = journalTransactions[167].Debit
	entryX205.Credit = journalTransactions[167].Credit
	if journalTransactions[167].NormalAccountBalance == "Debit" {
		entryX205.Balance = -(entryX205.Credit - entryX205.Debit) + entryX204.Balance
	} else {
		entryX205.Balance = entryX205.Credit - entryX205.Debit + entryX204.Balance
	}

	entryX206.ProductCode = journalTransactions[1].ProductCode
	entryX206.IFRS17Group = journalTransactions[1].IFRS17Group
	entryX206.LedgerName = "LIC_Risk_Adjustment"
	entryX206.MasterAccountType = journalTransactions[171].MasterAccountType
	entryX206.AccountNumber = 42002
	entryX206.PostingDate = ""
	entryX206.ContraAccountTransactionDescription = "LIC_Experience_Adjustments"
	entryX206.PostReference = "J00"
	entryX206.Debit = journalTransactions[171].Debit
	entryX206.Credit = journalTransactions[171].Credit
	if journalTransactions[171].NormalAccountBalance == "Debit" {
		entryX206.Balance = -(entryX206.Credit - entryX206.Debit) + entryX205.Balance
	} else {
		entryX206.Balance = entryX206.Credit - entryX206.Debit + entryX205.Balance
	}

	subledgerEntries = append(subledgerEntries, entryX101, entryX102, entryY101, entryY102, entryX201, entryX202, entryX203, entryX204, entryX205, entryX206)

	return subledgerEntries
}

func GenerateTrialBalanceReportByGroup(subLedgers []models.SubLedgerReportEntry, results []models.AOSStepResult, paaresults []models.PAAResult) []models.TrialBalanceReportEntry {
	// The assumption we will work with is that the stepResults have been
	// filtered by product_code and ifrs17_group

	var trialbalance []models.TrialBalanceReportEntry
	//For now we calculate each entry manually...
	if len(results) > 0 {

		var entry1, entry1999, entry2, entry3, entry4, entry6, entry7, entry8, entry9, entry10, entry11, entry12999, entry12, entry13999, entry13, entry14, entry14999, entry15, entry16, entry17, entry18, entry19, entry20, entry21, entry22, entry23, entry24 models.TrialBalanceReportEntry

		entry1.ProductCode = subLedgers[1].ProductCode
		entry1.IFRS17Group = subLedgers[1].IFRS17Group
		entry1.PostingDate = ""
		entry1.AccountName = "Bel_Inflow"
		entry1.AccountNumber = 1001
		entry1.Debit = math.Max(AccountBalance(entry1.AccountNumber, subLedgers), 0)
		entry1.Credit = math.Max(-AccountBalance(entry1.AccountNumber, subLedgers), 0)
		entry1.Amount = entry1.Debit - entry1.Credit

		//..999 offset
		entry1999.ProductCode = subLedgers[1].ProductCode
		entry1999.IFRS17Group = subLedgers[1].IFRS17Group
		entry1999.PostingDate = ""
		entry1999.AccountName = "Bel_Inflow_B/F(Offset)"
		entry1999.AccountNumber = 1001999
		entry1999.Debit = math.Max(AccountBalance(entry1999.AccountNumber, subLedgers), 0)
		entry1999.Credit = math.Max(-AccountBalance(entry1999.AccountNumber, subLedgers), 0)
		entry1999.Amount = entry1999.Debit - entry1999.Credit

		entry2.ProductCode = subLedgers[1].ProductCode
		entry2.IFRS17Group = subLedgers[1].IFRS17Group
		entry2.PostingDate = ""
		entry2.AccountName = "DAC"
		entry2.AccountNumber = 1002
		entry2.Debit = math.Max(AccountBalance(entry2.AccountNumber, subLedgers), 0)
		entry2.Credit = math.Max(-AccountBalance(entry2.AccountNumber, subLedgers), 0)
		entry2.Amount = entry2.Debit - entry2.Credit

		entry3.ProductCode = subLedgers[1].ProductCode
		entry3.IFRS17Group = subLedgers[1].IFRS17Group
		entry3.PostingDate = ""
		entry3.AccountName = "Loss_Component_Adjustment"
		entry3.AccountNumber = 202
		entry3.Debit = math.Max(AccountBalance(entry3.AccountNumber, subLedgers), 0)
		entry3.Credit = math.Max(-AccountBalance(entry3.AccountNumber, subLedgers), 0)
		entry3.Amount = entry3.Debit - entry3.Credit

		entry4.ProductCode = subLedgers[1].ProductCode
		entry4.IFRS17Group = subLedgers[1].IFRS17Group
		entry4.PostingDate = ""
		entry4.AccountName = "Actual_Mortality_Claims"
		entry4.AccountNumber = 205
		entry4.Debit = math.Max(AccountBalance(entry4.AccountNumber, subLedgers), 0)
		entry4.Credit = math.Max(-AccountBalance(entry4.AccountNumber, subLedgers), 0)
		entry4.Amount = entry4.Debit - entry4.Credit

		//entry5.ProductCode = subLedgers[1].ProductCode
		//entry5.IFRS17Group = subLedgers[1].IFRS17Group
		//entry5.ReportingDate = ""
		//entry5.AccountName = "Actual_Expenses"
		//entry5.AccountNumber = 206
		//entry5.Debit = math.Max(AccountBalance(entry5.AccountNumber, subLedgers), 0)
		//entry5.Credit = math.Max(-AccountBalance(entry5.AccountNumber, subLedgers), 0)
		//entry5.Amount = entry5.Debit - entry5.Credit

		entry6.ProductCode = subLedgers[1].ProductCode
		entry6.IFRS17Group = subLedgers[1].IFRS17Group
		entry6.PostingDate = ""
		entry6.AccountName = "Actual_Retrenchment_Claims"
		entry6.AccountNumber = 207
		entry6.Debit = math.Max(AccountBalance(entry6.AccountNumber, subLedgers), 0)
		entry6.Credit = math.Max(-AccountBalance(entry6.AccountNumber, subLedgers), 0)
		entry6.Amount = entry6.Debit - entry6.Credit

		entry7.ProductCode = subLedgers[1].ProductCode
		entry7.IFRS17Group = subLedgers[1].IFRS17Group
		entry7.PostingDate = ""
		entry7.AccountName = "Actual_Morbidity_Claims"
		entry7.AccountNumber = 208
		entry7.Debit = math.Max(AccountBalance(entry7.AccountNumber, subLedgers), 0)
		entry7.Credit = math.Max(-AccountBalance(entry7.AccountNumber, subLedgers), 0)
		entry7.Amount = entry7.Debit - entry7.Credit

		entry8.ProductCode = subLedgers[1].ProductCode
		entry8.IFRS17Group = subLedgers[1].IFRS17Group
		entry8.PostingDate = ""
		entry8.AccountName = "DAC_Release"
		entry8.AccountNumber = 203
		entry8.Debit = math.Max(AccountBalance(entry8.AccountNumber, subLedgers), 0)
		entry8.Credit = math.Max(-AccountBalance(entry8.AccountNumber, subLedgers), 0)
		entry8.Amount = entry8.Debit - entry8.Credit

		entry9.ProductCode = subLedgers[1].ProductCode
		entry9.IFRS17Group = subLedgers[1].IFRS17Group
		entry9.PostingDate = ""
		entry9.AccountName = "Loss_Component_NB"
		entry9.AccountNumber = 201
		entry9.Debit = math.Max(AccountBalance(entry9.AccountNumber, subLedgers), 0)
		entry9.Credit = math.Max(-AccountBalance(entry9.AccountNumber, subLedgers), 0)
		entry9.Amount = entry9.Debit - entry9.Credit

		entry10.ProductCode = subLedgers[1].ProductCode
		entry10.IFRS17Group = subLedgers[1].IFRS17Group
		entry10.PostingDate = ""
		entry10.AccountName = "Cash_and_Cash_Equivalent"
		entry10.AccountNumber = 1005
		entry10.Debit = math.Max(AccountBalance(entry10.AccountNumber, subLedgers), 0)
		entry10.Credit = math.Max(-AccountBalance(entry10.AccountNumber, subLedgers), 0)
		entry10.Amount = entry10.Debit - entry10.Credit

		entry11.ProductCode = subLedgers[1].ProductCode
		entry11.IFRS17Group = subLedgers[1].IFRS17Group
		entry11.PostingDate = ""
		entry11.AccountName = "Insurance_Finance_Income_or_Expense"
		entry11.AccountNumber = 420
		entry11.Debit = math.Max(AccountBalance(entry11.AccountNumber, subLedgers), 0)
		entry11.Credit = math.Max(-AccountBalance(entry11.AccountNumber, subLedgers), 0)
		entry11.Amount = entry11.Debit - entry11.Credit

		entry12999.ProductCode = subLedgers[1].ProductCode
		entry12999.IFRS17Group = subLedgers[1].IFRS17Group
		entry12999.PostingDate = ""
		entry12999.AccountName = "BEL_Outflow_B/F(Offset)"
		entry12999.AccountNumber = 2001999
		entry12999.Debit = math.Max(AccountBalance(entry12999.AccountNumber, subLedgers), 0)
		entry12999.Credit = math.Max(-AccountBalance(entry12999.AccountNumber, subLedgers), 0)
		entry12999.Amount = entry12999.Debit - entry12999.Credit

		entry12.ProductCode = subLedgers[1].ProductCode
		entry12.IFRS17Group = subLedgers[1].IFRS17Group
		entry12.PostingDate = ""
		entry12.AccountName = "BEL_Outflow"
		entry12.AccountNumber = 2001
		entry12.Debit = math.Max(AccountBalance(entry12.AccountNumber, subLedgers), 0)
		entry12.Credit = math.Max(-AccountBalance(entry12.AccountNumber, subLedgers), 0)
		entry12.Amount = entry12.Debit - entry12.Credit

		entry13999.ProductCode = subLedgers[1].ProductCode
		entry13999.IFRS17Group = subLedgers[1].IFRS17Group
		entry13999.PostingDate = ""
		entry13999.AccountName = "Risk_Adjustment_B/F(Offset)"
		entry13999.AccountNumber = 2002999
		entry13999.Debit = math.Max(AccountBalance(entry13999.AccountNumber, subLedgers), 0)
		entry13999.Credit = math.Max(-AccountBalance(entry13999.AccountNumber, subLedgers), 0)
		entry13999.Amount = entry13999.Debit - entry13999.Credit

		entry13.ProductCode = subLedgers[1].ProductCode
		entry13.IFRS17Group = subLedgers[1].IFRS17Group
		entry13.PostingDate = ""
		entry13.AccountName = "Risk_Adjustment"
		entry13.AccountNumber = 2002
		entry13.Debit = math.Max(AccountBalance(entry13.AccountNumber, subLedgers), 0)
		entry13.Credit = math.Max(-AccountBalance(entry13.AccountNumber, subLedgers), 0)
		entry13.Amount = entry13.Debit - entry13.Credit

		entry14.ProductCode = subLedgers[1].ProductCode
		entry14.IFRS17Group = subLedgers[1].IFRS17Group
		entry14.PostingDate = ""
		entry14.AccountName = "CSM"
		entry14.AccountNumber = 2003
		entry14.Debit = math.Max(AccountBalance(entry14.AccountNumber, subLedgers), 0)
		entry14.Credit = math.Max(-AccountBalance(entry14.AccountNumber, subLedgers), 0)
		entry14.Amount = entry14.Debit - entry14.Credit

		entry14999.ProductCode = subLedgers[1].ProductCode
		entry14999.IFRS17Group = subLedgers[1].IFRS17Group
		entry14999.PostingDate = ""
		entry14999.AccountName = "CSM_B/F(Offset)"
		entry14999.AccountNumber = 2003999
		entry14999.Debit = math.Max(AccountBalance(entry14999.AccountNumber, subLedgers), 0)
		entry14999.Credit = math.Max(-AccountBalance(entry14999.AccountNumber, subLedgers), 0)
		entry14999.Amount = entry14999.Debit - entry14999.Credit

		entry15.ProductCode = subLedgers[1].ProductCode
		entry15.IFRS17Group = subLedgers[1].IFRS17Group
		entry15.PostingDate = ""
		entry15.AccountName = "Risk_Adjustment_Release"
		entry15.AccountNumber = 101
		entry15.Debit = math.Max(AccountBalance(entry15.AccountNumber, subLedgers), 0)
		entry15.Credit = math.Max(-AccountBalance(entry15.AccountNumber, subLedgers), 0)
		entry15.Amount = entry15.Debit - entry15.Credit

		entry16.ProductCode = subLedgers[1].ProductCode
		entry16.IFRS17Group = subLedgers[1].IFRS17Group
		entry16.PostingDate = ""
		entry16.AccountName = "Expected_Mortality_Claims"
		entry16.AccountNumber = 102
		entry16.Debit = math.Max(AccountBalance(entry16.AccountNumber, subLedgers), 0)
		entry16.Credit = math.Max(-AccountBalance(entry16.AccountNumber, subLedgers), 0)
		entry16.Amount = entry16.Debit - entry16.Credit

		entry17.ProductCode = subLedgers[1].ProductCode
		entry17.IFRS17Group = subLedgers[1].IFRS17Group
		entry17.PostingDate = ""
		entry17.AccountName = "Expected_Expenses"
		entry17.AccountNumber = 103
		entry17.Debit = math.Max(AccountBalance(entry17.AccountNumber, subLedgers), 0)
		entry17.Credit = math.Max(-AccountBalance(entry17.AccountNumber, subLedgers), 0)
		entry17.Amount = entry17.Debit - entry17.Credit

		entry18.ProductCode = subLedgers[1].ProductCode
		entry18.IFRS17Group = subLedgers[1].IFRS17Group
		entry18.PostingDate = ""
		entry18.AccountName = "Expected_Retrenchment"
		entry18.AccountNumber = 104
		entry18.Debit = math.Max(AccountBalance(entry18.AccountNumber, subLedgers), 0)
		entry18.Credit = math.Max(-AccountBalance(entry18.AccountNumber, subLedgers), 0)
		entry18.Amount = entry18.Debit - entry18.Credit

		entry19.ProductCode = subLedgers[1].ProductCode
		entry19.IFRS17Group = subLedgers[1].IFRS17Group
		entry19.PostingDate = ""
		entry19.AccountName = "Expected_Morbidity"
		entry19.AccountNumber = 105
		entry19.Debit = math.Max(AccountBalance(entry19.AccountNumber, subLedgers), 0)
		entry19.Credit = math.Max(-AccountBalance(entry19.AccountNumber, subLedgers), 0)
		entry19.Amount = entry19.Debit - entry19.Credit

		entry20.ProductCode = subLedgers[1].ProductCode
		entry20.IFRS17Group = subLedgers[1].IFRS17Group
		entry20.PostingDate = ""
		entry20.AccountName = "CSM_Release"
		entry20.AccountNumber = 106
		entry20.Debit = math.Max(AccountBalance(entry20.AccountNumber, subLedgers), 0)
		entry20.Credit = math.Max(-AccountBalance(entry20.AccountNumber, subLedgers), 0)
		entry20.Amount = entry20.Debit - entry20.Credit

		entry21.ProductCode = subLedgers[1].ProductCode
		entry21.IFRS17Group = subLedgers[1].IFRS17Group
		entry21.PostingDate = ""
		entry21.AccountName = "Loss_Component_Reversal"
		entry21.AccountNumber = 107
		entry21.Debit = math.Max(AccountBalance(entry21.AccountNumber, subLedgers), 0)
		entry21.Credit = math.Max(-AccountBalance(entry21.AccountNumber, subLedgers), 0)
		entry21.Amount = entry21.Debit - entry21.Credit

		entry22.ProductCode = subLedgers[1].ProductCode
		entry22.IFRS17Group = subLedgers[1].IFRS17Group
		entry22.PostingDate = ""
		entry22.AccountName = "Amortisation_of_Acquisition_Cashflows"
		entry22.AccountNumber = 108
		entry22.Debit = math.Max(AccountBalance(entry22.AccountNumber, subLedgers), 0)
		entry22.Credit = math.Max(-AccountBalance(entry22.AccountNumber, subLedgers), 0)
		entry22.Amount = entry22.Debit - entry22.Credit

		entry23.ProductCode = subLedgers[1].ProductCode
		entry23.IFRS17Group = subLedgers[1].IFRS17Group
		entry23.PostingDate = ""
		entry23.AccountName = "Investment_Income"
		entry23.AccountNumber = 410
		entry23.Debit = math.Max(AccountBalance(entry23.AccountNumber, subLedgers), 0)
		entry23.Credit = math.Max(-AccountBalance(entry23.AccountNumber, subLedgers), 0)
		entry23.Amount = entry23.Debit - entry23.Credit

		entry24.ProductCode = subLedgers[1].ProductCode
		entry24.IFRS17Group = subLedgers[1].IFRS17Group
		entry24.PostingDate = ""
		entry24.AccountName = "Acquisition_Cashflows"
		entry24.AccountNumber = 2007
		entry24.Debit = math.Max(AccountBalance(entry24.AccountNumber, subLedgers), 0)
		entry24.Credit = math.Max(-AccountBalance(entry24.AccountNumber, subLedgers), 0)
		entry24.Amount = entry24.Debit - entry24.Credit

		trialbalance = append(trialbalance, entry1, entry1999, entry2, entry3, entry4, entry6, entry7, entry8, entry9, entry10, entry11, entry12999, entry12, entry13999, entry13, entry14, entry14999, entry15, entry16, entry17, entry18, entry19, entry20, entry21, entry22, entry23, entry24)

	}
	if len(paaresults) > 0 {
		var entry25, entry26, entry27, entry28 models.TrialBalanceReportEntry

		entry25.ProductCode = subLedgers[1].ProductCode
		entry25.IFRS17Group = subLedgers[1].IFRS17Group
		entry25.PostingDate = ""
		entry25.AccountName = "PAA_Loss_Component"
		entry25.AccountNumber = 209
		entry25.Debit = math.Max(AccountBalance(entry25.AccountNumber, subLedgers), 0)
		entry25.Credit = math.Max(-AccountBalance(entry25.AccountNumber, subLedgers), 0)
		entry25.Amount = entry25.Debit - entry25.Credit

		entry26.ProductCode = subLedgers[1].ProductCode
		entry26.IFRS17Group = subLedgers[1].IFRS17Group
		entry26.PostingDate = ""
		entry26.AccountName = "PAA_Unearned_Premium"
		entry26.AccountNumber = 2010
		entry26.Debit = math.Max(AccountBalance(entry26.AccountNumber, subLedgers), 0)
		entry26.Credit = math.Max(-AccountBalance(entry26.AccountNumber, subLedgers), 0)
		entry26.Amount = entry26.Debit - entry26.Credit

		entry27.ProductCode = subLedgers[1].ProductCode
		entry27.IFRS17Group = subLedgers[1].IFRS17Group
		entry27.PostingDate = ""
		entry27.AccountName = "PAA_Earned_Premium"
		entry27.AccountNumber = 109
		entry27.Debit = math.Max(AccountBalance(entry27.AccountNumber, subLedgers), 0)
		entry27.Credit = math.Max(-AccountBalance(entry27.AccountNumber, subLedgers), 0)
		entry27.Amount = entry27.Debit - entry27.Credit

		entry28.ProductCode = subLedgers[1].ProductCode
		entry28.IFRS17Group = subLedgers[1].IFRS17Group
		entry28.PostingDate = ""
		entry28.AccountName = "Premium_Defficiency_Reserve"
		entry28.AccountNumber = 2011
		entry28.Debit = math.Max(AccountBalance(entry28.AccountNumber, subLedgers), 0)
		entry28.Credit = math.Max(-AccountBalance(entry28.AccountNumber, subLedgers), 0)
		entry28.Amount = entry28.Debit - entry28.Credit

		//entry29.ProductCode = subLedgers[1].ProductCode
		//entry29.IFRS17Group = subLedgers[1].IFRS17Group
		//entry29.ReportingDate = ""
		//entry29.AccountName = "Acquisition_Cashflows"
		//entry29.AccountNumber = 2007
		//entry29.Debit = math.Max(AccountBalance(entry29.AccountNumber, subLedgers), 0)
		//entry29.Credit = math.Max(-AccountBalance(entry29.AccountNumber, subLedgers), 0)
		//entry29.Amount = entry29.Debit - entry29.Credit

		//entry30.ProductCode = subLedgers[1].ProductCode
		//entry30.IFRS17Group = subLedgers[1].IFRS17Group
		//entry30.ReportingDate = ""
		//entry30.AccountName = "Amortised_Acquisition_Cost"
		//entry30.AccountNumber = 108
		//entry30.Debit = math.Max(AccountBalance(entry30.AccountNumber, subLedgers), 0)
		//entry30.Credit = math.Max(-AccountBalance(entry30.AccountNumber, subLedgers), 0)
		//entry30.Amount = entry30.Debit - entry30.Credit

		//entry31.ProductCode = subLedgers[1].ProductCode
		//entry31.IFRS17Group = subLedgers[1].IFRS17Group
		//entry31.ReportingDate = ""
		//entry31.AccountName = "Cash"
		//entry31.AccountNumber = 1005
		//entry31.Debit = math.Max(AccountBalance(entry31.AccountNumber, subLedgers), 0)
		//entry31.Credit = math.Max(-AccountBalance(entry31.AccountNumber, subLedgers), 0)
		//entry31.Amount = entry31.Debit - entry31.Credit

		trialbalance = append(trialbalance, entry25, entry26, entry27, entry28)

	}

	var entry32, entry33, entry34, entry35, entry36, entry37, entry38, entry39, entry40, entry41 models.TrialBalanceReportEntry

	entry32.ProductCode = subLedgers[1].ProductCode
	entry32.IFRS17Group = subLedgers[1].IFRS17Group
	entry32.PostingDate = ""
	entry32.AccountName = "OCR"
	entry32.AccountNumber = 2012
	entry32.Debit = math.Max(AccountBalance(entry32.AccountNumber, subLedgers), 0)
	entry32.Credit = math.Max(-AccountBalance(entry32.AccountNumber, subLedgers), 0)
	entry32.Amount = entry32.Debit - entry32.Credit

	entry33.ProductCode = subLedgers[1].ProductCode
	entry33.IFRS17Group = subLedgers[1].IFRS17Group
	entry33.PostingDate = ""
	entry33.AccountName = "Actual_Expenses"
	entry33.AccountNumber = 206
	entry33.Debit = math.Max(AccountBalance(entry33.AccountNumber, subLedgers), 0)
	entry33.Credit = math.Max(-AccountBalance(entry33.AccountNumber, subLedgers), 0)
	entry33.Amount = entry33.Debit - entry33.Credit

	entry34.ProductCode = subLedgers[1].ProductCode
	entry34.IFRS17Group = subLedgers[1].IFRS17Group
	entry34.PostingDate = ""
	entry34.AccountName = "Non_Attributable_Expenses"
	entry34.AccountNumber = 212
	entry34.Debit = math.Max(AccountBalance(entry34.AccountNumber, subLedgers), 0)
	entry34.Credit = math.Max(-AccountBalance(entry34.AccountNumber, subLedgers), 0)
	entry34.Amount = entry34.Debit - entry34.Credit

	entry35.ProductCode = subLedgers[1].ProductCode
	entry35.IFRS17Group = subLedgers[1].IFRS17Group
	entry35.PostingDate = ""
	entry35.AccountName = "PAA_Incurred_Claims"
	entry35.AccountNumber = 210
	entry35.Debit = math.Max(AccountBalance(entry35.AccountNumber, subLedgers), 0)
	entry35.Credit = math.Max(-AccountBalance(entry35.AccountNumber, subLedgers), 0)
	entry35.Amount = entry35.Debit - entry35.Credit

	entry36.ProductCode = subLedgers[1].ProductCode
	entry36.IFRS17Group = subLedgers[1].IFRS17Group
	entry36.PostingDate = ""
	entry36.AccountName = "IBNR_Reserve"
	entry36.AccountNumber = 42001
	entry36.Debit = math.Max(AccountBalance(entry36.AccountNumber, subLedgers), 0)
	entry36.Credit = math.Max(-AccountBalance(entry36.AccountNumber, subLedgers), 0)
	entry36.Amount = entry36.Debit - entry36.Credit

	entry37.ProductCode = subLedgers[1].ProductCode
	entry37.IFRS17Group = subLedgers[1].IFRS17Group
	entry37.PostingDate = ""
	entry37.AccountName = "LIC_Risk_Adjustment"
	entry37.AccountNumber = 42002
	entry37.Debit = math.Max(AccountBalance(entry37.AccountNumber, subLedgers), 0)
	entry37.Credit = math.Max(-AccountBalance(entry37.AccountNumber, subLedgers), 0)
	entry37.Amount = entry37.Debit - entry37.Credit

	entry38.ProductCode = subLedgers[1].ProductCode
	entry38.IFRS17Group = subLedgers[1].IFRS17Group
	entry38.PostingDate = ""
	entry38.AccountName = "LIC_FCF_Changes"
	entry38.AccountNumber = 215
	entry38.Debit = math.Max(AccountBalance(entry38.AccountNumber, subLedgers), 0)
	entry38.Credit = math.Max(-AccountBalance(entry38.AccountNumber, subLedgers), 0)
	entry38.Amount = entry38.Debit - entry38.Credit

	entry39.ProductCode = subLedgers[1].ProductCode
	entry39.IFRS17Group = subLedgers[1].IFRS17Group
	entry39.PostingDate = ""
	entry39.AccountName = "LIC_Experience_Adjustments"
	entry39.AccountNumber = 216
	entry39.Debit = math.Max(AccountBalance(entry39.AccountNumber, subLedgers), 0)
	entry39.Credit = math.Max(-AccountBalance(entry39.AccountNumber, subLedgers), 0)
	entry39.Amount = entry39.Debit - entry39.Credit

	entry40.ProductCode = subLedgers[1].ProductCode
	entry40.IFRS17Group = subLedgers[1].IFRS17Group
	entry40.PostingDate = ""
	entry40.AccountName = "Premium_Debtors"
	entry40.AccountNumber = 1006
	entry40.Debit = math.Max(AccountBalance(entry40.AccountNumber, subLedgers), 0)
	entry40.Credit = math.Max(-AccountBalance(entry40.AccountNumber, subLedgers), 0)
	entry40.Amount = entry40.Debit - entry40.Credit

	entry41.ProductCode = subLedgers[1].ProductCode
	entry41.IFRS17Group = subLedgers[1].IFRS17Group
	entry41.PostingDate = ""
	entry41.AccountName = "Experience_Premium_Variance"
	entry41.AccountNumber = 113
	entry41.Debit = math.Max(AccountBalance(entry41.AccountNumber, subLedgers), 0)
	entry41.Credit = math.Max(-AccountBalance(entry41.AccountNumber, subLedgers), 0)
	entry41.Amount = entry41.Debit - entry41.Credit

	trialbalance = append(trialbalance, entry32, entry33, entry34, entry35, entry36, entry37, entry38, entry39, entry40, entry41)

	return trialbalance
}

func populateEntryFields(step int, mmodel string, accountTitle string, accountNumber int, masterAccountType string, accountBalance string, accountType string, reportBundle int, entry *models.JTReportEntry) {
	entry.AosStep = step
	entry.MeasurementModel = mmodel
	entry.AccountTitle = accountTitle
	entry.AccountNumber = accountNumber
	entry.MasterAccountType = masterAccountType
	entry.NormalAccountBalance = accountBalance
	entry.AccountType = accountType
	entry.ReportBundle = reportBundle
}
func AccountBalance(account int, subledgers []models.SubLedgerReportEntry) float64 {
	var result float64

	for _, subledger := range subledgers {
		if subledger.AccountNumber == account {
			result += subledger.Debit - subledger.Credit
		}
	}
	return result
}

func GenerateBelBuildupforAll(runDate string) map[string]interface{} {
	var res = make(map[string]interface{})
	var belBuildup []models.BelBuildupVariable
	var prodList []models.ProductList
	results := GetAosStepResultsForAllProductsForDownstreamCalcs(runDate)
	if len(results) == 0 {
		mockResults := getMockResults()
		results = append(results, mockResults...)
	}

	belBuildup = GenerateBelBuildup(results)
	//prodList := GetJournalEntryProducts(runDate)
	DB.Raw("SELECT distinct product_code FROM aos_step_results where run_date = ?", runDate).Scan(&prodList)

	res["steps"] = belBuildup
	res["products"] = prodList

	return res
}

func GenerateBelBuildupOneProduct(runDate, prodCode string) map[string]interface{} {
	var res = make(map[string]interface{})
	var belBuildup []models.BelBuildupVariable
	results := GetAosStepResultsForOneProductForDownstreamCalcs(runDate, prodCode)
	if len(results) == 0 {
		mockResults := getMockResults()
		results = append(results, mockResults...)
	}

	belBuildup = GenerateBelBuildup(results)
	groupList := GetGroups(runDate, prodCode)

	res["steps"] = belBuildup
	res["groups"] = groupList

	return res
}

func GenerateBelBuildupOneGroup(ifrs17Group, productCode, runDate string) map[string]interface{} {
	var res = make(map[string]interface{})
	var results []models.AOSStepResult
	var belBuildup []models.BelBuildupVariable
	if ifrs17Group == "All Groups" {
		//run product
		results = GetAosStepResultsForOneProductForDownstreamCalcs(runDate, productCode)
		if len(results) == 0 {
			mockResults := getMockResults()
			results = append(results, mockResults...)
		}
	} else {
		results = GetAosStepResultsForGroup(ifrs17Group, runDate)
		if len(results) == 0 {
			mockResults := getMockResults()
			results = append(results, mockResults...)
		}
	}

	belBuildup = GenerateBelBuildup(results)
	res["steps"] = belBuildup
	return res
}

func GenerateBelBuildup(stepResults []models.AOSStepResult) []models.BelBuildupVariable {
	// The assumption we will work with is that the stepResults have been
	// filtered by product_code and ifrs17_group

	var reportEntries []models.BelBuildupVariable
	//For now we calculate each entry manually...

	if len(stepResults) > 0 {
		var raRelease float64
		var entry1, entry2, entry3, entry4, entry5, entry6, entry7, entry8, entry9, entry10, entry11, entry12, entry13, entry14, entry15, entry16, entry17 models.BelBuildupVariable
		for _, step := range stepResults {
			if step.Name == "B/F_Current" {
				entry1.Name = step.Name
				entry1.VariableChange = 0
				entry1.BelBuildup = step.BEL
				entry1.RAChange = 0
				entry1.RABuildup = step.RiskAdjustment
				entry1.CSMBuildup = step.CSMBuildup
				entry1.LossComponentBuildup = step.LossComponentBuildup
				entry1.Notes = "Opening Balance-Current"
				entry1.ReBelChange = 0
				entry1.ReBelBuildup = step.ReinsuranceBel
				entry1.ReRAChange = 0
				entry1.ReRABuildup = step.ReinsuranceRiskAdjustment
				entry1.ReCSMBuildup = step.ReinsCSMBuildup
			}
			if step.Name == "B/F_Lockedin" {
				entry2.Name = step.Name
				entry2.VariableChange = step.BestEstimateLiabilityChange
				entry2.BelBuildup = entry1.BelBuildup + entry2.VariableChange
				entry2.RAChange = step.RiskAdjustmentChange
				entry2.RABuildup = entry1.RABuildup + entry2.RAChange
				entry2.CSMBuildup = step.CSMBuildup
				entry2.LossComponentBuildup = step.LossComponentBuildup
				entry2.Notes = "Opening Balance-Lockedin"
				entry2.ReBelChange = step.ReinsBelChange
				entry2.ReBelBuildup = entry1.ReBelBuildup + entry2.ReBelChange
				entry2.ReRAChange = step.ReinsRAChange
				entry2.ReRABuildup = entry1.ReRABuildup + entry2.ReRAChange
				entry2.ReCSMBuildup = step.ReinsCSMBuildup
			}
			if step.Name == "Initial_Recog" {
				entry3.Name = step.Name
				entry3.VariableChange = step.BestEstimateLiabilityChange
				entry3.BelBuildup = entry2.BelBuildup + entry3.VariableChange
				entry3.RAChange = step.RiskAdjustmentChange
				entry3.RABuildup = entry2.RABuildup + entry3.RAChange
				entry3.CSMBuildup = entry2.CSMBuildup + step.CSMChange
				entry3.LossComponentBuildup = entry2.LossComponentBuildup + step.LossComponentChange
				entry3.Notes = "New Business"
				entry3.ReBelChange = step.ReinsBelChange
				entry3.ReBelBuildup = entry2.ReBelBuildup + entry3.ReBelChange
				entry3.ReRAChange = step.ReinsRAChange
				entry3.ReRABuildup = entry2.ReRABuildup + entry3.ReRAChange
				entry3.ReCSMBuildup = entry2.ReCSMBuildup + step.ReinsCSMChange
			}

			if step.Name == "Exp_Mort" {
				entry4.Name = step.Name
				entry4.VariableChange = -step.PNLChange
				entry4.BelBuildup = entry3.BelBuildup + entry4.VariableChange
				entry4.RAChange = step.RiskAdjustmentChange
				entry4.RABuildup = entry3.RABuildup + entry4.RAChange
				entry4.CSMBuildup = step.CSMBuildup
				entry4.LossComponentBuildup = step.LossComponentBuildup
				entry4.Notes = "Cash outflow Unwind"
				entry4.ReBelChange = 0
				entry4.ReBelBuildup = entry3.ReBelBuildup + entry4.ReBelChange
				entry4.ReRAChange = 0
				entry4.ReRABuildup = entry3.ReRABuildup + entry4.ReRAChange
				entry4.ReCSMBuildup = step.ReinsCSMBuildup
			}

			if step.Name == "Exp_Exp" {
				entry5.Name = step.Name
				entry5.VariableChange = -step.PNLChange
				entry5.BelBuildup = entry4.BelBuildup + entry5.VariableChange
				entry5.RAChange = step.RiskAdjustmentChange
				entry5.RABuildup = entry4.RABuildup + entry5.RAChange
				entry5.CSMBuildup = step.CSMBuildup
				entry5.LossComponentBuildup = step.LossComponentBuildup
				entry5.Notes = "Cash outflow Unwind"
				entry5.ReBelChange = 0
				entry5.ReBelBuildup = entry4.ReBelBuildup + entry5.ReBelChange
				entry5.ReRAChange = 0
				entry5.ReRABuildup = entry4.ReRABuildup + entry5.ReRAChange
				entry5.ReCSMBuildup = step.ReinsCSMBuildup
			}

			if step.Name == "Exp_Retrenchment" {
				entry6.Name = step.Name
				entry6.VariableChange = -step.PNLChange
				entry6.BelBuildup = entry5.BelBuildup + entry6.VariableChange
				entry6.RAChange = step.RiskAdjustmentChange
				entry6.RABuildup = entry5.RABuildup + entry6.RAChange
				entry6.CSMBuildup = step.CSMBuildup
				entry6.LossComponentBuildup = step.LossComponentBuildup
				entry6.Notes = "Cash outflow Unwind"
				entry6.ReBelChange = 0
				entry6.ReBelBuildup = entry5.ReBelBuildup + entry6.ReBelChange
				entry6.ReRAChange = 0
				entry6.ReRABuildup = entry5.ReRABuildup + entry6.ReRAChange
				entry6.ReCSMBuildup = step.ReinsCSMBuildup
			}

			if step.Name == "Exp_Morbidity" {
				entry7.Name = step.Name
				entry7.VariableChange = -step.PNLChange
				entry7.BelBuildup = entry6.BelBuildup + entry7.VariableChange
				entry7.RAChange = step.RiskAdjustmentChange
				entry7.RABuildup = entry6.RABuildup + entry7.RAChange
				entry7.CSMBuildup = step.CSMBuildup
				entry7.LossComponentBuildup = step.LossComponentBuildup
				entry7.Notes = "Cash outflow Unwind"
				entry7.ReBelChange = 0
				entry7.ReBelBuildup = entry6.ReBelBuildup + entry7.ReBelChange
				entry7.ReRAChange = 0
				entry7.ReRABuildup = entry6.ReRABuildup + entry7.ReRAChange
				entry7.ReCSMBuildup = step.ReinsCSMBuildup
			}

			if step.Name == "Exp_RA_Release" {
				raRelease = step.RiskAdjustmentChange
			}

			if step.Name == "Interest_Accretion" {
				entry8.Name = "CashInflow_Unwind"
				entry8.VariableChange = step.ExpectedCashInflow
				entry8.BelBuildup = entry7.BelBuildup + entry8.VariableChange
				entry8.RAChange = 0
				entry8.RABuildup = entry7.RABuildup + entry8.RAChange
				entry8.CSMBuildup = entry7.CSMBuildup
				entry8.LossComponentBuildup = entry7.LossComponentBuildup
				entry8.Notes = "Cash inflow unwind"
				entry8.ReBelChange = step.ReinsBelChange
				entry8.ReBelBuildup = entry7.ReBelBuildup - entry8.ReBelChange
				entry8.ReRAChange = step.ReinsRAChange
				entry8.ReRABuildup = entry7.ReRABuildup - entry8.ReRAChange
				entry8.ReCSMBuildup = step.ReinsCSMBuildup

				entry9.Name = "Other_CashOutflows"
				entry9.VariableChange = -(step.ExpectedCashOutflow + entry4.VariableChange + entry5.VariableChange + entry6.VariableChange + entry7.VariableChange)
				entry9.BelBuildup = entry8.BelBuildup + entry9.VariableChange
				entry9.RAChange = raRelease
				entry9.RABuildup = entry8.RABuildup + entry9.RAChange
				entry9.CSMBuildup = entry7.CSMBuildup
				entry9.LossComponentUnwind = step.LossComponentUnwind
				entry9.LossComponentBuildup = entry7.LossComponentBuildup
				entry9.Notes = "Cash outflow Unwind"
				entry9.ReBelChange = step.ReinsBelChange
				entry9.ReBelBuildup = entry8.ReBelBuildup - entry9.ReBelChange
				entry9.ReRAChange = step.ReinsRAChange
				entry9.ReRABuildup = entry8.ReRABuildup - entry9.ReRAChange
				entry9.ReCSMBuildup = step.ReinsCSMBuildup

				//entry3.Name = "Clawbacks"
				//entry3.VariableChange = step.ExpectedCashOutflow
				//entry3.BelBuildup += entry1.VariableChange

				entry10.Name = step.Name
				entry10.VariableChange = step.BestEstimateLiabilityChange //step.PNLChange
				entry10.BelBuildup = entry9.BelBuildup + entry10.VariableChange
				entry10.RAChange = step.RiskAdjustmentChange //step.RiskAdjustmentChange
				entry10.RABuildup = entry9.RABuildup + entry10.RAChange
				entry10.CSMChange = step.CSMChange
				entry10.CSMBuildup = step.CSMBuildup
				entry10.LossComponentUnwind = step.LossComponentUnwind
				entry10.LossComponentChange = step.LossComponentChange
				entry10.LossComponentBuildup = step.LossComponentBuildup
				entry10.Notes = "Interest accrual"
				entry10.ReBelChange = step.ReinsBelChange
				entry10.ReBelBuildup = entry9.ReBelBuildup - entry10.ReBelChange
				entry10.ReRAChange = step.ReinsRAChange
				entry10.ReRABuildup = entry9.ReRABuildup - entry10.ReRAChange
				entry10.ReCSMBuildup = step.ReinsCSMBuildup
			}

			if step.Name == "Data_Change" {
				entry11.Name = step.Name
				entry11.VariableChange = step.BestEstimateLiabilityChange
				entry11.BelBuildup = entry10.BelBuildup + entry11.VariableChange
				entry11.RAChange = step.RiskAdjustmentChange
				entry11.RABuildup = entry10.RABuildup + entry11.RAChange
				entry11.CSMChange = step.CSMBuildup - entry10.CSMBuildup
				entry11.CSMBuildup = step.CSMBuildup
				entry11.LossComponentBuildup = step.LossComponentBuildup
				entry11.LossComponentChange = step.LossComponentChange
				entry11.Notes = "New Data(IF+NB)"
				entry11.ReBelChange = step.ReinsBelChange
				entry11.ReBelBuildup = entry10.ReBelBuildup - entry11.ReBelChange
				entry11.ReRAChange = step.ReinsRAChange
				entry11.ReRABuildup = entry10.ReRABuildup - entry11.ReRAChange
				entry11.ReCSMBuildup = step.ReinsCSMBuildup
			}

			if step.Name == "NFA_Mort" {
				entry12.Name = step.Name
				entry12.VariableChange = step.BestEstimateLiabilityChange
				entry12.BelBuildup = entry11.BelBuildup + entry12.VariableChange
				entry12.RAChange = step.RiskAdjustmentChange
				entry12.RABuildup = entry11.RABuildup + entry12.RAChange
				entry12.CSMBuildup = step.CSMBuildup
				entry12.CSMChange = step.CSMBuildup - entry11.CSMBuildup
				entry12.LossComponentBuildup = step.LossComponentBuildup
				entry12.LossComponentChange = step.LossComponentChange
				entry12.Notes = "Assumption Change"
				entry12.ReBelChange = step.ReinsBelChange
				entry12.ReBelBuildup = entry11.ReBelBuildup - entry12.ReBelChange
				entry12.ReRAChange = step.ReinsRAChange
				entry12.ReRABuildup = entry11.ReRABuildup - entry12.ReRAChange
				entry12.ReCSMBuildup = step.ReinsCSMBuildup

			}

			if step.Name == "NFA_Exp" {
				entry13.Name = step.Name
				entry13.VariableChange = step.BestEstimateLiabilityChange
				entry13.BelBuildup = entry12.BelBuildup + entry13.VariableChange
				entry13.RAChange = step.RiskAdjustmentChange
				entry13.RABuildup = entry12.RABuildup + entry13.RAChange
				entry13.CSMBuildup = step.CSMBuildup
				entry13.CSMChange = step.CSMBuildup - entry12.CSMBuildup
				entry13.LossComponentBuildup = step.LossComponentBuildup
				entry13.LossComponentChange = step.LossComponentChange
				entry13.Notes = "Assumption Change"
				entry13.ReBelChange = step.ReinsBelChange
				entry13.ReBelBuildup = entry12.ReBelBuildup - entry13.ReBelChange
				entry13.ReRAChange = step.ReinsRAChange
				entry13.ReRABuildup = entry12.ReRABuildup - entry13.ReRAChange
				entry13.ReCSMBuildup = step.ReinsCSMBuildup
			}

			if step.Name == "NFA_Retrenchment" {
				entry14.Name = step.Name
				entry14.VariableChange = step.BestEstimateLiabilityChange
				entry14.BelBuildup = entry13.BelBuildup + entry14.VariableChange
				entry14.RAChange = step.RiskAdjustmentChange
				entry14.RABuildup = entry13.RABuildup + entry14.RAChange
				entry14.CSMBuildup = step.CSMBuildup
				entry14.CSMChange = step.CSMBuildup - entry13.CSMBuildup
				entry14.LossComponentBuildup = step.LossComponentBuildup
				entry14.LossComponentChange = step.LossComponentChange
				entry14.Notes = "Assumption Change"
				entry14.ReBelChange = step.ReinsBelChange
				entry14.ReBelBuildup = entry13.ReBelBuildup - entry14.ReBelChange
				entry14.ReRAChange = step.ReinsRAChange
				entry14.ReRABuildup = entry13.ReRABuildup - entry14.ReRAChange
				entry14.ReCSMBuildup = step.ReinsCSMBuildup
			}

			if step.Name == "NFA_Morbidity" {
				entry15.Name = step.Name
				entry15.VariableChange = step.BestEstimateLiabilityChange
				entry15.BelBuildup = entry14.BelBuildup + entry15.VariableChange
				entry15.RAChange = step.RiskAdjustmentChange
				entry15.RABuildup = entry14.RABuildup + entry15.RAChange
				entry15.CSMBuildup = step.CSMBuildup
				entry15.CSMChange = step.CSMBuildup - entry14.CSMBuildup
				entry15.LossComponentBuildup = step.LossComponentBuildup
				entry15.LossComponentChange = step.LossComponentChange
				entry15.Notes = "Assumption Change"
				entry15.ReBelChange = step.ReinsBelChange
				entry15.ReBelBuildup = entry14.ReBelBuildup - entry15.ReBelChange
				entry15.ReRAChange = step.ReinsRAChange
				entry15.ReRABuildup = entry14.ReRABuildup - entry15.ReRAChange
				entry15.ReCSMBuildup = step.ReinsCSMBuildup

			}

			if step.Name == "FA" {
				entry16.Name = step.Name
				entry16.VariableChange = step.BestEstimateLiabilityChange
				entry16.BelBuildup = entry15.BelBuildup + entry16.VariableChange
				entry16.RAChange = step.RiskAdjustmentChange
				entry16.RABuildup = entry15.RABuildup + entry16.RAChange
				entry16.CSMBuildup = step.CSMBuildup
				entry16.CSMChange = step.CSMBuildup - entry15.CSMBuildup
				entry16.LossComponentBuildup = step.LossComponentBuildup
				entry16.LossComponentChange = step.LossComponentChange
				entry16.Notes = "Assumption Change"
				entry16.ReBelChange = step.ReinsBelChange
				entry16.ReBelBuildup = entry15.ReBelBuildup - entry16.ReBelChange
				entry16.ReRAChange = step.ReinsRAChange
				entry16.ReRABuildup = entry15.ReRABuildup - entry16.ReRAChange
				entry16.ReCSMBuildup = step.ReinsCSMBuildup
			}

			if step.Name == "C/F_Current" {
				entry17.Name = step.Name
				entry17.VariableChange = step.BestEstimateLiabilityChange
				entry17.BelBuildup = entry16.BelBuildup + entry17.VariableChange
				entry17.RAChange = step.RiskAdjustmentChange
				entry17.RABuildup = entry16.RABuildup + entry17.RAChange
				entry17.CSMBuildup = step.CSMBuildup
				entry17.CSMChange = step.CSMBuildup - entry16.CSMBuildup
				entry17.LossComponentBuildup = step.LossComponentBuildup
				entry17.LossComponentChange = step.LossComponentChange
				entry17.Notes = "Current Rates & CSM Release"
				entry17.ReBelChange = step.ReinsBelChange
				entry17.ReBelBuildup = entry16.ReBelBuildup + entry17.ReBelChange
				entry17.ReRAChange = step.ReinsRAChange
				entry17.ReRABuildup = entry16.ReRABuildup + entry17.ReRAChange
				entry17.ReCSMBuildup = step.ReinsCSMBuildup
			}

		}
		reportEntries = append(reportEntries, entry1, entry2, entry3, entry4, entry5, entry6, entry7, entry8, entry9, entry10, entry11, entry12, entry13, entry14, entry15, entry16, entry17)
	} else {
		//populateEntryFields(2, "GMM", "BEL_Inflow", 1001, "B/S", "Debit", "Asset", 1, &entry1)
	}
	//sort.Slice(reportEntries, func(i, j int) bool {
	//	var sortedByReportBundle bool
	//	sortedByReportBundle = reportEntries[i].ReportBundle < reportEntries[j].ReportBundle
	//	return sortedByReportBundle
	//})
	return reportEntries
}
