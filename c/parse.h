#ifndef PARSE_H_
#define PARSE_H_

#include <stdio.h>

// reads a comma-separated list of at most |max| integers from |f| into the
// array |result|. returns the number of integers read if successful, or -1
// and sets |errno| otherwise.
int parse_intcode(FILE* f, int result[], int max);

#endif  // PARSE_H_
