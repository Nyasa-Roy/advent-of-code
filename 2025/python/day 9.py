# from itertools import combinations

# def read_points(filename):
#     points = []
#     with open(filename, "r") as f:
#         for line in f:
#             line = line.strip()
#             if line:
#                 x, y = map(int, line.split(","))
#                 points.append((x, y))
#     return points


# def largest_rectangle_area(points):
#     max_area = 0

#     for (x1, y1), (x2, y2) in combinations(points, 2):
#         width = abs(x1 - x2) + 1
#         height = abs(y1 - y2) + 1
#         area = width * height
#         max_area = max(max_area, area)

#     return max_area


# if __name__ == "__main__":
#     points = read_points("inputday9.txt")
#     print(largest_rectangle_area(points))


from itertools import combinations
from collections import deque

def read_points(filename):
    points = []
    with open(filename) as f:
        for line in f:
            x, y = map(int, line.strip().split(","))
            points.append((x, y))
    return points


def build_allowed_tiles(red_points):
    red = set(red_points)
    green = set()

    n = len(red_points)
    for i in range(n):
        x1, y1 = red_points[i]
        x2, y2 = red_points[(i + 1) % n]

        if x1 == x2:
            for y in range(min(y1, y2), max(y1, y2) + 1):
                if (x1, y) not in red:
                    green.add((x1, y))
        else:
            for x in range(min(x1, x2), max(x1, x2) + 1):
                if (x, y1) not in red:
                    green.add((x, y1))

    boundary = red | green

    # Coordinate compression
    xs = sorted({x for x, y in boundary})
    ys = sorted({y for x, y in boundary})

    x_index = {x: i for i, x in enumerate(xs)}
    y_index = {y: i for i, y in enumerate(ys)}

    W, H = len(xs), len(ys)
    grid = [[0] * H for _ in range(W)]

    for x, y in boundary:
        grid[x_index[x]][y_index[y]] = 1  # wall

    # Flood fill from outside
    visited = [[False] * H for _ in range(W)]
    queue = deque()

    for i in range(W):
        queue.append((i, 0))
        queue.append((i, H - 1))
    for j in range(H):
        queue.append((0, j))
        queue.append((W - 1, j))

    while queue:
        x, y = queue.popleft()
        if not (0 <= x < W and 0 <= y < H):
            continue
        if visited[x][y] or grid[x][y] == 1:
            continue
        visited[x][y] = True
        for dx, dy in [(1,0),(-1,0),(0,1),(0,-1)]:
            queue.append((x + dx, y + dy))

    allowed = set(boundary)
    for x in range(W):
        for y in range(H):
            if not visited[x][y]:
                allowed.add((xs[x], ys[y]))

    return red, allowed


def largest_rectangle(red, allowed):
    best = 0
    for (x1, y1), (x2, y2) in combinations(red, 2):
        minx, maxx = sorted([x1, x2])
        miny, maxy = sorted([y1, y2])

        valid = True
        for x in range(minx, maxx + 1):
            for y in range(miny, maxy + 1):
                if (x, y) not in allowed:
                    valid = False
                    break
            if not valid:
                break

        if valid:
            area = (maxx - minx + 1) * (maxy - miny + 1)
            best = max(best, area)

    return best


if __name__ == "__main__":
    red_points = read_points("inputday9.txt")
    red, allowed = build_allowed_tiles(red_points)
    print(largest_rectangle(red, allowed))
