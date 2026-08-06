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

package diag

import (
	"strconv"
	"strings"
)

// Location identifies a position inside a translation source file. A zero Line
// means the position is unknown; a zero Column means the whole line is meant.
type Location struct {
	File   string
	Line   int
	Column int
}

func (l Location) String() string {
	if l.File == "" {
		return "<unknown>"
	}

	var b strings.Builder
	b.WriteString(l.File)

	if l.Line > 0 {
		b.WriteByte(':')
		b.WriteString(strconv.Itoa(l.Line))

		if l.Column > 0 {
			b.WriteByte(':')
			b.WriteString(strconv.Itoa(l.Column))
		}
	}

	return b.String()
}
