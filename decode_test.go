package ini

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		input []byte
		want  config
	}{
		{
			input: []byte(`
version=1.2.3

[owner]
name=John Doe
organization=Acme Widgets Inc.

[database]
server=192.0.2.62
port=143
file="payroll.dat"
enabled=true`),
			want: config{
				Version: "1.2.3",
				Owner: owner{
					Name:         "John Doe",
					Organization: "Acme Widgets Inc.",
				},
				Database: database{
					Server:  "192.0.2.62",
					Port:    143,
					File:    "\"payroll.dat\"",
					Enabled: true,
				},
			},
		},
	}

	for _, test := range tests {
		var got config
		if err := Unmarshal(test.input, &got); err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(got, test.want) {
			t.Errorf("%v != %v", got, test.want)
		}
	}
}
