package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/dhconnelly/advent-of-code-2019/ints"
)

type quant struct {
	amt  int
	chem string
}

type reaction struct {
	out quant
	ins []quant
}

func readQuant(s string) quant {
	var q quant
	_, err := fmt.Sscanf(s, "%d %s", &q.amt, &q.chem)
	if err != nil {
		log.Fatalf("failed to read quant %s: %s", s, err)
	}
	return q
}

func readReactions(r io.Reader) map[string]reaction {
	scan := bufio.NewScanner(r)
	reacts := make(map[string]reaction)
	for scan.Scan() {
		line := scan.Text()
		parts := strings.Split(line, "=>")
		out := readQuant(parts[1])
		var ins []quant
		for _, tok := range strings.Split(parts[0], ",") {
			ins = append(ins, readQuant(tok))
		}
		reacts[out.chem] = reaction{out, ins}
	}
	if err := scan.Err(); err != nil {
		log.Fatal(err)
	}
	return reacts
}

func divceil(m, n int) int {
	q := m / n
	if m%n == 0 {
		return q
	}
	return q + 1
}

func oreNeeded(
	chem string, amt int,
	reacts map[string]reaction,
	waste map[string]int,
) int {
	if chem == "ORE" {
		return amt
	}
	// reuse excess production before building more
	if avail := waste[chem]; avail > 0 {
		reclaimed := ints.Min(amt, avail)
		waste[chem] -= reclaimed
		amt -= reclaimed
	}
	if amt == 0 {
		return 0
	}
	// build as much as necessary and store the excess
	react := reacts[chem]
	k := 1
	if amt > react.out.amt {
		k = divceil(amt, react.out.amt)
	}
	waste[chem] += k*react.out.amt - amt
	ore := 0
	for _, in := range react.ins {
		ore += oreNeeded(in.chem, k*in.amt, reacts, waste)
	}
	return ore
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	reacts := readReactions(f)
	fmt.Println(oreNeeded("FUEL", 1, reacts, map[string]int{}))
}
