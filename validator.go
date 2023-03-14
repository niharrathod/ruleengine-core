package ruleenginecore

import (
	"fmt"
)

type fieldValidatorFunc func(fs Fields) *RuleEngineError
type conditionTypeValidatorFunc func(ct *ConditionType, fs Fields) *RuleEngineError
type ruleConditionValidatorFunc func(c *Condition) *RuleEngineError

var fieldValueTypeValidator = func() fieldValidatorFunc {
	return func(fs Fields) *RuleEngineError {
		for fieldName, fieldValueType := range fs {
			if !fieldValueType.isValid() {
				return newError(ErrCodeInvalidValueType,
					fmt.Sprintf("field: %v", fieldName),
					fmt.Sprintf("valid valueTypes are %v.", valueTypeList))
			}
		}
		return nil
	}
}

var operandCountValidator = func(count int) conditionTypeValidatorFunc {
	return func(ct *ConditionType, fs Fields) *RuleEngineError {
		if len(ct.Operands) != count {
			return newError(ErrCodeInvalidOperandsLength,
				fmt.Sprintf("expected operand count is  %v", count))
		}
		return nil
	}
}

var operandsWithSameValueTypeValidator = func() conditionTypeValidatorFunc {
	return func(ct *ConditionType, fs Fields) *RuleEngineError {
		var valueType *ValueType = nil
		for _, op := range ct.Operands {
			if valueType == nil {
				valueType = &op.ValueType
				continue
			}

			if *valueType != op.ValueType {
				return newError(ErrCodeInvalidOperand,
					"Expecting same valueType for all the operands")
			}
		}

		return nil
	}
}

var operandValueTypeValidator = func(supportedValueTypes ...ValueType) conditionTypeValidatorFunc {
	var supportedTypeSet = NewSet[ValueType]()
	var supportedTypeCommaSepStr string
	for _, valueType := range supportedValueTypes {
		supportedTypeSet.Add(valueType)
		if len(supportedTypeCommaSepStr) == 0 {
			supportedTypeCommaSepStr = valueType.String()
		} else {
			supportedTypeCommaSepStr = fmt.Sprintf("%v, %v", supportedTypeCommaSepStr, valueType)
		}
	}

	return func(ct *ConditionType, fs Fields) *RuleEngineError {
		for _, op := range ct.Operands {
			if !supportedTypeSet.Contains(op.ValueType) {
				return newError(ErrCodeInvalidOperand,
					fmt.Sprintf("Operand having invalid valueType. Supported valueType are %v", supportedTypeCommaSepStr))
			}
		}
		return nil
	}
}

var operandValidator = func() conditionTypeValidatorFunc {
	return func(ct *ConditionType, fs Fields) *RuleEngineError {
		for _, op := range ct.Operands {
			if err := validateAndParseOperand(op, fs); err != nil {
				return err
			}
		}
		return nil
	}
}

func validateAndParseOperand(operand *Operand, fs Fields) *RuleEngineError {

	if !operand.valid() {
		return newError(ErrCodeInvalidOperand, "Invalid ValueType or OperandType")
	}

	// Field operandType
	if operand.isField() {
		if !fs.exist(operand.Val, operand.ValueType) {
			return newError(ErrCodeFieldNotFound,
				fmt.Sprintf("Expecting field: %v with valueType: %v", operand.Val, operand.ValueType))
		}
		return nil
	}

	// Constant operandType
	typedValue, err := parseValue(operand.Val, operand.ValueType)
	if err != nil {
		err.addMsg(fmt.Sprintf("Constant operand with value: %v failed to parse as ValueType: %v",
			operand.Val, operand.ValueType))
		return err
	}

	operand.typedValue = typedValue
	return nil
}

var subConditionCountRuleConditionValidator = func(count int) ruleConditionValidatorFunc {
	return func(c *Condition) *RuleEngineError {
		if len(c.SubConditions) != count {
			return newError(ErrCodeInvalidSubConditionCount)
		}
		return nil
	}
}

var minSubConditionCountRuleConditionValidator = func(count int) ruleConditionValidatorFunc {
	return func(c *Condition) *RuleEngineError {
		if len(c.SubConditions) < count {
			return newError(ErrCodeInvalidSubConditionCount)
		}
		return nil
	}
}

type ruleEngineConfigValidator struct {
	fieldValidators         []fieldValidatorFunc
	condTypeValidators      map[string][]conditionTypeValidatorFunc
	ruleConditionValidators map[string][]ruleConditionValidatorFunc
}

func (v *ruleEngineConfigValidator) addFieldValidator(validators ...fieldValidatorFunc) {
	v.fieldValidators = validators
}

func (v *ruleEngineConfigValidator) addConditionTypeValidator(operator string, validators ...conditionTypeValidatorFunc) {
	v.condTypeValidators[operator] = validators
}

func (v *ruleEngineConfigValidator) addRuleConditionValidator(operator string, validators ...ruleConditionValidatorFunc) {
	v.ruleConditionValidators[operator] = validators
}

func (v *ruleEngineConfigValidator) validateConditionType(ct *ConditionType, fs Fields) *RuleEngineError {
	validatorFuncs := v.condTypeValidators[ct.Operator]

	for _, validatorFunc := range validatorFuncs {
		if err := validatorFunc(ct, fs); err != nil {
			return err
		}
	}
	return nil
}

func (v *ruleEngineConfigValidator) validateRule(rc *RuleConfig) *RuleEngineError {
	return validateRuleCondition(v, rc.RootCondition)
}

// recursive validation for Rule Condition
func validateRuleCondition(v *ruleEngineConfigValidator, c *Condition) *RuleEngineError {
	validatorFuncs, ok := v.ruleConditionValidators[c.Type]
	if ok {

		for _, validatorFunc := range validatorFuncs {
			if err := validatorFunc(c); err != nil {
				return err
			}
		}

	}

	for _, subCond := range c.SubConditions {
		// recursion
		if err := validateRuleCondition(v, subCond); err != nil {
			return err
		}
	}

	return nil
}

func (v *ruleEngineConfigValidator) validate(config *RuleEngineConfig) *RuleEngineError {
	for _, fieldValidator := range v.fieldValidators {
		if err := fieldValidator(config.Fields); err != nil {
			return err
		}
	}

	for conditionTypeName, conditionType := range config.ConditionTypes {
		if err := v.validateConditionType(conditionType, config.Fields); err != nil {
			err.addMsg(fmt.Sprintf("ConditionType: %v", conditionTypeName))
			return err
		}
	}

	for ruleName, rc := range config.Rules {
		if err := v.validateRule(rc); err != nil {
			err.addMsg(fmt.Sprintf("RuleName: %v", ruleName))
			return err
		}
	}

	return nil
}

var engineConfigValidator = ruleEngineConfigValidator{
	fieldValidators:         []fieldValidatorFunc{},
	condTypeValidators:      make(map[string][]conditionTypeValidatorFunc),
	ruleConditionValidators: make(map[string][]ruleConditionValidatorFunc),
}

func init() {
	engineConfigValidator.addFieldValidator(fieldValueTypeValidator())

	engineConfigValidator.addConditionTypeValidator(EqualOperator,
		operandCountValidator(2),
		operandsWithSameValueTypeValidator(),
		operandValueTypeValidator(Integer, Float, Boolean, String),
		operandValidator(),
	)
	engineConfigValidator.addConditionTypeValidator(NotEqualOperator,
		operandCountValidator(2),
		operandsWithSameValueTypeValidator(),
		operandValueTypeValidator(Integer, Float, Boolean, String),
		operandValidator(),
	)

	engineConfigValidator.addConditionTypeValidator(GreaterOperator,
		operandCountValidator(2),
		operandsWithSameValueTypeValidator(),
		operandValueTypeValidator(Integer, Float),
		operandValidator(),
	)

	engineConfigValidator.addConditionTypeValidator(GreaterEqualOperator,
		operandCountValidator(2),
		operandsWithSameValueTypeValidator(),
		operandValueTypeValidator(Integer, Float),
		operandValidator(),
	)

	engineConfigValidator.addConditionTypeValidator(LessOperator,
		operandCountValidator(2),
		operandsWithSameValueTypeValidator(),
		operandValueTypeValidator(Integer, Float),
		operandValidator(),
	)
	engineConfigValidator.addConditionTypeValidator(LessEqualOperator,
		operandCountValidator(2),
		operandsWithSameValueTypeValidator(),
		operandValueTypeValidator(Integer, Float),
		operandValidator(),
	)

	engineConfigValidator.addConditionTypeValidator(ContainOperator,
		operandCountValidator(2),
		operandsWithSameValueTypeValidator(),
		operandValueTypeValidator(String),
		operandValidator(),
	)

	engineConfigValidator.addRuleConditionValidator(OrCondition,
		minSubConditionCountRuleConditionValidator(2))

	engineConfigValidator.addRuleConditionValidator(AndCondition,
		minSubConditionCountRuleConditionValidator(2))

	engineConfigValidator.addRuleConditionValidator(NegationCondition,
		subConditionCountRuleConditionValidator(1))

}
