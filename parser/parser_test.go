package parser_test

import (
	"errors"
	"gmachine/ast"
	"gmachine/lexer"
	"gmachine/parser"
	"gmachine/token"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseProgram_ParsesLabelDefinitions(t *testing.T) {
	t.Parallel()

	input := `.test`
	l := newLexerFromString(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}

	want := []ast.Statement{
		ast.LabelDefinitionStatement{
			Token: token.Token{
				Type:    token.LABEL_DEFINITION,
				Literal: ".test",
				Line:    1,
			},
		},
	}
	got := program.Statements
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestParseProgram_ParsesConstantDefinition(t *testing.T) {
	t.Parallel()

	input := `CONS c 10`
	l := newLexerFromString(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}

	want := []ast.Statement{
		ast.ConstantDefinitionStatement{
			Token: token.Token{
				Type:    token.CONSTANT_DEFINITION,
				Literal: "CONS",
				Line:    1,
			},
			Name: ast.Identifier{
				Token: token.Token{
					Type:    token.IDENT,
					Literal: "c",
					Line:    1,
				},
				Value: "c",
			},
			Value: ast.IntegerLiteral{
				Token: token.Token{
					Type:    token.INT,
					Literal: "10",
					Line:    1,
				},
				Value: 10,
			},
		},
	}
	got := program.Statements
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestParseProgram_ParsesStringVariableDefinition(t *testing.T) {
	t.Parallel()

	input := `VARB msg "hello"`
	l := newLexerFromString(input)
	p := parser.New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}

	want := []ast.Statement{
		ast.VariableDefinitionStatement{
			Token: token.Token{
				Type:    token.VARIABLE_DEFINITION,
				Literal: "VARB",
				Line:    1,
			},
			Name: ast.Identifier{
				Token: token.Token{
					Type:    token.IDENT,
					Literal: "msg",
					Line:    1,
				},
				Value: "msg",
			},
			Value: ast.StringLiteral{
				Token: token.Token{
					Type:    token.STRING,
					Literal: "hello",
					Line:    1,
				},
				Value: "hello",
			},
		},
	}
	got := program.Statements
	if !cmp.Equal(want, got) {
		t.Error(want, got)
	}
}

func TestParseProgram_ParsesIntegerVariableDefinition(t *testing.T) {
	t.Parallel()

	input := `VARB num 100`
	l := newLexerFromString(input)
	p := parser.New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}

	want := []ast.Statement{
		ast.VariableDefinitionStatement{
			Token: token.Token{
				Type:    token.VARIABLE_DEFINITION,
				Literal: "VARB",
				Line:    1,
			},
			Name: ast.Identifier{
				Token: token.Token{
					Type:    token.IDENT,
					Literal: "num",
					Line:    1,
				},
				Value: "num",
			},
			Value: ast.IntegerLiteral{
				Token: token.Token{
					Type:    token.INT,
					Literal: "100",
					Line:    1,
				},
				Value: 100,
			},
		},
	}
	got := program.Statements
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestParseProgram_ParsesInstructionsWithoutOperand(t *testing.T) {
	t.Parallel()

	input := `
HALT
NOOP
OUTA
INCA
DECA
PSHA
POPA`

	l := newLexerFromString(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}

	want := []ast.Statement{
		ast.InstructionStatement{
			Token: token.Token{
				Type:    token.INSTRUCTION,
				Literal: "HALT",
				Line:    2,
			},
		},
		ast.InstructionStatement{
			Token: token.Token{
				Type:    token.INSTRUCTION,
				Literal: "NOOP",
				Line:    3,
			},
		},
		ast.InstructionStatement{
			Token: token.Token{
				Type:    token.INSTRUCTION,
				Literal: "OUTA",
				Line:    4,
			},
		},
		ast.InstructionStatement{
			Token: token.Token{
				Type:    token.INSTRUCTION,
				Literal: "INCA",
				Line:    5,
			},
		},
		ast.InstructionStatement{
			Token: token.Token{
				Type:    token.INSTRUCTION,
				Literal: "DECA",
				Line:    6,
			},
		},
		ast.InstructionStatement{
			Token: token.Token{
				Type:    token.INSTRUCTION,
				Literal: "PSHA",
				Line:    7,
			},
		},
		ast.InstructionStatement{
			Token: token.Token{
				Type:    token.INSTRUCTION,
				Literal: "POPA",
				Line:    8,
			},
		},
	}
	got := program.Statements
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestParseProgram_ReturnsErrorForInvalidOperand(t *testing.T) {
	t.Parallel()

	input := "SETA 2a"
	l := newLexerFromString(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}
	if len(p.Errors()) != 1 {
		t.Fatalf("parser returned %d errors. got=%d", len(p.Errors()), 1)
	}

	wantErr := parser.ErrInvalidIntegerLiteral
	err := p.Errors()[0]
	if !errors.Is(err, wantErr) {
		t.Fatalf("parser returned wrong error. got=%q, want=%q", err, wantErr)
	}
}

func TestParseProgram_ParsesInstructionsWithAnIntegerLiteralOperand(t *testing.T) {
	t.Parallel()

	input := `
SETA 42
JUMP 42
`
	l := newLexerFromString(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}

	want := []ast.Statement{
		ast.InstructionStatement{
			Token: token.Token{
				Type:    token.INSTRUCTION,
				Literal: "SETA",
				Line:    2,
			},
			Operand1: ast.IntegerLiteral{
				Token: token.Token{
					Type:    token.INT,
					Literal: "42",
					Line:    2,
				},
				Value: 42,
			},
		},
		ast.InstructionStatement{
			Token: token.Token{
				Type:    token.INSTRUCTION,
				Literal: "JUMP",
				Line:    3,
			},
			Operand1: ast.IntegerLiteral{
				Token: token.Token{
					Type:    token.INT,
					Literal: "42",
					Line:    3,
				},
				Value: 42,
			},
		},
	}
	got := program.Statements
	if !cmp.Equal(want, got) {
		t.Error(want, got)
	}
}

func TestParseProgram_ParsesInstructionWithAnIdentifierOperand(t *testing.T) {
	t.Parallel()

	input := "JUMP start"
	l := newLexerFromString(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}

	want := []ast.Statement{
		ast.InstructionStatement{
			Token: token.Token{
				Type:    token.INSTRUCTION,
				Literal: "JUMP",
				Line:    1,
			},
			Operand1: ast.Identifier{
				Token: token.Token{
					Type:    token.IDENT,
					Literal: "start",
					Line:    1,
				},
			},
		},
	}
	got := program.Statements
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestParseProgram_ParsesMoveInstructionWithRegisterAndIdentifierOperands(t *testing.T) {
	t.Parallel()

	input := `
MOVE A -> var
MOVE var -> A
`
	l := newLexerFromString(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}

	want := []ast.Statement{
		ast.InstructionStatement{
			Token: token.Token{
				Type:    token.INSTRUCTION,
				Literal: "MOVE",
				Line:    2,
			},
			Operand1: ast.RegisterLiteral{
				Token: token.Token{
					Type:    token.REGISTER,
					Literal: "A",
					Line:    2,
				},
			},
			Operand2: ast.Identifier{
				Token: token.Token{
					Type:    token.IDENT,
					Literal: "var",
					Line:    2,
				},
			},
		},
		ast.InstructionStatement{
			Token: token.Token{
				Type:    token.INSTRUCTION,
				Literal: "MOVE",
				Line:    3,
			},
			Operand1: ast.Identifier{
				Token: token.Token{
					Type:    token.IDENT,
					Literal: "var",
					Line:    3,
				},
			},
			Operand2: ast.RegisterLiteral{
				Token: token.Token{
					Type:    token.REGISTER,
					Literal: "A",
					Line:    3,
				},
			},
		},
	}
	got := program.Statements
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestParseProgram_ParsesInstructionsWithARegisterLiteralOperands(t *testing.T) {
	t.Parallel()

	input := `
MOVE A -> X
MOVE A -> Y
ADDA X
ADDA Y
`
	l := newLexerFromString(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}

	want := []ast.Statement{
		ast.InstructionStatement{
			Token: token.Token{
				Type:    token.INSTRUCTION,
				Literal: "MOVE",
				Line:    2,
			},
			Operand1: ast.RegisterLiteral{
				Token: token.Token{
					Type:    token.REGISTER,
					Literal: "A",
					Line:    2,
				},
			},
			Operand2: ast.RegisterLiteral{
				Token: token.Token{
					Type:    token.REGISTER,
					Literal: "X",
					Line:    2,
				},
			},
		},
		ast.InstructionStatement{
			Token: token.Token{
				Type:    token.INSTRUCTION,
				Literal: "MOVE",
				Line:    3,
			},
			Operand1: ast.RegisterLiteral{
				Token: token.Token{
					Type:    token.REGISTER,
					Literal: "A",
					Line:    3,
				},
			},
			Operand2: ast.RegisterLiteral{
				Token: token.Token{
					Type:    token.REGISTER,
					Literal: "Y",
					Line:    3,
				},
			},
		},
		ast.InstructionStatement{
			Token: token.Token{
				Type:    token.INSTRUCTION,
				Literal: "ADDA",
				Line:    4,
			},
			Operand1: ast.RegisterLiteral{
				Token: token.Token{
					Type:    token.REGISTER,
					Literal: "X",
					Line:    4,
				},
			},
		},
		ast.InstructionStatement{
			Token: token.Token{
				Type:    token.INSTRUCTION,
				Literal: "ADDA",
				Line:    5,
			},
			Operand1: ast.RegisterLiteral{
				Token: token.Token{
					Type:    token.REGISTER,
					Literal: "Y",
					Line:    5,
				},
			},
		},
	}

	got := program.Statements
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestParseProgram_ParsesInstructionWithACharacterLiteralOperand(t *testing.T) {
	t.Parallel()

	input := "SETA 'a'"
	l := newLexerFromString(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}
	want := []ast.Statement{
		ast.InstructionStatement{
			Token: token.Token{
				Type:    token.INSTRUCTION,
				Literal: "SETA",
				Line:    1,
			},
			Operand1: ast.CharacterLiteral{
				Token: token.Token{
					Type:    token.CHAR,
					Literal: "'a'",
					Line:    1,
				},
				Value: 'a',
			},
		},
	}
	got := program.Statements
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func newLexerFromString(input string) *lexer.Lexer {
	l, err := lexer.New(strings.NewReader(input))
	if err != nil {
		panic(err)
	}
	return l
}
