package compiler

import (
	"bufio"
	"c0_compiler/internal/instruction"
	"regexp"
	"strconv"
	"strings"
)

// TODO: fix current buggy implementation with lines

// Constants
var magics = []byte{0x43, 0x30, 0x3a, 0x29}
var version = []byte{0x0, 0x0, 0x0, 0x1}
var functionMatcher, _ = regexp.Compile("^\\.F{0-9}+")
var constantParser, _ = regexp.Compile("^([0-9]+) (.) (.+)")

var allLines []string
var currentLine = 0
var w *bufio.Writer

func nextLine() (res *string) {
	if !hasNextLine() {
		return nil
	}
	res = &allLines[currentLine]
	currentLine++
	return
}

func putBackOneLine() {
	if currentLine > 0 {
		currentLine--
	}
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

func isAnEmptyLine() bool {
	return len(strings.TrimSpace(allLines[currentLine])) == 0
}

func isADelimiterLine() bool {
	line := allLines[currentLine]
	if isAnEmptyLine() {
		return false
	}
	return line[0] == '.'
}

func writeValueWithWidth(value, width int) {
	padded := []byte{}
	for i := 0; i < width; i++ {
		padded = append(padded, 0)
	}
	for i := width - 1; i >= 0; i-- {
		padded[i] = byte(value % 0xff)
		value /= 0xff
	}
	_, _ = w.Write(padded)
}

func writeConstant(line string) {
	matches := constantParser.FindStringSubmatch(line)
	kind, value := matches[2], matches[3]
	if kind == "I" {
		literal, _ := strconv.Atoi(value)
		writeValueWithWidth(literal, 4)
	} else if kind == "D" {
		// TODO
	} else if kind == "S" {
		writeValueWithWidth(0, 1)
		writeValueWithWidth(len(value)-2, 2)
		for i := 1; i < len(value)-1; i++ {
			_ = w.WriteByte(value[i])
		}
	}
}

func writeConstants() {
	start, end := currentLine, currentLine
	for hasNextLine() && !isAnEmptyLine() && !isADelimiterLine() {
		end++
		nextLine()
	}
	writeValueWithWidth(end-start, 2)
	for i := start; i < end; i++ {
		line := allLines[i]
		writeConstant(line)
	}
}

func writeInstruction(n int) {
	line := allLines[n]
	fields := strings.Split(line, " ")
	currentInstruction := instruction.GetCodeFrom(fields[0])
	writeValueWithWidth(currentInstruction.Code, 1)
	for i := 1; i < len(fields); i++ {
		v, _ := strconv.Atoi(fields[i])
		writeValueWithWidth(v, currentInstruction.Operands[i-1])
	}
}

func writeStart() {
	start, end := currentLine, currentLine
	for hasNextLine() && !isAnEmptyLine() && !isADelimiterLine() {
		end++
		nextLine()
	}
	writeValueWithWidth(end-start, 2)
	for i := start; i < end; i++ {
		writeInstruction(i)
	}
}

func writeFunctionBriefings() {
	panic("Implement me!")
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
		if isAnEmptyLine() {
			nextLine()
		} else {
			line = strings.TrimSpace(*peekNextLine())
		}
		if line[0] != '.' {
			// TODO: in real implementation, switch to: panic("Unparsed line: " + line)
			nextLine()
		} else if line == ".constants:" {
			nextLine()
			writeConstants()
		} else if line == ".start:" {
			nextLine()
			writeStart()
		} else if line == ".functions:" {
			nextLine()
			writeFunctionBriefings()
		} else if functionMatcher.MatchString(line) {
			nextLine()
		} else {
			// TODO: in real implementation, switch to: panic("Unparsed line: " + line)
			nextLine()
		}
	}
}
