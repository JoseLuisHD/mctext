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

import (
	"strings"

	"github.com/JoseLuisHD/mctext/internal/catalog"
	"github.com/JoseLuisHD/mctext/internal/markup"
)

// snapshot is the immutable state a Manager serves lookups from. Publishing a
// new one is a single atomic pointer store, so readers never take a lock and
// never observe a half loaded catalog.
type snapshot struct {
	catalogs        map[Language]*catalog.Catalog
	resolution      map[Language][]Language
	defaultLanguage Language
	languages       []Language
	messages        int
}

// lookup resolves a key in a language, falling back when configured.
func (s *snapshot) lookup(language Language, key string) (*markup.Program, Language, bool) {
	chain, ok := s.resolution[language]
	if !ok {
		// An unknown language still resolves through the default chain so a
		// bad tag degrades to readable text instead of raw keys.
		chain = s.resolution[s.defaultLanguage]
	}

	for _, candidate := range chain {
		if c, exists := s.catalogs[candidate]; exists {
			if program, found := c.Lookup(key); found {
				return program, candidate, true
			}
		}
	}

	return nil, language, false
}

// nearestKey finds a declared key sharing the longest namespace prefix with the
// missing one, which usually points straight at the typo.
func (s *snapshot) nearestKey(language Language, key string) string {
	c, ok := s.catalogs[language]

	// The scan is linear in the catalog size, so it is only worth running on
	// catalogs small enough for the suggestion to stay cheap on an error path.
	if !ok || c.Len() > 5000 {
		return ""
	}

	prefix := key
	for {
		index := strings.LastIndex(prefix, "::")
		if index < 0 {
			return ""
		}

		prefix = prefix[:index]
		var best string
		for _, candidate := range c.Keys() {
			if strings.HasPrefix(candidate, prefix+"::") {
				if best == "" || len(candidate) < len(best) {
					best = candidate
				}
			}
		}

		if best != "" {
			return best
		}
	}
}
