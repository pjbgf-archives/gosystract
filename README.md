# gosystract
gosystract extracts all system calls that may be called within the execution path of a go application.

[![codecov](https://codecov.io/gh/pjbgf/gosystract/branch/master/graph/badge.svg?token=hPDXVgD96x)](https://codecov.io/gh/pjbgf/gosystract)
[![Dependabot Status](https://api.dependabot.com/badges/status?host=github&repo=pjbgf/gosystract)](https://dependabot.com)
[![GoReport](https://goreportcard.com/badge/github.com/pjbgf/gosystract)](https://goreportcard.com/report/github.com/pjbgf/gosystract)
[![GoDoc](https://godoc.org/github.com/pjbgf/gosystract?status.svg)](https://godoc.org/github.com/pjbgf/gosystract)
![build](https://github.com/pjbgf/gosystract/workflows/go/badge.svg)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](http://choosealicense.com/licenses/mit/)


## Installation:

### using container image
```console
docker run --rm -it paulinhu/gosystract gosystract
``` 

### using go environment
```console
go install github.com/pjbgf/gosystract
``` 

If you don't have $GOPATH/bin in your $PATH, prefix the command with:

`PATH=$PATH:$GOPATH/bin gosystract`

> Note that gosystract has a dependency to go tools when working against executable files. On those scenarios, ensure they are in your $PATH.

## Command-line Usage:

Syntax
```console
Usage:

	gosystrac [flags] filePath

Flags:
    --dumpfile, -d    Handles a dump file instead of go executable.
    --template        Defines a go template for the results.
                      Example: --template="{{- range . }}{{printf "%d - %s\n" .ID .Name}}{{- end}}"
```

Running against gosystract itself:
```console
$ gosystract $(which gosystract)

14 system calls found:
    sched_yield (24)
    futex (202)
    write (1)
    rt_sigprocmask (14)
    getpid (39)
    gettid (186)
    tgkill (234)
    rt_sigaction (13)
    exit_group (231)
    mmap (9)
    madvise (28)
    read (0)
    getpgrp (111)
    arch_prctl (158)
```

Running the sample dump file:
```console
$ gosystract --dump test/keyring.dump

16 system calls found:
    read (0)
    sched_yield (24)
    futex (202)
    write (1)
    rt_sigprocmask (14)
    getpid (39)
    gettid (186)
    tgkill (234)
    rt_sigaction (13)
    exit_group (231)
    mmap (9)
    madvise (28)
    getpgrp (111)
    arch_prctl (158)
    add_key (248)
    keyctl (250)
```

To generate a dump file from a go application use the go tool objdump: 
```console
$ go tool objdump goapp > goapp.dump
```

## Using it programmatically

```golang
package main

import "github.com/pjbgf/gosystract/cmd/systract"

func main() {
	source := systract.NewExeReader("goapp")
	syscalls, err := systract.Extract(source)
	if err != nil {
		panic(err)
	}

    for _, syscall := range syscalls {
        fmt.Printf("%s (%d)\n", syscall.Name, syscall.ID)
    }
}
```

## License

This application is licensed under the MIT License, you may obtain a copy of it [here](LICENSE).