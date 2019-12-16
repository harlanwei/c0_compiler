package analyzer

import (
	"c0_compiler/internal/cc0_error"
	"c0_compiler/internal/instruction"
	"c0_compiler/internal/token"
)

func analyzeLoopStatement() *Error {
	// <loop-statement> ::= 'while' '(' <condition> ')' <statement>

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
	offsetBeforeConditionEvaluation := currentFunction.GetCurrentOffset()
	if err := analyzeCondition(); err != nil {
		resetHeadTo(pos)
		return err
	}
	conditionLine := currentFunction.GetCurrentLine()

	next, err = getNextToken()
	if err != nil || next.Kind != token.RightParenthesis {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.InvalidStatement).On(currentLine, currentColumn)
	}
	if err := analyzeStatement(); err != nil {
		return err
	}
	currentFunction.Append(instruction.Jmp, offsetBeforeConditionEvaluation)
	conditionLine.SetFirstOperandTo(currentFunction.GetCurrentOffset())
	return nil
}
