#include "vm.h"

#include <assert.h>
#include <stdio.h>
#include <stdlib.h>

void init_vm(vm* vm) {
    for (int i = 0; i < MAX_MEM; i++) vm->mem[i] = 0;
    vm->mem_size = 0;
    vm->pc = 0;
    vm->state = VM_RUNNING;
}

vm new_vm(void) {
    vm vm;
    init_vm(&vm);
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

typedef enum {
    MODE_POS = 0,
    MODE_IMM = 1,
} mode;

typedef struct {
    opcode op;
    mode modes[3];
} instr;

instr parse_instr(int64_t val) {
    instr instr;
    instr.op = val % 100;
    instr.modes[0] = (val / 100) % 10;
    instr.modes[1] = (val / 1000) % 10;
    instr.modes[2] = (val / 10000) % 10;
    return instr;
}

void print_instr(int pc, instr instr) {
    printf("%08x\t%2d\t%d %d %d\n", pc, instr.op, instr.modes[0],
           instr.modes[1], instr.modes[2]);
}

int64_t eval_arg(vm* vm, mode mode, int64_t arg_ptr) {
    switch (mode) {
        case MODE_POS:
            return vm->mem[vm->mem[arg_ptr]];
        case MODE_IMM:
            return vm->mem[arg_ptr];
        default:
            vm->state = VM_ERROR;
            vm->error = INVALID_MODE;
            return 0;
    }
}

int64_t eval_dest(vm* vm, mode mode, int64_t arg_ptr) {
    switch (mode) {
        case MODE_POS:
            return vm->mem[arg_ptr];
        case MODE_IMM:
        default:
            vm->state = VM_ERROR;
            vm->error = INVALID_MODE;
            return 0;
    }
}

void step(vm* vm) {
    if (vm->pc < 0 || vm->pc >= vm->mem_size) {
        if (getenv("VM_TRACE")) printf("error: vm out of range\n");
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
        instr prev_instr = parse_instr(vm->mem[vm->pc]);
        int dest = eval_dest(vm, prev_instr.modes[0], vm->pc + 1);
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

    instr instr = parse_instr(vm->mem[vm->pc]);
    if (getenv("VM_TRACE")) print_instr(vm->pc, instr);
    switch (instr.op) {
        case ADD: {
            if (vm->pc > vm->mem_size - 4) {
                vm->error = PC_OUT_OF_RANGE;
                vm->state = VM_ERROR;
                return;
            }
            int64_t l = eval_arg(vm, instr.modes[0], vm->pc + 1);
            int64_t r = eval_arg(vm, instr.modes[1], vm->pc + 2);
            unsigned dest = eval_dest(vm, instr.modes[2], vm->pc + 3);
            vm->mem[dest] = l + r;
            if (getenv("VM_TRACE")) printf("%08x <- %lld\n", dest, l + r);
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
            int64_t l = eval_arg(vm, instr.modes[0], vm->pc + 1);
            int64_t r = eval_arg(vm, instr.modes[1], vm->pc + 2);
            unsigned dest = eval_dest(vm, instr.modes[2], vm->pc + 3);
            vm->mem[dest] = l * r;
            if (getenv("VM_TRACE")) printf("%08x <- %lld\n", dest, l * r);
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
            int64_t src = eval_arg(vm, instr.modes[0], vm->pc + 1);
            vm->output = src;
            vm->pc += 2;
            vm->state = VM_OUTPUT;
            return;
        }

        case JMP_IF: {
            if (vm->pc > vm->mem_size - 3) {
                vm->error = PC_OUT_OF_RANGE;
                vm->state = VM_ERROR;
                return;
            }
            int64_t arg = eval_arg(vm, instr.modes[0], vm->pc + 1);
            int64_t dest = eval_arg(vm, instr.modes[1], vm->pc + 2);
            if (arg)
                vm->pc = dest;
            else
                vm->pc += 3;
            return;
        }

        case JMP_NOT: {
            if (vm->pc > vm->mem_size - 3) {
                vm->error = PC_OUT_OF_RANGE;
                vm->state = VM_ERROR;
                return;
            }
            int64_t arg = eval_arg(vm, instr.modes[0], vm->pc + 1);
            int64_t dest = eval_arg(vm, instr.modes[1], vm->pc + 2);
            if (!arg)
                vm->pc = dest;
            else
                vm->pc += 3;
            return;
        }

        case LT: {
            if (vm->pc > vm->mem_size - 4) {
                vm->error = PC_OUT_OF_RANGE;
                vm->state = VM_ERROR;
                return;
            }
            int64_t l = eval_arg(vm, instr.modes[0], vm->pc + 1);
            int64_t r = eval_arg(vm, instr.modes[1], vm->pc + 2);
            unsigned dest = eval_dest(vm, instr.modes[2], vm->pc + 3);
            vm->mem[dest] = l < r;
            if (getenv("VM_TRACE")) printf("%08x <- %d\n", dest, l < r);
            vm->pc += 4;
            return;
        }

        case EQ: {
            if (vm->pc > vm->mem_size - 4) {
                vm->error = PC_OUT_OF_RANGE;
                vm->state = VM_ERROR;
                return;
            }
            int64_t l = eval_arg(vm, instr.modes[0], vm->pc + 1);
            int64_t r = eval_arg(vm, instr.modes[1], vm->pc + 2);
            unsigned dest = eval_dest(vm, instr.modes[2], vm->pc + 3);
            vm->mem[dest] = l == r;
            if (getenv("VM_TRACE")) printf("%08x <- %d\n", dest, l == r);
            vm->pc += 4;
            return;
        }

        case HALT: {
            vm->state = VM_HALTED;
            return;
        }

        default: {
            if (getenv("VM_TRACE")) printf("error: invalid opcode\n");
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
