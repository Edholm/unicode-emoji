package emoji

import (
	"reflect"
	"testing"
)

func TestEmoji_String(t *testing.T) {
	tests := []struct {
		name  string
		runes []rune
		want  string
	}{
		{
			name:  "Single rune",
			runes: []rune{'\U0001f413'},
			want:  "\U0001F413",
		},
		{
			name:  "Multiple runes",
			runes: []rune{'\U0001F3F4', '\U000E0067', '\U000E0062', '\U000E0065', '\U000E006E', '\U000E0067', '\U000E007F'},
			want:  "\U0001F3F4\U000E0067\U000E0062\U000E0065\U000E006E\U000E0067\U000E007F",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Emoji{
				Runes: tt.runes,
			}
			if got := e.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseAllEmojis(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		wantLength int
		wantErr    bool
	}{
		{
			name:       "3521 emojis",
			wantLength: 3521,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseAllEmojis()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAllEmojis() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantLength {
				t.Errorf("parseAllEmojis() got length = %d, wanted %d", len(got), tt.wantLength)
			}
		})
	}
}

func Test_extractRunes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		codePoints string
		want       []rune
		wantErr    bool
	}{
		{
			name:       "Range not supported",
			codePoints: "0000..0001",
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "Single 2 byte rune",
			codePoints: "23f0",
			want:       []rune{'\u23f0'},
			wantErr:    false,
		},
		{
			name:       "2.5 byte rune",
			codePoints: "1F30F",
			want:       []rune{'\U0001f30f'},
			wantErr:    false,
		},
		{
			name:       "Double 2 byte runes",
			codePoints: "00A9 FE0F",
			want:       []rune{'\u00a9', '\ufe0f'},
			wantErr:    false,
		},
		{
			name:       "One 2.5 byte and one 2 byte runes",
			codePoints: "1F170 FE0F",
			want:       []rune{'\U0001f170', '\ufe0f'},
			wantErr:    false,
		},
		{
			name:       "Three 2 byte runes",
			codePoints: "002A FE0F 20E3",
			want:       []rune{'\u002a', '\ufe0f', '\u20e3'},
			wantErr:    false,
		},
		{
			name:       "Seven 2.5 byte runes",
			codePoints: "1F3F4 E0067 E0062 E0065 E006E E0067 E007F",
			want:       []rune{'\U0001f3f4', '\U000e0067', '\U000e0062', '\U000e0065', '\U000e006e', '\U000e0067', '\U000e007f'},
			wantErr:    false,
		},
		{
			name:       "Invalid hex",
			codePoints: "000G",
			want:       nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := extractRunes(tt.codePoints)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractRunes() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractRunes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseEmoji(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		line           string
		want           Emoji
		wantErr        bool
		parseCorrectly bool
	}{
		{
			name: "Basic 2 byte rune",
			line: "23F0                                                   ; fully-qualified     # ⏰ E0.6 alarm clock",
			want: Emoji{
				Runes: []rune{'\u23f0'},
				Name:  "alarm clock",
			},
			wantErr:        false,
			parseCorrectly: true,
		},
		{
			name:
			"Basic 2.5 byte rune",
			line: "1F004                                                  ; fully-qualified     # 🀄 E0.6 mahjong red dragon",
			want: Emoji{
				Runes: []rune{'\U0001f004'},
				Name:  "mahjong red dragon",
			},
			wantErr:        false,
			parseCorrectly: true,
		},
		{
			name:
			"Unsupported range",
			line:           "23E9..23EC    ; fully-qualified     # 🎴 E0.6 flower playing cards",
			want:           Emoji{},
			wantErr:        true,
			parseCorrectly: false,
		},
		{
			name:
			"Unqualified code point",
			line:           "2666                                                   ; unqualified         # ♦ E0.6 diamond suit",
			want:           Emoji{},
			wantErr:        false,
			parseCorrectly: false,
		},
		{
			name:
			"Minimally qualified code point",
			line:           "1F9D4 200D 2642                                        ; minimally-qualified # 🧔‍♂ E13.1 man: beard",
			want:           Emoji{},
			wantErr:        false,
			parseCorrectly: false,
		},
		{
			name:
			"Invalid rune",
			line:           "1F0CG                                                  ; fully-qualified     # 🃏 E0.6 joker",
			want:           Emoji{},
			wantErr:        true,
			parseCorrectly: false,
		},
		{
			name:
			"Component code point",
			line: "1F3FC                                                  ; component           # 🏼 E1.0 medium-light skin tone",
			want: Emoji{
				Runes: []rune{'\U0001f3fc'},
				Name:  "medium-light skin tone",
			},
			wantErr:        false,
			parseCorrectly: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			parsed, got, err := parseEmoji(tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseEmoji() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.parseCorrectly != parsed {
				t.Errorf("parseEmoji() expected parse: %v, got %v", tt.parseCorrectly, parsed)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseEmoji() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractNameFromDescription(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		description string
		want        string
	}{
		{
			name:        "E13.0",
			description: " fully-qualified     # \U0001FAC2 E13.0 people hugging",
			want:        "people hugging",
		},
		{
			name:        "E0.6",
			description: " fully-qualified     # 🔯 E0.6 dotted six-pointed star",
			want:        "dotted six-pointed star",
		},
		{
			name:        "E3.0",
			description: " fully-qualified     # 🤚🏾 E3.0 raised back of hand: medium-dark skin tone",
			want:        "raised back of hand: medium-dark skin tone",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := extractNameFromDescription(tt.description); got != tt.want {
				t.Errorf("extractNameFromDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmojis_Search(t *testing.T) {
	t.Parallel()
	emojis := NewEmojis()
	tests := []struct {
		name        string
		query       string
		wantMatched []Emoji
		wantErr     bool
	}{
		{
			name:  "poo",
			query: "poo",
			wantMatched: []Emoji{
				{
					Runes: []rune{'\U0001F4A9'},
					Name:  "pile of poo",
				},
				{
					Runes: []rune{'\U0001F429'},
					Name:  "poodle",
				},
				{
					Runes: []rune{'\U0001F963'},
					Name:  "bowl with spoon",
				},
				{
					Runes: []rune{'\U0001F944'},
					Name:  "spoon",
				},
				{
					Runes: []rune{'\U0001F3B1'},
					Name:  "pool 8 ball",
				},
			},
			wantErr: false,
		},
		{
			name:  "Uppercase Poo",
			query: "POO",
			wantMatched: []Emoji{
				{
					Runes: []rune{'\U0001F4A9'},
					Name:  "pile of poo",
				},
				{
					Runes: []rune{'\U0001F429'},
					Name:  "poodle",
				},
				{
					Runes: []rune{'\U0001F963'},
					Name:  "bowl with spoon",
				},
				{
					Runes: []rune{'\U0001F944'},
					Name:  "spoon",
				},
				{
					Runes: []rune{'\U0001F3B1'},
					Name:  "pool 8 ball",
				},
			},
			wantErr: false,
		},
		{
			name:        "Whitespace only",
			query:       "\t",
			wantMatched: nil,
			wantErr:     false,
		},
		{
			name:        "Empty query",
			query:       "",
			wantMatched: nil,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatched, err := emojis.Search(tt.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMatched, tt.wantMatched) {
				t.Errorf("Search() gotMatched = %v, want %v", gotMatched, tt.wantMatched)
			}
		})
	}
}
