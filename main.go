package main

import (
	"fmt"
	"jayeze/config"
	vectorspace "jayeze/vector-space"
	"log"
	"net/http"
	"strings"

	"github.com/elahe-dastan/trunk/normalize"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/labstack/echo/v4"
)

var v *vectorspace.Vectorizer
var championVector *vectorspace.Vectorizer
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
	e.GET("/", query)
	e.GET("/champion/", championQuery)
	e.GET("/cluster/", clusterQuery)
	e.Logger.Fatal(e.Start(":1373"))
}

func vectorize(k *koanf.Koanf, f *file.File) {
	if err := k.Load(f, yaml.Parser()); err != nil {
		log.Fatal(err)
	}

	var cfg config.Config
	// Quick unmarshal.
	err := k.Unmarshal("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	// main vectorizer
	v = vectorspace.NewVectorizer(cfg.IndexPath, cfg.DocsSize)
	v.Vectorize()

	// champion vectorizer
	championVector = vectorspace.NewVectorizer(cfg.ChampionPath, cfg.DocsSize)
	championVector.Vectorize()

	// cluster vectorizer
	clusterVectors = make([]*vectorspace.Vectorizer, 5)
	clusters := cfg.Clusters
	for i, cluster := range clusters{
		vectorizer := vectorspace.NewVectorizer(cluster.Path, cluster.Size)
		vectorizer.Vectorize()
		clusterVectors[i] = vectorizer
	}
}

func query(c echo.Context) error {
	return c.JSON(http.StatusOK, v.Query(c.QueryParam("query"), 4))
}

func championQuery(c echo.Context) error {
	return c.JSON(http.StatusOK, championVector.Query(c.QueryParam("query"), 4))
}

func clusterQuery(c echo.Context) error {
	var vr *vectorspace.Vectorizer
	maxCosSimilarity := float64(0)
	queryTerms := strings.Split(c.QueryParam("query"), " ")
	normalizedQuery := make([]string, 0)
	for _, t := range queryTerms{
		normalizedQuery = append(normalizedQuery, normalize.Normalize(t)...)
	}
	for _, vector := range clusterVectors {
		cos := vector.CenterCosineSimilarity(normalizedQuery)
		if cos > maxCosSimilarity {
			maxCosSimilarity = cos
			vr = vector
		}
	}

	fmt.Println(vr.IndexPath)

	return c.JSON(http.StatusOK, vr.Query(c.QueryParam("query"), 4))
}
