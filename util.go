package ruleenginecore

import (
	"fmt"
	"strconv"
)

func parseValue(value string, toType ValueType) (any, *RuleEngineError) {
	switch toType {
	case Boolean:
		if val, err := strconv.ParseBool(value); err != nil {
			return nil, newError(ErrCodeParsingFailed, fmt.Sprintf("ParingError: %v", err))
		} else {
			return val, nil
		}
	case Integer:
		if val, err := strconv.ParseInt(value, 10, 64); err != nil {
			return nil, newError(ErrCodeParsingFailed, fmt.Sprintf("ParingError: %v", err))
		} else {
			return val, nil
		}
	case Float:
		if val, err := strconv.ParseFloat(value, 64); err != nil {
			return nil, newError(ErrCodeParsingFailed, fmt.Sprintf("ParingError: %v", err))
		} else {
			return val, nil
		}
	case String:
		return value, nil

	default:
		// no-op
		return nil, newError(ErrCodeInvalidValueType)
	}
}
