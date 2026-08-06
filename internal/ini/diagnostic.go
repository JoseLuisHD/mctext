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

package ini

import "github.com/JoseLuisHD/mctext/internal/diag"

// builder is a small helper that keeps the error construction sites terse.
type builder struct {
	location diag.Location
	snippet  string
	column   int
	width    int
}

func at(file string, line, column int) *builder {
	return &builder{location: diag.Location{File: file, Line: line, Column: column}}
}

func (b *builder) Quote(snippet string, column, width int) *builder {
	b.snippet, b.column, b.width = snippet, column, width
	return b
}

func (b *builder) Wrap(err error, summary string) *diag.Diagnostic {
	d := diag.New(err, summary).At(b.location)
	if b.snippet != "" {
		d.Quote(b.snippet, b.column, b.width)
	}

	return d
}
