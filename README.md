# sakuramml-go
mml compiler (text music) sakura by golang

テキスト音楽「サクラ」第五版をGo言語で実装しています。

## サクラのコンパイル

サクラを動かすには、golang が必要です。

リポジトリをクローン

```
git clone https://github.com/kujirahand/sakuramml-go.git
# あるいは ... git clone git@github.com:kujirahand/sakuramml-go.git
```

文法を改変したい場合は、goyaccをインストール

```
go get golang.org/x/tools/cmd/goyacc
go install golang.org/x/tools/cmd/goyacc
```

サクラをコンパイル

```
cd sakuramml-go
go build cmd/csakura/csakura.go
```

すると、csakuraというバイナリができる。ドレミのテキストで作曲して、MIDIファイルにコンパイルする。

```
./csakura a.mml
```

# なお完全に作成途中です

現在、基本的なMMLコマンド、およびストトン表記のものを変換できます。

# 完成度

- https://github.com/kujirahand/sakuramml-go/ (30%)
- https://github.com/kujirahand/sakuramml-c/ (40%)
- https://github.com/kujirahand/sakuramml-js/ (1%)
- https://github.com/kujirahand/sakuramml/ (100% --- 本家v2版 Pascal)


