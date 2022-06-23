package main

import (
	"time"

	textclient "word-histrogram/pkg/client"
)

func main() {
	client, err := textclient.New("localhost:1337")
	if err != nil {
		panic(err)
	}

	for i := 0; i < 1; i++ {
		go client.Send()
	}

	time.Sleep(3 * time.Second)
}
