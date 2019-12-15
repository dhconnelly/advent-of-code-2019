package main

import (
	"fmt"
	"github.com/dhconnelly/advent-of-code-2019/geom"
	"github.com/dhconnelly/advent-of-code-2019/intcode"
	"log"
	"os"
)

type direction int

const (
	INVALID direction = -1

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

var updates = map[direction]geom.Pt2{
	NORTH: geom.Pt2{0, 1},
	SOUTH: geom.Pt2{0, -1},
	WEST:  geom.Pt2{-1, 0},
	EAST:  geom.Pt2{1, 0},
}

func neighbors(nd node) []node {
	var nbrs []node
	for k := range updates {
		var nbr node
		nbr.p = nd.p.Add(updates[k])
		nbr.path = append(copied(nd.path), nbr.p)
		nbrs = append(nbrs, nbr)
	}
	return nbrs
}

type node struct {
	p    geom.Pt2
	path []geom.Pt2
}

func copied(path []geom.Pt2) []geom.Pt2 {
	path2 := make([]geom.Pt2, len(path))
	copy(path2, path)
	return path2
}

func dirOf(to, from geom.Pt2) direction {
	switch {
	case to.X < from.X:
		return WEST
	case to.X > from.X:
		return EAST
	case to.Y < from.Y:
		return SOUTH
	case to.Y > from.Y:
		return NORTH
	}
	return INVALID
}

func reversed(path []geom.Pt2) []geom.Pt2 {
	reversed := make([]geom.Pt2, len(path))
	for i, p := range path {
		reversed[len(path)-i-1] = p
	}
	return reversed
}

type droid struct {
	in   chan<- int64
	out  <-chan int64
	path map[geom.Pt2][]geom.Pt2
	stat map[geom.Pt2]status
	next []node
}

func NewDroid(in chan<- int64, out <-chan int64) *droid {
	return &droid{
		in:   in,
		out:  out,
		path: make(map[geom.Pt2][]geom.Pt2),
		stat: make(map[geom.Pt2]status),
	}
}

func (d *droid) step(to, from geom.Pt2) status {
	if to == from {
		return OK
	}
	if to.ManhattanDist(from) > 1 {
		log.Fatalf("ManhattanDist(%v, %v) > 1", to, from)
	}
	dir := dirOf(to, from)
	d.in <- int64(dir)
	return status(<-d.out)
}

func (d *droid) move(path []geom.Pt2) status {
	from := path[0]
	s := OK
	for _, to := range path {
		if s = d.step(to, from); s == WALL {
			break
		}
		from = to
	}
	return s
}

func (d *droid) visit(nd node) {
	fmt.Println("visiting:", nd)

	// try to move to the node if possible
	stat := d.move(nd.path)
	if d.stat[nd.p] = stat; stat == WALL {
		fmt.Println("wall:", nd.p)
		return
	}

	// store the path for this node
	fmt.Println("path:", nd.path)
	d.path[nd.p] = nd.path

	// add neighbors to the visit queue
	for _, nbr := range neighbors(nd) {
		if _, ok := d.path[nbr.p]; ok {
			continue
		}
		fmt.Println("next neighbor:", nbr)
		d.next = append(d.next, nbr)
	}

	// go back to where we came from
	if s := d.move(reversed(nd.path)); s != OK {
		log.Fatalf("can't return by path: %v, %s", nd.path, s)
	}
}

func (d *droid) findOxygen() int {
	for len(d.next) > 0 {
		fmt.Println()
		fmt.Println("q:", d.next)
		fmt.Println("stat:", d.stat)
		next := d.next[0]
		d.next = d.next[1:]
		d.visit(next)
		if d.stat[next.p] == OXGN {
			return len(d.path[next.p])
		}
	}
	log.Fatal("not found")
	return 0
}

func findOxygen(data []int64) int {
	in := make(chan int64)
	out := intcode.RunProgram(data, in)
	d := NewDroid(in, out)
	d.next = append(d.next, node{geom.Zero2, []geom.Pt2{geom.Zero2}})
	return d.findOxygen()
}

func main() {
	data, err := intcode.ReadProgram(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(findOxygen(data))
}
