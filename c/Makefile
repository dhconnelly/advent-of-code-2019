ifdef DEBUG
	CC_FLAGS=-g
else
	CC_FLAGS=
endif

CC=cc $(CC_FLAGS)

all: target/day2 target/day5 target/day7 target/day9 target/day11 target/test
runall: day2 day5 day7 day9 day11

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

target/test.o: test.c
	mkdir -p target
	$(CC) -c $< -o $@

target/pt2.o: pt2.c pt2.h hash.h
	mkdir -p target
	$(CC) -c $< -o $@

target/pt2_test.o: pt2_test.c pt2.h
	mkdir -p target
	$(CC) -c $< -o $@

target/hash.o: hash.c hash.h
	mkdir -p target
	$(CC) -c $< -o $@

target/hash_test.o: hash_test.c hash.h
	mkdir -p target
	$(CC) -c $< -o $@

target/test: target/test.o target/day2_test.o target/day5_test.o target/vm.o target/parse.o target/day7_test.o target/amp.o target/comb.o target/hash.o target/hash_test.o target/day9_test.o target/pt2.o target/pt2_test.o
	mkdir -p target
	$(CC) $^ -o $@

test: target/test
	./target/test

day%: target/day%
	./target/$@ ../$@/input.txt

.PHONY: clean

clean:
	rm -rf target

