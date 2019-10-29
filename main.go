package main

import (
	"bufio"
	"os"
	"sync"

	"github.com/golang-collections/collections/stack"
	"github.com/pjbgf/dump/pkg"
	"github.com/sirupsen/logrus"
)

var (
	dumpFileName string
	logger       *logrus.Logger

	processedNames map[string]bool
	namesMutex     sync.RWMutex
	processedIDs   map[uint16]bool
	idsMutex       sync.RWMutex

	symbols map[string]symbolDefinition
)

type symbolDefinition struct {
	syscallIDs []uint16
	subCalls   []string
}

const (
	callCaptureRegex         string = ".+CALL.(\b([a-zA-Z0-9_.\\/]|\\.|\\(\\*[a-z]+\\))+\b)+"
	callsFromEntryPointRegex string = "main.go" + callCaptureRegex
)

func init() {
	logger = logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	processedNames, processedIDs = make(map[string]bool), make(map[uint16]bool)
	symbols = make(map[string]symbolDefinition)
}

func main() {
	consume := func(id uint16) {
		name := pkg.SystemCalls[id]
		logger.Infof("syscall found: %s (%d)", name, id)
	}

	parseFile()

	symbolCalls := []string{"main.init.0", "main.init.1", "main.main"}
	for _, symbol := range symbolCalls {
		processPerStage(symbol, consume)
	}
}

func parseFile() {

	dumpFileName = "keyring" //os.Args[1]
	file, err := os.Open(dumpFileName)
	if err != nil {
		panic(file)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		symbolName, found := pkg.GetSymbolName(line)
		stack := stack.New()
		symbol := symbolDefinition{
			subCalls:   make([]string, 0),
			syscallIDs: make([]uint16, 0),
		}

		for found {
			if scanner.Scan() {
				line = scanner.Text()

				if pkg.IsEndOfSymbol(line) {
					break
				}

				if id, found := pkg.TryGetSyscallID(line, stack); found {
					symbol.syscallIDs = append(symbol.syscallIDs, id)
					continue
				}

				if subcall, found := pkg.GetCallTarget(line); found {
					symbol.subCalls = append(symbol.subCalls, subcall)
					continue
				}

				pkg.StackSyscallIDIfNecessary(line, stack)
			} else {
				break
			}
		}

		if len(symbol.subCalls) > 0 || len(symbol.syscallIDs) > 0 {
			symbols[symbolName] = symbol
			logger.Debugln(symbolName)
			logger.Debugln(symbol)
		} else {
			logger.Debugf("[%s] skipping empty symbol", symbolName)
		}
	}
}

func processPerStage(symbolName string, consume func(uint16)) {
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

			logger.Debugln(symbolName)
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

					processPerStage(name, consume)
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
