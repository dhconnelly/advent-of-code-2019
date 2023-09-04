#ifndef PT2_H_
#define PT2_H_

#include <stdint.h>

#include "hash.h"

typedef union {
    struct {
        int16_t x;
        int16_t y;
    } coords;
    uint32_t data;
} pt2;

pt2 make_pt(int16_t x, int16_t y);
pt2 pt_from_data(uint32_t data);
int pt_eq(pt2 a, pt2 b);
void get_nbrs(pt2 pt, pt2 nbrs[4]);

typedef struct {
    pt2 lo, hi;
} rect;

rect bounds(hashtable* map);

#endif  // PT2_H_
