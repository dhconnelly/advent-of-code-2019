#include <assert.h>
#include <stdio.h>
#include <stdlib.h>

#include "parse.h"
#include "vm.h"

int64_t part1(vm vm) {
    run(&vm);
    assert(vm.state == VM_INPUT);
    vm.input = 1;
    int64_t output;
    for (run(&vm); vm.state == VM_OUTPUT; run(&vm)) output = vm.output;
    assert(vm.state == VM_HALTED);
    return output;
}

int64_t part2(vm vm) { return 0; }

int main(int argc, char* argv[]) {
    if (argc != 2) {
        fprintf(stderr, "usage: day5 file\n");
        exit(1);
    }

    FILE* f = fopen(argv[1], "r");
    if (f == NULL) {
        perror("day5");
        exit(1);
    }

    vm vm = new_vm();
    if ((vm.mem_size = parse_intcode(f, vm.mem, MAX_MEM)) < 0) {
        perror("day5");
        exit(1);
    }

    printf("%lld\n", part1(vm));
    printf("%lld\n", part2(vm));

    return 0;
}
