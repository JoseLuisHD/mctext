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

import "strings"

// validateTag accepts a primary subtag of two or three letters followed by any
// number of alphanumeric subtags, e.g. "es", "pt-br", "zh-hant".
func validateTag(tag string) error {
	if tag == "" {
		return newError(ErrInvalidLanguage, "empty language tag")
	}

	subtags := strings.Split(tag, "-")
	for i, subtag := range subtags {
		if i == 0 {
			if len(subtag) < 2 || len(subtag) > 3 || !allLetters(subtag) {
				return newError(ErrInvalidLanguage,
					"language tag must start with a 2 or 3 letter code: "+tag)
			}

			continue
		}

		if len(subtag) < 2 || len(subtag) > 8 || !alphanumeric(subtag) {
			return newError(ErrInvalidLanguage, "malformed subtag in language tag: "+tag)
		}
	}

	return nil
}

func allLetters(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < 'a' || s[i] > 'z' {
			return false
		}
	}

	return true
}

func alphanumeric(s string) bool {
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if (ch < 'a' || ch > 'z') && (ch < '0' || ch > '9') {
			return false
		}
	}

	return true
}
