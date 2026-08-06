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
	"sync"
	"testing"
)

func TestManagerLoadsFromDisk(t *testing.T) {
	manager, err := New(
		WithDirectory("testdata/lang"),
		Open(English, "common.ini", "shop.ini"),
		Open(Spanish, "common.ini", "shop.ini"),
		WithDefault(English),
		WithStrict(true),
	)

	if err != nil {
		t.Fatalf("New returned an unexpected error:\n%v", err)
	}

	if got, want := manager.Size(), 12; got != want {
		t.Errorf("Size() = %d, want %d", got, want)
	}

	got := manager.Get(Spanish, "shop::purchase::success", Args{"item": "Espada", "price": 250})
	const want = "§aCompraste §eEspada§a por §6250§a monedas.§r"
	if got != want {
		t.Errorf("Get\n got: %q\nwant: %q", got, want)
	}

	if raw := manager.Raw(English, "ui::button::close"); raw != "<red>Close</red>" {
		t.Errorf("Raw() = %q, want the verbatim source", raw)
	}
}

func TestSectionsPrefixKeys(t *testing.T) {
	manager := memoryManager(t, map[string]string{
		"en/common.ini": "[a::b]\nc: hello\n",
		"es/common.ini": "[a::b]\nc: hola\n",
	})

	if !manager.Has(English, "a::b::c") {
		t.Fatalf("expected the section to prefix the key, got keys %v", manager.Keys(English))
	}
}

func TestFallbackResolvesMissingKeys(t *testing.T) {
	manager := memoryManager(t, map[string]string{
		"en/common.ini": "only::in::english: <red>Fallback</red>\nshared: en\n",
		"es/common.ini": "shared: es\n",
	}, WithFallback(English))

	if got := manager.Get(Spanish, "only::in::english", nil); got != "§cFallback§r" {
		t.Errorf("expected the English fallback, got %q", got)
	}

	if got := manager.Get(Spanish, "shared", nil); got != "es" {
		t.Errorf("the requested language must win over the fallback, got %q", got)
	}
}

func TestReloadIsAtomicAndKeepsServingOnFailure(t *testing.T) {
	files := map[string]string{
		"en/common.ini": "greet: Hello\n",
		"es/common.ini": "greet: Hola\n",
	}

	manager := memoryManager(t, files)

	broken := memoryManager(t, files)
	broken.cfg.source = MemorySource(map[string]string{"en/common.ini": "greet: <nope>x</nope>\n"})
	if err := broken.Reload(); err == nil {
		t.Fatal("expected the reload to fail")
	}

	if got := broken.Get(English, "greet", nil); got != "Hello" {
		t.Errorf("a failed reload must keep the previous snapshot, got %q", got)
	}

	_ = manager
}

func TestConcurrentReadsAreSafe(t *testing.T) {
	manager := memoryManager(t, map[string]string{
		"en/common.ini": "greet: <green>Hello {player}</green>\n",
		"es/common.ini": "greet: <green>Hola {player}</green>\n",
	})

	var wg sync.WaitGroup
	for worker := 0; worker < 16; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 500; i++ {
				if got := manager.Get(Spanish, "greet", Args{"player": "Neo"}); got != "§aHola Neo§r" {
					t.Errorf("Get = %q", got)

					return
				}
				_ = manager.Raw(English, "greet")
			}
		}()
	}

	wg.Wait()
}

func TestInitIsIdempotent(t *testing.T) {
	err := Init(
		WithSource(MemorySource(map[string]string{
			"en/common.ini": "greet: Hello\n",
			"es/common.ini": "greet: Hola\n",
		})),
		Open(English, "common.ini"),
		Open(Spanish, "common.ini"),
		WithStrict(true),
	)

	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	if !Initialised() {
		t.Fatal("Initialised() = false after a successful Init")
	}

	if got := Of(Spanish).Get("greet", nil); got != "Hola" {
		t.Errorf("Of(Spanish).Get = %q", got)
	}

	if err := Init(); !errors.Is(err, ErrAlreadyInitialised) {
		t.Errorf("a second Init must be refused, got %v", err)
	}
}
