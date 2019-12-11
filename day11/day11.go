package main

import (
	"fmt"
	"github.com/dhconnelly/advent-of-code-2019/geom"
	"github.com/dhconnelly/advent-of-code-2019/intcode"
	"log"
	"os"
)

type color int64

const (
	BLACK color = iota
	WHITE
)

type direction int64

const (
	TURN_LEFT direction = iota
	TURN_RIGHT
)

type orientation rune

const (
	LEFT orientation = iota
	UP
	RIGHT
	DOWN
)

func (o orientation) String() string {
	switch o {
	case LEFT:
		return "<"
	case UP:
		return "^"
	case RIGHT:
		return ">"
	case DOWN:
		return "v"
	default:
		log.Fatalf("bad orientation: %d", o)
	}
	return ""
}

func turn(cur orientation, dir direction) orientation {
	switch dir {
	case TURN_LEFT:
		if cur == LEFT {
			return DOWN
		}
		return cur - 1
	case TURN_RIGHT:
		return (cur + 1) % 4
	}
	log.Fatalf("bad dir: %d", dir)
	return -1
}

var diffs = map[orientation]geom.Pt2{
	LEFT:  geom.Pt2{-1, 0},
	UP:    geom.Pt2{0, 1},
	RIGHT: geom.Pt2{1, 0},
	DOWN:  geom.Pt2{0, -1},
}

func move(cur geom.Pt2, o orientation) geom.Pt2 {
	return cur.Add(diffs[o])
}

type grid map[geom.Pt2]color

func run(data []int64) grid {
	in := make(chan int64)
	out := intcode.RunProgram(data, in)
	g := grid(make(map[geom.Pt2]color))
	loc := geom.Zero2
	o := UP
loop:
	for {
		select {
		case c, ok := <-out:
			if !ok {
				break loop
			}
			g[loc] = color(c)
			dir := direction(<-out)
			o = turn(o, dir)
			loc = move(loc, o)
		case in <- int64(g[loc]):
		}
	}
	return g
}

func main() {
	data, err := intcode.ReadProgram(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	g := run(data)
	fmt.Println(len(g))
}
