package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/dhconnelly/advent-of-code-2019/geom"
)

const (
	empty   = ' '
	passage = '.'
	wall    = '#'
)

type grid struct {
	width  int
	height int
	g      map[geom.Pt2]rune
}

func readGrid(r io.Reader) grid {
	scan := bufio.NewScanner(r)
	g := grid{g: make(map[geom.Pt2]rune)}
	for scan.Scan() {
		g.width = 0
		line := scan.Text()
		for _, c := range line {
			g.g[geom.Pt2{g.width, g.height}] = c
			g.width++
		}
		g.height++
	}
	if err := scan.Err(); err != nil {
		log.Fatal(err)
	}
	return g
}

func mazeBounds(g grid) (outer, inner geom.Rect) {
outerLo:
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			p := geom.Pt2{x, y}
			if g.g[p] == wall {
				outer.Lo = p
				break outerLo
			}
		}
	}
outerHi:
	for y := g.height - 1; y >= 0; y-- {
		for x := g.width - 1; x >= 0; x-- {
			p := geom.Pt2{x, y}
			if g.g[p] == wall {
				outer.Hi = p
				break outerHi
			}
		}
	}
innerLo:
	for y := outer.Lo.Y; y <= outer.Hi.Y; y++ {
		for x := outer.Lo.X; x <= outer.Hi.X; x++ {
			p := geom.Pt2{x, y}
			if g.g[p] == empty {
				inner.Lo = p
				break innerLo
			}
		}
	}
innerHi:
	for y := outer.Hi.Y; y >= outer.Lo.Y; y-- {
		for x := outer.Hi.X; x >= outer.Lo.X; x-- {
			p := geom.Pt2{x, y}
			if g.g[p] == empty {
				inner.Hi = p
				break innerHi
			}
		}
	}
	return
}

type label struct{}

type maze struct {
	g      grid
	outer  geom.Rect
	inner  geom.Rect
	labels map[label][]geom.Pt2
}

func readMaze(g grid) maze {
	outer, inner := mazeBounds(g)
	fmt.Println(mazeBounds(g))
	return maze{g, inner, outer, nil}
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	g := readGrid(f)
	readMaze(g)
}
