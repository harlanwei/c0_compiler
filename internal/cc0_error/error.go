package cc0_error

import (
	"fmt"
	"os"
)

const (
	Source = 1
	Parser = 2
)

func PrintToStdErr(message string) {
	_, _ = fmt.Fprint(os.Stderr, message)
}

func PrintlnToStdErr(message string) {
	_, _ = fmt.Fprintln(os.Stderr, message)
}

func throw(source int) {
	var sourceMessage string
	switch source {
	case Source:
		sourceMessage = "Failed to open source file."
	case Parser:
		sourceMessage = "Parser failed."
	}
	PrintlnToStdErr(sourceMessage)
}

func ThrowButStayAlive(source int) {
	PrintToStdErr("Warning: ")
	throw(source)
	PrintlnToStdErr("")
}

func ThrowAndExit(source int) {
	PrintToStdErr("Fatal: ")
	throw(source)
	os.Exit(source)
}
