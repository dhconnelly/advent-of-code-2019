package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sort"

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

func allStepsFrom(g grid, from geom.Pt2, dx, dy int) []geom.Pt2 {
	var steps []geom.Pt2
	reachable := make(map[geom.Pt2]bool)
	for i := 0; i < g.height; i++ {
		for j := 0; j < g.width; j++ {
			add := false
			d := geom.Pt2{dx * j, dy * i}
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
	return steps
}

func allSteps(g grid) []geom.Pt2 {
	var steps []geom.Pt2
	steps = append(steps, allStepsFrom(g, geom.Pt2{0, g.height}, 1, -1)...)
	steps = append(steps, allStepsFrom(g, geom.Pt2{0, 0}, 1, 1)...)
	steps = append(steps, allStepsFrom(g, geom.Pt2{g.width, g.height}, -1, -1)...)
	steps = append(steps, allStepsFrom(g, geom.Pt2{g.width, 0}, -1, 1)...)
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

func countVisible(g grid, steps []geom.Pt2) map[geom.Pt2]int {
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

type byAngleFrom struct {
	p  geom.Pt2
	ps []geom.Pt2
}

func (points byAngleFrom) Len() int {
	return len(points.ps)
}

func angle(p1, p2 geom.Pt2) float64 {
	return math.Atan2(float64(p2.Y-p1.Y), float64(p2.X-p1.X))
}

func (points byAngleFrom) Less(i, j int) bool {
	to1, to2 := points.ps[i], points.ps[j]
	a1 := angle(points.p, to1)
	a2 := angle(points.p, to2)
	return a1 <= a2
}

func (points byAngleFrom) Swap(i, j int) {
	points.ps[j], points.ps[i] = points.ps[i], points.ps[j]
}

func reachableFrom(g grid, from geom.Pt2, steps []geom.Pt2) []geom.Pt2 {
	var to []geom.Pt2
	for _, d := range steps {
		var p geom.Pt2
		for p = from.Add(d); inBounds(g, p) && !g.points[p]; p = p.Add(d) {
		}
		if g.points[p] {
			to = append(to, p)
		}
	}
	sort.Sort(byAngleFrom{from, to})
	var j int
	for j = 0; j < len(to) && angle(from, to[j]) < -math.Pi/2.0; j++ {
	}
	ordered := make([]geom.Pt2, len(to))
	for i := 0; i < len(to); i++ {
		ix := (j + i) % len(to)
		ordered[i] = to[ix]
	}
	return ordered
}

func vaporizeAll(g grid, from geom.Pt2, steps []geom.Pt2) []geom.Pt2 {
	ordered := make([]geom.Pt2, len(g.points))
	for vaporized := 0; vaporized < len(g.points)-1; {
		toVaporize := reachableFrom(g, from, steps)
		for _, p := range toVaporize {
			g.points[p] = false
			ordered[vaporized] = p
			vaporized++
		}
	}
	return ordered
}

func main() {
	g := readGrid(os.Args[1])
	steps := allSteps(g)
	counts := countVisible(g, steps)
	best, count := bestPoint(g, counts)
	fmt.Println(count)
	vaporized := vaporizeAll(g, best, steps)
	winPt := vaporized[199]
	fmt.Println(winPt.X*100 + winPt.Y)
}
