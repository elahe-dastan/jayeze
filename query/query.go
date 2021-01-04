package query

import (
	"bufio"
	"jayeze/tokenize"
	"log"
	"os"
	"strings"
)

// returns related doc ids
func Query(query string) []string{
	ans := make([]string, 0)
	tokens := strings.Split(query, " ")

	for _, token := range tokens {
		docIds := findDoc(token)
		if docIds == nil {
			continue
		}

		ans = append(ans, docIds...)
	}

	return ans
}

func Query2(query string) []string{
	ans := make([]string, 0)
	tokens := strings.Split(query, " ")

	for _, token := range tokens {
		docIds := findDoc(token)
		if docIds == nil {
			continue
		}

		ans = append(ans, docIds...)
	}

	return ans
}

// memsize = 6
func findDoc(token string) []string {
	filePath := "./blocks2/1.txt"
	f, err := os.Open(filePath)
	//defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		termPostingList := tokenize.Unmarshal(scanner.Text())
		if termPostingList.Term != token {
			continue
		}

		return termPostingList.PostingList
	}

	return nil
}

