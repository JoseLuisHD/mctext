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

import (
	"strings"

	"github.com/JoseLuisHD/mctext/internal/diag"
)

type compiler struct {
	key    string
	raw    string
	origin diag.Location
	opts   Options

	buf      []byte
	segments []Segment
	args     []string
	stack    []openTag
	current  style
}

// Compile turns one source message into an immutable Program. Every colour tag
// is resolved to client codes here so that rendering stays allocation free.
//
// origin.Column must be the column at which the value starts inside its source
// line; reported positions are offset from it.
func Compile(key, raw string, origin diag.Location, opts Options) (*Program, error) {
	c := &compiler{
		key:    key,
		raw:    raw,
		origin: origin,
		opts:   opts.normalized(),
		buf:    make([]byte, 0, len(raw)+16),
	}

	return c.run()
}

func (c *compiler) run() (*Program, error) {
	for i := 0; i < len(c.raw); {
		switch c.raw[i] {
		case '\\':
			n, err := c.escape(i)
			if err != nil {
				return nil, err
			}

			i += n
		case '<':
			n, err := c.tag(i)
			if err != nil {
				return nil, err
			}

			i += n
		case '{':
			n, err := c.placeholder(i)
			if err != nil {
				return nil, err
			}

			i += n
		default:
			if strings.HasPrefix(c.raw[i:], Section) {
				n, err := c.legacyCode(i)
				if err != nil {
					return nil, err
				}

				i += n

				continue
			}
			i += c.literalRun(i)
		}
	}

	if len(c.stack) > 0 {
		return nil, c.unclosed()
	}

	return c.finish(), nil
}

// literalRun copies the longest span of ordinary text in one go, which keeps
// the common case down to a single append per run.
func (c *compiler) literalRun(start int) int {
	i := start
	for i < len(c.raw) {
		switch c.raw[i] {
		case '\\', '<', '{':
			c.buf = append(c.buf, c.raw[start:i]...)
			return i - start
		}

		if c.raw[i] == Section[0] && strings.HasPrefix(c.raw[i:], Section) {
			c.buf = append(c.buf, c.raw[start:i]...)
			return i - start
		}

		i++
	}

	c.buf = append(c.buf, c.raw[start:]...)

	return i - start
}

func (c *compiler) escape(i int) (int, error) {
	if i+1 >= len(c.raw) {
		return 0, c.fail(ErrMarkup, "dangling escape character at end of message", i, 1).
			Hint(`write \\ for a literal backslash`)
	}

	switch ch := c.raw[i+1]; ch {
	case '\\', '<', '>', '{', '}':
		c.buf = append(c.buf, ch)
	case 'n':
		c.buf = append(c.buf, '\n')
	case 't':
		c.buf = append(c.buf, '\t')
	default:
		return 0, c.fail(ErrMarkup, "unknown escape sequence \\"+string(ch), i, 2).
			Hint(`valid escapes are \\ \< \> \{ \} \n \t`)
	}

	return 2, nil
}

func (c *compiler) legacyCode(i int) (int, error) {
	if !c.opts.AllowLegacyCodes {
		width := len(Section)
		summary := "raw formatting code in message"
		d := c.fail(ErrMarkup, summary, i, width+1)

		if next := i + len(Section); next < len(c.raw) {
			if name, ok := TagFor(c.raw[next]); ok {
				d.Hint("write <" + name + ">…</" + name + "> instead of " + Section + string(c.raw[next]))
			}
		}

		return 0, d.Hint("raw codes escape the nesting model and break the surrounding style")
	}

	n := len(Section)
	if i+n < len(c.raw) {
		n++
	}

	c.buf = append(c.buf, c.raw[i:i+n]...)

	return n, nil
}

func (c *compiler) flushLiteral() {
	if len(c.buf) == 0 {
		return
	}

	c.segments = append(c.segments, Segment{Kind: SegmentLiteral, Text: string(c.buf)})
	c.buf = c.buf[:0]
}

func (c *compiler) finish() *Program {
	p := &Program{
		key:    c.key,
		raw:    c.raw,
		args:   c.args,
		origin: c.origin,
	}

	if len(c.args) == 0 {
		p.static = string(c.buf)
		p.sizeHint = len(p.static)
		return p
	}

	c.flushLiteral()
	p.segments = c.segments
	for i := range p.segments {
		p.sizeHint += len(p.segments[i].Text)
	}

	p.sizeHint += 12 * len(c.args)

	return p
}
