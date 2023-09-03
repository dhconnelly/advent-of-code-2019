#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "amp.h"
#include "comb.h"

static void day7_test1() {
    printf("day7_test1\n");
    int64_t prog[] = {3,  15, 3,  16, 1002, 16, 10, 16, 1,
                      16, 15, 15, 4,  15,   99, 0,  0};
    vm vm = make_vm(prog, sizeof(prog) / sizeof(int64_t));
    circuit circuit = make_circuit(vm);
    phase_sequence seq = {4, 3, 2, 1, 0};
    int64_t output = run_series(&circuit, seq, 0);
    assert(output == 43210);
}

static void day7_test2() {
    printf("day7_test2\n");
    int64_t prog[] = {3, 23, 3,  24, 1002, 24, 10, 24, 1002, 23, -1, 23, 101,
                      5, 23, 23, 1,  24,   23, 23, 4,  23,   99, 0,  0};
    vm vm = make_vm(prog, sizeof(prog) / sizeof(int64_t));
    circuit circuit = make_circuit(vm);
    phase_sequence seq = {0, 1, 2, 3, 4};
    int64_t output = run_series(&circuit, seq, 0);
    assert(output == 54321);
}

static void day7_test3() {
    printf("day7_test3\n");
    int64_t prog[] = {3,    31, 3,  32, 1002, 32, 10, 32, 1001, 31, -2, 31,
                      1007, 31, 0,  33, 1002, 33, 7,  33, 1,    33, 31, 31,
                      1,    32, 31, 31, 4,    31, 99, 0,  0,    0};
    vm vm = make_vm(prog, sizeof(prog) / sizeof(int64_t));
    circuit circuit = make_circuit(vm);
    phase_sequence seq = {1, 0, 4, 3, 2};
    int64_t output = run_series(&circuit, seq, 0);
    assert(output == 65210);
}

typedef struct perm_list {
    int64_t seq[CIRCUIT_SIZE];
    struct perm_list* next;
} perm_list;

static int perm_len(const perm_list* list) {
    int len = 0;
    for (const perm_list* node = list; node != NULL; node = node->next) ++len;
    return len;
}

static void* accumulate(int64_t perm[], void* data) {
    perm_list* list = malloc(sizeof(perm_list));
    for (int i = 0; i < CIRCUIT_SIZE; i++) list->seq[i] = perm[i];
    list->next = data;
    return list;
}

static void extract(const perm_list* list, char* strs[], int perm_len) {
    for (int i = 0; list != NULL; list = list->next, i++) {
        for (int j = 0; j < perm_len; j++) strs[i][j] = '0' + list->seq[j];
        strs[i][perm_len] = '\0';
    }
}

static int permcmp(const void* left, const void* right) {
    const char **s = (const char**)left, **t = (const char**)right;
    int x = atoi(*s), y = atoi(*t);
    return x - y;
}

void day7_test4() {
    int64_t seq[] = {0, 1, 2};
    perm_list* result = visit_permutations(seq, 3, accumulate, NULL);
    int len = perm_len(result);
    assert(len == 6);

    char* perms[6];
    for (int i = 0; i < 6; i++) perms[i] = malloc(4 * sizeof(char));
    extract(result, perms, 3);

    qsort(perms, 6, sizeof(char*), permcmp);
    char* expected[] = {"012", "021", "102", "120", "201", "210"};
    for (int i = 0; i < 6; i++) {
        assert(strcmp(perms[i], expected[i]) == 0);
        free(perms[i]);
    }
}

void day7_test() {
    day7_test1();
    day7_test2();
    day7_test3();
    day7_test4();
}
