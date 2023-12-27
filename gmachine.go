// Package gmachine implements a simple virtual CPU, known as the G-machine.
package gmachine

import (
	"bytes"
	"encoding/binary"
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
	RegY
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
	"Y": RegY,
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

// Floating Point Support
// Idea 1: Fixed point (32 bits for integer, 32 bits for fraction)
// Idea 2: Movable Point, select number of bits for integer part and fraction part
// Idea 3: Floating point (IEEE 754) (32 bits for mantissa, 8 bits for exponent, 1 bit for sign)

type Machine struct {
	P         Word
	S         Word
	A         Word
	X         Word
	Y         Word
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
		Y:         Word(0),
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
			binary.Write(g.Out, binary.BigEndian, g.A)
		case OpINCA:
			g.A++
		case OpDECA:
			g.A--
		case OpADDA:
			switch g.Next() {
			case RegX:
				g.A += g.X
			case RegY:
				g.A += g.Y
			}
		case OpMOVA:
			switch g.Next() {
			case RegX:
				g.X = g.A
			case RegY:
				g.Y = g.A
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
	Value   Word
}

type symbolTable struct {
	labels map[string]Word
	consts map[string]Word
}

func newSymbolTable() *symbolTable {
	return &symbolTable{
		labels: make(map[string]Word),
		consts: make(map[string]Word),
	}
}

func (t *symbolTable) defineLabel(name string, address Word) {
	t.labels[name] = address
}

func (t *symbolTable) defineConst(name string, value Word) {
	t.consts[name] = value
}

func (t *symbolTable) lookup(name string) (Word, bool) {
	if value, ok := t.labels[name]; ok {
		return value, ok
	}
	if value, ok := t.consts[name]; ok {
		return value, ok
	}
	return Word(0), false
}

func Assemble(reader io.Reader) ([]Word, error) {
	program := []Word{}
	refs := []Ref{}
	symbols := newSymbolTable()

	l, err := lexer.New(reader)
	if err != nil {
		return nil, err
	}
	p := parser.New(l)
	astProgram := p.ParseProgram()
	if astProgram == nil {
		return program, errors.New("failed to parse program")
	}
	if len(p.Errors()) > 0 {
		return program, p.Errors()[0]
	}

	// Assemble program
	for _, stmt := range astProgram.Statements {
		switch stmt := stmt.(type) {
		case *ast.ConstantDefinitionStatement:
			value := stmt.Value.(*ast.IntegerLiteral).Value
			symbols.defineConst(stmt.Name.Value, Word(value))
		case *ast.LabelDefinitionStatement:
			name := strings.TrimPrefix(stmt.TokenLiteral(), ".")
			symbols.defineLabel(name, Word(len(program)))
		case *ast.OpcodeStatement:
			program, refs, err = assembleOpcodeStatement(stmt, program, refs)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unknown statement type: %T", stmt)
		}
	}

	// Resolve references to labels and consts
	for _, refs := range refs {
		value, ok := symbols.lookup(refs.Name)
		if !ok {
			return nil, fmt.Errorf("%w: %s at line %d", ErrUnknownIdentifier, refs.Name, refs.Line)
		}
		program[refs.Address] = value
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
		if !slices.Contains([]Word{OpSETA, OpJUMP}, opcode) {
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

func (g *Machine) AssembleAndRun(r io.Reader) error {
	program, err := Assemble(r)
	if err != nil {
		return err
	}
	g.RunProgram(program)
	return nil
}

func RunFile(path string) int {
	content, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	g := New(os.Stdout)
	err = g.AssembleAndRun(content)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func Compile(in io.Reader, out io.Writer) error {
	program, err := Assemble(in)
	if err != nil {
		return err
	}

	err = binary.Write(out, binary.BigEndian, program)
	if err != nil {
		return err
	}

	return nil
}

func MainCompile() int {
	fileName := os.Args[1]
	outputFile := strings.TrimSuffix(fileName, ".g")

	in, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer in.Close()

	out, err := os.Create(outputFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer out.Close()

	err = Compile(in, out)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return 0
}

func MainRun() int {
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer f.Close()

	input, err := io.ReadAll(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	program := make([]Word, len(input)/8)
	err = binary.Read(bytes.NewReader(input), binary.BigEndian, &program)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	g := New(os.Stdout)
	g.RunProgram(program)
	if g.E != 0 {
		fmt.Fprintf(os.Stderr, "exception number: %d\n", g.E)
		return 1
	}

	return 0
}
