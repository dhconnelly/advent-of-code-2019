#include <assert.h>
#include <ctype.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

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

static void print_grid(grid* g) {
    for (int i = 0; i < g->rows; i++) {
        for (int j = 0; j < g->cols; j++) putchar(g->g[i][j].ch);
        putchar('\n');
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

static void dists(grid* g, pt2 from, hashtable* dists) {
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
            int d = front->dist + 1;
            if (t.typ == KEY || t.typ == DOOR) table_set(dists, at(g, nbr), d);
            if (t.typ == PASSAGE) append_q(&q, make_node(nbr, d));
        }
        free(front);
    }
}

static hashtable all_dists(grid* g) {
    hashtable table = make_table();
    for (int i = 0; i < g->rows; i++) {
        for (int j = 0; j < g->cols; j++) {
            pt2 from = make_pt(j, i);
            tile t = g->g[i][j];
            hashtable* d = new_table();
            if (t.typ == ENTRANCE || t.typ == DOOR || t.typ == KEY)
                dists(g, from, d);
            table_set(&table, at(g, from), (int64_t)d);  // hmmm
        }
    }
    return table;
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
    print_grid(&g);

    hashtable dists = all_dists(&g);
    uint32_t* keys = table_keys(&dists);
    for (int i = 0; i < table_size(&dists); i++) {
        hashtable* dest_dists = (hashtable*)*table_get(&dists, keys[i]);
        uint32_t* dest_keys = table_keys(dest_dists);
        for (int j = 0; j < table_size(dest_dists); j++) {
            int dist = *table_get(dest_dists, dest_keys[j]);
            printf("dist(%c, %c) = %d\n", keys[i], dest_keys[j], dist);
        }
    }
}
