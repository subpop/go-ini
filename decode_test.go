package ini

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUnmarshal(t *testing.T) {
	type database struct {
		Server    string
		Port      int
		Encrypted bool
	}

	tests := []struct {
		input []byte
		want  struct {
			Database database
		}
	}{
		{
			input: []byte("[Database]\nServer=192.0.2.62\nPort=143\nEncrypted=false"),
			want: struct {
				Database database
			}{
				Database: database{
					Server:    "192.0.2.62",
					Port:      143,
					Encrypted: false,
				},
			},
		},
	}

	for _, test := range tests {
		var got struct {
			Database database
		}
		if err := Unmarshal(test.input, &got); err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(got, test.want) {
			t.Errorf("%v != %v", got, test.want)
		}
	}
}
