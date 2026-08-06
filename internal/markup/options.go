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

// Options tunes the strictness and the budgets of the compiler.
type Options struct {
	// AllowLegacyCodes permits raw section codes in source files. It is off by
	// default because raw codes bypass the nesting model and silently break
	// surrounding styles.
	AllowLegacyCodes bool
	// MaxNesting bounds how deeply tags may be nested.
	MaxNesting int
	// MaxArguments bounds how many distinct placeholders a message may declare.
	MaxArguments int
}

// DefaultOptions returns the compiler budgets used when none are supplied.
func DefaultOptions() Options {
	return Options{MaxNesting: 16, MaxArguments: 32}
}

func (o Options) normalized() Options {
	if o.MaxNesting <= 0 {
		o.MaxNesting = DefaultOptions().MaxNesting
	}

	if o.MaxArguments <= 0 {
		o.MaxArguments = DefaultOptions().MaxArguments
	}

	return o
}
