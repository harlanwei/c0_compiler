package cc0_error

import (
	"fmt"
	"os"
)

const (
	Source = iota
	Parser
	Analyzer
	Assembler
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
	case Assembler:
		sourceMessage = "Failed to assemble."
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
	NoMain
	AssignmentToConstant
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

func (error *Error) DieAndReportPosition(from int) {
	ReportLineAndColumn(error.line, error.column)
	PrintlnToStdErr(error.Error())
	ThrowAndExit(from)
}

func (error *Error) Die(from int) {
	PrintlnToStdErr(error.Error())
	ThrowAndExit(from)
}

func (error *Error) Error() string {
	switch error.code {
	case IncompleteVariableDeclaration:
		return "The variable declaration is incomplete."
	case InvalidDeclaration:
		return "The declaration is not complying with the syntax."
	case Bug:
		return "There is a bug in the analyzer."
	case IncompleteExpression:
		return "The expression is not complete."
	case IllegalExpression:
		return "Unexpected components in the expression."
	case RedeclaredAnIdentifier:
		return "An identifier cannot be redeclared."
	case InvalidStatement:
		return "Encountered an illegal statement."
	case IncompleteFunctionCall:
		return "The function call is not complete."
	case UndefinedIdentifier:
		return "Cannot use an undefined identifier."
	case NoMain:
		return "No main function is defined."
	case AssignmentToConstant:
		return "Cannot assign a new value to a constant."
	default:
		return "An unknown error occurred."
	}
}
