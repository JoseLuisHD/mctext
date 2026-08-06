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
	"errors"
	"strconv"
	"strings"
)

func (c *compiler) placeholder(i int) (int, error) {
	end := strings.IndexByte(c.raw[i:], '}')
	if end < 0 {
		return 0, c.fail(ErrPlaceholder, "unterminated argument placeholder", i, 1).
			Hint(`close it with '}' or escape the brace as \{`)
	}

	end += i
	name := c.raw[i+1 : end]
	width := end - i + 1

	if err := validArgumentName(name); err != nil {
		return 0, c.fail(ErrPlaceholder, err.Error(), i, width).
			Hint("use letters, digits and underscores, e.g. {player_name}")
	}

	index := -1
	for k, existing := range c.args {
		if existing == name {
			index = k
			break
		}
	}

	if index < 0 {
		if len(c.args) >= c.opts.MaxArguments {
			return 0, c.fail(ErrLimit, "more than "+strconv.Itoa(c.opts.MaxArguments)+" distinct arguments", i, width)
		}

		index = len(c.args)
		c.args = append(c.args, name)
	}

	c.flushLiteral()
	c.segments = append(c.segments, Segment{Kind: SegmentArgument, Index: index})

	return width, nil
}

func validArgumentName(name string) error {
	if name == "" {
		return errors.New("empty argument name")
	}

	for i := 0; i < len(name); i++ {
		ch := name[i]
		switch {
		case ch >= 'a' && ch <= 'z', ch >= 'A' && ch <= 'Z', ch == '_':
		case ch >= '0' && ch <= '9':
			if i == 0 {
				return errors.New("argument name starts with a digit: " + name)
			}
		default:
			return errors.New("illegal character in argument name: " + name)
		}
	}

	return nil
}
