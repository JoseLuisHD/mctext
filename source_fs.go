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
	"io"
	"io/fs"
)

// FSSource adapts any fs.FS, including embed.FS, to the Source port.
func FSSource(fsys fs.FS, description string) Source {
	if description == "" {
		description = "fs"
	}

	return &fsSource{fsys: fsys, description: description}
}

type fsSource struct {
	fsys        fs.FS
	description string
}

func (s *fsSource) Describe() string {
	return s.description
}

func (s *fsSource) ReadFile(name string) ([]byte, error) {
	clean, err := sanitiseRelativePath(name)
	if err != nil {
		return nil, err
	}

	file, err := s.fsys.Open(clean)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	return io.ReadAll(file)
}
