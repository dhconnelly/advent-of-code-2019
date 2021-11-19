from dataclasses import dataclass


def grid(text: str) -> list[list[str]]:
    return [list(line) for line in text.splitlines()]


@dataclass
class Portal:
    label: str
    is_outer: bool
    tile: tuple[int, int]


@dataclass
class Maze:
    tiles: list[list[str | Portal]]
    links: dict[str, list[Portal]]


@dataclass
class Boundaries:
    top: int
    bottom: int
    left: int
    right: int


def find_boundaries(grid: list[list[str]]) -> Boundaries:
    top, bottom, left, right = -1, -1, -1, -1
    for row in range(len(grid)):
        if any(ch == "#" for ch in grid[row]):
            top = row
            break
    for row in range(row, len(grid)):
        if not any(ch == "#" for ch in grid[row]):
            bottom = row - 1
            break
    for col in range(len(grid[0])):
        if any(row[col] == "#" for row in grid):
            left = col
            break
    for col in range(col, len(grid[0])):
        if not any(row[col] == "#" for row in grid):
            right = col - 1
            break
    return Boundaries(top, bottom, left, right)


def tile(
    g: list[list[str]],
    bounds: Boundaries,
    row: int,
    col: int,
) -> str | Portal:
    # for empty, walls, and floor, just return that character
    ch = g[row][col]
    if ch in (" ", "#", "."):
        return ch

    # determine if this is a portal label that borders the maze
    # outside portals
    if row == bounds.top - 1:
        return Portal(g[row - 1][col] + g[row][col], True, (row + 1, col))
    elif row == bounds.bottom + 1:
        return Portal(g[row][col] + g[row + 1][col], True, (row - 1, col))
    elif col == bounds.left - 1:
        return Portal(g[row][col - 1] + g[row][col], True, (row, col + 1))
    elif col == bounds.right + 1:
        return Portal(g[row][col] + g[row][col + 1], True, (row, col - 1))
    # inside portals
    if row > 0 and g[row - 1][col] == ".":
        return Portal(g[row][col] + g[row + 1][col], False, (row - 1, col))
    elif row < len(g) - 1 and g[row + 1][col] == ".":
        return Portal(g[row - 1][col] + g[row][col], False, (row + 1, col))
    elif col < len(g[0]) - 1 and g[row][col + 1] == ".":
        return Portal(g[row][col - 1] + g[row][col], False, (row, col + 1))
    elif col > 0 and g[row][col - 1] == ".":
        return Portal(g[row][col] + g[row][col + 1], False, (row, col - 1))

    # just call it empty
    return " "


def find_links(tiles: list[list[str | Portal]]) -> dict[str, list[Portal]]:
    links: dict[str, list[Portal]] = {}
    for row in tiles:
        for tile in row:
            if isinstance(tile, Portal):
                links.setdefault(tile.label, []).append(tile)
    return links


def maze(grid: list[list[str]]) -> Maze:
    bounds = find_boundaries(grid)
    tiles = [
        [tile(grid, bounds, row, col) for col in range(len(grid[0]))]
        for row in range(len(grid))
    ]
    links = find_links(tiles)
    return Maze(tiles, links)


def main(args: list[str]):
    with open(args[0]) as f:
        s = f.read()
        g = grid(s)
        m = maze(g)
        print(m)


if __name__ == "__main__":
    import sys

    main(sys.argv[1:])
