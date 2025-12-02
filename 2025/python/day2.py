# total = 0

# with open("inputday2.txt") as f:
#     line = f.read().strip()

# ranges = line.split(",")

# def is_invalid_id(n):
#     s = str(n)
#     if len(s) % 2 != 0:
#         return False
#     if s[0] == '0':
#         return False

#     mid = len(s) // 2
#     return s[:mid] == s[mid:]

# for r in ranges:
#     start, end = map(int, r.split("-"))
#     for num in range(start, end + 1):
#         if is_invalid_id(num):
#             total += num

# print(total)


total = 0

with open("inputday2.txt") as f:
    line = f.read().strip()

ranges = line.split(",")

def is_invalid_id_part2(n):
    s = str(n)
    L = len(s)

    if s[0] == '0':
        return False

    for k in range(1, L // 2 + 1):
        if L % k != 0:
            continue
        r = L // k
        if r < 2:
            continue
        if s == s[:k] * r:
            return True

    return False

for r in ranges:
    start, end = map(int, r.split("-"))
    for num in range(start, end + 1):
        if is_invalid_id_part2(num):
            total += num

print(total)
