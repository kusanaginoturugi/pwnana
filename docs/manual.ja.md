# pwnana manual

[English](manual.md)

## Synopsis

```
pwnana [length] [count] [-d|--digit] [-s|--symbol] [--word-symbols] [-u|--uppercase] [-g|--goroawase] [-h|--help]
```

## Arguments

| 引数 | デフォルト | 説明 |
|---|---|---|
| `length` | 16 | パスワードの文字数。最小 8 |
| `count` | 端末高さの約半分 | 生成する個数。省略時は端末の行数から自動計算（出力時に偶数へ丸め） |
| `-d`, `--digit` | off | 数字を1つ含める |
| `-s`, `--symbol` | off | 記号を1つ含める |
| `--word-symbols` | off | 記号を `-+%~&/_` に絞る。単独指定でも記号を1つ含める |
| `-u`, `--uppercase` | off | 大文字を1つ含める |
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
6. `-u` 指定時は、英字1文字を大文字化する (`O` / `I` は避ける)

## Digit / symbol insertion

`-d`/`-s` は leet 置換を優先する。パスワード内に置換対象文字があればそれを使い、
なければ 4文字ごとの固定位置に挿入する。

`-g` 指定時の語呂合わせは **4桁のみ** 使用し、4文字単位の開始位置に埋め込む。

`--word-symbols` は `-s` の記号を `-+%~&/_` に切り替える。ターミナル上で
マウスのダブルクリック選択を壊しにくい記号だけを使いたいとき向け。
`--word-symbols` だけ指定した場合も、記号入りとして扱う。

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

`main.go` 冒頭の変数を編集する。

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

別配列 (Dvorak, Colemak など) を使う場合はここを書き換える。

### digraphs

```go
digraphs = []string{"oo", "ee", "ii", "ou", "ai", "au"}
```

発音上自然なダイグラフのリスト。20% の確率でランダムに選ばれる。

### leet

```go
leet = map[byte]byte{
    'i': '!',
    'a': '@',
    ...
}
```

値が数字なら `-d`、記号なら `-s` で使われる。エントリを追加・削除して好みに合わせる。
