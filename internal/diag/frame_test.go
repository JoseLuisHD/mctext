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
	"strings"
	"testing"
)

func TestBelongsTo(t *testing.T) {
	const prefix = "github.com/JoseLuisHD/mctext"
	cases := map[string]bool{
		"github.com/JoseLuisHD/mctext.Get":              true,
		"github.com/JoseLuisHD/mctext.(*Manager).Get":   true,
		"github.com/JoseLuisHD/mctext/internal/ini.New": true,
		"github.com/JoseLuisHD/mctext":                  true,
		"github.com/JoseLuisHD/mctext_test.TestGet":     false,
		"github.com/JoseLuisHD/mctexttool.Run":          false,
		"github.com/JoseLuisHD/game/shop.Buy":           false,
	}

	for function, want := range cases {
		if got := belongsTo(function, prefix); got != want {
			t.Errorf("belongsTo(%q) = %v, want %v", function, got, want)
		}
	}
}

func TestShortPathAndSymbol(t *testing.T) {
	if got := shortPath("/home/build/game/shop/handler.go"); got != "shop/handler.go" {
		t.Errorf("shortPath = %q, want the last two elements", got)
	}

	if got := shortPath("handler.go"); got != "handler.go" {
		t.Errorf("shortPath = %q, want the input unchanged", got)
	}

	if got := trimPackagePath("github.com/JoseLuisHD/mctext/shop.(*Handler).Buy"); got != "shop.(*Handler).Buy" {
		t.Errorf("trimPackagePath = %q", got)
	}
}

func TestCapturePointsAtTheCaller(t *testing.T) {
	frames := Capture(0, 3, "")
	if len(frames) == 0 {
		t.Fatal("Capture returned no frames")
	}

	if !strings.Contains(frames[0].Function, "TestCapturePointsAtTheCaller") {
		t.Errorf("the first frame should be the caller, got %q", frames[0].Function)
	}

	if frames[0].Line == 0 || frames[0].File == "" {
		t.Errorf("the frame carries no position: %+v", frames[0])
	}
}

func TestCaptureSkipsTheIgnoredPackage(t *testing.T) {
	frames := Capture(0, 3, "github.com/JoseLuisHD/mctext/internal/diag")
	for _, frame := range frames {
		if strings.Contains(frame.Function, "TestCaptureSkipsTheIgnoredPackage") {
			t.Fatalf("frames of the ignored package must be dropped: %+v", frames)
		}
	}

	if len(frames) == 0 {
		t.Fatal("skipping must not swallow the whole stack")
	}
}

func TestCaptureRespectsTheDepth(t *testing.T) {
	if frames := Capture(0, 0, ""); frames != nil {
		t.Errorf("a zero depth must disable capture, got %d frames", len(frames))
	}

	if frames := Capture(0, 2, ""); len(frames) > 2 {
		t.Errorf("captured %d frames, want at most 2", len(frames))
	}
}
