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

import (
	"errors"
	"strings"
)

// Separator is the namespace separator used inside keys.
const Separator = "::"

// findSeparator locates the delimiter between key and value. A run of exactly
// one ':' delimits; a run of two is part of the key namespace.
func findSeparator(line string) (index, width int) {
	for i := 0; i < len(line); i++ {
		switch line[i] {
		case '=':
			return i, 1
		case ':':
			run := 1
			for i+run < len(line) && line[i+run] == ':' {
				run++
			}

			if run == 1 {
				return i, 1
			}

			if run == 2 {
				i += run - 1
				continue
			}

			// Three or more colons in a row cannot be disambiguated.
			return -1, 0
		}
	}

	return -1, 0
}

func validateKey(key string) error {
	if key == "" {
		return errors.New("empty key")
	}

	if strings.HasPrefix(key, Separator) || strings.HasSuffix(key, Separator) {
		return errors.New("key must not start or end with '::': " + key)
	}

	for i := 0; i < len(key); i++ {
		ch := key[i]
		switch {
		case ch >= 'a' && ch <= 'z', ch >= 'A' && ch <= 'Z', ch >= '0' && ch <= '9':
		case ch == '_', ch == '-', ch == '.', ch == ':':
		default:
			return errors.New("illegal character in key: " + key)
		}
	}

	return nil
}
