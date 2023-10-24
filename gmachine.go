// Package gmachine implements a simple virtual CPU, known as the G-machine.
package gmachine

import (
	"errors"
	"fmt"
	"gmachine/ast"
	"gmachine/lexer"
	"gmachine/parser"
	"io"
	"os"
	"strings"

	"golang.org/x/exp/slices"
)

// TODO(adam): Research serial output to add support outputing characters from the gmachine

const MemSize = 1024
const StackSize = 256

const (
	OpHALT Word = iota + 1
	OpNOOP
	OpOUTA
	OpINCA
	OpDECA
	OpADDA
	OpMOVA
	OpSETA
	OpPSHA
	OpPOPA
	OpJUMP
)

const (
	RegA Word = iota
	RegX
)

const (
	ExceptionOK Word = iota
	ExceptionIllegalInstruction
	ExceptionOutOfMemory
)

var ErrInvalidOperand error = errors.New("invalid operand")
var ErrUnknownIdentifier error = errors.New("missing label")
var ErrInvalidNumber error = errors.New("invalid number")
var ErrInvalidRegister error = errors.New("invalid register")
var ErrUndefinedInstruction error = errors.New("undefined instruction")

var registers = map[string]Word{
	"A": RegA,
	"X": RegX,
}

var opcodes = map[string]Word{
	"HALT": OpHALT,
	"NOOP": OpNOOP,
	"OUTA": OpOUTA,
	"INCA": OpINCA,
	"DECA": OpDECA,
	"ADDA": OpADDA,
	"MOVA": OpMOVA,
	"SETA": OpSETA,
	"PSHA": OpPSHA,
	"POPA": OpPOPA,
	"JUMP": OpJUMP,
}

type Word uint64

type Machine struct {
	P         Word
	S         Word
	A         Word
	X         Word
	E         Word
	Out       io.Writer
	MemOffset Word
	Memory    []Word
}

func New(out io.Writer) *Machine {
	return &Machine{
		P:         Word(0),
		S:         Word(0),
		A:         Word(0),
		X:         Word(0),
		E:         Word(0),
		Out:       out,
		MemOffset: StackSize,
		Memory:    make([]Word, MemSize),
	}
}

func (g *Machine) Next() Word {
	word := g.Memory[g.MemOffset+g.P]
	g.P++
	return word
}

func (g *Machine) Run() {
	for {
		instruction := g.Next()
		if g.MemOffset+g.P >= MemSize {
			g.E = ExceptionOutOfMemory
			return
		}

		switch instruction {
		case OpHALT:
			return
		case OpNOOP:
			continue
		case OpOUTA:
			g.Out.Write([]byte{byte(g.A)})
		case OpINCA:
			g.A++
		case OpDECA:
			g.A--
		case OpADDA:
			switch g.Next() {
			case RegX:
				g.A += g.X
			}
		case OpMOVA:
			switch g.Next() {
			case RegX:
				g.X = g.A
			}
		case OpSETA:
			g.A = g.Next()
		case OpPSHA:
			g.Memory[g.S] = g.A
			g.S++
		case OpPOPA:
			g.S--
			g.A = g.Memory[g.S]
		case OpJUMP:
			g.P = g.Memory[g.MemOffset+g.P]
		default:
			g.E = ExceptionIllegalInstruction
			return
		}
	}
}

func (g *Machine) RunProgram(program []Word) {
	// Load program into machine
	copy(g.Memory[g.MemOffset:], program)
	g.Run()
}

type Ref struct {
	Name    string
	Line    int
	Address Word
}

func Assemble(input string) ([]Word, error) {
	program := []Word{}
	refs := []Ref{}
	labels := map[string]Word{}

	l := lexer.New(input)
	p := parser.New(l)
	astProgram := p.ParseProgram()
	if astProgram == nil {
		return program, errors.New("failed to parse program")
	}
	if len(p.Errors()) > 0 {
		return program, p.Errors()[0]
	}

	// Assemble program
	var err error
	for _, stmt := range astProgram.Statements {
		switch stmt := stmt.(type) {
		case *ast.LabelDefinitionStatement:
			labels[strings.TrimPrefix(stmt.TokenLiteral(), ".")] = Word(len(program))
		case *ast.OpcodeStatement:
			program, refs, err = assembleOpcodeStatement(stmt, program, refs)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unknown statement type: %T", stmt)
		}
	}

	// Resolve references
	for _, refs := range refs {
		addr, ok := labels[refs.Name]
		if !ok {
			return nil, fmt.Errorf("%w: %s at line %d", ErrUnknownIdentifier, refs.Name, refs.Line)
		}
		program[refs.Address] = addr
	}

	return program, nil
}

func assembleOpcodeStatement(stmt *ast.OpcodeStatement, program []Word, refs []Ref) ([]Word, []Ref, error) {
	opcode, ok := opcodes[stmt.TokenLiteral()]
	if !ok {
		return nil, nil, fmt.Errorf("%w: %s at line %d", ErrUndefinedInstruction, stmt.TokenLiteral(), stmt.Token.Line)
	}
	program = append(program, opcode)

	if stmt.Operand == nil {
		return program, refs, nil
	}

	switch operand := stmt.Operand.(type) {
	case *ast.RegisterLiteral:
		if !slices.Contains([]Word{OpADDA, OpMOVA}, opcode) {
			return nil, nil, fmt.Errorf("%w: %s at line %d", ErrInvalidOperand, stmt.TokenLiteral(), stmt.Token.Line)
		}
		register, ok := registers[operand.TokenLiteral()]
		if !ok {
			return nil, nil, fmt.Errorf("%w: %s at line %d", ErrInvalidRegister, operand.TokenLiteral(), operand.Token.Line)
		}
		program = append(program, register)
	case *ast.Identifier:
		if opcode != OpJUMP {
			return nil, nil, fmt.Errorf("%w: %s at line %d", ErrInvalidOperand, stmt.TokenLiteral(), stmt.Token.Line)
		}
		ref := Ref{
			Name:    operand.TokenLiteral(),
			Line:    operand.Token.Line,
			Address: Word(len(program)),
		}
		refs = append(refs, ref)
		program = append(program, Word(0))
	case *ast.IntegerLiteral:
		if !slices.Contains([]Word{OpSETA, OpJUMP}, opcode) {
			return nil, nil, fmt.Errorf("%w: %s at line %d", ErrInvalidOperand, stmt.TokenLiteral(), stmt.Token.Line)
		}
		program = append(program, Word(operand.Value))
	case *ast.CharacterLiteral:
		if opcode != OpSETA {
			return nil, nil, fmt.Errorf("%w: %s at line %d", ErrInvalidOperand, stmt.TokenLiteral(), stmt.Token.Line)
		}
		program = append(program, Word(operand.Value))
	default:
		return nil, nil, fmt.Errorf("%w: %s at line %d", ErrInvalidOperand, stmt.TokenLiteral(), stmt.Token.Line)
	}

	return program, refs, nil
}

func (g *Machine) AssembleAndRun(input string) error {
	program, err := Assemble(input)
	if err != nil {
		return err
	}
	g.RunProgram(program)
	return nil
}

func RunFile(path string) int {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	g := New(os.Stdout)
	err = g.AssembleAndRun(string(content))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
