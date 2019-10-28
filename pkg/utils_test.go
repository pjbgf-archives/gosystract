package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSyscallID(t *testing.T) {
	assertThat := func(assumption, assemblyLine string, expectedId uint8, expectedMatch bool) {
		assert := assert.New(t)

		id, ok := GetSyscallID(assemblyLine)

		assert.Equal(expectedMatch, ok)
		assert.Equal(expectedId, id)
	}

	assertThat("should support golang.org/x/sys/unix.Syscall calls", "zsyscall_linux_amd64.go:442	0x48bd75		48c704247d000000	MOVQ $0x7d, 0(SP)", 125, true)
	assertThat("should support SYSCALL calls", "sys_linux_amd64.s:625	0x453610		b818000000		MOVL $0x18, AX", 24, true)
}

func TestIsCallInstruction(t *testing.T) {

	assertThat := func(assumption, assemblyLine, expectedTarget string, expectedMatch bool) {
		assert := assert.New(t)
		target, ok := GetCallTarget(assemblyLine)

		assert.Equal(expectedMatch, ok)
		assert.Equal(expectedTarget, target)
	}

	assertThat("should match runtime funcs", "main.go:35		0x48c3d8		e8c334fcff		CALL runtime.morestack_noctxt(SB)", "runtime.morestack_noctxt", true)
	assertThat("should match composed funcs", "print.go:265		0x481553		e8c87b0000		CALL fmt.(*pp).doPrintln(SB)", "fmt.(*pp).doPrintln", true)
	assertThat("should match composed funcs 2", "print.go:134		0x480f9c		e84ff7fdff		CALL sync.(*Pool).Get(SB)", "sync.(*Pool).Get", true)
	assertThat("should not match funcs definition", "TEXT fmt.Fprintln(SB) /usr/local/go/src/fmt/print.go", "", false)
}

func TestIsSymbolDefinition(t *testing.T) {

	assertThat := func(assumption, assemblyLine, symbolName string, expected bool) {
		assert := assert.New(t)
		isSymbolDefinition := IsSymbolDefinition(assemblyLine, symbolName)

		assert.Equal(expected, isSymbolDefinition)
	}

	assertThat("should match main.main symbol", "TEXT main.main(SB) /media/pjb/src/git/learn-golang/caps/main.go", "main.main", true)
	assertThat("should match symbol with *", "TEXT sync.(*Pool).Get(SB) /usr/local/go/src/sync/pool.go", "sync.(*Pool).Get", true)
}

func TestContainsSyscall(t *testing.T) {

	assertThat := func(assumption, assemblyLine string, expected bool) {
		assert := assert.New(t)
		containsSyscall := ContainsSyscall(assemblyLine)

		assert.Equal(expected, containsSyscall)
	}

	assertThat("should return true for SYSCALL instruction", "sys_linux_amd64.s:535	0x4534f1		0f05			SYSCALL", true)
	assertThat("should return true for golang.org/x/sys/unix.Syscall instruction", "zsyscall_linux_amd64.go:442	0x48bd9a		e881030000		CALL golang.org/x/sys/unix.Syscall(SB)", true)
	assertThat("should return false for instructions containing syscall on their name", "proc.go:2853		0x430ab3		eb8b			JMP runtime.entersyscall_sysmon(SB)", false)
}
