package services

import (
	"api/models"
	"api/utils"
	"fmt"
	"github.com/dgraph-io/ristretto"
	"math"
	"strconv"
	"strings"
	"time"
)

var ExpAnalysisCache *ristretto.Cache

func RunExpAnalysis(runjobs []models.ExpAnalysisRunSetting, user models.AppUser) error {
	ExpAnalysisCache.Clear()
	fmt.Println("Running dataprep")
	startDataPrep := time.Now()
	var runGroup models.ExpRunGroup
	runGroup.CreationDate = time.Now()
	runGroup.UserEmail = user.UserEmail
	runGroup.UserName = user.UserName
	runGroup.Name = runjobs[0].RunName + "-Group"
	runGroup.ProcessingStatus = "processing"
	runGroup.RunDuration = 0
	runGroup.TotalRecords = 0
	runGroup.ProcessedRecords = 0
	err := DB.Create(&runGroup).Error
	if err != nil {
		return err
	}

	for i, _ := range runjobs {
		runjobs[i].ExpRunGroupID = runGroup.ID
		runjobs[i].CreationDate = time.Now()
		runjobs[i].ProcessingStatus = "queued"
		runjobs[i].ProcessedRecords = 0
		runjobs[i].ExpRunGroupName = runGroup.Name
		var count int64
		DB.Table("exp_exposure_data").Where("year = ? and exp_configuration_id = ? and version = ?", runjobs[i].ExposureDataYear, runjobs[i].ExpConfigurationId, runjobs[i].ExposureDataVersion).Count(&count)
		runGroup.TotalRecords += int(count)
		runjobs[i].TotalRecords = int(count)
		runjobs[i].UserName = user.UserName
		runjobs[i].UserEmail = user.UserEmail
		err := DB.Save(&runjobs[i]).Error
		if err != nil {
			fmt.Println("Error saving ExpAnalysisRunSetting to DB: ", err)
			return err
		}
	}
	DB.Save(&runGroup)

	for _, run := range runjobs {

		startRunTime := time.Now()
		run.ProcessingStatus = "processing"

		DB.Save(&run)

		var expAnalysisData []models.ExpExposureData
		var expmps []models.ExposureModelPoint

		err := DB.Where("year = ? and exp_configuration_id = ? and version = ?", run.ExposureDataYear, run.ExpConfigurationId, run.ExposureDataVersion).Find(&expAnalysisData).Error
		if err != nil {
			fmt.Println(err)
			//	return err
		}

		//preparing data, calculating mp data exposures, durations and agebands at the model point level
		for _, mp := range expAnalysisData {
			var expmp models.ExposureModelPoint
			var exitDate, dateReachingNextAge time.Time
			var ageAtStart, ageAtExit int
			var commencementDate time.Time
			var dateofBirth time.Time
			var claimsDate time.Time
			var disabilityDate time.Time
			var periodStartDate time.Time
			var periodEndDate time.Time
			var lapseDate time.Time
			var dateAtEnd time.Time
			var dateAtEndAgeX time.Time

			expmp.ExpRunID = run.ID
			expmp.ExpRunGroupID = run.ExpRunGroupID
			expmp.PolicyNumber = mp.PolicyNumber
			expmp.MemberType = mp.MemberType
			expmp.SumAssured = mp.SumAssured
			expmp.Gender = mp.Gender
			expmp.ProductName = mp.ProductName
			expmp.CommencementDate = mp.CommencementDate
			expmp.DateOfBirth = mp.DateOfBirth
			expmp.LapseDate = mp.LapseDate
			expmp.ClaimsDate = mp.ClaimsDate
			expmp.DisabilityBenefitCeaseDate = mp.DisabilityBenefitCeaseDate
			expmp.ExitCode = mp.ExitCode
			expmp.RelevantDisabilityClaimIndicator = 0
			expmp.RelevantDeathClaimIndicator = 0
			expmp.RelevantLapseIndicator = 0

			if mp.LapseDate != "0" && mp.LapseDate != "" {
				lapseDate, err = utils.ParseDateString(expmp.LapseDate)
				if err != nil {
					fmt.Println("Error saving ExpModelPoint to DB: ", err)
					return err
				}
			} else {
				lapseDate, err = utils.ParseDateString("1900-01-01")
			}

			if mp.CommencementDate != "0" && mp.CommencementDate != "" {
				commencementDate, err = utils.ParseDateString(expmp.CommencementDate)
			} else {
				commencementDate, err = utils.ParseDateString("1900-01-01")
			}
			if mp.DateOfBirth != "0" && mp.DateOfBirth != "" {
				dateofBirth, err = utils.ParseDateString(expmp.DateOfBirth)
			} else {
				dateofBirth, err = utils.ParseDateString("1900-01-01")
			}
			if mp.ClaimsDate != "0" && mp.ClaimsDate != "" {
				claimsDate, err = utils.ParseDateString(expmp.ClaimsDate)
			} else {
				claimsDate, err = utils.ParseDateString("1900-01-01")
			}
			if mp.DisabilityBenefitCeaseDate != "0" && mp.DisabilityBenefitCeaseDate != "" {
				disabilityDate, err = utils.ParseDateString(expmp.DisabilityBenefitCeaseDate)
			} else {
				disabilityDate, err = utils.ParseDateString("1900-01-01")
			}
			periodStartDate, err = utils.ParseDateString(run.PeriodStartDate)
			periodEndDate, err = utils.ParseDateString(run.PeriodEndDate)

			//Max Date
			dateAtStart := MaxDate(commencementDate, dateofBirth, periodStartDate)
			expmp.DateAtStart = DatetoStringConv(dateAtStart)

			if disabilityDate.Year() != 1900 {
				exitDate = disabilityDate
				expmp.ExitDate = expmp.DisabilityBenefitCeaseDate
				if disabilityDate.Before(periodEndDate) && disabilityDate.After(dateAtStart) {
					expmp.RelevantDisabilityClaimIndicator = 1
				}
			}
			if claimsDate.Year() != 1900 {
				exitDate = claimsDate
				expmp.ExitDate = expmp.ClaimsDate
				if claimsDate.Before(periodEndDate) && claimsDate.After(dateAtStart) {
					expmp.RelevantDeathClaimIndicator = 1
				}
			}
			if claimsDate.Year() == 1900 && disabilityDate.Year() == 1900 {
				if lapseDate.Year() != 1900 {
					exitDate = lapseDate
					expmp.ExitDate = expmp.LapseDate
					if lapseDate.Before(periodEndDate) && lapseDate.After(dateAtStart) {
						expmp.RelevantLapseIndicator = 1
					}
				} else {
					exitDate = periodEndDate
				}
			} // a policy could have lapsed due to a claim. claim date takes precedence else {
			//	if exitDate.After(lapseDate) {
			//		exitDate = lapseDate
			//		expmp.ExitDate = expmp.LapseDate
			//		if lapseDate.Before(periodEndDate) && lapseDate.After(dateAtStart) {
			//			expmp.RelevantLapseIndicator = 1
			//		}
			//	}
			//}

			if dateAtStart.Month() >= dateofBirth.Month() {
				if dateAtStart.Month() == dateofBirth.Month() && dateAtStart.Day() <= dateofBirth.Day() {
					dateReachingNextAge = time.Date(dateAtStart.Year(), dateofBirth.Month(), dateofBirth.Day(), 23, 0, 0, 0, time.UTC)
				} else {
					dateReachingNextAge = time.Date(dateAtStart.Year()+1, dateofBirth.Month(), dateofBirth.Day(), 23, 0, 0, 0, time.UTC)
				}
			} else {
				dateReachingNextAge = time.Date(dateAtStart.Year(), dateofBirth.Month(), dateofBirth.Day(), 23, 0, 0, 0, time.UTC)
			}
			expmp.DateReachingAgeNext = DatetoStringConv(dateReachingNextAge)

			if expmp.ExitDate == "0" {
				dateAtEndAgeX = MinDate(dateReachingNextAge, periodEndDate, periodEndDate)
			}
			if expmp.ExitDate != "0" {
				dateAtEndAgeX = MinDate(dateReachingNextAge, periodEndDate, exitDate)
			}
			expmp.DateAtEndAgeX = DatetoStringConv(dateAtEndAgeX)

			if expmp.ExitDate == "0" {
				dateAtEnd = MaxDate(dateAtStart, periodEndDate, periodEndDate)
			}
			if expmp.ExitDate != "0" {
				dateAtEnd = MinDate(exitDate, periodEndDate, periodEndDate)
			}
			expmp.DateAtEnd = DatetoStringConv(dateAtEnd)

			//dateAgeXExit := MinDate(dateReachingNextAge, exitDate, dateAtStart)

			ageAtStart = int(math.Max(float64(int(dateAtStart.Sub(dateofBirth).Hours()/24/365))+1, 0))
			expmp.AgeAtStart = ageAtStart
			ageAtExit = int(math.Max((float64(int(dateAtEnd.Sub(dateofBirth).Hours()/24/365)))+1, 0))
			expmp.AgeAtExit = ageAtExit
			expmp.Version = mp.Version

			if dateAtEnd.After(dateAtStart) {
				expmp.MaxDuration = utils.FloatPrecision(math.Ceil((dateAtEnd.Sub(dateAtStart).Hours()/24/365)*12), AccountingPrecision)
				expmp.AgeXExposure = utils.FloatPrecision(math.Ceil((dateAtEndAgeX.Sub(dateAtStart).Hours()/24/365)*12), AccountingPrecision)
				expmp.DurationIfAtStart = utils.FloatPrecision(math.Ceil((dateAtStart.Sub(commencementDate).Hours()/24/365)*12), AccountingPrecision)
				expmp.DurationIfATEnd = utils.FloatPrecision(math.Ceil((dateAtEnd.Sub(commencementDate).Hours()/24/365)*12), AccountingPrecision)

				expmp.TotalDurationExposure = expmp.DurationIfATEnd - expmp.DurationIfAtStart

				expmp.DurationIfAtEndY = math.Min(5.0, math.Ceil(expmp.DurationIfATEnd/12.0))

				expmp.AgeNextExposure = expmp.MaxDuration - expmp.AgeXExposure

				expmp.AmountAgeXExposure = expmp.AgeXExposure * expmp.SumAssured

				expmp.AmountAgeNextExposure = expmp.AgeNextExposure * expmp.SumAssured

				expmp.DurationInYear1 = math.Min(math.Max(float64(12-expmp.DurationIfAtStart), 0), float64(expmp.TotalDurationExposure))
				expmp.DurationInYear2 = math.Min(math.Max(float64(24-expmp.DurationIfAtStart-expmp.DurationInYear1), 0), float64(expmp.TotalDurationExposure)-float64(expmp.DurationInYear1))

				tempSum3 := expmp.DurationInYear1 + expmp.DurationInYear2
				expmp.DurationInYear3 = math.Min(math.Max(float64(36-expmp.DurationIfAtStart-tempSum3), 0), float64(expmp.TotalDurationExposure)-float64(tempSum3))

				tempSum4 := expmp.DurationInYear1 + expmp.DurationInYear2 + expmp.DurationInYear3
				expmp.DurationInYear4 = math.Min(math.Max(float64(48-expmp.DurationIfAtStart-tempSum4), 0), float64(expmp.TotalDurationExposure)-float64(tempSum4))

				tempSum5 := expmp.DurationInYear1 + expmp.DurationInYear2 + expmp.DurationInYear3 + expmp.DurationInYear4
				expmp.DurationInYear5Plus = math.Min(math.Max(float64(expmp.TotalDurationExposure-tempSum5), 0), float64(expmp.TotalDurationExposure)-float64(tempSum5))

				MortalityRateAtAgeX := 1.0 - math.Pow(1.0-GetExpMortalityRate(expmp.Gender[:1], expmp.AgeAtStart, run.MortalityDataVersion), 1.0/12.0)
				MortalityRateAtNextAge := 1.0 - math.Pow(1.0-GetExpMortalityRate(expmp.Gender[:1], expmp.AgeAtExit, run.MortalityDataVersion), 1.0/12.0)

				if expmp.Gender[:1] == Male {
					expmp.ExpectedClaimCountMale = expmp.AgeXExposure * MortalityRateAtAgeX
					expmp.ExpectedClaimAmountMale = expmp.AmountAgeXExposure * MortalityRateAtAgeX
					expmp.ExpectedAgeNextClaimCountMale = expmp.AgeNextExposure * MortalityRateAtNextAge
					expmp.ExpectedAgeNextClaimAmountMale = expmp.AmountAgeNextExposure * MortalityRateAtNextAge
				}
				if expmp.Gender[:1] == Female {
					expmp.ExpectedClaimCountFemale = expmp.AgeXExposure * MortalityRateAtAgeX
					expmp.ExpectedClaimAmountFemale = expmp.AmountAgeXExposure * MortalityRateAtAgeX
					expmp.ExpectedAgeNextClaimCountFemale = expmp.AgeNextExposure * MortalityRateAtNextAge
					expmp.ExpectedAgeNextClaimAmountFemale = expmp.AmountAgeNextExposure * MortalityRateAtNextAge
				}
			}
			expmps = append(expmps, expmp)
			runGroup.ProcessedRecords += 0.5

		}

		var crudeRateResults []models.ExpCrudeResult
		var crudeLapseResults []models.ExpLapseCrudeResult
		for i := 1; i <= 120; i++ {
			var crudeRateResult models.ExpCrudeResult
			crudeRateResult.Age = i
			crudeRateResult.ExpRunID = run.ID
			crudeRateResult.ExpRunGroupID = run.ExpRunGroupID
			crudeRateResults = append(crudeRateResults, crudeRateResult)
		}

		for j := 1; j <= 5; j++ {
			var crudeLapseRateResult models.ExpLapseCrudeResult
			AnnualLapseRate := GetExpLapseRate(j, run.LapseDataVersion) //GetExpLapseRate(j, run.ExpConfigurationId)
			crudeLapseRateResult.DurationYear = j
			crudeLapseRateResult.ExpRunID = run.ID
			crudeLapseRateResult.ExpRunGroupID = run.ExpRunGroupID
			crudeLapseRateResult.ExpectedAnnualRate = AnnualLapseRate
			crudeLapseRateResult.ExpectedMonthlyRate = 1.0 - math.Pow(1.0-AnnualLapseRate, 1.0/12.0)
			crudeLapseRateResult.ExpectedUx = -math.Log(1 - crudeLapseRateResult.ExpectedMonthlyRate)

			crudeLapseResults = append(crudeLapseResults, crudeLapseRateResult)
		}

		var totalMortalityExpAnalysisResult models.TotalMortalityExpAnalysisResult
		var totalLapseExpAnalysisResult models.TotalLapseExpAnalysisResult

		totalMortalityExpAnalysisResult.ExpRunID = run.ID
		totalMortalityExpAnalysisResult.ExpRunGroupID = run.ExpRunGroupID
		totalLapseExpAnalysisResult.ExpRunID = run.ID
		totalLapseExpAnalysisResult.ExpRunGroupID = run.ExpRunGroupID

		crudeRateResults[0].DataPointCount = 0
		run.ProcessedRecords = 0
		for _, expmp := range expmps {
			mutex.Lock()
			if len(expmp.Gender) > 0 {
				crudeRateResults[0].DataPointCount += 1
				for i := 1; i <= 120; i++ {
					if expmp.AgeAtStart == i && expmp.Gender[:1] == Male {
						crudeRateResults[i-1].ExposureCountMale += expmp.AgeXExposure
						crudeRateResults[i-1].ExposureAmountMale += expmp.AmountAgeXExposure
						crudeRateResults[i-1].ExpectedClaimCountMale += expmp.ExpectedClaimCountMale
						crudeRateResults[i-1].ExpectedClaimAmountMale += expmp.ExpectedClaimAmountMale
					}

					if expmp.AgeAtStart == i && expmp.Gender[:1] == Female {
						crudeRateResults[i-1].ExposureCountFemale += expmp.AgeXExposure
						crudeRateResults[i-1].ExposureAmountFemale += expmp.AmountAgeXExposure
						crudeRateResults[i-1].ExpectedClaimCountFemale += expmp.ExpectedClaimCountFemale
						crudeRateResults[i-1].ExpectedClaimAmountFemale += expmp.ExpectedClaimAmountFemale
					}

					if expmp.AgeAtExit == i && expmp.Gender[:1] == Male {
						crudeRateResults[i-1].ActualClaimCountMale += 1 * float64(expmp.RelevantDeathClaimIndicator)
						crudeRateResults[i-1].ActualClaimAmountMale += expmp.SumAssured * float64(expmp.RelevantDeathClaimIndicator)
						crudeRateResults[i-1].ExposureCountMale += expmp.AgeNextExposure
						crudeRateResults[i-1].ExposureAmountMale += expmp.AmountAgeNextExposure
						crudeRateResults[i-1].ExpectedClaimCountMale += expmp.ExpectedAgeNextClaimCountMale
						crudeRateResults[i-1].ExpectedClaimAmountMale += expmp.ExpectedAgeNextClaimAmountMale
					}

					if expmp.AgeAtExit == i && expmp.Gender[:1] == Female {
						crudeRateResults[i-1].ActualClaimCountFemale += 1 * float64(expmp.RelevantDeathClaimIndicator)
						crudeRateResults[i-1].ActualClaimAmountFemale += expmp.SumAssured * float64(expmp.RelevantDeathClaimIndicator)
						crudeRateResults[i-1].ExposureCountFemale += expmp.AgeNextExposure
						crudeRateResults[i-1].ExposureAmountFemale += expmp.AmountAgeNextExposure
						crudeRateResults[i-1].ExpectedClaimCountFemale += expmp.ExpectedAgeNextClaimCountFemale
						crudeRateResults[i-1].ExpectedClaimAmountFemale += expmp.ExpectedAgeNextClaimAmountFemale
					}
				}

				//Totals
				if expmp.Gender[:1] == Male {
					totalMortalityExpAnalysisResult.TotalExposureCountMale += expmp.AgeXExposure + expmp.AgeNextExposure
					totalMortalityExpAnalysisResult.TotalExposureAmountMale += expmp.AmountAgeXExposure + expmp.AmountAgeNextExposure
					totalMortalityExpAnalysisResult.TotalExpectedClaimCountMale += expmp.ExpectedClaimCountMale + expmp.ExpectedAgeNextClaimCountMale
					totalMortalityExpAnalysisResult.TotalExpectedClaimAmountMale += expmp.ExpectedClaimAmountMale + expmp.ExpectedAgeNextClaimAmountMale
					totalMortalityExpAnalysisResult.TotalActualClaimCountMale += 1 * float64(expmp.RelevantDeathClaimIndicator)
					totalMortalityExpAnalysisResult.TotalActualClaimAmountMale += expmp.SumAssured * float64(expmp.RelevantDeathClaimIndicator)
				}

				if expmp.Gender[:1] == Female {
					totalMortalityExpAnalysisResult.TotalExposureCountFemale += expmp.AgeXExposure + expmp.AgeNextExposure
					totalMortalityExpAnalysisResult.TotalExposureAmountFemale += expmp.AmountAgeXExposure + expmp.AmountAgeNextExposure
					totalMortalityExpAnalysisResult.TotalExpectedClaimCountFemale += expmp.ExpectedClaimCountFemale + expmp.ExpectedAgeNextClaimCountFemale
					totalMortalityExpAnalysisResult.TotalExpectedClaimAmountFemale += expmp.ExpectedClaimAmountFemale + expmp.ExpectedAgeNextClaimAmountFemale
					totalMortalityExpAnalysisResult.TotalActualClaimCountFemale += 1 * float64(expmp.RelevantDeathClaimIndicator)
					totalMortalityExpAnalysisResult.TotalActualClaimAmountFemale += expmp.SumAssured * float64(expmp.RelevantDeathClaimIndicator)
				}

				for j := 1; j <= 5; j++ {
					if expmp.MemberType == "MM" {
						if j == 1 {
							crudeLapseResults[j-1].CentralExposure += expmp.DurationInYear1
							if expmp.DurationIfAtEndY == 1 {
								crudeLapseResults[j-1].ActualLapses += float64(expmp.RelevantLapseIndicator)
								totalLapseExpAnalysisResult.TotalActualYear1Lapses += float64(expmp.RelevantLapseIndicator)
							}
						}
						if j == 2 {
							crudeLapseResults[j-1].CentralExposure += expmp.DurationInYear2
							if expmp.DurationIfAtEndY == 2 {
								crudeLapseResults[j-1].ActualLapses += float64(expmp.RelevantLapseIndicator)
								totalLapseExpAnalysisResult.TotalActualYear2Lapses += float64(expmp.RelevantLapseIndicator)
							}
						}
						if j == 3 {
							crudeLapseResults[j-1].CentralExposure += expmp.DurationInYear3
							if expmp.DurationIfAtEndY == 3 {
								crudeLapseResults[j-1].ActualLapses += float64(expmp.RelevantLapseIndicator)
								totalLapseExpAnalysisResult.TotalActualYear3Lapses += float64(expmp.RelevantLapseIndicator)
							}
						}
						if j == 4 {
							crudeLapseResults[j-1].CentralExposure += expmp.DurationInYear4
							if expmp.DurationIfAtEndY == 4 {
								crudeLapseResults[j-1].ActualLapses += float64(expmp.RelevantLapseIndicator)
								totalLapseExpAnalysisResult.TotalActualYear4Lapses += float64(expmp.RelevantLapseIndicator)
							}
						}
						if j == 5 {
							crudeLapseResults[j-1].CentralExposure += expmp.DurationInYear5Plus
							if expmp.DurationIfAtEndY == 5 {
								crudeLapseResults[j-1].ActualLapses += float64(expmp.RelevantLapseIndicator)
								totalLapseExpAnalysisResult.TotalActualYear5PlusLapses += float64(expmp.RelevantLapseIndicator)
							}
						}
					}
				}

				//Totals
				totalLapseExpAnalysisResult.TotalYear1Exposure += expmp.DurationInYear1
				totalLapseExpAnalysisResult.TotalYear2Exposure += expmp.DurationInYear2
				totalLapseExpAnalysisResult.TotalYear3Exposure += expmp.DurationInYear3
				totalLapseExpAnalysisResult.TotalYear4Exposure += expmp.DurationInYear4
				totalLapseExpAnalysisResult.TotalYear5PlusExposure += expmp.DurationInYear5Plus

			}
			mutex.Unlock()
			runGroup.ProcessedRecords += 0.5
			DB.Save(&runGroup)
			run.ProcessedRecords += 1
		}

		if totalMortalityExpAnalysisResult.TotalExposureCountMale > 0 {
			//Expected
			totalMortalityExpAnalysisResult.TotalExpectedCrudeRatesLivesMale = 1 - math.Exp(-totalMortalityExpAnalysisResult.TotalExpectedClaimCountMale/totalMortalityExpAnalysisResult.TotalExposureCountMale)*1000
			totalMortalityExpAnalysisResult.TotalExpectedCrudeRatesAmountMale = 1 - math.Exp(-totalMortalityExpAnalysisResult.TotalExpectedClaimAmountMale/totalMortalityExpAnalysisResult.TotalExposureAmountMale)*1000

			//Actual
			totalMortalityExpAnalysisResult.TotalActualCrudeRatesLivesMale = 1 - math.Exp(-totalMortalityExpAnalysisResult.TotalActualClaimCountMale/totalMortalityExpAnalysisResult.TotalExposureCountMale)*1000
			totalMortalityExpAnalysisResult.TotalActualCrudeRatesAmountMale = 1 - math.Exp(-totalMortalityExpAnalysisResult.TotalActualClaimAmountMale/totalMortalityExpAnalysisResult.TotalExposureAmountMale)*1000

		}

		if totalMortalityExpAnalysisResult.TotalExposureCountFemale > 0 {
			//Expected
			totalMortalityExpAnalysisResult.TotalExpectedCrudeRatesLivesFemale = 1 - math.Exp(-totalMortalityExpAnalysisResult.TotalExpectedClaimCountFemale/totalMortalityExpAnalysisResult.TotalExposureCountFemale)*1000
			totalMortalityExpAnalysisResult.TotalExpectedCrudeRatesAmountFemale = 1 - math.Exp(-totalMortalityExpAnalysisResult.TotalExpectedClaimAmountFemale/totalMortalityExpAnalysisResult.TotalExposureAmountFemale)*1000

			//Actual
			totalMortalityExpAnalysisResult.TotalActualCrudeRatesLivesFemale = 1 - math.Exp(-totalMortalityExpAnalysisResult.TotalActualClaimCountFemale/totalMortalityExpAnalysisResult.TotalExposureCountFemale)*1000
			totalMortalityExpAnalysisResult.TotalActualCrudeRatesAmountFemale = 1 - math.Exp(-totalMortalityExpAnalysisResult.TotalActualClaimAmountFemale/totalMortalityExpAnalysisResult.TotalExposureAmountFemale)*1000

		}
		totalMortalityExpAnalysisResult.Period = run.PeriodStartDate + "-" + run.PeriodEndDate

		for i, _ := range crudeRateResults {
			if crudeRateResults[i].ExposureCountMale > 0 {
				crudeRateResults[i].CrudeRatesLivesMale = math.Min(1-math.Exp(-crudeRateResults[i].ActualClaimCountMale/crudeRateResults[i].ExposureCountMale), 1)
			}
			if crudeRateResults[i].ExposureCountFemale > 0 {
				crudeRateResults[i].CrudeRatesLivesFemale = math.Min(1-math.Exp(-crudeRateResults[i].ActualClaimCountFemale/crudeRateResults[i].ExposureCountFemale), 1)
			}
			if crudeRateResults[i].ExposureAmountMale > 0 {
				crudeRateResults[i].CrudeRatesAmountMale = math.Min(1-math.Exp(-crudeRateResults[i].ActualClaimAmountMale/crudeRateResults[i].ExposureAmountMale), 1)
			}
			if crudeRateResults[i].ExposureAmountFemale > 0 {
				crudeRateResults[i].CrudeRatesAmountFemale = math.Min(1-math.Exp(-crudeRateResults[i].ActualClaimAmountFemale/crudeRateResults[i].ExposureAmountFemale), 1)
			}
		}

		for j, _ := range crudeLapseResults {
			crudeLapseResults[j].ExpectedLapses = crudeLapseResults[j].ExpectedUx * crudeLapseResults[j].CentralExposure

			if crudeLapseResults[j].CentralExposure > 0 {
				crudeLapseResults[j].ActualUx = crudeLapseResults[j].ActualLapses / crudeLapseResults[j].CentralExposure
			}
			crudeLapseResults[j].ActualMonthlyRate = 1 - math.Exp(-crudeLapseResults[j].ActualUx)
			crudeLapseResults[j].ActualAnnualRate = 1 - math.Pow(1.0-crudeLapseResults[j].ActualMonthlyRate, 12.0)
			sigma := math.Pow(crudeLapseResults[j].ActualUx*(1-crudeLapseResults[j].ActualUx), 1.0/2.0)
			expectedActualDiff := crudeLapseResults[j].ExpectedUx - crudeLapseResults[j].ActualUx
			n := crudeLapseResults[j].CentralExposure
			if sigma > 0 {
				crudeLapseResults[j].TestStatisticZ = expectedActualDiff * math.Sqrt(n) / sigma
			}

			mutex.Lock()
			if j == 0 {
				totalLapseExpAnalysisResult.TotalExpectedYear1Lapses += crudeLapseResults[j].ExpectedLapses
				totalLapseExpAnalysisResult.ExpectedYear1Ux = crudeLapseResults[j].ExpectedUx
			}
			if j == 1 {
				totalLapseExpAnalysisResult.TotalExpectedYear2Lapses += crudeLapseResults[j].ExpectedLapses
				totalLapseExpAnalysisResult.ExpectedYear2Ux = crudeLapseResults[j].ExpectedUx
			}
			if j == 2 {
				totalLapseExpAnalysisResult.TotalExpectedYear3Lapses += crudeLapseResults[j].ExpectedLapses
				totalLapseExpAnalysisResult.ExpectedYear3Ux = crudeLapseResults[j].ExpectedUx
			}
			if j == 3 {
				totalLapseExpAnalysisResult.TotalExpectedYear4Lapses += crudeLapseResults[j].ExpectedLapses
				totalLapseExpAnalysisResult.ExpectedYear4Ux = crudeLapseResults[j].ExpectedUx
			}
			if j == 4 {
				totalLapseExpAnalysisResult.TotalExpectedYear5PlusLapses += crudeLapseResults[j].ExpectedLapses
				totalLapseExpAnalysisResult.ExpectedYear5PlusUx = crudeLapseResults[j].ExpectedUx
			}
			mutex.Unlock()
		}

		//total ux

		totalLapseExpAnalysisResult.Period = run.PeriodStartDate + "-" + run.PeriodEndDate
		if totalLapseExpAnalysisResult.TotalYear1Exposure > 0 {
			totalLapseExpAnalysisResult.ActualYear1Ux = totalLapseExpAnalysisResult.TotalActualYear1Lapses / totalLapseExpAnalysisResult.TotalYear1Exposure
			totalLapseExpAnalysisResult.ActualYear1MonthlyRate = 1 - math.Exp(-totalLapseExpAnalysisResult.ActualYear1Ux)
			totalLapseExpAnalysisResult.ActualYear1AnnualRate = 1 - math.Pow(1-totalLapseExpAnalysisResult.ActualYear1MonthlyRate, 12.0)

			totalLapseExpAnalysisResult.ExpectedYear1MonthlyRate = 1 - math.Exp(-totalLapseExpAnalysisResult.ExpectedYear1Ux)
			totalLapseExpAnalysisResult.ExpectedYear1AnnualRate = 1 - math.Pow(1-totalLapseExpAnalysisResult.ExpectedYear1MonthlyRate, 12.0)

		}
		if totalLapseExpAnalysisResult.TotalYear2Exposure > 0 {
			totalLapseExpAnalysisResult.ActualYear2Ux = totalLapseExpAnalysisResult.TotalActualYear2Lapses / totalLapseExpAnalysisResult.TotalYear2Exposure
			totalLapseExpAnalysisResult.ActualYear2MonthlyRate = 1 - math.Exp(-totalLapseExpAnalysisResult.ActualYear2Ux)
			totalLapseExpAnalysisResult.ActualYear2AnnualRate = 1 - math.Pow(1-totalLapseExpAnalysisResult.ActualYear2MonthlyRate, 12.0)

			totalLapseExpAnalysisResult.ExpectedYear2MonthlyRate = 1 - math.Exp(-totalLapseExpAnalysisResult.ExpectedYear2Ux)
			totalLapseExpAnalysisResult.ExpectedYear2AnnualRate = 1 - math.Pow(1-totalLapseExpAnalysisResult.ExpectedYear2MonthlyRate, 12.0)
		}
		if totalLapseExpAnalysisResult.TotalYear3Exposure > 0 {
			totalLapseExpAnalysisResult.ActualYear3Ux = totalLapseExpAnalysisResult.TotalActualYear3Lapses / totalLapseExpAnalysisResult.TotalYear3Exposure
			totalLapseExpAnalysisResult.ActualYear3MonthlyRate = 1 - math.Exp(-totalLapseExpAnalysisResult.ActualYear3Ux)
			totalLapseExpAnalysisResult.ActualYear3AnnualRate = 1 - math.Pow(1-totalLapseExpAnalysisResult.ActualYear3MonthlyRate, 12.0)

			totalLapseExpAnalysisResult.ExpectedYear3MonthlyRate = 1 - math.Exp(-totalLapseExpAnalysisResult.ExpectedYear3Ux)
			totalLapseExpAnalysisResult.ExpectedYear3AnnualRate = 1 - math.Pow(1-totalLapseExpAnalysisResult.ExpectedYear3MonthlyRate, 12.0)
		}
		if totalLapseExpAnalysisResult.TotalYear4Exposure > 0 {
			totalLapseExpAnalysisResult.ActualYear4Ux = totalLapseExpAnalysisResult.TotalActualYear4Lapses / totalLapseExpAnalysisResult.TotalYear4Exposure
			totalLapseExpAnalysisResult.ActualYear4MonthlyRate = 1 - math.Exp(-totalLapseExpAnalysisResult.ActualYear4Ux)
			totalLapseExpAnalysisResult.ActualYear4AnnualRate = 1 - math.Pow(1-totalLapseExpAnalysisResult.ActualYear4MonthlyRate, 12.0)

			totalLapseExpAnalysisResult.ExpectedYear4MonthlyRate = 1 - math.Exp(-totalLapseExpAnalysisResult.ExpectedYear4Ux)
			totalLapseExpAnalysisResult.ExpectedYear4AnnualRate = 1 - math.Pow(1-totalLapseExpAnalysisResult.ExpectedYear4MonthlyRate, 12.0)
		}
		if totalLapseExpAnalysisResult.TotalYear5PlusExposure > 0 {
			totalLapseExpAnalysisResult.ActualYear5PlusUx = totalLapseExpAnalysisResult.TotalActualYear5PlusLapses / totalLapseExpAnalysisResult.TotalYear5PlusExposure
			totalLapseExpAnalysisResult.ActualYear5PlusMonthlyRate = 1 - math.Exp(-totalLapseExpAnalysisResult.ActualYear5PlusUx)
			totalLapseExpAnalysisResult.ActualYear5PlusAnnualRate = 1 - math.Pow(1-totalLapseExpAnalysisResult.ActualYear5PlusMonthlyRate, 12.0)

			totalLapseExpAnalysisResult.ExpectedYear5PlusMonthlyRate = 1 - math.Exp(-totalLapseExpAnalysisResult.ExpectedYear5PlusUx)
			totalLapseExpAnalysisResult.ExpectedYear5PlusAnnualRate = 1 - math.Pow(1-totalLapseExpAnalysisResult.ExpectedYear5PlusMonthlyRate, 12.0)
		}

		err = DB.CreateInBatches(&expmps, 100).Error
		if err != nil {
			fmt.Println("Error saving expmps to DB: ", err)
			return err
		}

		err = DB.CreateInBatches(&crudeRateResults, 100).Error
		if err != nil {
			fmt.Println("Error saving CrudeRateResults to DB: ", err)
			return err
		}

		err = DB.CreateInBatches(&crudeLapseResults, 100).Error
		if err != nil {
			fmt.Println("Error saving CrudeLapseResults to DB: ", err)
			return err
		}

		err = DB.CreateInBatches(&totalMortalityExpAnalysisResult, 100).Error
		if err != nil {
			fmt.Println("Error saving Total Mortality Experience Analysis Results to DB: ", err)
			return err
		}

		err = DB.CreateInBatches(&totalLapseExpAnalysisResult, 100).Error
		if err != nil {
			fmt.Println("Error saving Total Lapse Experience Analysis Results to DB: ", err)
			return err
		}
		endTime := time.Since(startRunTime)
		fmt.Println("Time taken to run ExperienceAnalysis: ", endTime.Seconds(), " seconds")
		run.RunDuration = endTime.Seconds()
		runGroup.RunDuration += endTime.Seconds()
		//runGroup.ProcessedRecords += len(expmps)
		run.ProcessingStatus = "completed"
		//run.ProcessedRecords = run.TotalRecords
		err = DB.Save(&run).Error

	}

	runGroup.ProcessingStatus = "completed"
	err = DB.Save(&runGroup).Error
	if err != nil {
		fmt.Println("Error saving ExpRunGroup to DB: ", err)
		//return
	}

	endDataPrep := time.Since(startDataPrep)
	fmt.Println("Time taken to run ExperienceAnalysis: ", endDataPrep.Seconds(), " seconds")
	return nil
}

func MaxDate(date1, date2, date3 time.Time) time.Time {
	if date1.After(date2) {
		if date1.After(date3) {
			return date1
		}
	} else {
		if date2.After(date3) {
			return date2
		}
	}
	return date3
}

func MinDate(date1, date2, date3 time.Time) time.Time {
	if date1.Before(date2) {
		if date1.Before(date3) {
			return date1
		}
	} else {
		if date2.Before(date3) {
			return date2
		}
	}
	return date3
}
func DatetoStringConv(date1 time.Time) string {
	var convertedstring string
	var monthi int
	var monthstring, daystring string
	monthi = getMonths(date1.Month().String())

	if monthi <= 9 {
		monthstring += "0"
		monthstring += strconv.Itoa(monthi)
	} else {
		monthstring += strconv.Itoa(monthi)
	}

	if date1.Day() <= 9 {
		daystring += "0"
		daystring += strconv.Itoa(date1.Day())
	} else {
		daystring += strconv.Itoa(date1.Day())
	}

	convertedstring += strconv.Itoa(date1.Year())
	convertedstring += "-"
	convertedstring += monthstring
	convertedstring += "-"
	convertedstring += daystring

	return convertedstring
}

func getMonths(month string) int {
	monthi := 1
	for i := time.January; i <= time.December; i++ {
		//month := time.Month(i).String()
		if month == i.String() {
			return monthi
		}

		monthi += 1
	}
	return monthi
}

func GetExpMortalityRate(gender string, age int, version string) float64 {
	tablename := strings.ToLower("exp_current_mortalities")

	ExpAnalysisCacheKey := tablename + "_" + gender + "_" + strconv.Itoa(age) + "_" + version
	cached, found := Cache.Get(ExpAnalysisCacheKey)

	if found {
		result := cached.(float64)
		//if result > 0 {
		return result
		//}
	} else {
		//fmt.Println("cache missed: ", key)
	}

	var MortalityRow models.ExpCurrentMortality
	var qx float64
	if tablename != "" {
		query := "gender = ? and anb=? and version = ?"
		err := DB.Table(tablename).Where(query, gender, age, version).Find(&MortalityRow).Error
		if err != nil {
			fmt.Println("Mortality Table: ", err)
			Cache.Set(ExpAnalysisCacheKey, 0, 1)
			return 0
		}

		qx = MortalityRow.Qx
		Cache.Set(ExpAnalysisCacheKey, qx, 1)
	}
	return qx

}

func GetExpLapseRate(duration int, version string) float64 {
	tablename := strings.ToLower("exp_current_lapses")

	ExpAnalysisCacheKey := tablename + "_" + strconv.Itoa(duration) + "_" + version
	cached, found := Cache.Get(ExpAnalysisCacheKey)

	if found {
		result := cached.(float64)
		//if result > 0 {
		return result
		//}
	} else {
		//fmt.Println("cache missed: ", key)
	}

	var lapseRateRow models.ExpCurrentLapse
	var lapse_rate float64
	if tablename != "" {
		query := "duration_in_force_year = ? and version = ?"
		err := DB.Table(tablename).Where(query, duration, version).Find(&lapseRateRow).Error
		if err != nil {
			fmt.Println("Lapse Table: ", err)
			Cache.Set(ExpAnalysisCacheKey, 0, 1)
			return 0
		}
		lapse_rate = lapseRateRow.LapseRate
		Cache.Set(ExpAnalysisCacheKey, lapse_rate, 1)
	}
	return lapse_rate

}
