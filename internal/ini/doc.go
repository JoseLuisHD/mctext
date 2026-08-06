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

// Package ini implements the translation file format: a line oriented
// "key: value" dialect whose keys may contain the :: namespace separator.
//
//	# comment            ; comment
//	[section::prefix]    optional key prefix until the next section
//	some::key: value     ':' or '=' separates key and value
//	long::key: first \   a trailing backslash continues on the next line
//	                     second
//
// A single ':' terminates the key, so keys may embed the '::' separator freely.
// Values keep their inner spacing; surrounding whitespace is trimmed unless the
// value is wrapped in double quotes.

package ini
