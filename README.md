# gosystract
gosystract returns the names and IDs of all system calls being called inside a go application.

[![codecov](https://codecov.io/gh/pjbgf/gosystract/branch/master/graph/badge.svg?token=v9BeEO6F0S)](https://codecov.io/gh/pjbgf/gosystract)
[![Dependabot Status](https://api.dependabot.com/badges/status?host=github&repo=pjbgf/gosystract)](https://dependabot.com)
[![GoReport](https://goreportcard.com/badge/github.com/pjbgf/gosystract)](https://goreportcard.com/badge/github.com/pjbgf/gosystract)
[![GoDoc](https://godoc.org/github.com/pjbgf/gosystract?status.svg)](https://godoc.org/github.com/pjbgf/gosystract)
[![Docker Pulls](https://img.shields.io/docker/pulls/paulinhu/gosystract.svg)](https://hub.docker.com/r/paulinhu/gosystract/tags)
![build](https://github.com/pjbgf/gosystract/workflows/go/badge.svg)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](http://choosealicense.com/licenses/mit/)


## Command-line Usage:

Running the sample dump file:
```console
$ go run main.go test/keyring.dump

16 syscalls found:
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

To generate a dump file from a go application use: 
```console
$ go tool objdump goapp > goapp.dump
```

## Using it programatically

```golang
import "github.com/pjbgf/gosystract/cmd/systract"

func printSyscalls() {
	syscalls, err := systract.Extract("goapp.dump")
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