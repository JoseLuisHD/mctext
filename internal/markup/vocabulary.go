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

import "sort"

// kindcthe three families of tag components.
type kind uint8

const (
	kindColour kind = iota
	kindFormat
	kindReset
)

// token is a single resolved tag component, e.g. the "bold" in <bold:red>.
type token struct {
	kind kind
	code byte  // client code for colours, e.g. 'c' for red
	flag flags // decoration bit for formats
}

var tokens = map[string]token{
	// Base colours.
	"black":        {kind: kindColour, code: '0'},
	"dark_blue":    {kind: kindColour, code: '1'},
	"dark_green":   {kind: kindColour, code: '2'},
	"dark_aqua":    {kind: kindColour, code: '3'},
	"dark_red":     {kind: kindColour, code: '4'},
	"dark_purple":  {kind: kindColour, code: '5'},
	"gold":         {kind: kindColour, code: '6'},
	"gray":         {kind: kindColour, code: '7'},
	"dark_gray":    {kind: kindColour, code: '8'},
	"blue":         {kind: kindColour, code: '9'},
	"green":        {kind: kindColour, code: 'a'},
	"aqua":         {kind: kindColour, code: 'b'},
	"red":          {kind: kindColour, code: 'c'},
	"light_purple": {kind: kindColour, code: 'd'},
	"yellow":       {kind: kindColour, code: 'e'},
	"white":        {kind: kindColour, code: 'f'},
	"light_blue":   {kind: kindColour, code: 'w'},

	// Extended and material colours.
	"minecoin_gold":      {kind: kindColour, code: 'g'},
	"material_quartz":    {kind: kindColour, code: 'h'},
	"material_iron":      {kind: kindColour, code: 'i'},
	"material_netherite": {kind: kindColour, code: 'j'},
	"material_redstone":  {kind: kindColour, code: 'm'},
	"material_copper":    {kind: kindColour, code: 'n'},
	"material_gold":      {kind: kindColour, code: 'p'},
	"material_emerald":   {kind: kindColour, code: 'q'},
	"material_diamond":   {kind: kindColour, code: 's'},
	"material_lapis":     {kind: kindColour, code: 't'},
	"material_amethyst":  {kind: kindColour, code: 'u'},
	"material_resin":     {kind: kindColour, code: 'v'},

	// Decorations.
	"obfuscated": {kind: kindFormat, flag: flagObfuscated},
	"bold":       {kind: kindFormat, flag: flagBold},
	"italic":     {kind: kindFormat, flag: flagItalic},

	// Control.
	"reset": {kind: kindReset},
}

// codeNames is the reverse index, used to translate a raw client code found in
// a source file back into the tag the author should have written.
var codeNames = func() map[byte]string {
	m := make(map[byte]string, len(tokens))
	for name, tok := range tokens {
		switch tok.kind {
		case kindColour:
			m[tok.code] = name
		case kindFormat:
			m[flagCode(tok.flag)] = name
		case kindReset:
			m['r'] = name
		}
	}

	return m
}()

// tagNames is the sorted vocabulary, used for "did you mean" suggestions.
var tagNames = func() []string {
	names := make([]string, 0, len(tokens))
	for name := range tokens {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}()

// lookupToken resolves a lower-cased tag component.
func lookupToken(name string) (token, bool) {
	tok, ok := tokens[name]
	return tok, ok
}

// TagFor returns the tag name matching a raw client code, e.g. 'c' -> "red".
func TagFor(code byte) (string, bool) {
	name, ok := codeNames[code]
	return name, ok
}
