#include <stdio.h>
#include <stdlib.h>

#include "parse.h"
#include "vm.h"

int main(int argc, char* argv[]) {
    if (argc != 2) {
        printf("usage: day2 file\n");
        exit(1);
    }

    FILE* f = fopen(argv[1], "r");
    if (f == NULL) {
        perror("day2");
        exit(1);
    }

    vm vm = new_vm();
    if ((vm.mem_size = parse_intcode(f, vm.mem, MAX_MEM)) < 0) {
        perror("day2");
        exit(1);
    }

    vm.mem[1] = 12;
    vm.mem[2] = 2;
    run(&vm);
    printf("%lld\n", vm.mem[0]);

    return 0;
}
