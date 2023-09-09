#ifndef HASH_H_
#define HASH_H_

#include <stdint.h>

#define HASHSIZE 101

typedef struct hashnode {
    uint32_t key;
    int64_t val;
    struct hashnode* next;
} hashnode;

typedef struct {
    hashnode* table[HASHSIZE];
    int size;
} hashtable;

void init_table(hashtable* table);
void fill_table(hashtable* table, int64_t arr[], int len);
int64_t* table_get(const hashtable* table, uint32_t key);
void table_set(hashtable* table, uint32_t key, int64_t val);
void table_copy(hashtable* into, const hashtable* from);
int table_size(hashtable* table);
uint32_t* table_keys(hashtable* table);
hashtable make_table(void);
hashtable* new_table(void);

#endif  // HASH_H_
