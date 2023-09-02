#include <assert.h>
#include <math.h>
#include <stdio.h>
#include <stdlib.h>

#include "amp.h"
#include "comb.h"
#include "parse.h"
#include "vm.h"

typedef struct {
    vm vm;
    int64_t max;
    int64_t input;
} acc;

void* find_max(int64_t seq[], void* data) {
    acc* acc = data;
    circuit circuit = make_circuit(acc->vm);
    int64_t output = run_series(&circuit, seq, acc->input);
    if (output > acc->max) acc->max = output;
    return acc;
}

int main(int argc, char* argv[]) {
    if (argc != 2) {
        fprintf(stderr, "usage: day7 file\n");
        exit(1);
    }
    FILE* f = fopen(argv[1], "r");
    if (f == NULL) {
        perror("day7");
        exit(1);
    }
    vm vm = new_vm();
    if ((vm.mem_size = parse_intcode(f, vm.mem, MAX_MEM)) < 0) {
        perror("day7: failed to parse intcode");
        exit(1);
    }

    acc acc;
    acc.vm = vm;
    acc.input = 0;
    acc.max = ~0;
    phase_sequence seq = {0, 1, 2, 3, 4};
    visit_permutations(seq, 5, find_max, &acc);

    printf("%lld\n", acc.max);
}
