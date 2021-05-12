package main

import (
	emoji "edholm.dev/unicode-emoji"
	"fmt"
)

func main() {
	emojis := emoji.NewEmojis()
	var query string
	for {
		fmt.Print("Search for emojis (q to quit): ")
		_, err := fmt.Scanln(&query)
		if err != nil {
			fmt.Printf("Failed to read your query: %v. Try again!\n", err)
			continue
		}

		if query == "q" {
			return
		}

		results, err := emojis.Search(query)
		if err != nil {
			fmt.Printf("Query failed: %v\n", err)
			continue
		}

		printResults(results)
	}
}

func printResults(emojis []emoji.Emoji) {
	fmt.Printf("Found %d matching emojis\n", len(emojis))
	for _, e := range emojis {
		fmt.Printf("%s (%U) - %s\n", e.String(), e.Runes, e.Name)
	}
}
