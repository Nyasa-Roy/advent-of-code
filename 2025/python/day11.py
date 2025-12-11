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

# DFS with memoization
memo = {}

def count_paths(node):
    if node == "out":
        return 1
    if node not in graph or len(graph[node]) == 0:
        return 0
    if node in memo:
        return memo[node]

    total = 0
    for nxt in graph[node]:
        total += count_paths(nxt)

    memo[node] = total
    return total

print(count_paths("you"))
