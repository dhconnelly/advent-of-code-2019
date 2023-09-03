#include <assert.h>
#include <stdio.h>
#include <stdlib.h>

#include "hash.h"
#include "parse.h"
#include "pt2.h"
#include "vm.h"

typedef enum { NORTH, EAST, SOUTH, WEST } direction;
typedef enum { BLACK = 0, WHITE = 1 } color;
typedef enum { LEFT = 0, RIGHT = 1 } turn;

typedef struct {
    direction dir;
    pt2 pos;
    hashtable tiles;  // pt -> color
    vm vm;
} robot;

static void init_robot(robot* r, int64_t data[], int data_len) {
    r->dir = NORTH;
    r->pos = make_pt(0, 0);
    init_table(&r->tiles);
    table_set(&r->tiles, r->pos.data, BLACK);
    load_vm(&r->vm, data, data_len);
}

static color get_tile(robot* r) {
    int64_t* val = table_get(&r->tiles, r->pos.data);
    return val == NULL ? BLACK : *val;
}

static direction apply_turn(direction dir, turn t) {
    static direction dirs[4][2] = {
        [NORTH] = {[LEFT] = WEST, [RIGHT] = EAST},
        [EAST] = {[LEFT] = NORTH, [RIGHT] = SOUTH},
        [SOUTH] = {[LEFT] = EAST, [RIGHT] = WEST},
        [WEST] = {[LEFT] = SOUTH, [RIGHT] = NORTH},
    };
    return dirs[dir][t];
}

static void forward(robot* r) {
    switch (r->dir) {
        case NORTH:
            r->pos.coords.y--;
            break;
        case EAST:
            r->pos.coords.x++;
            break;
        case SOUTH:
            r->pos.coords.y++;
            break;
        case WEST:
            r->pos.coords.x--;
            break;
        default:
            fprintf(stderr, "bad direction: %d\n", r->dir);
            exit(1);
    }
}

void run_robot(robot* r) {
    run(&r->vm);

    do {
        // reads current tile
        assert(r->vm.state == VM_INPUT);
        r->vm.input = get_tile(r);
        run(&r->vm);

        // outputs color to paint
        assert(r->vm.state == VM_OUTPUT);
        table_set(&r->tiles, r->pos.data, r->vm.output);
        run(&r->vm);

        // outputs turn direction
        assert(r->vm.state == VM_OUTPUT);
        r->dir = apply_turn(r->dir, r->vm.output);
        forward(r);
        run(&r->vm);
    } while (r->vm.state != VM_HALTED && r->vm.state != VM_ERROR);

    if (r->vm.state == VM_ERROR) {
        fprintf(stderr, "vm error: %d\n", r->vm.error);
        exit(1);
    }
}

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

    robot r;
    init_robot(&r, data, len);
    run_robot(&r);
    printf("%d\n", r.tiles.size);
}
