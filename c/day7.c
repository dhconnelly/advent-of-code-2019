#include <assert.h>
#include <math.h>
#include <stdio.h>
#include <stdlib.h>

#include "amp.h"
#include "comb.h"
#include "parse.h"
#include "vm.h"

typedef int64_t (*execute)(circuit*, phase_sequence, int64_t);

typedef struct {
    execute f;
    vm vm;
    int64_t max;
    int64_t input;
} acc;

void* find_max(int64_t seq[], void* data) {
    acc* acc = data;
    circuit circuit = make_circuit(acc->vm);
    int64_t output = acc->f(&circuit, seq, acc->input);
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

    int64_t data[1024];
    int len;
    if ((len = parse_intcode(f, data, 1024)) < 0) {
        perror("day7");
        exit(1);
    }

    vm base = make_vm(data, len);

    acc acc1 = {.f = &run_series, .vm = copy_vm(&base), .input = 0, .max = ~0};
    phase_sequence seq1 = {0, 1, 2, 3, 4};
    visit_permutations(seq1, 5, find_max, &acc1);
    printf("%lld\n", acc1.max);

    acc acc2 = {.f = &run_loop, .vm = copy_vm(&base), .input = 0, .max = ~0};
    phase_sequence seq2 = {5, 6, 7, 8, 9};
    visit_permutations(seq2, 5, find_max, &acc2);
    printf("%lld\n", acc2.max);
}
