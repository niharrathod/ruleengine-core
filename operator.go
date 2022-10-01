package ruleenginecore

import (
	"strings"
)


const (
	firstOperand  = 0
	secondOperand = 1

	valuePrepFail = "Could not get value for given field"
)

func greater[T int64 | float64](input map[string]any, operands []*Operand) bool {
	first, ok := operands[firstOperand].getValue(input).(T)
	if !ok {
		panic(valuePrepFail)
	}
	second, ok := operands[secondOperand].getValue(input).(T)
	if !ok {
		panic(valuePrepFail)
	}

	return first > second
}

func greaterAndEqual[T int64 | float64](input map[string]any, operands []*Operand) bool {
	first, ok := operands[firstOperand].getValue(input).(T)
	if !ok {
		panic(valuePrepFail)
	}
	second, ok := operands[secondOperand].getValue(input).(T)
	if !ok {
		panic(valuePrepFail)
	}
	return first >= second
}

func lesser[T int64 | float64](input map[string]any, operands []*Operand) bool {
	first, ok := operands[firstOperand].getValue(input).(T)
	if !ok {
		panic(valuePrepFail)
	}
	second, ok := operands[secondOperand].getValue(input).(T)
	if !ok {
		panic(valuePrepFail)
	}

	return first < second
}

func lesserAndEqual[T int64 | float64](input map[string]any, operands []*Operand) bool {
	first, ok := operands[firstOperand].getValue(input).(T)
	if !ok {
		panic(valuePrepFail)
	}
	second, ok := operands[secondOperand].getValue(input).(T)
	if !ok {
		panic(valuePrepFail)
	}

	return first <= second
}

func equal[T bool | int64 | float64 | string](input map[string]any, operands []*Operand) bool {
	first, ok := operands[firstOperand].getValue(input).(T)
	if !ok {
		panic(valuePrepFail)
	}
	second, ok := operands[secondOperand].getValue(input).(T)
	if !ok {
		panic(valuePrepFail)
	}

	return first == second
}

func notEqual[T bool | int64 | float64 | string](input map[string]any, operands []*Operand) bool {
	first, ok := operands[firstOperand].getValue(input).(T)
	if !ok {
		panic(valuePrepFail)
	}
	second, ok := operands[secondOperand].getValue(input).(T)
	if !ok {
		panic(valuePrepFail)
	}

	return first != second
}

func contain(input map[string]any, operands []*Operand) bool {
	first, ok := operands[firstOperand].getValue(input).(string)
	if !ok {
		panic(valuePrepFail)
	}
	second, ok := operands[secondOperand].getValue(input).(string)
	if !ok {
		panic(valuePrepFail)
	}
	return strings.Contains(first, second)
}
