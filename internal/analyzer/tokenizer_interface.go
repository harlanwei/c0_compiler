package analyzer

import (
	"c0_compiler/internal/cc0_error"
	"c0_compiler/internal/parser"
	"c0_compiler/internal/token"
)

type Parser = parser.Parser
type Token = token.Token

var currentLine = 0
var currentColumn = 0

// The only error this function will throw is NoMoreTokens so it's safe to check `err != nil` directly without
// specifying the kind of error.
func getNextToken() (res *Token, err error) {
	if !globalParser.HasNextToken() {
		res, err = nil, cc0_error.Of(cc0_error.NoMoreTokens).On(currentLine, currentColumn)
		return
	}
	res, err = globalParser.NextToken(), nil
	currentLine, currentColumn = res.Line, res.Column
	return
}

func getCurrentPos() int {
	return globalParser.CurrentHead()
}

func resetHeadTo(pos int) {
	thatToken := globalParser.ResetHeadTo(pos)
	currentColumn, currentLine = thatToken.Column, thatToken.Line
}
