package tfidf

import (
	"crypto/md5"
	"encoding/hex"
	"tfidf-test/lexer"
	"tfidf-test/util"
)

type TFIDF struct {
	docIndex  map[string]int         // train document index in TermFreqs
	termFreqs []map[string]int       // term frequency for each train document
	termDocs  map[string]int         // documents number for each term in train data
	n         int                    // number of documents in train data
	stopWords map[string]interface{} // words to be filtered
	tokenizer lexer.Tokenizer        // tokenizer, space is used as default
}

func New() *TFIDF {
	return &TFIDF{
		docIndex:  make(map[string]int),
		termFreqs: make([]map[string]int, 0),
		termDocs:  make(map[string]int),
		n:         0,
		tokenizer: &lexer.EnTokenizer{},
	}
}

// NewTokenizer new with specified tokenizer
func NewTokenizer(tokenizer lexer.Tokenizer) *TFIDF {
	return &TFIDF{
		docIndex:  make(map[string]int),
		termFreqs: make([]map[string]int, 0),
		termDocs:  make(map[string]int),
		n:         0,
		tokenizer: tokenizer,
	}
}

func (f *TFIDF) PrintFreq() {
	for _, value := range f.termFreqs {
		for o, a := range value {
			println(o, ":", a)
		}
	}
	// index := f.docHashPos("778a2407471e2afc2bf9459df37eb13a")

}

func (f *TFIDF) AddStopWords(words ...string) {
	if f.stopWords == nil {
		f.stopWords = make(map[string]interface{})
	}

	for _, word := range words {
		f.stopWords[word] = nil
	}
}

// AddStopWordsFile add stop words file to be filtered, with one word a line
func (f *TFIDF) AddStopWordsFile(file string) (err error) {
	lines, err := util.ReadLines(file)
	if err != nil {
		return
	}

	f.AddStopWords(lines...)
	return
}

// AddDocs add train documents
func (f *TFIDF) AddDocs(docs ...string) {
	for _, doc := range docs {
		h := hash(doc)
		if f.docHashPos(h) >= 0 {
			return
		}

		termFreq := f.termFreq(doc)
		if len(termFreq) == 0 {
			return
		}

		f.docIndex[h] = f.n
		f.n++

		f.termFreqs = append(f.termFreqs, termFreq)

		for term := range termFreq {
			f.termDocs[term]++
		}
	}
}

func (f *TFIDF) termFreq(doc string) (m map[string]int) {
	m = make(map[string]int)

	tokens := f.tokenizer.Tokenize(doc)
	if len(tokens) == 0 {
		return
	}

	for _, term := range tokens {
		if _, ok := f.stopWords[term]; ok {
			continue
		}

		m[term]++
	}

	return
}

func (f *TFIDF) docHashPos(hash string) int {
	if pos, ok := f.docIndex[hash]; ok {
		return pos
	}

	return -1
}

func hash(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}
