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
	"github.com/dhconnelly/advent-of-code-2019/ints"
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

func angle(dy, dx int) float64 {
	g := ints.Abs(ints.Gcd(dx, dy))
	if g == 0 {
		return math.NaN()
	}
	return math.Atan2(float64(dy/g), float64(dx/g))
}

func reachable(g grid, from geom.Pt2) map[float64]geom.Pt2 {
	ps := make(map[float64]geom.Pt2)
	for p1, ok1 := range g.points {
		if !ok1 {
			continue
		}
		if p1 == from {
			continue
		}
		s := angle(p1.Y-from.Y, p1.X-from.X)
		if p2, ok2 := ps[s]; !ok2 || from.Dist(p1) < from.Dist(p2) {
			ps[s] = p1
		}
	}
	return ps
}

func maxReachable(g grid) (geom.Pt2, int) {
	maxPt, max := geom.Zero2, 0
	for p, ok := range g.points {
		if !ok {
			continue
		}
		to := reachable(g, p)
		if l := len(to); l > max {
			max, maxPt = l, p
		}
	}
	return maxPt, max
}

func sortedByAngle(ps map[float64]geom.Pt2) []geom.Pt2 {
	sorted := make([]float64, 0, len(ps))
	for a := range ps {
		sorted = append(sorted, a)
	}
	sort.Float64s(sorted)
	var j int
	for j = 0; j < len(sorted) && sorted[j] < -math.Pi/2.0; j++ {
	}
	byAngle := make([]geom.Pt2, len(ps))
	for i := 0; i < len(sorted); i++ {
		a := sorted[(i+j)%len(sorted)]
		byAngle[i] = ps[a]
	}
	return byAngle
}

func vaporize(g grid, from geom.Pt2) []geom.Pt2 {
	vaporized := make([]geom.Pt2, len(g.points)-1)
	for i := 0; i < len(g.points)-1; {
		ps := reachable(g, from)
		sorted := sortedByAngle(ps)
		for _, p := range sorted {
			vaporized[i] = p
			g.points[p] = false
			i++
		}
	}
	return vaporized
}

func main() {
	g := readGrid(os.Args[1])
	pt, max := maxReachable(g)
	winner := vaporize(g, pt)[199]

	fmt.Println(max)
	fmt.Println(winner.X*100 + winner.Y)
}
