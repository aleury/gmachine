// Package gmachine implements a simple virtual CPU, known as the G-machine.
package gmachine

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// TODO(adam): Research serial output to add support outputing characters from the gmachine

const MemSize = 1024

const (
	OpHALT Word = iota + 1
	OpNOOP
	OpINCA
	OpDECA
	OpSETA
)

const (
	ExceptionOK Word = iota
	ExceptionIllegalInstruction
	ExceptionOutOfMemory
)

var ErrInvalidNumber error = errors.New("invalid number")
var ErrUndefinedInstruction error = errors.New("undefined instruction")

type Word uint64

type Machine struct {
	P      Word
	A      Word
	E      Word
	Memory []Word
}

func New() *Machine {
	return &Machine{
		P:      Word(0),
		A:      Word(0),
		E:      Word(0),
		Memory: make([]Word, MemSize),
	}
}

func (g *Machine) Run() {
	for {
		instruction := g.Memory[g.P]
		g.P++
		if g.P >= MemSize {
			g.E = ExceptionOutOfMemory
			return
		}

		switch instruction {
		case OpHALT:
			return
		case OpNOOP:
			continue
		case OpINCA:
			g.A++
		case OpDECA:
			g.A--
		case OpSETA:
			g.A = g.Memory[g.P]
			g.P++
		default:
			g.E = ExceptionIllegalInstruction
			return
		}
	}
}

func (g *Machine) RunProgram(program []Word) {
	copy(g.Memory[g.P:], program)
	g.Run()
}

func Assemble(input string) ([]Word, error) {
	program := []Word{}
	for lineNo, instruction := range strings.Split(strings.TrimSpace(input), "\n") {
		parts := strings.Split(instruction, " ")
		switch parts[0] {
		case "HALT":
			program = append(program, OpHALT)
		case "NOOP":
			program = append(program, OpNOOP)
		case "INCA":
			program = append(program, OpINCA)
		case "DECA":
			program = append(program, OpDECA)
		case "SETA":
			num, err := strconv.ParseUint(parts[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("%w: %s at line %d", ErrInvalidNumber, parts[1], lineNo+1)
			}
			program = append(program, OpSETA, Word(num))
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
	g := New()
	err = g.AssembleAndRun(string(content))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
