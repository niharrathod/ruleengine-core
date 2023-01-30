package ruleenginecore

import (
	"encoding/json"
	"fmt"
	"strings"
)

// 'OperandType' is either field or constant
type OperandType uint8

const (
	// 'FieldType' operand is considered as field and value is determined from the Input
	FieldType OperandType = iota + 1

	// 'ConstantType' operand is considered as constant and value is determined from Operand.Val
	ConstantType
)

var (
	OperandType_Name = map[uint8]string{
		1: "FieldType",
		2: "ConstantType",
	}
	OperandType_Value = map[string]uint8{
		"FieldType":    1,
		"fieldtype":    1,
		"field":        1,
		"ConstantType": 2,
		"constanttype": 2,
		"constant":     2,
	}
)

// IsValid check for valid operandType starting from 1("Field"),2("Constant")
func IsValid(operandType OperandType) bool {
	_, ok := OperandType_Name[uint8(operandType)]
	return ok
}

func (operandType OperandType) String() string {
	val, ok := OperandType_Name[uint8(operandType)]
	if !ok {
		panic(fmt.Sprint("Could not identify OperandType with value", uint8(operandType)))
	}
	return val
}

func ParseOperandType(s string) (OperandType, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	value, ok := OperandType_Value[s]
	if !ok {
		return OperandType(0), fmt.Errorf("%v is not the value OperandType", s)
	}
	return OperandType(value), nil
}

func (operandType OperandType) MarshalJSON() ([]byte, error) {
	if IsValid(operandType) {
		return json.Marshal(operandType.String())
	}
	return nil, fmt.Errorf("%v is not valid OperandType", operandType)
}

func (operandType *OperandType) UnmarshalJSON(data []byte) error {
	var val string
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}

	result, err := ParseOperandType(val)
	if err != nil {
		return err
	}

	*operandType = result
	return nil
}
