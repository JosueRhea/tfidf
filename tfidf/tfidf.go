package tfidf

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"tfidf-test/lexer"
)

type TFIDF struct {
	docIndex  map[string]int
	termFreqs []map[string]int
	termDocs  map[string]int
	n         int
	tokenizer lexer.Tokenizer
	stopWords map[string]struct{}
	data      []map[string]interface{}
}

func New() *TFIDF {
	return &TFIDF{
		docIndex:  make(map[string]int),
		termDocs:  make(map[string]int),
		n:         0,
		tokenizer: &lexer.EnTokenizer{},
		stopWords: make(map[string]struct{}),
	}
}

func (t *TFIDF) PrintDocsWithTermFreqs() {
	for key, value := range t.docIndex {
		fmt.Printf("%s:\n", key)
		docfreqs := t.termFreqs[value]
		for term, freq := range docfreqs {
			fmt.Printf("\t%s: %d\n", term, freq)
		}
	}
}

func jsonObjToString(data map[string]interface{}) (result string) {
	for _, value := range data {
		if str, ok := value.(string); ok {
			result += str + " "
		}
	}

	return result
}

func (t *TFIDF) AddDocs(docs []map[string]interface{}) {
	if t.termFreqs == nil {
		t.termFreqs = make([]map[string]int, len(docs))
	}
	if t.data == nil {
		t.data = make([]map[string]interface{}, len(docs))
	}
	for i, doc := range docs {
		id, ok := doc["id"].(string)
		if !ok {
			println("Id not provided for this document")
			continue
		}
		text := jsonObjToString(doc)

		tokens := t.tokenizer.Tokenize(text)
		freqs := make(map[string]int)
		for _, token := range tokens {
			if _, ok := t.stopWords[token]; ok {
				continue
			}
			term := strings.ToLower(token)
			freqs[term]++
			if freqs[term] == 1 {
				t.termDocs[term]++
			}
		}
		t.termFreqs[i] = freqs
		t.docIndex[id] = i
		t.n++
		t.data[i] = doc
	}
}

func (t *TFIDF) docFreq(term string) float64 {
	df := t.termDocs[term]
	return math.Log(float64(t.n) / float64(df))
}

func (t *TFIDF) termFreq(doc string, term string) float64 {
	id := t.docIndex[doc]
	docFreqs := t.termFreqs[id]
	tf := float64(docFreqs[term]) / float64(len(docFreqs))
	return tf
}

func (t *TFIDF) CalculateTFIDF(query string) []struct {
	ID   string
	Rank float64
} {
	tfidf := make(map[string]float64)

	tokens := t.tokenizer.Tokenize(query)

	for _, token := range tokens {
		if _, ok := t.stopWords[token]; ok {
			continue
		}

		term := strings.ToLower(token)

		for id, value := range t.docIndex {
			docFreqs := t.termFreqs[value]

			if _, ok := docFreqs[term]; ok {
				tf := t.termFreq(id, term)
				idf := t.docFreq(term)
				total := tf * idf
				tfidf[id] = total
			}
		}
	}

	// Convert map to slice
	var tfidfSlice []struct {
		ID   string
		Rank float64
	}
	for id, rank := range tfidf {
		tfidfSlice = append(tfidfSlice, struct {
			ID   string
			Rank float64
		}{
			ID:   id,
			Rank: rank,
		})
	}

	// Sort slice by rank in descending order
	sort.Slice(tfidfSlice, func(i, j int) bool {
		return tfidfSlice[i].Rank > tfidfSlice[j].Rank
	})

	return tfidfSlice
}
