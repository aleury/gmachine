package parser_test

import (
	"errors"
	"gmachine/ast"
	"gmachine/lexer"
	"gmachine/parser"
	"strings"
	"testing"
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

	wantedStatments := 1
	gotStatements := len(program.Statements)
	if wantedStatments != gotStatements {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", wantedStatments, gotStatements)
	}

	tests := []struct {
		expectedLabelDefn string
	}{
		{".test"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]

		if stmt.TokenLiteral() != tt.expectedLabelDefn {
			t.Fatalf("stmt.TokenLiteral not %s. got=%q", tt.expectedLabelDefn, stmt.TokenLiteral())
		}

		_, ok := stmt.(*ast.LabelDefinitionStatement)
		if !ok {
			t.Fatalf("stmt not *ast.LabelDefinitionStatement. got=%T", stmt)
		}
	}
}

func TestParseProgram_ParsesOpcodesWithoutOperand(t *testing.T) {
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

	wantedStatments := 7
	gotStatements := len(program.Statements)
	if wantedStatments != gotStatements {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", wantedStatments, gotStatements)
	}

	tests := []struct {
		expectedOpcode string
	}{
		{"HALT"},
		{"NOOP"},
		{"OUTA"},
		{"INCA"},
		{"DECA"},
		{"PSHA"},
		{"POPA"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]

		if stmt.TokenLiteral() != tt.expectedOpcode {
			t.Fatalf("stmt.TokenLiteral not %s. got=%q", tt.expectedOpcode, stmt.TokenLiteral())
		}

		_, ok := stmt.(*ast.OpcodeStatement)
		if !ok {
			t.Fatalf("stmt not *ast.OpcodeStatement. got=%T", stmt)
		}
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

func TestParseProgram_ParsesOpcodesWithAnIntegerLiteralOperand(t *testing.T) {
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

	wantedStatments := 2
	gotStatements := len(program.Statements)
	if wantedStatments != gotStatements {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", wantedStatments, gotStatements)
	}

	tests := []struct {
		expectedOpcode  string
		expectedOperand uint64
	}{
		{"SETA", 42},
		{"JUMP", 42},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]

		if stmt.TokenLiteral() != tt.expectedOpcode {
			t.Fatalf("stmt.TokenLiteral not %s. got=%q", tt.expectedOpcode, stmt.TokenLiteral())
		}

		opcodeStmt, ok := stmt.(*ast.OpcodeStatement)
		if !ok {
			t.Fatalf("stmt not *ast.OpcodeStatement. got=%T", stmt)
		}

		operand := opcodeStmt.Operand
		if operand == nil {
			t.Fatalf("operand is nil")
		}

		operandExpr, ok := operand.(*ast.IntegerLiteral)
		if !ok {
			t.Fatalf("operand not *ast.Integer. got=%T", operand)
		}
		if operandExpr.Value != tt.expectedOperand {
			t.Fatalf("operand.Value not %d. got=%d", tt.expectedOperand, operandExpr.Value)
		}
	}
}

func TestParseProgram_ParsesOpcodeWithAnIdentifierOperand(t *testing.T) {
	t.Parallel()

	input := "JUMP start"
	l := newLexerFromString(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}

	wantedStatments := 1
	gotStatements := len(program.Statements)
	if wantedStatments != gotStatements {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", wantedStatments, gotStatements)
	}

	tests := []struct {
		expectedOpcode  string
		expectedLiteral string
	}{
		{"JUMP", "start"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]

		if stmt.TokenLiteral() != tt.expectedOpcode {
			t.Fatalf("stmt.TokenLiteral not %s. got=%q", tt.expectedOpcode, stmt.TokenLiteral())
		}

		opcodeStmt, ok := stmt.(*ast.OpcodeStatement)
		if !ok {
			t.Fatalf("stmt not *ast.OpcodeStatement. got=%T", stmt)
		}

		operand := opcodeStmt.Operand
		if operand == nil {
			t.Fatalf("operand is nil")
		}

		operandExpr, ok := operand.(*ast.Identifier)
		if !ok {
			t.Fatalf("operand not *ast.Identifier. got=%T", operand)
		}
		if operandExpr.Value != "" {
			t.Fatalf("operand.Value is not nil")
		}
	}
}

func TestParseProgram_ParsesOpcodeWithARegisterLiteralOperand(t *testing.T) {
	t.Parallel()

	input := `
MOVA X
MOVA Y
SETA X
SETA Y
ADDA X
ADDA Y	
`
	l := newLexerFromString(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}

	wantedStatments := 6
	gotStatements := len(program.Statements)
	if wantedStatments != gotStatements {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", wantedStatments, gotStatements)
	}

	tests := []struct {
		expectedOpcode   string
		expectedRegister string
	}{
		{"MOVA", "X"},
		{"MOVA", "Y"},
		{"SETA", "X"},
		{"SETA", "Y"},
		{"ADDA", "X"},
		{"ADDA", "Y"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]

		if stmt.TokenLiteral() != tt.expectedOpcode {
			t.Fatalf("stmt.TokenLiteral not %s. got=%q", tt.expectedOpcode, stmt.TokenLiteral())
		}

		opcodeStmt, ok := stmt.(*ast.OpcodeStatement)
		if !ok {
			t.Fatalf("stmt not *ast.OpcodeStatement. got=%T", stmt)
		}

		operand := opcodeStmt.Operand
		if operand == nil {
			t.Fatalf("operand is nil")
		}

		operandExpr, ok := operand.(*ast.RegisterLiteral)
		if !ok {
			t.Fatalf("operand not *ast.RegisterLiteral. got=%T", operand)
		}

		if operandExpr.TokenLiteral() != tt.expectedRegister {
			t.Fatalf("operand.TokenLiteral not %s. got=%q", tt.expectedRegister, operandExpr.TokenLiteral())
		}
	}
}

func TestParseProgram_ParsesOpcodeWithACharacterLiteralOperand(t *testing.T) {
	t.Parallel()

	input := "SETA 'a'"
	l := newLexerFromString(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}

	wantedStatements := 1
	gotStatements := len(program.Statements)
	if wantedStatements != gotStatements {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", wantedStatements, gotStatements)
	}

	tests := []struct {
		expectedOpcode  string
		expectedLiteral string
		expectedValue   rune
	}{
		{"SETA", "'a'", rune('a')},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]

		if stmt.TokenLiteral() != tt.expectedOpcode {
			t.Fatalf("stmt.TokenLiteral not %s. got=%q", tt.expectedOpcode, stmt.TokenLiteral())
		}

		opcodeStmt, ok := stmt.(*ast.OpcodeStatement)
		if !ok {
			t.Fatalf("stmt not *ast.OpcodeStatement. got=%T", stmt)
		}

		operand := opcodeStmt.Operand
		if operand == nil {
			t.Fatalf("operand is nil")
		}

		operandExpr, ok := operand.(*ast.CharacterLiteral)
		if !ok {
			t.Fatalf("operand not *ast.CharacterLiteral. got=%T", operand)
		}
		if operandExpr.TokenLiteral() != tt.expectedLiteral {
			t.Fatalf("operand.TokenLiteral not %s. got=%q", tt.expectedLiteral, operandExpr.TokenLiteral())
		}
		if operandExpr.Value != tt.expectedValue {
			t.Fatalf("operand.Value not %d. got=%d", tt.expectedValue, operandExpr.Value)
		}
	}
}

func newLexerFromString(input string) *lexer.Lexer {
	l, err := lexer.New(strings.NewReader(input))
	if err != nil {
		panic(err)
	}
	return l
}
