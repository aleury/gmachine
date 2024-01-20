package lexer_test

import (
	"gmachine/lexer"
	"gmachine/token"
	"strings"
	"testing"
)

func TestNextToken_ReturnsIllegalTokenForUnknownToken(t *testing.T) {
	t.Parallel()
	l := newLexerFromString("~")
	wantLiteral := "~"
	var wantType token.TokenType = token.ILLEGAL
	got := l.NextToken()
	if wantLiteral != got.Literal {
		t.Errorf("token literal wrong - want=%q, got=%q", wantLiteral, got.Literal)
	}
	if wantType != got.Type {
		t.Errorf("token type wrong - want=%q, got=%q", wantType, got.Type)
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
MOVE A -> X
MOVE A -> Y
MOVE *A -> X
OUTA
HALT
"test"
""`
	tests := []struct {
		Type    token.TokenType
		Literal string
		Line    int
	}{
		{token.INSTRUCTION, "JUMP", 2},
		{token.INT, "2", 2},
		{token.INSTRUCTION, "JUMP", 3},
		{token.IDENT, "start", 3},
		{token.LABEL_DEFINITION, ".test", 7},
		{token.INSTRUCTION, "SETA", 8},
		{token.CHAR, "'a'", 8},
		{token.INSTRUCTION, "OUTA", 9},
		{token.INSTRUCTION, "HALT", 10},
		{token.LABEL_DEFINITION, ".start", 14},
		{token.INSTRUCTION, "NOOP", 15},
		{token.INSTRUCTION, "SETA", 16},
		{token.INT, "42", 16},
		{token.INSTRUCTION, "INCA", 17},
		{token.INSTRUCTION, "DECA", 18},
		{token.INSTRUCTION, "PSHA", 19},
		{token.INSTRUCTION, "POPA", 20},
		{token.INSTRUCTION, "MOVE", 21},
		{token.REGISTER, "A", 21},
		{token.ARROW, "->", 21},
		{token.REGISTER, "X", 21},
		{token.INSTRUCTION, "MOVE", 22},
		{token.REGISTER, "A", 22},
		{token.ARROW, "->", 22},
		{token.REGISTER, "Y", 22},
		{token.INSTRUCTION, "MOVE", 23},
		{token.ASTERISK, "*", 23},
		{token.REGISTER, "A", 23},
		{token.ARROW, "->", 23},
		{token.REGISTER, "X", 23},
		{token.INSTRUCTION, "OUTA", 24},
		{token.INSTRUCTION, "HALT", 25},
		{token.STRING, "test", 26},
		{token.STRING, "", 27},
		{token.EOF, "", 27},
	}

	l := newLexerFromString(input)
	lines := strings.Split(input, "\n")
	for _, want := range tests {
		line := lines[want.Line-1]
		t.Run(line, func(t *testing.T) {
			got := l.NextToken()
			if got.Type != want.Type || got.Literal != want.Literal {
				t.Fatalf("wanted=%q [%s], got=%q [%s]", want.Literal, want.Type, got.Literal, got.Type)
			}
			if got.Line != want.Line {
				t.Fatalf("line number wrong. wanted=%d, got=%d", want.Line, got.Line)
			}
		})
	}
}

func newLexerFromString(input string) *lexer.Lexer {
	l, err := lexer.New(strings.NewReader(input))
	if err != nil {
		panic(err)
	}
	return l
}
