package assembler

import (
	"c0_compiler/internal/cc0_error"
	"c0_compiler/internal/instruction"
	"fmt"
	"sort"
)

var sortedFunctions = &[]instruction.Symbol{}
var lines = &[]string{}

func appendLine(format string, params ...interface{}) {
	*lines = append(*lines, fmt.Sprintf(format, params...))
}

func appendEmptyLine() {
	*lines = append(*lines, "\n")
}

func printLine(line instruction.Line) {
	appendLine("%s ", line.I.Representation)
	for _, operand := range *line.Operands {
		appendLine("%d ", operand)
	}
	appendEmptyLine()
}

func assembleConstants(st *instruction.SymbolTable) {
	appendLine(".constants:\n")
	for index, sb := range *sortedFunctions {
		appendLine("%d S \"%s\"\n", index, sb.Name)
	}
	// TODO: add support for int and double constants
	appendEmptyLine()
}

func assembleFunctions() {
	appendLine(".functions:\n")
	for index, sb := range *sortedFunctions {
		appendLine("%d %d %d 1\t# %s\n", index, index, len(*sb.FnInfo.Parameters), sb.Name)
	}
	appendEmptyLine()
}

type By func(p1, p2 *instruction.Symbol) bool

type functionSorter struct {
	functions []instruction.Symbol
	by        By
}

func (f functionSorter) Len() int {
	return len(f.functions)
}

func (f functionSorter) Less(i, j int) bool {
	return f.by(&f.functions[i], &f.functions[j])
}

func (f functionSorter) Swap(i, j int) {
	f.functions[i], f.functions[j] = f.functions[j], f.functions[i]
}

func (by By) Sort(functions []instruction.Symbol) {
	ps := &functionSorter{
		functions: functions,
		by:        by,
	}
	sort.Sort(ps)
}

func sortFunctions(table *instruction.SymbolTable) {
	for _, sb := range table.Symbols {
		if !sb.IsCallable {
			continue
		}
		*sortedFunctions = append(*sortedFunctions, sb)
	}
	By(func(p1, p2 *instruction.Symbol) bool { return p1.Address < p2.Address }).Sort(*sortedFunctions)
}

func Run(globalSymbolTable *instruction.SymbolTable) *[]string {
	count := 0

	// Check if `main` function exists first
	if _, ok := globalSymbolTable.Symbols["main"]; !ok {
		cc0_error.Of(cc0_error.NoMain).Die(cc0_error.Assembler)
	}

	sortFunctions(globalSymbolTable)
	assembleConstants(globalSymbolTable)
	assembleFunctions()

	appendLine(".start:\n")
	for _, i := range *globalSymbolTable.RelatedFunction.GetLines() {
		printLine(i)
	}

	for name, sb := range globalSymbolTable.Symbols {
		if !sb.IsCallable {
			continue
		}
		appendLine("\n.F%d:\t# %s\n", count, name)
		for _, i := range *sb.FnInfo.GetLines() {
			printLine(i)
		}
		count++
	}

	return lines
}