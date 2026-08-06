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
	"sync"
	"sync/atomic"
)

// A Manager is safe for concurrent use by any number of goroutines. All of its
// mutable state lives behind one atomic pointer, replaced wholesale by Reload.
type Manager struct {
	mu   sync.Mutex // serialises loads, never taken on the read path
	cfg  config
	snap atomic.Pointer[snapshot]
	pool sync.Pool
}

func New(options ...Option) (*Manager, error) {
	cfg := defaultConfig()
	for _, option := range options {
		if option == nil {
			continue
		}

		if err := option(&cfg); err != nil {
			return nil, err
		}
	}

	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	m := &Manager{cfg: cfg}
	m.pool.New = func() any {
		buffer := make([]byte, 0, 256)
		return &buffer
	}

	if err := m.Reload(); err != nil {
		return nil, err
	}

	return m, nil
}

func validateConfig(cfg *config) error {
	if cfg.source == nil {
		return newError(ErrNoSource, "no source configured").
			Hint("pass WithDirectory(dir), WithSource(FSSource(embedded, \"lang\")) or WithSource(MemorySource(...))")
	}

	if len(cfg.order) == 0 {
		return newError(ErrNoSource, "no language declared").
			Hint("declare files with Open(mctext.English, \"common.ini\")")
	}

	declared := make(map[Language]bool, len(cfg.order))
	for _, entry := range cfg.order {
		declared[entry.language] = true
	}

	if !declared[cfg.defaultLanguage] {
		return newError(ErrInvalidLanguage,
			"default language "+string(cfg.defaultLanguage)+" has no declared files")
	}

	if cfg.fallback != "" && !declared[cfg.fallback] {
		return newError(ErrInvalidLanguage,
			"fallback language "+string(cfg.fallback)+" has no declared files")
	}

	return nil
}

// Reload rebuilds every catalog from the configured source and atomically
// replaces the served snapshot. Readers keep using the previous snapshot until
// the new one is fully built, so a failed reload leaves the Manager untouched.
func (m *Manager) Reload() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	next, err := m.load()
	if err != nil {
		return err
	}

	m.snap.Store(next)

	return nil
}

func (m *Manager) Languages() []Language {
	s := m.snap.Load()
	if s == nil {
		return nil
	}

	return append([]Language(nil), s.languages...)
}

func (m *Manager) DefaultLanguage() Language {
	if s := m.snap.Load(); s != nil {
		return s.defaultLanguage
	}

	return m.cfg.defaultLanguage
}

func (m *Manager) Size() int {
	if s := m.snap.Load(); s != nil {
		return s.messages
	}

	return 0
}

func (m *Manager) Keys(language Language) []string {
	s := m.snap.Load()
	if s == nil {
		return nil
	}

	if c, ok := s.catalogs[language]; ok {
		return c.Keys()
	}

	return nil
}

func (m *Manager) Sources(language Language) []string {
	s := m.snap.Load()
	if s == nil {
		return nil
	}

	if c, ok := s.catalogs[language]; ok {
		return append([]string(nil), c.Sources()...)
	}

	return nil
}

func (m *Manager) Has(language Language, key string) bool {
	s := m.snap.Load()
	if s == nil {
		return false
	}

	_, _, ok := s.lookup(language, key)

	return ok
}

func (m *Manager) Localizer(language Language) *Localizer {
	return &Localizer{manager: m, language: language}
}
