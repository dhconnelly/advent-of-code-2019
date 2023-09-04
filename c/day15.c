#include <assert.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

#include "hash.h"
#include "parse.h"
#include "pt2.h"
#include "vm.h"

typedef struct {
    pt2 pos;
    hashtable map;
    vm vm;
} explorer;

typedef enum { WALL = 0, FLOOR = 1, OXYGEN = 2 } tile;

static void init_explorer(explorer* explorer, int64_t data[], int data_len) {
    explorer->pos.coords.x = 0;
    explorer->pos.coords.y = 0;
    init_table(&explorer->map);
    table_set(&explorer->map, explorer->pos.data, FLOOR);
    load_vm(&explorer->vm, data, data_len);
    run(&explorer->vm);
}

typedef enum { NORTH = 1, SOUTH = 2, WEST = 3, EAST = 4 } direction;

static direction dir(pt2 to, pt2 from) {
    int dy = to.coords.y - from.coords.y;
    int dx = to.coords.x - from.coords.x;
    if (dy > 0) return NORTH;
    if (dy < 0) return SOUTH;
    if (dx > 0) return EAST;
    if (dx < 0) return WEST;
    assert(0);
}

static tile move(vm* vm, direction dir) {
    assert(vm->state == VM_INPUT);
    vm->input = dir;
    run(vm);
    assert(vm->state == VM_OUTPUT);
    int response = vm->output;
    run(vm);
    return response;
}

static direction backwards(direction d) {
    static direction dirs[] = {
        [NORTH] = SOUTH,
        [SOUTH] = NORTH,
        [WEST] = EAST,
        [EAST] = WEST,
    };
    return dirs[d];
}

static int can_move(tile t) { return t == FLOOR || t == OXYGEN; }

static void dfs(explorer* explorer) {
    pt2 nbrs[4], nbr, cur = explorer->pos;
    get_nbrs(cur, nbrs);
    for (int i = 0; i < 4; i++) {
        nbr = nbrs[i];
        if (table_get(&explorer->map, nbr.data)) continue;
        direction d = dir(nbr, cur);
        tile t = move(&explorer->vm, d);
        table_set(&explorer->map, nbr.data, t);
        if (can_move(t)) {
            explorer->pos = nbr;
            dfs(explorer);
            explorer->pos = cur;
            assert(can_move(move(&explorer->vm, backwards(d))));
        }
    }
}

static pt2 find(hashtable* map, tile t) {
    uint32_t* keys = table_keys(map);
    for (int i = 0; i < table_size(map); i++) {
        if (*table_get(map, keys[i]) == t) {
            pt2 pos = pt_from_data(keys[i]);
            free(keys);
            return pos;
        }
    }
    assert(0);
}

static void print_map(hashtable* map, pt2 cur) {
    static int ch[] = {
        [FLOOR] = ' ',
        [WALL] = '#',
        [OXYGEN] = 'O',
    };
    rect lohi = bounds(map);
    for (int row = lohi.hi.coords.y; row >= lohi.lo.coords.y; row--) {
        for (int col = lohi.lo.coords.x; col <= lohi.hi.coords.x; col++) {
            pt2 pos = make_pt(col, row);
            if (pt_eq(pos, cur))
                putchar('D');
            else {
                int64_t* c = table_get(map, pos.data);
                int tile = (c == NULL) ? '?' : ch[*c];
                putchar(tile);
            }
        }
        putchar('\n');
    }
}

typedef struct node {
    pt2 pos;
    int dist;
    struct node* next;
} node;

node* new_node(pt2 pos, int dist) {
    node* nd = malloc(sizeof(node));
    nd->pos = pos;
    nd->dist = dist;
    nd->next = NULL;
    return nd;
}

typedef struct {
    node *head, *tail;
} queue;

void init_q(queue* q) { q->head = q->tail = NULL; }

void free_q(queue* q) {
    node* nd = q->head;
    while (nd != NULL) {
        node* next = nd->next;
        free(nd);
        nd = next;
    }
}

int empty_q(queue* q) { return q->head == NULL; }

void append_q(queue* q, pt2 pos, int dist) {
    if (empty_q(q)) {
        q->head = q->tail = new_node(pos, dist);
    } else {
        q->tail->next = new_node(pos, dist);
        q->tail = q->tail->next;
    }
}

node* pop_q(queue* q) {
    assert(!empty_q(q));
    node* front = q->head;
    q->head = q->head->next;
    if (empty_q(q)) q->tail = NULL;
    return front;
}

static hashtable bfs(hashtable* map, pt2 from) {
    hashtable dists;
    init_table(&dists);
    table_set(&dists, from.data, 0);
    queue q;
    init_q(&q);
    append_q(&q, from, 0);
    hashtable v;
    init_table(&v);
    table_set(&v, from.data, 1);
    pt2 nbrs[4];
    while (!empty_q(&q)) {
        node* front = pop_q(&q);
        get_nbrs(front->pos, nbrs);
        for (int i = 0; i < 4; i++) {
            if (table_get(&v, nbrs[i].data)) continue;
            int64_t* tile = table_get(map, nbrs[i].data);
            if (!tile || !can_move(*tile)) continue;
            table_set(&v, nbrs[i].data, 1);
            table_set(&dists, nbrs[i].data, front->dist + 1);
            append_q(&q, nbrs[i], front->dist + 1);
        }
        free(front);
    }
    free_q(&q);
    return dists;
}

static int64_t max(hashtable* table) {
    uint32_t* keys = table_keys(table);
    int64_t max = INT64_MIN;
    for (int i = 0; i < table_size(table); i++) {
        int64_t val = *table_get(table, keys[i]);
        if (val > max) max = val;
    }
    free(keys);
    return max;
}

int main(int argc, char* argv[]) {
    if (argc != 2) {
        fprintf(stderr, "usage: day15 file\n");
        exit(1);
    }

    FILE* f = fopen(argv[1], "r");
    if (f == NULL) {
        perror("day15");
        exit(1);
    }

    int64_t data[2048];
    int len;
    if ((len = parse_intcode(f, data, 2048)) < 0) {
        perror("day15: failed to parse intcode");
        exit(1);
    }

    explorer explorer;
    init_explorer(&explorer, data, len);
    dfs(&explorer);
    pt2 oxy = find(&explorer.map, OXYGEN);
    hashtable dists = bfs(&explorer.map, make_pt(0, 0));
    printf("%ld\n", *table_get(&dists, oxy.data));

    hashtable oxy_dists = bfs(&explorer.map, oxy);
    printf("%ld\n", max(&oxy_dists));
}
