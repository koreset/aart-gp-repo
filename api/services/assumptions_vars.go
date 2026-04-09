package services

import (
	"api/models"
	"api/utils"
	"math"
)

// AccidentProportion is a proportion of the mortality assumption that arises from non-natural causes
// Is used to derive accidental mortality rates from mortality rate table
func AccidentProportion(projection *models.Projection, modelPoint models.ProductModelPoint, run models.RunParameters, states []models.ProductTransitionState, parameters models.ProductParameters) {
	if utils.StatesContains(&states, AccidentalDeath) && projection.ValuationTimeMonth <= parameters.CalculatedTerm {
		args := models.TransitionRateArguments{
			ProductId:           projection.ProductID,
			ProductCode:         modelPoint.ProductCode,
			Year:                run.MortalityYear,
			Age:                 projection.AgeNextBirthday,
			Gender:              modelPoint.Gender,
			SmokerStatus:        modelPoint.SmokerStatus,
			Income:              modelPoint.Income,
			SocioEconomicClass:  modelPoint.SocioEconomicClass,
			OccupationalClass:   modelPoint.OccupationalClass,
			SelectPeriod:        modelPoint.SelectPeriod,
			EducationLevel:      modelPoint.EducationLevel,
			DurationIfM:         modelPoint.DurationInForceMonths + projection.ProjectionMonth,
			ProjectionMonth:     projection.ProjectionMonth,
			DistributionChannel: modelPoint.DistributionChannel,
		}
		projection.AccidentProportion = GetMortalityRateAccidentProportion(args)
	} else {
		projection.AccidentProportion = 0
	}
}

// InflationFactor projects inflation factor at each projection period
// Is used to inflate expenses
func InflationFactor(i, valYear, valMonth int, projection *models.Projection, p *models.Projection, productMargins models.ProductMargins, run models.RunParameters, mp models.ProductModelPoint, parameters models.ProductParameters, features models.ProductFeatures, shock models.ProductShock) {
	if i == 0 {
		projection.InflationFactor = 1
		projection.InflationFactorAdjusted = 1
	} else if projection.ValuationTimeMonth <= parameters.CalculatedTerm {
		var inflationFactor float64 = 0
		if run.YieldCurveBasis == Current {
			inflationFactor = GetInflationFactor(projection.ProjectionMonth, run.YieldcurveYear, run.YieldcurveMonth, parameters.YieldCurveCode)
		} else {
			inflationFactor = GetLockedinInflationFactor(projection.ProjectionMonth, mp.LockedInYear, mp.LockedInMonth)
		}
		tempInflationFactor := inflationFactor
		if run.ShockSettings.Inflation {
			inflationFactor = tempInflationFactor + math.Max(tempInflationFactor*shock.MultiplicativeInflation, shock.AdditiveInflation)
		}

		projection.InflationFactor = utils.FloatPrecision(p.InflationFactor*math.Pow(1+inflationFactor, 1/12.0), defaultPrecision)
		projection.InflationFactorAdjusted = utils.FloatPrecision(p.InflationFactorAdjusted*math.Pow(1+inflationFactor*(1+productMargins.InflationMargin), 1/12.0), defaultPrecision)
		if features.AnnuityIncome {
			projection.AnnuityEscalationRate = math.Min(inflationFactor, mp.AnnuityEscalation)
		}
	} else {
		projection.InflationFactor = 0
		projection.InflationFactorAdjusted = 0
		projection.AnnuityEscalationRate = 0
	}
}

func IbnrInflationFactor(i int, rd int, projMonths int, run models.IBNRRunSetting, interval int, yieldCurveCode string) float64 {
	var inflationFactor float64
	projectionMonth := int(math.Max(float64(i-(projMonths-rd)), 0.0))
	if projectionMonth == 0 {
		inflationFactor = 0
	} else {
		inflationFactor = GetIbnrInflationFactor(projectionMonth, run.YieldCurveYear, run.YieldCurveMonth, yieldCurveCode)
	}
	factor := math.Pow(1+inflationFactor, float64(interval)/12.0)

	return factor
}

func IbnrDiscountFactor(i int, rd int, projMonths int, run models.IBNRRunSetting, interval int, yieldCurveCode string, basis string) float64 {
	var discountRate float64
	projectionMonth := int(math.Max(float64(i-(projMonths-rd)), 0.0))
	if projectionMonth == 0 {
		discountRate = 0
	} else {
		discountRate, _ = GetIbnrForwardRateWithError(projectionMonth, run.YieldCurveYear, run.YieldCurveMonth, yieldCurveCode)
	}
	rate := math.Pow(1+discountRate, -float64(interval)/12.0)

	return rate
}

// PremiumEscalation projects premium escalation factor and sum assured escalation factot
// Reads premium escalation rate and sum assured escalation rate from a model point file
func PremiumEscalation(i int, projection *models.Projection, p *models.Projection, modelPoint models.ProductModelPoint, parameters models.ProductParameters) {
	if projection.ProjectionMonth == 0 {
		projection.PremiumEscalation = 1.0
		projection.SumAssuredEscalation = 1.0
	} else if projection.ValuationTimeMonth <= parameters.CalculatedTerm {
		if projection.CalendarMonth == modelPoint.EscalationMonth && projection.ValuationTimeMonth > 12 {
			projection.PremiumEscalation = p.PremiumEscalation * (1.0 + modelPoint.PremiumEscalation)
			projection.SumAssuredEscalation = p.SumAssuredEscalation * (1.0 + modelPoint.SumAssuredEscalation)
		} else {
			projection.PremiumEscalation = p.PremiumEscalation
			projection.SumAssuredEscalation = p.SumAssuredEscalation
		}
	} else {
		projection.PremiumEscalation = 0
		projection.SumAssuredEscalation = 0
	}
}

func AnnuityEscalation(i int, projection *models.Projection, p *models.Projection, modelPoint models.ProductModelPoint, parameters models.ProductParameters, features models.ProductFeatures) {
	if i == 0 {
		projection.AnnuityEscalation = 1
	} else if projection.ValuationTimeMonth <= parameters.CalculatedTerm {
		if features.AnnuityCalendarMonthEscalations {
			if projection.CalendarMonth == modelPoint.AnnuityEscalationMonth {
				projection.AnnuityEscalation = p.AnnuityEscalation * (1 + projection.AnnuityEscalationRate)
			} else {
				projection.AnnuityEscalation = p.AnnuityEscalation
			}

		} else {
			if (projection.ValuationTimeMonth-1)%12 == 0 && projection.ValuationTimeMonth > 1 {
				projection.AnnuityEscalation = p.AnnuityEscalation * (1 + modelPoint.AnnuityEscalation)
			} else {
				projection.AnnuityEscalation = p.AnnuityEscalation
			}
		}
	} else {
		projection.AnnuityEscalation = 0
	}
}

// LapseMargin projects lapse margin; sign and value, over the projection period
// Reads the lapse margin by duration and product code
func LapseMargin(projection *models.Projection, modelPoint models.ProductModelPoint, run models.RunParameters, parameters models.ProductParameters) {
	if projection.ValuationTimeMonth <= parameters.CalculatedTerm {
		projection.LapseMargin = GetLapseMargin(projection.ValuationTimeMonth, run.LapseYear, modelPoint.ProductCode, run.RunBasis)
	} else {
		projection.LapseMargin = 0
	}

}

// PremiumWaiverOnFactor is a boolean(0,1) factor that indicates when a premium waiver is on
// Is dependent on policy term, premium waiver waiting period and whether or not the premium waiver benefit has been chosen
func PremiumWaiverOnFactor(projection *models.Projection, modelPoint models.ProductModelPoint, parameters models.ProductParameters) {
	if (projection.ValuationTimeMonth) > parameters.CalculatedTerm || !modelPoint.ContinuityOrPremiumWaiverOption || (modelPoint.DurationInForceMonths+projection.ProjectionMonth) <= parameters.PremiumWaiverWaitingPeriod || modelPoint.MemberType == "MM" {
		projection.PremiumWaiverOnFactor = 0
	} else {
		projection.PremiumWaiverOnFactor = 1
	}
}

// PaidUpOnFactor is a binary factor that indicates when a paid up benefit is on
// Is dependent on the policy term,whether or not the paid up benefit has been chosen, and waiting period
func PaidUpOnFactor(projection *models.Projection, modelPoint models.ProductModelPoint, parameters models.ProductParameters) {
	// PaidUpOnFactor
	if projection.ValuationTimeMonth <= parameters.CalculatedTerm && projection.MainMemberAgeNextBirthday >= parameters.PaidupEffectiveAge && modelPoint.PaidupOption && projection.ValuationTimeMonth > parameters.PaidUpOnSurvivalWaitingPeriod {
		projection.PaidUpOnFactor = 1
	} else {
		projection.PaidUpOnFactor = 0
	}
}

// MainMemberMortalityRate reads mortality rate by main member's current age over the projection period and main member's gender
// Applicable for a funeral product
func MainMemberMortalityRate(projection *models.Projection, modelPoint models.ProductModelPoint, productMargins models.ProductMargins, run models.RunParameters, states []models.ProductTransitionState, shock models.ProductShock, parameters models.ProductParameters, features models.ProductFeatures, mainMemberspecialMortalityMargin, mainMembermortalityTableProp float64) {
	if utils.StatesContains(&states, Death) && projection.ValuationTimeMonth <= parameters.CalculatedTerm {
		if projection.MainMemberAgeNextBirthday >= 120 {
			projection.MainMemberMortalityRate = 1
			projection.MainMemberMortalityRateAdjusted = 1
		} else if projection.MainMemberAgeNextBirthday < 120 {
			mortalityArgs := models.TransitionRateArguments{
				ProductId:           projection.ProductID,
				ProductCode:         modelPoint.ProductCode,
				Age:                 projection.MainMemberAgeNextBirthday,
				Gender:              modelPoint.MainMemberGender,
				SmokerStatus:        modelPoint.SmokerStatus,
				Income:              modelPoint.Income,
				SocioEconomicClass:  modelPoint.SocioEconomicClass,
				SelectPeriod:        modelPoint.SelectPeriod,
				EducationLevel:      modelPoint.EducationLevel,
				DurationIfM:         modelPoint.DurationInForceMonths + projection.ProjectionMonth,
				ProjectionMonth:     projection.ProjectionMonth,
				DistributionChannel: modelPoint.DistributionChannel,
				Year:                run.MortalityYear,
			}
			var resp1 float64 = 0
			resp1 = math.Min(GetMortalityRate(mortalityArgs)*mainMembermortalityTableProp, 1)
			if run.ShockSettings.Mortality {
				projection.MainMemberMortalityRate = utils.FloatPrecision(math.Max(0, math.Min(1, resp1*(1+shock.MultiplicativeMortality)+shock.AdditiveMortality)), defaultPrecision)
			} else if run.ShockSettings.MortalityCatastrophe {
				monthlyResp := 1.0 - math.Pow(1.0-resp1, 1/12.0)
				shockedMonthly := monthlyResp + shock.MortalityCatastropheMultiplier*math.Min(math.Max(shock.MultiplicativeMortality*monthlyResp*1000.0+shock.AdditiveMortality, shock.MortalityCatastropheFloor), shock.MortalityCatastropheCeiling)/1000.0
				shockedAnnual := 1.0 - math.Pow(1.0-shockedMonthly, 12)
				projection.MainMemberMortalityRate = utils.FloatPrecision(shockedAnnual, defaultPrecision)
			} else {
				projection.MainMemberMortalityRate = utils.FloatPrecision(resp1, defaultPrecision)
			}

			var resp2 float64 = 0
			if features.SpecialDecrementMargin {
				resp2 = projection.MainMemberMortalityRate * (1 + mainMemberspecialMortalityMargin)
			} else {
				resp2 = projection.MainMemberMortalityRate * (1 + productMargins.MortalityMargin)
			}

			projection.MainMemberMortalityRateAdjusted = utils.FloatPrecision(math.Min(resp2, 1), defaultPrecision)
		}
	} else {
		projection.MainMemberMortalityRate = 0
		projection.MainMemberMortalityRateAdjusted = 0
	}
}

// BaseLapses reads the base lapse rate by the applicable rating factors eg. reads by duration aka monthInView
// It computes best estimate rates and accrual lapse rates
func BaseLapses(projection *models.Projection, p models.Projection, modelPoint models.ProductModelPoint, run models.RunParameters, states []models.ProductTransitionState, shock models.ProductShock, parameters models.ProductParameters, lapseTableProp float64) {
	if utils.StatesContains(&states, Lapse) {
		if (modelPoint.TemporaryPremiumWaiverIndicator && modelPoint.TemporaryPremiumWaiverMonthExit <= projection.ValuationTimeMonth) || modelPoint.PremiumWaiverIndicator || modelPoint.PaidupIndicator || projection.ValuationTimeMonth > parameters.CalculatedTerm || projection.ProjectionMonth == 0 {
			projection.BaseLapse = 0
			projection.BaseLapseAdjusted = 0
		} else {
			args := models.TransitionRateArguments{
				ProductId:           projection.ProductID,
				ProductCode:         modelPoint.ProductCode,
				Age:                 projection.AgeNextBirthday,
				Gender:              modelPoint.Gender,
				SmokerStatus:        modelPoint.SmokerStatus,
				Income:              modelPoint.Income,
				SocioEconomicClass:  modelPoint.SocioEconomicClass,
				SelectPeriod:        modelPoint.SelectPeriod,
				EducationLevel:      modelPoint.EducationLevel,
				DurationIfM:         modelPoint.DurationInForceMonths + projection.ProjectionMonth,
				ProjectionMonth:     projection.ProjectionMonth,
				DistributionChannel: modelPoint.DistributionChannel,
				Year:                run.LapseYear,
			}
			monthInView := projection.ValuationTimeMonth
			var resp float64 = 0.0
			resp = GetLapseRate(args) * lapseTableProp
			if run.ShockSettings.Lapse {
				if shock.MassLapse > 0 && projection.ProjectionMonth == 1 {
					projection.BaseLapse = 1 - math.Pow(1-shock.MassLapse, 12) //Mass lapse 40% instantaneous rate.
				} else {
					projection.BaseLapse = utils.FloatPrecision(math.Max(0, math.Min(1, resp*(1+shock.MultiplicativeLapse)+shock.AdditiveLapse)), defaultPrecision)
				}

			} else {
				projection.BaseLapse = utils.FloatPrecision(resp, defaultPrecision)
			}

			projection.BaseLapseAdjusted = utils.FloatPrecision(projection.BaseLapse*(1+(GetLapseMargin(monthInView, run.LapseYear, projection.ProductCode, run.RunBasis))), defaultPrecision)
		}
	} else {
		projection.BaseLapse = 0
		projection.BaseLapseAdjusted = 0
	}
}

// BaseMortalityRate reads base mortality rate by rating factors(age,gender)
func BaseMortalityRate(projection *models.Projection, modelPoint models.ProductModelPoint, run models.RunParameters, productMargins models.ProductMargins, states []models.ProductTransitionState, shock models.ProductShock, parameters models.ProductParameters, features models.ProductFeatures, specialMortalityMargin, mortalityTableProp float64) {
	if utils.StatesContains(&states, Death) && projection.ValuationTimeMonth <= parameters.CalculatedTerm {
		//if modelPoint.MemberType == "MM" {
		//	projection.BaseMortalityRate = projection.MainMemberMortalityRate
		//	projection.BaseMortalityRateAdjusted = projection.MainMemberMortalityRateAdjusted
		//} else {
		if projection.AgeNextBirthday >= 120 {
			projection.BaseMortalityRate = 1
			projection.BaseMortalityRateAdjusted = 1
		} else {
			mortalityArgs := models.TransitionRateArguments{
				ProductId:           projection.ProductID,
				ProductCode:         modelPoint.ProductCode,
				Age:                 projection.AgeNextBirthday,
				Gender:              modelPoint.Gender,
				SmokerStatus:        modelPoint.SmokerStatus,
				Income:              modelPoint.Income,
				SocioEconomicClass:  modelPoint.SocioEconomicClass,
				SelectPeriod:        modelPoint.SelectPeriod,
				EducationLevel:      modelPoint.EducationLevel,
				DurationIfM:         modelPoint.DurationInForceMonths + projection.ProjectionMonth,
				ProjectionMonth:     projection.ProjectionMonth,
				DistributionChannel: modelPoint.DistributionChannel,
				Year:                run.MortalityYear,
			}
			var resp float64 = 0
			resp = math.Min(GetMortalityRate(mortalityArgs)*mortalityTableProp, 1)
			if run.ShockSettings.Mortality {
				projection.BaseMortalityRate = utils.FloatPrecision(math.Max(0, math.Min(1, resp*(1+shock.MultiplicativeMortality)+shock.AdditiveMortality)), defaultPrecision)
			} else if run.ShockSettings.MortalityCatastrophe {
				monthlyResp := 1.0 - math.Pow(1.0-resp, 1/12.0)
				shockedMonthly := monthlyResp + shock.MortalityCatastropheMultiplier*math.Min(math.Max(shock.MultiplicativeMortality*monthlyResp*1000.0+shock.AdditiveMortality, shock.MortalityCatastropheFloor), shock.MortalityCatastropheCeiling)/1000.0
				shockedAnnual := 1.0 - math.Pow(1.0-shockedMonthly, 12)
				projection.BaseMortalityRate = utils.FloatPrecision(shockedAnnual, defaultPrecision)
			} else {
				projection.BaseMortalityRate = utils.FloatPrecision(resp, defaultPrecision)
			}
			if features.SpecialDecrementMargin {
				projection.BaseMortalityRateAdjusted = utils.FloatPrecision(math.Min(projection.BaseMortalityRate*(1+(specialMortalityMargin)), 1), defaultPrecision)
			} else {
				projection.BaseMortalityRateAdjusted = utils.FloatPrecision(math.Min(projection.BaseMortalityRate*(1+(productMargins.MortalityMargin)), 1), defaultPrecision)
			}
		}
		//	}
	} else {
		projection.BaseMortalityRate = 0
		projection.BaseMortalityRateAdjusted = 0
	}
}

// BaseNonLifeRiskRate reads NonLifeRatings Table by Duration inforce
func NonLifeRiskRate(projection *models.Projection, modelPoint models.ProductModelPoint, run models.RunParameters, productMargins models.ProductMargins, states []models.ProductTransitionState, shock models.ProductShock, parameters models.ProductParameters) {
	if projection.ValuationTimeMonth <= parameters.CalculatedTerm {
		if projection.ValuationTimeMonth >= 120 {
			projection.NonLifeMonthlyRiskRate = 1
			projection.NonLifeMonthlyRiskRateAdjusted = 1
		} else {
			nonLifeArgs := models.TransitionRateArguments{
				ProductId:   projection.ProductID,
				ProductCode: modelPoint.ProductCode,
				DurationIfM: modelPoint.DurationInForceMonths + projection.ProjectionMonth,
				Year:        run.ParameterYear,
			}

			projection.NonLifeMonthlyRiskRate = GetNonLifeRiskRate(nonLifeArgs) / 12.0
			projection.NonLifeMonthlyRiskRateAdjusted = projection.NonLifeMonthlyRiskRate
		}
	} else {
		projection.NonLifeMonthlyRiskRate = 0
		projection.NonLifeMonthlyRiskRateAdjusted = 0
	}
}

// BaseRetrenchmentRate calculates the BaseRetrenchmentRate on two bases; best estimate and accrual
func BaseRetrenchmentRate(projection *models.Projection, modelPoint models.ProductModelPoint, productMargins models.ProductMargins, run models.RunParameters, states []models.ProductTransitionState, shock models.ProductShock, parameters models.ProductParameters) {
	if projection.ValuationTimeMonth <= parameters.CalculatedTerm {
		if (modelPoint.TemporaryPremiumWaiverIndicator && modelPoint.TemporaryPremiumWaiverMonthExit > projection.ValuationTimeMonth) || modelPoint.PremiumWaiverIndicator || modelPoint.PaidupIndicator {
			projection.BaseRetrenchmentRate = 0
		} else {
			var retrenchmentRate float64
			dfs := models.TransitionRateArguments{
				ProductId:           projection.ProductID,
				ProductCode:         modelPoint.ProductCode,
				Year:                run.RetrenchmentYear,
				Age:                 projection.AgeNextBirthday,
				Gender:              modelPoint.Gender,
				SmokerStatus:        modelPoint.SmokerStatus,
				Income:              modelPoint.Income,
				SocioEconomicClass:  modelPoint.SocioEconomicClass,
				OccupationalClass:   modelPoint.OccupationalClass,
				SelectPeriod:        modelPoint.SelectPeriod,
				EducationLevel:      modelPoint.EducationLevel,
				DurationIfM:         modelPoint.DurationInForceMonths + projection.ProjectionMonth,
				ProjectionMonth:     projection.ProjectionMonth,
				DistributionChannel: modelPoint.DistributionChannel,
			}
			retrenchmentRate = GetRetrenchmentRate(dfs)
			if run.ShockSettings.Retrenchment {
				projection.BaseRetrenchmentRate = math.Max(0, math.Min(1, retrenchmentRate*(1+shock.MultiplicativeRetrenchment)+shock.AdditiveRetrenchment))
				projection.BaseRetrenchmentRateAdjusted = math.Min(1, projection.BaseRetrenchmentRate*(1+productMargins.RetrenchmentMargin))
			} else {
				projection.BaseRetrenchmentRate = retrenchmentRate
				projection.BaseRetrenchmentRateAdjusted = math.Min(1, projection.BaseRetrenchmentRate*(1+productMargins.RetrenchmentMargin))
			}

		}
	} else {
		projection.BaseRetrenchmentRate = 0
		projection.BaseRetrenchmentRateAdjusted = 0
	}
}

// BaseDisabilityIncrement reads base independent disability incidence rates by rating factors
// Dependent on rating factors. Rating factors include gender, age and sec
func BaseDisabilityIncidenceRates(projection *models.Projection, modelPoint models.ProductModelPoint, margins models.ProductMargins, run models.RunParameters, states []models.ProductTransitionState, shock models.ProductShock, parameters models.ProductParameters, disabilityTableProp float64) {
	if projection.ValuationTimeMonth <= parameters.CalculatedTerm {
		if (modelPoint.TemporaryPremiumWaiverIndicator && modelPoint.TemporaryPremiumWaiverMonthExit > projection.ValuationTimeMonth) || modelPoint.PremiumWaiverIndicator || modelPoint.PaidupIndicator || modelPoint.MemberType != "MM" {
			projection.BaseDisabilityIncidenceRate = 0
			projection.BaseDisabilityIncidenceRateAdjusted = 0
		} else {
			if projection.AgeNextBirthday >= 120 {
				projection.BaseDisabilityIncidenceRate = 1
				projection.BaseDisabilityIncidenceRateAdjusted = 1
			} else {
				dfs := models.TransitionRateArguments{
					ProductId:           projection.ProductID,
					ProductCode:         modelPoint.ProductCode,
					Year:                run.MorbidityYear,
					Age:                 projection.AgeNextBirthday,
					Gender:              modelPoint.Gender,
					SmokerStatus:        modelPoint.SmokerStatus,
					Income:              modelPoint.Income,
					SocioEconomicClass:  modelPoint.SocioEconomicClass,
					OccupationalClass:   modelPoint.OccupationalClass,
					SelectPeriod:        modelPoint.SelectPeriod,
					EducationLevel:      modelPoint.EducationLevel,
					DurationIfM:         modelPoint.DurationInForceMonths + projection.ProjectionMonth,
					ProjectionMonth:     projection.ProjectionMonth,
					DistributionChannel: modelPoint.DistributionChannel,
				}

				if run.ShockSettings.Disability {
					var resp float64 = 0
					resp = GetDisabilityIncidenceRate(dfs) * disabilityTableProp
					projection.BaseDisabilityIncidenceRate = math.Max(0, math.Min(1, resp*(1+shock.MultiplicativeDisability)+shock.AdditiveDisability))

				} else if run.ShockSettings.MorbidityCatastrophe {
					var resp float64 = 0
					resp = GetDisabilityIncidenceRate(dfs) * disabilityTableProp
					monthlyResp := 1.0 - math.Pow(1.0-resp, 1/12.0)
					shockedMonthly := monthlyResp + shock.CATScalar*monthlyResp*1000.0*shock.MorbidityCatastropheMultiplier/1000.0
					shockedAnnual := 1.0 - math.Pow(1.0-shockedMonthly, 12)
					projection.BaseDisabilityIncidenceRate = math.Max(0, math.Min(1, shockedAnnual))

				} else {
					projection.BaseDisabilityIncidenceRate = GetDisabilityIncidenceRate(dfs) * disabilityTableProp
				}
				projection.BaseDisabilityIncidenceRateAdjusted = math.Min(projection.BaseDisabilityIncidenceRate*(1+margins.MorbidityMargin), 1)
			}
		}
	} else {
		projection.BaseDisabilityIncidenceRate = 0
		projection.BaseDisabilityIncidenceRateAdjusted = 0
	}
}
