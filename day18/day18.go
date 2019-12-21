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
	entr = '@'
	wall = '#'
	open = '.'
)

func doorFor(key rune) rune {
	return key - 32
}

func keyFor(door rune) rune {
	return door + 32
}

func isDoor(c rune) bool {
	return 'A' <= c && c <= 'Z'
}

func isKey(c rune) bool {
	return 'a' <= c && c <= 'z'
}

type maze map[geom.Pt2]rune

func (m maze) clone() maze {
	n := maze(make(map[geom.Pt2]rune))
	for k, v := range m {
		n[k] = v
	}
	return n
}

func (m maze) find(c rune) geom.Pt2 {
	for k, v := range m {
		if c == v {
			return k
		}
	}
	log.Fatal("not found:", c)
	return geom.Zero2
}

func (m maze) keys() []rune {
	var keys []rune
	for _, v := range m {
		if isKey(v) {
			keys = append(keys, v)
		}
	}
	return keys
}

func (m maze) adjacent(p geom.Pt2) []geom.Pt2 {
	var adj []geom.Pt2
	for _, q := range p.ManhattanNeighbors() {
		if c, ok := m[q]; ok && c != wall {
			adj = append(adj, q)
		}
	}
	return adj
}

type edge struct {
	c rune
	d int
}

type bfsNode struct {
	p geom.Pt2
	c rune
	d int
}

func neighbors(m maze, p geom.Pt2) []edge {
	var nbrs []edge
	var first bfsNode
	q := []bfsNode{{p, m[p], 0}}
	v := map[geom.Pt2]bool{p: true}
	for len(q) > 0 {
		first, q = q[0], q[1:]
		for _, adj := range m.adjacent(first.p) {
			if v[adj] {
				continue
			}
			v[adj] = true
			if c := m[adj]; c == entr || isDoor(c) || isKey(c) {
				nbrs = append(nbrs, edge{c, first.d + 1})
			} else {
				q = append(q, bfsNode{adj, c, first.d + 1})
			}
		}
	}
	return nbrs
}

type graph map[rune]map[rune]int

func (g graph) clone() graph {
	g2 := graph(make(map[rune]map[rune]int))
	for k, v := range g {
		g2[k] = make(map[rune]int)
		for k2, v2 := range v {
			g2[k][k2] = v2
		}
	}
	return g2
}

func (g graph) keys() []rune {
	var keys []rune
	for nbr := range g {
		if isKey(nbr) {
			keys = append(keys, nbr)
		}
	}
	return keys
}

func (g graph) connected(c1, c2 rune) bool {
	_, ok := g[c1][c2]
	return ok
}

func (g graph) takeKey(c rune) map[rune]int {
	g.remove(doorFor(c))
	return g.remove(c)
}

func (g graph) remove(c rune) map[rune]int {
	nbrs := g[c]
	// connect neighbors
	for nbr1, d1 := range nbrs {
		for nbr2, d2 := range nbrs {
			// all distinct neighbor pairs
			if nbr1 == nbr2 {
				continue
			}
			// that are unconnected or are now closer with c removed
			d := d1 + d2
			if g.connected(nbr1, nbr2) && d > g[nbr1][nbr2] {
				continue
			}
			g[nbr1][nbr2] = d
			g[nbr2][nbr1] = d
		}
	}
	// remove c from the graph
	for nbr := range nbrs {
		delete(g[nbr], c)
	}
	delete(g, c)
	return nbrs
}

func printGraph(g graph) {
	for c, nbrs := range g {
		fmt.Printf("%c: ", c)
		for nbr, d := range nbrs {
			fmt.Printf("[%c, dist=%d] ", nbr, d)
		}
		fmt.Println()
	}
}

func reachableGraph(m maze) graph {
	g := graph(make(map[rune]map[rune]int))
	for p, c := range m {
		if c == entr || isDoor(c) || isKey(c) {
			g[c] = make(map[rune]int)
			for _, e := range neighbors(m, p) {
				g[c][e.c] = e.d
			}
		}
	}
	return g
}

func readMaze(r io.Reader) maze {
	scan := bufio.NewScanner(r)
	m := maze(make(map[geom.Pt2]rune))
	y := 0
	for scan.Scan() {
		x := 0
		for _, c := range scan.Text() {
			m[geom.Pt2{x, y}] = c
			x++
		}
		y++
	}
	if err := scan.Err(); err != nil {
		log.Fatal(err)
	}
	return m
}

func copied(x []rune) []rune {
	y := make([]rune, len(x))
	for i, v := range x {
		y[i] = v
	}
	return y
}

func remove(rs []rune, i int) []rune {
	return append(copied(rs[:i]), rs[i+1:]...)
}

type memoKey struct {
	c rune
	s string
}

var memo = make(map[memoKey]int)

func shortestPathTaking(c rune, g graph, need []rune) int {
	if d, ok := memo[memoKey{c, string(need)}]; ok {
		return d
	}
	if len(need) == 0 {
		return 0
	}

	nbrs := g.takeKey(c)
	shortest := 0
	for i, next := range need {
		d, ok := nbrs[next]
		if !ok {
			continue
		}
		subSteps := shortestPathTaking(next, g.clone(), remove(need, i)) + d
		if shortest == 0 || subSteps < shortest {
			shortest = subSteps
		}
	}
	memo[memoKey{c, string(need)}] = shortest
	return shortest
}

func shortestKeyPath(g graph) int {
	return shortestPathTaking(entr, g, g.keys())
}

func splitMaze(m maze, atPoint geom.Pt2) {
	m[atPoint] = wall
	m[atPoint.Go(geom.Left)] = wall
	m[atPoint.Go(geom.Right)] = wall
	m[atPoint.Go(geom.Up)] = wall
	m[atPoint.Go(geom.Down)] = wall
	m[atPoint.Go(geom.Up).Go(geom.Left)] = entr
	m[atPoint.Go(geom.Up).Go(geom.Right)] = entr
	m[atPoint.Go(geom.Down).Go(geom.Left)] = entr
	m[atPoint.Go(geom.Down).Go(geom.Right)] = entr
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	m := readMaze(f)
	g := reachableGraph(m)
	steps := shortestKeyPath(g)
	fmt.Println(steps)

	splitMaze(m, m.find(entr))
	g = reachableGraph(m)
	printGraph(g)
}
