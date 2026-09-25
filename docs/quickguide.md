# pwnana quickguide

[日本語](quickguide.ja.md)

```
pwnana [length] [count] [-d] [-s] [--word-symbols] [-u] [-g]
```

| Example | Description |
|---|---|
| `pwnana` | 16 chars × about half the terminal height |
| `pwnana 20` | 20 chars × about half the terminal height |
| `pwnana 20 5` | 20 chars × 6 |
| `pwnana -d` | 16 chars, with one digit |
| `pwnana -s` | 16 chars, with one symbol |
| `pwnana --word-symbols` | 16 chars, symbols limited to `-+%~&/_` |
| `pwnana -u` | 16 chars, with one uppercase letter |
| `pwnana 24 -d -s` | 24 chars, digit + symbol |
| `pwnana 24 -d -s -u` | 24 chars, digit + symbol + uppercase |
| `pwnana -g` | 16 chars, with a 4-digit goroawase |

- The minimum `length` is 8 (anything smaller becomes 8)
- If `count` is omitted, enough passwords are shown to fill about half the terminal height
- `count` is rounded up to even (+1 when odd)
- Digits and symbols are inserted via leet substitution (e.g. `a→@`, `e→3`, `o→0`)
- If leet substitution isn't possible, digits and symbols are inserted at a fixed position every 4 characters
- `--word-symbols` includes a symbol even when used alone, and uses only `-+%~&/_`
- `-u` uppercases one letter
- As a default readability rule, uppercase `O` / `I` are never generated
