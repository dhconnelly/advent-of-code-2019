package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/dhconnelly/advent-of-code-2019/geom"
)

type bitset int64

func (b bitset) get(i int) bool {
	return (b & (1 << i)) > 0
}

func (b *bitset) set(i int, value bool) {
	if value {
		(*b) |= 1 << i
	} else {
		(*b) &= ^(1 << i)
	}
}

type layout struct {
	bits   bitset
	width  int
	height int
}

func (l *layout) alive(row, col int) bool {
	return l.bits.get(row*l.width + col)
}

func (l *layout) adjacent(row, col int) int {
	adj := 0
	if row > 0 && l.alive(row-1, col) {
		adj++
	}
	if row < l.height-1 && l.alive(row+1, col) {
		adj++
	}
	if col > 0 && l.alive(row, col-1) {
		adj++
	}
	if col < l.width-1 && l.alive(row, col+1) {
		adj++
	}
	return adj
}

func (l *layout) next() {
	next := l.bits
	for row := 0; row < l.height; row++ {
		for col := 0; col < l.width; col++ {
			adj := l.adjacent(row, col)
			n := row*l.width + col
			if l.bits.get(n) && adj != 1 {
				next.set(n, false)
			} else if !l.bits.get(n) && (adj == 1 || adj == 2) {
				next.set(n, true)
			}
		}
	}
	l.bits = next
}

func (l layout) String() string {
	n := l.width * l.height
	b := make([]byte, n)
	for j := 0; j < n; j++ {
		if l.bits.get(j) {
			b[j] = '#'
		} else {
			b[j] = '.'
		}
	}
	return string(b)
}

func readLayout(r io.Reader) layout {
	b := make([]byte, 1)
	var l layout
	i := 0
outer:
	for {
	inner:
		for {
			_, err := r.Read(b)
			if err == io.EOF {
				break outer
			}
			if err != nil {
				log.Fatal(err)
			}
			switch b[0] {
			case '\n':
				break inner
			case '#':
				l.bits.set(i, true)
				i++
			case '.':
				l.bits.set(i, false)
				i++
			default:
				log.Fatalf("bad char in layout: %c", b[0])
			}
		}
		l.height++
	}
	l.width = i / l.height
	return l
}

func findRepeat(l layout) bitset {
	m := map[bitset]bool{l.bits: true}
	for {
		l.next()
		if _, ok := m[l.bits]; ok {
			return l.bits
		}
		m[l.bits] = true
	}
}

type tile struct {
	p     geom.Pt2
	depth int
}

type grid struct {
	width, height int
	g             map[tile]bool
}

func toGrid(l layout) grid {
	g := make(map[tile]bool)
	for row := 0; row < l.height; row++ {
		for col := 0; col < l.width; col++ {
			t := tile{p: geom.Pt2{col, row}, depth: 0}
			g[t] = l.alive(row, col)
		}
	}
	delete(g, tile{p: geom.Pt2{2, 2}, depth: 0})
	return grid{width: l.width, height: l.height, g: g}
}

func (g grid) adjacent(t tile) []tile {
	var adj []tile

	// left
	if t.p.X > 0 && (t.p.X != 3 || t.p.Y != 2) {
		q := t.p.Go(geom.Left)
		adj = append(adj, tile{p: q, depth: t.depth})
	} else if t.p.X == 0 {
		q := geom.Pt2{1, 2}
		adj = append(adj, tile{p: q, depth: t.depth + 1})
	} else if t.p.X == 3 && t.p.Y == 2 {
		for y := 0; y < g.height; y++ {
			q := geom.Pt2{g.width - 1, y}
			adj = append(adj, tile{p: q, depth: t.depth - 1})
		}
	}

	// right
	if t.p.X < g.width-1 && (t.p.X != 1 || t.p.Y != 2) {
		q := t.p.Go(geom.Right)
		adj = append(adj, tile{p: q, depth: t.depth})
	} else if t.p.X == g.width-1 {
		q := geom.Pt2{3, 2}
		adj = append(adj, tile{p: q, depth: t.depth + 1})
	} else if t.p.X == 1 && t.p.Y == 2 {
		for y := 0; y < g.height; y++ {
			q := geom.Pt2{0, y}
			adj = append(adj, tile{p: q, depth: t.depth - 1})
		}
	}

	// down
	if t.p.Y > 0 && (t.p.X != 2 || t.p.Y != 3) {
		q := t.p.Go(geom.Down)
		adj = append(adj, tile{p: q, depth: t.depth})
	} else if t.p.Y == 0 {
		q := geom.Pt2{2, 1}
		adj = append(adj, tile{p: q, depth: t.depth + 1})
	} else if t.p.X == 2 && t.p.Y == 3 {
		for x := 0; x < g.width; x++ {
			q := geom.Pt2{x, g.height - 1}
			adj = append(adj, tile{p: q, depth: t.depth - 1})
		}
	}

	// up
	if t.p.Y < g.height-1 && (t.p.X != 2 || t.p.Y != 1) {
		q := t.p.Go(geom.Up)
		adj = append(adj, tile{p: q, depth: t.depth})
	} else if t.p.Y == g.height-1 {
		q := geom.Pt2{2, 3}
		adj = append(adj, tile{p: q, depth: t.depth + 1})
	} else if t.p.X == 2 && t.p.Y == 1 {
		for x := 0; x < g.width; x++ {
			q := geom.Pt2{x, 0}
			adj = append(adj, tile{p: q, depth: t.depth - 1})
		}
	}

	return adj
}

func countBugs(g grid, n int) grid {
	return g
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	l := readLayout(f)
	fmt.Println(findRepeat(l))

	g := toGrid(l)
	fmt.Println(g.adjacent(tile{p: geom.Pt2{3, 3}, depth: 1}))
	fmt.Println(g.adjacent(tile{p: geom.Pt2{1, 1}, depth: 0}))
	fmt.Println(g.adjacent(tile{p: geom.Pt2{3, 0}, depth: 0}))
	fmt.Println(g.adjacent(tile{p: geom.Pt2{4, 0}, depth: 0}))
	fmt.Println(g.adjacent(tile{p: geom.Pt2{3, 2}, depth: 1}))
	fmt.Println(g.adjacent(tile{p: geom.Pt2{3, 2}, depth: 0}))
	//fmt.Println(countBugs(g, 200))
}
