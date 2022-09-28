package ruleengine

import (
	"errors"
	"strconv"
)

func getTypedValue(value, toType string) (any, error) {
	switch toType {
	case BoolType:
		if val, err := strconv.ParseBool(value); err != nil {
			return nil, errors.New("Could not convert " + value + " to bool value")
		} else {
			return val, nil
		}
	case IntType:
		if val, err := strconv.ParseInt(value, 10, 64); err != nil {
			return nil, errors.New("Could not convert " + value + " to int value")
		} else {
			return val, nil
		}
	case FloatType:
		if val, err := strconv.ParseFloat(value, 64); err != nil {
			return nil, errors.New("Could not convert " + value + " to float value")
		} else {
			return val, nil
		}
	case StringType:
		return value, nil

	default:
		return nil, errors.New("Invalid value type:" + toType)
	}
}
