package main

import (
	"html/template"
	"os"

	cli "github.com/pjbgf/gosystract/cmd/cli"
	"github.com/pjbgf/gosystract/cmd/systract"
)

var (
	noSyscallsFoundMessage string = "no systems calls were found"
	usageMessage           string = `gosystract returns the names and IDs of all system calls being called inside a go application.
Usage: 
	gosystrac goapp.dump

To generate a dump file from a go application use: 
	go tool objdump goapp > goapp.dump`
)

func main() {
	err := cli.Run(os.Stdout, os.Args, systract.Extract)
	if err != nil {
		panic(err)
	}

	fileName := os.Args[1]
	syscalls, err := systract.Extract(fileName)
	if err != nil {
		panic(err)
	}

	var t *template.Template

	if len(os.Args) > 2 {
		t = template.Must(template.New("result.tmpl").Parse(os.Args[2]))
	} else {
		t = template.Must(template.ParseFiles("result.tmpl"))
	}

	err = t.Execute(os.Stdout, syscalls)
	if err != nil {
		panic(err)
	}
}
