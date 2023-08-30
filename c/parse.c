#include "parse.h"

#include <stdio.h>
#include <stdlib.h>

int parse_intcode(FILE *f, int result[], int max) {
    int n = 0;
    while (n < max && fscanf(f, "%d,", &result[n]) != EOF) n++;
    return n;
}
