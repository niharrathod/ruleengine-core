package ruleenginecore

import (
	"fmt"
	"strings"
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

// Supported default ConditionTypes
const (
	NegationCondition = "not"
	AndCondition      = "and"
	OrCondition       = "or"
)

// 'Input' defines an input for rule evaluation as map of fieldname as key and string representation of value as (map)value
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

func newOutput(ruleName string, priority int, result map[string]any) *Output {
	return &Output{
		Rulename: ruleName,
		Priority: priority,
		Result:   result,
	}
}

// 'fields' defines a mandatory input for RuleEngine evaluation, internally it represents as map of fieldname(string) as key and ValueType as value
//
// As part engine evaluation, for field value is picked from the input with fieldname.
type Fields map[string]ValueType

func (fs Fields) exist(fieldName string, expectedFieldType ValueType) bool {
	actualFieldType, ok := fs[fieldName]
	if !ok {
		return false
	}

	if actualFieldType != expectedFieldType {
		return false
	}

	return true
}

// 'Operand' defines an operand for custom ConditionType
//
// as part of evaluation process operand value is determined as following:
//
//	if operand.OperandType is 'field'
//		->operand.Val is fieldname, Operand value is determined from the input having fieldname as <operand.Val>
//	if operand.OperandType is 'constant'
//		->operand.Val is considered as operand value in a string form.
type Operand struct {
	// 'ValueType' defined as type of operand value
	ValueType ValueType `json:"valuetype"`

	// 'Type' define type of an operand as either field or constant
	Type OperandType `json:"type"`

	// 'Val' is value of an operand
	//
	// for OperandType as 'field'
	//		-> Val is considered as 'fieldname', while evaluation, value is picked from input as operand for evaluation
	// for OperandType as 'constant'
	//		-> Val is considered as value and picked as operand for evaluation.

	Val        string `json:"value"`
	typedValue any    `json:"-"`
}

func (op *Operand) isField() bool {
	return op.Type == Field
}

func (op *Operand) getValue(values parsedInput) any {
	switch op.Type {
	case Field:
		return values[op.Val]
	case Constant:
		return op.typedValue
	}
	// no-op
	panic(fmt.Sprintf("Invalid OperandType %v", op.ValueType))
}

// 'ConditionType' defines a custom condition type, which is be used while defining a rule
// Valid Operators are '>','>=','<','<=','==', '!=', 'contain'
//
//	'>','>=','<','<=' operators supports 'int', 'float' operand valueType
//	'==', '!=' operator support 'int','float','bool','string' operand valueType
//	'contain' operator supports 'string' operand valueType
type ConditionType struct {
	Operator string     `json:"operator"`
	Operands []*Operand `json:"operands"`
}

// 'Condition' define condition for a Rule which needs to be satisfy to consider rule a matched.
type Condition struct {
	// 'Type' sets type of a condition, either logical such as 'and','or','not' or types defined as 'ConditionTypes' with RuleEngineConfig
	Type          string       `json:"type"`
	SubConditions []*Condition `json:"subConditions"`
}

// 'RuleConfig' defines a rule for RuleEngine.
// Rule is internally n-ary tree, where every node is a Condition, inner nodes of a tree are type of 'and', 'or' or 'not' and leaf nodes are
// custom conditions from ConditionTypes defined by user
type RuleConfig struct {

	// 'Priority' is rule priority, evaluation operation prioritize the rule based of this value
	Priority int `json:"priority"`

	// 'RootCondition' defines n-ary tree of Rule
	RootCondition *Condition `json:"condition"`

	// 'Result' defines key-value container maintains values and returns as part of 'Output' if Rule matches.
	Result map[string]any `json:"result"`
}

// 'RuleEngineConfig' is a configuration for a RuleEngine. defines mandatory input fields, custom conditions and rule
type RuleEngineConfig struct {
	// 'Fields' defines mandatory as input for rule engine evaluation
	Fields Fields `json:"fields"`

	// 'ConditionTypes' defines custom condition as map having condition name as key, ConditionType as value
	ConditionTypes map[string]*ConditionType `json:"conditionTypes"`

	// 'Rules' defines set of rules for ruleengine, as map having rule name as key, RuleConfig as value
	Rules map[string]*RuleConfig `json:"rules"`
}

type parsedInput map[string]any

type RuleEngineError struct {
	ErrCode  uint
	ErrMsg   string
	OtherMsg string
}

func (ce *RuleEngineError) Error() string {
	return fmt.Sprintf("RuleEngineError: ErrCode:%v ErrMsg:%v. %v", ce.ErrCode, ce.ErrMsg, ce.OtherMsg)
}

func (ce *RuleEngineError) addMsg(msg string) {
	ce.OtherMsg = fmt.Sprintf("%v, %v", ce.OtherMsg, msg)
}

func newError(errCode uint, msgs ...string) *RuleEngineError {
	return &RuleEngineError{ErrCode: errCode, ErrMsg: errCodeToMessage[errCode], OtherMsg: strings.Join(msgs, ",")}
}

const (
	ErrCodeInvalidValueType = iota + 1
	ErrCodeInvalidOperator
	ErrCodeInvalidOperandType
	ErrCodeInvalidOperandsLength
	ErrCodeInvalidConditionType
	ErrCodeInvalidSubConditionCount
	ErrCodeConditionTypeNotFound
	ErrCodeFieldNotFound
	ErrCodeParsingFailed
	ErrCodeRuleNotFound
	ErrCodeInvalidEvaluateOperations
	ErrCodeContextCancelled
	ErrCodeInvalidOperand
)

var errCodeToMessage = map[uint]string{
	ErrCodeInvalidValueType:          "Invalid ValueType",
	ErrCodeInvalidOperator:           "Invalid operator",
	ErrCodeInvalidOperandType:        "Invalid operandType",
	ErrCodeInvalidOperandsLength:     "Invalid number of operands",
	ErrCodeInvalidConditionType:      "Invalid conditionType",
	ErrCodeInvalidSubConditionCount:  "Invalid sub-condition count",
	ErrCodeConditionTypeNotFound:     "Could not find conditionType",
	ErrCodeFieldNotFound:             "Field not found",
	ErrCodeParsingFailed:             "Could not parse value",
	ErrCodeRuleNotFound:              "Rule not found",
	ErrCodeInvalidEvaluateOperations: "Invalid evaluate options value n",
	ErrCodeContextCancelled:          "Context is cancelled",
	ErrCodeInvalidOperand:            "Invalid operandtype or valuetype",
}
