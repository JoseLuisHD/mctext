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
	"sort"
	"strconv"

	"github.com/JoseLuisHD/mctext/internal/catalog"
	"github.com/JoseLuisHD/mctext/internal/diag"
	"github.com/JoseLuisHD/mctext/internal/markup"
)

// verifyConsistency enforces, in strict mode, that every language declares the
// same keys and that every key declares the same arguments everywhere.
//
// This turns the most common production failure of a translation system, a
// message that only exists in one language, into a start-up error.
func verifyConsistency(catalogs map[Language]*catalog.Catalog, languages []Language, reference Language) error {
	base, ok := catalogs[reference]
	if !ok {
		return newError(ErrInconsistentCatalog,
			"reference language "+string(reference)+" was not loaded")
	}

	var failures []error
	baseKeys := base.Keys()
	baseArgs := make(map[string][]string, len(baseKeys))
	base.Each(func(program *markup.Program) {
		baseArgs[program.Key()] = program.Args()
	})

	for _, language := range languages {
		if language == reference {
			continue
		}

		other := catalogs[language]

		for _, key := range baseKeys {
			program, found := other.Lookup(key)
			if !found {
				failures = append(failures, missingTranslation(language, reference, key, base))

				continue
			}

			if expected := baseArgs[key]; !sameArguments(expected, program.Args()) {
				failures = append(failures, diag.New(ErrInconsistentCatalog,
					"key "+strconv.Quote(key)+" declares different arguments across languages").
					At(program.Origin()).
					Field("language", string(language)).
					Field("declared", listOrNone(sorted(program.Args()))).
					Field("expected", listOrNone(sorted(expected))).
					Field("reference", string(reference)).
					Hint("every language must fill the same placeholders"))
			}

			if len(failures) >= maxReportedErrors {
				return errors.Join(failures...)
			}
		}

		for _, key := range other.Keys() {
			if _, found := base.Lookup(key); !found {
				program, _ := other.Lookup(key)
				failures = append(failures, diag.New(ErrInconsistentCatalog,
					"key "+strconv.Quote(key)+" is not declared in the reference language").
					At(program.Origin()).
					Field("language", string(language)).
					Field("reference", string(reference)).
					Hint("add it to "+string(reference)+" or remove it here"))
			}

			if len(failures) >= maxReportedErrors {
				return errors.Join(failures...)
			}
		}
	}

	if len(failures) > 0 {
		return errors.Join(failures...)
	}

	return nil
}

func missingTranslation(language, reference Language, key string, base *catalog.Catalog) error {
	d := diag.New(ErrInconsistentCatalog,
		"key "+strconv.Quote(key)+" is missing from language "+string(language)).
		Field("reference", string(reference))

	if program, ok := base.Lookup(key); ok {
		d.At(program.Origin()).Quote(program.Raw(), 0, 0)
	}

	return d.Hint("translate it, or disable the check with WithStrict(false)")
}

func sameArguments(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	left, right := sorted(a), sorted(b)
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}

	return true
}

func sorted(values []string) []string {
	out := append([]string(nil), values...)
	sort.Strings(out)

	return out
}
