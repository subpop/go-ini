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
}

func ExampleUnmarshal() {
	type Database struct {
		Server string
		Port   int
		File   string
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
	File=payroll.dat`)

	if err := ini.Unmarshal(data, &config); err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(config)
	// Output:
	// {1.2.3 {John Doe Acme Widgets Inc.} {192.0.2.62 143 payroll.dat}}
}
