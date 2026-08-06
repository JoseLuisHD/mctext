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

package ini

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/JoseLuisHD/mctext/internal/diag"
)

func Parse(name string, data []byte, limits Limits) (*File, error) {
	limits = limits.normalized()

	if len(data) > limits.MaxFileSize {
		return nil, diag.New(ErrLimit, "file larger than "+strconv.Itoa(limits.MaxFileSize)+" bytes").
			At(diag.Location{File: name})
	}

	if !utf8.Valid(data) {
		return nil, diag.New(ErrSyntax, "file is not valid UTF-8").
			At(diag.Location{File: name}).
			Hint("save translation files as UTF-8 without BOM")
	}

	text := strings.TrimPrefix(string(data), "\ufeff")
	file := &File{Name: name, Entries: make([]Entry, 0, 64)}

	var prefix string
	lineNumber := 0
	for offset := 0; offset < len(text); {
		lineNumber++
		line, next := readLine(text, offset)
		offset = next

		trimmed := strings.TrimSpace(line)
		if trimmed == "" || trimmed[0] == '#' || trimmed[0] == ';' {
			continue
		}

		if len(line) > limits.MaxLineWidth {
			return nil, at(name, lineNumber, 0).
				Wrap(ErrLimit, "line longer than "+strconv.Itoa(limits.MaxLineWidth)+" bytes")
		}

		if trimmed[0] == '[' {
			section, err := parseSection(name, lineNumber, trimmed)
			if err != nil {
				return nil, err
			}

			prefix = section

			continue
		}

		// Join continuation lines before splitting so that a wrapped value is
		// indistinguishable from a single long line.
		startLine := lineNumber
		value := line
		for continues(value) {
			value = value[:len(value)-1]
			if offset >= len(text) {
				return nil, at(name, lineNumber, 0).
					Wrap(ErrSyntax, "continuation backslash at end of file").
					Hint(`write \\ for a literal trailing backslash`)
			}

			lineNumber++

			var cont string
			cont, offset = readLine(text, offset)
			value += strings.TrimLeft(cont, " \t")
		}

		entry, err := parseEntry(name, startLine, value, prefix, limits)
		if err != nil {
			return nil, err
		}

		if len(file.Entries) >= limits.MaxEntries {
			return nil, at(name, lineNumber, 0).
				Wrap(ErrLimit, "more than "+strconv.Itoa(limits.MaxEntries)+" entries")
		}

		file.Entries = append(file.Entries, entry)
	}

	return file, nil
}

func readLine(text string, offset int) (line string, next int) {
	end := strings.IndexByte(text[offset:], '\n')
	if end < 0 {
		return strings.TrimSuffix(text[offset:], "\r"), len(text)
	}

	end += offset

	return strings.TrimSuffix(text[offset:end], "\r"), end + 1
}

// continues reports whether a line ends with an unescaped continuation marker,
// that is, an odd number of trailing backslashes.
func continues(line string) bool {
	count := 0
	for i := len(line) - 1; i >= 0 && line[i] == '\\'; i-- {
		count++
	}

	return count%2 == 1
}

func parseSection(file string, line int, trimmed string) (string, error) {
	if !strings.HasSuffix(trimmed, "]") {
		return "", at(file, line, 1).
			Quote(trimmed, 1, len(trimmed)).
			Wrap(ErrSyntax, "unterminated section header").
			Hint("section headers look like [shop::purchase]")
	}

	section := strings.TrimSpace(trimmed[1 : len(trimmed)-1])
	if section == "" {
		return "", nil
	}

	if err := validateKey(section); err != nil {
		return "", at(file, line, 1).
			Quote(trimmed, 1, len(trimmed)).
			Wrap(ErrSyntax, err.Error())
	}

	return section + Separator, nil
}

func parseEntry(file string, line int, raw, prefix string, limits Limits) (Entry, error) {
	sep, sepLen := findSeparator(raw)
	if sep < 0 {
		column := len(strings.TrimRight(raw, " \t")) + 1
		return Entry{}, at(file, line, column).
			Quote(raw, column, 1).
			Wrap(ErrSyntax, "missing ':' between key and value").
			Hint("declare messages as some::key: message text")
	}

	key := prefix + strings.TrimSpace(raw[:sep])
	if err := validateKey(key); err != nil {
		return Entry{}, at(file, line, 1).
			Quote(raw, 1, sep).
			Wrap(ErrSyntax, err.Error())
	}

	if len(key) > limits.MaxKeyLength {
		return Entry{}, at(file, line, 1).
			Wrap(ErrLimit, "key longer than "+strconv.Itoa(limits.MaxKeyLength)+" bytes")
	}

	rest := raw[sep+sepLen:]
	trimmedLeft := strings.TrimLeft(rest, " \t")
	column := sep + sepLen + (len(rest) - len(trimmedLeft)) + 1
	value := strings.TrimRight(trimmedLeft, " \t")

	// Quoting is the escape hatch for values whose own whitespace matters.
	if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
		value = value[1 : len(value)-1]
		column++
	}

	return Entry{
		Key:      key,
		Value:    value,
		Location: diag.Location{File: file, Line: line, Column: column},
	}, nil
}
