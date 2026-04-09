package services

import (
	"api/models"
	"api/utils"
	"math"
)

const baseIndependentPrecision = 9
const contractingPartyPrecision = 30

// ContractingPartyAlivePortion computes a policy's lapse rate, allowing for the contracting party's mortality
// It computes dependency of the policy's lapse rate(for non-contracting party assured lives) on the contracting party's mortality.
func ContractingPartyAlivePortion(index int, projection *models.Projection, p models.Projection, modelPoint models.ProductModelPoint, features models.ProductFeatures, parameters models.ProductParameters) {
	if index == 0 || projection.ValuationTimeMonth > parameters.CalculatedTerm { // This is the first loop of the run
		if modelPoint.PremiumWaiverIndicator || projection.ValuationTimeMonth > parameters.CalculatedTerm {
			projection.ContractingPartyAlivePortion = 0
			projection.ContractingPartyAlivePortionAdjusted = 0
		} else {
			projection.ContractingPartyAlivePortion = 1
			projection.ContractingPartyAlivePortionAdjusted = 1
		}
		projection.ContractingPartyPolicyLapse = 0
		projection.ContractingPartyPolicyLapseAdjusted = 0
	} else if features.LapseDependentOnCpDeath && projection.ValuationTimeMonth <= parameters.CalculatedTerm {
		if modelPoint.MemberType != "MM" {
			projection.ContractingPartyAlivePortion = utils.FloatPrecision(p.ContractingPartyAlivePortion*math.Pow(1-projection.BaseLapse, 1/12.0)*math.Pow(1-projection.MainMemberMortalityRate, 1/12.0), contractingPartyPrecision)
			projection.ContractingPartyAlivePortionAdjusted = utils.FloatPrecision(p.ContractingPartyAlivePortionAdjusted*math.Pow(1-projection.BaseLapseAdjusted, 1/12.0)*math.Pow(1-projection.MainMemberMortalityRateAdjusted, 1/12.0), contractingPartyPrecision)

		} else {
			projection.ContractingPartyAlivePortion = utils.FloatPrecision(p.ContractingPartyAlivePortion*math.Pow(1-projection.BaseLapse, 1/12.0), contractingPartyPrecision)
			projection.ContractingPartyAlivePortionAdjusted = utils.FloatPrecision(p.ContractingPartyAlivePortionAdjusted*math.Pow(1-projection.BaseLapseAdjusted, 1/12.0), contractingPartyPrecision)
		}
	} else if projection.ValuationTimeMonth <= parameters.CalculatedTerm {
		projection.ContractingPartyAlivePortion = utils.FloatPrecision(p.ContractingPartyAlivePortion*math.Pow(1-projection.BaseLapse, 1/12.0), contractingPartyPrecision)
		projection.ContractingPartyAlivePortionAdjusted = utils.FloatPrecision(p.ContractingPartyAlivePortionAdjusted*math.Pow(1-projection.BaseLapseAdjusted, 1/12.0), contractingPartyPrecision)
	} else {
		projection.ContractingPartyAlivePortion = 0
		projection.ContractingPartyAlivePortionAdjusted = 0
	}
	if p.ContractingPartyAlivePortion == 0 || projection.ContractingPartyAlivePortion == 0 {
		projection.ContractingPartyPolicyLapse = 0
		projection.ContractingPartyPolicyLapseAdjusted = 0
	} else {
		if modelPoint.ContinuityOrPremiumWaiverOption && modelPoint.DurationInForceMonths+projection.ProjectionMonth > int(parameters.PremiumWaiverWaitingPeriod) {
			projection.ContractingPartyPolicyLapse = utils.FloatPrecision(1-math.Pow(1-projection.BaseLapse, 1/12.0), contractingPartyPrecision)
			projection.ContractingPartyPolicyLapseAdjusted = utils.FloatPrecision(1-math.Pow(1-projection.BaseLapseAdjusted, 1/12.0), contractingPartyPrecision)

		} else {
			if p.ContractingPartyAlivePortion > 0 {
				projection.ContractingPartyPolicyLapse = utils.FloatPrecision(1-projection.ContractingPartyAlivePortion/p.ContractingPartyAlivePortion, contractingPartyPrecision)
			}
			if p.ContractingPartyAlivePortionAdjusted > 0 {
				projection.ContractingPartyPolicyLapseAdjusted = utils.FloatPrecision(1-projection.ContractingPartyAlivePortionAdjusted/p.ContractingPartyAlivePortionAdjusted, contractingPartyPrecision)
			}
		}
	}
}

// BaseIndependentLapse computes independent annual lapse rate from the ContractingPartyPolicyLapse
func BaseIndependentLapse(projection *models.Projection, parameters models.ProductParameters) {
	if projection.ValuationTimeMonth <= parameters.CalculatedTerm {
		projection.BaseIndependentLapse = utils.FloatPrecision(1-math.Pow(1-projection.ContractingPartyPolicyLapse, 12.0), defaultPrecision)
		projection.BaseIndependentLapseAdjusted = utils.FloatPrecision(1-math.Pow(1-projection.ContractingPartyPolicyLapseAdjusted, 12.0), defaultPrecision)
	} else {
		projection.BaseIndependentLapse = 0
		projection.BaseIndependentLapseAdjusted = 0
	}

}

// IndependentLapseMonthly converts BaseIndependentLapse into independent monthly lapse rate
func IndependentLapseMonthly(projection *models.Projection, parameters models.ProductParameters) {
	if projection.ValuationTimeMonth <= parameters.CalculatedTerm {
		projection.IndependentLapseMonthly = utils.FloatPrecision(1-math.Pow(1-projection.BaseIndependentLapse, 1/12.0), defaultPrecision)
		projection.IndependentLapseMonthlyAdjusted = utils.FloatPrecision(1-math.Pow(1-projection.BaseIndependentLapseAdjusted, 1/12.0), defaultPrecision)
	} else {
		projection.IndependentLapseMonthly = 0
		projection.IndependentLapseMonthlyAdjusted = 0
	}

}

// MonthlyDependentLapse converts the IndependentLapseMonthly into dependent monthly lapse rate
func MonthlyDependentLapse(projection *models.Projection, parameters models.ProductParameters) {
	if projection.ValuationTimeMonth <= parameters.CalculatedTerm {
		projection.MonthlyDependentLapse = utils.FloatPrecision(projection.IndependentLapseMonthly*(1-TimingMortalityOne*projection.IndependentMortalityRateMonthly)*(1-TimingDisabilityOne*projection.IndependentDisabilityMonthly)*(1-TimingRetrenchmentOne*projection.IndependentRetrenchmentMonthly), defaultPrecision)
		projection.MonthlyDependentLapseAdjusted = utils.FloatPrecision(projection.IndependentLapseMonthlyAdjusted*(1-TimingMortalityOne*projection.IndependentMortalityRateAdjustedByMonth)*(1-TimingDisabilityOne*projection.IndependentDisabilityMonthlyAdjusted)*(1-TimingRetrenchmentOne*projection.IndependentRetrenchmentMonthlyAdjusted), defaultPrecision)
	} else {
		projection.MonthlyDependentLapse = 0
		projection.MonthlyDependentLapseAdjusted = 0
	}

}

// MainMemberMortalityRateByMonth reads base main member's mortality rate by rating factors
func MainMemberMortalityRateByMonth(projection *models.Projection, features models.ProductFeatures, parameters models.ProductParameters) {
	if projection.ValuationTimeMonth <= parameters.CalculatedTerm {
		projection.MainMemberMortalityRateByMonth = utils.FloatPrecision(1-math.Pow(1-projection.MainMemberMortalityRate, 1/12.0), defaultPrecision)
		projection.MainMemberMortalityRateAdjustedByMonth = utils.FloatPrecision(1-math.Pow(1-projection.MainMemberMortalityRateAdjusted, 1/12.0), defaultPrecision)
	} else {
		projection.MainMemberMortalityRateByMonth = 0
		projection.MainMemberMortalityRateAdjustedByMonth = 0
	}
}

// IndependentMortalityRateMonthly transforms annual base mortality rates into monthly mortality rates
func IndependentMortalityRateMonthly(projection *models.Projection, parameters models.ProductParameters) {
	if projection.ValuationTimeMonth <= parameters.CalculatedTerm {
		projection.IndependentMortalityRateMonthly = utils.FloatPrecision(1-math.Pow(1-projection.BaseMortalityRate, 1/12.0), defaultPrecision)
		projection.IndependentMortalityRateAdjustedByMonth = utils.FloatPrecision(1-math.Pow(1-projection.BaseMortalityRateAdjusted, 1/12.0), defaultPrecision)
	} else {
		projection.IndependentMortalityRateMonthly = 0
		projection.IndependentMortalityRateAdjustedByMonth = 0
	}
}

// MonthlyDependentMortality converts monthly independent mortality rates into monthly dependent mortality rates
func MonthlyDependentMortality(projection *models.Projection, parameters models.ProductParameters) {
	if projection.ValuationTimeMonth <= parameters.CalculatedTerm {
		projection.MonthlyDependentMortality = utils.FloatPrecision(projection.IndependentMortalityRateMonthly*(1-TimingLapseZero*projection.IndependentLapseMonthly)*(1-TimingDisabilityHalf*projection.IndependentDisabilityMonthly)*(1-TimingRetrenchmentZero*projection.IndependentRetrenchmentMonthly), defaultPrecision)
		projection.MonthlyDependentMortalityAdjusted = utils.FloatPrecision(projection.IndependentMortalityRateAdjustedByMonth*(1-TimingLapseZero*projection.IndependentLapseMonthlyAdjusted)*(1-TimingDisabilityHalf*projection.IndependentDisabilityMonthlyAdjusted)*(1-TimingRetrenchmentZero*projection.IndependentRetrenchmentMonthlyAdjusted), defaultPrecision)
	} else {
		projection.MonthlyDependentMortality = 0
		projection.MonthlyDependentMortalityAdjusted = 0
	}

}

// IndependentRetrenchmentMonthly converts annual independent retrenchment rate into monthly independent retrenchment rate
func IndependentRetrenchmentMonthly(projection *models.Projection, parameters models.ProductParameters) {
	if projection.ValuationTimeMonth <= parameters.CalculatedTerm {
		projection.IndependentRetrenchmentMonthly = utils.FloatPrecision(1-math.Pow(1-projection.BaseRetrenchmentRate, 1/12.0), defaultPrecision)
		projection.IndependentRetrenchmentMonthlyAdjusted = utils.FloatPrecision(1-math.Pow(1-projection.BaseRetrenchmentRateAdjusted, 1/12.0), defaultPrecision)
	} else {
		projection.IndependentRetrenchmentMonthly = 0
		projection.IndependentRetrenchmentMonthlyAdjusted = 0
	}
}

// MonthlyDependentRetrenchment converts independent monthly retrenchment rate into dependent monthly retrenchment rate
func MonthlyDependentRetrenchment(projection *models.Projection, parameters models.ProductParameters) {
	if projection.ValuationTimeMonth <= parameters.CalculatedTerm {
		projection.MonthlyDependentRetrenchment = utils.FloatPrecision(projection.IndependentRetrenchmentMonthly*(1-TimingLapseZero*projection.IndependentLapseMonthly)*(1-TimingMortalityOne*projection.IndependentMortalityRateMonthly)*(1-TimingDisabilityOne*projection.IndependentDisabilityMonthly), defaultPrecision)
		projection.MonthlyDependentRetrenchmentAdjusted = utils.FloatPrecision(projection.IndependentRetrenchmentMonthlyAdjusted*(1-TimingLapseZero*projection.IndependentLapseMonthlyAdjusted)*(1-TimingMortalityOne*projection.IndependentMortalityRateAdjustedByMonth)*(1-TimingDisabilityOne*projection.IndependentDisabilityMonthlyAdjusted), defaultPrecision)
	} else {
		projection.MonthlyDependentRetrenchment = 0
		projection.MonthlyDependentRetrenchmentAdjusted = 0
	}
}

// IndependentDisabilityMonthly computes monthly independent incidence rates from annual independent incidence rates
// References BaseDisabilityIncrement
func IndependentDisabilityMonthly(projection *models.Projection, parameters models.ProductParameters) {
	if projection.ValuationTimeMonth <= parameters.CalculatedTerm {
		projection.IndependentDisabilityMonthly = utils.FloatPrecision(1-math.Pow(1-projection.BaseDisabilityIncidenceRate, 1/12.0), defaultPrecision)
		projection.IndependentDisabilityMonthlyAdjusted = utils.FloatPrecision(1-math.Pow(1-projection.BaseDisabilityIncidenceRateAdjusted, 1/12.0), defaultPrecision)
	} else {
		projection.IndependentDisabilityMonthly = 0
		projection.IndependentDisabilityMonthlyAdjusted = 0
	}
}

// MonthlyDependentDisability converts monthly independent disability incidence rates into dependent disability monthly incidence rates
// Dependent on the assumed timings for other transition eg. lapse timing and mortality timing
func MonthlyDependentDisability(projection *models.Projection, parameters models.ProductParameters) {
	if projection.ValuationTimeMonth <= parameters.CalculatedTerm {
		projection.MonthlyDependentDisability = utils.FloatPrecision(projection.IndependentDisabilityMonthly*(1-TimingLapseZero*projection.IndependentLapseMonthly)*(1-TimingMortalityHalf*projection.IndependentMortalityRateMonthly)*(1-TimingRetrenchmentZero*projection.IndependentRetrenchmentMonthly), defaultPrecision)
		projection.MonthlyDependentDisabilityAdjusted = utils.FloatPrecision(projection.IndependentDisabilityMonthlyAdjusted*(1-TimingLapseZero*projection.IndependentLapseMonthlyAdjusted)*(1-TimingMortalityHalf*projection.IndependentMortalityRateAdjustedByMonth)*(1-TimingRetrenchmentZero*projection.IndependentRetrenchmentMonthlyAdjusted), defaultPrecision)
	} else {
		projection.MonthlyDependentDisability = 0
		projection.MonthlyDependentDisabilityAdjusted = 0
	}

}
