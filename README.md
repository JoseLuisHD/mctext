# mctext

`https://github.com/JoseLuisHD/mctext`

Multi-language message catalog for the Minecraft game. Messages live in INI files, are
declared explicitly at start-up, and are compiled once so that serving one at
runtime is a map lookup plus a byte concatenation.

- **Go 1.26.5**, zero external dependencies.
- **Colour markup as tags** — `<red>…</red>` instead of raw `§` codes, with real
  nesting.
- **Named arguments** — `{player}`, filled at lookup time.
- **Two access modes** — `Raw` (verbatim, zero allocations) and `Get`
  (colours + arguments resolved).
- **Loud diagnostics** — file, line, column, the offending text underlined, and
  a short call stack for runtime failures.
- **Lock-free reads**, safe for any number of goroutines.

---

## Quick start

```go
err := mctext.Init(
    mctext.WithDirectory("assets/lang"),
    mctext.Open(mctext.English, "common.ini", "shop.ini"),
    mctext.Open(mctext.Spanish, "common.ini", "shop.ini"),
    mctext.WithDefault(mctext.English),
    mctext.WithFallback(mctext.English),
    mctext.WithStrict(true),
    mctext.WithErrorHandler(func(err error) { log.Print(err) }),
)
if err != nil {
    log.Fatal(err)
}

msg := mctext.Get(mctext.Spanish, "shop::purchase::success",
    mctext.Args{"item": "Espada de Netherite", "price": 250})
```

With the default layout, `common.ini` is read from `assets/lang/en` for English
and from `assets/lang/es` for Spanish. Nothing is discovered by scanning the
directory: a file only takes part if it is listed in an `Open` option.

---

## File format

```ini
# comments start with # or ;
ui::button::close: <red>Close</red>
ui::plain::text:   Nothing special here

# a section prefixes every key that follows it, until the next one; [] clears it
[shop::purchase]
success: <green>Bought <yellow>{item}</yellow> for <gold>{price}</gold> coins.</green>
denied:  <bold:red>Not enough coins.</bold:red>

long::notice: first half \
              second half
```

- Keys use `::` as a namespace separator. A **single** `:` ends the key, so
  `id::some::message1: text` declares `id::some::message1`. `=` also works.
- A trailing backslash continues the line; write `\\` for a literal one.
- Wrap a value in double quotes to preserve its surrounding whitespace.
- One language may be split across as many files as you like. A key may only be
  declared once per language, across all of its files.

---

## Markup

| Written | Emitted |
|---|---|
| `<red>Denied</red>` | `§cDenied§r` |
| `<bold:red>Denied</bold:red>` | `§c§lDenied§r` |
| `<red:bold>…</bold:red>` | same — component order is irrelevant |
| `<gold>a<bold>b</bold>c</gold>` | `§6a§lb§6c§r` |
| `<green>a<red>b</red>c</green>` | `§aa§cb§ac§r` |
| `<reset:white>plain` | `§r§fplain` |
| `<aqua>x</>` | `§bx§r` — `</>` closes the innermost tag |

Two rules drive the emitted bytes:

1. **Colour first, decoration second.** On the client a colour code clears every
   active decoration, so `§l§c` would silently drop the bold. The compiler
   normalises the order, which is why `<bold:red>` emits `§c§l`.
2. **Closing restores the parent, it does not clear everything.** Tags nest with
   a state stack, so the inner `</bold>` above returns to gold rather than to
   the client default. The compiler also emits the shortest correct transition:
   adding a decoration on top of the same colour costs one code, and a colour
   change needs no reset because the colour code is itself a reset.

**Colours** — `black` `dark_blue` `dark_green` `dark_aqua` `dark_red`
`dark_purple` `gold` `gray` `dark_gray` `blue` `green` `aqua` `red`
`light_purple` `yellow` `white` `minecoin_gold` `light_blue` `material_quartz`
`material_iron` `material_netherite` `material_redstone` `material_copper`
`material_gold` `material_emerald` `material_diamond` `material_lapis`
`material_amethyst` `material_resin`

**Decorations** — `obfuscated` `bold` `italic`, plus the `reset` directive.

**Escapes** — `\<` `\>` `\{` `\}` `\\` for literals, `\n` `\t` for whitespace.

Raw `§` codes in a source file are rejected by default, because they bypass the
nesting model and corrupt the surrounding style. `WithLegacyCodes(true)` allows
them while migrating existing content.

---

## Diagnostics

Markup and file errors fail at `Init`, never in production:

```
translation: mismatched closing tag </gold>
  at:       en/ui.ini:1:39
  key:      ui::menu::title
  expected: </bold> opened at column 12
  │ <gold>Menu <bold>Main</gold>
  │                      ^^^^^^^
  hint:     close inner tags before outer ones, or use </> to close the innermost tag
```

Runtime failures carry a short call stack whose library frames are stripped, so
it points straight at the game code:

```
translation: no value supplied for "price"
  at:       en/ui.ini:1:26
  key:      shop::purchase::success
  language: en
  declared: [item price]
  supplied: [item]
  │ <green>Bought {item} for {price}</green>
  hint:     pass it as mctext.Args{"price": value}
  stack:
    shop.(*Handler).onPurchase (game/shop/handler.go:88)
    net.(*Session).dispatch (game/net/session.go:214)
```

Every error wraps a sentinel, so it can be classified precisely:

```go
if errors.Is(err, mctext.ErrMissingArgument) { … }
```

`ErrMarkup` · `ErrPlaceholder` · `ErrSyntax` · `ErrLimit` · `ErrMissingKey` ·
`ErrMissingArgument` · `ErrUnexpectedArgument` · `ErrDuplicateKey` ·
`ErrInconsistentCatalog` · `ErrSource` · `ErrInvalidLanguage` ·
`ErrNotInitialised` · `ErrAlreadyInitialised` · `ErrNoSource`

`Get` and `Raw` never fail loudly: a missing key degrades to the key itself and
a missing argument leaves its placeholder in place, so a translation mistake
never breaks the caller. The full diagnostic still reaches the configured error
handler, and `TryGet` / `TryRaw` return it directly.

### Strict mode

`WithStrict(true)` verifies at load time that every language declares the same
keys and that every key declares the same arguments in every language, and
rejects values no placeholder consumes. It turns the most common production
failure of a translation system — a message that only exists in one language —
into a start-up error. Recommended everywhere.

---

### Adding a language

```go
portuguese, err := mctext.RegisterLanguage("pt-br", "Português do Brasil")
…
mctext.Open(portuguese, "common.ini", "shop.ini"),
```

Registered default languages

| Constant | Language | Code |
|----------|----------|------|
| `English` | English | `en` |
| `Spanish` | Spanish | `es` |
| `French` | French | `fr` |
| `German` | German | `de` |
| `Italian` | Italian | `it` |
| `Portuguese` | Portuguese | `pt` |
| `Russian` | Russian | `ru` |
| `Japanese` | Japanese | `ja` |
| `Korean` | Korean | `ko` |
| `Chinese` | Chinese | `zh` |
| `Arabic` | Arabic | `ar` |
| `Hindi` | Hindi | `hi` |
| `Turkish` | Turkish | `tr` |
| `Dutch` | Dutch | `nl` |
| `Polish` | Polish | `pl` |
| `Swedish` | Swedish | `sv` |
| `Norwegian` | Norwegian | `no` |
| `Danish` | Danish | `da` |
| `Finnish` | Finnish | `fi` |
| `Czech` | Czech | `cs` |
| `Greek` | Greek | `el` |
| `Hungarian` | Hungarian | `hu` |
| `Romanian` | Romanian | `ro` |
| `Ukrainian` | Ukrainian | `uk` |
| `Thai` | Thai | `th` |
| `Vietnamese` | Vietnamese | `vi` |
| `Indonesian` | Indonesian | `id` |
| `Malay` | Malay | `ms` |
| `Hebrew` | Hebrew | `he` |
| `Persian` | Persian | `fa` |


### Embedding the catalogs in the binary

```go
//go:embed lang
var langFS embed.FS

mctext.WithSource(mctext.FSSource(langFS, "embedded lang")),
```

---

## Testing

```bash
go test ./...                                     # unit and integration tests
go test -race ./...                               # concurrency
go test -bench . -benchmem ./...                  # benchmarks
go test ./internal/markup -fuzz FuzzCompile       # parser fuzzing
```

The suite covers the golden byte output of every markup rule, the INI dialect,
fallback, strict validation, sanitisation, traversal rejection, atomic reload,
concurrent access, and the allocation budget of each access mode.
