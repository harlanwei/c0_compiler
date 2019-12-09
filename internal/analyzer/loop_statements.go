package analyzer

import (
	"c0_compiler/internal/cc0_error"
	"c0_compiler/internal/token"
)

func analyzeLoopStatement() *Error {
	// <loop-statement> ::= 'while' '(' <condition> ')' <statement>
	// TODO: generate instructions
	pos := getCurrentPos()
	next, err := getNextToken()
	if err != nil || next.Kind != token.While {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.InvalidStatement).On(currentLine, currentColumn)
	}
	next, err = getNextToken()
	if err != nil || next.Kind != token.LeftParenthesis {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.InvalidStatement).On(currentLine, currentColumn)
	}
	if err := analyzeCondition(); err != nil {
		resetHeadTo(pos)
		return err
	}
	next, err = getNextToken()
	if err != nil || next.Kind != token.RightParenthesis {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.InvalidStatement).On(currentLine, currentColumn)
	}
	return nil
}
