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

func PrintfToStdErr(formatString string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, formatString, args...)
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

const (
	Bug = iota
	NoMoreTokens
	IncompleteVariableDeclaration
	InvalidDeclaration
	IncompleteExpression
	IllegalExpression
	RedeclaredAnIdentifier
	InvalidStatement
	IncompleteFunctionCall
	UndefinedIdentifier
)

type Error struct {
	code   int
	line   int
	column int
}

func Of(code int) *Error {
	return &Error{code, 0, 0}
}

func (error *Error) On(line, column int) *Error {
	error.line = line
	error.column = column
	return error
}

func (error *Error) Fatal(from int) {
	ReportLineAndColumn(error.line, error.column)
	PrintlnToStdErr(error.Error())
	ThrowAndExit(from)
}

func (error *Error) Error() string {
	switch error.code {
	case IncompleteVariableDeclaration:
		return "the variable declaration is incomplete."
	case InvalidDeclaration:
		return "the declaration is not complying with the syntax."
	case Bug:
		return "there is a bug in the analyzer."
	case IncompleteExpression:
		return "the expression is not complete."
	case IllegalExpression:
		return "unexpected components in the expression."
	case RedeclaredAnIdentifier:
		return "an identifier cannot be redeclared."
	case InvalidStatement:
		return "encountered an illegal statement."
	case IncompleteFunctionCall:
		return "the function call is not complete."
	case UndefinedIdentifier:
		return "cannot use an undefined identifier."
	default:
		return "an unknown error occurred."
	}
}
