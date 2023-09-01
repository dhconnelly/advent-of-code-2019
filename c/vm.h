#ifndef VM_H_
#define VM_H_

#include <stdint.h>

#define MAX_MEM 1024

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
    HALT = 99,
} opcode;

typedef enum {
    NO_ERROR,
    PC_OUT_OF_RANGE,
    INVALID_OPCODE,
    INVALID_MODE,
} error;

typedef struct {
    state state;
    error error;
    int pc;
    int mem_size;
    int64_t mem[MAX_MEM];
    int64_t input;
    int64_t output;
} vm;

vm new_vm(void);
vm make_vm(int64_t mem[], int mem_size);
void print_vm(const vm* vm);
void step(vm* vm);
void run(vm* vm);

#endif  // VM_H_
