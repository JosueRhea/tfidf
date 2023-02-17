package main

import (
	"fmt"
	"log"
	"tfidf-test/tfidf"
	"tfidf-test/util"

	"github.com/goccy/go-json"
)

func main() {
	// f := tfidf.New()

	// f.AddStopWordsFile("stopwords-en.txt")
	// f.AddDocs("tokenizer support, contains english and", "jieba Chinese Tokenizer", "calculate tfidf value of giving document. tfidf", "used to manage go packages. tfidf")

	fileContent, err := util.ReadFileToString("data.json")
	if err != nil {
		log.Fatal("Error reading the file", err)
	}

	// myJsonString := `{{"some":1}}`

	var parsedData []map[string]interface{}

	json.Unmarshal([]byte(fileContent), &parsedData)
	t := tfidf.New()
	t.AddDocs(parsedData)
	// t.PrintDocsWithTermFreqs()
	search := t.CalculateTFIDF("fraccion")
	for _, value := range search {
		fmt.Printf("id: %s -> %f\n", value.ID, value.Rank)
	}
}
