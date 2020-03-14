package ini

import (
	"reflect"
	"testing"
)

func TestNewTag(t *testing.T) {
	tests := []struct {
		input reflect.StructField
		want  tag
	}{
		{
			input: reflect.StructField{
				Name: "Field",
				Tag:  reflect.StructTag(`ini:"Field"`),
			},
			want: tag{
				name:      "Field",
				omitempty: false,
			},
		},
		{
			input: reflect.StructField{
				Name: "Field",
				Tag:  reflect.StructTag(`ini:"Field,omitempty"`),
			},
			want: tag{
				name:      "Field",
				omitempty: true,
			},
		},
	}

	for _, test := range tests {
		got := newTag(test.input)

		if got != test.want {
			t.Errorf("%+v != %+v", got, test.want)
		}
	}
}
