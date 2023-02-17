package tfidf

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"tfidf/lexer"
	"tfidf/util"
)

type TFIDF struct {
	docIndex  map[string]int
	termFreqs []map[string]int
	termDocs  map[string]int
	n         int
	tokenizer lexer.Tokenizer
	stopWords map[string]interface{}
	data      []map[string]interface{}
}

func New() *TFIDF {
	return &TFIDF{
		docIndex:  make(map[string]int),
		termDocs:  make(map[string]int),
		n:         0,
		tokenizer: &lexer.EnTokenizer{},
		stopWords: make(map[string]interface{}),
	}
}

func (t *TFIDF) AddStopWords(words ...string) {
	if t.stopWords == nil {
		t.stopWords = make(map[string]interface{})
	}

	for _, word := range words {
		t.stopWords[word] = nil
	}
}

func (t *TFIDF) AddStopWordsFile(file string) (err error) {
	lines, err := util.ReadLines(file)
	if err != nil {
		return
	}

	t.AddStopWords(lines...)
	return
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

func jsonObjToString(data map[string]interface{}) string {
	var result string

	for _, value := range data {
		switch v := value.(type) {
		case string:
			result += v + " "
		case []interface{}:
			for _, elem := range v {
				if str, ok := elem.(string); ok {
					result += str + " "
				}
			}
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
			freqs[token]++
			if freqs[token] == 1 {
				t.termDocs[token]++
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

type TFIDFResult struct {
	ID   string
	Rank float64
	Data map[string]interface{}
}

func (t *TFIDF) CalculateTFIDF(query string) []TFIDFResult {
	var result []TFIDFResult

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
				result = append(result, TFIDFResult{ID: id, Rank: total, Data: t.data[value]})
			}
		}
	}

	// Sort slice by rank in descending order
	sort.Slice(result, func(i, j int) bool {
		return result[i].Rank > result[j].Rank
	})

	return result
}
