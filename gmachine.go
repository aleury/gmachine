// Package gmachine implements a simple virtual CPU, known as the G-machine.
package gmachine

import (
	"errors"
	"fmt"
	"gmachine/lexer"
	"gmachine/token"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
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
	l := lexer.New(input)
	program := []Word{}
	refs := []Ref{}
	labels := map[string]Word{}
loop:
	for {
		tok, err := l.NextToken()
		if err != nil {
			return nil, err
		}

		switch tok.Type {
		case token.ILLEGAL:
			return nil, errors.New("illegal token")
		case token.EOF:
			break loop
		case token.IDENT:
			if strings.HasPrefix(tok.Literal, ".") {
				labels[strings.TrimPrefix(tok.Literal, ".")] = Word(len(program))
			} else {
				ref := Ref{
					Name:    tok.Literal,
					Line:    tok.Line,
					Address: Word(len(program)),
				}
				refs = append(refs, ref)
				program = append(program, Word(0))
			}
		case token.OPCODE:
			opcode, ok := opcodes[tok.Literal]
			if !ok {
				return nil, fmt.Errorf("%w: %s at line %d", ErrUndefinedInstruction, tok.Literal, tok.Line)
			}
			program = append(program, opcode)

			switch opcode {
			case OpADDA, OpMOVA:
				operandTok, err := l.NextToken()
				if err != nil {
					return nil, err
				}
				reg, ok := registers[operandTok.Literal]
				if !ok {
					return nil, fmt.Errorf("%w: %s at line %d", ErrInvalidRegister, operandTok.Literal, operandTok.Line)
				}
				program = append(program, reg)
			case OpSETA:
				operandTok, err := l.NextToken()
				if err != nil {
					return nil, err
				}
				switch operandTok.Type {
				case token.INT:
					num, err := strconv.ParseUint(operandTok.Literal, 10, 64)
					if err != nil {
						return nil, fmt.Errorf("%w: %s at line %d", ErrInvalidNumber, operandTok.Literal, operandTok.Line)
					}
					program = append(program, Word(num))
				case token.CHAR:
					char, _ := utf8.DecodeRuneInString(strings.Trim(operandTok.Literal, "'"))
					program = append(program, Word(char))
				default:
					return nil, fmt.Errorf("%w: %s at line %d", ErrInvalidOperand, operandTok.Literal, operandTok.Line)
				}
			case OpJUMP:
				operandTok, err := l.NextToken()
				if err != nil {
					return nil, err
				}
				switch operandTok.Type {
				case token.INT:
					num, err := strconv.ParseUint(operandTok.Literal, 10, 64)
					if err != nil {
						return nil, fmt.Errorf("%w: %s at line %d", ErrInvalidNumber, operandTok.Literal, operandTok.Line)
					}
					program = append(program, Word(num))
				case token.IDENT:
					ref := Ref{
						Name:    operandTok.Literal,
						Line:    operandTok.Line,
						Address: Word(len(program)),
					}
					refs = append(refs, ref)
					program = append(program, Word(0))
				default:
					return nil, fmt.Errorf("%w: %s at line %d", ErrInvalidOperand, operandTok.Literal, operandTok.Line)
				}
			}
		}
	}

	// Resolve references
	for _, ref := range refs {
		addr, ok := labels[ref.Name]
		if !ok {
			return nil, fmt.Errorf("%w: %s at line %d", ErrUnknownIdentifier, ref.Name, ref.Line)
		}
		program[ref.Address] = addr
	}
	return program, nil
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
