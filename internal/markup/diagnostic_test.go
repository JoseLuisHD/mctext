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

package markup

import (
	"errors"
	"strings"
	"testing"

	"github.com/JoseLuisHD/mctext/internal/diag"
)

func TestCompileRejectsMalformedMarkup(t *testing.T) {
	cases := []struct {
		name     string
		raw      string
		sentinel error
		contains string
	}{
		{"unknown tag", "<redd>x</redd>", ErrMarkup, "did you mean <red>?"},
		{"unterminated tag", "<red>x", ErrMarkup, "unclosed tag"},
		{"mismatched close", "<gold>a<bold>b</gold>", ErrMarkup, "mismatched closing tag"},
		{"stray close", "a</red>", ErrMarkup, "without a matching opening tag"},
		{"closing a reset", "<reset>a</reset>", ErrMarkup, "must not be closed"},
		{"two colours", "<red:blue>x</>", ErrMarkup, "more than one colour"},
		{"empty tag", "<>x", ErrMarkup, "empty tag"},
		{"raw code", "§4Danger", ErrMarkup, "write <dark_red>"},
		{"unterminated placeholder", "hello {player", ErrPlaceholder, "unterminated"},
		{"illegal placeholder", "hello {pla yer}", ErrPlaceholder, "illegal character"},
		{"unknown escape", `a \q b`, ErrMarkup, "unknown escape sequence"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Compile("test::key", tc.raw, diag.Location{File: "test.ini", Line: 7, Column: 5}, Options{})
			if err == nil {
				t.Fatalf("Compile(%q) succeeded, want an error", tc.raw)
			}

			if !errors.Is(err, tc.sentinel) {
				t.Errorf("errors.Is(err, %v) = false for:\n%v", tc.sentinel, err)
			}

			if !strings.Contains(err.Error(), tc.contains) {
				t.Errorf("diagnostic does not mention %q:\n%v", tc.contains, err)
			}

			if !strings.Contains(err.Error(), "test.ini:7") {
				t.Errorf("diagnostic does not locate the failure:\n%v", err)
			}

			if !strings.Contains(err.Error(), "test::key") {
				t.Errorf("diagnostic does not name the key:\n%v", err)
			}
		})
	}
}

func TestCompileHonoursLimits(t *testing.T) {
	deep := strings.Repeat("<bold>", 5) + "x" + strings.Repeat("</bold>", 5)
	if _, err := Compile("k", deep, diag.Location{}, Options{MaxNesting: 3}); !errors.Is(err, ErrLimit) {
		t.Fatalf("expected a limit error for nesting depth 5 with MaxNesting=3, got %v", err)
	}

	if _, err := Compile("k", "{a}{b}{c}", diag.Location{}, Options{MaxArguments: 2}); !errors.Is(err, ErrLimit) {
		t.Fatalf("expected a limit error for 3 arguments with MaxArguments=2, got %v", err)
	}
}

func TestDiagnosticPointsAtTheOffendingColumn(t *testing.T) {
	_, err := Compile("ui::title", "<gold>Menu <bold>Main</gold>", diag.Location{File: "ui.ini", Line: 3, Column: 8}, Options{})
	if err == nil {
		t.Fatal("expected an error")
	}

	report := err.Error()
	// The closing tag starts at column 22 of the value, which itself begins at
	// column 8 of the source line.
	if !strings.Contains(report, "ui.ini:3:29") {
		t.Errorf("diagnostic should point at ui.ini:3:29:\n%s", report)
	}

	if !strings.Contains(report, "^") {
		t.Errorf("diagnostic should underline the offending text:\n%s", report)
	}
}

func FuzzCompile(f *testing.F) {
	seeds := []string{
		"<red>x</red>", "<bold:red>{a}</bold:red>", "<>", "</>", "{", "}", `\`,
		"§4", "<gold>a<bold>b</bold>c</gold>", "<reset:white>", "{a}{a}",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	// The compiler must either fail with a diagnostic or produce a program that
	// renders without panicking, for any input whatsoever.
	f.Fuzz(func(t *testing.T, raw string) {
		program, err := Compile("fuzz", raw, diag.Location{File: "fuzz.ini"}, Options{})
		if err != nil {
			var d *diag.Diagnostic
			if !errors.As(err, &d) {
				t.Fatalf("error is not a diagnostic: %v", err)
			}

			return
		}

		values := make([]string, len(program.Args()))
		_ = program.AppendTo(nil, values)
	})
}
