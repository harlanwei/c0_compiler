package main

import (
	"c0_compiler/internal/compiler"
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
	isDebugging := flag.Bool("debug", false, "Run the program in debugging mode")
	destination := flag.String("o", "out", "输出到指定的文件 file")
	flag.Parse()

	// cc0 [-h]
	if *shouldShowUsage {
		displayUsage(false)
	}

	// cc0 [options] input [-o file]
	remainingArgs := flag.Args()
	var source string
	if len(remainingArgs) != 1 {
		displayUsage(true)
	} else {
		source = remainingArgs[0]
	}

	if *shouldOutputBinary {
		compiler.Run(source, *destination, true, false)
	} else if *shouldOutputText {
		compiler.Run(source, *destination, false, *isDebugging)
	} else {
		displayUsage(true)
	}
}
