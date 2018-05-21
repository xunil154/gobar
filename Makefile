
all: gobar

clean:
	rm -rf bin/*

gobar: clean goget
	go build -o bin/gobar gobar.go

goget:
	go get

test: clean goget gobar
	go test github.com/xunil154/gobar/ui

run: gobar
	./bin/gobar
