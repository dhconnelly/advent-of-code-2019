#include "queue.h"

#include <assert.h>
#include <stdlib.h>

void init_q(queue* q) { q->head = q->tail = NULL; }

int empty_q(queue* q) { return q->head == NULL; }

static node* new_node(void* data) {
    node* nd = malloc(sizeof(node));
    nd->data = data;
    nd->next = NULL;
    return nd;
}

void append_q(queue* q, void* data) {
    if (empty_q(q)) {
        q->head = q->tail = new_node(data);
    } else {
        q->tail->next = new_node(data);
        q->tail = q->tail->next;
    }
}

void* pop_q(queue* q) {
    assert(!empty_q(q));
    node* head = q->head;
    void* data = head->data;
    q->head = head->next;
    free(head);
    if (empty_q(q)) q->tail = NULL;
    return data;
}
