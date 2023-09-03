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

#endif  // PT2_H_
