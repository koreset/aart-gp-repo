package services

import (
	"api/models"
	"api/utils"
	"github.com/rs/zerolog/log"
	"math"
	"strconv"
	"strings"
)

// PricingBenefitInForce indicates whether or not a benefit is on. 1 if is on and zero otherwise
func PricingBenefitInForce(pp *models.PricingPoint, modelPoint models.ProductPricingModelPoint, features models.ProductFeatures, parameters models.ProductPricingParameters) {
	//BenefitInForce - Credit Life
	if features.CreditLife {
		if pp.ProjectionMonth <= modelPoint.OutstandingTermMonths && pp.ValuationTimeMonth > modelPoint.WaitingPeriod {
			pp.BenefitInForce = 1
		} else {
			pp.BenefitInForce = 0
		}
	} else {
		if pp.ValuationTimeMonth <= parameters.CalculatedTerm {
			if modelPoint.PremiumWaiverIndicator && pp.ValuationTimeMonth > int(parameters.PremiumWaiverWaitingPeriod) && pp.ValuationTimeMonth <= int(parameters.PaidupEffectiveDuration) {
				pp.BenefitInForce = parameters.PremiumWaiverSumAssuredFactor
			} else {
				pp.BenefitInForce = 1
			}
		} else {
			pp.BenefitInForce = 0
		}
	}
}

// PricingAccidentProportion reads the proportion of the base mortality that arises from non-natural causes
// Used by downstream variables to calculate number of accidental deaths
func PricingAccidentProportion(pp *models.PricingPoint, mp models.ProductPricingModelPoint, states []models.ProductTransitionState, parameters models.ProductPricingParameters, columnname, tableName string) {
	// Accident Proportion
	if utils.StatesContains(&states, AccidentalDeath) && pp.ValuationTimeMonth <= parameters.CalculatedTerm {
		args := models.TransitionRateArguments{
			ProductId:           pp.JobProductID,
			ProductCode:         mp.ProductCode,
			Year:                0000,
			Age:                 pp.AgeNextBirthday,
			Gender:              mp.Gender,
			SmokerStatus:        mp.SmokerStatus,
			Income:              mp.Income,
			SocioEconomicClass:  mp.SocioEconomicClass,
			OccupationalClass:   mp.OccupationalClass,
			SelectPeriod:        mp.SelectPeriod,
			EducationLevel:      mp.EducationLevel,
			DurationIfM:         mp.DurationInForceMonths + pp.ProjectionMonth,
			ProjectionMonth:     pp.ProjectionMonth,
			DistributionChannel: mp.DistributionChannel,
		}
		pp.AccidentProportion = getPricingMortalityRateAccidentProportion(args, columnname, tableName)
	} else {
		pp.AccidentProportion = 0
	}
}

// PricingInflationFactor computes the expense inflation factor at each projection period
// It reads the expense inflation rate from the inflation curve
func PricingInflationFactor(i int, pp *models.PricingPoint, p models.PricingPoint, productMargins models.ProductPricingMargins, pricingParams models.PricingParameter) {
	// Inflation Factor and Adjusted
	if i == 0 {
		pp.InflationFactor = 1
		pp.InflationFactorAdjusted = 1
	} else {
		inflationFactor := getPricingInflationFactor(pp.ProjectionMonth, pricingParams.YieldCurveCode)
		pp.InflationFactor = utils.FloatPrecision(p.InflationFactor*math.Pow(1+inflationFactor, 1/12.0), defaultPrecision)
		pp.InflationFactorAdjusted = utils.FloatPrecision(p.InflationFactorAdjusted*math.Pow(1+inflationFactor*(1+productMargins.InflationMargin), 1/12.0), defaultPrecision)
	}
}

// PricingPremiumEscalation computes premium escalation factor at each projection period
// It reads premium escalation rate from the model point file
func PricingPremiumEscalation(i int, pp *models.PricingPoint, p models.PricingPoint, modelPoint models.ProductPricingModelPoint, features models.ProductFeatures, pricingConfig models.PricingConfig) {
	// Premium Escalation
	if i == 0 {
		pp.PremiumEscalation = math.Pow(1+modelPoint.PremiumEscalation, math.Max(math.Ceil(pp.ValuationTimeYear)-1, 0))
		pp.SumAssuredEscalation = math.Pow(1+modelPoint.SumAssuredEscalation, math.Max(math.Ceil(pp.ValuationTimeYear)-1, 0))

	} else {
		if (pp.ValuationTimeMonth-1)%12 == 0 && pp.ValuationTimeMonth > 1 {
			if features.ProductLevelEscalations && features.AgeRatedEscalations { //reads from the product level escalations table
				premiumEscalationRate := getEscalationRate(pp.MainMemberAgeNextBirthday, modelPoint.ProductCode, "PremiumEscalationRate", pricingConfig.ParameterBasis)
				sumAssuredEscalationRate := getEscalationRate(pp.MainMemberAgeNextBirthday, modelPoint.ProductCode, "SumAssuredEscalationRate", pricingConfig.ParameterBasis)
				pp.PremiumEscalation = p.PremiumEscalation * (1.0 + premiumEscalationRate)
				pp.SumAssuredEscalation = p.SumAssuredEscalation * (1.0 + sumAssuredEscalationRate)
			}
			if features.ProductLevelEscalations && !features.AgeRatedEscalations { //reads from the product level escalations table
				premiumEscalationRate := getEscalationRate(modelPoint.MainMemberAgeAtEntry, modelPoint.ProductCode, "PremiumEscalationRate", pricingConfig.ParameterBasis)
				sumAssuredEscalationRate := getEscalationRate(modelPoint.MainMemberAgeAtEntry, modelPoint.ProductCode, "SumAssuredEscalationRate", pricingConfig.ParameterBasis)
				pp.PremiumEscalation = math.Pow(1.0+premiumEscalationRate, math.Max(math.Ceil(pp.ValuationTimeYear)-1.0, 0.0))
				pp.SumAssuredEscalation = math.Pow(1.0+sumAssuredEscalationRate, math.Max(math.Ceil(pp.ValuationTimeYear)-1.0, 0.0))

			}
			if !features.ProductLevelEscalations { // Model Point escalation referenced
				pp.PremiumEscalation = math.Pow(1.0+modelPoint.PremiumEscalation, math.Max(math.Ceil(pp.ValuationTimeYear)-1.0, 0.0))
				pp.SumAssuredEscalation = math.Pow(1.0+modelPoint.SumAssuredEscalation, math.Max(math.Ceil(pp.ValuationTimeYear)-1.0, 0.0))

			}
		} else {
			pp.PremiumEscalation = p.PremiumEscalation
			pp.SumAssuredEscalation = p.SumAssuredEscalation
		}
		// To delete temp for kedari
		//pp.SumAssuredEscalation = math.Pow(1.0+modelPoint.SumAssuredEscalation, math.Max(float64(pp.ValuationTimeMonth-1.0)/12.0, 0.0))

	}
}

// PricingLapseMargin reads lapse margin at each projection period
func PricingLapseMargin(pp *models.PricingPoint, pricingParams models.PricingParameter, LapseMarginMonthCount int) {
	// LapseMargin
	pp.LapseMargin = getPricingLapseMargin(pp.ValuationTimeMonth, pp.ProductCode, pricingParams.Basis, LapseMarginMonthCount)
}

// PricingPremiumWaiverOnFactor indicates whether or not the premium waiver benefit is on. 1 if on and zero otherwise
func PricingPremiumWaiverOnFactor(pp *models.PricingPoint, modelPoint models.ProductPricingModelPoint, parameters models.ProductPricingParameters) {
	// PremiumWaiverOnFactor
	if pp.ProjectionYear > parameters.CalculatedTerm || !modelPoint.ContinuityOrPremiumWaiverOption || (modelPoint.DurationInForceMonths+pp.ProjectionMonth) < int(parameters.PremiumWaiverWaitingPeriod) || modelPoint.MemberType == MainMember {
		pp.PremiumWaiverOnFactor = 0
	} else {
		pp.PremiumWaiverOnFactor = 1
	}
}

// PricingPaidUpOnFactor indicates whether or not the paid up benefit is on. 1 if on and zero otherwise
func PricingPaidUpOnFactor(pp *models.PricingPoint, modelPoint models.ProductPricingModelPoint, parameters models.ProductPricingParameters) {
	// PaidUpOnFactor
	if pp.ProjectionMonth+modelPoint.DurationInForceMonths <= parameters.CalculatedTerm && pp.MainMemberAgeNextBirthday >= parameters.PaidupEffectiveAge && modelPoint.PaidupOption && pp.ValuationTimeMonth >= int(parameters.PremiumWaiverWaitingPeriod) {
		pp.PaidUpOnFactor = 1
	} else {
		pp.PaidUpOnFactor = 0
	}
}

// PricingMainMemberMortalityRate reads annual mortality rate by rating factors
func PricingMainMemberMortalityRate(pp *models.PricingPoint, mp models.ProductPricingModelPoint, productMargins models.ProductPricingMargins, features models.ProductFeatures, states []models.ProductTransitionState, parameters models.ProductPricingParameters, columnname, tableName string, shock models.ProductPricingShock, pricingShockBasis string) {
	if utils.StatesContains(&states, Death) && pp.ValuationTimeMonth <= parameters.CalculatedTerm {
		if pp.MainMemberAgeNextBirthday >= 120 {
			pp.MainMemberMortalityRate = 1
			pp.MainMemberMortalityRateAdjusted = 1
		} else if pp.MainMemberAgeNextBirthday < 120 {
			args := models.TransitionRateArguments{
				ProductId:           pp.JobProductID,
				ProductCode:         pp.ProductCode,
				Age:                 pp.MainMemberAgeNextBirthday,
				Gender:              mp.MainMemberGender,
				SmokerStatus:        mp.SmokerStatus,
				Income:              mp.Income,
				SocioEconomicClass:  mp.SocioEconomicClass,
				OccupationalClass:   mp.OccupationalClass,
				SelectPeriod:        mp.SelectPeriod,
				EducationLevel:      mp.EducationLevel,
				DurationIfM:         mp.DurationInForceMonths + pp.ProjectionMonth,
				ProjectionMonth:     pp.ProjectionMonth,
				DistributionChannel: mp.DistributionChannel,
				Year:                0000,
			}
			tempRate := getPricingMortalityRate(args, columnname, tableName)

			if pricingShockBasis != "N/A" {
				pp.MainMemberMortalityRate = math.Max(0, math.Min(1, tempRate*(1+shock.MultiplicativeMortality)+shock.AdditiveMortality))
			}
			if pricingShockBasis == "N/A" {
				pp.MainMemberMortalityRate = tempRate
			}

			var resp float64 = 0
			if pp.MainMemberAgeNextBirthday > 120 {
				resp = 1
			} else {
				resp = pp.MainMemberMortalityRate * (1 + productMargins.MortalityMargin)
			}
			pp.MainMemberMortalityRateAdjusted = utils.FloatPrecision(math.Min(resp, 1), defaultPrecision)
		}
	} else {
		pp.MainMemberMortalityRate = 0
		pp.MainMemberMortalityRateAdjusted = 0
	}
}

// PricingBaseLapse reads annual base lapse rate by the respective rating factors
func PricingBaseLapse(pp *models.PricingPoint, mp models.ProductPricingModelPoint, features models.ProductFeatures, prodCode string, parameters models.ProductPricingParameters, pricingparams models.PricingParameter, columnName, tableName string, shock models.ProductPricingShock, pricingShockBasis string) {
	if mp.TemporaryPremiumWaiverIndicator || mp.TemporaryPremiumWaiverMonthExit > pp.ValuationTimeMonth || mp.PremiumWaiverIndicator || mp.PaidupOption || pp.ValuationTimeMonth > parameters.CalculatedTerm || pp.ProjectionMonth == 0 {
		pp.BaseLapse = 0
		pp.BaseLapseAdjusted = 0
	} else {
		args := models.TransitionRateArguments{
			ProductId:           pp.JobProductID,
			ProductCode:         pp.ProductCode,
			Age:                 pp.MainMemberAgeNextBirthday,
			Gender:              mp.MainMemberGender,
			SmokerStatus:        mp.SmokerStatus,
			Income:              mp.Income,
			SocioEconomicClass:  mp.SocioEconomicClass,
			OccupationalClass:   mp.OccupationalClass,
			SelectPeriod:        mp.SelectPeriod,
			EducationLevel:      mp.EducationLevel,
			DurationIfM:         mp.DurationInForceMonths + pp.ProjectionMonth,
			ProjectionMonth:     pp.ProjectionMonth,
			DistributionChannel: mp.DistributionChannel,
			Year:                0000,
		}
		monthInView := pp.ValuationTimeMonth
		tempRate := getPricingLapseRate(args, columnName, tableName)

		if pricingShockBasis != "N/A" {
			pp.BaseLapse = math.Max(0, math.Min(1, tempRate*(1+shock.MultiplicativeLapse)+shock.AdditiveLapse))
		}
		if pricingShockBasis == "N/A" {
			pp.BaseLapse = tempRate
		}
		pp.BaseLapseAdjusted = pp.BaseLapse * (1 + (getPricingLapseMargin(monthInView, pp.ProductCode, pricingparams.Basis, pricingTables[mp.ProductCode].LapseMarginMonthCount)))
	}
}

// PricingContractingPartyAlivePortion computes a policy's lapse rate, allowing for the contracting party's mortality
// It computes dependency of the policy's lapse rate(for non-contracting party assured lives) on the contracting party's mortality.
func PricingContractingPartyAlivePortion(pp *models.PricingPoint, p models.PricingPoint, modelPoint models.ProductPricingModelPoint,
	features models.ProductFeatures, parameters models.ProductPricingParameters) {
	if pp.ProjectionMonth == 0 || pp.ValuationTimeMonth > parameters.CalculatedTerm { // This is the first loop of the run
		if modelPoint.PremiumWaiverIndicator || pp.ValuationTimeMonth > parameters.CalculatedTerm {
			pp.ContractingPartyAlivePortion = 0
			pp.ContractingPartyAlivePortionAdjusted = 0
		} else {
			pp.ContractingPartyAlivePortion = 1
			pp.ContractingPartyAlivePortionAdjusted = 1
		}
		pp.ContractingPartyPolicyLapse = 0
		pp.ContractingPartyPolicyLapseAdjusted = 0
	} else if features.LapseDependentOnCpDeath && pp.ValuationTimeMonth <= parameters.CalculatedTerm {
		pp.ContractingPartyAlivePortion = utils.FloatPrecision(p.ContractingPartyAlivePortion*math.Pow(1-pp.BaseLapse, 1/12.0)*math.Pow(1-pp.MainMemberMortalityRate, 1/12.0), contractingPartyPrecision)
		pp.ContractingPartyAlivePortionAdjusted = utils.FloatPrecision(p.ContractingPartyAlivePortionAdjusted*math.Pow(1-pp.BaseLapseAdjusted, 1/12.0)*math.Pow(1-pp.MainMemberMortalityRateAdjusted, 1/12.0), contractingPartyPrecision)

	} else if pp.ValuationTimeMonth <= parameters.CalculatedTerm {
		pp.ContractingPartyAlivePortion = utils.FloatPrecision(p.ContractingPartyAlivePortion*math.Pow(1-pp.BaseLapse, 1/12.0), contractingPartyPrecision)
		pp.ContractingPartyAlivePortionAdjusted = utils.FloatPrecision(p.ContractingPartyAlivePortionAdjusted*math.Pow(1-pp.BaseLapseAdjusted, 1/12.0), contractingPartyPrecision)

	} else {
		pp.ContractingPartyAlivePortion = 0
		pp.ContractingPartyAlivePortionAdjusted = 0
	}

	if p.ContractingPartyAlivePortion == 0 || pp.ContractingPartyAlivePortion == 0 {
		pp.ContractingPartyPolicyLapse = 0
		pp.ContractingPartyPolicyLapseAdjusted = 0
	} else {
		if modelPoint.ContinuityOrPremiumWaiverOption && modelPoint.DurationInForceMonths+pp.ProjectionMonth > int(parameters.PremiumWaiverWaitingPeriod) {
			pp.ContractingPartyPolicyLapse = utils.FloatPrecision(1-math.Pow(1-pp.BaseLapse, 1/12.0), contractingPartyPrecision)
			pp.ContractingPartyPolicyLapseAdjusted = utils.FloatPrecision(1-math.Pow(1-pp.BaseLapseAdjusted, 1/12.0), contractingPartyPrecision)

		} else {
			if p.ContractingPartyAlivePortion > 0 {
				pp.ContractingPartyPolicyLapse = utils.FloatPrecision(1-pp.ContractingPartyAlivePortion/p.ContractingPartyAlivePortion, contractingPartyPrecision)
			}
			if p.ContractingPartyAlivePortionAdjusted > 0 {
				pp.ContractingPartyPolicyLapseAdjusted = utils.FloatPrecision(1-pp.ContractingPartyAlivePortionAdjusted/p.ContractingPartyAlivePortionAdjusted, contractingPartyPrecision)
			}
		}
	}
}

// PricingContractingPartyPolicyLapse computes policy annual lapse rate allowing for the dependency of the lapse rate on the contracting party's mortality rate and other transitions where relevant
func PricingContractingPartyPolicyLapse(pp *models.PricingPoint, p models.PricingPoint, modelPoint models.ProductPricingModelPoint, features models.ProductFeatures, parameters models.ProductPricingParameters) {
	if features.FuneralCover {
		if pp.ProjectionMonth == 0 { // This is the first loop of the run
			pp.ContractingPartyPolicyLapse = 0
			pp.ContractingPartyPolicyLapseAdjusted = 0
		} else {
			if features.LapseDependentOnCpDeath {
				if p.ContractingPartyAlivePortion == 0 || pp.ContractingPartyAlivePortion == 0 {
					pp.ContractingPartyPolicyLapse = 0
					pp.ContractingPartyPolicyLapseAdjusted = 0
				} else {
					if modelPoint.ContinuityOrPremiumWaiverOption && modelPoint.DurationInForceMonths+pp.ProjectionMonth > int(parameters.PremiumWaiverWaitingPeriod) {
						pp.ContractingPartyPolicyLapse = utils.FloatPrecision(1-math.Pow(1-pp.BaseLapse, 1/12.0), defaultPrecision)
						pp.ContractingPartyPolicyLapseAdjusted = utils.FloatPrecision(1-math.Pow(1-pp.BaseLapseAdjusted, 1/12.0), defaultPrecision)

					} else {
						pp.ContractingPartyPolicyLapse = utils.FloatPrecision(1-pp.ContractingPartyAlivePortion/p.ContractingPartyAlivePortion, defaultPrecision)
						pp.ContractingPartyPolicyLapseAdjusted = utils.FloatPrecision(1-pp.ContractingPartyAlivePortionAdjusted/p.ContractingPartyAlivePortionAdjusted, defaultPrecision)
					}
				}
			}
		}
	} else {
		pp.ContractingPartyPolicyLapse = 0
		pp.ContractingPartyPolicyLapseAdjusted = 0
	}
}

// PricingBaseMortalityRate reads annual base mortality rate by age and gender
func PricingBaseMortalityRate(pp *models.PricingPoint, mp models.ProductPricingModelPoint, productMargins models.ProductPricingMargins, states []models.ProductTransitionState, parameters models.ProductPricingParameters, columnname, tableName string, shock models.ProductPricingShock, pricingShockBasis string) {
	if utils.StatesContains(&states, Death) {
		if mp.MemberType == "MM" {
			pp.BaseMortalityRate = pp.MainMemberMortalityRate
			pp.BaseMortalityRateAdjusted = pp.MainMemberMortalityRateAdjusted
		} else if pp.ValuationTimeMonth <= parameters.CalculatedTerm {
			if pp.AgeNextBirthday >= 120 {
				pp.BaseMortalityRate = 1
				pp.BaseMortalityRateAdjusted = 1
			} else {
				args := models.TransitionRateArguments{
					ProductId:           pp.JobProductID,
					ProductCode:         pp.ProductCode,
					Age:                 pp.AgeNextBirthday,
					Gender:              mp.Gender,
					SmokerStatus:        mp.SmokerStatus,
					Income:              mp.Income,
					SocioEconomicClass:  mp.SocioEconomicClass,
					OccupationalClass:   mp.OccupationalClass,
					SelectPeriod:        mp.SelectPeriod,
					EducationLevel:      mp.EducationLevel,
					DurationIfM:         mp.DurationInForceMonths + pp.ProjectionMonth,
					ProjectionMonth:     pp.ProjectionMonth,
					DistributionChannel: mp.DistributionChannel,
					Year:                0000,
				}
				tempRate := getPricingMortalityRate(args, columnname, tableName)

				if pricingShockBasis != "N/A" {
					pp.BaseMortalityRate = math.Max(0, math.Min(1, tempRate*(1+shock.MultiplicativeMortality)+shock.AdditiveMortality))
				}
				if pricingShockBasis == "N/A" {
					pp.BaseMortalityRate = tempRate
				}

				pp.BaseMortalityRateAdjusted = math.Min(pp.BaseMortalityRate*(1+(productMargins.MortalityMargin)), 1)
			}
		} else {
			pp.BaseMortalityRate = 0
			pp.BaseMortalityRateAdjusted = 0
		}
	}
}

// PricingBaseRetrenchmentRate  reads annual independent retrenchment rate
func PricingBaseRetrenchmentRate(pp *models.PricingPoint, features models.ProductFeatures, mp models.ProductPricingModelPoint, margins models.ProductPricingMargins, states []models.ProductTransitionState, parameters models.ProductPricingParameters, columnname, tableName string, shock models.ProductPricingShock, pricingShockBasis string) {
	if pp.ValuationTimeMonth <= parameters.CalculatedTerm {
		if (mp.TemporaryPremiumWaiverIndicator && mp.TemporaryPremiumWaiverMonthExit > pp.ValuationTimeMonth) || mp.PremiumWaiverIndicator || mp.PaidupIndicator {
			pp.BaseIndependentRetrenchmentRate = 0
		} else {

			args := models.TransitionRateArguments{
				ProductId:           pp.JobProductID,
				ProductCode:         pp.ProductCode,
				Age:                 pp.AgeNextBirthday,
				Gender:              mp.Gender,
				SmokerStatus:        mp.SmokerStatus,
				Income:              mp.Income,
				SocioEconomicClass:  mp.SocioEconomicClass,
				OccupationalClass:   mp.OccupationalClass,
				SelectPeriod:        mp.SelectPeriod,
				EducationLevel:      mp.EducationLevel,
				DurationIfM:         mp.DurationInForceMonths + pp.ProjectionMonth,
				ProjectionMonth:     pp.ProjectionMonth,
				DistributionChannel: mp.DistributionChannel,
				Year:                0000,
			}
			var retrenchmentRate models.PricingRetrenchmentRate
			//yearInView := int(math.Ceil(pp.ValuationTimeYear))
			retrenchmentRate = getPricingRetrenchmentRate(args, columnname, tableName)

			if pricingShockBasis != "N/A" {
				pp.BaseIndependentRetrenchmentRate = math.Max(0, math.Min(1, retrenchmentRate.Value*(1+shock.MultiplicativeRetrenchment)+shock.AdditiveRetrenchment))
			}
			if pricingShockBasis == "N/A" {
				pp.BaseIndependentRetrenchmentRate = retrenchmentRate.Value
			}

			pp.BaseIndependentRetrenchmentRateAdjusted = pp.BaseIndependentRetrenchmentRate * (1 + (margins.RetrenchmentMargin))
		}
	} else {
		pp.BaseIndependentRetrenchmentRate = 0
		pp.BaseIndependentRetrenchmentRateAdjusted = 0
	}
}

// PricingBaseDisabilityIncrement reads annual base disability incidence rate by the respective rating factors
func PricingBaseDisabilityIncrement(pp *models.PricingPoint, mp models.ProductPricingModelPoint, features models.ProductFeatures, margins models.ProductPricingMargins, states []models.ProductTransitionState, parameters models.ProductPricingParameters, columnname, tableName string, shock models.ProductPricingShock, pricingShockBasis string) {
	if pp.ValuationTimeMonth <= parameters.CalculatedTerm {
		if (mp.TemporaryPremiumWaiverIndicator && mp.TemporaryPremiumWaiverMonthExit > pp.ValuationTimeMonth) || mp.PremiumWaiverIndicator || mp.PaidupIndicator || mp.MemberType != "MM" {
			pp.BaseIndependentDisabilityIncrement = 0
			pp.BaseIndependentDisabilityIncrementAdjusted = math.Min(pp.BaseIndependentDisabilityIncrement*(1+margins.MorbidityMargin), 1)
		} else {
			if pp.AgeNextBirthday >= 120 {
				pp.BaseIndependentDisabilityIncrement = 1
				pp.BaseIndependentDisabilityIncrementAdjusted = math.Min(pp.BaseIndependentDisabilityIncrement*(1+margins.MorbidityMargin), 1)
			} else {
				args := models.TransitionRateArguments{
					ProductId:           pp.JobProductID,
					ProductCode:         pp.ProductCode,
					Age:                 pp.AgeNextBirthday,
					Gender:              mp.Gender,
					SmokerStatus:        mp.SmokerStatus,
					Income:              mp.Income,
					SocioEconomicClass:  mp.SocioEconomicClass,
					OccupationalClass:   mp.OccupationalClass,
					SelectPeriod:        mp.SelectPeriod,
					EducationLevel:      mp.EducationLevel,
					DurationIfM:         mp.DurationInForceMonths + pp.ProjectionMonth,
					ProjectionMonth:     pp.ProjectionMonth,
					DistributionChannel: mp.DistributionChannel,
					Year:                0000,
				}
				tempRate := getPricingDisabilityIncidenceRate(args, columnname, tableName)
				if pricingShockBasis != "N/A" {
					pp.BaseIndependentDisabilityIncrement = math.Max(0, math.Min(1, tempRate*(1+shock.MultiplicativeDisability)+shock.AdditiveDisability))
				}
				if pricingShockBasis == "N/A" {
					pp.BaseIndependentDisabilityIncrement = tempRate
				}
				pp.BaseIndependentDisabilityIncrementAdjusted = math.Min(pp.BaseIndependentDisabilityIncrement*(1+margins.MorbidityMargin), 1)
			}
		}
	} else {
		pp.BaseIndependentDisabilityIncrement = 0
		pp.BaseIndependentDisabilityIncrementAdjusted = 0
	}
}

// PricingMainMemberMortalityRateByMonth converts annual independent mortality rate into monthly mortality rate
func PricingMainMemberMortalityRateByMonth(pp *models.PricingPoint, features models.ProductFeatures, parameters models.ProductPricingParameters) {
	if pp.ValuationTimeMonth <= parameters.CalculatedTerm {
		pp.MainMemberMortalityRateByMonth = utils.FloatPrecision(1-math.Pow(1-pp.MainMemberMortalityRate, 1/12.0), defaultPrecision)
		pp.MainMemberMortalityRateAdjustedByMonth = utils.FloatPrecision(1-math.Pow(1-pp.MainMemberMortalityRateAdjusted, 1/12.0), defaultPrecision)
		pp.IndependentMortalityRateMonthly = utils.FloatPrecision(1-math.Pow(1-pp.BaseMortalityRate, 1/12.0), defaultPrecision)
		pp.IndependentMortalityRateAdjustedByMonth = utils.FloatPrecision(1-math.Pow(1-pp.BaseMortalityRateAdjusted, 1/12.0), defaultPrecision)
	} else {
		pp.MainMemberMortalityRateByMonth = 0
		pp.MainMemberMortalityRateAdjustedByMonth = 0
		pp.IndependentMortalityRateMonthly = 0
		pp.IndependentMortalityRateAdjustedByMonth = 0
	}

}

// PricingIndependentLapseMonthly converts annual independent lapse rate into monthly lapse rate
func PricingIndependentLapseMonthly(pp *models.PricingPoint, p models.PricingPoint, pricingParameter models.PricingParameter, parameters models.ProductPricingParameters) {
	if pp.ValuationTimeMonth <= parameters.CalculatedTerm {
		pp.IndependentLapseMonthly = 1 - math.Pow(1-pp.BaseLapse, 1/12.0)
		pp.IndependentLapseMonthlyAdjusted = 1 - math.Pow(1-pp.BaseLapseAdjusted, 1/12.0)
	} else {
		pp.IndependentLapseMonthly = 0
		pp.IndependentLapseMonthlyAdjusted = 0
	}
}

// PricingIndependentRetrenchmentMonthly converts annual independent retrenchment rate into monthly independent retrenchment rate
func PricingIndependentRetrenchmentMonthly(pp *models.PricingPoint, parameters models.ProductPricingParameters) {
	if pp.ValuationTimeMonth <= parameters.CalculatedTerm {
		pp.IndependentRetrenchmentMonthly = 1 - math.Pow(1-pp.BaseIndependentRetrenchmentRate, 1/12.0)
		pp.IndependentRetrenchmentMonthlyAdjusted = 1 - math.Pow(1-pp.BaseIndependentRetrenchmentRateAdjusted, 1/12.0)
	} else {
		pp.IndependentRetrenchmentMonthly = 0
		pp.IndependentRetrenchmentMonthlyAdjusted = 0
	}
}

// PricingIndependentDisabilityMonthly converts annual independent disability incidence rate into monthly independent incidence rate
func PricingIndependentDisabilityMonthly(pp *models.PricingPoint, parameters models.ProductPricingParameters) {
	if pp.ValuationTimeMonth <= parameters.CalculatedTerm {
		pp.IndependentDisabilityMonthly = 1 - math.Pow(1-pp.BaseIndependentDisabilityIncrement, 1/12.0)
		pp.IndependentDisabilityMonthlyAdjusted = 1 - math.Pow(1-pp.BaseIndependentDisabilityIncrementAdjusted, 1/12.0)
	} else {
		pp.IndependentDisabilityMonthly = 0
		pp.IndependentDisabilityMonthlyAdjusted = 0
	}
}

// PricingMonthlyDependentMortality converts independent monthly mortality rate into dependent monthly mortality rate
func PricingMonthlyDependentMortality(pp *models.PricingPoint, parameters models.ProductPricingParameters) {
	if pp.ValuationTimeMonth <= parameters.CalculatedTerm {
		pp.MonthlyDependentMortality = pp.IndependentMortalityRateMonthly * (1 - TimingLapseZero*pp.IndependentLapseMonthly) * (1 - TimingDisabilityHalf*pp.IndependentDisabilityMonthly) * (1 - TimingRetrenchmentZero*pp.IndependentRetrenchmentMonthly)
		pp.MonthlyDependentMortalityAdjusted = pp.IndependentMortalityRateAdjustedByMonth * (1 - TimingLapseZero*pp.IndependentLapseMonthlyAdjusted) * (1 - TimingDisabilityHalf*pp.IndependentDisabilityMonthlyAdjusted) * (1 - TimingRetrenchmentZero*pp.IndependentRetrenchmentMonthlyAdjusted)
	} else {
		pp.MonthlyDependentMortality = 0
		pp.MonthlyDependentMortalityAdjusted = 0
	}
}

// PricingMonthlyDependentLapse converts independent monthly lapse rate into dependent monthly lapse rate
func PricingMonthlyDependentLapse(pp *models.PricingPoint, parameters models.ProductPricingParameters) {
	if pp.ValuationTimeMonth <= parameters.CalculatedTerm {
		pp.MonthlyDependentLapse = utils.FloatPrecision(pp.IndependentLapseMonthly*(1-TimingMortalityOne*pp.IndependentMortalityRateMonthly)*(1-TimingDisabilityOne*pp.IndependentDisabilityMonthly)*(1-TimingRetrenchmentOne*pp.IndependentRetrenchmentMonthly), defaultPrecision)
		pp.MonthlyDependentLapseAdjusted = utils.FloatPrecision(pp.IndependentLapseMonthlyAdjusted*(1-TimingMortalityOne*pp.IndependentMortalityRateAdjustedByMonth)*(1-TimingDisabilityOne*pp.IndependentDisabilityMonthlyAdjusted)*(1-TimingRetrenchmentOne*pp.IndependentRetrenchmentMonthlyAdjusted), defaultPrecision)
	} else {
		pp.MonthlyDependentLapse = 0
		pp.MonthlyDependentLapseAdjusted = 0
	}
}

// PricingMonthlyDependentRetrenchment converts independent monthly retrenchment rate into dependent monthly retrenchment rate
func PricingMonthlyDependentRetrenchment(pp *models.PricingPoint, parameters models.ProductPricingParameters) {
	if pp.ValuationTimeMonth <= parameters.CalculatedTerm {
		pp.MonthlyDependentRetrenchment = parameters.RetrenchmentQualificationFactor * pp.IndependentRetrenchmentMonthly * (1 - TimingLapseZero*pp.IndependentLapseMonthly) * (1 - TimingMortalityOne*pp.IndependentMortalityRateMonthly) * (1 - TimingDisabilityOne*pp.IndependentDisabilityMonthly)
		pp.MonthlyDependentRetrenchmentAdjusted = parameters.RetrenchmentQualificationFactor * pp.IndependentRetrenchmentMonthlyAdjusted * (1 - TimingLapseZero*pp.IndependentLapseMonthlyAdjusted) * (1 - TimingMortalityOne*pp.IndependentMortalityRateAdjustedByMonth) * (1 - TimingDisabilityOne*pp.IndependentDisabilityMonthlyAdjusted)
	} else {
		pp.MonthlyDependentRetrenchment = 0
		pp.MonthlyDependentRetrenchmentAdjusted = 0
	}
}

// PricingMonthlyDependentDisability converts independent monthly disability incidence rate into dependent monthly disability incidence rate
func PricingMonthlyDependentDisability(pp *models.PricingPoint, parameters models.ProductPricingParameters) {
	if pp.ValuationTimeMonth <= parameters.CalculatedTerm {
		pp.MonthlyDependentDisability = pp.IndependentDisabilityMonthly * (1 - TimingLapseZero*pp.IndependentLapseMonthly) * (1 - TimingMortalityHalf*pp.IndependentMortalityRateMonthly) * (1 - TimingRetrenchmentZero*pp.IndependentRetrenchmentMonthly)
		pp.MonthlyDependentDisabilityAdjusted = pp.IndependentDisabilityMonthlyAdjusted * (1 - TimingLapseZero*pp.IndependentLapseMonthlyAdjusted) * (1 - TimingMortalityHalf*pp.IndependentMortalityRateAdjustedByMonth) * (1 - TimingRetrenchmentZero*pp.IndependentRetrenchmentMonthlyAdjusted)
	} else {
		pp.MonthlyDependentDisability = 0
		pp.MonthlyDependentDisabilityAdjusted = 0
	}
}

// PricingIncrementalPaidUp computes number of new paid ups at each projection period
func PricingIncrementalPaidUp(pp *models.PricingPoint, p models.PricingPoint, pricingConfig models.PricingConfig, parameters models.ProductPricingParameters, features models.ProductFeatures) {
	if pricingConfig.PaidupOnSurvival {
		if pp.ProjectionMonth == 0 {
			//if pricingConfig.Paidup {
			//	pp.IncrementalPaidUp = 1
			//	pp.IncrementalPaidUpAdjusted = 1
			//
			//} else {
			pp.IncrementalPaidUp = 0
			pp.IncrementalPaidUp = 0
			//}
		} else {
			if pp.ValuationTimeMonth >= parameters.PaidUpOnSurvivalWaitingPeriod && pp.AgeNextBirthday >= parameters.PaidupEffectiveAge && pp.ValuationTimeMonth <= parameters.CalculatedTerm {
				if p.NumberPaidUp == 0 { // incremental should only happen at one point throughout the projection period, moving from health state to paidup state. thereafter the movement will be out of the paidup state
					pp.IncrementalPaidUp = p.InitialPolicy - pp.TotalIncrementalNaturalDeaths - pp.TotalIncrementalAccidentalDeaths
					pp.IncrementalPaidUpAdjusted = p.InitialPolicyAdjusted - pp.TotalIncrementalNaturalDeathsAdjusted - pp.TotalIncrementalAccidentalDeathsAdjusted
				}
			} else {
				pp.IncrementalPaidUp = 0
				pp.IncrementalPaidUpAdjusted = 0
			}
		}
	}
}

// PricingNaturalDeathsInForce computes cumulative number of natural deaths over the projection period
func PricingNaturalDeathsInForce(pp *models.PricingPoint, p models.PricingPoint, parameters models.ProductPricingParameters) {
	if pp.ProjectionMonth == 0 {
		pp.NaturalDeathsInForce = 0
		pp.NaturalDeathsInForceAdjusted = 0
	} else {
		//pp.NaturalDeathsInForce
		if pp.ValuationTimeMonth <= parameters.CalculatedTerm {
			pp.NaturalDeathsInForce = math.Max(p.InitialPolicy-pp.NumberOfMaturities, 0)*pp.MonthlyDependentMortality*(1-pp.AccidentProportion) + p.NaturalDeathsInForce
			pp.NaturalDeathsInForceAdjusted = math.Max(p.InitialPolicyAdjusted-pp.NumberOfMaturitiesAdjusted, 0)*pp.MonthlyDependentMortalityAdjusted*(1-pp.AccidentProportion) + p.NaturalDeathsInForceAdjusted
		} else {
			pp.NaturalDeathsInForce = p.NaturalDeathsInForce
			pp.NaturalDeathsInForceAdjusted = p.NaturalDeathsInForceAdjusted

		}
	}
}

// PricingIncrementalNaturalDeaths computes number of natural deaths at each projection period
//func PricingIncrementalNaturalDeaths(pp *models.PricingPoint, p models.PricingPoint, features models.ProductFeatures, parameters models.ProductPricingParameters) {
//	if pp.ProjectionMonth == 0 || pp.ValuationTimeMonth > parameters.CalculatedTerm {
//		pp.IncrementalNaturalDeaths = 0
//		pp.IncrementalNaturalDeathsAdjusted = 0
//
//	} else {
//		//IncNaturalDeathsAndAdjusted
//		pp.IncrementalNaturalDeaths = pp.NaturalDeathsInForce - p.NaturalDeathsInForce
//		pp.IncrementalNaturalDeathsAdjusted = pp.NaturalDeathsInForceAdjusted - p.NaturalDeathsInForceAdjusted
//	}
//}

// PricingNumberOfDeathsAccident computes cumulative number of accidental deaths over the projection period
func PricingNumberOfDeathsAccident(pp *models.PricingPoint, p models.PricingPoint, parameters models.ProductPricingParameters) {
	if pp.ProjectionMonth == 0 {
		pp.NumberOfDeathsAccident = 0
		pp.NumberOfDeathsAccidentAdjusted = 0
	} else {
		//pp.NumberOfAccidentDeaths
		if pp.ValuationTimeMonth <= parameters.CalculatedTerm {
			pp.NumberOfDeathsAccident = math.Max(p.InitialPolicy-pp.NumberOfMaturities, 0)*pp.MonthlyDependentMortality*pp.AccidentProportion + p.NumberOfDeathsAccident
			pp.NumberOfDeathsAccidentAdjusted = math.Max(p.InitialPolicyAdjusted-pp.NumberOfMaturitiesAdjusted, 0)*pp.MonthlyDependentMortalityAdjusted*pp.AccidentProportion + p.NumberOfDeathsAccidentAdjusted
		} else {
			pp.NumberOfDeathsAccident = p.NumberOfDeathsAccident
			pp.NumberOfDeathsAccidentAdjusted = p.NumberOfDeathsAccidentAdjusted

		}
	}
}

// PricingNumberOfLapses computes cumulative number of lapses over the projection period
func PricingNumberOfLapses(pp *models.PricingPoint, p models.PricingPoint, parameters models.ProductPricingParameters, modelPoint models.ProductPricingModelPoint, pricingconfig models.PricingConfig) {
	if pp.ProjectionMonth == 0 {
		pp.NumberOfLapses = 0
	} else {
		//pp.NumberOfLapses

		if (pp.MainMemberAgeNextBirthday >= parameters.PaidupEffectiveAge && pricingconfig.PaidupOnSurvival) || pp.ValuationTimeMonth > parameters.CalculatedTerm {
			pp.NumberOfLapses = p.NumberOfLapses
			pp.NumberOfLapsesAdjusted = p.NumberOfLapsesAdjusted
		} else {
			pp.NumberOfLapses = math.Max(p.InitialPolicy-pp.IncrementalPaidUp-pp.NumberOfMaturities, 0)*pp.MonthlyDependentLapse + p.NumberOfLapses
			pp.NumberOfLapsesAdjusted = math.Max(p.InitialPolicyAdjusted-pp.IncrementalPaidUpAdjusted-pp.NumberOfMaturitiesAdjusted, 0)*pp.MonthlyDependentLapseAdjusted + p.NumberOfLapsesAdjusted
		}
	}
}

// PricingNumberOfDisabilities computes cumulative number of disabilities over the projection period
func PricingNumberOfDisabilities(pp *models.PricingPoint, p models.PricingPoint, parameters models.ProductPricingParameters) {
	if pp.ProjectionMonth == 0 {
		pp.NumberOfDisabilities = 0
		pp.NumberOfDisabilitiesAdjusted = 0
	} else {
		if parameters.MainMemberIndicator && pp.ValuationTimeMonth <= parameters.CalculatedTerm {
			//pp.NumberOfDisabilities
			pp.NumberOfDisabilities = math.Max(p.InitialPolicy-pp.NumberOfMaturities, 0)*pp.MonthlyDependentDisability + p.NumberOfDisabilities
			pp.NumberOfDisabilitiesAdjusted = math.Max(p.InitialPolicyAdjusted-pp.NumberOfMaturitiesAdjusted, 0)*pp.MonthlyDependentDisabilityAdjusted + p.NumberOfDisabilitiesAdjusted

		} else {
			//pp.NumberOfDisabilities
			pp.NumberOfDisabilities = p.NumberOfDisabilities
			pp.NumberOfDisabilitiesAdjusted = p.NumberOfDisabilitiesAdjusted
		}
	}
}

// PricingNumberOfRetrenchments computes cumulative number of retrenchments over the projection period
func PricingNumberOfRetrenchments(pp *models.PricingPoint, p models.PricingPoint, parameters models.ProductPricingParameters) {
	if pp.ProjectionMonth == 0 {
		pp.NumberOfRetrenchments = 0
		pp.NumberOfRetrenchmentsAdjusted = 0
	} else {
		//pp.NumberOfRetrenchments
		if pp.ValuationTimeMonth <= parameters.CalculatedTerm {
			pp.NumberOfRetrenchments = math.Max(p.InitialPolicyInclRetrenchmentDecrement-pp.NumberOfMaturities, 0)*pp.MonthlyDependentRetrenchment + p.NumberOfRetrenchments
			pp.NumberOfRetrenchmentsAdjusted = math.Max(p.InitialPolicyInclRetrenchmentDecrementAdjusted-pp.NumberOfMaturitiesAdjusted, 0)*pp.MonthlyDependentRetrenchmentAdjusted + p.NumberOfRetrenchmentsAdjusted
		} else {
			pp.NumberOfRetrenchments = p.NumberOfRetrenchments
			pp.NumberOfRetrenchmentsAdjusted = p.NumberOfRetrenchmentsAdjusted
		}
	}
}

// PricingNaturalDeathsPremiumWaiver computes cumulative number of natural deaths, in the premium waiver state, over the projection period
func PricingNaturalDeathsPremiumWaiver(pp *models.PricingPoint, p models.PricingPoint, parameters models.ProductPricingParameters) {
	if pp.ProjectionMonth == 0 {
		pp.NaturalDeathsPremiumWaiver = 0
		pp.NaturalDeathsPremiumWaiverAdjusted = 0
	} else { //pp.NaturalDeathsPremiumWaiver
		if pp.ValuationTimeMonth <= parameters.CalculatedTerm {
			pp.NaturalDeathsPremiumWaiver = p.IncrementalPremiumWaivers*pp.IndependentMortalityRateMonthly*(1-pp.AccidentProportion) + p.NaturalDeathsPremiumWaiver
			pp.NaturalDeathsPremiumWaiverAdjusted = p.IncrementalPremiumWaiversAdjusted*pp.IndependentMortalityRateAdjustedByMonth*(1-pp.AccidentProportion) + p.NaturalDeathsPremiumWaiverAdjusted
		} else {
			pp.NaturalDeathsPremiumWaiver = p.NaturalDeathsPremiumWaiver
			pp.NaturalDeathsPremiumWaiverAdjusted = p.NaturalDeathsPremiumWaiverAdjusted
		}
	}
}

// PricingNaturalDeathsPremiumWaiver computes cumulative number of natural deaths, in the premium waiver state, over the projection period
func PricingNumberPaidup(pp *models.PricingPoint, p models.PricingPoint, parameters models.ProductPricingParameters, config models.PricingConfig, params models.PricingParameter) {
	if pp.ProjectionMonth == 0 {
		pp.NumberPaidUp = 0
		pp.NumberPaidUpAdjusted = 0
	} else {
		if pp.ValuationTimeMonth <= parameters.CalculatedTerm {
			pp.NumberPaidUp = p.NumberPaidUp*(1.0-pp.IndependentMortalityRateMonthly) + pp.IncrementalPaidUp
			pp.NumberPaidUpAdjusted = p.NumberPaidUpAdjusted*(1.0-pp.IndependentMortalityRateAdjustedByMonth) + pp.IncrementalPaidUpAdjusted

			if config.SpouseIndicator && config.ChildIndicator {
				if int(math.Ceil(pp.ValuationTimeYear))+params.AverageAgeAtEntryPerChild <= parameters.ChildExitAge {
					// to add one more scenario...
					pp.PostPaidupPolicyCount = p.PostPaidupPolicyCount*(1.0-pp.IndependentMortalityRateMonthly) + pp.IncrementalPaidUp*(1.0-p.SpouseSurvivalRate)*(1.0-pp.ChildSurvivalRate)
					pp.PostPaidupPolicyCountAdjusted = p.PostPaidupPolicyCountAdjusted*(1.0-pp.IndependentMortalityRateAdjustedByMonth) + pp.IncrementalPaidUpAdjusted*(1.0-p.SpouseSurvivalRateAdjusted)*(1.0-pp.ChildSurvivalRateAdjusted)
					pp.LastSurvivorPostPaidupPolicyCount = p.LastSurvivorPostPaidupPolicyCount*(1.0-pp.IndependentMortalityRateMonthly*p.IndependentSpouseMortalityRateByMonth*math.Pow(p.ChildIndependentMortalityMonthly, params.NumberChildPerPolicy)) + pp.IncrementalPaidUp*pp.SpouseSurvivalRate
					pp.LastSurvivorPostPaidupPolicyCountAdjusted = p.LastSurvivorPostPaidupPolicyCountAdjusted*(1.0-pp.IndependentMortalityRateAdjustedByMonth*p.IndependentSpouseMortalityRateAdjustedByMonth*math.Pow(p.ChildIndependentMortalityAdjustedMonthly, params.NumberChildPerPolicy)) + pp.IncrementalPaidUpAdjusted*pp.SpouseSurvivalRateAdjusted

				} else {
					pp.PostPaidupPolicyCount = p.PostPaidupPolicyCount*(1.0-pp.IndependentMortalityRateMonthly) + pp.IncrementalPaidUp*(1.0-pp.SpouseSurvivalRate)
					pp.PostPaidupPolicyCountAdjusted = p.PostPaidupPolicyCountAdjusted*(1.0-pp.IndependentMortalityRateAdjustedByMonth) + pp.IncrementalPaidUpAdjusted*(1.0-pp.SpouseSurvivalRateAdjusted)
					pp.LastSurvivorPostPaidupPolicyCount = p.LastSurvivorPostPaidupPolicyCount*(1.0-pp.IndependentMortalityRateMonthly*p.IndependentSpouseMortalityRateByMonth) + pp.IncrementalPaidUp*pp.SpouseSurvivalRate
					pp.LastSurvivorPostPaidupPolicyCountAdjusted = p.LastSurvivorPostPaidupPolicyCountAdjusted*(1.0-pp.IndependentMortalityRateAdjustedByMonth*p.IndependentSpouseMortalityRateAdjustedByMonth) + pp.IncrementalPaidUpAdjusted*pp.SpouseSurvivalRateAdjusted

				}
			}

			if config.SpouseIndicator {
				//probability of a policy qualifying for a paidup status
				//and both MM and a Spouse survived(Last Survivor) or Only MM survived
				pp.PostPaidupPolicyCount = p.PostPaidupPolicyCount*(1.0-pp.IndependentMortalityRateMonthly) + pp.IncrementalPaidUp*(1.0-pp.SpouseSurvivalRate)
				pp.PostPaidupPolicyCountAdjusted = p.PostPaidupPolicyCountAdjusted*(1.0-pp.IndependentMortalityRateAdjustedByMonth) + pp.IncrementalPaidUpAdjusted*(1.0-pp.SpouseSurvivalRateAdjusted)
				pp.LastSurvivorPostPaidupPolicyCount = p.LastSurvivorPostPaidupPolicyCount*(1.0-pp.IndependentMortalityRateMonthly*pp.IndependentSpouseMortalityRateByMonth) + pp.IncrementalPaidUp*pp.SpouseSurvivalRate
				pp.LastSurvivorPostPaidupPolicyCountAdjusted = p.LastSurvivorPostPaidupPolicyCountAdjusted*(1.0-pp.IndependentMortalityRateAdjustedByMonth*pp.IndependentSpouseMortalityRateAdjustedByMonth) + pp.IncrementalPaidUpAdjusted*pp.SpouseSurvivalRateAdjusted
			}

			if config.ChildIndicator {
				if int(math.Ceil(pp.ValuationTimeYear))+params.AverageAgeAtEntryPerChild <= parameters.ChildExitAge {
					pp.PostPaidupPolicyCount = p.PostPaidupPolicyCount*(1.0-pp.IndependentMortalityRateMonthly*pp.ChildIndependentMortalityMonthly) + pp.IncrementalPaidUp
					pp.PostPaidupPolicyCountAdjusted = p.PostPaidupPolicyCountAdjusted*(1.0-pp.IndependentMortalityRateAdjustedByMonth*pp.ChildIndependentMortalityAdjustedMonthly) + pp.IncrementalPaidUpAdjusted
					pp.LastSurvivorPostPaidupPolicyCount = 0
					pp.LastSurvivorPostPaidupPolicyCountAdjusted = 0
				} else {
					pp.PostPaidupPolicyCount = p.PostPaidupPolicyCount*(1.0-pp.IndependentMortalityRateMonthly) + pp.IncrementalPaidUp
					pp.PostPaidupPolicyCountAdjusted = p.PostPaidupPolicyCountAdjusted*(1.0-pp.IndependentMortalityRateAdjustedByMonth) + pp.IncrementalPaidUpAdjusted
					pp.LastSurvivorPostPaidupPolicyCount = 0
					pp.LastSurvivorPostPaidupPolicyCountAdjusted = 0
				}
			}

			if !config.SpouseIndicator && !config.ChildIndicator {

				pp.PostPaidupPolicyCount = p.PostPaidupPolicyCount*(1.0-pp.IndependentMortalityRateMonthly) + pp.IncrementalPaidUp
				pp.PostPaidupPolicyCountAdjusted = p.PostPaidupPolicyCountAdjusted*(1.0-pp.IndependentMortalityRateAdjustedByMonth) + pp.IncrementalPaidUpAdjusted
				pp.LastSurvivorPostPaidupPolicyCount = 0
				pp.LastSurvivorPostPaidupPolicyCountAdjusted = 0
			}
		} else {
			pp.NumberPaidUp = p.NumberPaidUp
			pp.NumberPaidUpAdjusted = p.NumberPaidUpAdjusted
			pp.PostPaidupPolicyCount = 0
			pp.PostPaidupPolicyCount = 0
			pp.LastSurvivorPostPaidupPolicyCount = 0
			pp.LastSurvivorPostPaidupPolicyCountAdjusted = 0
		}
	}
}

// PricingNaturalDeathsTemporaryWaivers computes cumulative number of natural deaths, in the temporary premium waiver state, over the projection period
func PricingNaturalDeathsTemporaryWaivers(pp *models.PricingPoint, p models.PricingPoint, parameters models.ProductPricingParameters) {
	if pp.ProjectionMonth == 0 {
		pp.NaturalDeathsTemporaryWaivers = 0
		pp.NaturalDeathsTemporaryWaiversAdjusted = 0

	} else {
		//pp.NaturalDeathsTemporaryWaivers
		if pp.ValuationTimeMonth <= parameters.CalculatedTerm {
			pp.NaturalDeathsTemporaryWaivers = p.InitialTemporaryPremiumWaivers*pp.IndependentMortalityRateMonthly*(1-pp.AccidentProportion) + p.NaturalDeathsTemporaryWaivers
			pp.NaturalDeathsTemporaryWaiversAdjusted = p.InitialTemporaryPremiumWaiversAdjusted*pp.IndependentMortalityRateMonthly*(1-pp.AccidentProportion) + p.NaturalDeathsTemporaryWaiversAdjusted
		} else {
			pp.NaturalDeathsTemporaryWaivers = p.NaturalDeathsTemporaryWaivers
			pp.NaturalDeathsTemporaryWaiversAdjusted = p.NaturalDeathsTemporaryWaiversAdjusted
		}
	}
}

// PricingAccidentDeathsPremiumWaiver computes cumulative number of accidental deaths, in the premium waiver state, over the projection period
func PricingAccidentDeathsPremiumWaiver(pp *models.PricingPoint, p models.PricingPoint, parameters models.ProductPricingParameters) {
	if pp.ProjectionMonth == 0 {
		pp.AccidentDeathsPremiumWaiver = 0
		pp.AccidentDeathsPremiumWaiverAdjusted = 0
	} else {
		//pp.AccidentDeathsPremiumWaiver
		if pp.ValuationTimeMonth <= parameters.CalculatedTerm {
			pp.AccidentDeathsPremiumWaiver = p.IncrementalPremiumWaivers*pp.IndependentMortalityRateMonthly*pp.AccidentProportion + p.AccidentDeathsPremiumWaiver
			pp.AccidentDeathsPremiumWaiverAdjusted = p.IncrementalPremiumWaiversAdjusted*pp.IndependentMortalityRateAdjustedByMonth*pp.AccidentProportion + p.AccidentDeathsPremiumWaiverAdjusted
		} else {
			pp.AccidentDeathsPremiumWaiver = 0
			pp.AccidentDeathsPremiumWaiverAdjusted = 0
		}
	}
}

// PricingAccidentDeathsTemporaryPremiumWaiver computes cumulative number of accidental deaths, in the temporary premium waiver state, over the projection period
func PricingAccidentDeathsTemporaryPremiumWaiver(pp *models.PricingPoint, p models.PricingPoint, parameters models.ProductPricingParameters) {
	if pp.ProjectionMonth == 0 {
		pp.AccidentDeathsTemporaryPremiumWaiver = 0
		pp.AccidentDeathsTemporaryPremiumWaiverAdjusted = 0
	} else {
		//pp.AccidentDeathsTemporaryPremiumWaiver
		if pp.ValuationTimeMonth <= parameters.CalculatedTerm {
			pp.AccidentDeathsTemporaryPremiumWaiver = p.InitialTemporaryPremiumWaivers*pp.IndependentMortalityRateMonthly*pp.AccidentProportion + p.AccidentDeathsTemporaryPremiumWaiver
			pp.AccidentDeathsTemporaryPremiumWaiverAdjusted = p.InitialTemporaryPremiumWaiversAdjusted*pp.IndependentMortalityRateAdjustedByMonth*pp.AccidentProportion + p.AccidentDeathsTemporaryPremiumWaiverAdjusted
		} else {
			pp.AccidentDeathsTemporaryPremiumWaiver = p.AccidentDeathsTemporaryPremiumWaiver
			pp.AccidentDeathsTemporaryPremiumWaiverAdjusted = p.AccidentDeathsTemporaryPremiumWaiverAdjusted
		}
	}
}

// PricingLapseIncrements computes number of new lapses at each projection period
func PricingTotalIncrementalLapses(pp *models.PricingPoint, p models.PricingPoint, parameters models.ProductPricingParameters) {
	if pp.ProjectionMonth == 0 || pp.ValuationTimeMonth > parameters.CalculatedTerm {
		pp.TotalIncrementalLapses = 0
		pp.TotalIncrementalLapsesAdjusted = 0
	} else {
		//pp.TotalIncrementalLapses
		pp.TotalIncrementalLapses = pp.NumberOfLapses - p.NumberOfLapses
		pp.TotalIncrementalLapsesAdjusted = pp.NumberOfLapsesAdjusted - p.NumberOfLapsesAdjusted
	}
}

// PricingNumberOfMaturities computes cumulative number of maturities over the projection period
func PricingNumberOfMaturities(pp *models.PricingPoint, p models.PricingPoint, modelPoint models.ProductPricingModelPoint, features models.ProductFeatures, parameters models.ProductPricingParameters) {
	if pp.ProjectionMonth > 0 {
		if pp.ValuationTimeMonth == parameters.CalculatedTerm+1 && parameters.CalculatedTerm != 0 {
			pp.NumberOfMaturities = math.Max(startingInitialPolicy-p.NaturalDeathsInForce-p.AccidentDeathsInForce-p.NumberOfLapses-p.NumberOfDisabilities-p.NumberOfRetrenchments, 0) + p.NumberOfMaturities
			pp.NumberOfMaturitiesAdjusted = math.Max(startingInitialPolicy-p.NaturalDeathsInForceAdjusted-p.AccidentDeathsInForceAdjusted-p.NumberOfLapsesAdjusted-p.NumberOfDisabilitiesAdjusted-p.NumberOfRetrenchmentsAdjusted, 0) + p.NumberOfMaturitiesAdjusted
		} else {
			pp.NumberOfMaturities = p.NumberOfMaturities
			pp.NumberOfMaturitiesAdjusted = p.NumberOfMaturitiesAdjusted
		}
	} else {
		pp.NumberOfMaturities = 0
		pp.NumberOfMaturitiesAdjusted = 0
		if parameters.CalculatedTerm == 0 {
			pp.NumberOfMaturities = 1
			pp.NumberOfMaturitiesAdjusted = 1
		}
	}
}

// PricingInitialPolicy computes number of policies in the full premium paying healthy state over the projection period
func PricingInitialPolicy(pp *models.PricingPoint, modelPoint models.ProductPricingModelPoint, startingInitialPolicy *float64, startingInitialPolicyAdjusted *float64, pricingConfig models.PricingConfig, parameters models.ProductPricingParameters) {
	if pp.ProjectionMonth == 0 || pp.ValuationTimeMonth > parameters.CalculatedTerm {
		if pp.ValuationTimeMonth > parameters.CalculatedTerm {
			pp.InitialPolicy = 0
			pp.InitialPolicyAdjusted = 0
			pp.InitialPolicyInclRetrenchmentDecrement = 0
			pp.InitialPolicyInclRetrenchmentDecrementAdjusted = 0
		} else {
			pp.InitialPolicy = 1
			pp.InitialPolicyAdjusted = 1
			pp.InitialPolicyInclRetrenchmentDecrement = 1
			pp.InitialPolicyInclRetrenchmentDecrementAdjusted = 1
		}
		*startingInitialPolicy = pp.InitialPolicy
		*startingInitialPolicyAdjusted = pp.InitialPolicyAdjusted
	} else {
		pp.InitialPolicy = utils.FloatPrecision(math.Max(
			*startingInitialPolicy-
				pp.NumberOfMaturities-
				pp.NumberPaidUp-
				pp.NaturalDeathsInForce-
				pp.NaturalDeathsPaidUp-
				pp.NumberOfDeathsAccident-
				pp.NumberOfAccidentDeathsPaidUp-
				pp.NumberOfLapses-
				pp.NumberOfDisabilities,
			//pp.IncrementalPremiumWaivers-
			//pp.NaturalDeathsPremiumWaiver-
			//pp.NaturalDeathsTemporaryWaivers-
			//pp.AccidentDeathsPremiumWaiver-
			//pp.AccidentDeathsTemporaryPremiumWaiver,
			0), defaultdecrementPrecision)
		pp.InitialPolicyAdjusted = utils.FloatPrecision(math.Max(
			*startingInitialPolicyAdjusted-
				pp.NumberOfMaturities-
				pp.NumberPaidUpAdjusted-
				pp.NaturalDeathsInForceAdjusted-
				pp.NaturalDeathsPaidUpAdjusted-
				pp.NumberOfDeathsAccidentAdjusted-
				pp.NumberOfAccidentDeathsPaidUp-
				pp.NumberOfLapsesAdjusted-
				pp.NumberOfDisabilitiesAdjusted,
			//pp.IncrementalPremiumWaiversAdjusted-
			//pp.NaturalDeathsPremiumWaiverAdjusted-
			//pp.NaturalDeathsTemporaryWaiversAdjusted-
			//pp.AccidentDeathsPremiumWaiverAdjusted-
			//pp.AccidentDeathsTemporaryPremiumWaiverAdjusted,
			0), defaultdecrementPrecision)

		pp.InitialPolicyInclRetrenchmentDecrement = utils.FloatPrecision(math.Max(
			*startingInitialPolicy-
				pp.NumberOfMaturities-
				pp.NumberPaidUp-
				pp.NaturalDeathsInForce-
				pp.NaturalDeathsPaidUp-
				pp.NumberOfDeathsAccident-
				pp.NumberOfAccidentDeathsPaidUp-
				pp.NumberOfLapses-
				pp.NumberOfDisabilities-
				pp.NumberOfRetrenchments,
			//pp.IncrementalPremiumWaivers-
			//pp.NaturalDeathsPremiumWaiver-
			//pp.NaturalDeathsTemporaryWaivers-
			//pp.AccidentDeathsPremiumWaiver-
			//pp.AccidentDeathsTemporaryPremiumWaiver,
			0), defaultdecrementPrecision)
		pp.InitialPolicyInclRetrenchmentDecrementAdjusted = utils.FloatPrecision(math.Max(
			*startingInitialPolicyAdjusted-
				pp.NumberOfMaturities-
				pp.NumberPaidUpAdjusted-
				pp.NaturalDeathsInForceAdjusted-
				pp.NaturalDeathsPaidUpAdjusted-
				pp.NumberOfDeathsAccidentAdjusted-
				pp.NumberOfAccidentDeathsPaidUp-
				pp.NumberOfLapsesAdjusted-
				pp.NumberOfDisabilitiesAdjusted-
				pp.NumberOfRetrenchmentsAdjusted,
			//pp.IncrementalPremiumWaiversAdjusted-
			//pp.NaturalDeathsPremiumWaiverAdjusted-
			//pp.NaturalDeathsTemporaryWaiversAdjusted-
			//pp.AccidentDeathsPremiumWaiverAdjusted-
			//pp.AccidentDeathsTemporaryPremiumWaiverAdjusted,
			0), defaultdecrementPrecision)
	}
}

// PricingTotalIncrementalNaturalDeaths computes number of new natural deaths at each projection period
func PricingTotalIncrementalNaturalDeaths(pp *models.PricingPoint, p models.PricingPoint, parameters models.ProductPricingParameters) {
	if pp.ProjectionMonth == 0 || pp.ValuationTimeMonth > parameters.CalculatedTerm {
		pp.TotalIncrementalNaturalDeaths = 0
		pp.TotalIncrementalNaturalDeathsAdjusted = 0
	} else {
		pp.TotalIncrementalNaturalDeaths = pp.NaturalDeathsInForce - p.NaturalDeathsInForce
		pp.TotalIncrementalNaturalDeathsAdjusted = pp.NaturalDeathsInForceAdjusted - p.NaturalDeathsInForceAdjusted

	}
}

func PricingTotalIncrementalDisabilities(pp *models.PricingPoint, p models.PricingPoint, parameters models.ProductPricingParameters) {
	if pp.ProjectionMonth == 0 || pp.ValuationTimeMonth > parameters.CalculatedTerm {
		pp.TotalIncrementalDisabilities = 0
		pp.TotalIncrementalDisabilitiesAdjusted = 0
	} else {
		pp.TotalIncrementalDisabilities = utils.FloatPrecision(pp.NumberOfDisabilities-p.NumberOfDisabilities, defaultdecrementPrecision)
		pp.TotalIncrementalDisabilitiesAdjusted = utils.FloatPrecision(pp.NumberOfDisabilitiesAdjusted-p.NumberOfDisabilitiesAdjusted, defaultdecrementPrecision)
	}
}

func PricingTotalIncrementalRetrenchments(pp *models.PricingPoint, p models.PricingPoint, parameters models.ProductPricingParameters) {
	if pp.ProjectionMonth == 0 || pp.ValuationTimeMonth > parameters.CalculatedTerm {
		pp.TotalIncrementalRetrenchments = 0
		pp.TotalIncrementalRetrenchmentsAdjusted = 0
	} else {
		pp.TotalIncrementalRetrenchments = utils.FloatPrecision(pp.NumberOfRetrenchments-p.NumberOfRetrenchments, defaultdecrementPrecision)
		pp.TotalIncrementalRetrenchmentsAdjusted = utils.FloatPrecision(pp.NumberOfRetrenchmentsAdjusted-p.NumberOfRetrenchmentsAdjusted, defaultdecrementPrecision)
	}
}

// PricingIncrementalPremiumWaivers computes number of new premium waivers at each projection period
func PricingIncrementalPremiumWaivers(pp *models.PricingPoint, pricingConfig models.PricingConfig, params models.ProductPricingParameters, features models.ProductFeatures) {
	if pricingConfig.PremiumWaiverOnDeath {
		var premiumwaiverSumAssuredAdjustmentFac float64
		if pp.ProjectionMonth == 0 {
			//if pricingConfig.Paidup {
			//	pp.IncrementalPremiumWaivers = 1
			//	pp.IncrementalPremiumWaiversAdjusted = 1
			//
			//} else {
			pp.IncrementalPremiumWaivers = 0
			pp.IncrementalPremiumWaiversAdjusted = 0
			//}
		} else {
			if pp.ValuationTimeMonth > params.PremiumWaiverWaitingPeriod {
				if pp.ValuationTimeMonth <= (params.PremiumWaiverWaitingPeriod + int(params.PremiumWaiverAdjustedSumassuredTerm)) {
					premiumwaiverSumAssuredAdjustmentFac = params.PremiumWaiverSumAssuredFactor
				} else {
					premiumwaiverSumAssuredAdjustmentFac = 1
				}
			}
			if pp.ValuationTimeMonth > params.PremiumWaiverWaitingPeriod && pp.ValuationTimeMonth <= params.CalculatedTerm {
				pp.IncrementalPremiumWaivers = (pp.TotalIncrementalNaturalDeaths + pp.TotalIncrementalAccidentalDeaths) * premiumwaiverSumAssuredAdjustmentFac // premiumwaiverSumAssuredAdjustmentFac cover for remaining lives is adjusted if main member's death happens between t and t+x
				pp.IncrementalPremiumWaiversAdjusted = (pp.TotalIncrementalNaturalDeathsAdjusted + pp.TotalIncrementalAccidentalDeathsAdjusted) * premiumwaiverSumAssuredAdjustmentFac
			} else {
				pp.IncrementalPremiumWaivers = 0
				pp.IncrementalPremiumWaiversAdjusted = 0
			}
		}

	}
}

// PricingAccidentDeathsPaidUp computes cumulative number of accidental deaths, in the paid up state, over the projection period
func PricingNaturalDeathsPaidUp(pp *models.PricingPoint, p models.PricingPoint, parameters models.ProductPricingParameters) {
	if pp.ProjectionMonth == 0 || pp.ValuationTimeMonth > parameters.CalculatedTerm {
		pp.NaturalDeathsPaidUp = 0
		pp.NaturalDeathsPaidUpAdjusted = 0
		pp.TotalIncrementalPaidupNaturalDeaths = 0
		pp.TotalIncrementalPaidupNaturalDeathsAdjusted = 0
		return
	}
	pp.NaturalDeathsPaidUp = p.NumberPaidUp*pp.IndependentMortalityRateMonthly*(1-pp.AccidentProportion) + p.NaturalDeathsPaidUp
	pp.NaturalDeathsPaidUpAdjusted = p.NumberPaidUpAdjusted*pp.IndependentMortalityRateAdjustedByMonth*(1-pp.AccidentProportion) + p.NaturalDeathsPaidUpAdjusted
	pp.TotalIncrementalPaidupNaturalDeaths = pp.NaturalDeathsPaidUp - p.NaturalDeathsPaidUp
	pp.TotalIncrementalPaidupNaturalDeathsAdjusted = pp.NaturalDeathsPaidUpAdjusted - p.NaturalDeathsPaidUpAdjusted
}

// PricingAccidentDeathsPaidUp computes cumulative number of accidental deaths, in the paid up state, over the projection period
func PricingAccidentDeathsPaidUp(pp *models.PricingPoint, p models.PricingPoint, parameters models.ProductPricingParameters) {
	if pp.ProjectionMonth == 0 || pp.ValuationTimeMonth > parameters.CalculatedTerm {
		pp.NumberOfAccidentDeathsPaidUp = 0
		pp.NumberOfAccidentDeathsPaidUpAdjusted = 0
		pp.TotalIncrementalPaidupAccidentalDeaths = 0
		pp.TotalIncrementalPaidupAccidentalDeathsAdjusted = 0
		return
	}
	pp.NumberOfAccidentDeathsPaidUp = p.NumberPaidUp*pp.IndependentMortalityRateMonthly*pp.AccidentProportion + p.NumberOfAccidentDeathsPaidUp
	pp.NumberOfAccidentDeathsPaidUpAdjusted = p.NumberPaidUpAdjusted*pp.IndependentMortalityRateAdjustedByMonth*pp.AccidentProportion + p.NumberOfAccidentDeathsPaidUpAdjusted
	pp.TotalIncrementalPaidupAccidentalDeaths = pp.NumberOfAccidentDeathsPaidUp - p.NumberOfAccidentDeathsPaidUp
	pp.TotalIncrementalPaidupAccidentalDeathsAdjusted = pp.NumberOfAccidentDeathsPaidUpAdjusted - p.NumberOfAccidentDeathsPaidUpAdjusted
}

// PricingTotalIncrementalAccidentalDeaths computes number of new accidental deaths over the projection period
func PricingTotalIncrementalAccidentalDeaths(pp *models.PricingPoint, mp models.ProductPricingModelPoint, p models.PricingPoint, parameters models.ProductPricingParameters) {
	if pp.ProjectionMonth == 0 || pp.ValuationTimeMonth > parameters.CalculatedTerm {
		pp.TotalIncrementalAccidentalDeaths = 0
		pp.TotalIncrementalAccidentalDeathsAdjusted = 0

	} else {
		pp.TotalIncrementalAccidentalDeaths = pp.NumberOfDeathsAccident - p.NumberOfDeathsAccident
		pp.TotalIncrementalAccidentalDeathsAdjusted = pp.NumberOfDeathsAccidentAdjusted - p.NumberOfDeathsAccidentAdjusted
	}
}

// PricingSpouseAgeNextBirthday estimates spouses age next birthday from the main member's age next birthday
func PricingSpouseAgeNextBirthday(pp *models.PricingPoint, mp models.ProductPricingModelPoint, features models.ProductFeatures) {
	// Spouse Age
	if features.FuneralCover {
		if mp.MainMemberGender == Female {
			//pp.SpouseAgeNextBirthday = int(math.Min(math.Max(float64(pp.MainMemberAgeNextBirthday+2*0), 18), 65))
			pp.SpouseAgeNextBirthday = int(math.Ceil(math.Min(math.Max(float64(mp.MainMemberAgeAtEntry+2), 18), 65) + pp.ValuationTimeYear))
		} else {
			//pp.SpouseAgeNextBirthday = int(math.Min(math.Max(float64(pp.MainMemberAgeNextBirthday-2*0), 18), 65))
			pp.SpouseAgeNextBirthday = int(math.Ceil(math.Min(math.Max(float64(mp.MainMemberAgeAtEntry-2), 18), 65) + pp.ValuationTimeYear))
		}
	}
}

// PricingSpouseMortalityRate reads spouse's mortality rate, at each projection period, by the respective rating factors
func PricingSpouseMortalityRate(pp *models.PricingPoint, mp models.ProductPricingModelPoint, config models.PricingConfig, margins models.ProductPricingMargins, columnname, tableName string) {
	// Spouse mortalities
	if pp.MainMemberAgeNextBirthday < 120 && config.SpouseIndicator {
		if mp.MainMemberGender == Male {
			args := models.TransitionRateArguments{
				ProductId:           pp.JobProductID,
				ProductCode:         pp.ProductCode,
				Age:                 pp.SpouseAgeNextBirthday,
				Gender:              "F",
				SmokerStatus:        mp.SmokerStatus,
				Income:              mp.Income,
				SocioEconomicClass:  mp.SocioEconomicClass,
				SelectPeriod:        mp.SelectPeriod,
				EducationLevel:      mp.EducationLevel,
				DurationIfM:         mp.DurationInForceMonths,
				ProjectionMonth:     pp.ProjectionMonth,
				DistributionChannel: mp.DistributionChannel,
				Year:                0000,
			}
			pp.SpouseMortalityRate = getPricingMortalityRate(args, columnname, tableName)
		} else if mp.MainMemberGender == Female {
			args := models.TransitionRateArguments{
				ProductId:           pp.JobProductID,
				ProductCode:         pp.ProductCode,
				Age:                 pp.SpouseAgeNextBirthday,
				Gender:              "M",
				SmokerStatus:        mp.SmokerStatus,
				Income:              mp.Income,
				SocioEconomicClass:  mp.SocioEconomicClass,
				SelectPeriod:        mp.SelectPeriod,
				EducationLevel:      mp.EducationLevel,
				DurationIfM:         mp.DurationInForceMonths,
				ProjectionMonth:     pp.ProjectionMonth,
				DistributionChannel: mp.DistributionChannel,
				Year:                0000,
			}
			pp.SpouseMortalityRate = getPricingMortalityRate(args, columnname, tableName)
		}

		var respSpouse float64 = 0
		if pp.MainMemberAgeNextBirthday > 120 {
			respSpouse = 1
		} else {
			respSpouse = pp.SpouseMortalityRate * (1 + margins.MortalityMargin)
		}
		pp.SpouseMortalityRateAdjusted = utils.FloatPrecision(math.Min(respSpouse, 1), defaultPrecision)
	} else if pp.MainMemberAgeNextBirthday == 120 && config.SpouseIndicator {
		pp.SpouseMortalityRate = 1
		pp.SpouseMortalityRateAdjusted = 1
	} else {
		pp.SpouseMortalityRate = 0
		pp.SpouseMortalityRateAdjusted = 0
	}
}

// PricingBaseSpouseIndependentLapse converts spouse's annual lapse rate into monthly lapse rate
func PricingBaseSpouseIndependentLapse(pp *models.PricingPoint, config models.PricingConfig) {
	if config.SpouseIndicator {
		pp.BaseSpouseIndependentLapse = utils.FloatPrecision(1-math.Pow(1-pp.ContractingPartyPolicyLapse, 12.0), defaultPrecision)
		pp.BaseSpouseIndependentLapseAdjusted = utils.FloatPrecision(1-math.Pow(1-pp.ContractingPartyPolicyLapseAdjusted, 12.0), defaultPrecision)
	}
}

// PricingIndependentSpouseMortalityRateByMonth converts independent annual mortality rate into independent monthly mortality rate
func PricingIndependentSpouseMortalityRateByMonth(pp *models.PricingPoint) {
	// Spouse Monthly Independent MortalityRate
	pp.IndependentSpouseMortalityRateByMonth = utils.FloatPrecision(1-math.Pow(1-pp.SpouseMortalityRate, 1/12.0), defaultPrecision)
	pp.IndependentSpouseMortalityRateAdjustedByMonth = utils.FloatPrecision(1-math.Pow(1-pp.SpouseMortalityRateAdjusted, 1/12.0), defaultPrecision)
}

// PricingIndependentSpouseLapseMonthly converts spouse's annual lapse rate into monthly lapse rate
func PricingIndependentSpouseLapseMonthly(pp *models.PricingPoint) {
	// Spouse Monthly Independent Lapse
	pp.IndependentSpouseLapseMonthly = 1 - math.Pow(1-pp.BaseSpouseIndependentLapse, 1/12.0)
	pp.IndependentSpouseLapseMonthlyAdjusted = 1 - math.Pow(1-pp.BaseSpouseIndependentLapseAdjusted, 1/12.0)
}

// PricingMonthlySpouseDependentMortality converts spouse's independent monthly mortality rate into dependent monthly mortality rate
func PricingMonthlySpouseDependentMortality(pp *models.PricingPoint) {
	// Spouse Monthly Dependent Mortality
	pp.MonthlySpouseDependentMortality = utils.FloatPrecision(pp.IndependentSpouseMortalityRateByMonth*(1-TimingLapseZero*pp.IndependentSpouseLapseMonthly), defaultPrecision)
	pp.MonthlySpouseDependentMortalityAdjusted = utils.FloatPrecision(pp.IndependentSpouseMortalityRateAdjustedByMonth*(1-TimingLapseZero*pp.IndependentSpouseLapseMonthlyAdjusted), defaultPrecision)
}

// PricingMonthlySpouseDependentLapse converts spouse's independent monthly lapse rate into dependent monthly lapse rate
func PricingMonthlySpouseDependentLapse(pp *models.PricingPoint) {
	// Spouse Monthly Dependent Lapse
	pp.MonthlySpouseDependentLapse = utils.FloatPrecision(pp.IndependentSpouseLapseMonthly*(1-TimingMortalityOne*pp.IndependentSpouseMortalityRateByMonth), defaultPrecision)
	pp.MonthlySpouseDependentLapseAdjusted = utils.FloatPrecision(pp.IndependentSpouseLapseMonthlyAdjusted*(1-TimingMortalityOne*pp.IndependentSpouseMortalityRateAdjustedByMonth), defaultPrecision)
}

// PricingSpouseNumberOfPaidUps computes spouse's cumulative number of paid ups over the projection period
func PricingSpouseNumberOfPaidUps(pp *models.PricingPoint, p models.PricingPoint, config models.PricingConfig) {
	// Spouse Initial Paid Up
	if pp.ProjectionMonth == 0 {
		if config.PaidupOnSurvival {
			pp.SpouseNumberOfPaidUps = 0
			pp.SpouseNumberOfPaidUpsAdjusted = 0
		}
	} else {
		pp.SpouseNumberOfPaidUps = pp.IncrementalPaidUp*p.SpouseSurvivalRate + p.SpouseNumberOfPaidUps*(1-pp.SpouseMortalityRate)
		pp.SpouseNumberOfPaidUpsAdjusted = pp.IncrementalPaidUpAdjusted*p.SpouseSurvivalRateAdjusted + p.SpouseNumberOfPaidUpsAdjusted*(1-pp.SpouseMortalityRateAdjusted)
	}
}

// PricingSpouseNumberOfPremiumWaivers computes spouse's cumulative number of premium waivers over the projection period
func PricingSpouseNumberOfPremiumWaivers(pp *models.PricingPoint, p models.PricingPoint, config models.PricingConfig) {
	// Spouse initial PremiumWaiver
	if pp.ProjectionMonth == 0 {
		if config.PremiumWaiverOnDeath {
			pp.SpouseNumberOfPremiumWaivers = 0
			pp.SpouseNumberOfPremiumWaiversAdjusted = 0
		}
	} else {
		pp.SpouseNumberOfPremiumWaivers = pp.IncrementalPremiumWaivers + p.SpouseNumberOfPremiumWaivers*(1-pp.SpouseMortalityRate)
		pp.SpouseNumberOfPremiumWaiversAdjusted = pp.IncrementalPremiumWaiversAdjusted + p.SpouseNumberOfPremiumWaiversAdjusted*(1-pp.SpouseMortalityRateAdjusted)
	}
}

// PricingSpouseNumberPolicies computes number of spouses over the projection period
func PricingSpouseNumberPolicies(pp *models.PricingPoint, p models.PricingPoint, config models.PricingConfig) {
	// Spouse Number of Policies
	if config.SpouseIndicator {
		if pp.ProjectionMonth == 0 {
			pp.SpouseNumberPolicies = 1
			pp.SpouseNumberPoliciesAdjusted = 1
			pp.SpouseSurvivalRate = 1
			pp.SpouseSurvivalRateAdjusted = 1
		} else {
			if config.PaidupOnSurvival {
				if pp.IncrementalPaidUp == 0 {
					pp.SpouseNumberPolicies = utils.FloatPrecision(p.SpouseNumberPolicies*math.Max(1-pp.MonthlySpouseDependentLapse-pp.MonthlySpouseDependentMortality, 0), defaultdecrementPrecision)
					pp.SpouseNumberPoliciesAdjusted = utils.FloatPrecision(p.SpouseNumberPoliciesAdjusted*math.Max(1-pp.MonthlySpouseDependentLapseAdjusted-pp.MonthlySpouseDependentMortalityAdjusted, 0), defaultdecrementPrecision)
				}
				if pp.IncrementalPaidUp > 0 {
					pp.SpouseNumberPolicies = 0
					pp.SpouseNumberPoliciesAdjusted = 0
				}
			}
			if !config.PaidupOnSurvival {
				pp.SpouseNumberPolicies = utils.FloatPrecision(p.SpouseNumberPolicies*math.Max(1-pp.MonthlySpouseDependentLapse-pp.MonthlySpouseDependentMortality, 0), defaultdecrementPrecision)
				pp.SpouseNumberPoliciesAdjusted = utils.FloatPrecision(p.SpouseNumberPoliciesAdjusted*math.Max(1-pp.MonthlySpouseDependentLapseAdjusted-pp.MonthlySpouseDependentMortalityAdjusted, 0), defaultdecrementPrecision)
			}
			pp.SpouseSurvivalRate = utils.FloatPrecision(p.SpouseSurvivalRate*math.Max(1-pp.MonthlySpouseDependentMortality, 0), defaultdecrementPrecision)
			pp.SpouseSurvivalRateAdjusted = utils.FloatPrecision(p.SpouseSurvivalRateAdjusted*math.Max(1-pp.MonthlySpouseDependentMortalityAdjusted, 0), defaultdecrementPrecision)
		}
	}
}

// PricingTotalSpouseIncrementalNaturalDeaths computes cumulative number of spouse's new deaths at each projection period
func PricingTotalSpouseIncrementalNaturalDeaths(pp *models.PricingPoint, p models.PricingPoint, mp models.ProductPricingModelPoint, columnname, tableName string) {
	gender := ""
	if mp.MainMemberGender == Female {
		gender = Male
	} else {
		gender = Female
	}
	args := models.TransitionRateArguments{
		ProductId:           pp.JobProductID,
		ProductCode:         pp.ProductCode,
		Age:                 pp.SpouseAgeNextBirthday,
		Gender:              gender,
		SmokerStatus:        mp.SmokerStatus,
		Income:              mp.Income,
		SocioEconomicClass:  mp.SocioEconomicClass,
		SelectPeriod:        mp.SelectPeriod,
		EducationLevel:      mp.EducationLevel,
		DurationIfM:         mp.DurationInForceMonths,
		ProjectionMonth:     pp.ProjectionMonth,
		DistributionChannel: mp.DistributionChannel,
		Year:                0000,
	}

	if pp.ProjectionMonth == 0 {
		pp.TotalSpouseIncrementalNaturalDeaths = 0
		pp.TotalSpouseIncrementalNaturalDeathsAdjusted = 0

	} else {
		pp.TotalSpouseIncrementalNaturalDeaths = p.SpouseNumberPolicies * pp.MonthlySpouseDependentMortality * (1 - getPricingMortalityRateAccidentProportion(args, columnname, tableName))
		pp.TotalSpouseIncrementalNaturalDeathsAdjusted = p.SpouseNumberPoliciesAdjusted * pp.MonthlySpouseDependentMortalityAdjusted * (1 - getPricingMortalityRateAccidentProportion(args, columnname, tableName))
	}
}

// PricingSpouseNumberOfPaidUpNaturalDeaths computes cumulative number spouse's natural deaths over the projection period
func PricingSpouseNumberOfPaidUpNaturalDeaths(pp *models.PricingPoint, p models.PricingPoint, mp models.ProductPricingModelPoint, params models.ProductPricingParameters, columnname, tableName string) {
	gender := ""
	if mp.MainMemberGender == Female {
		gender = Male
	} else {
		gender = Female
	}
	args := models.TransitionRateArguments{
		ProductId:           pp.JobProductID,
		ProductCode:         pp.ProductCode,
		Age:                 pp.SpouseAgeNextBirthday,
		Gender:              gender,
		SmokerStatus:        mp.SmokerStatus,
		Income:              mp.Income,
		SocioEconomicClass:  mp.SocioEconomicClass,
		SelectPeriod:        mp.SelectPeriod,
		EducationLevel:      mp.EducationLevel,
		DurationIfM:         mp.DurationInForceMonths,
		ProjectionMonth:     pp.ProjectionMonth,
		DistributionChannel: mp.DistributionChannel,
		Year:                0000,
	}

	if pp.ProjectionMonth == 0 {
		pp.SpouseNumberOfPaidUpNaturalDeaths = 0
		pp.SpouseNumberOfPaidUpNaturalDeathsAdjusted = 0
		pp.TotalSpouseIncrementalPaidupNaturalDeaths = 0
		pp.TotalSpouseIncrementalPaidupNaturalDeathsAdjusted = 0
	} else {

		//Spouse Number of Paid Up Deaths and Adjusted
		pp.SpouseNumberOfPaidUpNaturalDeaths = p.SpouseNumberOfPaidUps*pp.SpouseMortalityRate*(1-getPricingMortalityRateAccidentProportion(args, columnname, tableName)) + p.SpouseNumberOfPaidUpNaturalDeaths                                 //+ pp.SpouseNumberOfPremiumWaivers*pp.SpouseMortalityRate*(1-getPricingMortalityRateAccidentProportion(args, columnname, tableName))*params.PremiumWaiverBenefitPayoutAdjusted
		pp.SpouseNumberOfPaidUpNaturalDeathsAdjusted = p.SpouseNumberOfPaidUpsAdjusted*pp.SpouseMortalityRateAdjusted*(1-getPricingMortalityRateAccidentProportion(args, columnname, tableName)) + p.SpouseNumberOfPaidUpNaturalDeathsAdjusted //+ pp.SpouseNumberOfPremiumWaiversAdjusted*pp.SpouseMortalityRateAdjusted*(1-getPricingMortalityRateAccidentProportion(args, columnname, tableName))*params.PremiumWaiverBenefitPayoutAdjusted
		pp.TotalSpouseIncrementalPaidupNaturalDeaths = pp.SpouseNumberOfPaidUpNaturalDeaths - p.SpouseNumberOfPaidUpNaturalDeaths
		pp.TotalSpouseIncrementalPaidupNaturalDeathsAdjusted = pp.SpouseNumberOfPaidUpNaturalDeathsAdjusted - p.SpouseNumberOfPaidUpNaturalDeathsAdjusted
	}
}

// PricingSpouseNumberOfPaidUpAccidentalDeaths computes cumulative number of spouse's accidental deaths over the projection period
func PricingSpouseNumberOfPaidUpAccidentalDeaths(pp *models.PricingPoint, p models.PricingPoint, mp models.ProductPricingModelPoint, params models.ProductPricingParameters, columnname, tableName string) {

	gender := ""
	if mp.MainMemberGender == Female {
		gender = Male
	} else {
		gender = Female
	}
	args := models.TransitionRateArguments{
		ProductId:           pp.JobProductID,
		ProductCode:         pp.ProductCode,
		Age:                 pp.SpouseAgeNextBirthday,
		Gender:              gender,
		SmokerStatus:        mp.SmokerStatus,
		Income:              mp.Income,
		SocioEconomicClass:  mp.SocioEconomicClass,
		SelectPeriod:        mp.SelectPeriod,
		EducationLevel:      mp.EducationLevel,
		DurationIfM:         mp.DurationInForceMonths,
		ProjectionMonth:     pp.ProjectionMonth,
		DistributionChannel: mp.DistributionChannel,
		Year:                0000,
	}
	// Spouse Total Natural Deaths and PaidUps
	if pp.ProjectionMonth == 0 {
		pp.SpouseNumberOfPaidUpAccidentalDeaths = 0
		pp.SpouseNumberOfPaidUpAccidentalDeathsAdjusted = 0
		pp.TotalSpouseIncrementalPaidupAccidentalDeaths = 0
		pp.TotalSpouseIncrementalPaidupAccidentalDeathsAdjusted = 0
	} else {
		pp.SpouseNumberOfPaidUpAccidentalDeaths = p.SpouseNumberOfPaidUps*pp.SpouseMortalityRate*(getPricingMortalityRateAccidentProportion(args, columnname, tableName)) + p.SpouseNumberOfPaidUpAccidentalDeaths                                 //+ pp.SpouseNumberOfPremiumWaivers*pp.SpouseMortalityRate*(getPricingMortalityRateAccidentProportion(args, columnname, tableName))*params.PremiumWaiverBenefitPayoutAdjusted
		pp.SpouseNumberOfPaidUpAccidentalDeathsAdjusted = p.SpouseNumberOfPaidUpsAdjusted*pp.SpouseMortalityRateAdjusted*(getPricingMortalityRateAccidentProportion(args, columnname, tableName)) + p.SpouseNumberOfPaidUpAccidentalDeathsAdjusted //+ pp.SpouseNumberOfPremiumWaiversAdjusted*pp.SpouseMortalityRateAdjusted*(getPricingMortalityRateAccidentProportion(args, columnname, tableName))*params.PremiumWaiverBenefitPayoutAdjusted
		pp.TotalSpouseIncrementalPaidupAccidentalDeaths = pp.SpouseNumberOfPaidUpAccidentalDeaths - p.SpouseNumberOfPaidUpAccidentalDeaths
		pp.TotalSpouseIncrementalPaidupAccidentalDeathsAdjusted = pp.SpouseNumberOfPaidUpAccidentalDeathsAdjusted - p.SpouseNumberOfPaidUpAccidentalDeathsAdjusted
	}
}

// PricingSpouseNumberOfPremiumWaiversNaturalDeaths computes cumulative number spouse's natural deaths over the projection period
func PricingSpouseNumberOfPremiumWaiverNaturalDeaths(pp *models.PricingPoint, p models.PricingPoint, mp models.ProductPricingModelPoint, params models.ProductPricingParameters, columnname, tableName string) {
	gender := ""
	if mp.MainMemberGender == Female {
		gender = Male
	} else {
		gender = Female
	}
	args := models.TransitionRateArguments{
		ProductId:           pp.JobProductID,
		ProductCode:         pp.ProductCode,
		Age:                 pp.SpouseAgeNextBirthday,
		Gender:              gender,
		SmokerStatus:        mp.SmokerStatus,
		Income:              mp.Income,
		SocioEconomicClass:  mp.SocioEconomicClass,
		SelectPeriod:        mp.SelectPeriod,
		EducationLevel:      mp.EducationLevel,
		DurationIfM:         mp.DurationInForceMonths,
		ProjectionMonth:     pp.ProjectionMonth,
		DistributionChannel: mp.DistributionChannel,
		Year:                0000,
	}

	if pp.ProjectionMonth == 0 {
		pp.SpouseNumberOfPremiumWaiverNaturalDeaths = 0
		pp.SpouseNumberOfPremiumWaiverNaturalDeathsAdjusted = 0
		pp.TotalSpouseIncrementalPwNaturalDeaths = 0
		pp.TotalSpouseIncrementalPwNaturalDeathsAdjusted = 0
	} else {

		//Spouse Number of PremiumWaivers Deaths and Adjusted
		pp.SpouseNumberOfPremiumWaiverNaturalDeaths = p.SpouseNumberOfPremiumWaivers*pp.SpouseMortalityRate*(1-getPricingMortalityRateAccidentProportion(args, columnname, tableName)) + p.SpouseNumberOfPremiumWaiverNaturalDeaths
		pp.SpouseNumberOfPremiumWaiverNaturalDeathsAdjusted = p.SpouseNumberOfPremiumWaiversAdjusted*pp.SpouseMortalityRateAdjusted*(1-getPricingMortalityRateAccidentProportion(args, columnname, tableName)) + p.SpouseNumberOfPremiumWaiverNaturalDeathsAdjusted
		pp.TotalSpouseIncrementalPwNaturalDeaths = pp.SpouseNumberOfPremiumWaiverNaturalDeaths - p.SpouseNumberOfPremiumWaiverNaturalDeaths
		pp.TotalSpouseIncrementalPwNaturalDeathsAdjusted = pp.SpouseNumberOfPremiumWaiverNaturalDeathsAdjusted - p.SpouseNumberOfPremiumWaiverNaturalDeathsAdjusted
	}
}

// PricingSpouseNumberOfPaidUpNaturalDeaths computes cumulative number spouse's natural deaths over the projection period
func PricingSpouseNumberOfPremiumWaiverAccidentalDeaths(pp *models.PricingPoint, p models.PricingPoint, mp models.ProductPricingModelPoint, params models.ProductPricingParameters, columnname, tableName string) {
	gender := ""
	if mp.MainMemberGender == Female {
		gender = Male
	} else {
		gender = Female
	}
	args := models.TransitionRateArguments{
		ProductId:           pp.JobProductID,
		ProductCode:         pp.ProductCode,
		Age:                 pp.SpouseAgeNextBirthday,
		Gender:              gender,
		SmokerStatus:        mp.SmokerStatus,
		Income:              mp.Income,
		SocioEconomicClass:  mp.SocioEconomicClass,
		SelectPeriod:        mp.SelectPeriod,
		EducationLevel:      mp.EducationLevel,
		DurationIfM:         mp.DurationInForceMonths,
		ProjectionMonth:     pp.ProjectionMonth,
		DistributionChannel: mp.DistributionChannel,
		Year:                0000,
	}

	if pp.ProjectionMonth == 0 {
		pp.SpouseNumberOfPremiumWaiverNaturalDeaths = 0
		pp.SpouseNumberOfPremiumWaiverNaturalDeathsAdjusted = 0
		pp.TotalSpouseIncrementalPwAccidentalDeaths = 0
		pp.TotalSpouseIncrementalPwAccidentalDeathsAdjusted = 0
	} else {
		pp.SpouseNumberOfPremiumWaiverAccidentalDeaths = p.SpouseNumberOfPremiumWaivers*pp.SpouseMortalityRate*getPricingMortalityRateAccidentProportion(args, columnname, tableName) + p.SpouseNumberOfPremiumWaiverAccidentalDeaths
		pp.SpouseNumberOfPremiumWaiveAccidentalDeathsAdjusted = p.SpouseNumberOfPremiumWaiversAdjusted*pp.SpouseMortalityRateAdjusted*getPricingMortalityRateAccidentProportion(args, columnname, tableName) + p.SpouseNumberOfPremiumWaiveAccidentalDeathsAdjusted
		pp.TotalSpouseIncrementalPwAccidentalDeaths = pp.SpouseNumberOfPremiumWaiverAccidentalDeaths - p.SpouseNumberOfPremiumWaiverAccidentalDeaths
		pp.TotalSpouseIncrementalPwAccidentalDeathsAdjusted = pp.TotalSpouseIncrementalPwAccidentalDeathsAdjusted - p.TotalSpouseIncrementalPwAccidentalDeathsAdjusted
	}
}

// PricingTotalSpouseIncrementalAccidentalDeaths computes spouse's number of new accidental deaths at each projection period
func PricingTotalSpouseIncrementalAccidentalDeaths(pp *models.PricingPoint, p models.PricingPoint, mp models.ProductPricingModelPoint, columnname, tableName string) {
	// Spouse total incremental Accident death

	gender := ""
	if mp.MainMemberGender == Female {
		gender = Male
	} else {
		gender = Female
	}
	args := models.TransitionRateArguments{
		ProductId:           pp.JobProductID,
		ProductCode:         pp.ProductCode,
		Age:                 pp.SpouseAgeNextBirthday,
		Gender:              gender,
		SmokerStatus:        mp.SmokerStatus,
		Income:              mp.Income,
		SocioEconomicClass:  mp.SocioEconomicClass,
		SelectPeriod:        mp.SelectPeriod,
		EducationLevel:      mp.EducationLevel,
		DurationIfM:         mp.DurationInForceMonths,
		ProjectionMonth:     pp.ProjectionMonth,
		DistributionChannel: mp.DistributionChannel,
		Year:                0000,
	}
	if pp.ProjectionMonth == 0 {
		pp.TotalSpouseIncrementalAccidentalDeaths = 0
		pp.TotalSpouseIncrementalAccidentalDeathsAdjusted = 0
	} else {
		pp.TotalSpouseIncrementalAccidentalDeaths = p.SpouseNumberPolicies * pp.MonthlySpouseDependentMortality * (getPricingMortalityRateAccidentProportion(args, columnname, tableName))
		pp.TotalSpouseIncrementalAccidentalDeathsAdjusted = p.SpouseNumberPoliciesAdjusted * pp.MonthlySpouseDependentMortalityAdjusted * (getPricingMortalityRateAccidentProportion(args, columnname, tableName))
	}
}

// PricingChildMortalityRate reads child mortality rate, at each projection period, by rating factors
func PricingChildMortalityRate(pp *models.PricingPoint, mp models.ProductPricingModelPoint, config models.PricingConfig, pricingParameter models.PricingParameter, margins models.ProductPricingMargins, columnname, tableName string) {
	// child Mortalities
	if config.ChildIndicator {
		if int(math.Ceil(pp.ValuationTimeYear))+pricingParameter.AverageAgeAtEntryPerChild >= 120 {
			pp.ChildMortalityRate = 1
			pp.ChildMortalityRateAdjusted = 1
		} else {
			argsmale := models.TransitionRateArguments{
				ProductId:           pp.JobProductID,
				ProductCode:         pp.ProductCode,
				Age:                 int(math.Ceil(pp.ValuationTimeYear)) + pricingParameter.AverageAgeAtEntryPerChild,
				Gender:              "M",
				SmokerStatus:        mp.SmokerStatus,
				Income:              mp.Income,
				SocioEconomicClass:  mp.SocioEconomicClass,
				SelectPeriod:        mp.SelectPeriod,
				EducationLevel:      mp.EducationLevel,
				DurationIfM:         mp.DurationInForceMonths,
				ProjectionMonth:     pp.ProjectionMonth,
				DistributionChannel: mp.DistributionChannel,
				Year:                0000,
			}
			argsfemale := models.TransitionRateArguments{
				ProductId:           pp.JobProductID,
				ProductCode:         pp.ProductCode,
				Age:                 int(math.Ceil(pp.ValuationTimeYear)) + pricingParameter.AverageAgeAtEntryPerChild,
				Gender:              "F",
				SmokerStatus:        mp.SmokerStatus,
				Income:              mp.Income,
				SocioEconomicClass:  mp.SocioEconomicClass,
				SelectPeriod:        mp.SelectPeriod,
				EducationLevel:      mp.EducationLevel,
				DurationIfM:         mp.DurationInForceMonths,
				ProjectionMonth:     pp.ProjectionMonth,
				DistributionChannel: mp.DistributionChannel,
				Year:                0000,
			}
			pp.ChildMortalityRate = getPricingMortalityRate(argsmale, columnname, tableName)*pricingParameter.DistributionMale + getPricingMortalityRate(argsfemale, columnname, tableName)*pricingParameter.DistributionFemale
			pp.ChildMortalityRateAdjusted = pp.ChildMortalityRate * (1 + margins.MortalityMargin)
		}
	}
}

// PricingChildAccidentalProportion  reads the proportion of the base mortality table that arises from non-natural causes
func PricingChildAccidentalProportion(pp *models.PricingPoint, mp models.ProductPricingModelPoint, config models.PricingConfig, pricingParameter models.PricingParameter, columnname, tableName string) {
	// child Mortalities

	if config.ChildIndicator {
		argsmale := models.TransitionRateArguments{
			ProductId:           pp.JobProductID,
			ProductCode:         pp.ProductCode,
			Age:                 int(math.Ceil(pp.ValuationTimeYear)) + pricingParameter.AverageAgeAtEntryPerChild,
			Gender:              "M",
			SmokerStatus:        mp.SmokerStatus,
			Income:              mp.Income,
			SocioEconomicClass:  mp.SocioEconomicClass,
			SelectPeriod:        mp.SelectPeriod,
			EducationLevel:      mp.EducationLevel,
			DurationIfM:         mp.DurationInForceMonths,
			ProjectionMonth:     pp.ProjectionMonth,
			DistributionChannel: mp.DistributionChannel,
			Year:                0000,
		}
		argsfemale := models.TransitionRateArguments{
			ProductId:           pp.JobProductID,
			ProductCode:         pp.ProductCode,
			Age:                 int(math.Ceil(pp.ValuationTimeYear)) + pricingParameter.AverageAgeAtEntryPerChild,
			Gender:              "F",
			SmokerStatus:        mp.SmokerStatus,
			Income:              mp.Income,
			SocioEconomicClass:  mp.SocioEconomicClass,
			SelectPeriod:        mp.SelectPeriod,
			EducationLevel:      mp.EducationLevel,
			DurationIfM:         mp.DurationInForceMonths,
			ProjectionMonth:     pp.ProjectionMonth,
			DistributionChannel: mp.DistributionChannel,
			Year:                0000,
		}
		pp.ChildAccidentalProportion = getPricingMortalityRateAccidentProportion(argsmale, columnname, tableName)*pricingParameter.DistributionMale + getPricingMortalityRateAccidentProportion(argsfemale, columnname, tableName)*pricingParameter.DistributionFemale
	} else {
		pp.ChildAccidentalProportion = 0
	}
}

// PricingChildMonthlyMortalityRate converts child's annual mortality rate into monthly mortality rate
func PricingIndependentChildMonthlyMortalityRate(pp *models.PricingPoint, config models.PricingConfig) {
	// child Mortalities
	if config.ChildIndicator && pp.ProjectionMonth > 0 {
		pp.ChildIndependentMortalityMonthly = 1 - math.Pow(1-pp.ChildMortalityRate, 1/12.0)
		pp.ChildIndependentMortalityAdjustedMonthly = 1 - math.Pow(1-pp.ChildMortalityRateAdjusted, 1/12.0)
	} else {
		pp.ChildIndependentMortalityMonthly = 0
		pp.ChildIndependentMortalityAdjustedMonthly = 0
	}
}

func PricingMonthlyChildDependentLapse(pp *models.PricingPoint) {
	// Child Monthly Dependent Lapse
	pp.MonthlyChildDependentLapse = utils.FloatPrecision(pp.ChildIndependentLapseMonthly*(1-TimingMortalityOne*pp.ChildIndependentMortalityMonthly), defaultPrecision)
	pp.MonthlyChildDependentLapseAdjusted = utils.FloatPrecision(pp.ChildIndependentLapseAdjustedMonthly*(1-TimingMortalityOne*pp.ChildIndependentMortalityAdjustedMonthly), defaultPrecision)
}

func PricingMonthlyChildDependentMortality(pp *models.PricingPoint) {
	// Child Monthly Dependent Mortality
	pp.MonthlyChildDependentMortality = utils.FloatPrecision(pp.ChildIndependentMortalityMonthly*(1-TimingLapseZero*pp.ChildIndependentLapseMonthly), defaultPrecision)
	pp.MonthlyChildDependentMortalityAdjusted = utils.FloatPrecision(pp.ChildIndependentMortalityAdjustedMonthly*(1-TimingLapseZero*pp.ChildIndependentLapseAdjustedMonthly), defaultPrecision)
}

func PricingChildIndependentLapse(pp *models.PricingPoint, config models.PricingConfig, parameters models.ProductPricingParameters) {
	if pp.ProjectionMonth == 0 || pp.ValuationTimeMonth > parameters.CalculatedTerm {
		pp.ChildIndependentLapse = 0
		pp.ChildIndependentLapseAdjusted = 0
	} else {
		if config.ChildIndicator {
			pp.ChildIndependentLapse = utils.FloatPrecision(1-math.Pow(1-pp.ContractingPartyPolicyLapse, 12.0), defaultPrecision)
			pp.ChildIndependentLapseAdjusted = utils.FloatPrecision(1-math.Pow(1-pp.ContractingPartyPolicyLapseAdjusted, 12.0), defaultPrecision)
		}
	}
}

func PricingIndependentChildLapseMonthly(pp *models.PricingPoint, config models.PricingConfig, parameters models.ProductPricingParameters) {
	if pp.ProjectionMonth == 0 || pp.ValuationTimeMonth > parameters.CalculatedTerm {
		pp.ChildIndependentLapseMonthly = 0
		pp.ChildIndependentLapseAdjustedMonthly = 0
	} else {
		if config.ChildIndicator {
			pp.ChildIndependentLapseMonthly = 1 - math.Pow(1-pp.ChildIndependentLapse, 1/12.0)
			pp.ChildIndependentLapseAdjustedMonthly = 1 - math.Pow(1-pp.ChildIndependentLapseAdjusted, 1/12.0)
		}
	}
}

// PricingChildNumberOfPaidUps computes cumulative number of child's paid ups over the projection period
func PricingChildNumberOfPaidUps(pp *models.PricingPoint, p models.PricingPoint, config models.PricingConfig, pricingParameter models.PricingParameter, params models.ProductPricingParameters) {
	// Child Initial Paid up
	if pp.ProjectionMonth == 0 {
		if config.PaidupOnSurvival {
			pp.ChildNumberOfPaidUps = 0
			pp.ChildNumberOfPaidUpsAdjusted = 0
		}
	} else {
		if int(math.Ceil(pp.ValuationTimeYear))+pricingParameter.AverageAgeAtEntryPerChild <= params.ChildExitAge {
			pp.ChildNumberOfPaidUps = pp.IncrementalPaidUp*p.ChildSurvivalRate + p.ChildNumberOfPaidUps*(1-pp.ChildIndependentMortalityMonthly)
			pp.ChildNumberOfPaidUpsAdjusted = pp.IncrementalPaidUpAdjusted*p.ChildSurvivalRateAdjusted + p.ChildNumberOfPaidUpsAdjusted*(1-pp.ChildIndependentMortalityAdjustedMonthly)
		} else {
			pp.ChildNumberOfPaidUps = 0
			pp.ChildNumberOfPaidUpsAdjusted = 0
		}
	}
}

// PricingChildNumberOfPremiumWaivers computes cumulative number of child's premium waiver over the projection period
func PricingChildNumberOfPremiumWaivers(pp *models.PricingPoint, p models.PricingPoint, config models.PricingConfig, pricingParameter models.PricingParameter, params models.ProductPricingParameters) {
	// child initial premium waiver
	if pp.ProjectionMonth == 0 {
		if config.PaidupOnSurvival {
			pp.ChildNumberOfPremiumWaivers = 0
			pp.ChildNumberOfPremiumWaiversAdjusted = 0
		}
	} else {
		if int(math.Ceil(pp.ValuationTimeYear))+pricingParameter.AverageAgeAtEntryPerChild <= params.ChildExitAge {
			pp.ChildNumberOfPremiumWaivers = pp.IncrementalPremiumWaivers + p.ChildNumberOfPremiumWaivers*(1-pp.ChildMortalityRate)
			pp.ChildNumberOfPremiumWaiversAdjusted = pp.IncrementalPremiumWaiversAdjusted + p.ChildNumberOfPremiumWaiversAdjusted*(1-pp.ChildMortalityRateAdjusted)

		} else {
			pp.ChildNumberOfPremiumWaivers = 0
			pp.ChildNumberOfPremiumWaiversAdjusted = 0
		}
	}
}

// PricingChildNumberPolicies computes number of child's over the projection period allow for decrements
func PricingChildNumberPolicies(pp *models.PricingPoint, p models.PricingPoint, config models.PricingConfig, pricingParameter models.PricingParameter, params models.ProductPricingParameters) {
	// Child Initial Policies
	if config.ChildIndicator {
		if pp.ProjectionMonth == 0 {
			pp.ChildNumberPolicies = pricingParameter.NumberChildPerPolicy
			pp.ChildNumberPoliciesAdjusted = pp.ChildNumberPolicies
			pp.ChildSurvivalRate = pp.ChildNumberPolicies
			pp.ChildSurvivalRateAdjusted = pp.ChildNumberPolicies
		} else {
			if int(math.Ceil(pp.ValuationTimeYear))+pricingParameter.AverageAgeAtEntryPerChild <= params.ChildExitAge {
				pp.ChildNumberPolicies = p.ChildNumberPolicies * math.Max(1-pp.MonthlyChildDependentMortality-pp.MonthlyChildDependentLapse, 0)
				pp.ChildNumberPoliciesAdjusted = p.ChildNumberPoliciesAdjusted * math.Max(1-pp.MonthlyChildDependentMortalityAdjusted-pp.MonthlyChildDependentLapseAdjusted, 0)
				pp.ChildSurvivalRate = p.ChildSurvivalRate * math.Max(1-pp.MonthlyChildDependentMortality, 0)
				pp.ChildSurvivalRateAdjusted = p.ChildSurvivalRateAdjusted * math.Max(1-pp.MonthlyChildDependentMortalityAdjusted, 0)

			} else {
				pp.ChildNumberPolicies = 0
				pp.ChildNumberPoliciesAdjusted = 0
				pp.ChildSurvivalRate = 0
				pp.ChildSurvivalRateAdjusted = 0
			}
		}
	}

}

// PricingTotalChildIncrementalNaturalDeaths computes child's number of new natural deaths at each projection period
func PricingTotalChildIncrementalNaturalDeaths(pp *models.PricingPoint, p models.PricingPoint, config models.PricingConfig, pricingParameter models.PricingParameter, params models.ProductPricingParameters) {
	if pp.ProjectionMonth == 0 {
		pp.TotalChildIncrementalNaturalDeaths = 0
		pp.TotalChildIncrementalNaturalDeathsAdjusted = 0
	} else {
		if int(math.Ceil(pp.ValuationTimeYear))+pricingParameter.AverageAgeAtEntryPerChild <= params.ChildExitAge {
			pp.TotalChildIncrementalNaturalDeaths = p.ChildNumberPolicies * pp.MonthlyChildDependentMortality * (1 - pp.ChildAccidentalProportion)
			pp.TotalChildIncrementalNaturalDeathsAdjusted = p.ChildNumberPoliciesAdjusted * pp.MonthlyChildDependentMortalityAdjusted * (1 - pp.ChildAccidentalProportion)
		}
	}
}

// PricingChildNumberOfPaidUpNaturalDeaths computes child's cumulative number of natural deaths, in the paid up state, over the projection period
func PricingChildNumberOfPaidUpNaturalDeaths(pp *models.PricingPoint, p models.PricingPoint, params models.ProductPricingParameters) {
	if pp.ProjectionMonth == 0 {
		pp.ChildNumberOfPaidUpNaturalDeaths = 0
		pp.ChildNumberOfPaidUpNaturalDeathsAdjusted = 0
		pp.TotalChildIncrementalPaidupNaturalDeaths = 0
		pp.TotalChildIncrementalPaidupNaturalDeathsAdjusted = 0
	} else {
		pp.ChildNumberOfPaidUpNaturalDeaths = p.ChildNumberOfPaidUps*pp.ChildIndependentMortalityMonthly*(1-pp.ChildAccidentalProportion) + p.ChildNumberOfPaidUpNaturalDeaths                                 //+ pp.ChildNumberOfPremiumWaivers*pp.ChildIndependentMortalityMonthly*(1-pp.ChildAccidentalProportion)*params.PremiumWaiverBenefitPayoutAdjusted
		pp.ChildNumberOfPaidUpNaturalDeathsAdjusted = p.ChildNumberOfPaidUpsAdjusted*pp.ChildIndependentMortalityAdjustedMonthly*(1-pp.ChildAccidentalProportion) + p.ChildNumberOfPaidUpNaturalDeathsAdjusted //+ pp.ChildNumberOfPremiumWaiversAdjusted*pp.ChildIndependentMortalityAdjustedMonthly*(1-pp.ChildAccidentalProportion)*params.PremiumWaiverBenefitPayoutAdjusted
		pp.TotalChildIncrementalPaidupNaturalDeaths = pp.ChildNumberOfPaidUpNaturalDeaths - p.ChildNumberOfPaidUpNaturalDeaths
		pp.TotalChildIncrementalPaidupNaturalDeathsAdjusted = pp.ChildNumberOfPaidUpNaturalDeathsAdjusted - p.ChildNumberOfPaidUpNaturalDeathsAdjusted
	}
}

// PricingChildNumberOfPaidUpAccidentalDeaths computes child's cumulative number of accidental deaths, in the paid up state, over the projection period
func PricingChildNumberOfPaidUpAccidentalDeaths(pp *models.PricingPoint, p models.PricingPoint, params models.ProductPricingParameters) {
	// child total natural deaths and paid ups
	if pp.ProjectionMonth == 0 {
		pp.ChildNumberOfPaidUpAccidentalDeaths = 0
		pp.ChildNumberOfPaidUpAccidentalDeathsAdjusted = 0
		pp.TotalChildIncrementalPaidupAccidentalDeaths = 0
		pp.TotalChildIncrementalPaidupAccidentalDeathsAdjusted = 0
	} else {
		pp.ChildNumberOfPaidUpAccidentalDeaths = p.ChildNumberOfPaidUps*pp.ChildMortalityRate*(pp.ChildAccidentalProportion) + p.ChildNumberOfPaidUpAccidentalDeaths                                 // + pp.ChildNumberOfPremiumWaivers*pp.ChildMortalityRate*(pp.ChildAccidentalProportion)*params.PremiumWaiverSumAssuredFactor
		pp.ChildNumberOfPaidUpAccidentalDeathsAdjusted = p.ChildNumberOfPaidUpsAdjusted*pp.ChildMortalityRateAdjusted*(pp.ChildAccidentalProportion) + p.ChildNumberOfPaidUpAccidentalDeathsAdjusted //pp.ChildNumberOfPremiumWaiversAdjusted*pp.ChildMortalityRateAdjusted*(pp.ChildAccidentalProportion)*params.PremiumWaiverSumAssuredFactor
		pp.TotalChildIncrementalPaidupAccidentalDeaths = pp.ChildNumberOfPaidUpAccidentalDeaths - p.ChildNumberOfPaidUpAccidentalDeaths
		pp.TotalChildIncrementalPaidupAccidentalDeathsAdjusted = pp.ChildNumberOfPaidUpAccidentalDeathsAdjusted - p.ChildNumberOfPaidUpAccidentalDeathsAdjusted
	}
}

func PricingChildNumberOfPremiumWaiverNaturalDeaths(pp *models.PricingPoint, p models.PricingPoint, params models.ProductPricingParameters) {
	if pp.ProjectionMonth == 0 {
		pp.ChildNumberOfPremiumWaiverNaturalDeaths = 0
		pp.ChildNumberOfPremiumWaiverNaturalDeathsAdjusted = 0
		pp.TotalChildIncrementalPwNaturalDeaths = 0
		pp.TotalChildIncrementalPwNaturalDeathsAdjusted = 0
	} else {
		pp.ChildNumberOfPremiumWaiverNaturalDeaths = p.ChildNumberOfPremiumWaivers*pp.ChildIndependentMortalityMonthly*(1-pp.ChildAccidentalProportion) + p.ChildNumberOfPremiumWaiverNaturalDeaths
		pp.ChildNumberOfPremiumWaiverNaturalDeathsAdjusted = p.ChildNumberOfPremiumWaiversAdjusted*pp.ChildIndependentMortalityAdjustedMonthly*(1-pp.ChildAccidentalProportion) + p.ChildNumberOfPremiumWaiverNaturalDeathsAdjusted
		pp.TotalChildIncrementalPwNaturalDeaths = pp.ChildNumberOfPremiumWaiverNaturalDeaths - p.ChildNumberOfPremiumWaiverNaturalDeaths
		pp.TotalChildIncrementalPwNaturalDeathsAdjusted = pp.ChildNumberOfPremiumWaiverNaturalDeathsAdjusted - p.ChildNumberOfPremiumWaiverNaturalDeathsAdjusted
	}
}

func PricingChildNumberOfPremiumWaiverAccidentalDeaths(pp *models.PricingPoint, p models.PricingPoint, params models.ProductPricingParameters) {
	// child total natural deaths and premium waiver
	if pp.ProjectionMonth == 0 {
		pp.ChildNumberOfPremiumWaiverAccidentalDeaths = 0
		pp.ChildNumberOfPremiumWaiverAccidentalDeathsAdjusted = 0
		pp.TotalChildIncrementalPwAccidentalDeaths = 0
		pp.TotalChildIncrementalPwAccidentalDeathsAdjusted = 0
	} else {
		pp.ChildNumberOfPremiumWaiverAccidentalDeaths = p.ChildNumberOfPremiumWaivers*pp.ChildMortalityRate*(pp.ChildAccidentalProportion) + p.ChildNumberOfPremiumWaiverAccidentalDeaths
		pp.ChildNumberOfPremiumWaiverAccidentalDeathsAdjusted = p.ChildNumberOfPremiumWaiversAdjusted*pp.ChildMortalityRateAdjusted*(pp.ChildAccidentalProportion) + p.ChildNumberOfPremiumWaiverAccidentalDeathsAdjusted
		pp.TotalChildIncrementalPwAccidentalDeaths = pp.ChildNumberOfPremiumWaiverAccidentalDeaths - p.ChildNumberOfPremiumWaiverAccidentalDeaths
		pp.TotalChildIncrementalPwAccidentalDeathsAdjusted = pp.ChildNumberOfPremiumWaiverAccidentalDeathsAdjusted - p.ChildNumberOfPremiumWaiverAccidentalDeathsAdjusted
	}
}

// PricingTotalChildIncrementalAccidentalDeaths computes child's number of new accidental deaths at each projection period
func PricingTotalChildIncrementalAccidentalDeaths(pp *models.PricingPoint, p models.PricingPoint, pricingParameter models.PricingParameter, params models.ProductPricingParameters) {
	if pp.ProjectionMonth > 0 {
		if int(math.Ceil(pp.ValuationTimeYear))+pricingParameter.AverageAgeAtEntryPerChild <= params.ChildExitAge {
			pp.TotalChildIncrementalAccidentalDeaths = p.ChildNumberPolicies * pp.MonthlyChildDependentMortality * (pp.ChildAccidentalProportion)
			pp.TotalChildIncrementalAccidentalDeathsAdjusted = p.ChildNumberPoliciesAdjusted * pp.MonthlyChildDependentMortalityAdjusted * (pp.ChildAccidentalProportion)
		}
	}
}

func PricingCalculatedInstalment(pp *models.PricingPoint, p *models.PricingPoint, mp models.ProductPricingModelPoint, parameters models.ProductPricingParameters) {
	if pp.ValuationTimeMonth <= parameters.CalculatedTerm {
		if pp.ProjectionMonth == 0 {
			if mp.Interest <= 0 {
				pp.CalculatedInstalment = mp.Instalment
			} else {
				var interestmonthly = math.Pow(1.0+mp.Interest, 1.0/12.0) - 1
				var numerator = 1 - math.Pow(1+interestmonthly, float64(-1.0*mp.OutstandingTermMonths))
				if numerator > 0 {
					pp.CalculatedInstalment = mp.OutstandingLoan * interestmonthly / numerator
				}
			}
		} else {
			pp.CalculatedInstalment = p.CalculatedInstalment
		}
	} else {
		pp.CalculatedInstalment = 0
	}
}

// OutstandingSumAssured computes outstanding loan amount at each projection period using two methodologies
// 1st methodology calculates outstanding sum assured as the present value of instalments over the projection period
// 2nd methodology computes outstanding loan amount off the outstanding sum assured after interest less the instalment amount
// Credit Life
func PricingOutstandingSumAssured(pp *models.PricingPoint, p *models.PricingPoint, mp models.ProductPricingModelPoint, params models.ProductPricingParameters, features models.ProductFeatures) {
	if pp.ValuationTimeMonth <= params.CalculatedTerm {
		if features.OsProjPvMethod {
			pp.OutstandingSumAssured = utils.FloatPrecision(CalculatePV(mp.Interest, math.Max(float64(mp.OutstandingTermMonths-pp.ValuationTimeMonth+1), 0), pp.CalculatedInstalment), defaultPrecision)
		} else {
			if pp.ProjectionMonth == 0 {
				pp.OutstandingSumAssured = math.Max(mp.OutstandingLoan, 0)
			} else {
				pp.OutstandingSumAssured = utils.FloatPrecision(math.Max(p.OutstandingSumAssured*(math.Pow(1+mp.Interest, 1/12.0))-pp.CalculatedInstalment, 0), defaultPrecision)
			}
		}
	} else {
		pp.OutstandingSumAssured = 0
	}
}

// PricingSumAssured projects sum assured over the projection period allowing for the sum assured escalation rate
func PricingSumAssured(pp *models.PricingPoint, mp models.ProductPricingModelPoint,
	params models.ProductPricingParameters, pricingParams models.PricingParameter, features models.ProductFeatures, pricingConfig models.PricingConfig) {
	if pp.ValuationTimeMonth <= params.CalculatedTerm {
		if mp.MemberType == Child {
			if params.ChildExitAge > pp.AgeNextBirthday {
				pp.SumAssured = getPricingChildSumAssured(pp.ProductCode, mp.Plan, int(math.Ceil(pp.ValuationTimeYear))+pricingParams.AverageAgeAtEntryPerChild)
			} else {
				pp.SumAssured = 0
			}
		} else {
			if features.SaOutstandingLoan {
				pp.SumAssured = pp.OutstandingSumAssured
			} else if features.SaFixedBaseLumpSum {
				pp.SumAssured = utils.FloatPrecision(mp.SumAssured*pp.SumAssuredEscalation, AccountingPrecision)
			} else {
				pp.SumAssured = 0
			}
		}
	} else {
		pp.SumAssured = 0
	}
}

// PricingChildSumAssured reads and projects child's sum assured over the projection period.
// It reads child's sum assured from the child's sum assured table by child's current age at each projection period
func PricingChildSumAssured(pp *models.PricingPoint, mp models.ProductPricingModelPoint, params models.ProductPricingParameters, pricingConfig models.PricingConfig, pricingParams models.PricingParameter) {
	if pricingConfig.ChildIndicator {
		if int(math.Ceil(pp.ValuationTimeYear))+pricingParams.AverageAgeAtEntryPerChild <= params.ChildExitAge {
			pp.ChildSumAssured = getPricingChildSumAssured(pp.ProductCode, mp.Plan, int(math.Ceil(pp.ValuationTimeYear))+pricingParams.AverageAgeAtEntryPerChild)
		} else {
			pp.ChildSumAssured = 0
		}
	}
}

// PricingAdditionalSumAssured reads and projects additionallump sum over the projection period
// PricingAdditionalSumAssured reads the from child's additional lump sum table if member type is child otherwise it reads from the additional child's lump Sum table
func PricingAdditionalSumAssured(pp *models.PricingPoint, mp models.ProductPricingModelPoint, params models.ProductPricingParameters, pricingParams models.PricingParameter) {
	if mp.MemberType == Child {
		if params.ChildExitAge > pp.AgeNextBirthday {
			pp.AdditionalSumAssured = getPricingChildAdditionalSumAssured(pp.ProductCode, mp.Plan, int(math.Ceil(pp.ValuationTimeYear))+pricingParams.AverageAgeAtEntryPerChild)
		} else {
			pp.AdditionalSumAssured = getPricingAdditionalSumAssured(pp.ProductCode, mp.Plan)
		}
	} else {
		pp.AdditionalSumAssured = params.StandardAdditionalSumAssured
	}
}

// PricingPremium reads annual premium from the model point file and projects over the projection period
func PricingPremium(pp *models.PricingPoint, mp models.ProductPricingModelPoint, features models.ProductFeatures, parameters models.ProductPricingParameters, pricingConfig models.PricingConfig, childFinalCalculatedPremium float64) {
	if pp.ValuationTimeMonth <= parameters.CalculatedTerm && mp.PremiumFrequency != 0 && pp.ProjectionMonth > 0 {
		if math.Mod(float64((pp.ValuationTimeMonth+11)*mp.PremiumFrequency), 12.0) == 0 {
			if features.CreditLifeFlatPremium {
				pp.Premium = mp.AnnualPremium / float64(mp.PremiumFrequency)
			} else if features.CreditLifeDecreasingPremium {
				pp.Premium = pp.OutstandingSumAssured * mp.PremiumRate * (1.0 / 1000.0) * (pp.PremiumEscalation / float64(mp.PremiumFrequency))
			} else {
				pp.Premium = mp.AnnualPremium * pp.PremiumEscalation / float64(mp.PremiumFrequency)
				pp.PremiumAdjusted = pp.Premium
				if pricingConfig.ChildIndicator {
					pp.ChildPremium = childFinalCalculatedPremium * pp.PremiumEscalation / float64(mp.PremiumFrequency)
				}
			}

		} else {
			pp.Premium = 0
			pp.PremiumAdjusted = 0
		}
	} else if mp.PremiumFrequency == 0 && pp.ProjectionMonth == 1 {
		pp.Premium = utils.FloatPrecision(mp.AnnualPremium, defaultPrecision)
		pp.PremiumAdjusted = pp.Premium
	} else {
		pp.Premium = 0
		pp.PremiumAdjusted = pp.Premium
	}
}

// PricingChildFuneralService reads child's funeral funeral service lump sum from the child's funeral service table
func PricingChildAdditionalSumAssured(pp *models.PricingPoint, mp models.ProductPricingModelPoint, parameters models.ProductPricingParameters, pricingConfig models.PricingConfig, pricingParams models.PricingParameter) {
	if pricingConfig.ChildIndicator {
		if int(math.Ceil(pp.ValuationTimeYear))+pricingParams.AverageAgeAtEntryPerChild <= parameters.ChildExitAge {
			pp.ChildAdditionalSumAssured = parameters.StandardAdditionalSumAssured //getPricingChildAdditionalSumAssured(pp.ProductCode, mp.Plan, int(math.Ceil(pp.ValuationTimeYear))+pricingParams.AverageAgeAtEntryPerChild)
		}
	}
}

// PricingPremiumIncome computes expected premium income at each projection period
func PricingPremiumIncome(pp *models.PricingPoint, p models.PricingPoint, mp models.ProductPricingModelPoint, parameters models.ProductPricingParameters, pricingConfig models.PricingConfig, pricingParams models.PricingParameter) {
	if mp.MemberType != "MM" || pp.ValuationTimeMonth > parameters.CalculatedTerm || pp.ProjectionMonth == 0 {
		pp.PremiumIncome = 0
		pp.PremiumIncomeAdjusted = 0
	} else {
		if pricingConfig.ChildIndicator {
			if int(math.Ceil(pp.ValuationTimeYear))+pricingParams.AverageAgeAtEntryPerChild <= parameters.ChildExitAge {
				pp.PremiumIncome = pp.Premium * p.InitialPolicy
				pp.PremiumIncomeAdjusted = pp.Premium * p.InitialPolicyAdjusted
			} else {
				pp.PremiumIncome = (pp.Premium - pp.ChildPremium) * p.InitialPolicy
				pp.PremiumIncomeAdjusted = (pp.Premium - pp.ChildPremium) * p.InitialPolicyAdjusted
			}

		} else {
			pp.PremiumIncome = pp.Premium * p.InitialPolicy
			pp.PremiumIncomeAdjusted = pp.Premium * p.InitialPolicyAdjusted
		}
	}
}

// PricingPremiumsNotReceived computes expected premium income not receipted during a grace period before a policy lapses
func PricingPremiumsNotReceived(pp *models.PricingPoint, mp models.ProductPricingModelPoint,
	parameters models.ProductPricingParameters, features models.ProductFeatures) {
	if features.PremiumHoliday && parameters.PremiumHolidayCycle != 0 && pp.ValuationTimeMonth > parameters.PremiumHolidayWaitingPeriod && pp.ValuationTimeMonth <= parameters.CalculatedTerm {
		pp.PremiumNotReceived = utils.FloatPrecision(math.Min(math.Max(float64(int(pp.ValuationTimeMonth/parameters.PremiumHolidayCycle)-mp.PremiumHolidayUsed), 0), float64(parameters.MaximumPremiumHolidays))*pp.PremiumIncome, defaultPrecision)
		pp.PremiumNotReceivedAdjusted = utils.FloatPrecision(math.Min(math.Max(float64(int(pp.ValuationTimeMonth/parameters.PremiumHolidayCycle)-mp.PremiumHolidayUsed), 0), float64(parameters.MaximumPremiumHolidays))*pp.PremiumIncomeAdjusted, defaultPrecision)
	} else {
		pp.PremiumNotReceived = 0
		pp.PremiumNotReceivedAdjusted = 0
	}
}

// PricingCommissions computes expected commission pay out using two methodologies
// 1st methodology computes initial commission
// 2nd methodology computes renewal commission
func PricingCommissions(mp models.ProductPricingModelPoint, pp *models.PricingPoint, p models.PricingPoint, parameters models.ProductPricingParameters, pricingParams models.PricingParameter) {
	if pp.ProjectionMonth == 0 || pp.ValuationTimeMonth > parameters.CalculatedTerm {
		pp.Commission = 0
		pp.CommissionAdjusted = 0
	} else {
		if mp.CommissionType == Initial {
			if pp.ValuationTimeMonth == 1 {
				pp.Commission = pricingParams.InitialCommissionPercentage1*mp.AnnualPremium + pricingParams.InitialCommissionRand
				pp.CommissionAdjusted = pricingParams.InitialCommissionPercentage1*mp.AnnualPremium + pricingParams.InitialCommissionRand

			} else {
				if pp.ValuationTimeMonth == 13 {
					pp.Commission = pricingParams.InitialCommissionPercentage2 * mp.AnnualPremium * p.InitialPolicy
					pp.CommissionAdjusted = pricingParams.InitialCommissionPercentage2 * mp.AnnualPremium * p.InitialPolicyAdjusted
				} else {
					pp.Commission = 0
					pp.CommissionAdjusted = 0
				}
			}
		} else if mp.CommissionType == Renewal {
			pp.Commission = pricingParams.RenewalCommissionPercentage*pp.PremiumIncome + pricingParams.RenewalCommissionRand*pp.InitialPolicy
			pp.CommissionAdjusted = pricingParams.RenewalCommissionPercentage*pp.PremiumIncome + pricingParams.RenewalCommissionRand*pp.InitialPolicyAdjusted
			pp.RenewalCommissionAnnuityFeeder = pricingParams.RenewalCommissionPercentage*pp.PremiumIncome + pricingParams.RenewalCommissionRand*pp.InitialPolicy
		} else if mp.CommissionType == Hybrid {
			var initialtemp float64
			var initialtempadjusted float64
			var initial12temp float64
			var initial12tempadjusted float64
			if pp.ValuationTimeMonth == 1 && pp.ProjectionMonth != 0 {
				initialtemp = utils.FloatPrecision(pricingParams.InitialCommissionPercentage1*mp.AnnualPremium+pricingParams.InitialCommissionRand, defaultPrecision)
				initialtempadjusted = utils.FloatPrecision(pricingParams.InitialCommissionPercentage1*mp.AnnualPremium+pricingParams.InitialCommissionRand, defaultPrecision)

			} else if pp.ValuationTimeMonth == 13 && pp.ProjectionMonth != 0 {
				initial12temp = utils.FloatPrecision(pricingParams.InitialCommissionPercentage2*mp.AnnualPremium*p.InitialPolicy, defaultPrecision)
				initial12tempadjusted = utils.FloatPrecision(pricingParams.InitialCommissionPercentage2*mp.AnnualPremium*p.InitialPolicyAdjusted, defaultPrecision)
			}
			if pp.ValuationTimeMonth >= pricingParams.HybridRenewalCommStartM && pp.ValuationTimeMonth <= pricingParams.HybridRenewalCommEndM {
				pp.Commission = initialtemp + initial12temp + pricingParams.RenewalCommissionPercentage*pp.PremiumIncome + pricingParams.RenewalCommissionRand*pp.InitialPolicy
				pp.CommissionAdjusted = initialtempadjusted + initial12tempadjusted + pricingParams.RenewalCommissionPercentage*pp.PremiumIncome + pricingParams.RenewalCommissionRand*pp.InitialPolicyAdjusted
				pp.RenewalCommissionAnnuityFeeder = pricingParams.RenewalCommissionPercentage*pp.PremiumIncome + pricingParams.RenewalCommissionRand*pp.InitialPolicy
			} else {
				pp.Commission = initialtemp + initial12temp                         //+ utils.FloatPrecision(parameters.RenewalCommissionPercentage*pp.PremiumIncome+parameters.RenewalCommissionRand*pp.InitialPolicy, defaultPrecision)
				pp.CommissionAdjusted = initialtempadjusted + initial12tempadjusted // + utils.FloatPrecision(parameters.RenewalCommissionPercentage*pp.PremiumIncome+parameters.RenewalCommissionRand*pp.InitialPolicyAdjusted, defaultPrecision)
			}
		} else {
			pp.Commission = 0
			pp.CommissionAdjusted = 0
		}
	}
}

func PricingClawback(mp models.ProductPricingModelPoint, pp *models.PricingPoint, parameters models.ProductPricingParameters, pricingParams models.PricingParameter) {
	if pp.ValuationTimeMonth > 0 && pp.ValuationTimeMonth <= pricingParams.ClawbackPeriod && mp.MemberType == "MM" && mp.CommissionType == Initial {
		var cb = getPricingClawback(pp.ValuationTimeMonth, mp)
		pp.ClawBack = cb.Year1InitialCommission*pp.TotalIncrementalLapses*pricingParams.InitialCommissionPercentage1*mp.AnnualPremium + cb.Year2InitialCommission*pp.TotalIncrementalLapses*pricingParams.InitialCommissionPercentage1*mp.AnnualPremium
		pp.ClawBackAdjusted = cb.Year1InitialCommission*pp.TotalIncrementalLapsesAdjusted*pricingParams.InitialCommissionPercentage1*mp.AnnualPremium + cb.Year2InitialCommission*pp.TotalIncrementalLapsesAdjusted*pricingParams.InitialCommissionPercentage1*mp.AnnualPremium
	} else if pp.ValuationTimeMonth > 0 && pp.ValuationTimeMonth <= pricingParams.ClawbackPeriod && mp.MemberType == "MM" && mp.CommissionType == Hybrid {
		var cb = getPricingClawback(pp.ValuationTimeMonth, mp)
		pp.ClawBack = cb.Year1InitialCommission*pp.TotalIncrementalLapses*pricingParams.InitialCommissionPercentage1*mp.AnnualPremium + cb.Year2InitialCommission*pp.TotalIncrementalLapses*pricingParams.InitialCommissionPercentage1*mp.AnnualPremium
		pp.ClawBackAdjusted = cb.Year1InitialCommission*pp.TotalIncrementalLapsesAdjusted*pricingParams.InitialCommissionPercentage1*mp.AnnualPremium + cb.Year2InitialCommission*pp.TotalIncrementalLapsesAdjusted*pricingParams.InitialCommissionPercentage1*mp.AnnualPremium
	} else {
		pp.ClawBack = 0
		pp.ClawBackAdjusted = 0
	}
}

// PricingDeathOutgo computes expected natural death pay out at each projection period
func PricingDeathOutgo(pp *models.PricingPoint, features models.ProductFeatures,
	modelPoint models.ProductPricingModelPoint, parameters models.ProductPricingParameters, pricingConfig models.PricingConfig) {
	if pp.ValuationTimeMonth <= modelPoint.WaitingPeriod || !features.DeathBenefit || pp.ValuationTimeMonth > parameters.CalculatedTerm {
		pp.DeathOutgo = 0
		pp.DeathOutgoAdjusted = 0
		pp.ChildDeathOutgo = 0
		pp.ChildDeathOutgoAdjusted = 0
		pp.SpouseDeathOutgo = 0
		pp.SpouseDeathOutgoAdjusted = 0
	} else {
		pp.DeathOutgo = (pp.TotalIncrementalNaturalDeaths + pp.TotalIncrementalPaidupNaturalDeaths) * (pp.SumAssured + parameters.StandardAdditionalSumAssured)                         //+pp.TotalSpouseIncrementalNaturalDeaths+pp.SpouseNumberOfPaidUpNaturalDeaths)*(pp.SumAssured+parameters.StandardAdditionalSumAssured) + (pp.TotalChildIncrementalNaturalDeaths+pp.ChildNumberOfPaidUpNaturalDeaths)*(pp.ChildSumAssured+parameters.StandardAdditionalSumAssured)
		pp.DeathOutgoAdjusted = (pp.TotalIncrementalNaturalDeathsAdjusted + pp.TotalIncrementalPaidupNaturalDeathsAdjusted) * (pp.SumAssured + parameters.StandardAdditionalSumAssured) //+pp.TotalSpouseIncrementalNaturalDeathsAdjusted+pp.SpouseNumberOfPaidUpNaturalDeathsAdjusted)*(pp.SumAssured+parameters.StandardAdditionalSumAssured) + (pp.TotalChildIncrementalNaturalDeathsAdjusted+pp.ChildNumberOfPaidUpNaturalDeathsAdjusted)*(pp.ChildSumAssured+parameters.StandardAdditionalSumAssured)
		pp.ChildDeathOutgo = (pp.TotalChildIncrementalNaturalDeaths + pp.TotalChildIncrementalPaidupNaturalDeaths + pp.TotalChildIncrementalPwNaturalDeaths) * (pp.ChildSumAssured + parameters.StandardAdditionalSumAssured)
		pp.ChildDeathOutgoAdjusted = (pp.TotalChildIncrementalNaturalDeathsAdjusted + pp.TotalChildIncrementalPaidupNaturalDeathsAdjusted + pp.TotalChildIncrementalPwNaturalDeathsAdjusted) * (pp.ChildSumAssured + parameters.StandardAdditionalSumAssured)
		pp.SpouseDeathOutgo = (pp.TotalSpouseIncrementalNaturalDeaths + pp.TotalSpouseIncrementalPaidupNaturalDeaths + pp.TotalSpouseIncrementalPwNaturalDeaths) * (pp.SumAssured + parameters.StandardAdditionalSumAssured)
		pp.SpouseDeathOutgoAdjusted = (pp.TotalSpouseIncrementalNaturalDeathsAdjusted + pp.TotalSpouseIncrementalPaidupNaturalDeathsAdjusted + pp.TotalSpouseIncrementalPwNaturalDeathsAdjusted) * (pp.SumAssured + parameters.StandardAdditionalSumAssured)
	}
}

// PricingAccidentalDeathOutgo computes expected accidental death benefit payout at each projection period
func PricingAccidentalDeathOutgo(pp *models.PricingPoint,
	parameters models.ProductPricingParameters,
	multipliers models.ProductPricingAccidentalBenefitMultiplier,
	features models.ProductFeatures, pricingConfig models.PricingConfig) {

	if pp.ProjectionMonth == 0 || pp.ValuationTimeMonth > parameters.CalculatedTerm || !features.AccidentalDeathBenefit {
		pp.AccidentalDeathOutgo = 0
		pp.AccidentalDeathOutgoAdjusted = 0
		pp.ChildAccidentalDeathOutgo = 0
		pp.ChildAccidentalDeathOutgoAdjusted = 0
		pp.SpouseAccidentalDeathOutgo = 0
		pp.SpouseAccidentalDeathOutgoAdjusted = 0
	} else {
		pp.AccidentalDeathOutgo = (pp.TotalIncrementalAccidentalDeaths + pp.TotalIncrementalPaidupAccidentalDeaths) * (multipliers.MainMember*pp.SumAssured + parameters.StandardAdditionalSumAssured)                         //+ (pp.TotalSpouseIncrementalAccidentalDeaths+pp.SpouseNumberOfPaidUpAccidentalDeaths)*(multipliers.Spouse*pp.SumAssured+parameters.StandardAdditionalSumAssured) + (pp.TotalChildIncrementalAccidentalDeaths+pp.ChildNumberOfPaidUpAccidentalDeaths)*(multipliers.Child*pp.ChildSumAssured+parameters.StandardAdditionalSumAssured)
		pp.AccidentalDeathOutgoAdjusted = (pp.TotalIncrementalAccidentalDeathsAdjusted + pp.TotalIncrementalPaidupAccidentalDeathsAdjusted) * (multipliers.MainMember*pp.SumAssured + parameters.StandardAdditionalSumAssured) //+ (pp.TotalSpouseIncrementalAccidentalDeathsAdjusted+pp.SpouseNumberOfPaidUpAccidentalDeathsAdjusted)*(multipliers.Spouse*pp.SumAssured+parameters.StandardAdditionalSumAssured) + (pp.TotalChildIncrementalAccidentalDeathsAdjusted+pp.ChildNumberOfPaidUpAccidentalDeathsAdjusted)*(multipliers.Child*pp.ChildSumAssured+parameters.StandardAdditionalSumAssured)
		pp.ChildAccidentalDeathOutgo = (pp.TotalChildIncrementalAccidentalDeaths + pp.TotalChildIncrementalPaidupAccidentalDeaths + pp.TotalChildIncrementalPwAccidentalDeaths) * (multipliers.Child*pp.ChildSumAssured + parameters.StandardAdditionalSumAssured)
		pp.ChildAccidentalDeathOutgoAdjusted = (pp.TotalChildIncrementalAccidentalDeathsAdjusted + pp.TotalChildIncrementalPaidupAccidentalDeathsAdjusted + pp.TotalChildIncrementalPwAccidentalDeathsAdjusted) * (multipliers.Child*pp.ChildSumAssured + parameters.StandardAdditionalSumAssured)
		pp.SpouseAccidentalDeathOutgo = (pp.TotalSpouseIncrementalAccidentalDeaths + pp.TotalSpouseIncrementalPaidupAccidentalDeaths + pp.TotalSpouseIncrementalPwAccidentalDeaths) * (multipliers.Spouse*pp.SumAssured + parameters.StandardAdditionalSumAssured)
		pp.SpouseAccidentalDeathOutgoAdjusted = (pp.TotalSpouseIncrementalAccidentalDeathsAdjusted + pp.TotalSpouseIncrementalPaidupAccidentalDeathsAdjusted + pp.TotalSpouseIncrementalPwAccidentalDeathsAdjusted) * (multipliers.Spouse*pp.SumAssured + parameters.StandardAdditionalSumAssured)

	}
}

func PricingDisabilityOutgo(pp *models.PricingPoint, mp models.ProductPricingModelPoint, parameters models.ProductPricingParameters, features models.ProductFeatures) {
	if !features.PermanentDisabilityBenefit || pp.ValuationTimeMonth > parameters.CalculatedTerm || pp.ValuationTimeMonth <= mp.WaitingPeriod {
		pp.DisabilityOutgo = 0
		pp.DisabilityOutgoAdjusted = 0
	} else {
		pp.DisabilityOutgo = utils.FloatPrecision(pp.TotalIncrementalDisabilities*pp.SumAssured, defaultPrecision)
		pp.DisabilityOutgoAdjusted = utils.FloatPrecision(pp.TotalIncrementalDisabilitiesAdjusted*pp.SumAssured, defaultPrecision)
	}
}

func PricingRetrenchmentOutgo(pp *models.PricingPoint, mp models.ProductPricingModelPoint, parameters models.ProductPricingParameters, features models.ProductFeatures) {
	if pp.ValuationTimeMonth <= mp.WaitingPeriod || !features.RetrenchmentBenefit || pp.ValuationTimeMonth > parameters.CalculatedTerm {
		pp.RetrenchmentOutgo = 0
		pp.RetrenchmentOutgoAdjusted = 0
	} else {
		pp.RetrenchmentOutgo = utils.FloatPrecision(pp.TotalIncrementalRetrenchments*parameters.RetrenchmentQualificationFactor*parameters.NumberRetrenchmentInstalments*(pp.CalculatedInstalment), defaultPrecision)
		pp.RetrenchmentOutgoAdjusted = utils.FloatPrecision(pp.TotalIncrementalRetrenchmentsAdjusted*parameters.RetrenchmentQualificationFactor*parameters.NumberRetrenchmentInstalments*(pp.CalculatedInstalment), defaultPrecision)
	}
}

// PricingEducator computes expected education benefit pay out at each projection period
// rider benefit and is applicable if educator indicator is 1
// payable upon a child surviving to a set age following the death of the Main Member
func PricingEducator(pp *models.PricingPoint, p models.PricingPoint, pricingParams models.PricingParameter, parameters models.ProductPricingParameters, mp models.ProductPricingModelPoint) {

	if pp.ValuationTimeMonth <= mp.EducatorWaitingPeriod {
		pp.NumberOfEducatorsInWp = pp.NaturalDeathsInForce + pp.NaturalDeathsPaidUp + pp.NaturalDeathsPremiumWaiver + pp.NaturalDeathsTemporaryWaivers + pp.NumberOfDeathsAccident + pp.NumberOfAccidentDeathsPaidUp + pp.AccidentDeathsPremiumWaiver + pp.AccidentDeathsTemporaryPremiumWaiver
		pp.NumberOfEducatorsInWpAdjusted = pp.NaturalDeathsInForceAdjusted + pp.NaturalDeathsPaidUpAdjusted + pp.NaturalDeathsPremiumWaiverAdjusted + pp.NaturalDeathsTemporaryWaiversAdjusted + pp.NumberOfDeathsAccidentAdjusted + pp.NumberOfAccidentDeathsPaidUpAdjusted + pp.AccidentDeathsPremiumWaiverAdjusted + pp.AccidentDeathsTemporaryPremiumWaiverAdjusted
	} else {
		pp.NumberOfEducatorsInWp = p.NumberOfEducatorsInWp
		pp.NumberOfEducatorsInWpAdjusted = p.NumberOfEducatorsInWpAdjusted
	}
	if (pp.ValuationTimeYear + float64(pricingParams.AverageAgeAtEntryPerChild)) == float64(parameters.EducatorSumAssuredPaymentAge) {
		if pp.ValuationTimeMonth > mp.EducatorWaitingPeriod {
			pp.Educator = (pp.NaturalDeathsInForce + pp.NaturalDeathsPaidUp + pp.NaturalDeathsPremiumWaiver + pp.NaturalDeathsTemporaryWaivers + pp.NumberOfDeathsAccident + pp.NumberOfAccidentDeathsPaidUp + pp.AccidentDeathsPremiumWaiver + pp.AccidentDeathsTemporaryPremiumWaiver - pp.NumberOfEducatorsInWp) * GetPricingRiderValue("Educator", mp.ProductCode, mp.Plan) * pp.ChildNumberPolicies
			pp.EducatorAdjusted = (pp.NaturalDeathsInForceAdjusted + pp.NaturalDeathsPaidUpAdjusted + pp.NaturalDeathsPremiumWaiverAdjusted + pp.NaturalDeathsTemporaryWaiversAdjusted + pp.NumberOfDeathsAccidentAdjusted + pp.NumberOfAccidentDeathsPaidUpAdjusted + pp.AccidentDeathsPremiumWaiverAdjusted + pp.AccidentDeathsTemporaryPremiumWaiverAdjusted - pp.NumberOfEducatorsInWpAdjusted) * GetPricingRiderValue("Educator", mp.ProductCode, mp.Plan) * pp.ChildNumberPoliciesAdjusted
		}
	}
}

// PricingCashBackOnSurvival computes expected cash back pay out at each projection period
// It computes cashback on survival and cashback on death of the main member
func PricingCashBackOnSurvival(pp *models.PricingPoint, mp models.ProductPricingModelPoint, parameters models.ProductPricingParameters, features models.ProductFeatures, pricingConfig models.PricingConfig) {
	if pp.ValuationTimeMonth <= parameters.CalculatedTerm && pp.ValuationTimeMonth > 1 {
		if parameters.CashbackOnSurvivalPeriod > 0 {
			if pp.ValuationTimeMonth%int(parameters.CashbackOnSurvivalPeriod) == 0 {
				if pp.ValuationTimeMonth <= int(parameters.CashbackOnSurvivalTerm) {
					pp.CashBackOnSurvival = parameters.CashbackOnSurvivalRatio * pp.InitialPolicy * mp.AnnualPremium
					pp.CashBackOnSurvivalAdjusted = parameters.CashbackOnSurvivalRatio * pp.InitialPolicyAdjusted * mp.AnnualPremium
				}
			} else {
				pp.CashBackOnSurvival = 0
				pp.CashBackOnSurvivalAdjusted = 0
			}
		}
	} else {
		pp.CashBackOnSurvival = 0
		pp.CashBackOnSurvivalAdjusted = 0
	}
}

// PricingCashBackOnDeath computes expected cash back pay out at each projection period
// It computes cashback on survival and cashback on death of the main member
// cashback of premiums paid todate on the death of the main member
func PricingCashBackOnDeath(pp *models.PricingPoint, mp models.ProductPricingModelPoint, parameters models.ProductPricingParameters, pricingConfig models.PricingConfig, pricingParams models.PricingParameter, features models.ProductFeatures) {
	// Add to this cashback above, the second part of the eq.
	if pp.ValuationTimeMonth <= parameters.CalculatedTerm && pp.ValuationTimeMonth <= int(parameters.CashbackOnDeathTerm) {
		pp.CashBackOnDeath = (pp.TotalIncrementalNaturalDeaths + pp.TotalIncrementalAccidentalDeaths) * float64(pp.ValuationTimeMonth/12) * mp.AnnualPremium * parameters.CashbackOnDeathRatio
		pp.CashBackOnDeathAdjusted = (pp.TotalIncrementalNaturalDeathsAdjusted + pp.TotalIncrementalAccidentalDeathsAdjusted) * float64(pp.ValuationTimeMonth/12) * mp.AnnualPremium * parameters.CashbackOnDeathRatio
	} else {
		pp.CashBackOnDeath = 0
		pp.CashBackOnDeathAdjusted = 0
	}
}

// computes funeral rider benefit payout
func PricingFuneralRider(pp *models.PricingPoint, mp models.ProductPricingModelPoint, parameters models.ProductPricingParameters, pricingConfig models.PricingConfig, pricingParams models.PricingParameter, features models.ProductFeatures) {
	// Add to this cashback above, the second part of the eq.
	if pp.ValuationTimeMonth <= parameters.CalculatedTerm && pp.ValuationTimeMonth <= int(parameters.CashbackOnDeathTerm) {
		pp.RiderFuneral = (pp.TotalIncrementalNaturalDeaths + pp.TotalIncrementalAccidentalDeaths) * math.Min(math.Max(parameters.RiderFuneralMinimumBenefit, parameters.RiderFuneralSaProportion*pp.SumAssured), parameters.RiderFuneralMaximumBenefit)
		pp.RiderFuneralAdjusted = (pp.TotalIncrementalNaturalDeathsAdjusted + pp.TotalIncrementalAccidentalDeathsAdjusted) * math.Min(math.Max(parameters.RiderFuneralMinimumBenefit, parameters.RiderFuneralSaProportion*pp.SumAssured), parameters.RiderFuneralMaximumBenefit)
	} else {
		pp.RiderFuneral = 0
		pp.RiderFuneralAdjusted = 0
	}
}

// PricingRider computes expected rider benefit pay out over the projection period
func PricingRider(pp *models.PricingPoint, features models.ProductFeatures,
	modelPoint models.ProductPricingModelPoint, pricingConfig models.PricingConfig, parameters models.ProductPricingParameters) {
	if features.RiderBenefit && pp.ValuationTimeMonth <= parameters.CalculatedTerm {
		if pp.ProjectionMonth == 0 || pp.ValuationTimeMonth <= modelPoint.WaitingPeriod {
			pp.Rider = 0
			pp.RiderAdjusted = 0
		} else {
			if pricingConfig.SpouseIndicator {
				pp.Rider = (pp.TotalIncrementalNaturalDeaths + pp.TotalIncrementalAccidentalDeaths + pp.TotalSpouseIncrementalNaturalDeaths + pp.TotalSpouseIncrementalAccidentalDeaths) * (getPricingRiders(modelPoint, features, pricingConfig, parameters))
				pp.RiderAdjusted = (pp.TotalIncrementalNaturalDeathsAdjusted + pp.TotalIncrementalAccidentalDeathsAdjusted + pp.TotalSpouseIncrementalNaturalDeathsAdjusted + pp.TotalSpouseIncrementalAccidentalDeathsAdjusted) * (getPricingRiders(modelPoint, features, pricingConfig, parameters))
			} else {
				pp.Rider = (pp.TotalIncrementalNaturalDeaths + pp.TotalIncrementalAccidentalDeaths) * (getPricingRiders(modelPoint, features, pricingConfig, parameters))
				pp.RiderAdjusted = (pp.TotalIncrementalNaturalDeathsAdjusted + pp.TotalIncrementalAccidentalDeathsAdjusted) * (getPricingRiders(modelPoint, features, pricingConfig, parameters))
			}
		}
	} else {
		pp.Rider = 0
		pp.RiderAdjusted = 0

	}
}

// PricingExpenses computes expected expense pay out at each projection period
func PricingExpenses(pp *models.PricingPoint, p models.PricingPoint, mp models.ProductPricingModelPoint, parameters models.ProductPricingParameters, margins models.ProductPricingMargins, features models.ProductFeatures, shock models.ProductPricingShock, pricingShockBasis string) {
	var exprand_Adjusted float64
	var polcount float64
	var polcountAdjusted float64
	var exprand float64 = 0
	var exprandShocked float64
	var exprandAdjustedShocked float64
	var exprandClaims float64
	var exprandClaimsadjusted float64
	var exprandClaimsShocked float64
	var exprandClaimsAdjustedShocked float64
	var polcountclaims float64
	var polcountclaimsAdjusted float64
	if pp.ProjectionMonth == 0 {
		exprand = 0
	} else {
		if pp.ValuationTimeMonth == 1 {
			exprand = parameters.InitialExpensePercentage*pp.Premium + parameters.InitialExpenseRand +
				parameters.RenewalExpensePercentage*pp.Premium + parameters.RenewalExpenseRand*pp.InflationFactor/12
			exprand_Adjusted = parameters.InitialExpensePercentage*pp.Premium + parameters.InitialExpenseRand +
				parameters.RenewalExpensePercentage*pp.Premium + parameters.RenewalExpenseRand*(1+margins.ExpenseMargin)*pp.InflationFactorAdjusted/12

			if pricingShockBasis != "N/A" {
				exprandShocked = exprand*(1+shock.MultiplicativeExpense) + shock.AdditiveExpense/12.0
				exprandAdjustedShocked = exprand_Adjusted*(1+shock.MultiplicativeExpense) + shock.AdditiveExpense/12.0
			}
			if pricingShockBasis == "N/A" {
				exprandShocked = exprand
				exprandAdjustedShocked = exprand_Adjusted
			}

		} else {
			exprand = parameters.RenewalExpensePercentage*pp.Premium + parameters.RenewalExpenseRand*pp.InflationFactor/12
			exprand_Adjusted = parameters.RenewalExpensePercentage*pp.Premium + parameters.RenewalExpenseRand*(1+margins.ExpenseMargin)*pp.InflationFactorAdjusted/12
			if pricingShockBasis != "N/A" {
				exprandShocked = exprand*(1+shock.MultiplicativeExpense) + shock.AdditiveExpense/12.0
				exprandAdjustedShocked = exprand_Adjusted*(1+shock.MultiplicativeExpense) + shock.AdditiveExpense/12.0
			}
			if pricingShockBasis == "N/A" {
				exprandShocked = exprand
				exprandAdjustedShocked = exprand_Adjusted
			}

		}
		exprandClaims = parameters.ClaimsExpensePercentage*pp.SumAssured + parameters.ClaimsExpenseRand*pp.InflationFactor
		exprandClaimsadjusted = parameters.ClaimsExpensePercentage*pp.SumAssured + parameters.ClaimsExpenseRand*(1+margins.ExpenseMargin)*pp.InflationFactorAdjusted
		if pricingShockBasis != "N/A" {
			exprandClaimsShocked = exprandClaims*(1+shock.MultiplicativeExpense) + shock.AdditiveExpense/12.0
			exprandClaimsAdjustedShocked = exprandClaimsadjusted*(1+shock.MultiplicativeExpense) + shock.AdditiveExpense/12.0
		}
		if pricingShockBasis == "N/A" {
			exprandClaimsShocked = exprandClaims
			exprandClaimsAdjustedShocked = exprandClaimsadjusted
		}
	}
	if mp.MemberType == "MM" {
		polcount = p.IncrementalPremiumWaivers + p.InitialPolicy + p.InitialTemporaryPremiumWaivers + p.PostPaidupPolicyCount + p.LastSurvivorPostPaidupPolicyCount
		polcountAdjusted = p.IncrementalPremiumWaiversAdjusted + p.InitialPolicyAdjusted + p.InitialTemporaryPremiumWaiversAdjusted + p.PostPaidupPolicyCountAdjusted + p.LastSurvivorPostPaidupPolicyCountAdjusted
		polcountclaims = pp.TotalIncrementalNaturalDeaths + pp.TotalIncrementalAccidentalDeaths //+ pp.TotalIncrementalDisabilities
		polcountclaimsAdjusted = pp.TotalIncrementalNaturalDeathsAdjusted + pp.TotalIncrementalAccidentalDeathsAdjusted
	} else {
		polcount = 0
		polcountAdjusted = 0
		polcountclaims = pp.TotalIncrementalNaturalDeaths + pp.TotalIncrementalAccidentalDeaths                         //+ pp.TotalIncrementalDisabilities
		polcountclaimsAdjusted = pp.TotalIncrementalNaturalDeathsAdjusted + pp.TotalIncrementalAccidentalDeathsAdjusted //+ pp.TotalIncrementalDisabilities
	}

	if pp.ValuationTimeMonth <= parameters.CalculatedTerm {
		pp.Expenses = exprandShocked*polcount + exprandClaimsShocked*polcountclaims                                         //utils.FloatPrecision(exprand*polcount+exprand_claims*polcountclaims, defaultPrecision)
		pp.ExpensesAdjusted = exprandAdjustedShocked*polcountAdjusted + exprandClaimsAdjustedShocked*polcountclaimsAdjusted //utils.FloatPrecision(exprand_Adjusted*polcountAdjusted+exprand_claimsadjusted*polcountclaimsAdjusted, defaultPrecision)
	}
}

// PricingNetCashFlow computes net cash flows from all the calculated cash flows arising on a policy
func PricingNetCashFlow(pp *models.PricingPoint, p models.PricingPoint, mp models.ProductPricingModelPoint, parameters models.ProductPricingParameters) {
	if pp.ProjectionMonth == 0 || pp.ValuationTimeMonth > parameters.CalculatedTerm {
		pp.NetCashFlow = 0
		pp.NetCashFlowAdjusted = 0
	} else {
		pp.NetCashFlow = pp.Commission + pp.DeathOutgo + pp.DisabilityOutgo + pp.RetrenchmentOutgo + pp.ChildDeathOutgo + pp.SpouseDeathOutgo + pp.AccidentalDeathOutgo + pp.ChildAccidentalDeathOutgo + pp.SpouseAccidentalDeathOutgo + pp.CashBackOnSurvival + pp.CashBackOnDeath + pp.RiderFuneral + pp.Rider + pp.Educator + pp.Expenses + pp.PremiumNotReceived - pp.PremiumIncome - pp.ClawBack
		pp.NetCashFlowAdjusted = pp.CommissionAdjusted + pp.DeathOutgoAdjusted + pp.DisabilityOutgoAdjusted + pp.RetrenchmentOutgoAdjusted + pp.ChildDeathOutgoAdjusted + pp.SpouseDeathOutgoAdjusted + pp.AccidentalDeathOutgoAdjusted + pp.ChildAccidentalDeathOutgoAdjusted + pp.SpouseAccidentalDeathOutgoAdjusted + pp.CashBackOnSurvivalAdjusted + pp.CashBackOnDeathAdjusted + pp.RiderFuneralAdjusted + pp.RiderAdjusted + pp.EducatorAdjusted + pp.ExpensesAdjusted + pp.PremiumNotReceivedAdjusted - pp.PremiumIncomeAdjusted - pp.ClawBackAdjusted
	}
}

func PricingkeyBuilder(dfs models.TransitionRateArguments, state string) string {
	var trans models.ProductTransition
	var err error

	associatedTableKey := strconv.Itoa(dfs.ProductId) + utils.Snakify(state)
	associatedTable, found := Cache.Get(associatedTableKey)
	if !found {
		err = DB.Where("product_id = ? and end_state =?", dfs.ProductId, state).First(&trans).Error
		if err != nil {
			log.Error().Msg(err.Error())
		} else {
			associatedTable = trans.AssociatedTable
			Cache.Set(associatedTableKey, associatedTable, 1)
		}
	}

	var ratingFactor models.ProductRatingFactor

	ratingFactorKey := "rf_" + strconv.Itoa(dfs.ProductId) + "_" + associatedTable.(string)

	cached, found := Cache.Get(ratingFactorKey)

	if !found {
		err = DB.Preload("Fds").Where("product_id = ? and transition_table =?", dfs.ProductId, associatedTable.(string)).First(&ratingFactor).Error
		if err != nil {
			log.Error().Msg(err.Error())
		} else {
			Cache.Set(ratingFactorKey, ratingFactor, 1)
		}
	} else {
		ratingFactor = cached.(models.ProductRatingFactor)
	}

	var keyString strings.Builder
	//keyString.WriteString(strconv.Itoa(dfs.Year) + "_")

	if utils.FactorsContains(&ratingFactor.Fds, "ANB") {
		keyString.WriteString(strconv.Itoa(dfs.Age) + "_")
	}

	if utils.FactorsContains(&ratingFactor.Fds, "GENDER") {
		keyString.WriteString(dfs.Gender[:1] + "_")
	}

	if utils.FactorsContains(&ratingFactor.Fds, "SMOKER_STATUS") {
		keyString.WriteString(dfs.SmokerStatus + "_")
	}

	if utils.FactorsContains(&ratingFactor.Fds, "INCOME") {
		keyString.WriteString(strconv.Itoa(dfs.Income))
		keyString.WriteString("_")
	}

	if utils.FactorsContains(&ratingFactor.Fds, "SEC") {
		keyString.WriteString(strconv.Itoa(dfs.SocioEconomicClass))
		keyString.WriteString("_")
	}

	if utils.FactorsContains(&ratingFactor.Fds, "OCC_CLASS") {
		keyString.WriteString(dfs.OccupationalClass)
		keyString.WriteString("_")
	}

	if utils.FactorsContains(&ratingFactor.Fds, "SELECT_PERIOD") {
		keyString.WriteString(strconv.Itoa(dfs.SelectPeriod))
		keyString.WriteString("_")
	}

	if utils.FactorsContains(&ratingFactor.Fds, "EDUCATION_LEVEL") {
		keyString.WriteString(strconv.Itoa(dfs.EducationLevel))
		keyString.WriteString("_")
	}

	if utils.FactorsContains(&ratingFactor.Fds, "DURATION_IF_M") {
		keyString.WriteString(strconv.Itoa(dfs.DurationIfM))
		keyString.WriteString("_")
	}

	if utils.FactorsContains(&ratingFactor.Fds, "PROJECTION_MONTH") {
		keyString.WriteString(strconv.Itoa(int(math.Min(float64(dfs.ProjectionMonth), 60.0))))
		keyString.WriteString("_")
	}

	if utils.FactorsContains(&ratingFactor.Fds, "DISTRIBUTION_CHANNEL") {
		keyString.WriteString(dfs.DistributionChannel)
		keyString.WriteString("_")
	}

	key := keyString.String()
	key = strings.Trim(key, "_")
	//cacheKey := dfs.ProductCode + "_" + trans.AssociatedTable + "_" + key

	return key
}
