#ifndef QUEUE_H_
#define QUEUE_H_

typedef struct node {
    void* data;
    struct node* next;
} node;

typedef struct {
    node *head, *tail;
} queue;

queue make_q(void);
void init_q(queue* q);
int empty_q(queue* q);
void append_q(queue* q, void* data);
void* pop_q(queue* q);

#endif  // QUEUE_H_
