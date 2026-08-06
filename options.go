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

// Option mutates the configuration of a Manager before it loads anything.
// Options are applied in order, so later ones win.
type Option func(*config) error

// WithSource sets where translation files are read from.
func WithSource(source Source) Option {
	return func(c *config) error {
		if source == nil {
			return newError(ErrNoSource, "nil source")
		}

		c.source = source

		return nil
	}
}

// WithDirectory is shorthand for WithSource(DirSource(dir)).
func WithDirectory(dir string) Option {
	return WithSource(DirSource(dir))
}

// WithLayout replaces the strategy that maps a language and a declared file
// name onto a path inside the source. The default is NestedLayout.
func WithLayout(resolve PathResolver) Option {
	return func(c *config) error {
		if resolve == nil {
			return newError(ErrNoSource, "nil path resolver")
		}

		c.resolve = resolve

		return nil
	}
}

// Open declares the files that make up a language, in load order. Later files
// may not redefine a key already declared by an earlier one.
//
// Files are always listed explicitly: nothing is discovered by scanning the
// directory, so a stray file can never silently join the build.
func Open(language Language, files ...string) Option {
	return func(c *config) error {
		if err := validateTag(string(language)); err != nil {
			return err
		}

		if len(files) == 0 {
			return newError(ErrNoSource, "no files declared for language "+string(language))
		}

		for i := range c.order {
			if c.order[i].language == language {
				c.order[i].files = append(c.order[i].files, files...)
				return nil
			}
		}

		c.order = append(c.order, languageFiles{
			language: language,
			files:    append([]string(nil), files...),
		})

		return nil
	}
}

// WithDefault sets the language used by the package level helpers that take no
// language, and by Localizer when none is supplied.
func WithDefault(language Language) Option {
	return func(c *config) error {
		if err := validateTag(string(language)); err != nil {
			return err
		}

		c.defaultLanguage = language

		return nil
	}
}

// WithFallback sets the language consulted when a key is missing from the
// requested one. Pass an empty language to disable fallback entirely.
func WithFallback(language Language) Option {
	return func(c *config) error {
		if language != "" {
			if err := validateTag(string(language)); err != nil {
				return err
			}
		}

		c.fallback = language

		return nil
	}
}

// WithStrict turns cross-language consistency into a load error: every language
// must declare the same keys, and every key must declare the same arguments in
// every language. It also rejects unexpected arguments at render time.
//
// Enabling it is strongly recommended: it converts a whole class of production
// bugs into a deployment failure.
func WithStrict(strict bool) Option {
	return func(c *config) error {
		c.strict = strict
		return nil
	}
}

// WithLegacyCodes allows raw section codes inside source files. It exists for
// migrating existing content and should be turned off once the migration is
// done, because raw codes bypass the nesting model.
func WithLegacyCodes(allow bool) Option {
	return func(c *config) error {
		c.allowLegacyCodes = allow
		return nil
	}
}

// WithArgumentSanitisation controls whether formatting codes are stripped from
// argument values before they are substituted.
//
// It defaults to true, and turning it off is a security decision: a player
// controlled value such as a nickname could otherwise inject colour codes and
// impersonate system messages.
func WithArgumentSanitisation(enabled bool) Option {
	return func(c *config) error {
		c.sanitiseValues = enabled
		return nil
	}
}

// WithStackDepth sets how many caller frames runtime diagnostics capture.
// Zero disables stack capture.
func WithStackDepth(depth int) Option {
	return func(c *config) error {
		if depth < 0 {
			depth = 0
		}

		c.stackDepth = depth

		return nil
	}
}

// WithErrorHandler installs a sink for the non fatal failures of Get and Raw,
// which return a best effort string rather than an error. Typically a logger.
func WithErrorHandler(handler func(error)) Option {
	return func(c *config) error {
		c.onError = handler
		return nil
	}
}

// WithLimits overrides the loader budgets. Zero valued fields keep their
// default.
func WithLimits(limits Limits) Option {
	return func(c *config) error {
		defaults := DefaultLimits()
		if limits.MaxFileSize <= 0 {
			limits.MaxFileSize = defaults.MaxFileSize
		}

		if limits.MaxEntries <= 0 {
			limits.MaxEntries = defaults.MaxEntries
		}

		if limits.MaxKeyLength <= 0 {
			limits.MaxKeyLength = defaults.MaxKeyLength
		}

		if limits.MaxLineWidth <= 0 {
			limits.MaxLineWidth = defaults.MaxLineWidth
		}

		if limits.MaxNesting <= 0 {
			limits.MaxNesting = defaults.MaxNesting
		}

		if limits.MaxArguments <= 0 {
			limits.MaxArguments = defaults.MaxArguments
		}

		c.limits = limits

		return nil
	}
}
