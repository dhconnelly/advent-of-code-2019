#include "vm.h"

#include <assert.h>
#include <stdio.h>
#include <stdlib.h>

vm new_vm(void) {
    vm vm;
    for (int i = 0; i < MAX_MEM; i++) vm.mem[i] = 0;
    vm.mem_size = 0;
    vm.pc = 0;
    vm.state = VM_RUNNING;
    return vm;
}

vm make_vm(int64_t mem[], int mem_size) {
    assert(mem_size < MAX_MEM);
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
    if (vm->pc < 0 || vm->pc >= vm->mem_size) {
        vm->state = VM_ERROR;
        vm->error = PC_OUT_OF_RANGE;
        return;
    }

    if (vm->state == VM_INPUT) {
        if (vm->pc > vm->mem_size - 4) {
            vm->error = PC_OUT_OF_RANGE;
            vm->state = VM_ERROR;
            return;
        }
        int dest = vm->mem[vm->pc + 1];
        vm->mem[dest] = vm->input;
        vm->pc += 2;
        vm->state = VM_RUNNING;
    }

    if (vm->state == VM_OUTPUT) {
        vm->state = VM_RUNNING;
    }

    if (vm->state == VM_HALTED || vm->state == VM_ERROR) {
        return;
    }

    opcode op = vm->mem[vm->pc];
    if (getenv("VM_TRACE")) printf("%08x\t%4d\n", vm->pc, op);
    switch (op) {
        case ADD: {
            if (vm->pc > vm->mem_size - 4) {
                vm->error = PC_OUT_OF_RANGE;
                vm->state = VM_ERROR;
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
                vm->state = VM_ERROR;
                return;
            }
            int l = vm->mem[vm->pc + 1];
            int r = vm->mem[vm->pc + 2];
            int dest = vm->mem[vm->pc + 3];
            vm->mem[dest] = vm->mem[l] * vm->mem[r];
            vm->pc += 4;
            return;
        }

        case IN: {
            if (vm->pc > vm->mem_size - 2) {
                vm->error = PC_OUT_OF_RANGE;
                vm->state = VM_ERROR;
                return;
            }
            vm->state = VM_INPUT;
            return;
        }

        case OUT: {
            if (vm->pc > vm->mem_size - 2) {
                vm->error = PC_OUT_OF_RANGE;
                vm->state = VM_ERROR;
                return;
            }
            int src = vm->mem[vm->pc + 1];
            vm->output = vm->mem[src];
            vm->pc += 2;
            vm->state = VM_OUTPUT;
            return;
        }

        case HALT: {
            vm->state = VM_HALTED;
            return;
        }

        default: {
            vm->state = VM_ERROR;
            vm->error = INVALID_OPCODE;
            return;
        }
    }
}

void run(vm* vm) {
    do {
        step(vm);
    } while (vm->state == VM_RUNNING);
}
