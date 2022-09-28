// rulengine-core is a strickly typed rule engine library, provding a simple interface to create ruleengine and evaluate rule for given input. 
package ruleengine

import "sort"

const (
	IntType    = "int"
	FloatType  = "float"
	BoolType   = "bool"
	StringType = "string"

	OperandAsField    = "field"
	OperandAsConstant = "constant"
)

type rule struct {
	name          string
	priority      int
	rootEvaluator Evaluator
	result        map[string]any
}

// <key> : <value> is <fieldname> : <typedvalue>
type typedValueMap map[string]any

func (r *rule) evaluate(input typedValueMap) bool {
	return r.rootEvaluator.evaluate(input)
}

type ruleEngine struct {
	fields map[string]string

	// map of rulename and rule
	ruleMap map[string]*rule

	// ascending ordered rules
	rules []*rule
}

type RuleEngine interface {
	Evaluate(input Input, op *evaluateOption) ([]*Output, *RuleEngineError)
	EvaluateHavingRulename(input Input, rulename string) (*Output, *RuleEngineError)
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

func prepareEvaluatorTree(cond *Condition, customConditions map[string]*ConditionType) Evaluator {
	switch cType := cond.ConditionType; cType {
	case AndOperator, OrOperator, NegationOperator:
		logicalEval := logicalEvaluator{
			operator:        cType,
			innerEvaluators: []Evaluator{},
		}

		for _, subCondition := range cond.SubConditions {
			logicalEval.innerEvaluators = append(logicalEval.innerEvaluators, prepareEvaluatorTree(subCondition, customConditions))
		}

		return &logicalEval

	default:
		custCondition := customConditions[cType]

		switch op := custCondition.Operator; op {
		case GreaterOperator:
			return &greaterEvaluator{operandType: custCondition.OperandType, operands: custCondition.Operands}
		case GreaterEqualOperator:
			return &greaterEqualEvaluator{operandType: custCondition.OperandType, operands: custCondition.Operands}
		case LessOpeartor:
			return &lessEvaluator{operandType: custCondition.OperandType, operands: custCondition.Operands}
		case LessEqualOpeartor:
			return &lessEqualEvaluator{operandType: custCondition.OperandType, operands: custCondition.Operands}
		case EqualOpearator:
			return &equalEvaluator{operandType: custCondition.OperandType, operands: custCondition.Operands}
		case NotEqualOperator:
			return &notEqualEvaluator{operandType: custCondition.OperandType, operands: custCondition.Operands}
		case ContainOperator:
			return &containEvaluator{operandType: custCondition.OperandType, operands: custCondition.Operands}
		}

	}
	panic("Invalid " + cond.ConditionType + " operator type.")
}

// Evaluate the input based on options
// example,
//
//  1. Evaludate(input, ruleenginecore.EvaluateOptions().Complete())
//     this would evaluate all  rules and returns output as a slice of succeeded rule result.
//
//  2. Evaludate(input, ruleenginecore.EvaluateOptions().AscendingPriorityBased(5))
//     this would evaluate rules in ascending priority order (ex : 1,2,3...) and returns top 5 succeeded rule result
//
//  3. Evaludate(input, ruleenginecore.EvaluateOptions().DescendingPriorityBased(5))
//     this would evaluate rules in descending priority order (ex: 10,9,8...) and returns top 5 succeeded rule result
func (re *ruleEngine) Evaluate(input Input, op *evaluateOption) ([]*Output, *RuleEngineError) {
	inputVals, err := input.Validate(re.fields)
	if err != nil {
		return nil, err
	}

	if op.considerPriority && op.findFirst <= 0 {
		return nil, NewError(ErrCodeInvalidEvaluateOperations, "EvaluateOption", "priority option with 'n' should be greater than 0")
	}

	result := []*Output{}

	if !op.considerPriority {
		for i := 0; i < len(re.rules); i++ {
			rule := re.rules[i]
			if rule.evaluate(inputVals) {
				out := Output{
					Rulename: rule.name,
					Priority: rule.priority,
					Result:   rule.result,
				}

				result = append(result, &out)
			}
		}

		return result, nil

	} else if op.ascPriority {
		for i := 0; i < len(re.rules); i++ {
			rule := re.rules[i]
			if rule.evaluate(inputVals) {
				out := Output{
					Rulename: rule.name,
					Priority: rule.priority,
					Result:   rule.result,
				}

				result = append(result, &out)

				op.findFirst--
				if op.findFirst == 0 {
					break
				}
			}
		}

		return result, nil
	} else {
		for i := len(re.rules) - 1; i >= 0; i-- {
			rule := re.rules[i]
			if rule.evaluate(inputVals) {
				out := Output{
					Rulename: rule.name,
					Priority: rule.priority,
					Result:   rule.result,
				}

				result = append(result, &out)

				op.findFirst--
				if op.findFirst == 0 {
					break
				}
			}
		}
		return result, nil
	}
}

// Evaluate the input.
// evaluates only one rule having given rulename
// returns an error(having error code ErrCodeRuleNotFound) if rule is not found with given rulename.
func (re *ruleEngine) EvaluateHavingRulename(input Input, rulename string) (*Output, *RuleEngineError) {
	inputVals, err := input.Validate(re.fields)
	if err != nil {
		return nil, err
	}

	if rule, ok := re.ruleMap[rulename]; ok {
		if rule.rootEvaluator.evaluate(inputVals) {
			result := Output{
				Rulename: rulename,
				Priority: rule.priority,
				Result:   rule.result,
			}
			return &result, nil
		}
		return nil, nil

	} else {
		return nil, NewError(ErrCodeRuleNotFound, "rulename:"+rulename, "")
	}
}
