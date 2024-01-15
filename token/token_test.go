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
		{"HALT", token.OPCODE},
		{"NOOP", token.OPCODE},
		{"OUTA", token.OPCODE},
		{"INCA", token.OPCODE},
		{"DECA", token.OPCODE},
		{"DECX", token.OPCODE},
		{"DECY", token.OPCODE},
		{"ADDA", token.OPCODE},
		{"MULA", token.OPCODE},
		{"MOVA", token.OPCODE},
		{"SETA", token.OPCODE},
		{"SETX", token.OPCODE},
		{"SETY", token.OPCODE},
		{"PSHA", token.OPCODE},
		{"POPA", token.OPCODE},
		{"JUMP", token.OPCODE},
		{"X", token.REGISTER},
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
