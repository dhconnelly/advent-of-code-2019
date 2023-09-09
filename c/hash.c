#include "hash.h"

#include <assert.h>
#include <stdio.h>
#include <stdlib.h>

void init_table(hashtable* table) {
    table->size = 0;
    for (int i = 0; i < HASHSIZE; i++) table->table[i] = NULL;
}

hashtable make_table(void) {
    hashtable table;
    init_table(&table);
    return table;
}

hashtable* new_table(void) {
    hashtable* table = malloc(sizeof(hashtable));
    init_table(table);
    return table;
}

void fill_table(hashtable* table, int64_t arr[], int len) {
    init_table(table);
    for (int i = 0; i < len; i++) table_set(table, i, arr[i]);
}

void table_copy(hashtable* into, const hashtable* from) {
    init_table(into);
    // TODO: optimize if needed by just cloning the lists instead of instering
    for (int i = 0; i < HASHSIZE; i++) {
        for (const hashnode* node = from->table[i]; node != NULL;
             node = node->next) {
            table_set(into, node->key, node->val);
        }
    }
}

static uint16_t hash(uint32_t key) {
    uint16_t idx;
    for (idx = 0; key != 0; key /= 10) idx = (key % 10) + 31 * idx;
    return idx % HASHSIZE;
}

int64_t* table_get(const hashtable* table, uint32_t key) {
    for (hashnode* node = table->table[hash(key)]; node != NULL;
         node = node->next) {
        if (node->key == key) {
            return &node->val;
        }
    }
    return NULL;
}

void table_set(hashtable* table, uint32_t key, int64_t val) {
    hashnode** node;
    for (node = &table->table[hash(key)]; *node != NULL && (*node)->key != key;
         node = &(*node)->next)
        ;
    if (*node == NULL) {
        *node = malloc(sizeof(hashnode));
        (*node)->key = key;
        (*node)->val = val;
        (*node)->next = NULL;
        table->size++;
    } else {
        assert((*node)->key == key);
        (*node)->val = val;
    }
}

int table_size(hashtable* table) { return table->size; }

uint32_t* table_keys(hashtable* table) {
    uint32_t* keys = malloc(table->size * sizeof(uint32_t));
    int i = 0;
    for (int bucket = 0; bucket < HASHSIZE; bucket++) {
        for (hashnode* node = table->table[bucket]; node != NULL;
             node = node->next) {
            keys[i++] = node->key;
        }
    }
    return keys;
}
