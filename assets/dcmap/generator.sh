#!/usr/bin/python

import sys

if len(sys.argv)<2:
	print("Usage: ./generator <mapsize>")
	sys.exit(1)

size = int(sys.argv[1])

print("{")
print("\"width\": "+str(size)+",")
print("\"height\": "+str(size)+",")
print("\"tiles\": [")
for x in range(3,size-4):
	for y in range(3,size-4):
		wall0 = ""
		if x==3:
			wall0= "wall"
			if y%4 == 0:
				wall0= "wallwindow"
		wall1 = ""
		if y == 3:
			wall1 = "wall"
			if x%4 == 0:
				wall1= "wallwindow"
		end=","
		if y==size-5 and x==size-5:
			end=""
		print("\t{\"x\":"+str(x)+", \"y\":"+str(y)+", \"wall0\":\""+wall0+"\",\"wall1\":\""+wall1+"\",\"floor\":\"inside\",\"rotation\":0, \"decoration\":\"\"}"+end)
	print("\t")
print("]")
print("}")


