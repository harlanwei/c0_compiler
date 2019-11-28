package analyzer

import "c0_compiler/internal/cc0_error"

const (
	Bug = iota
	NoMoreTokens
	IncompleteVariableDeclaration
	InvalidDeclaration
	IncompleteExpression
	IllegalExpression
	RedeclaredAnIdentifier
)

type Error struct {
	code   int
	line   int
	column int
}

func reportFatalError(err *Error) {
	cc0_error.ReportLineAndColumn(err.line, err.column)
	cc0_error.PrintlnToStdErr(err.Error())
	cc0_error.ThrowAndExit(cc0_error.Analyzer)
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
	default:
		return "an unknown error occurred."
	}
}

// Create an Error object referring to the current exception.
func errorOf(code int) *Error {
	return &Error{code, globalLineCount, globalColumnCount}
}
