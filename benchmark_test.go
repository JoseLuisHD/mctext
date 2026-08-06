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

// sink keeps the compiler from eliminating the benchmarked calls.
var sink string

func benchManager(b *testing.B) *Manager {
	b.Helper()
	manager, err := New(
		WithSource(MemorySource(map[string]string{
			"en/common.ini": "static: <gold:bold>HYCROW</gold:bold> <gray>Network</gray>\n" +
				"dynamic: <green>You bought <yellow>{item}</yellow> for <gold>{price}</gold> coins.</green>\n",
			"es/common.ini": "static: <gold:bold>HYCROW</gold:bold> <gray>Network</gray>\n" +
				"dynamic: <green>Compraste <yellow>{item}</yellow> por <gold>{price}</gold> monedas.</green>\n",
		})),
		Open(English, "common.ini"),
		Open(Spanish, "common.ini"),
		WithStrict(true),
	)

	if err != nil {
		b.Fatalf("New: %v", err)
	}

	return manager
}

func BenchmarkRaw(b *testing.B) {
	manager := benchManager(b)
	b.ReportAllocs()

	for b.Loop() {
		sink = manager.Raw(Spanish, "dynamic")
	}
}

func BenchmarkGetStatic(b *testing.B) {
	manager := benchManager(b)
	b.ReportAllocs()

	for b.Loop() {
		sink = manager.Get(Spanish, "static", nil)
	}
}

func BenchmarkGetWithArgs(b *testing.B) {
	manager := benchManager(b)
	args := Args{"item": "Netherite Sword", "price": 250}
	b.ReportAllocs()

	for b.Loop() {
		sink = manager.Get(Spanish, "dynamic", args)
	}
}

func BenchmarkGetWithPairs(b *testing.B) {
	manager := benchManager(b)
	b.ReportAllocs()

	for b.Loop() {
		sink = manager.GetPairs(Spanish, "dynamic", NewPair("item", "Netherite Sword"), NewPair("price", 250))
	}
}

func BenchmarkGetParallel(b *testing.B) {
	manager := benchManager(b)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		local := ""
		for pb.Next() {
			local = manager.Get(Spanish, "dynamic", Args{"item": "Sword", "price": 1})
		}

		sink = local
	})
}
