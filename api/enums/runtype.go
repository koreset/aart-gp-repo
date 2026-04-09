package enums

type RunType int

const (
	Valuation VariableType = iota
	Pricing
)

func (runType RunType) ToInt() int {
	mTypes := []int{
		0,
		1,
	}
	return mTypes[runType]
}
