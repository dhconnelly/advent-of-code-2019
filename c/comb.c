#include "comb.h"

#include <stdio.h>

static void* permute(int64_t xs[], int from, int len, visitor f, void* init) {
    if (from == len - 1) return f(xs, init);
    int64_t tmp;
    void* result = init;
    for (int to = from; to < len; to++) {
        tmp = xs[from];
        xs[from] = xs[to];
        xs[to] = tmp;
        result = permute(xs, from + 1, len, f, result);
        xs[to] = xs[from];
        xs[from] = tmp;
    }
    return result;
}

void* visit_permutations(int64_t xs[], int len, visitor f, void* init) {
    return permute(xs, 0, len, f, init);
}
