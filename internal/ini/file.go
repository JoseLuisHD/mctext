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

package ini

import "github.com/JoseLuisHD/mctext/internal/diag"

// Entry is one parsed declaration.
type Entry struct {
	Key      string
	Value    string
	Location diag.Location // Column points at the first character of the value
}

// File is the parsed content of one translation file, in declaration order.
type File struct {
	Name    string
	Entries []Entry
}
