#include "vm.h"

#include <assert.h>
#include <stdio.h>
#include <stdlib.h>

typedef enum {
    MODE_POS = 0,
    MODE_IMM = 1,
    MODE_REL = 2,
} mode;

typedef struct {
    opcode op;
    mode modes[3];
} instr;

void set_mem(vm* vm, int loc, int64_t val) { table_set(vm->mem, loc, val); }

int64_t get_mem(vm* vm, int loc) {
    int64_t* val = table_get(vm->mem, loc);
    return val == NULL ? 0 : *val;
}

int64_t eval_arg(vm* vm, mode mode, int64_t arg_ptr) {
    switch (mode) {
        case MODE_POS:
            return get_mem(vm, get_mem(vm, arg_ptr));
        case MODE_IMM:
            return get_mem(vm, arg_ptr);
        case MODE_REL:
            return get_mem(vm, vm->relbase + get_mem(vm, arg_ptr));
        default:
            vm->state = VM_ERROR;
            vm->error = INVALID_MODE;
            return 0;
    }
}

int64_t eval_dest(vm* vm, mode mode, int64_t arg_ptr) {
    switch (mode) {
        case MODE_POS:
            return get_mem(vm, arg_ptr);
        case MODE_REL:
            return vm->relbase + get_mem(vm, arg_ptr);
        case MODE_IMM:
        default:
            vm->state = VM_ERROR;
            vm->error = INVALID_MODE;
            return 0;
    }
}

void init_vm(vm* vm) {
    vm->pc = 0;
    vm->state = VM_RUNNING;
    init_table(vm->mem);
    vm->error = 0;
    vm->input = 0;
    vm->output = 0;
    vm->relbase = 0;
}

vm new_vm(void) {
    vm vm;
    init_vm(&vm);
    return vm;
}

vm make_vm(int64_t mem[], int mem_size) {
    vm vm = new_vm();
    fill_table(vm.mem, mem, mem_size);
    for (int i = 0; i < mem_size; i++) table_set(vm.mem, i, mem[i]);
    return vm;
}

vm copy_vm(const vm* base) {
    vm local = *base;
    table_copy(local.mem, base->mem);
    return local;
}

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

void step(vm* vm) {
    if (vm->state == VM_INPUT) {
        instr prev_instr = parse_instr(*table_get(vm->mem, vm->pc));
        int dest = eval_dest(vm, prev_instr.modes[0], vm->pc + 1);
        set_mem(vm, dest, vm->input);
        vm->pc += 2;
        vm->state = VM_RUNNING;
    }

    if (vm->state == VM_OUTPUT) {
        vm->state = VM_RUNNING;
    }

    if (vm->state == VM_HALTED || vm->state == VM_ERROR) {
        return;
    }

    instr instr = parse_instr(*table_get(vm->mem, vm->pc));
    if (getenv("VM_TRACE")) print_instr(vm->pc, instr);
    switch (instr.op) {
        case ADD: {
            int64_t l = eval_arg(vm, instr.modes[0], vm->pc + 1);
            int64_t r = eval_arg(vm, instr.modes[1], vm->pc + 2);
            unsigned dest = eval_dest(vm, instr.modes[2], vm->pc + 3);
            set_mem(vm, dest, l + r);
            if (getenv("VM_TRACE")) printf("%08x <- %lld\n", dest, l + r);
            vm->pc += 4;
            break;
            return;
        }

        case MUL: {
            int64_t l = eval_arg(vm, instr.modes[0], vm->pc + 1);
            int64_t r = eval_arg(vm, instr.modes[1], vm->pc + 2);
            unsigned dest = eval_dest(vm, instr.modes[2], vm->pc + 3);
            set_mem(vm, dest, l * r);
            if (getenv("VM_TRACE")) printf("%08x <- %lld\n", dest, l * r);
            vm->pc += 4;
            return;
        }

        case IN: {
            vm->state = VM_INPUT;
            return;
        }

        case OUT: {
            int64_t src = eval_arg(vm, instr.modes[0], vm->pc + 1);
            vm->output = src;
            vm->pc += 2;
            vm->state = VM_OUTPUT;
            return;
        }

        case JMP_IF: {
            int64_t arg = eval_arg(vm, instr.modes[0], vm->pc + 1);
            int64_t dest = eval_arg(vm, instr.modes[1], vm->pc + 2);
            if (arg)
                vm->pc = dest;
            else
                vm->pc += 3;
            return;
        }

        case JMP_NOT: {
            int64_t arg = eval_arg(vm, instr.modes[0], vm->pc + 1);
            int64_t dest = eval_arg(vm, instr.modes[1], vm->pc + 2);
            if (!arg)
                vm->pc = dest;
            else
                vm->pc += 3;
            return;
        }

        case LT: {
            int64_t l = eval_arg(vm, instr.modes[0], vm->pc + 1);
            int64_t r = eval_arg(vm, instr.modes[1], vm->pc + 2);
            unsigned dest = eval_dest(vm, instr.modes[2], vm->pc + 3);
            set_mem(vm, dest, l < r);
            if (getenv("VM_TRACE")) printf("%08x <- %d\n", dest, l < r);
            vm->pc += 4;
            return;
        }

        case EQ: {
            int64_t l = eval_arg(vm, instr.modes[0], vm->pc + 1);
            int64_t r = eval_arg(vm, instr.modes[1], vm->pc + 2);
            unsigned dest = eval_dest(vm, instr.modes[2], vm->pc + 3);
            set_mem(vm, dest, l == r);
            if (getenv("VM_TRACE")) printf("%08x <- %d\n", dest, l == r);
            vm->pc += 4;
            return;
        }

        case ADJREL: {
            int64_t offset = eval_arg(vm, instr.modes[0], vm->pc + 1);
            vm->relbase += offset;
            vm->pc += 2;
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
