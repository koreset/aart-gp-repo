package enums

type CommissionType int

const (
	Type1 CommissionType = iota + 1
	Type2
	Type3
	Type4
)

func (commissionType CommissionType) String() string {
	mTypes := []string{
		"Type1",
		"Type2",
		"Type3",
		"Type4",
	}
	return mTypes[commissionType-1]
}
