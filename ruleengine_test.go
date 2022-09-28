package ruleenginecore

import (
	"reflect"
	"testing"
)

var testOptionComplete *evaluateOption
var testOptionAscendingOne *evaluateOption
var testOptionDescendingOne *evaluateOption
var invalidEvaluateOptions *evaluateOption

func init() {
	testOptionComplete = EvaluateOptions().Complete()
	testOptionAscendingOne = EvaluateOptions().AscendingPriorityBased(1)
	testOptionDescendingOne = EvaluateOptions().DescendingPriorityBased(1)
	invalidEvaluateOptions = EvaluateOptions().DescendingPriorityBased(0)
}

var discount10TestRule = &rule{
	name:     "Discount10",
	priority: 1,
	rootEvaluator: &logicalEvaluator{
		operator: AndOperator,
		innerEvaluators: []Evaluator{
			&greaterEvaluator{
				operandType: IntType,
				operands: []*Operand{
					{
						OperandAs: OperandAsField,
						Val:       "totalAmount",
					},
					{
						OperandAs:  OperandAsConstant,
						Val:        "20000",
						typedValue: int64(20000),
					},
				},
			},
			&equalEvaluator{
				operandType: BoolType,
				operands: []*Operand{
					{
						OperandAs: OperandAsField,
						Val:       "IsHotelBooking",
					},
					{
						OperandAs:  OperandAsConstant,
						Val:        "true",
						typedValue: true,
					},
				},
			},
			&greaterEvaluator{
				operandType: IntType,
				operands: []*Operand{
					{
						OperandAs: OperandAsField,
						Val:       "PaxCount",
					},
					{
						OperandAs:  OperandAsConstant,
						Val:        "5",
						typedValue: int64(5),
					},
				},
			},
		},
	},
	result: map[string]any{
		"discount": 10,
	},
}

var discount5TestRule = &rule{
	name:     "Discount5",
	priority: 2,
	rootEvaluator: &logicalEvaluator{
		operator: AndOperator,
		innerEvaluators: []Evaluator{
			&greaterEvaluator{
				operandType: IntType,
				operands: []*Operand{
					{
						OperandAs: OperandAsField,
						Val:       "totalAmount",
					},
					{
						OperandAs:  OperandAsConstant,
						Val:        "20000",
						typedValue: int64(20000),
					},
				},
			},
			&equalEvaluator{
				operandType: BoolType,
				operands: []*Operand{
					{
						OperandAs: OperandAsField,
						Val:       "IsHotelBooking",
					},
					{
						OperandAs:  OperandAsConstant,
						Val:        "true",
						typedValue: true,
					},
				},
			},
			&logicalEvaluator{
				operator: NegationOperator,
				innerEvaluators: []Evaluator{
					&greaterEvaluator{
						operandType: IntType,
						operands: []*Operand{
							{
								OperandAs: OperandAsField,
								Val:       "PaxCount",
							},
							{
								OperandAs:  OperandAsConstant,
								Val:        "5",
								typedValue: int64(5),
							},
						},
					},
				},
			},
		},
	},
	result: map[string]any{
		"discount": 5,
	},
}

var discount2TestRule = &rule{
	name:     "Discount2",
	priority: 3,
	rootEvaluator: &logicalEvaluator{
		operator: AndOperator,
		innerEvaluators: []Evaluator{
			&greaterEvaluator{
				operandType: IntType,
				operands: []*Operand{
					{
						OperandAs: OperandAsField,
						Val:       "totalAmount",
					},
					{
						OperandAs:  OperandAsConstant,
						Val:        "10000",
						typedValue: int64(10000),
					},
				},
			},
			&equalEvaluator{
				operandType: BoolType,
				operands: []*Operand{
					{
						OperandAs: OperandAsField,
						Val:       "IsHotelBooking",
					},
					{
						OperandAs:  OperandAsConstant,
						Val:        "true",
						typedValue: true,
					},
				},
			},
			&greaterEvaluator{
				operandType: IntType,
				operands: []*Operand{
					{
						OperandAs: OperandAsField,
						Val:       "PaxCount",
					},
					{
						OperandAs:  OperandAsConstant,
						Val:        "2",
						typedValue: int64(2),
					},
				},
			},
		},
	},
	result: map[string]any{
		"discount": 2,
	},
}

func TestCreateRuleEngine(t *testing.T) {
	type args struct {
		engineConfig *RuleEngineConfig
	}
	tests := []struct {
		name    string
		args    args
		want    RuleEngine
		wantErr *RuleEngineError
	}{
		{
			name: "ValidSimpleRuleEngine",
			args: args{
				engineConfig: &RuleEngineConfig{
					Fields: Fields{
						"totalAmount":    IntType,
						"IsHotelBooking": BoolType,
						"PaxCount":       IntType,
					},
					ConditionTypes: map[string]*ConditionType{
						"amountMoreThan20k": {
							Operator:    GreaterOperator,
							OperandType: IntType,
							Operands: []*Operand{
								{
									OperandAs: OperandAsField,
									Val:       "totalAmount",
								},
								{
									OperandAs: OperandAsConstant,
									Val:       "20000",
								},
							},
						},
						"HotelBooking": {
							Operator:    EqualOpearator,
							OperandType: BoolType,
							Operands: []*Operand{
								{
									OperandAs: OperandAsField,
									Val:       "IsHotelBooking",
								},
								{
									OperandAs: OperandAsConstant,
									Val:       "true",
								},
							},
						},
						"PaxCountMoreThan5": {
							Operator:    GreaterOperator,
							OperandType: IntType,
							Operands: []*Operand{
								{
									OperandAs: OperandAsField,
									Val:       "PaxCount",
								},
								{
									OperandAs: OperandAsConstant,
									Val:       "5",
								},
							},
						},
					},
					Rules: map[string]*Rule{
						"Discount10": {
							Priority: 1,
							RootCondition: &Condition{
								ConditionType: AndOperator,
								SubConditions: []*Condition{
									{
										ConditionType: "amountMoreThan20k",
									},
									{
										ConditionType: "HotelBooking",
									},
									{
										ConditionType: "PaxCountMoreThan5",
									},
								},
							},
							Result: map[string]any{
								"discount": 10,
							},
						},
						"Discount5": {
							Priority: 2,
							RootCondition: &Condition{
								ConditionType: AndOperator,
								SubConditions: []*Condition{
									{
										ConditionType: "amountMoreThan20k",
									},
									{
										ConditionType: "HotelBooking",
									},
									{
										ConditionType: NegationOperator,
										SubConditions: []*Condition{
											{
												ConditionType: "PaxCountMoreThan5",
											},
										},
									},
								},
							},
							Result: map[string]any{
								"discount": 5,
							},
						},
					},
				},
			},
			want: &ruleEngine{
				fields: map[string]string{
					"totalAmount":    IntType,
					"IsHotelBooking": BoolType,
					"PaxCount":       IntType,
				},
				ruleMap: map[string]*rule{
					"Discount10": discount10TestRule,
					"Discount5":  discount5TestRule,
				},
				rules: []*rule{
					discount10TestRule,
					discount5TestRule,
				},
			},
			wantErr: nil,
		},
		{
			name: "InvalidFieldTypeInRuleEngineConfig",
			args: args{
				engineConfig: &RuleEngineConfig{
					Fields: Fields{
						"invalid": "asdf",
					},
				},
			},
			want: nil,
			wantErr: &RuleEngineError{
				ComponentName: "",
				ErrCode:       ErrCodeInvalidValueType,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := New(tt.args.engineConfig)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
			if tt.wantErr != nil && gotErr.ErrCode != tt.wantErr.ErrCode {
				t.Errorf("New() got1 = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func Test_ruleEngine_Evaluate(t *testing.T) {
	type fields struct {
		fields  map[string]string
		ruleMap map[string]*rule
		rules   []*rule
	}
	type args struct {
		input Input
		op    *evaluateOption
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*Output
		wantErr *RuleEngineError
	}{
		{
			name: "EvaluateComplete",
			fields: fields{
				fields: map[string]string{
					"totalAmount":    IntType,
					"IsHotelBooking": BoolType,
					"PaxCount":       IntType,
				},
				ruleMap: map[string]*rule{
					"Discount10": discount10TestRule,
					"Discount5":  discount5TestRule,
					"Discount2":  discount2TestRule,
				},
				rules: []*rule{
					discount10TestRule,
					discount5TestRule,
					discount2TestRule,
				},
			},
			args: args{
				input: Input{
					"totalAmount":    "25000",
					"IsHotelBooking": "true",
					"PaxCount":       "10",
				},
				op: testOptionComplete,
			},
			wantErr: nil,
			want: []*Output{
				{
					Rulename: "Discount10",
					Priority: 1,
					Result:   discount10TestRule.result,
				},
				{
					Rulename: "Discount2",
					Priority: 3,
					Result:   discount2TestRule.result,
				},
			},
		},
		{
			name: "AscendingPriorityEvaluate",
			fields: fields{
				fields: map[string]string{
					"totalAmount":    IntType,
					"IsHotelBooking": BoolType,
					"PaxCount":       IntType,
				},
				ruleMap: map[string]*rule{
					"Discount10": discount10TestRule,
					"Discount5":  discount5TestRule,
					"Discount2":  discount2TestRule,
				},
				rules: []*rule{
					discount10TestRule,
					discount5TestRule,
					discount2TestRule,
				},
			},
			args: args{
				input: Input{
					"totalAmount":    "25000",
					"IsHotelBooking": "true",
					"PaxCount":       "10",
				},
				op: testOptionAscendingOne,
			},
			wantErr: nil,
			want: []*Output{
				{
					Rulename: "Discount10",
					Priority: 1,
					Result:   discount10TestRule.result,
				},
			},
		},
		{
			name: "DescendingPriorityEvaluate",
			fields: fields{
				fields: map[string]string{
					"totalAmount":    IntType,
					"IsHotelBooking": BoolType,
					"PaxCount":       IntType,
				},
				ruleMap: map[string]*rule{
					"Discount10": discount10TestRule,
					"Discount5":  discount5TestRule,
					"Discount2":  discount2TestRule,
				},
				rules: []*rule{
					discount10TestRule,
					discount5TestRule,
					discount2TestRule,
				},
			},
			args: args{
				input: Input{
					"totalAmount":    "25000",
					"IsHotelBooking": "true",
					"PaxCount":       "10",
				},
				op: testOptionDescendingOne,
			},
			wantErr: nil,
			want: []*Output{
				{
					Rulename: "Discount2",
					Priority: 3,
					Result:   discount2TestRule.result,
				},
			},
		},
		{
			name: "InvalidFieldTypeInInput",
			fields: fields{
				fields: map[string]string{
					"totalAmount":    IntType,
					"IsHotelBooking": BoolType,
					"PaxCount":       "invalid",
				},
				ruleMap: map[string]*rule{
					"Discount10": discount10TestRule,
					"Discount5":  discount5TestRule,
					"Discount2":  discount2TestRule,
				},
				rules: []*rule{
					discount10TestRule,
					discount5TestRule,
					discount2TestRule,
				},
			},
			args: args{
				input: Input{
					"totalAmount":    "25000",
					"IsHotelBooking": "true",
					"PaxCount":       "10",
				},
				op: testOptionDescendingOne,
			},
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeFailedParsingInput,
			},
			want: nil,
		},
		{
			name: "InvalidEvaluateOption",
			fields: fields{
				fields: map[string]string{
					"totalAmount":    IntType,
					"IsHotelBooking": BoolType,
					"PaxCount":       IntType,
				},
				ruleMap: map[string]*rule{
					"Discount10": discount10TestRule,
					"Discount5":  discount5TestRule,
					"Discount2":  discount2TestRule,
				},
				rules: []*rule{
					discount10TestRule,
					discount5TestRule,
					discount2TestRule,
				},
			},
			args: args{
				input: Input{
					"totalAmount":    "25000",
					"IsHotelBooking": "true",
					"PaxCount":       "10",
				},
				op: invalidEvaluateOptions,
			},
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeInvalidEvaluateOperations,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re := &ruleEngine{
				fields:  tt.fields.fields,
				ruleMap: tt.fields.ruleMap,
				rules:   tt.fields.rules,
			}
			got, gotErr := re.Evaluate(tt.args.input, tt.args.op)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ruleEngine.Evaluate() got = %v, want %v", got, tt.want)
			}
			if tt.wantErr != nil && tt.wantErr.ErrCode != gotErr.ErrCode {
				t.Errorf("ruleEngine.Evaluate() got1 = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func Test_ruleEngine_RulenameBasedEvaluate(t *testing.T) {
	type fields struct {
		fields  map[string]string
		ruleMap map[string]*rule
		rules   []*rule
	}
	type args struct {
		input    Input
		rulename string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Output
		wantErr *RuleEngineError
	}{
		{
			name: "valid",
			fields: fields{
				fields: map[string]string{
					"totalAmount":    IntType,
					"IsHotelBooking": BoolType,
					"PaxCount":       IntType,
				},
				ruleMap: map[string]*rule{
					"Discount10": discount10TestRule,
					"Discount5":  discount5TestRule,
					"Discount2":  discount2TestRule,
				},
				rules: []*rule{
					discount10TestRule,
					discount5TestRule,
					discount2TestRule,
				},
			},
			args: args{
				input: Input{
					"totalAmount":    "25000",
					"IsHotelBooking": "true",
					"PaxCount":       "10",
				},
				rulename: "Discount10",
			},
			wantErr: nil,
			want: &Output{
				Rulename: "Discount10",
				Priority: 1,
				Result:   discount10TestRule.result,
			},
		},
		{
			name: "InvalidInput",
			fields: fields{
				fields: map[string]string{
					"totalAmount":    IntType,
					"IsHotelBooking": BoolType,
					"PaxCount":       IntType,
				},
				ruleMap: map[string]*rule{
					"Discount10": discount10TestRule,
					"Discount5":  discount5TestRule,
					"Discount2":  discount2TestRule,
				},
				rules: []*rule{
					discount10TestRule,
					discount5TestRule,
					discount2TestRule,
				},
			},
			args: args{
				input: Input{
					"totalAmount":    "25000",
					"IsHotelBooking": "true",
					"PaxCount":       "asdf",
				},
				rulename: "Discount10",
			},
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeFailedParsingInput,
			},
			want: nil,
		},
		{
			name: "InvalidRuleNameNotFound",
			fields: fields{
				fields: map[string]string{
					"totalAmount":    IntType,
					"IsHotelBooking": BoolType,
					"PaxCount":       IntType,
				},
				ruleMap: map[string]*rule{
					"Discount10": discount10TestRule,
					"Discount5":  discount5TestRule,
					"Discount2":  discount2TestRule,
				},
				rules: []*rule{
					discount10TestRule,
					discount5TestRule,
					discount2TestRule,
				},
			},
			args: args{
				input: Input{
					"totalAmount":    "25000",
					"IsHotelBooking": "true",
					"PaxCount":       "5",
				},
				rulename: "InvalidRulename",
			},
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeRuleNotFound,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re := &ruleEngine{
				fields:  tt.fields.fields,
				ruleMap: tt.fields.ruleMap,
				rules:   tt.fields.rules,
			}
			got, gotErr := re.EvaluateHavingRulename(tt.args.input, tt.args.rulename)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ruleEngine.RulenameBasedEvaluate() got = %v, want %v", got, tt.want)
			}
			if tt.wantErr != nil && gotErr.ErrCode != tt.wantErr.ErrCode {
				t.Errorf("ruleEngine.RulenameBasedEvaluate() gotErr = %v, wantErr %v", gotErr, tt.wantErr)
			}
		})
	}
}

func Test_rule_evaluate(t *testing.T) {
	type args struct {
		input typedValueMap
	}
	tests := []struct {
		name string
		r    rule
		args args
		want bool
	}{
		{
			name: "validRuleTrueEvaluate",
			r: rule{
				name:          "TestRule1",
				priority:      1,
				rootEvaluator: trueEvaluator,
				result:        map[string]any{},
			},
			args: args{
				typedValueMap{},
			},
			want: true,
		},
		{
			name: "validRuleFalseEvaluate",
			r: rule{
				name:          "TestRule2",
				priority:      1,
				rootEvaluator: falseEvaluator,
				result:        map[string]any{},
			},
			args: args{
				typedValueMap{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &rule{
				name:          tt.r.name,
				priority:      tt.r.priority,
				rootEvaluator: tt.r.rootEvaluator,
				result:        tt.r.result,
			}
			if got := r.evaluate(tt.args.input); got != tt.want {
				t.Errorf("rule.evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}
