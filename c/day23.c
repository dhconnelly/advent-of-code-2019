#include <assert.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

#include "parse.h"
#include "vm.h"

#define NUM_VMS 50
#define MAX_PACKETS 100

typedef struct {
    int64_t x, y;
} packet;

typedef struct {
    vm vm;
    packet q[MAX_PACKETS];
    int qp;
} netvm;

void init_netvm(netvm* vm, int addr, int64_t data[], int data_len) {
    load_vm(&vm->vm, data, data_len);
    for (int i = 0; i < MAX_PACKETS; i++) {
        vm->q[i].x = 0;
        vm->q[i].y = 0;
    }
    vm->qp = -1;
    run(&vm->vm);
    assert(vm->vm.state == VM_INPUT);
    vm->vm.input = addr;
    run(&vm->vm);
}

int network_idle(netvm vms[]) {
    for (int i = 0; i < NUM_VMS; i++)
        if (vms[i].vm.state != VM_INPUT || vms[i].qp >= 0) return 0;
    return 1;
}

int64_t part1(int64_t data[], int data_len, int enable_nat) {
    // initialize and prime the vms
    netvm vms[NUM_VMS];
    packet nat;
    int prev_idle = 0;
    int64_t last_y = -1;
    for (int i = 0; i < NUM_VMS; i++) init_netvm(&vms[i], i, data, data_len);

    for (int j = 0;; j++) {
        for (int i = 0; i < NUM_VMS; i++) {
            if (network_idle(vms)) {
                if (j != 0 && prev_idle && nat.y == last_y) return last_y;
                vms[0].qp = 0;
                vms[0].q[0].x = nat.x;
                vms[0].q[0].y = nat.y;
                last_y = nat.y;
                prev_idle = 1;
            }

            switch (vms[i].vm.state) {
                case VM_INPUT: {
                    if (vms[i].qp < 0) {
                        vms[i].vm.input = -1;
                        run(&vms[i].vm);
                    } else {
                        packet packet = vms[i].q[vms[i].qp--];
                        vms[i].vm.input = packet.x;
                        run(&vms[i].vm);
                        assert(vms[i].vm.state == VM_INPUT);
                        vms[i].vm.input = packet.y;
                        run(&vms[i].vm);
                    }
                    break;
                }

                case VM_OUTPUT: {
                    int64_t addr = vms[i].vm.output;
                    run(&vms[i].vm);
                    assert(vms[i].vm.state == VM_OUTPUT);
                    int64_t x = vms[i].vm.output;
                    run(&vms[i].vm);
                    assert(vms[i].vm.state == VM_OUTPUT);
                    int64_t y = vms[i].vm.output;
                    run(&vms[i].vm);
                    if (addr == 255) {
                        if (!enable_nat) return y;
                        nat.x = x;
                        nat.y = y;
                    } else {
                        assert(++vms[addr].qp < MAX_PACKETS - 1);
                        vms[addr].q[vms[addr].qp].x = x;
                        vms[addr].q[vms[addr].qp].y = y;
                    }
                    break;
                }

                default:
                    assert(0);
            }
        }
    }

    assert(0);
}

int main(int argc, char* argv[]) {
    if (argc != 2) {
        fprintf(stderr, "usage: day23 file\n");
        exit(1);
    }

    FILE* f = fopen(argv[1], "r");
    if (f == NULL) {
        perror("day23");
        exit(1);
    }

    int64_t data[4096];
    int len;
    if ((len = parse_intcode(f, data, 4096)) < 0) {
        perror("day23: failed to parse intcode");
        exit(1);
    }

    printf("%ld\n", part1(data, len, 0));
    printf("%ld\n", part1(data, len, 1));
}
