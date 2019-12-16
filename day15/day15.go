package main

import (
	"fmt"
	"github.com/dhconnelly/advent-of-code-2019/geom"
	"github.com/dhconnelly/advent-of-code-2019/intcode"
	"github.com/dhconnelly/advent-of-code-2019/ints"
	"log"
	"math"
	"os"
)

type direction int

const (
	NORTH direction = 1
	SOUTH direction = 2
	WEST  direction = 3
	EAST  direction = 4
)

func (dir direction) String() string {
	switch dir {
	case NORTH:
		return "NORTH"
	case SOUTH:
		return "SOUTH"
	case WEST:
		return "WEST"
	case EAST:
		return "EAST"
	}
	log.Fatal("bad direction:", dir)
	return ""
}

type status int

const (
	WALL status = 0
	OK   status = 1
	OXGN status = 2
)

func (stat status) String() string {
	switch stat {
	case WALL:
		return "WALL"
	case OK:
		return "OK"
	case OXGN:
		return "OXGN"
	}
	log.Fatalf("unknown status: %d", stat)
	return ""
}

var directions = map[direction]geom.Pt2{
	NORTH: geom.Pt2{0, 1},
	SOUTH: geom.Pt2{0, -1},
	WEST:  geom.Pt2{-1, 0},
	EAST:  geom.Pt2{1, 0},
}

func opposite(dir direction) direction {
	switch dir {
	case NORTH:
		return SOUTH
	case SOUTH:
		return NORTH
	case WEST:
		return EAST
	case EAST:
		return WEST
	}
	log.Fatal("bad direction:", dir)
	return 0
}

type droid struct {
	in  chan<- int64
	out <-chan int64
}

func (d *droid) step(dir direction) status {
	d.in <- int64(dir)
	return status(<-d.out)
}

func (d *droid) visit(p geom.Pt2, m map[geom.Pt2]status) {
	for dir, dp := range directions {
		next := p.Add(dp)
		if _, ok := m[next]; ok {
			continue
		}
		s := d.step(dir)
		if m[next] = s; s == WALL {
			continue
		}
		d.visit(next, m)
		d.step(opposite(dir))
	}
}

func bounds(m map[geom.Pt2]status) (minX, maxX, minY, maxY int) {
	minX, maxX = math.MaxInt64, math.MinInt64
	minY, maxY = math.MaxInt64, math.MinInt64
	for p := range m {
		minX = ints.Min(minX, p.X)
		maxX = ints.Max(maxX, p.X)
		minY = ints.Min(minY, p.Y)
		maxY = ints.Max(maxY, p.Y)
	}
	return
}

func printMap(m map[geom.Pt2]status) {
	minX, maxX, minY, maxY := bounds(m)
	for y := maxY; y >= minY; y-- {
		for x := minX; x <= maxX; x++ {
			if x == 0 && y == 0 {
				fmt.Print("O")
				continue
			}
			s, ok := m[geom.Pt2{x, y}]
			if !ok {
				fmt.Print(" ")
				continue
			}
			switch s {
			case WALL:
				fmt.Print("#")
			case OK:
				fmt.Print(".")
			case OXGN:
				fmt.Print("@")
			}
		}
		fmt.Println()
	}
}

func findOxygen(data []int64) int {
	in := make(chan int64)
	out := intcode.RunProgram(data, in)
	d := droid{in, out}
	m := map[geom.Pt2]status{geom.Zero2: OK}
	d.visit(geom.Zero2, m)
	printMap(m)
	return 0
}

func main() {
	data, err := intcode.ReadProgram(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(findOxygen(data))
}
