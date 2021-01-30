package main

import (
	"fmt"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	vectorspace "jayeze/vector-space"
	"log"
	"net/http"

	"github.com/elahe-dastan/trunk/normalize"
	"github.com/labstack/echo/v4"
)

var v *vectorspace.Vectorizer
var clusterVectors []*vectorspace.Vectorizer
func main() {
	k := koanf.New(".")
	f := file.Provider("config/config.yml")
	if err := f.Watch(func(event interface{}, err error) {
		if err != nil {
			log.Fatal(err)
		}
		vectorize(k, f)
	}); err != nil {
		log.Println(err)
	}

	vectorize(k, f)

	e := echo.New()
	e.GET("/:query", query)
	e.GET("/cluster/:query", clusterQuery)
	e.Logger.Fatal(e.Start(":1373"))
}

func vectorize(k *koanf.Koanf, f *file.File) {
	if err := k.Load(f, yaml.Parser()); err != nil {
		log.Fatal(err)
	}
	m := k.All()
	v = vectorspace.NewVectorizer(m["indexPath"].(string), int(m["docsNum"].(float64)))
	v.Vectorize()

	pa := m["clusters"]
	// kesafat
	//clusterVectors = make([]*vectorspace.Vectorizer, 5)
	//a := vectorspace.NewVectorizer(m["behdashtPath"].(string), int(m["behdashtNum"].(float64)))
	//a.Vectorize()
	//clusterVectors[0] = a
	//
	//b := vectorspace.NewVectorizer(m["tarikhPath"].(string), int(m["tarikhNum"].(float64)))
	//b.Vectorize()
	//clusterVectors[1] = b
	//
	//c := vectorspace.NewVectorizer(m["riaziatPath"].(string), int(m["riaziatNum"].(float64)))
	//c.Vectorize()
	//clusterVectors[2] = c
	//
	//d := vectorspace.NewVectorizer(m["fanavariPath"].(string), int(m["fanavariNum"].(float64)))
	//d.Vectorize()
	//clusterVectors[3] = d
	//
	//e := vectorspace.NewVectorizer(m["fizikPath"].(string), int(m["fizikNum"].(float64)))
	//e.Vectorize()
	//clusterVectors[4] = e
}

func query(c echo.Context) error {
	return c.JSON(http.StatusOK, v.Query(normalize.Normalize(c.Param("query"))[0]))
}

func clusterQuery(c echo.Context) error {
	var vr *vectorspace.Vectorizer
	maxCosSimilarity := float64(0)
	for _, vector := range clusterVectors {
		cos := vector.CenterCosineSimilarity(c.Param("query"))
		if cos > maxCosSimilarity {
			maxCosSimilarity = cos
			vr = vector
		}
	}

	return c.JSON(http.StatusOK, vr.Query(c.Param("query")))
}
