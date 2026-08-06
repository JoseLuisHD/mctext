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
	"os"
)

// maxSourceFileSize is a hard ceiling applied while reading, before the parser
// budgets are even consulted, so a huge file cannot be buffered in full.
const maxSourceFileSize = 64 << 20

// Every read is confined with os.Root, so a crafted file name can neither walk
// out of the directory with ".." nor follow a symbolic link that points
// outside of it. This is the reason the module requires Go 1.26.5: earlier
// patch releases allowed a trailing slash to escape the root (CVE-2026-39822).
//
// Files are opened one at a time; loading happens once at start-up, so the
// extra syscall per file is irrelevant.
type dirSource struct {
	dir string
}

func DirSource(dir string) Source {
	return &dirSource{dir: dir}
}

func (s *dirSource) Describe() string {
	return s.dir
}

func (s *dirSource) ReadFile(name string) ([]byte, error) {
	clean, err := sanitiseRelativePath(name)
	if err != nil {
		return nil, err
	}

	root, err := os.OpenRoot(s.dir)
	if err != nil {
		return nil, err
	}

	defer root.Close()

	file, err := root.Open(clean)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if !info.Mode().IsRegular() {
		return nil, newError(ErrSource, "not a regular file: "+name)
	}

	if info.Size() > maxSourceFileSize {
		return nil, newError(ErrSource, "translation file is too large: "+name)
	}

	return io.ReadAll(io.LimitReader(file, maxSourceFileSize+1))
}
