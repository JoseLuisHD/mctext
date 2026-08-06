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

package catalog

import (
	"sort"

	"github.com/JoseLuisHD/mctext/internal/markup"
)

// A Catalog is mutable only while it is being built. Once it has been published
// it is treated as immutable, which is what allows lookups to run without any
// synchronisation at all.
type Catalog struct {
	language string
	sources  []string
	messages map[string]*markup.Program
}

func New(language string, sizeHint int) *Catalog {
	if sizeHint < 16 {
		sizeHint = 16
	}

	return &Catalog{
		language: language,
		messages: make(map[string]*markup.Program, sizeHint),
	}
}

func (c *Catalog) Language() string {
	return c.language
}

func (c *Catalog) Sources() []string {
	return c.sources
}

func (c *Catalog) AddSource(name string) {
	c.sources = append(c.sources, name)
}

// Put stores a compiled message, returning the previous definition when the key
// was already declared so the caller can report the collision.
func (c *Catalog) Put(program *markup.Program) (previous *markup.Program, replaced bool) {
	previous, replaced = c.messages[program.Key()]
	c.messages[program.Key()] = program
	return previous, replaced
}

// Lookup returns the compiled message for a key.
func (c *Catalog) Lookup(key string) (*markup.Program, bool) {
	program, ok := c.messages[key]
	return program, ok
}

func (c *Catalog) Len() int {
	return len(c.messages)
}

func (c *Catalog) Keys() []string {
	keys := make([]string, 0, len(c.messages))
	for key := range c.messages {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}

func (c *Catalog) Each(fn func(program *markup.Program)) {
	for _, program := range c.messages {
		fn(program)
	}
}
