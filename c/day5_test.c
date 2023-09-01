#include <assert.h>
#include <stdio.h>

#include "vm.h"

void day5_test1() {
    printf("day5_test1\n");
    int64_t prog[] = {3, 0, 4, 0, 99};
    vm vm = make_vm(prog, sizeof(prog) / sizeof(int64_t));
    run(&vm);
    assert(vm.state == VM_INPUT);
    vm.input = 12345;
    run(&vm);
    assert(vm.state == VM_OUTPUT);
    assert(vm.output == 12345);
    run(&vm);
    assert(vm.state == VM_HALTED);
}

void day5_test() { day5_test1(); }
