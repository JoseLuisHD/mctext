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
	"github.com/JoseLuisHD/mctext/internal/ini"
	"github.com/JoseLuisHD/mctext/internal/markup"
)

// Limits bounds the resources the loader is allowed to consume.
type Limits struct {
	MaxFileSize  int // bytes per translation file
	MaxEntries   int // declarations per file
	MaxKeyLength int // bytes per key
	MaxLineWidth int // bytes per source line
	MaxNesting   int // nested tags per message
	MaxArguments int // distinct placeholders per message
}

// DefaultLimits returns the budgets applied when none are configured.
func DefaultLimits() Limits {
	fileLimits := ini.DefaultLimits()
	markupLimits := markup.DefaultOptions()
	return Limits{
		MaxFileSize:  fileLimits.MaxFileSize,
		MaxEntries:   fileLimits.MaxEntries,
		MaxKeyLength: fileLimits.MaxKeyLength,
		MaxLineWidth: fileLimits.MaxLineWidth,
		MaxNesting:   markupLimits.MaxNesting,
		MaxArguments: markupLimits.MaxArguments,
	}
}
