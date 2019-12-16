package main

import (
	"fmt"
	"github.com/dhconnelly/advent-of-code-2019/geom"
	"github.com/dhconnelly/advent-of-code-2019/intcode"
	"github.com/dhconnelly/advent-of-code-2019/ints"
	"log"
	"os"
)

type status int

const (
	WALL status = 0
	OK   status = 1
	OXGN status = 2
)

type direction int

const (
	NORTH direction = 1
	SOUTH direction = 2
	WEST  direction = 3
	EAST  direction = 4
)

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

var directions = map[direction]geom.Pt2{
	NORTH: geom.Pt2{0, 1},
	SOUTH: geom.Pt2{0, -1},
	WEST:  geom.Pt2{-1, 0},
	EAST:  geom.Pt2{1, 0},
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

func shortestPaths(from geom.Pt2, m map[geom.Pt2]status) map[geom.Pt2]int {
	q := []node{{from, 0}}
	visited := make(map[geom.Pt2]bool)
	dist := make(map[geom.Pt2]int)
	var nd node
	for len(q) > 0 {
		nd, q = q[0], q[1:]
		dist[nd.p] = nd.n
		for _, dp := range directions {
			nbr := nd.p.Add(dp)
			if visited[nbr] {
				continue
			}
			visited[nbr] = true
			if m[nbr] != WALL {
				q = append(q, node{nbr, nd.n + 1})
			}
		}
	}
	return dist
}

func shortestPath(from, to geom.Pt2, m map[geom.Pt2]status) int {
	return shortestPaths(from, m)[to]
}

func longestPath(from geom.Pt2, m map[geom.Pt2]status) int {
	max := 0
	for _, n := range shortestPaths(from, m) {
		max = ints.Max(max, n)
	}
	return max
}

func main() {
	data, err := intcode.ReadProgram(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	m := explore(data)
	p := findOxygen(m)
	fmt.Println(shortestPath(geom.Zero2, p, m))
	fmt.Println(longestPath(p, m))
}
