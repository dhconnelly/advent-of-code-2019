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

func (pt1 Pt2) Dist(pt2 Pt2) float64 {
	pt := pt1
	pt.TranslateBy(Pt2{-pt2.X, -pt2.Y})
	return pt.Norm()
}

func (pt Pt2) Norm() float64 {
	return math.Sqrt(math.Pow(float64(pt.X), 2.0) + math.Pow(float64(pt.Y), 2.0))
}

type Direction rune

const (
	None  Direction = 0
	Up    Direction = '^'
	Down  Direction = 'v'
	Right Direction = '>'
	Left  Direction = '<'
)

var Directions = map[Direction]Pt2{
	Up:    Pt2{0, 1},
	Down:  Pt2{0, -1},
	Left:  Pt2{-1, 0},
	Right: Pt2{1, 0},
}

func (from Pt2) DirectionTo(to Pt2) Direction {
	switch {
	case from.X < to.X:
		return Right
	case to.X < from.X:
		return Left
	case from.Y < to.Y:
		return Down
	case to.Y < from.Y:
		return Up
	}
	return 0
}

func (p Pt2) Go(dir Direction) Pt2 {
	return p.Add(Directions[dir])
}

func (p Pt2) ManhattanNeighbors() []Pt2 {
	return []Pt2{
		p.Add(Directions[Up]),
		p.Add(Directions[Down]),
		p.Add(Directions[Left]),
		p.Add(Directions[Right]),
	}
}
