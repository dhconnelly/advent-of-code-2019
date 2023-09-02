#ifndef AMP_H_
#define AMP_H_

#include "vm.h"

#define CIRCUIT_SIZE 5

typedef int64_t phase_sequence[CIRCUIT_SIZE];

typedef struct {
    vm vms[CIRCUIT_SIZE];
} circuit;

void init_circuit(circuit* circuit, vm vm);
circuit make_circuit(vm vm);
int64_t run_series(circuit* circuit, phase_sequence setting, int64_t input);

#endif  // AMP_H_
