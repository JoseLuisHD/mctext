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

// Args maps placeholder names to the values that fill them.
//
//	mctext.Args{"player": name, "amount": 250}
type Args map[string]any

// Pair is the allocation free alternative to Args for hot paths.
type Pair struct {
	Name  string
	Value any
}

func NewPair(name string, value any) Pair {
	return Pair{Name: name, Value: value}
}

// pairsView adapts a slice of pairs to the same lookup shape as Args without
// building a map.
type pairsView []Pair

func (p pairsView) lookup(name string) (any, bool) {
	for i := range p {
		if p[i].Name == name {
			return p[i].Value, true
		}
	}

	return nil, false
}

func (p pairsView) names() []string {
	names := make([]string, len(p))
	for i := range p {
		names[i] = p[i].Name
	}

	return names
}
