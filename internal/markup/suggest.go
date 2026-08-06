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

import "strings"

// suggestThreshold is the exclusive upper bound on the edit distance of an
// acceptable suggestion. Three keeps typos in reach without proposing an
// unrelated colour.
const suggestThreshold = 3

// suggest returns the closest tag name to the given unknown one, or an empty
// string when nothing is close enough to be worth suggesting.
func suggest(unknown string) string {
	unknown = strings.ToLower(unknown)
	best, bestDistance := "", suggestThreshold
	for _, candidate := range tagNames {
		if d := editDistance(unknown, candidate, bestDistance); d < bestDistance {
			best, bestDistance = candidate, d
		}
	}

	return best
}

// editDistance is a bounded Levenshtein distance: it gives up as soon as the
// result is known to exceed limit, which keeps lookups cheap over the whole
// vocabulary.
func editDistance(a, b string, limit int) int {
	if a == b {
		return 0
	}

	if abs(len(a)-len(b)) > limit {
		return limit + 1
	}

	previous := make([]int, len(b)+1)
	current := make([]int, len(b)+1)
	for j := range previous {
		previous[j] = j
	}

	for i := 1; i <= len(a); i++ {
		current[0] = i
		rowMinimum := current[0]
		for j := 1; j <= len(b); j++ {
			substitution := 1
			if a[i-1] == b[j-1] {
				substitution = 0
			}

			current[j] = min(previous[j]+1, current[j-1]+1, previous[j-1]+substitution)
			rowMinimum = min(rowMinimum, current[j])
		}

		if rowMinimum > limit {
			return limit + 1
		}

		previous, current = current, previous
	}

	return previous[len(b)]
}

func abs(value int) int {
	if value < 0 {
		return -value
	}

	return value
}
