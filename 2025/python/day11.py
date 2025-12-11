# # Read graph
# graph = {}

# with open("inputday11.txt") as f:
#     for line in f:
#         line = line.strip()
#         if not line:
#             continue
#         left, right = line.split(":")
#         src = left.strip()
#         outs = right.strip().split()
#         graph[src] = outs

# # DFS with memoization
# memo = {}

# def count_paths(node):
#     if node == "out":
#         return 1
#     if node not in graph or len(graph[node]) == 0:
#         return 0
#     if node in memo:
#         return memo[node]

#     total = 0
#     for nxt in graph[node]:
#         total += count_paths(nxt)

#     memo[node] = total
#     return total

# print(count_paths("you"))


# Read graph
graph = {}

with open("inputday11.txt") as f:
    for line in f:
        line = line.strip()
        if not line:
            continue
        left, right = line.split(":")
        src = left.strip()
        outs = right.strip().split()
        graph[src] = outs

SPECIAL1 = "dac"
SPECIAL2 = "fft"

from functools import lru_cache

@lru_cache(None)
def dfs(node, seen_dac, seen_fft):
    # Update special-node flags
    if node == SPECIAL1:
        seen_dac = True
    if node == SPECIAL2:
        seen_fft = True

    # If we reached OUT, check if conditions satisfied
    if node == "out":
        return 1 if (seen_dac and seen_fft) else 0

    # If node has no outgoing edges
    if node not in graph or len(graph[node]) == 0:
        return 0

    total = 0
    for nxt in graph[node]:
        total += dfs(nxt, seen_dac, seen_fft)

    return total

answer = dfs("svr", False, False)
print(answer)
