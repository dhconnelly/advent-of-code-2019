package geom

import (
	"github.com/dhconnelly/advent-of-code-2019/ints"
	"math"
)

type Pt2 struct {
	X, Y int
}

var Zero2 Pt2

func (pt1 Pt2) Add(pt2 Pt2) Pt2 {
	pt := pt1
	pt.TranslateBy(pt2)
	return pt
}

func (pt1 *Pt2) TranslateBy(pt2 Pt2) {
	pt1.X += pt2.X
	pt1.Y += pt2.Y
}

func (pt1 Pt2) ManhattanDist(pt2 Pt2) int {
	return ints.Abs(pt1.X-pt2.X) + ints.Abs(pt1.Y-pt2.Y)
}

func (pt Pt2) ManhattanNorm() int {
	return pt.ManhattanDist(Zero2)
}

func (pt Pt2) Norm() float64 {
	return math.Sqrt(math.Pow(float64(pt.X), 2.0) + math.Pow(float64(pt.Y), 2.0))
}
