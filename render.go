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

// It is the fastest possible lookup, a map access returning an existing string,
// and it allocates nothing. Use it for logging, diffing or re-encoding, not for
// anything a player will see.
func (m *Manager) Raw(language Language, key string) string {
	s := m.snap.Load()
	if s == nil {
		m.report(m.notInitialised(key, language))
		return key
	}

	program, _, ok := s.lookup(language, key)
	if !ok {
		m.report(m.missingKey(language, key, s))
		return key
	}

	return program.Raw()
}

func (m *Manager) TryRaw(language Language, key string) (string, error) {
	s := m.snap.Load()
	if s == nil {
		return key, m.notInitialised(key, language)
	}

	program, _, ok := s.lookup(language, key)
	if !ok {
		return key, m.missingKey(language, key, s)
	}

	return program.Raw(), nil
}

// It never fails loudly. A missing key returns the key itself and a missing
// argument leaves its placeholder in place, so a translation mistake degrades
// the text instead of breaking the caller. Every such failure is handed to the
// configured error handler with a full diagnostic, and TryGet exposes it
// directly.
func (m *Manager) Get(language Language, key string, args Args) string {
	text, err := m.render(language, key, args, nil)
	if err != nil {
		m.report(err)
	}

	return text
}

func (m *Manager) TryGet(language Language, key string, args Args) (string, error) {
	return m.render(language, key, args, nil)
}

func (m *Manager) GetPairs(language Language, key string, pairs ...Pair) string {
	text, err := m.render(language, key, nil, pairs)
	if err != nil {
		m.report(err)
	}

	return text
}

func (m *Manager) render(language Language, key string, args Args, pairs pairsView) (string, error) {
	s := m.snap.Load()
	if s == nil {
		return key, m.notInitialised(key, language)
	}

	program, resolved, ok := s.lookup(language, key)
	if !ok {
		return key, m.missingKey(language, key, s)
	}

	// Fast path: a message without arguments was fully rendered at load time,
	// so returning it costs nothing at all.
	if static, isStatic := program.Static(); isStatic {
		if err := m.checkUnexpected(program, resolved, key, args, pairs); err != nil {
			return static, err
		}

		return static, nil
	}

	names := program.Args()
	var scratch [8]string
	var values []string
	if len(names) <= len(scratch) {
		values = scratch[:len(names)]
	} else {
		values = make([]string, len(names))
	}

	var missing []string
	for i, name := range names {
		value, found := lookupArgument(args, pairs, name)
		if !found {
			missing = append(missing, name)
			values[i] = "{" + name + "}"

			continue
		}

		text := formatValue(value)
		if m.cfg.sanitiseValues {
			text = sanitiseValue(text)
		}

		values[i] = text
	}

	buffer := m.pool.Get().(*[]byte)
	*buffer = (*buffer)[:0]
	if cap(*buffer) < program.SizeHint() {
		*buffer = make([]byte, 0, program.SizeHint()+32)
	}

	*buffer = program.AppendTo(*buffer, values)
	text := string(*buffer)
	m.pool.Put(buffer)

	if len(missing) > 0 {
		return text, m.missingArguments(program, resolved, missing, args, pairs)
	}

	if err := m.checkUnexpected(program, resolved, key, args, pairs); err != nil {
		return text, err
	}

	return text, nil
}

func lookupArgument(args Args, pairs pairsView, name string) (any, bool) {
	if args != nil {
		value, ok := args[name]
		if ok {
			return value, true
		}
	}

	if pairs != nil {
		return pairs.lookup(name)
	}

	return nil, false
}
