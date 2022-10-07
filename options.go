package ruleenginecore

type evaluationType uint

const (
	Complete evaluationType = iota + 1
	AscendingPriorityBased
	DescendingPriorityBased
)

type EvaluateOption struct {
	evalType evaluationType
	limit    uint
}

var complete = EvaluateOption{
	evalType: Complete,
}

type evaluateOptionSelector bool

var evaluateOpSelector evaluateOptionSelector

// all rules are evaluated without any rule priority consideration
func (e *evaluateOptionSelector) Complete() *EvaluateOption {
	return &complete
}

// ascending priority based first n matched rules as outcome
func (e *evaluateOptionSelector) AscendingPriorityBased(n uint) *EvaluateOption {
	return &EvaluateOption{
		evalType: AscendingPriorityBased,
		limit:    n,
	}
}

// descending priority based first n matched rules as outcome
func (e *evaluateOptionSelector) DescendingPriorityBased(n uint) *EvaluateOption {
	return &EvaluateOption{
		evalType: DescendingPriorityBased,
		limit:    n,
	}
}

// options for rule engine evaluate operation
func EvaluateOptions() *evaluateOptionSelector {
	return &evaluateOpSelector
}
