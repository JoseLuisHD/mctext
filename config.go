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

type languageFiles struct {
	language Language
	files    []string
}

type config struct {
	source           Source
	resolve          PathResolver
	order            []languageFiles
	fallback         Language
	defaultLanguage  Language
	strict           bool
	allowLegacyCodes bool
	sanitiseValues   bool
	stackDepth       int
	onError          func(error)
	limits           Limits
}

func defaultConfig() config {
	return config{
		resolve:         NestedLayout,
		defaultLanguage: English,
		fallback:        English,
		sanitiseValues:  true,
		stackDepth:      4,
		limits:          DefaultLimits(),
	}
}

func (c *config) iniLimits() ini.Limits {
	return ini.Limits{
		MaxFileSize:  c.limits.MaxFileSize,
		MaxEntries:   c.limits.MaxEntries,
		MaxKeyLength: c.limits.MaxKeyLength,
		MaxLineWidth: c.limits.MaxLineWidth,
	}
}

func (c *config) markupOptions() markup.Options {
	return markup.Options{
		AllowLegacyCodes: c.allowLegacyCodes,
		MaxNesting:       c.limits.MaxNesting,
		MaxArguments:     c.limits.MaxArguments,
	}
}

func (c *config) files(language Language) []string {
	for i := range c.order {
		if c.order[i].language == language {
			return c.order[i].files
		}
	}

	return nil
}
