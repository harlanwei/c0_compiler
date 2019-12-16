package analyzer

import (
	"c0_compiler/internal/cc0_error"
	"c0_compiler/internal/instruction"
	"c0_compiler/internal/token"
)

func analyzeCondition() *Error {
	// <condition> ::= <expression>[<relational-operator><expression>]

	pos := getCurrentPos()
	if err := analyzeExpression(); err != nil {
		resetHeadTo(pos)
		return err
	}
	pos = getCurrentPos()
	next, err := getNextToken()
	if err != nil || !next.IsARelationalOperator() {
		resetHeadTo(pos)
		return nil
	}
	operator := next.Kind
	if err := analyzeExpression(); err != nil {
		resetHeadTo(pos)
		return err
	}
	currentFunction.Append(instruction.Icmp)
	// Only jump when the condition doesn't stand
	switch operator {
	case token.LessThan:
		currentFunction.Append(instruction.Jge, 0)
	case token.LessThanOrEqual:
		currentFunction.Append(instruction.Jg, 0)
	case token.EqualTo:
		currentFunction.Append(instruction.Jne, 0)
	case token.GreaterThanOrEqual:
		currentFunction.Append(instruction.Jl, 0)
	case token.GreaterThan:
		currentFunction.Append(instruction.Jle, 0)
	case token.NotEqualTo:
		currentFunction.Append(instruction.Je, 0)
	}
	return nil
}

func analyzeExpression() *Error {
	// <expression> ::= <additive-expression>
	return analyzeAdditiveExpression()
}

func analyzeAdditiveExpression() *Error {
	// <additive-expression> ::= <multiplicative-expression>{<additive-operator><multiplicative-expression>}

	// <multiplicative-expression>
	if err := analyzeMultiplicativeExpression(); err != nil {
		return err
	}

	// {<additive-operator><multiplicative-expression>}
	for {
		pos := getCurrentPos()
		next, err := getNextToken()
		if err != nil || !next.IsAnAdditiveOperator() {
			resetHeadTo(pos)
			return nil
		}
		operator := next.Kind
		if err := analyzeMultiplicativeExpression(); err != nil {
			resetHeadTo(pos)
			return err
		}

		if operator == token.PlusSign {
			currentFunction.Append(instruction.Iadd)
		} else {
			currentFunction.Append(instruction.Isub)
		}
	}
}

func analyzeMultiplicativeExpression() *Error {
	// <multiplicative-expression> ::= <unary-expression>{<multiplicative-operator><unary-expression>}

	// <unary-expression>
	if err := analyzeUnaryExpression(); err != nil {
		return err
	}

	// {<multiplicative-operator><unary-expression>}
	for {
		pos := getCurrentPos()
		next, err := getNextToken()
		if err != nil || !next.IsAMultiplicativeOperator() {
			resetHeadTo(pos)
			return nil
		}
		if err := analyzeUnaryExpression(); err != nil {
			resetHeadTo(pos)
			return err
		}
		if next.Kind == token.MultiplicationSign {
			currentFunction.Append(instruction.Imul)
		} else {
			currentFunction.Append(instruction.Idiv)
		}
	}
}

func analyzeUnaryExpression() *Error {
	// <unary-expression> ::= [<unary-operator>]<primary-expression>

	// [<unary-operator>]
	shouldBeNegated := false
	pos := getCurrentPos()
	next, err := getNextToken()
	if err != nil {
		return cc0_error.Of(cc0_error.IncompleteExpression).On(currentLine, currentColumn)
	}
	if next.IsAnUnaryOperator() {
		if next.Kind == token.MinusSign {
			shouldBeNegated = true
		}
	} else {
		resetHeadTo(pos)
	}

	// <primary-expression>
	if err := analyzePrimaryExpression(); err != nil {
		resetHeadTo(pos)
		return err
	}

	if shouldBeNegated {
		currentFunction.Append(instruction.Ineg)
	}
	return nil
}

func analyzePrimaryExpression() *Error {
	// <primary-expression> ::=
	//     '('<expression>')'
	//    | <identifier>
	//    | <integer-literal>
	//    | <function-call>

	pos := getCurrentPos()
	next, err := getNextToken()
	if err != nil {
		return cc0_error.Of(cc0_error.IncompleteExpression).On(currentLine, currentColumn)
	}

	// '('<expression>')'
	if next.Kind == token.LeftParenthesis {
		if err := analyzeExpression(); err != nil {
			return err
		}
		next, err = getNextToken()
		if err != nil {
			return cc0_error.Of(cc0_error.IncompleteExpression).On(currentLine, currentColumn)
		}
		if next.Kind != token.RightParenthesis {
			return cc0_error.Of(cc0_error.IllegalExpression).On(currentLine, currentColumn)
		}
	} else if next.Kind == token.Identifier { // <identifier> || <function-call>
		identifier := next.Value.(string)
		sb := currentSymbolTable.GetSymbolNamed(identifier)
		if sb == nil {
			return cc0_error.Of(cc0_error.UndefinedIdentifier).On(currentLine, currentColumn)
		}
		if sb.IsCallable {
			resetHeadTo(pos) // `analyzeFunctionCall` needs the identifier, thus the reset
			if err := analyzeFunctionCall(); err != nil {
				return err
			}
		} else {
			if sb.Kind == token.Double {
				currentFunction.Append(instruction.Ipush, 0)
			}
			currentFunction.Append(instruction.Loada, 0, sb.Address)
			currentFunction.Append(instruction.Iload)
		}
	} else if next.Kind == token.IntegerLiteral {
		// <integer-literal>
		currentFunction.Append(instruction.Ipush, int(next.Value.(int64)))
	} else {
		return cc0_error.Of(cc0_error.IllegalExpression).On(currentLine, currentColumn)
	}
	return nil
}

func analyzeAssignmentExpression() *Error {
	// <assignment-expression> ::= <identifier><assignment-operator><expression>
	pos := getCurrentPos()
	next, err := getNextToken()
	if err != nil || next.Kind != token.Identifier {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.IncompleteExpression).On(currentLine, currentColumn)
	}
	identifier := next.Value.(string)
	address := currentSymbolTable.GetAddressOf(identifier)
	if next, err := getNextToken(); err != nil || next.Kind != token.AssignmentSign {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.IncompleteExpression).On(currentLine, currentColumn)
	}
	currentFunction.Append(instruction.Loada, 0, address)
	if err := analyzeExpression(); err != nil {
		resetHeadTo(pos)
		return err
	}
	currentFunction.Append(instruction.Istore)
	return nil
}

func analyzeExpressionList() *Error {
	// <expression-list> ::= <expression>{','<expression>}
	pos := getCurrentPos()
	if err := analyzeExpression(); err != nil {
		resetHeadTo(pos)
		return err
	}
	for {
		pos = getCurrentPos()
		if next, err := getNextToken(); err != nil || next.Kind != token.Comma {
			resetHeadTo(pos)
			return nil
		}
		if err := analyzeExpression(); err != nil {
			return err
		}
	}
}

func analyzeFunctionCall() *Error {
	// <identifier> '(' [<expression-list>] ')'
	pos := getCurrentPos()
	next, err := getNextToken()
	if err != nil {
		return cc0_error.Of(cc0_error.IncompleteFunctionCall).On(currentLine, currentColumn)
	}
	if next.Kind != token.Identifier {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.IncompleteFunctionCall).On(currentLine, currentColumn)
	}
	identifier := next.Value.(string)
	sb := globalSymbolTable.GetSymbolNamed(identifier)
	if sb == nil {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.UndefinedIdentifier).On(currentLine, currentColumn)
	}
	if next, err := getNextToken(); err != nil || next.Kind != token.LeftParenthesis {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.IncompleteFunctionCall).On(currentLine, currentColumn)
	}
	// TODO: check parameters and push them on the stack
	_ = analyzeExpressionList()
	if next, err := getNextToken(); err != nil || next.Kind != token.RightParenthesis {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.IncompleteFunctionCall).On(currentLine, currentColumn)
	}
	currentFunction.Append(instruction.Call, sb.Address)
	return nil
}
