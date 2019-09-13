package ini_test

import (
	"fmt"
	"log"

	"github.com/subpop/go-ini"
)

type user struct {
	Email string `ini:"email"`
	Name  string `ini:"name"`
}

type gitConfig struct {
	User user `ini:"user"`
}

func ExampleUnmarshal() {
	gitconfig := `
	[user]
		email = gopher@golang.org
		name = Gopher
	`

	var gitCfg gitConfig
	if err := ini.Unmarshal([]byte(gitconfig), &gitCfg); err != nil {
		log.Fatal(err)
	}
	fmt.Println(gitCfg.User)

	// Output: {gopher@golang.org Gopher}
}

func ExampleMarshal() {
	gitCfg := gitConfig{
		User: user{
			Name:  "Gopher",
			Email: "gopher@golang.org",
		},
	}

	data, err := ini.Marshal(&gitCfg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))

	// Output:
	// [user]
	// email=gopher@golang.org
	// name=Gopher
}
