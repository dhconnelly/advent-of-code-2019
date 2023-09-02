#include <assert.h>
#include <stdio.h>
#include <stdlib.h>

#include "hash.h"
#include "parse.h"
#include "vm.h"

int64_t execute(vm* vm, int64_t noun, int64_t verb) {
    set_mem(vm, 1, noun);
    set_mem(vm, 2, verb);
    run(vm);
    return get_mem(vm, 0);
}

int64_t part1(vm* vm) { return execute(vm, 12, 2); }

int64_t part2(vm* base) {
    for (int noun = 0; noun <= 99; noun++) {
        for (int verb = 0; verb <= 99; verb++) {
            vm local = copy_vm(base);
            if (execute(&local, noun, verb) == 19690720) {
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

    int64_t data[1024];
    int len;
    if ((len = parse_intcode(f, data, 1024)) < 0) {
        perror("day2");
        exit(1);
    }

    vm vm = new_vm();
    fill_table(vm.mem, data, len);
    printf("%lld\n", part1(&vm));

    vm = new_vm();
    fill_table(vm.mem, data, len);
    printf("%lld\n", part2(&vm));

    return 0;
}
