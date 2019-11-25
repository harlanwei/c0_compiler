package token

const (
	INTEGER_LITERAL = 0
)

type any = interface{}

type Token struct {
	Kind  int
	Value any
}
