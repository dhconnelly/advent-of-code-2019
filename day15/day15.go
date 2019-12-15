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

func turnRight(prev direction) direction {
	switch prev {
	case NORTH:
		return EAST
	case SOUTH:
		return WEST
	case WEST:
		return NORTH
	case EAST:
		return SOUTH
	}
	log.Fatal("bad direction:", prev)
	return 0
}

var updates = map[direction]geom.Pt2{
	NORTH: geom.Pt2{0, 1},
	SOUTH: geom.Pt2{0, -1},
	WEST:  geom.Pt2{-1, 0},
	EAST:  geom.Pt2{1, 0},
}

func neighbors(p geom.Pt2) []geom.Pt2 {
	return []geom.Pt2{
		p.Add(updates[NORTH]), p.Add(updates[SOUTH]),
		p.Add(updates[WEST]), p.Add(updates[EAST]),
	}
}

func update(cur geom.Pt2, dir direction) geom.Pt2 {
	return cur.Add(updates[dir])
}

func copied(path []direction) []direction {
	path2 := make([]direction, len(path))
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

type droid struct {
	in   chan<- int64
	out  <-chan int64
	cur  geom.Pt2
	dist map[geom.Pt2]int
	path []geom.Pt2
}

func NewDroid(in chan<- int64, out <-chan int64) *droid {
	return &droid{in: in, out: out, dist: make(map[geom.Pt2]int)}
}

func (d *droid) move(dir direction) status {
	if dir == INVALID {
		return OK
	}
	d.in <- int64(dir)
	return status(<-d.out)
}

func (d *droid) visit(p geom.Pt2) {
	fmt.Println("visiting:", p)
	fmt.Println("path:", d.path)
	dir := dirOf(p, d.cur)
	switch d.move(dir) {
	case WALL:
		fmt.Println("wall at:", p)
		return
	case OK:
		fmt.Println("moved to:", p)
	case OXGN:
		fmt.Println("oxygen at:", p)
	}
	if dist, ok := d.dist[p]; !ok || len(d.path) < dist {
		d.dist[p] = len(d.path)
	} else {
		return
	}
	d.cur = p
	fmt.Println("dist:", d.dist[d.cur])
	for _, nbr := range neighbors(d.cur) {
		fmt.Println("next neighbor:", nbr)
		if len(d.path) > 0 && d.path[len(d.path)-1] == nbr {
			continue
		}
		d.path = append(d.path, nbr)
		d.visit(nbr)
		d.path = d.path[:len(d.path)-1]
	}
}

func findOxygen(data []int64) int {
	in := make(chan int64)
	out := intcode.RunProgram(data, in)
	d := NewDroid(in, out)
	d.visit(geom.Zero2)
	fmt.Println(d.dist)
	return 0
}

func main() {
	data, err := intcode.ReadProgram(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(findOxygen(data))
}
