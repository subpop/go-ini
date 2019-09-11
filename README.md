# About

`go-ini` is an INI-file codec package for Go. It mimics the Marshal/Unmarshal
pattern set forth by the `encoding/json` and `encoding/xml` packages. It borrows
heavily from the `encoding/json` package implementation.

# Usage

```go
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
	if err := Unmarshal([]byte(gitconfig), &gitCfg); err != nil {
		log.Fatal(err)
	}
	fmt.Println(gitCfg.User)
```