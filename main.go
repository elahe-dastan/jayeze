package main

import vectorspace "jayeze/vector-space"

func main() {
	v := vectorspace.NewVectorizer("/home/raha/go/src/shakhes/blocks4/1.txt", 100)
	v.Vectorize()
	v.Query("منصوریان")
}
