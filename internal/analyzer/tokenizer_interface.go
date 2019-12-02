package analyzer

import (
	"c0_compiler/internal/parser"
	"c0_compiler/internal/token"
	"fmt"
)

type Parser = parser.Parser
type Token = token.Token

var globalLineCount = 0
var globalColumnCount = 0

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

// TODO: change `putBackAToken` to `resetTo` for easier implementations
func putBackAToken() *Error {
	if prev := globalParser.UnreadToken(); prev != nil {
		globalLineCount, globalColumnCount = prev.Line, prev.Column
		return nil
	}
	return errorOf(Bug)
}

// Print all the tokens the parser generated directly to stdout.
func BindToStdOut(parser *Parser) {
	for parser.HasNextToken() {
		fmt.Println(parser.NextToken())
	}
}
