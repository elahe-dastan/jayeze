package main

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	vectorspace "jayeze/vector-space"
	"log"
)

func main() {
	k := koanf.New(".")
	f := file.Provider("config.yml")
	if err := k.Load(f, yaml.Parser()); err != nil{
		log.Fatal(err)
	}

	m := k.All()
	v := vectorspace.NewVectorizer(m["indexPath"].(string), m["docsNum"].(int))
	v.Vectorize()
	v.Query("منصوریان")
}
