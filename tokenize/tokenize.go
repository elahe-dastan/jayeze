package tokenize

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)


type Finger struct {
	FileSeek *bufio.Scanner
	TermPostingList     TermPostingList
}

type Fingers []Finger

func (f Fingers) Len() int           { return len(f) }
func (f Fingers) Less(i, j int) bool { return f[i].TermPostingList.Term < f[j].TermPostingList.Term }
func (f Fingers) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

type PostingList struct {
	DocId     int
	Frequency int
}

func (p PostingList) Marshal() string {
	return strconv.Itoa(p.DocId) + ":" + strconv.Itoa(p.Frequency)
}

type PostingLists []PostingList

func (postingLists PostingLists) Len() int { return len(postingLists) }
func (postingLists PostingLists) Less(i, j int) bool { return postingLists[i].DocId < postingLists[j].DocId }
func (postingLists PostingLists) Swap(i, j int)      { postingLists[i], postingLists[j] = postingLists[j], postingLists[i] }

func (postingLists PostingLists) Marshal() string {
	ans := ""
	for _, postingList := range postingLists {
		ans += postingList.Marshal() + ","
	}
	ans = strings.TrimSuffix(ans, ",")
	return ans
}

type TermPostingList struct {
	Term        string
	PostingList PostingLists // always sorted
}

func (t TermPostingList) Marshal() string {
	return t.Term + " " + t.PostingList.Marshal()
}

func Marshal(termPostingLists []TermPostingList) string {
	output := ""
	for _, termPostingList := range termPostingLists {
		output += termPostingList.Marshal() + "\n"
	}

	return output
}
//func (ftp FinalTermPostingList) Marshal() string {
//
//}

func Unmarshal(line string) TermPostingList {
	fmt.Println(line)
	termPostingList := strings.Split(line, " ")
	docIdsFrequencies := strings.Split(termPostingList[1], ",")
	postingLists := make([]PostingList, 0)
	for _, docIdFrequency := range docIdsFrequencies{
		df := strings.Split(docIdFrequency, ":")
		d, _ := strconv.Atoi(df[0])
		f, _ := strconv.Atoi(df[1])
		postingLists = append(postingLists, PostingList{
			DocId:     d,
			Frequency: f,
		})
	}

	return TermPostingList{
		Term:        termPostingList[0],
		PostingList: postingLists,
	}
}