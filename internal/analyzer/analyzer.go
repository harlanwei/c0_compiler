package analyzer

import (
	"bufio"
	"c0_compiler/internal/parser"
	"c0_compiler/internal/token"
	"c0_compiler/internal/vm"
	"fmt"
)

type Parser = parser.Parser

type Token = token.Token

var globalParser *Parser
var globalWriter *bufio.Writer
var globalLineCount = 0
var globalColumnCount = 0
var globalSymbolTable, currentSymbolTable *symbolTable

// Print all the tokens the parser generated directly to stdout.
func BindToStdOut(parser *Parser) {
	for parser.HasNextToken() {
		fmt.Println(parser.NextToken())
	}
}

// The only error this function will throw is NoMoreTokens so it's safe to check `err != nil` directly without
// specifying the kind of error.
func getNextToken() (res *Token, err error) {
	if !globalParser.HasNextToken() {
		res, err = nil, &Error{NoMoreTokens, globalLineCount, globalColumnCount}
		return
	}
	res, err = globalParser.NextToken(), nil
	globalLineCount, globalColumnCount = res.Line, res.Column
	return
}

func putBackAToken() *Error {
	if prev := globalParser.UnreadToken(); prev != nil {
		globalLineCount, globalColumnCount = prev.Line, prev.Column
	}
	return errorOf(Bug)
}

func analyzePrimaryExpression() *Error {
	// <primary-expression> ::=
	//     '('<expression>')'
	//    | <identifier>
	//    | <integer-literal>
	//    | <function-call>

	next, err := getNextToken()
	if err != nil {
		return errorOf(IncompleteExpression)
	}

	if next.Kind == token.LeftParenthesis {
		// '('<expression>')'
		if err := analyzeExpression(); err != nil {
			return err
		}
		next, err = getNextToken()
		if err != nil {
			return errorOf(IncompleteExpression)
		}
		if next.Kind != token.RightParenthesis {
			return errorOf(IllegalExpression)
		}
	} else if next.Kind == token.Identifier {
		// <identifier> || <function-call>

		another, err := getNextToken()

		if err != nil || another.Kind != token.LeftParenthesis {
			if err := putBackAToken(); err != nil {
				return err
			}
			// TODO: <identifier>
		} else {
			// TODO: <function-call>
		}
	} else if next.Kind == token.IntegerLiteral {
		// <integer-literal>
		index := vm.AddConstant(next.Value)
		vm.AddInstruction(vm.Loadc, index)
	}
	return errorOf(IllegalExpression)
}

func isUnaryOperator(t *Token) bool {
	return t.Kind == token.PlusSign || t.Kind == token.MinusSign
}

func analyzeUnaryExpression() *Error {
	// <unary-expression> ::= [<unary-operator>]<primary-expression>

	// [<unary-operator>]
	shouldBeNegated := false
	next, err := getNextToken()
	if err != nil {
		return errorOf(IncompleteExpression)
	}
	if isUnaryOperator(next) {
		if next.Kind == token.MinusSign {
			shouldBeNegated = true
		}
	} else if err := putBackAToken(); err != nil {
		return err
	}

	// <primary-expression>
	if err := analyzePrimaryExpression(); err != nil {
		return err
	}

	if shouldBeNegated {
		vm.AddInstruction(vm.Ineg)
	}
	return nil
}

func isMultiplicativeOperator(t *Token) bool {
	return t.Kind == token.MultiplicationSign || t.Kind == token.DivisionSign
}

func analyzeMultiplicativeExpression() *Error {
	// <multiplicative-expression> ::= <unary-expression>{<multiplicative-operator><unary-expression>}

	// <unary-expression>
	if err := analyzeUnaryExpression(); err != nil {
		return err
	}

	// {<multiplicative-operator><unary-expression>}
	for {
		next, err := getNextToken()
		if err != nil || !isMultiplicativeOperator(next) {
			if err := putBackAToken(); err != nil {
				return err
			}
			if next.Kind == token.MultiplicationSign {
				vm.AddInstruction(vm.Imul)
			} else {
				vm.AddInstruction(vm.Idiv)
			}
			return nil
		}
		if err := analyzeUnaryExpression(); err != nil {
			return err
		}
	}
}

func isAdditiveOperator(t *Token) bool {
	return t.Kind == token.PlusSign || t.Kind == token.MinusSign
}

func analyzeAdditiveExpression() *Error {
	// <additive-expression> ::= <multiplicative-expression>{<additive-operator><multiplicative-expression>}

	// <multiplicative-expression>
	if err := analyzeMultiplicativeExpression(); err != nil {
		return err
	}

	// {<additive-operator><multiplicative-expression>}
	for {
		next, err := getNextToken()
		if err != nil || !isAdditiveOperator(next) {
			if err := putBackAToken(); err != nil {
				return err
			}
			if next.Kind == token.PlusSign {
				vm.AddInstruction(vm.Iadd)
			} else {
				vm.AddInstruction(vm.Isub)
			}
			return nil
		}
		if err := analyzeMultiplicativeExpression(); err != nil {
			return err
		}
	}
}

func analyzeExpression() *Error {
	// <expression> ::= <additive-expression>
	return analyzeAdditiveExpression()
}

func analyzeDeclarator(isConstant bool, declaredType int) *Error {
	// <init-declarator> ::= <identifier>[<initializer>]

	// <identifier>
	next, err := getNextToken()
	if err != nil {
		return errorOf(IncompleteVariableDeclaration)
	}
	if next.Kind != token.Identifier {
		return errorOf(InvalidDeclaration)
	}

	if isConstant {
		if err := currentSymbolTable.addVariable(next.Value.(string), declaredType); err != nil {
			return err
		}
	} else {
		if err := currentSymbolTable.addConstant(next.Value.(string), declaredType); err != nil {
			return err
		}
	}

	// <initializer> ::= '='<expression>
	next, err = getNextToken()
	if err != nil {
		return nil
	}
	if next.Kind != token.AssignmentSign {
		if err := putBackAToken(); err != nil {
			return err
		}
		return nil
	}
	if err := analyzeExpression(); err != nil {
		return err
	}
	return nil
}

func analyzeDeclaratorList(isConstant bool, declaredType int) *Error {
	// <init-declarator-list> ::= <init-declarator>{','<init-declarator>}
	if err := analyzeDeclarator(isConstant, declaredType); err != nil {
		return err
	}
	for {
		next, err := getNextToken()
		if err != nil {
			return nil
		}
		if next.Kind != token.Comma {
			return errorOf(InvalidDeclaration)
		}
		if err := analyzeDeclarator(isConstant, declaredType); err != nil {
			return err
		}
	}
}

func analyzeVariableDeclaration() *Error {
	// <variable-declaration> ::= [<const-qualifier>]<type-specifier><init-declarator-list>';'

	// [<const-qualifier>]
	next, err := getNextToken()
	if err != nil {
		return errorOf(IncompleteVariableDeclaration)
	}
	isAConstant := false
	if next.Kind == token.Const {
		isAConstant = true
		next, err = getNextToken()
		if err != nil {
			return errorOf(IncompleteVariableDeclaration)
		}
	}

	// <type-specifier>
	// might need a little expansion
	if next.Kind != token.Int {
		return errorOf(IncompleteVariableDeclaration)
	}

	// <init-declarator-list>
	if err := analyzeDeclaratorList(isAConstant, TypeInt); err != nil {
		return err
	}

	// ;
	if next, err = getNextToken(); err != nil || next.Kind != token.Semicolon {
		return errorOf(IncompleteVariableDeclaration)
	}
	return nil
}

func analyzeVariableDeclarations() *Error {
	for {
		next, err := getNextToken()
		if backErr := putBackAToken(); backErr != nil {
			return backErr
		}
		if err != nil || (next.Kind != token.Identifier && next.Kind != token.Const) {
			return nil
		}
		if err := analyzeVariableDeclaration(); err != nil {
			return err
		}
	}
}

func analyzeCompoundStatement() *Error {
	// TODO
	return nil
}

func analyzeParameterClause() *Error {
	// TODO
	return nil
}

func analyzeFunctionDefinition() *Error {
	// <function-definition> ::= <type-specifier><identifier><parameter-clause><compound-statement>

	// TODO
	return nil
}

func analyzeFunctionDefinitions() *Error {
	// TODO
	return nil
}

func Run(parser *Parser, writer *bufio.Writer, shouldCompileToBinary bool) {
	globalParser, globalWriter = parser, writer
	globalSymbolTable = initGlobalSymbolTable()
	currentSymbolTable = globalSymbolTable

	if err := analyzeVariableDeclarations(); err != nil {
		reportFatalError(err)
	}
	if err := analyzeFunctionDefinitions(); err != nil {
		reportFatalError(err)
	}
}
