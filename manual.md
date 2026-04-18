# pwnana manual

## Synopsis

```
pwnana [length] [count] [-d|--digit] [-s|--symbol] [-h|--help]
```

## Arguments

| 引数 | デフォルト | 説明 |
|---|---|---|
| `length` | 16 | パスワードの文字数。最小 12 |
| `count` | 1 | 生成する個数 |
| `-d`, `--digit` | off | 数字を1つ含める |
| `-s`, `--symbol` | off | 記号を1つ含める |

## Generation algorithm

1. 左手 / 右手いずれかからランダムにスタート
2. 2文字チャンクを繰り返し生成
   - 80%: 左右交互の CV または VC ペア  
     例) 左子音+右母音 (`tu`, `ro`)、右子音+左母音 (`na`, `he`)
   - 20%: ダイグラフ (`oo`, `th`, `sh`, `ai` など)
3. 指定長になるまで連結し、超過分を切り捨て

## Digit / symbol insertion

`-d`/`-s` は leet 置換を優先する。パスワード内に置換対象文字があればそれを使い、
なければランダムな位置に挿入する。

デフォルトの leet マップ:

| 文字 | 置換 | 用途 |
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

スクリプト冒頭の定数を編集する。

### LAYOUT

```python
LAYOUT = {
    'left_consonants':  list('wrtsfgzxcvbq'),
    'left_vowels':      list('ea'),
    'right_consonants': list('yphjklnm'),
    'right_vowels':     list('uio'),
}
```

別配列 (Dvorak, Colemak など) を使う場合はここを書き換える。

### DIGRAPHS

```python
DIGRAPHS = ['oo', 'ee', 'ii', 'th', 'sh', 'wh', 'ou', 'ai', 'au', 'ng']
```

発音上自然なダイグラフのリスト。20% の確率でランダムに選ばれる。

### LEET

```python
LEET = {
    'i': '!',
    'a': '@',
    ...
}
```

値が数字なら `-d`、記号なら `-s` で使われる。エントリを追加・削除して好みに合わせる。
