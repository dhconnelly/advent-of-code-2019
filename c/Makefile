ifdef DEBUG
	CC_FLAGS=-g -pg -Werror
else
	CC_FLAGS=-Ofast -Werror
endif

CC=cc $(CC_FLAGS)

all: target/day2 target/day5 target/day7 target/day9 target/day11 target/day13 target/test target/day15 target/day23 target/day18
runall: day2 day5 day7 day9 day11 day13 day15 day23 day18

target/parse.o: parse.c parse.h
	mkdir -p target
	$(CC) -c $< -o $@

target/vm.o: vm.c vm.h
	mkdir -p target
	$(CC) -c $< -o $@

target/day2.o: day2.c vm.h parse.h hash.h
	mkdir -p target
	$(CC) -c $< -o $@

target/day2_test.o: day2_test.c vm.h hash.h
	mkdir -p target
	$(CC) -c $< -o $@

target/day2: target/day2.o target/vm.o target/parse.o target/hash.o
	mkdir -p target
	$(CC) $^ -o $@

target/day5.o: day5.c vm.h parse.h hash.h
	mkdir -p target
	$(CC) -c $< -o $@

target/day5_test.o: day5_test.c vm.h hash.h
	mkdir -p target
	$(CC) -c $< -o $@

target/day5: target/day5.o target/vm.o target/parse.o target/hash.o
	mkdir -p target
	$(CC) $^ -o $@

target/amp.o: amp.c amp.h vm.h hash.h
	mkdir -p target
	$(CC) -c $< -o $@

target/comb.o: comb.c comb.h
	mkdir -p target
	$(CC) -c $< -o $@

target/day7.o: day7.c vm.h parse.h amp.h comb.h hash.h
	mkdir -p target
	$(CC) -c $< -o $@

target/day7: target/day7.o target/vm.o target/parse.o target/amp.o target/comb.o target/hash.o
	mkdir -p target
	$(CC) $^ -o $@

target/day7_test.o: day7_test.c vm.h amp.h comb.h hash.h
	mkdir -p target
	$(CC) -c $< -o $@

target/day9.o: day9.c vm.h parse.h hash.h hash.h
	mkdir -p target
	$(CC) -c $< -o $@

target/day9: target/day9.o target/vm.o target/parse.o target/hash.o target/hash.o
	mkdir -p target
	$(CC) $^ -o $@

target/day9_test.o: day9_test.c vm.h
	mkdir -p target
	$(CC) -c $< -o $@

target/day11.o: day11.c vm.h parse.h
	mkdir -p target
	$(CC) -c $< -o $@

target/day11: target/day11.o target/vm.o target/parse.o target/hash.o target/pt2.o
	mkdir -p target
	$(CC) $^ -o $@

target/day13.o: day13.c vm.h parse.h
	mkdir -p target
	$(CC) -c $< -o $@

target/day13: target/day13.o target/vm.o target/parse.o target/hash.o
	mkdir -p target
	$(CC) $^ -o $@

target/test.o: test.c
	mkdir -p target
	$(CC) -c $< -o $@

target/pt2.o: pt2.c pt2.h hash.h
	mkdir -p target
	$(CC) -c $< -o $@

target/hash.o: hash.c hash.h
	mkdir -p target
	$(CC) -c $< -o $@

target/hash_test.o: hash_test.c hash.h
	mkdir -p target
	$(CC) -c $< -o $@

target/day15.o: day15.c hash.h pt2.h vm.h parse.h queue.h
	mkdir -p target
	$(CC) -c $< -o $@

target/day15: target/day15.o target/hash.o target/vm.o target/parse.o target/pt2.o target/queue.o
	mkdir -p target
	$(CC) $^ -o $@

target/day23.o: day23.c parse.h vm.h
	mkdir -p target
	$(CC) -c $< -o $@

target/day23: target/day23.o target/vm.o target/parse.o target/hash.o
	mkdir -p target
	$(CC) $^ -o $@

target/day25.o: day25.c parse.h vm.h
	mkdir -p target
	$(CC) -c $< -o $@

target/day25: target/day25.o target/vm.o target/parse.o target/hash.o
	mkdir -p target
	$(CC) $^ -o $@

target/day18.o: day18.c pt2.h hash.h queue.h
	mkdir -p target
	$(CC) -c $< -o $@

target/day18: target/day18.o target/pt2.o target/hash.o target/queue.o
	mkdir -p target
	$(CC) $^ -o $@

target/test: target/test.o target/day2_test.o target/day5_test.o target/vm.o target/parse.o target/day7_test.o target/amp.o target/comb.o target/hash.o target/hash_test.o target/day9_test.o target/pt2.o
	mkdir -p target
	$(CC) $^ -o $@

target/queue.o: queue.c queue.h
	mkdir -p target
	$(CC) -c $< -o $@

test: target/test
	./target/test

day%: target/day%
	./target/$@ inputs/$@.txt

.PHONY: clean

clean:
	rm -rf target

