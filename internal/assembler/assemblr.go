package assembler

import (
	"bytes"
	"c0_compiler/internal/cc0_error"
	"c0_compiler/internal/instruction"
	"c0_compiler/internal/token"
	"encoding/binary"
	"fmt"
	"sort"
)

var sortedFunctions = &[]instruction.Symbol{}
var lines = &[]string{}
var addressOffset int // for constants

func appendLine(format string, params ...interface{}) {
	*lines = append(*lines, fmt.Sprintf(format, params...))
}

func appendEmptyLine() {
	*lines = append(*lines, "\n")
}

func printLine(line instruction.Line) {
	if line.I.Code == instruction.Loadc && (*line.Operands)[0] <= 0 {
		(*line.Operands)[0] = addressOffset - (*line.Operands)[0]
	}
	str := line.I.Representation
	for _, operand := range *line.Operands {
		str += fmt.Sprintf(" %d", operand)
	}
	appendLine("%s\n", str)
}

func float64ToByte(f float64) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, f)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}

func assembleConstants(st *instruction.SymbolTable) {
	appendLine(".constants:\n")
	for index, sb := range *sortedFunctions {
		appendLine("%d S \"%s\"\n", index, sb.Name)
	}
	addressOffset = len(*sortedFunctions)
	for _, c := range *st.Constants {
		address := addressOffset - c.Address
		switch c.Kind {
		case instruction.ConstantKindInt:
			// not used
		case instruction.ConstantKindDouble:
			str := "0x"
			for _, b := range float64ToByte(c.Value.(float64)) {
				str += fmt.Sprintf("%02x", b)
			}
			appendLine("%d D %s\n", address, str)
		case instruction.ConstantKindString:
			appendLine("%d S \"%s\"\n", address, c.Value)
		}
	}
	appendEmptyLine()
}

func assembleFunctions() {
	appendLine(".functions:\n")
	for index, sb := range *sortedFunctions {
		requiredSize := 0
		for _, parameterName := range *sb.FnInfo.Parameters {
			switch sb.FnInfo.RelatedSymbolTable.GetSymbolNamed(parameterName).Kind {
			case token.Double:
				requiredSize += 2
			default:
				requiredSize += 1
			}
		}
		appendLine("%d %d %d 1\t# %s\n", index, index, requiredSize, sb.Name)
	}
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
		*sortedFunctions = append(*sortedFunctions, *sb)
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

	appendLine(".start:\n")
	for _, i := range *globalSymbolTable.RelatedFunction.GetLines() {
		printLine(i)
	}
	appendEmptyLine()

	assembleFunctions()

	for name, sb := range *sortedFunctions {
		appendLine("\n.F%d:\t# %s\n", count, name)
		for _, i := range *sb.FnInfo.GetLines() {
			printLine(i)
		}
		count++
	}

	return lines
}
