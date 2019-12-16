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
