package analyzer

import (
	"c0_compiler/internal/cc0_error"
	"c0_compiler/internal/instruction"
	"c0_compiler/internal/token"
)

type SymbolTable = instruction.SymbolTable
type Error = cc0_error.Error

var globalParser *Parser
var globalSymbolTable, currentSymbolTable *SymbolTable
var globalStart, currentFunction *instruction.Fn

func Run(parser *Parser) *SymbolTable {
	globalParser = parser
	globalStart = instruction.InitFn(token.Void)
	globalSymbolTable = instruction.InitSymbolTable(nil, globalStart)
	currentSymbolTable = globalSymbolTable
	currentFunction = globalStart

	if err := analyzeVariableDeclarations(); err != nil {
		err.DieAndReportPosition(cc0_error.Analyzer)
	}
	if err := analyzeFunctionDefinitions(); err != nil {
		err.DieAndReportPosition(cc0_error.Analyzer)
	}

	return globalSymbolTable
}
