package ruleenginecore

import (
	"testing"
)

type testTrueEvaluator struct{}

func (t *testTrueEvaluator) evaluate(input parsedInput) bool {
	return true
}

type testFalseEvaluator struct{}

func (t *testFalseEvaluator) evaluate(input parsedInput) bool {
	return false
}

var trueEvaluator *testTrueEvaluator = &testTrueEvaluator{}

var falseEvaluator *testFalseEvaluator = &testFalseEvaluator{}

func Test_logicalEvaluator_evaluate(t *testing.T) {

	type fields struct {
		operator        string
		innerEvaluators []evaluator
	}
	type args struct {
		input parsedInput
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      bool
		wantPanic bool
	}{
		{
			name: "validLogicalNegateEvaluatorTrue",
			fields: fields{
				operator: NegationCondition,
				innerEvaluators: []evaluator{
					trueEvaluator,
				},
			},
			args: args{
				map[string]any{},
			},
			want:      false,
			wantPanic: false,
		},
		{
			name: "validLogicalNegateEvaluatorFalse",
			fields: fields{
				operator: NegationCondition,
				innerEvaluators: []evaluator{
					falseEvaluator,
				},
			},
			args: args{
				map[string]any{},
			},
			want:      true,
			wantPanic: false,
		},
		{
			name: "validLogicalAndEvaluatorBothTrue",
			fields: fields{
				operator: AndCondition,
				innerEvaluators: []evaluator{
					trueEvaluator,
					trueEvaluator,
				},
			},
			args: args{
				map[string]any{},
			},
			want:      true,
			wantPanic: false,
		},
		{
			name: "validLogicalAndEvaluatorFalseTrue",
			fields: fields{
				operator: AndCondition,
				innerEvaluators: []evaluator{
					trueEvaluator,
					falseEvaluator,
				},
			},
			args: args{
				map[string]any{},
			},
			want:      false,
			wantPanic: false,
		},
		{
			name: "validLogicalOrEvaluatorBothFalse",
			fields: fields{
				operator: OrCondition,
				innerEvaluators: []evaluator{
					falseEvaluator,
					falseEvaluator,
				},
			},
			args: args{
				map[string]any{},
			},
			want:      false,
			wantPanic: false,
		},
		{
			name: "validLogicalOrEvaluatorFalseTrue",
			fields: fields{
				operator: OrCondition,
				innerEvaluators: []evaluator{
					falseEvaluator,
					trueEvaluator,
				},
			},
			args: args{
				map[string]any{},
			},
			want:      true,
			wantPanic: false,
		},
		{
			name: "InvalidLogicalEvaluatorPanic",
			fields: fields{
				operator: "panic",
				innerEvaluators: []evaluator{
					falseEvaluator,
					trueEvaluator,
				},
			},
			args: args{
				map[string]any{},
			},
			want:      true,
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

				t.Errorf("logicalEvaluator.evaluate() gotPanic:%v , want:%v wantPanic:%v", r != nil, tt.want, tt.wantPanic)
			}()

			le := &logicalEvaluator{
				operator:        tt.fields.operator,
				innerEvaluators: tt.fields.innerEvaluators,
			}
			if got := le.evaluate(tt.args.input); got != tt.want {
				t.Errorf("logicalEvaluator.evaluate() got:%v, want:%v wantPanic:%v", got, tt.want, tt.wantPanic)
			}
		})
	}
}

func Test_greaterEvaluator_evaluate(t *testing.T) {
	type testCondition struct {
		operands []*Operand
	}
	type args struct {
		input map[string]any
	}
	tests := []struct {
		name      string
		condition testCondition
		args      args
		want      bool
		wantPanic bool
	}{
		{
			name: "IntType_valid",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Integer,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Integer,
						Type:       Constant,
						Val:        "10",
						typedValue: int64(10),
					},
				},
			},
			args: args{
				map[string]any{
					"count": int64(20),
				},
			},
			want:      true,
			wantPanic: false,
		},
		{
			name: "IntType_PassedInvalidFieldType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Integer,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Integer,
						Type:       Constant,
						Val:        "10",
						typedValue: int64(10),
					},
				},
			},
			args: args{
				map[string]any{
					"count": "asdf", // invalid
				},
			},
			want:      false,
			wantPanic: true,
		},
		{
			name: "IntType_PassedInvalidConstantType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Integer,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Integer,
						Type:       Constant,
						Val:        "asdf",
						typedValue: "asdf", // invalid
					},
				},
			},
			args: args{
				map[string]any{
					"count": int64(20),
				},
			},
			want:      false,
			wantPanic: true,
		},
		{
			name: "FloatType_valid",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Float,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Float,
						Type:       Constant,
						Val:        "10.1",
						typedValue: float64(10.1),
					},
				},
			},
			args: args{
				map[string]any{
					"count": float64(20.1),
				},
			},
			want:      true,
			wantPanic: false,
		},
		{
			name: "FloatType_PassedInvalidFieldType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Float,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Float,
						Type:       Constant,
						Val:        "10.1",
						typedValue: float64(10.1),
					},
				},
			},
			args: args{
				map[string]any{
					"count": "asdf", // invalid
				},
			},
			want:      false,
			wantPanic: true,
		},
		{
			name: "FloatType_PassedInvalidConstantType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Float,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Float,
						Type:       Constant,
						Val:        "asdf",
						typedValue: "asdf", // invalid
					},
				},
			},
			args: args{
				map[string]any{
					"count": float64(20.1),
				},
			},
			want:      false,
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

				t.Errorf("greaterEvaluator.evaluate() gotPanic:%v , want:%v wantPanic:%v", r != nil, tt.want, tt.wantPanic)
			}()

			ge := &greaterEvaluator{
				operands: tt.condition.operands,
			}
			if got := ge.evaluate(tt.args.input); got != tt.want {
				t.Errorf("greaterEvaluator.evaluate() got:%v, want:%v wantPanic:%v", got, tt.want, tt.wantPanic)
			}
		})
	}
}

func Test_greaterEqualEvaluator_evaluate(t *testing.T) {
	type testCondition struct {
		operands []*Operand
	}
	type args struct {
		input map[string]any
	}
	tests := []struct {
		name      string
		condition testCondition
		args      args
		want      bool
		wantPanic bool
	}{
		{
			name: "IntType_valid",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Integer,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Integer,
						Type:       Constant,
						Val:        "10",
						typedValue: int64(10),
					},
				},
			},
			args: args{
				map[string]any{
					"count": int64(20),
				},
			},
			want:      true,
			wantPanic: false,
		},
		{
			name: "IntType_PassedInvalidFieldType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Integer,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Integer,
						Type:       Constant,
						Val:        "10",
						typedValue: int64(10),
					},
				},
			},
			args: args{
				map[string]any{
					"count": "asdf", // invalid
				},
			},
			want:      false,
			wantPanic: true,
		},
		{
			name: "IntType_PassedInvalidConstantType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Integer,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Integer,
						Type:       Constant,
						Val:        "asdf",
						typedValue: "asdf", // invalid
					},
				},
			},
			args: args{
				map[string]any{
					"count": int64(20),
				},
			},
			want:      false,
			wantPanic: true,
		},
		{
			name: "FloatType_valid",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Float,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Float,
						Type:       Constant,
						Val:        "10.1",
						typedValue: float64(10.1),
					},
				},
			},
			args: args{
				map[string]any{
					"count": float64(20.1),
				},
			},
			want:      true,
			wantPanic: false,
		},
		{
			name: "FloatType_PassedInvalidFieldType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Float,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Float,
						Type:       Constant,
						Val:        "10.1",
						typedValue: float64(10.1),
					},
				},
			},
			args: args{
				map[string]any{
					"count": "asdf", // invalid
				},
			},
			want:      false,
			wantPanic: true,
		},
		{
			name: "FloatType_PassedInvalidConstantType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Float,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Float,
						Type:       Constant,
						Val:        "asdf",
						typedValue: "asdf", // invalid
					},
				},
			},
			args: args{
				map[string]any{
					"count": float64(20.1),
				},
			},
			want:      false,
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

				t.Errorf("greaterEqualEvaluator.evaluate() gotPanic:%v , want:%v wantPanic:%v", r != nil, tt.want, tt.wantPanic)
			}()

			gte := &greaterEqualEvaluator{
				operands: tt.condition.operands,
			}
			if got := gte.evaluate(tt.args.input); got != tt.want {
				t.Errorf("greaterEqualEvaluator.evaluate() got:%v, want:%v wantPanic:%v", got, tt.want, tt.wantPanic)
			}
		})
	}
}

func Test_lessEvaluator_evaluate(t *testing.T) {
	type testCondition struct {
		operands []*Operand
	}
	type args struct {
		input map[string]any
	}
	tests := []struct {
		name      string
		condition testCondition
		args      args
		want      bool
		wantPanic bool
	}{
		{
			name: "IntType_valid",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Integer,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Integer,
						Type:       Constant,
						Val:        "20",
						typedValue: int64(20),
					},
				},
			},
			args: args{
				map[string]any{
					"count": int64(10),
				},
			},
			want:      true,
			wantPanic: false,
		},
		{
			name: "IntType_PassedInvalidFieldType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Integer,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Integer,
						Type:       Constant,
						Val:        "20",
						typedValue: int64(20),
					},
				},
			},
			args: args{
				map[string]any{
					"count": "asdf", // invalid
				},
			},
			want:      false,
			wantPanic: true,
		},
		{
			name: "IntType_PassedInvalidConstantType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Integer,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Integer,
						Type:       Constant,
						Val:        "asdf",
						typedValue: "asdf", // invalid
					},
				},
			},
			args: args{
				map[string]any{
					"count": int64(10),
				},
			},
			want:      false,
			wantPanic: true,
		},
		{
			name: "FloatType_valid",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Float,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Float,
						Type:       Constant,
						Val:        "20.1",
						typedValue: float64(20.1),
					},
				},
			},
			args: args{
				map[string]any{
					"count": float64(10.1),
				},
			},
			want:      true,
			wantPanic: false,
		},
		{
			name: "FloatType_PassedInvalidFieldType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Float,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Float,
						Type:       Constant,
						Val:        "20.1",
						typedValue: float64(20.1),
					},
				},
			},
			args: args{
				map[string]any{
					"count": "asdf", // invalid
				},
			},
			want:      false,
			wantPanic: true,
		},
		{
			name: "FloatType_PassedInvalidConstantType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Float,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Float,
						Type:       Constant,
						Val:        "asdf",
						typedValue: "asdf", // invalid
					},
				},
			},
			args: args{
				map[string]any{
					"count": float64(10.1),
				},
			},
			want:      false,
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

				t.Errorf("lessEvaluator.evaluate() gotPanic:%v , want:%v wantPanic:%v", r != nil, tt.want, tt.wantPanic)
			}()

			lt := &lessEvaluator{
				operands: tt.condition.operands,
			}
			if got := lt.evaluate(tt.args.input); got != tt.want {
				t.Errorf("lessEvaluator.evaluate() got:%v, want:%v wantPanic:%v", got, tt.want, tt.wantPanic)
			}
		})
	}
}

func Test_lessEqualEvaluator_evaluate(t *testing.T) {
	type testCondition struct {
		operands []*Operand
	}
	type args struct {
		input map[string]any
	}
	tests := []struct {
		name      string
		condition testCondition
		args      args
		want      bool
		wantPanic bool
	}{
		{
			name: "IntType_valid",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Integer,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Integer,
						Type:       Constant,
						Val:        "20",
						typedValue: int64(20),
					},
				},
			},
			args: args{
				map[string]any{
					"count": int64(10),
				},
			},
			want:      true,
			wantPanic: false,
		},
		{
			name: "IntType_PassedInvalidFieldType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Integer,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Integer,
						Type:       Constant,
						Val:        "20",
						typedValue: int64(20),
					},
				},
			},
			args: args{
				map[string]any{
					"count": "asdf", // invalid
				},
			},
			want:      false,
			wantPanic: true,
		},
		{
			name: "IntType_PassedInvalidConstantType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Integer,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Integer,
						Type:       Constant,
						Val:        "asdf",
						typedValue: "asdf", // invalid
					},
				},
			},
			args: args{
				map[string]any{
					"count": int64(10),
				},
			},
			want:      false,
			wantPanic: true,
		},
		{
			name: "FloatType_valid",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Float,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Float,
						Type:       Constant,
						Val:        "20.1",
						typedValue: float64(20.1),
					},
				},
			},
			args: args{
				map[string]any{
					"count": float64(10.1),
				},
			},
			want:      true,
			wantPanic: false,
		},
		{
			name: "FloatType_PassedInvalidFieldType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Float,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Float,
						Type:       Constant,
						Val:        "20.1",
						typedValue: float64(20.1),
					},
				},
			},
			args: args{
				map[string]any{
					"count": "asdf", // invalid
				},
			},
			want:      false,
			wantPanic: true,
		},
		{
			name: "FloatType_PassedInvalidConstantType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Float,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Float,
						Type:       Constant,
						Val:        "asdf",
						typedValue: "asdf", // invalid
					},
				},
			},
			args: args{
				map[string]any{
					"count": float64(10.1),
				},
			},
			want:      false,
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

				t.Errorf("lessEqualEvaluator.evaluate() gotPanic:%v , want:%v wantPanic:%v", r != nil, tt.want, tt.wantPanic)
			}()

			lte := &lessEqualEvaluator{
				operands: tt.condition.operands,
			}
			if got := lte.evaluate(tt.args.input); got != tt.want {
				t.Errorf("lessEqualEvaluator.evaluate() got:%v, want:%v wantPanic:%v", got, tt.want, tt.wantPanic)
			}
		})
	}
}

func Test_equalEvaluator_evaluate(t *testing.T) {
	type testCondition struct {
		operands []*Operand
	}
	type args struct {
		input map[string]any
	}
	tests := []struct {
		name      string
		condition testCondition
		args      args
		want      bool
		wantPanic bool
	}{
		{
			name: "IntType_valid",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Integer,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Integer,
						Type:       Constant,
						Val:        "20",
						typedValue: int64(20),
					},
				},
			},
			args: args{
				map[string]any{
					"count": int64(20),
				},
			},
			want:      true,
			wantPanic: false,
		},
		{
			name: "IntType_PassedInvalidFieldType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Integer,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Integer,
						Type:       Constant,
						Val:        "20",
						typedValue: int64(20),
					},
				},
			},
			args: args{
				map[string]any{
					"count": "asdf", // invalid
				},
			},
			want:      false,
			wantPanic: true,
		},
		{
			name: "IntType_PassedInvalidConstantType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Integer,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Integer,
						Type:       Constant,
						Val:        "asdf",
						typedValue: "asdf", // invalid
					},
				},
			},
			args: args{
				map[string]any{
					"count": int64(10),
				},
			},
			want:      false,
			wantPanic: true,
		},
		{
			name: "FloatType_valid",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Float,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Float,
						Type:       Constant,
						Val:        "20.1",
						typedValue: float64(20.1),
					},
				},
			},
			args: args{
				map[string]any{
					"count": float64(20.1),
				},
			},
			want:      true,
			wantPanic: false,
		},
		{
			name: "FloatType_PassedInvalidFieldType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Float,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Float,
						Type:       Constant,
						Val:        "20.1",
						typedValue: float64(20.1),
					},
				},
			},
			args: args{
				map[string]any{
					"count": "asdf", // invalid
				},
			},
			want:      false,
			wantPanic: true,
		},
		{
			name: "FloatType_PassedInvalidConstantType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Float,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Float,
						Type:       Constant,
						Val:        "asdf",
						typedValue: "asdf", // invalid
					},
				},
			},
			args: args{
				map[string]any{
					"count": float64(20.1),
				},
			},
			want:      false,
			wantPanic: true,
		},
		{
			name: "StringType_valid",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: String,
						Type:      Field,
						Val:       "firstname",
					},
					{
						ValueType:  String,
						Type:       Constant,
						Val:        "ironman",
						typedValue: "ironman",
					},
				},
			},
			args: args{
				map[string]any{
					"firstname": "ironman",
				},
			},
			want:      true,
			wantPanic: false,
		},
		{
			name: "StringType_PassedInvalidFieldType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: String,
						Type:      Field,
						Val:       "firstname",
					},
					{
						ValueType:  String,
						Type:       Constant,
						Val:        "ironman",
						typedValue: "ironman",
					},
				},
			},
			args: args{
				map[string]any{
					"firstname": 1, // invalid
				},
			},
			want:      false,
			wantPanic: true,
		},
		{
			name: "StringType_PassedInvalidConstantType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: String,
						Type:      Field,
						Val:       "firstname",
					},
					{
						ValueType:  String,
						Type:       Constant,
						Val:        "1",
						typedValue: 1, // invalid
					},
				},
			},
			args: args{
				map[string]any{
					"firstname": "ironman",
				},
			},
			want:      false,
			wantPanic: true,
		},
		{
			name: "BoolType_valid",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Boolean,
						Type:      Field,
						Val:       "IsHoliday",
					},
					{
						ValueType:  Boolean,
						Type:       Constant,
						Val:        "true",
						typedValue: true,
					},
				},
			},
			args: args{
				map[string]any{
					"IsHoliday": true,
				},
			},
			want:      true,
			wantPanic: false,
		},
		{
			name: "BoolType_PassedInvalidFieldType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Boolean,
						Type:      Field,
						Val:       "IsHoliday",
					},
					{
						ValueType:  Boolean,
						Type:       Constant,
						Val:        "true",
						typedValue: true,
					},
				},
			},
			args: args{
				map[string]any{
					"IsHoliday": 1, // invalid
				},
			},
			want:      false,
			wantPanic: true,
		},
		{
			name: "BoolType_PassedInvalidConstantType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Boolean,
						Type:      Field,
						Val:       "IsHoliday",
					},
					{
						ValueType:  Boolean,
						Type:       Constant,
						Val:        "1",
						typedValue: 1, // invalid
					},
				},
			},
			args: args{
				map[string]any{
					"IsHoliday": true,
				},
			},
			want:      false,
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

				t.Errorf("equalEvaluator.evaluate() gotPanic:%v , want:%v wantPanic:%v", r != nil, tt.want, tt.wantPanic)
			}()
			eq := &equalEvaluator{
				operands: tt.condition.operands,
			}
			if got := eq.evaluate(tt.args.input); got != tt.want {
				t.Errorf("equalEvaluator.evaluate() got:%v, want:%v wantPanic:%v", got, tt.want, tt.wantPanic)
			}
		})
	}
}

func Test_notEqualEvaluator_evaluate(t *testing.T) {
	type testCondition struct {
		operands []*Operand
	}
	type args struct {
		input map[string]any
	}
	tests := []struct {
		name      string
		condition testCondition
		args      args
		want      bool
		wantPanic bool
	}{
		{
			name: "IntType_valid",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Integer,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Integer,
						Type:       Constant,
						Val:        "20",
						typedValue: int64(20),
					},
				},
			},
			args: args{
				map[string]any{
					"count": int64(10),
				},
			},
			want:      true,
			wantPanic: false,
		},
		{
			name: "IntType_PassedInvalidFieldType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Integer,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Integer,
						Type:       Constant,
						Val:        "20",
						typedValue: int64(20),
					},
				},
			},
			args: args{
				map[string]any{
					"count": "asdf", // invalid
				},
			},
			want:      false,
			wantPanic: true,
		},
		{
			name: "IntType_PassedInvalidConstantType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Integer,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Integer,
						Type:       Constant,
						Val:        "asdf",
						typedValue: "asdf", // invalid
					},
				},
			},
			args: args{
				map[string]any{
					"count": int64(10),
				},
			},
			want:      false,
			wantPanic: true,
		},
		{
			name: "FloatType_valid",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Float,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Float,
						Type:       Constant,
						Val:        "20.1",
						typedValue: float64(20.1),
					},
				},
			},
			args: args{
				map[string]any{
					"count": float64(10.1),
				},
			},
			want:      true,
			wantPanic: false,
		},
		{
			name: "FloatType_PassedInvalidFieldType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Float,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Float,
						Type:       Constant,
						Val:        "20.1",
						typedValue: float64(20.1),
					},
				},
			},
			args: args{
				map[string]any{
					"count": "asdf", // invalid
				},
			},
			want:      false,
			wantPanic: true,
		},
		{
			name: "FloatType_PassedInvalidConstantType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Float,
						Type:      Field,
						Val:       "count",
					},
					{
						ValueType:  Float,
						Type:       Constant,
						Val:        "asdf",
						typedValue: "asdf", // invalid
					},
				},
			},
			args: args{
				map[string]any{
					"count": float64(20.1),
				},
			},
			want:      false,
			wantPanic: true,
		},
		{
			name: "StringType_valid",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: String,
						Type:      Field,
						Val:       "firstname",
					},
					{
						ValueType:  String,
						Type:       Constant,
						Val:        "ironman",
						typedValue: "ironman",
					},
				},
			},
			args: args{
				map[string]any{
					"firstname": "spiderman",
				},
			},
			want:      true,
			wantPanic: false,
		},
		{
			name: "StringType_PassedInvalidFieldType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: String,
						Type:      Field,
						Val:       "firstname",
					},
					{
						ValueType:  String,
						Type:       Constant,
						Val:        "ironman",
						typedValue: "ironman",
					},
				},
			},
			args: args{
				map[string]any{
					"firstname": 1, // invalid
				},
			},
			want:      false,
			wantPanic: true,
		},
		{
			name: "StringType_PassedInvalidConstantType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: String,
						Type:      Field,
						Val:       "firstname",
					},
					{
						ValueType:  String,
						Type:       Constant,
						Val:        "1",
						typedValue: 1, // invalid
					},
				},
			},
			args: args{
				map[string]any{
					"firstname": "ironman",
				},
			},
			want:      false,
			wantPanic: true,
		},
		{
			name: "BoolType_valid",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Boolean,
						Type:      Field,
						Val:       "IsHoliday",
					},
					{
						ValueType:  Boolean,
						Type:       Constant,
						Val:        "true",
						typedValue: true,
					},
				},
			},
			args: args{
				map[string]any{
					"IsHoliday": false,
				},
			},
			want:      true,
			wantPanic: false,
		},
		{
			name: "BoolType_PassedInvalidFieldType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Boolean,
						Type:      Field,
						Val:       "IsHoliday",
					},
					{
						ValueType:  Boolean,
						Type:       Constant,
						Val:        "true",
						typedValue: true,
					},
				},
			},
			args: args{
				map[string]any{
					"IsHoliday": 1, // invalid
				},
			},
			want:      false,
			wantPanic: true,
		},
		{
			name: "BoolType_PassedInvalidConstantType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: Boolean,
						Type:      Field,
						Val:       "IsHoliday",
					},
					{
						ValueType:  Boolean,
						Type:       Constant,
						Val:        "1",
						typedValue: 1, // invalid
					},
				},
			},
			args: args{
				map[string]any{
					"IsHoliday": true,
				},
			},
			want:      false,
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

				t.Errorf("notEqualEvaluator.evaluate() gotPanic:%v , want:%v wantPanic:%v", r != nil, tt.want, tt.wantPanic)
			}()
			neq := &notEqualEvaluator{
				operands: tt.condition.operands,
			}
			if got := neq.evaluate(tt.args.input); got != tt.want {
				t.Errorf("notEqualEvaluator.evaluate() got:%v, want:%v wantPanic:%v", got, tt.want, tt.wantPanic)
			}
		})
	}
}

func Test_containEvaluator_evaluate(t *testing.T) {
	type testCondition struct {
		operands []*Operand
	}
	type args struct {
		input map[string]any
	}
	tests := []struct {
		name      string
		condition testCondition
		args      args
		want      bool
		wantPanic bool
	}{
		{
			name: "StringType_valid",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: String,
						Type:      Field,
						Val:       "story",
					},
					{
						ValueType:  String,
						Type:       Constant,
						Val:        "dog",
						typedValue: "dog",
					},
				},
			},
			args: args{
				map[string]any{
					"story": "dogs are loyal",
				},
			},
			want:      true,
			wantPanic: false,
		},
		{
			name: "StringType_PassedInvalidFieldType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: String,
						Type:      Field,
						Val:       "story",
					},
					{
						ValueType:  String,
						Type:       Constant,
						Val:        "dog",
						typedValue: "dog",
					},
				},
			},
			args: args{
				map[string]any{
					"story": 1, // invalid
				},
			},
			want:      false,
			wantPanic: true,
		},
		{
			name: "StringType_PassedInvalidConstantType",
			condition: testCondition{
				operands: []*Operand{
					{
						ValueType: String,
						Type:      Field,
						Val:       "story",
					},
					{
						ValueType:  String,
						Type:       Constant,
						Val:        "1",
						typedValue: 1, // invalid
					},
				},
			},
			args: args{
				map[string]any{
					"story": "dogs are loyal",
				},
			},
			want:      false,
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

				t.Errorf("containEvaluator.evaluate() gotPanic:%v , want:%v wantPanic:%v", r != nil, tt.want, tt.wantPanic)
			}()
			ce := &containEvaluator{
				operands: tt.condition.operands,
			}
			if got := ce.evaluate(tt.args.input); got != tt.want {
				t.Errorf("containEvaluator.evaluate() got:%v, want:%v wantPanic:%v", got, tt.want, tt.wantPanic)
			}
		})
	}
}
