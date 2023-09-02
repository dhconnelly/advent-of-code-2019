#include "hash.h"

#include <assert.h>
#include <stdio.h>

void hash_test(void) {
    printf("test_hash\n");

    hashtable table;
    init_table(table);

    assert(table_get(table, 0) == NULL);
    assert(table_get(table, -17) == NULL);

    table_set(table, -17, 12345);
    table_set(table, 0, 67890);

    assert(*table_get(table, -17) == 12345);
    assert(*table_get(table, 0) == 67890);
    assert(table_get(table, 1) == NULL);

    table_set(table, 0, 19);
    assert(table_get(table, 17) == NULL);
    assert(*table_get(table, 0) == 19);
}
