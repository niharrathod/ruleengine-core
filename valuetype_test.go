package ruleenginecore

import (
	"reflect"
	"testing"
)

func TestValueType_isValid(t *testing.T) {
	tests := []struct {
		name      string
		valueType ValueType
		want      bool
	}{
		{
			name:      "Valid_ValueType_Boolean",
			valueType: Boolean,
			want:      true,
		},
		{
			name:      "Valid_ValueType_Integer",
			valueType: Integer,
			want:      true,
		},
		{
			name:      "Valid_ValueType_Float",
			valueType: Float,
			want:      true,
		},
		{
			name:      "Valid_ValueType_String",
			valueType: String,
			want:      true,
		},
		{
			name:      "Invalid_ValueType",
			valueType: unknownValueType,
			want:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.valueType.isValid(); got != tt.want {
				t.Errorf("ValueType.isValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueType_String(t *testing.T) {
	tests := []struct {
		name      string
		valueType ValueType
		want      string
	}{
		{
			name:      "ValueType_Boolean",
			valueType: Boolean,
			want:      "Boolean",
		},
		{
			name:      "ValueType_Integer",
			valueType: Integer,
			want:      "Integer",
		},
		{
			name:      "ValueType_Float",
			valueType: Float,
			want:      "Float",
		},
		{
			name:      "ValueType_String",
			valueType: String,
			want:      "String",
		},
		{
			name:      "ValueType_Unknown",
			valueType: unknownValueType,
			want:      "ValueType(0)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.valueType.String(); got != tt.want {
				t.Errorf("ValueType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseValueType(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    ValueType
		wantErr bool
	}{
		{
			name: "valid_bool",
			args: args{
				s: "bool",
			},
			want:    Boolean,
			wantErr: false,
		},
		{
			name: "valid_Bool",
			args: args{
				s: "Bool",
			},
			want:    Boolean,
			wantErr: false,
		},
		{
			name: "valid_boolean",
			args: args{
				s: "boolean",
			},
			want:    Boolean,
			wantErr: false,
		},
		{
			name: "valid_Boolean",
			args: args{
				s: "Boolean",
			},
			want:    Boolean,
			wantErr: false,
		},
		{
			name: "valid_string",
			args: args{
				s: "string",
			},
			want:    String,
			wantErr: false,
		},
		{
			name: "valid_String",
			args: args{
				s: "String",
			},
			want:    String,
			wantErr: false,
		},
		{
			name: "valid_int",
			args: args{
				s: "int",
			},
			want:    Integer,
			wantErr: false,
		},
		{
			name: "valid_Int",
			args: args{
				s: "Int",
			},
			want:    Integer,
			wantErr: false,
		},
		{
			name: "valid_integer",
			args: args{
				s: "integer",
			},
			want:    Integer,
			wantErr: false,
		},
		{
			name: "valid_Integer",
			args: args{
				s: "Integer",
			},
			want:    Integer,
			wantErr: false,
		},
		{
			name: "valid_float",
			args: args{
				s: "float",
			},
			want:    Float,
			wantErr: false,
		},
		{
			name: "valid_Float",
			args: args{
				s: "Float",
			},
			want:    Float,
			wantErr: false,
		},
		{
			name: "invalid_valueType",
			args: args{
				s: "asdf",
			},
			want:    unknownValueType,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseValueType(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseValueType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseValueType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueType_MarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		valueType ValueType
		want      []byte
		wantErr   bool
	}{
		{
			name:      "valid_boolean",
			valueType: Boolean,
			want:      []byte("\"Boolean\""),
			wantErr:   false,
		},
		{
			name:      "valid_integer",
			valueType: Integer,
			want:      []byte("\"Integer\""),
			wantErr:   false,
		},
		{
			name:      "valid_float",
			valueType: Float,
			want:      []byte("\"Float\""),
			wantErr:   false,
		},
		{
			name:      "valid_string",
			valueType: String,
			want:      []byte("\"String\""),
			wantErr:   false,
		},
		{
			name:      "invalid_marshal",
			valueType: unknownValueType,
			want:      nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.valueType.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValueType.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValueType.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueType_UnmarshalJSON(t *testing.T) {
	var tempValueType = unknownValueType
	type args struct {
		data []byte
	}
	tests := []struct {
		name      string
		valueType *ValueType
		args      args
		wantErr   bool
		want      ValueType
	}{
		{
			name:      "invalid_value",
			valueType: nil,
			args: args{
				[]byte("unknown"),
			},
			wantErr: true,
		},
		{
			name:      "valid_boolean",
			valueType: &tempValueType,
			args: args{
				[]byte("\"Bool\""),
			},
			wantErr: false,
			want:    Boolean,
		},
		{
			name:      "valid_integer",
			valueType: &tempValueType,
			args: args{
				[]byte("\"Int\""),
			},
			wantErr: false,
			want:    Integer,
		},
		{
			name:      "valid_float",
			valueType: &tempValueType,
			args: args{
				[]byte("\"Float\""),
			},
			wantErr: false,
			want:    Float,
		},
		{
			name:      "valid_string",
			valueType: &tempValueType,
			args: args{
				[]byte("\"String\""),
			},
			wantErr: false,
			want:    String,
		},
		{
			name:      "valid_string",
			valueType: &tempValueType,
			args: args{
				[]byte("\"Unknown\""),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.valueType.UnmarshalJSON(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValueType.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && *tt.valueType != tt.want {
				t.Errorf("ValueType.UnmarshalJSON() failed want %v got %v", tt.want, tt.valueType)
			}
		})
	}
}

func TestValueType_isBoolean(t *testing.T) {
	tests := []struct {
		name      string
		valueType ValueType
		want      bool
	}{
		{
			name:      "valid_boolean",
			valueType: Boolean,
			want:      true,
		},
		{
			name:      "invalid_valuetype",
			valueType: unknownValueType,
			want:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.valueType.isBoolean(); got != tt.want {
				t.Errorf("ValueType.isBoolean() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueType_isInteger(t *testing.T) {
	tests := []struct {
		name      string
		valueType ValueType
		want      bool
	}{
		{
			name:      "valid_Integer",
			valueType: Integer,
			want:      true,
		},
		{
			name:      "invalid_valuetype",
			valueType: unknownValueType,
			want:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.valueType.isInteger(); got != tt.want {
				t.Errorf("ValueType.isInteger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueType_isFloat(t *testing.T) {
	tests := []struct {
		name      string
		valueType ValueType
		want      bool
	}{
		{
			name:      "valid_Float",
			valueType: Float,
			want:      true,
		},
		{
			name:      "invalid_valuetype",
			valueType: unknownValueType,
			want:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.valueType.isFloat(); got != tt.want {
				t.Errorf("ValueType.isFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueType_isString(t *testing.T) {
	tests := []struct {
		name      string
		valueType ValueType
		want      bool
	}{
		{
			name:      "valid_String",
			valueType: String,
			want:      true,
		},
		{
			name:      "invalid_valuetype",
			valueType: unknownValueType,
			want:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.valueType.isString(); got != tt.want {
				t.Errorf("ValueType.isString() = %v, want %v", got, tt.want)
			}
		})
	}
}
