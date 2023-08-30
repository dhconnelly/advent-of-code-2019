#include <assert.h>
#include <stdio.h>
#include <stdlib.h>

#include "parse.h"
#include "vm.h"

int64_t execute(vm vm, int64_t noun, int64_t verb) {
    vm.mem[1] = noun;
    vm.mem[2] = verb;
    run(&vm);
    return vm.mem[0];
}

int64_t part1(vm vm) { return execute(vm, 12, 2); }

int64_t part2(vm vm) {
    for (int noun = 0; noun <= 99; noun++) {
        for (int verb = 0; verb <= 99; verb++) {
            if (execute(vm, noun, verb) == 19690720) {
                return 100 * noun + verb;
            }
        }
    }
    assert(0);
}

int main(int argc, char* argv[]) {
    if (argc != 2) {
        printf("usage: day2 file\n");
        exit(1);
    }

    FILE* f = fopen(argv[1], "r");
    if (f == NULL) {
        perror("day2");
        exit(1);
    }

    vm vm = new_vm();
    if ((vm.mem_size = parse_intcode(f, vm.mem, MAX_MEM)) < 0) {
        perror("day2");
        exit(1);
    }

    printf("%lld\n", part1(vm));
    printf("%lld\n", part2(vm));

    return 0;
}
