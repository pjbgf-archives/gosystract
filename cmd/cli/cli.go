package internal

import (
	"errors"
	"io"
	"text/template"

	"github.com/pjbgf/gosystract/cmd/systract"
)

var (
	invalidSyntaxMessage string = "invalid systax"

	usageMessage string = `gosystract returns the names and IDs of all system calls being called inside a go application.
Usage: 
	gosystrac goapp.dump
	gosystrac goapp.dump "{{- range . }}{{printf "%d - %s\n" .ID .Name}}{{- end}}"

To generate a dump file from a go application use: 
	go tool objdump goapp > goapp.dump`

	resultGoTemplate string = `{{if . -}}
{{- len . }} system calls found:
{{- range . }}
    {{ .Name }} ({{.ID}})
{{- end}}
{{- else}}no systems calls were found{{- end}}
`
)

// Run does basic handling of user input
func Run(output io.Writer, args []string, extract func(dumpFileName string) ([]systract.SystemCall, error)) error {

	if len(args) < 2 || len(args) > 3 {
		_, _ = output.Write([]byte(usageMessage))
		return errors.New(invalidSyntaxMessage)
	}

	fileName := args[1]
	customFormat := ""
	if len(args) == 3 {
		customFormat = args[2]
	}

	syscalls, err := extract(fileName)
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
