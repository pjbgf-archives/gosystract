package pkg

import (
	"regexp"
	"strconv"

	"github.com/golang-collections/collections/stack"
)

const (
	symbolDefinitionRegex   string = "TEXT.(\\b([a-zA-Z0-9_.\\/]|\\.|\\(\\*[a-zA-Z0-9_.\\/]+\\))+\\b)\\("
	golangSyscallHexIDRegex string = "MOVQ.\\$0x([0-9a-fA-F]+)"
	syscallHexIDRegex       string = "MOVL.\\$0x([0-9a-fA-F]+)"
	callCaptureRegex        string = ".+CALL.(\\b([a-zA-Z0-9_.\\/]|\\.|\\(\\*[a-zA-Z0-9_.\\/]+\\))+\\b)+"
)

func StackSyscallIDIfNecessary(assemblyLine string, s *stack.Stack) {

	if id, ok := getSyscallID(assemblyLine, syscallHexIDRegex); ok {
		s.Push(id)
	}

	if id, ok := getSyscallID(assemblyLine, golangSyscallHexIDRegex); ok {
		s.Push(id)
	}
}

func TryGetSyscallID(assemblyLine string, s *stack.Stack) (uint8, bool) {

	if s.Len() > 0 && ContainsSyscall(assemblyLine) {
		val1 := s.Pop()
		val2 := s.Pop()

		syscallID := val1
		if val2 != nil {
			syscallID = val2
		}
		return syscallID.(uint8), true
	}

	return 0, false
}

func GetSyscallID(assemblyLine string) (uint8, bool) {

	if id, ok := getSyscallID(assemblyLine, syscallHexIDRegex); ok {
		return id, ok
	}

	return getSyscallID(assemblyLine, golangSyscallHexIDRegex)
}

func getSyscallID(assemblyLine, regex string) (uint8, bool) {
	re := regexp.MustCompile(regex)
	captures := re.FindStringSubmatch(assemblyLine)

	if captures != nil && len(captures) > 0 {
		if id, err := strconv.ParseUint(captures[1], 16, 64); err == nil {
			return uint8(id), true
		}
	}

	return 0, false
}

func GetSymbolName(assemblyLine string) (string, bool) {
	re := regexp.MustCompile(symbolDefinitionRegex)
	captures := re.FindStringSubmatch(assemblyLine)

	if captures != nil && len(captures) > 0 {
		return captures[1], true
	}

	return "", false
}

func GetCallTarget(assemblyLine string) (string, bool) {
	re := regexp.MustCompile(callCaptureRegex)
	captures := re.FindStringSubmatch(assemblyLine)

	if captures != nil && len(captures) > 0 {
		return captures[1], true
	}

	return "", false
}

func IsSymbolDefinition(assemblyLine string) bool {
	re := regexp.MustCompile(symbolDefinitionRegex)
	captures := re.FindStringSubmatch(assemblyLine)

	return (captures != nil && len(captures) > 0)
}

func IsSymbolDefinition2(assemblyLine, symbolName string) bool {
	re := regexp.MustCompile("TEXT." + regexp.QuoteMeta(symbolName) + "\\(")
	captures := re.FindStringSubmatch(assemblyLine)

	return (captures != nil && len(captures) > 0)
}

func ContainsSyscall(assemblyLine string) bool {
	re := regexp.MustCompile("SYSCALL|golang.org/x/sys/unix.Syscall")
	captures := re.FindStringSubmatch(assemblyLine)

	return (captures != nil && len(captures) > 0)
}

func IsEndOfSymbol(line string) bool {
	return (line == "" || line == "\n")
}
