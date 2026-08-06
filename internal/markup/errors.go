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

import (
	"errors"

	"github.com/JoseLuisHD/mctext/internal/diag"
)

// Sentinel errors reported by the compiler. Every returned error wraps one of
// them, so callers can classify failures with errors.Is while still printing
// the full diagnostic.
var (
	// ErrMarkup reports malformed colour or format markup.
	ErrMarkup = errors.New("malformed markup")
	// ErrPlaceholder reports a malformed argument placeholder.
	ErrPlaceholder = errors.New("malformed argument placeholder")
	// ErrLimit reports that a compile-time budget was exhausted.
	ErrLimit = diag.ErrLimit
)
