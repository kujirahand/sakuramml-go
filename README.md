# sakuramml-go
mml compiler (text music) sakura by golang

テキスト音楽「サクラ」第五版をGo言語で実装しています。

## サクラのコンパイル

サクラを動かすには、golang / goyacc が必要。

```
brew install golang
```

環境変数GOPATHを確認してgithub.com/kujirhand/sakuramml-goにリポジトリをクローン。

```
cd `go env GOPATH`
pwd
mkdir -p src/github.com/kujirahand
cd src/github.com/kujirahand
```

リポジトリをクローン

```
git clone https://github.com/kujirahand/sakuramml-go.git
# あるいは ... git clone git@github.com:kujirahand/sakuramml-go.git
```

goyaccをインストール

```
go get golang.org/x/tools/cmd/goyacc
go install golang.org/x/tools/cmd/goyacc
```

サクラをコンパイル

```
cd sakuramml-go
go build csakura.go
```

すると、csakuraというバイナリができる。ドレミのテキストで作曲して、MIDIファイルにコンパイルする。

```
./csakura a.mml
```


## Setup

```
$ go get github.com/kujirahand/sakuramml-go
$ go install github.com/kujirahand/sakuramml-go
```

# なお完全に作成途中

- https://github.com/kujirahand/sakuramml-go/ (30%)
- https://github.com/kujirahand/sakuramml-c/ (40%)
- https://github.com/kujirahand/sakuramml-js/ (1%)
- https://github.com/kujirahand/sakuramml/ (100% --- Pascal)


