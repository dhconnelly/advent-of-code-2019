Advent of Code 2019
===================

Go solutions to [Advent of Code](https://adventofcode.com) problems for
2019. Written by [Daniel Connelly](https://dhconnelly.com). There's also a
corresponding
[blog post](https://dhconnelly.com/advent-of-code-2019-commentary.html).

To build everything:

    go build

To run the test cases:

    go test

The input files that were generated for me are stored in the `input/`
subdirectory. The input file for day N is named `dayN`. If it's named
differently, the program won't be able to find it!

To run all solutions over all input files stored in the directory `input`:

    go build
    ./advent-of-code-2019 input/

To run just one day (e.g. day 7):

    ./advent-of-code-2019 input/ day7

Day specs are defined in day_specs.go. These specify the input filename, the
test cases for each part of the problem, and the functions that solve parts 1
and 2 of the problem. Each function implements the interface `func (io.Reader)
string`.

Licensed under the MIT License, if for some reason that's interesting to you.
