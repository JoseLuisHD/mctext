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
	"strconv"
	"strings"

	"github.com/JoseLuisHD/mctext/internal/diag"
	"github.com/JoseLuisHD/mctext/internal/markup"
)

// checkUnexpected reports values that no placeholder consumes. It only runs in
// strict mode, because a spurious argument is harmless at runtime but almost
// always signals a renamed placeholder.
func (m *Manager) checkUnexpected(program *markup.Program, language Language, key string, args Args, pairs pairsView) error {
	if !m.cfg.strict || (len(args) == 0 && len(pairs) == 0) {
		return nil
	}

	var unexpected []string
	for name := range args {
		if program.ArgIndex(name) < 0 {
			unexpected = append(unexpected, name)
		}
	}

	for _, name := range pairs.names() {
		if program.ArgIndex(name) < 0 {
			unexpected = append(unexpected, name)
		}
	}

	if len(unexpected) == 0 {
		return nil
	}

	return diag.New(ErrUnexpectedArgument, "message does not declare "+quoteList(unexpected)).
		At(program.Origin()).
		Field("key", key).
		Field("language", string(language)).
		Field("declared", listOrNone(program.Args())).
		Quote(program.Raw(), 0, 0).
		Hint("remove the value or add the placeholder to the message").
		Stack(m.cfg.stackDepth, stackIgnorePrefix)
}

func (m *Manager) missingArguments(program *markup.Program, language Language, missing []string, args Args, pairs pairsView) error {
	provided := make([]string, 0, len(args)+len(pairs))
	for name := range args {
		provided = append(provided, name)
	}

	provided = append(provided, pairs.names()...)

	return diag.New(ErrMissingArgument, "no value supplied for "+quoteList(missing)).
		At(program.Origin()).
		Field("key", program.Key()).
		Field("language", string(language)).
		Field("declared", listOrNone(program.Args())).
		Field("supplied", listOrNone(provided)).
		Quote(program.Raw(), 0, 0).
		Hint("pass it as mctext.Args{\""+missing[0]+"\": value}").
		Stack(m.cfg.stackDepth, stackIgnorePrefix)
}

func (m *Manager) missingKey(language Language, key string, s *snapshot) error {
	d := diag.New(ErrMissingKey, "no message declared for key "+strconv.Quote(key)).
		Field("language", string(language)).
		Stack(m.cfg.stackDepth, stackIgnorePrefix)

	if _, ok := s.catalogs[language]; !ok {
		d.Field("loaded", languageList(s.languages)).
			Hint("declare the language with Open(" + strconv.Quote(string(language)) + ", …)")
	} else if alt := s.nearestKey(language, key); alt != "" {
		d.Hint("did you mean " + strconv.Quote(alt) + "?")
	}

	return d
}

func (m *Manager) notInitialised(key string, language Language) error {
	return diag.New(ErrNotInitialised, "lookup before the catalogs were loaded").
		Field("key", key).
		Field("language", string(language)).
		Hint("call mctext.Init(…) once during start-up").
		Stack(m.cfg.stackDepth, stackIgnorePrefix)
}

func (m *Manager) report(err error) {
	if err != nil && m.cfg.onError != nil {
		m.cfg.onError(err)
	}
}

func quoteList(items []string) string {
	var b strings.Builder
	for i, item := range items {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(strconv.Quote(item))
	}

	return b.String()
}

func listOrNone(items []string) string {
	if len(items) == 0 {
		return "none"
	}

	return "[" + strings.Join(items, " ") + "]"
}

func languageList(languages []Language) string {
	names := make([]string, len(languages))
	for i, language := range languages {
		names[i] = string(language)
	}

	return "[" + strings.Join(names, " ") + "]"
}
