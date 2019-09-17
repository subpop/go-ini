package ini

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type owner struct {
	Name         string `ini:"name"`
	Organization string `ini:"organization"`
}
type database struct {
	Server  string `ini:"server"`
	Port    int    `ini:"port"`
	File    string `ini:"file"`
	Enabled bool   `ini:"enabled,omitempty"`
}

type config struct {
	Version  string   `ini:"version,omitempty"`
	Owner    owner    `ini:"owner"`
	Database database `ini:"database"`
}

func TestMarshal(t *testing.T) {
	tests := []struct {
		input config
		want  []byte
	}{
		{
			input: config{
				Version: "1.2.3",
				Owner: owner{
					Name:         "John Doe",
					Organization: "Acme Widgets Inc.",
				},
				Database: database{
					Server: "192.0.2.62",
					Port:   143,
					File:   "\"payroll.dat\"",
				},
			},
			want: []byte(`version=1.2.3

[owner]
name=John Doe
organization=Acme Widgets Inc.

[database]
server=192.0.2.62
port=143
file="payroll.dat"`),
		},
		{
			input: config{
				Owner: owner{
					Name:         "John Doe",
					Organization: "Acme Widgets Inc.",
				},
				Database: database{
					Server: "192.0.2.62",
					Port:   143,
					File:   "\"payroll.dat\"",
				},
			},
			want: []byte(`[owner]
name=John Doe
organization=Acme Widgets Inc.

[database]
server=192.0.2.62
port=143
file="payroll.dat"`),
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
