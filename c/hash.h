#ifndef HASH_H_
#define HASH_H_

#include <stdint.h>

#define HASHSIZE 101

typedef struct hashnode {
    int key;
    int64_t val;
    struct hashnode* next;
} hashnode;

typedef hashnode* hashtable[HASHSIZE];

void init_table(hashtable table);
void fill_table(hashtable table, int64_t arr[], int len);
int64_t* table_get(hashtable table, int key);
void table_set(hashtable table, int key, int64_t val);

#endif  // HASH_H_
