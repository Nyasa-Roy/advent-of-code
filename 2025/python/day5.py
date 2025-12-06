# ranges = []
# ids = []

# with open("inputday5.txt") as f:
#     lines = [line.strip() for line in f]

# i = 0

# while lines[i] != "":
#     a, b = lines[i].split("-")
#     ranges.append((int(a), int(b)))
#     i += 1

# i += 1

# while i < len(lines) and lines[i] != "":
#     ids.append(int(lines[i]))
#     i += 1

# fresh = 0

# for x in ids:
#     ok = False
#     for a, b in ranges:
#         if a <= x <= b:
#             ok = True
#             break
#     if ok:
#         fresh += 1

# print(fresh)


ranges = []

with open("inputday5.txt") as f:
    lines = [line.strip() for line in f]

i = 0

while lines[i] != "":
    a, b = lines[i].split("-")
    ranges.append((int(a), int(b)))
    i += 1

ranges.sort()

merged = []
start, end = ranges[0]

for a, b in ranges[1:]:
    if a <= end + 1:       
        end = max(end, b)   
    else:
        merged.append((start, end))
        start, end = a, b

merged.append((start, end))

total = 0
for a, b in merged:
    total += (b - a + 1)

print(total)
