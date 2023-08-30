#include "parse.h"

#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

int parse_intcode(FILE *f, int64_t result[], int max) {
    int n = 0;
    while (n < max && fscanf(f, "%lld,", &result[n]) != EOF) n++;
    return n;
}
