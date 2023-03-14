package ruleenginecore

import "testing"

func TestFields_exist(t *testing.T) {
	type args struct {
		fieldName         string
		expectedFieldType ValueType
	}
	tests := []struct {
		name string
		fs   Fields
		args args
		want bool
	}{
		{
			name: "valid_IntegerField",
			fs: Fields{
				"test": Integer,
			},
			args: args{
				fieldName:         "test",
				expectedFieldType: Integer,
			},
			want: true,
		},
		{
			name: "invalid_FieldNotExist",
			fs: Fields{
				"test": Integer,
			},
			args: args{
				fieldName:         "invalid",
				expectedFieldType: Integer,
			},
			want: false,
		},
		{
			name: "invalid_WrongFieldType",
			fs: Fields{
				"test": Integer,
			},
			args: args{
				fieldName:         "test",
				expectedFieldType: Float,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fs.exist(tt.args.fieldName, tt.args.expectedFieldType); got != tt.want {
				t.Errorf("Fields.exist() = %v, want %v", got, tt.want)
			}
		})
	}
}
