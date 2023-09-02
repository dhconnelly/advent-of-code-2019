CC=cc

all: target/day2 target/day5 target/day7

target/parse.o: parse.c parse.h
	mkdir -p target
	$(CC) -c $< -o $@

target/vm.o: vm.c vm.h
	mkdir -p target
	$(CC) -c $< -o $@

target/day2.o: day2.c vm.h parse.h
	mkdir -p target
	$(CC) -c $< -o $@

target/day2_test.o: day2_test.c vm.h
	mkdir -p target
	$(CC) -c $< -o $@

target/day2: target/day2.o target/vm.o target/parse.o
	mkdir -p target
	$(CC) $^ -o $@

target/day5.o: day5.c vm.h parse.h
	mkdir -p target
	$(CC) -c $< -o $@

target/day5_test.o: day5_test.c vm.h
	mkdir -p target
	$(CC) -c $< -o $@

target/day5: target/day5.o target/vm.o target/parse.o
	mkdir -p target
	$(CC) $^ -o $@

target/amp.o: amp.c amp.h vm.h
	mkdir -p target
	$(CC) -c $< -o $@

target/comb.o: comb.c comb.h
	mkdir -p target
	$(CC) -c $< -o $@

target/day7.o: day7.c vm.h parse.h amp.h comb.h
	mkdir -p target
	$(CC) -c $< -o $@

target/day7: target/day7.o target/vm.o target/parse.o target/amp.o target/comb.o
	mkdir -p target
	$(CC) $^ -o $@

target/day7_test.o: day7_test.c vm.h amp.h comb.h
	mkdir -p target
	$(CC) -c $< -o $@

target/test.o: test.c
	mkdir -p target
	$(CC) -c $< -o $@

target/hash.o: hash.c hash.h
	mkdir -p target
	$(CC) -c $< -o $@

target/hash_test.o: hash_test.c hash.h
	mkdir -p target
	$(CC) -c $< -o $@

target/test: target/test.o target/day2_test.o target/day5_test.o target/vm.o target/parse.o target/day7_test.o target/amp.o target/comb.o target/hash.o target/hash_test.o
	mkdir -p target
	$(CC) $^ -o $@

test: target/test
	./target/test

.PHONY: clean

clean:
	rm -rf target

