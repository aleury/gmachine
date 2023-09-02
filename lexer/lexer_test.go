package lexer_test

import (
	"errors"
	"gmachine/lexer"
	"gmachine/token"
	"testing"
)

func TestNextToken_ReturnsIllegalTokenForUnknownToken(t *testing.T) {
	t.Parallel()
	l := lexer.New("~")
	wantLiteral := "~"
	var wantType token.TokenType = token.ILLEGAL
	got, err := l.NextToken()
	if err != nil {
		t.Fatal("didn't expect an error:", err)
	}
	if wantLiteral != got.Literal {
		t.Errorf("token literal wrong - want=%q, got=%q", wantLiteral, got.Literal)
	}
	if wantType != got.Type {
		t.Errorf("token type wrong - want=%q, got=%q", wantType, got.Type)
	}
}

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
	input := `; this is a comment
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
HALT`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
	}{
		{token.OPCODE, "JUMP", 2},
		{token.INT, "2", 2},
		{token.OPCODE, "JUMP", 3},
		{token.IDENT, "start", 3},
		{token.IDENT, ".test", 7},
		{token.OPCODE, "SETA", 8},
		{token.CHAR, "'a'", 8},
		{token.OPCODE, "OUTA", 9},
		{token.OPCODE, "HALT", 10},
		{token.IDENT, ".start", 14},
		{token.OPCODE, "NOOP", 15},
		{token.OPCODE, "SETA", 16},
		{token.INT, "42", 16},
		{token.OPCODE, "INCA", 17},
		{token.OPCODE, "DECA", 18},
		{token.OPCODE, "PSHA", 19},
		{token.OPCODE, "POPA", 20},
		{token.OPCODE, "MOVA", 21},
		{token.REGISTER, "X", 21},
		{token.OPCODE, "OUTA", 22},
		{token.OPCODE, "HALT", 23},
		{token.EOF, "", 23},
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
		if tok.Line != tt.expectedLine {
			t.Fatalf("tests[%d] - line number wrong. wanted=%d, got=%d", i, tt.expectedLine, tok.Line)
		}
	}
}
