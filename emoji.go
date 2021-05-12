package emoji

import (
	"bufio"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// Emoji represents a slice of unicode code points in the "emoji" range
type Emoji struct {
	Runes []rune // one or more runes representing a single emoji
	Name  string // The official name of the emoji, e.g. grinning face
}

// String converts the Emoji to a string representation
func (e Emoji) String() string {
	return string(e.Runes)
}

type Emojis struct {
	parseMu sync.Mutex // Sync parsing
	cache   []Emoji    // In-mem cache of already parsed emojis
	// TODO: expire cache after some time?
}

func NewEmojis() *Emojis {
	return &Emojis{}
}

var (
	// ErrInvalidCodePoint denotes that the unicode code point is invalid or un-parsable.
	ErrInvalidCodePoint = errors.New("invalid unicode code point")
	// ErrParsingFailed denotes failure to parse the embedded unicode data for some reason.
	ErrParsingFailed = errors.New("unable to parse unicode emoji data")
	//go:embed emoji-test.txt
	unicodeData string
	// descJunkRe is regexp to remove junk from the emoji description to extract only the name.
	descJunkRe = regexp.MustCompile(".+ E[0-9]+.[0-9]+ ")
)

// All returns all available emojis in the unicode specification. It will download and parse them from
// official sources, but will return a cached slice on subsequent calls
func (e *Emojis) All() (emojis []Emoji, err error) {
	if len(e.cache) > 0 {
		return e.cache, nil
	}

	e.parseMu.Lock()
	emojis, err = parseAllEmojis()
	e.cache = emojis
	e.parseMu.Unlock()

	return
}

// Random returns a random emoji from the set of all. A side-effect is that it may cache all emojis in-mem
func (e *Emojis) Random() (Emoji, error) {
	emojis, err := e.All()
	if err != nil {
		return Emoji{}, err
	}

	rndIndex := rand.Intn(len(emojis))

	return emojis[rndIndex], nil
}

func (e *Emojis) Search(query string) (matched []Emoji, err error) {
	if strings.TrimSpace(query) == "" {
		return
	}

	all, err := e.All()
	if err != nil {
		return nil, err
	}

	query = strings.ToLower(query)
	for _, emoji := range all {
		name := strings.ToLower(emoji.Name)
		if strings.Contains(name, query) {
			matched = append(matched, emoji)
		}
	}

	return
}

// parseAllEmojis parses the embedded unicode emoji list
func parseAllEmojis() ([]Emoji, error) {
	// 3521 is a magic number and and is the number of emojis in the 13.1 emoji-test.txt (fully-qualified + component)
	allEmojis := make([]Emoji, 0, 3521)
	scanner := bufio.NewScanner(strings.NewReader(unicodeData))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}
		parsed, emojis, err := parseEmoji(line)
		if err != nil {
			log.Println(err)

			continue
		} else if !parsed {
			continue
		}

		allEmojis = append(allEmojis, emojis)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("%w: error scanning unicode data: %v", ErrParsingFailed, err)
	}

	return allEmojis, nil
}

// parseEmoji takes a line from the unicode list and returns a slice of Emoji
func parseEmoji(line string) (bool, Emoji, error) {
	split := strings.Split(line, ";")
	if len(split) != 2 {
		return false, Emoji{}, fmt.Errorf("%w: malformed line: %q", ErrInvalidCodePoint, line)
	}

	desc := strings.TrimSpace(split[1])
	if !strings.HasPrefix(desc, "fully-qualified") && !strings.HasPrefix(desc, "component") {
		return false, Emoji{}, nil
	}

	// Represents one or more unicode code points, e.g. 00A9 FE0F but _NOT_ a range like 1F334..1F335
	codePoints := strings.TrimSpace(split[0])

	runes, err := extractRunes(codePoints)
	if err != nil {
		return false, Emoji{}, err
	}

	return true, Emoji{
		Runes: runes,
		Name:  extractNameFromDescription(desc),
	}, nil
}

// extractRunes takes a code point string (e.g 00A9 FE0F, but not 1F334..1F335) and returns them parsed as runes
func extractRunes(codePoints string) ([]rune, error) {
	if strings.ContainsRune(codePoints, '.') {
		return nil, fmt.Errorf("%w: code point range not supported yet", ErrInvalidCodePoint)
	}

	// E.g. 00A9 FE0F
	split := strings.Split(codePoints, " ")
	runes := make([]rune, 0, len(split))
	for _, s := range split {
		codePoint, err := strconv.ParseInt(s, 16, 32)
		if err != nil {
			return nil, fmt.Errorf("%w: strconv.ParseInt: %v", ErrInvalidCodePoint, err)
		}
		runes = append(runes, rune(codePoint))
	}

	return runes, nil
}

func extractNameFromDescription(desc string) string {
	split := descJunkRe.Split(desc, 2)
	return split[1]
}
