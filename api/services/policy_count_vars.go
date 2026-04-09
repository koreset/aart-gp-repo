package services

import (
	"api/models"
	"api/utils"
	"math"
)

var startingInitialPolicy float64 = 0
var startingInitialPolicyAdjusted float64 = 0

// InitialPaidUp projects number of paid up policies over the projection period
// dependent on other decrements or transitions applicable on a product
// Paidup no death of the main member
func InitialPaidUp(projection *models.Projection, p *models.Projection, modelPoint models.ProductModelPoint, features models.ProductFeatures, index int, parameters models.ProductParameters) {
	if features.FuneralCover {
		if index == 0 || projection.ValuationTimeMonth > parameters.CalculatedTerm {
			if modelPoint.PaidupIndicator && modelPoint.PaidupOption && projection.ValuationTimeMonth <= parameters.CalculatedTerm {
				projection.InitialPaidUp = 1
				projection.InitialPaidUpAdjusted = 1
			} else {
				projection.InitialPaidUp = 0
				projection.InitialPaidUpAdjusted = 0
			}
		}
		if modelPoint.PaidupOption {
			projection.InitialPaidUp = utils.FloatPrecision(math.Max(p.InitialPolicy-(projection.NaturalDeathsInForce-p.NaturalDeathsInForce)-
				(projection.NumberOfAccidentDeaths-p.NumberOfAccidentDeaths)-(projection.NumberOfLapses-p.NumberOfLapses), 0)*projection.PaidUpOnFactor+
				math.Max(p.InitialPaidUp-(projection.NaturalDeathsPaidUp-p.NaturalDeathsPaidUp)-(projection.AccidentDeathsPaidUp-p.AccidentDeathsPaidUp), 0), defaultdecrementPrecision)
			projection.InitialPaidUpAdjusted = utils.FloatPrecision(math.Max(p.InitialPolicyAdjusted-(projection.NaturalDeathsInForceAdjusted-p.NaturalDeathsInForceAdjusted)-
				(projection.NumberOfAccidentDeathsAdjusted-p.NumberOfAccidentDeathsAdjusted)-(projection.NumberOfLapsesAdjusted-p.NumberOfLapsesAdjusted), 0)*projection.PaidUpOnFactor+
				math.Max(p.InitialPaidUpAdjusted-(projection.NaturalDeathsPaidUpAdjusted-p.NaturalDeathsPaidUpAdjusted)-(projection.AccidentDeathsPaidUpAdjusted-p.AccidentDeathsPaidUpAdjusted), 0), defaultdecrementPrecision)

		} else {
			projection.InitialPaidUp = 0
			projection.InitialPaidUpAdjusted = 0
		}

	}
}

// InitialPolicy projects number of policies in healthy state over the projection period
// dependent on decrements or transitions from the healthy state to other applicable states
// uses dependent transition rates to project number of policies
func InitialPolicy(index int, projection *models.Projection, modelPoint models.ProductModelPoint, parameters models.ProductParameters) {
	if index == 0 {
		if !modelPoint.PremiumWaiverIndicator && !modelPoint.PaidupIndicator {
			projection.InitialPolicy = 1
			projection.InitialPolicyAdjusted = 1
			if modelPoint.MemberType == "MM" {
				projection.PolicyCount = 1
			}
		} else {
			projection.InitialPolicy = 0
			projection.InitialPolicyAdjusted = 0
			projection.InitialPaidUp = 1
			projection.InitialPaidUpAdjusted = 1
			if modelPoint.MemberType == "MM" {
				projection.PolicyCount = 1
			}
		}
		startingInitialPolicy = 1         // projection.InitialPolicy
		startingInitialPolicyAdjusted = 1 // projection.InitialPolicyAdjusted
		if parameters.CalculatedTerm == 0 || modelPoint.DurationInForceMonths > parameters.CalculatedTerm {
			projection.NumberOfMaturities = 1
			projection.NumberOfMaturitiesAdjusted = 1
			projection.InitialPolicy = 0
			projection.InitialPolicyAdjusted = 0
		}
	} else {
		if projection.ValuationTimeMonth <= parameters.CalculatedTerm {
			projection.InitialPolicy = utils.FloatPrecision(math.Max(
				startingInitialPolicy- //projection.InitialPremiumWaivers-projection.InitialPaidUp-
					projection.NumberOfMaturities-
					projection.NaturalDeathsInForce-
					projection.NumberOfAccidentDeaths-
					projection.NumberOfLapses-
					projection.NumberOfDisabilities-
					projection.NumberOfRetrenchments,
				0), defaultdecrementPrecision)
			projection.InitialPolicyAdjusted = utils.FloatPrecision(math.Max(
				startingInitialPolicyAdjusted- //projection.InitialPremiumWaiversAdjusted-projection.InitialPaidUpAdjusted-
					projection.NumberOfMaturitiesAdjusted-
					projection.NaturalDeathsInForceAdjusted-
					projection.NumberOfAccidentDeathsAdjusted-
					projection.NumberOfLapsesAdjusted-
					projection.NumberOfDisabilitiesAdjusted-
					projection.NumberOfRetrenchmentsAdjusted,
				0), defaultdecrementPrecision)
			if modelPoint.MemberType == "MM" {
				projection.PolicyCount = projection.InitialPolicy
			}
		} else {
			projection.InitialPolicy = 0
			projection.InitialPolicyAdjusted = 0
		}
	}
}

// NaturalDeathsInForce computes cumulative number of natural deaths over the projection period
func NaturalDeathsInForce(index int, projection *models.Projection, p *models.Projection, parameters models.ProductParameters) {
	if index == 0 {
		projection.NaturalDeathsInForce = 0
		projection.NaturalDeathsInForceAdjusted = 0
	} else {
		if projection.ValuationTimeMonth <= parameters.CalculatedTerm {
			//projection.NaturalDeathsInForce
			projection.NaturalDeathsInForce = utils.FloatPrecision(math.Max(p.InitialPolicy-projection.NumberOfMaturities, 0)*projection.MonthlyDependentMortality*(1-(projection.AccidentProportion))+p.NaturalDeathsInForce, defaultdecrementPrecision)
			projection.NaturalDeathsInForceAdjusted = utils.FloatPrecision(math.Max(p.InitialPolicyAdjusted-projection.NumberOfMaturitiesAdjusted, 0)*projection.MonthlyDependentMortalityAdjusted*(1-(projection.AccidentProportion))+p.NaturalDeathsInForceAdjusted, defaultdecrementPrecision)
		} else {
			projection.NaturalDeathsInForce = p.NaturalDeathsInForce
			projection.NaturalDeathsInForceAdjusted = p.NaturalDeathsInForce
		}
	}
}

// NumberOfDeathsAccident computes cumulative number of accidental deaths over the projection period
func NumberOfDeathsAccident(index int, projection *models.Projection, p *models.Projection, parameters models.ProductParameters) {
	if index == 0 {
		projection.NumberOfAccidentDeaths = 0
		projection.NumberOfAccidentDeathsAdjusted = 0
	} else {
		if projection.ValuationTimeMonth <= parameters.CalculatedTerm {
			projection.NumberOfAccidentDeaths = utils.FloatPrecision(math.Max(p.InitialPolicy-projection.NumberOfMaturities, 0)*projection.MonthlyDependentMortality*projection.AccidentProportion+p.NumberOfAccidentDeaths, defaultdecrementPrecision)
			projection.NumberOfAccidentDeathsAdjusted = utils.FloatPrecision(math.Max(p.InitialPolicyAdjusted-projection.NumberOfMaturitiesAdjusted, 0)*projection.MonthlyDependentMortalityAdjusted*projection.AccidentProportion+p.NumberOfAccidentDeathsAdjusted, defaultdecrementPrecision)
		} else {
			projection.NumberOfAccidentDeaths = p.NumberOfAccidentDeaths
			projection.NumberOfAccidentDeathsAdjusted = p.NumberOfAccidentDeathsAdjusted
		}
	}
}

// NumberOfLapses computes number of cumulative lapses over the projection period
func NumberOfLapses(index int, projection *models.Projection, p *models.Projection, modelPoint models.ProductModelPoint, parameters models.ProductParameters) {
	if index == 0 {
		projection.NumberOfLapses = 0
		projection.NumberOfLapsesAdjusted = 0
	} else {
		if (projection.MainMemberAgeNextBirthday >= parameters.PaidupEffectiveAge && modelPoint.PaidupOption) || projection.ValuationTimeMonth > parameters.CalculatedTerm {
			projection.NumberOfLapses = p.NumberOfLapses
			projection.NumberOfLapsesAdjusted = p.NumberOfLapsesAdjusted
		} else {
			projection.NumberOfLapses = utils.FloatPrecision(math.Max(p.InitialPolicy-projection.NumberOfMaturities, 0)*projection.MonthlyDependentLapse+p.NumberOfLapses, defaultdecrementPrecision)
			projection.NumberOfLapsesAdjusted = utils.FloatPrecision(math.Max(p.InitialPolicyAdjusted-projection.NumberOfMaturitiesAdjusted, 0)*projection.MonthlyDependentLapseAdjusted+p.NumberOfLapsesAdjusted, defaultdecrementPrecision)
		}
	}
}

// NumberOfDisabilities computes number of cumulative disabilities over the projection period
func NumberOfDisabilities(index int, projection *models.Projection, p *models.Projection, mp models.ProductModelPoint, parameters models.ProductParameters) {
	if index == 0 || mp.MemberType != "MM" {
		projection.NumberOfDisabilities = 0
		projection.NumberOfDisabilitiesAdjusted = 0
	} else {
		if projection.ValuationTimeMonth <= parameters.CalculatedTerm {
			projection.NumberOfDisabilities = utils.FloatPrecision(math.Max(p.InitialPolicy-projection.NumberOfMaturities, 0)*projection.MonthlyDependentDisability+p.NumberOfDisabilities, defaultdecrementPrecision)
			projection.NumberOfDisabilitiesAdjusted = utils.FloatPrecision(math.Max(p.InitialPolicyAdjusted-projection.NumberOfMaturitiesAdjusted, 0)*projection.MonthlyDependentDisabilityAdjusted+p.NumberOfDisabilitiesAdjusted, defaultdecrementPrecision)
		} else {
			projection.NumberOfDisabilities = p.NumberOfDisabilities
			projection.NumberOfDisabilitiesAdjusted = p.NumberOfDisabilitiesAdjusted
		}
	}
}

// NumberOfRetrenchments Computes number of cumulative retrenchments over the projection period
func NumberOfRetrenchments(index int, projection *models.Projection, p *models.Projection, mp models.ProductModelPoint, parameters models.ProductParameters) {
	if index == 0 || mp.MemberType != "MM" {
		projection.NumberOfRetrenchments = 0
		projection.NumberOfRetrenchmentsAdjusted = 0
	} else {
		if projection.ValuationTimeMonth <= parameters.CalculatedTerm {
			projection.NumberOfRetrenchments = utils.FloatPrecision(math.Max(p.InitialPolicy-projection.NumberOfMaturities, 0)*projection.MonthlyDependentRetrenchment*(1-p.NumberOfRetrenchments)+p.NumberOfRetrenchments, defaultdecrementPrecision)
			projection.NumberOfRetrenchmentsAdjusted = utils.FloatPrecision(math.Max(p.InitialPolicyAdjusted-projection.NumberOfMaturitiesAdjusted, 0)*projection.MonthlyDependentRetrenchmentAdjusted*(1-p.NumberOfRetrenchmentsAdjusted)+p.NumberOfRetrenchmentsAdjusted, defaultdecrementPrecision)
		} else {
			projection.NumberOfRetrenchments = p.NumberOfRetrenchments
			projection.NumberOfRetrenchmentsAdjusted = p.NumberOfRetrenchmentsAdjusted
		}

	}
}

// InitialPremiumWaivers computes number cumulative number of premium waivers over the projection period
func InitialPremiumWaivers(index int, projection *models.Projection, p *models.Projection, modelPoint models.ProductModelPoint, features models.ProductFeatures, parameters models.ProductParameters) {
	if features.FuneralCover {
		if index == 0 || projection.ValuationTimeMonth > parameters.CalculatedTerm {
			if modelPoint.PremiumWaiverIndicator && modelPoint.ContinuityOrPremiumWaiverOption && projection.ValuationTimeMonth <= parameters.CalculatedTerm {
				projection.InitialPremiumWaivers = 1
				projection.InitialPremiumWaiversAdjusted = 1
			} else {
				projection.InitialPremiumWaivers = 0
				projection.InitialPremiumWaiversAdjusted = 0
			}

		} else {
			if modelPoint.PremiumWaiverIndicator && modelPoint.ContinuityOrPremiumWaiverOption {
				projection.InitialPremiumWaivers = utils.FloatPrecision(math.Max(p.InitialPremiumWaivers-(projection.NaturalDeathsPremiumWaiver-p.NaturalDeathsPremiumWaiver)-(projection.AccidentDeathsPremiumWaiver-p.AccidentDeathsPremiumWaiver), 0)+math.Max(p.InitialPolicy-(projection.NaturalDeathsInForce-p.NaturalDeathsInForce)-(projection.NumberOfAccidentDeaths-p.NumberOfAccidentDeaths)-(projection.NumberOfLapses-p.NumberOfLapses), 0)*(p.ContractingPartyAlivePortion*projection.MainMemberMortalityRateByMonth)*projection.PremiumWaiverOnFactor, defaultdecrementPrecision)
				projection.InitialPremiumWaiversAdjusted = utils.FloatPrecision(math.Max(p.InitialPremiumWaiversAdjusted-(projection.NaturalDeathsPremiumWaiverAdjusted-p.NaturalDeathsPremiumWaiverAdjusted)-(projection.AccidentDeathsPremiumWaiverAdjusted-p.AccidentDeathsPremiumWaiverAdjusted), 0)+math.Max(p.InitialPolicyAdjusted-(projection.NaturalDeathsInForceAdjusted-p.NaturalDeathsInForceAdjusted)-(projection.NumberOfAccidentDeathsAdjusted-p.NumberOfAccidentDeathsAdjusted)-(projection.NumberOfLapsesAdjusted-p.NumberOfLapsesAdjusted), 0)*(p.ContractingPartyAlivePortionAdjusted*projection.MainMemberMortalityRateAdjustedByMonth)*projection.PremiumWaiverOnFactor, defaultdecrementPrecision)
			} else {
				projection.InitialPremiumWaivers = 0
				projection.InitialPremiumWaiversAdjusted = 0
			}
		}
	} else {
		projection.InitialPremiumWaivers = 0
		projection.InitialPremiumWaiversAdjusted = 0
	}
}

// NaturalDeathsPremiumWaiver computes cumulative number of natural deaths,in the premium waiver state, over the projection period
func NaturalDeathsPremiumWaiver(index int, projection *models.Projection, p *models.Projection, parameters models.ProductParameters) {
	if index == 0 || projection.ValuationTimeMonth > parameters.CalculatedTerm {
		projection.NaturalDeathsPremiumWaiver = 0
		projection.NaturalDeathsPremiumWaiverAdjusted = 0

	} else {
		projection.NaturalDeathsPremiumWaiver = utils.FloatPrecision(p.InitialPremiumWaivers*projection.MonthlyDependentMortality*(1-projection.AccidentProportion)+p.NaturalDeathsPremiumWaiver, defaultdecrementPrecision)
		projection.NaturalDeathsPremiumWaiverAdjusted = utils.FloatPrecision(p.InitialPremiumWaiversAdjusted*projection.MonthlyDependentMortalityAdjusted*(1-projection.AccidentProportion)+p.NaturalDeathsPremiumWaiverAdjusted, defaultdecrementPrecision)
	}
}

// NaturalDeathsTemporaryWaivers computes cumulative number of natural deaths, in the temporary premium waiver, state over the projection period
func NaturalDeathsTemporaryWaivers(index int, projection *models.Projection, p *models.Projection, parameters models.ProductParameters) {
	if index == 0 || projection.ValuationTimeMonth > parameters.CalculatedTerm {
		projection.NaturalDeathsTemporaryWaivers = 0
		projection.NaturalDeathsTemporaryWaiversAdjusted = 0
	} else {
		projection.NaturalDeathsTemporaryWaivers = utils.FloatPrecision(p.InitialTemporaryPremiumWaivers*projection.MonthlyDependentMortality*(1-projection.AccidentProportion)+p.NaturalDeathsTemporaryWaivers, defaultdecrementPrecision)
		projection.NaturalDeathsTemporaryWaiversAdjusted = utils.FloatPrecision(p.InitialTemporaryPremiumWaiversAdjusted*projection.MonthlyDependentMortalityAdjusted*(1-projection.AccidentProportion)+p.NaturalDeathsTemporaryWaiversAdjusted, defaultdecrementPrecision)
	}
}

// NaturalDeathsPaidUp computes cumulative number of natural deaths, in the PaidUP, state over the projection period
func NaturalDeathsPaidUp(index int, projection *models.Projection, p *models.Projection, parameters models.ProductParameters) {
	if index == 0 || projection.ValuationTimeMonth > parameters.CalculatedTerm {
		projection.NaturalDeathsPaidUp = 0
		projection.NaturalDeathsPaidUpAdjusted = 0
	} else {
		projection.NaturalDeathsPaidUp = utils.FloatPrecision(p.InitialPaidUp*projection.MonthlyDependentMortality*(1-projection.AccidentProportion)+p.NaturalDeathsPaidUp, defaultdecrementPrecision)
		projection.NaturalDeathsPaidUpAdjusted = utils.FloatPrecision(p.InitialPaidUpAdjusted*projection.MonthlyDependentMortalityAdjusted*(1-projection.AccidentProportion)+p.NaturalDeathsPaidUpAdjusted, defaultdecrementPrecision)
	}
}

// AccidentDeathsPremiumWaiver computes cumulative number of accidental deaths,in the premium waiver state, over the projection period
func AccidentDeathsPremiumWaiver(index int, projection *models.Projection, p *models.Projection, parameters models.ProductParameters) {
	if index == 0 || projection.ValuationTimeMonth > parameters.CalculatedTerm {
		projection.AccidentDeathsPremiumWaiver = 0
		projection.AccidentDeathsPremiumWaiverAdjusted = 0
	} else {
		projection.AccidentDeathsPremiumWaiver = utils.FloatPrecision(p.InitialPremiumWaivers*projection.MonthlyDependentMortality*projection.AccidentProportion+p.AccidentDeathsPremiumWaiver, defaultdecrementPrecision)
		projection.AccidentDeathsPremiumWaiverAdjusted = utils.FloatPrecision(p.InitialPremiumWaiversAdjusted*projection.MonthlyDependentMortalityAdjusted*projection.AccidentProportion+p.AccidentDeathsPremiumWaiverAdjusted, defaultdecrementPrecision)
	}
}

// AccidentDeathsTemporaryPremiumWaiver computes cumulative number of accidental deaths, in temporary premium waiver state, over the projection period
func AccidentDeathsTemporaryPremiumWaiver(index int, projection *models.Projection, p *models.Projection, parameters models.ProductParameters) {
	if index == 0 || projection.ValuationTimeMonth > parameters.CalculatedTerm {
		projection.AccidentDeathsTemporaryPremiumWaiver = 0
		projection.AccidentDeathsTemporaryPremiumWaiverAdjusted = 0

	} else {
		projection.AccidentDeathsTemporaryPremiumWaiver = utils.FloatPrecision(p.InitialTemporaryPremiumWaivers*projection.MonthlyDependentMortality*projection.AccidentProportion+p.AccidentDeathsTemporaryPremiumWaiver, defaultdecrementPrecision)
		projection.AccidentDeathsTemporaryPremiumWaiverAdjusted = utils.FloatPrecision(p.InitialTemporaryPremiumWaiversAdjusted*projection.MonthlyDependentMortalityAdjusted*projection.AccidentProportion+p.AccidentDeathsTemporaryPremiumWaiverAdjusted, defaultdecrementPrecision)

	}
}

// TotalIncrementalLapses computes number of new lapses over the projection period
func TotalIncrementalLapses(index int, projection *models.Projection, p *models.Projection, parameters models.ProductParameters) {
	if index == 0 || projection.ValuationTimeMonth > parameters.CalculatedTerm {
		projection.TotalIncrementalLapses = 0
		projection.TotalIncrementalLapsesAdjusted = 0
	} else {
		projection.TotalIncrementalLapses = utils.FloatPrecision(projection.NumberOfLapses-p.NumberOfLapses, defaultdecrementPrecision)
		projection.TotalIncrementalLapsesAdjusted = utils.FloatPrecision(projection.NumberOfLapsesAdjusted-p.NumberOfLapsesAdjusted, defaultdecrementPrecision)
	}
}

// NumberOfMaturities computes number of cumulative maturities over the projection period
func NumberOfMaturities(index int, projection *models.Projection, p *models.Projection, mp models.ProductModelPoint, features models.ProductFeatures, states []models.ProductTransitionState, parameters models.ProductParameters, run models.RunParameters) {
	//if utils.StatesContains(&states, Maturity) {
	if index == 0 {
		projection.NumberOfMaturities = 0
		projection.NumberOfMaturitiesAdjusted = 0
		if parameters.CalculatedTerm == 0 {
			projection.NumberOfMaturities = 1
			projection.NumberOfMaturitiesAdjusted = 1
		}
	} else {
		if projection.ValuationTimeMonth <= parameters.CalculatedTerm+1 {
			if features.MaturityBenefitPattern {
				if (parameters.CalculatedTerm + 1 - projection.ValuationTimeMonth) > 216 { //assumes maximum number of benefits can be 18(216= 18*12)
					projection.NumberOfMaturities = p.NumberOfMaturities
					projection.NumberOfMaturitiesAdjusted = p.NumberOfMaturitiesAdjusted
				} else {
					if p.NumberOfMaturities > 0 { //assumes maturity benefits are guaranteed from the initial maturity benefit.so other decrements apply
						projection.NumberOfMaturities = p.NumberOfMaturities
						projection.NumberOfMaturitiesAdjusted = p.NumberOfMaturitiesAdjusted
					} else {
						remainingTermYear := int(utils.RoundUp(math.Max(float64(parameters.CalculatedTerm+1-projection.ValuationTimeMonth), 0.0) / 12.0))
						matRate := GetMaturityBenefitPattern(remainingTermYear, mp, run)
						if matRate.PartMaturityRate > 0 {
							projection.NumberOfMaturities = utils.FloatPrecision(math.Max(startingInitialPolicy-p.NumberOfMaturities-p.NaturalDeathsInForce-p.NumberOfAccidentDeaths-p.NumberOfLapses-p.NumberOfDisabilities-p.NumberOfRetrenchments, 0), defaultdecrementPrecision)
							projection.NumberOfMaturitiesAdjusted = utils.FloatPrecision(math.Max(startingInitialPolicy-p.NumberOfMaturitiesAdjusted-p.NaturalDeathsInForceAdjusted-p.NumberOfAccidentDeathsAdjusted-p.NumberOfLapsesAdjusted-p.NumberOfDisabilitiesAdjusted-p.NumberOfRetrenchments, 0), defaultdecrementPrecision)
						} else {
							projection.NumberOfMaturities = 0
							projection.NumberOfMaturitiesAdjusted = 0
						}
					}
				}
			} else {
				if projection.ValuationTimeMonth == parameters.CalculatedTerm+1 && parameters.CalculatedTerm != 0 {
					projection.NumberOfMaturities = utils.FloatPrecision(math.Max(startingInitialPolicy-p.NumberOfMaturities-p.NaturalDeathsInForce-p.NumberOfAccidentDeaths-p.NumberOfLapses-p.NumberOfDisabilities-p.NumberOfRetrenchments, 0)+p.NumberOfMaturities, defaultdecrementPrecision)
					projection.NumberOfMaturitiesAdjusted = utils.FloatPrecision(math.Max(startingInitialPolicy-p.NumberOfMaturitiesAdjusted-p.NaturalDeathsInForceAdjusted-p.NumberOfAccidentDeathsAdjusted-p.NumberOfLapsesAdjusted-p.NumberOfDisabilitiesAdjusted-p.NumberOfRetrenchments, 0)+p.NumberOfMaturitiesAdjusted, defaultdecrementPrecision)
				} else {
					projection.NumberOfMaturities = p.NumberOfMaturities
					projection.NumberOfMaturitiesAdjusted = p.NumberOfMaturitiesAdjusted
				}
			}
		} else {
			projection.NumberOfMaturities = p.NumberOfMaturities
			projection.NumberOfMaturitiesAdjusted = p.NumberOfMaturitiesAdjusted
		}
	}
}

// TotalIncrementalNaturalDeaths computes number of natural deaths at each projection period
func TotalIncrementalNaturalDeaths(projection *models.Projection, p *models.Projection, index int, parameters models.ProductParameters) {
	if index == 0 || projection.ValuationTimeMonth > parameters.CalculatedTerm {
		projection.TotalIncrementalNaturalDeaths = 0
		projection.TotalIncrementalNaturalDeathsAdjusted = 0
		return
	}
	projection.TotalIncrementalNaturalDeaths = utils.FloatPrecision((projection.NaturalDeathsInForce+projection.NaturalDeathsPaidUp+projection.NaturalDeathsPremiumWaiver+projection.NaturalDeathsTemporaryWaivers)-
		(p.NaturalDeathsInForce+p.NaturalDeathsPaidUp+p.NaturalDeathsPremiumWaiver+p.NaturalDeathsTemporaryWaivers), defaultdecrementPrecision)
	projection.TotalIncrementalNaturalDeathsAdjusted = utils.FloatPrecision((projection.NaturalDeathsInForceAdjusted+projection.NaturalDeathsPaidUpAdjusted+projection.NaturalDeathsPremiumWaiverAdjusted+projection.NaturalDeathsTemporaryWaiversAdjusted)-
		(p.NaturalDeathsInForceAdjusted+p.NaturalDeathsPaidUpAdjusted+p.NaturalDeathsPremiumWaiverAdjusted+p.NaturalDeathsTemporaryWaiversAdjusted), defaultdecrementPrecision)
}

// TotalIncrementalDisabilities computes number of disability at each projection period
func TotalIncrementalDisabilities(projection *models.Projection, p *models.Projection, index int, parameters models.ProductParameters) {
	if index == 0 || projection.ValuationTimeMonth > parameters.CalculatedTerm {
		projection.TotalIncrementalDisabilities = 0
		projection.TotalIncrementalDisabilitiesAdjusted = 0
		return
	}
	projection.TotalIncrementalDisabilities = utils.FloatPrecision((projection.NumberOfDisabilities)-
		(p.NumberOfDisabilities), defaultdecrementPrecision)
	projection.TotalIncrementalDisabilitiesAdjusted = utils.FloatPrecision((projection.NumberOfDisabilitiesAdjusted)-
		(p.NumberOfDisabilitiesAdjusted), defaultdecrementPrecision)
}

// TotalIncrementalRetrenchments computes number of retrenchment at each projection period
func TotalIncrementalRetrenchments(projection *models.Projection, p *models.Projection, index int, parameters models.ProductParameters) {
	if index == 0 || projection.ValuationTimeMonth > parameters.CalculatedTerm {
		projection.TotalIncrementalRetrenchments = 0
		projection.TotalIncrementalRetrenchmentsAdjusted = 0
		return
	}
	projection.TotalIncrementalRetrenchments = utils.FloatPrecision((projection.NumberOfRetrenchments)-
		(p.NumberOfRetrenchments), defaultdecrementPrecision)
	projection.TotalIncrementalRetrenchmentsAdjusted = utils.FloatPrecision((projection.NumberOfRetrenchmentsAdjusted)-
		(p.NumberOfRetrenchmentsAdjusted), defaultdecrementPrecision)
}

// AccidentDeathsPaidUp computes cumulative number of accidental deaths, in the paid up state, over the projection period
func AccidentDeathsPaidUp(projection *models.Projection, p *models.Projection, index int) {
	if index == 0 {
		projection.AccidentDeathsPaidUp = 0
		projection.AccidentDeathsPaidUpAdjusted = 0
		return
	}
	projection.AccidentDeathsPaidUp = utils.FloatPrecision(p.InitialPaidUp*projection.MonthlyDependentMortality*projection.AccidentProportion+p.AccidentDeathsPaidUp, defaultdecrementPrecision)
	projection.AccidentDeathsPaidUpAdjusted = utils.FloatPrecision(p.InitialPaidUpAdjusted*projection.MonthlyDependentMortalityAdjusted*projection.AccidentProportion+p.AccidentDeathsPaidUpAdjusted, defaultdecrementPrecision)
}

// TotalIncrementalAccidentalDeaths computes number of accidental deaths at each projection period
func TotalIncrementalAccidentalDeaths(projection *models.Projection, p *models.Projection, index int, parameters models.ProductParameters) {
	if index == 0 || projection.ValuationTimeMonth > parameters.CalculatedTerm {
		projection.TotalIncrementalAccidentalDeaths = 0
		projection.TotalIncrementalAccidentalDeathsAdjusted = 0
		return
	}
	projection.TotalIncrementalAccidentalDeaths = utils.FloatPrecision((projection.NumberOfAccidentDeaths+projection.AccidentDeathsPaidUp+projection.AccidentDeathsPremiumWaiver+projection.AccidentDeathsTemporaryPremiumWaiver)-
		(p.NumberOfAccidentDeaths+p.AccidentDeathsPaidUp+p.AccidentDeathsPremiumWaiver+p.AccidentDeathsTemporaryPremiumWaiver), defaultdecrementPrecision)
	projection.TotalIncrementalAccidentalDeathsAdjusted = utils.FloatPrecision((projection.NumberOfAccidentDeathsAdjusted+projection.AccidentDeathsPaidUpAdjusted+projection.AccidentDeathsPremiumWaiverAdjusted+projection.AccidentDeathsTemporaryPremiumWaiverAdjusted)-
		(p.NumberOfAccidentDeathsAdjusted+p.AccidentDeathsPaidUpAdjusted+p.AccidentDeathsPremiumWaiverAdjusted+p.AccidentDeathsTemporaryPremiumWaiverAdjusted), defaultdecrementPrecision)
}
