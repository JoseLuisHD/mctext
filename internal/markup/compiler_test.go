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

import (
	"testing"

	"github.com/JoseLuisHD/mctext/internal/diag"
)

func compile(t *testing.T, raw string) *Program {
	t.Helper()
	program, err := Compile("test::key", raw, diag.Location{File: "test.ini", Line: 1, Column: 1}, Options{})
	if err != nil {
		t.Fatalf("Compile(%q) returned an unexpected error:\n%v", raw, err)
	}

	return program
}

func render(t *testing.T, program *Program, values ...string) string {
	t.Helper()
	return string(program.AppendTo(nil, values))
}

func TestCompileColours(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want string
	}{
		{"plain text", "Some Message", "Some Message"},
		{"single colour", "<red>Some Message</red>", "§cSome Message§r"},
		{"colour before format", "<bold:red>Some Message</bold:red>", "§c§lSome Message§r"},
		{"component order is irrelevant", "<red:bold>x</red:bold>", "§c§lx§r"},
		{"closing tag order is irrelevant", "<bold:red>x</red:bold>", "§c§lx§r"},
		{"generic close", "<aqua>x</>", "§bx§r"},
		{"nested additive format", "<gold>a<bold>b</bold>c</gold>", "§6a§lb§6c§r"},
		{"nested colour override", "<green>a<red>b</red>c</green>", "§aa§cb§ac§r"},
		{"material colour", "<material_netherite>x</material_netherite>", "§jx§r"},
		{"reset is self contained", "<reset:white>plain", "§r§fplain"},
		{"multiple decorations", "<italic:obfuscated:blue>x</>", "§9§k§ox§r"},
		{"escapes", `\<not a tag\> \{not an arg\}`, "<not a tag> {not an arg}"},
		{"adjacent runs", "<red>a</red><blue>b</blue>", "§ca§r§9b§r"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			program := compile(t, tc.raw)
			static, ok := program.Static()
			if !ok {
				t.Fatalf("expected a static program for %q", tc.raw)
			}

			if static != tc.want {
				t.Errorf("Compile(%q)\n got: %q\nwant: %q", tc.raw, static, tc.want)
			}

			if program.Raw() != tc.raw {
				t.Errorf("Raw() = %q, want the verbatim source %q", program.Raw(), tc.raw)
			}
		})
	}
}

func TestCompileArguments(t *testing.T) {
	program := compile(t, "<green>Hi <yellow>{player}</yellow>, you owe {amount} to {player}.</green>")

	want := []string{"player", "amount"}
	if got := program.Args(); len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("Args() = %v, want %v (declared once, in order of appearance)", got, want)
	}

	if _, ok := program.Static(); ok {
		t.Fatal("a message with arguments must not be static")
	}

	got := render(t, program, "Neo", "250")
	const expected = "§aHi §eNeo§a, you owe 250 to Neo.§r"
	if got != expected {
		t.Errorf("AppendTo\n got: %q\nwant: %q", got, expected)
	}
}

func TestCompileAllowsLegacyCodesWhenEnabled(t *testing.T) {
	program, err := Compile("k", "§4Danger", diag.Location{}, Options{AllowLegacyCodes: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if static, _ := program.Static(); static != "§4Danger" {
		t.Errorf("legacy codes must pass through verbatim, got %q", static)
	}
}
