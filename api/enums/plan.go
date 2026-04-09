package enums

type Plan int

const (
	A Plan = iota + 1
	B
	C
	D
)

func (planType Plan) String() string {
	mTypes := []string{
		"A",
		"B",
		"C",
		"D",
	}
	return mTypes[planType-1]
}
