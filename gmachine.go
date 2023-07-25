// Package gmachine implements a simple virtual CPU, known as the G-machine.
package gmachine

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
