# pwnana

Pronounceable, keyboard-alternating password generator.

pwgen ライクだけど、**長くて打ちやすい**パスワードを生成する。

## Features

- QWERTY の左右交互打鍵を意識した文字配置
- 発音可能 (CV/VC パターン + oo, ee, th, sh などのダイグラフ)
- leet 置換による数字・記号挿入 (i→!, a→@, e→3, o→0 など)
- 文字数・個数・オプションをコマンドラインで指定

## Install

```sh
ln -s /path/to/pwnana ~/.local/bin/pwnana
```

## Quick usage

```sh
pwnana              # 16文字 × 1個
pwnana 20 5         # 20文字 × 5個
pwnana 24 -d -s     # 24文字、数字+記号入り
```

## Customization

スクリプト冒頭の `LAYOUT`、`DIGRAPHS`、`LEET` を編集するだけ。  
詳細は [manual.md](manual.md) を参照。
