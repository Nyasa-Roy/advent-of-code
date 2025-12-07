import sys
sys.setrecursionlimit(1000000)


with open("inputday7.txt") as f:
    grid = [list(line.rstrip("\n")) for line in f]

R = len(grid)
C = len(grid[0])


start_r = start_c = None
for r in range(R):
    for c in range(C):
        if grid[r][c] == "S":
            start_r, start_c = r, c
            break
    if start_r is not None:
        break

memo = {}

def dfs(r, c):
    
    if c < 0 or c >= C:
        return 1

    
    if r >= R - 1:
        return 1

    if (r, c) in memo:
        return memo[(r, c)]

    below = grid[r + 1][c]

    
    if below == "^":
        total = dfs(r, c - 1) + dfs(r, c + 1)
    else:
        
        total = dfs(r + 1, c)

    memo[(r, c)] = total
    return total

answer = dfs(start_r, start_c)
print(answer)
