package main

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/golang-collections/collections/stack"
	"github.com/pjbgf/dump/pkg"
	"github.com/sirupsen/logrus"
)

var dumpFileName string
var logger *logrus.Logger
var preValidation chan string

const (
	callCaptureRegex         string = ".+CALL.(\b([a-zA-Z0-9_.\\/]|\\.|\\(\\*[a-z]+\\))+\b)+"
	callsFromEntryPointRegex string = "main.go" + callCaptureRegex
	goFileNameRegex          string = "[a-zA-Z0-9_]+.go"
)

func init() {
	logger = logrus.New()
	logger.SetLevel(logrus.InfoLevel)
}

type RESULT struct {
	symbolName string
	subCalls   []string
	sysCalls   []uint8
}

func main() {
	dumpFileName = "caps-dump"

	// symbolCalls := []string{"main.init.0"}
	symbolCalls := []string{"main.init.0", "main.init.1", "main.main"}

	getInitialSymbols := func(done <-chan bool, symbols ...string) <-chan RESULT {
		symbolStream := make(chan RESULT)
		go func() {
			defer close(symbolStream)
			for _, symbol := range symbols {
				logger.Debugf("getInitialSymbols [%s]", symbol)
				select {
				case <-done:
					return
				case symbolStream <- RESULT{subCalls: []string{symbol}}:
				}
			}
		}()
		return symbolStream
	}

	processedSymbols := make(map[string]bool)
	unique := func(done <-chan bool, symbolStream <-chan RESULT) <-chan RESULT {
		uniqueStream := make(chan RESULT)
		go func() {
			defer close(uniqueStream)
			exists := func(name string) bool {
				if _, exists := processedSymbols[name]; !exists {
					processedSymbols[name] = true
					return false
				}
				return true
			}

			for result := range symbolStream {
				for _, symbolCall := range result.subCalls {
					if !exists(symbolCall) {
						logger.Debugf("unique [%s] %s", result.symbolName, symbolCall)
						select {
						case <-done:
							return
						case uniqueStream <- getSyscallsAndSubCalls(symbolCall):
						}
					} else {
						logger.Debugf("ignore: %s", symbolCall)
					}
				}
			}
		}()
		return uniqueStream
	}

	processSymbolCalls := func(done <-chan bool, symbolStream <-chan RESULT) <-chan RESULT {
		processedStream := make(chan RESULT)
		go func() {
			defer close(processedStream)
			for result := range symbolStream {
				for _, symbolCall := range result.subCalls {
					logger.Debugf("processSymbolCalls [%s] %s", result.symbolName, symbolCall)

					select {
					case <-done:
						return
					case processedStream <- getSyscallsAndSubCalls(symbolCall):
					}
				}
			}
		}()
		return processedStream
	}

	loop := func(done <-chan bool, symbolStream <-chan RESULT) <-chan RESULT {
		loopStream := make(chan RESULT)
		go func() {
			defer close(loopStream)
			for result := range symbolStream {
				logger.Debugf("loop [%s] len: %d", result.symbolName, len(result.subCalls))
				loopThrough(done, result, loopStream)
			}
		}()
		return loopStream
	}

	AAAAAAAA = make(map[string]bool)
	start := time.Now()
	done := make(chan bool)
	defer close(done)

	symbolStream := getInitialSymbols(done, symbolCalls...)
	result := loop(done, unique(done, processSymbolCalls(done, symbolStream)))
	for p := range result {
		for _, subCall := range p.subCalls {
			logger.Debugf("[%s] subcall: %s", p.symbolName, subCall)
		}
		for _, id := range p.sysCalls {
			logger.Infof("[%s] syscall ID: %d", p.symbolName, id)
		}
	}

	diff := start.Sub(time.Now())
	logger.Infof("Total seconds: %d", diff.Seconds)
}

var AAAAAAAA map[string]bool

func loopThrough(done <-chan bool, result RESULT, c chan RESULT) {
	for _, symbolCall := range result.subCalls {

		if _, exists := AAAAAAAA[symbolCall]; !exists {
			AAAAAAAA[symbolCall] = true

			subResult := getSyscallsAndSubCalls(symbolCall)

			select {
			case <-done:
				return
			case c <- subResult:
			}

			if len(subResult.subCalls) > 0 {
				loopThrough(done, subResult, c)
			}
		}
	}
}

func getSyscallsAndSubCalls(symbolName string) RESULT {
	file, err := os.Open(dumpFileName)
	if err != nil {
		panic(file)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	subCalls := make([]string, 0)
	sysCallIDs := make([]uint8, 0)
	processedSymbols := make(map[string]bool)
	for scanner.Scan() {
		line := scanner.Text()
		if pkg.IsSymbolDefinition(line, symbolName) {
			logger.Debugf("-- %s\n", symbolName)
			stack := stack.New()
			for scanner.Scan() {
				line = scanner.Text()
				if line == "" || line == "\n" || strings.Contains(line, "TEXT ") {
					return RESULT{symbolName: symbolName, subCalls: subCalls, sysCalls: sysCallIDs}
				}

				if id, found := pkg.GetSyscallID(line); found {
					stack.Push(id)
				}

				if pkg.ContainsSyscall(line) {
					val1 := stack.Pop()
					val2 := stack.Pop()

					syscallID := val1
					if val2 != nil {
						syscallID = val2
					}
					sysCallIDs = append(sysCallIDs, syscallID.(uint8))
					// logger.Infof("[%s] FOUND SYSCALL %d", symbolName, syscallID)
				} else {
					target, found := pkg.GetCallTarget(line)
					if found {

						if _, exists := processedSymbols[target]; !exists {
							processedSymbols[target] = true
							subCalls = append(subCalls, target)
						}
					}
				}
			}
		}
	}
	return RESULT{symbolName: symbolName, subCalls: subCalls, sysCalls: sysCallIDs}
}

func getSourceFileName(symbolName string) (string, error) {
	file, err := os.Open(dumpFileName)
	if err != nil {
		panic(file)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		line := scanner.Text()
		re := regexp.MustCompile(symbolName)
		captures := re.FindStringSubmatch(line)

		if captures != nil && len(captures) > 0 {
			if scanner.Scan() {
				re = regexp.MustCompile(goFileNameRegex)
				captures = re.FindStringSubmatch(scanner.Text())

				if captures != nil && len(captures) > 0 {
					return captures[0], nil
				}
			}
		}
	}

	return "", errors.New("no source file name found")
}

func getMatchingLines(regex string) []string {
	file, err := os.Open(dumpFileName)
	if err != nil {
		panic(file)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	lines := make([]string, 0)
	for scanner.Scan() {

		line := scanner.Text()
		re := regexp.MustCompile(regex)
		captures := re.FindStringSubmatch(line)

		if captures != nil && len(captures) > 0 {
			lines = append(lines, captures[1])
		}
	}
	return lines
}

// TEXT golang.org/x/sys/unix.Capget(SB) /media/pjb/src/go/pkg/mod/golang.org/x/sys@v0.0.0-20191027211539-f8518d3b3627/unix/zsyscall_linux_amd64.go
//   zsyscall_linux_amd64.go:441	0x48bd40		64488b0c25f8ffffff	MOVQ FS:0xfffffff8, CX
//   zsyscall_linux_amd64.go:441	0x48bd49		483b6110		CMPQ 0x10(CX), SP
//   zsyscall_linux_amd64.go:441	0x48bd4d		0f86ca000000		JBE 0x48be1d
//   zsyscall_linux_amd64.go:441	0x48bd53		4883ec50		SUBQ $0x50, SP
//   zsyscall_linux_amd64.go:441	0x48bd57		48896c2448		MOVQ BP, 0x48(SP)
//   zsyscall_linux_amd64.go:441	0x48bd5c		488d6c2448		LEAQ 0x48(SP), BP
//   zsyscall_linux_amd64.go:442	0x48bd61		488b442458		MOVQ 0x58(SP), AX
//   zsyscall_linux_amd64.go:442	0x48bd66		4889442440		MOVQ AX, 0x40(SP)
//   zsyscall_linux_amd64.go:442	0x48bd6b		488b442460		MOVQ 0x60(SP), AX
//   zsyscall_linux_amd64.go:442	0x48bd70		4889442438		MOVQ AX, 0x38(SP)
//   zsyscall_linux_amd64.go:442	0x48bd75		48c704247d000000	MOVQ $0x7d, 0(SP)
//   zsyscall_linux_amd64.go:442	0x48bd7d		488b442440		MOVQ 0x40(SP), AX
//   zsyscall_linux_amd64.go:442	0x48bd82		4889442408		MOVQ AX, 0x8(SP)
//   zsyscall_linux_amd64.go:442	0x48bd87		488b442438		MOVQ 0x38(SP), AX
//   zsyscall_linux_amd64.go:442	0x48bd8c		4889442410		MOVQ AX, 0x10(SP)
//   zsyscall_linux_amd64.go:442	0x48bd91		48c744241800000000	MOVQ $0x0, 0x18(SP)
//   zsyscall_linux_amd64.go:442	0x48bd9a		e881030000		CALL golang.org/x/sys/unix.Syscall(SB)
