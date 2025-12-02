total = 0

with open("inputday2.txt") as f:
    line = f.read().strip()

ranges = line.split(",")

def is_invalid_id(n):
    s = str(n)
    if len(s) % 2 != 0:
        return False
    if s[0] == '0':
        return False

    mid = len(s) // 2
    return s[:mid] == s[mid:]

for r in ranges:
    start, end = map(int, r.split("-"))
    for num in range(start, end + 1):
        if is_invalid_id(num):
            total += num

print(total)
