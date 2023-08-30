#include <assert.h>

#include "vm.h"

static void test1() {
    int data[] = {1, 9, 10, 3, 2, 3, 11, 0, 99, 30, 40, 50};
    vm vm = make_vm(data, sizeof(data) / sizeof(data[0]));
    run(&vm);
    assert(vm.mem[0] == 3500);
}

static void test2() {
    int data[] = {1, 0, 0, 0, 99};
    vm vm = make_vm(data, sizeof(data) / sizeof(data[0]));
    run(&vm);
    assert(vm.mem[0] == 2);
}

static void test3() {
    int data[] = {2, 3, 0, 3, 99};
    vm vm = make_vm(data, sizeof(data) / sizeof(data[0]));
    run(&vm);
    assert(vm.mem[3] == 6);
}

static void test4() {
    int data[] = {2, 4, 4, 5, 99, 0};
    vm vm = make_vm(data, sizeof(data) / sizeof(data[0]));
    run(&vm);
    assert(vm.mem[5] == 9801);
}

static void test5() {
    int data[] = {1, 1, 1, 4, 99, 5, 6, 0, 99};
    vm vm = make_vm(data, sizeof(data) / sizeof(data[0]));
    run(&vm);
    assert(vm.mem[0] == 30);
}

void day2_test() {
    test1();
    test2();
    test3();
    test4();
    test5();
}
