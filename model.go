package ruleenginecore

import (
	"fmt"
)

type Input map[string]string

// defines an output of a successful rule evaluation
type Output struct {
	Rulename string `json:"rulename"`
	Priority int    `json:"priority"`
	// the same result as an entity defined as part rule engine config for every rule.
	Result map[string]any `json:"result"`
}

// defines a mandatory fields (fieldname and fieldtype) for rule engine input.
// as part engine evaluation field value is picked from the input via given name.
//
// valid field types are 'int', 'float', 'bool', 'string'
// 'int'  is 64 bit signed integer
// 'float' is 64 bit signed float
// 'bool' is boolean true or false
// 'string' is UTF-8 based string
type Fields map[string]string

// define an Operand for any operator
// as part of evaluation process operand value is determined as following:
// if operand.OperandAs is 'field'
//
//	operand.Val is fieldname, Operand value is determined from the input having fieldname as <operand.Val>
//
// if operand.OperandAs is 'constant'
//
//	operand.Val is considered as operand value in a string form.
type Operand struct {
	// define operand as field or constant. valid values are 'field' and 'constant'
	OperandAs  string `json:"operandAs"`
	Val        string `json:"val"`
	typedValue any    `json:"-"`
}

func (op *Operand) getValue(input typedValueMap) any {
	switch op.OperandAs {
	case OperandAsField:
		return input[op.Val]
	case OperandAsConstant:
		return op.typedValue
	}
	panic("Invalid OperandAs " + op.OperandAs)
}

// defines a custom condition type which can be used as part of rule definition
// Valid Operators are '>','>=','<','<=','==', '!=', 'contain'
//
//	'>','>=','<','<=' operators supports 'int', 'float' operandType
//	'>','>=','<','<=','==', '!=' operator support 'int','float','bool','string' operandType
//	'contain' operator supports 'string' operandType
//	Operator, OperandType and Operands can be defined as following
type ConditionType struct {
	Operator    string     `json:"operator"`
	OperandType string     `json:"operandType"`
	Operands    []*Operand `json:"operands"`
}

type Condition struct {
	ConditionType string       `json:"conditionType"`
	SubConditions []*Condition `json:"subConditions"`
}

type Rule struct {
	Priority      int            `json:"priority"`
	RootCondition *Condition     `json:"condition"`
	Result        map[string]any `json:"result"`
}

type RuleEngineConfig struct {
	Fields         Fields                    `json:"fields"`
	ConditionTypes map[string]*ConditionType `json:"conditionTypes"`
	Rules          map[string]*Rule          `json:"rules"`
}

type RuleEngineError struct {
	ComponentName string
	ErrMsg        string
	ErrCode       uint
	OtherMsg      string
}

func (ce RuleEngineError) Error() string {
	return fmt.Sprintf("RuleEngineError: for %v, %v , code:%v otherMsg:%v.", ce.ComponentName, ce.ErrMsg, ce.ErrCode, ce.OtherMsg)
}

func NewError(errCode uint, componentName string, otherMsg string) *RuleEngineError {
	return &RuleEngineError{ErrCode: errCode, ComponentName: componentName, ErrMsg: errCodeToMessage[errCode], OtherMsg: otherMsg}
}

const (
	ErrCodeNone = iota
	ErrCodeInvalidValueType
	ErrCodeInvalidOperator
	ErrCodeInvalidOperandType
	ErrCodeInvalidOperandsLength
	ErrCodeInvalidConditionType
	ErrCodeInvalidSubConditionCount
	ErrCodeConditionTypeNotFound
	ErrCodeInvalidOperandAs
	ErrCodeFieldNotFound
	ErrCodeFailedParsingInput
	ErrCodeRuleNotFound
	ErrCodeInvalidEvaluateOperations
)

var errCodeToMessage = map[uint]string{
	ErrCodeInvalidValueType:          "Passed value type is invalid",
	ErrCodeInvalidOperator:           "Invalid operator",
	ErrCodeInvalidOperandType:        "Invalid operandType",
	ErrCodeInvalidOperandsLength:     "Invalid number of operands",
	ErrCodeInvalidConditionType:      "Invalid conditionType",
	ErrCodeInvalidSubConditionCount:  "Invalid sub-condition count",
	ErrCodeConditionTypeNotFound:     "Could not find conditionType",
	ErrCodeInvalidOperandAs:          "Invalid operandAs",
	ErrCodeFieldNotFound:             "Field not found",
	ErrCodeFailedParsingInput:        "Could not parse input",
	ErrCodeRuleNotFound:              "Rule not found",
	ErrCodeInvalidEvaluateOperations: "Invalid evaluate options value n",
}
