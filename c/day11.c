#include <stdio.h>
#include <stdlib.h>

#include "hash.h"
#include "parse.h"
#include "pt2.h"

typedef enum { N, E, S, W } direction;

typedef struct {
    direction dir;
    pt2 pos;
    hashtable tiles;
} robot;

int main(int argc, char* argv[]) {
    if (argc != 2) {
        fprintf(stderr, "usage: day11 file\n");
        exit(1);
    }

    FILE* f = fopen(argv[1], "r");
    if (f == NULL) {
        perror("day11");
        exit(1);
    }

    int64_t data[1024];
    int len;
    if ((len = parse_intcode(f, data, 1024)) < 0) {
        perror("day11: failed to parse intcode");
        exit(1);
    }

    printf("read %d elements\n", len);
}
