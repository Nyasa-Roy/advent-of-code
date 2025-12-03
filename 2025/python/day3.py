with open("inputday3.txt") as f:
    inputs = [line.strip() for line in f.readlines()]
sum = 0
for i in inputs:
    max1="0"
    max2="0"
    idx1=-1
    idx2=-1
    for idx,j in enumerate(i):
        if j>max1 and idx<len(i)-1:
            max1=j
            idx1=idx
            max2="0"
            idx2=-1
        elif j>max2 and idx>idx1:
            max2=j
            idx2=idx
    print(max1+max2)
    sum+= int(max1+max2)
print(sum)