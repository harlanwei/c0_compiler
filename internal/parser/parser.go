package parser

import (
	"bufio"
	"c0_compiler/internal/token"
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

func parseInteger(el string) *Token {
	matched, _ := regexp.MatchString("[0-9]+", el)
	if matched {
		value, _ := strconv.Atoi(el)
		return &Token{Kind: token.INTEGER_LITERAL, Value: value}
	}
	return nil
}

func parse(buffer []string) (res []Token) {
	for _, el := range buffer {
		if unicode.IsNumber(rune(el[0])) {
			if parsedToken := parseInteger(el); parsedToken != nil {
				res = append(res, *parsedToken)
			}
		}
	}
	return
}

func CreateInstance(scanner *bufio.Scanner) (parser *Parser) {
	buffer := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Fields(line)
		buffer = append(buffer, words...)
	}
	parser = &Parser{parse(buffer), 0}

	return
}
