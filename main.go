package main

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/labstack/echo/v4"
	vectorspace "jayeze/vector-space"
	"log"
	"net/http"
)

var v *vectorspace.Vectorizer
func main() {
	k := koanf.New(".")
	f := file.Provider("config.yml")
	if err := f.Watch(func(event interface{}, err error) {
		if err != nil{
			log.Fatal(err)
		}
		vectorize(k, f)
	}); err != nil{
		log.Println(err)
	}

	vectorize(k, f)

	e := echo.New()
	e.GET("/:query", query)
	e.Logger.Fatal(e.Start(":1373"))
}

func vectorize(k *koanf.Koanf, f *file.File) {
	if err := k.Load(f, yaml.Parser()); err != nil{
		log.Fatal(err)
	}
	m := k.All()
	v = vectorspace.NewVectorizer(m["indexPath"].(string), int(m["docsNum"].(float64)))
	v.Vectorize()
}

func query(c echo.Context) error {
	return c.JSON(http.StatusOK, v.Query(c.Param("query")))
}