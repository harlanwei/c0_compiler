package compiler

import (
	"bufio"
	"c0_compiler/internal/analyzer"
	"c0_compiler/internal/error"
	"c0_compiler/internal/parser"
	"fmt"
	"os"
)

func Run(in, out string, shouldCompileToBinary, isDebugging bool) {
	reader, err := os.Open(in)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Can't open specified source file: %s", in)
		os.Exit(error.CC0_SOURCE_FAILED)
	}
	scanner := bufio.NewScanner(reader)
	p := parser.CreateInstance(scanner)

	var writer *bufio.Writer
	if isDebugging {
		writer = bufio.NewWriter(os.Stdout)
	} else {
		fwriter, err := os.Create(out)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Can't open specified output file: %s", out)
			os.Exit(error.CC0_SOURCE_FAILED)
		}
		writer = bufio.NewWriter(fwriter)
	}
	analyzer.Run(p, writer, shouldCompileToBinary)
}
