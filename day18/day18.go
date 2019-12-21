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

func (m maze) maybeFind(c rune) (geom.Pt2, bool) {
	for k, v := range m {
		if c == v {
			return k, true
		}
	}
	return geom.Zero2, false
}

func (m maze) find(c rune) geom.Pt2 {
	p, ok := m.maybeFind(c)
	if !ok {
		log.Fatalf("not found: %c", c)
		return geom.Zero2
	}
	return p
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

type bfsNode struct {
	p geom.Pt2
	c rune
	d int
}

func neighbors(m maze, p geom.Pt2) map[rune]bfsNode {
	nbrs := make(map[rune]bfsNode)
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
			c := m[adj]
			nd := bfsNode{adj, c, first.d + 1}
			if isDoor(c) || isKey(c) {
				nbrs[c] = nd
			} else {
				q = append(q, nd)
			}
		}
	}
	return nbrs
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

func copiedPts(ps []geom.Pt2) []geom.Pt2 {
	qs := make([]geom.Pt2, len(ps))
	copy(qs, ps)
	return qs
}

func remove(rs []rune, i int) []rune {
	return append(copied(rs[:i]), rs[i+1:]...)
}

func replace(ps []geom.Pt2, i int, p geom.Pt2) []geom.Pt2 {
	ps = copiedPts(ps)
	ps[i] = p
	return ps
}

var memo = make(map[string]int)

func key(pos []geom.Pt2, need []rune) string {
	return fmt.Sprintf("%v-%s", pos, string(need))
}

func multiShortestPath(pos []geom.Pt2, m maze, need []rune) int {
	mk := key(pos, need)
	if d, ok := memo[mk]; ok {
		return d
	}
	if len(need) == 0 {
		return 0
	}
	shortest := 0
	for j, from := range pos {
		nbrs := neighbors(m, from)
		for i, key := range need {
			nd, ok := nbrs[key]
			if !ok {
				continue
			}
			m2 := m.clone()
			m2[nd.p] = open
			if door, ok := m.maybeFind(doorFor(nd.c)); ok {
				m2[door] = open
			}
			subSteps := multiShortestPath(replace(pos, j, nd.p), m2, remove(need, i)) + nd.d
			if subSteps >= 0 && (shortest == 0 || subSteps < shortest) {
				shortest = subSteps
			}
		}
	}
	if shortest > 0 {
		memo[mk] = shortest
		return shortest
	}
	return -1
}

func splitMaze(m maze, atPoint geom.Pt2) []geom.Pt2 {
	m[atPoint] = wall
	m[atPoint.Go(geom.Left)] = wall
	m[atPoint.Go(geom.Right)] = wall
	m[atPoint.Go(geom.Up)] = wall
	m[atPoint.Go(geom.Down)] = wall
	p1 := atPoint.Go(geom.Up).Go(geom.Left)
	p2 := atPoint.Go(geom.Up).Go(geom.Right)
	p3 := atPoint.Go(geom.Down).Go(geom.Left)
	p4 := atPoint.Go(geom.Down).Go(geom.Right)
	m[p1] = '1'
	m[p2] = '2'
	m[p3] = '3'
	m[p4] = '4'
	return []geom.Pt2{p1, p2, p3, p4}
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	m := readMaze(f)

	steps1 := multiShortestPath([]geom.Pt2{m.find(entr)}, m, m.keys())
	fmt.Println(steps1)

	pos := splitMaze(m, m.find(entr))
	steps2 := multiShortestPath(pos, m, m.keys())
	fmt.Println(steps2)
}
