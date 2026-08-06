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

func newManager(t *testing.T, options ...Option) *Manager {
	t.Helper()
	base := []Option{
		WithSource(MemorySource(map[string]string{
			"en/common.ini": "greet: Hello {player}, you have {count} messages\n",
			"es/common.ini": "greet: Hola {player}, tienes {count} mensajes\n",
		})),
		Open(English, "common.ini"),
		Open(Spanish, "common.ini"),
	}

	manager, err := New(append(base, options...)...)
	if err != nil {
		t.Fatalf("New returned an unexpected error:\n%v", err)
	}

	return manager
}

func TestStackCaptureCanBeDisabled(t *testing.T) {
	manager := newManager(t, WithStackDepth(0))
	_, err := manager.TryGet(English, "greet", nil)
	if err == nil {
		t.Fatal("expected an error")
	}

	if strings.Contains(err.Error(), "stack:") {
		t.Errorf("WithStackDepth(0) must suppress the stack:\n%v", err)
	}
}
