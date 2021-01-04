package tokenize

import (
	"bufio"
	"strings"
)

type TermPostingList struct {
	Term        string
	PostingList []string // always sorted
}

func (t TermPostingList) Marshal() string {
	return t.Term + " " + strings.Join(t.PostingList, ",")
}

func Marshal(list []TermPostingList) string {
	output := ""
	for i := range list {
		output += list[i].Marshal() + "\n"
	}

	return output
}

func Unmarshal(data string) TermPostingList {
	termPostingList := strings.Split(data, " ")
	//postingList := make([]int, 0)
	postings := strings.Split(termPostingList[1], ",")
	//for i := 0; i < len(postings); i++ {
	//	posting, err := strconv.Atoi(postings[i])
	//	if err != nil {
	//		log.Fatal()
	//	}
	//
	//	postingList = append(postingList, posting)
	//}

	return TermPostingList{
		Term:        termPostingList[0],
		PostingList: postings,
	}
}

type Finger struct {
	FileSeek *bufio.Scanner
	TermPostingList     TermPostingList
}

type Fingers []Finger

func (f Fingers) Len() int           { return len(f) }
func (f Fingers) Less(i, j int) bool { return f[i].TermPostingList.Term < f[j].TermPostingList.Term }
func (f Fingers) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
