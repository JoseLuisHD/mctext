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

import "github.com/JoseLuisHD/mctext/internal/diag"

// Program is the immutable result of compiling one source message. Every tag
// and every placeholder is resolved once, at load time, so rendering never
// parses anything: it walks a flat slice and concatenates bytes.
//
// A Program is safe for concurrent use because nothing mutates it after
// Compile returns.
type Program struct {
	key      string
	raw      string
	static   string
	segments []Segment
	args     []string
	sizeHint int
	origin   diag.Location
}

func (p *Program) Key() string {
	return p.key
}

func (p *Program) Raw() string {
	return p.raw
}

func (p *Program) Origin() diag.Location {
	return p.origin
}

func (p *Program) Args() []string {
	return p.args
}

func (p *Program) Static() (string, bool) {
	return p.static, len(p.args) == 0
}

func (p *Program) SizeHint() int {
	return p.sizeHint
}

// ArgIndex returns the position of a named argument, or -1 when the message
// does not declare it.
func (p *Program) ArgIndex(name string) int {
	for i, arg := range p.args {
		if arg == name {
			return i
		}
	}

	return -1
}

// AppendTo renders the message into dst. values must hold one entry per
// declared argument, in the order reported by Args.
func (p *Program) AppendTo(dst []byte, values []string) []byte {
	if len(p.args) == 0 {
		return append(dst, p.static...)
	}

	for i := range p.segments {
		seg := &p.segments[i]
		if seg.Kind == SegmentLiteral {
			dst = append(dst, seg.Text...)
			continue
		}

		if seg.Index < len(values) {
			dst = append(dst, values[seg.Index]...)
		}
	}

	return dst
}
