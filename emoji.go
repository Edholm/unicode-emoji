package emoji

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

// Emoji represents a slice of unicode code points in the "emoji" range
type Emoji struct {
	Runes []rune // one or more runes representing a single emoji
}

// String converts the Emoji to a string representation
func (e Emoji) String() string {
	return string(e.Runes)
}

var (
	url        = "https://www.unicode.org/Public/emoji/13.1/emoji-sequences.txt"
	emojiCache []Emoji
)

// AllEmojis returns all available emojis in the unicode specification. It will download and parse them from
// official sources, but will return a cached slice on subsequent calls
func AllEmojis() (emojis []Emoji, err error) {
	if len(emojiCache) > 0 {
		return emojiCache, nil
	}

	emojis, err = downloadAndParseEmojis(url)
	emojiCache = emojis
	return
}

// RandomEmoji returns a random emoji from all of them. A side-effect is that it will cache all emojis in-mem
func RandomEmoji() (Emoji, error) {
	emojis, err := AllEmojis()
	if err != nil {
		return Emoji{}, err
	}

	rndIndex := rand.Intn(len(emojis))
	return emojis[rndIndex], nil
}

// downloadAndParseEmojis downloads the unicode sequence list and parses emojis from it.
func downloadAndParseEmojis(url string) ([]Emoji, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("resp.Body.Close(): %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got %s, expected 200 OK from %q", resp.Status, url)
	}

	// 2053 is a magic number and and is the number of emojis in the 13.1 emoji-sequences.txt
	allEmojis := make([]Emoji, 0, 2053)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}
		emojis, err := parseEmoji(line)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, emoji := range emojis {
			allEmojis = append(allEmojis, emoji)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	return allEmojis, nil
}

// parseEmoji takes a line from the unicode sequence list and returns a slice of Emojis
func parseEmoji(line string) ([]Emoji, error) {
	endIndex := strings.Index(line, ";")
	// Represents one or more unicode code points, e.g. 00A9 FE0F or a range: 1F334..1F335
	codePointRange := strings.TrimSpace(line[:endIndex])

	if !strings.ContainsRune(codePointRange, '.') {
		runes, err := extractRunes(codePointRange)
		return []Emoji{
				{
					Runes: runes,
				}},
			err
	}

	runes, err := expandCodePointRange(codePointRange)
	if err != nil {
		return nil, err
	}
	// In this case each rune is a different emoji
	emojis := make([]Emoji, len(runes))
	for i, r := range runes {
		emojis[i] = Emoji{
			Runes: []rune{r},
		}
	}
	return emojis, nil
}

// extractRunes takes a code point string (e.g 00A9 FE0F, but not 1F334..1F335) and returns them parsed as runes
func extractRunes(codePoints string) ([]rune, error) {
	if strings.ContainsRune(codePoints, '.') {
		return nil, fmt.Errorf("code point range not supported yet")
	}

	// E.g. 00A9 FE0F
	split := strings.Split(codePoints, " ")
	runes := make([]rune, 0, len(split))
	for _, s := range split {
		codePoint, err := strconv.ParseInt(s, 16, 32)
		if err != nil {
			return nil, fmt.Errorf("strconv.ParseInt: %v", err)
		}
		runes = append(runes, rune(codePoint))
	}

	return runes, nil
}

// expandCodePointRange will take a range like 1F380..1F393 and return all code points as runes between the 'start' and 'end'
func expandCodePointRange(cpRange string) ([]rune, error) {
	split := strings.Split(cpRange, "..")
	if len(split) != 2 {
		return nil, fmt.Errorf("%q does not look like a code point range", cpRange)
	}

	first, err := strconv.ParseInt(split[0], 16, 32)
	if err != nil {
		return nil, fmt.Errorf("%q is not a valid code point: %w", split[0], err)
	}
	end, err := strconv.ParseInt(split[1], 16, 32)
	if err != nil {
		return nil, fmt.Errorf("%q is not a valid code point: %w", split[1], err)
	}

	runes := make([]rune, 0, 10)
	current := first
	for current != end {
		runes = append(runes, rune(current))
		current = current + 1
	}

	return runes, nil
}
