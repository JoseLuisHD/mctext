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

// Limits bounds the resources a single file may consume. They exist so that a
// corrupted or hostile file fails loudly instead of exhausting memory.
type Limits struct {
	MaxFileSize  int
	MaxEntries   int
	MaxKeyLength int
	MaxLineWidth int
}

// DefaultLimits returns limits that are generous for real content and still
// small enough to contain accidents.
func DefaultLimits() Limits {
	return Limits{
		MaxFileSize:  8 << 20,
		MaxEntries:   50_000,
		MaxKeyLength: 512,
		MaxLineWidth: 16 << 10,
	}
}

func (l Limits) normalized() Limits {
	d := DefaultLimits()
	if l.MaxFileSize <= 0 {
		l.MaxFileSize = d.MaxFileSize
	}

	if l.MaxEntries <= 0 {
		l.MaxEntries = d.MaxEntries
	}

	if l.MaxKeyLength <= 0 {
		l.MaxKeyLength = d.MaxKeyLength
	}

	if l.MaxLineWidth <= 0 {
		l.MaxLineWidth = d.MaxLineWidth
	}

	return l
}
