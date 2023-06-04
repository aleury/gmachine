// Package gmachine implements a simple virtual CPU, known as the G-machine.
package gmachine

const DefaultMemSize = 256

const (
	OpHALT = 0
	OpNOOP = 1
)

type Word byte

type Machine struct {
	P      Word
	Memory []Word
}

func New() *Machine {
	return &Machine{
		P:      Word(0),
		Memory: make([]Word, DefaultMemSize),
	}
}

func (m *Machine) Run() {
	for {
		instruction := m.Memory[m.P]
		m.P += 1

		if instruction == OpHALT {
			return
		}
		if instruction == OpNOOP {
			continue
		}
	}
}
