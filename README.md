# go-ini

[![GoDoc](https://godoc.org/github.com/subpop/go-ini?status.svg)](https://godoc.org/github.com/subpop/go-ini)
[![Build Status](https://travis-ci.org/subpop/go-ini.svg?branch=master)](https://travis-ci.org/subpop/go-ini)
[![Go Report Card](https://goreportcard.com/badge/github.com/subpop/go-ini)](https://goreportcard.com/report/github.com/subpop/go-ini)


A Go package that encodes and decodes INI-files.

# Usage

```go
data := `[settings]
username=root
password=swordfish
`

var config struct {
    Settings struct {
        Username string `ini:"username"`
        Password string `ini:"password"`
    } `ini:"settings"`
}

if err := ini.Unmarshal(data, &config); err != nil {
    fmt.Println(err)
}
fmt.Println(config)
```