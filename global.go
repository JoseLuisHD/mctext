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
	"sync"
	"sync/atomic"
)

// The process wide instance. It is the convenience layer described in the
// requirements: one initialisation at start-up and free functions everywhere
// else. Everything it does is delegated to a plain Manager, which remains
// directly constructible with New for tests and for embedding several
// independent catalogs in one process.
var (
	globalOnce    sync.Once
	globalManager atomic.Pointer[Manager]
)

// Init loads every declared translation file exactly once for the whole
// process. Subsequent calls do not reload anything and return
// ErrAlreadyInitialised.
func Init(options ...Option) error {
	var (
		initialised bool
		err         error
	)

	globalOnce.Do(func() {
		initialised = true
		manager, loadErr := New(options...)
		if loadErr != nil {
			err = loadErr
			return
		}
		globalManager.Store(manager)
	})

	if !initialised {
		return newError(ErrAlreadyInitialised, "Init was already called").
			Hint("use Reload to refresh, or New to build an independent manager")
	}

	return err
}

func MustInit(options ...Option) {
	if err := Init(options...); err != nil {
		panic(err)
	}
}

func Initialised() bool {
	return globalManager.Load() != nil
}

func Global() *Manager {
	return globalManager.Load()
}

func Reload() error {
	manager := globalManager.Load()
	if manager == nil {
		return newError(ErrNotInitialised, "Reload before Init")
	}

	return manager.Reload()
}

func Get(language Language, key string, args Args) string {
	manager := globalManager.Load()
	if manager == nil {
		return key
	}

	return manager.Get(language, key, args)
}

func GetPairs(language Language, key string, pairs ...Pair) string {
	manager := globalManager.Load()
	if manager == nil {
		return key
	}

	return manager.GetPairs(language, key, pairs...)
}

func TryGet(language Language, key string, args Args) (string, error) {
	manager := globalManager.Load()
	if manager == nil {
		return key, newError(ErrNotInitialised, "TryGet before Init").
			Field("key", key).
			Hint("call mctext.Init(…) once during start-up")
	}

	return manager.TryGet(language, key, args)
}

func Raw(language Language, key string) string {
	manager := globalManager.Load()
	if manager == nil {
		return key
	}

	return manager.Raw(language, key)
}

func TryRaw(language Language, key string) (string, error) {
	manager := globalManager.Load()
	if manager == nil {
		return key, newError(ErrNotInitialised, "TryRaw before Init")
	}

	return manager.TryRaw(language, key)
}

func Has(language Language, key string) bool {
	manager := globalManager.Load()
	return manager != nil && manager.Has(language, key)
}

func Of(language Language) *Localizer {
	manager := globalManager.Load()
	if manager == nil {
		return &Localizer{manager: nil, language: language}
	}

	return manager.Localizer(language)
}

func Languages() []Language {
	manager := globalManager.Load()
	if manager == nil {
		return nil
	}

	return manager.Languages()
}
