#include "pt2.h"

#include <limits.h>
#include <stdlib.h>

pt2 make_pt(int16_t x, int16_t y) {
    pt2 pt;
    pt.coords.x = x;
    pt.coords.y = y;
    return pt;
}

pt2 pt_from_data(uint32_t data) {
    pt2 pt;
    pt.data = data;
    return pt;
}

int pt_eq(pt2 a, pt2 b) { return a.data == b.data; }

void get_nbrs(pt2 pt, pt2 nbrs[4]) {
    nbrs[0] = make_pt(pt.coords.x, pt.coords.y - 1);
    nbrs[1] = make_pt(pt.coords.x, pt.coords.y + 1);
    nbrs[2] = make_pt(pt.coords.x - 1, pt.coords.y);
    nbrs[3] = make_pt(pt.coords.x + 1, pt.coords.y);
}

rect bounds(hashtable* map) {
    uint32_t* keys = table_keys(map);
    rect lohi;
    lohi.lo.coords.y = INT16_MAX, lohi.hi.coords.y = INT16_MIN,
    lohi.lo.coords.x = INT16_MAX, lohi.hi.coords.x = INT16_MIN;
    for (int i = 0; i < table_size(map); i++) {
        pt2 pt = pt_from_data(keys[i]);
        if (pt.coords.y < lohi.lo.coords.y) lohi.lo.coords.y = pt.coords.y;
        if (pt.coords.y > lohi.hi.coords.y) lohi.hi.coords.y = pt.coords.y;
        if (pt.coords.x < lohi.lo.coords.x) lohi.lo.coords.x = pt.coords.x;
        if (pt.coords.x > lohi.hi.coords.x) lohi.hi.coords.x = pt.coords.x;
    }
    free(keys);
    return lohi;
}
