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

// Field is an extra labelled value rendered as part of a diagnostic report.
type Field struct {
	Label string
	Value string
}

// Diagnostic is the error type returned by every operation that can fail while
// loading or rendering a message. It renders as a compact, self contained
// report: what went wrong, where, and how to fix it.
type Diagnostic struct {
	Summary  string
	Loc      Location
	Snippet  string
	Caret    int // 1-based rune column inside Snippet
	CaretLen int
	Fields   []Field
	Hints    []string
	Frames   []Frame
	Err      error
}

func New(err error, summary string) *Diagnostic {
	return &Diagnostic{Summary: summary, Err: err}
}

func (d *Diagnostic) At(loc Location) *Diagnostic {
	d.Loc = loc
	return d
}

// Quote attaches the offending source text and the column to underline.
func (d *Diagnostic) Quote(snippet string, column, length int) *Diagnostic {
	d.Snippet = snippet
	d.Caret = column
	d.CaretLen = length
	return d
}

func (d *Diagnostic) Field(label, value string) *Diagnostic {
	if value != "" {
		d.Fields = append(d.Fields, Field{Label: label, Value: value})
	}

	return d
}

func (d *Diagnostic) Hint(hint string) *Diagnostic {
	if hint != "" {
		d.Hints = append(d.Hints, hint)
	}

	return d
}

func (d *Diagnostic) Stack(depth int, ignorePrefix string) *Diagnostic {
	d.Frames = Capture(1, depth, ignorePrefix)
	return d
}

func (d *Diagnostic) Unwrap() error {
	return d.Err
}
