package parser

import (
	"bufio"
	"c0_compiler/internal/cc0_error"
	"c0_compiler/internal/token"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type Token = token.Token

type Parser struct {
	buffer [] Token
	pos    int
}

func (p *Parser) HasNextToken() bool {
	return p.pos < len(p.buffer)
}

func (p *Parser) NextToken() (res Token) {
	res = p.buffer[p.pos]
	p.pos++
	return
}

func (p *Parser) UnreadToken() bool {
	if p.pos == 0 {
		return false
	}
	p.pos--
	return true
}

func reportPosition(token *Token) {
	_, _ = fmt.Fprintf(os.Stderr, "At line %d, column %d: ", token.Line, token.Column)
}

type regexEngine struct {
	decInteger *regexp.Regexp
	hexInteger *regexp.Regexp
}

func createDefaultEngine() *regexEngine {
	decInteger, _ := regexp.Compile("^[1-9][0-9]*$")
	hexInteger, _ := regexp.Compile("^0[xX][0-9a-fA-F]+$")
	return &regexEngine{
		decInteger: decInteger,
		hexInteger: hexInteger,
	}
}

var globalEngine *regexEngine

func parseDecimalInteger(word string) int64 {
	if globalEngine.decInteger.MatchString(word) {
		value, err := strconv.ParseInt(word, 10, 64)
		if err != nil {
			return -1
		}
		return value
	}
	return -1
}

func parseHexInteger(word string) int64 {
	if globalEngine.hexInteger.MatchString(word) {
		value, err := strconv.ParseInt(word, 16, 64)
		if err != nil {
			return -1
		}
		return value
	}
	return -1
}

func parseIntegerLiteral(currentToken *Token) {
	word := currentToken.Value.(string)

	if word[0] == '0' {
		if parsedValue := parseHexInteger(word); parsedValue >= 0 {
			currentToken.Kind = token.IntegerLiteral
			currentToken.Value = parsedValue
		} else {
			reportPosition(currentToken)
			cc0_error.PrintlnToStdErr("Illegal hexadecimal integer literal.")
			cc0_error.ThrowAndExit(cc0_error.Parser)
		}
	} else if parsedValue := parseDecimalInteger(word); parsedValue >= 0 {
		currentToken.Kind = token.IntegerLiteral
		currentToken.Value = parsedValue
	} else {
		reportPosition(currentToken)
		cc0_error.PrintlnToStdErr("Illegal integer literal.")
		cc0_error.ThrowAndExit(cc0_error.Parser)
	}
}

func parseOperator(currentToken *Token) {
	kind := &currentToken.Kind
	word := currentToken.Value.(string)

	switch word {
	case "+":
		*kind = token.PlusSign
	case "-":
		*kind = token.MinusSign
	case "*":
		*kind = token.MultiplicationSign
	case "/":
		*kind = token.DivisionSign
	case "=":
		*kind = token.AssignmentSign
	case "(":
		*kind = token.LeftParenthesis
	case ")":
		*kind = token.RightParenthesis
	case "{":
		*kind = token.LeftBracket
	case "}":
		*kind = token.RightBracket
	case ">":
		*kind = token.GreaterThan
	case ">=":
		*kind = token.GreaterThanOrEqual
	case "==":
		*kind = token.EqualTo
	case "<=":
		*kind = token.LessThanOrEqual
	case "<":
		*kind = token.LessThan
	case "!=":
		*kind = token.NotEqualTo
	case ";":
		*kind = token.Semicolon
	}
}

func parse(buffer []Token) {
	for ind, _ := range buffer {
		currentToken := &buffer[ind]
		word := currentToken.Value.(string)

		if unicode.IsNumber(rune(word[0])) {
			parseIntegerLiteral(currentToken)
		} else if unicode.IsLetter(rune(word[0])) {
			// TODO: Identifier or keywords
		} else {
			parseOperator(currentToken)
		}
	}
}

func CreateInstance(scanner *bufio.Scanner) (parser *Parser) {
	globalEngine = createDefaultEngine()

	buffer := make([]Token, 0)
	lineCount := 0

	// TODO: fix the bugs related with dividers
	for scanner.Scan() {
		lineCount++
		columnCount := 0
		tokens := make([]Token, 0)

		line := scanner.Text()
		for _, word := range strings.Fields(line) {
			columnCount += strings.Index(line[columnCount:], word)
			tokens = append(tokens, Token{
				Kind:   token.NotParsed,
				Value:  word,
				Line:   lineCount,
				Column: columnCount + 1,
			})
			columnCount += len(word)
		}
		buffer = append(buffer, tokens...)
	}
	parse(buffer)
	parser = &Parser{buffer, 0}

	return
}
