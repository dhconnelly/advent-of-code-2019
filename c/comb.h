#ifndef COMB_H_
#define COMB_H_

#include <stdint.h>

typedef void* visitor(int64_t xs[], void* data);
void* visit_permutations(int64_t xs[], int len, visitor f, void* init);

#endif  // COMB_H_
