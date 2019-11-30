package main

import (
	"os"

	cli "github.com/pjbgf/gosystract/cmd/cli"
	"github.com/pjbgf/gosystract/cmd/systract"
)

func main() {
	cli.Run(os.Stdout, os.Stderr, os.Args, systract.Extract, os.Exit)
}
