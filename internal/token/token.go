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
	DoubleLiteral
	Const
	Void
	Char // the type
	Int
	Double
	CharLiteral
	StringLiteral
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
	case Int, Char, Double, Void:
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

func (t *Token) IsAnUnaryOperator() bool {
	return t.Kind == PlusSign || t.Kind == MinusSign
}

func (t *Token) IsAMultiplicativeOperator() bool {
	return t.Kind == MultiplicationSign || t.Kind == DivisionSign
}

func (t *Token) IsAnAdditiveOperator() bool {
	return t.Kind == PlusSign || t.Kind == MinusSign
}
