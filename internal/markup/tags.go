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
	"sort"
	"strconv"
	"strings"
)

// openTag records an unclosed tag together with the appearance to restore when
// it is closed.
type openTag struct {
	name   string
	parent style
	column int
}

func (c *compiler) tag(i int) (int, error) {
	end := strings.IndexByte(c.raw[i:], '>')
	if end < 0 {
		return 0, c.fail(ErrMarkup, "unterminated tag", i, 1).
			Hint(`close it with '>' or escape the angle bracket as \<`)
	}

	end += i
	inner := c.raw[i+1 : end]
	width := end - i + 1

	if inner == "" {
		return 0, c.fail(ErrMarkup, "empty tag", i, width).
			Hint(`escape the angle bracket as \< to write it literally`)
	}

	if closing := strings.HasPrefix(inner, "/"); closing {
		if err := c.closeTag(inner[1:], i, width); err != nil {
			return 0, err
		}

		return width, nil
	}

	if err := c.openTag(inner, i, width); err != nil {
		return 0, err
	}

	return width, nil
}

func (c *compiler) openTag(inner string, offset, width int) error {
	st, canonical, void, err := c.resolve(inner, offset, width)
	if err != nil {
		return err
	}

	if void {
		// A reset is a standalone directive: it clears everything and then
		// applies whatever was combined with it, e.g. <reset:white>.
		c.buf = append(c.buf, Reset...)
		c.buf = appendStyle(c.buf, st)
		c.current = st

		return nil
	}

	if len(c.stack) >= c.opts.MaxNesting {
		return c.fail(ErrLimit, "tag nesting deeper than "+strconv.Itoa(c.opts.MaxNesting), offset, width).
			Hint("flatten the message or raise WithLimits(Limits{MaxNesting: …})")
	}

	parent := c.current
	next := parent.merge(st)
	c.buf = appendTransition(c.buf, parent, next)
	c.current = next
	c.stack = append(c.stack, openTag{name: canonical, parent: parent, column: c.columnAt(offset)})

	return nil
}

func (c *compiler) closeTag(inner string, offset, width int) error {
	if len(c.stack) == 0 {
		d := c.fail(ErrMarkup, "closing tag without a matching opening tag", offset, width)
		if inner == "reset" {
			d.Hint("<reset> is self contained and must not be closed")
		}

		return d.Hint("every </tag> needs an earlier <tag>")
	}

	top := c.stack[len(c.stack)-1]
	if inner != "" {
		canonical, err := c.canonicalise(inner, offset, width)
		if err != nil {
			return err
		}

		if canonical != top.name {
			return c.fail(ErrMarkup, "mismatched closing tag </"+inner+">", offset, width).
				Field("expected", "</"+top.name+"> opened at column "+strconv.Itoa(top.column)).
				Hint("close inner tags before outer ones, or use </> to close the innermost tag")
		}
	}

	c.buf = appendTransition(c.buf, c.current, top.parent)
	c.current = top.parent
	c.stack = c.stack[:len(c.stack)-1]

	return nil
}

// resolve turns the body of a tag into a style. Components are classified by
// the vocabulary rather than by their position, so <bold:red> and <red:bold>
// are the same tag.
func (c *compiler) resolve(inner string, offset, width int) (st style, canonical string, void bool, err error) {
	parts := strings.Split(inner, ":")
	names := make([]string, 0, len(parts))
	colourSet := false

	for _, part := range parts {
		if part == "" {
			return style{}, "", false, c.fail(ErrMarkup, "empty component in tag <"+inner+">", offset, width).
				Hint("combine components as <bold:red>, without empty segments")
		}

		name := strings.ToLower(part)
		tok, ok := lookupToken(name)
		if !ok {
			d := c.fail(ErrMarkup, "unknown tag component \""+part+"\"", offset, width)
			if alt := suggest(name); alt != "" {
				d.Hint("did you mean <" + alt + ">?")
			}

			return style{}, "", false, d
		}

		switch tok.kind {
		case kindColour:
			if colourSet {
				return style{}, "", false, c.fail(ErrMarkup, "more than one colour in tag <"+inner+">", offset, width).
					Hint("a run of text can only carry one colour")
			}
			st.colour = tok.code
			colourSet = true
		case kindFormat:
			st.flags |= tok.flag
		case kindReset:
			void = true
		}
		names = append(names, name)
	}

	sort.Strings(names)

	return st, strings.Join(names, ":"), void, nil
}

// canonicalise validates a closing tag body and returns its canonical form.
func (c *compiler) canonicalise(inner string, offset, width int) (string, error) {
	_, canonical, _, err := c.resolve(inner, offset, width)
	return canonical, err
}
