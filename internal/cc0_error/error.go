package cc0_error

import (
	"fmt"
	"os"
)

const (
	Source = iota
	Parser
	Analyzer
)

func ReportLineAndColumn(line, column int) {
	PrintfToStdErr("At line %d, column %d: ", line, column)
}

func PrintfToStdErr(formatString string, args... interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, formatString, args)
}

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
		sourceMessage = "Parser encountered a problem. See output messages above."
	case Analyzer:
		sourceMessage = "Incorrect syntax encountered. See output messages above."
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
