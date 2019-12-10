package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/dhconnelly/advent-of-code-2019/geom"
)

type grid struct {
	width, height int
	points        map[geom.Pt2]bool
}

func readGrid(path string) grid {
	txt, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	txt = bytes.TrimSpace(txt)
	lines := bytes.Split(txt, []byte("\n"))
	g := grid{len(lines[0]), len(lines), make(map[geom.Pt2]bool)}
	for i, line := range lines {
		for j, ch := range line {
			if ch == '#' {
				g.points[geom.Pt2{j, i}] = true
			}
		}
	}
	return g
}

func inBounds(g grid, p geom.Pt2) bool {
	return p.X < g.width && p.Y < g.height && p.X >= 0 && p.Y >= 0
}

func allStepsFrom(g grid, from geom.Pt2) []geom.Pt2 {
	var steps []geom.Pt2
	reachable := make(map[geom.Pt2]bool)
	for i := 0; i < g.height; i++ {
		for j := 0; j < g.width; j++ {
			for _, di := range []int{-1, 1} {
				for _, dj := range []int{-1, 1} {
					add := false
					d := geom.Pt2{dj * j, di * i}
					if d == geom.Zero2 {
						continue
					}
					for p := from.Add(d); inBounds(g, p); p = p.Add(d) {
						if !reachable[p] {
							add = true
							reachable[p] = true
						}
					}
					if add {
						steps = append(steps, d)
					}
				}
			}
		}
	}
	return steps
}

func allSteps(g grid) []geom.Pt2 {
	steps := allStepsFrom(g, geom.Zero2)
	steps = append(steps, allStepsFrom(g, geom.Pt2{0, g.height})...)
	steps = append(steps, allStepsFrom(g, geom.Pt2{g.width, g.height})...)
	steps = append(steps, allStepsFrom(g, geom.Pt2{g.width, 0})...)
	return steps
}

func visit(visible map[geom.Pt2]int, g grid, p geom.Pt2, steps []geom.Pt2) {
	for _, d := range steps {
		var cur geom.Pt2
		for cur = p.Add(d); inBounds(g, cur) && !g.points[cur]; cur = cur.Add(d) {
		}
		if g.points[cur] {
			visible[p]++
		}
	}
}

func countVisible(g grid) map[geom.Pt2]int {
	steps := allSteps(g)
	counts := make(map[geom.Pt2]int)
	for p, _ := range g.points {
		visit(counts, g, p, steps)
	}
	return counts
}

func bestPoint(g grid, counts map[geom.Pt2]int) (geom.Pt2, int) {
	var best geom.Pt2
	count := 0
	for p, _ := range g.points {
		if c := counts[p]; c > count {
			best = p
			count = c
		}
	}
	return best, count
}

func printCounts(g grid, counts map[geom.Pt2]int) {
	for i := 0; i < g.height; i++ {
		for j := 0; j < g.width; j++ {
			if c := counts[geom.Pt2{j, i}]; c > 0 {
				fmt.Printf("%3d ", c)
			} else {
				fmt.Print("... ")
			}
		}
		fmt.Println()
	}
}

func main() {
	g := readGrid(os.Args[1])
	counts := countVisible(g)
	fmt.Println(counts)
	printCounts(g, counts)
	fmt.Println(bestPoint(g, counts))
}
