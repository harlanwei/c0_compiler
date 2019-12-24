package analyzer

import (
	"c0_compiler/internal/cc0_error"
	"c0_compiler/internal/instruction"
	"c0_compiler/internal/token"
)

var currentInitializationType int

func analyzeDeclaratorList(isConstant bool) *Error {
	// <init-declarator-list> ::= <init-declarator>{','<init-declarator>}
	if err := analyzeInitDeclarator(isConstant); err != nil {
		return err
	}
	for {
		pos := getCurrentPos()
		next, err := getNextToken()
		if !(err == nil && next.Kind == token.Comma) {
			resetHeadTo(pos)
			return nil
		}
		if err := analyzeInitDeclarator(isConstant); err != nil {
			return err
		}
	}
}

func analyzeVariableDeclarations() *Error {
	for {
		pos := getCurrentPos()
		next, err := getNextToken()
		if err != nil || (!next.IsATypeSpecifier() && next.Kind != token.Const) {
			resetHeadTo(pos)
			return nil
		}
		if next.Kind == token.Const {
			if next, err := getNextToken(); err != nil || !next.IsATypeSpecifier() {
				resetHeadTo(pos)
				return nil
			}
		}
		if next, err := getNextToken(); err != nil || next.Kind != token.Identifier {
			resetHeadTo(pos)
			return nil
		}
		if next, err := getNextToken(); err != nil || (next.Kind != token.AssignmentSign &&
			next.Kind != token.Comma && next.Kind != token.Semicolon) {
			resetHeadTo(pos)
			return nil
		}
		resetHeadTo(pos)
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
	if !next.IsATypeSpecifier() {
		resetHeadTo(pos)
		return cc0_error.Of(cc0_error.IncompleteVariableDeclaration).On(currentLine, currentColumn)
	}
	currentInitializationType = next.Kind

	// <init-declarator-list>
	if err := analyzeDeclaratorList(isAConstant); err != nil {
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

func analyzeInitDeclarator(isConstant bool) *Error {
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
		if err := currentSymbolTable.AddAConstant(identifier, currentInitializationType); err != nil {
			return err
		}
	} else {
		if err := currentSymbolTable.AddAVariable(identifier, currentInitializationType); err != nil {
			return err
		}
	}

	// <initializer> ::= '='<expression>
	pos = getCurrentPos()
	next, err = getNextToken()
	if err != nil {
		if isConstant {
			return cc0_error.Of(cc0_error.IncompleteExpression).On(currentLine, currentColumn)
		}
		return nil
	}
	if next.Kind != token.AssignmentSign {
		resetHeadTo(pos)
		if isConstant {
			return cc0_error.Of(cc0_error.IncompleteExpression).On(currentLine, currentColumn)
		}
		if currentInitializationType == token.Double {
			currentFunction.Append(instruction.Snew, 2)
		} else {
			currentFunction.Append(instruction.Snew, 1)
		}
		return nil
	}

	address := currentSymbolTable.GetAddressOf(identifier)
	if currentKind := currentSymbolTable.GetSymbolNamed(identifier).Kind; currentKind == token.Double {
		currentFunction.Append(instruction.Snew, 2)
	} else {
		currentFunction.Append(instruction.Ipush, 0)
	}
	currentFunction.Append(instruction.Loada, currentSymbolTable.GetLevelDiff(identifier), address)
	kind, anotherErr := analyzeExpression()
	if anotherErr != nil {
		return anotherErr
	}

	if kind != currentInitializationType {
		convertType(kind, currentInitializationType)
	}

	switch currentInitializationType {
	case token.Void:
		cc0_error.ThrowAndExit(cc0_error.Parser)
	case token.Double:
		currentFunction.Append(instruction.Dstore)
	case token.Int, token.Char:
		currentFunction.Append(instruction.Istore)
	}
	return nil
}
