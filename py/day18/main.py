import collections


def parse(f):
    return [[ch for ch in line.strip()] for line in f]


DIRS = [(-1, 0), (1, 0), (0, -1), (0, 1)]


def dists(maze, row, col):
    d = {}
    v = {(row, col)}
    q = collections.deque()
    q.append(((row, col), 0))
    while q:
        ((row, col), dist) = q.popleft()
        for (dr, dc) in DIRS:
            nr, nc = row + dr, col + dc
            if nr < 0 or nr >= len(maze) or nc < 0 or nc >= len(maze[nr]):
                continue
            if (nr, nc) in v:
                continue
            v.add((nr, nc))
            ch = maze[nr][nc]
            if ch not in ('.', '#'):
                d[ch] = dist+1
            if ch in ('.', '@') or ch.islower():
                q.append(((nr, nc), dist+1))
    return d


def cache_key(ch, keys):
    return ch + str(keys)


def connect_across(g, ch):
    nbrs = g[ch]
    nbr_keys = list(nbrs.keys())
    for (i, x) in enumerate(nbr_keys):
        for y in nbr_keys[i+1:]:
            dist = g[x][ch] + g[ch][y]
            if y not in g[x]:
                g[x][y] = dist
                g[y][x] = dist
            else:
                g[x][y] = min(g[x][y], dist)
                g[y][x] = min(g[y][x], dist)


def collect_key(g, key):
    g2 = {k: v.copy() for (k, v) in g.items()}
    connect_across(g2, key)
    if (door := key.upper()) in g:
        connect_across(g2, door)
    return g2


def bit_in(ch, bits):
    i = ord(ch) - ord('a')
    return ((1 << i) & bits) > 0


def bit_remove(ch, bits):
    i = ord(ch) - ord('a')
    return (~(1 << i)) & bits


def bit_make(chs):
    bits = 0
    for ch in chs:
        bits = bits | (1 << (ord(ch) - ord('a')))
    return bits


def remove(g, key):
    for nbr in g[key]:
        del g[nbr][key]
    del g[key]
    if key.islower() and (door := key.upper()) in g:
        for nbr in g[door]:
            del g[nbr][door]
        del g[door]


def collect_all(g, from_key, want_keys, cache):
    if want_keys == 0:
        return 0
    ck = cache_key(from_key, want_keys)
    if ck in cache:
        return cache[ck]
    min_steps = None
    dists = g[from_key]
    remove(g, from_key)
    for key in (ch for ch in dists if ch.islower() and bit_in(ch, want_keys)):
        want_keys2 = bit_remove(key, want_keys)
        g2 = collect_key(g, key)
        steps = dists[key] + collect_all(g2, key, want_keys2, cache)
        if min_steps is None or steps < min_steps:
            min_steps = steps
    cache[ck] = min_steps
    return min_steps


def all_dists(maze):
    g = {}
    for (i, row) in enumerate(maze):
        for (j, ch) in enumerate(row):
            if ch not in ('.', '#'):
                g[ch] = dists(maze, i, j)
    return g


def main(args):
    with open(args[1]) as f:
        maze = parse(f)
    dists = all_dists(maze)
    want_keys = bit_make({ch for ch in dists if ch.islower()})
    print(collect_all(dists, '@', want_keys, {}))


if __name__ == '__main__':
    import sys
    main(sys.argv)
