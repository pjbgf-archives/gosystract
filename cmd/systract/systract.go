// Package systract provides libraries to extract syscalls from go applications programmatically.
package systract

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
	"sync"

	"github.com/golang-collections/collections/stack"
	"github.com/pkg/errors"
)

var executableSymbolCalls = []string{"main.main", "main.init.0", "main.init.1"}

const (
	symbolDefinitionRegex string = "TEXT.((\\%|\\(|\\)|\\*|[a-zA-Z0-9_.\\/])+)\\b\\("
	syscallHexIDRegex     string = "MOV(Q|L).\\$0x([0-9a-fA-F]+)"
	callCaptureRegex      string = ".+CALL.(\\b([a-zA-Z0-9_.\\/]|\\.|\\(\\*[a-zA-Z0-9_.\\/]+\\))+\\b)+"
)

// SystemCall represents a system call
type SystemCall struct {
	ID   uint16
	Name string
}

type symbolDefinition struct {
	name       string
	syscallIDs []uint16
	subCalls   []string
}

// SourceReader defines the interface for source readers
type SourceReader interface {
	GetReader() (io.Reader, error)
}

// Extract returns all system calls made in the execution path of the dumpFile provided.
func Extract(source SourceReader) ([]SystemCall, error) {
	reader, err := source.GetReader()
	if err != nil {
		return nil, err
	}

	symbols := parseDump(reader)
	if _, isExe := symbols[executableSymbolCalls[0]]; !isExe {
		return nil, errors.New("could not find entry point")
	}
	syscalls := extractSyscalls(symbols)

	return syscalls, nil
}

// kick off process from executable key entry points.
func extractSyscalls(symbols map[string]symbolDefinition) []SystemCall {
	syscallID := make(chan uint16)

	var wg sync.WaitGroup
	wg.Add(len(executableSymbolCalls))
	for _, symbol := range executableSymbolCalls {
		go func(s string) {
			dumpWalker(symbols, s, syscallID)
			wg.Done()
		}(symbol)
	}

	go func() {
		wg.Wait()
		close(syscallID)
	}()

	syscalls := make([]SystemCall, 0)
	unique := make(map[uint16]bool)

	for id := range syscallID {
		if _, exists := unique[id]; !exists {
			unique[id] = true
			syscalls = append(syscalls, SystemCall{
				ID:   id,
				Name: systemCalls[id],
			})
		}
	}

	return syscalls
}

func parseDump(reader io.Reader) map[string]symbolDefinition {
	symbols := make(map[string]symbolDefinition)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()

		stack := stack.New()
		symbol := symbolDefinition{
			subCalls:   make([]string, 0),
			syscallIDs: make([]uint16, 0),
		}
		symbolName, found := getSymbolName(line)

		for found {
			if scanner.Scan() {
				line = scanner.Text()
				if isEndOfSymbol(line) {
					break
				}

				if id, found := tryPopSyscallID(line, stack); found {
					symbol.syscallIDs = append(symbol.syscallIDs, id)
					continue
				}

				if subcall, found := getCallTarget(line); found {
					symbol.subCalls = append(symbol.subCalls, subcall)
					continue
				}

				stackSyscallIDIfNecessary(line, stack)
			} else {
				break
			}
		}

		if len(symbol.subCalls) > 0 || len(symbol.syscallIDs) > 0 {
			symbols[symbolName] = symbol
		}
	}

	return symbols
}

func dumpWalker(symbols map[string]symbolDefinition, symbolName string, syscallID chan<- uint16) {
	var walk func(symbol string)
	processed := make(map[string]bool)

	walk = func(symbol string) {
		if _, exists := processed[symbol]; !exists {
			processed[symbol] = true
			if s, found := symbols[symbol]; found {
				for _, id := range s.syscallIDs {
					syscallID <- id
				}

				for _, name := range s.subCalls {
					walk(name)
				}
			}
		}
	}

	walk(symbolName)
}

func stackSyscallIDIfNecessary(assemblyLine string, s *stack.Stack) {
	if id, ok := getSyscallID(assemblyLine); ok {
		s.Push(id)
	}
}

func tryPopSyscallID(assemblyLine string, s *stack.Stack) (uint16, bool) {
	if s.Len() > 0 && containsSyscall(assemblyLine) {
		val1 := s.Pop()
		val2 := s.Pop()

		syscallID := val1
		if val2 != nil {
			syscallID = val2
		}
		return syscallID.(uint16), true
	}

	return 0, false
}

func getSyscallID(assemblyLine string) (uint16, bool) {
	re := regexp.MustCompile(syscallHexIDRegex)
	captures := re.FindStringSubmatch(assemblyLine)

	if captures != nil && len(captures) > 0 {
		if n, err := strconv.ParseUint(captures[2], 16, 16); err == nil {
			id := uint16(n)
			if _, exists := systemCalls[id]; exists {
				return id, true
			}
		}
	}

	return 0, false
}

func getSymbolName(assemblyLine string) (string, bool) {
	return extract(assemblyLine, symbolDefinitionRegex)
}

func getCallTarget(assemblyLine string) (string, bool) {
	return extract(assemblyLine, callCaptureRegex)
}

func extract(assemblyLine, regex string) (string, bool) {
	re := regexp.MustCompile(regex)
	captures := re.FindStringSubmatch(assemblyLine)

	if captures != nil && len(captures) > 0 {
		return captures[1], true
	}

	return "", false
}

func containsSyscall(assemblyLine string) bool {
	re := regexp.MustCompile("SYSCALL|golang.org/x/sys/unix.Syscall|syscall.Syscall")
	captures := re.FindStringSubmatch(assemblyLine)

	return (captures != nil && len(captures) > 0)
}

func isEndOfSymbol(line string) bool {
	return (line == "" || line == "\n")
}
