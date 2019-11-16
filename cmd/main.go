package main

import (
	"fmt"
	"os"

	cli "github.com/pjbgf/gosystract/cmd/cli"
	"github.com/pjbgf/gosystract/cmd/systract"
)

func main() {
	err := cli.Run(os.Stdout, os.Args, systract.Extract)
	if err != nil {
		fmt.Printf("\nerror: %s\n", err)
		os.Exit(1)
	}
}
