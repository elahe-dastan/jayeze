package main

import (
	"container/heap"
	"encoding/json"
	"fmt"
	heap2 "jayeze/heap"
	vectorspace "jayeze/vector-space"
	"net/http"

	"github.com/elahe-dastan/trunk/normalize"
	"github.com/labstack/echo/v4"
)

var v *vectorspace.Vectorizer
var clusterVectors []*vectorspace.Vectorizer
func main() {
	h := &heap2.SimilarityHeap{}
	heap.Init(h)

	h.Push(heap2.Similarity{
		DocId: 0,
		Cos:   0.23,
	})
	h.Push(heap2.Similarity{
		DocId: 1,
		Cos:   0,
	})
	h.Push(heap2.Similarity{
		DocId: 2,
		Cos:   0.1,
	})
	h.Push(heap2.Similarity{
		DocId: 3,
		Cos:   0.7,
	})

	ans, _ := json.Marshal(h)
	fmt.Println(string(ans))

	fmt.Println(h.Pop())
	fmt.Println(h.Pop())
	fmt.Println(h.Pop())
	fmt.Println(h.Pop())

//	k := koanf.New(".")
//	f := file.Provider("config.yml")
//	if err := f.Watch(func(event interface{}, err error) {
//		if err != nil{
//			log.Fatal(err)
//		}
//		vectorize(k, f)
//	}); err != nil{
//		log.Println(err)
//	}
//
//	vectorize(k, f)
//
//	e := echo.New()
//	e.GET("/:query", query)
//	e.GET("/cluster/:query", clusterQuery)
//	e.Logger.Fatal(e.Start(":1373"))
//}
//
//func vectorize(k *koanf.Koanf, f *file.File) {
//	if err := k.Load(f, yaml.Parser()); err != nil{
//		log.Fatal(err)
//	}
//	m := k.All()
//	v = vectorspace.NewVectorizer(m["indexPath"].(string), int(m["docsNum"].(float64)))
//	v.Vectorize()
//
//
//	// kesafat
//	clusterVectors = make([]*vectorspace.Vectorizer, 5)
//	a := vectorspace.NewVectorizer(m["behdashtPath"].(string), int(m["behdashtNum"].(float64)))
//	a.Vectorize()
//	clusterVectors[0] = a
//
//	b := vectorspace.NewVectorizer(m["tarikhPath"].(string), int(m["tarikhNum"].(float64)))
//	b.Vectorize()
//	clusterVectors[1] = b
//
//	c := vectorspace.NewVectorizer(m["riaziatPath"].(string), int(m["riaziatNum"].(float64)))
//	c.Vectorize()
//	clusterVectors[2] = c
//
//	d := vectorspace.NewVectorizer(m["fanavariPath"].(string), int(m["fanavariNum"].(float64)))
//	d.Vectorize()
//	clusterVectors[3] = d
//
//	e := vectorspace.NewVectorizer(m["fizikPath"].(string), int(m["fizikNum"].(float64)))
//	e.Vectorize()
//	clusterVectors[4] = e
}

func query(c echo.Context) error {
	return c.JSON(http.StatusOK, v.Query(normalize.Normalize(c.Param("query"))[0]))
}

func clusterQuery(c echo.Context) error {
	var vr *vectorspace.Vectorizer
	maxCosSimilarity := float64(0)
	for _, vector := range clusterVectors{
		cos := vector.CenterCosineSimilarity(c.Param("query"))
		if cos > maxCosSimilarity {
			maxCosSimilarity = cos
			vr = vector
		}
	}

	return c.JSON(http.StatusOK, vr.Query(c.Param("query")))
}