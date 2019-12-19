package parser

import (
	"bufio"
	"c0_compiler/internal/cc0_error"
	"c0_compiler/internal/token"
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var isInACommentBlock = false

type Token = token.Token

type Parser struct {
	buffer [] Token
	pos    int
}

func (p *Parser) HasNextToken() bool {
	return p.pos < len(p.buffer)
}

func (p *Parser) NextToken() (res *Token) {
	res = &p.buffer[p.pos]
	p.pos++
	return
}

func (p *Parser) CurrentHead() int {
	return p.pos
}

func (p *Parser) ResetHeadTo(pos int) *Token {
	p.pos = pos
	if p.HasNextToken() {
		return &p.buffer[p.pos]
	} else {
		return &Token{
			Kind:   0,
			Value:  nil,
			Line:   0,
			Column: 0,
		}
	}
}

type parserError struct {
	message string
	fatal   bool
}

func reportPosition(token *Token) {
	cc0_error.ReportLineAndColumn(token.Line, token.Column)
}

var decMatcher, _ = regexp.Compile("^(0|([1-9][0-9]*))$")
var hexMatcher, _ = regexp.Compile("^[0-9a-fA-F]+$")
var doubleMatcher, _ = regexp.Compile("^[0-9]+\\.[0-9]+$")

// Accepts a regex matcher used to pre-check the format of the given word.
func parseIntegerValue(matcher *regexp.Regexp, base int, word string) (int64, *parserError) {
	if matcher.MatchString(word) {
		value, err := strconv.ParseInt(word, base, 64)
		if err != nil {
			if _, ok := err.(*strconv.NumError); ok {
				return 0, &parserError{"value out of range; set to 0", false}
			}
			return 0, &parserError{err.Error(), true}
		}
		return value, nil
	}
	return 0, &parserError{"illegal integer literal", true}
}

func parseFloat(str string) (float64, error) {
	val, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return val, nil
	}

	// Some number is specifed in scientific notation
	pos := strings.IndexAny(str, "eE")
	if pos < 0 {
		return strconv.ParseFloat(str, 64)
	}

	var baseVal float64
	var expVal int64

	baseStr := str[0:pos]
	baseVal, err = strconv.ParseFloat(baseStr, 64)
	if err != nil {
		return 0, err
	}

	expStr := str[(pos + 1):]
	expVal, err = strconv.ParseInt(expStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return baseVal * math.Pow10(int(expVal)), nil
}

func parseDoubleValue(matcher *regexp.Regexp, word string) (float64, *parserError) {
	value, err := parseFloat(word)
	if err != nil {
		return 0, &parserError{"Failed to parse double value", true}
	}
	return value, nil
}

func decIntegerParser(word string) (int64, *parserError) {
	return parseIntegerValue(decMatcher, 10, word)
}

func hexIntegerParser(word string) (int64, *parserError) {
	// `strconv.ParseInt(.., 16, ..)` does not accept integers in the format of `0x1234`.
	// Must remove the prefixes before handing the tokens over.
	return parseIntegerValue(hexMatcher, 16, word[2:])
}

func doubleParser(word string) (float64, *parserError) {
	return parseDoubleValue(doubleMatcher, word)
}

func parseIntegerLiteral(currentToken *Token) {
	var integerParser func(string) (int64, *parserError)
	word := currentToken.Value.(string)

	if len(word) > 1 && (word[0:2] == "0x" || word[0:2] == "0X") {
		integerParser = hexIntegerParser
	} else {
		integerParser = decIntegerParser
	}

	parsedValue, err := integerParser(word)
	currentToken.Kind = token.IntegerLiteral
	currentToken.Value = parsedValue
	if err != nil {
		reportPosition(currentToken)
		cc0_error.PrintlnToStdErr(err.message)
		if err.fatal {
			cc0_error.ThrowAndExit(cc0_error.Parser)
		} else {
			cc0_error.ThrowButStayAlive(cc0_error.Parser)
		}
	}
}

func parseDoubleLiteral(currentToken *Token) {
	parsedValue, err := doubleParser(currentToken.Value.(string))
	if err != nil {
		reportPosition(currentToken)
		cc0_error.PrintlnToStdErr(err.message)
		cc0_error.ThrowAndExit(cc0_error.Parser)
	}
	currentToken.Kind = token.DoubleLiteral
	currentToken.Value = parsedValue
}

func parseOperator(currentToken *Token) {
	kind := &currentToken.Kind

	switch word := currentToken.Value.(string); word {
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
	case ",":
		*kind = token.Comma
	case ";":
		*kind = token.Semicolon
	default:
		reportPosition(currentToken)
		cc0_error.PrintfToStdErr("Unrecognized character '%s'\n", word)
		cc0_error.ThrowAndExit(cc0_error.Parser)
	}
}

func parseKeywords(currentToken *Token) {
	kind := &currentToken.Kind

	switch word := currentToken.Value; word {
	case "const":
		*kind = token.Const
	case "void":
		*kind = token.Void
	case "int":
		*kind = token.Int
	case "char":
		*kind = token.Char
	case "double":
		*kind = token.Double
	case "struct":
		*kind = token.Struct
	case "if":
		*kind = token.If
	case "else":
		*kind = token.Else
	case "switch":
		*kind = token.Switch
	case "case":
		*kind = token.Case
	case "default":
		*kind = token.Default
	case "while":
		*kind = token.While
	case "for":
		*kind = token.For
	case "do":
		*kind = token.Do
	case "return":
		*kind = token.Return
	case "break":
		*kind = token.Break
	case "continue":
		*kind = token.Continue
	case "print":
		*kind = token.Print
	case "scan":
		*kind = token.Scan
	default:
		*kind = token.Identifier
	}
}

func parseAllTheTokensIn(buffer []Token) {
	for ind := range buffer {
		currentToken := &buffer[ind]

		if currentToken.Kind == token.StringLiteral || currentToken.Kind == token.CharLiteral {
			continue
		}
		word := currentToken.Value.(string)
		if strings.ContainsRune(word, '.') ||
			(unicode.IsNumber(rune(word[0])) &&
				len(word) > 2 &&
				strings.ToLower(word[0:2]) != "0x" &&
				strings.ContainsRune(strings.ToLower(word), 'e')) {
			parseDoubleLiteral(currentToken)
		} else if unicode.IsNumber(rune(word[0])) {
			parseIntegerLiteral(currentToken)
		} else if unicode.IsLetter(rune(word[0])) {
			parseKeywords(currentToken)
		} else {
			parseOperator(currentToken)
		}
	}
}

func isDigitOrLetter(character rune) bool {
	return unicode.IsNumber(character) || unicode.IsLetter(character) || character == '.'
}

func isAnOperatorWithTwoCharacters(operator string) bool {
	switch operator {
	case "<=", ">=", "==", "!=", "//", "/*", "*/":
		return true
	}
	return false
}

func parseCharSequence(line string, start int) (result rune, end int) {
	lineLength := len(line)
	if line[start] == '\\' {
		// escaped char
		if start+1 >= lineLength {
			return -1, 0
		}
		if line[start+1] == 'x' {
			// represented by hex
			if start+3 >= lineLength {
				return -1, 0
			}
			resultAsInt64, _ := strconv.ParseInt(line[start+2:start+4], 16, 32)
			result = rune(resultAsInt64)
			end = start + 4
		} else if line[start+1] == 'n' {
			result = '\n'
			end = start + 2
		} else if line[start+1] == 'r' {
			result = '\r'
			end = start + 2
		} else if line[start+1] == 't' {
			result = '\t'
			end = start + 2
		} else {
			result = rune(line[start+1])
			end = start + 2
		}
	} else {
		result = rune(line[start])
		end = start + 1
	}
	return
}

// `line[start]` should be a single quote.
func parseCharLiteral(line string, start int) (result rune, end int) {
	lineLength := len(line)
	if start+1 >= lineLength {
		return -1, 0
	}
	result, end = parseCharSequence(line, start+1)
	if end >= len(line) || line[end] != '\'' {
		return -1, 0
	}
	end++
	return
}

// `line[start]` should be a double quote.
func parseStringLiteral(line string, start int) ([]rune, int, bool) {
	end, lineLength := start+1, len(line)
	res := []rune{}
	for end < lineLength {
		b, nextEnd := parseCharSequence(line, end)
		if b < 0 {
			return res, -1, false
		}
		if nextEnd-end == 1 && b == '"' {
			return res, nextEnd, true
		}
		res = append(res, b)
		end = nextEnd
	}
	return res, -1, false
}

func reportIllegalCharLiteral(lineCount, columnCount int) {
	cc0_error.ReportLineAndColumn(lineCount, columnCount)
	cc0_error.PrintlnToStdErr("Illegal character literal.")
	cc0_error.ThrowAndExit(cc0_error.Parser)
}

func reportIllegalStringLiteral(lineCount, columnCount int) {
	cc0_error.ReportLineAndColumn(lineCount, columnCount)
	cc0_error.PrintlnToStdErr("Illegal string literal.")
	cc0_error.ThrowAndExit(cc0_error.Parser)
}

func divideTokens(lineCount int, line string, buffer *[]Token) {
	columnCount := 0
	for columnCount < len(line) {
		character := rune(line[columnCount])

		if unicode.IsSpace(character) {
			columnCount++
			continue
		}
		end := columnCount
		if isDigitOrLetter(character) {
			previousCharacter := '0'
			for end < len(line) && (isDigitOrLetter(rune(line[end])) ||
				((previousCharacter == 'e' || previousCharacter == 'E') && line[end] == '-')) {
				previousCharacter = rune(line[end])
				end++
			}
		} else {
			if columnCount+1 < len(line) && isAnOperatorWithTwoCharacters(line[columnCount:columnCount+2]) {
				end = columnCount + 2
			} else {
				end = columnCount + 1
			}
		}

		currentTokenString := line[columnCount:end]
		columnCount = end

		if currentTokenString == "//" {
			// Discard rest of this line
			return
		}

		if currentTokenString == "/*" {
			isInACommentBlock = true
			continue
		} else if currentTokenString == "*/" {
			if !isInACommentBlock {
				cc0_error.ReportLineAndColumn(lineCount, columnCount)
				cc0_error.PrintlnToStdErr("encountered an illegal comment block.")
				cc0_error.ThrowAndExit(cc0_error.Parser)
			}
			isInACommentBlock = false
			continue
		}

		if isInACommentBlock {
			continue
		}
		if currentTokenString == "'" {
			if columnCount+2 >= len(line) {
				reportIllegalCharLiteral(lineCount, columnCount)
			}
			parsed, end := parseCharLiteral(line, columnCount-1)
			if parsed < 0 {
				reportIllegalCharLiteral(lineCount, columnCount)
			}
			columnCount = end
			*buffer = append(*buffer, Token{
				Kind:   token.CharLiteral,
				Value:  parsed,
				Line:   lineCount,
				Column: columnCount,
			})
		} else if currentTokenString == "\"" {
			parsed, end, ok := parseStringLiteral(line, columnCount-1)
			if !ok {
				reportIllegalStringLiteral(lineCount, columnCount)
			}
			*buffer = append(*buffer, Token{
				Kind:   token.StringLiteral,
				Value:  string(parsed),
				Line:   lineCount,
				Column: columnCount,
			})
			columnCount = end
		} else {
			*buffer = append(*buffer, Token{
				Kind:   token.NotParsed,
				Value:  currentTokenString,
				Line:   lineCount,
				Column: columnCount + 1,
			})
		}
	}
}

func Parse(scanner *bufio.Scanner) (parser *Parser) {
	buffer := make([]Token, 0)
	lineCount := 0

	for scanner.Scan() {
		lineCount++
		line := scanner.Text()
		divideTokens(lineCount, line, &buffer)
	}
	parseAllTheTokensIn(buffer)
	parser = &Parser{buffer, 0}

	return
}
