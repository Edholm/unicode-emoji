package emoji

import (
	"errors"
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

func Test_expandCodePointRange(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		cpRange string
		want    []rune
	}{
		{
			name:    "Range with two runes",
			cpRange: "231A..231B",
			want:    []rune{'\u231a', '\u231b'},
		},
		{
			name:    "Range with four runes",
			cpRange: "23E9..23EC",
			want:    []rune{'\u23e9', '\u23ea', '\u23eb', '\u23ec'},
		},
		{
			name:    "3 byte range",
			cpRange: "1FAC0..1FAC2",
			want:    []rune{'\U0001fac0', '\U0001fac1', '\U0001fac2'},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := expandCodePointRange(tt.cpRange)
			if err != nil {
				t.Errorf("expandCodePointRange() got unexpected error = %v", err)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("expandCodePointRange() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_expandCodePointRangeValidation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		cpRange    string
		wantErrMsg string
	}{
		{
			name:       "No range supplied",
			cpRange:    "23F0",
			wantErrMsg: "invalid unicode code point: \"23F0\" does not look like a code point range",
		},
		{
			name:       "Faulty start code point",
			cpRange:    "23G0..",
			wantErrMsg: "invalid unicode code point: \"23G0\" because strconv.ParseInt: parsing \"23G0\": invalid syntax",
		},
		{
			name:       "Faulty end code point",
			cpRange:    "23F0..23G0",
			wantErrMsg: "invalid unicode code point: \"23G0\" because strconv.ParseInt: parsing \"23G0\": invalid syntax",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := expandCodePointRange(tt.cpRange)
			if err == nil {
				t.Errorf("expandCodePointRange() didn't produce an error when expected")

				return
			}
			if got != nil {
				t.Errorf("expandCodePointRange() returned a rune slice when not expected: %v", got)

				return
			}

			if !reflect.DeepEqual(tt.wantErrMsg, err.Error()) {
				t.Errorf("expandCodePointRange() got = %v, want %v", err, tt.wantErrMsg)
			}
			if !errors.Is(err, ErrInvalidCodePoint) {
				t.Errorf("expandCodePointRange() didn't get expected ErrInvalidCodePoint, got: %v", err)
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
		name    string
		line    string
		want    []Emoji
		wantErr bool
	}{
		{
			name: "Basic 2 byte rune",
			line: "23F0          ; Basic_Emoji                  ; alarm clock                                                    # E0.6   [1] (⏰)",
			want: []Emoji{
				{
					Runes: []rune{'\u23f0'},
				},
			},
			wantErr: false,
		},
		{
			name: "Basic 2.5 byte rune",
			line: "1F004         ; Basic_Emoji                  ; mahjong red dragon                                             # E0.6   [1] (🀄)",
			want: []Emoji{
				{
					Runes: []rune{'\U0001f004'},
				},
			},
			wantErr: false,
		},
		{
			name: "Basic 2 byte rune range",
			line: "23E9..23EC    ; Basic_Emoji                  ; fast-forward button                                            # E0.6   [4] (⏩..⏬)",
			want: []Emoji{
				{
					Runes: []rune{'\u23e9'},
				},
				{
					Runes: []rune{'\u23ea'},
				},
				{
					Runes: []rune{'\u23eb'},
				},
				{
					Runes: []rune{'\u23ec'},
				},
			},
			wantErr: false,
		},
		{
			name: "2.5 byte rune range",
			line: "1F4F0..1F4F4  ; Basic_Emoji                  ; newspaper                                                      # E0.6   [5] (📰..📴)",
			want: []Emoji{
				{
					Runes: []rune{'\U0001f4f0'},
				},
				{
					Runes: []rune{'\U0001f4f1'},
				},
				{
					Runes: []rune{'\U0001f4f2'},
				},
				{
					Runes: []rune{'\U0001f4f3'},
				},
				{
					Runes: []rune{'\U0001f4f4'},
				},
			},
			wantErr: false,
		},
		{
			name:    "Invalid rune range",
			line:    "1F4F0..1F4FG  ; Basic_Emoji                  ; newspaper                                                      # E0.6   [5] (📰..📴)",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseEmoji(tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseEmoji() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseEmoji() got = %v, want %v", got, tt.want)
			}
		})
	}
}
