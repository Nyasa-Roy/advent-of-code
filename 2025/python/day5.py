ranges = []
ids = []

with open("inputday5.txt") as f:
    lines = [line.strip() for line in f]

i = 0

while lines[i] != "":
    a, b = lines[i].split("-")
    ranges.append((int(a), int(b)))
    i += 1

i += 1

while i < len(lines) and lines[i] != "":
    ids.append(int(lines[i]))
    i += 1

fresh = 0

for x in ids:
    ok = False
    for a, b in ranges:
        if a <= x <= b:
            ok = True
            break
    if ok:
        fresh += 1

print(fresh)
