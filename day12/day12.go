package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/dhconnelly/advent-of-code-2019/geom"
)

type moon struct {
	p, v geom.Pt3
}

func readPoints(path string) []moon {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	var ms []moon
	var m moon
	for scan.Scan() {
		line := scan.Text()
		_, err = fmt.Sscanf(
			line, "<x=%d, y=%d, z=%d>",
			&m.p.X, &m.p.Y, &m.p.Z,
		)
		if err != nil {
			log.Fatalf("bad point: %s", err)
		}
		ms = append(ms, m)
	}
	if err = scan.Err(); err != nil {
		log.Fatal(err)
	}
	return ms
}

func applyGravityPair(m1, m2 *moon) {
	if m1.p.X < m2.p.X {
		m1.v.X, m2.v.X = m1.v.X+1, m2.v.X-1
	} else if m1.p.X > m2.p.X {
		m1.v.X, m2.v.X = m1.v.X-1, m2.v.X+1
	}

	if m1.p.Y < m2.p.Y {
		m1.v.Y, m2.v.Y = m1.v.Y+1, m2.v.Y-1
	} else if m1.p.Y > m2.p.Y {
		m1.v.Y, m2.v.Y = m1.v.Y-1, m2.v.Y+1
	}

	if m1.p.Z < m2.p.Z {
		m1.v.Z, m2.v.Z = m1.v.Z+1, m2.v.Z-1
	} else if m1.p.Z > m2.p.Z {
		m1.v.Z, m2.v.Z = m1.v.Z-1, m2.v.Z+1
	}
}

func applyGravity(ms []moon) {
	for i := 0; i < len(ms)-1; i++ {
		for j := i + 1; j < len(ms); j++ {
			applyGravityPair(&ms[i], &ms[j])
		}
	}
}

func applyVelocity(ms []moon) {
	for i := range ms {
		ms[i].p.TranslateBy(ms[i].v)
	}
}

func simulate(ms []moon, n int) {
	for i := 0; i < n; i++ {
		applyGravity(ms)
		applyVelocity(ms)
	}
}

func moonEnergy(m moon) int {
	return m.p.ManhattanNorm() * m.v.ManhattanNorm()
}

func energy(ms []moon) int {
	total := 0
	for _, m := range ms {
		total += moonEnergy(m)
	}
	return total
}

func main() {
	ms := readPoints(os.Args[1])
	simulate(ms, 1000)
	fmt.Println(energy(ms))
}
