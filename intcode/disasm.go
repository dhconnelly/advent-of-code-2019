package intcode

import (
	"fmt"
	"strings"
)

type Kind int

const (
	RawData Kind = iota + 1
	Instr
)

type Instruction struct {
	Opcode Opcode
	Modes  []Mode
}

type Line struct {
	Offset int
	Width  int

	Which Kind
	Data  []int64
	Instr Instruction
}

func joinInts(ints []int64) string {
	strs := make([]string, len(ints))
	for i, v := range ints {
		strs[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(strs, ", ")
}

func (l Line) String() string {
	switch l.Which {
	case Instr:
		s := fmt.Sprintf("[%4d] %8s    ", l.Offset, l.Instr.Opcode)
		args := make([]string, len(l.Instr.Modes))
		for i, md := range l.Instr.Modes {
			args[i] = fmt.Sprintf("%8s", fmt.Sprintf("%s(%d)", md, l.Data[i+1]))
		}
		return s + strings.Join(args, " ")
	default:
		return fmt.Sprintf("[%4d] %v", l.Offset, l.Data)
	}
}

func Disassemble(data []int64) []Line {
	var lines []Line
	var line Line
	for i := 0; i < len(data); {
		line.Offset = i
		instr := parseInstruction(data[i])
		if !instr.op.isValid() {
			line.Which = RawData
			line.Width = 1
		} else {
			line.Which = Instr
			line.Instr.Opcode = instr.op
			line.Instr.Modes = instr.modes
			line.Width = int(instr.arity) + 1
		}
		line.Data = data[i : i+line.Width]
		lines = append(lines, line)
		i += line.Width
	}
	return lines
}
