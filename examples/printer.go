package main

import (
	emoji "edholm.dev/unicode-emoji"
	"fmt"
	"log"
)

func main() {
	emojis := emoji.NewEmojis()
	all, err := emojis.All()
	if err != nil {
		log.Fatalln(err)
	}

	for _, e := range all {
		fmt.Printf("%s | %s\n", e.String(), e.Name)
	}

	fmt.Printf("\n\n - Got %d emojis", len(all))
}
