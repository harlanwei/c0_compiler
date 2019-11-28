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

func putBackTokens() *Error {
	// TODO: refactor this so that multiple tokens read by the current analysis function can be put back
	// TODO: also, move error handling directly into this function instead of having each individual function do it
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

	// '('<expression>')'
	if next.Kind == token.LeftParenthesis {
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
	} else if next.Kind == token.Identifier { // <identifier> || <function-call>
		another, err := getNextToken()

		if err != nil || another.Kind != token.LeftParenthesis {
			_ = putBackTokens()
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
	} else {
		_ = putBackTokens()
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
			_ = putBackTokens()
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
			_ = putBackTokens()
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
		_ = putBackTokens()
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
		_ = putBackTokens()
		if err != nil || (next.Kind != token.Identifier && next.Kind != token.Const) {
			return nil
		}
		if err := analyzeVariableDeclaration(); err != nil {
			return err
		}
	}
}

func analyzeReturnStatement() *Error {
	// <return-statement> ::= 'return' [<expression>] ';'

	if next, err := getNextToken(); err != nil || next.Kind != token.Return {
		if err == nil {
			_ = putBackTokens()
		}
		return errorOf(InvalidStatement)
	}
	_ = analyzeExpression()
	if next, err := getNextToken(); err != nil || next.Kind != token.Semicolon {
		if err == nil {
			_ = putBackTokens()
		}
		return errorOf(InvalidStatement)
	}
	return nil
}

func analyzeJumpStatement() *Error {
	// <jump-statement> ::= <return-statement>
	return analyzeReturnStatement()
}

func analyzePrintable() *Error {
	// <printable> ::= <expression>
	return analyzeExpression()
}

func analyzePrintableList() *Error {
	// <printable-list>  ::= <printable> {',' <printable>}
	if err := analyzePrintable(); err != nil {
		return err
	}
	for {
		if next, err := getNextToken(); err != nil || next.Kind != token.Comma {
			if err == nil {
				_ = putBackTokens()
			}
			return nil
		}
		if err := analyzePrintable(); err != nil {
			return err
		}
	}
}

func analyzeIOStatement() *Error {
	// <scan-statement>  ::= 'scan' '(' <identifier> ')' ';'
	// <print-statement> ::= 'print' '(' [<printable-list>] ')' ';'
	next, err := getNextToken()
	if err != nil {
		return errorOf(InvalidStatement)
	}
	if next.Kind == token.Scan {
		if next, err := getNextToken(); err != nil || next.Kind != token.LeftParenthesis {
			return errorOf(InvalidStatement)
		}
		next, err = getNextToken()
		if err != nil || next.Kind != token.Identifier {
			return errorOf(InvalidStatement)
		}
		// identifier := next.Value.(string)
		// TODO: generate expression
		if next, err := getNextToken(); err != nil || next.Kind != token.RightParenthesis {
			return errorOf(InvalidStatement)
		}
	} else if next.Kind == token.Print {
		if next, err := getNextToken(); err != nil || next.Kind != token.LeftParenthesis {
			return errorOf(InvalidStatement)
		}
		if err := analyzePrintableList(); err != nil {
			return err
		}
		// TODO: generate expression
		if next, err := getNextToken(); err != nil || next.Kind != token.RightParenthesis {
			return errorOf(InvalidStatement)
		}
	} else {
		return errorOf(InvalidStatement)
	}
	if next, err := getNextToken(); err != nil || next.Kind != token.Semicolon {
		return errorOf(InvalidStatement)
	}
	return nil
}

func analyzeLoopStatement() *Error {
	// <loop-statement> ::= 'while' '(' <condition> ')' <statement>
	// TODO: generate instructions
	next, err := getNextToken()
	if err != nil || next.Kind != token.While {
		return errorOf(InvalidStatement)
	}
	next, err = getNextToken()
	if err != nil || next.Kind != token.LeftParenthesis {
		return errorOf(InvalidStatement)
	}
	if err := analyzeCondition(); err != nil {
		return err
	}
	next, err = getNextToken()
	if err != nil || next.Kind != token.RightParenthesis {
		return errorOf(InvalidStatement)
	}
	return nil
}

func analyzeCondition() *Error {
	// <condition> ::= <expression>[<relational-operator><expression>]

	if err := analyzeExpression(); err != nil {
		return err
	}
	next, err := getNextToken()
	if err != nil {
		return nil
	}
	if !next.IsARelationalOperator() {
		_ = putBackTokens()
		return nil
	}
	if err := analyzeExpression(); err != nil {
		return err
	}
	// TODO: generate instructions
	return nil
}

func analyzeConditionStatement() *Error {
	// <condition-statement> ::=  'if' '(' <condition> ')' <statement> ['else' <statement>]
	// TODO: generate instructions

	if next, err := getNextToken(); err != nil || next.Kind != token.If {
		return errorOf(InvalidStatement)
	}
	if next, err := getNextToken(); err != nil || next.Kind != token.LeftParenthesis {
		return errorOf(InvalidStatement)
	}
	if err := analyzeCondition(); err != nil {
		return err
	}
	if next, err := getNextToken(); err != nil || next.Kind != token.RightParenthesis {
		return errorOf(InvalidStatement)
	}
	if err := analyzeStatement(); err != nil {
		return err
	}
	if next, err := getNextToken(); err != nil || next.Kind != token.Else {
		// Only when a token is actually read is there the need to put back a token
		if err == nil {
			_ = putBackTokens()
		}
		return nil
	}
	if err := analyzeStatement(); err != nil {
		return err
	}

	return nil
}

func analyzeAssignmentExpression() *Error {
	// <assignment-expression> ::= <identifier><assignment-operator><expression>
	next, err := getNextToken()
	if err != nil {
		return errorOf(IncompleteExpression)
	}
	if next.Kind != token.Identifier {
		_ = putBackTokens()
		return errorOf(IncompleteExpression)
	}
	if next, err := getNextToken(); err != nil || next.Kind != token.AssignmentSign {
		return errorOf(IncompleteExpression)
	}
	if err := analyzeExpression(); err != nil {
		return err
	}
	return nil
}

func analyzeExpressionList() *Error {
	// <expression-list> ::= <expression>{','<expression>}
	if err := analyzeExpression(); err != nil {
		return err
	}
	for {
		if next, err := getNextToken(); err != nil || next.Kind != token.Comma {
			if err == nil {
				_ = putBackTokens()
			}
			return nil
		}
		if err := analyzeExpression(); err != nil {
			return err
		}
	}
}

func analyzeFunctionCall() *Error {
	// <identifier> '(' [<expression-list>] ')'
	next, err := getNextToken()
	if err != nil {
		return errorOf(IncompleteFunctionCall)
	}
	if next.Kind != token.Identifier {
		_ = putBackTokens()
		return errorOf(IncompleteFunctionCall)
	}
	identifier := next.Value.(string)
	if !currentSymbolTable.hasDeclared(identifier) {
		return errorOf(UndefinedIdentifier)
	}
	if next, err := getNextToken(); err != nil || next.Kind != token.LeftParenthesis {
		return errorOf(IncompleteFunctionCall)
	}
	_ = analyzeExpressionList()
	if next, err := getNextToken(); err != nil || next.Kind != token.RightParenthesis {
		return errorOf(IncompleteFunctionCall)
	}

	return nil
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
	next, err := getNextToken()
	if err == nil && next.Kind == token.LeftBracket {
		if err := analyzeStatementSeq(); err != nil {
			return err
		}
		if next, err := getNextToken(); err != nil || next.Kind != token.RightBracket {
			return errorOf(InvalidStatement)
		}
		return nil
	}
	_ = putBackTokens()

	// <condition-statement>
	if err := analyzeConditionStatement(); err == nil {
		return nil
	}

	// <loop-statement>
	if err := analyzeLoopStatement(); err == nil {
		return nil
	}

	// <jump-statement>
	if err := analyzeJumpStatement(); err == nil {
		return nil
	}

	// <print-statement> | <scan-statement>
	if err := analyzeIOStatement(); err == nil {
		return nil
	}

	// <assignment-expression>';'
	if err := analyzeAssignmentExpression(); err == nil {
		if next, err := getNextToken(); err != nil || next.Kind != token.Semicolon {
			return errorOf(InvalidStatement)
		}
		return nil
	}

	// <function-call>';'
	if err := analyzeFunctionCall(); err == nil {
		if next, err := getNextToken(); err != nil || next.Kind != token.Semicolon {
			return errorOf(InvalidStatement)
		}
		return nil
	}

	// ';'
	if next, err := getNextToken(); err != nil || next.Kind != token.Semicolon {
		return errorOf(InvalidStatement)
	}

	return nil
}

func analyzeStatementSeq() *Error {
	// {<statement>}
	for {
		if err := analyzeStatement; err != nil {
			return nil
		}
	}
}

func analyzeCompoundStatement() *Error {
	// '{' {<variable-declaration>} <statement-seq> '}'
	if next, err := getNextToken(); err != nil || next.Kind != token.LeftBracket {
		return errorOf(InvalidStatement)
	}
	if err := analyzeVariableDeclarations(); err != nil {
		return err
	}
	if err := analyzeStatementSeq(); err != nil {
		return err
	}
	if next, err := getNextToken(); err != nil || next.Kind != token.RightBracket {
		return errorOf(InvalidStatement)
	}
	return nil
}

func analyzeParameterDeclaration(functionName string) *Error {
	// [<const-qualifier>]<type-specifier><identifier>
	next, err := getNextToken()
	if err != nil {
		return errorOf(InvalidDeclaration)
	}
	// isConst := false
	if next.Kind == token.Const {
		// isConst = true
		next, err = getNextToken()
		if err != nil {
			return errorOf(InvalidDeclaration)
		}
	}
	if !next.IsATypeSpecifier() {
		return errorOf(InvalidDeclaration)
	}
	// kind := tokenKindToType(next)
	next, err = getNextToken()
	if err != nil || next.Kind != token.Identifier {
		return errorOf(InvalidDeclaration)
	}

	// TODO: generate instructions
	return nil
}

func analyzeParameterDeclarationList(functionName string) *Error {
	// <parameter-declaration-list> ::= <parameter-declaration>{','<parameter-declaration>}

	// <parameter-declaration>
	if err := analyzeParameterDeclaration(functionName); err != nil {
		return err
	}

	// {','<parameter-declaration>}
	for {
		next, err := getNextToken()
		if err != nil {
			return nil
		}
		if next.Kind != token.Comma {
			_ = putBackTokens()
			return nil
		}
		if err := analyzeParameterDeclaration(functionName); err != nil {
			return err
		}
	}
}

func analyzeParameterClause(functionName string) *Error {
	// <parameter-clause> ::= '(' [<parameter-declaration-list>] ')'
	next, err := getNextToken()
	if err != nil || next.Kind != token.LeftParenthesis {
		return errorOf(InvalidDeclaration)
	}
	if err := analyzeParameterDeclarationList(functionName); err != nil {
		return err
	}
	next, err = getNextToken()
	if err != nil || next.Kind != token.RightParenthesis {
		return errorOf(IncompleteExpression)
	}
	return nil
}

func analyzeFunctionDefinition() *Error {
	// <function-definition> ::= <type-specifier><identifier><parameter-clause><compound-statement>

	next, err := getNextToken()
	if err != nil || !next.IsATypeSpecifier() {
		return errorOf(InvalidDeclaration)
	}
	// kind := next.Kind
	next, err = getNextToken()
	if err != nil || next.Kind != token.Identifier {
		return errorOf(InvalidDeclaration)
	}
	identifier := next.Value.(string)
	if err := analyzeParameterClause(identifier); err != nil {
		return err
	}
	if err := analyzeCompoundStatement(); err != nil {
		return err
	}
	// TODO: generate definition
	return nil
}

func analyzeFunctionDefinitions() *Error {
	for {
		next, err := getNextToken()
		_ = putBackTokens()
		if err != nil || !next.IsATypeSpecifier() {
			return nil
		}
		if err := analyzeFunctionDefinition(); err != nil {
			return err
		}
	}
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
