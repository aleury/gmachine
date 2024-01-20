package token_test

import (
	"gmachine/token"
	"testing"
)

func TestLookupIdent(t *testing.T) {
	t.Parallel()
	tests := []struct {
		given string
		want  token.TokenType
	}{
		{"CONS", token.CONSTANT_DEFINITION},
		{"VARB", token.VARIABLE_DEFINITION},
		{"HALT", token.INSTRUCTION},
		{"NOOP", token.INSTRUCTION},
		{"MOVE", token.INSTRUCTION},
		{"OUTA", token.INSTRUCTION},
		{"INCA", token.INSTRUCTION},
		{"INCX", token.INSTRUCTION},
		{"INCY", token.INSTRUCTION},
		{"DECA", token.INSTRUCTION},
		{"DECX", token.INSTRUCTION},
		{"DECY", token.INSTRUCTION},
		{"ADDA", token.INSTRUCTION},
		{"MULA", token.INSTRUCTION},
		{"SETA", token.INSTRUCTION},
		{"SETX", token.INSTRUCTION},
		{"SETY", token.INSTRUCTION},
		{"PSHA", token.INSTRUCTION},
		{"POPA", token.INSTRUCTION},
		{"JUMP", token.INSTRUCTION},
		{"JXNZ", token.INSTRUCTION},
		{"A", token.REGISTER},
		{"X", token.REGISTER},
		{"Y", token.REGISTER},
		{"test", token.IDENT},
		{".test", token.IDENT},
	}
	for i, tt := range tests {
		got := token.LookupIdent(tt.given)
		if tt.want != got {
			t.Fatalf("tests[%d] - given=%q, wanted=%q, got=%q", i, tt.given, tt.want, got)
		}
	}
}
