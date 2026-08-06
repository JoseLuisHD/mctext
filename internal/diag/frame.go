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

package diag

import (
	"runtime"
	"strings"
)

// Frame is a single entry of a captured call stack.
type Frame struct {
	Function string
	File     string
	Line     int
}

// Capture walks the call stack and returns at most depth frames, skipping the
// innermost ones whose package path starts with ignorePrefix.
func Capture(skip, depth int, ignorePrefix string) []Frame {
	if depth <= 0 {
		return nil
	}

	// Walk a generous window so the ignored library frames can be dropped
	// without truncating the frames the caller actually cares about.
	pcs := make([]uintptr, depth+16)
	n := runtime.Callers(skip+2, pcs)
	if n == 0 {
		return nil
	}

	frames := runtime.CallersFrames(pcs[:n])
	out := make([]Frame, 0, depth)
	skipping := ignorePrefix != ""
	for len(out) < depth {
		frame, more := frames.Next()
		if frame.Function == "" && !more {
			break
		}

		if skipping {
			if belongsTo(frame.Function, ignorePrefix) {
				if !more {
					break
				}

				continue
			}

			skipping = false
		}

		out = append(out, Frame{
			Function: trimPackagePath(frame.Function),
			File:     shortPath(frame.File),
			Line:     frame.Line,
		})

		if !more {
			break
		}
	}

	return out
}

// belongsTo reports whether a symbol lives in the given package path. The
// boundary check keeps sibling packages such as "…/translation_test" visible.
func belongsTo(function, prefix string) bool {
	if !strings.HasPrefix(function, prefix) {
		return false
	}

	rest := function[len(prefix):]

	return rest == "" || rest[0] == '.' || rest[0] == '/'
}

// trimPackagePath keeps the last path element of a fully qualified symbol so
// stacks stay readable: "github.com/hycrow-network/game/shop.(*Handler).Buy" becomes
// "shop.(*Handler).Buy".
func trimPackagePath(fn string) string {
	if i := strings.LastIndexByte(fn, '/'); i >= 0 {
		return fn[i+1:]
	}

	return fn
}

// shortPath keeps the last two path elements of a file name.
func shortPath(path string) string {
	if i := strings.LastIndexByte(path, '/'); i >= 0 {
		if j := strings.LastIndexByte(path[:i], '/'); j >= 0 {
			return path[j+1:]
		}
	}

	return path
}
