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
	"errors"

	"github.com/JoseLuisHD/mctext/internal/diag"
	"github.com/JoseLuisHD/mctext/internal/ini"
	"github.com/JoseLuisHD/mctext/internal/markup"
)

// Sentinel errors. Every error returned by this package wraps exactly one of
// them, so callers classify failures with errors.Is and still get the full
// diagnostic report from Error.
var (
	// ErrNotInitialised is returned when the manager is used before Init.
	ErrNotInitialised = errors.New("manager is not initialised")
	// ErrAlreadyInitialised is returned by a second call to Init.
	ErrAlreadyInitialised = errors.New("manager is already initialised")
	// ErrInvalidLanguage reports a malformed or unregistered language tag.
	ErrInvalidLanguage = errors.New("invalid language")
	// ErrNoSource reports that no source of translation files was configured.
	ErrNoSource = errors.New("no source configured")
	// ErrMissingKey reports a lookup for a key no catalog declares.
	ErrMissingKey = errors.New("missing translation key")
	// ErrMissingArgument reports a placeholder that received no value.
	ErrMissingArgument = errors.New("missing argument")
	// ErrUnexpectedArgument reports a supplied value the message never uses.
	ErrUnexpectedArgument = errors.New("unexpected argument")
	// ErrDuplicateKey reports a key declared twice inside the same language.
	ErrDuplicateKey = errors.New("duplicate translation key")
	// ErrInconsistentCatalog reports a mismatch between languages in strict
	// mode: a key or an argument present in one language and absent in another.
	ErrInconsistentCatalog = errors.New("inconsistent catalogs")
	// ErrSource reports that a translation file could not be read.
	ErrSource = errors.New("cannot read translation file")

	// ErrMarkup reports malformed colour or format markup. It is raised by the
	// compiler and surfaces from Init.
	ErrMarkup = markup.ErrMarkup
	// ErrPlaceholder reports a malformed argument placeholder.
	ErrPlaceholder = markup.ErrPlaceholder
	// ErrSyntax reports a malformed translation file.
	ErrSyntax = ini.ErrSyntax
	// ErrLimit reports an exhausted safety budget.
	ErrLimit = diag.ErrLimit
)

// Diagnostic is the rich error value produced by this package. It carries the
// source location, the offending text and a short call stack.
type Diagnostic = diag.Diagnostic

// Location identifies a position inside a translation file.
type Location = diag.Location

// Frame is one entry of the short call stack attached to runtime diagnostics.
type Frame = diag.Frame

func newError(sentinel error, summary string) *Diagnostic {
	return diag.New(sentinel, summary)
}

// stackIgnorePrefix is the package path whose frames are dropped from captured
// stacks so that reports point at game code rather than at this library.
const stackIgnorePrefix = "github.com/JoseLuisHD/mctext"
