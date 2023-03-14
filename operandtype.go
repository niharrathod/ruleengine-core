package ruleenginecore

import (
	"encoding/json"
	"fmt"
	"strings"
)

// 'OperandType' defines type of operand either 'Field' or 'Constant'
type OperandType uint8

const (
	unknownOperandType OperandType = iota

	// 'Field' is OperandType where operand value is derived from the Input using name as Operand.Val
	Field

	// 'Constant' is OperandType where operand value is considered as Operand.Val
	Constant
)

var (
	operandType_Name = map[OperandType]string{
		1: "Field",
		2: "Constant",
	}
	operandType_Value = map[string]OperandType{
		"field":    1,
		"Field":    1,
		"constant": 2,
		"Constant": 2,
	}
)

// 'isValid' check for valid operandType starting from 1("Field"),2("Constant")
func (operandType OperandType) isValid() bool {
	_, ok := operandType_Name[operandType]
	return ok
}

func (operandType OperandType) String() string {
	val, ok := operandType_Name[operandType]
	if !ok {
		return fmt.Sprintf("OperandType(%v)", uint8(operandType))
	}
	return val
}

func parseOperandType(s string) (OperandType, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	value, ok := operandType_Value[s]
	if !ok {
		return unknownOperandType, fmt.Errorf("invalid OperandType(%v)", s)
	}
	return value, nil
}

func (operandType OperandType) MarshalJSON() ([]byte, error) {
	if operandType.isValid() {
		return json.Marshal(operandType.String())
	}
	return nil, fmt.Errorf("invalid OperandType(%v)", uint8(operandType))
}

func (operandType *OperandType) UnmarshalJSON(data []byte) error {
	var val string
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}

	result, err := parseOperandType(val)
	if err != nil {
		return err
	}

	*operandType = result
	return nil
}

// comma separated operandType list
var operandTypeList string

func init() {
	validOperandType := []string{}

	for _, value := range operandType_Name {
		validOperandType = append(validOperandType, value)
	}

	operandTypeList = strings.Join(validOperandType[:], ", ")
}
