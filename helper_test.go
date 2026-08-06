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

import "testing"

func memoryManager(t *testing.T, files map[string]string, extra ...Option) *Manager {
	t.Helper()
	options := append([]Option{
		WithSource(MemorySource(files)),
		Open(English, "common.ini"),
		Open(Spanish, "common.ini"),
	}, extra...)

	manager, err := New(options...)
	if err != nil {
		t.Fatalf("New returned an unexpected error:\n%v", err)
	}

	return manager
}
