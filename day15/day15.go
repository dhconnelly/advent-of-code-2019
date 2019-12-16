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

type status int

const (
	WALL status = 0
	OK   status = 1
	OXGN status = 2
)

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

func explore(prog []int64) map[geom.Pt2]status {
	in := make(chan int64)
	out := intcode.RunProgram(prog, in)
	d := droid{in, out}
	m := map[geom.Pt2]status{geom.Zero2: OK}
	d.visit(geom.Zero2, m)
	return m
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

func findOxygen(m map[geom.Pt2]status) geom.Pt2 {
	for p, s := range m {
		if s == OXGN {
			return p
		}
	}
	log.Fatal("oxygen not found")
	return geom.Zero2
}

type node struct {
	p geom.Pt2
	n int
}

func shortestPath(from, to geom.Pt2, m map[geom.Pt2]status) int {
	q := []node{{geom.Zero2, 0}}
	visited := make(map[geom.Pt2]bool)
	for len(q) > 0 {
		nd := q[0]
		q = q[1:]
		for _, dp := range directions {
			nbr := nd.p.Add(dp)
			if nbr == to {
				return nd.n + 1
			}
			if visited[nbr] {
				continue
			}
			visited[nbr] = true
			if m[nbr] != WALL {
				q = append(q, node{nbr, nd.n + 1})
			}
		}
	}
	log.Fatalf("no path from %v to %v", from, to)
	return 0
}

func main() {
	data, err := intcode.ReadProgram(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	m := explore(data)
	p := findOxygen(m)
	fmt.Println(shortestPath(geom.Zero2, p, m))
}
