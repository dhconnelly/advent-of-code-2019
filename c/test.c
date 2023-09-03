#include <stdio.h>

extern void day2_test(void);
extern void day5_test(void);
extern void day7_test(void);
extern void hash_test(void);
extern void day9_test(void);
extern void pt2_test(void);

int main() {
    printf("running all tests...\n");
    day2_test();
    day5_test();
    day7_test();
    hash_test();
    day9_test();
    pt2_test();
    printf("done.\n");
}
