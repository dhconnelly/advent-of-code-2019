package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
)

type movement struct {
	dir  rune
	dist int
}

type path struct {
	movs []movement
}

func readPaths(filePath string) []path {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var paths []path
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		var p path
		r := bufio.NewReader(strings.NewReader(scan.Text()))
		var done bool
		for !done {
			tok, err := r.ReadString(',')
			if err == io.EOF {
				done = true
			} else if err != nil {
				log.Fatal(err)
			}
			var mov movement
			_, err = fmt.Sscanf(tok, "%c%d", &mov.dir, &mov.dist)
			if err != nil {
				log.Fatalf("failed to parse movement: %s", tok)
			}
			p.movs = append(p.movs, mov)
		}
		paths = append(paths, p)
	}
	if err := scan.Err(); err != nil {
		log.Fatal(err)
	}
	return paths
}

type coord struct {
	x, y int
}

func toCoords(p path) []coord {
	var coords []coord
	var cur coord
	for _, mov := range p.movs {
		var dim *int
		d, n := 0, 0
		switch mov.dir {
		case 'U':
			dim = &cur.y
			d, n = +1, mov.dist
		case 'D':
			dim = &cur.y
			d, n = -1, mov.dist
		case 'R':
			dim = &cur.x
			d, n = +1, mov.dist
		case 'L':
			dim = &cur.x
			d, n = -1, mov.dist
		default:
			log.Fatalf("bad direction: %s", mov.dir)
		}
		for ; n > 0; n-- {
			*dim += d
			coords = append(coords, cur)
		}
	}
	return coords
}

type coordSet = []coord

func allToCoords(paths []path) []coordSet {
	var s []coordSet
	for _, p := range paths {
		s = append(s, toCoords(p))
	}
	return s
}

func findIntersects(coordSets []coordSet) []coord {
	var intersects []coord
	// track coord -> which coordSets are occupying it
	m := make(map[coord]map[int]bool)
	for i, coordSet := range coordSets {
		for _, coord := range coordSet {
			if m[coord] == nil {
				// not yet occupied
				m[coord] = make(map[int]bool)
				m[coord][i] = true
				continue
			}
			m[coord][i] = true
			if len(m[coord]) > 1 {
				// if more than one set occupying: intersection
				intersects = append(intersects, coord)
			}
		}
	}
	return intersects
}

func dist(c coord) int {
	return int(math.Abs(float64(c.x)) + math.Abs(float64(c.y)))
}

func closestIntersect(intersects []coord) coord {
	closest, closestDist := intersects[0], dist(intersects[0])
	for _, coord := range intersects[1:] {
		if d := dist(coord); d < closestDist {
			closest, closestDist = coord, d
		}
	}
	return closest
}

func stepsToCoord(c coord, p path) int {
	steps := 0
	for _, pc := range toCoords(p) {
		steps++
		if pc == c {
			break
		}
	}
	return steps
}

func fastestIntersect(intersects []coord, paths []path) int {
	var speed int
	for _, c := range intersects {
		sum := 0
		for _, p := range paths {
			sum += stepsToCoord(c, p)
		}
		if speed == 0 || sum < speed {
			speed = sum
		}
	}
	return speed
}

func main() {
	paths := readPaths(os.Args[1])
	coords := allToCoords(paths)
	intersects := findIntersects(coords)
	fmt.Println(dist(closestIntersect(intersects)))
	fmt.Println(fastestIntersect(intersects, paths))
}
