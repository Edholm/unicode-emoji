package main

import (
	emoji "edholm.dev/unicode-emoji"
	"log"
	"time"
)

func main() {
	for {
		randomEmoji, err := emoji.RandomEmoji()
		if err != nil {
			log.Panicln(err)
		}
		println(randomEmoji.String())
		time.Sleep(1 * time.Second)
	}
}
