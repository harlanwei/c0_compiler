package analyzer

import (
	"c0_compiler/internal/cc0_error"
	"c0_compiler/internal/token"
)

type c0Type = int

const (
	TypeInt = iota
	TypeVoid
)

func tokenKindToType(t *token.Token) c0Type {
	switch t.Kind {
	case token.Int:
		return TypeInt
	case token.Void:
		return TypeVoid
	default:
		cc0_error.PrintlnToStdErr("Failed to resolve a certain type. This might be a bug in the analyzer.")
		cc0_error.ThrowButStayAlive(cc0_error.Analyzer)
		return TypeVoid
	}
}
