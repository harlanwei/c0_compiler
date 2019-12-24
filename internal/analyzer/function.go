package analyzer

import (
	"c0_compiler/internal/cc0_error"
	"c0_compiler/internal/instruction"
	"c0_compiler/internal/token"
)

func analyzeFunctionDefinitions() *Error {
	for {
		pos := getCurrentPos()
		next, err := getNextToken()
		resetHeadTo(pos)
		if err != nil || !next.IsATypeSpecifier() {
			return nil
		}
		if err := analyzeFunctionDefinition(); err != nil {
			resetHeadTo(pos)
			return err
		}
	}
}

func analyzeFunctionDefinition() *Error {
	// <function-definition> ::= <type-specifier><identifier><parameter-clause><compound-statement>

	pos := getCurrentPos()
	next, err := getNextToken()
	if err != nil || !next.IsATypeSpecifier() {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.InvalidDeclaration).On(currentLine, currentColumn)
	}
	kind := next.Kind
	currentFunction = instruction.InitFn(kind)
	currentSymbolTable = currentSymbolTable.AppendChildSymbolTable(currentFunction)

	next, err = getNextToken()
	if err != nil || next.Kind != token.Identifier {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.InvalidDeclaration).On(currentLine, currentColumn)
	}
	identifier := next.Value.(string)
	if err := globalSymbolTable.AddAFunction(identifier, kind, currentFunction); err != nil {
		resetHeadTo(pos)
		return err
	}
	if err := analyzeParameterClause(); err != nil {
		resetHeadTo(pos)
		return err
	}
	if err := analyzeCompoundStatement(); err != nil {
		resetHeadTo(pos)
		return err
	}

	switch currentFunction.ReturnType {
	case token.Void:
		currentFunction.Append(instruction.Ret)
	case token.Int, token.Char:
		currentFunction.Append(instruction.Ipush, 0)
		currentFunction.Append(instruction.Iret)
	case token.Double:
		currentFunction.Append(instruction.Snew, 2)
		currentFunction.Append(instruction.Dret)
	}

	if funSymbol := globalSymbolTable.GetSymbolNamed(identifier); funSymbol != nil {
		// Have to do the assignment this way thanks to all the trivia of golang
		funSymbol.FnInfo = currentFunction
		globalSymbolTable.Symbols[identifier] = funSymbol
	}

	currentFunction = globalStart
	currentSymbolTable = globalSymbolTable
	return nil
}

func analyzeParameterClause() *Error {
	// <parameter-clause> ::= '(' [<parameter-declaration-list>] ')'
	pos := getCurrentPos()
	next, err := getNextToken()
	if err != nil || next.Kind != token.LeftParenthesis {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.InvalidDeclaration).On(currentLine, currentColumn)
	}

	next, err = getNextToken()
	if err != nil {
		return cc0_error.Of(cc0_error.InvalidStatement).On(currentLine, currentColumn)
	}
	if next.Kind == token.RightParenthesis {
		return nil
	}

	// Put back the token previously read in
	resetHeadTo(pos + 1)

	if err := analyzeParameterDeclarationList(); err != nil {
		resetHeadTo(pos)
		return err
	}
	next, err = getNextToken()
	if err != nil || next.Kind != token.RightParenthesis {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.IncompleteExpression)
	}

	return nil
}

func analyzeParameterDeclarationList() *Error {
	// <parameter-declaration-list> ::= <parameter-declaration>{','<parameter-declaration>}

	// <parameter-declaration>
	if err := analyzeParameterDeclaration(); err != nil {
		return err
	}

	// {','<parameter-declaration>}
	for {
		pos := getCurrentPos()
		next, err := getNextToken()
		if err != nil {
			return nil
		}
		if next.Kind != token.Comma {
			resetHeadTo(pos)
			return nil
		}
		if err := analyzeParameterDeclaration(); err != nil {
			resetHeadTo(pos)
			return err
		}
	}
}

func analyzeParameterDeclaration() *Error {
	// [<const-qualifier>]<type-specifier><identifier>
	next, err := getNextToken()
	if err != nil {
		return cc0_error.Of(cc0_error.InvalidDeclaration).On(currentLine, currentColumn)
	}
	isConst := false
	if next.Kind == token.Const {
		isConst = true
		next, err = getNextToken()
		if err != nil {
			return cc0_error.Of(cc0_error.InvalidDeclaration).On(currentLine, currentColumn)
		}
	}
	if !next.IsATypeSpecifier() || next.Kind == token.Void {
		return cc0_error.Of(cc0_error.InvalidDeclaration).On(currentLine, currentColumn)
	}
	kind := next.Kind
	next, err = getNextToken()
	if err != nil || next.Kind != token.Identifier {
		return cc0_error.Of(cc0_error.InvalidDeclaration).On(currentLine, currentColumn)
	}
	identifier := next.Value.(string)

	if isConst {
		_ = currentSymbolTable.AddAConstant(identifier, kind)
	} else {
		_ = currentSymbolTable.AddAVariable(identifier, kind)
	}
	*currentFunction.Parameters = append(*currentFunction.Parameters, identifier)
	return nil
}
