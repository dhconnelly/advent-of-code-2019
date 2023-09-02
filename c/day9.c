#include <assert.h>
#include <stdio.h>
#include <stdlib.h>

#include "parse.h"
#include "vm.h"

void execute(vm vm, int64_t input) {
    do {
        run(&vm);
        switch (vm.state) {
            case VM_ERROR:
                fprintf(stderr, "vm error: %d\n", vm.error);
                return;
            case VM_OUTPUT:
                printf("%lld\n", vm.output);
                break;
            case VM_INPUT:
                vm.input = input;
                break;
            case VM_RUNNING:
                break;
            case VM_HALTED:
                break;
        }
    } while (vm.state != VM_HALTED);
}

int main(int argc, char* argv[]) {
    if (argc != 2) {
        fprintf(stderr, "usage: day9 file\n");
        exit(1);
    }
    FILE* f = fopen(argv[1], "r");
    if (f == NULL) {
        perror("day9");
        exit(1);
    }
    vm vm = new_vm();
    int64_t data[1024];
    int len;
    if ((len = parse_intcode(f, data, 1024)) < 0) {
        perror("day9");
        exit(1);
    }
    fill_table(vm.mem, data, 1024);
    execute(copy_vm(&vm), 1);
    execute(copy_vm(&vm), 2);
}
