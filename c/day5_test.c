#include <assert.h>
#include <stdio.h>

#include "vm.h"

static void day5_test1() {
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

static void day5_test2() {
    printf("day5_test2\n");
    int64_t prog[] = {1002, 4, 3, 4, 33};
    vm vm = make_vm(prog, sizeof(prog) / sizeof(int64_t));
    run(&vm);
    assert(get_mem(&vm, 4) == 99);
}

static int64_t evaluate(vm base, int64_t input) {
    vm f = copy_vm(&base);
    run(&f);
    assert(f.state == VM_INPUT);
    f.input = input;
    run(&f);
    assert(f.state == VM_OUTPUT);
    int64_t output = f.output;
    run(&f);
    assert(f.state == VM_HALTED);
    return output;
}

static void test_equals() {
    printf("day5_test3 eq\n");
    int64_t prog1[] = {3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8};
    int64_t prog2[] = {3, 3, 1108, -1, 8, 3, 4, 3, 99};
    vm vms[2] = {
        make_vm(prog1, sizeof(prog1) / sizeof(int64_t)),
        make_vm(prog2, sizeof(prog2) / sizeof(int64_t)),
    };
    for (int i = 0; i < 2; i++) {
        assert(evaluate(vms[i], 7) == 0);
        assert(evaluate(vms[i], 8) == 1);
        assert(evaluate(vms[i], 9) == 0);
    }
}

static void test_lt() {
    printf("day5_test3 lt\n");
    int64_t prog1[] = {3, 9, 7, 9, 10, 9, 4, 9, 99, -1, 8};
    int64_t prog2[] = {3, 3, 1107, -1, 8, 3, 4, 3, 99};
    vm vms[2] = {
        make_vm(prog1, sizeof(prog1) / sizeof(int64_t)),
        make_vm(prog2, sizeof(prog2) / sizeof(int64_t)),
    };
    for (int i = 0; i < 2; i++) {
        assert(evaluate(vms[i], 7));
        assert(!evaluate(vms[i], 8));
        assert(!evaluate(vms[i], 9));
    }
}

static void test_jmp() {
    printf("day5_test3 jmp\n");
    int64_t prog1[] = {3, 12, 6, 12, 15, 1, 13, 14, 13, 4, 13, 99, -1, 0, 1, 9};
    int64_t prog2[] = {3, 3, 1105, -1, 9, 1101, 0, 0, 12, 4, 12, 99, 1};
    vm vms[2] = {
        make_vm(prog1, sizeof(prog1) / sizeof(int64_t)),
        make_vm(prog2, sizeof(prog2) / sizeof(int64_t)),
    };
    for (int i = 0; i < 2; i++) {
        assert(evaluate(vms[i], 0) == 0);
        assert(evaluate(vms[i], -17) == 1);
        assert(evaluate(vms[i], 1) == 1);
        assert(evaluate(vms[i], 9) == 1);
    }
}

static void day5_test3() {
    printf("day5_test3\n");

    test_equals();
    test_lt();
    test_jmp();

    int64_t prog[] = {3,  21,  1008, 21,   8,   20, 1005, 20,   22,   107,
                      8,  21,  20,   1006, 20,  31, 1106, 0,    36,   98,
                      0,  0,   1002, 21,   125, 20, 4,    20,   1105, 1,
                      46, 104, 999,  1105, 1,   46, 1101, 1000, 1,    20,
                      4,  20,  1105, 1,    46,  98, 99};
    vm vm = make_vm(prog, sizeof(prog) / sizeof(int64_t));
    assert(evaluate(vm, 7) == 999);
    assert(evaluate(vm, 8) == 1000);
    assert(evaluate(vm, 9) == 1001);
}

void day5_test() {
    day5_test1();
    day5_test2();
    day5_test3();
}
