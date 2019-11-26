package token

const (
	NotParsed = -1
	PlusSign  = iota
	MinusSign
	MultiplicationSign
	DivisionSign
	LessThan
	LessThanOrEqual
	EqualTo
	GreaterThanOrEqual
	GreaterThan
	NotEqualTo
	AssignmentSign
	LeftBracket
	RightBracket
	LeftParenthesis
	RightParenthesis
	Comma
	Semicolon
	IntegerLiteral
	Const
	Void
	Int
	Char
	Double
	Struct
	If
	Else
	Switch
	Case
	Default
	While
	For
	Do
	Return
	Break
	Continue
	Print
	Scan
	Identifier
)

type any = interface{}

type Token struct {
	Kind   int
	Value  any
	Line   int
	Column int
}
