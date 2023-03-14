package ruleenginecore

type evaluationType uint

const (
	complete evaluationType = iota + 1
	ascendingPriorityBased
	descendingPriorityBased
)

type evaluateOption struct {
	evalType evaluationType
	limit    int
}

var completeEvalOption = evaluateOption{
	evalType: complete,
}

type evaluateOptionSelector bool

var evaluateOpSelector evaluateOptionSelector

// Evaluates all rules for given input
func (e *evaluateOptionSelector) Complete() *evaluateOption {
	return &completeEvalOption
}

// Evaluates rules in ascending priority(ex. 1,2,3,...) order and considers first n matched rules as outcome
func (e *evaluateOptionSelector) AscendingPriorityBased(n int) *evaluateOption {
	return &evaluateOption{
		evalType: ascendingPriorityBased,
		limit:    n,
	}
}

// Evaluates rules in descending priority(ex. 10,9,8,...) order and considers first n matched rules as outcome
func (e *evaluateOptionSelector) DescendingPriorityBased(n int) *evaluateOption {
	return &evaluateOption{
		evalType: descendingPriorityBased,
		limit:    n,
	}
}

// Evaluation options for rule engine
//
//  1. EvaluateOptions().Complete()
//     evaluates all rules and returns output as a slice of matched rule.
//
//  2. EvaluateOptions().AscendingPriorityBased(5)
//     evaluates rules in ascending priority order (ex : 1,2,3...) and returns top 5 matched rule as output
//
//  3. EvaluateOptions().DescendingPriorityBased(5)
//     evaluate rules in descending priority order (ex: 10,9,8...) and returns top 5 matched rule as output
func EvaluateOptions() *evaluateOptionSelector {
	return &evaluateOpSelector
}
