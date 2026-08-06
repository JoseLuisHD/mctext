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

// It is a thin, immutable view: creating one allocates a two word struct and
// copying it is free.
type Localizer struct {
	manager  *Manager
	language Language
}

func (l *Localizer) Language() Language {
	return l.language
}

func (l *Localizer) Get(key string, args Args) string {
	if l.manager == nil {
		return key
	}

	return l.manager.Get(l.language, key, args)
}

func (l *Localizer) GetPairs(key string, pairs ...Pair) string {
	if l.manager == nil {
		return key
	}

	return l.manager.GetPairs(l.language, key, pairs...)
}

func (l *Localizer) TryGet(key string, args Args) (string, error) {
	if l.manager == nil {
		return key, newError(ErrNotInitialised, "localizer is not bound to a manager")
	}

	return l.manager.TryGet(l.language, key, args)
}

func (l *Localizer) Raw(key string) string {
	if l.manager == nil {
		return key
	}

	return l.manager.Raw(l.language, key)
}

func (l *Localizer) Has(key string) bool {
	return l.manager != nil && l.manager.Has(l.language, key)
}

func (l *Localizer) With(language Language) *Localizer {
	return &Localizer{manager: l.manager, language: language}
}

func (l *Localizer) Manager() *Manager {
	return l.manager
}
