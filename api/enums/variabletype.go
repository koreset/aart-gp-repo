package enums

import (
	"bytes"
	"encoding/json"
)

type VariableType int

const (
	Integer VariableType = iota + 1
	Float
	String
	Boolean
)

var toString = map[VariableType]string{
	Integer: "Integer",
	Float:   "Float",
	String:  "String",
	Boolean: "Boolean",
}

var toID = map[string]VariableType{
	"Integer": Integer,
	"Float":   Float,
	"String":  String,
	"Boolean": Boolean,
}

func (v VariableType) String() string {
	return toString[v]
}

// MarshalJSON marshals the enum as a quoted json string
func (v VariableType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toString[v])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (v *VariableType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*v = toID[j]
	return nil
}
