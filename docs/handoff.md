# Handoff

## Plan

- CLI を sh 版から Go 版 (`main.go`) に移行済み。以後は Go 版を正とする
- ドキュメントは英語を `*.md`、日本語を `*.ja.md` として両方維持する
- Arch 向けに `packaging/arch/PKGBUILD` (`pwnana-git`) を用意
- TODO: ライセンスを決めて `LICENSE` を追加し、PKGBUILD の `license=()` を更新する (現状 `LicenseRef-unknown`、namcap がエラーを出す)
- TODO: AUR 公開するか検討

## Log

- 2026-09-25: CLI を Go に書き換え、sh 版 `pwnana` を削除
- 2026-09-25: `tldr/pwnana.ja.page.md` を追加
- 2026-09-25: `README.md` を英語化し、日本語版を `README.ja.md` に分離。相互リンク
- 2026-09-25: `manual.md` / `quickguide.md` を `docs/` に移動し、英語版 (`*.md`) と日本語版 (`*.ja.md`) に分割
- 2026-09-25: `packaging/arch/PKGBUILD` を追加。作業ツリーのスナップショットで makepkg ビルドと動作を確認

## Handoff notes

- ビルド: `go build -o ~/.local/bin/pwnana .`
- Arch パッケージ: `cd packaging/arch && makepkg -si`
- tldr (tealdeer) でローカル表示: `ln -sf "$PWD/tldr/pwnana.ja.page.md" ~/.local/share/tealdeer/pages/pwnana.page.md`
- Web 版 (Cloudflare Workers assets) で配信されるのは `.assetsignore` により `index.html` と `goroawase.json` のみ。md の移動は Web 版に影響しない
- `goroawase.json` は Go バイナリに `go:embed` で埋め込まれ、Web 版とも共有している
