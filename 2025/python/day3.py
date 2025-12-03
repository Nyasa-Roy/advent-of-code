# with open("inputday3.txt") as f:
#     inputs = [line.strip() for line in f.readlines()]
# sum = 0
# for i in inputs:
#     max1="0"
#     max2="0"
#     idx1=-1
#     idx2=-1
#     for idx,j in enumerate(i):
#         if j>max1 and idx<len(i)-1:
#             max1=j
#             idx1=idx
#             max2="0"
#             idx2=-1
#         elif j>max2 and idx>idx1:
#             max2=j
#             idx2=idx
#     print(max1+max2)
#     sum+= int(max1+max2)
# print(sum)


with open("inputday3.txt") as f:
    lines = [line.strip() for line in f]

total = 0

for s in lines:
    take = 12
    result = []
    start = 0
    n = len(s)

    while len(result) < take:
        left_needed = take - len(result)
        end = n - left_needed

        best_digit = '0'
        best_pos = start

        for i in range(start, end + 1):
            if s[i] > best_digit:
                best_digit = s[i]
                best_pos = i

        result.append(best_digit)
        start = best_pos + 1

    total += int("".join(result))

print(total)
