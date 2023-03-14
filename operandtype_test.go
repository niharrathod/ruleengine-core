package ruleenginecore

import (
	"reflect"
	"testing"
)

func TestOperandType_isValid(t *testing.T) {
	tests := []struct {
		name        string
		operandType OperandType
		want        bool
	}{
		{
			name:        "valid_Field",
			operandType: Field,
			want:        true,
		},
		{
			name:        "valid_Constant",
			operandType: Constant,
			want:        true,
		},
		{
			name:        "invalid",
			operandType: unknownOperandType,
			want:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.operandType.isValid(); got != tt.want {
				t.Errorf("OperandType.isValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperandType_String(t *testing.T) {
	tests := []struct {
		name        string
		operandType OperandType
		want        string
	}{
		{
			name:        "valid_field",
			operandType: Field,
			want:        "Field",
		},
		{
			name:        "valid_constant",
			operandType: Constant,
			want:        "Constant",
		},
		{
			name:        "invalid",
			operandType: unknownOperandType,
			want:        "OperandType(0)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.operandType.String(); got != tt.want {
				t.Errorf("OperandType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseOperandType(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    OperandType
		wantErr bool
	}{
		{
			name: "valid_field",
			args: args{
				s: "field",
			},
			want:    Field,
			wantErr: false,
		},
		{
			name: "valid_Field",
			args: args{
				s: "Field",
			},
			want:    Field,
			wantErr: false,
		},
		{
			name: "valid_constant",
			args: args{
				s: "constant",
			},
			want:    Constant,
			wantErr: false,
		},
		{
			name: "valid_Constant",
			args: args{
				s: "Constant",
			},
			want:    Constant,
			wantErr: false,
		},
		{
			name: "invalid",
			args: args{
				s: "Unknown",
			},
			want:    unknownOperandType,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseOperandType(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseOperandType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseOperandType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperandType_MarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		operandType OperandType
		want        []byte
		wantErr     bool
	}{
		{
			name:        "valid_Field",
			operandType: Field,
			want:        []byte("\"Field\""),
			wantErr:     false,
		},
		{
			name:        "valid_Constant",
			operandType: Constant,
			want:        []byte("\"Constant\""),
			wantErr:     false,
		},
		{
			name:        "invalid",
			operandType: unknownOperandType,
			want:        nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.operandType.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("OperandType.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OperandType.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperandType_UnmarshalJSON(t *testing.T) {
	var tempOperandType = unknownOperandType
	type args struct {
		data []byte
	}
	tests := []struct {
		name        string
		operandType *OperandType
		args        args
		wantErr     bool
		want        OperandType
	}{
		{
			name:        "valid_field",
			operandType: &tempOperandType,
			args: args{
				data: []byte("\"field\""),
			},
			wantErr: false,
			want:    Field,
		},
		{
			name:        "valid_Field",
			operandType: &tempOperandType,
			args: args{
				data: []byte("\"Field\""),
			},
			wantErr: false,
			want:    Field,
		},
		{
			name:        "valid_constant",
			operandType: &tempOperandType,
			args: args{
				data: []byte("\"constant\""),
			},
			wantErr: false,
			want:    Constant,
		},
		{
			name:        "valid_Constant",
			operandType: &tempOperandType,
			args: args{
				data: []byte("\"Constant\""),
			},
			wantErr: false,
			want:    Constant,
		},
		{
			name:        "invalid_string",
			operandType: &tempOperandType,
			args: args{
				data: []byte("unknown"),
			},
			wantErr: true,
		},
		{
			name:        "invalid_string",
			operandType: &tempOperandType,
			args: args{
				data: []byte("\"Unknown\""),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.operandType.UnmarshalJSON(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("OperandType.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && *tt.operandType != tt.want {
				t.Errorf("OperandType.UnmarshalJSON() got = %v, want %v", *tt.operandType, tt.want)
			}

		})
	}
}
