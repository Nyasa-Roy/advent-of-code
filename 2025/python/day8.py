# import math

# points = []
# with open("inputday8.txt") as f:
#     for line in f:
#         line = line.strip()
#         if not line:
#             continue
#         x, y, z = map(int, line.split(","))
#         points.append((x, y, z))

# n = len(points)

# edges = []
# for i in range(n):
#     x1, y1, z1 = points[i]
#     for j in range(i+1, n):
#         x2, y2, z2 = points[j]
#         d2 = (x1-x2)**2 + (y1-y2)**2 + (z1-z2)**2  # squared distance (good enough)
#         edges.append((d2, i, j))


# edges.sort(key=lambda x: x[0])




# parent = list(range(n))
# size = [1] * n

# def find(a):
#     while parent[a] != a:
#         parent[a] = parent[parent[a]]
#         a = parent[a]
#     return a

# def union(a, b):
#     ra, rb = find(a), find(b)
#     if ra == rb:
#         return
#     if size[ra] < size[rb]:
#         ra, rb = rb, ra
#     parent[rb] = ra
#     size[ra] += size[rb]




# for k in range(1000):
#     d2, i, j = edges[k]
#     union(i, j)




# from collections import Counter

# component_sizes = Counter()

# for i in range(n):
#     r = find(i)
#     component_sizes[r] += 1


# largest = sorted(component_sizes.values(), reverse=True)[:3]

# answer = largest[0] * largest[1] * largest[2]
# print(answer)


# Read input
points = []
with open("inputday8.txt") as f:
    for line in f:
        line = line.strip()
        if not line:
            continue
        x, y, z = map(int, line.split(","))
        points.append((x, y, z))

n = len(points)


edges = []
for i in range(n):
    x1, y1, z1 = points[i]
    for j in range(i+1, n):
        x2, y2, z2 = points[j]
        d2 = (x1-x2)**2 + (y1-y2)**2 + (z1-z2)**2
        edges.append((d2, i, j))

edges.sort(key=lambda x: x[0])


parent = list(range(n))
size = [1] * n
components = n

def find(a):
    while parent[a] != a:
        parent[a] = parent[parent[a]]
        a = parent[a]
    return a

def union(a, b):
    global components
    ra, rb = find(a), find(b)
    if ra == rb:
        return False
    if size[ra] < size[rb]:
        ra, rb = rb, ra
    parent[rb] = ra
    size[ra] += size[rb]
    return True


last_edge = None

for d2, i, j in edges:
    if union(i, j):
        components -= 1
        last_edge = (i, j)
        if components == 1:
            break


x1 = points[last_edge[0]][0]
x2 = points[last_edge[1]][0]

print(x1 * x2)
