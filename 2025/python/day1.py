with open("inputday1.txt") as f:
    inputs = [line.strip() for line in f.readlines()]

curr = 50
out = 0   
ans = 0   

for step in inputs:
    dir = step[0]
    num = int(step[1:])

    if dir == "L":
        first = curr if curr != 0 else 100
        if num >= first:
            hits = 1 + (num - first) // 100
        else:
            hits = 0
        curr = (curr - num) % 100

    else: 
        first = (100 - curr) if curr != 0 else 100
        if num >= first:
            hits = 1 + (num - first) // 100
        else:
            hits = 0
        curr = (curr + num) % 100

    if curr == 0:
        out += 1
        if hits > 0:
            ans += hits - 1
    else:
        ans += hits

print(out)          
print(ans + out)    