package geom

import (
	"github.com/dhconnelly/advent-of-code-2019/ints"
)

type Pt3 struct {
	X, Y, Z int
}

var Zero3 Pt3

func (p1 *Pt3) TranslateBy(p2 Pt3) {
	p1.X += p2.X
	p1.Y += p2.Y
	p1.Z += p2.Z
}

func (p1 Pt3) Add(p2 Pt3) Pt3 {
	p3 := p1
	p3.TranslateBy(p2)
	return p3
}

func (p Pt3) IsZero() bool {
	return p.Eq(Zero3)
}

func (p Pt3) ManhattanNorm() int {
	return ints.Abs(p.X) + ints.Abs(p.Y) + ints.Abs(p.Z)
}

func (p1 Pt3) Eq(p2 Pt3) bool {
	return p1.X == p2.X && p1.Y == p2.Y && p1.Z == p2.Z
}
