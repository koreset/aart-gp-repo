package enums

type MemberType int

const (
	MainMember MemberType = iota + 1
	Spouse
	Parent
	Extended
	Child
)

var toMTString = map[MemberType]string{
	MainMember: "MM",
	Spouse:     "SP",
	Parent:     "PAR",
	Extended:   "EXt",
	Child:      "CH",
}

var toMTID = map[string]MemberType{
	"MM":  MainMember,
	"SP":  Spouse,
	"PAR": Parent,
	"EXT": Extended,
	"CH":  Child,
}

func (memberType MemberType) String() string {
	return toMTString[memberType]
}

//type VariableType int
//
//const (
//	Integer VariableType = iota + 1
//	Float
//	String
//	Boolean
//)
//
//var toString = map[VariableType]string{
//	Integer: "Integer",
//	Float:   "Float",
//	String:  "String",
//	Boolean: "Boolean",
//}
//
//var toID = map[string]VariableType{
//	"1": Integer,
//	"2": Float,
//	"3": String,
//	"4": Boolean,
//}
//
//func (v VariableType) String() string {
//	return toString[v]
//}
//
//// MarshalJSON marshals the enum as a quoted json string
//func (v VariableType) MarshalJSON() ([]byte, error) {
//	buffer := bytes.NewBufferString(`"`)
//	buffer.WriteString(toString[v])
//	buffer.WriteString(`"`)
//	return buffer.Bytes(), nil
//}
//
//// UnmarshalJSON unmashals a quoted json string to the enum value
//func (v *VariableType) UnmarshalJSON(b []byte) error {
//	var j string
//	err := json.Unmarshal(b, &j)
//	if err != nil {
//		return err
//	}
//	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
//	*v = toID[j]
//	return nil
//}
