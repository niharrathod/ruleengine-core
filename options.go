package ruleenginecore

type evaluateOption struct {
	considerPriority bool
	ascPriority      bool
	findFirst        uint
}

var complete = evaluateOption{
	considerPriority: false,
}

type evaluateOptionSelector bool

var evaluateOpSelector evaluateOptionSelector

// all rules are evaluated without any rule priority consideration
func (e *evaluateOptionSelector) Complete() *evaluateOption {
	return &complete
}

// ascending priority based first n matched rule outcome
func (e *evaluateOptionSelector) AscendingPriorityBased(n uint) *evaluateOption {
	return &evaluateOption{
		considerPriority: true,
		ascPriority:      true,
		findFirst:        n,
	}
}

// descending priority based first n matched rule outcome
func (e *evaluateOptionSelector) DescendingPriorityBased(n uint) *evaluateOption {
	return &evaluateOption{
		considerPriority: true,
		ascPriority:      false,
		findFirst:        n,
	}
}

// options for rule engine evaluate operation
func EvaluateOptions() *evaluateOptionSelector {
	return &evaluateOpSelector
}
