package ruleenginecore

import (
	"encoding/json"
	"fmt"
	"strings"
)

// 'ValueType' defines supported value types for rule engine
type ValueType uint8

const (
	unknownValueType ValueType = iota

	//	'Boolean' is boolean true or false
	Boolean

	//	'String' is UTF-8 based string
	String

	//	'Integer' is 64 bit signed integer
	Integer

	//	'Float' is 64 bit signed float
	Float
)

var (
	valueType_Name = map[ValueType]string{
		1: "Boolean",
		2: "String",
		3: "Integer",
		4: "Float",
	}
	valueType_Value = map[string]ValueType{
		"bool":    1,
		"Bool":    1,
		"boolean": 1,
		"Boolean": 1,
		"string":  2,
		"String":  2,
		"int":     3,
		"Int":     3,
		"integer": 3,
		"Integer": 3,
		"float":   4,
		"Float":   4,
	}
)

// 'isValid' check for valid ValueType
func (valueType ValueType) isValid() bool {
	_, ok := valueType_Name[valueType]
	return ok
}

// 'String' gives string representation of ValueType
func (valueType ValueType) String() string {
	val, ok := valueType_Name[valueType]
	if !ok {
		return fmt.Sprintf("ValueType(%v)", uint8(valueType))
	}
	return val
}

func parseValueType(s string) (ValueType, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	value, ok := valueType_Value[s]
	if !ok {
		return unknownValueType, fmt.Errorf("invalid ValueType(%v)", s)
	}
	return ValueType(value), nil
}

func (valueType ValueType) MarshalJSON() ([]byte, error) {
	if valueType.isValid() {
		return json.Marshal(valueType.String())
	}
	return nil, fmt.Errorf("invalid ValueType(%v)", valueType)
}

func (valueType *ValueType) UnmarshalJSON(data []byte) error {
	var val string
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}

	result, err := parseValueType(val)
	if err != nil {
		return err
	}

	*valueType = result
	return nil
}

func (valueType ValueType) isBoolean() bool {
	return valueType == Boolean
}

func (valueType ValueType) isInteger() bool {
	return valueType == Integer
}

func (valueType ValueType) isFloat() bool {
	return valueType == Float
}

func (valueType ValueType) isString() bool {
	return valueType == String
}

// comma separated value type list
var valueTypeList string

func init() {
	validTypes := []string{}

	for _, valueType := range valueType_Name {
		validTypes = append(validTypes, valueType)
	}

	valueTypeList = strings.Join(validTypes[:], ", ")
}
