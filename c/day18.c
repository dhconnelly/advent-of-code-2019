#include <assert.h>
#include <ctype.h>
#include <limits.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

#include "hash.h"
#include "pt2.h"
#include "queue.h"

typedef enum {
    ENTRANCE = '@',
    WALL = '#',
    PASSAGE = '.',
    KEY = 'K',
    DOOR = 'D',
} tile_type;

typedef struct {
    tile_type typ;
    int ch;
} tile;

#define MAX_ROWS 100
#define MAX_COLS 100

typedef struct {
    tile g[MAX_ROWS][MAX_COLS];
    int rows, cols;
} grid;

static void init_grid(grid* g) { g->rows = g->cols = 0; }

static void parse_grid(FILE* f, grid* g) {
    int col = 0;
    tile t;
    while ((t.ch = getc(f)) != EOF && !ferror(f)) {
        if (t.ch == '\n') {
            ++g->rows;
            col = 0;
            continue;
        }
        switch (t.ch) {
            case ENTRANCE:
            case WALL:
            case PASSAGE:
                t.typ = t.ch;
                break;
            default:
                if (!isalpha(t.ch)) {
                    fprintf(stderr, "day18: invalid tile: %c\n", t.ch);
                    exit(EXIT_FAILURE);
                }
                t.typ = islower(t.ch) ? KEY : DOOR;
                break;
        }
        g->g[g->rows][col] = t;
        g->cols = ++col;
    }
    if (col != 0) ++g->rows;
    if (ferror(f)) {
        perror("day18: can't parse grid");
        exit(EXIT_FAILURE);
    }
}

typedef struct {
    pt2 pos;
    int dist;
} bfs_node;

static bfs_node* make_node(pt2 pos, int dist) {
    bfs_node* node = malloc(sizeof(bfs_node));
    node->pos = pos;
    node->dist = dist;
    return node;
}

static char at(grid* g, pt2 from) {
    return g->g[from.coords.y][from.coords.x].ch;
}

#define ADJ_LEN (26 + 26 + 1)

static int adjindex(char ch) {
    if (ch >= 'a' && ch <= 'z') return ch - 'a';
    if (ch >= 'A' && ch <= 'Z') return ch - 'A' + 26;
    return 52;
}

static char adjval(int idx) {
    if (idx >= 0 && idx < 26) return 'a' + idx;
    if (idx >= 26 && idx < 52) return 'A' + idx - 26;
    return '@';
}

typedef struct {
    int d[ADJ_LEN][ADJ_LEN];
} adjmat;

static void clone(adjmat* into, adjmat* adj) {
    for (int i = 0; i < ADJ_LEN; i++) {
        for (int j = i; j < ADJ_LEN; j++) {
            into->d[i][j] = adj->d[i][j];
            into->d[j][i] = adj->d[j][i];
        }
    }
}

static adjmat empty_adjmat() {
    adjmat adj;
    for (int i = 0; i < ADJ_LEN; i++) {
        for (int j = 0; j < ADJ_LEN; j++) adj.d[i][j] = -1;
    }
    return adj;
}

static void shortest_dists(grid* g, pt2 from, adjmat* adj) {
    int* dists = adj->d[adjindex(at(g, from))];
    queue q = make_q();
    append_q(&q, make_node(from, 0));
    hashtable v = make_table();
    table_set(&v, from.data, 0);
    pt2 nbrs[4];
    while (!empty_q(&q)) {
        bfs_node* front = pop_q(&q);
        get_nbrs(front->pos, nbrs);
        for (int i = 0; i < 4; i++) {
            pt2 nbr = nbrs[i];
            if (table_get(&v, nbr.data)) continue;
            table_set(&v, nbr.data, 1);
            tile t = g->g[nbr.coords.y][nbr.coords.x];
            int dist = front->dist + 1;
            // TODO: don't care about distance to doors
            if (t.typ == KEY || t.typ == DOOR)
                dists[adjindex(at(g, nbr))] = dist;
            if (t.typ == PASSAGE || t.typ == ENTRANCE)
                append_q(&q, make_node(nbr, dist));
        }
        free(front);
    }
}

static adjmat all_dists(grid* g) {
    adjmat adj = empty_adjmat();
    for (int i = 0; i < g->rows; i++) {
        for (int j = 0; j < g->cols; j++) {
            pt2 from = make_pt(j, i);
            tile t = g->g[i][j];
            if (t.typ == ENTRANCE || t.typ == DOOR || t.typ == KEY)
                shortest_dists(g, from, &adj);
        }
    }
    return adj;
}

void collect(adjmat* adj, int idx) {
    int id, jd;
    for (int i = 0; i < ADJ_LEN - 1; i++) {
        if (i == idx || (id = adj->d[idx][i]) < 0) continue;
        for (int j = i + 1; j < ADJ_LEN; j++) {
            if (j == idx || (jd = adj->d[idx][j]) < 0) continue;
            if (adj->d[i][j] < 0 || id + jd < adj->d[i][j]) {
                adj->d[i][j] = adj->d[j][i] = id + jd;
            }
        }
    }
    for (int i = 0; i < ADJ_LEN; i++) {
        adj->d[i][idx] = -1;
    }
}

uint32_t all_keys(grid* g) {
    uint32_t keys = 0;
    for (int row = 0; row < g->rows; row++) {
        for (int col = 0; col < g->cols; col++) {
            if (g->g[row][col].typ == KEY)
                keys |= (1 << adjindex(g->g[row][col].ch));
        }
    }
    return keys;
}

static int is_key(int idx) { return islower(adjval(idx)); }
static int is_door(int idx) { return isupper(adjval(idx)); }
static int need_key(int door_idx, uint32_t keys_needed) {
    char door = adjval(door_idx);
    char key = 'a' + (door - 'A');
    int key_idx = adjindex(key);
    return ((1 << key_idx) & keys_needed) > 0;
}
static uint32_t collect_key(uint32_t keys_needed, int key_idx) {
    return keys_needed & ~(1 << key_idx);
}

static uint64_t memo_key(uint32_t l, uint32_t r) {
    return (((uint64_t)l) << 32) | (uint64_t)r;
}

static int64_t collect_all(adjmat* adj, uint32_t from_idx, uint32_t keys_needed,
                           hashtable* memo) {
    if (keys_needed == 0) return 0;
    uint64_t key = memo_key(from_idx, keys_needed);
    int64_t* memo_val = table_get(memo, key);
    if (memo_val != NULL) return *memo_val;
    adjmat scratch;
    int64_t min = INT64_MAX;
    for (int64_t d, i = 0; i < ADJ_LEN; i++) {
        if (i == from_idx || (d = adj->d[from_idx][i]) < 0) continue;
        if (is_door(i) && need_key(i, keys_needed)) continue;
        clone(&scratch, adj);
        collect(&scratch, i);
        int64_t sub_dist = collect_all(
            &scratch, i, is_key(i) ? collect_key(keys_needed, i) : keys_needed,
            memo);
        if (sub_dist < 0) continue;
        int64_t total_dist = sub_dist + d;
        if (total_dist < min) min = total_dist;
    }
    int64_t result = (min == INT64_MAX) ? -1 : min;
    table_set(memo, key, result);
    return result;
}

int main(int argc, char* argv[]) {
    if (argc != 2) {
        fprintf(stderr, "usage: day18 file\n");
        exit(EXIT_FAILURE);
    }

    FILE* f = fopen(argv[1], "r");
    if (f == NULL) {
        perror("day18");
        exit(EXIT_FAILURE);
    }

    grid g;
    init_grid(&g);
    parse_grid(f, &g);
    adjmat dists = all_dists(&g);
    hashtable memo = make_table();
    printf("%lld\n", collect_all(&dists, adjindex('@'), all_keys(&g), &memo));
}
