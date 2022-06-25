#!/bin/sh
goyacc parser.go.y
go build cmd/csakura/csakura.go
./csakura -d b.mml

