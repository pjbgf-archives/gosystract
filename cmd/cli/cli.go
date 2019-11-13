// Package cli provides a command-line interface for gosystract.
package cli

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/pjbgf/gosystract/cmd/systract"
)

var (
	// gitcommit is set at compilation time through ldflags
	gitcommit string = "[ not set ]"

	invalidSyntaxMessage string = "syntax invalid"

	usageMessage string = `Usage:
	gosystrac [flags] filePath

Flags:
	--dumpfile, -d    Handles a dump file instead of go executable.
	--template	  Defines a go template for the results.
			  Example: --template="{{- range . }}{{printf "%d - %s\n" .ID .Name}}{{- end}}"
`

	resultGoTemplate string = `{{if . -}}
{{- len . }} system calls found:
{{- range . }}
    {{ .Name }} ({{.ID}})
{{- end}}
{{- else}}no systems calls were found{{- end}}
`
)

func parseInputValues(args []string) (
	inputIsDumpFile bool, customFormat string, fileName string, err error) {

	if len(args) < 2 {
		err = errors.New(invalidSyntaxMessage)
		return
	}

	fileName = args[len(args)-1]
	for _, arg := range args[1:] {
		if arg == "--dumpfile" || arg == "-d" {
			inputIsDumpFile = true
			continue
		}

		if strings.HasPrefix(arg, "--template=") {
			customFormat = strings.TrimPrefix(arg, "--template=")

			if strings.HasPrefix(customFormat, "\"") {
				customFormat = strings.TrimPrefix(customFormat, "\"")
			}

			if strings.HasSuffix(customFormat, "\"") {
				customFormat = strings.TrimSuffix(customFormat, "\"")
			}

			continue
		}
	}

	return
}

/*
Run processes the source and writes the found syscalls into output.
The parameter args contains the executable name, the optional flags followed by the filepath.

Example:
[]string{ "gosystract", "--dumpfile", "filename"}

Flag options:

--dumpfile, -d    Handles a dump file instead of go executable.

--template        Defines a go template for the results.
*/
func Run(output io.Writer, args []string, extract func(source systract.SourceReader) ([]systract.SystemCall, error)) error {

	inputIsDumpFile, customFormat, fileName, err := parseInputValues(args)
	if err != nil {
		usage := fmt.Sprintf("gosystract version %s\n%s", gitcommit, usageMessage)
		_, _ = output.Write([]byte(usage))
		return errors.New(invalidSyntaxMessage)
	}

	var sourceReader systract.SourceReader
	if inputIsDumpFile {
		sourceReader = systract.NewDumpReader(fileName)
	} else {
		sourceReader = systract.NewExeReader(fileName)
	}

	syscalls, err := extract(sourceReader)
	if err != nil {
		return err
	}

	return writeResults(output, syscalls, customFormat)
}

func writeResults(output io.Writer, syscalls []systract.SystemCall, customFormat string) error {
	var t *template.Template
	if customFormat != "" {
		t = template.Must(template.New("result").Parse(customFormat))
	} else {
		t = template.Must(template.New("result").Parse(resultGoTemplate))
	}

	return t.Execute(output, syscalls)
}
