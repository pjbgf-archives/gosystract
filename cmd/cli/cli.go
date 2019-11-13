package internal

import (
	"errors"
	"io"
	"text/template"

	"github.com/pjbgf/gosystract/cmd/systract"
)

var (
	invalidSyntaxMessage string = "invalid systax"
	usageMessage         string = `gosystract returns the names and IDs of all system calls being called inside a go application.
Usage: 
	gosystrac goapp.dump
	gosystrac goapp.dump "{{- range . }}{{printf "%d - %s\n" .ID .Name}}{{- end}}"

To generate a dump file from a go application use: 
	go tool objdump goapp > goapp.dump`
)

// Run does basic handling of user input
func Run(output io.Writer, args []string, extract func(dumpFileName string) ([]systract.SystemCall, error)) error {

	if len(args) < 2 || len(args) > 3 {
		output.Write([]byte(usageMessage))
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

	writeResults(output, syscalls, customFormat)

	return nil
}

func writeResults(output io.Writer, syscalls []systract.SystemCall, customFormat string) error {
	var t *template.Template
	if customFormat != "" {
		t = template.Must(template.New("result.tmpl").Parse(customFormat))
	} else {
		t = template.Must(template.ParseFiles("result.tmpl"))
	}

	return t.Execute(output, syscalls)
}
