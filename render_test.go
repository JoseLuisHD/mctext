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
	"strings"
	"testing"
)

func TestArgumentValuesAreSanitised(t *testing.T) {
	files := map[string]string{
		"en/common.ini": "chat: <gray>{player}</gray>: {message}\n",
		"es/common.ini": "chat: <gray>{player}</gray>: {message}\n",
	}

	manager := memoryManager(t, files)
	got := manager.Get(English, "chat", Args{"player": "§4Admin§r", "message": "hi §lthere"})
	if strings.Contains(got, "§4") || strings.Contains(got, "§l") {
		t.Errorf("player supplied codes must be stripped, got %q", got)
	}

	if got != "§7Admin§r: hi there" {
		t.Errorf("got %q", got)
	}

	unsafe := memoryManager(t, files, WithArgumentSanitisation(false))
	if got := unsafe.Get(English, "chat", Args{"player": "§4Admin", "message": "x"}); !strings.Contains(got, "§4") {
		t.Errorf("sanitisation must be disablable, got %q", got)
	}
}

func TestValueFormatting(t *testing.T) {
	manager := memoryManager(t, map[string]string{
		"en/common.ini": "v: {value}\n",
		"es/common.ini": "v: {value}\n",
	})

	cases := []struct {
		value any
		want  string
	}{
		{"text", "text"},
		{42, "42"},
		{int64(-7), "-7"},
		{uint8(3), "3"},
		{2.5, "2.5"},
		{true, "true"},
		{nil, ""},
	}

	for _, tc := range cases {
		if got := manager.Get(English, "v", Args{"value": tc.value}); got != tc.want {
			t.Errorf("Get(%v) = %q, want %q", tc.value, got, tc.want)
		}
	}
}

func TestPairsAvoidTheMap(t *testing.T) {
	manager := memoryManager(t, map[string]string{
		"en/common.ini": "greet: Hello {player}\n",
		"es/common.ini": "greet: Hola {player}\n",
	})

	if got := manager.GetPairs(Spanish, "greet", NewPair("player", "Neo")); got != "Hola Neo" {
		t.Errorf("GetPairs = %q", got)
	}
}

func TestLocalizerBindsALanguage(t *testing.T) {
	manager := memoryManager(t, map[string]string{
		"en/common.ini": "greet: Hello\n",
		"es/common.ini": "greet: Hola\n",
	})

	if got := manager.Localizer(Spanish).Get("greet", nil); got != "Hola" {
		t.Errorf("Localizer.Get = %q", got)
	}

	if got := manager.Localizer(Spanish).With(English).Get("greet", nil); got != "Hello" {
		t.Errorf("Localizer.With = %q", got)
	}
}

func TestAllocationBudget(t *testing.T) {
	manager := memoryManager(t, map[string]string{
		"en/common.ini": "static: <red>Denied</red>\ndynamic: Hello {player}\n",
		"es/common.ini": "static: <red>Denegado</red>\ndynamic: Hola {player}\n",
	})

	if allocs := testing.AllocsPerRun(200, func() {
		_ = manager.Raw(English, "static")
	}); allocs != 0 {
		t.Errorf("Raw must not allocate, got %.1f allocations", allocs)
	}

	if allocs := testing.AllocsPerRun(200, func() {
		_ = manager.Get(English, "static", nil)
	}); allocs != 0 {
		t.Errorf("a message without arguments must not allocate, got %.1f allocations", allocs)
	}

	args := Args{"player": "Neo"}
	if allocs := testing.AllocsPerRun(200, func() {
		_ = manager.Get(English, "dynamic", args)
	}); allocs > 2 {
		t.Errorf("rendering should stay within 2 allocations, got %.1f", allocs)
	}
}
