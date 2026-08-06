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
	"sort"
	"strings"
	"sync"
)

var registry = struct {
	mu    sync.RWMutex
	names map[Language]string
}{
	names: map[Language]string{
		English: "English",
		Spanish: "Español",
	},
}

func RegisterLanguage(tag, displayName string) (Language, error) {
	normalized := strings.ToLower(strings.TrimSpace(tag))
	if err := validateTag(normalized); err != nil {
		return "", err
	}

	if strings.TrimSpace(displayName) == "" {
		displayName = normalized
	}

	registry.mu.Lock()
	defer registry.mu.Unlock()

	if existing, ok := registry.names[Language(normalized)]; ok && existing != displayName {
		return Language(normalized), newError(ErrInvalidLanguage,
			"language "+normalized+" is already registered as "+existing)
	}

	registry.names[Language(normalized)] = displayName

	return Language(normalized), nil
}

func RegisteredLanguages() []Language {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	languages := make([]Language, 0, len(registry.names))
	for language := range registry.names {
		languages = append(languages, language)
	}

	sort.Slice(languages, func(i, j int) bool { return languages[i] < languages[j] })

	return languages
}
