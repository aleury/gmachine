package lexer_test

import (
	"errors"
	"gmachine/lexer"
	"gmachine/token"
	"testing"
)

func TestNextToken_ReturnsErrorForInvalidCharacterLiteral(t *testing.T) {
	t.Parallel()
	l := lexer.New("'c")
	_, err := l.NextToken()
	if err == nil {
		t.Fatal("expected an error, but didn't receive one")
	}
	if !errors.Is(err, lexer.ErrInvalidCharacterLiteral) {
		t.Errorf("error wrong, wanted=%q, got=%q", lexer.ErrInvalidCharacterLiteral, err)
	}
}

func TestNextToken_TokenizesValidCode(t *testing.T) {
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
		tok, err := l.NextToken()
		if err != nil {
			t.Fatalf("tests[%d] - didn't expect an error: %q", i, err)
		}
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. wanted=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. wanted=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
