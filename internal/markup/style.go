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

// style is the fully resolved appearance of a run of text.
type style struct {
	colour byte // zero means "no colour of its own", i.e. the client default
	flags  flags
}

func (s style) empty() bool {
	return s.colour == 0 && s.flags == 0
}

// merge layers other on top of s: a colour overrides, decorations accumulate.
func (s style) merge(other style) style {
	if other.colour != 0 {
		s.colour = other.colour
	}

	s.flags |= other.flags

	return s
}

// appendStyle writes the full sequence for s, colour first.
func appendStyle(dst []byte, s style) []byte {
	if s.colour != 0 {
		dst = append(dst, Section...)
		dst = append(dst, s.colour)
	}

	return appendFlags(dst, s.flags)
}

// appendFlags writes decorations in the canonical order.
func appendFlags(dst []byte, f flags) []byte {
	for _, bit := range orderedFlags {
		if f&bit != 0 {
			dst = append(dst, Section...)
			dst = append(dst, flagCode(bit))
		}
	}

	return dst
}

// appendTransition writes the shortest correct sequence that moves the client
// from the "from" appearance to the "to" appearance.
func appendTransition(dst []byte, from, to style) []byte {
	if from == to {
		return dst
	}

	if to.empty() {
		return append(dst, Reset...)
	}

	if from.colour == to.colour && to.flags&from.flags == from.flags {
		// Same colour and only decorations added: they stack for free.
		return appendFlags(dst, to.flags&^from.flags)
	}

	if to.colour != 0 {
		// Emitting a colour clears the previous decorations by itself, so no
		// explicit reset is needed before re-applying the new style.
		return appendStyle(dst, to)
	}

	// The destination has no colour of its own, so the only way to reach it is
	// a full reset followed by its decorations.
	dst = append(dst, Reset...)

	return appendFlags(dst, to.flags)
}
