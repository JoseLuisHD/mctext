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
	"strconv"
	"unicode/utf8"

	"github.com/JoseLuisHD/mctext/internal/diag"
)

func (c *compiler) unclosed() error {
	top := c.stack[len(c.stack)-1]
	d := c.fail(ErrMarkup, "unclosed tag <"+top.name+">", 0, 0)
	d.Loc.Column = c.origin.Column + top.column - 1
	d.Caret = top.column
	d.CaretLen = len(top.name) + 2

	if len(c.stack) > 1 {
		d.Field("open tags", strconv.Itoa(len(c.stack)))
	}

	return d.Hint("add </" + top.name + "> or </> before the end of the message")
}

// fail builds a diagnostic anchored at a byte offset inside the raw value.
func (c *compiler) fail(err error, summary string, offset, width int) *diag.Diagnostic {
	column := c.columnAt(offset)
	loc := c.origin

	if loc.Column > 0 {
		loc.Column += column - 1
	} else {
		loc.Column = column
	}

	return diag.New(err, summary).
		At(loc).
		Field("key", c.key).
		Quote(c.raw, column, width)
}

// columnAt converts a byte offset inside the value into a 1-based rune column.
func (c *compiler) columnAt(offset int) int {
	if offset <= 0 {
		return 1
	}

	if offset > len(c.raw) {
		offset = len(c.raw)
	}

	return utf8.RuneCountInString(c.raw[:offset]) + 1
}
