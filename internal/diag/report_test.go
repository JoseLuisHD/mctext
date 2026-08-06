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

package diag

import (
	"errors"
	"strings"
	"testing"
)

func TestLocationString(t *testing.T) {
	cases := []struct {
		location Location
		want     string
	}{
		{Location{}, "<unknown>"},
		{Location{File: "ui.ini"}, "ui.ini"},
		{Location{File: "ui.ini", Line: 12}, "ui.ini:12"},
		{Location{File: "ui.ini", Line: 12, Column: 9}, "ui.ini:12:9"},
		{Location{Line: 12, Column: 9}, "<unknown>"},
	}

	for _, tc := range cases {
		if got := tc.location.String(); got != tc.want {
			t.Errorf("Location%+v.String() = %q, want %q", tc.location, got, tc.want)
		}
	}
}

func TestReportUnderlinesTheOffendingText(t *testing.T) {
	const snippet = "<gold>Menu</red>"

	report := New(ErrLimit, "something went wrong").
		At(Location{File: "lang/en/ui.ini", Line: 12, Column: 9}).
		Field("key", "ui::menu::title").
		Quote(snippet, 11, len("</red>")).
		Hint("close the tag you opened").
		Error()

	lines := strings.Split(report, "\n")
	if lines[0] != "translation: something went wrong" {
		t.Errorf("first line = %q", lines[0])
	}

	// The caret must sit exactly under the eleventh rune of the quoted text.
	wantCaret := "  │ " + strings.Repeat(" ", 10) + strings.Repeat("^", 6)
	if !strings.Contains(report, "  │ "+snippet+"\n"+wantCaret) {
		t.Errorf("the underline is misplaced:\n%s", report)
	}

	// Every labelled field must line its value up at the same column.
	var column = -1
	for _, line := range lines {
		for _, label := range []string{"at:", "key:", "hint:"} {
			if strings.HasPrefix(line, "  "+label) {
				start := len(line) - len(strings.TrimLeft(line[len("  "+label):], " "))
				if column == -1 {
					column = start
				} else if start != column {
					t.Errorf("field %q is not aligned with the others:\n%s", label, report)
				}
			}
		}
	}

	if column == -1 {
		t.Fatalf("no labelled field was rendered:\n%s", report)
	}
}

func TestReportOmitsEmptyParts(t *testing.T) {
	report := New(ErrLimit, "bare").Field("key", "").Hint("").Error()
	if report != "translation: bare" {
		t.Errorf("empty fields and hints must be dropped, got %q", report)
	}
}

func TestDiagnosticWrapsItsSentinel(t *testing.T) {
	sentinel := errors.New("sentinel")
	err := error(New(sentinel, "boom"))

	if !errors.Is(err, sentinel) {
		t.Error("errors.Is must find the wrapped sentinel")
	}

	var d *Diagnostic
	if !errors.As(err, &d) || d.Summary != "boom" {
		t.Error("errors.As must recover the diagnostic")
	}
}

func TestReportExpandsTabs(t *testing.T) {
	report := New(ErrLimit, "x").Quote("a\tb", 0, 0).Error()
	if strings.Contains(report, "\t") {
		t.Errorf("tabs must be expanded so the underline stays aligned:\n%q", report)
	}
}
