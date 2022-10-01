package ruleenginecore

import (
	"context"
	"reflect"
	"testing"
)

var discount10TestRule = &rule{
	name:     "Discount10",
	priority: 1,
	rootEvaluator: &logicalEvaluator{
		operator: AndOperator,
		innerEvaluators: []evaluator{
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
		innerEvaluators: []evaluator{
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
				innerEvaluators: []evaluator{
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
		innerEvaluators: []evaluator{
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
							Operator:    EqualOperator,
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

			if tt.wantErr == nil && gotErr == nil && reflect.DeepEqual(got, tt.want) {
				return
			}

			if tt.wantErr != nil && gotErr != nil && gotErr.ErrCode == tt.wantErr.ErrCode {
				return
			}

			t.Errorf("ruleenginecore.New() got:%v gotErr:%v, want:%v wantErr:%v", got, gotErr, tt.want, tt.wantErr)
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
				op: EvaluateOptions().Complete(),
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
				op: EvaluateOptions().AscendingPriorityBased(1),
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
				op: EvaluateOptions().DescendingPriorityBased(1),
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
				op: EvaluateOptions().DescendingPriorityBased(1),
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
				op: EvaluateOptions().AscendingPriorityBased(0),
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
			got, gotErr := re.Evaluate(context.TODO(), tt.args.input, tt.args.op)

			if tt.wantErr == nil && gotErr == nil && reflect.DeepEqual(got, tt.want) {
				return
			}

			if tt.wantErr != nil && gotErr != nil && gotErr.ErrCode == tt.wantErr.ErrCode {
				return
			}

			t.Errorf("ruleEngine.Evaluate() got:%v gotErr:%v, want:%v wantErr:%v", got, gotErr, tt.want, tt.wantErr)
		})
	}
}

func Test_ruleEngine_Evaluate_ContextAware(t *testing.T) {
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
		name                  string
		fields                fields
		args                  args
		immediateCtxCancelled bool
		want                  []*Output
		wantErr               *RuleEngineError
	}{
		{
			name: "EvaluateComplete_Normal",
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
				op: EvaluateOptions().Complete(),
			},
			immediateCtxCancelled: false,
			wantErr:               nil,
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
			name: "EvaluateComplete_CtxCancelled",
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
				op: EvaluateOptions().Complete(),
			},
			immediateCtxCancelled: true,
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeContextCancelled,
			},
			want: nil,
		},
		{
			name: "AscendingPriorityEvaluate_Normal",
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
				op: EvaluateOptions().AscendingPriorityBased(1),
			},
			immediateCtxCancelled: false,
			wantErr:               nil,
			want: []*Output{
				{
					Rulename: "Discount10",
					Priority: 1,
					Result:   discount10TestRule.result,
				},
			},
		},
		{
			name: "AscendingPriorityEvaluate_ctxCancelled",
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
				op: EvaluateOptions().AscendingPriorityBased(1),
			},
			immediateCtxCancelled: true,
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeContextCancelled,
			},
			want: nil,
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
				op: EvaluateOptions().DescendingPriorityBased(1),
			},
			immediateCtxCancelled: false,
			wantErr:               nil,
			want: []*Output{
				{
					Rulename: "Discount2",
					Priority: 3,
					Result:   discount2TestRule.result,
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
				op: EvaluateOptions().DescendingPriorityBased(1),
			},
			immediateCtxCancelled: true,
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeContextCancelled,
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
			ctx, cancel := context.WithCancel(context.Background())
			if tt.immediateCtxCancelled {
				cancel()
			} else {
				defer cancel()
			}
			got, gotErr := re.Evaluate(ctx, tt.args.input, tt.args.op)

			if tt.wantErr == nil && gotErr == nil && reflect.DeepEqual(got, tt.want) {
				return
			}
			if tt.wantErr != nil && gotErr != nil && tt.wantErr.ErrCode == gotErr.ErrCode {
				return
			}

			t.Errorf("ruleEngine.Evaluate() got:%v gotErr:%v, want:%v wantErr:%v", got, gotErr, tt.want, tt.wantErr)
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
			name: "noMatch",
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
					"IsHotelBooking": "false",
					"PaxCount":       "10",
				},
				rulename: "Discount10",
			},
			wantErr: nil,
			want:    nil,
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
			got, gotErr := re.EvaluateHavingRulename(context.TODO(), tt.args.input, tt.args.rulename)

			if tt.wantErr == nil && gotErr == nil && reflect.DeepEqual(got, tt.want) {
				return
			}

			if tt.wantErr != nil && gotErr != nil && gotErr.ErrCode == tt.wantErr.ErrCode {
				return
			}
			t.Errorf("ruleEngine.RulenameBasedEvaluate() got:%v gotErr:%v, want:%v wantErr:%v", got, gotErr, tt.want, tt.wantErr)
		})
	}
}

func Test_ruleEngine_RulenameBasedEvaluate_ContextAware(t *testing.T) {
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
		name                  string
		fields                fields
		args                  args
		immediateCtxCancelled bool
		want                  *Output
		wantErr               *RuleEngineError
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
			immediateCtxCancelled: false,
			wantErr:               nil,
			want: &Output{
				Rulename: "Discount10",
				Priority: 1,
				Result:   discount10TestRule.result,
			},
		},
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
			immediateCtxCancelled: true,
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeContextCancelled,
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
			ctx, cancel := context.WithCancel(context.Background())
			if tt.immediateCtxCancelled {
				cancel()
			} else {
				defer cancel()
			}

			got, gotErr := re.EvaluateHavingRulename(ctx, tt.args.input, tt.args.rulename)

			if tt.wantErr == nil && gotErr == nil && reflect.DeepEqual(got, tt.want) {
				return
			}
			if tt.wantErr != nil && gotErr != nil && tt.wantErr.ErrCode == gotErr.ErrCode {
				return
			}

			t.Errorf("ruleEngine.RulenameBasedEvaluate() got:%v gotErr:%v, want:%v wantErr:%v", got, gotErr, tt.want, tt.wantErr)
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
			if got, _ := r.evaluate(context.TODO(), tt.args.input); got != tt.want {
				t.Errorf("rule.evaluate() got:%v, want:%v", got, tt.want)
			}
		})
	}
}

func Test_rule_evaluateContextAware(t *testing.T) {
	type args struct {
		input typedValueMap
	}
	tests := []struct {
		name                  string
		r                     rule
		args                  args
		immediateCtxCancelled bool
		want                  bool
		wantErr               *RuleEngineError
	}{
		{
			name: "contextCancelled",
			r: rule{
				name:          "TestRule1",
				priority:      1,
				rootEvaluator: trueEvaluator,
				result:        map[string]any{},
			},
			args: args{
				input: typedValueMap{},
			},
			immediateCtxCancelled: true,
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeContextCancelled,
			},
		},
		{
			name: "doNotCancelContext",
			r: rule{
				name:          "TestRule1",
				priority:      1,
				rootEvaluator: trueEvaluator,
				result:        map[string]any{},
			},
			args: args{
				input: typedValueMap{},
			},
			immediateCtxCancelled: false,
			want:                  true,
			wantErr:               nil,
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

			ctx, cancel := context.WithCancel(context.Background())
			if tt.immediateCtxCancelled {
				cancel()
			} else {
				defer cancel()
			}

			got, gotErr := r.evaluate(ctx, tt.args.input)

			if tt.wantErr == nil && got == tt.want && gotErr == nil {
				return
			}

			if tt.wantErr != nil && gotErr != nil && gotErr.ErrCode == tt.wantErr.ErrCode {
				return
			}

			t.Errorf("rule.evaluate() got:%v gotErr:%v, want:%v wantErr:%v", got, gotErr, tt.want, tt.wantErr)
		})
	}
}

func Test_prepareEvaluatorTree(t *testing.T) {
	type args struct {
		cond             *Condition
		customConditions map[string]*ConditionType
	}
	tests := []struct {
		name      string
		args      args
		want      evaluator
		wantPanic bool
	}{
		{
			name: "validGreaterEqual",
			args: args{
				customConditions: map[string]*ConditionType{
					"GTE": {
						Operator:    GreaterEqualOperator,
						OperandType: IntType,
					},
				},
				cond: &Condition{
					ConditionType: "GTE",
				},
			},
			want: &greaterEqualEvaluator{
				operandType: IntType,
			},
			wantPanic: false,
		},
		{
			name: "validLess",
			args: args{
				customConditions: map[string]*ConditionType{
					"LT": {
						Operator:    LessOperator,
						OperandType: IntType,
					},
				},
				cond: &Condition{
					ConditionType: "LT",
				},
			},
			want: &lessEvaluator{
				operandType: IntType,
			},
			wantPanic: false,
		},
		{
			name: "validLessEqual",
			args: args{
				customConditions: map[string]*ConditionType{
					"LTE": {
						Operator:    LessEqualOperator,
						OperandType: IntType,
					},
				},
				cond: &Condition{
					ConditionType: "LTE",
				},
			},
			want: &lessEqualEvaluator{
				operandType: IntType,
			},
			wantPanic: false,
		},
		{
			name: "validNotEqual",
			args: args{
				customConditions: map[string]*ConditionType{
					"NEQ": {
						Operator:    NotEqualOperator,
						OperandType: IntType,
					},
				},
				cond: &Condition{
					ConditionType: "NEQ",
				},
			},
			want: &notEqualEvaluator{
				operandType: IntType,
			},
			wantPanic: false,
		},
		{
			name: "validContain",
			args: args{
				customConditions: map[string]*ConditionType{
					"CONT": {
						Operator:    ContainOperator,
						OperandType: StringType,
					},
				},
				cond: &Condition{
					ConditionType: "CONT",
				},
			},
			want: &containEvaluator{
				operandType: StringType,
			},
			wantPanic: false,
		},
		{
			name: "InvalidConditionType",
			args: args{
				customConditions: map[string]*ConditionType{
					"CONT": {
						Operator:    ContainOperator,
						OperandType: StringType,
					},
				},
				cond: &Condition{
					ConditionType: "Invalid",
				},
			},
			want:      nil,
			wantPanic: true,
		},
		{
			name: "InvalidCustomConditionOperator",
			args: args{
				customConditions: map[string]*ConditionType{
					"Test": {
						Operator:    "Invalid",
						OperandType: StringType,
					},
				},
				cond: &Condition{
					ConditionType: "Test",
				},
			},
			want:      nil,
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if !tt.wantPanic && r == nil {
					return
				}
				if tt.wantPanic && r != nil {
					return
				}

				t.Errorf("prepareEvaluatorTree() gotPanic:%v , want:%T wantPanic:%v", r != nil, tt.want, tt.wantPanic)
			}()

			got := prepareEvaluatorTree(tt.args.cond, tt.args.customConditions)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("prepareEvaluatorTree() got:%T, want:%T wantPanic:%v", got, tt.want, tt.wantPanic)
			}
		})
	}
}
