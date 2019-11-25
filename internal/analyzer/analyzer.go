package analyzer

import (
	"bufio"
	"c0_compiler/internal/parser"
	"fmt"
)

type Parser = parser.Parser

func Run(parser *Parser, writer *bufio.Writer, shouldCompileToBinary bool) {
	for parser.HasNextToken() {
		_, _ = fmt.Fprintln(writer, parser.NextToken())
	}
	_ = writer.Flush()
}
