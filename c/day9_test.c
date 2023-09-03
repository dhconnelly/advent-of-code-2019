#include <assert.h>
#include <stdio.h>
#include <stdlib.h>

#include "vm.h"

static int execute(vm* vm, int64_t outputs[], int max_outputs) {
    int i = 0;
    do {
        run(vm);
        switch (vm->state) {
            case VM_INPUT:
                printf("unexpected INPUT\n");
                exit(1);
            case VM_OUTPUT:
                assert(i < max_outputs);
                outputs[i++] = vm->output;
                break;
            case VM_RUNNING:
            case VM_HALTED:
                break;
            case VM_ERROR:
                fprintf(stderr, "error: %d\n", vm->error);
                assert(0);
        }
    } while (vm->state != VM_HALTED);
    return i;
}

static void day9_test1(void) {
    printf("day9_test1\n");
    int64_t prog[] = {109,  1,   204, -1,  1001, 100, 1, 100,
                      1008, 100, 16,  101, 1006, 101, 0, 99};
    vm vm = make_vm(prog, sizeof(prog) / sizeof(int64_t));
    int64_t outputs[100];
    int n = execute(&vm, outputs, 100);
    for (int i = 0; i < n; i++) {
        assert(outputs[i] == prog[i]);
    }
}

static void day9_test2(void) {
    printf("day9_test2\n");
    int64_t prog[] = {1102, 34915192, 34915192, 7, 4, 7, 99, 0};
    vm vm = make_vm(prog, sizeof(prog) / sizeof(int64_t));
    int64_t outputs[1];
    int n = execute(&vm, outputs, 1);
    assert(outputs[0] == 1219070632396864);
}

static void day9_test3(void) {
    printf("day9_test3\n");
    int64_t prog[] = {104, 1125899906842624, 99};
    vm vm = make_vm(prog, sizeof(prog) / sizeof(int64_t));
    int64_t outputs[100];
    int n = execute(&vm, outputs, 100);
    assert(outputs[0] == prog[1]);
}

void day9_test(void) {
    printf("day9_test\n");
    day9_test1();
    day9_test2();
    day9_test3();
}
