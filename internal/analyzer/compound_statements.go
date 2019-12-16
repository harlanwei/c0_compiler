package analyzer

import (
	"c0_compiler/internal/cc0_error"
	"c0_compiler/internal/token"
)

func analyzeCompoundStatement() *Error {
	// '{' {<variable-declaration>} <statement-seq> '}'
	pos := getCurrentPos()
	if next, err := getNextToken(); err != nil || next.Kind != token.LeftBracket {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.InvalidStatement).On(currentLine, currentColumn)
	}
	if err := analyzeVariableDeclarations(); err != nil {
		resetHeadTo(pos)
		return err
	}
	if err := analyzeStatementSeq(); err != nil {
		resetHeadTo(pos)
		return err
	}
	if next, err := getNextToken(); err != nil || next.Kind != token.RightBracket {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.InvalidStatement).On(currentLine, currentColumn)
	}
	return nil
}

func analyzeStatementSeq() *Error {
	// {<statement>}
	for {
		if err := analyzeStatement(); err != nil {
			return nil
		}
	}
}

func analyzeStatement() *Error {
	// <statement> ::=
	//		'{' <statement-seq> '}'
	// 		|<condition-statement>
	// 		|<loop-statement>
	// 		|<jump-statement>
	// 		|<print-statement>
	// 		|<scan-statement>
	// 		|<assignment-expression>';'
	// 		|<function-call>';'
	// 		|';'

	// '{' <statement-seq> '}'
	pos := getCurrentPos()
	next, err := getNextToken()
	if err == nil && next.Kind == token.LeftBracket {
		if err := analyzeStatementSeq(); err != nil {
			resetHeadTo(pos)
			return err
		}
		if next, err := getNextToken(); err != nil || next.Kind != token.RightBracket {
			resetHeadTo(pos)
			return cc0_error.Of(cc0_error.InvalidStatement).On(currentLine, currentColumn)
		}
		return nil
	}
	resetHeadTo(pos)

	// <condition-statement>
	if err := analyzeConditionStatement(); err == nil {
		return nil
	}
	resetHeadTo(pos)

	// <loop-statement>
	if err := analyzeLoopStatement(); err == nil {
		return nil
	}
	resetHeadTo(pos)

	// <jump-statement>
	if err := analyzeJumpStatement(); err == nil {
		return nil
	}
	resetHeadTo(pos)

	// <print-statement> | <scan-statement>
	if err := analyzeIOStatement(); err == nil {
		return nil
	}
	resetHeadTo(pos)

	// <assignment-expression>';'
	if err := analyzeAssignmentExpression(); err == nil {
		if next, err := getNextToken(); err != nil || next.Kind != token.Semicolon {
			return cc0_error.Of(cc0_error.InvalidStatement)
		}
		return nil
	}
	resetHeadTo(pos)

	// <function-call>';'
	if err := analyzeFunctionCall(); err == nil {
		if next, err := getNextToken(); err != nil || next.Kind != token.Semicolon {
			return cc0_error.Of(cc0_error.InvalidStatement)
		}
		return nil
	}
	resetHeadTo(pos)

	// ';'
	if next, err := getNextToken(); err != nil || next.Kind != token.Semicolon {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.InvalidStatement).On(currentLine, currentColumn)
	}

	return nil
}
