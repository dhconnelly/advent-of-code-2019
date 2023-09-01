#include <stdio.h>

extern void day2_test(void);
extern void day5_test(void);

int main() {
    printf("running all tests...\n");
    day2_test();
    day5_test();
    printf("done.\n");
}
