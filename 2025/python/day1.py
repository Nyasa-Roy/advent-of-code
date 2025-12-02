with open("inputday1.txt") as f:
    inputs = [line.strip() for line in f.readlines()]
l = len(inputs)
s = 50
curr=s
out=0

for i in range(l):
    dir = inputs[i][0]
    num = int(inputs[i][1:])
    print(inputs[i])
    print(dir,num)
    if(dir=="L"):
        curr=(curr-num)%100
    if(dir=="R"):
        curr=(curr+num)%100
    print(curr)
    if(curr==0):
        out+=1
print(out)