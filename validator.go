package ruleenginecore

import (
	"strings"
)

var validValueTypeMap = map[string]bool{IntType: true, FloatType: true, BoolType: true, StringType: true}
var printableValidValueType string

var validOperatorMap = map[string]bool{GreaterOperator: true, GreaterEqualOperator: true, LessOperator: true, LessEqualOperator: true, EqualOperator: true, NotEqualOperator: true, ContainOperator: true}

func init() {
	validFields := []string{}
	for k := range validValueTypeMap {
		validFields = append(validFields, k)
	}
	printableValidValueType = strings.Join(validFields[:], ", ")
}

// check for
// 1. is field exist with <name>
// 2. does field having <name> have type <fType>
func (fs Fields) isValid(name, fType string) *RuleEngineError {
	if ft, ok := fs[name]; !ok {
		return newError(ErrCodeFieldNotFound, "Field:"+name, "")
	} else if ft != fType {
		return newError(ErrCodeInvalidValueType, "field:"+name, "field is having type:"+ft+" condition is expecting type:"+fType)
	}
	return nil
}

// validates fields for valid type
func (fs Fields) validate() *RuleEngineError {
	for name, fType := range fs {
		if _, ok := validValueTypeMap[fType]; !ok {
			return newError(ErrCodeInvalidValueType, "field:"+name, "fieldType is invalid. valid types are "+printableValidValueType)
		}
	}
	return nil
}

// validating condition type for following
// 1. valid condition operator
// 2. valid number of operands for given operator
// 3. valid operand type for given operator
// 4. for OperandAsField, match operandType with field type.
func (c *ConditionType) validateAndParseValues(name string, fs Fields) *RuleEngineError {

	if _, ok := validOperatorMap[c.Operator]; !ok {
		return newError(ErrCodeInvalidOperator, "condition:"+name, c.Operator+" is invalid operator")
	}

	if c.Operator == ContainOperator {
		if len(c.Operands) != 2 {
			return newError(ErrCodeInvalidOperandsLength, "condition:"+name, c.Operator+" operator expects exactly two operands")
		}

		if c.OperandType != StringType {
			return newError(ErrCodeInvalidOperandType, "condition:"+name, c.Operator+" operator only supports '"+StringType+"' type")
		}

		for _, operand := range c.Operands {
			if operand.Type == FieldType {
				if err := fs.isValid(operand.Val, c.OperandType); err != nil {
					return newError(err.ErrCode, "condition:"+name, err.OtherMsg)
				}
			} else if operand.Type == ConstantType {
				if typedValue, err := getTypedValue(operand.Val, c.OperandType); err != nil {
					return newError(ErrCodeFailedParsingInput, "condition:"+name, err.Error())
				} else {
					operand.typedValue = typedValue
				}
			} else {
				return newError(ErrCodeInvalidOperandAs, "condition:"+name, operand.Type.String()+" is invalid. valid operandAs are "+FieldType.String()+", "+ConstantType.String())
			}
		}
	}

	if c.Operator == EqualOperator || c.Operator == NotEqualOperator {
		if len(c.Operands) != 2 {
			return newError(ErrCodeInvalidOperandsLength, "For condition:"+name, c.Operator+" operator expects exactly two operands")
		}

		if _, ok := validValueTypeMap[c.OperandType]; !ok {
			return newError(ErrCodeInvalidOperandType, "condition:"+name, c.Operator+" operator supports "+printableValidValueType+" types")
		}

		for _, operand := range c.Operands {
			if operand.Type == FieldType {
				if err := fs.isValid(operand.Val, c.OperandType); err != nil {
					return newError(err.ErrCode, "condition:"+name, err.OtherMsg)
				}
			} else if operand.Type == ConstantType {
				if typedValue, err := getTypedValue(operand.Val, c.OperandType); err != nil {
					return newError(ErrCodeFailedParsingInput, "condition:"+name, err.Error())
				} else {
					operand.typedValue = typedValue
				}
			} else {
				return newError(ErrCodeInvalidOperandAs, "condition:"+name, operand.Type.String()+" is invalid. valid operandAs are "+FieldType.String()+", "+ConstantType.String())
			}
		}
	}

	if c.Operator == GreaterOperator || c.Operator == GreaterEqualOperator || c.Operator == LessOperator || c.Operator == LessEqualOperator {

		if len(c.Operands) != 2 {
			return newError(ErrCodeInvalidOperandsLength, "condition:"+name, c.Operator+" operator expects exactly two operands")
		}

		if c.OperandType != IntType && c.OperandType != FloatType {
			return newError(ErrCodeInvalidOperandType, "condition:"+name, c.Operator+" operator supports "+IntType+", "+FloatType+" types")
		}

		for _, operand := range c.Operands {
			if operand.Type == FieldType {
				if err := fs.isValid(operand.Val, c.OperandType); err != nil {
					return newError(err.ErrCode, "condition:"+name, err.OtherMsg)
				}
			} else if operand.Type == ConstantType {
				if typedValue, err := getTypedValue(operand.Val, c.OperandType); err != nil {
					return newError(ErrCodeFailedParsingInput, "condition:"+name, err.Error())
				} else {
					operand.typedValue = typedValue
				}
			} else {
				return newError(ErrCodeInvalidOperandAs, "condition:"+name, operand.Type.String()+" is invalid. valid operandAs are "+FieldType.String()+", "+ConstantType.String())
			}
		}
	}

	return nil
}

// validating rule for following
// 1. valid condition logical operator ( 'and' | 'or' )
// 2. validate innerCondition for either logical condition or customConditionTypes
func (r *Rule) validate(name string, custConditionType map[string]*ConditionType) *RuleEngineError {
	return validateCondition(name, r.RootCondition, custConditionType)
}

func validateCondition(name string, c *Condition, custConditionType map[string]*ConditionType) *RuleEngineError {

	if c.ConditionType == OrOperator || c.ConditionType == AndOperator {

		if len(c.SubConditions) < 2 {
			return newError(ErrCodeInvalidSubConditionCount, "rule:"+name, "conditionType:"+c.ConditionType+" expects at-least 2 sub-conditions")
		}

		for _, cond := range c.SubConditions {
			if err := validateCondition(name, cond, custConditionType); err != nil {
				return err
			}
		}
		return nil
	}

	if c.ConditionType == NegationOperator {

		if len(c.SubConditions) != 1 {
			return newError(ErrCodeInvalidSubConditionCount, "rule:"+name, "conditionType:"+c.ConditionType+" expects exactly one sub-condition")
		}

		for _, cond := range c.SubConditions {
			if err := validateCondition(name, cond, custConditionType); err != nil {
				return err
			}
		}
		return nil
	}

	if _, ok := custConditionType[c.ConditionType]; ok {
		return nil
	}

	return newError(ErrCodeConditionTypeNotFound, "rule:"+name, "conditionType:"+c.ConditionType+" not found")
}

// validates for valid types, operators, condition definition and rule definition.
func (c *RuleEngineConfig) Validate() *RuleEngineError {
	err := c.Fields.validate()
	if err != nil {
		return err
	}

	for conditionName, customCondition := range c.ConditionTypes {
		if err := customCondition.validateAndParseValues(conditionName, c.Fields); err != nil {
			return err
		}
	}

	for ruleName, rule := range c.Rules {
		if err := rule.validate(ruleName, c.ConditionTypes); err != nil {
			return err
		}
	}

	return nil
}

// validates for mandatory fields and type conversion from string to respective field type
func (input Input) validateAndParseValues(fs Fields) (typedValueMap, *RuleEngineError) {
	ret := typedValueMap{}
	for fieldname, fieldtype := range fs {
		strVal, found := input[fieldname]
		if !found {
			return nil, newError(ErrCodeFieldNotFound, "input", "field:"+fieldname+" with type:"+fieldtype+" is mandatory")
		}

		if val, err := getTypedValue(strVal, fieldtype); err != nil {
			return nil, newError(ErrCodeFailedParsingInput, "input:"+fieldname, err.Error())
		} else {
			ret[fieldname] = val
		}
	}

	return ret, nil
}
