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

// Package markup compiles the tag based colour syntax of translation messages
// into the formatting codes the game client understands.
//
// The vocabulary lives in vocabulary.go, the appearance model in style.go and
// the wire format in codes.go; the compiler files turn one source message into
// an immutable Program that renders without parsing anything.

package markup
