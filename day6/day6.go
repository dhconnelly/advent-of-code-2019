package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type orbit struct {
	orbiter, orbited string
}

func read(path string) []orbit {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	var orbits []orbit
	for scan.Scan() {
		line := scan.Text()
		objs := strings.Split(line, ")")
		orbits = append(orbits, orbit{objs[1], objs[0]})
	}
	if err = scan.Err(); err != nil {
		log.Fatal(err)
	}
	return orbits
}

func orbitMap(orbits []orbit) map[string]string {
	m := make(map[string]string)
	for _, o := range orbits {
		m[o.orbiter] = o.orbited
	}
	return m
}

func chain(k string, m map[string]string) []string {
	var chain []string
	for v, ok := m[k]; ok; v, ok = m[v] {
		chain = append(chain, v)
	}
	return chain
}

func countOrbits(orbits map[string]string) int {
	n := 0
	for k, _ := range orbits {
		n += len(chain(k, orbits))
	}
	return n
}

func closestAncestor(chain1, chain2 []string) (string, int, int) {
	m := make(map[string]int)
	for i, k := range chain1 {
		m[k] = i
	}
	for i, o := range chain2 {
		if j, ok := m[o]; ok {
			return o, i, j
		}
	}
	return "", 0, 0
}

func transfers(orbits map[string]string, from, to string) int {
	fromChain, toChain := chain(from, orbits), chain(to, orbits)
	_, dist1, dist2 := closestAncestor(fromChain, toChain)
	return dist1 + dist2
}

func main() {
	m := orbitMap(read(os.Args[1]))
	fmt.Println(countOrbits(m))
	fmt.Println(transfers(m, "YOU", "SAN"))
}
