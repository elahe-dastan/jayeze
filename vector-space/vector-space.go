package vector_space

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"github.com/elahe-dastan/trunk/normalize"
	"io/ioutil"
	heap2 "jayeze/heap"
	"jayeze/set"
	"jayeze/tokenize"
	"log"
	"math"
	"strings"
)

type Vectorizer struct {
	termPostingLists []tokenize.TermPostingList
	docsNum          int
	termsNum         int
	idf              []float64
	tf               [][]int
	tfIdf            [][]float64
	termIndex        map[string]int
	heap             *heap2.SimilarityHeap
	center           []float64
	postingList      map[string][]int
}

func NewVectorizer(indexPath string, docsNum int) *Vectorizer {
	dat, err := ioutil.ReadFile(indexPath)
	if err != nil {
		log.Fatal(err)
	}

	tmp := strings.Split(string(dat), "\n")
	lines := tmp[:len(tmp)-1]
	termPostingLists := make([]tokenize.TermPostingList, len(lines))

	termIndex := make(map[string]int)
	postingList := make(map[string][]int)

	for i, l := range lines {
		termPostingList := tokenize.Unmarshal(l)
		termPostingLists[i] = termPostingList
		termIndex[termPostingList.Term] = i
		finalTermPostingList := tokenize.Unmarshal(l)
		docIds := make([]int, len(finalTermPostingList.PostingList))
		for j, p := range finalTermPostingList.PostingList {
			docIds[j] = p.DocId
		}

		postingList[finalTermPostingList.Term] = docIds
	}

	tf := make([][]int, docsNum)
	for i := 0; i < docsNum; i++ {
		tf[i] = make([]int, len(lines))
	}

	tfIdf := make([][]float64, docsNum)
	for i := 0; i < docsNum; i++ {
		tfIdf[i] = make([]float64, len(lines))
	}

	center := make([]float64, len(lines))

	h := &heap2.SimilarityHeap{}
	heap.Init(h)

	return &Vectorizer{
		termPostingLists: termPostingLists,
		docsNum:          docsNum,
		termsNum:         len(lines),
		tf:               tf,
		tfIdf:            tfIdf,
		termIndex:        termIndex,
		heap:             h,
		center:           center,
		postingList:      postingList,
	}
}

func (v *Vectorizer) Vectorize() {
	v.calculateIDF()
	v.calculateTF()
	v.calculateTFIDF()
	v.calculateCenter()
}

func (v *Vectorizer) calculateIDF() {
	v.idf = make([]float64, v.termsNum)

	for i, t := range v.termPostingLists {
		// the formula is the log of [docsNum / (number of docs containing the term)] but in the case that
		// all the documents contain the word the answer will be ZERO so I'll use `docsNum + 1` instead of docsNum
		if i == 2899{
			fmt.Println()
		}
		v.idf[i] = math.Log10(float64(v.docsNum + 1)/ float64(len(t.PostingList)))
	}
}

func (v *Vectorizer) calculateTF() {
	// i expresses term index
	for i, t := range v.termPostingLists {
		for j := 0; j < len(t.PostingList); j++ {
			docId := t.PostingList[j].DocId
			v.tf[docId-1][i] = t.PostingList[j].Frequency
		}
	}
}

func (v *Vectorizer) calculateTFIDF() {
	for i := 0; i < v.docsNum; i++ {
		for j := 0; j < v.termsNum; j++ {
			v.tfIdf[i][j] = (math.Log10(1 + float64(v.tf[i][j]))) * v.idf[j]
		}
	}
}

func (v *Vectorizer) calculateCenter() {
	for i := 0; i < v.docsNum; i++ {
		vector := v.tfIdf[i]
		for j := 0; j < v.termsNum; j++ {
			v.center[j] += vector[j]
		}
	}

	for i := 0; i < v.termsNum; i++ {
		v.center[i] /= float64(v.docsNum)
	}
}

func (v *Vectorizer) Query(query string, k int) string {
	queryTerms := strings.Split(query, " ")
	normalizedQuery := make([]string, 0)
	for _, t := range queryTerms{
		normalizedQuery = append(normalizedQuery, normalize.Normalize(t)...)
	}

	queryVector := v.queryVectorizer(normalizedQuery)
	heapSize := v.cosineSimilarity(queryVector, normalizedQuery)
	answer := ""
	m := heapSize
	if k < m {
		m = k
	}
	for i := 0; i < m; i++ {
		docSimilarity := heap.Pop(v.heap).(heap2.Similarity)
		ans, err := json.Marshal(docSimilarity)
		//fmt.Println(ans)
		if err != nil {
			log.Fatal(err)
		}
		answer += string(ans)
	}

	return answer
}

func (v *Vectorizer) CenterCosineSimilarity(query []string) float64 {
	queryVector := v.queryVectorizer(query)
	// query vector is not normalized and it's vector is just tf not tf-idf
	innerProduct := float64(0)
	norm := float64(0) // this is norm powered by two
	for i, tfIdf := range v.center {
		innerProduct += tfIdf * queryVector[i]
		norm += math.Pow(tfIdf, 2)
	}
	cos := innerProduct / math.Sqrt(norm)

	return cos
}

func (v *Vectorizer) queryVectorizer(query []string) []float64 {
	vector := make([]float64, v.termsNum)
	//fmt.Println(tokens)
	for _, t := range query {
		index, ok := v.termIndex[t]
		if !ok {
			continue
		}
		vector[index]++
	}
	return vector
}

func (v *Vectorizer) indexElimination(query []string) []int{
	docIds := set.MakeSet()
	for _, q := range query{
		postingList, ok := v.postingList[q]
		if ok {
			for _, p := range postingList{
				docIds.Add(p)
			}
		}
	}

	ans := make([]int, 0)
	for k, _ := range docIds.Container{
		ans = append(ans, k)
	}

	return ans
}

// read only the p_docs in the posting list -- first read only the p_docs in the champion list
func (v *Vectorizer) cosineSimilarity(queryVector []float64, query []string) int {
	heapSize := 0
	docIds := v.indexElimination(query)
	// query vector is not normalized and it's vector is just tf not tf-idf
	for _, docId := range docIds {
		doc := v.tfIdf[docId - 1]
		//fmt.Println(docId)
		innerProduct := float64(0)
		norm := float64(0) // this is norm powered by two
		for i, tfIdf := range doc {
			innerProduct += tfIdf * queryVector[i]
			norm += math.Pow(tfIdf, 2)
		}
		cos := innerProduct / math.Sqrt(norm)
		heap.Push(v.heap, heap2.Similarity{
			DocId: docId,
			Cos:   cos,
		})
		heapSize++
	}

	return heapSize
}
