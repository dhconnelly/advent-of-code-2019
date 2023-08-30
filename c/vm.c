#include "vm.h"

#include <stdio.h>

vm new_vm(void) {
    vm vm;
    for (int i = 0; i < MAX_MEM; i++) vm.mem[i] = 0;
    vm.mem_size = 0;
    vm.pc = 0;
    vm.state = RUNNING;
    return vm;
}

vm make_vm(int64_t mem[], int mem_size) {
    vm vm = new_vm();
    for (int i = 0; i < mem_size && i < MAX_MEM; i++) vm.mem[i] = mem[i];
    vm.mem_size = mem_size;
    return vm;
}

void print_vm(const vm* vm) {
    printf("vm {\n");
    printf("  pc = %d\n", vm->pc);
    printf("  mem_size = %d\n", vm->mem_size);
    printf("  mem = [");
    for (int i = 0; i < vm->mem_size; i++) {
        if (i != 0) printf(", ");
        printf("%lld", vm->mem[i]);
    }
    printf("]\n");
    printf("}\n");
}

void step(vm* vm) {
    if (vm->state != RUNNING) {
        return;
    }
    if (vm->pc < 0 || vm->pc >= vm->mem_size) {
        vm->state = ERROR;
        vm->error = PC_OUT_OF_RANGE;
        return;
    }

    opcode op = vm->mem[vm->pc];
    switch (op) {
        case ADD: {
            if (vm->pc > vm->mem_size - 4) {
                vm->error = PC_OUT_OF_RANGE;
                vm->state = ERROR;
                return;
            }
            int l = vm->mem[vm->pc + 1];
            int r = vm->mem[vm->pc + 2];
            int dest = vm->mem[vm->pc + 3];
            vm->mem[dest] = vm->mem[l] + vm->mem[r];
            vm->pc += 4;
            break;
            return;
        }

        case MUL: {
            if (vm->pc > vm->mem_size - 4) {
                vm->error = PC_OUT_OF_RANGE;
                vm->state = ERROR;
                return;
            }
            int l = vm->mem[vm->pc + 1];
            int r = vm->mem[vm->pc + 2];
            int dest = vm->mem[vm->pc + 3];
            vm->mem[dest] = vm->mem[l] * vm->mem[r];
            vm->pc += 4;
            return;
        }

        case HALT: {
            vm->state = HALTED;
            return;
        }

        default: {
            vm->state = ERROR;
            vm->error = INVALID_OPCODE;
            return;
        }
    }
}

void run(vm* vm) {
    while (vm->state == RUNNING) step(vm);
}
