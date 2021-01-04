package main

import vector_space "jayeze/vector-space"

func main() {
	v := vector_space.NewVectorizer("indexFile", 3)
	v.Vectorize()
	v.Query("نشست کمیسیون")
}
