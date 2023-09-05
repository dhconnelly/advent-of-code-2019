#include <assert.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

#include "parse.h"
#include "vm.h"

void play(vm* vm) {
    while (vm->state != VM_HALTED) {
        switch (vm->state) {
            case VM_INPUT: {
                char* line = NULL;
                size_t len;
                int result = getline(&line, &len, stdin);
                if (feof(stdin)) {
                    free(line);
                    printf("Goodbye\n");
                    return;
                } else if (result < 0) {
                    free(line);
                    perror("day25: failed to read input");
                    exit(1);
                }
                for (int i = 0; i < len && vm->state == VM_INPUT; i++) {
                    vm->input = line[i];
                    run(vm);
                }
                free(line);
                break;
            }
            case VM_OUTPUT:
                while (vm->state == VM_OUTPUT) {
                    putchar(vm->output);
                    run(vm);
                }
                break;
            default:
                assert(0);
        }
    }
}

int main(int argc, char* argv[]) {
    if (argc != 2) {
        fprintf(stderr, "usage: day25 file\n");
        exit(1);
    }

    FILE* f = fopen(argv[1], "r");
    if (f == NULL) {
        perror("day25");
        exit(1);
    }

    int64_t data[8192];
    int len;
    if ((len = parse_intcode(f, data, 8192)) < 0) {
        perror("day25: failed to parse intcode");
        exit(1);
    }

    vm vm = make_vm(data, len);
    run(&vm);
    play(&vm);
}
