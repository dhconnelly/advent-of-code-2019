package geom

import (
	"github.com/dhconnelly/advent-of-code-2019/ints"
)

type Pt2 struct {
	X, Y int
}

var Zero2 Pt2

func ManhattanDist2(pt1, pt2 Pt2) int {
	return ints.Abs(pt1.X-pt2.X) + ints.Abs(pt1.Y-pt2.Y)
}

func ManhattanNorm2(pt Pt2) int {
	return ManhattanDist2(pt, Zero2)
}
