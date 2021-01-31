package vector_space

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"io/ioutil"
	heap2 "jayeze/heap"
	"jayeze/tokenize"
	"log"
	"math"
	"strconv"
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
		finalTermPostingList := tokenize.UnmarshalFinal(l)
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

	fmt.Println(postingList)

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
		count := 1
		for j := 1; j < len(t.PostingList); j++ {
			if t.PostingList[j] != t.PostingList[j-1] {
				count++
			}
		}

		v.idf[i] = math.Log10(float64(v.docsNum / count))
	}
}

func (v *Vectorizer) calculateTF() {
	// i expresses term index
	for i, t := range v.termPostingLists {
		for j := 0; j < len(t.PostingList); j++ {
			docId, err := strconv.Atoi(t.PostingList[j])
			if err != nil {
				log.Fatal(err)
			}

			v.tf[docId-1][i]++
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

	//fmt.Println(v.center)
}

func (v *Vectorizer) Query(query []string) string {
	queryVector := v.queryVectorizer(query)
	v.cosineSimilarity(queryVector)
	answer := ""
	for i := 0; i < 100; i++ {
		docSimilarity := heap.Pop(v.heap).(heap2.Similarity)
		if docSimilarity.Cos == 0 {
			break
		}
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

// read only the p_docs in the posting list -- first read only the p_docs in the champion list
func (v *Vectorizer) cosineSimilarity(queryVector []float64) {
	// query vector is not normalized and it's vector is just tf not tf-idf
	for docId, doc := range v.tfIdf {
		fmt.Println(docId)
		innerProduct := float64(0)
		norm := float64(0) // this is norm powered by two
		for i, tfIdf := range doc {
			innerProduct += tfIdf * queryVector[i]
			norm += math.Pow(tfIdf, 2)
		}
		cos := innerProduct / math.Sqrt(norm)
		heap.Push(v.heap, heap2.Similarity{
			DocId: docId + 1,
			Cos:   cos,
		})
	}
}
