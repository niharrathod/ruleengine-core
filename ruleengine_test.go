// rulengine-core is a strictly typed rule engine library, providing a simple interface to create ruleengine and evaluate rule for given input.
package ruleenginecore

import (
	"context"
	"reflect"
	"testing"
)

func Test_rule_evaluate(t *testing.T) {
	type fields struct {
		name          string
		priority      int
		rootEvaluator evaluator
		result        map[string]any
	}
	type args struct {
		ctx   context.Context
		input parsedInput
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  *RuleEngineError
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &rule{
				name:          tt.fields.name,
				priority:      tt.fields.priority,
				rootEvaluator: tt.fields.rootEvaluator,
				result:        tt.fields.result,
			}
			got, got1 := r.evaluate(tt.args.ctx, tt.args.input)
			if got != tt.want {
				t.Errorf("rule.evaluate() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("rule.evaluate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_ruleEngine_validateAndParseInput(t *testing.T) {
	type fields struct {
		fields Fields
	}
	type args struct {
		input Input
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    parsedInput
		wantErr *RuleEngineError
	}{
		{
			name: "valid_InputParse",
			fields: fields{
				fields: Fields{
					"totalAmount": Integer,
				},
			},
			args: args{
				input: Input{
					"totalAmount": "100",
				},
			},
			want: parsedInput{
				"totalAmount": int64(100),
			},
			wantErr: nil,
		},
		{
			name: "invalid_FieldNotFound",
			fields: fields{
				fields: Fields{
					"totalAmount": Integer,
				},
			},
			args: args{
				input: Input{
					"invalid": "100",
				},
			},
			want:    nil,
			wantErr: newError(ErrCodeFieldNotFound),
		},
		{
			name: "invalid_ParsingFailed",
			fields: fields{
				fields: Fields{
					"totalAmount": Integer,
				},
			},
			args: args{
				input: Input{
					"totalAmount": "invalid",
				},
			},
			want:    nil,
			wantErr: newError(ErrCodeParsingFailed),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re := &ruleEngine{
				fields: tt.fields.fields,
			}
			got, gotErr := re.validateAndParseInput(tt.args.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ruleEngine.validateAndParseInput() got = %v, want %v", got, tt.want)
			}
			if tt.wantErr != nil && gotErr.ErrCode != tt.wantErr.ErrCode {
				t.Errorf("ruleEngine.validateAndParseInput() gotErr = %v, wantErr %v", gotErr, tt.wantErr)
			}
		})
	}
}

func Test_ruleEngine_Evaluate(t *testing.T) {
	type fields struct {
		fields  Fields
		ruleMap map[string]*rule
		rules   []*rule
	}
	type args struct {
		ctx   context.Context
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
			name: "valid_EvaluateComplete",
			fields: fields{
				fields: map[string]ValueType{
					"totalAmount":    Integer,
					"IsHotelBooking": Boolean,
					"PaxCount":       Integer,
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
				op:  EvaluateOptions().Complete(),
				ctx: context.TODO(),
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
			name: "valid_AscendingPriorityEvaluate",
			fields: fields{
				fields: map[string]ValueType{
					"totalAmount":    Integer,
					"IsHotelBooking": Boolean,
					"PaxCount":       Integer,
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
				op:  EvaluateOptions().AscendingPriorityBased(1),
				ctx: context.TODO(),
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
			name: "valid_DescendingPriorityEvaluate",
			fields: fields{
				fields: map[string]ValueType{
					"totalAmount":    Integer,
					"IsHotelBooking": Boolean,
					"PaxCount":       Integer,
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
				op:  EvaluateOptions().DescendingPriorityBased(1),
				ctx: context.TODO(),
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
			name: "invalid_InputParsingFailed",
			fields: fields{
				fields: map[string]ValueType{
					"totalAmount":    Integer,
					"IsHotelBooking": Boolean,
					"PaxCount":       Integer,
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
					"totalAmount":    "invalid",
					"IsHotelBooking": "true",
					"PaxCount":       "10",
				},
				op:  EvaluateOptions().DescendingPriorityBased(1),
				ctx: context.TODO(),
			},
			wantErr: newError(ErrCodeParsingFailed),
			want:    nil,
		},
		{
			name: "invalid_AscendingPriorityEvaluate_ContextCancelled",
			fields: fields{
				fields: map[string]ValueType{
					"totalAmount":    Integer,
					"IsHotelBooking": Boolean,
					"PaxCount":       Integer,
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
				op:  EvaluateOptions().AscendingPriorityBased(1),
				ctx: cancelledTestContext,
			},
			wantErr: newError(ErrCodeContextCancelled),
			want:    nil,
		},
		{
			name: "valid_DescendingPriorityEvaluate_ContextCancelled",
			fields: fields{
				fields: map[string]ValueType{
					"totalAmount":    Integer,
					"IsHotelBooking": Boolean,
					"PaxCount":       Integer,
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
				op:  EvaluateOptions().DescendingPriorityBased(1),
				ctx: cancelledTestContext,
			},
			wantErr: newError(ErrCodeContextCancelled),
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re := &ruleEngine{
				fields:  tt.fields.fields,
				ruleMap: tt.fields.ruleMap,
				rules:   tt.fields.rules,
			}
			got, gotErr := re.Evaluate(tt.args.ctx, tt.args.input, tt.args.op)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ruleEngine.Evaluate() got = %v, want %v", got, tt.want)
			}
			if tt.wantErr != nil && gotErr.ErrCode != tt.wantErr.ErrCode {
				t.Errorf("ruleEngine.Evaluate() gotErr = %v, wantErr %v", gotErr, tt.wantErr)
			}
		})
	}
}

func Test_ruleEngine_EvaluateSingleRule(t *testing.T) {
	type fields struct {
		fields  Fields
		ruleMap map[string]*rule
		rules   []*rule
	}
	type args struct {
		ctx      context.Context
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
			name: "valid_SimpleRuleEvaluation",
			fields: fields{
				fields: map[string]ValueType{
					"totalAmount":    Integer,
					"IsHotelBooking": Boolean,
					"PaxCount":       Integer,
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
				ctx:      context.TODO(),
			},
			wantErr: nil,
			want: &Output{
				Rulename: "Discount10",
				Priority: 1,
				Result:   discount10TestRule.result,
			},
		},
		{
			name: "invalid_InputFieldMissing",
			fields: fields{
				fields: map[string]ValueType{
					"totalAmount":    Integer,
					"IsHotelBooking": Boolean,
					"PaxCount":       Integer,
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
					"invalid": "true",
				},
				rulename: "Discount10",
				ctx:      context.TODO(),
			},
			wantErr: newError(ErrCodeFieldNotFound),
			want:    nil,
		},
		{
			name: "invalid_RuleNotFound",
			fields: fields{
				fields: map[string]ValueType{
					"totalAmount":    Integer,
					"IsHotelBooking": Boolean,
					"PaxCount":       Integer,
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
				rulename: "invalid",
				ctx:      context.TODO(),
			},
			wantErr: newError(ErrCodeRuleNotFound),
			want:    nil,
		},
		{
			name: "valid_ContextCancelled",
			fields: fields{
				fields: map[string]ValueType{
					"totalAmount":    Integer,
					"IsHotelBooking": Boolean,
					"PaxCount":       Integer,
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
				ctx:      cancelledTestContext,
			},
			wantErr: newError(ErrCodeContextCancelled),
			want:    nil,
		},
		{
			name: "valid_NoRuleMatched",
			fields: fields{
				fields: map[string]ValueType{
					"totalAmount":    Integer,
					"IsHotelBooking": Boolean,
					"PaxCount":       Integer,
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
			args: args{
				input: Input{
					"totalAmount":    "25000",
					"IsHotelBooking": "true",
					"PaxCount":       "1",
				},
				rulename: "Discount10",
				ctx:      context.TODO(),
			},
			wantErr: nil,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re := &ruleEngine{
				fields:  tt.fields.fields,
				ruleMap: tt.fields.ruleMap,
				rules:   tt.fields.rules,
			}
			got, gotErr := re.EvaluateSingleRule(tt.args.ctx, tt.args.input, tt.args.rulename)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ruleEngine.EvaluateSingleRule() got = %v, want %v", got, tt.want)
			}
			if tt.wantErr != nil && gotErr.ErrCode != tt.wantErr.ErrCode {
				t.Errorf("ruleEngine.EvaluateSingleRule() got1 = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
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
						"totalAmount":    Integer,
						"IsHotelBooking": Boolean,
						"PaxCount":       Integer,
					},
					ConditionTypes: map[string]*ConditionType{
						"amountMoreThan20k": {
							Operator: GreaterOperator,
							Operands: []*Operand{
								{
									Type:      Field,
									ValueType: Integer,
									Val:       "totalAmount",
								},
								{
									Type:      Constant,
									ValueType: Integer,
									Val:       "20000",
								},
							},
						},
						"HotelBooking": {
							Operator: EqualOperator,
							Operands: []*Operand{
								{
									Type:      Field,
									ValueType: Boolean,
									Val:       "IsHotelBooking",
								},
								{
									Type:      Constant,
									ValueType: Boolean,
									Val:       "true",
								},
							},
						},
						"PaxCountMoreThan5": {
							Operator: GreaterOperator,
							Operands: []*Operand{
								{
									Type:      Field,
									ValueType: Integer,
									Val:       "PaxCount",
								},
								{
									Type:      Constant,
									ValueType: Integer,
									Val:       "5",
								},
							},
						},
					},
					Rules: map[string]*RuleConfig{
						"Discount10": {
							Priority: 1,
							RootCondition: &Condition{
								Type: AndCondition,
								SubConditions: []*Condition{
									{
										Type: "amountMoreThan20k",
									},
									{
										Type: "HotelBooking",
									},
									{
										Type: "PaxCountMoreThan5",
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
								Type: AndCondition,
								SubConditions: []*Condition{
									{
										Type: "amountMoreThan20k",
									},
									{
										Type: "HotelBooking",
									},
									{
										Type: NegationCondition,
										SubConditions: []*Condition{
											{
												Type: "PaxCountMoreThan5",
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
				fields: map[string]ValueType{
					"totalAmount":    Integer,
					"IsHotelBooking": Boolean,
					"PaxCount":       Integer,
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
			name: "Invalid_FieldTypeInRuleEngineConfig",
			args: args{
				engineConfig: &RuleEngineConfig{
					Fields: Fields{
						"invalid": unknownValueType,
					},
				},
			},
			want:    nil,
			wantErr: newError(ErrCodeInvalidValueType),
		},
		{
			name: "Invalid_ConditionTypeNotFound",
			args: args{
				engineConfig: &RuleEngineConfig{
					Fields: Fields{
						"totalAmount":    Integer,
						"IsHotelBooking": Boolean,
						"PaxCount":       Integer,
					},
					ConditionTypes: map[string]*ConditionType{},
					Rules: map[string]*RuleConfig{
						"Discount10": {
							Priority: 1,
							RootCondition: &Condition{
								Type: AndCondition,
								SubConditions: []*Condition{
									{
										Type: "amountMoreThan20k",
									},
									{
										Type: "HotelBooking",
									},
									{
										Type: "PaxCountMoreThan5",
									},
								},
							},
							Result: map[string]any{
								"discount": 10,
							},
						},
					},
				},
			},
			want:    nil,
			wantErr: newError(ErrCodeConditionTypeNotFound),
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

var discount10TestRule = &rule{
	name:     "Discount10",
	priority: 1,
	rootEvaluator: &logicalEvaluator{
		operator: AndCondition,
		innerEvaluators: []evaluator{
			&greaterEvaluator{
				operands: []*Operand{
					{
						Type:      Field,
						ValueType: Integer,
						Val:       "totalAmount",
					},
					{
						Type:       Constant,
						ValueType:  Integer,
						Val:        "20000",
						typedValue: int64(20000),
					},
				},
			},
			&equalEvaluator{
				operands: []*Operand{
					{
						Type:      Field,
						ValueType: Boolean,
						Val:       "IsHotelBooking",
					},
					{
						Type:       Constant,
						ValueType:  Boolean,
						Val:        "true",
						typedValue: true,
					},
				},
			},
			&greaterEvaluator{
				operands: []*Operand{
					{
						Type:      Field,
						ValueType: Integer,
						Val:       "PaxCount",
					},
					{
						Type:       Constant,
						ValueType:  Integer,
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
		operator: AndCondition,
		innerEvaluators: []evaluator{
			&greaterEvaluator{
				operands: []*Operand{
					{
						Type:      Field,
						ValueType: Integer,
						Val:       "totalAmount",
					},
					{
						Type:       Constant,
						ValueType:  Integer,
						Val:        "20000",
						typedValue: int64(20000),
					},
				},
			},
			&equalEvaluator{
				operands: []*Operand{
					{
						Type:      Field,
						ValueType: Boolean,
						Val:       "IsHotelBooking",
					},
					{
						Type:       Constant,
						ValueType:  Boolean,
						Val:        "true",
						typedValue: true,
					},
				},
			},
			&logicalEvaluator{
				operator: NegationCondition,
				innerEvaluators: []evaluator{
					&greaterEvaluator{
						operands: []*Operand{
							{
								Type:      Field,
								ValueType: Integer,
								Val:       "PaxCount",
							},
							{
								Type:       Constant,
								ValueType:  Integer,
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
		operator: AndCondition,
		innerEvaluators: []evaluator{
			&greaterEvaluator{
				operands: []*Operand{
					{
						Type:      Field,
						ValueType: Integer,
						Val:       "totalAmount",
					},
					{
						Type:       Constant,
						ValueType:  Integer,
						Val:        "10000",
						typedValue: int64(10000),
					},
				},
			},
			&equalEvaluator{
				operands: []*Operand{
					{
						Type:      Field,
						ValueType: Boolean,
						Val:       "IsHotelBooking",
					},
					{
						Type:       Constant,
						ValueType:  Boolean,
						Val:        "true",
						typedValue: true,
					},
				},
			},
			&greaterEvaluator{
				operands: []*Operand{
					{
						Type:      Field,
						ValueType: Integer,
						Val:       "PaxCount",
					},
					{
						Type:       Constant,
						ValueType:  Integer,
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

var cancelledTestContext context.Context

func init() {
	testContext, canFunc := context.WithCancel(context.Background())
	canFunc()
	cancelledTestContext = testContext
}
