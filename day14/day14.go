package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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
	reactions := make(map[string]reaction)
	for scan.Scan() {
		line := scan.Text()
		parts := strings.Split(line, "=>")
		out := readQuant(parts[1])
		var ins []quant
		for _, tok := range strings.Split(parts[0], ",") {
			ins = append(ins, readQuant(tok))
		}
		reactions[out.chem] = reaction{out, ins}
	}
	if err := scan.Err(); err != nil {
		log.Fatal(err)
	}
	return reactions
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	reactions := readReactions(f)
	fmt.Println(reactions)
}
