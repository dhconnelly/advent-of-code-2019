#ifndef VM_H_
#define VM_H_

#include <stdint.h>

#include "hash.h"

typedef enum {
    VM_HALTED,
    VM_RUNNING,
    VM_ERROR,
    VM_INPUT,
    VM_OUTPUT,
} state;

typedef enum {
    ADD = 1,
    MUL = 2,
    IN = 3,
    OUT = 4,
    JMP_IF = 5,
    JMP_NOT = 6,
    LT = 7,
    EQ = 8,
    ADJREL = 9,
    HALT = 99,
} opcode;

typedef enum {
    ERR_NONE,
    ERR_PC_OUT_OF_RANGE,
    ERR_INVALID_OPCODE,
    ERR_INVALID_MODE,
} error;

typedef struct {
    state state;
    error error;
    int pc;
    int relbase;
    hashtable mem;
    int64_t input;
    int64_t output;
    int trace;
} vm;

vm new_vm(void);
vm make_vm(int64_t mem[], int mem_size);
vm copy_vm(const vm* vm);
void run(vm* vm);
int64_t get_mem(vm* vm, int addr);
void set_mem(vm* vm, int addr, int64_t val);

#endif  // VM_H_
