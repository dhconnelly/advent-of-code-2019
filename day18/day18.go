package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"

	"github.com/dhconnelly/advent-of-code-2019/geom"
	"github.com/dhconnelly/advent-of-code-2019/ints"
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

func shouldVisitNeighbor(nbr rune, nbrs map[rune]int, after map[rune]map[rune]bool) bool {
	if !isKey(nbr) {
		return false
	}
	for nbr2 := range nbrs {
		if after[nbr2][nbr] {
			return false
		}
	}
	return true
}

type byAfter struct {
	nodes []bfsNode
	after map[rune]map[rune]bool
}

func (ba *byAfter) Len() int {
	return len(ba.nodes)
}

func (ba *byAfter) Swap(i, j int) {
	ba.nodes[i], ba.nodes[j] = ba.nodes[j], ba.nodes[i]
}

func (ba *byAfter) Less(i, j int) bool {
	return ba.after[ba.nodes[i].c][ba.nodes[j].c]
}

func prioritySort(nbrs map[rune]int, after map[rune]map[rune]bool) []bfsNode {
	nodes := byAfter{after: after}
	for nbr, d := range nbrs {
		if isKey(nbr) {
			nodes.nodes = append(nodes.nodes, bfsNode{c: nbr, d: d})
		}
	}
	sort.Sort(&nodes)
	return nodes.nodes
}

func shortestPathTaking(c rune, g graph, after map[rune]map[rune]bool, from []rune, remainingKeys, fromSteps, limit int) ([]rune, int) {
	//fmt.Printf("taking %c (steps=%d, limit = %d, remaining = %d)\n", c, fromSteps, limit, remainingKeys)
	if fromSteps >= limit {
		//fmt.Println("over limit!")
		return nil, 0
	}
	if remainingKeys == 0 {
		fmt.Printf("done! length %d, path %v\n", fromSteps, from)
		return from, fromSteps
	}

	nbrs := g.takeKey(c)
	//fmt.Printf("neighbors of %c: %v\n", c, nbrs)

	var shortestPath []rune
	shortest := 0
	for _, nd := range prioritySort(nbrs, after) {
		if fromSteps+nd.d >= limit {
			continue
		}
		//fmt.Printf("branch: taking %c\n", nbr)
		path, steps := shortestPathTaking(nd.c, g.clone(), after, append(copied(from), nd.c), remainingKeys-1, fromSteps+nd.d, limit)
		//fmt.Printf("branch for %c done. steps = %d\n", nbr, steps)
		if len(path) > 0 && (shortest == 0 || steps < shortest) {
			shortest, shortestPath = steps, path
			limit = ints.Min(shortest, limit)
		}
	}
	return shortestPath, shortest
}

func shortestKeyPath(g graph, after map[rune]map[rune]bool) ([]rune, int) {
	return shortestPathTaking(entr, g, after, nil, len(g.keys()), 0, math.MaxInt64)
}

type pathNode struct {
	c    rune
	path []rune
}

func shortestPath(from, to rune, g graph) []rune {
	var first pathNode
	q := []pathNode{{from, nil}}
	v := map[rune]bool{from: true}
	for len(q) > 0 {
		first, q = q[0], q[1:]
		for nbr := range g[first.c] {
			if v[nbr] {
				continue
			}
			v[nbr] = true
			path := append(copied(first.path), nbr)
			if nbr == to {
				return path
			}
			q = append(q, pathNode{nbr, path})
		}
	}
	return nil
}

func shortestPaths(to rune, g graph) map[rune][]rune {
	paths := make(map[rune][]rune)
	for c := range g {
		paths[c] = shortestPath(c, to, g)
	}
	return paths
}

func printPaths(p map[rune][]rune) {
	for c, path := range p {
		fmt.Printf("[%c] ", c)
		for _, n := range path {
			fmt.Printf("%c ", n)
		}
		fmt.Println()
	}
}

func doorPaths(p map[rune][]rune) map[rune][]rune {
	dp := make(map[rune][]rune)
	for c, path := range p {
		for _, n := range path {
			if isDoor(n) {
				dp[c] = append(dp[c], n)
			}
		}
	}
	return dp
}

func dependents(paths map[rune][]rune) map[rune]map[rune]bool {
	after := make(map[rune]map[rune]bool)
	for c, path := range paths {
		if !isDoor(c) {
			continue
		}
		k := keyFor(c)
		after[k] = make(map[rune]bool)
		for _, next := range path {
			if next == entr {
				break
			}
			if isKey(next) {
				after[k][next] = true
			}
		}
	}
	return after
}

func requirements(g graph) map[rune][]rune {
	paths := shortestPaths(entr, g)
	reqs := make(map[rune][]rune)
	for _, key := range g.keys() {
		path := paths[key]
		for _, c := range path {
			if isDoor(c) {
				reqs[key] = append(reqs[key], keyFor(c))
			}
		}
	}
	return reqs
}

func printPath(path []rune) {
	fmt.Println(string(path))
}

func avail(reqs []rune, keys map[rune]bool) bool {
	for _, req := range reqs {
		if keys[req] {
			return false
		}
	}
	return true
}

func orders(keys map[rune]bool, reqs map[rune][]rune, path []rune) [][]rune {
	if len(keys) == len(path) {
		return [][]rune{path}
	}
	var all [][]rune
	for key, ok := range keys {
		if !ok || !avail(reqs[key], keys) {
			continue
		}
		keys[key] = false
		all = append(all, orders(keys, reqs, append(copied(path), key))...)
		keys[key] = true
	}
	return all
}

func keySet(keys []rune) map[rune]bool {
	set := make(map[rune]bool)
	for _, key := range keys {
		set[key] = true
	}
	return set
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	m := readMaze(f)
	g := reachableGraph(m)
	printGraph(g)

	p := shortestPaths(entr, g)
	printPaths(p)

	after := dependents(p)
	fmt.Println(after)

	reqs := requirements(g)
	fmt.Println(reqs)
	fmt.Println(orders(keySet(g.keys()), reqs, nil))

	//order := flatten(g, after)
	//printPath(order)

	//path, steps := shortestKeyPath(g, after)
	//printPath(path)
	//fmt.Println(steps)
}
