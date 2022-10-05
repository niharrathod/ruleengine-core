// rulengine-core is a strickly typed rule engine library, provding a simple interface to create ruleengine and evaluate rule for given input.
package ruleenginecore

import (
	"context"
	"sort"
)

type rule struct {
	name          string
	priority      int
	rootEvaluator evaluator
	result        map[string]any
}

// <key> : <value> is <fieldname> : <typedvalue>
type typedValueMap map[string]any

func (r *rule) evaluate(ctx context.Context, input typedValueMap) (bool, *RuleEngineError) {
	out := make(chan bool)
	ctxCancelled := false

	go func(ctx context.Context, input typedValueMap, resultChan chan<- bool) {
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
		return false, newError(ErrCodeContextCancelled, "ruleEvaluation : "+r.name, "")
	}
	return result, nil
}

type ruleEngine struct {
	fields map[string]string

	// map of rulename and rule
	ruleMap map[string]*rule

	// ascending ordered rules
	rules []*rule
}

type RuleEngine interface {
	// Evaluates the input based on options
	// example,
	//
	//  1. Evaluate(input, ruleenginecore.EvaluateOptions().Complete())
	//     this would evaluate all rules and returns output as a slice of matched rule.
	//
	//  2. Evaluate(input, ruleenginecore.EvaluateOptions().AscendingPriorityBased(5))
	//     this would evaluate rules in ascending priority order (ex : 1,2,3...) and returns top 5 matched rule as output
	//
	//  3. Evaluate(input, ruleenginecore.EvaluateOptions().DescendingPriorityBased(5))
	//     this would evaluate rules in descending priority order (ex: 10,9,8...) and returns top 5 matched rule as output
	Evaluate(ctx context.Context, input Input, op *evaluateOption) ([]*Output, *RuleEngineError)

	// Evaluate the input but only one rule having given 'rulename'
	EvaluateHavingRulename(ctx context.Context, input Input, rulename string) (*Output, *RuleEngineError)
}

// creates new rule engine with given engine configuration
func New(engineConfig *RuleEngineConfig) (RuleEngine, *RuleEngineError) {
	err := engineConfig.Validate()
	if err != nil {
		return nil, err
	}

	engine := ruleEngine{fields: map[string]string{}, ruleMap: map[string]*rule{}, rules: []*rule{}}

	for fieldname, fieldtype := range engineConfig.Fields {
		engine.fields[fieldname] = fieldtype
	}

	for rname, r := range engineConfig.Rules {
		ru := rule{
			name:          rname,
			priority:      r.Priority,
			result:        r.Result,
			rootEvaluator: prepareEvaluatorTree(r.RootCondition, engineConfig.ConditionTypes),
		}
		engine.ruleMap[rname] = &ru
		engine.rules = append(engine.rules, &ru)
	}

	sort.Slice(engine.rules, func(i, j int) bool {
		return engine.rules[i].priority < engine.rules[j].priority
	})

	return &engine, nil
}

func prepareEvaluatorTree(cond *Condition, customConditions map[string]*ConditionType) evaluator {
	switch cType := cond.ConditionType; cType {
	case AndOperator, OrOperator, NegationOperator:
		logicalEval := logicalEvaluator{
			operator:        cType,
			innerEvaluators: []evaluator{},
		}

		for _, subCondition := range cond.SubConditions {
			logicalEval.innerEvaluators = append(logicalEval.innerEvaluators, prepareEvaluatorTree(subCondition, customConditions))
		}

		return &logicalEval

	default:
		custCondition, ok := customConditions[cType]

		if !ok {
			panic("Could not find condition for " + cType + " type.")
		}

		switch op := custCondition.Operator; op {
		case GreaterOperator:
			return &greaterEvaluator{operandType: custCondition.OperandType, operands: custCondition.Operands}
		case GreaterEqualOperator:
			return &greaterEqualEvaluator{operandType: custCondition.OperandType, operands: custCondition.Operands}
		case LessOperator:
			return &lessEvaluator{operandType: custCondition.OperandType, operands: custCondition.Operands}
		case LessEqualOperator:
			return &lessEqualEvaluator{operandType: custCondition.OperandType, operands: custCondition.Operands}
		case EqualOperator:
			return &equalEvaluator{operandType: custCondition.OperandType, operands: custCondition.Operands}
		case NotEqualOperator:
			return &notEqualEvaluator{operandType: custCondition.OperandType, operands: custCondition.Operands}
		case ContainOperator:
			return &containEvaluator{operandType: custCondition.OperandType, operands: custCondition.Operands}
		default:
			panic("Could not find condition for " + cType + " type.")
		}
	}
}

func (re *ruleEngine) Evaluate(ctx context.Context, input Input, op *evaluateOption) ([]*Output, *RuleEngineError) {
	inputVals, err := input.validateAndParseValues(re.fields)
	if err != nil {
		return nil, err
	}

	if op.evalType != Complete && op.limit <= 0 {
		return nil, newError(ErrCodeInvalidEvaluateOperations, "EvaluateOption", "priority option with 'n' should be greater than 0")
	}

	result := []*Output{}

	if op.evalType == Complete {
		for i := 0; i < len(re.rules); i++ {
			rule := re.rules[i]
			matched, err := rule.evaluate(ctx, inputVals)

			if err != nil {
				return nil, err
			}

			if matched {
				out := Output{
					Rulename: rule.name,
					Priority: rule.priority,
					Result:   rule.result,
				}

				result = append(result, &out)
			}
		}

		return result, nil

	} else if op.evalType == AscendingPriorityBased {
		for i := 0; i < len(re.rules); i++ {
			rule := re.rules[i]
			matched, err := rule.evaluate(ctx, inputVals)
			if err != nil {
				return nil, err
			}

			if matched {
				out := Output{
					Rulename: rule.name,
					Priority: rule.priority,
					Result:   rule.result,
				}

				result = append(result, &out)

				op.limit--
				if op.limit == 0 {
					break
				}
			}
		}

		return result, nil
	} else {
		for i := len(re.rules) - 1; i >= 0; i-- {
			rule := re.rules[i]
			matched, err := rule.evaluate(ctx, inputVals)
			if err != nil {
				return nil, err
			}

			if matched {
				out := Output{
					Rulename: rule.name,
					Priority: rule.priority,
					Result:   rule.result,
				}

				result = append(result, &out)

				op.limit--
				if op.limit == 0 {
					break
				}
			}
		}
		return result, nil
	}
}

func (re *ruleEngine) EvaluateHavingRulename(ctx context.Context, input Input, rulename string) (*Output, *RuleEngineError) {
	inputVals, err := input.validateAndParseValues(re.fields)
	if err != nil {
		return nil, err
	}

	if rule, ok := re.ruleMap[rulename]; ok {
		matched, err := rule.evaluate(ctx, inputVals)
		if err != nil {
			return nil, err
		}

		if matched {
			result := Output{
				Rulename: rulename,
				Priority: rule.priority,
				Result:   rule.result,
			}
			return &result, nil
		}
		return nil, nil

	} else {
		return nil, newError(ErrCodeRuleNotFound, "rulename:"+rulename, "")
	}
}
