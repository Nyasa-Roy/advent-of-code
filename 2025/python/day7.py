from collections import deque

grid = [list(line.rstrip("\n")) for line in open("inputday7.txt")]
R = len(grid)
C = len(grid[0])


start_col = grid[0].index("S")


q = deque()
q.append((0, start_col))
visited = set()
splits = 0

while q:
    r, c = q.popleft()
    nr = r + 1  

    
    if nr >= R or c < 0 or c >= C:
        continue

    if (nr, c) in visited:
        continue
    visited.add((nr, c))

    cell = grid[nr][c]

    if cell == "^":
        splits += 1
        
        q.append((nr, c - 1))
        q.append((nr, c + 1))
    else:
        
        q.append((nr, c))

print(splits)
