package services

import (
	"api/models"
	"api/utils"
	"math"
	"strings"
)

// SumAssured  reads sum assured from the model point file and projects the sum assured using the sum assured escalation rate
func SumAssured(projection *models.Projection, p *models.Projection, mp models.ProductModelPoint, parameters models.ProductParameters, features models.ProductFeatures) {
	if projection.ValuationTimeMonth <= parameters.CalculatedTerm {
		if mp.MemberType == "CH" {
			if projection.AgeNextBirthday <= parameters.ChildExitAge {
				projection.SumAssured = GetChildSumAssured(projection.ProductCode, strings.ToUpper(mp.Plan), projection.AgeNextBirthday)
			} else {
				projection.SumAssured = 0
			}

		} else {
			if features.SaOutstandingLoan {
				projection.SumAssured = projection.OutstandingSumAssured
			} else if features.SaFixedBaseLumpSum {
				projection.SumAssured = utils.FloatPrecision(mp.SumAssured*projection.SumAssuredEscalation, defaultPrecision)
			} else {
				projection.SumAssured = 0
			}
		}
	} else {
		projection.SumAssured = 0
	}
}

// OutstandingSumAssured computes outstanding loan amount at each projection period using two methodologies
// 1st methodology calculates outstanding sum assured as the present value of instalments over the projection period
// 2nd methodology computes outstanding loan amount off the outstanding sum assured after interest less the instalment amount
// Credit Life
func OutstandingSumAssured(projection *models.Projection, p *models.Projection, mp models.ProductModelPoint, params models.ProductParameters, features models.ProductFeatures) {
	if projection.ValuationTimeMonth <= params.CalculatedTerm {
		if features.OsProjPvMethod {
			projection.OutstandingSumAssured = utils.FloatPrecision(CalculatePV(mp.Interest, math.Max(float64(mp.OutstandingTermMonths-projection.ValuationTimeMonth+1), 0), projection.CalculatedInstalment), defaultPrecision)
		} else {
			if projection.ProjectionMonth == 0 {
				projection.OutstandingSumAssured = math.Max(mp.OutstandingLoan, 0)
			} else {
				projection.OutstandingSumAssured = utils.FloatPrecision(math.Max(p.OutstandingSumAssured*(math.Pow(1.0+mp.Interest, 1.0/12.0))-projection.CalculatedInstalment, 0), defaultPrecision)
			}
		}
	} else {
		projection.OutstandingSumAssured = 0
	}
}

// AdditionalSumAssured reads additional sum assured from the funeral service sum assured table over the projection period
// It reads from child additional sum assured table for member type Child
func RiderSumAssured(projection *models.Projection, modelPoint models.ProductModelPoint, params models.ProductParameters, features models.ProductFeatures) {
	if projection.ValuationTimeMonth <= params.CalculatedTerm {
		if modelPoint.MemberType == "CH" {
			if projection.AgeNextBirthday <= params.ChildExitAge && projection.ValuationTimeMonth <= params.CalculatedTerm {
				projection.RiderSumAssured = GetChildAdditionalSumAssured(projection.ProductCode, modelPoint.Plan, projection.AgeNextBirthday)
			} else {
				projection.RiderSumAssured = 0
			}
		} else {
			if features.CreditLifeFuneralOption {
				projection.RiderSumAssured = params.RiderSumAssured + math.Min(math.Max(projection.SumAssured*params.RiderFuneralSaProportion, params.RiderFuneralMinimumBenefit), params.RiderFuneralMaximumBenefit)
			} else {
				projection.RiderSumAssured = params.RiderSumAssured
			}
		}
	} else {
		projection.RiderSumAssured = 0
	}
}

func StandardAdditionalLumpSum(projection *models.Projection, modelPoint models.ProductModelPoint, params models.ProductParameters) {
	if projection.ValuationTimeMonth <= params.CalculatedTerm {
		if modelPoint.MemberType == "CH" {
			if params.ChildExitAge >= projection.AgeNextBirthday && projection.ValuationTimeMonth <= params.CalculatedTerm {
				projection.StandardAdditionalLumpSum = params.StandardAdditionalSumAssured
			} else {
				projection.StandardAdditionalLumpSum = 0
			}
		} else {
			projection.StandardAdditionalLumpSum = params.StandardAdditionalSumAssured
		}
	} else {
		projection.StandardAdditionalLumpSum = 0
	}
}

func AnnuityIncome(projection *models.Projection, modelPoint models.ProductModelPoint, params models.ProductParameters) {
	if projection.ValuationTimeMonth <= params.CalculatedTerm {
		projection.AnnuityIncome = modelPoint.AnnuityIncome * projection.AnnuityEscalation
	} else {
		projection.AnnuityIncome = 0
	}
}

// Premium projects premium over the projection period allowing for premium escalation
func Premium(projection *models.Projection, modelPoint models.ProductModelPoint, features models.ProductFeatures, parameters models.ProductParameters) {
	if projection.ValuationTimeMonth <= parameters.CalculatedTerm && modelPoint.PremiumFrequency != 0 && projection.ProjectionMonth > 0 {
		if math.Mod(float64(projection.ProjectionMonth+11)*float64(modelPoint.PremiumFrequency), 12) == 0 {
			if features.CreditLifeFlatPremium {
				projection.Premium = utils.FloatPrecision(modelPoint.AnnualPremium*projection.PremiumEscalation/float64(modelPoint.PremiumFrequency), defaultPrecision)
			} else if features.CreditLifeDecreasingPremium {
				projection.Premium = projection.OutstandingSumAssured * modelPoint.PremiumRate * (12.0 / 1000.0) * (projection.PremiumEscalation / float64(modelPoint.PremiumFrequency))
			} else if features.NonLife {
				if modelPoint.Term > 0 {
					projection.Premium = utils.FloatPrecision((modelPoint.AnnualPremium*12.0/float64(modelPoint.Term))/float64(modelPoint.PremiumFrequency), defaultPrecision)
				}
			} else {
				projection.Premium = utils.FloatPrecision(modelPoint.AnnualPremium*projection.PremiumEscalation/float64(modelPoint.PremiumFrequency), defaultPrecision)
			}
		} else {
			projection.Premium = 0
		}
	} else if modelPoint.PremiumFrequency == 0 && projection.ValuationTimeMonth == 1 {
		projection.Premium = utils.FloatPrecision(modelPoint.AnnualPremium, defaultPrecision)
	} else {
		projection.Premium = 0
	}
}

// PremiumIncome computes premium income at each projection period
func PremiumIncome(projection *models.Projection, p *models.Projection, mp models.ProductModelPoint, parameters models.ProductParameters) {
	if projection.ValuationTimeMonth > parameters.CalculatedTerm || projection.ProjectionMonth == 0 {
		projection.PremiumIncome = 0
		projection.PremiumIncomeAdjusted = 0
	} else {
		projection.PremiumIncome = utils.FloatPrecision(projection.Premium*p.InitialPolicy, defaultPrecision)
		projection.PremiumIncomeAdjusted = utils.FloatPrecision(projection.Premium*p.InitialPolicyAdjusted, defaultPrecision)
	}
}

// PremiumNotReceived computes premium not receipted either as a result of a grace period or premium holiday benefit
func PremiumNotReceived(projection *models.Projection, features models.ProductFeatures, parameters models.ProductParameters, modelPoint models.ProductModelPoint) {
	if features.PremiumHoliday && parameters.PremiumHolidayCycle != 0 && projection.ValuationTimeMonth > parameters.PremiumHolidayWaitingPeriod && projection.ValuationTimeMonth <= parameters.CalculatedTerm {
		projection.PremiumNotReceived = utils.FloatPrecision(math.Min(math.Max(float64(int(projection.ValuationTimeMonth/parameters.PremiumHolidayCycle)-modelPoint.PremiumHolidayUsed), 0), float64(parameters.MaximumPremiumHolidays))*projection.PremiumIncome, defaultPrecision)
		projection.PremiumNotReceivedAdjusted = utils.FloatPrecision(math.Min(math.Max(float64(int(projection.ValuationTimeMonth/parameters.PremiumHolidayCycle)-modelPoint.PremiumHolidayUsed), 0), float64(parameters.MaximumPremiumHolidays))*projection.PremiumIncomeAdjusted, defaultPrecision)
	} else {
		projection.PremiumNotReceived = 0
		projection.PremiumNotReceivedAdjusted = 0
	}
}

// Commission Computes commission paid at each projection period based on commission type specified in the model point file
// Commission Type 1 computes initial commission
// Commission Type 2 computes renewal commission
func Commission(projection *models.Projection, p *models.Projection, modelPoint models.ProductModelPoint, features models.ProductFeatures, parameters models.ProductParameters) {
	var annualPremium float64

	if features.NonLife {
		if modelPoint.Term > 0 {
			annualPremium = modelPoint.AnnualPremium * 12.0 / float64(modelPoint.Term)
		}
	} else {
		annualPremium = modelPoint.AnnualPremium
	}

	if features.CreditLife && features.CreditLifeDecreasingPremium {
		annualPremium = modelPoint.OriginalLoan * modelPoint.PremiumRate * (12.0 / 1000.0)
	}

	if projection.ProjectionMonth == 0 || projection.ValuationTimeMonth > parameters.CalculatedTerm || modelPoint.MemberType != "MM" {
		projection.Commission = 0
		projection.CommissionAdjusted = 0

	} else {
		var commission = GetCommissionStructure(modelPoint)
		var initialtemp float64
		var initialtempadjusted float64
		var initial12temp float64
		var initial12tempadjusted float64

		if projection.ValuationTimeMonth == commission.InitialYear1CommissionPaymentMonth && projection.ProjectionMonth != 0 {
			initialtemp = utils.FloatPrecision(commission.InitialCommissionPercentage1*annualPremium+commission.InitialCommissionRand, defaultPrecision)
			initialtempadjusted = utils.FloatPrecision(commission.InitialCommissionPercentage1*annualPremium+commission.InitialCommissionRand, defaultPrecision)
		} else if projection.ValuationTimeMonth == commission.InitialYear2CommissionPaymentMonth && projection.ProjectionMonth != 0 {
			initial12temp = utils.FloatPrecision(commission.InitialCommissionPercentage2*annualPremium*p.InitialPolicy, defaultPrecision)
			initial12tempadjusted = utils.FloatPrecision(commission.InitialCommissionPercentage2*annualPremium*p.InitialPolicyAdjusted, defaultPrecision)
		}

		if projection.ValuationTimeMonth >= commission.HybridRenewalCommStartM && projection.ValuationTimeMonth <= commission.HybridRenewalCommEndM {
			projection.RenewalCommission = utils.FloatPrecision(commission.RenewalCommissionPercentage*projection.PremiumIncome+commission.RenewalCommissionRand*projection.InitialPolicy, defaultPrecision)
			projection.RenewalCommissionAdjusted = utils.FloatPrecision(commission.RenewalCommissionPercentage*projection.PremiumIncome+commission.RenewalCommissionRand*projection.InitialPolicyAdjusted, defaultPrecision)
		}
		projection.Commission = initialtemp + initial12temp
		projection.CommissionAdjusted = initialtempadjusted + initial12tempadjusted

		//if modelPoint.CommissionType == Initial {
		//	if projection.ValuationTimeMonth == 1 {
		//		projection.Commission = utils.FloatPrecision(parameters.InitialCommissionPercentage1*annualPremium+parameters.InitialCommissionRand, defaultPrecision)
		//		projection.CommissionAdjusted = utils.FloatPrecision(parameters.InitialCommissionPercentage1*annualPremium+parameters.InitialCommissionRand, defaultPrecision)
		//
		//	} else {
		//		if projection.ValuationTimeMonth == 13 {
		//			projection.Commission = utils.FloatPrecision(parameters.InitialCommissionPercentage2*annualPremium*p.InitialPolicy, defaultPrecision)
		//			projection.CommissionAdjusted = utils.FloatPrecision(parameters.InitialCommissionPercentage2*annualPremium*p.InitialPolicyAdjusted, defaultPrecision)
		//		} else {
		//			projection.Commission = 0
		//			projection.CommissionAdjusted = 0
		//		}
		//	}
		//} else if modelPoint.CommissionType == Renewal {
		//	projection.RenewalCommission = utils.FloatPrecision(parameters.RenewalCommissionPercentage*projection.PremiumIncome+parameters.RenewalCommissionRand*projection.InitialPolicy, defaultPrecision)
		//	projection.RenewalCommissionAdjusted = utils.FloatPrecision(parameters.RenewalCommissionPercentage*projection.PremiumIncome+parameters.RenewalCommissionRand*projection.InitialPolicyAdjusted, defaultPrecision)
		//} else if modelPoint.CommissionType == Hybrid {
		//	var initialtemp float64
		//	var initialtempadjusted float64
		//	var initial12temp float64
		//	var initial12tempadjusted float64
		//	if projection.ValuationTimeMonth == 1 && projection.ProjectionMonth != 0 {
		//		initialtemp = utils.FloatPrecision(parameters.InitialCommissionPercentage1*annualPremium+parameters.InitialCommissionRand, defaultPrecision)
		//		initialtempadjusted = utils.FloatPrecision(parameters.InitialCommissionPercentage1*annualPremium+parameters.InitialCommissionRand, defaultPrecision)
		//	} else if projection.ValuationTimeMonth == 13 && projection.ProjectionMonth != 0 {
		//		initial12temp = utils.FloatPrecision(parameters.InitialCommissionPercentage2*annualPremium*p.InitialPolicy, defaultPrecision)
		//		initial12tempadjusted = utils.FloatPrecision(parameters.InitialCommissionPercentage2*annualPremium*p.InitialPolicyAdjusted, defaultPrecision)
		//	}
		//
		//	if projection.ValuationTimeMonth >= parameters.HybridRenewalCommStartM && projection.ValuationTimeMonth <= parameters.HybridRenewalCommEndM {
		//		projection.RenewalCommission = utils.FloatPrecision(parameters.RenewalCommissionPercentage*projection.PremiumIncome+parameters.RenewalCommissionRand*projection.InitialPolicy, defaultPrecision)
		//		projection.RenewalCommissionAdjusted = utils.FloatPrecision(parameters.RenewalCommissionPercentage*projection.PremiumIncome+parameters.RenewalCommissionRand*projection.InitialPolicyAdjusted, defaultPrecision)
		//	}
		//	projection.Commission = initialtemp + initial12temp
		//	projection.CommissionAdjusted = initialtempadjusted + initial12tempadjusted
		//
		//} else {
		//	projection.Commission = 0
		//	projection.CommissionAdjusted = 0
		//	projection.RenewalCommission = 0
		//	projection.RenewalCommissionAdjusted = 0
		//}
	}
}

// ClawBack computes clawback for commission type 1; initial commission
func ClawBack(mp models.ProductModelPoint, projection *models.Projection, features models.ProductFeatures, parameters models.ProductParameters) {
	var annualPremium float64
	if features.NonLife {
		if mp.Term > 0 {
			annualPremium = mp.AnnualPremium * 12.0 / float64(mp.Term)
		}
	} else {
		annualPremium = mp.AnnualPremium
	}

	if features.CreditLife && features.CreditLifeDecreasingPremium {
		annualPremium = mp.OriginalLoan * mp.PremiumRate * (12.0 / 1000.0)
	}

	var commission = GetCommissionStructure(mp)

	if projection.ProjectionMonth > 0 && projection.ValuationTimeMonth <= commission.ClawbackPeriod && mp.MemberType == "MM" {
		var cb = GetClawback(projection.ValuationTimeMonth, mp)
		projection.ClawBack = utils.FloatPrecision(cb.Year1InitialCommission*projection.TotalIncrementalLapses*commission.InitialCommissionPercentage1*annualPremium+cb.Year2InitialCommission*projection.TotalIncrementalLapses*commission.InitialCommissionPercentage2*annualPremium, defaultPrecision)
		projection.ClawBackAdjusted = utils.FloatPrecision(cb.Year1InitialCommission*projection.TotalIncrementalLapsesAdjusted*commission.InitialCommissionPercentage1*annualPremium+cb.Year2InitialCommission*projection.TotalIncrementalLapsesAdjusted*commission.InitialCommissionPercentage2*annualPremium, defaultPrecision)
	} else {
		projection.ClawBack = 0
		projection.ClawBackAdjusted = 0
	}

	//if projection.ProjectionMonth > 0 && projection.ValuationTimeMonth <= parameters.ClawbackPeriod && mp.MemberType == "MM" && mp.CommissionType == Initial {
	//	var cb = GetClawback(projection.ValuationTimeMonth, mp)
	//	projection.ClawBack = utils.FloatPrecision(cb.Year1InitialCommission*projection.TotalIncrementalLapses*parameters.InitialCommissionPercentage1*annualPremium+cb.Year2InitialCommission*projection.TotalIncrementalLapses*parameters.InitialCommissionPercentage2*annualPremium, defaultPrecision)
	//	projection.ClawBackAdjusted = utils.FloatPrecision(cb.Year1InitialCommission*projection.TotalIncrementalLapsesAdjusted*parameters.InitialCommissionPercentage1*annualPremium+cb.Year2InitialCommission*projection.TotalIncrementalLapsesAdjusted*parameters.InitialCommissionPercentage2*annualPremium, defaultPrecision)
	//} else if projection.ProjectionMonth > 0 && projection.ValuationTimeMonth <= parameters.ClawbackPeriod && mp.MemberType == "MM" && mp.CommissionType == Hybrid {
	//	var cb = GetClawback(projection.ValuationTimeMonth, mp)
	//	projection.ClawBack = utils.FloatPrecision(cb.Year1InitialCommission*projection.TotalIncrementalLapses*parameters.InitialCommissionPercentage1*annualPremium+cb.Year2InitialCommission*projection.TotalIncrementalLapses*parameters.InitialCommissionPercentage2*annualPremium, defaultPrecision)
	//	projection.ClawBackAdjusted = utils.FloatPrecision(cb.Year1InitialCommission*projection.TotalIncrementalLapsesAdjusted*parameters.InitialCommissionPercentage1*annualPremium+cb.Year2InitialCommission*projection.TotalIncrementalLapsesAdjusted*parameters.InitialCommissionPercentage2*annualPremium, defaultPrecision)
	//} else {
	//	projection.ClawBack = 0
	//	projection.ClawBackAdjusted = 0
	//}
}

// DeathOutgo computes natural death outgo at each projection period
func DeathOutgo(projection *models.Projection, p *models.Projection, features models.ProductFeatures, params models.ProductParameters, modelPoint models.ProductModelPoint) {
	if projection.ValuationTimeMonth <= modelPoint.WaitingPeriod || !features.DeathBenefit || projection.ValuationTimeMonth > params.CalculatedTerm {
		projection.DeathOutgo = 0
		projection.DeathOutgoAdjusted = 0
	} else {
		if features.UnitFund || features.WithProfit {
			projection.DeathOutgo = projection.SumAtRisk * projection.TotalIncrementalNaturalDeaths
			projection.DeathOutgoAdjusted = projection.SumAtRisk * projection.TotalIncrementalNaturalDeathsAdjusted
		} else {
			projection.DeathOutgo = utils.FloatPrecision(projection.TotalIncrementalNaturalDeaths*(projection.SumAssured+params.StandardAdditionalSumAssured), defaultPrecision)
			projection.DeathOutgoAdjusted = utils.FloatPrecision(projection.TotalIncrementalNaturalDeathsAdjusted*(projection.SumAssured+params.StandardAdditionalSumAssured), defaultPrecision)
		}
	}
}

// AccidentalDeathOutgo computes accidental death outgo at each projection period
// The accidental benefit multiplier must greater than 1 if the Double Accident benefit feature is selected...
func AccidentalDeathOutgo(projection *models.Projection, features models.ProductFeatures, modelPoint models.ProductModelPoint, multipliers models.ProductAccidentalBenefitMultiplier, parameters models.ProductParameters) {
	if projection.ProjectionMonth == 0 || projection.ValuationTimeMonth > parameters.CalculatedTerm || !features.AccidentalDeathBenefit {
		projection.AccidentalDeathOutgo = 0
		projection.AccidentalDeathOutgoAdjusted = 0
	} else {
		if features.UnitFund || features.WithProfit {
			projection.AccidentalDeathOutgo = projection.SumAtRisk * projection.TotalIncrementalAccidentalDeaths
			projection.AccidentalDeathOutgoAdjusted = projection.SumAtRisk * projection.TotalIncrementalAccidentalDeathsAdjusted
		} else {
			projection.AccidentalDeathOutgo = utils.FloatPrecision(projection.TotalIncrementalAccidentalDeaths*(parameters.AccidentalBenefitMultiplier*projection.SumAssured+parameters.StandardAdditionalSumAssured), defaultPrecision)
			projection.AccidentalDeathOutgoAdjusted = utils.FloatPrecision(projection.TotalIncrementalAccidentalDeathsAdjusted*(parameters.AccidentalBenefitMultiplier*projection.SumAssured+parameters.StandardAdditionalSumAssured), defaultPrecision)
		}
	}
}

// CashBackOnSurvival computes cashback benefit outgo at each projection period
// cashback indicator is read from the model point file; non zero if applicable otherwise zero
func CashBackOnSurvival(projection *models.Projection, mp models.ProductModelPoint, parameters models.ProductParameters, features models.ProductFeatures) {
	if features.CashBackSurvival && mp.CashbackIndicator && projection.ValuationTimeMonth <= parameters.CalculatedTerm && projection.ValuationTimeMonth > 1 {
		if projection.ValuationTimeMonth%int(parameters.CashbackOnSurvivalPeriod) == 0 && projection.ValuationTimeMonth <= int(parameters.CashbackOnSurvivalTerm) && projection.ValuationTimeMonth > 1 {
			projection.CashBackOnSurvival = utils.FloatPrecision(parameters.CashbackOnSurvivalRatio*projection.InitialPolicy*mp.AnnualPremium, defaultPrecision)
			projection.CashBackOnSurvivalAdjusted = utils.FloatPrecision(parameters.CashbackOnSurvivalRatio*projection.InitialPolicyAdjusted*mp.AnnualPremium, defaultPrecision)
		} else {
			projection.CashBackOnSurvival = 0
			projection.CashBackOnSurvivalAdjusted = 0
		}
	} else {
		projection.CashBackOnSurvival = 0
		projection.CashBackOnSurvivalAdjusted = 0
	}
}

func DisabilityOutgo(projection *models.Projection, p *models.Projection, mp models.ProductModelPoint, parameters models.ProductParameters, features models.ProductFeatures) {
	if projection.ValuationTimeMonth <= mp.WaitingPeriod || !features.PermanentDisabilityBenefit || projection.ValuationTimeMonth > parameters.CalculatedTerm || projection.ValuationTimeMonth <= mp.DeferredPeriod {
		projection.DisabilityOutgo = 0
		projection.DisabilityOutgoAdjusted = 0
	} else {
		projection.DisabilityOutgo = utils.FloatPrecision(projection.TotalIncrementalDisabilities*projection.SumAssured, defaultPrecision)
		projection.DisabilityOutgoAdjusted = utils.FloatPrecision(projection.TotalIncrementalDisabilitiesAdjusted*projection.SumAssured, defaultPrecision)
	}
}

func AnnuityOutgo(projection *models.Projection, p *models.Projection, mp models.ProductModelPoint, parameters models.ProductParameters) {
	if projection.ValuationTimeMonth > parameters.CalculatedTerm || projection.ValuationTimeMonth <= mp.DeferredPeriod {
		projection.AnnuityOutgo = 0
		projection.AnnuityOutgoAdjusted = 0
	} else {
		projection.AnnuityOutgo = utils.FloatPrecision(projection.InitialPolicy*projection.AnnuityIncome*parameters.BenefitMultiplier, defaultPrecision)
		projection.AnnuityOutgoAdjusted = utils.FloatPrecision(projection.InitialPolicyAdjusted*projection.AnnuityIncome*parameters.BenefitMultiplier, defaultPrecision)
	}
}

// calculates instalment if interest rate is provided
func CalculatedInstalment(projection *models.Projection, p *models.Projection, mp models.ProductModelPoint, parameters models.ProductParameters) {
	if projection.ValuationTimeMonth <= parameters.CalculatedTerm {
		if projection.ProjectionMonth == 0 {
			if mp.Interest <= 0 {
				projection.CalculatedInstalment = mp.Instalment
			} else {
				var interestmonthly = math.Pow(1.0+mp.Interest, 1.0/12.0) - 1
				var numerator = 1 - math.Pow(1+interestmonthly, float64(-1.0*mp.OriginalTerm)) //float64(-1.0*mp.OutstandingTermMonths))
				if numerator > 0 {
					projection.CalculatedInstalment = mp.OriginalLoan * interestmonthly / numerator //mp.OutstandingLoan * interestmonthly / numerator
				}
			}
		} else {
			projection.CalculatedInstalment = p.CalculatedInstalment
		}
	} else {
		projection.CalculatedInstalment = 0
	}
}
func RetrenchmentOutgo(projection *models.Projection, mp models.ProductModelPoint, parameters models.ProductParameters, features models.ProductFeatures) {
	if projection.ValuationTimeMonth <= mp.WaitingPeriod || !features.RetrenchmentBenefit || projection.ValuationTimeMonth > parameters.CalculatedTerm {
		projection.RetrenchmentOutgo = 0
		projection.RetrenchmentOutgoAdjusted = 0
	} else {
		projection.RetrenchmentOutgo = utils.FloatPrecision(projection.TotalIncrementalRetrenchments*parameters.RetrenchmentQualificationFactor*parameters.NumberRetrenchmentInstalments*projection.CalculatedInstalment, defaultPrecision)
		projection.RetrenchmentOutgoAdjusted = utils.FloatPrecision(projection.TotalIncrementalRetrenchmentsAdjusted*parameters.RetrenchmentQualificationFactor*parameters.NumberRetrenchmentInstalments*projection.CalculatedInstalment, defaultPrecision)
	}
}

// CashBackOnDeath computes cashback benefit outgo at each projection period
// cashback indicator is read from the model point file; non zero if applicable otherwise zero
func CashBackOnDeath(projection *models.Projection, mp models.ProductModelPoint, parameters models.ProductParameters, features models.ProductFeatures) {
	if features.CashBackOnDeath {
		if mp.CashbackIndicator && projection.ValuationTimeMonth <= int(parameters.CashbackOnDeathTerm) {
			projection.CashBackOnDeath = utils.FloatPrecision(parameters.CashbackOnDeathRatio*(projection.TotalIncrementalNaturalDeaths+projection.NumberOfAccidentDeaths)*(float64(projection.ValuationTimeMonth)/12)*mp.AnnualPremium, defaultPrecision)
			projection.CashBackOnDeathAdjusted = utils.FloatPrecision(parameters.CashbackOnDeathRatio*(projection.TotalIncrementalNaturalDeathsAdjusted+projection.NumberOfAccidentDeathsAdjusted)*(float64(projection.ValuationTimeMonth)/12)*mp.AnnualPremium, defaultPrecision)
		} else {
			projection.CashBackOnDeath = 0
			projection.CashBackOnDeathAdjusted = 0
		}
	} else {
		projection.CashBackOnDeath = 0
		projection.CashBackOnDeathAdjusted = 0
	}
}

// Rider computes rider benefit cash outflow over the projection period
// The calculation happens if the rider is chosen as part of the product features
// It reads the lump sum rider benefit from the rider table
func Rider(projection *models.Projection, features models.ProductFeatures, modelPoint models.ProductModelPoint, parameters models.ProductParameters) {
	if features.RiderBenefit && projection.ValuationTimeMonth <= parameters.CalculatedTerm {
		if projection.ProjectionMonth == 0 || projection.ValuationTimeMonth <= modelPoint.WaitingPeriod {
			projection.Rider = 0
			projection.RiderAdjusted = 0

		} else {
			projection.Rider = utils.FloatPrecision((projection.TotalIncrementalNaturalDeaths+projection.TotalIncrementalAccidentalDeaths)*projection.RiderSumAssured, defaultPrecision)
			projection.RiderAdjusted = utils.FloatPrecision((projection.TotalIncrementalNaturalDeathsAdjusted+projection.TotalIncrementalAccidentalDeathsAdjusted)*projection.RiderSumAssured, defaultPrecision)
		}
	} else {
		projection.Rider = 0
		projection.RiderAdjusted = 0

	}
}

// Expenses computes expense cash out flow at each projection period.
func Expenses(projection *models.Projection, p *models.Projection, mp models.ProductModelPoint, parameters models.ProductParameters, features models.ProductFeatures, run models.RunParameters, margins models.ProductMargins, shock models.ProductShock) {
	var exprand_shocked float64
	var exprand_shockedAdjusted float64
	var polcount float64
	var polcountAdjusted float64
	var exprand float64 = 0
	var exprandAdjusted float64 = 0
	var exprand_claims float64
	var exprand_claimsadjusted float64
	var exprand_shockedclaims float64
	var exprand_shockedclaimsinflationAdjusted float64
	var exprand_shockedClaimsAdjusted float64
	var polcountclaims float64
	var polcountclaimsAdjusted float64
	var initial_attr_exp_prop float64
	var renewal_attr_exp_prop float64
	var claims_attr_exp_prop float64

	if run.IFRS17Indicator {
		initial_attr_exp_prop = parameters.InitialAttributableExpenseProportion
		renewal_attr_exp_prop = parameters.RenewalAttributableExpenseProportion
		claims_attr_exp_prop = parameters.ClaimsAttributableExpenseProportion
	} else {
		initial_attr_exp_prop = 1
		renewal_attr_exp_prop = 1
		claims_attr_exp_prop = 1
	}

	if projection.ValuationTimeMonth == 0 {
		exprand = (parameters.InitialExpensePercentage*projection.Premium + parameters.InitialExpenseRand) * initial_attr_exp_prop
		if run.ShockSettings.Expense {
			exprand_shocked = exprand*(1+shock.MultiplicativeExpense) + shock.AdditiveExpense/12.0
		} else {
			exprand_shocked = exprand
		}
		exprand_shockedAdjusted = exprand_shocked * (1 + margins.ExpenseMargin)
	} else {
		if projection.ProjectionMonth == 0 {
			exprand_shocked = 0
		} else {
			exprand = (parameters.RenewalExpensePercentage*projection.Premium + parameters.RenewalExpenseRand*projection.InflationFactor/12.0) * renewal_attr_exp_prop
			exprand_claims = (parameters.ClaimsExpensePercentage*projection.SumAssured + parameters.ClaimsExpenseRand*projection.InflationFactor) * claims_attr_exp_prop
			exprandAdjusted = (parameters.RenewalExpensePercentage*projection.Premium + parameters.RenewalExpenseRand*projection.InflationFactorAdjusted/12.0) * renewal_attr_exp_prop
			exprand_claimsadjusted = (parameters.ClaimsExpensePercentage*projection.SumAssured + parameters.ClaimsExpenseRand*projection.InflationFactorAdjusted) * claims_attr_exp_prop
			if run.ShockSettings.Expense {
				exprand_shocked = exprand*(1+shock.MultiplicativeExpense) + shock.AdditiveExpense*projection.InflationFactor/12.0
				exprand_shockedclaims = exprand_claims * (1 + shock.MultiplicativeExpense) //+ shock.AdditiveExpense * projection.InflationFactor/12.0  *claims expense can only be shocked proportionally
				exprand_shockedAdjusted = exprandAdjusted*(1+shock.MultiplicativeExpense) + shock.AdditiveExpense*projection.InflationFactorAdjusted/12.0
				exprand_shockedclaimsinflationAdjusted = exprand_claimsadjusted * (1 + shock.MultiplicativeExpense) //+ shock.AdditiveExpense * projection.InflationFactorAdjusted/12.0 * claims expense can only be shocked proportionally
			} else {
				exprand_shocked = exprand
				exprand_shockedclaims = exprand_claims
				exprand_shockedAdjusted = exprandAdjusted
				exprand_shockedclaimsinflationAdjusted = exprand_claimsadjusted
			}
			exprand_shockedAdjusted = exprand_shockedAdjusted * (1 + margins.ExpenseMargin)
			exprand_shockedClaimsAdjusted = exprand_shockedclaimsinflationAdjusted * (1 + margins.ExpenseMargin)
		}
	}
	if mp.MemberType != "MM" && features.FuneralCover {
		polcount = 0
	} else {

		polcount = p.InitialPolicy + p.InitialPremiumWaivers + p.InitialTemporaryPremiumWaivers + p.InitialPaidUp
		polcountAdjusted = p.InitialPolicyAdjusted + p.InitialPremiumWaiversAdjusted + p.InitialTemporaryPremiumWaiversAdjusted + p.InitialPaidUpAdjusted
		polcountclaims = projection.TotalIncrementalNaturalDeaths + projection.TotalIncrementalAccidentalDeaths + projection.TotalIncrementalDisabilities
		polcountclaimsAdjusted = projection.TotalIncrementalNaturalDeathsAdjusted + projection.TotalIncrementalAccidentalDeathsAdjusted + projection.TotalIncrementalDisabilitiesAdjusted
	}
	if features.AnnuityIncome {
		if projection.ValuationTimeMonth <= parameters.CalculatedTerm {
			projection.RenewalExpenses = utils.FloatPrecision((exprand_shocked*p.InitialPolicy+exprand_shockedclaims*p.InitialPolicy)*parameters.BenefitMultiplier, defaultPrecision)
			projection.RenewalExpensesAdjusted = utils.FloatPrecision((exprand_shockedAdjusted*p.InitialPolicyAdjusted+exprand_shockedClaimsAdjusted*p.InitialPolicyAdjusted)*parameters.BenefitMultiplier, defaultPrecision)
		}
	} else {
		if projection.ValuationTimeMonth <= parameters.CalculatedTerm {
			projection.RenewalExpenses = utils.FloatPrecision(exprand_shocked*polcount+exprand_shockedclaims*polcountclaims, defaultPrecision)
			projection.RenewalExpensesAdjusted = utils.FloatPrecision(exprand_shockedAdjusted*polcountAdjusted+exprand_shockedClaimsAdjusted*polcountclaimsAdjusted, defaultPrecision)
		}
	}
}

// NetCashFlow computes net cash flow at each projection period, referencing all the calculated cash flows.
func NetCashFlow(projection *models.Projection, p *models.Projection, mp models.ProductModelPoint, parameters models.ProductParameters, features models.ProductFeatures) {
	if projection.ProjectionMonth == 0 || projection.ValuationTimeMonth > parameters.CalculatedTerm+1 {
		projection.NetCashFlow = 0
		projection.NetCashFlowAdjusted = 0
	} else {
		if features.UnitFund || features.WithProfit {
			projection.NetCashFlow = projection.Commission + projection.RenewalCommission + projection.DeathOutgo + projection.AccidentalDeathOutgo + projection.CashBackOnSurvival + projection.CashBackOnDeath + projection.DisabilityOutgo + projection.RetrenchmentOutgo + projection.AnnuityOutgo + projection.Rider + projection.RenewalExpenses + projection.PremiumNotReceived + projection.NonLifeClaimsOutgo + projection.EPremiumAdvisoryCost + projection.EFundAdvisoryCost + projection.MaturityOutgo + projection.EGuaranteeCost - (projection.PremiumIncome - projection.EAllocatedPremiumIncome) - projection.ClawBack - projection.EPolicyFee - projection.EFundAssetManagementCharge - projection.EFundRiskCharge - projection.ESurrenderPenaltyCharge - projection.EPremiumAdvisoryFee - projection.EFundAdvisoryFee - projection.EBsaShareholderCharge
			projection.NetCashFlowAdjusted = projection.CommissionAdjusted + projection.RenewalCommissionAdjusted + projection.DeathOutgoAdjusted + projection.AccidentalDeathOutgoAdjusted + projection.CashBackOnSurvivalAdjusted + projection.CashBackOnDeathAdjusted + projection.DisabilityOutgoAdjusted + projection.RetrenchmentOutgoAdjusted + projection.AnnuityOutgoAdjusted + projection.RiderAdjusted + projection.RenewalExpensesAdjusted + projection.PremiumNotReceivedAdjusted + projection.NonLifeClaimsOutgoAdjusted + projection.EPremiumAdvisoryCostAdjusted + projection.EFundAdvisoryCostAdjusted + projection.MaturityOutgo + projection.EGuaranteeCost - (projection.PremiumIncomeAdjusted - projection.EAllocatedPremiumIncomeAdjusted) - projection.ClawBackAdjusted - projection.EPolicyFeeAdjusted - projection.EFundAssetManagementChargeAdjusted - projection.EFundRiskChargeAdjusted - projection.ESurrenderPenaltyChargeAdjusted - projection.EPremiumAdvisoryFeeAdjusted - projection.EFundAdvisoryFeeAdjusted - projection.EBsaShareholderChargeAdjusted
		} else {
			projection.NetCashFlow = projection.Commission + projection.RenewalCommission + projection.DeathOutgo + projection.AccidentalDeathOutgo + projection.CashBackOnSurvival + projection.CashBackOnDeath + projection.DisabilityOutgo + projection.RetrenchmentOutgo + projection.AnnuityOutgo + projection.Rider + projection.RenewalExpenses + projection.PremiumNotReceived + projection.NonLifeClaimsOutgo - projection.PremiumIncome - projection.ClawBack
			projection.NetCashFlowAdjusted = projection.CommissionAdjusted + projection.RenewalCommissionAdjusted + projection.DeathOutgoAdjusted + projection.AccidentalDeathOutgoAdjusted + projection.CashBackOnSurvivalAdjusted + projection.CashBackOnDeathAdjusted + projection.DisabilityOutgoAdjusted + projection.RetrenchmentOutgoAdjusted + projection.AnnuityOutgoAdjusted + projection.RiderAdjusted + projection.RenewalExpensesAdjusted + projection.PremiumNotReceivedAdjusted + projection.NonLifeClaimsOutgoAdjusted - projection.PremiumIncomeAdjusted - projection.ClawBackAdjusted

		}

	}
}

// CoverageUnits computes coverage units as quantity of benefits multiplied by in-force policies at each projection period
func CoverageUnits(projection *models.Projection, mp models.ProductModelPoint, features models.ProductFeatures, parameters models.ProductParameters) {
	if projection.ValuationTimeMonth > parameters.CalculatedTerm {
		projection.CoverageUnits = 0
	} else {
		projection.CoverageUnits = (projection.SumAssured + (GetRiders(mp, features))) * projection.InitialPolicy
	}
}

// Reinsurance
func Reinsurance(projection *models.Projection, p *models.Projection, features models.ProductFeatures, modelPoint models.ProductModelPoint, parameters models.ProductParameters) {
	if features.ProportionalReinsurance && projection.ValuationTimeMonth <= parameters.CalculatedTerm {
		if projection.ProjectionMonth == 0 || projection.ValuationTimeMonth <= modelPoint.WaitingPeriod {
			projection.ReinsurancePremium = 0
			projection.ReinsurancePremiumAdjusted = 0
			projection.ReinsuranceCedingCommission = 0
			projection.ReinsuranceCedingCommissionAdjusted = 0
			projection.ReinsuranceClaims = 0
			projection.ReinsuranceClaimsAdjusted = 0

		} else {
			var level1value float64
			var level2value float64
			var level3value float64

			if features.NonLife {
				level1value = math.Min(modelPoint.AnnualPremium*parameters.LossRatio, parameters.Level1Upperbound-parameters.Level1Lowerbound)
				level2value = math.Min(modelPoint.AnnualPremium*parameters.LossRatio-level1value, parameters.Level2Upperbound-parameters.Level2Lowerbound)
				level3value = math.Min(modelPoint.AnnualPremium*parameters.LossRatio-level1value-level2value, parameters.Level3Upperbound-parameters.Level3Lowerbound)
				projection.CededSumAssured = level1value*parameters.Level1CededProp + level2value*parameters.Level2CededProp + level3value*parameters.Level3CededProp
			} else if features.AnnuityIncome {
				level1value = math.Min(projection.AnnuityIncome, parameters.Level1Upperbound-parameters.Level1Lowerbound)
				level2value = math.Min(projection.AnnuityIncome-level1value, parameters.Level2Upperbound-parameters.Level2Lowerbound)
				level3value = math.Min(projection.AnnuityIncome-level1value-level2value, parameters.Level3Upperbound-parameters.Level3Lowerbound)
				projection.CededSumAssured = level1value*parameters.Level1CededProp + level2value*parameters.Level2CededProp + level3value*parameters.Level3CededProp
			} else {
				level1value = math.Min(projection.SumAssured, parameters.Level1Upperbound-parameters.Level1Lowerbound)
				level2value = math.Min(projection.SumAssured-level1value, parameters.Level2Upperbound-parameters.Level2Lowerbound)
				level3value = math.Min(projection.SumAssured-level1value-level2value, parameters.Level3Upperbound-parameters.Level3Lowerbound)
				projection.CededSumAssured = level1value*parameters.Level1CededProp + level2value*parameters.Level2CededProp + level3value*parameters.Level3CededProp
			}

			if modelPoint.SumAssured == 0 {
				projection.ReinsurancePremium = 0
				projection.ReinsurancePremiumAdjusted = 0
			} else {
				projection.ReinsurancePremium = projection.CededSumAssured * parameters.FlatAnnualReinsPremRate * (p.InitialPolicy - projection.NumberOfMaturities) / 12.0
				projection.ReinsurancePremiumAdjusted = projection.CededSumAssured * parameters.FlatAnnualReinsPremRate * (p.InitialPolicyAdjusted - projection.NumberOfMaturitiesAdjusted) / 12.0
			}
			projection.ReinsuranceCedingCommission = projection.ReinsurancePremium * parameters.CedingCommission
			projection.ReinsuranceCedingCommissionAdjusted = projection.ReinsurancePremiumAdjusted * parameters.CedingCommission
			if features.NonLife {
				projection.ReinsuranceClaims = projection.CededSumAssured * projection.NonLifeMonthlyRiskRate
				projection.ReinsuranceClaimsAdjusted = projection.CededSumAssured * projection.NonLifeMonthlyRiskRate

			} else if features.AnnuityIncome {
				projection.ReinsuranceClaims = projection.CededSumAssured * p.InitialPolicy * parameters.BenefitMultiplier
				projection.ReinsuranceClaimsAdjusted = projection.CededSumAssured * p.InitialPolicyAdjusted * parameters.BenefitMultiplier

			} else {
				projection.ReinsuranceClaims = projection.CededSumAssured * (projection.TotalIncrementalNaturalDeaths + projection.TotalIncrementalAccidentalDeaths + projection.TotalIncrementalDisabilities)
				projection.ReinsuranceClaimsAdjusted = projection.CededSumAssured * (projection.TotalIncrementalNaturalDeathsAdjusted + projection.TotalIncrementalAccidentalDeathsAdjusted + projection.TotalIncrementalDisabilitiesAdjusted)
			}

		}
		projection.NetReinsuranceCashflow = projection.ReinsuranceClaims + projection.ReinsuranceCedingCommission - projection.ReinsurancePremium
		projection.NetReinsuranceCashflowAdjusted = projection.ReinsuranceClaimsAdjusted + projection.ReinsuranceCedingCommissionAdjusted - projection.ReinsurancePremiumAdjusted

	} else {
		projection.ReinsurancePremium = 0
		projection.ReinsurancePremiumAdjusted = 0
		projection.ReinsuranceCedingCommission = 0
		projection.ReinsuranceCedingCommissionAdjusted = 0
		projection.ReinsuranceClaims = 0
		projection.ReinsuranceClaimsAdjusted = 0
		projection.NetReinsuranceCashflow = 0
		projection.NetReinsuranceCashflowAdjusted = 0

	}
}

func NonLifeRiskOutgo(projection *models.Projection, p *models.Projection, features models.ProductFeatures, params models.ProductParameters, modelPoint models.ProductModelPoint) {
	if projection.ValuationTimeMonth <= modelPoint.WaitingPeriod || !features.NonLife || projection.ValuationTimeMonth > params.CalculatedTerm || projection.ProjectionMonth == 0 {
		projection.NonLifeClaimsOutgo = 0
		projection.NonLifeClaimsOutgoAdjusted = 0
	} else {
		projection.NonLifeClaimsOutgo = utils.FloatPrecision(modelPoint.AnnualPremium*params.LossRatio*projection.NonLifeMonthlyRiskRate, defaultPrecision)
		projection.NonLifeClaimsOutgoAdjusted = projection.NonLifeClaimsOutgo
	}
}

func UnitFund(projection *models.Projection, p *models.Projection, mp models.ProductModelPoint, parameters models.ProductParameters, features models.ProductFeatures, valYear, valMonth int, run models.RunParameters) {
	var unitFundCharges models.ProductUnitFundCharge
	var fundReturns models.ProductInvestmentReturn
	var bsaReturns models.ProductInvestmentReturn
	var assetDistribution models.ProductFundAssetDistribution
	var surrenderValueCoeff models.ProductSurrenderValueCoefficient
	var maturityBenefitPattern models.ProductMaturityPattern
	var remainingTermYear int
	var remainingTermMonth int

	if features.MaturityBenefit {
		if projection.ProjectionMonth == 0 {
			projection.MaturityValue = 0
		} else {
			if projection.ValuationTimeMonth <= parameters.CalculatedTerm+1 {
				if features.UnitFund {
					if features.MaturityBenefitPattern {
						if (parameters.CalculatedTerm + 1 - projection.ValuationTimeMonth) <= 216 { //assumes maximum number of benefits can be 18(216= 18*12)
							remainingTermMonth = int(math.Max(float64(parameters.CalculatedTerm+1-projection.ValuationTimeMonth), 0.0))
							remainingTermYear = int(utils.RoundUp(math.Max(float64(parameters.CalculatedTerm+1-projection.ValuationTimeMonth), 0.0) / 12.0))
							if math.Mod(float64(remainingTermMonth), 12.0) == 0 {
								maturityBenefitPattern = GetMaturityBenefitPattern(remainingTermYear, mp, run)
								projection.MaturityValue = math.Max(mp.GuaranteedMaturityBenefit-p.UnfundedUnitFundEom, 0) * maturityBenefitPattern.PartMaturityRate // maturity value shortfall, revisit.. add guaranteed unit fund projection
							}
						} else {
							projection.MaturityValue = 0
						}
					} else {
						if projection.ValuationTimeMonth == parameters.CalculatedTerm+1 {
							projection.MaturityValue = math.Max(mp.GuaranteedMaturityBenefit-p.UnfundedUnitFundEom, 0)
						}
					}
				} else {
					if features.MaturityBenefitPattern {
						if (parameters.CalculatedTerm + 1 - projection.ValuationTimeMonth) <= 216 { //assumes maximum number of benefits can be 18(216= 18*12)
							remainingTermMonth = int(math.Max(float64(parameters.CalculatedTerm+1-projection.ValuationTimeMonth), 0.0))
							remainingTermYear = int(utils.RoundUp(math.Max(float64(parameters.CalculatedTerm+1-projection.ValuationTimeMonth), 0.0) / 12.0))
							if math.Mod(float64(remainingTermMonth), 12.0) == 0 {
								maturityBenefitPattern = GetMaturityBenefitPattern(remainingTermYear, mp, run)
								if remainingTermYear == 0 {
									projection.MaturityValue = math.Max(p.SumAssured*maturityBenefitPattern.PartMaturityRate+p.ReversionaryBonus+p.TerminalBonus, mp.GuaranteedMaturityBenefit*maturityBenefitPattern.PartMaturityRate)
								} else {
									projection.MaturityValue = math.Max(p.SumAssured, mp.GuaranteedMaturityBenefit) * maturityBenefitPattern.PartMaturityRate
								}
							}
						} else {
							projection.MaturityValue = 0
						}
					} else {
						if projection.ValuationTimeMonth == parameters.CalculatedTerm+1 {
							projection.MaturityValue = math.Max(p.SumAssured+p.ReversionaryBonus+p.TerminalBonus, mp.GuaranteedMaturityBenefit)
						}
					}
				}
			}
		}
	}

	if projection.ProjectionMonth == 0 {
		if features.UnitFund {
			projection.UnfundedUnitFundSom = mp.UnitFund
			projection.UnfundedUnitFundEom = mp.UnitFund
		}
		projection.BonusStabilisationAccount = mp.BonusStabilisationAccount
		assetDistribution = GetAssetDistribution(valYear, mp, run)
	} else if projection.ValuationTimeMonth <= parameters.CalculatedTerm {
		unitFundCharges = GetUnitFundCharge(projection.ValuationTimeMonth, mp, run)
		runYear := valYear + int((valMonth+projection.ProjectionMonth)/12)
		fundReturns = GetFundInvestmentReturns(runYear, mp, run)
		bsaReturns = GetBsaFundInvestmentReturns(runYear, mp, run)
		assetDistribution = GetAssetDistribution(valYear, mp, run)
		unitFundGrowthMargin := assetDistribution.Cash*fundReturns.CashYieldGap + assetDistribution.CorporateBond*fundReturns.CorporateBondRiskPremium + assetDistribution.Equity*fundReturns.EquityRiskPremium + assetDistribution.Property*fundReturns.PropertyRiskPremium
		riskFreeRate, _ := GetForwardRateWithError(projection.ProjectionMonth, valYear, run.YieldcurveMonth, parameters.YieldCurveCode)
		if features.UnitFund {
			projection.UnitGrowthRiskMargin = unitFundGrowthMargin

		}
		monthlyInvestmentGrowth := math.Pow(1+riskFreeRate+unitFundGrowthMargin+fundReturns.FundReversionaryBonus, 1.0/12.0) - 1

		if features.WithProfit {
			projection.AllocatedPremium = 0
		} else {
			projection.AllocatedPremium = projection.Premium * unitFundCharges.PremiumAllocationRate * (1 - unitFundCharges.BidOfferSpread)
		}

		projection.PolicyFee = projection.Premium*unitFundCharges.PremiumPolicyFeeRate + unitFundCharges.AnnualPolicyFeeAmount/12.0

		if features.AdvisoryFeeCharge {
			projection.PremiumAdvisoryFee = projection.Premium*unitFundCharges.PremiumAdvisoryFeeRate + unitFundCharges.AnnualAdvisoryFeeAmount/12.0
			projection.FundAdvisoryFee = math.Max(p.UnfundedUnitFundEom+projection.AllocatedPremium-projection.PolicyFee, 0)*unitFundCharges.PremiumAdvisoryFeeRate + unitFundCharges.AnnualAdvisoryFeeAmount/12.0
		}

		if features.UnitFund {
			projection.UnfundedUnitFundSom = p.UnfundedUnitFundEom + projection.AllocatedPremium - projection.PolicyFee - projection.PremiumAdvisoryFee - projection.FundAdvisoryFee - projection.MaturityValue
			projection.FundInvestmentIncome = projection.UnfundedUnitFundSom * monthlyInvestmentGrowth
			projection.FundAssetManagementCharge = (projection.UnfundedUnitFundSom + projection.FundInvestmentIncome) * unitFundCharges.FundManagementChargeRate
		}

		if features.WithProfit {
			projection.ReversionaryBonus = (p.SumAssured + p.ReversionaryBonus) * (math.Pow(1+fundReturns.FundReversionaryBonus, 1.0/12.0) - 1)
			projection.TerminalBonus = projection.ReversionaryBonus * fundReturns.FundTerminalBonus
			projection.SumAtRisk = mp.SumAssured + projection.ReversionaryBonus + projection.TerminalBonus
		} else {
			projection.SumAtRisk = math.Max(mp.SumAssured-(projection.UnfundedUnitFundSom+projection.FundInvestmentIncome-projection.FundAssetManagementCharge), 0)
		}

		if features.FundRiskCharge {
			projection.FundRiskCharge = projection.SumAtRisk * (projection.MonthlyDependentMortality)
		}

		if features.UnitFund {
			projection.UnfundedUnitFundEom = projection.UnfundedUnitFundSom + projection.FundInvestmentIncome - projection.FundAssetManagementCharge - projection.FundRiskCharge
			projection.EAllocatedPremiumIncome = projection.AllocatedPremium * p.InitialPolicy
			projection.EAllocatedPremiumIncomeAdjusted = projection.AllocatedPremium * p.InitialPolicyAdjusted
		}
		projection.EPolicyFee = projection.PolicyFee * p.InitialPolicy
		projection.EPolicyFeeAdjusted = projection.PolicyFee * p.InitialPolicyAdjusted

		projection.EPremiumAdvisoryFee = projection.PremiumAdvisoryFee * p.InitialPolicy
		projection.EPremiumAdvisoryFeeAdjusted = projection.PremiumAdvisoryFee * p.InitialPolicyAdjusted
		projection.EFundAdvisoryFee = projection.FundAdvisoryFee * p.InitialPolicy
		projection.EFundAdvisoryFeeAdjusted = projection.FundAdvisoryFee * p.InitialPolicyAdjusted

		projection.BonusStabilisationAccount = p.BonusStabilisationAccount * math.Pow(1+riskFreeRate+bsaReturns.FundGrowthMargin, 1.0/12.0) * (1 - unitFundCharges.BsaShareholderFeeRate)
		projection.EBsaShareholderCharge = math.Pow(p.BonusStabilisationAccount*(1+riskFreeRate+bsaReturns.FundGrowthMargin), 1.0/12.0) * unitFundCharges.BsaShareholderFeeRate * projection.InitialPolicy
		projection.EBsaShareholderChargeAdjusted = math.Pow(p.BonusStabilisationAccount*(1+riskFreeRate+bsaReturns.FundGrowthMargin), 1.0/12.0) * unitFundCharges.BsaShareholderFeeRate * projection.InitialPolicyAdjusted
		projection.EFundAssetManagementCharge = projection.FundAssetManagementCharge * projection.InitialPolicy
		projection.EFundAssetManagementChargeAdjusted = projection.FundAssetManagementCharge * projection.InitialPolicyAdjusted
		if features.FundRiskCharge {
			projection.EFundRiskCharge = projection.SumAtRisk * (projection.TotalIncrementalNaturalDeaths + projection.TotalIncrementalAccidentalDeaths)
			projection.EFundRiskChargeAdjusted = projection.SumAtRisk * (projection.TotalIncrementalNaturalDeathsAdjusted + projection.TotalIncrementalAccidentalDeathsAdjusted)
		}

		if features.SurrenderBenefit {
			if features.UnitFund {
				projection.SurrenderPenalty = math.Min(projection.UnfundedUnitFundEom*unitFundCharges.MarketValueAdjustment*unitFundCharges.FundSurrenderPenaltyRate, projection.UnfundedUnitFundEom*unitFundCharges.MarketValueAdjustment*(1-unitFundCharges.FundMinimumSurrenderValueRate))
				projection.SurrenderValue = math.Max(projection.UnfundedUnitFundEom*unitFundCharges.MarketValueAdjustment-projection.SurrenderPenalty, 0)
			} else {
				if features.SurrenderValueQuadraticFormula {
					surrenderValueCoeff = GetSurrenderQuadraticCoefficients(mp, run)
					projection.SurrenderValue = (projection.SumAssured + projection.ReversionaryBonus) * math.Min(1, surrenderValueCoeff.A*math.Pow(float64(projection.ValuationTimeMonth), 2.0)+surrenderValueCoeff.B*float64(projection.ValuationTimeMonth)+surrenderValueCoeff.C)
				} else {
					projection.SurrenderValue = (projection.SumAssured + projection.ReversionaryBonus) * (1 - unitFundCharges.FundSurrenderPenaltyRate)
				}
			}
		}

		if features.SurrenderValueQuadraticFormula {
			if projection.ValuationTimeMonth <= surrenderValueCoeff.SurrenderPayoutWaitingPeriodMonths {
				projection.SurrenderOutgo = 0
				projection.SurrenderOutgoAdjusted = 0
			} else {
				projection.SurrenderOutgo = projection.SurrenderValue * projection.TotalIncrementalLapses
				projection.SurrenderOutgoAdjusted = projection.SurrenderValue * projection.TotalIncrementalLapsesAdjusted
			}
		} else {
			projection.SurrenderOutgo = projection.SurrenderValue * projection.TotalIncrementalLapses
			projection.SurrenderOutgoAdjusted = projection.SurrenderValue * projection.TotalIncrementalLapsesAdjusted
		}
		projection.ESurrenderPenaltyCharge = projection.SurrenderPenalty * projection.TotalIncrementalLapses
		projection.ESurrenderPenaltyChargeAdjusted = projection.SurrenderPenalty * projection.TotalIncrementalLapsesAdjusted

		// same as premium and fund advisory fee
		projection.EPremiumAdvisoryCost = (projection.Premium*unitFundCharges.PremiumAdvisoryFeeRate + unitFundCharges.AnnualAdvisoryFeeAmount/12.0) * p.InitialPolicy
		projection.EPremiumAdvisoryCostAdjusted = (projection.Premium*unitFundCharges.PremiumAdvisoryFeeRate + unitFundCharges.AnnualAdvisoryFeeAmount/12.0) * p.InitialPolicyAdjusted
		projection.EFundAdvisoryCost = (math.Max(p.UnfundedUnitFundEom+projection.AllocatedPremium-projection.PolicyFee, 0)*unitFundCharges.PremiumAdvisoryFeeRate + unitFundCharges.AnnualAdvisoryFeeAmount/12.0) * p.InitialPolicy
		projection.EFundAdvisoryCostAdjusted = (math.Max(p.UnfundedUnitFundEom+projection.AllocatedPremium-projection.PolicyFee, 0)*unitFundCharges.PremiumAdvisoryFeeRate + unitFundCharges.AnnualAdvisoryFeeAmount/12.0) * p.InitialPolicyAdjusted
		//if features.UnitFund && features.DeathBenefit {
		//	if projection.ValuationTimeMonth <= mp.WaitingPeriod {
		//		projection.ERiskCost = projection.SumAtRisk * projection.TotalIncrementalAccidentalDeaths
		//		projection.ERiskCostAdjusted = projection.SumAtRisk * projection.TotalIncrementalAccidentalDeathsAdjusted
		//	} else {
		//		projection.ERiskCost = projection.SumAtRisk * (projection.TotalIncrementalNaturalDeaths + projection.TotalIncrementalAccidentalDeaths)
		//		projection.ERiskCostAdjusted = projection.SumAtRisk * (projection.TotalIncrementalNaturalDeathsAdjusted + projection.TotalIncrementalAccidentalDeathsAdjusted)
		//	}
		//}
	}
}

func MaturityOutgo(projection *models.Projection, p *models.Projection, mp models.ProductModelPoint, parameters models.ProductParameters, features models.ProductFeatures, valYear, valMonth int, run models.RunParameters) {
	if projection.ValuationTimeMonth <= parameters.CalculatedTerm+1 {
		if projection.ProjectionMonth > 0 {
			if features.UnitFund {
				projection.EGuaranteeCost = projection.MaturityValue * projection.NumberOfMaturities
				projection.EGuaranteeCostAdjusted = projection.MaturityValue * projection.NumberOfMaturitiesAdjusted
			} else {
				projection.MaturityOutgo = projection.MaturityValue * projection.NumberOfMaturities
				projection.MaturityOutgoAdjusted = projection.MaturityValue * projection.NumberOfMaturitiesAdjusted
			}
		}
	}

}
