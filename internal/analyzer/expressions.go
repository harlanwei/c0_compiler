package analyzer

import (
	"c0_compiler/internal/cc0_error"
	"c0_compiler/internal/instruction"
	"c0_compiler/internal/token"
)

func getConvertInstruction(source, dest int) int {
	if source == token.Void || dest == token.Void {
		cc0_error.ThrowAndExit(cc0_error.Analyzer)
	}
	if (source == token.Int || source == token.Char) && dest == token.Double {
		return instruction.I2d
	} else if source == token.Double && (source == token.Int || source == token.Char) {
		return instruction.D2i
	}
	panic("An unexpected type convert took place here.")
	return 0
}

func convertType(source, dest int) {
	currentFunction.Append(getConvertInstruction(source, dest))
}

func convergeToLargerType(lhs, rhs int) int {
	if lhs > rhs {
		return lhs
	}
	return rhs
}

func analyzeCondition() *Error {
	// <condition> ::= <expression>[<relational-operator><expression>]

	pos := getCurrentPos()
	kind, err := analyzeExpression()
	if err != nil {
		resetHeadTo(pos)
		return err
	}
	previousOffset := currentFunction.GetCurrentOffset()
	pos = getCurrentPos()
	next, anotherErr := getNextToken()
	if anotherErr != nil || !next.IsARelationalOperator() {
		resetHeadTo(pos)
		return nil
	}
	operator := next.Kind
	anotherKind, err := analyzeExpression()
	if err != nil {
		resetHeadTo(pos)
		return err
	}

	if kind != anotherKind {
		convergedKind := convergeToLargerType(kind, anotherKind)
		if kind != convergedKind {
			currentFunction.InsertInstructionAt(previousOffset, getConvertInstruction(kind, convergedKind))
		} else {
			convertType(anotherKind, convergedKind)
		}
		kind = convergedKind
	}

	if kind == token.Int {
		currentFunction.Append(instruction.Icmp)
	} else {
		currentFunction.Append(instruction.Dcmp)
	}

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

func analyzeExpression() (int, *Error) {
	// <expression> ::= <additive-expression>
	return analyzeAdditiveExpression()
}

func analyzeAdditiveExpression() (int, *Error) {
	// <additive-expression> ::= <multiplicative-expression>{<additive-operator><multiplicative-expression>}

	// <multiplicative-expression>
	kind, err := analyzeMultiplicativeExpression()
	if err != nil {
		return 0, err
	}
	previousOffset := currentFunction.GetCurrentOffset()

	// {<additive-operator><multiplicative-expression>}
	for {
		pos := getCurrentPos()
		next, err := getNextToken()
		if err != nil || !next.IsAnAdditiveOperator() {
			resetHeadTo(pos)
			return 0, nil
		}
		operator := next.Kind
		anotherKind, anotherErr := analyzeMultiplicativeExpression()
		if anotherErr != nil {
			resetHeadTo(pos)
			return 0, anotherErr
		}

		kind := convergeToLargerType(kind, anotherKind)
		if kind != anotherKind {
			convergedKind := convergeToLargerType(kind, anotherKind)
			if kind != convergedKind {
				currentFunction.InsertInstructionAt(previousOffset, getConvertInstruction(kind, convergedKind))
			} else {
				convertType(anotherKind, convergedKind)
			}
			kind = convergedKind
		}

		if operator == token.PlusSign {
			if kind == token.Double {
				currentFunction.Append(instruction.Dadd)
			} else {
				currentFunction.Append(instruction.Iadd)
			}
		} else {
			if kind == token.Double {
				currentFunction.Append(instruction.Dsub)
			} else {
				currentFunction.Append(instruction.Isub)
			}
		}
	}
}

func analyzeCastExpression() (int, *Error) {
	// <cast-expression> ::= {'('<type-specifier>')'}<unary-expression>
	pos := getCurrentPos()
	kind := 0
	next, err := getNextToken()
	if err != nil {
		resetHeadTo(pos)
		return 0, cc0_error.Of(cc0_error.IncompleteExpression)
	}
	if next.Kind == token.LeftParenthesis {
		next, err := getNextToken()
		if err != nil {
			return 0, cc0_error.Of(cc0_error.IncompleteExpression)
		}
		kind = next.Kind
		for {
			anotherPos := getCurrentPos()
			next, err = getNextToken()
			if err != nil {
				resetHeadTo(anotherPos)
				return 0, cc0_error.Of(cc0_error.IncompleteExpression)
			}
			if !next.IsATypeSpecifier() {
				resetHeadTo(anotherPos)
				break
			}
		}
	} else {
		resetHeadTo(pos)
	}
	unaryKind, anotherErr := analyzeUnaryExpression()
	if anotherErr != nil {
		resetHeadTo(pos)
		return 0, anotherErr
	}
	convertType(unaryKind, kind)
	return kind, nil
}

func analyzeMultiplicativeExpression() (int, *Error) {
	// <multiplicative-expression> ::= <cast-expression>{<multiplicative-operator><cast-expression>}

	// <cast-expression>
	kind, err := analyzeCastExpression()
	if err != nil {
		return 0, err
	}

	// {<multiplicative-operator><unary-expression>}
	for {
		pos := getCurrentPos()
		next, err := getNextToken()
		if err != nil || !next.IsAMultiplicativeOperator() {
			resetHeadTo(pos)
			return 0, nil
		}
		anotherKind, anotherErr := analyzeUnaryExpression()
		if err != nil {
			resetHeadTo(pos)
			return 0, anotherErr
		}
		kind = convergeToLargerType(kind, anotherKind)
		if next.Kind == token.MultiplicationSign {
			if kind == token.Double {
				currentFunction.Append(instruction.Dmul)
			} else {
				currentFunction.Append(instruction.Imul)
			}
		} else {
			if kind == token.Double {
				currentFunction.Append(instruction.Ddiv)
			} else {
				currentFunction.Append(instruction.Idiv)
			}
		}
	}
}

func analyzeUnaryExpression() (int, *Error) {
	// <unary-expression> ::= [<unary-operator>]<primary-expression>

	// [<unary-operator>]
	shouldBeNegated := false
	pos := getCurrentPos()
	next, err := getNextToken()
	if err != nil {
		return 0, cc0_error.Of(cc0_error.IncompleteExpression).On(currentLine, currentColumn)
	}
	if next.IsAnUnaryOperator() {
		if next.Kind == token.MinusSign {
			shouldBeNegated = true
		}
	} else {
		resetHeadTo(pos)
	}

	// <primary-expression>
	kind, anotherErr := analyzePrimaryExpression()
	if anotherErr != nil {
		resetHeadTo(pos)
		return 0, anotherErr
	}

	if shouldBeNegated {
		currentFunction.Append(instruction.Ineg)
	}
	return kind, nil
}

func analyzePrimaryExpression() (int, *Error) {
	// <primary-expression> ::=
	//     '('<expression>')'
	//    | <identifier>
	//    | <integer-literal>
	//    | <function-call>

	pos := getCurrentPos()
	next, err := getNextToken()
	kind := 0
	if err != nil {
		return 0, cc0_error.Of(cc0_error.IncompleteExpression).On(currentLine, currentColumn)
	}

	// '('<expression>')'
	if next.Kind == token.LeftParenthesis {
		var anotherErr *Error
		kind, anotherErr = analyzeExpression()
		if anotherErr != nil {
			return 0, anotherErr
		}
		next, err = getNextToken()
		if err != nil {
			return 0, cc0_error.Of(cc0_error.IncompleteExpression).On(currentLine, currentColumn)
		}
		if next.Kind != token.RightParenthesis {
			return 0, cc0_error.Of(cc0_error.IllegalExpression).On(currentLine, currentColumn)
		}
	} else if next.Kind == token.Identifier { // <identifier> || <function-call>
		identifier := next.Value.(string)
		sb := currentSymbolTable.GetSymbolNamed(identifier)
		if sb == nil {
			return 0, cc0_error.Of(cc0_error.UndefinedIdentifier).On(currentLine, currentColumn)
		}
		kind = sb.Kind
		if sb.IsCallable {
			resetHeadTo(pos) // `analyzeFunctionCall` needs the identifier, thus the reset
			if err := analyzeFunctionCall(); err != nil {
				return 0, err
			}
		} else {
			currentFunction.Append(instruction.Loada, currentSymbolTable.GetLevelDiff(identifier), sb.Address)
			if sb.Kind == token.Double {
				currentFunction.Append(instruction.Dload)
			} else {
				currentFunction.Append(instruction.Iload)
			}
		}
	} else if next.Kind == token.IntegerLiteral {
		// <integer-literal>
		currentFunction.Append(instruction.Ipush, int(next.Value.(int64)))
		kind = token.Int
	} else {
		return 0, cc0_error.Of(cc0_error.IllegalExpression).On(currentLine, currentColumn)
	}
	return kind, nil
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
	if currentSymbolTable.GetSymbolNamed(identifier).IsConstant {
		cc0_error.ReportLineAndColumn(currentLine, currentColumn)
		cc0_error.PrintfToStdErr("Cannot assign a new value to the constant: %s\n", identifier)
		cc0_error.ThrowAndExit(cc0_error.Analyzer)
	}
	address := currentSymbolTable.GetAddressOf(identifier)
	if next, err := getNextToken(); err != nil || next.Kind != token.AssignmentSign {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.IncompleteExpression).On(currentLine, currentColumn)
	}
	currentFunction.Append(instruction.Loada, currentSymbolTable.GetLevelDiff(identifier), address)
	kind, anotherErr := analyzeExpression()
	convertType(kind, currentSymbolTable.GetSymbolNamed(identifier).Kind)
	if anotherErr != nil {
		resetHeadTo(pos)
		return anotherErr
	}
	currentFunction.Append(instruction.Istore)
	return nil
}

func analyzeExpressionList() *Error {
	// <expression-list> ::= <expression>{','<expression>}
	pos := getCurrentPos()
	if _, err := analyzeExpression(); err != nil {
		resetHeadTo(pos)
		return err
	}
	for {
		pos = getCurrentPos()
		if next, err := getNextToken(); err != nil || next.Kind != token.Comma {
			resetHeadTo(pos)
			return nil
		}
		if _, err := analyzeExpression(); err != nil {
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
	// TODO: check parameters
	_ = analyzeExpressionList()
	if next, err := getNextToken(); err != nil || next.Kind != token.RightParenthesis {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.IncompleteFunctionCall).On(currentLine, currentColumn)
	}
	currentFunction.Append(instruction.Call, sb.Address)
	return nil
}
