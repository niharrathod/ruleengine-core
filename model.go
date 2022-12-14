package ruleenginecore

import (
	"fmt"
)

// Supported value types
const (
	//	'int'  is 64 bit signed integer
	IntType = "int"

	// 'float' is 64 bit signed float
	FloatType = "float"

	// 'bool' is boolean true or false
	BoolType = "bool"

	// 'string' is UTF-8 based string
	StringType = "string"
)

// Supported OperandAs
const (
	// operand is considered as field and value is determined from the Input
	OperandAsField = "field"

	// operand is considered as constant and value is determined from Operand.Val
	OperandAsConstant = "constant"
)

// Supported operators to define custom ConditionType
const (
	GreaterOperator      = ">"
	GreaterEqualOperator = ">="
	LessOperator         = "<"
	LessEqualOperator    = "<="
	EqualOperator        = "=="
	NotEqualOperator     = "!="
	ContainOperator      = "contain"
)

// Supported logical operators to defining condition for rule
const (
	NegationOperator = "not"
	AndOperator      = "and"
	OrOperator       = "or"
)

// defines an input for rule evaluation as map of fieldname as key and string representation of value as (map)value
type Input map[string]string

// defines an output for evaluation result
type Output struct {
	// matched rulename
	Rulename string `json:"rulename"`

	// priority of matched rule
	Priority int `json:"priority"`

	// matched rule result, defined as part of RuleEngineConfig for every rule
	Result map[string]any `json:"result"`
}

// defines a mandatory fields for RuleEngine Input,  as map of fieldname as key and fieldtype as value
//
// As part engine evaluation, for field value is picked from the input with fieldname.
//
//	valid fieldtypes are 'int', 'float', 'bool', 'string'
//	'int'  is 64 bit signed integer
//	'float' is 64 bit signed float
//	'bool' is boolean true or false
//	'string' is UTF-8 based string
type Fields map[string]string

// define an Operand for custom ConditionType
//
// as part of evaluation process operand value is determined as following:
//
//	if operand.OperandAs is 'field'
//		->operand.Val is fieldname, Operand value is determined from the input having fieldname as <operand.Val>
//	if operand.OperandAs is 'constant'
//		->operand.Val is considered as operand value in a string form.
type Operand struct {
	// define operand as field or constant. valid values are 'field' and 'constant'
	OperandAs string `json:"operandAs"`

	// value of an operand
	//
	// for operandAs 'field'
	//		-> Val is considered as 'fieldname', while evaluation, value is picked from input
	// for operandAs 'constant'
	//		-> Val is considered as value of operand for evaluation.

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

// defines a custom condition type
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

// define condition for a Rule which needs to be satisfy to consider rule a matched. Every Condition either logical such as 'and','or','not' or type of custom conditions defined as 'ConditionTypes'
type Condition struct {
	ConditionType string       `json:"conditionType"`
	SubConditions []*Condition `json:"subConditions"`
}

// defines a rule for RuleEngine.
type Rule struct {
	Priority      int            `json:"priority"`
	RootCondition *Condition     `json:"condition"`
	Result        map[string]any `json:"result"`
}

// Configuration for a RuleEngine. defines mandatory input fields, custom conditions and rule
type RuleEngineConfig struct {
	// defines fields, mandatory values for input as map having fieldname as key, fieldtype as value
	Fields Fields `json:"fields"`

	// defines custom ConditionTypes as map having conditionTypeName as key, ConditionType as value
	ConditionTypes map[string]*ConditionType `json:"conditionTypes"`

	// defines Rules for ruleengine, as map having rulename as key, rule as value
	Rules map[string]*Rule `json:"rules"`
}

type RuleEngineError struct {
	ComponentName string
	ErrMsg        string
	ErrCode       uint
	OtherMsg      string
}

func (ce *RuleEngineError) Error() string {
	return fmt.Sprintf("RuleEngineError: for %v, %v , code:%v otherMsg:%v.", ce.ComponentName, ce.ErrMsg, ce.ErrCode, ce.OtherMsg)
}

func newError(errCode uint, componentName string, otherMsg string) *RuleEngineError {
	return &RuleEngineError{ErrCode: errCode, ComponentName: componentName, ErrMsg: errCodeToMessage[errCode], OtherMsg: otherMsg}
}

const (
	ErrCodeInvalidValueType = iota + 1
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
	ErrCodeContextCancelled
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
	ErrCodeContextCancelled:          "context is cancelled",
}
