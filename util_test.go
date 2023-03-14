package ruleenginecore

import (
	"reflect"
	"testing"
)

func Test_parseValue(t *testing.T) {
	type args struct {
		value  string
		toType ValueType
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr *RuleEngineError
	}{
		{
			name: "valid_boolean",
			args: args{
				value:  "true",
				toType: Boolean,
			},
			want:    true,
			wantErr: nil,
		},
		{
			name: "invalid_boolean",
			args: args{
				value:  "invalid",
				toType: Boolean,
			},
			want:    false,
			wantErr: newError(ErrCodeParsingFailed),
		},
		{
			name: "valid_Int_1",
			args: args{
				value:  "1",
				toType: Integer,
			},
			want:    int64(1),
			wantErr: nil,
		},
		{
			name: "valid_Int_2",
			args: args{
				value:  "-1",
				toType: Integer,
			},
			want:    int64(-1),
			wantErr: nil,
		},
		{
			name: "invalid_Int",
			args: args{
				value:  "invalid",
				toType: Integer,
			},
			want:    int64(-1),
			wantErr: newError(ErrCodeParsingFailed),
		},
		{
			name: "valid_float_1",
			args: args{
				value:  "1.2",
				toType: Float,
			},
			want:    float64(1.2),
			wantErr: nil,
		},
		{
			name: "valid_float_2",
			args: args{
				value:  "-1.2",
				toType: Float,
			},
			want:    float64(-1.2),
			wantErr: nil,
		},
		{
			name: "invalid_float",
			args: args{
				value:  "invalid",
				toType: Float,
			},
			want:    float64(-1.2),
			wantErr: newError(ErrCodeParsingFailed),
		},
		{
			name: "valid_string",
			args: args{
				value:  "asdf",
				toType: String,
			},
			want:    "asdf",
			wantErr: nil,
		},
		{
			name: "invalid_valueType",
			args: args{
				value:  "invalid",
				toType: ValueType(0),
			},
			want:    "invalid",
			wantErr: newError(ErrCodeInvalidValueType),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := parseValue(tt.args.value, tt.args.toType)

			if tt.wantErr == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseValue() = %v, want %v", got, tt.want)
			}
			if !isErrorEqual(gotErr, tt.wantErr) {
				t.Errorf("parseValue() gotErr %v, wantErr %v", gotErr, tt.wantErr)
				return
			}

		})
	}
}
