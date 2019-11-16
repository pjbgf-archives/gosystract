package main

import (
	"fmt"
	"os"

	cli "github.com/pjbgf/gosystract/cmd/cli"
	"github.com/pjbgf/gosystract/cmd/systract"
)

var onError func(error) = func(err error) {
	fmt.Printf("\nerror: %s\n", err)
	os.Exit(1)
}

func main() {
	cli.Run(os.Stdout, os.Args, systract.Extract, onError)
}
