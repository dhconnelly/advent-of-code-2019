#include <assert.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

#include "parse.h"
#include "vm.h"

static int readint(vm* vm) {
    assert(vm->state == VM_OUTPUT);
    int output = vm->output;
    run(vm);
    return output;
}

enum {
    TILE_EMPTY = 0,
    TILE_WALL,
    TILE_BLOCK,
    TILE_PADDLE,
    TILE_BALL,
};

char tile(int of) {
    static char tiles[] = {
        [TILE_EMPTY] = ' ',   //
        [TILE_WALL] = '$',    //
        [TILE_BLOCK] = '#',   //
        [TILE_PADDLE] = '-',  //
        [TILE_BALL] = 'o',    //
    };
    return tiles[of];
}

static int part1(int64_t data[], int len) {
    vm vm = make_vm(data, len);
    run(&vm);
    int blocks = 0;
    do {
        int x = readint(&vm);
        int y = readint(&vm);
        int c = readint(&vm);
        if (c == TILE_BLOCK) blocks++;
    } while (vm.state != VM_HALTED);
    printf("%d\n", blocks);
    return EXIT_SUCCESS;
}

static void writeint(vm* vm, int val) {
    assert(vm->state == VM_INPUT);
    vm->input = val;
    run(vm);
}

static int sign(int x) {
    if (x < 0) return -1;
    if (x > 0) return 1;
    return 0;
}

static int part2(int64_t data[], int len) {
    vm vm = make_vm(data, len);
    set_mem(&vm, 0, 2);
    run(&vm);
    int ball_x = -1, paddle_x = -1;
    int score = 0;
    do {
        switch (vm.state) {
            case VM_INPUT:
                writeint(&vm, ball_x != -1 ? sign(ball_x - paddle_x) : 0);
                break;
            case VM_OUTPUT: {
                int x = readint(&vm);
                int y = readint(&vm);
                int c = readint(&vm);
                if (x == -1 && y == 0)
                    score = c;
                else if (c == TILE_PADDLE)
                    paddle_x = x;
                else if (c == TILE_BALL)
                    ball_x = x;
                break;
            }
            default:
                assert(0);
        }
    } while (vm.state != VM_HALTED);
    printf("%d\n", score);
    return EXIT_SUCCESS;
}

int main(int argc, char* argv[]) {
    if (argc != 2) {
        fprintf(stderr, "usage: day13 file\n");
        exit(1);
    }

    FILE* f = fopen(argv[1], "r");
    if (f == NULL) {
        perror("day13");
        exit(1);
    }

    int64_t data[4096];
    int len;
    if ((len = parse_intcode(f, data, 4096)) < 0) {
        perror("day13: failed to parse intcode");
        exit(1);
    }

    part1(data, len);
    part2(data, len);
}
