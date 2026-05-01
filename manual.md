# pwnana manual

## Synopsis

```
pwnana [length] [count] [-d|--digit] [-s|--symbol] [-g|--goroawase] [-h|--help]
```

## Arguments

| 引数 | デフォルト | 説明 |
|---|---|---|
| `length` | 16 | パスワードの文字数。最小 8 |
| `count` | 端末高さの約半分 | 生成する個数。省略時は端末の行数から自動計算（出力時に偶数へ丸め） |
| `-d`, `--digit` | off | 数字を1つ含める |
| `-s`, `--symbol` | off | 記号を1つ含める |
| `-g`, `--goroawase` | off | 4桁の語呂合わせを埋め込み、ヒントを表示する |

## Generation algorithm

1. 左手 / 右手いずれかからランダムにスタート
2. 2文字チャンクを繰り返し生成
   - 80%: 左右交互の CV または VC ペア  
     例) 左子音+右母音 (`tu`, `ro`)、右子音+左母音 (`na`, `he`)
   - 20%: 母音中心ダイグラフ (`oo`, `ee`, `ii`, `ou`, `ai`, `au`)
3. 指定長になるまで連結し、超過分を切り捨て
4. 可読性ルールとして、大文字 `O` / `I` はデフォルトで使用しない
5. 語尾が子音で終わる場合、`n` 以外は母音に置換して発音しやすくする

## Digit / symbol insertion

`-d`/`-s` は leet 置換を優先する。パスワード内に置換対象文字があればそれを使い、
なければ 4文字ごとの固定位置に挿入する。

`-g` 指定時の語呂合わせは **4桁のみ** 使用し、4文字単位の開始位置に埋め込む。

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
DIGRAPHS = ['oo', 'ee', 'ii', 'ou', 'ai', 'au']
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
