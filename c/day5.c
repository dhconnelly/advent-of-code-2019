#include <assert.h>
#include <stdio.h>
#include <stdlib.h>

#include "parse.h"
#include "vm.h"

int64_t evaluate(vm* vm, int64_t input) {
    run(vm);
    assert(vm->state == VM_INPUT);
    vm->input = input;
    int64_t output;
    for (run(vm); vm->state == VM_OUTPUT; run(vm)) output = vm->output;
    assert(vm->state == VM_HALTED);
    return output;
}

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

    vm base = new_vm();
    int64_t data[1024];
    int len;
    if ((len = parse_intcode(f, data, 1024)) < 0) {
        perror("day5");
        exit(1);
    }
    fill_table(base.mem, data, 1024);

    vm local = copy_vm(&base);
    printf("%lld\n", evaluate(&local, 1));

    local = copy_vm(&base);
    printf("%lld\n", evaluate(&local, 5));

    return 0;
}
