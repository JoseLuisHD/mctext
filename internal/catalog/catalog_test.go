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

package catalog

import (
	"testing"

	"github.com/JoseLuisHD/mctext/internal/diag"
	"github.com/JoseLuisHD/mctext/internal/markup"
)

func program(t *testing.T, key, raw string) *markup.Program {
	t.Helper()
	compiled, err := markup.Compile(key, raw, diag.Location{File: "test.ini", Line: 1}, markup.Options{})
	if err != nil {
		t.Fatalf("Compile(%q): %v", raw, err)
	}

	return compiled
}

func TestPutAndLookup(t *testing.T) {
	c := New("en", 0)
	c.Put(program(t, "ui::menu::title", "<red>Title</red>"))

	found, ok := c.Lookup("ui::menu::title")
	if !ok {
		t.Fatal("Lookup did not find the message that was just stored")
	}

	if found.Raw() != "<red>Title</red>" {
		t.Errorf("Lookup returned %q", found.Raw())
	}

	if _, ok := c.Lookup("ui::menu::missing"); ok {
		t.Error("Lookup found a key that was never stored")
	}

	if c.Len() != 1 || c.Language() != "en" {
		t.Errorf("Len = %d, Language = %q", c.Len(), c.Language())
	}
}

func TestPutReportsTheReplacedMessage(t *testing.T) {
	c := New("en", 0)

	if previous, replaced := c.Put(program(t, "k", "first")); replaced || previous != nil {
		t.Error("the first Put must not report a replacement")
	}

	previous, replaced := c.Put(program(t, "k", "second"))
	if !replaced {
		t.Fatal("the second Put must report a replacement so the loader can refuse it")
	}

	if previous.Raw() != "first" {
		t.Errorf("the replaced message was %q, want the original", previous.Raw())
	}
}

func TestKeysAreSorted(t *testing.T) {
	c := New("en", 0)
	for _, key := range []string{"b::two", "a::one", "c::three"} {
		c.Put(program(t, key, "x"))
	}

	want := []string{"a::one", "b::two", "c::three"}
	got := c.Keys()
	if len(got) != len(want) {
		t.Fatalf("Keys() = %v", got)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Keys() = %v, want %v", got, want)
		}
	}
}

func TestSourcesKeepLoadOrder(t *testing.T) {
	c := New("es", 0)
	c.AddSource("es/common.ini")
	c.AddSource("es/shop.ini")

	sources := c.Sources()
	if len(sources) != 2 || sources[0] != "es/common.ini" || sources[1] != "es/shop.ini" {
		t.Errorf("Sources() = %v, want them in load order", sources)
	}
}

func TestEachVisitsEveryMessage(t *testing.T) {
	c := New("en", 0)
	c.Put(program(t, "a", "one"))
	c.Put(program(t, "b", "two"))

	visited := 0
	c.Each(func(*markup.Program) { visited++ })
	if visited != 2 {
		t.Errorf("Each visited %d messages, want 2", visited)
	}
}
