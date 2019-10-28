package pkg

import (
	"regexp"
	"strconv"
)

const (
	golangSyscallHexIDRegex string = "MOVQ.\\$0x([0-9a-fA-F]+)"
	syscallHexIDRegex       string = "MOVL.\\$0x([0-9a-fA-F]+)"
	callCaptureRegex        string = ".+CALL.(\\b([a-zA-Z0-9_.\\/]|\\.|\\(\\*[a-zA-Z0-9_.\\/]+\\))+\\b)+"
)

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

func GetCallTarget(assemblyLine string) (string, bool) {
	re := regexp.MustCompile(callCaptureRegex)
	captures := re.FindStringSubmatch(assemblyLine)

	if captures != nil && len(captures) > 0 {
		return captures[1], true
	}

	return "", false
}

func IsSymbolDefinition(assemblyLine, symbolName string) bool {
	re := regexp.MustCompile("TEXT." + regexp.QuoteMeta(symbolName) + "\\(")
	captures := re.FindStringSubmatch(assemblyLine)

	return (captures != nil && len(captures) > 0)
}

func ContainsSyscall(assemblyLine string) bool {
	re := regexp.MustCompile("SYSCALL|golang.org/x/sys/unix.Syscall")
	captures := re.FindStringSubmatch(assemblyLine)

	return (captures != nil && len(captures) > 0)
}
