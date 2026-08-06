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
	"strings"
	"testing"
)

func TestMissingKeyReturnsTheKeyAndReports(t *testing.T) {
	var reported []error
	manager := memoryManager(t, map[string]string{
		"en/common.ini": "ui::menu::title: Title\n",
		"es/common.ini": "ui::menu::title: Titulo\n",
	}, WithFallback(""), WithErrorHandler(func(err error) { reported = append(reported, err) }))

	if got := manager.Get(English, "ui::menu::titel", nil); got != "ui::menu::titel" {
		t.Errorf("a missing key must degrade to the key itself, got %q", got)
	}

	if len(reported) != 1 {
		t.Fatalf("expected exactly one reported failure, got %d", len(reported))
	}

	if !errors.Is(reported[0], ErrMissingKey) {
		t.Errorf("expected ErrMissingKey, got %v", reported[0])
	}

	if !strings.Contains(reported[0].Error(), "did you mean") {
		t.Errorf("expected a suggestion in:\n%v", reported[0])
	}
}

func TestStrictModeRejectsUnexpectedArguments(t *testing.T) {
	manager := memoryManager(t, map[string]string{
		"en/common.ini": "greet: Hello {player}\n",
		"es/common.ini": "greet: Hola {player}\n",
	}, WithStrict(true))

	_, err := manager.TryGet(English, "greet", Args{"player": "Neo", "plyer": "typo"})
	if !errors.Is(err, ErrUnexpectedArgument) {
		t.Fatalf("expected ErrUnexpectedArgument, got %v", err)
	}
}

func TestStrictModeRejectsInconsistentCatalogs(t *testing.T) {
	_, err := New(
		WithSource(MemorySource(map[string]string{
			"en/common.ini": "a: one\nb: two {x}\n",
			"es/common.ini": "b: dos {y}\n",
		})),
		Open(English, "common.ini"),
		Open(Spanish, "common.ini"),
		WithStrict(true),
	)

	if err == nil {
		t.Fatal("expected the load to fail")
	}

	if !errors.Is(err, ErrInconsistentCatalog) {
		t.Fatalf("expected ErrInconsistentCatalog, got %v", err)
	}

	report := err.Error()
	if !strings.Contains(report, `"a" is missing from language es`) {
		t.Errorf("expected the missing key to be reported:\n%s", report)
	}

	if !strings.Contains(report, "different arguments") {
		t.Errorf("expected the argument mismatch to be reported:\n%s", report)
	}
}

func TestDuplicateKeysAreRejected(t *testing.T) {
	_, err := New(
		WithSource(MemorySource(map[string]string{
			"en/a.ini": "dup: one\n",
			"en/b.ini": "dup: two\n",
		})),
		Open(English, "a.ini", "b.ini"),
		WithDefault(English),
		WithFallback(English),
	)

	if !errors.Is(err, ErrDuplicateKey) {
		t.Fatalf("expected ErrDuplicateKey, got %v", err)
	}
}

func TestSourceRejectsTraversal(t *testing.T) {
	_, err := New(
		WithSource(MemorySource(map[string]string{"en/common.ini": "a: b\n"})),
		WithLayout(FlatLayout),
		Open(English, "../../etc/passwd"),
		WithDefault(English),
		WithFallback(English),
	)

	if !errors.Is(err, ErrSource) {
		t.Fatalf("expected the traversal to be rejected, got %v", err)
	}
}

func TestRegisterLanguage(t *testing.T) {
	portuguese, err := RegisterLanguage("pt-br", "Português do Brasil")
	if err != nil {
		t.Fatalf("RegisterLanguage: %v", err)
	}

	if !portuguese.Registered() || portuguese.DisplayName() != "Português do Brasil" {
		t.Errorf("the language was not registered correctly")
	}

	if _, err := RegisterLanguage("english!", ""); !errors.Is(err, ErrInvalidLanguage) {
		t.Errorf("expected a malformed tag to be rejected, got %v", err)
	}

	manager, err := New(
		WithSource(MemorySource(map[string]string{
			"en/common.ini":    "greet: Hello\n",
			"pt-br/common.ini": "greet: Olá\n",
		})),
		Open(English, "common.ini"),
		Open(portuguese, "common.ini"),
		WithStrict(true),
	)

	if err != nil {
		t.Fatalf("adding a language must not require any change to the package:\n%v", err)
	}

	if got := manager.Get(portuguese, "greet", nil); got != "Olá" {
		t.Errorf("Get = %q", got)
	}
}
