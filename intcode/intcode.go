package intcode

import (
	"fmt"
	"log"
)

type Machine struct {
	b []int
}

func NewMachine(b []int) *Machine {
	return &Machine{b}
}

func (m *Machine) Get(i int) int {
	return m.b[i]
}

func (m *Machine) run() error {
	for i := 0; i < len(m.b); i += 4 {
		switch m.b[i] {
		case 1:
			m.b[m.b[i+3]] = m.b[m.b[i+1]] + m.b[m.b[i+2]]
		case 2:
			m.b[m.b[i+3]] = m.b[m.b[i+1]] * m.b[m.b[i+2]]
		case 99:
			return nil
		default:
			return fmt.Errorf("undefined opcode at position %d: %d", i, m.b[i])
		}
	}
	return nil
}

func (m *Machine) Execute(noun, verb int) (int, error) {
	b := make([]int, len(m.b))
	copy(b, m.b)
	m.b[1], m.b[2] = noun, verb
	if err := m.run(); err != nil {
		return 0, err
	}
	result := m.b[0]
	m.b = b
	return result, nil
}

func (m *Machine) ExecuteOrDie(noun, verb int) int {
	result, err := m.Execute(noun, verb)
	if err != nil {
		log.Fatal(err)
	}
	return result
}
