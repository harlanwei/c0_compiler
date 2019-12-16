package analyzer

import (
	"c0_compiler/internal/cc0_error"
	"c0_compiler/internal/instruction"
	"c0_compiler/internal/token"
)

func analyzeConditionStatement() *Error {
	// <condition-statement> ::=  'if' '(' <condition> ')' <statement> ['else' <statement>]

	pos := getCurrentPos()
	if next, err := getNextToken(); err != nil || next.Kind != token.If {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.InvalidStatement).On(currentLine, currentColumn)
	}
	if next, err := getNextToken(); err != nil || next.Kind != token.LeftParenthesis {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.InvalidStatement).On(currentLine, currentColumn)
	}
	if err := analyzeCondition(); err != nil {
		resetHeadTo(pos)
		return err
	}
	if next, err := getNextToken(); err != nil || next.Kind != token.RightParenthesis {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.InvalidStatement).On(currentLine, currentColumn)
	}
	conditionalJumpLine := currentFunction.GetCurrentLine()

	if err := analyzeStatement(); err != nil {
		return err
	}
	ifOffset := currentFunction.GetCurrentOffset()
	conditionalJumpLine.SetFirstOperandTo(ifOffset)

	pos = getCurrentPos()
	if next, err := getNextToken(); err != nil || next.Kind != token.Else {
		resetHeadTo(pos)
		return nil
	}

	if err := analyzeStatement(); err != nil {
		resetHeadTo(pos)
		return err
	}
	return nil
}

func analyzeJumpStatement() *Error {
	// <jump-statement> ::= <return-statement>
	return analyzeReturnStatement()
}

func analyzeReturnStatement() *Error {
	// <return-statement> ::= 'return' [<expression>] ';'
	pos := getCurrentPos()
	if next, err := getNextToken(); err != nil || next.Kind != token.Return {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.InvalidStatement).On(currentLine, currentColumn)
	}

	currentFunction.Append(instruction.Loada, 0, 0)
	_ = analyzeExpression()
	reservedStackSize := 0
	switch currentFunction.ReturnType {
	case token.Double:
		reservedStackSize = 2
		currentFunction.Append(instruction.Dstore)
	case token.Int:
		reservedStackSize = 1
		currentFunction.Append(instruction.Istore)
	}

	if next, err := getNextToken(); err != nil || next.Kind != token.Semicolon {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.InvalidStatement).On(currentLine, currentColumn)
	}

	currentFunction.PopStack(reservedStackSize)
	currentFunction.Append(instruction.Ret)
	return nil
}
