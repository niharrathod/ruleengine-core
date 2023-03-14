package ruleenginecore

import "fmt"

type evaluator interface {
	evaluate(input parsedInput) bool
}

type customEvaluator struct {
	operands []*Operand
}

type logicalEvaluator struct {
	operator        string
	innerEvaluators []evaluator
}

func (le *logicalEvaluator) evaluate(input parsedInput) bool {
	switch operator := le.operator; operator {
	case OrCondition:
		for _, evaluator := range le.innerEvaluators {
			if evaluator.evaluate(input) {
				return true
			}
		}
		return false
	case AndCondition:
		for _, evaluator := range le.innerEvaluators {
			if !evaluator.evaluate(input) {
				return false
			}
		}
		return true
	case NegationCondition:
		return !le.innerEvaluators[0].evaluate(input)
	}

	// no-op
	panic("operator:" + le.operator + " is invalid")
}

type greaterEvaluator customEvaluator

func (ge *greaterEvaluator) evaluate(input parsedInput) bool {
	switch ge.operands[0].ValueType {
	case Integer:
		return greater[int64](input, ge.operands)
	case Float:
		return greater[float64](input, ge.operands)
	}
	// no-op
	panic("Invalid operandType for '" + GreaterOperator + "' operator")
}

type greaterEqualEvaluator customEvaluator

func (gte *greaterEqualEvaluator) evaluate(input parsedInput) bool {
	switch gte.operands[0].ValueType {
	case Integer:
		return greaterAndEqual[int64](input, gte.operands)
	case Float:
		return greaterAndEqual[float64](input, gte.operands)
	}

	// no-op
	panic("Invalid operandType for '" + GreaterEqualOperator + "' operator")
}

type lessEvaluator customEvaluator

func (lt *lessEvaluator) evaluate(input parsedInput) bool {

	switch lt.operands[0].ValueType {
	case Integer:
		return lesser[int64](input, lt.operands)
	case Float:
		return lesser[float64](input, lt.operands)
	}

	// no-op
	panic("Invalid operandType for '" + LessOperator + "' operator")
}

type lessEqualEvaluator customEvaluator

func (lte *lessEqualEvaluator) evaluate(input parsedInput) bool {
	switch lte.operands[0].ValueType {
	case Integer:
		return lesserAndEqual[int64](input, lte.operands)
	case Float:
		return lesserAndEqual[float64](input, lte.operands)
	}

	// no-op
	panic("Invalid operandType for '" + LessEqualOperator + "' operator")
}

type equalEvaluator customEvaluator

func (eq *equalEvaluator) evaluate(input parsedInput) bool {
	switch eq.operands[0].ValueType {
	case Integer:
		return equal[int64](input, eq.operands)
	case Float:
		return equal[float64](input, eq.operands)
	case Boolean:
		return equal[bool](input, eq.operands)
	case String:
		return equal[string](input, eq.operands)
	}

	// no-op
	panic("Invalid operandType for '" + EqualOperator + "' operator")
}

type notEqualEvaluator customEvaluator

func (neq *notEqualEvaluator) evaluate(input parsedInput) bool {
	switch neq.operands[0].ValueType {
	case Integer:
		return notEqual[int64](input, neq.operands)
	case Float:
		return notEqual[float64](input, neq.operands)
	case Boolean:
		return notEqual[bool](input, neq.operands)
	case String:
		return notEqual[string](input, neq.operands)
	}

	// no-op
	panic("Invalid operandType for '" + NotEqualOperator + "' operator")
}

type containEvaluator customEvaluator

func (ce *containEvaluator) evaluate(input parsedInput) bool {
	switch ce.operands[0].ValueType {
	case String:
		return contain(input, ce.operands)
	}

	// no-op
	panic("Invalid operandType: for '" + ContainOperator + "' operator")
}

type evaluatorBuilderFunc func(operands []*Operand) evaluator

type evaluatorFactory struct {
	evaluatorBuilders map[string]evaluatorBuilderFunc
}

func (ef *evaluatorFactory) AddEvaluator(operator string, evalBuilderFunc evaluatorBuilderFunc) {
	ef.evaluatorBuilders[operator] = evalBuilderFunc
}

func (ef *evaluatorFactory) build(ct *ConditionType) (evaluator, *RuleEngineError) {
	evalBuilderFunc, ok := ef.evaluatorBuilders[ct.Operator]
	if !ok {
		return nil, newError(ErrCodeInvalidOperator,
			fmt.Sprintf("Operator: %v", ct.Operands))
	}

	return evalBuilderFunc(ct.Operands), nil
}

var evalFactory = evaluatorFactory{
	evaluatorBuilders: make(map[string]evaluatorBuilderFunc),
}

func addNewEvaluator(operator string, evalBuilderFunc evaluatorBuilderFunc) {
	evalFactory.AddEvaluator(operator, evalBuilderFunc)
}

func init() {
	addNewEvaluator(GreaterOperator, func(operands []*Operand) evaluator {
		return &greaterEvaluator{operands: operands}
	})
	addNewEvaluator(GreaterEqualOperator, func(operands []*Operand) evaluator {
		return &greaterEqualEvaluator{operands: operands}
	})
	addNewEvaluator(LessOperator, func(operands []*Operand) evaluator {
		return &lessEvaluator{operands: operands}
	})
	addNewEvaluator(LessEqualOperator, func(operands []*Operand) evaluator {
		return &lessEqualEvaluator{operands: operands}
	})
	addNewEvaluator(EqualOperator, func(operands []*Operand) evaluator {
		return &equalEvaluator{operands: operands}
	})
	addNewEvaluator(NotEqualOperator, func(operands []*Operand) evaluator {
		return &notEqualEvaluator{operands: operands}
	})
	addNewEvaluator(ContainOperator, func(operands []*Operand) evaluator {
		return &containEvaluator{operands: operands}
	})
}

func ruleEvaluatorBuild(rootCondition *Condition, conditionTypes map[string]*ConditionType, fs Fields) (evaluator, *RuleEngineError) {
	switch c := rootCondition.Type; c {
	case AndCondition, OrCondition, NegationCondition:

		logicalEval := logicalEvaluator{
			operator:        c,
			innerEvaluators: []evaluator{},
		}

		for _, subCondition := range rootCondition.SubConditions {
			eval, err := ruleEvaluatorBuild(subCondition, conditionTypes, fs)
			if err != nil {
				return nil, err
			}
			logicalEval.innerEvaluators = append(logicalEval.innerEvaluators, eval)
		}

		return &logicalEval, nil
	}

	ct, ok := conditionTypes[rootCondition.Type]
	if !ok {
		return nil, newError(ErrCodeConditionTypeNotFound,
			fmt.Sprintf("ConditionTypeName: %v", rootCondition.Type))
	}

	return evalFactory.build(ct)
}
