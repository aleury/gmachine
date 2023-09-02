package lexer_test

import (
	"gmachine/lexer"
	"gmachine/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
; this is a comment
JUMP 2
JUMP start

; this is another comment

.test
SETA 'a'
OUTA
HALT 

; this is yet another comment 

.start
NOOP
SETA 42
INCA
DECA
PSHA
POPA
MOVA X
OUTA
HALT
`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.OPCODE, "JUMP"},
		{token.INT, "2"},
		{token.OPCODE, "JUMP"},
		{token.IDENT, "start"},
		{token.IDENT, ".test"},
		{token.OPCODE, "SETA"},
		{token.CHAR, "'a'"},
		{token.OPCODE, "OUTA"},
		{token.OPCODE, "HALT"},
		{token.IDENT, ".start"},
		{token.OPCODE, "NOOP"},
		{token.OPCODE, "SETA"},
		{token.INT, "42"},
		{token.OPCODE, "INCA"},
		{token.OPCODE, "DECA"},
		{token.OPCODE, "PSHA"},
		{token.OPCODE, "POPA"},
		{token.OPCODE, "MOVA"},
		{token.REGISTER, "X"},
		{token.OPCODE, "OUTA"},
		{token.OPCODE, "HALT"},
		{token.EOF, ""},
	}

	l := lexer.New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. wanted=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. wanted=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
