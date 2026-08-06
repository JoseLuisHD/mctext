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
	"io/fs"
	"path"
)

// MemorySource serves files from a map, which keeps tests hermetic and fast.
func MemorySource(files map[string]string) Source {
	copied := make(map[string]string, len(files))
	for name, content := range files {
		copied[path.Clean(name)] = content
	}

	return &memorySource{files: copied}
}

type memorySource struct {
	files map[string]string
}

func (s *memorySource) Describe() string {
	return "memory"
}

func (s *memorySource) ReadFile(name string) ([]byte, error) {
	clean, err := sanitiseRelativePath(name)
	if err != nil {
		return nil, err
	}

	content, ok := s.files[clean]
	if !ok {
		return nil, fs.ErrNotExist
	}

	return []byte(content), nil
}
