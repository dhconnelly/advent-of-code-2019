#include "hash.h"

#include <assert.h>
#include <stdlib.h>

void init_table(hashtable table) {
    for (int i = 0; i < HASHSIZE; i++) table[i] = NULL;
}

void fill_table(hashtable table, int64_t arr[], int len) {
    init_table(table);
    for (int i = 0; i < len; i++) table_set(table, i, arr[i]);
}

static int hash(int key) {
    unsigned idx;
    for (idx = 0; key != 0; key /= 10) idx = (key % 10) + 31 * idx;
    return idx % HASHSIZE;
}

int64_t* table_get(hashtable table, int key) {
    for (hashnode* node = table[hash(key)]; node != NULL; node = node->next) {
        if (node->key == key) return &node->val;
    }
    return NULL;
}

void table_set(hashtable table, int key, int64_t val) {
    hashnode** node;
    for (node = &table[hash(key)]; *node != NULL && (*node)->key != key;
         node = &(*node)->next)
        ;
    if (*node == NULL) {
        *node = malloc(sizeof(hashnode));
        (*node)->key = key;
        (*node)->val = val;
        (*node)->next = NULL;
    } else {
        assert((*node)->key == key);
        (*node)->val = val;
    }
}
