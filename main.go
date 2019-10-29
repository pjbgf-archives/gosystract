package main

import (
	"fmt"
	"os"

	"github.com/pjbgf/gosystract/cmd/systract"
)

var (
	noSyscallsFoundMessage string = "no syscalls found in %s\n"
	usageMessage           string = `gosystract returns the names and IDs of all system calls being called inside a go application.
Usage: 
	gosystrac goapp.dump

To generate a dump file from a go application use: 
	go tool objdump goapp > goapp.dump`
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println(usageMessage)
		os.Exit(1)
	}

	fileName := os.Args[1]
	syscalls, err := systract.Extract(fileName)
	if err != nil {
		panic(err)
	}

	if len(syscalls) == 0 {
		fmt.Printf(noSyscallsFoundMessage, fileName)
	} else {
		fmt.Printf("%d syscalls found:\n\n", len(syscalls))
		for _, syscall := range syscalls {
			fmt.Printf("%s (%d)\n", syscall.Name, syscall.ID)
		}
	}
}
