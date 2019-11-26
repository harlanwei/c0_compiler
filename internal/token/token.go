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
	Const              = 40
	Void               = 41
	Int                = 42
	Char               = 43
	Double             = 44
	Struct             = 45
	If                 = 46
	Else               = 47
	Switch             = 48
	Case               = 49
	Default            = 50
	While              = 51
	For                = 52
	Do                 = 53
	Return             = 54
	Break              = 55
	Continue           = 56
	Print              = 57
	Scan               = 58
	Identifier         = 1000
)

type any = interface{}

type Token struct {
	Kind   int
	Value  any
	Line   int
	Column int
}
