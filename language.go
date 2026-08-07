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

// Language is a BCP-47 style tag such as "en", "es" or "pt-br". It is a plain
// comparable value so it can be used as a map key and passed around freely.
type Language string

// The languages shipped by default. Any other tag can be added at runtime with
// RegisterLanguage, before Init is called.
const (
	English    Language = "en"
	Spanish    Language = "es"
	French     Language = "fr"
	German     Language = "de"
	Italian    Language = "it"
	Portuguese Language = "pt"
	Russian    Language = "ru"
	Japanese   Language = "ja"
	Korean     Language = "ko"
	Chinese    Language = "zh"
	Arabic     Language = "ar"
	Hindi      Language = "hi"
	Turkish    Language = "tr"
	Dutch      Language = "nl"
	Polish     Language = "pl"
	Swedish    Language = "sv"
	Norwegian  Language = "no"
	Danish     Language = "da"
	Finnish    Language = "fi"
	Czech      Language = "cs"
	Greek      Language = "el"
	Hungarian  Language = "hu"
	Romanian   Language = "ro"
	Ukrainian  Language = "uk"
	Thai       Language = "th"
	Vietnamese Language = "vi"
	Indonesian Language = "id"
	Malay      Language = "ms"
	Hebrew     Language = "he"
	Persian    Language = "fa"
)

func (l Language) String() string {
	return string(l)
}

func (l Language) DisplayName() string {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	if name, ok := registry.names[l]; ok {
		return name
	}

	return string(l)
}

func (l Language) Registered() bool {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	_, ok := registry.names[l]

	return ok
}
