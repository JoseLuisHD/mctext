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

// Section is the formatting prefix understood by the client. It is two bytes
// wide in UTF-8, which every size computation accounts for.
const Section = "\u00a7"

// Reset is the sequence that clears colour and decorations at once.
const Reset = Section + "r"

// flags is the bit set of the additive text decorations. Decorations are
// additive because, unlike colours, emitting one does not clear the others.
type flags uint8

const (
	flagObfuscated flags = 1 << iota
	flagBold
	flagItalic
)

// orderedFlags fixes the emission order of decorations. A canonical order means
// two identical styles always produce byte-identical output, which is what
// makes the golden tests meaningful.
var orderedFlags = [...]flags{flagObfuscated, flagBold, flagItalic}

// flagCode maps a decoration bit onto its client code.
func flagCode(f flags) byte {
	switch f {
	case flagObfuscated:
		return 'k'
	case flagBold:
		return 'l'
	case flagItalic:
		return 'o'
	default:
		return 0
	}
}
