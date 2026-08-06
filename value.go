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
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/JoseLuisHD/mctext/internal/markup"
)

// formatValue renders an argument value. The concrete types that dominate game
// messages are handled without reflection; anything else falls back to fmt.
func formatValue(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case nil:
		return ""
	case fmt.Stringer:
		return v.String()
	case error:
		return v.Error()
	default:
		return fmt.Sprint(v)
	}
}

// sanitiseValue removes formatting codes from an argument value.
//
// Argument values are substituted after compilation, so tags inside them are
// never interpreted; raw codes, however, would reach the client verbatim and
// let a player forge coloured output. Stripping them closes that hole.
func sanitiseValue(value string) string {
	index := strings.Index(value, markup.Section)
	if index < 0 {
		return value
	}

	var b strings.Builder
	b.Grow(len(value))
	for index >= 0 {
		b.WriteString(value[:index])
		value = value[index+len(markup.Section):]
		if value != "" {
			// Drop the code character that follows the section sign.
			_, size := utf8.DecodeRuneInString(value)
			value = value[size:]
		}

		index = strings.Index(value, markup.Section)
	}

	b.WriteString(value)

	return b.String()
}
