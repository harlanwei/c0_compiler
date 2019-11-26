package token

const (
	NotParsed          = -1
	PlusSign           = 1
	MinusSign          = 2
	MultiplicationSign = 3
	DivisionSign       = 4
	LessThan           = 5
	LessThanOrEqual    = 6
	EqualTo            = 7
	GreaterThanOrEqual = 8
	GreaterThan        = 9
	NotEqualTo         = 10
	AssignmentSign     = 20
	LeftBracket        = 21
	RightBracket       = 22
	LeftParenthesis    = 23
	RightParenthesis   = 24
	Semicolon          = 25
	IntegerLiteral     = 30
)

type any = interface{}

type Token struct {
	Kind   int
	Value  any
	Line   int
	Column int
}
