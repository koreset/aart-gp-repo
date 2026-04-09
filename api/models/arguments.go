package models

type TransitionRateArguments struct {
	ProductId          	int
	ProductCode        	string
	Year               	int
	Age                	int
	Gender             	string
	SmokerStatus       	string
	Income             	int
	SocioEconomicClass 	int
	OccupationalClass  	string
	SelectPeriod       	int
	EducationLevel     	int
	DurationIfM		   	int
	ProjectionMonth	   	int
	DistributionChannel string
}
