package ini_test

import (
	"fmt"
	"os"

	"github.com/subpop/go-ini"
)

func ExampleMarshal() {
	type Database struct {
		Server string
		Port   int
		File   string
		Path   map[string]string
	}

	type Person struct {
		Name         string
		Organization string
	}

	type Config struct {
		Version  string
		Owner    Person
		Database Database
	}

	config := Config{
		Version: "1.2.3",
		Owner: Person{
			Name:         "John Doe",
			Organization: "Acme Widgets Inc.",
		},
		Database: Database{
			Server: "192.0.2.62",
			Port:   143,
			File:   "payroll.dat",
			Path:   map[string]string{"unix": "/var/db"},
		},
	}

	b, err := ini.Marshal(config)
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
	// Output:
	// Version=1.2.3
	//
	// [Owner]
	// Name=John Doe
	// Organization=Acme Widgets Inc.
	//
	// [Database]
	// Server=192.0.2.62
	// Port=143
	// File=payroll.dat
	// Path[unix]=/var/db
}

func ExampleUnmarshal() {
	type Database struct {
		Server string
		Port   int
		File   string
		Path   map[string]string
	}

	type Person struct {
		Name         string
		Organization string
	}

	type Config struct {
		Version  string
		Owner    Person
		Database Database
	}

	var config Config

	data := []byte(`Version=1.2.3

	[Owner]
	Name=John Doe
	Organization=Acme Widgets Inc.
	
	[Database]
	Server=192.0.2.62
	Port=143
	File=payroll.dat
	Path[unix]=/var/db
	Path[win32]=C:\db`)

	if err := ini.Unmarshal(data, &config); err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(config)
	// Output:
	// {1.2.3 {John Doe Acme Widgets Inc.} {192.0.2.62 143 payroll.dat map[unix:/var/db win32:C:\db]}}
}
