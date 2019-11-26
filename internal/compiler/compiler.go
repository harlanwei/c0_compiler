package compiler

import (
	"bufio"
	"c0_compiler/internal/analyzer"
	"c0_compiler/internal/cc0_error"
	"c0_compiler/internal/parser"
	"os"
)

// Setup parser and analyzer and try to compile the given file.
func Run(in, out string, shouldCompileToBinary, isDebugging bool) {
	reader, err := os.Open(in)
	if err != nil {
		cc0_error.PrintfToStdErr("Can't open specified source file: %s\n", in)
		cc0_error.ThrowAndExit(cc0_error.Source)
	}
	scanner := bufio.NewScanner(reader)
	p := parser.Parse(scanner)

	var writer *bufio.Writer
	if isDebugging {
		writer = bufio.NewWriter(os.Stdout)
	} else {
		fwriter, err := os.Create(out)
		if err != nil {
			cc0_error.PrintfToStdErr("Can't open specified output file: %s\n", out)
			cc0_error.ThrowAndExit(cc0_error.Source)
		}
		writer = bufio.NewWriter(fwriter)
	}
	analyzer.Run(p, writer, shouldCompileToBinary)
}
