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

func (d *Diagnostic) Error() string {
	var b strings.Builder
	b.Grow(256)
	b.WriteString("translation: ")
	b.WriteString(d.Summary)

	if d.Loc.File != "" {
		writeField(&b, "at", d.Loc.String())
	}

	for _, f := range d.Fields {
		writeField(&b, f.Label, f.Value)
	}

	if d.Snippet != "" {
		b.WriteString("\n  │ ")
		b.WriteString(expandTabs(d.Snippet))

		if d.Caret > 0 {
			b.WriteString("\n  │ ")
			b.WriteString(strings.Repeat(" ", d.Caret-1))
			b.WriteString(strings.Repeat("^", max(d.CaretLen, 1)))
		}
	}

	for _, h := range d.Hints {
		writeField(&b, "hint", h)
	}

	if len(d.Frames) > 0 {
		b.WriteString("\n  stack:")
		for _, f := range d.Frames {
			b.WriteString("\n    ")
			b.WriteString(f.Function)
			b.WriteString(" (")
			b.WriteString(f.File)
			b.WriteByte(':')
			b.WriteString(strconv.Itoa(f.Line))
			b.WriteByte(')')
		}
	}

	return b.String()
}

func writeField(b *strings.Builder, label, value string) {
	b.WriteString("\n  ")
	b.WriteString(label)
	b.WriteByte(':')

	// Align the values of the well known short labels.
	for i := len(label); i < 9; i++ {
		b.WriteByte(' ')
	}

	b.WriteString(value)
}

func expandTabs(s string) string {
	if !strings.ContainsRune(s, '\t') {
		return s
	}

	return strings.ReplaceAll(s, "\t", "    ")
}
