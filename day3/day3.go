package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"
)

type vec struct {
	dir  rune
	dist int
}

type coord struct {
	x, y int
}

func readVecPaths(filePath string) [][]vec {
	txt, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.TrimSpace(string(txt))
	var vecPaths [][]vec
	for _, line := range strings.Split(lines, "\n") {
		var vecPath []vec
		for _, tok := range strings.Split(line, ",") {
			var v vec
			_, err := fmt.Sscanf(tok, "%c%d", &v.dir, &v.dist)
			if err != nil {
				log.Fatalf("failed to parse vec: %s", tok)
			}
			vecPath = append(vecPath, v)
		}
		vecPaths = append(vecPaths, vecPath)
	}
	return vecPaths
}

func toPath(vecPath []vec) []coord {
	var path []coord
	var cur coord
	for _, v := range vecPath {
		var dim *int
		d := 0
		switch v.dir {
		case 'U':
			dim = &cur.y
			d = +1
		case 'D':
			dim = &cur.y
			d = -1
		case 'R':
			dim = &cur.x
			d = +1
		case 'L':
			dim = &cur.x
			d = -1
		default:
			log.Fatalf("bad direction: %s", v.dir)
		}
		for n := v.dist; n > 0; n-- {
			*dim += d
			path = append(path, cur)
		}
	}
	return path
}

func findIntersects(coords1, coords2 []coord) []coord {
	coords := make(map[coord]bool)
	for _, c := range coords1 {
		coords[c] = true
	}
	var intersections []coord
	for _, c := range coords2 {
		if coords[c] {
			intersections = append(intersections, c)
		}
	}
	return intersections
}

func dist(c coord) int {
	return int(math.Abs(float64(c.x)) + math.Abs(float64(c.y)))
}

func closestIntersect(intersects []coord) int {
	var closestDist int
	for _, coord := range intersects {
		if d := dist(coord); closestDist == 0 || d < closestDist {
			closestDist = d
		}
	}
	return closestDist
}

func stepsTo(to coord, path []coord) int {
	for i, cur := range path {
		if cur == to {
			return i + 1
		}
	}
	return 0
}

func fastestIntersect(coords []coord, path1, path2 []coord) int {
	var speed int
	for _, c := range coords {
		sum := stepsTo(c, path1) + stepsTo(c, path2)
		if speed == 0 || sum < speed {
			speed = sum
		}
	}
	return speed
}

func main() {
	vecPaths := readVecPaths(os.Args[1])
	path1, path2 := toPath(vecPaths[0]), toPath(vecPaths[1])
	intersects := findIntersects(path1, path2)
	fmt.Println(closestIntersect(intersects))
	fmt.Println(fastestIntersect(intersects, path1, path2))
}
