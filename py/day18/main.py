import collections


def parse(f):
    return [[ch for ch in line.strip()] for line in f]


DIRS = [(-1, 0), (1, 0), (0, -1), (0, 1)]


def dists(maze, targets, held_keys, row, col):
    d = {}
    v = set()
    v.add((row, col))
    q = collections.deque()
    q.append(((row, col), 0))
    while len(q) > 0:
        (row, col), dist = q.popleft()
        if maze[row][col] in targets:
            d[maze[row][col]] = ((row, col), dist)
        for (dr, dc) in DIRS:
            nr, nc = row + dr, col + dc
            if nr < 0 or nr >= len(maze) or nc < 0 or nc >= len(maze[nr]):
                continue
            if (nr, nc) in v:
                continue
            ch = maze[nr][nc]
            if ch == '#' or ch.isupper() and ch.lower() not in held_keys:
                continue
            v.add((nr, nc))
            q.append(((nr, nc), dist+1))
    return d


def cache_key(keys, row, col):
    return f"{''.join(sorted(keys))},{row},{col}"


def collect(maze, want_keys, held_keys, row, col, cache={}):
    if len(want_keys) == 0:
        return 0
    ck = cache_key(held_keys, row, col)
    if ck in cache:
        return cache[ck]
    min_steps = None
    reachable = dists(maze, want_keys, held_keys, row, col)
    for key, ((row, col), dist) in reachable.items():
        want_keys.remove(key)
        held_keys.add(key)
        steps = dist + collect(maze, want_keys, held_keys, row, col)
        if min_steps is None or steps < min_steps:
            min_steps = steps
        want_keys.add(key)
        held_keys.remove(key)
    cache[ck] = min_steps
    return min_steps


def all_keys(maze):
    return {ch for row in maze for ch in row if ch.islower()}


def find_start(maze):
    for i, row in enumerate(maze):
        for j, ch in enumerate(row):
            if ch == '@':
                return (i, j)


def main(args):
    with open(args[1]) as f:
        maze = parse(f)
    (row, col) = find_start(maze)
    print(collect(maze, all_keys(maze), set(), row, col))


if __name__ == '__main__':
    import sys
    main(sys.argv)
