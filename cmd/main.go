package main

import (
	"os"

	cli "github.com/pjbgf/gosystract/cmd/cli"
	"github.com/pjbgf/gosystract/cmd/systract"
)

var onError func(error) = func(err error) {
	defer func() {
		if rec := recover(); rec != nil {
			os.Exit(1)
		}
	}()

	panic(err)
}

func main() {
	cli.Run(os.Stdout, os.Args, systract.Extract, onError)
}
