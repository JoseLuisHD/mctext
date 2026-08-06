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

package mctext

import "path"

// Source is the port through which translation files are read. Depending on an
// interface rather than on the file system is what lets the loader run against
// an embedded FS in production and against memory in tests, with no change to
// any other layer.
type Source interface {
	// ReadFile returns the raw content of a slash separated relative path.
	ReadFile(name string) ([]byte, error)
	// Describe returns a human readable origin, used in diagnostics.
	Describe() string
}

// PathResolver maps a language and a declared file name onto a path inside the
// Source. The default resolver expects one directory per language.
type PathResolver func(language Language, file string) string

// NestedLayout resolves "common.ini" for Spanish as "es/common.ini".
func NestedLayout(language Language, file string) string {
	return path.Join(string(language), file)
}

// FlatLayout resolves file names relative to the source root unchanged, so the
// language has to be part of the declared name.
func FlatLayout(_ Language, file string) string {
	return file
}
