# pwnana manual

[日本語](manual.ja.md)

## Synopsis

```
pwnana [length] [count] [-d|--digit] [-s|--symbol] [--word-symbols] [-u|--uppercase] [-g|--goroawase] [-h|--help]
```

## Arguments

| Argument | Default | Description |
|---|---|---|
| `length` | 16 | Password length. Minimum 8 |
| `count` | About half the terminal height | Number of passwords. If omitted, computed from the terminal rows (rounded up to even on output) |
| `-d`, `--digit` | off | Include one digit |
| `-s`, `--symbol` | off | Include one symbol |
| `--word-symbols` | off | Limit symbols to `-+%~&/_`. Includes one symbol even when used alone |
| `-u`, `--uppercase` | off | Include one uppercase letter |
| `-g`, `--goroawase` | off | Embed a 4-digit goroawase (Japanese number wordplay) and show its hint |

## Generation algorithm

1. Start randomly from either the left or the right hand
2. Repeatedly generate 2-character chunks
   - 80%: an alternating-hand CV or VC pair  
     e.g.) left consonant + right vowel (`tu`, `ro`), right consonant + left vowel (`na`, `he`)
   - 20%: a vowel-centered digraph (`oo`, `ee`, `ii`, `ou`, `ai`, `au`)
3. Concatenate until the requested length is reached, then truncate the excess
4. As a readability rule, uppercase `O` / `I` are not used by default
5. If the password ends with a consonant other than `n`, it is replaced with a vowel to keep it pronounceable
6. With `-u`, one letter is uppercased (avoiding `O` / `I`)

## Digit / symbol insertion

`-d`/`-s` prefer leet substitution. If the password contains a substitutable character, it is used;
otherwise the digit/symbol is inserted at a fixed position every 4 characters.

With `-g`, goroawase uses **4 digits only** and is embedded at a position aligned to 4-character boundaries.

`--word-symbols` switches the symbols used by `-s` to `-+%~&/_`. Use it when you only want
symbols that are unlikely to break double-click word selection in the terminal.
Specifying `--word-symbols` alone also counts as requesting a symbol.

Default leet map:

| Char | Replacement | Used by |
|---|---|---|
| `i` | `!` | symbol |
| `a` | `@` | symbol |
| `s` | `$` | symbol |
| `e` | `3` | digit |
| `o` | `0` | digit |
| `t` | `7` | digit |
| `q` | `9` | digit |
| `b` | `6` | digit |
| `g` | `8` | digit |

## Customization

Edit the variables at the top of `main.go`.

### layout

```go
layout = struct {
    leftConsonants  string
    leftVowels      string
    rightConsonants string
    rightVowels     string
}{
    leftConsonants:  "wrtsfgzxcvbq",
    leftVowels:      "ea",
    rightConsonants: "yphjklnm",
    rightVowels:     "uio",
}
```

Rewrite this if you use another layout (Dvorak, Colemak, etc.).

### digraphs

```go
digraphs = []string{"oo", "ee", "ii", "ou", "ai", "au"}
```

A list of digraphs that sound natural. One is picked at random 20% of the time.

### leet

```go
leet = map[byte]byte{
    'i': '!',
    'a': '@',
    ...
}
```

Digit values are used by `-d`, symbol values by `-s`. Add or remove entries to suit your taste.
