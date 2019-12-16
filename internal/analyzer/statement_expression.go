package analyzer

import (
	"c0_compiler/internal/cc0_error"
	"c0_compiler/internal/instruction"
	"c0_compiler/internal/token"
)

func analyzeDeclaratorList(isConstant bool, declaredType int) *Error {
	// <init-declarator-list> ::= <init-declarator>{','<init-declarator>}
	if err := analyzeInitDeclarator(isConstant, declaredType); err != nil {
		return err
	}
	for {
		pos := getCurrentPos()
		next, err := getNextToken()
		if !(err == nil && next.Kind == token.Comma) {
			resetHeadTo(pos)
			return nil
		}
		if err := analyzeInitDeclarator(isConstant, declaredType); err != nil {
			return err
		}
	}
}

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
	currentFunction.Append(instruction.Ipush, address)
	if err := analyzeExpression(); err != nil {
		resetHeadTo(pos)
		return err
	}
	currentFunction.Append(instruction.Istore)
	return nil
}

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
		currentFunction.Append(instruction.Loada, 0, currentSymbolTable.GetAddressOf(identifier))
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
	// <printable> ::= <expression>
	if err := analyzeExpression(); err != nil {
		return err
	}
	currentFunction.Append(instruction.Iprint)
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
	_ = analyzeExpression()
	if next, err := getNextToken(); err != nil || next.Kind != token.Semicolon {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.InvalidStatement).On(currentLine, currentColumn)
	}
	// TODO: generate instructions
	currentFunction.Append(instruction.Ret)
	return nil
}

func analyzeVariableDeclarations() *Error {
	for {
		pos := getCurrentPos()
		next, err := getNextToken()
		resetHeadTo(pos)
		if err != nil || (next.Kind != token.Int && next.Kind != token.Const) {
			return nil
		}
		if err := analyzeVariableDeclaration(); err != nil {
			return err
		}
	}
}

func analyzeVariableDeclaration() *Error {
	// <variable-declaration> ::= [<const-qualifier>]<type-specifier><init-declarator-list>';'

	// [<const-qualifier>]
	pos := getCurrentPos()
	next, err := getNextToken()
	if err != nil {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.IncompleteVariableDeclaration).On(currentLine, currentColumn)
	}
	isAConstant := false
	if next.Kind == token.Const {
		isAConstant = true
		next, err = getNextToken()
		if err != nil {
			return cc0_error.Of(cc0_error.IncompleteVariableDeclaration).On(currentLine, currentColumn)
		}
	}

	// <type-specifier>
	if next.Kind != token.Int {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.IncompleteVariableDeclaration).On(currentLine, currentColumn)
	}
	kind := tokenKindToType(next)

	// <init-declarator-list>
	if err := analyzeDeclaratorList(isAConstant, kind); err != nil {
		resetHeadTo(pos)
		return err
	}

	// ;
	if next, err = getNextToken(); err != nil || next.Kind != token.Semicolon {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.IncompleteVariableDeclaration).On(currentLine, currentColumn)
	}
	return nil
}

func analyzeInitDeclarator(isConstant bool, declaredType int) *Error {
	// <init-declarator> ::= <identifier>[<initializer>]

	// <identifier>
	pos := getCurrentPos()
	next, err := getNextToken()
	if err != nil {
		return cc0_error.Of(cc0_error.IncompleteVariableDeclaration).On(currentLine, currentColumn)
	}
	if next.Kind != token.Identifier {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.InvalidDeclaration).On(currentLine, currentColumn)
	}
	identifier := next.Value.(string)

	if isConstant {
		if err := currentSymbolTable.AddAConstant(identifier, declaredType); err != nil {
			return err
		}
	} else {
		if err := currentSymbolTable.AddAVariable(identifier, declaredType); err != nil {
			return err
		}
	}

	// <initializer> ::= '='<expression>
	pos = getCurrentPos()
	next, err = getNextToken()
	if err != nil {
		return nil
	}
	if next.Kind != token.AssignmentSign {
		currentFunction.Append(instruction.Snew, 1)
		resetHeadTo(pos)
		return nil
	}

	address := currentSymbolTable.GetAddressOf(identifier)
	currentFunction.Append(instruction.Ipush, 0)
	currentFunction.Append(instruction.Loada, 0, address)
	if err := analyzeExpression(); err != nil {
		return err
	}
	currentFunction.Append(instruction.Istore)
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
