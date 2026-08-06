// mctext (https://github.com/JoseLuisHD/mctext)
//
// Copyright 2026 JoseLuisHD
//
// Licensed under the Apache License, Version 2.0 (the "License").
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package markup

import "testing"

func TestAppendTransition(t *testing.T) {
	var (
		none      = style{}
		red       = style{colour: 'c'}
		gold      = style{colour: '6'}
		goldBold  = style{colour: '6', flags: flagBold}
		redBold   = style{colour: 'c', flags: flagBold}
		boldOnly  = style{flags: flagBold}
		bothMarks = style{flags: flagBold | flagItalic}
	)

	cases := []struct {
		name string
		from style
		to   style
		want string
	}{
		{"no change", gold, gold, ""},
		{"enter a colour", none, red, "§c"},
		{"leave everything", red, none, "§r"},
		{"add a decoration to the same colour", gold, goldBold, "§l"},
		{"drop a decoration keeping the colour", goldBold, gold, "§6"},
		{"change colour", gold, red, "§c"},
		{"change colour carrying a decoration", goldBold, redBold, "§c§l"},
		{"colour to decoration only", gold, boldOnly, "§r§l"},
		{"decoration only, additive", boldOnly, bothMarks, "§o"},
		{"decoration only, subtractive", bothMarks, boldOnly, "§r§l"},
		{"enter colour and decorations at once", none, redBold, "§c§l"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := string(appendTransition(nil, tc.from, tc.to)); got != tc.want {
				t.Errorf("appendTransition(%+v -> %+v) = %q, want %q", tc.from, tc.to, got, tc.want)
			}
		})
	}
}

func TestStyleMerge(t *testing.T) {
	parent := style{colour: '6', flags: flagBold}

	if got := parent.merge(style{flags: flagItalic}); got != (style{colour: '6', flags: flagBold | flagItalic}) {
		t.Errorf("decorations must accumulate, got %+v", got)
	}

	if got := parent.merge(style{colour: 'c'}); got != (style{colour: 'c', flags: flagBold}) {
		t.Errorf("a colour must override while decorations survive, got %+v", got)
	}

	if !(style{}).empty() || (style{flags: flagBold}).empty() {
		t.Error("empty() must only hold for the zero style")
	}
}

func TestSuggest(t *testing.T) {
	cases := map[string]string{
		"redd":       "red",
		"gren":       "green",
		"bld":        "bold",
		"dark_bleu":  "dark_blue",
		"reset":      "reset",
		"nonsense42": "",
	}
	for unknown, want := range cases {
		if got := suggest(unknown); got != want {
			t.Errorf("suggest(%q) = %q, want %q", unknown, got, want)
		}
	}
}

func TestVocabularyIsBidirectional(t *testing.T) {
	for name, tok := range tokens {
		var code byte
		switch tok.kind {
		case kindColour:
			code = tok.code
		case kindFormat:
			code = flagCode(tok.flag)
		case kindReset:
			code = 'r'
		}
		if code == 0 {
			t.Errorf("tag %q resolves to no client code", name)
			continue
		}
		if got, ok := TagFor(code); !ok || got != name {
			t.Errorf("TagFor(%q) = %q, want %q", string(code), got, name)
		}
	}
}
