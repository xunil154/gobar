
all: gobar

gobar: goget
	go build -o bin/gobar gobar.go

goget:
	go get

test: goget gobar
	go test github.com/xunil154/gobar/ui

run: gobar
	./gobar
