# unicode-emoji

A simple Go package that allows you to download and list all available emojis in the unicode space.

## Dependencies

* Go 1.16

## Building & Running

A simple `go build` will do. As for running you can run a super simple example in the `printer` directory:

```
cd printer
go run printer.go
```

## Example usage

Below you can find simple usages of the emoji package.

### Random emoji every second printed to stdout

```go
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

```

### Print all emojis

```go
package main

import (
	emoji "edholm.dev/unicode-emoji"
)

func main() {
	emojis, err := emoji.AllEmojis()
	if err != nil {
		panic(err)
	}

	for _, e := range emojis {
		println(e.String())
	}
}
```
