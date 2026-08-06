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
	"errors"
	"strings"
	"testing"
)

func parse(t *testing.T, content string) *File {
	t.Helper()
	file, err := Parse("test.ini", []byte(content), Limits{})
	if err != nil {
		t.Fatalf("Parse returned an unexpected error:\n%v", err)
	}

	return file
}

func TestParseKeysAndValues(t *testing.T) {
	file := parse(t, strings.Join([]string{
		"# a comment",
		"; another comment",
		"",
		"id::some::message1: Hello world",
		"with_equals = value",
		`quoted: "  padded  "`,
		"colon::in::value: ratio 16:9",
		"[shop::purchase]",
		"success: bought",
		"[]",
		"top::level: back to no prefix",
	}, "\n"))

	want := []struct{ key, value string }{
		{"id::some::message1", "Hello world"},
		{"with_equals", "value"},
		{"quoted", "  padded  "},
		{"colon::in::value", "ratio 16:9"},
		{"shop::purchase::success", "bought"},
		{"top::level", "back to no prefix"},
	}

	if len(file.Entries) != len(want) {
		t.Fatalf("got %d entries, want %d: %+v", len(file.Entries), len(want), file.Entries)
	}

	for i, expected := range want {
		got := file.Entries[i]
		if got.Key != expected.key || got.Value != expected.value {
			t.Errorf("entry %d = {%q: %q}, want {%q: %q}", i, got.Key, got.Value, expected.key, expected.value)
		}
	}
}

func TestParseValueColumnPointsAtTheValue(t *testing.T) {
	file := parse(t, "key:    value\n")
	if got := file.Entries[0].Location.Column; got != 9 {
		t.Errorf("value column = %d, want 9", got)
	}
}

func TestParseContinuation(t *testing.T) {
	file := parse(t, "long::key: first part \\\n           second part\nnext: x\n")
	if got := file.Entries[0].Value; got != "first part second part" {
		t.Errorf("continuation = %q", got)
	}

	if got := file.Entries[0].Location.Line; got != 1 {
		t.Errorf("a continued entry must report its first line, got %d", got)
	}

	if got := file.Entries[1].Location.Line; got != 3 {
		t.Errorf("the following entry must report line 3, got %d", got)
	}
}

func TestParseKeepsTrailingBackslashWhenEscaped(t *testing.T) {
	file := parse(t, `key: path\\`+"\n")
	if got := file.Entries[0].Value; got != `path\\` {
		t.Errorf("an escaped backslash must not continue the line, got %q", got)
	}
}

func TestParseRejectsMalformedInput(t *testing.T) {
	cases := []struct {
		name     string
		content  string
		sentinel error
		contains string
	}{
		{"missing separator", "just a line\n", ErrSyntax, "missing ':'"},
		{"empty key", ": value\n", ErrSyntax, "empty key"},
		{"dangling separator", "a::: value\n", ErrSyntax, "missing ':'"},
		{"illegal key", "bad key!: value\n", ErrSyntax, "illegal character"},
		{"trailing separator", "a::: b\n", ErrSyntax, "missing ':'"},
		{"unterminated section", "[shop\nk: v\n", ErrSyntax, "unterminated section"},
		{"dangling continuation", "k: v \\\n", ErrSyntax, "continuation backslash"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse("test.ini", []byte(tc.content), Limits{})
			if err == nil {
				t.Fatalf("Parse(%q) succeeded, want an error", tc.content)
			}

			if !errors.Is(err, tc.sentinel) {
				t.Errorf("errors.Is(err, %v) = false for:\n%v", tc.sentinel, err)
			}

			if !strings.Contains(err.Error(), tc.contains) {
				t.Errorf("diagnostic does not mention %q:\n%v", tc.contains, err)
			}
		})
	}
}

func TestParseRejectsInvalidUTF8(t *testing.T) {
	_, err := Parse("test.ini", []byte{'k', ':', ' ', 0xff, 0xfe}, Limits{})
	if !errors.Is(err, ErrSyntax) {
		t.Fatalf("expected ErrSyntax, got %v", err)
	}
}

func TestParseEnforcesLimits(t *testing.T) {
	content := strings.Repeat("k: v\n", 10)
	if _, err := Parse("test.ini", []byte(content), Limits{MaxEntries: 3}); !errors.Is(err, ErrLimit) {
		t.Errorf("expected an entry limit error, got %v", err)
	}

	if _, err := Parse("test.ini", []byte(content), Limits{MaxFileSize: 4}); !errors.Is(err, ErrLimit) {
		t.Errorf("expected a file size limit error, got %v", err)
	}

	if _, err := Parse("test.ini", []byte("k: v\n"), Limits{MaxKeyLength: 0}); err != nil {
		t.Errorf("a zero limit must fall back to the default, got %v", err)
	}
}

func TestParseStripsBOMAndCRLF(t *testing.T) {
	file := parse(t, "\ufeffkey: value\r\nother: x\r\n")
	if file.Entries[0].Key != "key" || file.Entries[0].Value != "value" {
		t.Errorf("BOM or CRLF leaked into the entry: %+v", file.Entries[0])
	}
}
