package ruleenginecore

import (
	"testing"
)

func isErrorEqual(err1, err2 *RuleEngineError) bool {
	if err1 == nil && err2 == nil {
		return true
	}

	if err1 != nil && err2 != nil && err1.ErrCode == err2.ErrCode {
		return true
	}

	return false
}

func Test_fieldValueTypeValidator(t *testing.T) {
	type args struct {
		fs Fields
	}
	tests := []struct {
		name          string
		args          args
		validatorFunc fieldValidatorFunc
		wantErr       *RuleEngineError
	}{
		{
			name: "valid",
			args: args{
				fs: Fields{
					"testInteger": Integer,
					"testFloat":   Float,
					"testString":  String,
					"testBoolean": Boolean,
				},
			},
			validatorFunc: fieldValueTypeValidator(),
			wantErr:       nil,
		},
		{
			name: "invalid",
			args: args{
				fs: Fields{
					"invalid": unknownValueType,
				},
			},
			validatorFunc: fieldValueTypeValidator(),
			wantErr:       newError(ErrCodeInvalidValueType),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if gotErr := tt.validatorFunc(tt.args.fs); !isErrorEqual(gotErr, tt.wantErr) {
				t.Errorf("fieldValueTypeValidator() = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func Test_operandCountValidator(t *testing.T) {
	type args struct {
		fs Fields
		ct *ConditionType
	}
	tests := []struct {
		name          string
		args          args
		validatorFunc conditionTypeValidatorFunc
		wantErr       *RuleEngineError
	}{
		{
			name: "valid",
			args: args{
				ct: &ConditionType{
					Operator: EqualOperator,
					Operands: []*Operand{
						{
							Type: Field,
						},
						{
							Type: Constant,
						},
					},
				},
			},
			validatorFunc: operandCountValidator(2),
			wantErr:       nil,
		},
		{
			name: "invalid",
			args: args{
				ct: &ConditionType{
					Operator: EqualOperator,
					Operands: []*Operand{
						{
							Type: Field,
						},
						{
							Type: Constant,
						},
					},
				},
			},
			validatorFunc: operandCountValidator(1),
			wantErr:       newError(ErrCodeInvalidOperandsLength),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if gotErr := tt.validatorFunc(tt.args.ct, tt.args.fs); !isErrorEqual(gotErr, tt.wantErr) {
				t.Errorf("operandCountValidator() = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func Test_operandsWithSameValueTypeValidator(t *testing.T) {
	type args struct {
		fs Fields
		ct *ConditionType
	}
	tests := []struct {
		name          string
		args          args
		validatorFunc conditionTypeValidatorFunc
		wantErr       *RuleEngineError
	}{
		{
			name: "valid",
			args: args{
				ct: &ConditionType{
					Operator: EqualOperator,
					Operands: []*Operand{
						{
							Type:      Field,
							ValueType: Integer,
						},
						{
							Type:      Constant,
							ValueType: Integer,
						},
					},
				},
			},
			validatorFunc: operandsWithSameValueTypeValidator(),
			wantErr:       nil,
		},
		{
			name: "invalid",
			args: args{
				ct: &ConditionType{
					Operator: EqualOperator,
					Operands: []*Operand{
						{
							Type:      Field,
							ValueType: Integer,
						},
						{
							Type:      Constant,
							ValueType: Float,
						},
					},
				},
			},
			validatorFunc: operandsWithSameValueTypeValidator(),
			wantErr:       newError(ErrCodeInvalidOperand),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if gotErr := tt.validatorFunc(tt.args.ct, tt.args.fs); !isErrorEqual(gotErr, tt.wantErr) {
				t.Errorf("operandsWithSameValueTypeValidator() = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func Test_operandValueTypeValidator(t *testing.T) {
	type args struct {
		fs Fields
		ct *ConditionType
	}
	tests := []struct {
		name          string
		args          args
		validatorFunc conditionTypeValidatorFunc
		wantErr       *RuleEngineError
	}{
		{
			name: "valid",
			args: args{
				ct: &ConditionType{
					Operator: EqualOperator,
					Operands: []*Operand{
						{
							Type:      Field,
							ValueType: Integer,
						},
						{
							Type:      Constant,
							ValueType: Float,
						},
					},
				},
			},
			validatorFunc: operandValueTypeValidator(Integer, Float),
			wantErr:       nil,
		},
		{
			name: "invalid",
			args: args{
				ct: &ConditionType{
					Operator: EqualOperator,
					Operands: []*Operand{
						{
							Type:      Field,
							ValueType: String,
						},
						{
							Type:      Constant,
							ValueType: Float,
						},
					},
				},
			},
			validatorFunc: operandValueTypeValidator(Integer, Float),
			wantErr:       newError(ErrCodeInvalidOperand),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if gotErr := tt.validatorFunc(tt.args.ct, tt.args.fs); !isErrorEqual(gotErr, tt.wantErr) {
				t.Errorf("operandValueTypeValidator() = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func Test_operandValidator(t *testing.T) {
	type args struct {
		fs Fields
		ct *ConditionType
	}
	tests := []struct {
		name          string
		args          args
		validatorFunc conditionTypeValidatorFunc
		wantErr       *RuleEngineError
	}{
		{
			name: "valid",
			args: args{
				fs: Fields{
					"testField": Integer,
				},
				ct: &ConditionType{
					Operator: EqualOperator,
					Operands: []*Operand{
						{
							Type:      Field,
							ValueType: Integer,
							Val:       "testField",
						},
						{
							Type:      Constant,
							ValueType: Integer,
							Val:       "10",
						},
					},
				},
			},
			validatorFunc: operandValidator(),
			wantErr:       nil,
		},
		{
			name: "invalid_wrongOperandType",
			args: args{
				fs: Fields{
					"testField": Integer,
				},
				ct: &ConditionType{
					Operands: []*Operand{
						{
							Type: unknownOperandType,
						},
					},
				},
			},
			validatorFunc: operandValidator(),
			wantErr:       newError(ErrCodeInvalidOperandType),
		},
		{
			name: "invalid_wrongOperandValueType",
			args: args{
				fs: Fields{
					"testField": Integer,
				},
				ct: &ConditionType{
					Operands: []*Operand{
						{
							Type:      Field,
							ValueType: unknownValueType,
						},
					},
				},
			},
			validatorFunc: operandValidator(),
			wantErr:       newError(ErrCodeInvalidValueType),
		},
		{
			name: "invalid_FieldNotFound",
			args: args{
				ct: &ConditionType{
					Operands: []*Operand{
						{
							Type:      Field,
							ValueType: Integer,
							Val:       "testField",
						},
					},
				},
			},
			validatorFunc: operandValidator(),
			wantErr:       newError(ErrCodeFieldNotFound),
		},
		{
			name: "invalid_ConstantValueParsingFailed",
			args: args{
				ct: &ConditionType{
					Operands: []*Operand{
						{
							Type:      Constant,
							ValueType: Integer,
							Val:       "invalidInteger",
						},
					},
				},
			},
			validatorFunc: operandValidator(),
			wantErr:       newError(ErrCodeParsingFailed),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotErr := tt.validatorFunc(tt.args.ct, tt.args.fs); !isErrorEqual(gotErr, tt.wantErr) {
				t.Errorf("operandValidator() = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func Test_RuleValidator_SubConditionCount(t *testing.T) {
	type args struct {
		c *Condition
	}
	tests := []struct {
		name          string
		args          args
		validatorFunc ruleConditionValidatorFunc
		wantErr       *RuleEngineError
	}{
		{
			name: "valid",
			args: args{
				c: &Condition{
					Type: AndCondition,
					SubConditions: []*Condition{
						{
							Type: OrCondition,
						},
						{
							Type: OrCondition,
						},
					},
				},
			},
			validatorFunc: subConditionCountRuleConditionValidator(2),
			wantErr:       nil,
		},
		{
			name: "invalid",
			args: args{
				c: &Condition{
					Type: AndCondition,
					SubConditions: []*Condition{
						{
							Type: OrCondition,
						},
					},
				},
			},
			validatorFunc: subConditionCountRuleConditionValidator(2),
			wantErr:       newError(ErrCodeInvalidSubConditionCount),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotErr := tt.validatorFunc(tt.args.c); !isErrorEqual(gotErr, tt.wantErr) {
				t.Errorf("operandValidator() = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func Test_RuleValidator_MinSubConditionCount(t *testing.T) {
	type args struct {
		c *Condition
	}
	tests := []struct {
		name          string
		args          args
		validatorFunc ruleConditionValidatorFunc
		wantErr       *RuleEngineError
	}{
		{
			name: "valid_exact",
			args: args{
				c: &Condition{
					Type: AndCondition,
					SubConditions: []*Condition{
						{
							Type: OrCondition,
						},
						{
							Type: OrCondition,
						},
					},
				},
			},
			validatorFunc: minSubConditionCountRuleConditionValidator(2),
			wantErr:       nil,
		},
		{
			name: "valid_more",
			args: args{
				c: &Condition{
					Type: AndCondition,
					SubConditions: []*Condition{
						{
							Type: OrCondition,
						},
						{
							Type: OrCondition,
						},
						{
							Type: OrCondition,
						},
					},
				},
			},
			validatorFunc: minSubConditionCountRuleConditionValidator(2),
			wantErr:       nil,
		},
		{
			name: "invalid",
			args: args{
				c: &Condition{
					Type: AndCondition,
					SubConditions: []*Condition{
						{
							Type: OrCondition,
						},
					},
				},
			},
			validatorFunc: minSubConditionCountRuleConditionValidator(2),
			wantErr:       newError(ErrCodeInvalidSubConditionCount),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotErr := tt.validatorFunc(tt.args.c); !isErrorEqual(gotErr, tt.wantErr) {
				t.Errorf("operandValidator() = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func Test_ruleEngineConfigValidator_validate(t *testing.T) {
	type testValidators struct {
		fieldValidators         []fieldValidatorFunc
		condTypeValidators      map[string][]conditionTypeValidatorFunc
		ruleConditionValidators map[string][]ruleConditionValidatorFunc
	}
	type args struct {
		config *RuleEngineConfig
	}
	tests := []struct {
		name       string
		validators testValidators
		args       args
		wantErr    *RuleEngineError
	}{
		{
			name: "valid_fields",
			validators: testValidators{
				fieldValidators: engineConfigValidator.fieldValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					Fields: Fields{
						"testInteger": Integer,
						"testFloat":   Float,
						"testString":  String,
						"testBoolean": Boolean,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "invalid_fields",
			validators: testValidators{
				fieldValidators: engineConfigValidator.fieldValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					Fields: Fields{
						"testInteger": unknownValueType,
						"testFloat":   Float,
						"testString":  String,
						"testBoolean": Boolean,
					},
				},
			},
			wantErr: newError(ErrCodeInvalidValueType),
		},
		{
			name: "valid_EqualOperatorBasedConditionType_Integer",
			validators: testValidators{
				condTypeValidators: engineConfigValidator.condTypeValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					Fields: Fields{
						"testInteger": Integer,
					},
					ConditionTypes: map[string]*ConditionType{
						"IntegerEqual": {
							Operator: EqualOperator,
							Operands: []*Operand{
								{
									Type:      Field,
									ValueType: Integer,
									Val:       "testInteger",
								},
								{
									Type:      Constant,
									ValueType: Integer,
									Val:       "10",
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "valid_EqualOperatorBasedConditionType_Float",
			validators: testValidators{
				condTypeValidators: engineConfigValidator.condTypeValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					Fields: Fields{
						"testFloat": Float,
					},
					ConditionTypes: map[string]*ConditionType{
						"IntegerEqual": {
							Operator: EqualOperator,
							Operands: []*Operand{
								{
									Type:      Field,
									ValueType: Float,
									Val:       "testFloat",
								},
								{
									Type:      Constant,
									ValueType: Float,
									Val:       "10.1",
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "valid_EqualOperatorBasedConditionType_Boolean",
			validators: testValidators{
				condTypeValidators: engineConfigValidator.condTypeValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					Fields: Fields{
						"testBoolean": Boolean,
					},
					ConditionTypes: map[string]*ConditionType{
						"IntegerEqual": {
							Operator: EqualOperator,
							Operands: []*Operand{
								{
									Type:      Field,
									ValueType: Boolean,
									Val:       "testBoolean",
								},
								{
									Type:      Constant,
									ValueType: Boolean,
									Val:       "true",
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "valid_EqualOperatorBasedConditionType_String",
			validators: testValidators{
				condTypeValidators: engineConfigValidator.condTypeValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					Fields: Fields{
						"testString": String,
					},
					ConditionTypes: map[string]*ConditionType{
						"IntegerEqual": {
							Operator: EqualOperator,
							Operands: []*Operand{
								{
									Type:      Field,
									ValueType: String,
									Val:       "testString",
								},
								{
									Type:      Constant,
									ValueType: String,
									Val:       "testString",
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "invalid_EqualOperatorBasedConditionType_InvalidOperandValueType",
			validators: testValidators{
				condTypeValidators: engineConfigValidator.condTypeValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					Fields: Fields{
						"testString": String,
					},
					ConditionTypes: map[string]*ConditionType{
						"IntegerEqual": {
							Operator: EqualOperator,
							Operands: []*Operand{
								{
									Type:      Field,
									ValueType: unknownValueType,
									Val:       "testString",
								},
								{
									Type:      Constant,
									ValueType: unknownValueType,
									Val:       "testString",
								},
							},
						},
					},
				},
			},
			wantErr: newError(ErrCodeInvalidOperand),
		},

		{
			name: "valid_NotEqualOperatorBasedConditionType_Integer",
			validators: testValidators{
				condTypeValidators: engineConfigValidator.condTypeValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					Fields: Fields{
						"testInteger": Integer,
					},
					ConditionTypes: map[string]*ConditionType{
						"IntegerEqual": {
							Operator: NotEqualOperator,
							Operands: []*Operand{
								{
									Type:      Field,
									ValueType: Integer,
									Val:       "testInteger",
								},
								{
									Type:      Constant,
									ValueType: Integer,
									Val:       "10",
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "valid_NotEqualOperatorBasedConditionType_Float",
			validators: testValidators{
				condTypeValidators: engineConfigValidator.condTypeValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					Fields: Fields{
						"testFloat": Float,
					},
					ConditionTypes: map[string]*ConditionType{
						"IntegerEqual": {
							Operator: NotEqualOperator,
							Operands: []*Operand{
								{
									Type:      Field,
									ValueType: Float,
									Val:       "testFloat",
								},
								{
									Type:      Constant,
									ValueType: Float,
									Val:       "10.1",
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "valid_NotEqualOperatorBasedConditionType_Boolean",
			validators: testValidators{
				condTypeValidators: engineConfigValidator.condTypeValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					Fields: Fields{
						"testBoolean": Boolean,
					},
					ConditionTypes: map[string]*ConditionType{
						"IntegerEqual": {
							Operator: NotEqualOperator,
							Operands: []*Operand{
								{
									Type:      Field,
									ValueType: Boolean,
									Val:       "testBoolean",
								},
								{
									Type:      Constant,
									ValueType: Boolean,
									Val:       "true",
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "valid_NotEqualOperatorBasedConditionType_String",
			validators: testValidators{
				condTypeValidators: engineConfigValidator.condTypeValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					Fields: Fields{
						"testString": String,
					},
					ConditionTypes: map[string]*ConditionType{
						"IntegerEqual": {
							Operator: NotEqualOperator,
							Operands: []*Operand{
								{
									Type:      Field,
									ValueType: String,
									Val:       "testString",
								},
								{
									Type:      Constant,
									ValueType: String,
									Val:       "testString",
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "invalid_NotEqualOperatorBasedConditionType_InvalidOperandValueType",
			validators: testValidators{
				condTypeValidators: engineConfigValidator.condTypeValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					Fields: Fields{
						"testString": String,
					},
					ConditionTypes: map[string]*ConditionType{
						"IntegerEqual": {
							Operator: NotEqualOperator,
							Operands: []*Operand{
								{
									Type:      Field,
									ValueType: unknownValueType,
									Val:       "testString",
								},
								{
									Type:      Constant,
									ValueType: unknownValueType,
									Val:       "testString",
								},
							},
						},
					},
				},
			},
			wantErr: newError(ErrCodeInvalidOperand),
		},

		{
			name: "valid_GreaterOperatorBasedConditionType_Integer",
			validators: testValidators{
				condTypeValidators: engineConfigValidator.condTypeValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					Fields: Fields{
						"testInteger": Integer,
					},
					ConditionTypes: map[string]*ConditionType{
						"IntegerEqual": {
							Operator: GreaterOperator,
							Operands: []*Operand{
								{
									Type:      Field,
									ValueType: Integer,
									Val:       "testInteger",
								},
								{
									Type:      Constant,
									ValueType: Integer,
									Val:       "10",
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "valid_GreaterOperatorBasedConditionType_Float",
			validators: testValidators{
				condTypeValidators: engineConfigValidator.condTypeValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					Fields: Fields{
						"testFloat": Float,
					},
					ConditionTypes: map[string]*ConditionType{
						"IntegerEqual": {
							Operator: GreaterOperator,
							Operands: []*Operand{
								{
									Type:      Field,
									ValueType: Float,
									Val:       "testFloat",
								},
								{
									Type:      Constant,
									ValueType: Float,
									Val:       "10.1",
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "invalid_GreaterOperatorBasedConditionType_InvalidOperandValueType",
			validators: testValidators{
				condTypeValidators: engineConfigValidator.condTypeValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					Fields: Fields{
						"testString": String,
					},
					ConditionTypes: map[string]*ConditionType{
						"IntegerEqual": {
							Operator: GreaterOperator,
							Operands: []*Operand{
								{
									Type:      Field,
									ValueType: unknownValueType,
									Val:       "testString",
								},
								{
									Type:      Constant,
									ValueType: unknownValueType,
									Val:       "testString",
								},
							},
						},
					},
				},
			},
			wantErr: newError(ErrCodeInvalidOperand),
		},

		{
			name: "valid_GreaterEqualOperatorBasedConditionType_Integer",
			validators: testValidators{
				condTypeValidators: engineConfigValidator.condTypeValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					Fields: Fields{
						"testInteger": Integer,
					},
					ConditionTypes: map[string]*ConditionType{
						"IntegerEqual": {
							Operator: GreaterEqualOperator,
							Operands: []*Operand{
								{
									Type:      Field,
									ValueType: Integer,
									Val:       "testInteger",
								},
								{
									Type:      Constant,
									ValueType: Integer,
									Val:       "10",
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "valid_GreaterEqualOperatorBasedConditionType_Float",
			validators: testValidators{
				condTypeValidators: engineConfigValidator.condTypeValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					Fields: Fields{
						"testFloat": Float,
					},
					ConditionTypes: map[string]*ConditionType{
						"IntegerEqual": {
							Operator: GreaterEqualOperator,
							Operands: []*Operand{
								{
									Type:      Field,
									ValueType: Float,
									Val:       "testFloat",
								},
								{
									Type:      Constant,
									ValueType: Float,
									Val:       "10.1",
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "invalid_GreaterEqualOperatorBasedConditionType_InvalidOperandValueType",
			validators: testValidators{
				condTypeValidators: engineConfigValidator.condTypeValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					Fields: Fields{
						"testString": String,
					},
					ConditionTypes: map[string]*ConditionType{
						"IntegerEqual": {
							Operator: GreaterEqualOperator,
							Operands: []*Operand{
								{
									Type:      Field,
									ValueType: unknownValueType,
									Val:       "testString",
								},
								{
									Type:      Constant,
									ValueType: unknownValueType,
									Val:       "testString",
								},
							},
						},
					},
				},
			},
			wantErr: newError(ErrCodeInvalidOperand),
		},

		{
			name: "valid_LessOperatorBasedConditionType_Integer",
			validators: testValidators{
				condTypeValidators: engineConfigValidator.condTypeValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					Fields: Fields{
						"testInteger": Integer,
					},
					ConditionTypes: map[string]*ConditionType{
						"IntegerEqual": {
							Operator: LessOperator,
							Operands: []*Operand{
								{
									Type:      Field,
									ValueType: Integer,
									Val:       "testInteger",
								},
								{
									Type:      Constant,
									ValueType: Integer,
									Val:       "10",
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "valid_LessOperatorBasedConditionType_Float",
			validators: testValidators{
				condTypeValidators: engineConfigValidator.condTypeValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					Fields: Fields{
						"testFloat": Float,
					},
					ConditionTypes: map[string]*ConditionType{
						"IntegerEqual": {
							Operator: LessOperator,
							Operands: []*Operand{
								{
									Type:      Field,
									ValueType: Float,
									Val:       "testFloat",
								},
								{
									Type:      Constant,
									ValueType: Float,
									Val:       "10.1",
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "invalid_LessOperatorBasedConditionType_InvalidOperandValueType",
			validators: testValidators{
				condTypeValidators: engineConfigValidator.condTypeValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					Fields: Fields{
						"testString": String,
					},
					ConditionTypes: map[string]*ConditionType{
						"IntegerEqual": {
							Operator: LessOperator,
							Operands: []*Operand{
								{
									Type:      Field,
									ValueType: unknownValueType,
									Val:       "testString",
								},
								{
									Type:      Constant,
									ValueType: unknownValueType,
									Val:       "testString",
								},
							},
						},
					},
				},
			},
			wantErr: newError(ErrCodeInvalidOperand),
		},

		{
			name: "valid_LessEqualOperatorBasedConditionType_Integer",
			validators: testValidators{
				condTypeValidators: engineConfigValidator.condTypeValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					Fields: Fields{
						"testInteger": Integer,
					},
					ConditionTypes: map[string]*ConditionType{
						"IntegerEqual": {
							Operator: LessEqualOperator,
							Operands: []*Operand{
								{
									Type:      Field,
									ValueType: Integer,
									Val:       "testInteger",
								},
								{
									Type:      Constant,
									ValueType: Integer,
									Val:       "10",
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "valid_LessEqualOperatorBasedConditionType_Float",
			validators: testValidators{
				condTypeValidators: engineConfigValidator.condTypeValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					Fields: Fields{
						"testFloat": Float,
					},
					ConditionTypes: map[string]*ConditionType{
						"IntegerEqual": {
							Operator: LessOperator,
							Operands: []*Operand{
								{
									Type:      Field,
									ValueType: Float,
									Val:       "testFloat",
								},
								{
									Type:      Constant,
									ValueType: Float,
									Val:       "10.1",
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "invalid_LessEqualOperatorBasedConditionType_InvalidOperandValueType",
			validators: testValidators{
				condTypeValidators: engineConfigValidator.condTypeValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					Fields: Fields{
						"testString": String,
					},
					ConditionTypes: map[string]*ConditionType{
						"IntegerEqual": {
							Operator: LessOperator,
							Operands: []*Operand{
								{
									Type:      Field,
									ValueType: unknownValueType,
									Val:       "testString",
								},
								{
									Type:      Constant,
									ValueType: unknownValueType,
									Val:       "testString",
								},
							},
						},
					},
				},
			},
			wantErr: newError(ErrCodeInvalidOperand),
		},

		{
			name: "valid_ContainOperatorBasedConditionType_Integer",
			validators: testValidators{
				condTypeValidators: engineConfigValidator.condTypeValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					Fields: Fields{
						"testString": String,
					},
					ConditionTypes: map[string]*ConditionType{
						"IntegerEqual": {
							Operator: ContainOperator,
							Operands: []*Operand{
								{
									Type:      Field,
									ValueType: String,
									Val:       "testString",
								},
								{
									Type:      Constant,
									ValueType: String,
									Val:       "10",
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "invalid_ContainOperatorBasedConditionType_InvalidOperandValueType",
			validators: testValidators{
				condTypeValidators: engineConfigValidator.condTypeValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					Fields: Fields{
						"testString": String,
					},
					ConditionTypes: map[string]*ConditionType{
						"IntegerEqual": {
							Operator: LessOperator,
							Operands: []*Operand{
								{
									Type:      Field,
									ValueType: unknownValueType,
									Val:       "testString",
								},
								{
									Type:      Constant,
									ValueType: unknownValueType,
									Val:       "testString",
								},
							},
						},
					},
				},
			},
			wantErr: newError(ErrCodeInvalidOperand),
		},

		{
			name: "valid_RuleCondition_singleLevel",
			validators: testValidators{
				ruleConditionValidators: engineConfigValidator.ruleConditionValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					ConditionTypes: map[string]*ConditionType{
						"testCondition": {
							Operator: GreaterOperator},
					},
					Rules: map[string]*RuleConfig{
						"testRule": {
							Priority: 1,
							RootCondition: &Condition{
								Type: OrCondition,
								SubConditions: []*Condition{
									{
										Type: "testCondition",
									},
									{
										Type: "testCondition",
									},
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},

		{
			name: "valid_RuleCondition_multiLevel",
			validators: testValidators{
				ruleConditionValidators: engineConfigValidator.ruleConditionValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					ConditionTypes: map[string]*ConditionType{
						"testCondition": {
							Operator: GreaterOperator},
					},
					Rules: map[string]*RuleConfig{
						"testRule": {
							Priority: 1,
							RootCondition: &Condition{
								Type: OrCondition,
								SubConditions: []*Condition{
									{
										Type: AndCondition,
										SubConditions: []*Condition{
											{Type: "testCondition"},
											{Type: "testCondition"},
										},
									},
									{
										Type: AndCondition,
										SubConditions: []*Condition{
											{Type: "testCondition"},
											{Type: "testCondition"},
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},

		{
			name: "invalid_RuleCondition_singleLevel",
			validators: testValidators{
				ruleConditionValidators: engineConfigValidator.ruleConditionValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					ConditionTypes: map[string]*ConditionType{
						"testCondition": {
							Operator: GreaterOperator},
					},
					Rules: map[string]*RuleConfig{
						"testRule": {
							Priority: 1,
							RootCondition: &Condition{
								Type: OrCondition,
								SubConditions: []*Condition{
									{
										Type: "testCondition",
									},
								},
							},
						},
					},
				},
			},
			wantErr: newError(ErrCodeInvalidSubConditionCount),
		},

		{
			name: "invalid_RuleCondition_multiLevel",
			validators: testValidators{
				ruleConditionValidators: engineConfigValidator.ruleConditionValidators,
			},
			args: args{
				config: &RuleEngineConfig{
					ConditionTypes: map[string]*ConditionType{
						"testCondition": {
							Operator: GreaterOperator},
					},
					Rules: map[string]*RuleConfig{
						"testRule": {
							Priority: 1,
							RootCondition: &Condition{
								Type: OrCondition,
								SubConditions: []*Condition{
									{
										Type: AndCondition,
										SubConditions: []*Condition{
											{Type: "testCondition"},
											{Type: "testCondition"},
										},
									},
									{
										Type: AndCondition,
										SubConditions: []*Condition{
											{Type: "testCondition"},
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: newError(ErrCodeInvalidSubConditionCount),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &ruleEngineConfigValidator{
				fieldValidators:         tt.validators.fieldValidators,
				condTypeValidators:      tt.validators.condTypeValidators,
				ruleConditionValidators: tt.validators.ruleConditionValidators,
			}

			if gotErr := v.validate(tt.args.config); !isErrorEqual(gotErr, tt.wantErr) {
				t.Errorf("ruleEngineConfigValidator.validate() gotErr %v, wantErr %v", gotErr, tt.wantErr)
			}
		})
	}
}
