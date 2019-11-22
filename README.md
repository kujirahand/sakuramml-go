# sakuramml-go
mml compiler (text music) sakura by golang

テキスト音楽「サクラ」をGo言語で実装しています。
完全に作成途中です。

- https://github.com/kujirahand/sakuramml-go/ (30%)
- https://github.com/kujirahand/sakuramml-c/ (40%)
- https://github.com/kujirahand/sakuramml-js/ (1%)
- https://github.com/kujirahand/sakuramml/ (100% --- Pascal)

## Setup

環境変数を手軽に書き換えるdirenvを利用しています。
macOSなら``brew install direnv``でインストールしておいてください。

```
direnv allow
```

## Compile

```
go build src/csakura.go
```
