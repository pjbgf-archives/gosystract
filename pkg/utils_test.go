package pkg

import (
	"testing"

	"github.com/golang-collections/collections/stack"
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

	assertThat := func(assumption, assemblyLine string, expected bool) {
		assert := assert.New(t)
		actual := IsSymbolDefinition(assemblyLine)

		assert.Equal(expected, actual)
	}

	assertThat("should match main.main symbol", "TEXT main.main(SB) /media/pjb/src/git/learn-golang/caps/main.go", true)
	assertThat("should match symbol with *", "TEXT sync.(*Pool).Get(SB) /usr/local/go/src/sync/pool.go", true)
}

func TestGetSymbolName(t *testing.T) {

	assertThat := func(assumption, assemblyLine, expectedName string, expectedMatch bool) {
		assert := assert.New(t)
		actual, ok := GetSymbolName(assemblyLine)

		assert.Equal(expectedMatch, ok)
		assert.Equal(expectedName, actual)
	}

	assertThat("should match main.main symbol", "TEXT main.main(SB) /media/pjb/src/git/learn-golang/caps/main.go", "main.main", true)
	assertThat("should match symbol with *", "TEXT sync.(*Pool).Get(SB) /usr/local/go/src/sync/pool.go", "sync.(*Pool).Get", true)
}

func TestIsEndOfSymbolDefinition(t *testing.T) {

	assertThat := func(assumption, assemblyLine string, expected bool) {
		assert := assert.New(t)
		actual := IsEndOfSymbol(assemblyLine)

		assert.Equal(expected, actual)
	}

	assertThat("should match line break", "\n", true)
	assertThat("should match empty line", "", true)
	assertThat("should not match empty symbol definition line", "TEXT sync.(*Pool).Get(SB) /usr/local/go/src/sync/pool.go", false)
	assertThat("should not match empty symbol definition line", "TEXT sync.(*Pool).Get(SB) /usr/local/go/src/sync/pool.go", false)
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

func TestTryGetSyscallID(t *testing.T) {

	assertThat := func(assumption, assemblyLine string, expectedID uint8, s *stack.Stack, expectedMatch bool) {
		assert := assert.New(t)
		actual, ok := TryGetSyscallID(assemblyLine, s)

		assert.Equal(expectedMatch, ok)
		assert.Equal(expectedID, actual)
	}

	stack := stack.New()
	StackSyscallIDIfNecessary("zsyscall_linux_amd64.go:442	0x48bd75		48c704247d000000	MOVQ $0x7d, 0(SP)				", stack)
	StackSyscallIDIfNecessary("zsyscall_linux_amd64.go:442	0x48bd7d		488b442440		MOVQ 0x40(SP), AX				", stack)
	StackSyscallIDIfNecessary("zsyscall_linux_amd64.go:442	0x48bd82		4889442408		MOVQ AX, 0x8(SP)					", stack)
	StackSyscallIDIfNecessary("zsyscall_linux_amd64.go:442	0x48bd91		48c744241800000000	MOVQ $0x0, 0x18(SP)				", stack)

	assertThat("should match main.main symbol", "zsyscall_linux_amd64.go:442	0x48bd9a		e881030000		CALL golang.org/x/sys/unix.Syscall(SB)", uint8(125), stack, true)
}

func TestStackSyscallIDIfNecessary(t *testing.T) {

	assert := assert.New(t)
	stack := stack.New()
	StackSyscallIDIfNecessary("zsyscall_linux_amd64.go:442	0x48bd75		48c704247d000000	MOVQ $0x7d, 0(SP)", stack)

	assert.Equal(1, stack.Len())
	assert.Equal(uint8(125), stack.Pop())
}
