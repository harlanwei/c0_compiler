package analyzer

import (
	"c0_compiler/internal/cc0_error"
	"c0_compiler/internal/instruction"
	"c0_compiler/internal/token"
)

func analyzeIOStatement() *Error {
	// <scan-statement>  ::= 'scan' '(' <identifier> ')' ';'
	// <print-statement> ::= 'print' '(' [<printable-list>] ')' ';'
	pos := getCurrentPos()
	next, err := getNextToken()
	if err != nil {
		return cc0_error.Of(cc0_error.InvalidStatement).On(currentLine, currentColumn)
	}
	if next.Kind == token.Scan {
		if next, err := getNextToken(); err != nil || next.Kind != token.LeftParenthesis {
			resetHeadTo(pos)
			return cc0_error.Of(cc0_error.InvalidStatement).On(currentLine, currentColumn)
		}
		next, err = getNextToken()
		if err != nil || next.Kind != token.Identifier {
			resetHeadTo(pos)
			return cc0_error.Of(cc0_error.InvalidStatement).On(currentLine, currentColumn)
		}
		identifier := next.Value.(string)
		if next, err := getNextToken(); err != nil || next.Kind != token.RightParenthesis {
			resetHeadTo(pos)
			return cc0_error.Of(cc0_error.InvalidStatement).On(currentLine, currentColumn)
		}
		currentFunction.Append(instruction.Loada, currentSymbolTable.GetLevelDiff(identifier), currentSymbolTable.GetAddressOf(identifier))
		currentFunction.Append(instruction.Iscan)
		currentFunction.Append(instruction.Istore)
	} else if next.Kind == token.Print {
		if next, err := getNextToken(); err != nil || next.Kind != token.LeftParenthesis {
			resetHeadTo(pos)
			return cc0_error.Of(cc0_error.InvalidStatement).On(currentLine, currentColumn)
		}
		if err := analyzePrintableList(); err != nil {
			resetHeadTo(pos)
			return err
		}
		if next, err := getNextToken(); err != nil || next.Kind != token.RightParenthesis {
			resetHeadTo(pos)
			return cc0_error.Of(cc0_error.InvalidStatement).On(currentLine, currentColumn)
		}
		currentFunction.Append(instruction.Printl)
	} else {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.InvalidStatement).On(currentLine, currentColumn)
	}
	if next, err := getNextToken(); err != nil || next.Kind != token.Semicolon {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.InvalidStatement).On(currentLine, currentColumn)
	}
	return nil
}

func analyzePrintableList() *Error {
	// <printable-list>  ::= <printable> {',' <printable>}
	pos := getCurrentPos()
	if err := analyzePrintable(); err != nil {
		resetHeadTo(pos)
		return err
	}
	for {
		pos = getCurrentPos()
		if next, err := getNextToken(); err != nil || next.Kind != token.Comma {
			resetHeadTo(pos)
			return nil
		}
		currentFunction.Append(instruction.Bipush, 32)
		currentFunction.Append(instruction.Cprint)
		if err := analyzePrintable(); err != nil {
			return err
		}
	}
}

func analyzePrintable() *Error {
	// <printable> ::= <expression> | <string-literal> | <char-literal>
	pos := getCurrentPos()
	kind, err := analyzeExpression()
	if err != nil {
		resetHeadTo(pos)
	} else {
		if kind == token.Double {
			currentFunction.Append(instruction.Dprint)
		} else if kind == token.Char {
			currentFunction.Append(instruction.Cprint)
		} else {
			currentFunction.Append(instruction.Iprint)
		}
		return nil
	}
	next, anotherErr := getNextToken()
	if anotherErr != nil {
		return cc0_error.Of(cc0_error.IncompleteExpression)
	}
	if next.Kind == token.StringLiteral {
		address := globalSymbolTable.AddALiteral(instruction.ConstantKindString, next.Value.(string))
		currentFunction.Append(instruction.Loadc, -address)
		currentFunction.Append(instruction.Sprint)
		return nil
	} else if next.Kind == token.CharLiteral {
		currentFunction.Append(instruction.Bipush, int(next.Value.(int32)))
		currentFunction.Append(instruction.Cprint)
	} else {
		return cc0_error.Of(cc0_error.IncompleteExpression)
	}
	return nil
}
