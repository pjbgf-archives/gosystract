package systract

import (
	"bufio"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"

	"github.com/golang-collections/collections/stack"
	"github.com/pkg/errors"
)

var (
	executableSymbolCalls = []string{"main.main", "main.init.0", "main.init.1"}
	processedNames        map[string]bool
	namesMutex            sync.RWMutex
	processedIDs          map[uint16]bool
	idsMutex              sync.RWMutex

	symbols map[string]symbolDefinition
)

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
	syscallIDs []uint16
	subCalls   []string
}

func init() {
	processedNames, processedIDs = make(map[string]bool), make(map[uint16]bool)
	symbols = make(map[string]symbolDefinition)
}

// Extract returns all system calls made in the execution path of the dumpFile provided.
func Extract(dumpFileName string) ([]SystemCall, error) {

	filePath, err := sanitiseFileName(dumpFileName)
	if err != nil {
		return nil, err
	}
	if !fileExists(filePath) {
		return nil, errors.New("file does not exist or permission denied")
	}

	syscalls := make([]SystemCall, 0)
	consume := func(id uint16) {
		syscalls = append(syscalls, SystemCall{
			ID:   id,
			Name: systemCalls[id],
		})
	}

	parseFile(filePath)
	if !isExecutable() {
		return nil, errors.New("libraries are currently not supported")
	}
	processExecutable(consume)

	return syscalls, nil
}

func isExecutable() bool {
	_, ok := symbols[executableSymbolCalls[0]]
	return ok
}

// kick off process from executable key entry points.
func processExecutable(consume func(id uint16)) {
	for _, symbol := range executableSymbolCalls {
		processDump(symbol, consume)
	}
}

func fileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil {
		return true
	}

	return false
}

func sanitiseFileName(input string) (string, error) {
	if !path.IsAbs(input) {
		base, err := os.Getwd()
		if err != nil {
			return "", errors.Wrap(err, "error getting current folder")
		}

		return filepath.Join(base, filepath.Clean(input)), nil
	}
	return input, nil
}

func parseFile(filePath string) {
	/* #nosec filePath is pre-processed by sanitiseFileName */
	file, err := os.Open(filePath)
	if err != nil {
		panic(file)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		symbolName, found := getSymbolName(line)
		stack := stack.New()
		symbol := symbolDefinition{
			subCalls:   make([]string, 0),
			syscallIDs: make([]uint16, 0),
		}

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
}

func processDump(symbolName string, consume func(uint16)) {
	var (
		sysCallIDs = make(chan uint16)
		done       = make(chan struct{})
	)

	go func() {
		defer close(sysCallIDs)

		namesMutex.RLock()
		_, exists := processedNames[symbolName]
		namesMutex.RUnlock()

		if !exists {
			namesMutex.Lock()
			processedNames[symbolName] = true
			namesMutex.Unlock()

			if s, found := symbols[symbolName]; found {
				for _, id := range s.syscallIDs {
					idsMutex.RLock()
					_, exists := processedIDs[id]
					idsMutex.RUnlock()

					if !exists {
						idsMutex.Lock()
						processedIDs[id] = true
						idsMutex.Unlock()

						sysCallIDs <- id
					}
				}

				for _, name := range s.subCalls {
					processDump(name, consume)
				}
			}
		}
	}()

	go func() {
		for i := range sysCallIDs {
			consume(i)
		}

		close(done)
	}()

	<-done
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
	re := regexp.MustCompile(symbolDefinitionRegex)
	captures := re.FindStringSubmatch(assemblyLine)

	if captures != nil && len(captures) > 0 {
		return captures[1], true
	}

	return "", false
}

func getCallTarget(assemblyLine string) (string, bool) {
	re := regexp.MustCompile(callCaptureRegex)
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
