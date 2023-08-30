#ifndef VM_H_
#define VM_H_

#include <stdint.h>

#define MAX_MEM 1024

typedef enum {
    HALTED,
    RUNNING,
    ERROR,
} state;

typedef enum {
    ADD = 1,
    MUL = 2,
    HALT = 99,
} opcode;

typedef enum {
    PC_OUT_OF_RANGE = 1,
    INVALID_OPCODE = 2,
} error;

typedef struct {
    state state;
    error error;
    int pc;
    int mem_size;
    int64_t mem[MAX_MEM];
} vm;

vm new_vm(void);
vm make_vm(int64_t mem[], int mem_size);
void print_vm(const vm* vm);
void step(vm* vm);
void run(vm* vm);

#endif  // VM_H_
