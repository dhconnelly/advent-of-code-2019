package intcode

type Mode int

const (
	pos Mode = iota
	imm
	rel
)

func (md Mode) String() string {
	switch md {
	case pos:
		return "pos"
	case imm:
		return "imm"
	case rel:
		return "rel"
	}
	return ""
}

type Opcode int

type instruction struct {
	op    Opcode
	arity int64
	modes []Mode
}

func parseInstruction(i int64) instruction {
	var in instruction
	in.op = Opcode(i % 100)
	in.arity = opcodeToArity[in.op]
	for i /= 100; int64(len(in.modes)) < in.arity; i /= 10 {
		in.modes = append(in.modes, Mode(i%10))
	}
	return in
}
