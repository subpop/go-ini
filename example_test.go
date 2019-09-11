package ini_test

import (
	"fmt"
	"log"

	"github.com/subpop/go-ini"
)

func ExampleUnmarshal() {
	type GitConfig struct {
		User struct {
			Email string `ini:"email"`
			Name  string `ini:"name"`
		} `ini:"user"`
	}

	gitconfig := `
	[user]
		email = gopher@golang.org
		name = Gopher
	`

	var gitCfg GitConfig
	if err := ini.Unmarshal([]byte(gitconfig), &gitCfg); err != nil {
		log.Fatal(err)
	}
	fmt.Println(gitCfg.User)

	// Output: {gopher@golang.org Gopher}
}
