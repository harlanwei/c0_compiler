package compiler

import (
	"bufio"
	"c0_compiler/internal/instruction"
	"regexp"
	"strconv"
	"strings"
)

const lower32BitsMask = 0xffffffff

// Constants
var magics = []byte{0x43, 0x30, 0x3a, 0x29}
var version = []byte{0x0, 0x0, 0x0, 0x1}
var functionMatcher, _ = regexp.Compile("\\.F[0-9]+:")
var constantParser, _ = regexp.Compile("(?s)^([0-9]+) (.) (.+)")

// Global variables
var allLines []string
var paramLengths = []int{}
var currentLine = 0
var currentFn = 0
var w *bufio.Writer

func nextLine() (res *string) {
	if !hasNextLine() {
		return nil
	}
	res = &allLines[currentLine]
	currentLine++
	return
}

func collectSection() (int, int) {
	start, end := currentLine, currentLine
	for hasNextLine() && !nextLineIsEmpty() && !nextLineIsADelimiter() {
		end++
		nextLine()
	}
	return start, end
}

func peekNextLine() (res *string) {
	if !hasNextLine() {
		return nil
	}
	return &allLines[currentLine]
}

func hasNextLine() bool {
	return currentLine < len(allLines)
}

func nextLineIsEmpty() bool {
	return len(strings.TrimSpace(allLines[currentLine])) == 0
}

func nextLineIsADelimiter() bool {
	line := strings.TrimSpace(allLines[currentLine])
	if nextLineIsEmpty() {
		return false
	}
	return line[0] == '.'
}

func writeI32WithWidth(value, width int) {
	padded := []byte{}
	for i := 0; i < width; i++ {
		padded = append(padded, 0)
	}
	for i := width - 1; i >= 0; i-- {
		padded[i] = byte(value & 0xff)
		value >>= 8
	}
	_, _ = w.Write(padded)
}

func writeConstant(line string) {
	matches := constantParser.FindStringSubmatch(line)
	kind, value := matches[2], matches[3]
	value = strings.TrimRight(value, "\n")
	if kind == "I" {
		literal, _ := strconv.Atoi(value)
		writeI32WithWidth(1, 1)
		writeI32WithWidth(literal, 4)
	} else if kind == "D" {
		writeI32WithWidth(2, 1)
		for i := 2; i < len(value); i += 2 {
			substr := value[i : i+2]
			parsed, _ := strconv.ParseInt(substr, 16, 16)
			writeI32WithWidth(int(parsed), 1)
		}
	} else if kind == "S" {
		writeI32WithWidth(0, 1)
		writeI32WithWidth(len(value)-2, 2)
		for i := 1; i < len(value)-1; i++ {
			_ = w.WriteByte(value[i])
		}
	}
}

func compileConstants() {
	start, end := collectSection()
	writeI32WithWidth(end-start, 2)
	for i := start; i < end; i++ {
		line := allLines[i]
		writeConstant(line)
	}
}

func compileInstruction(n int) {
	line := strings.TrimSpace(allLines[n])
	fields := strings.Split(line, " ")
	currentInstruction := instruction.GetCodeFrom(fields[0])
	writeI32WithWidth(currentInstruction.Code, 1)
	for i := 1; i < len(fields); i++ {
		v, _ := strconv.Atoi(fields[i])
		writeI32WithWidth(v, currentInstruction.Operands[i-1])
	}
}

func compileGlobalStartSection() {
	start, end := collectSection()
	writeI32WithWidth(end-start, 2)
	for i := start; i < end; i++ {
		compileInstruction(i)
	}
}

func compileFunctionBriefings() {
	start, end := collectSection()
	writeI32WithWidth(end-start, 2)
	for i := start; i < end; i++ {
		line := allLines[i]
		fields := strings.Split(line, " ")
		nParams, _ := strconv.Atoi(fields[2])
		paramLengths = append(paramLengths, nParams)
	}
}

func compileFunction() {
	start, end := collectSection()
	writeI32WithWidth(currentFn, 2) // nameIndex
	writeI32WithWidth(paramLengths[currentFn], 2)
	writeI32WithWidth(1, 2)         // level
	writeI32WithWidth(end-start, 2) // nInstructions
	for i := start; i < end; i++ {
		compileInstruction(i)
	}
	currentFn++
}

func Run(lines *[]string, destination *bufio.Writer) {
	w = destination
	if _, err := w.Write(magics); err != nil {
		panic(err)
	}
	if _, err := w.Write(version); err != nil {
		panic(err)
	}

	allLines = *lines

	for hasNextLine() {
		var line string
		if nextLineIsEmpty() {
			nextLine()
			continue
		} else {
			line = strings.TrimSpace(*peekNextLine())
		}
		if line[0] != '.' {
			panic("Unparsed line: " + line)
		} else if line == ".constants:" {
			nextLine()
			compileConstants()
		} else if line == ".start:" {
			nextLine()
			compileGlobalStartSection()
		} else if line == ".functions:" {
			nextLine()
			compileFunctionBriefings()
		} else if functionMatcher.MatchString(line) {
			nextLine()
			compileFunction()
		} else {
			panic("Unparsed line: " + line)
		}
	}
}
