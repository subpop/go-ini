package ini

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMarshal(t *testing.T) {
	type database struct {
		Server    string
		Port      int
		Encrypted bool
	}

	tests := []struct {
		input struct {
			Database database
		}
		want []byte
	}{
		{
			input: struct {
				Database database
			}{
				Database: database{
					Server:    "192.0.2.62",
					Port:      143,
					Encrypted: false,
				},
			},
			want: []byte("[Database]\nServer=192.0.2.62\nPort=143\nEncrypted=false"),
		},
	}

	for _, test := range tests {
		got, err := Marshal(test.input)
		if err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(got, test.want) {
			t.Errorf("%q != %q", string(got), string(test.want))
		}
	}
}
