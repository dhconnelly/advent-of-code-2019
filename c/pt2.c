#include "pt2.h"

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
