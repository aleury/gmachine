// Package gmachine implements a simple virtual CPU, known as the G-machine.
package gmachine

import (
	"errors"
	"fmt"
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

var ErrInvalidNumber error = errors.New("invalid number")
var ErrInvalidRegister error = errors.New("invalid register")
var ErrUndefinedInstruction error = errors.New("undefined instruction")

var registers = map[string]Word{
	"A": RegA,
	"X": RegX,
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
	copy(g.Memory[g.MemOffset+g.P:], program)
	g.Run()
}

func Assemble(input string) ([]Word, error) {
	program := []Word{}
	lines := strings.Split(strings.TrimSpace(input), "\n")
	for lineNo, line := range lines {
		if strings.HasPrefix(line, ";") {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		switch parts[0] {
		case "HALT":
			program = append(program, OpHALT)
		case "NOOP":
			program = append(program, OpNOOP)
		case "OUTA":
			program = append(program, OpOUTA)
		case "INCA":
			program = append(program, OpINCA)
		case "DECA":
			program = append(program, OpDECA)
		case "PSHA":
			program = append(program, OpPSHA)
		case "POPA":
			program = append(program, OpPOPA)
		case "ADDA":
			reg, ok := registers[parts[1]]
			if !ok {
				return nil, fmt.Errorf("%w: %s at line %d", ErrInvalidRegister, parts[1], lineNo+1)
			}
			program = append(program, OpADDA, reg)
		case "MOVA":
			reg, ok := registers[parts[1]]
			if !ok {
				return nil, fmt.Errorf("%w: %s at line %d", ErrInvalidRegister, parts[1], lineNo+1)
			}
			program = append(program, OpMOVA, reg)
		case "SETA":
			var operand Word
			if strings.HasPrefix(parts[1], "'") && strings.HasSuffix(parts[1], "'") {
				char, _ := utf8.DecodeRuneInString(strings.Trim(parts[1], "'"))
				operand = Word(char)
			} else {
				num, err := strconv.ParseUint(parts[1], 10, 64)
				if err != nil {
					return nil, fmt.Errorf("%w: %s at line %d", ErrInvalidNumber, parts[1], lineNo+1)
				}
				operand = Word(num)
			}
			program = append(program, OpSETA, operand)
		case "JUMP":
			loc, err := strconv.ParseUint(parts[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("%w: %s at line %d", ErrInvalidNumber, parts[1], lineNo+1)
			}
			program = append(program, OpJUMP, Word(loc))
		default:
			return nil, fmt.Errorf("%w: %s at line %d", ErrUndefinedInstruction, parts[0], lineNo+1)
		}
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
