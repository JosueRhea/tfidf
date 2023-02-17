package main

import (
	"fmt"
	"log"
	"tfidf/tfidf"
	"tfidf/util"

	"github.com/goccy/go-json"
)

func main() {
	fileContent, err := util.ReadFileToString("test-data.json")
	if err != nil {
		log.Fatal("Error reading the file", err)
	}

	var parsedData []map[string]interface{}

	json.Unmarshal([]byte(fileContent), &parsedData)
	t := tfidf.New()
	t.AddDocs(parsedData)
	{
		stopWordsFilePath := "stopwords-en.txt"
		error := t.AddStopWordsFile(stopWordsFilePath)
		if error != nil {
			log.Printf("Cannot add %s because %s", stopWordsFilePath, error.Error())
		}
	}
	{
		stopWordsFilePath := "stopwords-es.txt"
		error := t.AddStopWordsFile(stopWordsFilePath)
		if error != nil {
			log.Printf("Cannot add %s because %s", stopWordsFilePath, error.Error())
		}
	}
	t.AddStopWordsFile("stopwords-es.txt")
	// t.PrintDocsWithTermFreqs()
	search := t.CalculateTFIDF("marvel")
	fmt.Println(len(search))
	for _, value := range search {
		fmt.Printf("id: %s -> %f\n", value.ID, value.Rank)
		fmt.Printf("value.Data: %v\n", value.Data["description"])
	}
}
