package intcode

const (
	add    opcode = 1
	mul           = 2
	read          = 3
	print         = 4
	jmpif         = 5
	jmpnot        = 6
	lt            = 7
	eq            = 8
	adjrel        = 9
	halt          = 99
)

func (o opcode) String() string {
	switch o {
	case add:
		return "add"
	case mul:
		return "mul"
	case read:
		return "read"
	case print:
		return "print"
	case jmpif:
		return "jmpif"
	case jmpnot:
		return "jmpnot"
	case lt:
		return "lt"
	case eq:
		return "eq"
	case adjrel:
		return "adjrel"
	case halt:
		return "halt"
	}
	return ""
}

var opcodeToArity = map[opcode]int64{
	add:    3,
	mul:    3,
	read:   1,
	print:  1,
	jmpif:  2,
	jmpnot: 2,
	lt:     3,
	eq:     3,
	adjrel: 1,
	halt:   0,
}

type handler func(m *machine, instr instruction) bool

var handlers = map[opcode]handler{
	add: func(m *machine, instr instruction) bool {
		l := m.get(m.pc+1, instr.modes[0])
		r := m.get(m.pc+2, instr.modes[1])
		m.set(m.pc+3, l+r, instr.modes[2])
		m.pc += instr.arity + 1
		return true
	},

	mul: func(m *machine, instr instruction) bool {
		l := m.get(m.pc+1, instr.modes[0])
		r := m.get(m.pc+2, instr.modes[1])
		m.set(m.pc+3, l*r, instr.modes[2])
		m.pc += instr.arity + 1
		return true
	},

	read: func(m *machine, instr instruction) bool {
		m.set(m.pc+1, m.read(), instr.modes[0])
		m.pc += instr.arity + 1
		return true
	},

	print: func(m *machine, instr instruction) bool {
		v := m.get(m.pc+1, instr.modes[0])
		m.write(v)
		m.pc += instr.arity + 1
		return true
	},

	jmpif: func(m *machine, instr instruction) bool {
		l := m.get(m.pc+1, instr.modes[0])
		r := m.get(m.pc+2, instr.modes[1])
		if l != 0 {
			m.pc = r
		} else {
			m.pc += instr.arity + 1
		}
		return true
	},

	jmpnot: func(m *machine, instr instruction) bool {
		l := m.get(m.pc+1, instr.modes[0])
		r := m.get(m.pc+2, instr.modes[1])
		if l == 0 {
			m.pc = r
		} else {
			m.pc += instr.arity + 1
		}
		return true
	},

	lt: func(m *machine, instr instruction) bool {
		l := m.get(m.pc+1, instr.modes[0])
		r := m.get(m.pc+2, instr.modes[1])
		var val int64
		if l < r {
			val = 1
		} else {
			val = 0
		}
		m.set(m.pc+3, val, instr.modes[2])
		m.pc += instr.arity + 1
		return true
	},

	eq: func(m *machine, instr instruction) bool {
		l := m.get(m.pc+1, instr.modes[0])
		r := m.get(m.pc+2, instr.modes[1])
		var val int64
		if l == r {
			val = 1
		} else {
			val = 0
		}
		m.set(m.pc+3, val, instr.modes[2])
		m.pc += instr.arity + 1
		return true
	},

	adjrel: func(m *machine, instr instruction) bool {
		v := m.get(m.pc+1, instr.modes[0])
		m.relbase += v
		m.pc += instr.arity + 1
		return true
	},

	halt: func(m *machine, instr instruction) bool {
		return false
	},
}
