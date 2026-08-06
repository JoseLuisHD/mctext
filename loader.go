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
	"errors"
	"strconv"

	"github.com/JoseLuisHD/mctext/internal/catalog"
	"github.com/JoseLuisHD/mctext/internal/diag"
	"github.com/JoseLuisHD/mctext/internal/ini"
	"github.com/JoseLuisHD/mctext/internal/markup"
)

// maxReportedErrors bounds how many load failures are reported at once. Showing
// every broken message in a single run is useful; showing ten thousand is not.
const maxReportedErrors = 25

// load builds a complete snapshot. It touches no shared state, so a failure
// leaves the Manager serving whatever it was serving before.
func (m *Manager) load() (*snapshot, error) {
	var failures []error
	catalogs := make(map[Language]*catalog.Catalog, len(m.cfg.order))
	languages := make([]Language, 0, len(m.cfg.order))
	total := 0

	for _, entry := range m.cfg.order {
		built, errs := m.loadLanguage(entry)
		failures = append(failures, errs...)

		if len(failures) >= maxReportedErrors {
			return nil, errors.Join(failures...)
		}

		catalogs[entry.language] = built
		languages = append(languages, entry.language)
		total += built.Len()
	}

	if len(failures) > 0 {
		return nil, errors.Join(failures...)
	}

	if m.cfg.strict {
		if err := verifyConsistency(catalogs, languages, m.cfg.defaultLanguage); err != nil {
			return nil, err
		}
	}

	return &snapshot{
		catalogs:        catalogs,
		resolution:      buildResolution(languages, m.cfg.fallback),
		defaultLanguage: m.cfg.defaultLanguage,
		languages:       languages,
		messages:        total,
	}, nil
}

func (m *Manager) loadLanguage(entry languageFiles) (*catalog.Catalog, []error) {
	built := catalog.New(string(entry.language), 256)
	var failures []error

	for _, declared := range entry.files {
		path := m.cfg.resolve(entry.language, declared)
		data, err := m.cfg.source.ReadFile(path)
		if err != nil {
			failures = append(failures, diag.New(ErrSource, "cannot read "+strconv.Quote(path)).
				At(diag.Location{File: path}).
				Field("language", string(entry.language)).
				Field("source", m.cfg.source.Describe()).
				Field("cause", err.Error()))

			continue
		}

		file, err := ini.Parse(path, data, m.cfg.iniLimits())
		if err != nil {
			failures = append(failures, err)
			continue
		}

		built.AddSource(path)

		for i := range file.Entries {
			declaration := &file.Entries[i]
			program, err := markup.Compile(
				declaration.Key,
				declaration.Value,
				declaration.Location,
				m.cfg.markupOptions(),
			)

			if err != nil {
				failures = append(failures, err)
				if len(failures) >= maxReportedErrors {
					return built, failures
				}

				continue
			}

			if previous, replaced := built.Put(program); replaced {
				failures = append(failures, diag.New(ErrDuplicateKey,
					"key "+strconv.Quote(program.Key())+" is declared twice").
					At(program.Origin()).
					Field("language", string(entry.language)).
					Field("first", previous.Origin().String()).
					Hint("a key may only be declared once per language, across all its files"))
			}
		}
	}

	return built, failures
}

// buildResolution precomputes, for every language, the ordered list of catalogs
// a lookup walks. Doing it once at load time keeps the read path branch free.
func buildResolution(languages []Language, fallback Language) map[Language][]Language {
	resolution := make(map[Language][]Language, len(languages))
	for _, language := range languages {
		chain := []Language{language}
		if fallback != "" && fallback != language {
			chain = append(chain, fallback)
		}

		resolution[language] = chain
	}

	return resolution
}
