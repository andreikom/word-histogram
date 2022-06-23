package tasks

import (
	"log"
	"strings"

	"word-histrogram/pkg/storage"
)

// WordsItem parses and saves words to a storage
type WordsItem struct {
	storage.WordStore
	Buffer []byte
}

// Execute parses and saves words to a storage
func (task *WordsItem) Execute() {
	converted := string(task.Buffer)
	words := strings.Split(converted, ",")

	wordCount := make(map[string]int)
	for _, word := range words {
		_, matched := wordCount[word]
		if matched {
			wordCount[word]++

			continue
		}

		wordCount[word] = 1
	}

	task.WordStore.Add(wordCount)
	log.Println("added new result to storage: " + converted)
}
