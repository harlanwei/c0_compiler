package main

import (
	"bufio"
	"c0_compiler/internal/analyzer"
	"c0_compiler/internal/assembler"
	"c0_compiler/internal/cc0_error"
	"c0_compiler/internal/compiler"
	"c0_compiler/internal/parser"
	"flag"
	"fmt"
	"os"
)

const usage = `Usage:
cc0 [options] input [-o file]
cc0 [-h]

Options:
	-s        将输入的 c0 源代码翻译为文本汇编文件
	-c        将输入的 c0 源代码翻译为二进制目标文件
	-h        显示关于编译器使用的帮助
	-o file   输出到指定的文件 file，默认为 out`

func displayUsage(toStdErr bool) {
	if toStdErr {
		_, _ = fmt.Fprintf(os.Stderr, usage)
	} else {
		fmt.Println(usage)
	}
	os.Exit(0)
}

func main() {
	shouldShowUsage := flag.Bool("h", false, "显示关于编译器使用的帮助")
	shouldOutputText := flag.Bool("s", false, "将输入的 c0 源代码翻译为文本汇编文件")
	shouldOutputBinary := flag.Bool("c", false, "将输入的 c0 源代码翻译为二进制目标文件")
	destination := flag.String("o", "out", "输出到指定的文件 file")

	// cc0 [options] input [-o file]
	flag.Parse()
	remainingArgs := flag.Args()
	var source string
	if len(remainingArgs) != 1 {
		displayUsage(true)
	} else {
		source = remainingArgs[0]
	}

	// cc0 [-h]
	if *shouldShowUsage {
		displayUsage(false)
	}

	reader, err := os.Open(source)
	if err != nil {
		cc0_error.PrintfToStdErr("Can't open specified source file: %s\n", source)
		cc0_error.ThrowAndExit(cc0_error.Source)
	}
	defer func() {
		if err := reader.Close(); err != nil {
			panic(err)
		}
	}()
	scanner := bufio.NewScanner(reader)
	p := parser.Parse(scanner)
	globalSymbolTable := analyzer.Run(p)
	lines := assembler.Run(globalSymbolTable)

	var outfile *os.File
	outfile, err = os.Create(*destination)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := outfile.Close(); err != nil {
			panic(err)
		}
	}()
	w := bufio.NewWriter(outfile)

	if *shouldOutputText {
		for _, line := range *lines {
			if _, err := w.WriteString(line); err != nil {
				panic(err)
			}
		}
	} else if *shouldOutputBinary {
		compiler.Run(lines, w)
	} else {
		displayUsage(true)
	}
	if err := w.Flush(); err != nil {
		panic(err)
	}
}
