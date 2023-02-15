package main

import (
	"fmt"
	"log"
	"tfidf-test/util"

	"github.com/goccy/go-json"
)

type json_type struct {
	Title       string
	Description string
	Data        string
	Id          string
}

func main() {
	// f := tfidf.New()

	// f.AddStopWordsFile("stopwords-en.txt")
	// f.AddDocs("tokenizer support, contains english and", "jieba Chinese Tokenizer", "calculate tfidf value of giving document. tfidf", "used to manage go packages. tfidf")

	fileContent, err := util.ReadFileToString("data.json")
	if err != nil {
		log.Fatal("Error reading the file", err)
	}

	// myJsonString := `{{"some":1}}`

	var parsedData []json_type

	json.Unmarshal([]byte(fileContent), &parsedData)
	fmt.Println(parsedData[2].Data)
}
