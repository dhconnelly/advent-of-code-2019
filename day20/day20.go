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

type label [2]rune

func (lbl label) String() string {
	return string(lbl[:])
}

func isLabel(g grid, p geom.Pt2) bool {
	c := g.g[p]
	return 'A' <= c && c <= 'Z'
}

func flip(lbl label) label {
	return label([2]rune{lbl[1], lbl[0]})
}

func addLabel(
	g grid,
	from geom.Pt2,
	dir geom.Direction,
	reversed bool,
	adjs map[label][]geom.Pt2,
	lbls map[geom.Pt2]label,
) {
	p := from.Go(dir)
	if !isLabel(g, p) {
		return
	}
	a, b := g.g[p], g.g[p.Go(dir)]
	lbl := label([2]rune{a, b})
	if dir == geom.Up || dir == geom.Left {
		lbl = flip(lbl)
	}
	if reversed {
		lbl = flip(lbl)
	}
	adjs[lbl] = append(adjs[lbl], from)
	lbls[p] = lbl
}

func findOuterLabels(
	g grid,
	r geom.Rect,
	adjs map[label][]geom.Pt2,
	lbls map[geom.Pt2]label,
) {
	for x := r.Lo.X; x <= r.Hi.X; x++ {
		addLabel(g, geom.Pt2{x, r.Lo.Y}, geom.Up, false, adjs, lbls)
	}
	for x := r.Lo.X; x <= r.Hi.X; x++ {
		addLabel(g, geom.Pt2{x, r.Hi.Y}, geom.Down, false, adjs, lbls)
	}
	for y := r.Lo.Y; y <= r.Hi.Y; y++ {
		addLabel(g, geom.Pt2{r.Lo.X, y}, geom.Left, false, adjs, lbls)
	}
	for y := r.Lo.Y; y <= r.Hi.Y; y++ {
		addLabel(g, geom.Pt2{r.Hi.X, y}, geom.Right, false, adjs, lbls)
	}
}

func findInnerLabels(
	g grid,
	r geom.Rect,
	adjs map[label][]geom.Pt2,
	lbls map[geom.Pt2]label,
) {
	for x := r.Lo.X; x <= r.Hi.X; x++ {
		addLabel(g, geom.Pt2{x, r.Lo.Y - 1}, geom.Up, true, adjs, lbls)
	}
	for x := r.Lo.X; x <= r.Hi.X; x++ {
		addLabel(g, geom.Pt2{x, r.Hi.Y + 1}, geom.Down, true, adjs, lbls)
	}
	for y := r.Lo.Y; y <= r.Hi.Y; y++ {
		addLabel(g, geom.Pt2{r.Lo.X - 1, y}, geom.Left, true, adjs, lbls)
	}
	for y := r.Lo.Y; y <= r.Hi.Y; y++ {
		addLabel(g, geom.Pt2{r.Hi.X + 1, y}, geom.Right, true, adjs, lbls)
	}
}

func findLabels(
	g grid, outer, inner geom.Rect,
) (map[label][]geom.Pt2, map[geom.Pt2]label) {
	adjs := make(map[label][]geom.Pt2)
	lbls := make(map[geom.Pt2]label)
	findOuterLabels(g, outer, adjs, lbls)
	findInnerLabels(g, inner, adjs, lbls)
	return adjs, lbls
}

type maze struct {
	g     grid
	outer geom.Rect
	inner geom.Rect
	adjs  map[label][]geom.Pt2
	lbls  map[geom.Pt2]label
}

func readMaze(g grid) maze {
	outer, inner := mazeBounds(g)
	adjs, lbls := findLabels(g, outer, inner)
	return maze{g, inner, outer, adjs, lbls}
}

func (m maze) adjacent(from geom.Pt2) []geom.Pt2 {
	if m.g.g[from] == wall {
		return nil
	}
	var nbrs []geom.Pt2
	for _, nbr := range from.ManhattanNeighbors() {
		c := m.g.g[nbr]
		// don't go through walls
		if c == wall {
			continue
		}
		// go into passages
		if c == passage {
			nbrs = append(nbrs, nbr)
			continue
		}
		// go through portals
		for _, adj := range m.adjs[m.lbls[nbr]] {
			if from != adj {
				nbrs = append(nbrs, adj)
			}
		}
	}
	return nbrs
}

func printAdjacent(m maze, x, y int) {
	p := geom.Pt2{x, y}
	fmt.Println(p, m.adjacent(p))
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	g := readGrid(f)
	m := readMaze(g)
}
