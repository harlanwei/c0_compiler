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

func (t *Token) IsATypeSpecifier() bool {
	k := t.Kind
	switch k {
	case Int, Char, Void:
		return true
	}
	return false
}

func (t *Token) IsARelationalOperator() bool {
	k := t.Kind
	switch k {
	case LessThan, LessThanOrEqual, GreaterThan, GreaterThanOrEqual, NotEqualTo, EqualTo:
		return true
	}
	return false
}
