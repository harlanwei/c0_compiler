package analyzer

import (
	"c0_compiler/internal/cc0_error"
	"c0_compiler/internal/instruction"
	"c0_compiler/internal/token"
)

func getConvertInstruction(source, dest int) []int {
	if source == token.Void || dest == token.Void {
		cc0_error.ThrowAndExit(cc0_error.Analyzer)
	}
	if (source == token.Int || source == token.Char) && dest == token.Double {
		return []int{instruction.I2d}
	} else if source == token.Double && dest == token.Int {
		return []int{instruction.D2i}
	} else if source == token.Double && dest == token.Char {
		return []int{instruction.D2i, instruction.I2c}
	} else if source == token.Int && dest == token.Char {
		return []int{instruction.I2c}
	}
	return []int{0}
}

func convertType(source, dest int) {
	if source == dest {
		return
	}
	for _, inst := range getConvertInstruction(source, dest) {
		currentFunction.Append(inst)
	}
}

func convergeToLargerType(lhs, rhs int) int {
	if lhs > rhs {
		return lhs
	}
	return rhs
}

func computeType(lhs, rhs, previousOffset int) int {
	if lhs != rhs {
		convergedKind := convergeToLargerType(lhs, rhs)
		if lhs != convergedKind {
			currentFunction.ReplaceNopAt(previousOffset, getConvertInstruction(lhs, convergedKind)[0])
		} else {
			convertType(rhs, convergedKind)
		}
		return convergedKind
	}
	return lhs
}

func analyzeCondition() *Error {
	// <condition> ::= <expression>[<relational-operator><expression>]

	pos := getCurrentPos()
	kind, err := analyzeExpression()
	if err != nil {
		resetHeadTo(pos)
		return err
	}
	currentFunction.Append(instruction.Nop)
	previousOffset := currentFunction.GetCurrentOffset() - 1
	pos = getCurrentPos()
	next, anotherErr := getNextToken()
	if anotherErr != nil || !next.IsARelationalOperator() {
		resetHeadTo(pos)
		currentFunction.ChangeInstructionTo(previousOffset, instruction.Je, 0)
		return nil
	}
	operator := next.Kind
	anotherKind, err := analyzeExpression()
	if err != nil {
		resetHeadTo(pos)
		return err
	}

	kind = computeType(kind, anotherKind, previousOffset)
	if kind == token.Int || kind == token.Char {
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

	currentFunction.Append(instruction.Nop)
	previousOffset := currentFunction.GetCurrentOffset() - 1

	// {<additive-operator><multiplicative-expression>}
	for {
		pos := getCurrentPos()
		next, err := getNextToken()
		if err != nil || !next.IsAnAdditiveOperator() {
			resetHeadTo(pos)
			return kind, nil
		}
		operator := next.Kind
		anotherKind, anotherErr := analyzeMultiplicativeExpression()
		if anotherErr != nil {
			resetHeadTo(pos)
			return 0, anotherErr
		}

		kind = computeType(kind, anotherKind, previousOffset)
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

func analyzeCastExpressionHelper(pos int, kindStack []int) (int, *Error) {
	unaryKind, anotherErr := analyzeUnaryExpression()
	if anotherErr != nil {
		resetHeadTo(pos)
		return 0, anotherErr
	}
	kind := 0
	if len(kindStack) == 0 {
		kind = unaryKind
	} else {
		lastKind := unaryKind
		for len(kindStack) > 0 {
			lastPos := len(kindStack) - 1
			kind = kindStack[lastPos]
			if kind != lastKind {
				convertType(lastKind, kind)
			}
			lastKind = kind
			kindStack = kindStack[0:lastPos]
		}
	}
	return kind, nil
}

func analyzeCastExpression() (int, *Error) {
	// <cast-expression> ::= {'('<type-specifier>')'}<unary-expression>
	pos := getCurrentPos()
	kind := 0
	next, err := getNextToken()
	typeStack := []int{}
	if err != nil {
		resetHeadTo(pos)
		return 0, cc0_error.Of(cc0_error.IncompleteExpression)
	}
	if next.Kind == token.LeftParenthesis {
		next, err := getNextToken()
		if err != nil || !next.IsATypeSpecifier() {
			resetHeadTo(pos)
			return analyzeCastExpressionHelper(pos, typeStack)
		}
		kind = next.Kind
		typeStack = append(typeStack, kind)

		next, err = getNextToken()
		if err != nil || next.Kind != token.RightParenthesis {
			return 0, cc0_error.Of(cc0_error.IncompleteExpression)
		}
		for {
			anotherPos := getCurrentPos()
			next, err = getNextToken()
			if err != nil {
				resetHeadTo(anotherPos)
				return 0, cc0_error.Of(cc0_error.IncompleteExpression)
			}
			if next.Kind != token.LeftParenthesis {
				resetHeadTo(anotherPos)
				break
			}
			next, err = getNextToken()
			if err != nil || !next.IsATypeSpecifier() {
				resetHeadTo(anotherPos)
				break
			}
			nextKind := next.Kind
			typeStack = append(typeStack, nextKind)

			next, err = getNextToken()
			if err != nil || next.Kind != token.RightParenthesis {
				resetHeadTo(anotherPos)
				return 0, cc0_error.Of(cc0_error.IncompleteExpression)
			}
		}
	} else {
		resetHeadTo(pos)
	}
	return analyzeCastExpressionHelper(pos, typeStack)
}

func analyzeMultiplicativeExpression() (int, *Error) {
	// <multiplicative-expression> ::= <cast-expression>{<multiplicative-operator><cast-expression>}

	// <cast-expression>
	kind, err := analyzeCastExpression()
	if err != nil && kind == 0 {
		return 0, err
	}

	// {<multiplicative-operator><unary-expression>}
	for {
		pos := getCurrentPos()
		next, err := getNextToken()
		if err != nil || !next.IsAMultiplicativeOperator() {
			resetHeadTo(pos)
			return kind, nil
		}
		currentFunction.Append(instruction.Nop)
		previousOffset := currentFunction.GetCurrentOffset() - 1
		anotherKind, anotherErr := analyzeUnaryExpression()
		if err != nil {
			resetHeadTo(pos)
			return 0, anotherErr
		}
		kind = computeType(kind, anotherKind, previousOffset)

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
		if kind == token.Double {
			currentFunction.Append(instruction.Dneg)
		} else {
			currentFunction.Append(instruction.Ineg)
		}
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
	} else if next.Kind == token.DoubleLiteral {
		address := globalSymbolTable.AddALiteral(instruction.ConstantKindDouble, next.Value)
		currentFunction.Append(instruction.Loadc, -address)
		kind = token.Double
	} else if next.Kind == token.CharLiteral {
		currentFunction.Append(instruction.Bipush, int(next.Value.(int32)))
		kind = token.Char
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

	// pre read
	preReadPos := getCurrentPos()
	theOneAfterNext, err := getNextToken()
	resetHeadTo(preReadPos)
	if err == nil && theOneAfterNext.Kind != token.AssignmentSign {
		return cc0_error.Of(cc0_error.IncompleteExpression)
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
	if currentVariableKind := currentSymbolTable.GetSymbolNamed(identifier).Kind; kind != currentVariableKind {
		convertType(kind, currentSymbolTable.GetSymbolNamed(identifier).Kind)
		kind = currentVariableKind
	}
	if anotherErr != nil {
		resetHeadTo(pos)
		return anotherErr
	}
	if kind == token.Double {
		currentFunction.Append(instruction.Dstore)
	} else {
		currentFunction.Append(instruction.Istore)
	}
	return nil
}

var currentFnTotalParams = 0
var currentDeclaredCount = 0

func analyzeExpressionList(fn *instruction.Fn) *Error {
	// <expression-list> ::= <expression>{','<expression>}
	currentDeclaredCount = 0
	currentFnTotalParams = len(*fn.Parameters)

	pos := getCurrentPos()
	kind, err := analyzeExpression()
	if err != nil {
		resetHeadTo(pos)
		return err
	}
	currentDeclaredCount++
	if currentDeclaredCount > currentFnTotalParams {
		return cc0_error.Of(cc0_error.IllegalExpression).On(currentLine, currentColumn)
	}
	if paramKind := fn.RelatedSymbolTable.GetSymbolNamed((*fn.Parameters)[currentDeclaredCount-1]).Kind; kind != paramKind {
		convertType(kind, paramKind)
	}

	for {
		pos = getCurrentPos()
		if next, err := getNextToken(); err != nil || next.Kind != token.Comma {
			resetHeadTo(pos)
			if currentDeclaredCount != currentFnTotalParams {
				return cc0_error.Of(cc0_error.IllegalExpression).On(currentLine, currentColumn)
			}
			return nil
		}
		anotherKind, err := analyzeExpression()
		if err != nil {
			return err
		}
		currentDeclaredCount++
		if currentDeclaredCount > currentFnTotalParams {
			return cc0_error.Of(cc0_error.IllegalExpression).On(currentLine, currentColumn)
		}
		if paramKind := fn.RelatedSymbolTable.GetSymbolNamed((*fn.Parameters)[currentDeclaredCount-1]).Kind; anotherKind != paramKind {
			convertType(anotherKind, paramKind)
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
	_ = analyzeExpressionList(sb.FnInfo)
	if next, err := getNextToken(); err != nil || next.Kind != token.RightParenthesis {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.IncompleteFunctionCall).On(currentLine, currentColumn)
	}
	currentFunction.Append(instruction.Call, sb.Address)
	if currentDeclaredCount != currentFnTotalParams {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.IllegalExpression).On(currentLine, currentColumn)
	}
	return nil
}
