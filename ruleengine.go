// rulengine-core is a strictly typed rule engine library, providing a simple interface to create ruleengine and evaluate rule for given input.
package ruleenginecore

import (
	"context"
	"fmt"
	"sort"
)

type RuleEngine interface {
	// 'Evaluate' evaluates the input based on options
	Evaluate(ctx context.Context, input Input, op *evaluateOption) ([]*Output, *RuleEngineError)

	// 'EvaluateSingleRule' evaluates the input for one rule having given 'rulename'
	EvaluateSingleRule(ctx context.Context, input Input, rulename string) (*Output, *RuleEngineError)
}

type rule struct {
	name          string
	priority      int
	rootEvaluator evaluator
	result        map[string]any
}

func newRule(ruleName string, r *RuleConfig, customConditionType map[string]*ConditionType, fs Fields) (*rule, *RuleEngineError) {
	rootEvaluator, err := ruleEvaluatorBuild(r.RootCondition, customConditionType, fs)
	if err != nil {
		return nil, err
	}
	ru := &rule{
		name:          ruleName,
		priority:      r.Priority,
		result:        r.Result,
		rootEvaluator: rootEvaluator,
	}
	return ru, nil
}

func (r *rule) evaluate(ctx context.Context, input parsedInput) (bool, *RuleEngineError) {
	out := make(chan bool)
	ctxCancelled := false

	go func(ctx context.Context, input parsedInput, resultChan chan<- bool) {
		select {
		case <-ctx.Done():
			ctxCancelled = true
			resultChan <- false
		default:
			resultChan <- r.rootEvaluator.evaluate(input)
		}
	}(ctx, input, out)

	result := <-out

	if ctxCancelled {
		return false, newError(ErrCodeContextCancelled,
			fmt.Sprintf("Context cancelled while evaluating RuleName: %v", r.name))
	}
	return result, nil
}

type ruleEngine struct {
	fields Fields

	// map of rulename and rule
	ruleMap map[string]*rule

	// ascending ordered rules
	rules []*rule
}

func (re *ruleEngine) validateAndParseInput(input Input) (parsedInput, *RuleEngineError) {
	ret := parsedInput{}
	for fieldname, fieldtype := range re.fields {
		strVal, found := input[fieldname]
		if !found {
			return nil, newError(ErrCodeFieldNotFound,
				fmt.Sprintf("Expecting input with name: %v and valueType: %v", fieldname, fieldtype))
		}

		if val, err := parseValue(strVal, fieldtype); err != nil {
			err.addMsg(fmt.Sprintf("Input parsing failed for field: %v having type %v", fieldname, fieldtype))
			return nil, err
		} else {
			ret[fieldname] = val
		}
	}

	return ret, nil
}

func (re *ruleEngine) Evaluate(ctx context.Context, input Input, op *evaluateOption) ([]*Output, *RuleEngineError) {
	parsedInput, err := re.validateAndParseInput(input)
	if err != nil {
		return nil, err
	}

	if op.evalType == complete {
		return re.ascendingEvaluation(ctx, parsedInput, len(re.rules))
	} else if op.evalType == ascendingPriorityBased {
		return re.ascendingEvaluation(ctx, parsedInput, op.limit)
	} else {
		return re.descendingEvaluation(ctx, parsedInput, op.limit)
	}
}

func (re *ruleEngine) ascendingEvaluation(ctx context.Context, input parsedInput, limit int) ([]*Output, *RuleEngineError) {
	result := []*Output{}
	for i := 0; i < len(re.rules); i++ {
		rule := re.rules[i]
		matched, err := rule.evaluate(ctx, input)
		if err != nil {
			return nil, err
		}

		if matched {
			result = append(result, newOutput(rule.name, rule.priority, rule.result))

			if len(result) == limit {
				break
			}
		}
	}
	return result, nil
}

func (re *ruleEngine) descendingEvaluation(ctx context.Context, input parsedInput, limit int) ([]*Output, *RuleEngineError) {
	result := []*Output{}
	for i := len(re.rules) - 1; i >= 0; i-- {
		rule := re.rules[i]
		matched, err := rule.evaluate(ctx, input)
		if err != nil {
			return nil, err
		}

		if matched {
			result = append(result, newOutput(rule.name, rule.priority, rule.result))

			if len(result) == limit {
				break
			}
		}
	}
	return result, nil
}

func (re *ruleEngine) EvaluateSingleRule(ctx context.Context, input Input, rulename string) (*Output, *RuleEngineError) {
	parsedInput, err := re.validateAndParseInput(input)
	if err != nil {
		return nil, err
	}

	rule, ok := re.ruleMap[rulename]
	if !ok {
		return nil, newError(ErrCodeRuleNotFound, fmt.Sprintf("RuleName: %v", rulename))
	}

	matched, err := rule.evaluate(ctx, parsedInput)
	if err != nil {
		return nil, err
	}

	if matched {
		return newOutput(rulename, rule.priority, rule.result), nil
	}
	return nil, nil
}

// creates new rule engine using provided configuration
func New(engineConfig *RuleEngineConfig) (RuleEngine, *RuleEngineError) {

	if err := engineConfigValidator.validate(engineConfig); err != nil {
		return nil, err
	}

	engine := ruleEngine{
		fields:  engineConfig.Fields,
		ruleMap: map[string]*rule{},
		rules:   []*rule{},
	}

	for ruleName, r := range engineConfig.Rules {
		ru, err := newRule(ruleName, r, engineConfig.ConditionTypes, engine.fields)
		if err != nil {
			return nil, err
		}
		engine.ruleMap[ruleName] = ru
		engine.rules = append(engine.rules, ru)
	}

	sort.Slice(engine.rules, func(i, j int) bool {
		return engine.rules[i].priority < engine.rules[j].priority
	})

	return &engine, nil
}
