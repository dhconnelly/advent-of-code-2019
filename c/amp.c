#include "amp.h"

#include <assert.h>
#include <stdio.h>

void init_circuit(circuit* circuit, vm vm) {
    for (int i = 0; i < CIRCUIT_SIZE; i++) circuit->vms[i] = vm;
}

circuit make_circuit(vm vm) {
    circuit circuit;
    init_circuit(&circuit, vm);
    return circuit;
}

typedef int64_t phase_sequence[CIRCUIT_SIZE];

int64_t run_series(circuit* circuit, phase_sequence setting, int64_t input) {
    int64_t prev_output = input;
    for (int i = 0; i < CIRCUIT_SIZE; i++) {
        vm* vm = &circuit->vms[i];
        run(vm);
        assert(vm->state == VM_INPUT);
        vm->input = setting[i];
        run(vm);
        assert(vm->state == VM_INPUT);
        vm->input = prev_output;
        run(vm);
        assert(vm->state == VM_OUTPUT);
        prev_output = vm->output;
    }
    return prev_output;
}

int64_t run_loop(circuit* circuit, phase_sequence setting, int64_t input) {
    int64_t prev_output = input;
    for (int iter = 0, done = 0; !done; iter++) {
        for (int i = 0; i < CIRCUIT_SIZE; i++) {
            vm* vm = &circuit->vms[i];
            if (iter == 0) {
                run(vm);
                assert(vm->state == VM_INPUT);
                vm->input = setting[i];
            }
            run(vm);
            if (vm->state == VM_HALTED) {
                done = 1;
                break;
            }
            assert(vm->state == VM_INPUT);
            vm->input = prev_output;
            run(vm);
            assert(vm->state == VM_OUTPUT);
            prev_output = vm->output;
        }
    }
    return prev_output;
}
