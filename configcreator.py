#!/usr/bin/python

import os


a=[]

file = open("cluster.conf")
while 1:
	line = file.readline()
	if not line:
		break
	line  = line [0:len(line)-1]
	a.append(line)


start = 0

for ff in a:
	name = ff.split(":")
	x = "cluster"+name[2]+".conf"
	fo = open(x, "w+")
	
	b = a[start:start+1]
	start = start+1
	c = []
	for em in a:
		if b[0] != em:
			c.append(em)
	
	for s in b:
		fo.write( s+"\n");
	for s in c:
		fo.write( s+"\n");

	fo.close()


