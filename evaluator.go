package ruleenginecore

type Evaluator interface {
	evaluate(input typedValueMap) bool
}

type customEvaluator struct {
	operandType string
	operands    []*Operand
}

type logicalEvaluator struct {
	operator        string
	innerEvaluators []Evaluator
}

func (le *logicalEvaluator) evaluate(input typedValueMap) bool {
	switch operator := le.operator; operator {
	case OrOperator:
		for _, evaluator := range le.innerEvaluators {
			if evaluator.evaluate(input) {
				return true
			}
		}
		return false
	case AndOperator:
		for _, evaluator := range le.innerEvaluators {
			if !evaluator.evaluate(input) {
				return false
			}
		}
		return true
	case NegationOperator:
		return !le.innerEvaluators[0].evaluate(input)
	}

	panic("operator:" + le.operator + " is invalid")
}

type greaterEvaluator customEvaluator

func (ge *greaterEvaluator) evaluate(input typedValueMap) bool {
	if ge.operandType == IntType {
		return greater[int64](input, ge.operands)
	}
	if ge.operandType == FloatType {
		return greater[float64](input, ge.operands)
	}
	panic("Invalid operandType:" + ge.operandType + " for '" + GreaterOperator + "' operator")
}

type greaterEqualEvaluator customEvaluator

func (gte *greaterEqualEvaluator) evaluate(input typedValueMap) bool {
	if gte.operandType == IntType {
		return greaterAndEqual[int64](input, gte.operands)
	}
	if gte.operandType == FloatType {
		return greaterAndEqual[float64](input, gte.operands)
	}
	panic("Invalid operandType:" + gte.operandType + " for '" + GreaterEqualOperator + "' operator")
}

type lessEvaluator customEvaluator

func (lt *lessEvaluator) evaluate(input typedValueMap) bool {
	if lt.operandType == IntType {
		return lesser[int64](input, lt.operands)
	}
	if lt.operandType == FloatType {
		return lesser[float64](input, lt.operands)
	}
	panic("Invalid operandType:" + lt.operandType + " for '" + LessOperator + "' operator")
}

type lessEqualEvaluator customEvaluator

func (lte *lessEqualEvaluator) evaluate(input typedValueMap) bool {
	if lte.operandType == IntType {
		return lesserAndEqual[int64](input, lte.operands)
	}
	if lte.operandType == FloatType {
		return lesserAndEqual[float64](input, lte.operands)
	}
	panic("Invalid operandType:" + lte.operandType + " for '" + LessEqualOperator + "' operator")
}

type equalEvaluator customEvaluator

func (eq *equalEvaluator) evaluate(input typedValueMap) bool {
	if eq.operandType == IntType {
		return equal[int64](input, eq.operands)
	}
	if eq.operandType == FloatType {
		return equal[float64](input, eq.operands)
	}
	if eq.operandType == BoolType {
		return equal[bool](input, eq.operands)
	}
	if eq.operandType == StringType {
		return equal[string](input, eq.operands)
	}
	panic("Invalid operandType:" + eq.operandType + " for '" + EqualOperator + "' operator")
}

type notEqualEvaluator customEvaluator

func (neq *notEqualEvaluator) evaluate(input typedValueMap) bool {
	if neq.operandType == IntType {
		return notEqual[int64](input, neq.operands)
	}
	if neq.operandType == FloatType {
		return notEqual[float64](input, neq.operands)
	}
	if neq.operandType == BoolType {
		return notEqual[bool](input, neq.operands)
	}
	if neq.operandType == StringType {
		return notEqual[string](input, neq.operands)
	}
	panic("Invalid operandType:" + neq.operandType + " for '" + NotEqualOperator + "' operator")
}

type containEvaluator customEvaluator

func (ce *containEvaluator) evaluate(input typedValueMap) bool {
	if ce.operandType == StringType {
		return contain(input, ce.operands)
	}
	panic("Invalid operandType:" + ce.operandType + " for '" + ContainOperator + "' operator")
}
