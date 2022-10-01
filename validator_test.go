package ruleenginecore

import (
	"reflect"
	"testing"
)

func TestFields_IsValid(t *testing.T) {
	type args struct {
		fieldname string
		fieldtype string
	}
	tests := []struct {
		name    string
		fs      Fields
		args    args
		wantErr *RuleEngineError
	}{
		{
			name: "ValidNameAndType",
			fs:   Fields{"dog": "int", "cat": "string"},
			args: args{
				fieldname: "dog",
				fieldtype: "int",
			},
			wantErr: nil,
		},
		{
			name: "FieldNotFound",
			fs:   Fields{"dog": "int", "cat": "string"},
			args: args{
				fieldname: "rat",
				fieldtype: "string",
			},
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeFieldNotFound,
			},
		},
		{
			name: "InvalidFieldType",
			fs:   Fields{"dog": "int", "cat": "string"},
			args: args{
				fieldname: "dog",
				fieldtype: "string",
			},
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeInvalidValueType,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := tt.fs.isValid(tt.args.fieldname, tt.args.fieldtype)

			if (tt.wantErr == nil) && (gotErr == nil) {
				return
			}

			if gotErr != nil && tt.wantErr != nil && gotErr.ErrCode == tt.wantErr.ErrCode {
				return
			}

			t.Errorf("Fields.IsValid() gotErr:%v, wantErr:%v ", gotErr, tt.wantErr)
		})
	}
}

func TestFields_validate(t *testing.T) {
	tests := []struct {
		name    string
		fs      Fields
		wantErr *RuleEngineError
	}{
		{
			name:    "ValidNameAndType",
			fs:      Fields{"dog": "int", "cat": "string"},
			wantErr: nil,
		},
		{
			name: "ValidNameAndType",
			fs:   Fields{"dog": "int", "cat": "invalid"},
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeInvalidValueType,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := tt.fs.validate()
			if (tt.wantErr == nil) && (gotErr == nil) {
				return
			}

			if gotErr != nil && tt.wantErr != nil && gotErr.ErrCode == tt.wantErr.ErrCode {
				return
			}

			t.Errorf("fs.validate() gotErr:%v, wantErr:%v ", gotErr, tt.wantErr)
		})
	}
}

func TestConditionType_validateAndParseValues(t *testing.T) {
	type args struct {
		name string
		fs   Fields
	}
	tests := []struct {
		name          string
		conditionType ConditionType
		args          args
		wantErr       *RuleEngineError
	}{
		{
			name: "validConditionGreaterInt",
			conditionType: ConditionType{Operator: GreaterOperator,
				OperandType: IntType,
				Operands:    []*Operand{{OperandAs: OperandAsField, Val: "discount"}, {OperandAs: OperandAsConstant, Val: "10"}}},

			args:    args{name: "DiscountGreaterThan10", fs: Fields{"discount": IntType}},
			wantErr: nil,
		},
		{
			name: "InvalidOperandLength",
			conditionType: ConditionType{Operator: GreaterOperator,
				OperandType: IntType,
				Operands:    []*Operand{{OperandAs: OperandAsField, Val: "discount"}, {OperandAs: OperandAsConstant, Val: "10"}, {OperandAs: OperandAsConstant, Val: "10"}}},

			args: args{name: "DiscountGreaterThan10", fs: Fields{"discount": IntType}},
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeInvalidOperandsLength,
			},
		},
		{
			name: "InvalidOperandType",
			conditionType: ConditionType{Operator: GreaterOperator,
				OperandType: StringType,
				Operands:    []*Operand{{OperandAs: OperandAsField, Val: "discount"}, {OperandAs: OperandAsConstant, Val: "10"}}},

			args: args{name: "DiscountGreaterThan10", fs: Fields{"discount": IntType}},
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeInvalidOperandType,
			},
		},
		{
			name: "ParsingConstantOperandFailed",
			conditionType: ConditionType{Operator: GreaterOperator,
				OperandType: IntType,
				Operands:    []*Operand{{OperandAs: OperandAsField, Val: "discount"}, {OperandAs: OperandAsConstant, Val: "asdf"}}},

			args: args{name: "DiscountGreaterThan10", fs: Fields{"discount": IntType}},
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeFailedParsingInput,
			},
		},
		{
			name: "InvalidOperandAs",
			conditionType: ConditionType{Operator: GreaterOperator,
				OperandType: IntType,
				Operands:    []*Operand{{OperandAs: OperandAsField, Val: "discount"}, {OperandAs: "asdf", Val: "asdf"}}},

			args: args{name: "DiscountGreaterThan10", fs: Fields{"discount": IntType}},
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeInvalidOperandAs,
			},
		},
		{
			name: "InvalidOperator",
			conditionType: ConditionType{Operator: ")",
				OperandType: IntType,
				Operands:    []*Operand{{OperandAs: OperandAsField, Val: "discount"}, {OperandAs: "asdf", Val: "asdf"}}},

			args: args{name: "DiscountGreaterThan10", fs: Fields{"discount": IntType}},
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeInvalidOperator,
			},
		},
		{
			name: "ValidContain",
			conditionType: ConditionType{Operator: ContainOperator,
				OperandType: StringType,
				Operands:    []*Operand{{OperandAs: OperandAsField, Val: "story"}, {OperandAs: OperandAsConstant, Val: "dog"}}},

			args:    args{name: "StoryHasDogCondition", fs: Fields{"story": StringType}},
			wantErr: nil,
		},
		{
			name: "InvalidContainOperandAs",
			conditionType: ConditionType{Operator: ContainOperator,
				OperandType: StringType,
				Operands:    []*Operand{{OperandAs: "Invalid", Val: "story"}, {OperandAs: OperandAsConstant, Val: "dog"}}},

			args: args{name: "StoryHasDogCondition", fs: Fields{"story": StringType}},
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeInvalidOperandAs,
			},
		},
		{
			name: "ValidContain",
			conditionType: ConditionType{Operator: ContainOperator,
				OperandType: StringType,
				Operands:    []*Operand{{OperandAs: OperandAsField, Val: "story"}, {OperandAs: OperandAsConstant, Val: "dog"}}},

			args:    args{name: "StoryHasDogCondition", fs: Fields{"story": StringType}},
			wantErr: nil,
		},
		{
			name: "InvalidContainOperandLength",
			conditionType: ConditionType{Operator: ContainOperator,
				OperandType: StringType,
				Operands:    []*Operand{{OperandAs: OperandAsField, Val: "story"}}},

			args: args{name: "StoryHasDogCondition", fs: Fields{"story": StringType}},
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeInvalidOperandsLength,
			},
		},
		{
			name: "InvalidContainOperandType",
			conditionType: ConditionType{Operator: ContainOperator,
				OperandType: "Invalid",
				Operands:    []*Operand{{OperandAs: OperandAsField, Val: "story"}, {OperandAs: OperandAsField, Val: "story"}}},

			args: args{name: "StoryHasDogCondition", fs: Fields{"story": StringType}},
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeInvalidOperandType,
			},
		},
		{
			name: "InvalidContainOperandType",
			conditionType: ConditionType{Operator: ContainOperator,
				OperandType: StringType,
				Operands:    []*Operand{{OperandAs: OperandAsField, Val: "story"}, {OperandAs: OperandAsConstant, Val: "dog"}}},

			args: args{name: "StoryHasDogCondition", fs: Fields{"story": IntType}},
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeInvalidValueType,
			},
		},
		{
			name: "ValidEqual",
			conditionType: ConditionType{Operator: EqualOperator,
				OperandType: StringType,
				Operands:    []*Operand{{OperandAs: OperandAsField, Val: "story"}, {OperandAs: OperandAsField, Val: "story"}}},

			args:    args{name: "StoryHasDogCondition", fs: Fields{"story": StringType}},
			wantErr: nil,
		},
		{
			name: "InvalidEqualOperandLength",
			conditionType: ConditionType{Operator: EqualOperator,
				OperandType: StringType,
				Operands:    []*Operand{{OperandAs: OperandAsField, Val: "story"}}},

			args: args{name: "StoryHasDogCondition", fs: Fields{"story": StringType}},
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeInvalidOperandsLength,
			},
		},
		{
			name: "InvalidEqualOperandType",
			conditionType: ConditionType{Operator: EqualOperator,
				OperandType: "asdf",
				Operands:    []*Operand{{OperandAs: OperandAsField, Val: "story"}, {OperandAs: OperandAsField, Val: "story"}}},

			args: args{name: "StoryHasDogCondition", fs: Fields{"story": StringType}},
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeInvalidOperandType,
			},
		},
		{
			name: "InvalidEqualOperandValueType",
			conditionType: ConditionType{Operator: EqualOperator,
				OperandType: IntType,
				Operands:    []*Operand{{OperandAs: OperandAsField, Val: "story"}, {OperandAs: OperandAsField, Val: "story"}}},

			args: args{name: "StoryHasDogCondition", fs: Fields{"story": StringType}},
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeInvalidValueType,
			},
		},
		{
			name: "InvalidEqualConstantValue",
			conditionType: ConditionType{Operator: EqualOperator,
				OperandType: IntType,
				Operands:    []*Operand{{OperandAs: OperandAsField, Val: "count"}, {OperandAs: OperandAsConstant, Val: "asdf"}}},

			args: args{name: "CountIs100", fs: Fields{"count": IntType}},
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeFailedParsingInput,
			},
		},
		{
			name: "InvalidEqualOperandAs",
			conditionType: ConditionType{Operator: EqualOperator,
				OperandType: IntType,
				Operands:    []*Operand{{OperandAs: "asdf", Val: "count"}, {OperandAs: OperandAsConstant, Val: "asdf"}}},

			args: args{name: "CountIs100", fs: Fields{"count": IntType}},
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeInvalidOperandAs,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConditionType{
				Operator:    tt.conditionType.Operator,
				OperandType: tt.conditionType.OperandType,
				Operands:    tt.conditionType.Operands,
			}

			gotErr := c.validateAndParseValues(tt.args.name, tt.args.fs)
			if tt.wantErr == nil && gotErr == nil {
				return
			}

			if tt.wantErr != nil && gotErr != nil && gotErr.ErrCode == tt.wantErr.ErrCode {
				return
			}

			for _, op := range c.Operands {
				if op.OperandAs == OperandAsConstant {
					wantVal, _ := getTypedValue(op.Val, c.OperandType)
					if !reflect.DeepEqual(op.typedValue, wantVal) {
						t.Errorf("ConditionType.validate() gotVal:%v, wantVal:%v ", op.typedValue, wantVal)
					}
				}
			}

			t.Errorf("ConditionType.validate() gotErr:%v, wantErr:%v ", gotErr, tt.wantErr)
		})
	}
}

func TestRule_validate(t *testing.T) {

	type args struct {
		name              string
		custConditionType map[string]*ConditionType
	}
	tests := []struct {
		name    string
		rule    Rule
		args    args
		wantErr *RuleEngineError
	}{
		{
			name: "ValidOrRule",
			rule: Rule{
				Priority: 1,
				RootCondition: &Condition{
					ConditionType: OrOperator,
					SubConditions: []*Condition{
						{
							ConditionType: "conditionType1",
						},
						{
							ConditionType: "conditionType2",
						},
					},
				},
				Result: map[string]any{"asdf": "asdf"},
			},
			args: args{name: "OrRule", custConditionType: map[string]*ConditionType{
				"conditionType1": nil,
				"conditionType2": nil,
				"conditionType3": nil,
			}},
			wantErr: nil,
		},
		{
			name: "ValidAndRule",
			rule: Rule{
				Priority: 1,
				RootCondition: &Condition{
					ConditionType: AndOperator,
					SubConditions: []*Condition{
						{
							ConditionType: "conditionType1",
						},
						{
							ConditionType: "conditionType2",
						},
					},
				},
				Result: map[string]any{"asdf": "asdf"},
			},
			args: args{name: "AndRule", custConditionType: map[string]*ConditionType{
				"conditionType1": nil,
				"conditionType2": nil,
				"conditionType3": nil,
			}},
			wantErr: nil,
		},
		{
			name: "InvalidRuleConditionType",
			rule: Rule{
				Priority: 1,
				RootCondition: &Condition{
					ConditionType: AndOperator,
					SubConditions: []*Condition{
						{
							ConditionType: "conditionTypeXX",
						},
						{
							ConditionType: "conditionType2",
						},
					},
				},
				Result: map[string]any{"asdf": "asdf"},
			},
			args: args{name: "AndRule", custConditionType: map[string]*ConditionType{
				"conditionType1": nil,
				"conditionType2": nil,
				"conditionType3": nil,
			}},
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeConditionTypeNotFound,
			},
		},
		{
			name: "InvalidRuleConditionCount",
			rule: Rule{
				Priority: 1,
				RootCondition: &Condition{
					ConditionType: AndOperator,
					SubConditions: []*Condition{
						{
							ConditionType: "conditionTypeXX",
						},
					},
				},
				Result: map[string]any{"asdf": "asdf"},
			},
			args: args{name: "AndRule", custConditionType: map[string]*ConditionType{
				"conditionType1": nil,
				"conditionType2": nil,
				"conditionType3": nil,
			}},
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeInvalidSubConditionCount,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Rule{
				Priority:      tt.rule.Priority,
				RootCondition: tt.rule.RootCondition,
				Result:        tt.rule.Result,
			}
			gotErr := r.validate(tt.args.name, tt.args.custConditionType)
			if (tt.wantErr == nil) && (gotErr == nil) {
				return
			}

			if gotErr != nil && tt.wantErr != nil && gotErr.ErrCode == tt.wantErr.ErrCode {
				return
			}

			t.Errorf("Rule.validate() gotErr:%v, wantErr:%v ", gotErr, tt.wantErr)
		})
	}
}

func TestRuleEngineConfig_validate(t *testing.T) {

	tests := []struct {
		name             string
		ruleEngineConfig RuleEngineConfig
		wantErr          *RuleEngineError
	}{
		{
			name: "ValidRuleEngineConfig",
			ruleEngineConfig: RuleEngineConfig{
				Fields: Fields{
					"discount":    IntType,
					"totalAmount": IntType,
				},
				ConditionTypes: map[string]*ConditionType{
					"IsDiscountGreaterThan10": {
						Operator:    GreaterOperator,
						OperandType: IntType,
						Operands: []*Operand{
							{
								OperandAs: OperandAsField,
								Val:       "discount",
							},
							{
								OperandAs: OperandAsConstant,
								Val:       "10",
							},
						},
					},
					"amountGreaterThan1000": {
						Operator:    GreaterOperator,
						OperandType: IntType,
						Operands: []*Operand{
							{
								OperandAs: OperandAsField,
								Val:       "totalAmount",
							},
							{
								OperandAs: OperandAsConstant,
								Val:       "1000",
							},
						},
					},
				},
				Rules: map[string]*Rule{
					"FirstRule": {
						Priority: 1,
						RootCondition: &Condition{
							ConditionType: AndOperator,
							SubConditions: []*Condition{
								{
									ConditionType: "IsDiscountGreaterThan10",
								},
								{
									ConditionType: "amountGreaterThan1000",
								},
							},
						},
						Result: map[string]any{
							"test": "success",
						},
					},
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RuleEngineConfig{
				Fields:         tt.ruleEngineConfig.Fields,
				ConditionTypes: tt.ruleEngineConfig.ConditionTypes,
				Rules:          tt.ruleEngineConfig.Rules,
			}

			gotErr := c.validate()

			if (tt.wantErr == nil) && (gotErr == nil) {
				return
			}

			if gotErr != nil && tt.wantErr != nil && gotErr.ErrCode == tt.wantErr.ErrCode {
				return
			}

			t.Errorf("RuleEngineConfig.Validate() gotErr:%v, wantErr:%v", gotErr, tt.wantErr)
		})
	}
}

func TestInput_validateAndParseValues(t *testing.T) {
	type args struct {
		fs Fields
	}
	tests := []struct {
		name    string
		input   Input
		args    args
		want    typedValueMap
		wantErr *RuleEngineError
	}{
		{
			name:  "ValidInput",
			args:  args{fs: Fields{"dog": "int", "cat": "float", "rat": "bool", "ant": "string"}},
			input: Input{"dog": "1", "cat": "1.1", "rat": "false", "ant": "sugar"},
			want: typedValueMap{
				"dog": int64(1),
				"cat": float64(1.1),
				"rat": false,
				"ant": "sugar",
			},
			wantErr: nil,
		},
		{
			name:  "MissingMandatoryFieldAsInput",
			args:  args{fs: Fields{"dog": "int", "cat": "float", "rat": "bool", "ant": "string"}},
			input: Input{"dog": "1", "cat": "1.1", "rat": "false"},
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeFieldNotFound,
			},
		},
		{
			name:  "IntParsingFailed",
			args:  args{fs: Fields{"dog": "int", "cat": "float", "rat": "bool", "ant": "string"}},
			input: Input{"dog": "1.1", "cat": "1.1", "rat": "false", "ant": "sugar"},
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeFailedParsingInput,
			},
		},
		{
			name:  "FloatParsingFailed",
			args:  args{fs: Fields{"dog": "int", "cat": "float", "rat": "bool", "ant": "string"}},
			input: Input{"dog": "1", "cat": "1.1asdf", "rat": "false", "ant": "sugar"},
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeFailedParsingInput,
			},
		},
		{
			name:  "BoolParsingFailed",
			args:  args{fs: Fields{"dog": "int", "cat": "float", "rat": "bool", "ant": "string"}},
			input: Input{"dog": "1", "cat": "1.1", "rat": "asdf", "ant": "sugar"},
			wantErr: &RuleEngineError{
				ErrCode: ErrCodeFailedParsingInput,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := tt.input.validateAndParseValues(tt.args.fs)

			if tt.wantErr == nil && gotErr == nil {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Input.Validate() got:%v gotErr:%v, want:%v wantErr:%v", got, gotErr, tt.want, tt.wantErr)
				}
				return
			}

			if tt.wantErr != nil && gotErr != nil && gotErr.ErrCode == tt.wantErr.ErrCode {
				return
			}

			t.Errorf("Input.Validate() got:%v gotErr:%v, want:%v wantErr:%v", got, gotErr, tt.want, tt.wantErr)
		})
	}
}
