grid = [line.rstrip("\n") for line in open("inputday6.txt")]

rows = len(grid)
cols = len(grid[0])

problems = []
col = 0

while col < cols:
    
    empty = True
    for r in range(rows):
        if grid[r][col] != " ":
            empty = False
            break

    if empty:
        col += 1
        continue

    
    group_cols = []

    
    while col < cols:
        # Check if this column is empty
        is_empty = True
        for r in range(rows):
            if grid[r][col] != " ":
                is_empty = False
                break

        if is_empty:
            break

        group_cols.append(col)
        col += 1

    
    numbers = []
    operator = None

    
    for r in range(rows - 1):  # last row is operator
        num = ""
        for c in group_cols:
            ch = grid[r][c]
            if ch != " ":
                num += ch

        if num != "":
            numbers.append(int(num))

    
    for c in group_cols:
        if grid[rows - 1][c] in "+*":
            operator = grid[rows - 1][c]
            break

    problems.append((numbers, operator))

grand_total = 0

for nums, op in problems:
    if op == "+":
        total = sum(nums)
    else:  
        total = 1
        for n in nums:
            total *= n
    grand_total += total

print(grand_total)
